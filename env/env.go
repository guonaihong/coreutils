package env

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"os/exec"
	"strings"
)

type Env struct {
	IgnoreEnvironment bool   `opt:"i, ignore-environment" usage:"start with an empty environment"`
	Unset             string `opt:"u, unset" usage:"remove variable from the environment"`
	//Chdir             string `opt:"C, chdir" usage:"change working directory to DIR"`
	Null bool `opt:"0, null" usage:"end each output line with NUL, not newline"`
}

func Main(argv []string) {
	var env Env
	delimiter := byte('\n')

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	command.ParseStruct(argv[1:], &env)
	args := command.Args()

	if env.IgnoreEnvironment {
		os.Clearenv()
	}

	if env.Null {
		delimiter = byte(0)
	}

	for k, v := range args {
		if strings.IndexByte(v, '=') > 0 {
			kv := strings.SplitN(v, "=", 2)
			if err := os.Setenv(kv[0], kv[1]); err != nil {
				fmt.Printf("%s\n", err)
			}
			continue
		}

		cmd := exec.Command(args[k], args[k+1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	allArgs := os.Environ()
	for _, v := range allArgs {
		fmt.Fprintf(os.Stdout, "%s%c", v, delimiter)
	}
}
