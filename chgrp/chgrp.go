package chgrp

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
)

type Chgrp struct {
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
	Gid int
	W   io.Writer
}

func (u *User) Init(w io.Writer) {
	u.Gid = -1
	u.W = w
}

func New(argv []string) (*Chgrp, []string) {
	c := &Chgrp{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	//-f, --silent, --quiet
	c.Quiet = command.Opt("f, silent, quiet",
		"suppress most error messages").
		Flags(flag.PosixShort).NewBool(false)

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

type chgrpGroup user.Group

func (g *chgrpGroup) String() string {
	if g == nil {
		return ""
	}
	return g.Name
}

type group struct {
	group *chgrpGroup
	name  string
}

func getGroupFromName(name string) (*user.Group, error) {

	var err error
	var g *user.Group

	g, err = user.LookupGroup(name)
	if err != nil {
		if _, ok := err.(user.UnknownGroupError); ok {
			err = fmt.Errorf("chgrp: invalid group: '%s'", name)
		}
		return nil, err
	}

	return g, nil
}

func checkArgs(args []string) error {
	if len(args) == 0 {
		return errors.New("chgrp: missing operand")
	}

	if len(args) == 1 {
		return fmt.Errorf("chgrp: missing operand after '%s'", args[0])
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
			// The error in gnu chgrp is 'No such file or directory'
			return fmt.Errorf("chgrp: cannot access '%s': No such file or directory",
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
		return fmt.Errorf("chgrp: changing ownership of '%s': %s", fileName, err)
	}

	return err
}

func (c *Chgrp) IsNoDereference() bool {
	return c.NoDereference != nil && *c.NoDereference
}

func (c *Chgrp) IsFrom() bool {
	return c.From != nil && len(*c.From) > 0
}

func (c *Chgrp) IsChanges() bool {
	return c.Changes != nil && *c.Changes
}

func (c *Chgrp) IsVerbose() bool {
	return c.Verbose != nil && *c.Verbose
}

func (c *Chgrp) IsPreserveRoot() bool {
	return c.PreserveRoot != nil && *c.PreserveRoot
}

func (c *Chgrp) IsReference() bool {
	return c.Reference != nil && len(*c.Reference) > 0
}

func noChanages(gu *group, fileGroup *user.Group) bool {
	return (gu.group == nil || gu.group.Gid == fileGroup.Gid)

}

func genToName(gu *group) string {
	return gu.group.Name
}

func (c *Chgrp) printVerbse(
	fileName string,
	canChanges bool,
	st *unix.Stat_t,
	gu *group) (err error) {

	if gu == nil || gu.group == nil {
		if c.IsVerbose() {
			err = fmt.Errorf("ownership of '%s' retained\n", fileName)
		}
		return err
	}

	fileGroup, err := user.LookupGroupId(fmt.Sprintf("%d", st.Gid))
	if err != nil {
		return err
	}

	if !canChanges {
		return fmt.Errorf("ownership of '%s' retained as %s\n",
			fileName, genToName(gu))
	}

	if c.IsChanges() {
		if !noChanages(gu, fileGroup) {
			goto next
		}
		return nil
	}

	if noChanages(gu, fileGroup) {

		return fmt.Errorf("ownership of '%s' retained as %s\n",
			fileName, genToName(gu))
	}

next:

	from := fmt.Sprintf("%s", fileGroup.Name)

	return fmt.Errorf("changed ownership of '%s' from %s to %s\n",
		fileName, from,
		genToName(gu))
}

func parseGid(g *user.Group) (gid int) {
	gid = -1

	if g != nil {
		gid, _ = strconv.Atoi(g.Gid)
	}

	return
}

func (c *Chgrp) getGidFromName() (gid int, err error) {

	gid = -1

	g, err := getGroupFromName(*c.From)
	if err != nil {
		return
	}

	gid = parseGid(g)
	return
}

func (c *Chgrp) genGidFromFile(fileName string) (gid string, err error) {
	var st unix.Stat_t

	stat := unix.Stat
	if c.IsNoDereference() {
		stat = unix.Lstat
	}

	err = stat(fileName, &st)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", st.Gid), nil
}

func writeErrorToUser(u *User, err error) {
	if u != nil && u.W != nil {
		d := []byte{}
		if err != nil {
			d = []byte(err.Error())
		}
		u.W.Write(d)
	}
}

func (c *Chgrp) Chgrp(name string, fileName string, u *User) (err error) {

	//w := u.W
	defer func() {
		writeErrorToUser(u, err)
	}()

	if c.IsPreserveRoot() && fileName == "/" {
		return fmt.Errorf("chgrp: it is dangerous to operate recursively on '/'\n" +
			"chgrp: use --no-preserve-root to override this failsafe\n")
	}

	g, err := getGroupFromName(name)
	if err != nil {
		return err
	}

	gid := parseGid(g)

	// Convenient to write test procedures
	u.Gid = gid

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
		gid2, err := c.getGidFromName()
		if err != nil {
			//todo format error
			return err
		}

		if !(gid2 != -1 && uint32(gid2) == st.Gid) {
			canChanges = false
		}
	}

	if c.IsChanges() || c.IsVerbose() {
		err = c.printVerbse(fileName, canChanges, &st, &group{group: (*chgrpGroup)(g)})
		if err != nil {
			return formatError(fileName, err)
		}
	}

	if !canChanges {
		return nil
	}

	chown := os.Chown

	if c.IsNoDereference() {
		chown = os.Lchown
	}

	err = chown(fileName, -1, gid)
	if err != nil {
		return formatError(fileName, err)
	}

	if c.Recursive != nil && *c.Recursive {
		stat := os.Lstat
		if c.L != nil && *c.L {
			stat = os.Stat
		}

		walk := utils.NewWalk(c.H != nil && *c.H, stat)
		walk.Walk(name, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			return chown(fileName, -1, gid)
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
		gid, err := c.genGidFromFile(*c.Reference)
		if err != nil {
			fmt.Printf("%s\n", err)
			errCode = 1
		}

		args = append([]string{gid}, args...)
	}

	for _, a := range args[1:] {
		err := c.Chgrp(args[0], a, &u)
		if err != nil {
			if !*c.Quiet {
				fmt.Printf("%s\n", err)
			}
			errCode = 1
		}
	}
}
