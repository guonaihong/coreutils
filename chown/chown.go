package chown

import (
	"errors"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"os/user"
	"strconv"
	"strings"
)

type Chown struct {
	Changes        *bool
	Quiet          *bool
	Verbose        *bool
	Dereference    *bool
	NoDereference  *bool
	From           *string
	NoPreserveRoot *bool
	PreserveRoot   *bool
	Reference      *string
	Recursive      *bool
	H              *bool
	L              *bool
	P              *bool
}

type User struct {
	Uid int
	Gid int
	W   io.Writer
}

func (u *User) Init(w io.Writer) {
	u.Uid = -1
	u.Gid = -1
	u.W = w
}

func (u *User) writeMsg(msg string) {
	if u == nil {
		return
	}
	if u.W != nil {
		u.W.Write([]byte(msg))
	}
}

func (u *User) writeError(err error) {
	if u == nil {
		return
	}

	if u.W != nil {
		d := []byte{}
		if err != nil {
			d = []byte(err.Error())
		}
		u.W.Write(d)
	}
}

func New(argv []string) (*Chown, []string) {
	c := &Chown{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	c.Changes = command.Opt("c, changes",
		"like verbose but report only when a change is made").
		Flags(flag.PosixShort).NewBool(false)

	//-f, --silent, --quiet
	c.Quiet = command.Opt("f, silent, quiet",
		"suppress most error messages").
		Flags(flag.PosixShort).NewBool(false)

	c.Verbose = command.Opt("v, verbose",
		"output a diagnostic for every file processed").
		Flags(flag.PosixShort).NewBool(false)

	c.Dereference = command.Opt("dereference",
		"affect the referent of each symbolic link (this is\n"+
			"the default), rather than the symbolic link itself").
		Flags(flag.PosixShort).NewBool(false)

	c.NoDereference = command.Opt("h, no-dereference",
		"affect symbolic links instead of any referenced file\n"+
			"(useful only on systems that can change the\n"+
			"ownership of a symlink)").
		Flags(flag.PosixShort).NewBool(false)

	c.From = command.Opt("from",
		"change the owner and/or group of each file only if\n"+
			"its current owner and/or group match those specified\n"+
			"here.  Either may be omitted, in which case a match\n"+
			"is not required for the omitted attribute").
		Flags(flag.PosixShort).NewString("")

	c.NoPreserveRoot = command.Opt("no-preserve-root",
		"do not treat '/' specially (the default)").
		Flags(flag.PosixShort).NewBool(false)

	c.PreserveRoot = command.Opt("preserve-root",
		"fail to operate recursively on '/'").
		Flags(flag.PosixShort).NewBool(false)

	c.Reference = command.Opt("reference", "use RFILE's owner and group rather than\n"+
		"specifying OWNER:GROUP values").
		Flags(flag.PosixShort).NewString("")

	c.Recursive = command.Opt("R, recursive", "operate on files and directories recursively").
		Flags(flag.PosixShort).NewBool(false)

	c.H = command.Opt("H", "if a command line argument is a symbolic link\n"+
		"to a directory, traverse it").
		Flags(flag.PosixShort).NewBool(false)

	c.L = command.Opt("L", "traverse every symbolic link to a directory\n"+
		"encountered").
		Flags(flag.PosixShort).NewBool(false)

	c.P = command.Opt("P", "do not traverse any symbolic links (default)").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	return c, command.Args()
}

type chownUser user.User

func (u *chownUser) String() string {
	if u == nil {
		return ""
	}

	return u.Username
}

type chownGroup user.Group

func (g *chownGroup) String() string {
	if g == nil {
		return ""
	}
	return g.Name
}

type userGroup struct {
	group *chownGroup
	user  *chownUser
	name  string
}

func getUserGroupFromName(name string) (*userGroup, error) {

	var names []string

	if name == ":" || name == "." {
		return &userGroup{name: name}, nil
	}

	names = strings.Split(name, ":")

	if strings.Index(name, ".") != -1 {
		names = strings.Split(name, ".")
	}

	var u *user.User
	var err error

	if len(names[0]) > 0 {
		u, err = user.Lookup(names[0])
		if err != nil {
			if _, ok := err.(user.UnknownUserError); ok {
				err = fmt.Errorf("chown: invalid user: '%s'", name)
			}
			return nil, err
		}
	}

	var g *user.Group
	if len(names) > 1 && len(names[1]) > 0 {
		g, err = user.LookupGroup(names[1])
		if err != nil {
			if _, ok := err.(user.UnknownGroupError); ok {
				err = fmt.Errorf("chown: invalid group: '%s'", name)
			}
			return nil, err
		}
	}

	return &userGroup{user: (*chownUser)(u), group: (*chownGroup)(g), name: name}, nil
}

func checkArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("chown: missing operand")
	}

	if len(args) == 1 {
		return fmt.Errorf("chown: missing operand after '%s'", args[0])
	}

	return nil
}

