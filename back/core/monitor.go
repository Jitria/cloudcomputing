package core

import (
	"assign/common"
	"assign/config"
	"assign/logger"
	"assign/types"
	"context"
	"fmt"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerimage "github.com/docker/docker/api/types/image"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Monitor() {
	go monitorPod()
	go monitorService()
	go makeImage()
}

////////////////
/*  snapshot   */
////////////////

type snapShot struct {
	containerName string
	commitName    string
}

func makeImage() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			exists := logger.GetpodNames()
			deleteImage(exists)
			snapShotInfos := getCommitInfo(exists)
			commitContainer(snapShotInfos)
		}
	}
}

func deleteImage(exists []string) {
	var imageList []dockerimage.Summary
	var err error
	if imageList, err = config.GlobalConfig.DockerClient.ImageList(context.Background(), dockerimage.ListOptions{}); err != nil {
		common.StopProgram(err)
	}
	for _, image := range imageList {
		if len(image.RepoTags) > 0 {
			imageName := strings.Split(strings.Trim(image.RepoTags[0], "[]"), ":latest")[0]
			for _, exist := range exists {
				if imageName == exist {
					config.GlobalConfig.DockerClient.ImageRemove(context.Background(), image.ID, dockerimage.RemoveOptions{})
				}
			}
		}
	}
}

func getCommitInfo(exists []string) []snapShot {
	var snapShotInfos []snapShot
	var containerList []dockertypes.Container
	var err error

	if containerList, err = config.GlobalConfig.DockerClient.ContainerList(context.Background(), dockercontainer.ListOptions{}); err != nil {
		common.StopProgram(err)
	}

	for _, container := range containerList {
		containerName := strings.Trim(container.Names[0], "[]")
		commitName := strings.Split(containerName, "_")[1]
		for _, exist := range exists {
			if commitName == exist {
				snapShotInfo := snapShot{
					containerName: containerName,
					commitName:    commitName,
				}
				snapShotInfos = append(snapShotInfos, snapShotInfo)
			}
		}

	}

	return snapShotInfos
}

func commitContainer(snapShotInfos []snapShot) {
	for _, snapShotInfo := range snapShotInfos {
		_, err := config.GlobalConfig.DockerClient.ContainerCommit(context.Background(), snapShotInfo.containerName, dockercontainer.CommitOptions{
			Reference: snapShotInfo.commitName,
		})
		if err != nil {
			common.StopProgram(err)
		}
	}
}

////////////////
/*  monitor   */
////////////////

func monitorPod() {
	watcher, err := config.GlobalConfig.KubernetesClient.CoreV1().Pods("default").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		common.StopProgram(err)
	}
	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			fmt.Println("Unexpected object type")
			continue
		}

		switch event.Type {
		case watch.Added:
			fmt.Printf("Pod %s added\n", pod.Name)
		case watch.Modified:
			fmt.Printf("Pod %s modified\n", pod.Name)
			for _, container := range pod.Spec.Containers {
				fmt.Printf("Pod name: %s\n", container.Name)
			}
		case watch.Deleted:
			fmt.Printf("Pod %s deleted\n", pod.Name)
			RebuildUbuntuContainer(config.GlobalConfig.KubernetesClient, pod.Name)
		case watch.Error:
			fmt.Printf("Error occurred: %v\n", event.Object)
		}
	}
}

func monitorService() {
	var info types.Info
	watcher, err := config.GlobalConfig.KubernetesClient.CoreV1().Services("default").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		common.StopProgram(err)
	}
	for event := range watcher.ResultChan() {
		service, ok := event.Object.(*corev1.Service)
		if !ok {
			fmt.Println("Unexpected object type")
			continue
		}

		switch event.Type {
		case watch.Added, watch.Modified:
			var id string
			if id, err = logger.GetStudentID(info, "serviceName"); id == "" {
				continue
			}
			info.StudentID = id
			info.ServiceName = service.Name
			info.NodePort = int(service.Spec.Ports[0].NodePort)
			logger.UpdateInfo(info, "serviceName")
			logger.UpdateInfo(info, "nodePort")
		case watch.Deleted:
			var id string
			if id, err = logger.GetStudentID(info, "serviceName"); id == "" {
				continue
			}
			info.StudentID = id
			logger.UpdateInfo(info, "serviceName")
			logger.UpdateInfo(info, "nodePort")
		case watch.Error:
			fmt.Printf("Error occurred: %v\n", event.Object)
		}
	}
}

////////////////////
/*  checkStatus   */
////////////////////
func CheckDeploymentStatus(deploymentName, namespace string) error {
	for i := 0; i < 10; i++ {
		deployment, err := config.GlobalConfig.KubernetesClient.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if deployment.Status.AvailableReplicas > 0 {
			fmt.Printf("Deployment is running.\n")
			return nil
		}

		fmt.Println("Deployment is not running yet. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("Deployment is not running.")
}

// TODO
func CheckServiceStatus(serviceName, namespace string) error {
	time.Sleep(10 * time.Second)

	fmt.Printf("Service is running.\n")
	return nil
}
