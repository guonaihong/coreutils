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
	Changes *bool
	//-f, --silent, --quiet
	Verbose        *bool
	Dereference    *bool
	NoDereference  *bool
	Form           *string
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

func New(argv []string) (*Chown, []string) {
	c := &Chown{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	c.Changes = command.Opt("c, changes",
		"like verbose but report only when a change is made").
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

	c.Form = command.Opt("from",
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

type groupUser struct {
	group *chownGroup
	user  *chownUser
	name  string
}

func getGroupUser(name string) (*groupUser, error) {

	var names []string

	if name == ":" || name == "." {
		return &groupUser{name: name}, nil
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

	return &groupUser{user: (*chownUser)(u), group: (*chownGroup)(g), name: name}, nil
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
	}

	return fmt.Errorf("chown: changing ownership of '%s': %s", fileName, err)
}

func (c *Chown) IsChanges() bool {
	return c.Changes != nil && *c.Changes
}

func (c *Chown) IsVerbose() bool {
	return c.Verbose != nil && *c.Verbose
}

func noChanages(gu *groupUser, fileUser *user.User, fileGroup *user.Group) bool {
	return (gu.group == nil || gu.group.Gid == fileUser.Gid) &&
		(gu.user == nil || gu.user.Uid == fileUser.Uid)
}

func findGroup(name string) bool {
	return !(strings.Index(name, ":") == -1 && strings.Index(name, ".") == -1)
}

func genToName(gu *groupUser) (to string) {
	to = fmt.Sprintf("%s:%s", gu.user, gu.group)
	if !findGroup(gu.name) || len(gu.name) > 0 && gu.name[0] == ':' {
		to = gu.name
		fmt.Printf("to = %s\n", to)
	}
	return
}

func (c *Chown) printVerbse(fileName string, gu *groupUser, w io.Writer) error {
	var st unix.Stat_t

	err := unix.Lstat(fileName, &st)
	if err != nil {
		//todo
		return err
	}

	if gu.group == nil && gu.user == nil {
		if c.IsVerbose() {
			fmt.Fprintf(w, "ownership of '%s' retained\n", fileName)
		}
		return nil
	}

	fileUser, err := user.LookupId(fmt.Sprintf("%d", st.Uid))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return nil
	}

	fileGroup, err := user.LookupGroupId(fmt.Sprintf("%d", st.Gid))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return nil
	}

	if c.IsChanges() {
		if !noChanages(gu, fileUser, fileGroup) {
			goto next
		}
		return nil
	}

	if noChanages(gu, fileUser, fileGroup) {

		fmt.Fprintf(w,
			"ownership of '%s' retained as %s\n",
			fileName, genToName(gu))
		return nil
	}

next:

	from := fmt.Sprintf("%s:%s", fileUser.Username, fileGroup.Name)
	if !findGroup(gu.name) {
		from = fileUser.Username
	}

	fmt.Fprintf(w,
		"changed ownership of '%s' from %s to %s\n",
		fileName, from,
		genToName(gu))

	return nil
}

// user: or user.
func isEndSplitter(name string) bool {
	last := len(name) - 1
	return len(name) > 0 && (name[last] == ':' || name[last] == '.')
}

func (c *Chown) Chown(name string, fileName string, u *User) error {

	w := u.W
	groupUser, err := getGroupUser(name)
	if err != nil {
		return err
	}

	uid, gid := -1, -1
	if groupUser.user != nil {
		uid, _ = strconv.Atoi(groupUser.user.Uid)
	}

	if groupUser.group != nil {
		gid, _ = strconv.Atoi(groupUser.group.Gid)
	}

	if gid == -1 && uid != -1 && isEndSplitter(name) {
		g, err := user.LookupGroupId(groupUser.user.Uid)
		if err != nil {
			return err
		}
		groupUser.group = (*chownGroup)(g)
		gid = uid
	}

	if c.IsChanges() || c.IsVerbose() {
		err = c.printVerbse(fileName, groupUser, w)
		if err != nil {
			return formatError(fileName, err)
		}
	}

	u.Uid = uid
	u.Gid = gid

	chown := os.Chown

	if c.NoDereference != nil && *c.NoDereference {
		chown = os.Lchown
	}

	err = chown(fileName, uid, gid)
	if err != nil {
		return formatError(fileName, err)
	}

	return nil
}

func Main(argv []string) {

	c, args := New(argv)

	err := checkArgs(args)
	if err != nil {
		utils.Die("%s\n", err)
	}

	u := User{}
	u.Init(os.Stdout)
	for _, a := range args[1:] {
		err := c.Chown(args[0], a, &u)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
