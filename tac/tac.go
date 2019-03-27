package tac

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"io/ioutil"
	"os"
)

const bufSize = 8092

type Tac struct {
	Before    *bool
	Regex     *string
	Separator *string
}

func New(argv []string) (*Tac, []string) {
	t := Tac{}
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	t.Before = command.Opt("b, before", "attach the separator before instead of after").
		Flags(flag.PosixShort).
		NewBool(false)
	t.Regex = command.Opt("r, regex", "interpret the separator as a regular expression").
		Flags(flag.PosixShort).
		NewString("")
	t.Separator = command.Opt("s, separator", "use STRING as the separator instead of newline").
		Flags(flag.PosixShort).
		NewString("\n")

	command.Parse(argv[1:])
	args := command.Args()

	args = utils.NewArgs(args)
	return &t, args
}

func printOffset(rs io.ReadSeeker, w io.Writer, buf []byte, start, end int64) error {

	curPos, err := rs.Seek(0, 1)
	if err != nil {
		return err
	}

	_, err = rs.Seek(start, 0)
	if err != nil {
		return err
	}

	defer rs.Seek(curPos, 0)

	for {

		if start >= end {
			break
		}

		needRead := end - start
		if int(needRead) > len(buf) {
			needRead = int64(len(buf))
		}

		n, e := rs.Read(buf[:needRead])
		if e != nil {
			break
		}

		w.Write(buf[:n])
		start += int64(n)
	}
	return nil
}

func readFromTailStdin(r io.Reader, w io.Writer, sep []byte) error {
	all, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	offset := make([]int, 0, 50)

	for i := 0; i < len(all); {
		pos := bytes.Index(all[i:], sep)
		if pos == -1 {
			break
		}

		offset = append(offset, i+pos)
		i += pos + 1
	}

	if len(offset) == 0 {
		offset = append(offset, len(all))
		sep = []byte("")
	}

	right := len(all)

	for i := len(offset) - 1; i >= 0; i-- {
		start := offset[i] + len(sep)

		w.Write(all[start:right])
		right = offset[i] + len(sep)
	}

	w.Write(all[0:right])
	return nil
}

func readFromTail(rs io.ReadSeeker, w io.Writer, sep []byte) error {

	tail, err := rs.Seek(0, 2)
	if err != nil {
		return err
	}

	head := tail

	buf := make([]byte, bufSize+len(sep))
	buf2 := make([]byte, bufSize)

	for head > 0 {

		minRead := head
		if minRead > bufSize {
			minRead = bufSize
		}

		_, err := rs.Seek(-minRead, 1)
		if err != nil {
			return err
		}

		n, err := rs.Read(buf[:minRead])
		if err != nil {
			return err
		}

		head -= minRead
		rs.Seek(-minRead, 1)

		t := n
		h := n

		for {
			pos := bytes.LastIndex(buf[:h], sep)

			if pos == -1 {
				//not found
				break
			}

			if pos >= 0 {

				w.Write(buf[pos+len(sep) : t])
				if l := t - pos - len(sep); l > 0 {
					tail -= int64(l)
				}

				if !bytes.Equal(buf[pos:t], sep) {
					t = pos - len(sep)
				}

				h = pos - 1

				if tail > head+minRead {
					err = printOffset(rs, w, buf2, head+minRead, tail)
					if err != nil {
						return err
					}
					tail = head + minRead
				}

				if pos == 0 {
					break
				}
			}

		}
	}

	if tail > 0 {
		printOffset(rs, w, buf2, 0, tail)
	}
	return nil
}

func (t *Tac) Tac(rs io.ReadSeeker, w io.Writer) error {
	if t.Separator != nil {
		err := readFromTail(rs, w, []byte(*t.Separator))
		if err != nil {
			return err
		}
	}
	return nil
}

func Main(argv []string) {

	tac, args := New(argv)

	for _, fileName := range args {
		f, err := utils.OpenInputFd(fileName)
		if err != nil {
			utils.Die("tac: %s\n", err)
		}

		err = tac.Tac(f, os.Stdout)
		if err != nil {
			utils.Die("tac: %s\n", err)
		}
		utils.CloseInputFd(f)
	}
}
