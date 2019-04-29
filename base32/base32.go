package base32

import (
	"encoding/base32"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type Base32 struct {
	Decode        *bool
	IgnoreGarbage *bool
	Wrap          *int
}

func New(argv []string) (*Base32, []string) {
	b := Base32{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	b.Decode = command.Opt("d, decode", "decode data").
		Flags(flag.PosixShort).NewBool(false)

	b.IgnoreGarbage = command.Opt("i, ignore-garbage", "when decoding, ignore non-alphabet characters").
		Flags(flag.PosixShort).NewBool(false)

	b.Wrap = command.Opt("w, wrap", "wrap encoded lines after COLS character (default 76).\n"+
		"Use 0 to disable line wrapping").
		Flags(flag.PosixShort).NewInt(76)

	command.Parse(argv[1:])
	args := command.Args()

	return &b, args
}

func (b *Base32) Base32(rs io.ReadSeeker, w io.Writer) {
	encoder := base32.NewEncoder(base32.StdEncoding, w)
	io.Copy(encoder, rs)
	encoder.Close()
}

func Main(argv []string) {
	c, args := New(argv)
	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, fileName := range args {
		f, err := utils.OpenInputFd(fileName)
		if err != nil {
			utils.Die("base32: %s\n", err)
		}

		c.Base32(f, os.Stdout)
		utils.CloseInputFd(f)
	}
}
