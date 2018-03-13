// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Windows system calls.

package windows

import (
	errorspkg "errors"
	"sync"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type Handle uintptr

const (
	InvalidHandle = ^Handle(0)

	// Flags for DefineDosDevice.
	DDD_EXACT_MATCH_ON_REMOVE = 0x00000004
	DDD_NO_BROADCAST_SYSTEM   = 0x00000008
	DDD_RAW_TARGET_PATH       = 0x00000001
	DDD_REMOVE_DEFINITION     = 0x00000002

	// Return values for GetDriveType.
	DRIVE_UNKNOWN     = 0
	DRIVE_NO_ROOT_DIR = 1
	DRIVE_REMOVABLE   = 2
	DRIVE_FIXED       = 3
	DRIVE_REMOTE      = 4
	DRIVE_CDROM       = 5
	DRIVE_RAMDISK     = 6

	// File system flags from GetVolumeInformation and GetVolumeInformationByHandle.
	FILE_CASE_SENSITIVE_SEARCH        = 0x00000001
	FILE_CASE_PRESERVED_NAMES         = 0x00000002
	FILE_FILE_COMPRESSION             = 0x00000010
	FILE_DAX_VOLUME                   = 0x20000000
	FILE_NAMED_STREAMS                = 0x00040000
	FILE_PERSISTENT_ACLS              = 0x00000008
	FILE_READ_ONLY_VOLUME             = 0x00080000
	FILE_SEQUENTIAL_WRITE_ONCE        = 0x00100000
	FILE_SUPPORTS_ENCRYPTION          = 0x00020000
	FILE_SUPPORTS_EXTENDED_ATTRIBUTES = 0x00800000
	FILE_SUPPORTS_HARD_LINKS          = 0x00400000
	FILE_SUPPORTS_OBJECT_IDS          = 0x00010000
	FILE_SUPPORTS_OPEN_BY_FILE_ID     = 0x01000000
	FILE_SUPPORTS_REPARSE_POINTS      = 0x00000080
	FILE_SUPPORTS_SPARSE_FILES        = 0x00000040
	FILE_SUPPORTS_TRANSACTIONS        = 0x00200000
	FILE_SUPPORTS_USN_JOURNAL         = 0x02000000
	FILE_UNICODE_ON_DISK              = 0x00000004
	FILE_VOLUME_IS_COMPRESSED         = 0x00008000
	FILE_VOLUME_QUOTAS                = 0x00000020
)

// StringToUTF16 is deprecated. Use UTF16FromString instead.
// If s contains a NUL byte this function panics instead of
// returning an error.
func StringToUTF16(s string) []uint16 ***REMOVED***
	a, err := UTF16FromString(s)
	if err != nil ***REMOVED***
		panic("windows: string with NUL passed to StringToUTF16")
	***REMOVED***
	return a
***REMOVED***

// UTF16FromString returns the UTF-16 encoding of the UTF-8 string
// s, with a terminating NUL added. If s contains a NUL byte at any
// location, it returns (nil, syscall.EINVAL).
func UTF16FromString(s string) ([]uint16, error) ***REMOVED***
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] == 0 ***REMOVED***
			return nil, syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return utf16.Encode([]rune(s + "\x00")), nil
***REMOVED***

// UTF16ToString returns the UTF-8 encoding of the UTF-16 sequence s,
// with a terminating NUL removed.
func UTF16ToString(s []uint16) string ***REMOVED***
	for i, v := range s ***REMOVED***
		if v == 0 ***REMOVED***
			s = s[0:i]
			break
		***REMOVED***
	***REMOVED***
	return string(utf16.Decode(s))
***REMOVED***

// StringToUTF16Ptr is deprecated. Use UTF16PtrFromString instead.
// If s contains a NUL byte this function panics instead of
// returning an error.
func StringToUTF16Ptr(s string) *uint16 ***REMOVED*** return &StringToUTF16(s)[0] ***REMOVED***

// UTF16PtrFromString returns pointer to the UTF-16 encoding of
// the UTF-8 string s, with a terminating NUL added. If s
// contains a NUL byte at any location, it returns (nil, syscall.EINVAL).
func UTF16PtrFromString(s string) (*uint16, error) ***REMOVED***
	a, err := UTF16FromString(s)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &a[0], nil
***REMOVED***

func Getpagesize() int ***REMOVED*** return 4096 ***REMOVED***

