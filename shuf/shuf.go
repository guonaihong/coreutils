package shuf

import (
	"bufio"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"math/rand"
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

type result struct {
	m map[int]struct{}
	a [][]byte
}

func (r *result) init() {
	if r.m == nil {
		r.m = map[int]struct{}{}
		r.a = make([][]byte, 0, 10)
	}
}

func (r *result) add(no int, padding int, a []byte) {
	r.m[no-padding] = struct{}{}
	r.a = append(r.a, a)
}

func (r *result) del(no int) {
	delete(r.m, no)
}

func (s *Shuf) readFromFile(name string, r *result) error {

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

	for no := 0; ; no++ {

		l, e := br.ReadBytes(s.lineDelim)
		if e != nil && len(l) == 0 {
			break
		}

		r.add(no, 0, l)
	}

	return nil
}

func (s *Shuf) Shuf(args []string, w io.Writer) error {

	r := &result{}

	r.init()

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

		i := start
		for ; start <= end; start++ {
			r.add(start, i, []byte(fmt.Sprintf("%d\n", start)))
		}

	} else if s.IsEcho() {

		for k, v := range args {
			r.add(k, 0, append([]byte(v), '\n'))
		}

	} else {
		s.readFromFile(args[0], r)
	}

	n := len(r.a)

	if s.HeadCount != nil && *s.HeadCount >= 0 {
		n = *s.HeadCount
	}

	var rnd *rand.Rand
	var err error

	if s.IsRandSource() {
		rnd, err = utils.GetRandSource(*s.RandomSource)
		if err != nil {
			return err
		}
	}

	for i := 0; ; i++ {

		if !s.IsRepeat() && i >= n {
			break
		}

		k := 0

		if rnd != nil {
			k = int(rnd.Int63n(int64(len(r.a))))
			w.Write([]byte(r.a[k]))
			goto next
		}

		for k = range r.m {

			w.Write([]byte(r.a[k]))

			break
		}

	next:
		if !s.IsRepeat() {
			r.del(k)
		}
	}

	return nil
}

func (s *Shuf) IsRepeat() bool {
	return s.Repeat != nil && *s.Repeat
}

func (s *Shuf) IsRandSource() bool {
	return s.RandomSource != nil && len(*s.RandomSource) > 0
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
		utils.Die(fmt.Sprintf("shuf: extra operand '%s'\n", args[1]))
	}

	if len(*shuf.InputRange) > 0 && len(args) > 0 {
		utils.Die(fmt.Sprintf("shuf: extra operand '%s'\n", args[0]))
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
