package main

import (
	"fmt"
	"github.com/guoanihong/flag"
	"os"
)

type Pwd struct {
	L *bool
	P *bool
}

func New(argv []string) (*Pwd, []string) {
	p := Pwd{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	command.Opt()
}

func Main(argv []string) {
}
