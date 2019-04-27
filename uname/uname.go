package uname

import (
	"bytes"
	"fmt"
	"github.com/guonaihong/coreutils/utils"
	"github.com/guonaihong/flag"
	"io"
	"syscall"
	"unsafe"
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
}

type utsname struct {
	Sysname    []int8
	Nodename   []int8
	Release    []int8
	Version    []int8
	Machine    []int8
	Domainname []int8
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
	all := command.Opt("print all information, in the following order,\n" +
		"except omit -p and -i if unknown:")

	u.KernelName = command.Opt("s, kernel-name", "print the kernel name")
	u.NodeName = command.Opt("n, nodename", "print the network node hostname")
	u.KernelRelease = command.Opt("r, kernel-release", "print the kernel release")
	u.KernelVersion = command.Opt("v, kernel-version", "print the kernel version")
	u.Machine = command.Opt("m, machine", "print the machine hardware name")
	u.Processor = command.Opt("p, processor", "print the processor type (non-portable)")
	u.HardwarePlatform = command.Opt("i, hardware-platform", "print the hardware platform (non-portable)")
	u.OperatingSystem = command.Opt("o, operating-system", "print the operating system")

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
}

func (u *Uname) shouldBindUname(name *utsname) {
	buf := syscall.Utsname{}

	syscall.Uname(&buf)

	if u.isKernelName() {
		name.Sysname = truncated0Bytes(buf.Sysname)
	}

	if u.isNodeName() {
		name.NodeName = truncated0Bytes(buf.NodeName)
	}

	if u.isKernelRelease() {
		name.Release = truncated0Bytes(buf.Release)
	}

	if u.isKernelVersion() {
		name.Version = truncated0Bytes(buf.Version)
	}

	if u.isMachine() {
		name.Machine = truncated0Bytes(buf.Machine)
	}

	name.Domainname = truncated0Bytes(buf.Domainname)
}

func (u *Uname) Uname(w io.Writer) {
	var name utsname

	u.shouldBindUname(&name)
	if u.isKernelName() {
		w.Write(name.Sysname)
	}

	if u.isNodeName() {
		w.Write(name.NodeName)
	}

	if u.isKernelRelease() {
		w.Write(name.Release)
	}

	if u.isKernelVersion() {
		w.Write(name.Version)
	}

	if u.isMachine() {
		w.Write(name.Machine)
	}

	w.Write(name.Domainname)
}

func truncated0Bytes(b [65]uint8) []byte {
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
