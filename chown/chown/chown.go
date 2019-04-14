package main

import (
	"github.com/guonaihong/coreutils/chown"
	"os"
)

func main() {
	chown.Main(os.Args)
}
