package hashcore

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"hash"
)

type Type int

const (
	Md5 Type = iota
	Sha1
	Sha256
	Sha224
	Sha384
	Sha512
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
	hash          hash.Hash
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

	command.Parse(argv[1:])

	return &command, command.Args()
}

func (h *HashCore) Check() {
}

func (h *HashCore) Hash(t Type, fileName string) error {

	fd, err := utils.OpenInFile(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	switch t {
	case Md5:
		h.hashVal = md5.New()
	case Sha1:
		h.hashVal = sha1.New()
	case Sha224:
		h.hashVal = sha256.New224()
	case Sha256:
		h.hashVal = sha256.New()
	case Sha384:
		h.hashVal = sha521.New384()
	case Sha512:
		h.hashVal = sha512.New()
	}

	io.Copy(h.hash, fd)
}

func Main(argv []string) {
}
