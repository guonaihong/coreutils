package main

import (
	"fmt"
	"github.com/guonaihong/flag"
	"os"
)

func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func main() {

	interactive := flag.Bool("i, interactive", false, "prompt before overwrite")
	force := flag.Bool("f, force", false, "do not prompt before overwriting")
	backupSuffix := flag.String("backup", "~", "make a backup of each existing destination file")
	backup := flag.Bool("b", false, "like --backup but does not accept an argument")
	noClobber := flag.Bool("n, no-clobber", false, "do not overwrite an existing file")
	stripTrailingSlashes := flag.Bool("strip-trailing-slashes", "", "remove any trailing slashes from each SOURCE argument")
	suffx := flag.String("S, suffix", "", "override the usual backup suffix")
	targetDirectory := flag.String("t, target-directory", "", "move all SOURCE arguments into DIRECTORY")
	noTargetDirectory := flag.String("T, no-target-directory", false, "treat DEST as a normal file")
	update := flag.Bool("u, update", false, "move only when the SOURCE file is newer than the destination file or when the destination file is missing")
	verbose := flag.String("v, verbose", false, "explain what is being done")
	context := flag.Bool("Z, context", false, "set SELinux security context of destination file to default type")
	version := flag.Bool("version", false, "output version information and exit")

	flag.Parse()

	args := flag.Args()

	if len(args) <= 1 {
		flag.Usage()
		os.Exit(64)
	}

	suffix := "~"
	if len(*backupSuffix) > 0 {
		suffix = *backupSuffix
		*backup = true
	}

	source := args[:len(args)-1]
	target := args[len(args)-1]

	b := make([]byte, 512)
	for _, s := range source {
		newTarget := target + "/" + s
		if *interactive {
			if *force {
				goto rename
			}

			if exist(newTarget) {
				fmt.Printf("overwrite %s ? (y/n [n])", newTarget)
				n, err := os.Stdin.Read(b)
				if err != nil {
					return
				}

				if string(b[:n]) != "y\n" {
					return
				}
			}

		}

		if *backup {
			if exist(newTarget) {
				err := os.Rename(newTarget, newTarget+suffix)
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
