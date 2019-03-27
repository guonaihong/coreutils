package tac

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"strings"
	"testing"
)

func TestBefore(t *testing.T) {
}

func TestRegex(t *testing.T) {
}

// tac -s string
func testReadFromTailStdin(src, dst string, sep string, t *testing.T) {
	rs := strings.NewReader(src)
	w := &bytes.Buffer{}

	readFromTailStdin(rs, w, []byte(sep))

	if w.String() != dst {
		t.Fatalf("tac -s fail(%s, %d), need(%s)\n", w.String(), w.Len(), dst)
	}
}

func testSeparator(src, dst string, sep string, t *testing.T) {
	tac := Tac{}
	tac.Separator = utils.String(sep)
	rs := strings.NewReader(src)
	w := &bytes.Buffer{}

	tac.Tac(rs, w)

	if w.String() != dst {
		t.Fatalf("tac -s fail(%s), need(%s)\n", w.String(), dst)
	}
}

func TestSeparator(t *testing.T) {
	src := `1
2
3
4
5
6
7
8
9
10
`
	dst := `10
9
8
7
6
5
4
3
2
1
`
	testReadFromTailStdin(src, dst, "\n", t)

	//testSeparator(src, dst, "\n", t)

	src = "123aaa456aaa789aaa\n"

	dst = `
789aaa456aaa123aaa`
	testReadFromTailStdin(src, dst, "aaa", t)
	//testSeparator(src, dst, "aaa", t)

	src = "wwwwwwwwwwww"
	testReadFromTailStdin(src, src, "aaa", t)

	src = "1,2\n"
	testReadFromTailStdin(src, src, "\n", t)

	dst = `2
1,`
	testReadFromTailStdin(src, dst, ",", t)
}
