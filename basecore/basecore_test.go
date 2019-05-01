package basecore

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"strings"
	"testing"
)

func testEncode64(src, dst string, wrap *int, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base64"
	b.Wrap = wrap

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		var wrapMessage string
		if wrap == nil {
			wrapMessage = "wrap is nil"
		} else {
			wrapMessage = fmt.Sprintf("wrap = %d", *wrap)
		}

		t.Errorf("base64 fail(%s), need(%s), %s\n", w.String(), dst, wrapMessage)
	}
}

func TestEncode64(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	dst := "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=\n"
	testEncode64(src, dst, nil, t)
	testEncode64(src, dst, utils.Int(0), t)
	testEncode64(src, dst, utils.Int(76), t)

	src = "hello china"
	dst = "aGVsbG8gY2hpbmE=\n"
	testEncode64(src, dst, nil, t)
	testEncode64(src, dst, utils.Int(0), t)
	testEncode64(src, dst, utils.Int(76), t)

	src = "hello china\n"
	dst = "aGVsbG8gY2hpbmEK\n"
	testEncode64(src, dst, nil, t)
	testEncode64(src, dst, utils.Int(0), t)
	testEncode64(src, dst, utils.Int(76), t)
}

func testEncode32(src, dst string, wrap *int, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base32"

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		var wrapMessage string
		if wrap == nil {
			wrapMessage = "wrap is nil"
		} else {
			wrapMessage = fmt.Sprintf("wrap = %d", *wrap)
		}

		t.Errorf("base32 fail(%s), need(%s), src(%s) wrap(%s)\n", w.String(), dst, src, wrapMessage)
	}
}

func TestEncode32(t *testing.T) {
	src := "abcdefghijklmnopqrstuvwxyz"
	dst := "MFRGGZDFMZTWQ2LKNNWG23TPOBYXE43UOV3HO6DZPI======\n"
	testEncode32(src, dst, nil, t)
	testEncode32(src, dst, utils.Int(0), t)
	testEncode32(src, dst, utils.Int(76), t)

	src = "hello china\n"
	dst = "NBSWY3DPEBRWQ2LOMEFA====\n"
	testEncode32(src, dst, nil, t)
	testEncode32(src, dst, utils.Int(0), t)
	testEncode32(src, dst, utils.Int(76), t)

	src = "hello china"
	dst = "NBSWY3DPEBRWQ2LOME======\n"
	testEncode32(src, dst, nil, t)
	testEncode32(src, dst, utils.Int(0), t)
	testEncode32(src, dst, utils.Int(76), t)

	src = "hello\n"
	dst = "NBSWY3DPBI======\n"
	testEncode32(src, dst, nil, t)
	testEncode32(src, dst, utils.Int(0), t)
	testEncode32(src, dst, utils.Int(76), t)
}

func testDecode32(src, dst string, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base32"
	b.OpenDecode = utils.Bool(true)

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base32 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}

func TestDecode32(t *testing.T) {
	src := "MFRGGZDFMZTWQ2LKNNWG23TPOBYXE43UOV3HO6DZPI======\n"
	dst := "abcdefghijklmnopqrstuvwxyz"
	testDecode32(src, dst, t)

	src = "NBSWY3DPEBRWQ2LOMEFA====\n"
	dst = "hello china\n"
	testDecode32(src, dst, t)

	src = "NBSWY3DPEBRWQ2LOME======\n"
	dst = "hello china"
	testDecode32(src, dst, t)
}

func testDecode64(src, dst string, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base64"
	b.OpenDecode = utils.Bool(true)

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base64 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}

func TestDecode64(t *testing.T) {
	src := "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=\n"
	dst := "abcdefghijklmnopqrstuvwxyz"
	testDecode64(src, dst, t)

	dst = "hello china"
	src = "aGVsbG8gY2hpbmE=\n"
	testDecode64(src, dst, t)

	dst = "hello china\n"
	src = "aGVsbG8gY2hpbmEK\n"
	testDecode64(src, dst, t)
}

func testDecode32IgnoreGarbage(src, dst string, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base32"
	b.OpenDecode = utils.Bool(true)
	b.IgnoreGarbage = utils.Bool(true)

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base32 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}

func TestDecode32IgnoreGarbage(t *testing.T) {
	src := "MFRGGZDFMZTWQ2LKNNWG23TPOBYXE43UOV3HO6DZPI======\n"
	dst := "abcdefghijklmnopqrstuvwxyz"
	testDecode32IgnoreGarbage(src, dst, t)

	src = "NBSWY3DPEBRWQ2LOMEFA====\n"
	dst = "hello china\n"
	testDecode32IgnoreGarbage(src, dst, t)

	src = "NBSWY3DPEBRWQ2LOME======\n"
	dst = "hello china"
	testDecode32IgnoreGarbage(src, dst, t)
}

func testDecode64IgnoreGarbage(src, dst string, t *testing.T) {
	w := &bytes.Buffer{}

	b := Base{}
	b.baseName = "base64"
	b.OpenDecode = utils.Bool(true)
	b.IgnoreGarbage = utils.Bool(true)

	b.Base(strings.NewReader(src), w)

	if w.String() != dst {
		t.Errorf("base64 -d fail(%s), need(%s)\n", w.String(), dst)
	}
}

func TestDecode64IgnoreGarbage(t *testing.T) {
	src := "YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=\n"
	dst := "abcdefghijklmnopqrstuvwxyz"
	testDecode64IgnoreGarbage(src, dst, t)

	dst = "hello china"
	src = "aGVsbG8gY2hpbmE=\n"
	testDecode64IgnoreGarbage(src, dst, t)

	dst = "hello china\n"
	src = "aGVsbG8gY2hpbmEK\n"
	testDecode64IgnoreGarbage(src, dst, t)
}
