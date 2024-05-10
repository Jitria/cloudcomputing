package types

import (
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type GlobalConfig struct {
	CRI              string
	Gin              *gin.Engine
	KubernetesClient *kubernetes.Clientset
	DockerClient     *client.Client
}
