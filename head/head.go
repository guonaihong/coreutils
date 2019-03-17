package head

import (
	"bufio"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type Head struct {
	Bytes     *int
	Lines     *int
	Quiet     *bool
	Verbose   *bool
	LineDelim byte
}

func New(argv []string) (*Head, []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	h := Head{}

	nbytes := command.Opt("c, bytes", "print the first NUM bytes of each file;"+
		" with the leading '-', print all but the last NUM bytes of each file").
		Flags(flag.PosixShort).
		NewString("0")

	h.Lines = command.OptOpt(
		flag.Flag{
			Regex: `^\d+$`,
			Short: []string{"n"},
			Long:  []string{"lines"},
			Usage: "print the first NUM lines instead of the first 10;" +
				"with the leading '-', print all but the last" +
				"NUM lines of each file"}).
		Flags(flag.RegexKeyIsValue | flag.PosixShort).
		NewInt(10)

	h.Quiet = command.Opt("q, quiet, silent", "never print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	h.Verbose = command.Opt("v, verbose", "always print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	zeroTerminated := command.Opt("z, zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).
		NewBool(false)

	command.Parse(argv[1:])

	n, err := utils.HeadParseSize(*nbytes)
	if err != nil {
		utils.Die("head:%s\n", err)
	}

	h.Bytes = n.IntPtr()

	h.LineDelim = '\n'
	if *zeroTerminated {
		h.LineDelim = '\000'
	}

	args := command.Args()
	if len(args) == 0 {
		args = append(args, "-")
	}

	return &h, args
}

func (h *Head) PrintBytes(rs io.ReadSeeker, w io.Writer) {
	size := *h.Bytes
	buf := make([]byte, 1024*8)

	if size < 0 {
		if n, err := rs.Seek(0, 2); err != nil {
			utils.Die("head: %s\n", err)
			return
		} else {
			size = int(n) + *h.Bytes
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

func (h *Head) PrintLines(rs io.ReadSeeker, w io.Writer) {
	br := bufio.NewReader(rs)

	lineNo := 10
	if h.Lines != nil {
		lineNo = *h.Lines
	}

	if lineNo < 0 {
		no := 0
		for no = 0; ; no++ {
			l, e := br.ReadBytes(h.LineDelim)
			if e != nil && len(l) == 0 {
				break
			}

		}

		rs.Seek(0, 0)
		lineNo = no + lineNo
	}

	for i := 0; i < lineNo; i++ {
		l, e := br.ReadBytes(h.LineDelim)
		if e != nil && len(l) == 0 {
			break
		}

		w.Write(l)
	}
}

func (h *Head) main(rs io.ReadSeeker, w io.Writer, name string) {
	if h.Verbose != nil && *h.Verbose {
		h.PrintTitle(w, name)
	}

	if h.Bytes != nil && *h.Bytes != 0 {
		h.PrintBytes(rs, w)
		return
	}

	h.PrintLines(rs, w)
}

func (h *Head) PrintTitle(w io.Writer, name string) {
	fmt.Fprintf(w, "==> %s <==\n", name)
}

func Main(argv []string) {

	h, args := New(argv)

	for _, v := range args {
		fd, err := utils.OpenInputFd(v)
		if err != nil {
			utils.Die("head:%s\n", err)
		}

		h.main(fd, os.Stdout, v)
		utils.CloseInputFd(fd)
	}
}
