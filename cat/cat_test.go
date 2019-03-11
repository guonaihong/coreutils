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
