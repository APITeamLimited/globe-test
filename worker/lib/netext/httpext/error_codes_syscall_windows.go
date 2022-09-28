package httpext

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func getOSSyscallErrorCode(e *net.OpError, se *os.SyscallError) (errCode, string) ***REMOVED***
	switch se.Unwrap() ***REMOVED***
	case syscall.WSAECONNRESET:
		return tcpResetByPeerErrorCode, fmt.Sprintf(tcpResetByPeerErrorCodeMsg, e.Op)
	***REMOVED***
	return 0, ""
***REMOVED***
