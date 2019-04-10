package touch

import (
	"github.com/guonaihong/coreutils/utils"
	"os"
	"testing"
	"time"
)

const fileName1 = "touch-test-file1"
const fileName2 = "touch-test-file2"

func testParseTime(tm string, needTm time.Time, t *testing.T) {
	touch := Touch{}
	//touch.debug = utils.Bool(true)
	t1, err := touch.parseTime(tm)
	if err != nil {
		t.Errorf("touch -t fail:%s\n", err)
	}

	if !t1.Equal(needTm) {
		t.Errorf("needTm(%v), tm(%v)\n", needTm, t1)
	}
}

func TestParseTime(t *testing.T) {
	now := time.Now()
	//15
	testParseTime("201212101830.55", time.Date(2012, 12, 10, 18, 30, 55, 0, time.UTC), t)

	//12
	testParseTime("201812101730", time.Date(2018, 12, 10, 17, 30, 00, 0, time.UTC), t)

	//13
	testParseTime("1712101730.44", time.Date(2017, 12, 10, 17, 30, 44, 0, time.UTC), t)
	testParseTime("6712101730.44", time.Date(2067, 12, 10, 17, 30, 44, 0, time.UTC), t)

	//10
	testParseTime("1712101730", time.Date(2017, 12, 10, 17, 30, 00, 0, time.UTC), t)

	//11
	testParseTime("12101730.33", time.Date(now.Year(), 12, 10, 17, 30, 33, 0, time.UTC), t)

	//8
	testParseTime("12101730", time.Date(now.Year(), 12, 10, 17, 30, 0, 0, time.UTC), t)
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
	testParseDate("1 May 2005 10:22", time.Date(2005, 5, 1, 10, 22, 0, 0, time.UTC), t)
	testParseDate("14 May", time.Date(now.Year(), 5, 14, 0, 0, 0, 0, time.UTC), t)
	testParseDate("14:24", time.Date(now.Year(), now.Month(), now.Day(), 14, 24, 0, 0, time.UTC), t)
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
