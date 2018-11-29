package main

import (
	"bufio"
	"bytes"
	"github.com/guonaihong/flag"
	"os"
	"sync"
)

func main() {
	delim := flag.String("d, delimiters", "\t", "reuse characters from LIST instead of TABs")
	flag.Parse()

	args := flag.Args()

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

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				resultChan[id] <- scanner.Text()
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
