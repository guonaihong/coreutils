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

	c.Number = SetBool(true)

	in := strings.NewReader(testNumber)
	out := &bytes.Buffer{}

	c.main(in, out)

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

	c.main(in, out)

	outStr := out.String()
	outStr = strings.Replace(outStr, "^I", "", -1)
	outStr = strings.Replace(outStr, "\n", "", -1)

	if len(outStr) != 0 {
		t.Fatalf("cat -T fail (%s)", outStr)

	}
}
