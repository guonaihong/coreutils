package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/guonaihong/flag"
	"io"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type sortLine struct {
	line        []byte
	number      int
	floatNumber float64
	setNumber   bool
	setFloat    bool
}

type fieldSep map[rune]struct{}

func (f *fieldSep) init(s string) {
	*f = make(fieldSep, 10)
	for _, v := range s {
		(*f)[v] = struct{}{}
	}

}

func (f *fieldSep) is(r rune) bool {
	_, ok := (*f)[r]
	return ok
}

func isAlpha(r rune) bool {
	return unicode.IsLower(r) || unicode.IsUpper(r)
}

func isAlnum(r rune) bool {
	return isAlpha(r) || unicode.IsDigit(r)
}

func parsePrint(aLine []byte) (b []byte) {
	for _, v := range aLine {
		if unicode.IsPrint(rune(v)) {
			b = append(b, v)
		}
	}
	return
}

func parseDict(aLine []byte, f fieldSep) (b []byte) {

	for _, v := range aLine {
		if isAlnum(rune(v)) || f.is(rune(v)) {
			b = append(b, v)
		}
	}

	return
}

type Month int8

func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

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

func delLineBreaks(line []byte) []byte {

	if len(line) > 0 && line[len(line)-1] == '\n' {
		return line[:len(line)-1]
	}

	return line
}

