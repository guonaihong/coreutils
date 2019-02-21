package dirname

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"path/filepath"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	zero := command.Bool("z, zero", false, "end each output line with NUL, not newline")
	verion := command.Bool("version", false, "output version information and exit")

	command.Parse(argv[1:])

	if *verion {
		fmt.Printf("todo output version \n")
		os.Exit(0)
	}

	args := command.Args()
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
