package main

import (
	"fmt"
	"os/exec"
)

func main() {
	command := []string{
		"base32",
		"base64",
		"cat",
		"cut",
		"tail",
		"head",
	}

	for _, c := range command {
		runCmd := fmt.Sprintf("env GOPATH=`pwd` go test github.com/guonaihong/coreutils/%s", c)

		cmd := exec.Command("bash", "-c", runCmd)

		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(out))
	}

}
