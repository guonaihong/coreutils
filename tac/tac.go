package tac

import (
	"bytes"
	"fmt"
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
	BufSize   int
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

		fmt.Printf("#######= (%s)\n", buf[:n])
		w.Write(buf[:n])
		start += int64(n)
	}
	return nil
}

func readFromTailStdin(r io.Reader, w io.Writer, sep []byte, before bool) error {
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
		i += pos + len(sep)
	}

	if len(offset) == 0 {
		offset = append(offset, len(all))
		sep = []byte("")
	}

	right := len(all)

	for i := len(offset) - 1; i >= 0; i-- {
		start := offset[i]
		if !before {
			start += len(sep)
		}

		w.Write(all[start:right])

		right = start
	}

	w.Write(all[0:right])
	return nil
}

func (t *Tac) readFromTail(rs io.ReadSeeker, w io.Writer, sep []byte, before bool) error {

	tail, err := rs.Seek(0, 2)
	if err != nil {
		return readFromTailStdin(rs, w, sep, before)
	}

	head := tail

	buf := make([]byte, t.BufSize+len(sep))
	buf2 := make([]byte, t.BufSize)

	fmt.Printf("##################%d\n", tail)
	for head > 0 {

		minRead := head
		if minRead > int64(t.BufSize) {
			minRead = int64(t.BufSize)
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

		right := n
		h := n

		fmt.Printf("======================:%d\n", minRead)
		for {
			pos := bytes.LastIndex(buf[:h], sep)

			if pos == -1 {
				//not found
				break
			}

			if pos >= 0 {

				start := pos + len(sep)

				/*
					l := right - start
						fmt.Printf("%p, head = %d, pos = %d, start = %d, right = %d, right - start = %d, (%s)\n",
							t, head, pos, start, right, right-start, buf[start:])
				*/

				if !bytes.Equal(buf[start:right], sep) {
					right = pos + len(sep)
				}

				h = pos

				/*
					fmt.Printf("1.l = %d, tail = %d, head = %d, start = %d pos = %d, buf(%s)\n",
						l, tail, head, int(head)+pos+len(sep), pos, buf[:n])
				*/

				start = int(head) + pos

				if !before {
					start += len(sep)
				}

				if tail > int64(start) {
					err = printOffset(rs, w, buf2, int64(start), tail)
					if err != nil {
						return err
					}
					//Move tail position
					tail = int64(start)
				}
				//fmt.Printf("2.l = %d, tail = %d, head = %d, start = %d\n", l, tail, head, int(head)+pos+len(sep))

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
	before := false
	if t.Before != nil {
		before = *t.Before
	}

	if t.BufSize == 0 {
		t.BufSize = bufSize
	}

	if t.Separator != nil && 2*len(*t.Separator) > t.BufSize {
		t.BufSize = 2 * len(*t.Separator)
	}

	if t.Separator != nil {
		err := t.readFromTail(rs, w, []byte(*t.Separator), before)
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
