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

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	showAll := command.Opt("A, show-all", "equivalent to -vET").
		Flags(flag.PosixShort).NewBool(false)

	numberNonblank := command.Opt("b, number-nonblank",
		"number nonempty output lines, overrides -n").
		Flags(flag.PosixShort).NewBool(false)

	e := command.Opt("e", "equivalent to -vE").
		Flags(flag.PosixShort).NewBool(false)

	showEnds := command.Opt("E, show-end", "display $ at end of each line").
		Flags(flag.PosixShort).NewBool(false)

	number := command.Opt("n, numbe", "number all output line").
		Flags(flag.PosixShort).NewBool(false)

	squeezeBlank := command.Opt("s, squeeze-blank",
		"suppress repeated empty output lines").
		Flags(flag.PosixShort).NewBool(false)

	t := command.Opt("t", "equivalent to -vT").
		Flags(flag.PosixShort).NewBool(false)

	showTabs := command.Opt("T, show-tabs", "display TAB characters as ^I").
		Flags(flag.PosixShort).NewBool(false)

	showNonprinting := command.Opt("v, show-nonprinting",
		"use ^ and M- notation, except for LFD and TAB").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	var oldNew []string

	if *showAll {
		*showNonprinting = true
		*showEnds = true
		*showTabs = true
	}

	if *e {
		*showNonprinting = true
		*showEnds = true
	}

	if *t {
		*showNonprinting = true
		*showTabs = true
	}

	if *showEnds {
		oldNew = append(oldNew, "\n", "$\n")
	}

	if *showTabs {
		oldNew = append(oldNew, "\t", "^I")
	}

	catCore := func(r io.Reader) {
		br := bufio.NewReader(r)
		replacer := strings.NewReplacer(oldNew...)

		isSpace := 0
		for count := 1; ; count++ {
			l, e := br.ReadBytes('\n')
			if e != nil && len(l) == 0 {
				break
			}

			if *squeezeBlank {
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

			if len(oldNew) > 0 {
				l = []byte(replacer.Replace(string(l)))
			}

			if *showNonprinting {
				l = writeNonblank(l)
			}

			if *numberNonblank || *number {
				if *numberNonblank {
					count--
				}

				if !(*numberNonblank && len(l) == 1) {
					l = append([]byte(fmt.Sprintf("%6d  ", count)), l...)
				}
			}

			os.Stdout.Write(l)
		}
	}

	args := command.Args()
	if len(args) > 0 {
		for _, fileName := range args {
			f, err := utils.OpenInputFd(fileName)
			if err != nil {
				utils.Die("cat: %s\n", err)
			}

			catCore(f)
			utils.CloseInputFd(f)
		}
		return
	}
	catCore(os.Stdin)
}
