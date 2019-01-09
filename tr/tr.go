package main

import (
	"bufio"
	"fmt"
	"github.com/guonaihong/flag"
	"os"
)

type tr struct {
	tab map[byte]byte
}

func (t *tr) init(set1, set2 string) {

	t.tab = map[byte]byte{}

	findRange := func(i int, set string) bool {
		return i-1 >= 0 && set[i] == '-' && i+1 < len(set)
	}

	loopSet1, loopSet2 := false, false
	b1, b2 := byte(0), byte(0)
	lastB1, lastB2 := byte(0), byte(0)

	haveStart := false
	for i, j := 0, 0; i < len(set1); {

		if !loopSet1 && i-1 >= 0 {
			b1 = set1[i-1]
			haveStart = true
		}

		if !loopSet2 && j-1 >= 0 {
			if j < len(set2) {
				b2 = set2[j-1]
			}
			haveStart = true
		}

		if findRange(i, set1) {
			loopSet1 = true
			i++ //skip -
			lastB1 = set1[i]
		}

		if findRange(j, set2) {
			loopSet2 = true
			j++ //skip -
			lastB2 = set2[j]
		}

		if !haveStart {
			//goto next
		}

		fmt.Printf("b1:%c, b2:%c, lastB1:%c, lastB2:%c, %d\n", b1, b2, lastB1, lastB2, ',')
		t.tab[byte(b1)] = byte(b2)

		if loopSet1 {
			if b1 <= lastB1 {
				b1++
			} else {
				loopSet1 = false
			}
		}

		if loopSet2 {
			if b2 <= lastB2 {
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

func (t *tr) convert(b byte) byte {
	b0, ok := t.tab[b]
	if ok {
		return b0
	}

	return b
}

func main() {
	//complement := flag.Bool("c, C, complement", false, "use the complement of SET1")
	//delete := flag.Bool("d, delete", false, "delete characters in SET1, do not translate")
	//squeezeRepeats := flag.Bool("s, squeeze-repeats", false, "replace each sequence of a repeated character that is listed in the last specified SET, with a single occurrence of that character")
	//truncateSet1 := flag.Bool("t, truncate-set1", false, "first truncate SET1 to length of SET2")

	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		return
	}

	set1 := args[0]
	set2 := args[1]

	tab := tr{}

	tab.init(set1, set2)

	stdin := bufio.NewReader(os.Stdin)

	var oneByte [1]byte
	var outByte [1]byte

	for {

		_, err := stdin.Read(oneByte[:])
		if err != nil {
			break
		}

		outByte[0] = tab.convert(oneByte[0])
		os.Stdout.Write(outByte[:])
	}
}
