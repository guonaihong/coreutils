package wc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"os"
)

type wcCmd struct {
	bytesCounts     *bool
	characterCounts *bool
	lines           *bool
	files0From      *string
	maxLineLength   *bool
	words           *bool
}

func (w *wcCmd) main(fd *os.File) {
	br := bufio.NewReader(fd)

	totalBytes, totalLine, maxLineLength := 0, 0, 0
	totalWords := 0

	for {
		l, e := br.ReadBytes('\n')
		if e != nil && len(l) == 0 {
			break
		}

		ls := bytes.Split(l, []byte{' '})
		for _, v := range ls {
			if len(bytes.TrimSpace(v)) == 0 {
				continue
			}
			totalWords++
		}

		totalBytes += len(l)
		totalLine++
		if len(l) > maxLineLength {
			maxLineLength = len(l)
		}
	}

	fmt.Printf("%4d  %4d  %4d  %s\n", totalLine, totalWords, totalBytes, fd.Name())

}

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ContinueOnError)

	wc := wcCmd{}

	wc.bytesCounts = command.Opt("c, bytes", "print the byte counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.characterCounts = command.Opt("m, chars", "print the character counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.lines = command.Opt("l, lines", "print the newline counts").
		Flags(flag.PosixShort).NewBool(false)

	wc.files0From = command.Opt("files0-from", "read input from the files specified by\n"+
		"NUL-terminated names in file F;\n"+
		"If F is - then read names from standard input").
		Flags(flag.PosixShort).NewString("")

	wc.maxLineLength = command.Opt("L, max-line-length", "print the maximum display width").
		Flags(flag.PosixShort).NewBool(false)

	wc.words = command.Opt("w, words", "print the word counts").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	args := command.Args()
	for _, v := range args {
		fd, err := utils.OpenInputFd(v)
		if err != nil {
			utils.Die("wc: %s\n", err)
		}
		wc.main(fd)
		utils.CloseInputFd(fd)
	}
}
