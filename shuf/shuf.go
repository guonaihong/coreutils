package shuf

import (
	"bufio"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Shuf struct {
	Echo           *bool
	InputRange     *string
	HeadCount      *int
	Output         *string
	RandomSource   *string
	Repeat         *bool
	ZeroTerminated *bool
	lineDelim      byte
}

func New(argv []string) (*Shuf, []string) {

	shuf := &Shuf{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	shuf.Echo = command.Opt("e, echo",
		"treat each ARG as an input line").
		Flags(flag.PosixShort).NewBool(false)

	shuf.InputRange = command.Opt("i, input-range",
		"treat each number LO through HI as an input line").
		Flags(flag.PosixShort).NewString("")

	shuf.HeadCount = command.Opt("n, head-count",
		"output at most COUNT lines").
		Flags(flag.PosixShort).NewInt(-1)

	shuf.Output = command.Opt("o, output",
		"write result to FILE instead of standard output").
		Flags(flag.PosixShort).NewString("")

	shuf.RandomSource = command.Opt("random-source",
		"get random bytes from FILE").
		Flags(flag.PosixShort).NewString("")

	shuf.Repeat = command.Opt("r, repeat",
		"output lines can be repeated").
		Flags(flag.PosixShort).NewBool(false)

	shuf.ZeroTerminated = command.Opt("z, zero-terminated",
		"line delimiter is NUL, not newline").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	args := command.Args()

	return shuf, args
}

func (s *Shuf) readFromFile(name string, m map[string]struct{}) error {

	s.lineDelim = byte('\n')

	if s.ZeroTerminated != nil && *s.ZeroTerminated {
		s.lineDelim = byte('\000')
	}

	fd, err := utils.OpenInputFd(name)
	if err != nil {
		return err
	}

	defer utils.CloseInputFd(fd)

	br := bufio.NewReader(fd)

	for {

		l, e := br.ReadBytes(s.lineDelim)
		if e != nil && len(l) == 0 {
			break
		}

		m[string(l)] = struct{}{}
	}
	return nil
}

func (s *Shuf) Shuf(args []string, w io.Writer) error {

	m := map[string]struct{}{}

	if s.IsInputRange() {
		rs := strings.Split(*s.InputRange, "-")
		start, err := strconv.Atoi(rs[0])
		if err != nil {
			return fmt.Errorf(`invalid input range: "%s"`, rs[0])
		}

		end, err := strconv.Atoi(rs[1])
		if err != nil {
			return fmt.Errorf(`invalid input range: "%s"`, rs[1])
		}

		if start > end {
			return fmt.Errorf(`invalid input range: "%s"`, *s.InputRange)
		}

		for ; start <= end; start++ {
			m[fmt.Sprintf("%d\n", start)] = struct{}{}
		}
	}

	if s.IsEcho() {

		for _, v := range args {
			m[v] = struct{}{}
		}

	} else {
		s.readFromFile(args[0], m)
	}

	n := math.MaxInt64

	if s.HeadCount != nil && *s.HeadCount >= 0 {
		n = *s.HeadCount
	}

	for i := 0; i < n; i++ {

		for k := range m {
			w.Write([]byte(k))
			if s.Repeat != nil && !*s.Repeat {
				delete(m, k)
			}
			break
		}
	}

	return nil
}

func (s *Shuf) IsEcho() bool {
	return s.Echo != nil && *s.Echo
}

func (s *Shuf) IsInputRange() bool {
	return s.InputRange != nil && len(*s.InputRange) > 0
}

func (s *Shuf) CheckInputRange(rangeStr string) error {
	if strings.Index(rangeStr, "-") == -1 {
		return fmt.Errorf(`invalid input range: "%s"`, rangeStr)
	}
	return nil
}

func Main(argv []string) {

	shuf, args := New(argv)

	if !*shuf.Echo && len(args) > 1 {
		utils.Die(fmt.Sprintf("shuf: extra operand '%s'", args[1]))
	}

	w := os.Stdout

	if len(*shuf.Output) > 0 {
		fd, err := os.Open(*shuf.Output)
		if err != nil {
			utils.Die("shuf: %s", *shuf.Output)
		}

		w = fd
	}

	if shuf.IsInputRange() {
		if err := shuf.CheckInputRange(*shuf.InputRange); err != nil {
			utils.Die("shuf: " + err.Error())
		}
	}

	shuf.Shuf(args, w)
	utils.CloseOutputFd(w)
}
