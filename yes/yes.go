package yes

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"strings"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	v := command.Bool("version", false, "output version information and exit")
	if *v {
		fmt.Printf("todo output version\n")
		os.Exit(0)
	}

	command.Parse(argv[1:])
	args := command.Args()

	output := "y\n"

	if len(args) > 0 {
		output = strings.Join(args, " ")
		output += "\n"
	}

	for {
		os.Stdout.Write([]byte(output))
	}
}
