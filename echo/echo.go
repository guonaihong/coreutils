package main

import (
	"fmt"
	"github.com/guonaihong/flag"
	"strconv"
)

func main() {
	newLine := flag.Bool("n", false, "do not output the trailing newline")
	enable := flag.Bool("e", false, "enable interpretation of backslash escapes")
	disable := flag.Bool("E", true, "disable interpretation of backslash escapes (default)")
	flag.Parse()

	args := flag.Args()

	c0 := uint64(0)
	var err error

	defer func() {
		if *newLine == false {
			fmt.Printf("\n")
		}
	}()

	if *enable {
		for _, s := range args {
			for i := 0; i < len(s); i++ {
				c := s[i]

				if c == '\\' && i < len(s) {
					i++
					if i >= len(s) {
						fmt.Printf("\\")
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
							goto notAnEscape
						}
						i++
						c0, err = strconv.ParseUint(s[i:], 16, 8)
						if err != nil {
							fmt.Print("\\")
							goto notAnEscape
						}
						i = len(s)
						c = byte(c0)
					case '0':
						if i+1 >= len(s) {
							goto notAnEscape
						}
						i++
						c0, err = strconv.ParseUint(s[i:], 8, 9)
						if err != nil {
							fmt.Print("\\")
							goto notAnEscape
						}

						i = len(s)
						c = byte(c0)
					default:
						fmt.Print("\\")
					}

				}

			notAnEscape:
				fmt.Printf("%c", c)
			}
			fmt.Print(" ")
		}
		return
	}

	if *disable {
		for i, s := range args {
			fmt.Print(s)
			if i+1 != len(args) {
				fmt.Printf(" ")
			}
		}
	}

}
