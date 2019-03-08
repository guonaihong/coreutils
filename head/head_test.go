package head

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintBytes(t *testing.T) {
	//test head -c 3
	h := head{}
	v := 3
	h.bytes = &v

	str := "123456789"
	rs := strings.NewReader(str)
	w := &bytes.Buffer{}

	h.printBytes(rs, w)
	if w.String() != "123" {
		t.Errorf("need to be 123:(%s)\n", w.String())
	}

	//test head -c -3
	v = -3
	h.bytes = &v

	rs = strings.NewReader(str)
	w = &bytes.Buffer{}

	h.printBytes(rs, w)
	if w.String() != "123456" {
		t.Errorf("need to be 789:(%s)\n", w.String())
	}

	//test head -c 3
	src := append([]byte("abc"), make([]byte, 1024*8)...)
	src = append(src, []byte("def")...)

	v = 3
	h.bytes = &v

	brs := bytes.NewReader(src)
	w = &bytes.Buffer{}

	h.printBytes(brs, w)
	if w.String() != "abc" {
		t.Errorf("need to be abc:(%s)\n", w.String())
	}

	//test head -c 0
	v = 0
	h.bytes = &v
	brs = bytes.NewReader(src)
	w = &bytes.Buffer{}

	h.printBytes(brs, w)
	if w.String() != "" {
		t.Errorf("need to be \"\":(%s)\n", w.String())
	}
}
