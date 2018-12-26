package main

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
	"strings"
)

func main() {
	v := flag.Bool("version", false, "output version information and exit")
	if *v {
		fmt.Printf("todo output version\n")
		os.Exit(0)
	}

	flag.Parse()
	args := flag.Args()

	output := "y\n"

	if len(args) > 0 {
		output = strings.Join(args, " ")
		output += "\n"
	}

	for {
		os.Stdout.Write([]byte(output))
	}
}
