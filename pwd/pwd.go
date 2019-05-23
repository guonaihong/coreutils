package main

import (
	"fmt"
	"github.com/guoanihong/flag"
	"os"
)

type Pwd struct {
	Logical  bool
	Physical bool
}

func New(argv []string) (*Pwd, []string) {
	p := Pwd{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	command.Opt("L, logical", "print the value of $PWD if it names the current working\n"+
		"directory").Flags(flag.PosixShort).Var(&p.Logical)

	command.Opt("P, physical", "print the physical directory, without any symbolic links").
		Flags(flag.PosixShort).Var(&p.Physical)

	command.Parse(argv[1:])
	return &p, command.Args()
}

func Main(argv []string) {
	p, args := New(argv)
}