// NewCallback converts a Go function to a function pointer conforming to the stdcall calling convention.
// This is useful when interoperating with Windows code requiring callbacks.
func NewCallback(fn interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return syscall.NewCallback(fn)
***REMOVED***

// NewCallbackCDecl converts a Go function to a function pointer conforming to the cdecl calling convention.
// This is useful when interoperating with Windows code requiring callbacks.
func NewCallbackCDecl(fn interface***REMOVED******REMOVED***) uintptr ***REMOVED***
	return syscall.NewCallbackCDecl(fn)
***REMOVED***

// windows api calls

//sys	GetLastError() (lasterr error)
//sys	LoadLibrary(libname string) (handle Handle, err error) = LoadLibraryW
//sys	LoadLibraryEx(libname string, zero Handle, flags uintptr) (handle Handle, err error) = LoadLibraryExW
//sys	FreeLibrary(handle Handle) (err error)
//sys	GetProcAddress(module Handle, procname string) (proc uintptr, err error)
//sys	GetVersion() (ver uint32, err error)
//sys	FormatMessage(flags uint32, msgsrc uintptr, msgid uint32, langid uint32, buf []uint16, args *byte) (n uint32, err error) = FormatMessageW
//sys	ExitProcess(exitcode uint32)
//sys	CreateFile(name *uint16, access uint32, mode uint32, sa *SecurityAttributes, createmode uint32, attrs uint32, templatefile int32) (handle Handle, err error) [failretval==InvalidHandle] = CreateFileW
//sys	ReadFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (err error)
//sys	WriteFile(handle Handle, buf []byte, done *uint32, overlapped *Overlapped) (err error)
//sys	SetFilePointer(handle Handle, lowoffset int32, highoffsetptr *int32, whence uint32) (newlowoffset uint32, err error) [failretval==0xffffffff]
//sys	CloseHandle(handle Handle) (err error)
//sys	GetStdHandle(stdhandle uint32) (handle Handle, err error) [failretval==InvalidHandle]
//sys	SetStdHandle(stdhandle uint32, handle Handle) (err error)
//sys	findFirstFile1(name *uint16, data *win32finddata1) (handle Handle, err error) [failretval==InvalidHandle] = FindFirstFileW
//sys	findNextFile1(handle Handle, data *win32finddata1) (err error) = FindNextFileW
//sys	FindClose(handle Handle) (err error)
//sys	GetFileInformationByHandle(handle Handle, data *ByHandleFileInformation) (err error)
//sys	GetCurrentDirectory(buflen uint32, buf *uint16) (n uint32, err error) = GetCurrentDirectoryW
//sys	SetCurrentDirectory(path *uint16) (err error) = SetCurrentDirectoryW
//sys	CreateDirectory(path *uint16, sa *SecurityAttributes) (err error) = CreateDirectoryW
//sys	RemoveDirectory(path *uint16) (err error) = RemoveDirectoryW
//sys	DeleteFile(path *uint16) (err error) = DeleteFileW
//sys	MoveFile(from *uint16, to *uint16) (err error) = MoveFileW
//sys	MoveFileEx(from *uint16, to *uint16, flags uint32) (err error) = MoveFileExW
//sys	GetComputerName(buf *uint16, n *uint32) (err error) = GetComputerNameW
//sys	GetComputerNameEx(nametype uint32, buf *uint16, n *uint32) (err error) = GetComputerNameExW
//sys	SetEndOfFile(handle Handle) (err error)
//sys	GetSystemTimeAsFileTime(time *Filetime)
//sys	GetSystemTimePreciseAsFileTime(time *Filetime)
//sys	GetTimeZoneInformation(tzi *Timezoneinformation) (rc uint32, err error) [failretval==0xffffffff]
//sys	CreateIoCompletionPort(filehandle Handle, cphandle Handle, key uint32, threadcnt uint32) (handle Handle, err error)
//sys	GetQueuedCompletionStatus(cphandle Handle, qty *uint32, key *uint32, overlapped **Overlapped, timeout uint32) (err error)
//sys	PostQueuedCompletionStatus(cphandle Handle, qty uint32, key uint32, overlapped *Overlapped) (err error)
//sys	CancelIo(s Handle) (err error)
//sys	CancelIoEx(s Handle, o *Overlapped) (err error)
//sys	CreateProcess(appName *uint16, commandLine *uint16, procSecurity *SecurityAttributes, threadSecurity *SecurityAttributes, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *StartupInfo, outProcInfo *ProcessInformation) (err error) = CreateProcessW
//sys	OpenProcess(da uint32, inheritHandle bool, pid uint32) (handle Handle, err error)
//sys	TerminateProcess(handle Handle, exitcode uint32) (err error)
//sys	GetExitCodeProcess(handle Handle, exitcode *uint32) (err error)
//sys	GetStartupInfo(startupInfo *StartupInfo) (err error) = GetStartupInfoW
//sys	GetCurrentProcess() (pseudoHandle Handle, err error)
//sys	GetProcessTimes(handle Handle, creationTime *Filetime, exitTime *Filetime, kernelTime *Filetime, userTime *Filetime) (err error)
//sys	DuplicateHandle(hSourceProcessHandle Handle, hSourceHandle Handle, hTargetProcessHandle Handle, lpTargetHandle *Handle, dwDesiredAccess uint32, bInheritHandle bool, dwOptions uint32) (err error)
//sys	WaitForSingleObject(handle Handle, waitMilliseconds uint32) (event uint32, err error) [failretval==0xffffffff]
//sys	GetTempPath(buflen uint32, buf *uint16) (n uint32, err error) = GetTempPathW
//sys	CreatePipe(readhandle *Handle, writehandle *Handle, sa *SecurityAttributes, size uint32) (err error)
//sys	GetFileType(filehandle Handle) (n uint32, err error)
//sys	CryptAcquireContext(provhandle *Handle, container *uint16, provider *uint16, provtype uint32, flags uint32) (err error) = advapi32.CryptAcquireContextW
//sys	CryptReleaseContext(provhandle Handle, flags uint32) (err error) = advapi32.CryptReleaseContext
//sys	CryptGenRandom(provhandle Handle, buflen uint32, buf *byte) (err error) = advapi32.CryptGenRandom
//sys	GetEnvironmentStrings() (envs *uint16, err error) [failretval==nil] = kernel32.GetEnvironmentStringsW
//sys	FreeEnvironmentStrings(envs *uint16) (err error) = kernel32.FreeEnvironmentStringsW
//sys	GetEnvironmentVariable(name *uint16, buffer *uint16, size uint32) (n uint32, err error) = kernel32.GetEnvironmentVariableW
//sys	SetEnvironmentVariable(name *uint16, value *uint16) (err error) = kernel32.SetEnvironmentVariableW
//sys	SetFileTime(handle Handle, ctime *Filetime, atime *Filetime, wtime *Filetime) (err error)
//sys	GetFileAttributes(name *uint16) (attrs uint32, err error) [failretval==INVALID_FILE_ATTRIBUTES] = kernel32.GetFileAttributesW
//sys	SetFileAttributes(name *uint16, attrs uint32) (err error) = kernel32.SetFileAttributesW
//sys	GetFileAttributesEx(name *uint16, level uint32, info *byte) (err error) = kernel32.GetFileAttributesExW
//sys	GetCommandLine() (cmd *uint16) = kernel32.GetCommandLineW
//sys	CommandLineToArgv(cmd *uint16, argc *int32) (argv *[8192]*[8192]uint16, err error) [failretval==nil] = shell32.CommandLineToArgvW
//sys	LocalFree(hmem Handle) (handle Handle, err error) [failretval!=0]
//sys	SetHandleInformation(handle Handle, mask uint32, flags uint32) (err error)
//sys	FlushFileBuffers(handle Handle) (err error)
//sys	GetFullPathName(path *uint16, buflen uint32, buf *uint16, fname **uint16) (n uint32, err error) = kernel32.GetFullPathNameW
//sys	GetLongPathName(path *uint16, buf *uint16, buflen uint32) (n uint32, err error) = kernel32.GetLongPathNameW
//sys	GetShortPathName(longpath *uint16, shortpath *uint16, buflen uint32) (n uint32, err error) = kernel32.GetShortPathNameW
//sys	CreateFileMapping(fhandle Handle, sa *SecurityAttributes, prot uint32, maxSizeHigh uint32, maxSizeLow uint32, name *uint16) (handle Handle, err error) = kernel32.CreateFileMappingW
//sys	MapViewOfFile(handle Handle, access uint32, offsetHigh uint32, offsetLow uint32, length uintptr) (addr uintptr, err error)
//sys	UnmapViewOfFile(addr uintptr) (err error)
//sys	FlushViewOfFile(addr uintptr, length uintptr) (err error)
//sys	VirtualLock(addr uintptr, length uintptr) (err error)
//sys	VirtualUnlock(addr uintptr, length uintptr) (err error)
//sys	VirtualAlloc(address uintptr, size uintptr, alloctype uint32, protect uint32) (value uintptr, err error) = kernel32.VirtualAlloc
//sys	VirtualFree(address uintptr, size uintptr, freetype uint32) (err error) = kernel32.VirtualFree
//sys	VirtualProtect(address uintptr, size uintptr, newprotect uint32, oldprotect *uint32) (err error) = kernel32.VirtualProtect
//sys	TransmitFile(s Handle, handle Handle, bytesToWrite uint32, bytsPerSend uint32, overlapped *Overlapped, transmitFileBuf *TransmitFileBuffers, flags uint32) (err error) = mswsock.TransmitFile
//sys	ReadDirectoryChanges(handle Handle, buf *byte, buflen uint32, watchSubTree bool, mask uint32, retlen *uint32, overlapped *Overlapped, completionRoutine uintptr) (err error) = kernel32.ReadDirectoryChangesW
//sys	CertOpenSystemStore(hprov Handle, name *uint16) (store Handle, err error) = crypt32.CertOpenSystemStoreW
//sys   CertOpenStore(storeProvider uintptr, msgAndCertEncodingType uint32, cryptProv uintptr, flags uint32, para uintptr) (handle Handle, err error) [failretval==InvalidHandle] = crypt32.CertOpenStore
//sys	CertEnumCertificatesInStore(store Handle, prevContext *CertContext) (context *CertContext, err error) [failretval==nil] = crypt32.CertEnumCertificatesInStore
//sys   CertAddCertificateContextToStore(store Handle, certContext *CertContext, addDisposition uint32, storeContext **CertContext) (err error) = crypt32.CertAddCertificateContextToStore
//sys	CertCloseStore(store Handle, flags uint32) (err error) = crypt32.CertCloseStore
//sys   CertGetCertificateChain(engine Handle, leaf *CertContext, time *Filetime, additionalStore Handle, para *CertChainPara, flags uint32, reserved uintptr, chainCtx **CertChainContext) (err error) = crypt32.CertGetCertificateChain
//sys   CertFreeCertificateChain(ctx *CertChainContext) = crypt32.CertFreeCertificateChain
//sys   CertCreateCertificateContext(certEncodingType uint32, certEncoded *byte, encodedLen uint32) (context *CertContext, err error) [failretval==nil] = crypt32.CertCreateCertificateContext
//sys   CertFreeCertificateContext(ctx *CertContext) (err error) = crypt32.CertFreeCertificateContext
//sys   CertVerifyCertificateChainPolicy(policyOID uintptr, chain *CertChainContext, para *CertChainPolicyPara, status *CertChainPolicyStatus) (err error) = crypt32.CertVerifyCertificateChainPolicy
//sys	RegOpenKeyEx(key Handle, subkey *uint16, options uint32, desiredAccess uint32, result *Handle) (regerrno error) = advapi32.RegOpenKeyExW
//sys	RegCloseKey(key Handle) (regerrno error) = advapi32.RegCloseKey
//sys	RegQueryInfoKey(key Handle, class *uint16, classLen *uint32, reserved *uint32, subkeysLen *uint32, maxSubkeyLen *uint32, maxClassLen *uint32, valuesLen *uint32, maxValueNameLen *uint32, maxValueLen *uint32, saLen *uint32, lastWriteTime *Filetime) (regerrno error) = advapi32.RegQueryInfoKeyW
//sys	RegEnumKeyEx(key Handle, index uint32, name *uint16, nameLen *uint32, reserved *uint32, class *uint16, classLen *uint32, lastWriteTime *Filetime) (regerrno error) = advapi32.RegEnumKeyExW
//sys	RegQueryValueEx(key Handle, name *uint16, reserved *uint32, valtype *uint32, buf *byte, buflen *uint32) (regerrno error) = advapi32.RegQueryValueExW
//sys	getCurrentProcessId() (pid uint32) = kernel32.GetCurrentProcessId
//sys	GetConsoleMode(console Handle, mode *uint32) (err error) = kernel32.GetConsoleMode
//sys	SetConsoleMode(console Handle, mode uint32) (err error) = kernel32.SetConsoleMode
//sys	GetConsoleScreenBufferInfo(console Handle, info *ConsoleScreenBufferInfo) (err error) = kernel32.GetConsoleScreenBufferInfo
//sys	WriteConsole(console Handle, buf *uint16, towrite uint32, written *uint32, reserved *byte) (err error) = kernel32.WriteConsoleW
//sys	ReadConsole(console Handle, buf *uint16, toread uint32, read *uint32, inputControl *byte) (err error) = kernel32.ReadConsoleW
//sys	CreateToolhelp32Snapshot(flags uint32, processId uint32) (handle Handle, err error) [failretval==InvalidHandle] = kernel32.CreateToolhelp32Snapshot
//sys	Process32First(snapshot Handle, procEntry *ProcessEntry32) (err error) = kernel32.Process32FirstW
//sys	Process32Next(snapshot Handle, procEntry *ProcessEntry32) (err error) = kernel32.Process32NextW
//sys	DeviceIoControl(handle Handle, ioControlCode uint32, inBuffer *byte, inBufferSize uint32, outBuffer *byte, outBufferSize uint32, bytesReturned *uint32, overlapped *Overlapped) (err error)
// This function returns 1 byte BOOLEAN rather than the 4 byte BOOL.
//sys	CreateSymbolicLink(symlinkfilename *uint16, targetfilename *uint16, flags uint32) (err error) [failretval&0xff==0] = CreateSymbolicLinkW
//sys	CreateHardLink(filename *uint16, existingfilename *uint16, reserved uintptr) (err error) [failretval&0xff==0] = CreateHardLinkW
//sys	GetCurrentThreadId() (id uint32)
//sys	CreateEvent(eventAttrs *SecurityAttributes, manualReset uint32, initialState uint32, name *uint16) (handle Handle, err error) = kernel32.CreateEventW
//sys	CreateEventEx(eventAttrs *SecurityAttributes, name *uint16, flags uint32, desiredAccess uint32) (handle Handle, err error) = kernel32.CreateEventExW
//sys	OpenEvent(desiredAccess uint32, inheritHandle bool, name *uint16) (handle Handle, err error) = kernel32.OpenEventW
//sys	SetEvent(event Handle) (err error) = kernel32.SetEvent
//sys	ResetEvent(event Handle) (err error) = kernel32.ResetEvent
//sys	PulseEvent(event Handle) (err error) = kernel32.PulseEvent

// Volume Management Functions
//sys	DefineDosDevice(flags uint32, deviceName *uint16, targetPath *uint16) (err error) = DefineDosDeviceW
//sys	DeleteVolumeMountPoint(volumeMountPoint *uint16) (err error) = DeleteVolumeMountPointW
//sys	FindFirstVolume(volumeName *uint16, bufferLength uint32) (handle Handle, err error) [failretval==InvalidHandle] = FindFirstVolumeW
//sys	FindFirstVolumeMountPoint(rootPathName *uint16, volumeMountPoint *uint16, bufferLength uint32) (handle Handle, err error) [failretval==InvalidHandle] = FindFirstVolumeMountPointW
//sys	FindNextVolume(findVolume Handle, volumeName *uint16, bufferLength uint32) (err error) = FindNextVolumeW
//sys	FindNextVolumeMountPoint(findVolumeMountPoint Handle, volumeMountPoint *uint16, bufferLength uint32) (err error) = FindNextVolumeMountPointW
//sys	FindVolumeClose(findVolume Handle) (err error)
//sys	FindVolumeMountPointClose(findVolumeMountPoint Handle) (err error)
//sys	GetDriveType(rootPathName *uint16) (driveType uint32) = GetDriveTypeW
//sys	GetLogicalDrives() (drivesBitMask uint32, err error) [failretval==0]
//sys	GetLogicalDriveStrings(bufferLength uint32, buffer *uint16) (n uint32, err error) [failretval==0] = GetLogicalDriveStringsW
//sys	GetVolumeInformation(rootPathName *uint16, volumeNameBuffer *uint16, volumeNameSize uint32, volumeNameSerialNumber *uint32, maximumComponentLength *uint32, fileSystemFlags *uint32, fileSystemNameBuffer *uint16, fileSystemNameSize uint32) (err error) = GetVolumeInformationW
//sys	GetVolumeInformationByHandle(file Handle, volumeNameBuffer *uint16, volumeNameSize uint32, volumeNameSerialNumber *uint32, maximumComponentLength *uint32, fileSystemFlags *uint32, fileSystemNameBuffer *uint16, fileSystemNameSize uint32) (err error) = GetVolumeInformationByHandleW
//sys	GetVolumeNameForVolumeMountPoint(volumeMountPoint *uint16, volumeName *uint16, bufferlength uint32) (err error) = GetVolumeNameForVolumeMountPointW
//sys	GetVolumePathName(fileName *uint16, volumePathName *uint16, bufferLength uint32) (err error) = GetVolumePathNameW
//sys	GetVolumePathNamesForVolumeName(volumeName *uint16, volumePathNames *uint16, bufferLength uint32, returnLength *uint32) (err error) = GetVolumePathNamesForVolumeNameW
//sys	QueryDosDevice(deviceName *uint16, targetPath *uint16, max uint32) (n uint32, err error) [failretval==0] = QueryDosDeviceW
//sys	SetVolumeLabel(rootPathName *uint16, volumeName *uint16) (err error) = SetVolumeLabelW
//sys	SetVolumeMountPoint(volumeMountPoint *uint16, volumeName *uint16) (err error) = SetVolumeMountPointW

// syscall interface implementation for other packages

// GetProcAddressByOrdinal retrieves the address of the exported
// function from module by ordinal.
func GetProcAddressByOrdinal(module Handle, ordinal uintptr) (proc uintptr, err error) ***REMOVED***
	r0, _, e1 := syscall.Syscall(procGetProcAddress.Addr(), 2, uintptr(module), ordinal, 0)
	proc = uintptr(r0)
	if proc == 0 ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = errnoErr(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func Exit(code int) ***REMOVED*** ExitProcess(uint32(code)) ***REMOVED***

func makeInheritSa() *SecurityAttributes ***REMOVED***
	var sa SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
***REMOVED***

func Open(path string, mode int, perm uint32) (fd Handle, err error) ***REMOVED***
	if len(path) == 0 ***REMOVED***
		return InvalidHandle, ERROR_FILE_NOT_FOUND
	***REMOVED***
	pathp, err := UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return InvalidHandle, err
	***REMOVED***
	var access uint32
	switch mode & (O_RDONLY | O_WRONLY | O_RDWR) ***REMOVED***
	case O_RDONLY:
		access = GENERIC_READ
	case O_WRONLY:
		access = GENERIC_WRITE
	case O_RDWR:
		access = GENERIC_READ | GENERIC_WRITE
	***REMOVED***
	if mode&O_CREAT != 0 ***REMOVED***
		access |= GENERIC_WRITE
	***REMOVED***
	if mode&O_APPEND != 0 ***REMOVED***
		access &^= GENERIC_WRITE
		access |= FILE_APPEND_DATA
	***REMOVED***
	sharemode := uint32(FILE_SHARE_READ | FILE_SHARE_WRITE)
	var sa *SecurityAttributes
	if mode&O_CLOEXEC == 0 ***REMOVED***
		sa = makeInheritSa()
	***REMOVED***
	var createmode uint32
	switch ***REMOVED***
	case mode&(O_CREAT|O_EXCL) == (O_CREAT | O_EXCL):
		createmode = CREATE_NEW
	case mode&(O_CREAT|O_TRUNC) == (O_CREAT | O_TRUNC):
		createmode = CREATE_ALWAYS
	case mode&O_CREAT == O_CREAT:
		createmode = OPEN_ALWAYS
	case mode&O_TRUNC == O_TRUNC:
		createmode = TRUNCATE_EXISTING
	default:
		createmode = OPEN_EXISTING
	***REMOVED***
	h, e := CreateFile(pathp, access, sharemode, sa, createmode, FILE_ATTRIBUTE_NORMAL, 0)
	return h, e
***REMOVED***

func Read(fd Handle, p []byte) (n int, err error) ***REMOVED***
	var done uint32
	e := ReadFile(fd, p, &done, nil)
	if e != nil ***REMOVED***
		if e == ERROR_BROKEN_PIPE ***REMOVED***
			// NOTE(brainman): work around ERROR_BROKEN_PIPE is returned on reading EOF from stdin
			return 0, nil
		***REMOVED***
		return 0, e
	***REMOVED***
	if raceenabled ***REMOVED***
		if done > 0 ***REMOVED***
			raceWriteRange(unsafe.Pointer(&p[0]), int(done))
		***REMOVED***
		raceAcquire(unsafe.Pointer(&ioSync))
	***REMOVED***
	return int(done), nil
***REMOVED***

func Write(fd Handle, p []byte) (n int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	var done uint32
	e := WriteFile(fd, p, &done, nil)
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	if raceenabled && done > 0 ***REMOVED***
		raceReadRange(unsafe.Pointer(&p[0]), int(done))
	***REMOVED***
	return int(done), nil
***REMOVED***

var ioSync int64

func Seek(fd Handle, offset int64, whence int) (newoffset int64, err error) ***REMOVED***
	var w uint32
	switch whence ***REMOVED***
	case 0:
		w = FILE_BEGIN
	case 1:
		w = FILE_CURRENT
	case 2:
		w = FILE_END
	***REMOVED***
	hi := int32(offset >> 32)
	lo := int32(offset)
	// use GetFileType to check pipe, pipe can't do seek
	ft, _ := GetFileType(fd)
	if ft == FILE_TYPE_PIPE ***REMOVED***
		return 0, syscall.EPIPE
	***REMOVED***
	rlo, e := SetFilePointer(fd, lo, &hi, w)
	if e != nil ***REMOVED***
		return 0, e
	***REMOVED***
	return int64(hi)<<32 + int64(rlo), nil
***REMOVED***

func Close(fd Handle) (err error) ***REMOVED***
	return CloseHandle(fd)
***REMOVED***

var (
	Stdin  = getStdHandle(STD_INPUT_HANDLE)
	Stdout = getStdHandle(STD_OUTPUT_HANDLE)
	Stderr = getStdHandle(STD_ERROR_HANDLE)
)

func getStdHandle(stdhandle uint32) (fd Handle) ***REMOVED***
	r, _ := GetStdHandle(stdhandle)
	CloseOnExec(r)
	return r
***REMOVED***

const ImplementsGetwd = true

func Getwd() (wd string, err error) ***REMOVED***
	b := make([]uint16, 300)
	n, e := GetCurrentDirectory(uint32(len(b)), &b[0])
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	return string(utf16.Decode(b[0:n])), nil
***REMOVED***

func Chdir(path string) (err error) ***REMOVED***
	pathp, err := UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return SetCurrentDirectory(pathp)
***REMOVED***

func Mkdir(path string, mode uint32) (err error) ***REMOVED***
	pathp, err := UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return CreateDirectory(pathp, nil)
***REMOVED***

func Rmdir(path string) (err error) ***REMOVED***
	pathp, err := UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return RemoveDirectory(pathp)
***REMOVED***

func Unlink(path string) (err error) ***REMOVED***
	pathp, err := UTF16PtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return DeleteFile(pathp)
***REMOVED***

func Rename(oldpath, newpath string) (err error) ***REMOVED***
	from, err := UTF16PtrFromString(oldpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	to, err := UTF16PtrFromString(newpath)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return MoveFileEx(from, to, MOVEFILE_REPLACE_EXISTING)
***REMOVED***

func ComputerName() (name string, err error) ***REMOVED***
	var n uint32 = MAX_COMPUTERNAME_LENGTH + 1
	b := make([]uint16, n)
	e := GetComputerName(&b[0], &n)
	if e != nil ***REMOVED***
		return "", e
	***REMOVED***
	return string(utf16.Decode(b[0:n])), nil
***REMOVED***

func Ftruncate(fd Handle, length int64) (err error) ***REMOVED***
	curoffset, e := Seek(fd, 0, 1)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	defer Seek(fd, curoffset, 0)
	_, e = Seek(fd, length, 0)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	e = SetEndOfFile(fd)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	return nil
***REMOVED***

func Gettimeofday(tv *Timeval) (err error) ***REMOVED***
	var ft Filetime
	GetSystemTimeAsFileTime(&ft)
	*tv = NsecToTimeval(ft.Nanoseconds())
	return nil
***REMOVED***

func Pipe(p []Handle) (err error) ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	var r, w Handle
	e := CreatePipe(&r, &w, makeInheritSa(), 0)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	p[0] = r
	p[1] = w
	return nil
***REMOVED***

func Utimes(path string, tv []Timeval) (err error) ***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	pathp, e := UTF16PtrFromString(path)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	h, e := CreateFile(pathp,
		FILE_WRITE_ATTRIBUTES, FILE_SHARE_WRITE, nil,
		OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	defer Close(h)
	a := NsecToFiletime(tv[0].Nanoseconds())
	w := NsecToFiletime(tv[1].Nanoseconds())
	return SetFileTime(h, nil, &a, &w)
***REMOVED***

func UtimesNano(path string, ts []Timespec) (err error) ***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	pathp, e := UTF16PtrFromString(path)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	h, e := CreateFile(pathp,
		FILE_WRITE_ATTRIBUTES, FILE_SHARE_WRITE, nil,
		OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	defer Close(h)
	a := NsecToFiletime(TimespecToNsec(ts[0]))
	w := NsecToFiletime(TimespecToNsec(ts[1]))
	return SetFileTime(h, nil, &a, &w)
***REMOVED***

func Fsync(fd Handle) (err error) ***REMOVED***
	return FlushFileBuffers(fd)
***REMOVED***

func Chmod(path string, mode uint32) (err error) ***REMOVED***
	if mode == 0 ***REMOVED***
		return syscall.EINVAL
	***REMOVED***
	p, e := UTF16PtrFromString(path)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	attrs, e := GetFileAttributes(p)
	if e != nil ***REMOVED***
		return e
	***REMOVED***
	if mode&S_IWRITE != 0 ***REMOVED***
		attrs &^= FILE_ATTRIBUTE_READONLY
	***REMOVED*** else ***REMOVED***
		attrs |= FILE_ATTRIBUTE_READONLY
	***REMOVED***
	return SetFileAttributes(p, attrs)
***REMOVED***

func LoadGetSystemTimePreciseAsFileTime() error ***REMOVED***
	return procGetSystemTimePreciseAsFileTime.Find()
***REMOVED***

func LoadCancelIoEx() error ***REMOVED***
	return procCancelIoEx.Find()
***REMOVED***

func LoadSetFileCompletionNotificationModes() error ***REMOVED***
	return procSetFileCompletionNotificationModes.Find()
***REMOVED***

// net api calls

const socket_error = uintptr(^uint32(0))

//sys	WSAStartup(verreq uint32, data *WSAData) (sockerr error) = ws2_32.WSAStartup
//sys	WSACleanup() (err error) [failretval==socket_error] = ws2_32.WSACleanup
//sys	WSAIoctl(s Handle, iocc uint32, inbuf *byte, cbif uint32, outbuf *byte, cbob uint32, cbbr *uint32, overlapped *Overlapped, completionRoutine uintptr) (err error) [failretval==socket_error] = ws2_32.WSAIoctl
//sys	socket(af int32, typ int32, protocol int32) (handle Handle, err error) [failretval==InvalidHandle] = ws2_32.socket
//sys	Setsockopt(s Handle, level int32, optname int32, optval *byte, optlen int32) (err error) [failretval==socket_error] = ws2_32.setsockopt
//sys	Getsockopt(s Handle, level int32, optname int32, optval *byte, optlen *int32) (err error) [failretval==socket_error] = ws2_32.getsockopt
//sys	bind(s Handle, name unsafe.Pointer, namelen int32) (err error) [failretval==socket_error] = ws2_32.bind
//sys	connect(s Handle, name unsafe.Pointer, namelen int32) (err error) [failretval==socket_error] = ws2_32.connect
//sys	getsockname(s Handle, rsa *RawSockaddrAny, addrlen *int32) (err error) [failretval==socket_error] = ws2_32.getsockname
//sys	getpeername(s Handle, rsa *RawSockaddrAny, addrlen *int32) (err error) [failretval==socket_error] = ws2_32.getpeername
//sys	listen(s Handle, backlog int32) (err error) [failretval==socket_error] = ws2_32.listen
//sys	shutdown(s Handle, how int32) (err error) [failretval==socket_error] = ws2_32.shutdown
//sys	Closesocket(s Handle) (err error) [failretval==socket_error] = ws2_32.closesocket
//sys	AcceptEx(ls Handle, as Handle, buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, recvd *uint32, overlapped *Overlapped) (err error) = mswsock.AcceptEx
//sys	GetAcceptExSockaddrs(buf *byte, rxdatalen uint32, laddrlen uint32, raddrlen uint32, lrsa **RawSockaddrAny, lrsalen *int32, rrsa **RawSockaddrAny, rrsalen *int32) = mswsock.GetAcceptExSockaddrs
//sys	WSARecv(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32, overlapped *Overlapped, croutine *byte) (err error) [failretval==socket_error] = ws2_32.WSARecv
//sys	WSASend(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, overlapped *Overlapped, croutine *byte) (err error) [failretval==socket_error] = ws2_32.WSASend
//sys	WSARecvFrom(s Handle, bufs *WSABuf, bufcnt uint32, recvd *uint32, flags *uint32,  from *RawSockaddrAny, fromlen *int32, overlapped *Overlapped, croutine *byte) (err error) [failretval==socket_error] = ws2_32.WSARecvFrom
//sys	WSASendTo(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, to *RawSockaddrAny, tolen int32,  overlapped *Overlapped, croutine *byte) (err error) [failretval==socket_error] = ws2_32.WSASendTo
//sys	GetHostByName(name string) (h *Hostent, err error) [failretval==nil] = ws2_32.gethostbyname
//sys	GetServByName(name string, proto string) (s *Servent, err error) [failretval==nil] = ws2_32.getservbyname
//sys	Ntohs(netshort uint16) (u uint16) = ws2_32.ntohs
//sys	GetProtoByName(name string) (p *Protoent, err error) [failretval==nil] = ws2_32.getprotobyname
//sys	DnsQuery(name string, qtype uint16, options uint32, extra *byte, qrs **DNSRecord, pr *byte) (status error) = dnsapi.DnsQuery_W
//sys	DnsRecordListFree(rl *DNSRecord, freetype uint32) = dnsapi.DnsRecordListFree
//sys	DnsNameCompare(name1 *uint16, name2 *uint16) (same bool) = dnsapi.DnsNameCompare_W
//sys	GetAddrInfoW(nodename *uint16, servicename *uint16, hints *AddrinfoW, result **AddrinfoW) (sockerr error) = ws2_32.GetAddrInfoW
//sys	FreeAddrInfoW(addrinfo *AddrinfoW) = ws2_32.FreeAddrInfoW
//sys	GetIfEntry(pIfRow *MibIfRow) (errcode error) = iphlpapi.GetIfEntry
//sys	GetAdaptersInfo(ai *IpAdapterInfo, ol *uint32) (errcode error) = iphlpapi.GetAdaptersInfo
//sys	SetFileCompletionNotificationModes(handle Handle, flags uint8) (err error) = kernel32.SetFileCompletionNotificationModes
//sys	WSAEnumProtocols(protocols *int32, protocolBuffer *WSAProtocolInfo, bufferLength *uint32) (n int32, err error) [failretval==-1] = ws2_32.WSAEnumProtocolsW
//sys	GetAdaptersAddresses(family uint32, flags uint32, reserved uintptr, adapterAddresses *IpAdapterAddresses, sizePointer *uint32) (errcode error) = iphlpapi.GetAdaptersAddresses
//sys	GetACP() (acp uint32) = kernel32.GetACP
//sys	MultiByteToWideChar(codePage uint32, dwFlags uint32, str *byte, nstr int32, wchar *uint16, nwchar int32) (nwrite int32, err error) = kernel32.MultiByteToWideChar

// For testing: clients can set this flag to force
// creation of IPv6 sockets to return EAFNOSUPPORT.
var SocketDisableIPv6 bool

type RawSockaddrInet4 struct ***REMOVED***
	Family uint16
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
***REMOVED***

type RawSockaddrInet6 struct ***REMOVED***
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

type RawSockaddr struct ***REMOVED***
	Family uint16
	Data   [14]int8
***REMOVED***

type RawSockaddrAny struct ***REMOVED***
	Addr RawSockaddr
	Pad  [96]int8
***REMOVED***

type Sockaddr interface ***REMOVED***
	sockaddr() (ptr unsafe.Pointer, len int32, err error) // lowercase; only we can define Sockaddrs
***REMOVED***

type SockaddrInet4 struct ***REMOVED***
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
***REMOVED***

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, int32, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, syscall.EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), int32(unsafe.Sizeof(sa.raw)), nil
***REMOVED***

type SockaddrInet6 struct ***REMOVED***
	Port   int
	ZoneId uint32
	Addr   [16]byte
	raw    RawSockaddrInet6
***REMOVED***

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, int32, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, syscall.EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), int32(unsafe.Sizeof(sa.raw)), nil
***REMOVED***

type SockaddrUnix struct ***REMOVED***
	Name string
***REMOVED***

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, int32, error) ***REMOVED***
	// TODO(brainman): implement SockaddrUnix.sockaddr()
	return nil, 0, syscall.EWINDOWS
***REMOVED***

func (rsa *RawSockaddrAny) Sockaddr() (Sockaddr, error) ***REMOVED***
	switch rsa.Addr.Family ***REMOVED***
	case AF_UNIX:
		return nil, syscall.EWINDOWS

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		for i := 0; i < len(sa.Addr); i++ ***REMOVED***
			sa.Addr[i] = pp.Addr[i]
		***REMOVED***
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		for i := 0; i < len(sa.Addr); i++ ***REMOVED***
			sa.Addr[i] = pp.Addr[i]
		***REMOVED***
		return sa, nil
	***REMOVED***
	return nil, syscall.EAFNOSUPPORT
***REMOVED***

func Socket(domain, typ, proto int) (fd Handle, err error) ***REMOVED***
	if domain == AF_INET6 && SocketDisableIPv6 ***REMOVED***
		return InvalidHandle, syscall.EAFNOSUPPORT
	***REMOVED***
	return socket(int32(domain), int32(typ), int32(proto))
***REMOVED***

func SetsockoptInt(fd Handle, level, opt int, value int) (err error) ***REMOVED***
	v := int32(value)
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&v)), int32(unsafe.Sizeof(v)))
***REMOVED***

