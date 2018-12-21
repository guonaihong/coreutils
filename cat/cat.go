package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"strings"
)

func main() {
	numberNonblank := flag.Bool("b, number-nonblank", false, "number nonempty output lines")
	showEnds := flag.Bool("E, show-ends", false, "display $ at end of each line")
	number := flag.Bool("n, number", false, "number all output lines")
	squeezeBlank := flag.Bool("s, squeeze-blank", false, "suppress repeated empty output lines")
	showTables := flag.Bool("T, show-tables", false, "display TAB characters as ^I")

	flag.Parse()

	var r io.Reader

	r = os.Stdin

	args := flag.Args()
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("cat: %s\n", err)
			os.Exit(1)
		}

		defer f.Close()
        r = f
	}

	br := bufio.NewReader(r)

	var oldNew []string

	if *showEnds {
		oldNew = append(oldNew, "\n", "$\n")
	}

	if *showTables {
		oldNew = append(oldNew, "\t", "^I")
	}

	replacer := strings.NewReplacer(oldNew...)
	count := 1

	isSpace := 0
	for ; ; count++ {
		l, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}

        l = append(l, '\n')
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

		if *numberNonblank || *number {
			if *numberNonblank {
				count--
			}

			if !(*numberNonblank && len(l) == 1) {
				newLine := append([]byte{}, []byte(fmt.Sprintf("%6d  ", count))...)
				l = append(newLine, l...)
			}
		}

        os.Stdout.Write(l)
	}

}
