//go:build linux || darwin || dragonfly || freebsd || netbsd || openbsd || solaris || illumos
// +build linux darwin dragonfly freebsd netbsd openbsd solaris illumos

package pool

import (
	"errors"
	"io"
	"net"
	"syscall"
	"time"
)

var errUnexpectedRead = errors.New("unexpected read from socket")

func connCheck(conn net.Conn) error ***REMOVED***
	// Reset previous timeout.
	_ = conn.SetDeadline(time.Time***REMOVED******REMOVED***)

	sysConn, ok := conn.(syscall.Conn)
	if !ok ***REMOVED***
		return nil
	***REMOVED***
	rawConn, err := sysConn.SyscallConn()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var sysErr error

	if err := rawConn.Read(func(fd uintptr) bool ***REMOVED***
		var buf [1]byte
		n, err := syscall.Read(int(fd), buf[:])
		switch ***REMOVED***
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = errUnexpectedRead
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		***REMOVED***
		return true
	***REMOVED***); err != nil ***REMOVED***
		return err
	***REMOVED***

	return sysErr
***REMOVED***
