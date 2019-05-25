package main

import (
	"fmt"
	"os/exec"
)

func main() {
	command := []string{
		"base32",
		"base64",
		"basename",
		"cat",
		"cut",
		"dirname",
		"echo",
		"head",
		"paste",
		"pwd",
		"tee",
		"tail",
		"tr",
		"true",
		"uniq",
		"whoami",
		"yes",
		"sleep",
		"tac",
	}

	for _, c := range command {
		runCmd := fmt.Sprintf("env GOPATH=`pwd` go build github.com/guonaihong/coreutils/%s/%s", c, c)

		fmt.Printf("%s\n", runCmd)
		cmd := exec.Command("bash", "-c", runCmd)

		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(out))
	}

}