func formatError(fileName string, err error) error {
	needError := true
	switch e := err.(type) {
	case *os.PathError:
		err = e.Err
		if os.IsNotExist(err) {
			// The error in golang is 'no such file or directory'
			// The error in gnu chown is 'No such file or directory'
			return fmt.Errorf("chown: cannot access '%s': No such file or directory",
				fileName)
		}

		if os.IsPermission(err) {
			err = errors.New("Operation not permitted")
		}
	case *os.LinkError:
		err = e.Err
	case *os.SyscallError:
		err = e.Err
	default:
		needError = false
	}

	if needError {
		return fmt.Errorf("chown: changing ownership of '%s': %s", fileName, err)
	}

	return err
}

func (c *Chown) IsNoDereference() bool {
	return c.NoDereference != nil && *c.NoDereference
}

func (c *Chown) IsFrom() bool {
	return c.From != nil && len(*c.From) > 0
}

func (c *Chown) IsChanges() bool {
	return c.Changes != nil && *c.Changes
}

func (c *Chown) IsVerbose() bool {
	return c.Verbose != nil && *c.Verbose
}

func (c *Chown) IsPreserveRoot() bool {
	return c.PreserveRoot != nil && *c.PreserveRoot
}

func (c *Chown) IsReference() bool {
	return c.Reference != nil && len(*c.Reference) > 0
}

func noChanages(gu *userGroup, fileUser *user.User, fileGroup *user.Group) bool {
	return (gu.group == nil || gu.group.Gid == fileUser.Gid) &&
		(gu.user == nil || gu.user.Uid == fileUser.Uid)
}

func findGroup(name string) bool {
	return !(strings.Index(name, ":") == -1 && strings.Index(name, ".") == -1)
}

func genToName(gu *userGroup) (to string) {
	to = fmt.Sprintf("%s:%s", gu.user, gu.group)
	if !findGroup(gu.name) || len(gu.name) > 0 && gu.name[0] == ':' {
		to = gu.name
	}
	return
}

func (c *Chown) printVerbse(
	fileName string,
	canChanges bool,
	st *unix.Stat_t,
	gu *userGroup,
	u *User) (err error, msg string) {

	if gu.group == nil && gu.user == nil {
		if c.IsVerbose() {
			msg = fmt.Sprintf("ownership of '%s' retained\n", fileName)
		}
		return
	}

	fileUser, err := user.LookupId(fmt.Sprintf("%d", st.Uid))
	if err != nil {
		return err, ""
	}

	fileGroup, err := user.LookupGroupId(fmt.Sprintf("%d", st.Gid))
	if err != nil {
		return err, ""
	}

	if !canChanges {
		return nil, fmt.Sprintf("ownership of '%s' retained as %s\n",
			fileName, genToName(gu))
	}

	if c.IsChanges() {
		if !noChanages(gu, fileUser, fileGroup) {
			goto next
		}
		return nil, ""
	}

	if noChanages(gu, fileUser, fileGroup) {

		return nil, fmt.Sprintf("ownership of '%s' retained as %s\n",
			fileName, genToName(gu))
	}

next:

	from := fmt.Sprintf("%s:%s", fileUser.Username, fileGroup.Name)
	if !findGroup(gu.name) {
		from = fileUser.Username
	}

	u.writeMsg(fmt.Sprintf("changed ownership of '%s' from %s to %s\n",
		fileName, from,
		genToName(gu)))

	return
}

