package main

import (
	"github.com/guonaihong/coreutils/base32"
	"github.com/guonaihong/coreutils/base64"
	"github.com/guonaihong/coreutils/basename"
	"github.com/guonaihong/coreutils/cat"
	"github.com/guonaihong/coreutils/chgrp"
	"github.com/guonaihong/coreutils/chown"
	"github.com/guonaihong/coreutils/cut"
	"github.com/guonaihong/coreutils/dirname"
	"github.com/guonaihong/coreutils/echo"
	"github.com/guonaihong/coreutils/head"
	"github.com/guonaihong/coreutils/link"
	"github.com/guonaihong/coreutils/md5sum"
	"github.com/guonaihong/coreutils/paste"
	"github.com/guonaihong/coreutils/pwd"
	"github.com/guonaihong/coreutils/rmdir"
	"github.com/guonaihong/coreutils/seq"
	"github.com/guonaihong/coreutils/sha1sum"
	"github.com/guonaihong/coreutils/sha224sum"
	"github.com/guonaihong/coreutils/sha256sum"
	"github.com/guonaihong/coreutils/sha384sum"
	"github.com/guonaihong/coreutils/sha512sum"
	"github.com/guonaihong/coreutils/shuf"
	"github.com/guonaihong/coreutils/sleep"
	"github.com/guonaihong/coreutils/tac"
	"github.com/guonaihong/coreutils/tail"
	"github.com/guonaihong/coreutils/tee"
	"github.com/guonaihong/coreutils/touch"
	"github.com/guonaihong/coreutils/tr"
	"github.com/guonaihong/coreutils/true"
	"github.com/guonaihong/coreutils/uname"
	"github.com/guonaihong/coreutils/uniq"
	"github.com/guonaihong/coreutils/unlink"
	"github.com/guonaihong/coreutils/whoami"
	"github.com/guonaihong/coreutils/yes"
	"github.com/guonaihong/flag"
	"os"
)

func main() {
	parent := flag.NewParentCommand(os.Args[0])

	parent.SubCommand("base32", "Use the base32 subcommand", func() {
		base32.Main(os.Args[1:])
	})

	parent.SubCommand("base64", "Use the base64 subcommand", func() {
		base64.Main(os.Args[1:])
	})

	parent.SubCommand("cat", "Use the cat subcommand", func() {
		cat.Main(os.Args[1:])
	})

	parent.SubCommand("chown", "Use the chown subcommand", func() {
		chown.Main(os.Args[1:])
	})

	parent.SubCommand("chgrp", "Use the chgrp subcommand", func() {
		chgrp.Main(os.Args[1:])
	})

	parent.SubCommand("cut", "Use the cut subcommand", func() {
		cut.Main(os.Args[1:])
	})

	parent.SubCommand("paste", "Use the paste subcommand", func() {
		paste.Main(os.Args[1:])
	})

	parent.SubCommand("pwd", "Use the pwd subcommand", func() {
		pwd.Main(os.Args[1:])
	})

	parent.SubCommand("rmdir", "Use the rmdir subcommand", func() {
		rmdir.Main(os.Args[1:])
	})

	parent.SubCommand("basename", "Use the basename subcommand", func() {
		basename.Main(os.Args[1:])
	})

	parent.SubCommand("dirname", "Use the dirname subcommand", func() {
		dirname.Main(os.Args[1:])
	})

	parent.SubCommand("echo", "Use the echo subcommand", func() {
		echo.Main(os.Args[1:])
	})

	parent.SubCommand("head", "Use the head subcommand", func() {
		head.Main(os.Args[1:])
	})

	parent.SubCommand("link", "Use the link subcommand", func() {
		link.Main(os.Args[1:])
	})

	parent.SubCommand("md5sum", "Use the md5sum subcommand", func() {
		md5sum.Main(os.Args[0:])
	})

	parent.SubCommand("seq", "Use the seq subcommand", func() {
		seq.Main(os.Args[1:])
	})

	parent.SubCommand("sha1sum", "Use the sha1sum subcommand", func() {
		sha1sum.Main(os.Args[1:])
	})

	parent.SubCommand("sha224sum", "Use the sha224sum subcommand", func() {
		sha224sum.Main(os.Args[1:])
	})

	parent.SubCommand("sha256sum", "Use the sha256sum subcommand", func() {
		sha256sum.Main(os.Args[1:])
	})

	parent.SubCommand("sha384sum", "Use the sha384sum subcommand", func() {
		sha384sum.Main(os.Args[1:])
	})

	parent.SubCommand("sha512sum", "Use the sha512sum subcommand", func() {
		sha512sum.Main(os.Args[1:])
	})

	parent.SubCommand("shuf", "Use the shuf subcommand", func() {
		shuf.Main(os.Args[1:])
	})

	parent.SubCommand("sleep", "Use the sleep subcommand", func() {
		sleep.Main(os.Args[1:])
	})

	parent.SubCommand("tac", "Use the tac subcommand", func() {
		tac.Main(os.Args[1:])
	})

	parent.SubCommand("tail", "Use the tail subcommand", func() {
		tail.Main(os.Args[1:])
	})

	parent.SubCommand("tee", "Use the tee subcommand", func() {
		tee.Main(os.Args[1:])
	})

	parent.SubCommand("touch", "Use the touch subcommand", func() {
		touch.Main(os.Args[1:])
	})

	parent.SubCommand("tr", "Use the tr subcommand", func() {
		tr.Main(os.Args[1:])
	})

	parent.SubCommand("true", "Use the true subcommand", func() {
		true.Main(os.Args[1:])
	})

	parent.SubCommand("uname", "Use the uname subcommand", func() {
		uname.Main(os.Args[1:])
	})

	parent.SubCommand("uniq", "Use the uniq subcommand", func() {
		uniq.Main(os.Args[1:])
	})

	parent.SubCommand("unlink", "Use the unlink subcommand", func() {
		unlink.Main(os.Args[1:])
	})

	parent.SubCommand("whoami", "Use the whoami subcommand", func() {
		whoami.Main(os.Args[1:])
	})

	parent.SubCommand("yes", "Use the yes subcommand", func() {
		yes.Main(os.Args[1:])
	})

	parent.Parse(os.Args[1:])
}
