package chgrp

import (
	"bytes"
	"github.com/guonaihong/coreutils/utils"
	"os"
	"os/user"
	"testing"
)

const (
	binGid  = 2
	rootGid = 0
)

func testGetGroupUserError(name string, needErr string, t *testing.T) {
	_, err := getGroupFromName(name)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func TestGetGroupUserError(t *testing.T) {
	testGetGroupUserError("wwwwwww", "chgrp: invalid group: 'wwwwwww'", t)
}

func testChgrp(name string, fileName string, needErr string, rootRun bool, t *testing.T) {
	u, err := user.Current()
	if err != nil {
		return
	}

	if u.Username == "root" && !rootRun {
		return
	}

	user := User{}
	user.Init(os.Stdout)
	c := Chgrp{}
	err = c.Chgrp(name, fileName, &user)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func testChgrpVerbose(name string, out string, gid int, t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}

	if u.Username != "root" {
		t.Errorf("need root user\n")
		return
	}

	c := Chgrp{}

	var w bytes.Buffer

	user := User{}
	user.Init(&w)
	os.Chown("test.dat", 2, 2)

	c.Verbose = utils.Bool(true)
	err = c.Chgrp(name, "test.dat", &user)
	if w.String() != out || w.String() == "" {
		t.Errorf("need(%s), actual(%s), rv(%v), name(%s)\n",
			out, w.String(), err, name)
	}

	if user.Gid != gid {
		t.Errorf("name (%s) need gid(%d), actual gid(%d)\n",
			name, gid, user.Gid)
	}
}

// need root user to run
func TestChgrpVerbose(t *testing.T) {
	testChgrpVerbose("bin", "group of 'test.dat' retained as bin\n", -1, t)

	testChgrpVerbose("root", "changed ownership of 'test.dat' from bin to root\n", -1, t)
}

func testChgrpChanges(name string, out string, gid int, t *testing.T) {
	u, err := user.Current()
	if err != nil {
		t.Errorf("%s\n", err)
		return
	}

	if u.Username != "root" {
		t.Errorf("need root user\n")
		return
	}

	c := Chgrp{}

	var w bytes.Buffer

	user := User{}
	user.Init(&w)
	os.Chown("test.dat", 2, 2)

	c.Changes = utils.Bool(true)
	err = c.Chgrp(name, "test.dat", &user)
	if w.String() != out {
		t.Errorf("need(%s), actual(%s), rv(%v), name(%s)\n",
			out, w.String(), err, name)
	}

	if user.Gid != gid {
		t.Errorf("name (%s) need gid(%d), actual gid(%d)\n",
			name, gid, user.Gid)
	}
}

func TestChgrpChanges(t *testing.T) {
	testChgrpChanges("bin", "", -1, t)

	testChgrpChanges("root", "changed ownership of 'test.dat' from bin to root\n", -1, t)
}

func TestChgrp(t *testing.T) {
	testChgrp("root", "chgrp_test.go", "chgrp: changing ownership of 'chgrp_test.go': Operation not permitted", false, t)
}
