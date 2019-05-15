package hashcore

import (
	"strings"
	"testing"
)

const errMsg = `md5sum: md52.sumx: No such file or directory
md5sum: kao: No such file or directory
kao: FAILED open or read
md5sum: WARNING: 1 listed file could not be read
`

func testHashCheck(need string, t *testing.T) {
	var s strings.Builder
	h := HashCore{}
	fileNames := []string{"md52.sumx", "md5sum.sum"}

	for _, f := range fileNames {
		h.CheckHash(Md5, f, &s)
	}

	if need != s.String() {
		t.Errorf("need(%s) actual(%s)\n", need, s.String())
	}
}

func TestHashCheck(t *testing.T) {
	testHashCheck(errMsg, t)
}