func Bind(fd Handle, sa Sockaddr) (err error) ***REMOVED***
	ptr, n, err := sa.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return bind(fd, ptr, n)
***REMOVED***

func Connect(fd Handle, sa Sockaddr) (err error) ***REMOVED***
	ptr, n, err := sa.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return connect(fd, ptr, n)
***REMOVED***

func Getsockname(fd Handle) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	l := int32(unsafe.Sizeof(rsa))
	if err = getsockname(fd, &rsa, &l); err != nil ***REMOVED***
		return
	***REMOVED***
	return rsa.Sockaddr()
***REMOVED***

func Getpeername(fd Handle) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	l := int32(unsafe.Sizeof(rsa))
	if err = getpeername(fd, &rsa, &l); err != nil ***REMOVED***
		return
	***REMOVED***
	return rsa.Sockaddr()
***REMOVED***

func Listen(s Handle, n int) (err error) ***REMOVED***
	return listen(s, int32(n))
***REMOVED***

func Shutdown(fd Handle, how int) (err error) ***REMOVED***
	return shutdown(fd, int32(how))
***REMOVED***

func WSASendto(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, to Sockaddr, overlapped *Overlapped, croutine *byte) (err error) ***REMOVED***
	rsa, l, err := to.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return WSASendTo(s, bufs, bufcnt, sent, flags, (*RawSockaddrAny)(unsafe.Pointer(rsa)), l, overlapped, croutine)
