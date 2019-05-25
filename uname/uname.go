package uname

import (
	"bytes"
	_ "fmt"
	"github.com/guonaihong/flag"
	"golang.org/x/sys/unix"
	"io"
	"os"
	"runtime"
)

type Uname struct {
	KernelName       bool
	NodeName         bool
	KernelRelease    bool
	KernelVersion    bool
	Machine          bool
	Processor        bool
	HardwarePlatform bool
	OperatingSystem  bool
	count            int
}

type utsname struct {
	Sysname         []byte
	Nodename        []byte
	Release         []byte
	Version         []byte
	Machine         []byte
	OperatingSystem []byte
}

func New(argv []string) (*Uname, []string) {
	u := &Uname{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	all := command.Opt("a, all", "print all information, in the following order,\n"+
		"except omit -p and -i if unknown:").
		Flags(flag.PosixShort).NewBool(false)

	command.Opt("s, kernel-name", "print the kernel name").
		Flags(flag.PosixShort).Var(&u.KernelName)

	command.Opt("n, nodename", "print the network node hostname").
		Flags(flag.PosixShort).Var(&u.NodeName)

	command.Opt("r, kernel-release", "print the kernel release").
		Flags(flag.PosixShort).Var(&u.KernelRelease)

	command.Opt("v, kernel-version", "print the kernel version").
		Flags(flag.PosixShort).Var(&u.KernelVersion)

	command.Opt("m, machine", "print the machine hardware name").
		Flags(flag.PosixShort).Var(&u.Machine)

	command.Opt("p, processor", "print the processor type (non-portable)").
		Flags(flag.PosixShort).Var(&u.Processor)

	command.Opt("i, hardware-platform", "print the hardware platform (non-portable)").
		Flags(flag.PosixShort).Var(&u.HardwarePlatform)

	command.Opt("o, operating-system", "print the operating system").
		Flags(flag.PosixShort).Var(&u.OperatingSystem)

	command.Parse(argv[1:])

	if all != nil && *all {
		u.KernelName = true
		u.NodeName = true
		u.KernelRelease = true
		u.KernelVersion = true
		u.Machine = true
		u.Processor = true
		u.HardwarePlatform = true
		u.OperatingSystem = true
	}

	return u, command.Args()

}

func truncated0Bytes(b []byte) []byte {
	pos := bytes.IndexByte(b, 0)
	if pos == -1 {
		return b
	}
	return b[:pos]
}

func (u *Uname) shouldBindUname(name *utsname) {
	buf := unix.Utsname{}

	unix.Uname(&buf)

try:
	if u.KernelName {
		name.Sysname = truncated0Bytes(buf.Sysname[:])
		u.count++
	}

	if u.NodeName {
		name.Nodename = truncated0Bytes(buf.Nodename[:])
		u.count++
	}

	if u.KernelRelease {
		name.Release = truncated0Bytes(buf.Release[:])
		u.count++
	}

	if u.KernelVersion {
		name.Version = truncated0Bytes(buf.Version[:])
		u.count++
	}

	name.Machine = truncated0Bytes(buf.Machine[:])

	if u.Machine {
		u.count++
	}

	if u.HardwarePlatform {
		u.count++
	}

	if u.Processor {
		u.count++
	}

	if u.OperatingSystem {
		name.OperatingSystem = getOsName()
		u.count++
	}

	if u.count == 0 {
		u.KernelName = true
		goto try
	}
}

func (u *Uname) writeSpace(w io.Writer) {
	if u.count--; u.count > 0 {
		w.Write([]byte{' '})
	}
}

func (u *Uname) Uname(w io.Writer) {
	var name utsname

	u.shouldBindUname(&name)

	needLine := u.count > 0

	if u.KernelName {
		w.Write(name.Sysname)
		u.writeSpace(w)
	}

	if u.NodeName {
		w.Write(name.Nodename)
		u.writeSpace(w)
	}

	if u.KernelRelease {
		w.Write(name.Release)
		u.writeSpace(w)
	}

	if u.KernelVersion {
		w.Write(name.Version)
		u.writeSpace(w)
	}

	if u.Machine {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.Processor {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.HardwarePlatform {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.OperatingSystem {
		w.Write(getOsName())
		u.writeSpace(w)
	}

	if needLine {
		w.Write([]byte{'\n'})
	}
}

func getOsName() []byte {
	osName := runtime.GOOS
	switch osName {
	case "linux":
		osName = "GNU/Linux"
	case "windows":
		osName = "Windows NT"
	case "freebsd":
		osName = "FreeBSD"
	case "openbsd":
		osName = "OpenBSD"
	case "darwin":
		osName = "Darwin"
	case "android":
		osName = "Android"
	}

	return []byte(osName)
}

func Main(argv []string) {
	u, _ := New(argv)

	u.Uname(os.Stdout)
}
