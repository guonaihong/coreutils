package basecore

import (
	"bytes"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type Base struct {
	OpenDecode    *bool
	IgnoreGarbage *bool
	Wrap          *int
	rs            io.ReadSeeker
	w             io.Writer
	buffer        bytes.Buffer
	buf           []byte
	encodeFlush   bool
	needChar      func(b byte) bool
	baseName      string
}

func New(argv []string) (*Base, []string) {
	b := Base{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	b.OpenDecode = command.Opt("d, decode", "decode data").
		Flags(flag.PosixShort).NewBool(false)

	b.IgnoreGarbage = command.Opt("i, ignore-garbage",
		"when decoding, ignore non-alphabet characters").
		Flags(flag.PosixShort).NewBool(false)

	b.Wrap = command.Opt("w, wrap", "wrap encoded lines after COLS character (default 76).\n"+
		"Use 0 to disable line wrapping").
		Flags(flag.PosixShort).NewInt(76)

	command.Parse(argv[1:])
	args := command.Args()

	return &b, args
}

func (b *Base) IsWrap() bool {
	return !(b.Wrap == nil || *b.Wrap == 0)
}

func (b *Base) Write(p []byte) (n int, err error) {

	if !b.IsWrap() {
		return b.w.Write(p)
	}

	b.buffer.Write(p)
	for b.buffer.Len() >= len(b.buf) {
		n1, err := b.buffer.Read(b.buf)
		if err != nil {
			fmt.Printf("read fail:%s\n", err)
			break
		}

		b.w.Write(b.buf[:n1])
		b.w.Write([]byte{'\n'})
	}

	if b.encodeFlush && b.buffer.Len() > 0 {
		defer func() {
			b.buffer.Reset()
			b.encodeFlush = false
			b.w.Write([]byte{'\n'})
		}()

		return b.w.Write(b.buffer.Bytes())
	}

	return
}

// todo < 0
func (b *Base) checkArgs() error {
	if b.Wrap != nil && *b.Wrap < 0 {
		return fmt.Errorf(`%s: invalid wrap size: "%d"`, b.baseName, *b.Wrap)
	}
	return nil
}

func (b *Base) Encode(rs io.ReadSeeker, w io.Writer) {
	b.w = w

	var encoder io.WriteCloser
	if b.baseName == "base32" {
		encoder = base32.NewEncoder(base32.StdEncoding, b)
	} else {
		encoder = base64.NewEncoder(base64.StdEncoding, b)
	}

	io.Copy(encoder, rs)

	b.encodeFlush = true
	encoder.Close()
	b.Write([]byte{}) //encodeFlush

	if !b.IsWrap() {
		w.Write([]byte{'\n'})
	}
}

func (b *Base) Read(p []byte) (n int, err error) {
	if b.IgnoreGarbage != nil && *b.IgnoreGarbage {
		if b.buf == nil || len(p) > len(b.buf) {
			b.buf = make([]byte, len(p))
		}

		n1, err := b.rs.Read(b.buf)
		if err != nil {
			return n, err
		}

		for _, v := range b.buf[:n1] {
			if !b.needChar(v) {
				continue
			}

			p[n] = v
			n++
		}

		copy(p, b.buf[:n+1])
		return n, nil
	}

	return b.rs.Read(p)
}

func (b *Base) Decode(rs io.ReadSeeker, w io.Writer) {
	b.rs = rs
	// base32
	// ABCDEFGHIJKLMNOPQRSTUVWXYZ234567=
	b.needChar = func(b byte) bool {
		return b >= 'A' && b <= 'Z' ||
			b >= '2' && b <= '7' ||
			b == '='
	}

	if b.baseName == "base64" {
		// base64
		// abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789=+/
		b.needChar = func(b byte) bool {
			return b >= 'a' && b <= 'z' ||
				b >= 'A' && b <= 'Z' ||
				b >= '0' && b <= '9' ||
				b == '=' ||
				b == '+' ||
				b == '/'
		}
	}

	var decode io.Reader

	if b.baseName == "base32" {
		decode = base32.NewDecoder(base32.StdEncoding, b)
	} else {
		decode = base64.NewDecoder(base64.StdEncoding, b)
	}

	io.Copy(w, decode)

}

func (b *Base) Base(rs io.ReadSeeker, w io.Writer) {
	if b.IsWrap() {
		b.buf = make([]byte, *b.Wrap)
	}

	if b.OpenDecode != nil && *b.OpenDecode {
		b.Decode(rs, w)
		return
	}

	b.Encode(rs, w)
}

func Main(argv []string, baseName string) {
	b, args := New(argv)
	if len(args) == 0 {
		args = append(args, "-")
	}

	b.baseName = baseName
	err := b.checkArgs()
	if err != nil {
		utils.Die("%s\n", err)
	}

	for _, fileName := range args {
		f, err := utils.OpenFile(fileName)
		if err != nil {
			utils.Die("base32: %s\n", err)
		}

		b.Base(f, os.Stdout)
		f.Close()
	}
}