***REMOVED***

func LoadGetAddrInfo() error ***REMOVED***
	return procGetAddrInfoW.Find()
***REMOVED***

var connectExFunc struct ***REMOVED***
	once sync.Once
	addr uintptr
	err  error
***REMOVED***

func LoadConnectEx() error ***REMOVED***
	connectExFunc.once.Do(func() ***REMOVED***
		var s Handle
		s, connectExFunc.err = Socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)
		if connectExFunc.err != nil ***REMOVED***
			return
		***REMOVED***
		defer CloseHandle(s)
		var n uint32
		connectExFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_CONNECTEX)),
			uint32(unsafe.Sizeof(WSAID_CONNECTEX)),
			(*byte)(unsafe.Pointer(&connectExFunc.addr)),
			uint32(unsafe.Sizeof(connectExFunc.addr)),
			&n, nil, 0)
	***REMOVED***)
	return connectExFunc.err
***REMOVED***

func connectEx(s Handle, name unsafe.Pointer, namelen int32, sendBuf *byte, sendDataLen uint32, bytesSent *uint32, overlapped *Overlapped) (err error) ***REMOVED***
	r1, _, e1 := syscall.Syscall9(connectExFunc.addr, 7, uintptr(s), uintptr(name), uintptr(namelen), uintptr(unsafe.Pointer(sendBuf)), uintptr(sendDataLen), uintptr(unsafe.Pointer(bytesSent)), uintptr(unsafe.Pointer(overlapped)), 0, 0)
	if r1 == 0 ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = error(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func ConnectEx(fd Handle, sa Sockaddr, sendBuf *byte, sendDataLen uint32, bytesSent *uint32, overlapped *Overlapped) error ***REMOVED***
	err := LoadConnectEx()
	if err != nil ***REMOVED***
		return errorspkg.New("failed to find ConnectEx: " + err.Error())
	***REMOVED***
	ptr, n, err := sa.sockaddr()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return connectEx(fd, ptr, n, sendBuf, sendDataLen, bytesSent, overlapped)
***REMOVED***

var sendRecvMsgFunc struct ***REMOVED***
	once     sync.Once
	sendAddr uintptr
	recvAddr uintptr
	err      error
***REMOVED***

func loadWSASendRecvMsg() error ***REMOVED***
	sendRecvMsgFunc.once.Do(func() ***REMOVED***
		var s Handle
		s, sendRecvMsgFunc.err = Socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)
		if sendRecvMsgFunc.err != nil ***REMOVED***
			return
		***REMOVED***
		defer CloseHandle(s)
		var n uint32
		sendRecvMsgFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSARECVMSG)),
			uint32(unsafe.Sizeof(WSAID_WSARECVMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.recvAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.recvAddr)),
			&n, nil, 0)
		if sendRecvMsgFunc.err != nil ***REMOVED***
			return
		***REMOVED***
		sendRecvMsgFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSASENDMSG)),
			uint32(unsafe.Sizeof(WSAID_WSASENDMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.sendAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.sendAddr)),
			&n, nil, 0)
	***REMOVED***)
	return sendRecvMsgFunc.err
