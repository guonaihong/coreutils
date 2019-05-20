package main

import (
	"github.com/guonaihong/coreutils/sha512sum"
	"os"
)

func main() {
	sha512sum.Main(os.Args)
}
