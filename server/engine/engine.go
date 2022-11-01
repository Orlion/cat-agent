package engine

import (
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func a() error {
	syscall.ForkLock.RLock()
	listenFd, err := syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, 0)
	if err == nil {
		syscall.CloseOnExec(listenFd)
	}
	syscall.ForkLock.RUnlock()
	if err != nil {
		return os.NewSyscallError("socket", err)
	}
	if err = syscall.SetNonblock(listenFd, true); err != nil {
		return os.NewSyscallError("setnonblock", err)
	}

	err = syscall.Bind(listenFd, nil)
	if err != nil {
		return err
	}

	err = syscall.Listen(listenFd, 0)
	if err != nil {
		return err
	}

	// EpollCreate需要传size参数，但是在linux 2.6.8之后size被忽略，但是需要大于0
	// 所以EpollCreate1与EpollCreate其实没什么区别
	// unix.EPOLL_CLOEXEC 标识当fork出子进程且子进程执行exec之后会自动关闭从父进程继承过来的文件描述符
	epfd, err := unix.EpollCreate1(unix.EPOLL_CLOEXEC)
	if err != nil {
		return os.NewSyscallError("epoll_create1", err)
	}

	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, listenFd, &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLPRI, // EPOLLIN标识文件描述符可读，EPOLLPRI表示文件描述符有紧急的数据可读
		Fd:     int32(listenFd),
	})
	if err != nil {
		return os.NewSyscallError("epoll_ctl add", err)
	}

	// 第三个参数-1表示永久阻塞
	events := make([]unix.EpollEvent, 128)
	n, err := unix.EpollWait(epfd, events, -1)
	if err != nil {
		return os.NewSyscallError("epoll_wait", err)
	}

	return nil
}
