package tee

import (
	"os/signal"
	"syscall"
)

func ignore() {
	signal.Ignore(syscall.SIGABRT)
	signal.Ignore(syscall.SIGALRM)
	signal.Ignore(syscall.SIGBUS)
	signal.Ignore(syscall.SIGCHLD)

	signal.Ignore(syscall.SIGCLD)
	signal.Ignore(syscall.SIGPOLL)
	signal.Ignore(syscall.SIGPWR)

	signal.Ignore(syscall.SIGCONT)
	signal.Ignore(syscall.SIGFPE)
	signal.Ignore(syscall.SIGHUP)
	signal.Ignore(syscall.SIGILL)
	signal.Ignore(syscall.SIGINT)
	signal.Ignore(syscall.SIGIO)
	signal.Ignore(syscall.SIGIOT)
	signal.Ignore(syscall.SIGKILL)
	signal.Ignore(syscall.SIGPIPE)
	signal.Ignore(syscall.SIGPROF)
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
