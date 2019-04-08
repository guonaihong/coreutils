package touch

import (
	"github.com/guonaihong/coreutils/utils"
	"os"
	"testing"
	"time"
)

const fileName1 = "touch-test-file1"
const fileName2 = "touch-test-file2"

func testParseTime(t *testing.T) {
}

func TestParseTime(t *testing.T) {
}

func testParseDate(date string, needDate time.Time, t *testing.T) {
	touch := Touch{}
	t1, err := touch.parseDate(date)
	if err != nil {
		t.Errorf("touch -d fail:%s\n", err)
	}

	if !t1.Equal(needDate) {
		t.Errorf("needDate(%v), date(%v)\n", needDate, t1)
	}
}

func TestParseDate(t *testing.T) {
	now := time.Now()
	testParseDate("1 May 2005 10:22", time.Date(2005, 5, 1, 10, 22, 0, 0, time.Local))
	testParseDate("14 May", time.Date(now.Year(), 5, 14, 0, 0, 0, 0, time.Local))
	testParseDate("14:24", time.Date(now.Year(), now.Month(), now.Day(), 14, 24, 0, 0, time.Local))
}

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
