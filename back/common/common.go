package common

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
)

func GetCommandOutputWithoutErr(cmd string, args []string) string {
	res := exec.Command(cmd, args...)
	stdin, err := res.StdinPipe()
	if err != nil {
		return ""
	}

	go func() {
		defer func() {
			if err = stdin.Close(); err != nil {
				fmt.Printf(err.Error())
			}
		}()
		_, _ = io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := res.CombinedOutput()
	if err != nil {
		return ""
	}

	return string(out)
}

func IsK8sLocal() bool {
	home := os.Getenv("HOME")
	if _, err := os.Stat(filepath.Clean(home + "/.kube/config")); err == nil {
		return true
	}

	return false
}

func StopProgram(err error) {
	if err != nil {
		log.Printf("Error occurred: %s\n", err)
	}
	os.Exit(0)
}

var currentServer = "a"

func GetServerName() string {
	defer func() {
		var carry int
		for i := len(currentServer) - 1; i >= 0; i-- {
			if currentServer[i] < 'z' {
				currentServer = currentServer[:i] + string(currentServer[i]+1) + currentServer[i+1:]
				break
			} else {
				carry++
				currentServer = currentServer[:i] + "a" + currentServer[i+1:]
			}
		}
		if carry > 0 {
			currentServer = "a" + currentServer
		}
	}()
	return currentServer
}

var currentPod = "a"

func GetPodName() string {
	defer func() {
		var carry int
		for i := len(currentPod) - 1; i >= 0; i-- {
			if currentPod[i] < 'z' {
				currentPod = currentPod[:i] + string(currentPod[i]+1) + currentPod[i+1:]
				break
			} else {
				carry++
				currentPod = currentPod[:i] + "a" + currentPod[i+1:]
			}
		}
		if carry > 0 {
			currentPod = "a" + currentPod
		}
	}()
	return currentPod
}

func GetServerIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}
