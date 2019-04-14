package chown

import (
	"errors"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
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

type groupUser struct {
	group *user.Group
	user  *user.User
}

func getGroupUser(name string) (*groupUser, error) {

	var names []string

	if name == ":" || name == "." {
		return &groupUser{}, nil
	}

	if strings.Index(name, ":") != -1 {
		names = strings.Split(name, ":")
	}

	if strings.Index(name, ".") != -1 {
		names = strings.Split(name, ".")
	}

	name1 := name
	if len(names) > 0 {
		name1 = names[0]
	}

	u, err := user.Lookup(name1)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			err = fmt.Errorf("chown: invalid user: '%s'", name)
		}
		return nil, err
	}

	if len(names) == 0 {
		return &groupUser{user: u}, nil
	}

	var g *user.Group
	if len(names[1]) > 0 {
		g, err = user.LookupGroup(names[1])
		if err != nil {
			if _, ok := err.(user.UnknownGroupError); ok {
				err = fmt.Errorf("chown: invalid group: '%s'", name)
			}
			return nil, err
		}
	}

	return &groupUser{user: u, group: g}, nil
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

func (c *Chown) Chown(name string, fileName string) error {

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

	err = os.Chown(fileName, uid, gid)
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

	for _, a := range args[1:] {
		err := c.Chown(args[0], a)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
