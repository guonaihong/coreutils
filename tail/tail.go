package tail

import (
	"bufio"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
)

type Tail struct {
	Bytes             *string
	Follow            *string
	Lines             *string
	Quiet             *bool
	Verbose           *bool
	SleepInterval     *float64
	Retry             *bool
	MaxUnchangedStats *string
	LineDelim         byte
	Pid               *string
}

func New(argv []string) (*Tail, []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	t := Tail{}
	t.Bytes = command.Opt("c, bytes", "output the last NUM bytes; or use -c +NUM to\n"+
		"output starting with byte NUM of each file").
		Flags(flag.PosixShort).
		NewString("0")

	t.Follow = command.Opt("f, follow", "output appended data as the file grows;\n"+
		"an absent option argument means 'descriptor'").
		Flags(flag.PosixShort).
		NewString("")

	_ = command.Opt("F", "same as --follow=name --retry").
		Flags(flag.PosixShort).
		NewBool(false)

	t.Lines = command.OptOpt(
		flag.Flag{
			Regex: `^\d+$`,
			Short: []string{"n"},
			Long:  []string{"lines"},
			Usage: "output appended data as the file grows;\n" +
				"an absent option argument means 'descriptor'"}).
		Flags(flag.RegexKeyIsValue).
		NewString("10")

	t.MaxUnchangedStats = command.Opt("max-unchanged-stats", "with --follow=name, reopen a FILE which has not\n"+
		"changed size after N (default 5) iterations\n"+
		"to see if it has been unlinked or renamed\n"+
		"(this is the usual case of rotated log files);\n"+
		"with inotify, this option is rarely useful").
		NewString("")

	t.Pid = command.Opt("pid", "with -f, terminate after process ID, PID dies").
		Flags(flag.PosixShort).
		NewString("")

	t.Quiet = command.Opt("q, quiet, silent", "never print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	t.Retry = command.Opt("retry", "keep trying to open a file if it is inaccessible").
		NewBool(false)

	t.SleepInterval = command.Opt("s, sleep-interval", "with -f, sleep for approximately N seconds\n"+
		"(default 1.0) between iterations;\n"+
		"with inotify and --pid=P, check process P at\n"+
		"least once every N seconds").
		Flags(flag.PosixShort).
		NewFloat64(0.0)

	t.Verbose = command.Opt("v, verbose", "always print headers giving file names").
		Flags(flag.PosixShort).
		NewBool(false)

	zeroTerminated := command.Opt("z, zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).
		NewBool(false)

	command.Parse(argv[1:])

	t.LineDelim = '\n'
	if *zeroTerminated {
		t.LineDelim = '\000'
	}

	args := command.Args()
	if len(args) == 0 {
		args = append(args, "-")
	}

	return &t, args
}

func (t *Tail) PrintBytes(rs io.ReadSeeker, w io.Writer) error {
	nBytes := 0

	n, err := utils.HeadParseSize(*t.Bytes)
	if err != nil {
		return err
	}

	nBytes = int(n)
	bytes := *t.Bytes

	readTail := false
	if bytes[0] != '+' || nBytes < 0 {
		if nBytes > 0 {
			nBytes = -nBytes
		}
		readTail = true
	}

	if readTail {
		offset, err := rs.Seek(0, 2)
		if err != nil {
			return err
		}

		rs.Seek(offset+int64(nBytes), 0)

		_, err = io.Copy(w, rs)
		return err
	}

	if nBytes > 0 {
		nBytes--
	}

	if _, err = rs.Seek(int64(nBytes), 0); err != nil {
		return err
	}

	_, err = io.Copy(w, rs)

	return err
}

func (t *Tail) printTailLines(rs io.ReadSeeker, w io.Writer, n int) error {
	no := 0
	totalMap := map[int]int{}
	br := bufio.NewReader(rs)

	for no = 0; ; no++ {

		l, e := br.ReadBytes(t.LineDelim)

		if e != nil && len(l) == 0 {
			no--
			break
		}

		if no == 0 {
			totalMap[no] = len(l)
			continue
		}

		totalMap[no] = totalMap[no-1] + len(l)
	}

	rs.Seek(int64(totalMap[no+n]), 0)

	for {

		l, e := br.ReadBytes(t.LineDelim)
		if e != nil && len(l) == 0 {
			break
		}

		w.Write(l)
	}

	return nil
}

func (t *Tail) PrintLines(rs io.ReadSeeker, w io.Writer) error {
	nLines := 0

	n0, err := utils.HeadParseSize(*t.Lines)
	if err != nil {
		return err
	}

	nLines = int(n0)

	lines := *t.Lines
	readLast := true

	//+10
	if lines[0] != '+' && lines[0] != '-' || nLines < 0 {
		readLast = false
		if nLines > 0 {
			nLines = -nLines
		}
	}

	if readLast {

		br := bufio.NewReader(rs)
		for i := 1; ; i++ {

			l, e := br.ReadBytes(t.LineDelim)
			if e != nil && len(l) == 0 {
				break
			}

			if i < nLines {
				continue
			}

			w.Write(l)
		}
		return nil
	}

	return t.printTailLines(rs, w, nLines)
}

func (t *Tail) main(rs io.ReadSeeker, w io.Writer, name string) {
}

func Main(argv []string) {
	t, args := New(argv)

	for _, v := range args {
		fd, err := utils.OpenInputFd(v)
		if err != nil {
			utils.Die("head:%s\n", err)
		}

		t.main(fd, os.Stdout, v)
		utils.CloseInputFd(fd)
	}
}
