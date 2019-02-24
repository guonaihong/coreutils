package paste

import (
	"bufio"
	"bytes"
	"github.com/guonaihong/flag"
	"os"
	"sync"
)

func Main(argv []string) {

	var lineDelim = '\n'

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	delim := command.Opt("d, delimiters", "reuse characters from LIST instead of TABs").
		NewString("\t")

	zeroTerminated := command.Opt("z, zero-terminated", "line delimiter is NUL, not newline").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	if *zeroTerminated {
		lineDelim = '\000'
	}

	args := command.Args()

	wg := sync.WaitGroup{}

	wg.Add(len(args))

	resultChan := make([]chan string, len(args))

	for k, _ := range resultChan {
		resultChan[k] = make(chan string, 1)
	}

	defer wg.Wait()

	for id, fileName := range args {
		go func(id int, fileName string) {
			defer wg.Done()
			defer close(resultChan[id])

			file, err := os.Open(fileName)
			if err != nil {
				os.Exit(1)
			}

			defer file.Close()

			br := bufio.NewReader(file)

			for {

				l, err := br.ReadBytes(byte(lineDelim))
				if err != nil && len(l) == 0 {
					break
				}

				resultChan[id] <- string(l)
			}

		}(id, fileName)

	}

	out := &bytes.Buffer{}

	for {
		allQuit := true
		for k, v := range resultChan {
			line, ok := <-v

			if ok {
				allQuit = false
			}
			if k != 0 {
				out.WriteString(*delim)
			}
			out.WriteString(line)
		}

		if allQuit {
			break
		}

		out.WriteString("\n")
		os.Stdout.Write(out.Bytes())
		out.Reset()
	}
}
