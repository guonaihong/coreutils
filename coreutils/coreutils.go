package main

import (
	"github.com/guonaihong/coreutils/cat"
	"github.com/guonaihong/coreutils/paste"
	"github.com/guonaihong/coreutils/tr"
	"github.com/guonaihong/coreutils/uniq"
	"github.com/guonaihong/flag"
	"os"
)

func main() {
	parent := flag.NewParentCommand(os.Args[0])

	parent.SubCommand("cat", "Use the cat subcommand", func() {
		cat.Main(os.Args[1:])
	})

	parent.SubCommand("paste", "Use the paste subcommand", func() {
		paste.Main(os.Args[1:])
	})

	parent.SubCommand("uniq", "Use the uniq subcommand", func() {
		uniq.Main(os.Args[1:])
	})

	parent.SubCommand("tr", "Use the tr subcommand", func() {
		tr.Main(os.Args[1:])
	})

	parent.Parse(os.Args[1:])
}
