package tail

import (
	"bytes"
	"strings"
	"testing"
)

func TestLines(t *testing.T) {
	tail := Tail{}

	one := "1\n"
	lines := `2
3
4
5
6
7
8
9
10
11`

	l := "10"
	tail.Lines = &l
	tail.LineDelim = '\n'

	rs := strings.NewReader(one + lines)
	w := &bytes.Buffer{}

	tail.PrintLines(rs, w)
	if w.String() != lines {
		t.Errorf("tail -n 10 fail (%s)\n", w.String())
	}

	l = "-10"
	tail.Lines = &l
	rs.Seek(0, 0)
	w.Reset()

	tail.PrintLines(rs, w)
	if w.String() != lines {
		t.Errorf("tail -n -10 fail (%s)\n", w.String())
	}

	lines = `1
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
	last := "11"

	rs = strings.NewReader(lines + last)
	w.Reset()

	if w.String() != last {
		t.Errorf("tail -n +10 fail (%s)\n", w.String())
	}
}
