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
	Decode        *bool
	IgnoreGarbage *bool
	Wrap          *int
	w             io.Writer
	buffer        bytes.Buffer
	buf           []byte
	flush         bool
	baseName      string
}

func New(argv []string) (*Base, []string) {
	b := Base{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	b.Decode = command.Opt("d, decode", "decode data").
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

	if b.flush {
		defer b.buffer.Reset()
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

func (b *Base) Base(rs io.ReadSeeker, w io.Writer) {
	if b.IsWrap() {
		b.buf = make([]byte, *b.Wrap)
	}

	b.w = w
	var encoder io.WriteCloser
	if b.baseName == "base32" {
		encoder = base32.NewEncoder(base32.StdEncoding, b)
	} else {
		encoder = base64.NewEncoder(base64.StdEncoding, b)
	}
	io.Copy(encoder, rs)

	b.flush = true
	b.Write([]byte{}) //flush

	encoder.Close()
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
		f, err := utils.OpenInputFd(fileName)
		if err != nil {
			utils.Die("base32: %s\n", err)
		}

		b.Base(f, os.Stdout)
		utils.CloseInputFd(f)
	}
}
