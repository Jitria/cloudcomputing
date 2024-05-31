package core

import (
	"assign/common"
	"assign/config"
	"assign/types"
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const startPort = 30000
const endPort = 32000

func MakeServer() types.Info {
	var info types.Info
	var err error

	if info.NodePort, err = findAvailableNodePort(config.GlobalConfig.KubernetesClient); err != nil {
		common.StopProgram(err)
	}

	if info.PodName, err = createUbuntuContainer(config.GlobalConfig.KubernetesClient); err != nil {
		common.StopProgram(err)
	}
	if info.ServiceName, err = createUbuntuService(config.GlobalConfig.KubernetesClient, info.NodePort, info.PodName); err != nil {
		common.StopProgram(err)
	}

	return info
}

//////////////////////////////
/*  findAvailableNodePort   */
//////////////////////////////
func findAvailableNodePort(clientset *kubernetes.Clientset) (int, error) {
	for port := startPort; port <= endPort; port++ {
		services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return 0, err
		}

		portUsed := false
		for _, service := range services.Items {
			for _, servicePort := range service.Spec.Ports {
				if servicePort.NodePort == int32(port) {
					portUsed = true
					break
				}
			}
			if portUsed {
				break
			}
		}

		if !portUsed {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in the range %d-%d", startPort, endPort)
}

/////////////////////////////
/*  createUbuntu&Service   */
/////////////////////////////
func createUbuntuService(clientset *kubernetes.Clientset, port int, podName string) (string, error) {
	name := common.GetServerName()
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": podName,
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       int32(22),
					TargetPort: intstr.FromInt(22),
					NodePort:   int32(port),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	_, err = clientset.CoreV1().Services("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return name, nil
}

func createUbuntuContainer(clientset *kubernetes.Clientset) (string, error) {
	name := common.GetPodName()
	podsClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  name,
					Image: "jitria/selfubuntu:latest",
					Command: []string{
						"/bin/bash",
						"-c",
						"service ssh start && sleep infinity",
					},
				},
			},
		},
	}

	_, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return name, nil
}

func RebuildUbuntuContainer(clientset *kubernetes.Clientset, imagename string) (string, error) {
	name := imagename
	imagename = imagename + ":latest"
	podsClient := clientset.CoreV1().Pods(corev1.NamespaceDefault)
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:            name,
					Image:           imagename,
					ImagePullPolicy: corev1.PullNever,
					Command: []string{
						"/bin/bash",
						"-c",
						"service ssh start && sleep infinity",
					},
				},
			},
		},
	}

	_, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}

	return name, nil
}

func int32Ptr(i int32) *int32 { return &i }
