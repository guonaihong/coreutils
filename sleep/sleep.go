package sleep

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"time"
)

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	version := command.Bool("version", false, "output version information and exit")
	command.Parse(argv[1:])

	if *version {
		fmt.Printf("todo: output version\n")
		os.Exit(0)
	}

	args := command.Args()
	sleepCore := func(arg string) {

		i := 0
		rt := time.Second
		for ; i < len(arg); i++ {
			if arg[i] >= '0' && arg[i] <= '9' || arg[i] == '.' {
				continue
			}

			i++
			break
		}

		i--
		if len(arg)-i > 1 {
			fmt.Printf("Invalid interval:%s\n", arg)
			os.Exit(1)
		}

		t := 0.0
		fmt.Sscanf(arg, "%f", &t)

		if len(arg)-i == 1 {
			switch arg[len(arg)-1] {
			case 's', 'S':
				rt = time.Second
			case 'm', 'M':
				rt = time.Minute
			case 'h', 'H':
				rt = time.Hour
			case 'd', 'D':
				rt = time.Hour * 24
			}

		}

		rt = time.Duration(float64(rt) * t)
		time.Sleep(rt)
	}

	for _, v := range args {
		sleepCore(v)
	}
}
