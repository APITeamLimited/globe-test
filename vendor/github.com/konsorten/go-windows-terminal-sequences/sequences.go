// +build windows

package sequences

import (
	"syscall"
	"unsafe"
)

var (
	kernel32Dll    *syscall.LazyDLL  = syscall.NewLazyDLL("Kernel32.dll")
	setConsoleMode *syscall.LazyProc = kernel32Dll.NewProc("SetConsoleMode")
)

func EnableVirtualTerminalProcessing(stream syscall.Handle, enable bool) error ***REMOVED***
	const ENABLE_VIRTUAL_TERMINAL_PROCESSING uint32 = 0x4

	var mode uint32
	err := syscall.GetConsoleMode(syscall.Stdout, &mode)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if enable ***REMOVED***
		mode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	***REMOVED*** else ***REMOVED***
		mode &^= ENABLE_VIRTUAL_TERMINAL_PROCESSING
	***REMOVED***

	ret, _, err := setConsoleMode.Call(uintptr(unsafe.Pointer(stream)), uintptr(mode))
	if ret == 0 ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
