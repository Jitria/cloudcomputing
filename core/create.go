package core

import (
	"assign/common"
	"assign/config"
	"assign/types"
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

const startPort = 30000
const endPort = 32000

func MakeServer() types.Info {
	var port int
	var err error
	var info types.Info
	if port, err = findAvailableNodePort(config.GlobalConfig.KubernetesClient); err != nil {
		common.StopProgram(err)
	}
	name := common.GetNSName()
	if err = createUbuntuService(config.GlobalConfig.KubernetesClient, name, port); err != nil {
		common.StopProgram(err)
	}
	info.NodePort = port
	info.ServiceName = name
	if err = createUbuntuContainer(config.GlobalConfig.KubernetesClient, name); err != nil {
		common.StopProgram(err)
	}
	info.DeploymentName = name

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
func createUbuntuService(clientset *kubernetes.Clientset, name string, port int) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "ubuntu",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   "TCP",
					Port:       22,
					TargetPort: intstr.FromInt(22),
					NodePort:   int32(port),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}

	_, err := clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	_, err = clientset.CoreV1().Services("default").Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	return nil
}

func createUbuntuContainer(clientset *kubernetes.Clientset, name string) error {
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "ubuntu",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ubuntu",
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
			},
		},
	}

	_, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
