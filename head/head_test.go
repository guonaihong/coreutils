package head

import (
	"bytes"
	"strings"
	"testing"
)

func TestCommandOption(t *testing.T) {
	h, _ := New([]string{"head", "-3", "./wc.go"})
	if h.Lines == nil || *h.Lines != 3 {
		t.Error("Command line parsing error, the required value is -3")
	}

	h, _ = New([]string{"head", "-n", "4", "./wc.go"})
	if h.Lines == nil || *h.Lines != 4 {
		t.Error("Command line parsing error, the required value is -n 4")
	}
}

func TestPrintBytes(t *testing.T) {
	//test head -c 3
	h := Head{}
	v := 3
	h.Bytes = &v

	str := "123456789"
	rs := strings.NewReader(str)
	w := &bytes.Buffer{}

	h.PrintBytes(rs, w)
	if w.String() != "123" {
		t.Errorf("need to be 123:(%s)\n", w.String())
	}

	//test head -c -3
	v = -3
	h.Bytes = &v

	rs = strings.NewReader(str)
	w = &bytes.Buffer{}

	h.PrintBytes(rs, w)
	if w.String() != "123456" {
		t.Errorf("need to be 789:(%s)\n", w.String())
	}

	//test head -c 3
	src := append([]byte("abc"), make([]byte, 1024*8)...)
	src = append(src, []byte("def")...)

	v = 3
	h.Bytes = &v

	brs := bytes.NewReader(src)
	w = &bytes.Buffer{}

	h.PrintBytes(brs, w)
	if w.String() != "abc" {
		t.Errorf("need to be abc:(%s)\n", w.String())
	}

	//test head -c 0
	v = 0
	h.Bytes = &v
	brs = bytes.NewReader(src)
	w = &bytes.Buffer{}

	h.PrintBytes(brs, w)
	if w.String() != "" {
		t.Errorf("need to be \"\":(%s)\n", w.String())
	}
}

func TestPrintLines(t *testing.T) {
	h := Head{LineDelim: '\n'}

	h10 := `1
2
3
4
5
6
7
8
9
10`

	rs := strings.NewReader(h10)
	w := &bytes.Buffer{}

	h.PrintLines(rs, w)
	if w.String() != h10 {
		t.Errorf("need to be \n(%s):\n(%s)\n", h10, w.String())
	}

	aZ := `a
	b
	c
	d
`
	last := `e`

	lines := -1
	h.Lines = &lines

	rs = strings.NewReader(aZ + last)
	w.Reset()

	h.PrintLines(rs, w)
	if w.String() != aZ {
		t.Errorf("need to be \n(%s):\n(%s)\n", aZ, w.String())
	}

}

func TestPrintTitle(t *testing.T) {
	verbose := true
	h := Head{Verbose: &verbose}
	w := &bytes.Buffer{}

	h.PrintTitle(w, "test.go")

	needStr := "==> test.go <==\n"
	if w.String() != needStr {
		t.Errorf("need to be (%s):(%s)\n", needStr, w.String())
	}
}
