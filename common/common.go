package common

import (
	"fmt"
	"io"
	"log"
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
