package main

import (
	"bytes"
	"github.com/guonaihong/flag"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	multiple := flag.Bool("a, multiple", false, "support multiple arguments and treat each as a NAME")
	suffix := flag.String("s, suffix", "", "remove a trailing SUFFIX; implies -a")
	zero := flag.Bool("z, zero", false, "end each output line with NUL, not newline")
	flag.Parse()

	args := flag.Args()

	var out bytes.Buffer

	baseNameCore := func(baseName string) {
		if strings.HasSuffix(baseName, *suffix) {
			baseName = baseName[:len(baseName)-len(*suffix)]
		}

		out.WriteString(baseName)
		if *zero {
			out.WriteByte(0)
		} else {
			out.WriteByte('\n')
		}
	}

	if *multiple {
		for _, a := range args {
			baseName := filepath.Base(a)

			baseNameCore(baseName)
		}
	}

	if out.Len() == 0 && len(args) > 0 {
		baseName := filepath.Base(args[0])

		if len(args) == 2 {
			*suffix = args[1]
		}

		baseNameCore(baseName)
	}

	os.Stdout.Write(out.Bytes())
}
