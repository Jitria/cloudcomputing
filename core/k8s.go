package core

import (
	"assign/common"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitK8sClient() error {
	if common.IsK8sLocal() {
		return InitLocalAPIClient()
	}
	return nil
}

func InitInclusterAPIClient() error {
	read, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return err
	}
	K8sToken := string(read)

	kubeConfig := &rest.Config{
		Host:        "https://127.0.0.1:8001", // default address for kube-proxy
		BearerToken: K8sToken,
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: true,
		},
	}

	Kubeinfo.ClientSet, err = kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

func InitLocalAPIClient() error {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	if _, err := os.Stat(filepath.Clean(kubeconfig)); err != nil {
		return err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	Kubeinfo.ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	return nil
}
