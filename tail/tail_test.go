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

func TestPrintBytes(t *testing.T) {
	testBytes := `123456789abcdef`

	tail := Tail{}
	rs := strings.NewReader(testBytes)
	w := &bytes.Buffer{}

	b := "1"
	tail.Bytes = &b

	err := tail.PrintBytes(rs, w)
	if w.String() != "f" {
		t.Errorf("tail -c 1 fail(%s)\n", w.String())
	}

	if err != nil {
		t.Errorf("PrintBytes %s\n", err)
	}

	b = "-1"
	tail.Bytes = &b
	rs.Seek(0, 0)
	w.Reset()
	err = tail.PrintBytes(rs, w)
	if w.String() != "f" {
		t.Errorf("tail -c -1 (%s)\n", w.String())
	}

	b = "0"
	tail.Bytes = &b
	rs.Seek(0, 0)
	w.Reset()

	err = tail.PrintBytes(rs, w)
	if w.String() != "" {
		t.Errorf("tail -c 0 fail(%s)\n", w.String())
	}

	b = "-0"
	tail.Bytes = &b
	rs.Seek(0, 0)
	w.Reset()

	err = tail.PrintBytes(rs, w)
	if w.String() != "" {
		t.Errorf("tail -c -0 fail(%s)\n", w.String())
	}
}

func TestPrintBytesPlus(t *testing.T) {
	testBytes := `123456789abcdef`
	rs := strings.NewReader(testBytes)
	w := &bytes.Buffer{}

	tail := Tail{}

	b := "+1"
	tail.Bytes = &b
	err := tail.PrintBytes(rs, w)
	if w.String() != "123456789abcdef" {
		t.Errorf("tail -c +1 fail(%s)\n", w.String())
	}

	if err != nil {
		t.Errorf("PrintLines %s\n", err)
	}

	b = "+2"
	tail.Bytes = &b
	rs.Seek(0, 0)
	w.Reset()

	err = tail.PrintBytes(rs, w)
	if w.String() != "23456789abcdef" {
		t.Errorf("tail -c + 2 fail(%s)\n", w.String())
	}

	if err != nil {
		t.Errorf("PrintLines %s\n", err)
	}

	b = "+0"
	tail.Bytes = &b
	rs.Seek(0, 0)
	w.Reset()

	err = tail.PrintBytes(rs, w)
	if w.String() != "123456789abcdef" {
		t.Errorf("tail -c +0 fail(%s)\n", err)
	}
}
