package types

import (
	"database/sql"

	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
)

type GlobalConfig struct {
	CRI              string
	Gin              *gin.Engine
	KubernetesClient *kubernetes.Clientset
	DockerClient     *client.Client
	DB               *sql.DB
}

type ContainerNService struct {
	Service   string
	Container int
}
