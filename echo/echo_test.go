package echo

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"testing"
)

func testOptione(in []string, dst string, t *testing.T) {
	echo := Echo{}
	echo.Enable = utils.Bool(true)
	echo.NewLine = utils.Bool(true)
	w := &bytes.Buffer{}

	echo.Echo(in, w)

	if dst != w.String() {
		t.Fatalf("echo -e fail(%s,%x l:%d) need(%s, l:%d)\n",
			w.String(), w.String(), w.Len(), dst, len(dst))
	}
}

// echo -e
func TestOptione(t *testing.T) {

	in := []string{
		`\n\r`,
	}

	testOptione(in, "\n\r", t)
	in = []string{
		"hello",
		"world",
	}

	testOptione(in, "hello world", t)

	in = []string{
		`\a\b\f\n\r\t\v\000\xff`,
	}

	//todo
	/*
		\e     escape
	*/
	testOptione(in, "\a\b\f\n\r\t\v\000\xff", t)

	in = []string{
		`\e`,
	}
	testOptione(in, "\x1b", t)
}

func testOptionn(in []string, dst string, t *testing.T) {
	echo := Echo{}
	echo.NewLine = utils.Bool(true)
	w := &bytes.Buffer{}

	echo.Echo(in, w)

	if dst != w.String() {
		t.Fatalf("echo -n fail(%s,%x l:%d) need(%s, l:%d)\n",
			w.String(), w.String(), w.Len(), dst, len(dst))
	}
}

// echo -n
func TestOptionn(t *testing.T) {
	in := []string{
		"hello", "world", "12345",
	}

	testOptionn(in, "hello world 12345", t)
}

func testOptionE(in []string, dst string, t *testing.T) {
	echo := Echo{}
	echo.Disable = utils.Bool(true)
	w := &bytes.Buffer{}

	echo.Echo(in, w)

	if dst != w.String() {
		t.Fatalf("echo -E fail(%s,%x l:%d) need(%s, l:%d)\n",
			w.String(), w.String(), w.Len(), dst, len(dst))
	}
}

// echo -E
func TestOptionE(t *testing.T) {
	in := []string{
		"hello", `\n\n\n\n`,
	}

	testOptionn(in, `hello \n\n\n\n`, t)
}
