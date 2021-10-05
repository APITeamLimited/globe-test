// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Fork, exec, wait, etc.

package windows

import (
	errorspkg "errors"
	"unsafe"
)

// EscapeArg rewrites command line argument s as prescribed
// in http://msdn.microsoft.com/en-us/library/ms880421.
// This function returns "" (2 double quotes) if s is empty.
// Alternatively, these transformations are done:
// - every back slash (\) is doubled, but only if immediately
//   followed by double quote (");
// - every double quote (") is escaped by back slash (\);
// - finally, s is wrapped with double quotes (arg -> "arg"),
//   but only if there is space or tab inside s.
func EscapeArg(s string) string ***REMOVED***
	if len(s) == 0 ***REMOVED***
		return "\"\""
	***REMOVED***
	n := len(s)
	hasSpace := false
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		case '"', '\\':
			n++
		case ' ', '\t':
			hasSpace = true
		***REMOVED***
	***REMOVED***
	if hasSpace ***REMOVED***
		n += 2
	***REMOVED***
	if n == len(s) ***REMOVED***
		return s
	***REMOVED***

	qs := make([]byte, n)
	j := 0
	if hasSpace ***REMOVED***
		qs[j] = '"'
		j++
	***REMOVED***
	slashes := 0
	for i := 0; i < len(s); i++ ***REMOVED***
		switch s[i] ***REMOVED***
		default:
			slashes = 0
			qs[j] = s[i]
		case '\\':
			slashes++
			qs[j] = s[i]
		case '"':
			for ; slashes > 0; slashes-- ***REMOVED***
				qs[j] = '\\'
				j++
			***REMOVED***
			qs[j] = '\\'
			j++
			qs[j] = s[i]
		***REMOVED***
		j++
	***REMOVED***
	if hasSpace ***REMOVED***
		for ; slashes > 0; slashes-- ***REMOVED***
			qs[j] = '\\'
			j++
		***REMOVED***
		qs[j] = '"'
		j++
	***REMOVED***
	return string(qs[:j])
***REMOVED***

func CloseOnExec(fd Handle) ***REMOVED***
	SetHandleInformation(Handle(fd), HANDLE_FLAG_INHERIT, 0)
***REMOVED***

// FullPath retrieves the full path of the specified file.
func FullPath(name string) (path string, err error) ***REMOVED***
	p, err := UTF16PtrFromString(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	n := uint32(100)
	for ***REMOVED***
		buf := make([]uint16, n)
		n, err = GetFullPathName(p, uint32(len(buf)), &buf[0], nil)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
		if n <= uint32(len(buf)) ***REMOVED***
			return UTF16ToString(buf[:n]), nil
		***REMOVED***
	***REMOVED***
***REMOVED***

// NewProcThreadAttributeList allocates a new ProcThreadAttributeList, with the requested maximum number of attributes.
func NewProcThreadAttributeList(maxAttrCount uint32) (*ProcThreadAttributeList, error) ***REMOVED***
	var size uintptr
	err := initializeProcThreadAttributeList(nil, maxAttrCount, 0, &size)
	if err != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
		if err == nil ***REMOVED***
			return nil, errorspkg.New("unable to query buffer size from InitializeProcThreadAttributeList")
		***REMOVED***
		return nil, err
	***REMOVED***
	const psize = unsafe.Sizeof(uintptr(0))
	// size is guaranteed to be ≥1 by InitializeProcThreadAttributeList.
	al := (*ProcThreadAttributeList)(unsafe.Pointer(&make([]unsafe.Pointer, (size+psize-1)/psize)[0]))
	err = initializeProcThreadAttributeList(al, maxAttrCount, 0, &size)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return al, err
***REMOVED***

// Update modifies the ProcThreadAttributeList using UpdateProcThreadAttribute.
func (al *ProcThreadAttributeList) Update(attribute uintptr, flags uint32, value unsafe.Pointer, size uintptr, prevValue unsafe.Pointer, returnedSize *uintptr) error ***REMOVED***
	return updateProcThreadAttribute(al, flags, attribute, value, size, prevValue, returnedSize)
***REMOVED***

// Delete frees ProcThreadAttributeList's resources.
func (al *ProcThreadAttributeList) Delete() ***REMOVED***
	deleteProcThreadAttributeList(al)
***REMOVED***
