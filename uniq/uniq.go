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

type uniq struct {
	input  map[string]int
	output map[string]struct{}
}

func (u *uniq) init() {
	u.input = map[string]int{}
	u.output = map[string]struct{}{}
}

func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func getSkipFields(skipFields int, l []byte) []byte {
	fields := bytes.Split(l, []byte(" "))
	// Index starts from 1
	if skipFields-1 >= 0 && skipFields-1 < len(fields) {
		return bytes.Join(fields[skipFields-1:], []byte(" "))
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

func main() {
	count := flag.Bool("c, count", false, "prefix lines by the number of occurrences")
	repeated := flag.Bool("d, repeated", false, "only print duplicate lines")
	flag.String("D, all-repeated", "", "print all duplicate lines delimit-method={none(default),prepend,separate} Delimiting is done with blank lines")
	skipFields := flag.Int("f, skip-fields", math.MinInt32, "avoid comparing the first N fields")
	ignoreCase := flag.Bool("i, ignore-case", false, "ignore differences in case when comparing")
	skipChars := flag.Int("s, skip-chars", math.MinInt32, "avoid comparing the first N characters")
	unique := flag.Bool("u, unique", false, "only print unique lines")
	zeroTerminated := flag.Bool("z, zero-terminated", false, "end lines with 0 byte, not newline")
	checkChars := flag.Int("w, check-chars", math.MinInt32, "compare no more than N characters in lines")

	flag.Parse()

	args := flag.Args()
	uniqHead := uniq{}
	uniqHead.init()

	if *skipFields < 0 {
		die("uniq: invalid number of fields to skip\n")
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
		}

		if *skipFields != math.MinInt32 && *skipChars >= 0 {
			l = getSkipFields(*skipChars, l)
		}

		if *checkChars != math.MinInt32 && *checkChars >= 0 {
		}
		return string(l)
	}

	uniqCore := func(r io.Reader, w io.Writer) {
		br := bufio.NewReader(r)
		var allLine [][]byte

		for lineNo := 0; ; lineNo++ {
			l, e := br.ReadBytes(lineDelim)
			if e != nil && len(l) == 0 {
				break
			}

			allLine = append(allLine, append([]byte{}, l...))

			key := getKey(l)
			uniqHead.input[key]++
		}

		key := ""
		for _, l := range allLine {

			key = getKey(l)

			v, ok := uniqHead.input[key]
			if !ok {
				panic("not foud" + string(l))
			}

			if *count {
				l = append([]byte(fmt.Sprintf("%6d  ", v)), l...)
			}

			if *unique {
				if v == 1 {
					goto write
				}
				goto next
			}

			if *repeated {
				if v > 1 {
					goto write
				}
				goto next
			}

		write:
			if _, ok := uniqHead.output[key]; ok {
				continue
			}

			replaceEndLineDelim(l)
			w.Write(l)
		next:
			uniqHead.output[key] = struct{}{}
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