***REMOVED***

func WSASendMsg(fd Handle, msg *WSAMsg, flags uint32, bytesSent *uint32, overlapped *Overlapped, croutine *byte) error ***REMOVED***
	err := loadWSASendRecvMsg()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.sendAddr, 6, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(flags), uintptr(unsafe.Pointer(bytesSent)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
	if r1 == socket_error ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = errnoErr(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

func WSARecvMsg(fd Handle, msg *WSAMsg, bytesReceived *uint32, overlapped *Overlapped, croutine *byte) error ***REMOVED***
	err := loadWSASendRecvMsg()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.recvAddr, 5, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(unsafe.Pointer(bytesReceived)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0)
	if r1 == socket_error ***REMOVED***
		if e1 != 0 ***REMOVED***
			err = errnoErr(e1)
		***REMOVED*** else ***REMOVED***
			err = syscall.EINVAL
		***REMOVED***
	***REMOVED***
	return err
***REMOVED***

// Invented structures to support what package os expects.
type Rusage struct ***REMOVED***
	CreationTime Filetime
	ExitTime     Filetime
	KernelTime   Filetime
	UserTime     Filetime
***REMOVED***

type WaitStatus struct ***REMOVED***
	ExitCode uint32
***REMOVED***

func (w WaitStatus) Exited() bool ***REMOVED*** return true ***REMOVED***

func (w WaitStatus) ExitStatus() int ***REMOVED*** return int(w.ExitCode) ***REMOVED***

func (w WaitStatus) Signal() Signal ***REMOVED*** return -1 ***REMOVED***

func (w WaitStatus) CoreDump() bool ***REMOVED*** return false ***REMOVED***

func (w WaitStatus) Stopped() bool ***REMOVED*** return false ***REMOVED***

func (w WaitStatus) Continued() bool ***REMOVED*** return false ***REMOVED***

func (w WaitStatus) StopSignal() Signal ***REMOVED*** return -1 ***REMOVED***

func (w WaitStatus) Signaled() bool ***REMOVED*** return false ***REMOVED***

func (w WaitStatus) TrapCause() int ***REMOVED*** return -1 ***REMOVED***

// Timespec is an invented structure on Windows, but here for
// consistency with the corresponding package for other operating systems.
type Timespec struct ***REMOVED***
	Sec  int64
	Nsec int64
***REMOVED***

func TimespecToNsec(ts Timespec) int64 ***REMOVED*** return int64(ts.Sec)*1e9 + int64(ts.Nsec) ***REMOVED***

func NsecToTimespec(nsec int64) (ts Timespec) ***REMOVED***
	ts.Sec = nsec / 1e9
	ts.Nsec = nsec % 1e9
	return
***REMOVED***

// TODO(brainman): fix all needed for net

func Accept(fd Handle) (nfd Handle, sa Sockaddr, err error) ***REMOVED*** return 0, nil, syscall.EWINDOWS ***REMOVED***
func Recvfrom(fd Handle, p []byte, flags int) (n int, from Sockaddr, err error) ***REMOVED***
	return 0, nil, syscall.EWINDOWS
***REMOVED***
func Sendto(fd Handle, p []byte, flags int, to Sockaddr) (err error)       ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func SetsockoptTimeval(fd Handle, level, opt int, tv *Timeval) (err error) ***REMOVED*** return syscall.EWINDOWS ***REMOVED***

// The Linger struct is wrong but we only noticed after Go 1.
// sysLinger is the real system call structure.

// BUG(brainman): The definition of Linger is not appropriate for direct use
// with Setsockopt and Getsockopt.
// Use SetsockoptLinger instead.

type Linger struct ***REMOVED***
	Onoff  int32
	Linger int32
***REMOVED***

type sysLinger struct ***REMOVED***
	Onoff  uint16
	Linger uint16
***REMOVED***

type IPMreq struct ***REMOVED***
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
***REMOVED***

type IPv6Mreq struct ***REMOVED***
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
***REMOVED***

func GetsockoptInt(fd Handle, level, opt int) (int, error) ***REMOVED*** return -1, syscall.EWINDOWS ***REMOVED***

func SetsockoptLinger(fd Handle, level, opt int, l *Linger) (err error) ***REMOVED***
	sys := sysLinger***REMOVED***Onoff: uint16(l.Onoff), Linger: uint16(l.Linger)***REMOVED***
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&sys)), int32(unsafe.Sizeof(sys)))
***REMOVED***

