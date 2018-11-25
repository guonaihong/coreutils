package main

import (
	"fmt"
	"github.com/guonaihong/flag"
	"github.com/guonaihong/gutil/file"
	"github.com/guonaihong/log"
	"io"
	"os"
)

func main() {
	append := flag.Bool("a, append", false, "append to the given FILEs, do not overwrite")
	//ignoreInterrupts := flag.Bool("i, ignore-interrupts", false, "ignore interrupt signals")//todo
	gzip := flag.Bool("g, gzip", false, "compressed archived log files")
	maxSize := flag.String("s, max-size", "0", "current file maximum write size")
	maxArchive := flag.Int("A, max-archive", 0, "How many archive files are saved")

	flag.Parse()

	var fileArch *log.File
	var buffer [4096]byte

	var w io.Writer

	args := flag.Args()
	fileName := args[0]
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

		fileArch = log.NewFile("", fileName, compress, int(size), *maxArchive)

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

	for {
		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			break
		}
		w.Write(buffer[:n])
	}
}
