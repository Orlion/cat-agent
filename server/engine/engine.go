package engine

import (
	"os"
	"syscall"
)

func a() error {
	syscall.ForkLock.RLock()
	fd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, 0)
	if err == nil {
		syscall.CloseOnExec(fd)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return os.NewSyscallError("socket", err)
	}
	if err = syscall.SetNonblock(fd, true); err != nil {
		return os.NewSyscallError("setnonblock", err)
	}

	err = syscall.Bind(fd, nil)
	if err != nil {
		return err
	}

	err = syscall.Listen(fd, 0)
	if err != nil {
		return err
	}

	return nil
}
