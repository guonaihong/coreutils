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

	err := tail.PrintLines(rs, w)
	if w.String() != lines {
		t.Errorf("tail -n 10 fail (%s)\n", w.String())
	}

	l = "-10"
	tail.Lines = &l
	rs.Seek(0, 0)
	w.Reset()

	err = tail.PrintLines(rs, w)
	if w.String() != lines {
		t.Errorf("tail -n -10 fail (%s)\n", w.String())
	}

	if err != nil {
		t.Errorf("printLines %s\n", err)
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
`
	last := "10\n11"

	l = "+10"
	tail.Lines = &l
	rs = strings.NewReader(lines + last)
	w.Reset()
	err = tail.PrintLines(rs, w)

	if w.String() != last {
		t.Errorf("tail -n +10 fail (%s)\n", w.String())
	}
	if err != nil {
		t.Errorf("printLines %s \n", err)
	}
}

func TestLinesPlus(t *testing.T) {
	tail := Tail{}
	lines := `1
2
3
4
5
6
7
8
9
10
11
`
	l := "+0"
	tail.Lines = &l
	rs := strings.NewReader(lines)
	w := &bytes.Buffer{}
	err := tail.PrintLines(rs, w)

	if w.String() != lines {
		t.Errorf("tail -n +0 fail (%s)\n", w.String())
	}
	if err != nil {
		t.Errorf("printLines %s \n", err)
	}

	l = "+1"
	tail.Lines = &l
	rs = strings.NewReader(lines)
	w = &bytes.Buffer{}
	err = tail.PrintLines(rs, w)

	if w.String() != lines {
		t.Errorf("tail -n +0 fail (%s)\n", w.String())
	}
	if err != nil {
		t.Errorf("printLines %s \n", err)
	}
}
