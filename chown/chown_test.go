package chown

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"os"
	"testing"
)

func testGetGroupUserError(name string, needErr string, t *testing.T) {
	_, err := getGroupUser(name)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func TestGetGroupUserError(t *testing.T) {
	testGetGroupUserError("wwwwwww:root", "chown: invalid user: 'wwwwwww:root'", t)
	testGetGroupUserError("root:wwwwwww", "chown: invalid group: 'root:wwwwwww'", t)
}

func testChown(name string, fileName string, needErr string, t *testing.T) {
	c := Chown{}
	err := c.Chown(name, fileName, os.Stdout)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func testChownChanges(name string, out string, t *testing.T) {
	c := Chown{}

	var w bytes.Buffer

	os.Chown("tst.dat", 2, 2)

	c.Verbose = utils.Bool(true)
	err := c.Chown(name, "tst.dat", &w)
	if w.String() != out || w.String() == "" {
		t.Errorf("need(%s), actual(%s), rv(%v)\n", out, w.String(), err)
	}

}

// need root user to run
func TestChownChanges(t *testing.T) {
	testChownChanges(":", "ownership of 'tst.dat' retained\n", t)
	testChownChanges("bin:", "", t)
}

func TestChown(t *testing.T) {
	testChown("root", "chown_test.go", "chown: changing ownership of 'chown_test.go': Operation not permitted", t)
	testChown(":", "yy", "chown: cannot access 'yy': No such file or directory", t)
}