func SetsockoptInet4Addr(fd Handle, level, opt int, value [4]byte) (err error) ***REMOVED***
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&value[0])), 4)
***REMOVED***
func SetsockoptIPMreq(fd Handle, level, opt int, mreq *IPMreq) (err error) ***REMOVED***
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(mreq)), int32(unsafe.Sizeof(*mreq)))
***REMOVED***
func SetsockoptIPv6Mreq(fd Handle, level, opt int, mreq *IPv6Mreq) (err error) ***REMOVED***
	return syscall.EWINDOWS
***REMOVED***

func Getpid() (pid int) ***REMOVED*** return int(getCurrentProcessId()) ***REMOVED***

func FindFirstFile(name *uint16, data *Win32finddata) (handle Handle, err error) ***REMOVED***
	// NOTE(rsc): The Win32finddata struct is wrong for the system call:
	// the two paths are each one uint16 short. Use the correct struct,
	// a win32finddata1, and then copy the results out.
	// There is no loss of expressivity here, because the final
	// uint16, if it is used, is supposed to be a NUL, and Go doesn't need that.
	// For Go 1.1, we might avoid the allocation of win32finddata1 here
	// by adding a final Bug [2]uint16 field to the struct and then
	// adjusting the fields in the result directly.
	var data1 win32finddata1
	handle, err = findFirstFile1(name, &data1)
	if err == nil ***REMOVED***
		copyFindData(data, &data1)
	***REMOVED***
	return
