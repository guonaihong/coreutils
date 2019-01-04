package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"math"
	"os"
	"strconv"
	"strings"
)

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

func main() {
	bytes0 := flag.String("b, bytes", "", "select only these bytes")
	characters := flag.String("c, characters", "", "select only these characters")
	delimiter := flag.String("d, delimiter", "\t", "use DELIM instead of TAB for field delimiter")
	fields := flag.String("f, fields", "", "select only these fields;  also print any line that contains no delimiter character, unless the -s option is specified")
	complement := flag.Bool("complement", false, "complement the set of selected bytes, characters or fields")
	onlyDelimited := flag.Bool("s, only-delimited", false, "do not print lines not containing delimiters")
	outputDelimiter := flag.String("output-delimiter", "", "use STRING as the output delimiter the default is to use the input delimiter")
	zeroTerminated := flag.Bool("zero-terminated", false, "line delimiter is NUL, not newline")

	flag.Parse()

	args := flag.Args()

	lineDelim := byte('\n')

	checkFiledsNum := func() {
		filedsCount := 0

		if len(*bytes0) > 0 {
			filedsCount++
		}

		if len(*characters) > 0 {
			filedsCount++
		}

		if len(*fields) > 0 {
			filedsCount++
		}

		if filedsCount >= 2 {
			fmt.Printf("only one type of list may be specified\n")
			os.Exit(1)
		}

		if filedsCount == 0 {
			fmt.Printf("you must specify a list of bytes, characters, or fields\n")
			os.Exit(1)
		}
	}

	checkFiledsNum()

	if *zeroTerminated {
		lineDelim = byte(0)
	}

	if len(*bytes0) > 0 {
		*characters = *bytes0
	}

	filterFilter := filterCtrl{}
	if len(*fields) > 0 {
		filterFilter.init(*fields)
		*outputDelimiter = *delimiter
	} else {
		filterFilter.init(*characters)
	}

	cutCore := func(file *os.File) {
		reader := bufio.NewReader(file)
		buf := bytes.Buffer{}
		output := [][]byte{}
		byteOutput := []byte{}

		defer func() {
			if buf.Len() > 0 {
				os.Stdout.Write(buf.Bytes())
			}
		}()

		for {
			line, err := reader.ReadBytes(lineDelim)
			if err != nil {
				break
			}

			have := false
			if len(*fields) > 0 {
				ls := bytes.Split(line, []byte(*delimiter))
				if len(ls) == 1 {
					if *onlyDelimited {
						continue
					}
				}

				for i, v := range ls {
					checkOk := filterFilter.check(i + 1)
					if *complement {
						checkOk = !checkOk
					}

					if checkOk {
						have = true
						output = append(output, v)
					}
				}
				//todo
				buf.Write(bytes.Join(output, []byte(*outputDelimiter)))
				output = output[:0]

				goto write
			}

			for i, v := range line {
				checkOk := filterFilter.check(i + 1)
				if *complement {
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
				if buf.Bytes()[buf.Len()-1] != lineDelim {
					buf.WriteByte(lineDelim)
				}
			}

			if buf.Len() >= 0 {
				os.Stdout.Write(buf.Bytes())
				buf.Reset()
			}
		}
	}

	if len(args) == 0 {
		cutCore(os.Stdin)
		return
	}

	for _, v := range args {
		func() {
			fd, err := os.Open(v)
			if err != nil {
				fmt.Printf("%s\n", err)
				os.Exit(1)
			}

			defer fd.Close()
			cutCore(fd)
		}()
	}
}
