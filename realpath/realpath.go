package realpath

import (
	"fmt"
	"github.com/guonaihong/flag"
)

type RealPath struct {
	CanonicalizeExisting bool
	CanonicalizeMissing  bool
	Logical              bool
	Physical             bool
	Quiet                bool
	RelativeTo           string
	RelativeBase         bool
	NoSymlinks           bool
	Zero                 bool
}

func New(argv []string) (*RealPath, []string) {

	realPath := &RealPath{}

	command := flag.NewFlagSet(argv[0])

	command.Opt("e, canonicalize-existing", "all components of the path must exist").
		Flags(flag.PosixShort).Var(&realPath.CanonicalizeExisting)

	command.Opt("m, canonicalize-missing", "no path components need exist or be a directory").
		Flags(flag.PosixShort).Var(&realPath.CanonicalizeMissing)

	command.Opt("L, logical", "resolve '..' components before symlinks").
		Flags(flag.PosixShort).Var(&realPath.Logical)

	command.Opt("P, physical", "resolve symlinks as encountered (default)").
		Flags(flag.PosixShort).Var(&realPath.Physical)

	command.Opt("q, quiet", "suppress most error messages").
		Flags(flag.PosixShort).Var(&realPath.Quiet)

	command.Opt("relative-to", "print the resolved path relative to DIR").
		Flags(flag.PosixShort).Var(&realPath.RelativeTo)

	command.Opt("relative-base", "print absolute paths unless paths below DIR").
		Flags(flag.PosixShort).Var(&realPath.RelativeBase)

	command.Opt("s, strip, no-symlinks", "don't expand symlinks").
		Flags(flag.PosixShort).Var(&realPath.NoSymlinks)

	command.Opt("z, zero", "end each output line with NUL, not newline").
		Flags(flag.PosixShort).Var(&realPath.Zero)

	command.Parse(argv[1:])

	return realPath, command.Args()
}
