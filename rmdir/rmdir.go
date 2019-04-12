package rmdir

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
)

type Rmdir struct {
	Parent               *bool
	IgnoreFailOnNonEmpty *bool
}

func New(argv []string) (*Rmdir, []string) {

	rmdir := &Rmdir{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	rmdir.IgnoreFailOnNonEmpty = command.Opt("ignore-fail-on-non-empty",
		"ignore each failure that is solely because a directory\n"+
			"is non-empty").
		NewBool(false)

	rmdir.Parent = command.Opt("p, parent",
		"remove DIRECTORY and its ancestors; e.g., 'rmdir -p a/b/c' is\n"+
			"similar to 'rmdir a/b/c a/b a'").
		NewBool(false)

	command.Parse(argv[1:])

	return rmdir, command.Args()
}

func (r *Rmdir) Rmdir(name string) error {

	if r.Parent != nil && *r.Parent {
		return os.RemoveAll(name)
	}

	return os.Remove(name)
}

func (r *Rmdir) isDir(name string) (bool, error) {
	fi, err := os.Stat(name)

	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func (r *Rmdir) IsIgnoreFailOnNonEmpty() bool {
	return r.IgnoreFailOnNonEmpty != nil &&
		*r.IgnoreFailOnNonEmpty
}

func Main(argv []string) {

	rmdir, args := New(argv)

	for _, v := range args {

		isDir, err := rmdir.isDir(v)
		if err != nil {
			fmt.Printf("rmdir: %s\n", err)
			continue
		}

		if isDir == false {
			fmt.Printf("rmdir: failed to remove '%s': Not a directory\n", v)
			continue
		}

		err = rmdir.Rmdir(v)

		if err != nil {
			if rmdir.IsIgnoreFailOnNonEmpty() {
				continue
			}

			fmt.Printf("rmdir: %s\n", err)
		}
	}
}
