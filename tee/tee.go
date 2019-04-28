package tee

import (
	"fmt"
	"github.com/guonaihong/flag"
	"github.com/guonaihong/gutil/file"
	"github.com/guonaihong/log"
	"io"
	"os"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	append := command.Opt("a, append", "append to the given FILEs, do not overwrite").
		Flags(flag.PosixShort).NewBool(false)

	ignoreInterrupts := command.Opt("i, ignore-interrupts", "ignore interrupt signals").
		Flags(flag.PosixShort).NewBool(false)

	gzip := command.Opt("g, gzip", "compressed archived log files").
		Flags(flag.PosixShort).NewBool(false)

	maxSize := command.Opt("s, max-size", "current file maximum write size").
		NewString("0")

	maxArchive := command.Opt("A, max-archive", "How many archive files are saved").
		NewInt64(0)

	command.Parse(argv[1:])

	var fileArch *log.File
	var buffer [4096]byte
	var fileName string
	var w io.Writer

	if *ignoreInterrupts {
		ignore()
	}

	args := command.Args()
	if len(args) == 0 {
		w = os.Stdout
		goto end
	}

	fileName = args[0]
	if *maxSize != "" {
		size, err := file.ParseSize(*maxSize)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
			return
		}

		var compress log.CompressType
		compress = log.NotCompress
		if *gzip {
			compress = log.Gzip
		}

		fileArch = log.NewFile("", fileName, compress, int(size), int(*maxArchive))

		w = fileArch
	} else {
		isAppend := 0
		if *append {
			isAppend = os.O_APPEND
		}

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|isAppend, 0644)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		w = file
	}

end:
	for {
		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			break
		}
		w.Write(buffer[:n])
		if len(args) != 0 {
			os.Stdout.Write(buffer[:n])
		}
	}
}
