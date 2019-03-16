package main

import (
	"encoding/base64"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type Base64 struct {
	Decode        *bool
	IgnoreGarbage *bool
	Wrap          *int
}

func New(argv []string) (*Base64, []string) {
	b := Base64{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	b.Decode = command.Opt("d, decode", "decode data").
		Flags(flag.PosixShort).NewBool(false)

	b.IgnoreGarbage = command.Opt("i, ignore-garbage", "when decoding, ignore non-alphabet characters").
		Flags(flag.PosixShort).NewBool(false)

	b.Wrap = command.Opt("w, wrap", "wrap encoded lines after COLS character (default 76).\n"+
		"Use 0 to disable line wrapping").
		Flags(flag.PosixShort).NewInt(0)

	command.Parse(argv[1:])
	args := command.Args()

	return &b, args
}

func (b *Base64) Base64(rs io.ReadSeeker, w io.Writer) {
	encoder := base64.NewEncoder(base64.StdEncoding, w)
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
			utils.Die("base64: %s\n", err)
		}

		c.Base64(f, os.Stdout)
		utils.CloseInputFd(f)
	}
}
