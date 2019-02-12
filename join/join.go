package main

import (
	"gihub.com/guonaihong/flag"
)

func main() {
	flag.Int("a", 0, "also print unpairable lines from file FILENUM, "+
		"where FILENUM is 1 or 2, corresponding to FILE1 or FILE2")
	flag.String("e", "", "replace missing input fields with EMPTY")
	flag.Bool("i, ignore-case", false, "ignore differences in case "+
		"when comparing fields")
	flag.Int("j", 0, "equivalent to '-1 FIELD -2 FIELD'")
}
