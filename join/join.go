package join

import (
	"bufio"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
)

func Main(args []string) {
	command := flag.NewFlagSet(args[0], flag.ExitOnError)
	alsoNum := command.Int("a", 0, "also print unpairable lines from file FILENUM, "+
		"where FILENUM is 1 or 2, corresponding to FILE1 or FILE2")
	command.String("e", "", "replace missing input fields with EMPTY")
	command.Bool("i, ignore-case", false, "ignore differences in case "+
		"when comparing fields")
	command.Int("j", 0, "equivalent to '-1 FIELD -2 FIELD'")
	command.StringSlice("o", []string{}, "obey FORMAT while constructing output line")
	command.String("t", " ", "use CHAR as input and output field separator")
	command.String("v", "", "like -a FILENUM, but suppress joined output lines")
	command.String("1", "", "join on this FIELD of file 1")
	command.String("2", "", "join on this FIELD of file 2")
	command.Bool("check-order", false, "check that the input is correctly sorted, "+
		"even if all input lines are pairable")
	command.Bool("nocheck-order", false, "do not check that the input is correctly sorted")
	command.Bool("header", false, "treat the first line in each file as field headers,"+
		" print them without trying to pair them")
	command.Bool("z, zero-terminated", false, "line delimiter is NUL, not newline")

	args := command.Parse(args[1:])

	if len(args) != 2 {
		if len(args) == 0 {
			utils.Die("uniq: missing operand \n")
		} else {
			utils.Die("uniq:missing operand after %s\n", args[len(args)])
		}
	}
}
