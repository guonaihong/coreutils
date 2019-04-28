package uname

import (
	"bytes"
	_ "fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

type Uname struct {
	KernelName       *bool
	NodeName         *bool
	KernelRelease    *bool
	KernelVersion    *bool
	Machine          *bool
	Processor        *bool
	HardwarePlatform *bool
	OperatingSystem  *bool
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

func (u *Uname) isKernelName() bool {
	return u.KernelName != nil && *u.KernelName
}

func (u *Uname) isNodeName() bool {
	return u.NodeName != nil && *u.NodeName
}

func (u *Uname) isKernelRelease() bool {
	return u.KernelRelease != nil && *u.KernelRelease
}

func (u *Uname) isKernelVersion() bool {
	return u.KernelVersion != nil && *u.KernelVersion
}

func (u *Uname) isMachine() bool {
	return u.Machine != nil && *u.Machine
}

func (u *Uname) isProcessor() bool {
	return u.Processor != nil && *u.Processor
}

func (u *Uname) isHardwarePlatform() bool {
	return u.HardwarePlatform != nil && *u.HardwarePlatform
}

func (u *Uname) isOperatingSystem() bool {
	return u.OperatingSystem != nil && *u.OperatingSystem
}

func New(argv []string) (*Uname, []string) {
	u := &Uname{}

	command := flag.NewFlagSet(argv[0], flag.ExitOnError)

	all := command.Opt("a, all", "print all information, in the following order,\n"+
		"except omit -p and -i if unknown:").
		Flags(flag.PosixShort).NewBool(false)

	u.KernelName = command.Opt("s, kernel-name", "print the kernel name").
		Flags(flag.PosixShort).NewBool(false)

	u.NodeName = command.Opt("n, nodename", "print the network node hostname").
		Flags(flag.PosixShort).NewBool(false)

	u.KernelRelease = command.Opt("r, kernel-release", "print the kernel release").
		Flags(flag.PosixShort).NewBool(false)

	u.KernelVersion = command.Opt("v, kernel-version", "print the kernel version").
		Flags(flag.PosixShort).NewBool(false)

	u.Machine = command.Opt("m, machine", "print the machine hardware name").
		Flags(flag.PosixShort).NewBool(false)

	u.Processor = command.Opt("p, processor", "print the processor type (non-portable)").
		Flags(flag.PosixShort).NewBool(false)

	u.HardwarePlatform = command.Opt("i, hardware-platform", "print the hardware platform (non-portable)").
		Flags(flag.PosixShort).NewBool(false)

	u.OperatingSystem = command.Opt("o, operating-system", "print the operating system").
		Flags(flag.PosixShort).NewBool(false)

	command.Parse(argv[1:])

	if all != nil && *all {
		u.KernelName = utils.Bool(true)
		u.NodeName = utils.Bool(true)
		u.KernelRelease = utils.Bool(true)
		u.KernelVersion = utils.Bool(true)
		u.Machine = utils.Bool(true)
		u.Processor = utils.Bool(true)
		u.HardwarePlatform = utils.Bool(true)
		u.OperatingSystem = utils.Bool(true)
	}

	return u, command.Args()

}

func (u *Uname) shouldBindUname(name *utsname) {
	buf := syscall.Utsname{}

	syscall.Uname(&buf)

try:
	if u.isKernelName() {
		name.Sysname = truncated0Bytes(buf.Sysname)
		u.count++
	}

	if u.isNodeName() {
		name.Nodename = truncated0Bytes(buf.Nodename)
		u.count++
	}

	if u.isKernelRelease() {
		name.Release = truncated0Bytes(buf.Release)
		u.count++
	}

	if u.isKernelVersion() {
		name.Version = truncated0Bytes(buf.Version)
		u.count++
	}

	name.Machine = truncated0Bytes(buf.Machine)
	if u.isMachine() {
		u.count++
	}

	if u.isHardwarePlatform() {
		u.count++
	}

	if u.isProcessor() {
		u.count++
	}

	if u.isOperatingSystem() {
		name.OperatingSystem = getOsName()
		u.count++
	}

	if u.count == 0 {
		u.KernelName = utils.Bool(true)
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

	if u.isKernelName() {
		w.Write(name.Sysname)
		u.writeSpace(w)
	}

	if u.isNodeName() {
		w.Write(name.Nodename)
		u.writeSpace(w)
	}

	if u.isKernelRelease() {
		w.Write(name.Release)
		u.writeSpace(w)
	}

	if u.isKernelVersion() {
		w.Write(name.Version)
		u.writeSpace(w)
	}

	if u.isMachine() {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.isProcessor() {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.isHardwarePlatform() {
		w.Write(name.Machine)
		u.writeSpace(w)
	}

	if u.isOperatingSystem() {
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

func truncated0Bytes(b [65]int8) []byte {
	var rv []byte

	addr := (*[len(b)]byte)(unsafe.Pointer(&b[0]))
	rv = addr[:]
	pos := bytes.IndexByte(rv, 0)
	if pos != -1 {
		rv = rv[:pos]
	}

	return rv
}

func Main(argv []string) {
	u, _ := New(argv)

	u.Uname(os.Stdout)
}