// user: or user.
func isEndSplitter(name string) bool {
	last := len(name) - 1
	return len(name) > 0 && (name[last] == ':' || name[last] == '.')
}

func parseUidGid(userGroup *userGroup) (uid, gid int) {
	uid, gid = -1, -1
	if userGroup.user != nil {
		uid, _ = strconv.Atoi(userGroup.user.Uid)
	}

	if userGroup.group != nil {
		gid, _ = strconv.Atoi(userGroup.group.Gid)
	}

	return
}

func (c *Chown) getUidGidFromName() (uid, gid int, err error) {

	uid, gid = -1, -1

	userGroup, err := getUserGroupFromName(*c.From)
	if err != nil {
		return
	}

	uid, gid = parseUidGid(userGroup)
	return
}

func (c *Chown) genUidGidFromFile(fileName string) (uidGid string, err error) {
	var st unix.Stat_t

	stat := unix.Stat
	if c.IsNoDereference() {
		stat = unix.Lstat
	}

	err = stat(fileName, &st)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d:%d", st.Uid, st.Gid), nil
}

func (c *Chown) Chown(name string, fileName string, u *User) (err error) {

	defer func() {
		if c.Quiet != nil && *c.Quiet {
			return
		}
		//The default err is output to stdout, or io.Writer (for testing convenience)
		u.writeError(err)
	}()

	if c.IsPreserveRoot() && fileName == "/" {
		return fmt.Errorf("chown: it is dangerous to operate recursively on '/'\n" +
			"chown: use --no-preserve-root to override this failsafe\n")
	}

	userGroup, err := getUserGroupFromName(name)
	if err != nil {
		return err
	}

	uid, gid := parseUidGid(userGroup)

	if gid == -1 && uid != -1 && isEndSplitter(name) {
		g, err := user.LookupGroupId(userGroup.user.Uid)
		if err != nil {
			return err
		}

		userGroup.group = (*chownGroup)(g)
		gid = uid
	}

	// Convenient to write test procedures
	u.Uid, u.Gid = uid, gid

	var st unix.Stat_t

	stat := unix.Stat
	if c.IsNoDereference() {
		stat = unix.Lstat
	}

	if c.IsFrom() || c.IsChanges() || c.IsVerbose() {
		err := stat(fileName, &st)
		if err != nil {
			// todo format error
			return err
		}

	}

	canChanges := true
	if c.IsFrom() {
		uid2, gid2, err := c.getUidGidFromName()
		if err != nil {
			//todo format error
			return err
		}

		if !(uid2 != -1 && uint32(uid2) == st.Uid || gid2 != -1 && uint32(gid2) == st.Gid) {
			canChanges = false
		}
	}

	if c.IsChanges() || c.IsVerbose() {
		err, msg := c.printVerbse(fileName, canChanges, &st, userGroup, u)
		if err != nil {
			return formatError(fileName, err)
		}
		if len(msg) > 0 {
			u.writeMsg(msg)
		}
	}

	if !canChanges {
		return nil
	}

	chown := os.Chown

	if c.IsNoDereference() {
		chown = os.Lchown
	}

	err = chown(fileName, uid, gid)
	if err != nil {
		return formatError(fileName, err)
	}

	if c.Recursive != nil && *c.Recursive {
		stat := os.Lstat
		if c.L != nil && *c.L {
			stat = os.Stat
		}

		walk := utils.NewWalk(c.H != nil && *c.H, stat)
		walk.Walk(fileName, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return chown(path, uid, gid)
		})
	}

	return nil
}

func Main(argv []string) {

	c, args := New(argv)

	err := checkArgs(args)
	if err != nil {
		utils.Die("%s\n", err)
	}

	errCode := 0
	defer func(errCode *int) {
		os.Exit(*errCode)
	}(&errCode)

	u := User{}
	u.Init(os.Stdout)

	if c.IsReference() {
		uidGid, err := c.genUidGidFromFile(*c.Reference)
		if err != nil {
			fmt.Printf("%s\n", err)
			errCode = 1
		}

		args = append([]string{uidGid}, args...)
	}

	for _, a := range args[1:] {
		err := c.Chown(args[0], a, &u)
		if err != nil {
			errCode = 1
		}
	}
}
