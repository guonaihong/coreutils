package main

import "fmt"

type Readlink struct {
	Canonicalize         bool
	CanonicalizeExisting bool
	CanonicalizeMissing  bool
	NoNewline            bool
	Quiet                bool
	Silent               bool
	Verbose              bool
	Zero                 bool
}

func New() {
}
