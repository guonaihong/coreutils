package link

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	command.Parse(argv[1:])

	args := command.Args()

	switch len(args) {
	case 0:
		fmt.Printf("link: missing operand\n")
	case 1:
		fmt.Printf("link: missing operand after '%s'\n", args[0])
	case 2:
		if err := os.Link(args[0], args[1]); err != nil {
			fmt.Printf("%s\n", err)
		}
	default:
		fmt.Printf("link: extra operand '%s'\n", args[2])

	}
}
