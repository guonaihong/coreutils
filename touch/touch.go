package touch

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"golang.org/x/sys/unix"
	"os"
	"regexp"
	"strings"
	"syscall"
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
	AccessTime       *bool
	NoCreate         *bool
	Date             *string
	NoDereference    *bool
	ModifyTime       *bool
	Reference        *string
	Stamp            *string
	Time             *string
	r                *regexp.Regexp
	month2NumReplace *strings.Replacer
	debug            *bool
}

var shortMonthNames = []string{
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
}

var longMonthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

func (t *Touch) openDebug() bool {
	return t.debug != nil && *t.debug
}

func (t *Touch) parseDate(d string) (time.Time, error) {
	if t.r == nil {
		// parse-->1 May 2005 10:22

		p := `(\d+ (?:%s))?\s*(\d{4})?\s*(\d{2}:\d{2})?`

		p = fmt.Sprintf(p, strings.Join(append(shortMonthNames,
			longMonthNames...), "|"))

		t.r = regexp.MustCompile(p)

		month2Num := make([]string, 0,
			len(shortMonthNames)*4)

		if t.month2NumReplace == nil {

			for k := range shortMonthNames {

				n := fmt.Sprintf("%02d", k+1)
				month2Num = append(month2Num, longMonthNames[k], n)
				month2Num = append(month2Num, shortMonthNames[k], n)
			}

			t.month2NumReplace = strings.NewReplacer(month2Num...)
		}

		if t.openDebug() {
			fmt.Printf("regexp-->(%s)\n", p) //debug
		}
	}

	res := t.r.FindStringSubmatch(d)
	now := time.Now()
	year, month, day := now.Date()

	//month-day
	if len(res[1]) == 0 {
		res[1] = fmt.Sprintf("%02d-%02d", month, day)
	} else {
		res[1] = t.month2NumReplace.Replace(res[1])
		rs := strings.Split(res[1], " ")
		//swap day month to month-day
		res[1] = fmt.Sprintf("%02s-%02s", rs[1], rs[0])
	}

	// year
	if len(res[2]) == 0 {
		res[2] = fmt.Sprintf("%d", year)
	}

	// hour:minute
	if len(res[3]) == 0 {
		res[3] = fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	}

	todoParse := fmt.Sprintf("%s-%sT%s:00Z", res[2], res[1], res[3])
	if t.openDebug() {
		fmt.Printf("todoParse = %s\n", todoParse)
	}

	return time.ParseInLocation("2006-01-02T15:04:05Z",
		todoParse,
		time.Local)
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

func (t *Touch) IsNoDereference() bool {
	return t.NoDereference != nil && *t.NoDereference
}

func (t *Touch) IsReference() bool {
	return t.Reference != nil && len(*t.Reference) > 0
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
		Flags(flag.PosixShort).NewString("")

	touch.Stamp = command.Opt("t", "use [[CC]YY]MMDDhhmm[.ss] instead of current time").
		Flags(flag.PosixShort).NewString("")

	touch.Time = command.Opt("time", "change the specified time:\n"+
		"WORD is access, atime, or use: equivalent to -a\n"+
		"WORD is modify or mtime: equivalent to -m").
		NewString("")

	touch.debug = command.Opt("debug", "debug mode").NewBool(false)

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

func statTimes(name string) (atime, mtime, ctime time.Time, err error) {
	fi, err := os.Stat(name)
	if err != nil {
		return
	}

	mtime = fi.ModTime()
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		now := time.Now()
		return now, now, now, nil
	}

	atime = time.Unix(int64(stat.Atim.Sec), int64(stat.Atim.Nsec))
	ctime = time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec))
	return
}

func (t *Touch) parseTimeOpt() {
	if t.Time == nil || len(*t.Time) == 0 {
		return
	}

	switch *t.Time {
	case "atime", "access", "use":
		t.AccessTime = utils.Bool(true)
	case "mtime", "modify":
		t.ModifyTime = utils.Bool(true)
	}
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

	if t.Date != nil && len(*t.Date) > 0 {
		t, err := t.parseDate(*t.Date)
		if err != nil {
			return err
		}
		atime, mtime = t, t
	}

	// -time
	t.parseTimeOpt()

	var st unix.Stat_t

	if t.IsReference() {
		name = *t.Reference
	}

	var a, m time.Time

	if t.IsNoDereference() {
		err := unix.Lstat(name, &st)
		if err != nil {
			return err
		}

		a = time.Unix(st.Atim.Sec, st.Atim.Nsec)
		m = time.Unix(st.Atim.Sec, st.Atim.Nsec)
	} else {
		a, m, _, _ = statTimes(name)
	}

	// get atime mtime
	if t.AccessTime != nil && *t.AccessTime || t.IsReference() {
		atime = a
	}

	if t.ModifyTime != nil && *t.ModifyTime || t.IsReference() {
		mtime = m
	}

	if t.IsNoDereference() {

		return unix.Lutimes(name,
			[]unix.Timeval{
				Time2Timeval(atime),
				Time2Timeval(mtime)})
	}

	os.Chtimes(name, atime, mtime)
	return nil
}

func Time2Timeval(t time.Time) unix.Timeval {
	return unix.Timeval{Sec: int64(t.Second()),
		Usec: int64(t.Nanosecond() / 1000)}
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
