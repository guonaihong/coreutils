package utils

import (
	"testing"
)

func TestHeadParseSize(t *testing.T) {
	n, _ := HeadParseSize("1")
	if n != Size(1) {
		t.Errorf("Value is not equal to -1\n")
	}

	n, err := HeadParseSize("1kB")
	if n != KB || err != nil {
		t.Errorf("Value is not equal to %d:%d:err(%v)\n", int(KB), n, err)
	}

	n, err = HeadParseSize("1MB")
	if n != MB || err != nil {
		t.Errorf("Value is not equal to %d:%d:err(%v)\n", int(MB), n, err)
	}

	n, err = HeadParseSize("1K")
	if n != K || err != nil {
		t.Errorf("Value is not equal to %d:%d:err(%v)\n", int(K), n, err)
	}

	n, err = HeadParseSize("123MB")
	if n != 123*1000*1000 || err != nil {
		t.Errorf("Value is not equal to %d:%d:err(%v)\n", int(123*MB), n, err)
	}
}
