package main

import (
	"github.com/guonaihong/coreutils/sha1sum"
	"os"
)

func main() {
	sha1sum.Main(os.Args)
}
