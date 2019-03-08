package head

import (
	"bufio"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type head struct {
	bytes     *int
	lines     *int
	quiet     *bool
	verbose   *bool
	lineDelim byte
}

func newHead(argv []string) (*head, []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	h := head{}

	h.bytes = command.Opt("c, bytes", "print the first NUM bytes of each file;"+
		" with the leading '-', print all but the last NUM bytes of each file").
		Flags(flag.PosixShort).
		NewInt(0)

	h.lines = command.OptOpt(
		flag.Flag{
			Regex: `^\d+$`,
			Short: []string{"l"},
			Long:  []string{"lines"},
			Usage: "print the first NUM lines instead of the first 10;" +
				"with the leading '-', print all but the last" +
				"NUM lines of each file"}).
		Flags(flag.RegexKeyIsValue).
		NewInt(10)

	h.quiet = command.Opt("q, quiet, silent", "never print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	h.verbose = command.Opt("v, verbose", "always print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	zeroTerminated := command.Opt("z, zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).
		NewBool(false)

	if *zeroTerminated {
		h.lineDelim = '\000'
	}

	command.Parse(argv[1:])

	args := command.Args()
	if len(args) == 0 {
		args = append(args, "-")
	}

	return &h, args
}

func (h *head) printBytes(rs io.ReadSeeker, w io.Writer) {
	size := *h.bytes
	buf := make([]byte, 1024*8)

	if size < 0 {
		if n, err := rs.Seek(0, 2); err != nil {
			utils.Die("head: %s\n", err)
			return
		} else {
			size = int(n) + *h.bytes
			rs.Seek(0, 0)
		}

	}

	for size > 0 {
		needRead := len(buf)
		if needRead > size {
			needRead = size
		}
		n, err := rs.Read(buf[:needRead])
		if err == io.EOF {
			break
		}

		w.Write(buf[:n])
		size -= n
	}
}

func (h *head) printLines(rs io.ReadSeeker, w io.Writer) {
	br := bufio.NewReader(rs)

	lineNo := *h.lines
	if lineNo < 0 {
		lineMap := map[int]int{}
		no := 0
		for no = 0; ; no++ {
			l, e := br.ReadBytes(h.lineDelim)
			if e == io.EOF {
				break
			}

			no++
			lineMap[no] = len(l) + lineMap[no-1]
		}

		lineNo = lineMap[no+lineNo]
	}

	for i := 0; i < lineNo; {
		l, e := br.ReadBytes(h.lineDelim)
		if e == io.EOF {
			break
		}

		w.Write(l)
	}
}

func (h *head) main(rs io.ReadSeeker, w io.Writer) {
	if *h.bytes != 0 {
		h.printBytes(rs, w)
		return
	}

	h.printLines(rs, w)
}

func Main(argv []string) {

	h, args := newHead(argv)

	for _, v := range args {
		fd, err := utils.OpenInputFd(v)
		if err != nil {
			utils.Die("head:%s\n", err)
		}

		h.main(fd, os.Stdout)
		utils.CloseInputFd(fd)
	}
}
