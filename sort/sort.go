package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"sort"
)

type sortLine struct {
	line []byte
}

func main() {
	flag.String("b, ignore-leading-blanks", "", "ignore leading blanks")
	flag.String("d, dictionary-order", "", "consider only blanks and alphanumeric characters")
	flag.String("f, ignore-case", "", "fold lower case to upper case characters")
	flag.String("g, general-numeric-sort", "", "compare according to general numerical value")
	flag.String("i, ignore-nonprinting", "", "consider only printable characters")
	flag.String("M, month-sort", "", "compare (unknown) < 'JAN' < ... < 'DEC'")
	flag.String("h, human-numeric-sort", "", "compare human readable numbers (e.g., 2K 1G)")
	flag.String("n, numeric-sort", "", "compare according to string numerical value")
	flag.String("R, random-sort", "", "shuffle, but group identical keys.  See shuf(1)")
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
	flag.String("o, output", "", "write result to FILE instead of standard output")
	flag.String("s, stable", "", "stabilize sort by disabling last-resort comparison")
	flag.String("S, buffer-size", "", "use SIZE for main memory buffer")
	flag.String("t, field-separator=SEP", "", "use SEP instead of non-blank to blank transition")
	flag.String("T, temporary-directory=DIR", "", "use DIR for temporaries, not $TMPDIR or /tmp; multiple options specify multiple directories")
	flag.String("parallel", "", "change the number of sorts run concurrently to N")
	flag.String("u, unique", "", "with -c, check for strict ordering; without -c, output only the first of an equal run")
	flag.String("z, zero-terminated", "", "line delimiter is NUL, not newline")

	args := flag.Args()

	defaultCmp := func(allLine []sortLine, i, j int) bool {
		return bytes.Compare(allLine[i].line, allLine[j].line) < 0
	}

	sort := func(r io.Reader) {
		br := bufio.NewReader(r)

		var allLine []sortLine

		for {
			l, _, e := br.ReadLine()
			if e == io.EOF {
				break
			}

			allLine = append(allLine, sortLine{line: append([]byte{}, l...)})
			last := allLine[len(allLine)-1]
			last.line = append(last.line, '\n')
			allLine[len(allLine)-1] = last
		}

		sort.SliceStable(allLine, func(i, j int) bool {
			if *reverse {
				return !defaultCmp(allLine, i, j)
			}

			return defaultCmp(allLine, i, j)
		})

		for _, v := range allLine {
			os.Stdout.Write(v.line)
		}
	}

	if len(args) == 0 {
		sort(os.Stdin)
		return
	}

	for _, a := range args {
		fd, err := os.Open(a)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}

		defer fd.Close()

		sort(fd)
	}
}
