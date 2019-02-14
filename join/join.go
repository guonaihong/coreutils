package join

import (
	"bufio"
	"bytes"
	_ "fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
)

var lineDelim byte = '\n'

func delEndLineDelim(l *[]byte) {
	if (*l)[len(*l)-1] == lineDelim {
		*l = (*l)[:len(*l)-1]
	}
}

type joinCmd struct {
	key1            *int
	key2            *int
	printUnpairable *int
	separator       *string

	r1, r2 io.Reader
	w      io.Writer
}

func (j *joinCmd) addOutLine(outLine *bytes.Buffer, ls [][]byte, key int) {
	for k, v := range ls {
		if k == key {
			continue
		}

		outLine.Write(v)
		outLine.WriteString(*j.separator)
	}

}

func (j *joinCmd) getField(l1, l2 []byte) {
	delEndLineDelim(&l1)
	ls1 := bytes.Split(l1, []byte(*j.separator))
	ls2 := bytes.Split(l2, []byte(*j.separator))
	outLine := bytes.NewBuffer(nil)

	key1 := *j.key1 - 1
	key2 := *j.key2 - 1
	if bytes.Equal(ls1[key1], ls2[key2]) {
		outLine.Write(ls1[key1])
		outLine.WriteString(*j.separator)

		j.addOutLine(outLine, ls1, key1)
		j.addOutLine(outLine, ls2, key2)
		line := outLine.Bytes()
		j.w.Write(line[:len(line)-len(*j.separator)])
	}
}

func (j *joinCmd) main() {
	br1 := bufio.NewReader(j.r1)
	br2 := bufio.NewReader(j.r2)

	fileEof1, fileEof2 := false, false
	var l1, l2 []byte
	var err error

	for {

		if !fileEof1 {
			l1, err = br1.ReadBytes(lineDelim)
			if err != nil && len(l1) == 0 {
				fileEof1 = true
			}
		}

		if !fileEof2 {
			l2, err = br2.ReadBytes(lineDelim)
			if err != nil && len(l2) == 0 {
				fileEof2 = true
			}
		}

		if fileEof1 && fileEof2 {
			break
		}

		j.getField(l1, l2)
	}
}

func Main(argv []string) {

	cmdOpt := joinCmd{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	cmdOpt.printUnpairable = command.Int("a", 0, "also print unpairable lines from file FILENUM, "+
		"where FILENUM is 1 or 2, corresponding to FILE1 or FILE2")
	command.String("e", "", "replace missing input fields with EMPTY")
	command.Bool("i, ignore-case", false, "ignore differences in case "+
		"when comparing fields")
	command.Int("j", 0, "equivalent to '-1 FIELD -2 FIELD'")
	command.StringSlice("o", []string{}, "obey FORMAT while constructing output line")
	cmdOpt.separator = command.String("t", " ", "use CHAR as input and output field separator")
	command.String("v", "", "like -a FILENUM, but suppress joined output lines")
	cmdOpt.key1 = command.Int("1", 1, "join on this FIELD of file 1")
	cmdOpt.key2 = command.Int("2", 1, "join on this FIELD of file 2")
	command.Bool("check-order", false, "check that the input is correctly sorted, "+
		"even if all input lines are pairable")
	command.Bool("nocheck-order", false, "do not check that the input is correctly sorted")
	command.Bool("header", false, "treat the first line in each file as field headers,"+
		" print them without trying to pair them")
	command.Bool("z, zero-terminated", false, "line delimiter is NUL, not newline")

	command.Parse(argv[1:])

	args := command.Args()
	if len(args) != 2 {
		if len(args) == 0 {
			utils.Die("uniq: missing operand \n")
		} else {
			utils.Die("uniq:missing operand after %s\n", args[len(args)])
		}
	}

	fd1, err := utils.OpenInputFd(args[0])
	if err != nil {
		utils.Die("join: %s\n", err)
	}

	fd2, err := utils.OpenInputFd(args[1])
	if err != nil {
		utils.Die("join: %s\n", err)
	}

	outFd, err := utils.OpenOutputFd("-")
	if err != nil {
		utils.Die("join: %s\n", err)
	}

	cmdOpt.r1 = fd1
	cmdOpt.r2 = fd2
	cmdOpt.w = outFd

	cmdOpt.main()

	utils.CloseInputFd(fd1)
	utils.CloseInputFd(fd2)
	utils.CloseOutputFd(outFd)

}
