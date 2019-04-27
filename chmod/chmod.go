package chmod

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type Chmod struct {
	Changes        *bool
	Quiet          *bool
	Verbose        *bool
	NoPreserveRoot *bool
	PreserveRoot   *bool
	Reference      *string
	Recursive      *bool
}

func New(argv []string) (*Chmod, []string) {
	c := &Chmod{}

	command := flag.NewFlagSet(argv[0])

	c.Changes = command.Opt("c, changes",
		"like verbose but report only when a change is made")

	c.Quiet = command.Opt("f, silent, quiet",
		"suppress most error messages")

	c.Verbose = command.Opt("v, verbose",
		"output a diagnostic for every file processed")

	c.NoPreserveRoot = command.Opt("no-preserve-root",
		"do not treat '/' specially (the default)")

	c.PreserveRoot = command.Opt("preserve-root",
		"fail to operate recursively on '/'")

	c.Reference = command.Opt("reference",
		"use RFILE's mode instead of MODE values")

	c.Recursive = command.Opt("R, recursive",
		"change files and directories recursively")

	command.Parse(argv[1:])

	return c, command.Args()
}

func main() {
	fmt.Println("vim-go")
}
