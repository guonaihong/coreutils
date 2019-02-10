package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"io"
	"math"
	"os"
)

var lineDelim byte = '\n'
var endLineDelim byte = '\n'

type ioCount struct {
	input  int
	output int
}

type uniq struct {
	count map[string]ioCount
}

func (u *uniq) init() {
	u.count = map[string]ioCount{}
}

func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func getSkipFields(skipFields int, l []byte) []byte {
	fields := bytes.Split(l, []byte(" "))
	// Index starts from 1
	if skipFields >= 0 && skipFields < len(fields) {
		return bytes.Join(fields[skipFields:], []byte(" "))
	}

	return []byte{}
}

func getSkipChars(skipChars int, l []byte) []byte {

	if skipChars >= 0 && skipChars < len(l) {
		return l[skipChars:]
	}

	return []byte{}
}

func getCheckChars(checkChars int, l []byte) []byte {
	if checkChars >= 0 && checkChars < len(l) {
		return l[:checkChars]
	}
	return l
}
func openInputFd(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdin, nil
	}

	return os.Open(fileName)
}

func openOutputFd(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdout, nil
	}

	return os.Create(fileName)
}

func closeInputFd(fd *os.File) {
	if fd != os.Stdin {
		fd.Close()
	}
}

func closeOutputFd(fd *os.File) {
	if fd != os.Stdout {
		fd.Close()
	}
}

func replaceEndLineDelim(l []byte) {
	if len(l) > 0 && l[len(l)-1] == '\n' && l[len(l)-1] != endLineDelim {
		l[len(l)-1] = endLineDelim
	}

}

func checkAllRepeated(arg string) {
	switch arg {
	case "none", "prepend", "separate":
	default:
		die("uniq: invalid argument `%s' for `--all-repeated' \nValid arguments are:\n  - `none'\n  - `prepend'\n  - `separate'\n", arg)
	}
}

func checkGroup(arg string) {
	switch arg {
	case "prepend", "append", "separate", "both":
	default:
		die("uniq: invalid argument `%s' for `--group' \nValid arguments are:\n  - `prepend'\n  - `append'\n  - `separate'\n  - `separate'\n", arg)
	}
}

func writeLine(w io.Writer, l []byte) {
	replaceEndLineDelim(l)
	w.Write(l)
}

func main() {
	count := flag.Bool("c, count", false, "prefix lines by the number of occurrences")
	repeated := flag.Bool("d, repeated", false, "only print duplicate lines")
	dup := flag.Bool("D", false, "print all duplicate lines")
	allRepeated := flag.String("all-repeated", "", "like -D, but allow separating groups with an empty line; METHOD={none(default),prepend,separate}")
	skipFields := flag.Int("f, skip-fields", math.MinInt32, "avoid comparing the first N fields")
	group := flag.String("group", "", "show all items, separating groups with an empty line; METHOD={separate(default),prepend,append,both}")
	ignoreCase := flag.Bool("i, ignore-case", false, "ignore differences in case when comparing")
	skipChars := flag.Int("s, skip-chars", math.MinInt32, "avoid comparing the first N characters")
	unique := flag.Bool("u, unique", false, "only print unique lines")
	zeroTerminated := flag.Bool("z, zero-terminated", false, "end lines with 0 byte, not newline")
	checkChars := flag.Int("w, check-chars", math.MinInt32, "compare no more than N characters in lines")

	flag.Parse()

	args := flag.Args()
	uniqHead := uniq{}
	uniqHead.init()

	if *dup && *allRepeated == "" {
		*allRepeated = "none"
	}

	if len(*group) > 0 {
		checkGroup(*group)
		if *count || *dup || *repeated || *unique {
			die("uniq: --group is mutually exclusive with -c/-d/-D/-u")
		}

		switch *group {
		case "prepend", "separate":
			*allRepeated = *group
		case "append":
			*allRepeated = "separate"
		case "both":
			*allRepeated = "separate"
		}
	}

	if len(*allRepeated) > 0 {
		checkAllRepeated(*allRepeated)
	}

	if *zeroTerminated {
		endLineDelim = '\000'
	}

	getKey := func(l []byte) string {
		if *ignoreCase {
			l = bytes.ToUpper(l)
		}

		if *skipFields != math.MinInt32 && *skipFields >= 0 {
			l = getSkipFields(*skipFields, l)
			//fmt.Printf("skip fields after l = (%s)\n", l)
		}

		if *skipChars != math.MinInt32 && *skipChars >= 0 {
			l = getSkipChars(*skipChars, l)
		}

		if *checkChars != math.MinInt32 && *checkChars >= 0 {
			l = getCheckChars(*checkChars, l)
		}

		//fmt.Printf("skipChars = %d, checkChars = %d\n", *skipChars, *checkChars)
		return string(l)
	}

	uniqCore := func(r io.Reader, w io.Writer) {
		br := bufio.NewReader(r)
		var allLine [][]byte
		var preLine []byte

		for lineNo := 0; ; lineNo++ {
			l, e := br.ReadBytes(lineDelim)
			if e != nil && len(l) == 0 {
				break
			}

			allLine = append(allLine, append([]byte{}, l...))

			key := getKey(l)
			ioCount, _ := uniqHead.count[key]
			ioCount.input++
			uniqHead.count[key] = ioCount
		}

		key := ""
		if len(allLine) > 0 && *group == "both" {
			w.Write([]byte{'\n'})
		}

		for _, l := range allLine {

			key = getKey(l)

			ioCount, ok := uniqHead.count[key]
			if !ok {
				panic("not foud" + string(l))
			}

			if *count {
				l = append([]byte(fmt.Sprintf("%6d  ", ioCount.input)), l...)
			}

			if *unique {
				if ioCount.input == 1 {
					goto write
				}
				goto next
			}

			if *repeated {
				if ioCount.input > 1 {
					goto write
				}
				goto next
			}

		write:
			if len(*allRepeated) == 0 {
				if ioCount, ok = uniqHead.count[key]; ok && ioCount.output > 0 {
					continue
				}
			} else {
				if *group == "" && ioCount.input == 1 {
					continue
				}
			}

			if preLine != nil {
				writeLine(w, append(preLine, '\n'))
				preLine = nil
			}

			switch *allRepeated {
			case "none":
			case "prepend":
				if ioCount.output == 0 {
					l = append([]byte{'\n'}, l...)
				}
			case "separate":
				if ioCount.output == ioCount.input-1 {

					preLine = append([]byte{}, l...)
					goto next
				}

			}

			writeLine(w, l)
		next:
			ioCount, _ = uniqHead.count[key]
			ioCount.output++
			uniqHead.count[key] = ioCount
		}

		if preLine != nil {
			if *group == "append" || *group == "both" {
				preLine = append(preLine, '\n')
			}
			writeLine(w, preLine)
		}
	}

	if len(args) == 0 {
		args = append(args, "-")
	}

	var r io.Reader
	var err error

	r, err = openInputFd(args[0])
	if err != nil {
		die("uniq: %s\n", err)
	}

	w := os.Stdout
	if len(args) == 2 {
		w, err = openOutputFd(args[1])
		if err != nil {
			die("uniq: %s\n", err)
		}
	}

	uniqCore(r, w)

	closeOutputFd(w)
}
