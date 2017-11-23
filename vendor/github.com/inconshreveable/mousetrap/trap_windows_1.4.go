// +build windows
// +build go1.4

package mousetrap

import (
	"os"
	"syscall"
	"unsafe"
)

func getProcessEntry(pid int) (*syscall.ProcessEntry32, error) ***REMOVED***
	snapshot, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer syscall.CloseHandle(snapshot)
	var procEntry syscall.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = syscall.Process32First(snapshot, &procEntry); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for ***REMOVED***
		if procEntry.ProcessID == uint32(pid) ***REMOVED***
			return &procEntry, nil
		***REMOVED***
		err = syscall.Process32Next(snapshot, &procEntry)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

// StartedByExplorer returns true if the program was invoked by the user double-clicking
// on the executable from explorer.exe
//
// It is conservative and returns false if any of the internal calls fail.
// It does not guarantee that the program was run from a terminal. It only can tell you
// whether it was launched from explorer.exe
func StartedByExplorer() bool ***REMOVED***
	pe, err := getProcessEntry(os.Getppid())
	if err != nil ***REMOVED***
		return false
	***REMOVED***
	return "explorer.exe" == syscall.UTF16ToString(pe.ExeFile[:])
***REMOVED***
