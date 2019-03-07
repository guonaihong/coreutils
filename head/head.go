package head

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type headCmd struct {
	bytes     *int
	lines     *int
	quiet     *bool
	verbose   *bool
	lineDelim byte
}

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	head := headCmd{zeroTerminated}

	head.bytes = command.Opt("c, bytes", "print the first NUM bytes of each file;"+
		" with the leading '-', print all but the last NUM bytes of each file").
		Flags(flag.PosixShort).
		NewInt(0)

	head.lines = command.OptOpt(
		Flag{
			Regex: `^\d+$`,
			Short: []string{"l"},
			Long:  []string{"lines"},
			Usage: "print the first NUM lines instead of the first 10;" +
				"with the leading '-', print all but the last" +
				"NUM lines of each file"}).
		Flags(RegexKeyIsValue).
		NewInt(0)

	head.quiet = command.Opt("q, quiet, silent", "never print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	head.verbose = command.Opt("v, verbose", "always print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	zeroTerminated := command.Opt("z, zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).
		NewBool(false)

	if *zeroTerminated {
		head.lineDelim = '\000'
	}

	command.Parse(args[1:])
}
