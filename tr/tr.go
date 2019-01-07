package main

import (
	"bufio"
	"github.com/guonaihong/flag"
	"os"
)

type tr struct {
	tab map[byte]byte
}

func (t *tr) init(set1, set2 string) {

	t.tab = map[byte]byte{}

	for i, j := 0, 0; i < len(set1); i++ {

		if i+1 < len(set1) && i+2 < len(set2) && set1[i] == '-' {
			j0 := j
			for i0 := set1[i]; i0 < set1[i+2]; i0++ {
				t.tab[set1[i0]] = set2[j0]
			}
			continue
		}

		t.tab[set1[i]] = set2[j]

		if j < len(set1) {
			j++
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
