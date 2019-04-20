package chown

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"os"
	"os/user"
	"testing"
)

const (
	binUid  = 2
	binGid  = 2
	rootUid = 0
	rootGid = 0
)

func testGetGroupUserError(name string, needErr string, t *testing.T) {
	_, err := getUserGroupFromName(name)
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

	user := User{}
	user.Init(os.Stdout)
	c := Chown{}
	err = c.Chown(name, fileName, &user)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func testChownVerbose(name string, out string, uid int, gid int, t *testing.T) {
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

	user := User{}
	user.Init(&w)
	os.Chown("test.dat", 2, 2)

	c.Verbose = utils.Bool(true)
	err = c.Chown(name, "test.dat", &user)
	if w.String() != out || w.String() == "" {
		t.Errorf("need(%s), actual(%s), rv(%v), name(%s)\n",
			out, w.String(), err, name)
	}

	if user.Uid != uid || user.Gid != gid {
		t.Errorf("name (%s) need uid(%d) gid(%d), actual uid(%d) gid(%d)\n",
			name, uid, gid, user.Uid, user.Gid)
	}
}

// need root user to run
func TestChownVerbose(t *testing.T) {
	testChownVerbose(":", "ownership of 'test.dat' retained\n", -1, -1, t)
	testChownVerbose("bin", "ownership of 'test.dat' retained as bin\n", binUid, -1, t)
	testChownVerbose("bin:", "ownership of 'test.dat' retained as bin:bin\n", binUid, binGid, t)
	testChownVerbose(":bin", "ownership of 'test.dat' retained as :bin\n", -1, binUid, t)

	testChownVerbose("root:", "changed ownership of 'test.dat' from bin:bin to root:root\n", rootUid, rootGid, t)
	testChownVerbose(":root", "changed ownership of 'test.dat' from bin:bin to :root\n", -1, rootGid, t)

	testChownVerbose("root", "changed ownership of 'test.dat' from bin to root\n", rootUid, -1, t)
}

func testChownChanges(name string, out string, uid int, gid int, t *testing.T) {
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

	user := User{}
	user.Init(&w)
	os.Chown("test.dat", 2, 2)

	c.Changes = utils.Bool(true)
	err = c.Chown(name, "test.dat", &user)
	if w.String() != out {
		t.Errorf("need(%s), actual(%s), rv(%v), name(%s)\n",
			out, w.String(), err, name)
	}

	if user.Uid != uid || user.Gid != gid {
		t.Errorf("name (%s) need uid(%d) gid(%d), actual uid(%d) gid(%d)\n",
			name, uid, gid, user.Uid, user.Gid)
	}
}

func TestChownChanges(t *testing.T) {
	testChownChanges(":", "", -1, -1, t)
	testChownChanges("bin", "", binUid, -1, t)
	testChownChanges("bin:", "", binUid, binGid, t)
	testChownChanges(":bin", "", -1, binUid, t)

	testChownChanges("root:", "changed ownership of 'test.dat' from bin:bin to root:root\n", rootUid, rootGid, t)
	testChownChanges(":root", "changed ownership of 'test.dat' from bin:bin to :root\n", -1, rootGid, t)

	testChownChanges("root", "changed ownership of 'test.dat' from bin to root\n", rootUid, -1, t)
}

func TestChown(t *testing.T) {
	testChown("root", "chown_test.go", "chown: changing ownership of 'chown_test.go': Operation not permitted", false, t)
	testChown(":", "yy", "chown: cannot access 'yy': No such file or directory", true, t)
}
