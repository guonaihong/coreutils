package env

import (
	"github.com/guonaihong/flag"
	"os"
)

type Env struct {
	IgnoreEnvironment bool   `opt:"i, ignore-environment" usage:"start with an empty environment"`
	Unset             string `opt:"u, unset" usage:"remove variable from the environment"`
	Chdir             string `opt:"C, chdir" usage:"change working directory to DIR"`
	Null              byte   `opt:"0, null" usage:"end each output line with NUL, not newline"`
}

func Main(argv []string) {
	var env Env

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	command.ParseStruct(argv[1:], &env)
	args := command.Args()

	if env.IgnoreEnvironment {
		os.Clearenv()
	}

	for _, v := range args {
		kv, err := strings.Split(v, "=")
		if err != nil {
			//todo
		}
		if err = os.Setenv(kv[0], kv[1]); err != nil {
			//todo
		}
	}

	allArgs := os.Environ()
	for _, v := range allArgs {
		fmt.Fprintf(os.Stdout, "%s\n", v)
	}
}
