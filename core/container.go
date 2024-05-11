package core

import (
	"assign/types"
	"context"
	"fmt"

	"github.com/docker/docker/client"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

var (
	cNs              = make(map[int]types.ContainerNService) //name:count
	currentcontainer int
)

const totalcontainer = 100

func SetupEnvironment() {
	makeMap()
}

func makeMap() {
	cNs[0] = types.ContainerNService{Service: "service0", Container: 0}
	cNs[1] = types.ContainerNService{Service: "service1", Container: 0}
	cNs[2] = types.ContainerNService{Service: "service2", Container: 0}
	cNs[3] = types.ContainerNService{Service: "service3", Container: 0}
	cNs[4] = types.ContainerNService{Service: "service4", Container: 0}
	cNs[5] = types.ContainerNService{Service: "service5", Container: 0}
	cNs[6] = types.ContainerNService{Service: "service6", Container: 0}
	cNs[7] = types.ContainerNService{Service: "service7", Container: 0}
	cNs[8] = types.ContainerNService{Service: "service8", Container: 0}
	cNs[9] = types.ContainerNService{Service: "service9", Container: 0}
}

func CreateMariaDBContainer(cli *client.Client) error {
	//pv
	//pvc
	//dbdeploy
	//deservice
	return nil
}
func createServiceContainer(clientset *kubernetes.Clientset) error {
	for _, item := range cNs {

		service := &v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name: item.Service,
			},
			Spec: v1.ServiceSpec{
				Selector: map[string]string{
					"app": "ubuntu",
				},
				Ports: []v1.ServicePort{
					{
						Protocol:   "TCP",
						Port:       80,
						TargetPort: intstr.FromInt(8080),
					},
				},
				Type: v1.ServiceTypeLoadBalancer,
			},
		}

		_, err := clientset.CoreV1().Services("default").Create(context.TODO(), service, metav1.CreateOptions{})
		if err != nil {
			return err
		}

		// serviceinfo, err := clientset.CoreV1().Services("default").Get(context.TODO(), item.Service, metav1.GetOptions{})
		// if err != nil {
		// 	return err
		// }
		// serviceinfo.Spec.ClusterIP

	}
	return nil
}

func findAvailableNodePort(clientset *kubernetes.Clientset, startPort, endPort int32) (int32, error) {
	for port := startPort; port <= endPort; port++ {
		services, err := clientset.CoreV1().Services("default").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return 0, err
		}

		portUsed := false
		for _, service := range services.Items {
			for _, servicePort := range service.Spec.Ports {
				if servicePort.NodePort == port {
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

func CreateUbuntuContainer(clientset *kubernetes.Clientset) error {
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ubuntu-deployment",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "ubuntu",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ubuntu",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "ubuntu-container",
							Image: "ubuntu",
							Command: []string{
								"/bin/bash",
								"-c",
								"while true; do sleep 3600; done",
							},
						},
					},
				},
			},
		},
	}

	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return nil
}

func int32Ptr(i int32) *int32 { return &i }
