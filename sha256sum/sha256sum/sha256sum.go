package main

import (
	"github.com/guonaihong/coreutils/sha256sum"
	"os"
)

func main() {
	sha256sum.Main(os.Args)
}
