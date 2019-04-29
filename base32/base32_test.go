package base32

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	dst := "MFRGGZDFMZTWQ2LKNNWG23TPOBYXE43UOV3HO6DZPI======"
	w := &bytes.Buffer{}

	b := Base32{}
	decode := true
	b.Decode = &decode

	b.Base32(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base32 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}
