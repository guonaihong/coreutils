package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"
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

func readFile(fileName string, body []byte) (n int, err error) {

	var fd *os.File

	fd, err = os.Open(fileName)
	if err != nil {
		return 0, err
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	n1, err := fd.Seek(int64(0), 2)
	if err != nil {
		return 0, err
	}

	if n1 <= 0 {
		n1 = math.MaxInt32
	}

	fd.Seek(r.Int63n(n1), os.SEEK_SET)
	return fd.Read(body)
}

func GetRandSource(fileName string) (*rand.Rand, error) {

	seed := int64(0)
	buf := make([]byte, 8)

	_, err := readFile(fileName, buf)
	if err != nil {
		return nil, err
	}

	read := bytes.NewReader(buf)

	err = binary.Read(read, binary.LittleEndian, &seed)
	if err != nil {
		return nil, err
	}

	return rand.New(rand.NewSource(seed)), nil
}
