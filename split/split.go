package main

import (
	"github.com/guonaihong/flag"
)

func main() {
	flag.Int("a, suffix-length", 0, "use suffixes of length N (default 2)")
	flag.String("additional-suffix", "", "append an additional SUFFIX to file names")
	flag.String("b, bytes", "", "put SIZE bytes per output file")
	flag.String("C, line-bytes", "", "put at most SIZE bytes of lines per output file")
	flag.String("d", "", "use numeric suffixes starting at 0, not alphabetic")
	flag.String("numeric-suffixes", "", "same as -d, but allow setting the start value")
	flag.String("x", "", "use hex suffixes starting at 0, not alphabetic")
	flag.String("hex-suffixes", "", "same as -x, but allow setting the start value")
	flag.String("e, elide-empty-files", "", "do not generate empty output files with '-n'")
	flag.String("filter", "", "write to shell COMMAND; file name is $FILE")
	flag.String("l, lines", "", "put NUMBER lines per output file")
	flag.String("n, number", "", "generate CHUNKS output files; see explanation below")
	flag.String("t, separator", "", "use SEP instead of newline as the record separator; '\000' (zero) specifies the NUL character")
	flag.Bool("u, unbuffered", false, "immediately copy input to output with '-n r/...'")
	flag.Bool("verbose", false, "print a diagnostic just before each output file is opened")

	flag.Parse()
}
