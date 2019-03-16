package cut

import (
	"bytes"
	"strings"
	"testing"
)

func testBytes(src, dst, b string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Bytes = &b
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -c fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

//todo
func TestBytes(t *testing.T) {
	src := "12345678910"
	dst := "1234567"

	testBytes(src, dst, "1-5,3,6-7", t)
	//====================

	src = `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`
	dst = `And
Aru
Ass
Bih
Chh`
	testBytes(src, dst, "1,2,3", t)
	//========================

	dst = `Andra 
Aruach
Assm
Bihr
Chhtti`

	testBytes(src, dst, "1-3,5-7", t)
}

//todo
func TestCharacters(t *testing.T) {
}

//todo
func TestDelimiter(t *testing.T) {
}

//todo
func TestFields(t *testing.T) {
}

//todo
func TestComplement(t *testing.T) {
}

//todo
func TestOutputDelimiter(t *testing.T) {
}

//todo
func TestCmdOption(t *testing.T) {
}
