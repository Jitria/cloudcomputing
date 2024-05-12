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

type Person struct {
	ID string `json:"id"`

	// == //
	Position string `json:"position"`
}

type Info struct {
	NodePort       int    `json:"nodePort"`
	DeploymentName string `json:"deploymentName"`
	ServiceName    string `json:"serviceName"`

	// == //
	StudentID string `json:"studentID"`
	Ip        string `json:"ip"`
}
