package chown

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type Chown struct {
	Changes *string
	//-f, --silent, --quiet
	Verbose        *bool
	Dereference    *bool
	NoDereference  *bool
	Form           *string
	NoPreserveRoot bool
	PreserveRoot   bool
	Reference      *string
	Recursive      *bool
	H              *bool
	L              *bool
	P              *bool
}

func New() (Chown, []string) {
	c := &Chown{}

	command := flag.NewFlagSet()

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
		Flags(flag.PosixShort).NewBool(false)

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

func (c *Chown) Chown(args []string) {
}

func Main(argv []string) {

	c, args := New(argv)

	c.Chown(args)
}
