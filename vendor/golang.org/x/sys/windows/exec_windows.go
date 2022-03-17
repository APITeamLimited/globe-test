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

// ComposeCommandLine escapes and joins the given arguments suitable for use as a Windows command line,
// in CreateProcess's CommandLine argument, CreateService/ChangeServiceConfig's BinaryPathName argument,
// or any program that uses CommandLineToArgv.
func ComposeCommandLine(args []string) string ***REMOVED***
	var commandLine string
	for i := range args ***REMOVED***
		if i > 0 ***REMOVED***
			commandLine += " "
		***REMOVED***
		commandLine += EscapeArg(args[i])
	***REMOVED***
	return commandLine
***REMOVED***

// DecomposeCommandLine breaks apart its argument command line into unescaped parts using CommandLineToArgv,
// as gathered from GetCommandLine, QUERY_SERVICE_CONFIG's BinaryPathName argument, or elsewhere that
// command lines are passed around.
func DecomposeCommandLine(commandLine string) ([]string, error) ***REMOVED***
	if len(commandLine) == 0 ***REMOVED***
		return []string***REMOVED******REMOVED***, nil
	***REMOVED***
	var argc int32
	argv, err := CommandLineToArgv(StringToUTF16Ptr(commandLine), &argc)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer LocalFree(Handle(unsafe.Pointer(argv)))
	var args []string
	for _, v := range (*argv)[:argc] ***REMOVED***
		args = append(args, UTF16ToString((*v)[:]))
	***REMOVED***
	return args, nil
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

// NewProcThreadAttributeList allocates a new ProcThreadAttributeListContainer, with the requested maximum number of attributes.
func NewProcThreadAttributeList(maxAttrCount uint32) (*ProcThreadAttributeListContainer, error) ***REMOVED***
	var size uintptr
	err := initializeProcThreadAttributeList(nil, maxAttrCount, 0, &size)
	if err != ERROR_INSUFFICIENT_BUFFER ***REMOVED***
		if err == nil ***REMOVED***
			return nil, errorspkg.New("unable to query buffer size from InitializeProcThreadAttributeList")
		***REMOVED***
		return nil, err
	***REMOVED***
	alloc, err := LocalAlloc(LMEM_FIXED, uint32(size))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// size is guaranteed to be ≥1 by InitializeProcThreadAttributeList.
	al := &ProcThreadAttributeListContainer***REMOVED***data: (*ProcThreadAttributeList)(unsafe.Pointer(alloc))***REMOVED***
	err = initializeProcThreadAttributeList(al.data, maxAttrCount, 0, &size)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return al, err
***REMOVED***

// Update modifies the ProcThreadAttributeList using UpdateProcThreadAttribute.
func (al *ProcThreadAttributeListContainer) Update(attribute uintptr, value unsafe.Pointer, size uintptr) error ***REMOVED***
	al.pointers = append(al.pointers, value)
	return updateProcThreadAttribute(al.data, 0, attribute, value, size, nil, nil)
***REMOVED***

// Delete frees ProcThreadAttributeList's resources.
func (al *ProcThreadAttributeListContainer) Delete() ***REMOVED***
	deleteProcThreadAttributeList(al.data)
	LocalFree(Handle(unsafe.Pointer(al.data)))
	al.data = nil
	al.pointers = nil
***REMOVED***

// List returns the actual ProcThreadAttributeList to be passed to StartupInfoEx.
func (al *ProcThreadAttributeListContainer) List() *ProcThreadAttributeList ***REMOVED***
	return al.data
***REMOVED***
