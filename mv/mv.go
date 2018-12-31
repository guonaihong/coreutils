package main

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
"strings"
)

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func deleteTrailingSlashes(s string) string {
	return strings.TrimRight(s, "/\\")
}

func readYes() bool {
    b := make([]byte, 512)
	n, err := os.Stdin.Read(b)
	if err != nil {
		return false
	}

	if string(b[:n]) != "y\n" {
		return false
	}
	return true
}

func main() {

	interactive := flag.Bool("i, interactive", false, "prompt before overwrite")
	force := flag.Bool("f, force", false, "do not prompt before overwriting")
	//backupSuffix := flag.String("backup", "", "make a backup of each existing destination file")
	backup := flag.Bool("b", false, "like --backup but does not accept an argument")
	noClobber := flag.Bool("n, no-clobber", false, "do not overwrite an existing file")
	stripTrailingSlashes := flag.Bool("strip-trailing-slashes", false, "remove any trailing slashes from each SOURCE argument")
	suffix := flag.String("S, suffix", "~", "override the usual backup suffix")
	targetDirectory := flag.String("t, target-directory", "", "move all SOURCE arguments into DIRECTORY")
	noTargetDirectory := flag.Bool("T, no-target-directory", false, "treat DEST as a normal file")
	//update := flag.Bool("u, update", false, "move only when the SOURCE file is newer than the destination file or when the destination file is missing")
	//verbose := flag.Bool("v, verbose", false, "explain what is being done")
	//context := flag.Bool("Z, context", false, "set SELinux security context of destination file to default type")
	//version := flag.Bool("version", false, "output version information and exit")

	flag.Parse()

	args := flag.Args()

	if len(args) <= 1 {
		flag.Usage()
		os.Exit(64)
	}

	source := args[:len(args)-1]
	target := args[len(args)-1]

	if len(*targetDirectory) > 0 {
		source = args
		target = *targetDirectory
	}

	if *noTargetDirectory {
		fmt.Println("icannot combine --target-directory (-t)",
			"and --no-target-directory (-T)")
		os.Exit(0)
	}

	for _, s := range source {

		if *stripTrailingSlashes {
			s = deleteTrailingSlashes(s)
		}
		// if target exists and is a directory, s is moved into it.
		newTarget := target + "/" + s
		if fi, err := os.Lstat(target); !os.IsNotExist(err) {
			mode := fi.Mode()
			if mode.IsRegular() {

				if *noClobber {
					continue
				}

				newTarget = target
			}
		} else {
			// if target does not exist, s is renamed target
			newTarget = target

            mode := fi.Mode()
			perm := mode.Perm()
			if !*force && int(perm)&os.O_WRONLY == 0 {
				fmt.Printf("mv: replace '%s', overriding mode %x (%v)?", newTarget, perm, perm)
				readYes()
			}
		}

		if *force {
			goto rename
		}

		if *interactive {
			if *force {
				goto rename
			}

			if exist(newTarget) {
				fmt.Printf("overwrite %s ? (y/n [n])", newTarget)
				if !readYes() {
					return
				}
			}

		}

		if *backup {
			if exist(newTarget) {
				err := os.Rename(newTarget, newTarget+*suffix)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}

		}

	rename:
		err := os.Rename(s, target)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
	}
}
