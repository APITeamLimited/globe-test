// +build windows
// +build !go1.4

package mousetrap

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	// defined by the Win32 API
	th32cs_snapprocess uintptr = 0x2
)

var (
	kernel                   = syscall.MustLoadDLL("kernel32.dll")
	CreateToolhelp32Snapshot = kernel.MustFindProc("CreateToolhelp32Snapshot")
	Process32First           = kernel.MustFindProc("Process32FirstW")
	Process32Next            = kernel.MustFindProc("Process32NextW")
)

// ProcessEntry32 structure defined by the Win32 API
type processEntry32 struct ***REMOVED***
	dwSize              uint32
	cntUsage            uint32
	th32ProcessID       uint32
	th32DefaultHeapID   int
	th32ModuleID        uint32
	cntThreads          uint32
	th32ParentProcessID uint32
	pcPriClassBase      int32
	dwFlags             uint32
	szExeFile           [syscall.MAX_PATH]uint16
***REMOVED***

func getProcessEntry(pid int) (pe *processEntry32, err error) ***REMOVED***
	snapshot, _, e1 := CreateToolhelp32Snapshot.Call(th32cs_snapprocess, uintptr(0))
	if snapshot == uintptr(syscall.InvalidHandle) ***REMOVED***
		err = fmt.Errorf("CreateToolhelp32Snapshot: %v", e1)
		return
	***REMOVED***
	defer syscall.CloseHandle(syscall.Handle(snapshot))

	var processEntry processEntry32
	processEntry.dwSize = uint32(unsafe.Sizeof(processEntry))
	ok, _, e1 := Process32First.Call(snapshot, uintptr(unsafe.Pointer(&processEntry)))
	if ok == 0 ***REMOVED***
		err = fmt.Errorf("Process32First: %v", e1)
		return
	***REMOVED***

	for ***REMOVED***
		if processEntry.th32ProcessID == uint32(pid) ***REMOVED***
			pe = &processEntry
			return
		***REMOVED***

		ok, _, e1 = Process32Next.Call(snapshot, uintptr(unsafe.Pointer(&processEntry)))
		if ok == 0 ***REMOVED***
			err = fmt.Errorf("Process32Next: %v", e1)
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func getppid() (pid int, err error) ***REMOVED***
	pe, err := getProcessEntry(os.Getpid())
	if err != nil ***REMOVED***
		return
	***REMOVED***

	pid = int(pe.th32ParentProcessID)
	return
***REMOVED***

// StartedByExplorer returns true if the program was invoked by the user double-clicking
// on the executable from explorer.exe
//
// It is conservative and returns false if any of the internal calls fail.
// It does not guarantee that the program was run from a terminal. It only can tell you
// whether it was launched from explorer.exe
func StartedByExplorer() bool ***REMOVED***
	ppid, err := getppid()
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	pe, err := getProcessEntry(ppid)
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	name := syscall.UTF16ToString(pe.szExeFile[:])
	return name == "explorer.exe"
***REMOVED***
