package cat

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestNumber(t *testing.T) {
	testNumber :=
		`1
	2
	3
	4
	5
	6`

	var lastLine []byte
	c := Cat{}

	c.Number = true

	in := strings.NewReader(testNumber)
	out := &bytes.Buffer{}

	c.Cat(in, out)

	br := bufio.NewReader(strings.NewReader(out.String()))
	number := 0

	for {
		l, err := br.ReadBytes('\n')
		if err != nil && len(l) == 0 {
			break
		}
		lastLine = l
	}

	fmt.Sscanf(string(lastLine), "%d", &number)
	if number != 6 {
		t.Fatalf("cat -n fail lastLine:%d", number)
	}

}

func TestTab(t *testing.T) {
	testTab := "\t\t\t\t\t\t\t\t\n\t\t\t\t\t\n"
	c := Cat{}
	c.SetTab()

	in := strings.NewReader(testTab)
	out := &bytes.Buffer{}

	c.Cat(in, out)

	outStr := out.String()
	outStr = strings.Replace(outStr, "^I", "", -1)
	outStr = strings.Replace(outStr, "\n", "", -1)

	if len(outStr) != 0 {
		t.Fatalf("cat -T fail (%s)", outStr)

	}
}

func TestEnds(t *testing.T) {
	testEnds := `1
2
3`
	c := Cat{}
	c.SetEnds()

	in := strings.NewReader(testEnds)
	out := &bytes.Buffer{}

	c.Cat(in, out)

	ls := strings.Split(out.String(), "\n")
	for _, v := range ls {
		if v == "3" {
			continue
		}

		if !strings.HasSuffix(v, "$") {
			t.Fatalf("cat -E fail (%s)\n", v)
		}
	}
}

func TestSqueezeBlank(t *testing.T) {
	c := Cat{}
	c.SqueezeBlank = true

	rs := strings.NewReader("\n\n\n12\n\n\n\n34\n\n\n")
	w := &bytes.Buffer{}
	c.Cat(rs, w)

	outStr := `
12

34

`
	if w.String() != outStr {
		t.Fatalf("cat -s fail(%s)\n", w.String())
	}
}

func TestNumberNonblank(t *testing.T) {
	c := Cat{}
	c.NumberNonblank = true

	rs := strings.NewReader("\n\n\n12\n\n\n\n34\n\n\n")
	w := &bytes.Buffer{}
	c.Cat(rs, w)
	outStr := `


     1	12



     2	34


`
	if w.String() != outStr {
		t.Fatalf("cat -b fail(%s), need(%s)\n", w.String(), outStr)
	}
}
