package basename

import (
	"bytes"
	"github.com/guonaihong/flag"
	"os"
	"path/filepath"
	"strings"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	multiple := command.Opt("a, multiple", "support multiple arguments and treat each as a NAME").
		Flags(flag.PosixShort).NewBool(false)

	suffix := command.Opt("s, suffix", "remove a trailing SUFFIX; implies -a").
		NewString("")

	zero := command.Opt("z, zero", "end each output line with NUL, not newline").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	args := command.Args()

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
