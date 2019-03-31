package touch

import (
	"github.com/guonaihong/flag"
	"os"
	"time"
)

type Touch struct {
	AccessTime    *bool
	NoCreate      *bool
	Date          *string
	NoDereference *bool
	ModifyTime    *bool
	Reference     *bool
	Stamp         *bool
	Time          *string
}

func New(argv []string) (*Touch, []string) {
	touch := Touch{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	touch.AccessTime = command.Opt("a", "change only the access time").
		Flags(flag.PosixShort).NewBool(false)
	touch.NoCreate = command.Opt("c, no-create", "do not create any files").
		Flags(flag.PosixShort).NewBool(false)
	touch.Date = command.Opt("d, date",
		"parse STRING and use it instead of current time").
		Flags(flag.PosixShort).NewString("")

	touch.NoDereference = command.Opt("h, no-dereference",
		"affect each symbolic link instead of any referenced\n"+
			"file (useful only on systems that can change the\n"+
			"timestamps of a symlink)").
		Flags(flag.PosixShort).NewBool(false)
	touch.ModifyTime = command.Opt("m", "change only the modification time").
		Flags(flag.PosixShort).NewBool(false)

	touch.Reference = command.Opt("r, referenced",
		"use this file's times instead of current time").
		Flags(flag.PosixShort).NewBool(false)
	touch.Stamp = command.Opt("t", "use [[CC]YY]MMDDhhmm[.ss] instead of current time").
		Flags(flag.PosixShort).NewBool(false)
	touch.Time = command.Opt("time", "change the specified time:\n"+
		"WORD is access, atime, or use: equivalent to -a\n"+
		"WORD is modify or mtime: equivalent to -m").
		Flags(flag.PosixShort).NewString("")

	command.Parse(argv[1:])

	args := command.Args()

	return &touch, args
}

func isNotExist(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return true
	}
	return false
}

func (t *Touch) Touch(name string) {
	if isNotExist(name) && name != "-" {

		if t.NoCreate == nil || t.NoCreate != nil && !*t.NoCreate {
			fd, err := os.Create(name)
			if err != nil {
				fd.Close()
			}
		}
	}

	now := time.Now()
	atime := now
	mtime := now
	os.Chtimes(name, atime, mtime)
}

func Main(argv []string) {
	t, args := New(argv)
	for _, v := range args {
		t.Touch(v)
	}
}
