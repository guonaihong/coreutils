package hashcore

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"hash"
	"io"
	"os"
	"strings"
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

func (h *HashCore) formatError(t Type, fileName string, err error) string {
	switch e := err.(type) {
	case *os.PathError:
		err = e.Err
		if os.IsNotExist(err) {

			return fmt.Sprintf("%ssum: %s: No such file or directory\n",
				strings.ToLower(t.String()), fileName)
		}
	}
	return err.Error()
}

func (h *HashCore) CheckHash(t Type, formatFail *int, fileName string, w io.Writer) error {

	fd, err := utils.OpenFile(fileName)
	if err != nil {
		if h.IsStatus() {
			os.Exit(1)
		}
		return errors.New(h.formatError(t, fileName, err))
	}

	defer fd.Close()

	br := bufio.NewReader(fd)

	var out bytes.Buffer

	var checkFail, readFail, done int

	hashName := strings.ToLower(t.String())
	defer func(f, c, read, done *int) {
		if *f > 0 {
			fmt.Fprintf(w, "%ssum: WARNING: %d line is improperly formatted\n",
				hashName, *f)
		}

		if *c > 0 {
			fmt.Fprintf(w, "%ssum: WARNING: %d computed checksums did NOT match\n",
				hashName, *c)
		}

		if *read > 0 {
			fmt.Fprintf(w, "%ssum: WARNING: %d listed file could not be read\n",
				hashName, *read)
		}

		if h.IsIgnoreMissing() && *done == 0 {
			fmt.Fprintf(w, "%ssum: %s: no file was verified\n",
				hashName, fileName)
		}
	}(formatFail, &checkFail, &readFail, &done)

	for no := 1; ; no++ {

		l, e := br.ReadBytes('\n')

		if e != nil && len(l) == 0 {
			break
		}

		hashAndFile := bytes.Fields(l)

		var fileName string

		switch {
		case len(hashAndFile) == 2:
			fileName = string(hashAndFile[1])
			h.Binary = utils.Bool(false)
			if len(fileName) > 0 && fileName[0] == '*' {
				fileName = fileName[1:]
				h.Binary = utils.Bool(true)
			}
		case len(hashAndFile) == 3:
			_, err := fmt.Sscanf(string(hashAndFile[1]), "(%s)", &fileName)
			if err != nil {
				return err
			}

			if h.Warn != nil && *h.Warn {
				fmt.Fprintf(w, "%sum: %s: %d: improperly formatted MD5 checksum line\n",
					hashName, fileName, no)
			}
		default:
			(*formatFail)++
		}

		err := h.Hash(t, fileName, &out)
		if err != nil {
			if !h.IsIgnoreMissing() {
				fmt.Fprintf(w, h.formatError(t, fileName, err))
				fmt.Fprintf(w, "%s: FAILED open or read\n", fileName)
				readFail++
			}
			if h.IsStatus() {
				os.Exit(1)
			}
			continue
		}

		if bytes.Equal(out.Bytes(), l) {
			if !h.IsQuiet() {
				fmt.Fprintf(w, "%s: OK\n", fileName)
			}
			done++
		} else {
			if !h.IsIgnoreMissing() {
				fmt.Fprintf(w, "%s: FAILED\n", fileName)
			}
		}
	}
	return nil
}

func (h *HashCore) IsStatus() bool {
	return h.Status != nil && *h.Status
}

func (h *HashCore) IsIgnoreMissing() bool {
	return h.IgnoreMissing != nil && *h.IgnoreMissing
}

func (h *HashCore) IsQuiet() bool {
	return h.Quiet != nil && *h.Quiet
}

func (h *HashCore) IsCheck() bool {
	return h.Check != nil && len(*h.Check) > 0
}

func (h *HashCore) IsTag() bool {
	return h.Tag != nil && *h.Tag
}

func (h *HashCore) Hash(t Type, fileName string, w io.Writer) error {

	asterisk := "  "
	if h.Binary != nil && *h.Binary {
		asterisk = " *"
	}

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
		fmt.Fprintf(w, "%x%s%s\n", h.hash.Sum(nil), asterisk, fileName)
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

	var isFormatFail bool
	for _, fileName := range args {
		var err error
		var formatFail int

		if hash.IsCheck() {
			err = hash.CheckHash(t, &formatFail, fileName, os.Stdout)
			isFormatFail = formatFail > 0
		} else {
			hash.Hash(t, fileName, os.Stdout)
		}

		if err != nil {
			os.Stdout.Write([]byte(err.Error()))
		}
	}

	if hash.Strict != nil && *hash.Strict {
		if isFormatFail {
			os.Exit(1)
		}
	}
}
