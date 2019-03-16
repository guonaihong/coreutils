package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	dst := "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo="
	w := &bytes.Buffer{}

	b := Base64{}
	decode := true
	b.Decode = &decode

	b.Base64(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base64 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}
