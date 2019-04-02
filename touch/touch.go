package touch

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"os"
	"time"
)

const (
	YY = 2
	Y  = 4
	m  = 2
	DD = 2
	HH = 2
	MM = 2
	S  = 2
)

type Touch struct {
	AccessTime    *bool
	NoCreate      *bool
	Date          *string
	NoDereference *bool
	ModifyTime    *bool
	Reference     *bool
	Stamp         *string
	Time          *string
}

func parseTime(s string) (time.Time, error) {

	timeBuf := bytes.NewBuffer(make([]byte, 0, 15))

	y := Y
	switch len(s) {
	case 15, 12, 13, 10, 11, 8:
		if len(s) == 13 || len(s) == 10 {
			y = YY
		}

		if len(s) == 11 || len(s) == 8 {
			y = 0
			timeBuf.WriteString(fmt.Sprintf("%d", time.Now().Year()+1900))
		}

		timeBuf.WriteString(s[0:y])
		timeBuf.WriteByte('-')
		timeBuf.WriteString(s[y : y+m])
		timeBuf.WriteByte('-')
		timeBuf.WriteString(s[y+m : y+m+DD])
		timeBuf.WriteByte('T')
		timeBuf.WriteString(s[y+m+DD : y+m+DD+HH])
		timeBuf.WriteByte(':')
		timeBuf.WriteString(s[y+m+DD+HH : y+m+DD+HH+MM])
		timeBuf.WriteByte(':')

		if len(s) == 15 || len(s) == 13 || len(s) == 11 || len(s) == 8 {
			timeBuf.WriteString(s[y+m+DD+HH+MM+1 : len(s)])
		} else {
			timeBuf.WriteString("00")
		}
		timeBuf.WriteString("Z")
	default:
		return time.Time{}, fmt.Errorf("invalid date format: %s", s)
	}

	return time.Parse(time.RFC3339, timeBuf.String())
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
		Flags(flag.PosixShort).NewString("")

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

func (t *Touch) Touch(name string) error {
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

	if t.Stamp != nil && len(*t.Stamp) > 0 {
		t, err := parseTime(*t.Stamp)
		if err != nil {
			return err
		}

		atime, mtime = t, t
	}

	os.Chtimes(name, atime, mtime)
	return nil
}

func Main(argv []string) {
	t, args := New(argv)
	for _, v := range args {
		err := t.Touch(v)
		if err != nil {
			utils.Die("touch: %s\n", err)
		}
	}
}
