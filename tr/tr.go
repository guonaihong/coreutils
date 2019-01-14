package main

import (
	"bufio"
	"fmt"
	"github.com/guonaihong/flag"
	"math"
	"os"
	"strconv"
)

type tr struct {
	set1Tab       map[byte]byte
	set2Tab       map[byte]byte
	squeezeRepeat int64
	delete        bool
}

const unused = math.MaxUint16

func isoctal(b byte) bool {
	if b >= '0' && b <= '7' {
		return true
	}

	return false
}

func isoctalStr(s string, max int) (i int, haveOctal bool) {
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

func unquote(s string, i *int) byte {
	c := s[*i]
	if c == '\\' && *i < len(s) {
		(*i)++
		if *i >= len(s) {
			c = '\\'
			goto notAnEscape
		}

		c = s[*i]
		switch c {
		case 'a':
			c = '\a'
		case 'b':
			c = '\b'
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

		case '0', '1', '2', '3', '4', '5', '6', '7':
			if *i+1 >= len(s) {
				goto notAnEscape
			}

			n, haveOctal := isoctalStr(s[*i+1:], 3)
			if !haveOctal {
				goto notAnEscape
			}

			c0, err := strconv.ParseUint(s[*i+1:*i+1+n], 8, 32)
			if err != nil {
				goto notAnEscape
			}

			*i = *i + 1 + n - 1
			c = byte(c0)
		case '\\':
		default:
			c = '\\'
		}

	}

notAnEscape:
	return c
}

func parseSet(setTab map[byte]byte, set1, set2 string) {

	findRange := func(i int, set string) bool {
		return i+1 < len(set) && i+2 < len(set) && set[i+1] == '-'
	}

	loopSet1, loopSet2 := false, false
	b1, b2 := byte(0), uint16(unused)
	lastB1, lastB2 := byte(0), byte(0)

	for i, j := 0, 0; i < len(set1); {

		if !loopSet1 {
			//b1 = set1[i]
			b1 = unquote(set1, &i)
		}

		if !loopSet2 {
			if j < len(set2) {
				b2 = uint16(unquote(set2, &j))
				//b2 = uint16(set2[j])
			}
		}

		if findRange(i, set1) {
			loopSet1 = true
			i += 2 //skip start-
			lastB1 = unquote(set1, &i)
		}

		if findRange(j, set2) {
			loopSet2 = true
			j += 2 //skip start-
			lastB2 = unquote(set2, &j)
		}

		setTab[byte(b1)] = byte(b1)
		if b2 != unused {
			setTab[byte(b1)] = byte(b2)
		}
		//fmt.Printf("b1:%c, b2:(%c), lastB1:(%c), lastB2:(%c)\n", b1, setTab[b1], lastB1, lastB2)

		if loopSet1 {
			if b1 < lastB1 {
				b1++
			} else {
				loopSet1 = false
			}
		}

		if loopSet2 {
			if byte(b2) < lastB2 {
				b2++
			} else {
				loopSet2 = false
			}
		}

		//next:
		if !loopSet2 && j < len(set2) {
			j++
		}

		if !loopSet1 {
			i++
		}
	}
}

func (t *tr) init(set1, set2 string) {
	t.set1Tab = map[byte]byte{}
	t.set2Tab = map[byte]byte{}
	parseSet(t.set1Tab, set1, set2)
	parseSet(t.set2Tab, set2, "")
}

func (t *tr) convert(b byte) byte {
	b0, ok := t.set1Tab[b]
	if ok {
		return byte(b0)
	}

	return b
}

func (t *tr) needDelete(b byte) bool {
	_, ok := t.set1Tab[b]
	return ok
}

func (t *tr) squeezeRepeats(b byte) (byte, bool) {
	outByte, ok := t.set1Tab[b]
	if !ok {
		if outByte2, ok := t.set2Tab[b]; ok {
			outByte = outByte2
			goto next
		}
		outByte = b
		goto set
	}

next:
	if t.squeezeRepeat == math.MaxInt64 {
		t.squeezeRepeat = int64(outByte)
		return outByte, false
	}

	if byte(t.squeezeRepeat) == outByte {
		return 0, true
	}

set:
	t.squeezeRepeat = int64(outByte)

	return outByte, false
}

func main() {
	//complement := flag.Bool("c, C, complement", false, "use the complement of SET1")
	delete := flag.Bool("d, delete", false, "delete characters in SET1, do not translate")
	squeezeRepeats := flag.Bool("s, squeeze-repeats", false, "replace each sequence of a repeated character that is listed in the last specified SET, with a single occurrence of that character")
	//truncateSet1 := flag.Bool("t, truncate-set1", false, "first truncate SET1 to length of SET2")

	flag.Parse()

	args := flag.Args()

	var set1, set2 string
	if len(args) >= 1 {
		set1 = args[0]
	}

	if len(args) >= 2 {
		set2 = args[1]
	}

	tab := tr{
		delete:        *delete,
		squeezeRepeat: math.MaxInt64,
	}

	if *delete && len(args) >= 2 {
		if !*squeezeRepeats {
			fmt.Printf("extra operand:%s\n", args[len(args)-1])
			fmt.Println("Only one string may be given when ",
				"deleting without squeezing repeats")
			os.Exit(1)
		}
	}

	tab.init(set1, set2)

	stdin := bufio.NewReader(os.Stdin)

	var oneByte [1]byte

	for {

		_, err := stdin.Read(oneByte[:])
		if err != nil {
			break
		}

		if *squeezeRepeats {
			if outByte, ok := tab.squeezeRepeats(oneByte[0]); ok {
				continue
			} else {
				oneByte[0] = outByte
			}

			goto output
		}

		if *delete {
			if tab.needDelete(oneByte[0]) {
				continue
			}
			goto output
		}

		oneByte[0] = tab.convert(oneByte[0])

	output:
		os.Stdout.Write(oneByte[:])
	}
}
