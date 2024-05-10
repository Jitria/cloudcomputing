package core

import (
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
	cNs = make(map[string]int) //name:count
)

func CreateMariaDBContainer(cli *client.Client) error {
	makeMap()
	//pv
	//pvc
	//dbdeploy
	//deservice
	return nil
}
func makeMap() {
	cNs["service1"] = 0
	cNs["service2"] = 0
	cNs["service3"] = 0
	cNs["service4"] = 0
	cNs["service5"] = 0
	cNs["service6"] = 0
	cNs["service7"] = 0
	cNs["service8"] = 0
	cNs["service9"] = 0
	cNs["service10"] = 0
}
func CreateServiceContainer(clientset *kubernetes.Clientset) error {
	service := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "service1",
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

	fmt.Println("Service 'service1' created successfully.")
	return nil
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
