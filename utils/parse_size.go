package utils

import (
	"errors"
	"fmt"
)

func IsDecimal(b byte) bool {
	return b >= '0' && b <= '9'
}

func IsDecimalStr(s string, max int) (i int, have bool) {
	return IsNeedStr(s, max, IsDecimal)
}

func IsNeedStr(s string, max int, is func(b byte) bool) (i int, have bool) {
	for i = 0; i < len(s); i++ {
		if i >= max {
			return i, have
		}

		if !is(s[i]) {
			return i, have
		}

		have = true
	}

	return i, have
}

type Size int

const (
	Byte Size = 1
	B         = 512
	K         = 1024 * Byte
	M         = 1024 * K
	G         = 1024 * M
	T         = 1024 * G
	P         = 1024 * T
	E         = 1024 * P
	//Z         = 1024 * E
	//Y         = 1024 * Z
)

const (
	KB = 1000
	MB = 1000 * KB
	GB = 1000 * MB
	TB = 1000 * GB
	PB = 1000 * TB
	EB = 1000 * PB
	//ZB = 1000 * EB
	//YB = 1000 * ZB
)

func (s Size) IntPtr() *int {
	n := int(s)
	return &n
}

func HeadParseSize(s string) (Size, error) {
	i, have := IsDecimalStr(s, len(s))
	if !have {
		return Size(0), fmt.Errorf("invalid number of bytes: '%s'", s)
	}

	n := 0
	fmt.Sscanf(s, "%d", &n)

	if i >= len(s) {
		return Size(n), nil
	}

	suffix := Size(0)
	switch s[i:] {
	case "B":
		suffix = B
	case "kB":
		suffix = KB
	case "MB":
		suffix = MB
	case "GB":
		suffix = GB
	case "TB":
		suffix = TB
	case "PB":
		suffix = PB
	case "EB":
		suffix = EB
		/*
			case "ZB":
				suffix = ZB
			case "YB":
				suffix = YB
		*/
	case "K":
		suffix = K
	case "M":
		suffix = M
	case "G":
		suffix = G
	case "T":
		suffix = T
	case "P":
		suffix = P
	case "E":
		suffix = E
		/*
			case "Z":
				suffix = Z
			case "Y":
				suffix = Y
		*/
	default:
		return 0, errors.New("Unsupported suffix")
	}

	return Size(n * int(suffix)), nil
}
