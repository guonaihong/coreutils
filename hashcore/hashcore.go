package hashcore

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"hash"
	"io"
	"os"
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

func (t Type) String() string {
	switch t {
	case Md5:
		return "MD5"
	case Sha1:
		return "SHA1"
	case Sha256:
		return "SHA256"
	case Sha224:
		return "SHA224"
	case Sha384:
		return "SHA384"
	case Sha512:
		return "SHA512"
	}

	panic("unkown")
}

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

	return hash, command.Args()
}

func (h *HashCore) CheckHash(fileName string) error {

	fd, err := utils.OpenFile(fileName)
	if err != nil {
		return err
	}

	defer fd.Close()

	br := bufio.NewReader(fd)

	for {

		l, e := br.ReadBytes('\n')

		if e != nil && len(l) == 0 {
			break
		}

		hashAndFile := bytes.Fields(l)
		if len(hashAndFile) == 2 {

		}
	}
}

func (h *HashCore) IsTag() bool {
	return h.Tag != nil && *h.Tag
}

func (h *HashCore) Hash(t Type, fileName string, w io.Writer) error {

	fd, err := utils.OpenFile(fileName)
	if err != nil {
		return err
	}
	defer fd.Close()

	switch t {
	case Md5:
		h.hash = md5.New()
	case Sha1:
		h.hash = sha1.New()
	case Sha224:
		h.hash = sha256.New224()
	case Sha256:
		h.hash = sha256.New()
	case Sha384:
		h.hash = sha512.New384()
	case Sha512:
		h.hash = sha512.New()
	}

	io.Copy(h.hash, fd)

	if !h.IsTag() {
		fmt.Fprintf(w, "%x  %s\n", h.hash.Sum(nil), fileName)
	} else {
		fmt.Fprintf(w, "%s (%s) = %x\n", t, fileName, h.hash.Sum(nil))
	}

	return nil
}

func Main(argv []string, t Type) {
	hash, args := New(argv, t.String())
	if len(args) == 0 {
		args = append(args, "-")
	}

	for _, a := range args {
		hash.Hash(t, a, os.Stdout)
	}
}
