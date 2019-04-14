package chown

import (
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
	err := c.Chown(name, fileName)
	if err.Error() != needErr {
		t.Errorf("need error(%s), actual error(%s)\n", needErr, err.Error())
	}
}

func TestChown(t *testing.T) {
	testChown("root", "chown_test.go", "chown: changing ownership of 'chown_test.go': Operation not permitted", t)
	testChown(":", "yy", "chown: cannot access 'yy': No such file or directory", t)
}
