package main

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"path/filepath"
)

func main() {
	zero := flag.Bool("z, zero", false, "end each output line with NUL, not newline")
	verion := flag.Bool("version", false, "output version information and exit")

	flag.Parse()

	if *verion {
		fmt.Printf("todo output version \n")
		os.Exit(0)
	}

	args := flag.Args()
	out := bytes.Buffer{}

	for _, v := range args {
		if len(v) > 1 && v[len(v)-1] == '/' || v[len(v)-1] == '\\' {
			v = v[:len(v)-1]
		}

		newDir := filepath.Dir(v)

		out.Write([]byte(newDir))

		if *zero {
			out.WriteByte(0)
		} else {
			out.WriteByte('\n')
		}
	}

	os.Stdout.Write(out.Bytes())
}
