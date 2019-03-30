package seq

import (
	"bytes"
	"testing"
)

func testFirstEndStep(dst string, first, step, last string, t *testing.T) {
	s := Seq{}
	s.first = "1"
	s.step = "1"
	s.last = "10"

	w := &bytes.Buffer{}

	s.Seq(w)

	if dst != w.String() {
		t.Errorf("seq fail (%s) need(%s)\n", w.String(), dst)
	}
}

func TestFirstEndStep(t *testing.T) {
	dst := `1
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
	testFirstEndStep(dst, "1", "1", "10", t)
}
