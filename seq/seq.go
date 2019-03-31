package seq

import (
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"strconv"
	"strings"
)

type Seq struct {
	Format     *string
	Separator  *string
	EqualWidth *string

	first string
	step  string
	last  string
}

func New(argv []string) (*Seq, []string) {

	seq := Seq{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	seq.Format = command.Opt("f, format", "use printf style floating-point FORMAT").
		Flags(flag.PosixShort).NewString("")
	seq.Separator = command.Opt("s, separator",
		"use STRING to separate numbers (default: \\n)").
		Flags(flag.PosixShort).NewString("\n")
	seq.EqualWidth = command.Opt("w, equal-width",
		"equalize width by padding with leading zeroes").
		Flags(flag.PosixShort).NewString("")

	command.Parse(argv[1:])

	args := command.Args()
	if len(args) == 0 {
		utils.Die("seq: missing operand\n")
	}

	seq.first = "1"
	seq.step = "1"
	if len(args) >= 3 {
		seq.first = args[0]
		seq.step = args[1]
		seq.last = args[2]
	} else if len(args) >= 2 {
		seq.first = args[0]
		seq.last = args[1]
	} else if len(args) >= 1 {
		seq.last = args[0]
	}

	return &seq, args
}

func (s *Seq) check(format string) error {
	count := strings.Count(format, "%")
	if count == 0 {
		return fmt.Errorf(`seq: format '%s' has no %% directive`, format)
	}

	if count > 1 {
		return fmt.Errorf(`seq: format '%s' has too many %% directive`, format)
	}

	pos := strings.Index(format, "%")
	switch format[pos+1 : pos+2] {
	case "e", "f", "E", "F":
	//case "a", "A": //todo
	case "G", "g":
	default:
		return fmt.Errorf("format %s has unknown %%%c directive",
			format, format[pos+1])
	}
	return nil
}

func (s *Seq) Seq(w io.Writer) error {

	first, err := strconv.ParseFloat(s.first, 0)
	if err != nil {
		return fmt.Errorf("seq: invalid floating point argument: %s", s.first)
	}

	step, err := strconv.ParseFloat(s.step, 0)
	if err != nil {
		return fmt.Errorf("seq: invalid floating point argument: %s", s.step)
	}

	last, err := strconv.ParseFloat(s.last, 0)
	if err != nil {
		return fmt.Errorf("seq: invalid floating point argument: %s", s.last)
	}

	format := "%g"
	if s.Format != nil && len(*s.Format) > 0 {
		format = *s.Format
	}

	if err = s.check(format); err != nil {
		return err
	}

	var v interface{}
	separator := "\n"
	if s.Separator != nil {
		separator = *s.Separator
	}

	for ; first <= last; first += step {
		v = first
		fmt.Fprintf(w, format+separator, v)
	}

	return nil
}

func Main(argv []string) {
	s, _ := New(argv)

	err := s.Seq(os.Stdout)
	if err != nil {
		utils.Die("%s\n", err)
	}
}
