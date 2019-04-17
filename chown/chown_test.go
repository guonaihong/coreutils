package chown

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"os"
	"os/user"
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

func testChown(name string, fileName string, needErr string, rootRun bool, t *testing.T) {
	u, err := user.Current()
	if err != nil {
		return
	}

	if u.Username == "root" && !rootRun {
		return
	}

	c := Chown{}
	err = c.Chown(name, fileName, os.Stdout)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func testChownVerbose(name string, out string, t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}

	if u.Username != "root" {
		t.Errorf("need root user\n")
		return
	}

	c := Chown{}

	var w bytes.Buffer

	os.Chown("test.dat", 2, 2)

	c.Verbose = utils.Bool(true)
	err = c.Chown(name, "test.dat", &w)
	if w.String() != out || w.String() == "" {
		t.Errorf("need(%s), actual(%s), rv(%v)\n", out, w.String(), err)
	}

}

// need root user to run
func TestChownVerbose(t *testing.T) {
	testChownVerbose(":", "ownership of 'test.dat' retained\n", t)
	testChownVerbose("bin:", "ownership of 'test.dat' retained as bin:bin\n", t)
	testChownVerbose(":bin", "ownership of 'test.dat' retained as bin:bin\n", t)
	testChownVerbose("root:", "changed ownership of 'test.dat' from bin:bin to root:root\n", t)
	testChownVerbose(":root", "changed ownership of 'test.dat' from bin:bin to :root\n", t)
	testChownVerbose("root", "changed ownership of 'test.dat' from bin to root\n", t)
}

func TestChown(t *testing.T) {
	testChown("root", "chown_test.go", "chown: changing ownership of 'chown_test.go': Operation not permitted", false, t)
	testChown(":", "yy", "chown: cannot access 'yy': No such file or directory", true, t)
}
