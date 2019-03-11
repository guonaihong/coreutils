package tail

import (
	"fmt"
)

type Tail struct {
	Bytes             *int
	Follow            *string
	Lines             *int
	Quiet             *bool
	Verbose           *bool
	SleepInterval     *float64
	Retry             *bool
	MaxUnchangedStats *string
	LineDelim         byte
}

func New(argv []string) (*Tail, []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	t := Tail{}
	nbytes := command.Opt("c, bytes", "output the last NUM bytes; or use -c +NUM to\n"+
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
		NewInt(10)

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

	t.Retry = command.Opt("retry", "keep trying to open a file if it is inaccessible")

	t.SleepInterval = command.Opt("s, sleep-interval", "with -f, sleep for approximately N seconds\n"+
		"(default 1.0) between iterations;\n"+
		"with inotify and --pid=P, check process P at\n"+
		"least once every N seconds").
		Flags(flag.PosixShort).
		NewBool(false)

	t.Verbose = command.Opt("v, verbose", "always print headers giving file names").
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

	t.Bytes = n.IntPtr()

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

func (t *Tail) PrintBytes(rs io.ReadSeeker, w io.Writer) {
}

func (t *Tail) PrintLines(rs io.ReadSeeker, w io.Writer) {
}

func Main(argv []string) {
	t, args := New(argv)

	for _, v := range args {
		fd, err := utils.OpenInputFd(v)
		if err != nil {
			utils.Die("head:%s\n", err)
		}

		h.main(fd, os.Stdout, v)
		utils.CloseInputFd(fd)
	}
}
