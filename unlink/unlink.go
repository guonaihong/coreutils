package unlink

import (
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"syscall"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	command.Parse(argv[1:])

	args := command.Args()

	if len(args) == 0 {
		utils.Die("unlink: missing operand\n")
	}

	err := syscall.Unlink(args[0])
	if err != nil {
		fmt.Printf("unlink: cannot unlink '%s': %s\n", args[0], err)
	}
}
