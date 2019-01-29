package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"sort"
	"strconv"
)

type sortLine struct {
	line      []byte
	number    int
	setNumber bool
}

func isoctal(b byte) bool {
	if b >= '0' && b <= '7' {
		return true
	}

	return false
}

type Month int8

const (
	Unknown Month = iota
	January
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

const emptyStr = "\t\n\v\f\r \x85\xA0"

var lineDelim byte = '\n'

func parseMonth(s []byte) Month {
	m := string(bytes.ToUpper(bytes.TrimLeft(s, emptyStr)))
	if len(m) >= 3 {
		m = m[:3]
	}

	switch m {
	case "JAN":
		return January
	case "FEB":
		return February
	case "MAR":
		return March
	case "APR":
		return April
	case "MAY":
		return May
	case "JUN":
		return June
	case "JUL":
		return July
	case "AUG":
		return August
	case "SEP":
		return September
	case "OCT":
		return October
	case "NOV":
		return November
	case "DEC":
		return December
	}

	return Unknown
}

func isOctalStr(s string, max int) (i int, haveOctal bool) {
	for i = 0; i < len(s); i++ {
		if i >= max {
			return i, haveOctal
		}

		if !isoctal(s[i]) {
			return i, haveOctal
		}

		haveOctal = true
	}

	return i, haveOctal
}

func (sl *sortLine) parseNumber() {
	if !sl.setNumber {
		defer func() { sl.setNumber = true }()

		line := sl.line
		nstr := string(line[:len(line)-1])

		n, haveOctal := isOctalStr(nstr, len(nstr))
		if !haveOctal {
			return
		}

		sl.number, _ = strconv.Atoi(nstr[:n])
		//fmt.Printf("%d--->%s\n", sl.number, sl.numberErr)
	}
}

func (sl *sortLine) isNumber() bool {
	return sl.setNumber
}

func main() {
	ignoreLeadingBlanks := flag.Bool("b, ignore-leading-blanks", false, "ignore leading blanks")
	dictionaryOrder := flag.Bool("d, dictionary-order", false, "consider only blanks and alphanumeric characters")
	ignoreCase := flag.Bool("f, ignore-case", false, "fold lower case to upper case characters")
	flag.String("g, general-numeric-sort", "", "compare according to general numerical value")
	flag.String("i, ignore-nonprinting", "", "consider only printable characters")
	monthSort := flag.Bool("M, month-sort", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	humanNumericSort := flag.Bool("h, human-numeric-sort", false, "compare human readable numbers (e.g., 2K 1G)")
	numericSort := flag.Bool("n, numeric-sort", false, "compare according to string numerical value")
	randomSort := flag.Bool("R, random-sort", false, "shuffle, but group identical keys.  See shuf(1)")
	flag.String("random-source", "", "get random bytes from FILE")
	reverse := flag.Bool("r, reverse", false, "reverse the result of comparisons")
	flag.String("sort", "", "sort according to WORD: general-numeric -g, human-numeric -h, month -M, numeric -n, random -R, version -V")
	flag.Bool("V, version-sort", false, "natural sort of (version) numbers within text")
	flag.String("batch-size", "", "merge at most NMERGE inputs at once; for more use temp files")
	flag.String("c, check, check", "", "check for sorted input; do not sort")
	flag.String("C, check=quiet, check=silent", "", "like -c, but do not report first bad line")
	flag.String("compress-program=PROG", "", "compress temporaries with PROG; decompress them with PROG -d")
	flag.String("debug", "", "annotate the part of the line used to sort, and warn about questionable usage to stderr")
	flag.String("files0-from=F", "", "read input from the files specified by NUL-terminated names in file F; If F is - then read names from standard input")
	flag.String("k, key=KEYDEF", "", "sort via a key; KEYDEF gives location and type")
	flag.String("m, merge", "", "merge already sorted files; do not sort")
	output := flag.String("o, output", "", "write result to FILE instead of standard output")
	flag.String("s, stable", "", "stabilize sort by disabling last-resort comparison")
	flag.String("S, buffer-size", "", "use SIZE for main memory buffer")
	flag.String("t, field-separator", "", "use SEP instead of non-blank to blank transition")
	flag.String("T, temporary-directory=DIR", "", "use DIR for temporaries, not $TMPDIR or /tmp; multiple options specify multiple directories")
	flag.String("parallel", "", "change the number of sorts run concurrently to N")
	unique := flag.Bool("u, unique", false, "with -c, check for strict ordering; without -c, output only the first of an equal run")
	zeroTerminated := flag.Bool("z, zero-terminated", false, "line delimiter is NUL, not newline")

	flag.Parse()
	args := flag.Args()

	if *zeroTerminated {
		lineDelim = byte(0)
	}

	defaultCmp := func(allLine []sortLine, i, j int) bool {
		cmp := func(allLine []sortLine, i, j int) bool {
			aLine, bLine := allLine[i].line, allLine[j].line

			if *ignoreCase {
				aLine = bytes.ToUpper(aLine)
				bLine = bytes.ToUpper(bLine)
			}

			if *ignoreLeadingBlanks {
				aLine = bytes.TrimLeft(aLine, emptyStr)
				bLine = bytes.TrimLeft(bLine, emptyStr)
			}

			if *humanNumericSort {
			}

			if *dictionaryOrder {
			}

			diff := 0
			if *monthSort {
				diff = int(parseMonth(aLine)) - int(parseMonth(bLine))
				if diff != 0 {
					return diff < 0
				}
			}

			if *numericSort {
				//fmt.Printf("i = %d, %d\n", i, len(allLine))
				allLine[i].parseNumber()
				allLine[j].parseNumber()

				if allLine[i].isNumber() && allLine[j].isNumber() {
					if allLine[i].number != allLine[j].number {
						return allLine[i].number < allLine[j].number
					}
				}

			}

			return bytes.Compare(aLine, bLine) < 0
		}

		if *reverse {
			return !cmp(allLine, i, j)
		}

		return cmp(allLine, i, j)
	}

	sort := func(r io.Reader, w io.Writer) {
		br := bufio.NewReader(r)

		var allLine []sortLine
		lineMap := map[int]struct{}{}
		uniqueMap := map[string]struct{}{}

		for count := 0; ; count++ {
			l, e := br.ReadBytes(lineDelim)
			if e == io.EOF {
				break
			}

			l = append(l, '\n')

			if *randomSort {
				lineMap[count] = struct{}{}
			}

			if *unique {
				key := string(l)
				if _, ok := uniqueMap[key]; ok {
					continue
				}
				uniqueMap[key] = struct{}{}
			}

			allLine = append(allLine, sortLine{line: append([]byte{}, l...)})
		}

		sort.Slice(allLine, func(i, j int) bool { return defaultCmp(allLine, i, j) })

		if len(lineMap) > 0 {
			for k, _ := range lineMap {
				w.Write(allLine[k].line)
			}

			return
		}

		for _, v := range allLine {
			w.Write(v.line)
		}
	}

	w := os.Stdout
	if len(*output) > 0 {
		fd, err := os.Create(*output)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		w = fd
		defer fd.Close()
	}

	if len(args) == 0 {
		sort(os.Stdin, w)
		return
	}

	for _, a := range args {
		fd, err := os.Open(a)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer fd.Close()

		sort(fd, w)
	}
}