func parseMonth(b []byte) Month {
	m := string(bytes.ToUpper(bytes.TrimLeft(b, emptyStr)))
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

type size int

const (
	KB size = 1024
	MB      = KB * 1024
	GB      = MB * 1024
	TB      = GB * 1024
	PB      = TB * 1024
)

func parseHumanNumberic(b []byte) int {
	s := string(bytes.TrimLeft(b, emptyStr))

	i, haveDecimal := isDecimalStr(s, len(s))
	if !haveDecimal {
		return 0
	}

	n, _ := strconv.Atoi(s[:i])

	suffix := strings.ToLower(s[i:])
	switch suffix {
	case "kb":
		return n * int(KB)
	case "mb":
		return n * int(MB)
	case "gb":
		return n * int(GB)
	case "tb":
		return n * int(TB)
	case "pb":
		return n * int(PB)
	default:
		return n
	}

	return n

}

func (sl *sortLine) parseGeneralNumericSort(b []byte) float64 {
	if !sl.setFloat {
		defer func() { sl.setFloat = true }()
		line := string(b)
		n, haveFloat := isFloatStr(line, len(line))
		if !haveFloat {
			return 0.0
		}

		sl.floatNumber, _ = strconv.ParseFloat(line[:n], 64)
		return sl.floatNumber
	}

	return sl.floatNumber
}

func isTypeStr(s string, max int, cb func(b byte) bool) (i int, haveStr bool) {
	for i = 0; i < len(s); i++ {
		if i >= max {
			return i, haveStr
		}

		if !cb(s[i]) {
			return i, haveStr
		}

		haveStr = true
	}

	return i, haveStr
}

func isFloatStr(s string, max int) (i int, haveFloat bool) {
	return isTypeStr(s, max, func(b byte) bool { return b >= '0' && b <= '9' || b == '.' || b == 'e' })
}

func isDecimalStr(s string, max int) (i int, haveDecimal bool) {

	return isTypeStr(s, max, func(b byte) bool { return b >= '0' && b <= '9' })
}

func (sl *sortLine) parseNumber(b []byte) int {
	if !sl.setNumber {
		defer func() { sl.setNumber = true }()

		line := string(b)

		n, haveDecimal := isDecimalStr(line, len(line))
		if !haveDecimal {
			return n
		}

		sl.number, _ = strconv.Atoi(line[:n])
	}

	return sl.number
}

func readFile(fileName string, body []byte) (n int, err error) {
	var fd *os.File
	fd, err = os.Open(fileName)
	if err != nil {
		die("sort: open fail: %s \n", err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fd.Seek(r.Int63n(math.MaxInt32), os.SEEK_SET)
	return fd.Read(body)
}

func getRandSource(fileName string) *rand.Rand {
	seed := int64(0)
	buf := make([]byte, 8)
	readFile(fileName, buf)
	read := bytes.NewReader(buf)
	binary.Read(read, binary.LittleEndian, &seed)
	return rand.New(rand.NewSource(seed))
}

func main() {
	ignoreLeadingBlanks := flag.Bool("b, ignore-leading-blanks", false, "ignore leading blanks")
	dictionaryOrder := flag.Bool("d, dictionary-order", false, "consider only blanks and alphanumeric characters")
	ignoreCase := flag.Bool("f, ignore-case", false, "fold lower case to upper case characters")
	generalNumericSort := flag.Bool("g, general-numeric-sort", false, "compare according to general numerical value")
	ignoreNonprinting := flag.Bool("i, ignore-nonprinting", false, "consider only printable characters")
	monthSort := flag.Bool("M, month-sort", false, "compare (unknown) < 'JAN' < ... < 'DEC'")
	humanNumericSort := flag.Bool("h, human-numeric-sort", false, "compare human readable numbers (e.g., 2K 1G)")
	numericSort := flag.Bool("n, numeric-sort", false, "compare according to string numerical value")
	randomSort := flag.Bool("R, random-sort", false, "shuffle, but group identical keys.  See shuf(1)")
	randomSource := flag.String("random-source", "", "get random bytes from FILE")
	reverse := flag.Bool("r, reverse", false, "reverse the result of comparisons")
	sortFlag := flag.Bool("sort", false, "sort according to WORD: general-numeric -g, human-numeric -h, month -M, numeric -n, random -R, version -V")
	versionSort := flag.Bool("V, version-sort", false, "natural sort of (version) numbers within text")
	flag.String("batch-size", "", "merge at most NMERGE inputs at once; for more use temp files")
	check := flag.Bool("c", false, "check for sorted input; do not sort")
	check2 := flag.Bool("C", false, "like -c, but do not report first bad line")
	flag.String("compress-program=PROG", "", "compress temporaries with PROG; decompress them with PROG -d")
	flag.String("debug", "", "annotate the part of the line used to sort, and warn about questionable usage to stderr")
	flag.String("files0-from=F", "", "read input from the files specified by NUL-terminated names in file F; If F is - then read names from standard input")
	flag.String("k, key=KEYDEF", "", "sort via a key; KEYDEF gives location and type")
	flag.String("m, merge", "", "merge already sorted files; do not sort")
	output := flag.String("o, output", "", "write result to FILE instead of standard output")
	stable := flag.Bool("s, stable", false, "stabilize sort by disabling last-resort comparison")
	flag.String("S, buffer-size", "", "use SIZE for main memory buffer")
	fieldSeparator := flag.String("t, field-separator", "", "use SEP instead of non-blank to blank transition")
	flag.String("T, temporary-directory=DIR", "", "use DIR for temporaries, not $TMPDIR or /tmp; multiple options specify multiple directories")
	flag.String("parallel", "", "change the number of sorts run concurrently to N")
	unique := flag.Bool("u, unique", false, "with -c, check for strict ordering; without -c, output only the first of an equal run")
	zeroTerminated := flag.Bool("z, zero-terminated", false, "line delimiter is NUL, not newline")

	flag.Parse()
	args := flag.Args()

	fieldSep0 := fieldSep{}
	fieldSep0.init(*fieldSeparator)

	if *zeroTerminated {
		lineDelim = byte(0)
	}

	if *sortFlag {
		*generalNumericSort = true
		*humanNumericSort = true
		*monthSort = true
		*numericSort = true
		*randomSort = true
		*versionSort = true
	}
	defaultCmp := func(allLine []sortLine, i, j int) bool {
		cmp := func(allLine []sortLine, i, j int) bool {
			aLine, bLine := allLine[i].line, allLine[j].line

			aLine = delLineBreaks(aLine)
			bLine = delLineBreaks(bLine)

			if *ignoreCase {
				aLine = bytes.ToUpper(aLine)
				bLine = bytes.ToUpper(bLine)
			}

			if *ignoreLeadingBlanks {
				aLine = bytes.TrimLeft(aLine, emptyStr)
				bLine = bytes.TrimLeft(bLine, emptyStr)
			}

			diff := 0
			if *humanNumericSort {
				diff = int(parseHumanNumberic(aLine)) - int(parseHumanNumberic(bLine))
				if diff != 0 {
					return diff < 0
				}
			}

			if *generalNumericSort {
				diff = int(allLine[i].parseGeneralNumericSort(aLine) - allLine[j].parseGeneralNumericSort(bLine))
				if diff != 0 {
					return diff < 0
				}
			}

			if *ignoreNonprinting {
				diff = bytes.Compare(parsePrint(aLine), parsePrint(bLine))
				if diff != 0 {
					return diff < 0
				}
			}

			if *dictionaryOrder {
				diff = bytes.Compare(parseDict(aLine, fieldSep0), parseDict(bLine, fieldSep0))
				if diff != 0 {
					return diff < 0
				}
			}

			if *monthSort {
				diff = int(parseMonth(aLine)) - int(parseMonth(bLine))
				if diff != 0 {
					return diff < 0
				}
			}

			if *numericSort {
				diff = allLine[i].parseNumber(aLine) - allLine[j].parseNumber(bLine)
				//fmt.Printf("%s, %s, %d, %d\n", aLine, bLine, allLine[i].number, allLine[j].number)
				if diff != 0 {
					return diff < 0
				}

			}

			return bytes.Compare(aLine, bLine) < 0
		}

		if *reverse {
			return !cmp(allLine, i, j)
		}

		return cmp(allLine, i, j)
	}

	sort := func(r io.Reader, w io.Writer, name string) {
		br := bufio.NewReader(r)
		allLine := []sortLine{}
		lineMap := map[int]struct{}{}
		uniqueMap := map[string]struct{}{}

		for count := 0; ; count++ {
			l, e := br.ReadBytes(lineDelim)
			if e != nil && len(l) == 0 {
				break
			}

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
			if (*check || *check2) && count >= 1 {
				if !defaultCmp(allLine, count-1, count) {
					if *check2 {
						os.Exit(1)
					}
					die("sort: %s:%d: disorder: %s", name, count+1, l)
				}
			}

			if e != nil {
				break
			}
		}

		sortFunc := sort.Slice
		if *stable {
			sortFunc = sort.SliceStable
		}

		sortFunc(allLine, func(i, j int) bool { return defaultCmp(allLine, i, j) })

		if len(*randomSource) > 0 && len(lineMap) > 0 {
			for len(lineMap) > 0 {
				r := getRandSource(*randomSource)
				k := 0
				for {
					k = int(r.Int63n(63))
					_, ok := lineMap[k]
					if !ok {
						continue
					}
					break

				}

				w.Write(allLine[k].line)
				delete(lineMap, k)
			}
			return
		}

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
		sort(os.Stdin, w, "-")
		return
	}

	for _, a := range args {
		fd, err := os.Open(a)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer fd.Close()

		sort(fd, w, a)
	}
}
