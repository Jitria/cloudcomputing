package config

import (
	"assign/types"
	"os"

	"github.com/gin-gonic/gin"
)

var GlobalConfig types.GlobalConfig

func LoadConfig() {
	GlobalConfig.CRI = getCRISocket()
	GlobalConfig.Gin = setupRouter()
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	return r
}
