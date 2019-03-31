package touch

import (
	"github.com/guonaihong/coreutils/utils"
	"os"
	"testing"
)

const fileName1 = "touch-test-file1"
const fileName2 = "touch-test-file2"

//touch -c file
func TestCreateFileNoCreate(t *testing.T) {
	touch := Touch{}
	touch.NoCreate = utils.Bool(true)

	touch.Touch(fileName1)
	if !isNotExist(fileName1) {
		t.Errorf("touch -c found file %s\n", fileName1)
	}
}

//touch file
func TestCreateFile(t *testing.T) {
	touch := Touch{}
	touch.Touch(fileName2)
	defer os.Remove(fileName2)

	if isNotExist(fileName2) {
		t.Errorf("not found %s\n", fileName2)
	}
}
