package main

import (
	"github.com/guonaihong/coreutils/tee"
	"os"
)

func main() {
	tee.Main(os.Args)
}
