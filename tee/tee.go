package tee

import (
	"fmt"
	"github.com/guonaihong/flag"
	"github.com/guonaihong/gutil/file"
	"github.com/guonaihong/log"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func ignore() {
	signal.Ignore(syscall.SIGABRT)
	signal.Ignore(syscall.SIGALRM)
	signal.Ignore(syscall.SIGBUS)
	signal.Ignore(syscall.SIGCHLD)
	signal.Ignore(syscall.SIGCLD)
	signal.Ignore(syscall.SIGCONT)
	signal.Ignore(syscall.SIGFPE)
	signal.Ignore(syscall.SIGHUP)
	signal.Ignore(syscall.SIGILL)
	signal.Ignore(syscall.SIGINT)
	signal.Ignore(syscall.SIGIO)
	signal.Ignore(syscall.SIGIOT)
	signal.Ignore(syscall.SIGKILL)
	signal.Ignore(syscall.SIGPIPE)
	signal.Ignore(syscall.SIGPOLL)
	signal.Ignore(syscall.SIGPROF)
	signal.Ignore(syscall.SIGPWR)
	signal.Ignore(syscall.SIGQUIT)
	signal.Ignore(syscall.SIGSEGV)
	signal.Ignore(syscall.SIGSTKFLT)
	signal.Ignore(syscall.SIGSTOP)
	signal.Ignore(syscall.SIGSYS)
	signal.Ignore(syscall.SIGTERM)
	signal.Ignore(syscall.SIGTRAP)
	//signal.Ignore(syscall.SIGTSTP) //ctrl-z
	signal.Ignore(syscall.SIGTTIN)
	signal.Ignore(syscall.SIGTTOU)
	signal.Ignore(syscall.SIGUNUSED)
	signal.Ignore(syscall.SIGURG)
	signal.Ignore(syscall.SIGUSR1)
	signal.Ignore(syscall.SIGUSR2)
	signal.Ignore(syscall.SIGVTALRM)
	signal.Ignore(syscall.SIGWINCH)
	signal.Ignore(syscall.SIGXCPU)
	signal.Ignore(syscall.SIGXFSZ)
}

func Main(argv []string) {
	command := flag.NewFlagSet(argv[0], flag.ExitOnError)
	append := command.Bool("a, append", false, "append to the given FILEs, do not overwrite")
	ignoreInterrupts := command.Bool("i, ignore-interrupts", false, "ignore interrupt signals") //todo
	gzip := command.Bool("g, gzip", false, "compressed archived log files")
	maxSize := command.String("s, max-size", "0", "current file maximum write size")
	maxArchive := command.Int("A, max-archive", 0, "How many archive files are saved")

	command.Parse(argv[1:])

	var fileArch *log.File
	var buffer [4096]byte
	var fileName string
	var w io.Writer

	if *ignoreInterrupts {
		ignore()
	}

	args := command.Args()
	if len(args) == 0 {
		w = os.Stdout
		goto end
	}

	fileName = args[0]
	if *maxSize != "" {
		size, err := file.ParseSize(*maxSize)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
			return
		}

		var compress log.CompressType
		compress = log.NotCompress
		if *gzip {
			compress = log.Gzip
		}

		fileArch = log.NewFile("", fileName, compress, int(size), *maxArchive)

		w = fileArch
	} else {
		isAppend := 0
		if *append {
			isAppend = os.O_APPEND
		}

		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|isAppend, 0644)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		w = file
	}

end:
	for {
		n, err := os.Stdin.Read(buffer[:])
		if err != nil {
			break
		}
		w.Write(buffer[:n])
		if len(args) != 0 {
			os.Stdout.Write(buffer[:n])
		}
	}
}
