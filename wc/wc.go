package wc

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type wcCmd struct {
	bytesCounts     *bool
	characterCounts *bool
	lines           *bool
	files0From      *bool
	maxLineLength   *bool
	words           *bool
}

func Main(args []string) {
	command := flag.NewFlagSet(args[0], flag.ContinueOnError)

	wc := wcCmd{}

	wc.bytesCounts = command.Opt("c, bytes", "print the byte counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.characterCounts = command.Opt("m, chars", "print the character counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.lines = command.Opt("print the newline counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.files0From = command.Opt("files0-from", "read input from the files specified by\n"+
		"NUL-terminated names in file F;\n"+
		"If F is - then read names from standard input").
		Flags(flag.PosixShort).NewString("")

	wc.maxLineLength = command.Opt("L, max-line-length", "print the maximum display width").
		Flags(flag.PosixShort).NewBool(false)

	wc.words = command.Opt("w, words", "print the word counts").
		Flags(flag.PosixShort).NewBool(false)

}
