package tr

import (
	"bytes"
	"strings"
	"testing"
)

func testSet(dst, src string, set1, set2 string, test *testing.T) {
	t := Tr{}
	t.Init(set1, set2)
	rs := strings.NewReader(src)
	w := &bytes.Buffer{}

	t.Tr(rs, w)

	if w.String() != dst {
		test.Errorf("tr %s %s fail, (%s)need:%s\n", set1, set2, w.String(), dst)
	}
}

func TestSet(t *testing.T) {
	src := `HELLO WORLD`
	dst := `hello world`

	testSet(dst, src, "A-Z", "a-z", t)

	src = `12345`
	dst = `87654`

	testSet(dst, src, "0-9", "9876543210", t)

	testSet(src, dst, "9876543210", "0-9", t)

	src = `tr came, tr saw, tr conquered.`
	dst = `ge pnzr, ge fnj, ge pbadhrerq.`
	testSet(dst, src, "a-zA-Z", "n-za-mN-ZA-M", t)

	src = `WELCOME TO 
shanghai`
	dst = `WELCOME TO 
SHANGHAI`
	testSet(dst, src, "[a-z]", "[A-Z]", t)

	testSet(dst, src, "[:lower:]", "[:upper:]", t)

	src = `Welcome To shanghai`
	dst = `Welcome	To	shanghai`
	testSet(dst, src, "[:space:]", `\t`, t)

	src = `{WELCOME TO}`
	dst = `(WELCOME TO)`
	testSet(dst, src, "{}", "()", t)
}

func testDelete(dst, src string, set1, set2 string, delete bool, test *testing.T) {
	t := Tr{}
	t.Delete = delete
	t.Init(set1, set2)
	rs := strings.NewReader(src)
	w := &bytes.Buffer{}

	t.Tr(rs, w)

	if w.String() != dst {
		test.Errorf("tr -d %s fail, (%s)need:%s\n", set1, w.String(), dst)
	}
}

func TestDelete(t *testing.T) {
	src := `welcome`
	dst := `elcome`

	testDelete(dst, src, "w", "", true, t)

	src = "my ID is 73535"
	dst = "my ID is "

	testDelete(dst, src, "[:digit:]", "", true, t)
}

func testComplement(dst, src string, set1, set2 string, complement bool, test *testing.T) {
	t := Tr{}
	t.Complement = complement
	t.Init(set1, set2)
	rs := strings.NewReader(src)
	w := &bytes.Buffer{}

	t.Tr(rs, w)

	if w.String() != dst {
		test.Errorf("tr -d %s fail, (%s)need:%s\n", set1, w.String(), dst)
	}
}

func TestComplement(t *testing.T) {
	src := `unix`
	dst := `uaaa`

	testComplement(dst, src, "u", "a", true, t)
}
