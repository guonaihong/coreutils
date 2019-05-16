package hashcore

import (
	"github.com/guonaihong/coreutils/utils"
	"strings"
	"testing"
)

const errMsgFileDoesNotExist = `md5sum: not.sum: No such file or directory
md5sum: kao: No such file or directory
kao: FAILED open or read
md5sum: WARNING: 1 listed file could not be read
`
const okMsg = `md5sum.sum: OK
`

const errMsgIgnoreMissing = `md5sum: not.sum: No such file or directory
md5sum: md5sum.sum: no file was verified
`

func testHashCheck(hashType Type, need string, fileNames []string, quiet bool, ignoreMissing bool, t *testing.T) {
	var s strings.Builder
	var formatFail int
	h := HashCore{}
	h.Quiet = utils.Bool(quiet)
	h.IgnoreMissing = utils.Bool(ignoreMissing)
	for _, f := range fileNames {
		err := h.CheckHash(hashType, &formatFail, f, &s)
		if err != nil {
			s.WriteString(err.Error())
		}
	}

	if need != s.String() {
		t.Errorf("need(%s) actual(%s)\n", need, s.String())
	}
}

func TestHashCheckFileDoesNotExist(t *testing.T) {
	testHashCheck(Md5, errMsgFileDoesNotExist, []string{"not.sum", "md5sum.sum"}, false, false, t)
	testHashCheck(Md5, errMsgFileDoesNotExist, []string{"not.sum", "md5sum.sum"}, true, false, t)
	testHashCheck(Md5, errMsgIgnoreMissing, []string{"not.sum", "md5sum.sum"}, false, true, t)
	testHashCheck(Md5, errMsgIgnoreMissing, []string{"not.sum", "md5sum.sum"}, true, true, t)
}

func TestHashFileExist2(t *testing.T) {
	testHashCheck(Md5, okMsg, []string{"ok.md5.sum"}, false, false, t)
	testHashCheck(Md5, "", []string{"ok.md5.sum"}, true, false, t)
}

func TestHashFileExist1(t *testing.T) {
	testHashCheck(Md5, okMsg, []string{"ok.binary.md5.sum"}, false, false, t)
	testHashCheck(Md5, "", []string{"ok.md5.sum"}, true, false, t)
}