***REMOVED***

func FindNextFile(handle Handle, data *Win32finddata) (err error) ***REMOVED***
	var data1 win32finddata1
	err = findNextFile1(handle, &data1)
	if err == nil ***REMOVED***
		copyFindData(data, &data1)
	***REMOVED***
	return
***REMOVED***

func getProcessEntry(pid int) (*ProcessEntry32, error) ***REMOVED***
	snapshot, err := CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer CloseHandle(snapshot)
	var procEntry ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = Process32First(snapshot, &procEntry); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	for ***REMOVED***
		if procEntry.ProcessID == uint32(pid) ***REMOVED***
			return &procEntry, nil
		***REMOVED***
		err = Process32Next(snapshot, &procEntry)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

func Getppid() (ppid int) ***REMOVED***
	pe, err := getProcessEntry(Getpid())
	if err != nil ***REMOVED***
		return -1
	***REMOVED***
	return int(pe.ParentProcessID)
***REMOVED***

// TODO(brainman): fix all needed for os
func Fchdir(fd Handle) (err error)             ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func Link(oldpath, newpath string) (err error) ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func Symlink(path, link string) (err error)    ***REMOVED*** return syscall.EWINDOWS ***REMOVED***

func Fchmod(fd Handle, mode uint32) (err error)        ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func Chown(path string, uid int, gid int) (err error)  ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func Lchown(path string, uid int, gid int) (err error) ***REMOVED*** return syscall.EWINDOWS ***REMOVED***
func Fchown(fd Handle, uid int, gid int) (err error)   ***REMOVED*** return syscall.EWINDOWS ***REMOVED***

