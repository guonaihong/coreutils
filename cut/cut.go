package cut

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Cut struct {
	Bytes           string
	Characters      string
	Delimiter       string
	Fields          string
	Complement      bool
	OnlyDelimited   bool
	OutputDelimiter string
	LineDelim       byte
	filterCtrl
}

func New(argv []string) (*Cut, []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	c := Cut{}

	command.Opt("b, bytes", "select only these bytes").
		Flags(flag.PosixShort).
		Var(&c.Bytes)

	command.Opt("c, characters", "select only these characters").
		Flags(flag.PosixShort).
		Var(&c.Characters)

	command.Opt("d, delimiter", "use DELIM instead of TAB for field delimiter").
		Flags(flag.PosixShort).
		DefaultVar(&c.Delimiter, "\t")

	command.Opt("f, fields", "select only these fields;  also print any line\n"+
		" that contains no delimiter character, unless\n"+
		" the -s option is specified").
		Flags(flag.PosixShort).
		Var(&c.Fields)

	command.Opt("complement", "complement the set of selected bytes, characters\n"+
		"or fields").
		Flags(flag.PosixShort).
		Var(&c.Complement)

	command.Opt("s, only-delimited", "do not print lines not containing delimiters").
		Flags(flag.PosixShort).
		Var(&c.OnlyDelimited)

	command.Opt("output-delimiter", "use STRING as the output delimiter\n"+
		"the default is to use the input delimiter").
		Flags(flag.PosixShort).
		Var(&c.OutputDelimiter)

	zeroTerminated := command.Opt("zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).
		NewBool(false)

	command.Parse(argv[1:])

	args := command.Args()

	c.LineDelim = byte('\n')
	if *zeroTerminated {
		c.LineDelim = '\000'
	}

	return &c, args
}

type paragraph struct {
	start, end int
}

type filterCtrl struct {
	filter map[int]struct{}
	p      []paragraph
}

func die(filter string) {
	fmt.Printf(`cut: invalid field value "%s"\n`, filter)
	os.Exit(1)
}

func (f *filterCtrl) init(filter string) {
	var err error
	f.filter = map[int]struct{}{}
	s := strings.Split(filter, ",")

	start, end := 0, 0

	for _, v := range s {
		startEnd := strings.Split(v, "-")

		if len(startEnd) == 2 {
			p := paragraph{}

			if len(startEnd[0]) == 0 {
				p.start = 1
			}

			if len(startEnd[1]) == 0 {
				p.end = math.MaxInt64
			}

			if len(startEnd[0]) > 0 {
				start, err = strconv.Atoi(startEnd[0])
				if err != nil {
					die(filter)
				}

				if p.start == 0 && start > 0 {
					p.start = start
				}

				if start > 0 && start < p.start {
					p.start = start
				}
			}

			if len(startEnd[1]) > 0 {
				end, err = strconv.Atoi(startEnd[1])
				if err != nil {
					die(filter)
				}
				if end > p.end {
					p.end = end
				}
			}

			f.p = append(f.p, p)
			continue
		}

		n, err := strconv.Atoi(v)
		if err != nil {
			die(v)
		}

		f.filter[n] = struct{}{}
	}
}

func (f *filterCtrl) check(index int) (ok bool) {

	for _, p := range f.p {
		if index >= p.start && index <= p.end {
			return true
		}
	}

	_, ok = f.filter[index]
	return ok
}

func (c *Cut) Cut(rs io.ReadSeeker, w io.Writer) {
	reader := bufio.NewReader(rs)
	buf := bytes.Buffer{}
	output := [][]byte{}
	byteOutput := []byte{}

	defer func() {
		if buf.Len() > 0 {
			w.Write(buf.Bytes())
		}
	}()

	for {
		line, err := reader.ReadBytes(c.LineDelim)
		if err != nil && len(line) == 0 {
			break
		}

		have := false
		if c.isFields() {
			ls := bytes.Split(line, []byte(c.Delimiter))
			if len(ls) == 1 {
				if c.OnlyDelimited {
					continue
				}
				buf.Write(line)
				goto write
			}

			for i, v := range ls {
				checkOk := c.check(i + 1)
				if c.Complement {
					checkOk = !checkOk
				}

				if checkOk {
					have = true
					output = append(output, v)
				}
			}

			buf.Write(bytes.Join(output, []byte(c.OutputDelimiter)))
			output = output[:0]

			goto write
		}

		for i, v := range line {
			checkOk := c.check(i + 1)
			if c.Complement {
				checkOk = !checkOk
			}

			if checkOk {
				have = true
				byteOutput = append(byteOutput, v)
			}
		}

		buf.Write(byteOutput)
		byteOutput = byteOutput[:0]

	write:
		if have {
			if buf.Bytes()[buf.Len()-1] != c.LineDelim && line[len(line)-1] == c.LineDelim {
				buf.WriteByte(c.LineDelim)
			}
		}

		if buf.Len() >= 0 {
			w.Write(buf.Bytes())
			buf.Reset()
		}
	}
}

func (c *Cut) isBytes() bool {
	return len(c.Bytes) > 0
}

func (c *Cut) isCharacters() bool {
	return len(c.Characters) > 0
}

func (c *Cut) isFields() bool {
	return len(c.Fields) > 0
}

func (c *Cut) isDelimiter() bool {
	return len(c.Delimiter) > 0
}

func (c *Cut) isOutputDelimiter() bool {
	return len(c.OutputDelimiter) > 0
}

func (c *Cut) Init() {
	checkFiledsNum := func() {
		filedsCount := 0

		if c.isBytes() {
			filedsCount++
		}

		if c.isCharacters() {
			filedsCount++
		}

		if c.isFields() {
			filedsCount++
		}

		if filedsCount >= 2 {
			utils.Die("only one type of list may be specified\n")
		}

		if filedsCount == 0 {
			utils.Die("you must specify a list of bytes, characters, or fields\n")
		}
	}

	checkFiledsNum()

	if c.isBytes() {
		c.Characters = c.Bytes
	}

	if c.isFields() {
		c.init(c.Fields)
		if c.isDelimiter() && c.OutputDelimiter == "" {
			c.OutputDelimiter = c.Delimiter
		}
		return
	}

	c.init(c.Characters)

}

func Main(argv []string) {

	c, args := New(argv)

	c.Init()
	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, v := range args {
		fd, err := utils.OpenFile(v)
		if err != nil {
			utils.Die("cut: %s\n", err)
		}

		c.Cut(fd, os.Stdout)
		fd.Close()
	}
}
