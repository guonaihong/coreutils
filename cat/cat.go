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

	var oldNew []string

	if *showEnds {
		oldNew = append(oldNew, "\n", "$\n")
	}

	if *showTables {
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

			if *numberNonblank || *number {
				if *numberNonblank {
					count--
				}

				if !(*numberNonblank && len(l) == 1) {
					l = append([]byte(fmt.Sprintf("%6d  ", count)), l...)
				}
			}

			os.Stdout.Write(l)
			if e != nil {
				break
			}
		}
	}

	args := flag.Args()
	if len(args) > 0 {
		for _, fileName := range args {
			f, err := os.Open(fileName)
			if err != nil {
				fmt.Printf("cat: %s\n", err)
				os.Exit(1)
			}

			catCore(f)
			f.Close()
		}
		return
	}
	catCore(os.Stdin)
}