func Getuid() (uid int)                  ***REMOVED*** return -1 ***REMOVED***
func Geteuid() (euid int)                ***REMOVED*** return -1 ***REMOVED***
func Getgid() (gid int)                  ***REMOVED*** return -1 ***REMOVED***
func Getegid() (egid int)                ***REMOVED*** return -1 ***REMOVED***
func Getgroups() (gids []int, err error) ***REMOVED*** return nil, syscall.EWINDOWS ***REMOVED***

type Signal int

func (s Signal) Signal() ***REMOVED******REMOVED***

func (s Signal) String() string ***REMOVED***
	if 0 <= s && int(s) < len(signals) ***REMOVED***
		str := signals[s]
		if str != "" ***REMOVED***
			return str
		***REMOVED***
	***REMOVED***
	return "signal " + itoa(int(s))
***REMOVED***

func LoadCreateSymbolicLink() error ***REMOVED***
	return procCreateSymbolicLinkW.Find()
***REMOVED***

// Readlink returns the destination of the named symbolic link.
func Readlink(path string, buf []byte) (n int, err error) ***REMOVED***
	fd, err := CreateFile(StringToUTF16Ptr(path), GENERIC_READ, 0, nil, OPEN_EXISTING,
		FILE_FLAG_OPEN_REPARSE_POINT|FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***
	defer CloseHandle(fd)

	rdbbuf := make([]byte, MAXIMUM_REPARSE_DATA_BUFFER_SIZE)
	var bytesReturned uint32
	err = DeviceIoControl(fd, FSCTL_GET_REPARSE_POINT, nil, 0, &rdbbuf[0], uint32(len(rdbbuf)), &bytesReturned, nil)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	rdb := (*reparseDataBuffer)(unsafe.Pointer(&rdbbuf[0]))
	var s string
	switch rdb.ReparseTag ***REMOVED***
	case IO_REPARSE_TAG_SYMLINK:
		data := (*symbolicLinkReparseBuffer)(unsafe.Pointer(&rdb.reparseBuffer))
		p := (*[0xffff]uint16)(unsafe.Pointer(&data.PathBuffer[0]))
		s = UTF16ToString(p[data.PrintNameOffset/2 : (data.PrintNameLength-data.PrintNameOffset)/2])
	case IO_REPARSE_TAG_MOUNT_POINT:
		data := (*mountPointReparseBuffer)(unsafe.Pointer(&rdb.reparseBuffer))
		p := (*[0xffff]uint16)(unsafe.Pointer(&data.PathBuffer[0]))
		s = UTF16ToString(p[data.PrintNameOffset/2 : (data.PrintNameLength-data.PrintNameOffset)/2])
	default:
		// the path is not a symlink or junction but another type of reparse
		// point
		return -1, syscall.ENOENT
	***REMOVED***
	n = copy(buf, []byte(s))

	return n, nil
***REMOVED***
