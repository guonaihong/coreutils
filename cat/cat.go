package cat

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"strings"
)

type Cat struct {
	NumberNonblank  *bool
	ShowEnds        *bool
	Number          *bool
	SqueezeBlank    *bool
	ShowTabs        *bool
	ShowNonprinting *bool
	oldNew          []string
}

func writeNonblank(l []byte) []byte {
	var out bytes.Buffer

	for _, c := range l {
		switch {
		case c == 9: // '\t'
			out.WriteByte(c)
		case c >= 0 && c <= 8 || c > 10 && c <= 31:
			out.Write([]byte{'^', c + 64})
		case c >= 32 && c <= 126 || c == 10: // 10 is '\n'
			out.WriteByte(c)
		case c == 127:
			out.Write([]byte{'^', c - 64})
		case c >= 128 && c <= 159:
			out.Write([]byte{'M', '-', '^', c - 64})
		case c >= 160 && c <= 254:
			out.Write([]byte{'M', '-', c - 128})
		default:
			out.Write([]byte{'M', '-', '^', 63})
		}
	}

	return out.Bytes()
}

func New(argv []string) (*Cat, []string) {
	c := Cat{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	showAll := command.Opt("A, show-all", "equivalent to -vET").
		Flags(flag.PosixShort).NewBool(false)

	c.NumberNonblank = command.Opt("b, number-nonblank",
		"number nonempty output lines, overrides -n").
		Flags(flag.PosixShort).NewBool(false)

	e := command.Opt("e", "equivalent to -vE").
		Flags(flag.PosixShort).NewBool(false)

	c.ShowEnds = command.Opt("E, show-end", "display $ at end of each line").
		Flags(flag.PosixShort).NewBool(false)

	c.Number = command.Opt("n, number", "number all output line").
		Flags(flag.PosixShort).NewBool(false)

	c.SqueezeBlank = command.Opt("s, squeeze-blank",
		"suppress repeated empty output lines").
		Flags(flag.PosixShort).NewBool(false)

	t := command.Opt("t", "equivalent to -vT").
		Flags(flag.PosixShort).NewBool(false)

	c.ShowTabs = command.Opt("T, show-tabs", "display TAB characters as ^I").
		Flags(flag.PosixShort).NewBool(false)

	c.ShowNonprinting = command.Opt("v, show-nonprinting",
		"use ^ and M- notation, except for LFD and TAB").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])
	args := command.Args()

	if *showAll {
		*c.ShowNonprinting = true
		*c.ShowEnds = true
		*c.ShowTabs = true
	}

	if *e {
		*c.ShowNonprinting = true
		*c.ShowEnds = true
	}
	if *t {
		*c.ShowNonprinting = true
		*c.ShowTabs = true
	}

	return &c, args
}

func SetBool(v bool) *bool {
	return &v
}

func (c *Cat) SetTab() {
	c.oldNew = append(c.oldNew, "\t", "^I")
}

func (c *Cat) SetEnds() {
	c.oldNew = append(c.oldNew, "\n", "$\n")
}

func (c *Cat) isNumberNonblank() bool {
	return c.NumberNonblank != nil && *c.NumberNonblank
}

func (c *Cat) isNumber() bool {
	return c.Number != nil && *c.Number
}

func (c *Cat) isShowNonprinting() bool {
	return c.ShowNonprinting != nil && *c.ShowNonprinting
}

func (c *Cat) isSqueezeBlank() bool {
	return c.SqueezeBlank != nil && *c.SqueezeBlank
}

func (c *Cat) Cat(rs io.ReadSeeker, w io.Writer) {
	br := bufio.NewReader(rs)
	replacer := strings.NewReplacer(c.oldNew...)
	isSpace := 0

	for count := 1; ; count++ {

		l, e := br.ReadBytes('\n')
		if e != nil && len(l) == 0 {
			break
		}

		if c.isSqueezeBlank() {
			if len(bytes.TrimSpace(l)) == 0 {
				isSpace++
			} else {
				isSpace = 0
			}

			if isSpace > 1 {
				count--
				continue
			}
		}

		if len(c.oldNew) > 0 {
			l = []byte(replacer.Replace(string(l)))
		}

		if c.isShowNonprinting() {
			l = writeNonblank(l)
		}

		if c.isNumberNonblank() || c.isNumber() {

			if c.isNumberNonblank() && len(l) == 1 {
				count--
			}

			if !(c.isNumberNonblank() && len(l) == 1) {
				l = append([]byte(fmt.Sprintf("%6d\t", count)), l...)
			}
		}

		w.Write(l)
	}
}

func Main(argv []string) {

	c, args := New(argv)

	if *c.ShowEnds {
		c.SetEnds()
	}

	if *c.ShowTabs {
		c.SetTab()
	}

	if len(args) > 0 {
		for _, fileName := range args {
			f, err := utils.OpenInputFd(fileName)
			if err != nil {
				utils.Die("cat: %s\n", err)
			}

			c.Cat(f, os.Stdout)
			utils.CloseInputFd(f)
		}
		return
	}
	c.Cat(os.Stdin, os.Stdout)
}
