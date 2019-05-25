package pwd

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"path/filepath"
)

type Pwd struct {
	Logical  bool
	Physical bool
}

func New(argv []string) (*Pwd, []string) {
	p := Pwd{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	command.Opt("L, logical", "print the value of $PWD if it names the current working\n"+
		"directory").Flags(flag.PosixShort).DefaultVar(&p.Logical, true)

	command.Opt("P, physical", "print the physical directory, without any symbolic links").
		Flags(flag.PosixShort).Var(&p.Physical)

	command.Parse(argv[1:])
	return &p, command.Args()
}

func Main(argv []string) {
	p, _ := New(argv)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("pwd: %s\n", err)
		return
	}

	if p.Physical {
		dir, err = filepath.EvalSymlinks(dir)
		if err != nil {
			fmt.Printf("pwd: %s\n", dir)
		}
	}

	if p.Logical {
		fmt.Printf("%s\n", dir)
	}
}
