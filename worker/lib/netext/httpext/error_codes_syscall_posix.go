//go:build !windows
// +build !windows

package httpext

import (
	"net"
	"os"
)

func getOSSyscallErrorCode(e *net.OpError, se *os.SyscallError) (errCode, string) ***REMOVED***
	return 0, ""
***REMOVED***
