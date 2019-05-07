package hashcore

import (
	"github.com/guonaihong/flag"
)

type HashCore struct {
	Binary        *bool
	Check         *string
	Tag           *bool
	Text          *bool
	IgnoreMissing *bool
	Quiet         *bool
	Status        *bool
	Strict        *bool
	Warn          *bool
}

func New(argv []string, hashName string) (*HashCore, []string) {
	hash := &HashCore{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	hash.Binary = command.Opt("b, binary", "read in binary mode").
		Flags(flag.PosixShort).NewBool(false)

	hash.Check = command.Opt("c, check", "read MD5 sums from the FILEs and check them").
		Flags(flag.PosixShort).NewString("")

	hash.Tag = command.Opt("tag", "create a BSD-style checksum").
		Flags(flag.PosixShort).NewBool(false)

	hash.IgnoreMissing = command.Opt("ignore-missing", "don't fail or report status for missing files").
		Flags(flag.PosixShort).NewBool(false)

	hash.Quiet = command.Opt("quiet", "don't print OK for each successfully verified file").
		Flags(flag.PosixShort).NewBool(false)

	hash.Status = command.Opt("status", "don't output anything, status code shows success").
		Flags(flag.PosixShort).NewBool(false)

	hash.Strict = command.Opt("strict", "exit non-zero for improperly formatted checksum lines").
		Flags(flag.PosixShort).NewBool(false)

	hash.Warn = command.Opt("w, warn", "warn about improperly formatted checksum lines").
		Flags(flag.PosixShort).NewBool(false)

}
