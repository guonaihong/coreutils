package echo

import (
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"strconv"
)

type Echo struct {
	NewLine *bool
	Enable  *bool
	Disable *bool
}

func New(argv []string) (*Echo, []string) {

	e := Echo{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	e.NewLine = command.Opt("n", "do not output the trailing newline").
		Flags(flag.PosixShort).NewBool(false)

	e.Enable = command.Opt("e", "enable interpretation of backslash escapes").
		Flags(flag.PosixShort).NewBool(false)

	e.Disable = command.Opt("E", "disable interpretation of backslash escapes (default)").
		Flags(flag.PosixShort).NewBool(true)

	command.Parse(argv[1:])

	args := command.Args()

	return &e, args
}

func (e *Echo) Echo(args []string, w io.Writer) {

	c0 := uint64(0)
	var err error

	defer func() {
		if e.NewLine != nil && *e.NewLine == false {
			w.Write([]byte{'\n'})
		}
	}()

	if e.Enable != nil && *e.Enable {

		printSlash := false
		for k, s := range args {
			for i := 0; i < len(s); i++ {
				c := s[i]

				if c == '\\' && i < len(s) {
					i++
					if i >= len(s) {
						w.Write([]byte{'\\'})
						goto notAnEscape
					}

					c = s[i]
					switch c {
					case 'a':
						c = '\a'
					case 'b':
						c = '\b'
					case 'c':
						return
					case 'e':
						c = '\x1B'
					case 'f':
						c = '\f'
					case 'n':
						c = '\n'
					case 'r':
						c = '\r'
					case 't':
						c = '\t'
					case 'v':
						c = '\v'
					case 'x':
						if i+1 >= len(s) {
							printSlash = true
							goto notAnEscape
						}

						n, haveHex := utils.IsXdigitStr(s[i+1:], 2)
						if !haveHex {
							printSlash = true
							goto notAnEscape
						}

						c0, err = strconv.ParseUint(s[i+1:i+1+n], 16, 32)
						if err != nil {
							printSlash = true
							goto notAnEscape
						}

						i = i + 1 + n - 1
						c = byte(c0)

					case '0':
						if i+1 >= len(s) {
							printSlash = true
							goto notAnEscape
						}

						n, haveOctal := utils.IsOctalStr(s[i+1:], 3)
						if !haveOctal {
							printSlash = true
							goto notAnEscape
						}

						c0, err = strconv.ParseUint(s[i+1:i+1+n], 8, 32)
						if err != nil {
							printSlash = true
							goto notAnEscape
						}

						i = i + 1 + n - 1
						c = byte(c0)
					case '\\':
					default:
						w.Write([]byte{'\\'})
					}

				}

			notAnEscape:
				if printSlash {
					w.Write([]byte{'\\'})
					printSlash = false
				}

				// fmt.Printf("%c") is not the same as the putchar output in c
				// in go fmt.Printf("%c\n", 172) -->  ¬
				// in c  putchar(172)            -->  ?
				w.Write([]byte{c})
			}
			if k+1 != len(args) {
				w.Write([]byte{' '})
			}
		}
		return
	}

	if e.Disable == nil || *e.Disable {
		for i, s := range args {
			w.Write([]byte(s))
			if i+1 != len(args) {
				w.Write([]byte{' '})
			}
		}
	}

}

func Main(argv []string) {
	echo, args := New(argv)
	echo.Echo(args, os.Stdout)
}
