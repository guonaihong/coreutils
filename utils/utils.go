package utils

import (
	"fmt"
	"os"
)

func Die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func OpenInputFd(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdin, nil
	}

	return os.Open(fileName)
}

func OpenOutputFd(fileName string) (*os.File, error) {
	if fileName == "-" {
		return os.Stdout, nil
	}

	return os.Create(fileName)
}

func CloseInputFd(fd *os.File) {
	if fd != os.Stdin {
		fd.Close()
	}
}

func CloseOutputFd(fd *os.File) {
	if fd != os.Stdout {
		fd.Close()
	}
}
