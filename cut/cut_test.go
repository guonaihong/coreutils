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
		t.Errorf("cut -b fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

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

	dst = `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`

	testBytes(src, dst, "1-", t)

	dst = `And
Aru
Ass
Bih
Chh`

	testBytes(src, dst, "-3", t)
}

func testCharacters(src, dst, cstr string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Characters = &cstr
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -c fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

func TestCharacters(t *testing.T) {
	src := `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`
	dst := `nr 
rah
sm
ir
hti`

	testCharacters(src, dst, "2,5,7", t)

	dst = `Andhra 
Arunach
Assam
Bihar
Chhatti`

	testCharacters(src, dst, "1-7", t)

	dst = `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`

	testCharacters(src, dst, "1-", t)

	dst = `Andhr
Aruna
Assam
Bihar
Chhat`
	testCharacters(src, dst, "-5", t)

}

func testFieldsDelimiter(src, dst string, delimiter string, fields string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Delimiter = &delimiter
	c.Fields = &fields
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -f -d fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

func TestFieldsDelimiter(t *testing.T) {
	src := `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`
	dst := `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`

	testFieldsDelimiter(src, dst, "\t", "1", t)

	dst = `Andhra
Arunachal
Assam
Bihar
Chhattisgarh`
	testFieldsDelimiter(src, dst, " ", "1", t)
	//=======================================

	dst = `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`
	testFieldsDelimiter(src, dst, " ", "1-4", t)
}

func testComplement(complement bool, src, dst string, delimiter string, fields string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Complement = &complement
	c.Delimiter = &delimiter
	c.Fields = &fields
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -complement fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

func testComplement2(complement bool, src, dst string, characters string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Complement = &complement
	c.Characters = &characters
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -complement fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}
func TestComplement(t *testing.T) {
	src := `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`

	dst := `Pradesh
Pradesh
Assam
Bihar
Chhattisgarh`
	testComplement(true, src, dst, " ", "1", t)
	//=========

	dst = `Andha Pradesh
Arunchal Pradesh
Assa
Biha
Chhatisgarh`
	testComplement2(true, src, dst, "5", t)
}

func testOutputDelimiter(src, dst string, delimiter string, fields string, outputDelimiter string, t *testing.T) {
	w := &bytes.Buffer{}

	c := Cut{}
	c.Delimiter = &delimiter
	c.Fields = &fields
	c.OutputDelimiter = &outputDelimiter
	c.LineDelim = '\n'
	c.Init()
	c.Cut(strings.NewReader(src), w)
	if w.String() != dst {
		t.Errorf("cut -output-delimiter fail(%s)(%x)(%x), need(%s)", w.String(), w.String(), w.Len(), dst)
	}
}

func TestOutputDelimiter(t *testing.T) {
	src := `Andhra Pradesh
Arunachal Pradesh
Assam
Bihar
Chhattisgarh`

	dst := `Andhra%Pradesh
Arunachal%Pradesh
Assam
Bihar
Chhattisgarh`

	testOutputDelimiter(src, dst, " ", "1,2", "%", t)
}

func TestCmdOption(t *testing.T) {
	c, _ := New([]string{"cut", "-f3"})
	if c.Fields == nil || *c.Fields != "3" {
		t.Errorf("cut -f options fail\n")
	}

	c, _ = New([]string{"cut", "-c1-5"})
	if c.Characters == nil || *c.Characters != "1-5" {
		t.Errorf("cut -c options fail\n")
	}

	c, _ = New([]string{"cut", "-c-5"})
	if c.Characters == nil || *c.Characters != "-5" {
		t.Errorf("cut -c options fail\n")
	}
}
