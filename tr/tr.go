package tr

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"math"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type tr struct {
	set1Tab        map[byte]byte
	set2Tab        map[byte]byte
	set1Complement map[byte]byte
	squeezeRepeat  int64
	delete         bool
	truncateSet1   bool
	complement     bool
	class          setClass
}

const unused = math.MaxUint16

type setClass struct {
	tab map[string]*bytes.Buffer
}

func (s *setClass) init() {

	isAlpha := func(r rune) bool {
		return unicode.IsLower(r) || unicode.IsUpper(r)
	}

	isAlnum := func(r rune) bool {
		return isAlpha(r) || unicode.IsDigit(r)
	}

	isXdigit := func(r rune) bool {
		return unicode.IsDigit(r) || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F'
	}

	classFunc := map[string]func(r rune) bool{
		"alnum": isAlnum,
		"alpha": isAlpha,
		"cntrl": unicode.IsControl,
		//"digit":  unicode.IsDigit,
		"graph": unicode.IsGraphic,
		//"lower":  unicode.IsLower,
		"print": unicode.IsPrint,
		"punct": unicode.IsPunct,
		"space": unicode.IsSpace,
		//"upper":  unicode.IsUpper,
		"xdigit": isXdigit,
	}

	s.tab = map[string]*bytes.Buffer{}

	tab2 := map[string]*bytes.Buffer{
		"digit": bytes.NewBufferString("0123456789"),
		"lower": bytes.NewBufferString("abcdefghijklmnopqrstuvwxyz"),
		"upper": bytes.NewBufferString("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	}

	for className, f := range classFunc {

		for j := 0; j < 255; j++ {
			if !f(rune(j)) {
				continue
			}

			oldSet, ok := s.tab[className]
			if !ok {
				oldSet = &bytes.Buffer{}
				s.tab[className] = oldSet
			}
			oldSet.WriteByte(byte(j))
		}
	}

	for k, v := range tab2 {
		s.tab[k] = v
	}

	//fmt.Printf("#tab :%s\n", s.tab)
}

func (s *setClass) get(set string) string {
	for {
		leftSb := strings.Index(set, "[:")
		if leftSb == -1 {
			return set
		}

		rightSb := strings.Index(set, ":]")
		if rightSb == -1 {
			return set
		}

		className := set[leftSb+2 : rightSb]

		if newSet, ok := s.tab[className]; !ok {
			fmt.Printf("invalid character class:%s\n", className)
			os.Exit(1)
		} else {
			set = strings.Replace(set, set[leftSb:rightSb+2], newSet.String(), 1)
		}
	}

	return set
}

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

func parseSet(setTab map[byte]byte, set1, set2 string, truncateSet1 bool, parseRange1 bool) {

	findRange := func(i int, set string) bool {
		return i+1 < len(set) && i+2 < len(set) && set[i+1] == '-'
	}

	loopSet1, loopSet2 := false, false
	b1, b2 := byte(0), uint16(unused)
	lastB1, lastB2 := byte(0), byte(0)
	for i, j := 0, 0; i < len(set1); {

		if truncateSet1 && j >= len(set2) {
			return
		}

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

		if parseRange1 && findRange(i, set1) {
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

		/*
			fmt.Printf("b1:%c, b2:(%c), lastB1:(%c), lastB2:(%c), \n",
				b1, setTab[b1], lastB1, lastB2)
		*/

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

		if !loopSet1 {
			i++
		}
		if !loopSet2 && j < len(set2) {
			j++
		}

	}
}

func (t *tr) init(set1, set2 string) {
	t.class.init()

	t.class.get(set1)
	//fmt.Printf("set1:%s\n", t.class.get(set1))
	//fmt.Printf("set2:%s\n", t.class.get(set2))
	set1 = t.class.get(set1)
	set2 = t.class.get(set2)
	t.set1Tab = map[byte]byte{}
	t.set2Tab = map[byte]byte{}
	t.set1Complement = map[byte]byte{}
	parseSet(t.set1Tab, set1, set2, t.truncateSet1, true)
	parseSet(t.set2Tab, set2, "", t.truncateSet1, true)

	if t.complement {
		buf := bytes.Buffer{}
		for i := 0; i < 255; i++ {
			if _, ok := t.set1Tab[byte(i)]; ok {
				continue
			}
			buf.WriteByte(byte(i))
		}
		parseSet(t.set1Complement, buf.String(), set2, t.truncateSet1, false)
	}

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

func (t *tr) getComplement(b byte) byte {
	if _, ok := t.set1Tab[b]; ok {
		return b
	}

	if outByte, ok := t.set1Complement[b]; ok {
		return outByte
	}
	panic(fmt.Sprintf("unkown error:%c", b))
}

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	complement := command.Bool("c, C, complement", false, "use the complement of SET1")
	delete := command.Bool("d, delete", false, "delete characters in SET1, do not translate")
	squeezeRepeats := command.Bool("s, squeeze-repeats", false, "replace each sequence of a repeated character that is listed in the last specified SET, with a single occurrence of that character")
	truncateSet1 := command.Bool("t, truncate-set1", false, "first truncate SET1 to length of SET2")

	command.Parse(argv[1:])

	args := command.Args()

	var set1, set2 string
	if len(args) >= 1 {
		set1 = args[0]
	}

	if len(args) >= 2 {
		set2 = args[1]
	}

	tab := tr{
		delete:        *delete,
		truncateSet1:  *truncateSet1,
		squeezeRepeat: math.MaxInt64,
		complement:    *complement,
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

		if *delete {
			if tab.needDelete(oneByte[0]) {
				continue
			}
			goto output
		}

		if *complement {
			outByte := tab.getComplement(oneByte[0])
			oneByte[0] = outByte
			goto output
		}

		if *squeezeRepeats {
			if outByte, ok := tab.squeezeRepeats(oneByte[0]); ok {
				continue
			} else {
				oneByte[0] = outByte
			}

			goto output
		}

		oneByte[0] = tab.convert(oneByte[0])

	output:
		os.Stdout.Write(oneByte[:])
	}
}
