package config

import (
	"assign/common"
	"assign/types"
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/docker/docker/client"
)

var GlobalConfig types.GlobalConfig

func LoadConfig() {
	var err error
	if GlobalConfig.CRI = getCRISocket(); GlobalConfig.CRI == "" {
		common.StopProgram(fmt.Errorf("No CRI"))
	}
	GlobalConfig.Gin = gin.Default()
	if GlobalConfig.KubernetesClient, err = initLocalAPIClient(); err != nil {
		common.StopProgram(err)
	}
	if GlobalConfig.DockerClient, err = initDockerClient(); err != nil {
		common.StopProgram(err)
	}

}

/**********************/
/*   UpdateHostInfo   */
/**********************/

var ContainerRuntimeSocketMap = map[string][]string{
	"docker": {
		"/var/run/docker.sock",
		"/run/docker.sock",
	},
	"containerd": {
		"/var/snap/microk8s/common/run/containerd.sock",
		"/run/k3s/containerd/containerd.sock",
		"/run/containerd/containerd.sock",
		"/var/run/containerd/containerd.sock",
	},
}

func getCRISocket() string {
	for k := range ContainerRuntimeSocketMap {
		for _, candidate := range ContainerRuntimeSocketMap[k] {
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
	}
	return ""
}

func initLocalAPIClient() (*kubernetes.Clientset, error) {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	if _, err := os.Stat(filepath.Clean(kubeconfig)); err != nil {
		return nil, err
	}

	cf, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	ncf, err := kubernetes.NewForConfig(cf)
	if err != nil {
		return nil, err
	}

	return ncf, nil
}

func initDockerClient() (*client.Client, error) {
	var err error
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	client.NegotiateAPIVersion(context.Background())

	return client, nil
}
