package common

import (
	"fmt"
	"io"
	"os/exec"
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
