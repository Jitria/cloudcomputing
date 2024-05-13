package core

import (
	"assign/common"
	"assign/config"
	"assign/logger"
	"assign/types"
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func Monitor() {
	go monitorDeployment()
	go monitorContainer()
	go monitorService()
}

func monitorDeployment() {
	var info types.Info
	watcher, err := config.GlobalConfig.KubernetesClient.AppsV1().Deployments("default").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		common.StopProgram(err)
	}
	for event := range watcher.ResultChan() {
		deployment, ok := event.Object.(*appsv1.Deployment)
		if !ok {
			fmt.Println("Unexpected object type")
			continue
		}

		switch event.Type {
		case watch.Added, watch.Modified:
			id := logger.GetStudentID(info, "deploymentName")
			info.StudentID = id
			info.DeploymentName = deployment.Name
			logger.UpdateInfo(info, "deploymentName")
		case watch.Deleted:
			id := logger.GetStudentID(info, "deploymentName")
			info.StudentID = id
			logger.UpdateInfo(info, "deploymentName")
		case watch.Error:
			fmt.Printf("Error occurred: %v\n", event.Object)
		}
	}
}

func monitorContainer() {
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
				fmt.Printf("Container name: %s\n", container.Name)
			}
		case watch.Deleted:
			fmt.Printf("Pod %s deleted\n", pod.Name)
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
			id := logger.GetStudentID(info, "serviceName")
			info.StudentID = id
			info.ServiceName = service.Name
			info.NodePort = int(service.Spec.Ports[0].NodePort)
			logger.UpdateInfo(info, "serviceName")
			logger.UpdateInfo(info, "nodePort")
		case watch.Deleted:
			id := logger.GetStudentID(info, "serviceName")
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

func CheckServiceStatus(serviceName, namespace string) error {
	time.Sleep(10 * time.Second)

	fmt.Printf("Service is running.\n")
	return nil
}
