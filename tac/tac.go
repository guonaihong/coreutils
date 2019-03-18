package tac

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type Tac struct {
	Before    *bool
	Regex     *string
	Separator *string
}

func New(argv []string) (*Tac, []string) {
	t := Tac{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	t.Before = command.Opt("b, before", "attach the separator before instead of after").
		Flags(flag.PosixShort).
		NewBool(false)
	t.Regex = command.Opt("r, regex", "interpret the separator as a regular expression").
		Flags(flag.PosixShort).
		NewString("")
	t.Separator = command.Opt("s, separator", "use STRING as the separator instead of newline").
		Flags(flag.PosixShort).
		NewString("")

	command.Parse(argv[1:])
	args := command.Args()

	return &t, args
}

func (t *Tac) Tac(rs io.ReadSeeker, w io.Writer) {
}

func Main(argv []string) {
}
