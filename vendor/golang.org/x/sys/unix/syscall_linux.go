// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Linux system calls.
// This file is compiled as ordinary Go code,
// but it is also input to mksyscall,
// which parses the //sys lines and generates system call stubs.
// Note that sometimes we use a lowercase //sys name and
// wrap it in our own nicer implementation.

package unix

import (
	"syscall"
	"unsafe"
)

/*
 * Wrapped
 */

func Access(path string, mode uint32) (err error) ***REMOVED***
	return Faccessat(AT_FDCWD, path, mode, 0)
***REMOVED***

func Chmod(path string, mode uint32) (err error) ***REMOVED***
	return Fchmodat(AT_FDCWD, path, mode, 0)
***REMOVED***

func Chown(path string, uid int, gid int) (err error) ***REMOVED***
	return Fchownat(AT_FDCWD, path, uid, gid, 0)
***REMOVED***

func Creat(path string, mode uint32) (fd int, err error) ***REMOVED***
	return Open(path, O_CREAT|O_WRONLY|O_TRUNC, mode)
***REMOVED***

//sys	fchmodat(dirfd int, path string, mode uint32) (err error)

func Fchmodat(dirfd int, path string, mode uint32, flags int) (err error) ***REMOVED***
	// Linux fchmodat doesn't support the flags parameter. Mimick glibc's behavior
	// and check the flags. Otherwise the mode would be applied to the symlink
	// destination which is not what the user expects.
	if flags&^AT_SYMLINK_NOFOLLOW != 0 ***REMOVED***
		return EINVAL
	***REMOVED*** else if flags&AT_SYMLINK_NOFOLLOW != 0 ***REMOVED***
		return EOPNOTSUPP
	***REMOVED***
	return fchmodat(dirfd, path, mode)
***REMOVED***

//sys	ioctl(fd int, req uint, arg uintptr) (err error)

// ioctl itself should not be exposed directly, but additional get/set
// functions for specific types are permissible.

// IoctlSetInt performs an ioctl operation which sets an integer value
// on fd, using the specified request number.
func IoctlSetInt(fd int, req uint, value int) error ***REMOVED***
	return ioctl(fd, req, uintptr(value))
***REMOVED***

func IoctlSetWinsize(fd int, req uint, value *Winsize) error ***REMOVED***
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
***REMOVED***

func IoctlSetTermios(fd int, req uint, value *Termios) error ***REMOVED***
	return ioctl(fd, req, uintptr(unsafe.Pointer(value)))
***REMOVED***

// IoctlGetInt performs an ioctl operation which gets an integer value
// from fd, using the specified request number.
func IoctlGetInt(fd int, req uint) (int, error) ***REMOVED***
	var value int
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return value, err
***REMOVED***

func IoctlGetWinsize(fd int, req uint) (*Winsize, error) ***REMOVED***
	var value Winsize
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
***REMOVED***

func IoctlGetTermios(fd int, req uint) (*Termios, error) ***REMOVED***
	var value Termios
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&value)))
	return &value, err
***REMOVED***

//sys	Linkat(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) (err error)

func Link(oldpath string, newpath string) (err error) ***REMOVED***
	return Linkat(AT_FDCWD, oldpath, AT_FDCWD, newpath, 0)
***REMOVED***

func Mkdir(path string, mode uint32) (err error) ***REMOVED***
	return Mkdirat(AT_FDCWD, path, mode)
***REMOVED***

func Mknod(path string, mode uint32, dev int) (err error) ***REMOVED***
	return Mknodat(AT_FDCWD, path, mode, dev)
***REMOVED***

func Open(path string, mode int, perm uint32) (fd int, err error) ***REMOVED***
	return openat(AT_FDCWD, path, mode|O_LARGEFILE, perm)
***REMOVED***

//sys	openat(dirfd int, path string, flags int, mode uint32) (fd int, err error)

func Openat(dirfd int, path string, flags int, mode uint32) (fd int, err error) ***REMOVED***
	return openat(dirfd, path, flags|O_LARGEFILE, mode)
***REMOVED***

//sys	ppoll(fds *PollFd, nfds int, timeout *Timespec, sigmask *Sigset_t) (n int, err error)

func Ppoll(fds []PollFd, timeout *Timespec, sigmask *Sigset_t) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return ppoll(nil, 0, timeout, sigmask)
	***REMOVED***
	return ppoll(&fds[0], len(fds), timeout, sigmask)
***REMOVED***

//sys	Readlinkat(dirfd int, path string, buf []byte) (n int, err error)

func Readlink(path string, buf []byte) (n int, err error) ***REMOVED***
	return Readlinkat(AT_FDCWD, path, buf)
***REMOVED***

func Rename(oldpath string, newpath string) (err error) ***REMOVED***
	return Renameat(AT_FDCWD, oldpath, AT_FDCWD, newpath)
***REMOVED***

func Rmdir(path string) error ***REMOVED***
	return Unlinkat(AT_FDCWD, path, AT_REMOVEDIR)
***REMOVED***

//sys	Symlinkat(oldpath string, newdirfd int, newpath string) (err error)

func Symlink(oldpath string, newpath string) (err error) ***REMOVED***
	return Symlinkat(oldpath, AT_FDCWD, newpath)
***REMOVED***

func Unlink(path string) error ***REMOVED***
	return Unlinkat(AT_FDCWD, path, 0)
***REMOVED***

//sys	Unlinkat(dirfd int, path string, flags int) (err error)

//sys	utimes(path string, times *[2]Timeval) (err error)

func Utimes(path string, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		err := utimensat(AT_FDCWD, path, nil, 0)
		if err != ENOSYS ***REMOVED***
			return err
		***REMOVED***
		return utimes(path, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var ts [2]Timespec
	ts[0] = NsecToTimespec(TimevalToNsec(tv[0]))
	ts[1] = NsecToTimespec(TimevalToNsec(tv[1]))
	err := utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

//sys	utimensat(dirfd int, path string, times *[2]Timespec, flags int) (err error)

func UtimesNano(path string, ts []Timespec) error ***REMOVED***
	if ts == nil ***REMOVED***
		err := utimensat(AT_FDCWD, path, nil, 0)
		if err != ENOSYS ***REMOVED***
			return err
		***REMOVED***
		return utimes(path, nil)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	err := utimensat(AT_FDCWD, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), 0)
	if err != ENOSYS ***REMOVED***
		return err
	***REMOVED***
	// If the utimensat syscall isn't available (utimensat was added to Linux
	// in 2.6.22, Released, 8 July 2007) then fall back to utimes
	var tv [2]Timeval
	for i := 0; i < 2; i++ ***REMOVED***
		tv[i] = NsecToTimeval(TimespecToNsec(ts[i]))
	***REMOVED***
	return utimes(path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

func UtimesNanoAt(dirfd int, path string, ts []Timespec, flags int) error ***REMOVED***
	if ts == nil ***REMOVED***
		return utimensat(dirfd, path, nil, flags)
	***REMOVED***
	if len(ts) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return utimensat(dirfd, path, (*[2]Timespec)(unsafe.Pointer(&ts[0])), flags)
***REMOVED***

//sys	futimesat(dirfd int, path *byte, times *[2]Timeval) (err error)

func Futimesat(dirfd int, path string, tv []Timeval) error ***REMOVED***
	pathp, err := BytePtrFromString(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if tv == nil ***REMOVED***
		return futimesat(dirfd, pathp, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return futimesat(dirfd, pathp, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
***REMOVED***

func Futimes(fd int, tv []Timeval) (err error) ***REMOVED***
	// Believe it or not, this is the best we can do on Linux
	// (and is what glibc does).
	return Utimes("/proc/self/fd/"+itoa(fd), tv)
***REMOVED***

const ImplementsGetwd = true

//sys	Getcwd(buf []byte) (n int, err error)

func Getwd() (wd string, err error) ***REMOVED***
	var buf [PathMax]byte
	n, err := Getcwd(buf[0:])
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// Getcwd returns the number of bytes written to buf, including the NUL.
	if n < 1 || n > len(buf) || buf[n-1] != 0 ***REMOVED***
		return "", EINVAL
	***REMOVED***
	return string(buf[0 : n-1]), nil
***REMOVED***

func Getgroups() (gids []int, err error) ***REMOVED***
	n, err := getgroups(0, nil)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if n == 0 ***REMOVED***
		return nil, nil
	***REMOVED***

	// Sanity check group count. Max is 1<<16 on Linux.
	if n < 0 || n > 1<<20 ***REMOVED***
		return nil, EINVAL
	***REMOVED***

	a := make([]_Gid_t, n)
	n, err = getgroups(n, &a[0])
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	gids = make([]int, n)
	for i, v := range a[0:n] ***REMOVED***
		gids[i] = int(v)
	***REMOVED***
	return
***REMOVED***

func Setgroups(gids []int) (err error) ***REMOVED***
	if len(gids) == 0 ***REMOVED***
		return setgroups(0, nil)
	***REMOVED***

	a := make([]_Gid_t, len(gids))
	for i, v := range gids ***REMOVED***
		a[i] = _Gid_t(v)
	***REMOVED***
	return setgroups(len(a), &a[0])
***REMOVED***

type WaitStatus uint32

// Wait status is 7 bits at bottom, either 0 (exited),
// 0x7F (stopped), or a signal number that caused an exit.
// The 0x80 bit is whether there was a core dump.
// An extra number (exit code, signal causing a stop)
// is in the high bits. At least that's the idea.
// There are various irregularities. For example, the
// "continued" status is 0xFFFF, distinguishing itself
// from stopped via the core dump bit.

const (
	mask    = 0x7F
	core    = 0x80
	exited  = 0x00
	stopped = 0x7F
	shift   = 8
)

func (w WaitStatus) Exited() bool ***REMOVED*** return w&mask == exited ***REMOVED***

func (w WaitStatus) Signaled() bool ***REMOVED*** return w&mask != stopped && w&mask != exited ***REMOVED***

func (w WaitStatus) Stopped() bool ***REMOVED*** return w&0xFF == stopped ***REMOVED***

func (w WaitStatus) Continued() bool ***REMOVED*** return w == 0xFFFF ***REMOVED***

func (w WaitStatus) CoreDump() bool ***REMOVED*** return w.Signaled() && w&core != 0 ***REMOVED***

func (w WaitStatus) ExitStatus() int ***REMOVED***
	if !w.Exited() ***REMOVED***
		return -1
	***REMOVED***
	return int(w>>shift) & 0xFF
***REMOVED***

func (w WaitStatus) Signal() syscall.Signal ***REMOVED***
	if !w.Signaled() ***REMOVED***
		return -1
	***REMOVED***
	return syscall.Signal(w & mask)
***REMOVED***

func (w WaitStatus) StopSignal() syscall.Signal ***REMOVED***
	if !w.Stopped() ***REMOVED***
		return -1
	***REMOVED***
	return syscall.Signal(w>>shift) & 0xFF
***REMOVED***

func (w WaitStatus) TrapCause() int ***REMOVED***
	if w.StopSignal() != SIGTRAP ***REMOVED***
		return -1
	***REMOVED***
	return int(w>>shift) >> 8
***REMOVED***

//sys	wait4(pid int, wstatus *_C_int, options int, rusage *Rusage) (wpid int, err error)

func Wait4(pid int, wstatus *WaitStatus, options int, rusage *Rusage) (wpid int, err error) ***REMOVED***
	var status _C_int
	wpid, err = wait4(pid, &status, options, rusage)
	if wstatus != nil ***REMOVED***
		*wstatus = WaitStatus(status)
	***REMOVED***
	return
***REMOVED***

func Mkfifo(path string, mode uint32) error ***REMOVED***
	return Mknod(path, mode|S_IFIFO, 0)
***REMOVED***

func Mkfifoat(dirfd int, path string, mode uint32) error ***REMOVED***
	return Mknodat(dirfd, path, mode|S_IFIFO, 0)
***REMOVED***

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet4, nil
***REMOVED***

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Port < 0 || sa.Port > 0xFFFF ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrInet6, nil
***REMOVED***

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	name := sa.Name
	n := len(name)
	if n >= len(sa.raw.Path) ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ ***REMOVED***
		sa.raw.Path[i] = int8(name[i])
	***REMOVED***
	// length is family (uint16), name, NUL.
	sl := _Socklen(2)
	if n > 0 ***REMOVED***
		sl += _Socklen(n) + 1
	***REMOVED***
	if sa.raw.Path[0] == '@' ***REMOVED***
		sa.raw.Path[0] = 0
		// Don't count trailing NUL for abstract address.
		sl--
	***REMOVED***

	return unsafe.Pointer(&sa.raw), sl, nil
***REMOVED***

type SockaddrLinklayer struct ***REMOVED***
	Protocol uint16
	Ifindex  int
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]byte
	raw      RawSockaddrLinklayer
***REMOVED***

func (sa *SockaddrLinklayer) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_PACKET
	sa.raw.Protocol = sa.Protocol
	sa.raw.Ifindex = int32(sa.Ifindex)
	sa.raw.Hatype = sa.Hatype
	sa.raw.Pkttype = sa.Pkttype
	sa.raw.Halen = sa.Halen
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Addr[i] = sa.Addr[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrLinklayer, nil
***REMOVED***

type SockaddrNetlink struct ***REMOVED***
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
	raw    RawSockaddrNetlink
***REMOVED***

func (sa *SockaddrNetlink) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_NETLINK
	sa.raw.Pad = sa.Pad
	sa.raw.Pid = sa.Pid
	sa.raw.Groups = sa.Groups
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNetlink, nil
***REMOVED***

type SockaddrHCI struct ***REMOVED***
	Dev     uint16
	Channel uint16
	raw     RawSockaddrHCI
***REMOVED***

func (sa *SockaddrHCI) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_BLUETOOTH
	sa.raw.Dev = sa.Dev
	sa.raw.Channel = sa.Channel
	return unsafe.Pointer(&sa.raw), SizeofSockaddrHCI, nil
***REMOVED***

// SockaddrCAN implements the Sockaddr interface for AF_CAN type sockets.
// The RxID and TxID fields are used for transport protocol addressing in
// (CAN_TP16, CAN_TP20, CAN_MCNET, and CAN_ISOTP), they can be left with
// zero values for CAN_RAW and CAN_BCM sockets as they have no meaning.
//
// The SockaddrCAN struct must be bound to the socket file descriptor
// using Bind before the CAN socket can be used.
//
//      // Read one raw CAN frame
//      fd, _ := Socket(AF_CAN, SOCK_RAW, CAN_RAW)
//      addr := &SockaddrCAN***REMOVED***Ifindex: index***REMOVED***
//      Bind(fd, addr)
//      frame := make([]byte, 16)
//      Read(fd, frame)
//
// The full SocketCAN documentation can be found in the linux kernel
// archives at: https://www.kernel.org/doc/Documentation/networking/can.txt
type SockaddrCAN struct ***REMOVED***
	Ifindex int
	RxID    uint32
	TxID    uint32
	raw     RawSockaddrCAN
***REMOVED***

func (sa *SockaddrCAN) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_CAN
	sa.raw.Ifindex = int32(sa.Ifindex)
	rx := (*[4]byte)(unsafe.Pointer(&sa.RxID))
	for i := 0; i < 4; i++ ***REMOVED***
		sa.raw.Addr[i] = rx[i]
	***REMOVED***
	tx := (*[4]byte)(unsafe.Pointer(&sa.TxID))
	for i := 0; i < 4; i++ ***REMOVED***
		sa.raw.Addr[i+4] = tx[i]
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrCAN, nil
***REMOVED***

// SockaddrALG implements the Sockaddr interface for AF_ALG type sockets.
// SockaddrALG enables userspace access to the Linux kernel's cryptography
// subsystem. The Type and Name fields specify which type of hash or cipher
// should be used with a given socket.
//
// To create a file descriptor that provides access to a hash or cipher, both
// Bind and Accept must be used. Once the setup process is complete, input
// data can be written to the socket, processed by the kernel, and then read
// back as hash output or ciphertext.
//
// Here is an example of using an AF_ALG socket with SHA1 hashing.
// The initial socket setup process is as follows:
//
//      // Open a socket to perform SHA1 hashing.
//      fd, _ := unix.Socket(unix.AF_ALG, unix.SOCK_SEQPACKET, 0)
//      addr := &unix.SockaddrALG***REMOVED***Type: "hash", Name: "sha1"***REMOVED***
//      unix.Bind(fd, addr)
//      // Note: unix.Accept does not work at this time; must invoke accept()
//      // manually using unix.Syscall.
//      hashfd, _, _ := unix.Syscall(unix.SYS_ACCEPT, uintptr(fd), 0, 0)
//
// Once a file descriptor has been returned from Accept, it may be used to
// perform SHA1 hashing. The descriptor is not safe for concurrent use, but
// may be re-used repeatedly with subsequent Write and Read operations.
//
// When hashing a small byte slice or string, a single Write and Read may
// be used:
//
//      // Assume hashfd is already configured using the setup process.
//      hash := os.NewFile(hashfd, "sha1")
//      // Hash an input string and read the results. Each Write discards
//      // previous hash state. Read always reads the current state.
//      b := make([]byte, 20)
//      for i := 0; i < 2; i++ ***REMOVED***
//          io.WriteString(hash, "Hello, world.")
//          hash.Read(b)
//          fmt.Println(hex.EncodeToString(b))
//      ***REMOVED***
//      // Output:
//      // 2ae01472317d1935a84797ec1983ae243fc6aa28
//      // 2ae01472317d1935a84797ec1983ae243fc6aa28
//
// For hashing larger byte slices, or byte streams such as those read from
// a file or socket, use Sendto with MSG_MORE to instruct the kernel to update
// the hash digest instead of creating a new one for a given chunk and finalizing it.
//
//      // Assume hashfd and addr are already configured using the setup process.
//      hash := os.NewFile(hashfd, "sha1")
//      // Hash the contents of a file.
//      f, _ := os.Open("/tmp/linux-4.10-rc7.tar.xz")
//      b := make([]byte, 4096)
//      for ***REMOVED***
//          n, err := f.Read(b)
//          if err == io.EOF ***REMOVED***
//              break
//          ***REMOVED***
//          unix.Sendto(hashfd, b[:n], unix.MSG_MORE, addr)
//      ***REMOVED***
//      hash.Read(b)
//      fmt.Println(hex.EncodeToString(b))
//      // Output: 85cdcad0c06eef66f805ecce353bec9accbeecc5
//
// For more information, see: http://www.chronox.de/crypto-API/crypto/userspace-if.html.
type SockaddrALG struct ***REMOVED***
	Type    string
	Name    string
	Feature uint32
	Mask    uint32
	raw     RawSockaddrALG
***REMOVED***

func (sa *SockaddrALG) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	// Leave room for NUL byte terminator.
	if len(sa.Type) > 13 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	if len(sa.Name) > 63 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***

	sa.raw.Family = AF_ALG
	sa.raw.Feat = sa.Feature
	sa.raw.Mask = sa.Mask

	typ, err := ByteSliceFromString(sa.Type)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***
	name, err := ByteSliceFromString(sa.Name)
	if err != nil ***REMOVED***
		return nil, 0, err
	***REMOVED***

	copy(sa.raw.Type[:], typ)
	copy(sa.raw.Name[:], name)

	return unsafe.Pointer(&sa.raw), SizeofSockaddrALG, nil
***REMOVED***

// SockaddrVM implements the Sockaddr interface for AF_VSOCK type sockets.
// SockaddrVM provides access to Linux VM sockets: a mechanism that enables
// bidirectional communication between a hypervisor and its guest virtual
// machines.
type SockaddrVM struct ***REMOVED***
	// CID and Port specify a context ID and port address for a VM socket.
	// Guests have a unique CID, and hosts may have a well-known CID of:
	//  - VMADDR_CID_HYPERVISOR: refers to the hypervisor process.
	//  - VMADDR_CID_HOST: refers to other processes on the host.
	CID  uint32
	Port uint32
	raw  RawSockaddrVM
***REMOVED***

func (sa *SockaddrVM) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_VSOCK
	sa.raw.Port = sa.Port
	sa.raw.Cid = sa.CID

	return unsafe.Pointer(&sa.raw), SizeofSockaddrVM, nil
***REMOVED***

func anyToSockaddr(rsa *RawSockaddrAny) (Sockaddr, error) ***REMOVED***
	switch rsa.Addr.Family ***REMOVED***
	case AF_NETLINK:
		pp := (*RawSockaddrNetlink)(unsafe.Pointer(rsa))
		sa := new(SockaddrNetlink)
		sa.Family = pp.Family
		sa.Pad = pp.Pad
		sa.Pid = pp.Pid
		sa.Groups = pp.Groups
		return sa, nil

	case AF_PACKET:
		pp := (*RawSockaddrLinklayer)(unsafe.Pointer(rsa))
		sa := new(SockaddrLinklayer)
		sa.Protocol = pp.Protocol
		sa.Ifindex = int(pp.Ifindex)
		sa.Hatype = pp.Hatype
		sa.Pkttype = pp.Pkttype
		sa.Halen = pp.Halen
		for i := 0; i < len(sa.Addr); i++ ***REMOVED***
			sa.Addr[i] = pp.Addr[i]
		***REMOVED***
		return sa, nil

	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		if pp.Path[0] == 0 ***REMOVED***
			// "Abstract" Unix domain socket.
			// Rewrite leading NUL as @ for textual display.
			// (This is the standard convention.)
			// Not friendly to overwrite in place,
			// but the callers below don't care.
			pp.Path[0] = '@'
		***REMOVED***

		// Assume path ends at NUL.
		// This is not technically the Linux semantics for
		// abstract Unix domain sockets--they are supposed
		// to be uninterpreted fixed-size binary blobs--but
		// everyone uses this convention.
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 ***REMOVED***
			n++
		***REMOVED***
		bytes := (*[10000]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, nil

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

	case AF_VSOCK:
		pp := (*RawSockaddrVM)(unsafe.Pointer(rsa))
		sa := &SockaddrVM***REMOVED***
			CID:  pp.Cid,
			Port: pp.Port,
		***REMOVED***
		return sa, nil
	***REMOVED***
	return nil, EAFNOSUPPORT
***REMOVED***

func Accept(fd int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept(fd, &rsa, &len)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	sa, err = anyToSockaddr(&rsa)
	if err != nil ***REMOVED***
		Close(nfd)
		nfd = 0
	***REMOVED***
	return
***REMOVED***

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	if len > SizeofSockaddrAny ***REMOVED***
		panic("RawSockaddrAny too small")
	***REMOVED***
	sa, err = anyToSockaddr(&rsa)
	if err != nil ***REMOVED***
		Close(nfd)
		nfd = 0
	***REMOVED***
	return
***REMOVED***

func Getsockname(fd int) (sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	if err = getsockname(fd, &rsa, &len); err != nil ***REMOVED***
		return
	***REMOVED***
	return anyToSockaddr(&rsa)
***REMOVED***

func GetsockoptInet4Addr(fd, level, opt int) (value [4]byte, err error) ***REMOVED***
	vallen := _Socklen(4)
	err = getsockopt(fd, level, opt, unsafe.Pointer(&value[0]), &vallen)
	return value, err
***REMOVED***

func GetsockoptIPMreq(fd, level, opt int) (*IPMreq, error) ***REMOVED***
	var value IPMreq
	vallen := _Socklen(SizeofIPMreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) ***REMOVED***
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptIPv6Mreq(fd, level, opt int) (*IPv6Mreq, error) ***REMOVED***
	var value IPv6Mreq
	vallen := _Socklen(SizeofIPv6Mreq)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptIPv6MTUInfo(fd, level, opt int) (*IPv6MTUInfo, error) ***REMOVED***
	var value IPv6MTUInfo
	vallen := _Socklen(SizeofIPv6MTUInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptICMPv6Filter(fd, level, opt int) (*ICMPv6Filter, error) ***REMOVED***
	var value ICMPv6Filter
	vallen := _Socklen(SizeofICMPv6Filter)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptUcred(fd, level, opt int) (*Ucred, error) ***REMOVED***
	var value Ucred
	vallen := _Socklen(SizeofUcred)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptTCPInfo(fd, level, opt int) (*TCPInfo, error) ***REMOVED***
	var value TCPInfo
	vallen := _Socklen(SizeofTCPInfo)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
***REMOVED***

// Keyctl Commands (http://man7.org/linux/man-pages/man2/keyctl.2.html)

// KeyctlInt calls keyctl commands in which each argument is an int.
// These commands are KEYCTL_REVOKE, KEYCTL_CHOWN, KEYCTL_CLEAR, KEYCTL_LINK,
// KEYCTL_UNLINK, KEYCTL_NEGATE, KEYCTL_SET_REQKEY_KEYRING, KEYCTL_SET_TIMEOUT,
// KEYCTL_ASSUME_AUTHORITY, KEYCTL_SESSION_TO_PARENT, KEYCTL_REJECT,
// KEYCTL_INVALIDATE, and KEYCTL_GET_PERSISTENT.
//sys	KeyctlInt(cmd int, arg2 int, arg3 int, arg4 int, arg5 int) (ret int, err error) = SYS_KEYCTL

// KeyctlBuffer calls keyctl commands in which the third and fourth
// arguments are a buffer and its length, respectively.
// These commands are KEYCTL_UPDATE, KEYCTL_READ, and KEYCTL_INSTANTIATE.
//sys	KeyctlBuffer(cmd int, arg2 int, buf []byte, arg5 int) (ret int, err error) = SYS_KEYCTL

// KeyctlString calls keyctl commands which return a string.
// These commands are KEYCTL_DESCRIBE and KEYCTL_GET_SECURITY.
func KeyctlString(cmd int, id int) (string, error) ***REMOVED***
	// We must loop as the string data may change in between the syscalls.
	// We could allocate a large buffer here to reduce the chance that the
	// syscall needs to be called twice; however, this is unnecessary as
	// the performance loss is negligible.
	var buffer []byte
	for ***REMOVED***
		// Try to fill the buffer with data
		length, err := KeyctlBuffer(cmd, id, buffer, 0)
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***

		// Check if the data was written
		if length <= len(buffer) ***REMOVED***
			// Exclude the null terminator
			return string(buffer[:length-1]), nil
		***REMOVED***

		// Make a bigger buffer if needed
		buffer = make([]byte, length)
	***REMOVED***
***REMOVED***

// Keyctl commands with special signatures.

// KeyctlGetKeyringID implements the KEYCTL_GET_KEYRING_ID command.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_get_keyring_ID.3.html
func KeyctlGetKeyringID(id int, create bool) (ringid int, err error) ***REMOVED***
	createInt := 0
	if create ***REMOVED***
		createInt = 1
	***REMOVED***
	return KeyctlInt(KEYCTL_GET_KEYRING_ID, id, createInt, 0, 0)
***REMOVED***

// KeyctlSetperm implements the KEYCTL_SETPERM command. The perm value is the
// key handle permission mask as described in the "keyctl setperm" section of
// http://man7.org/linux/man-pages/man1/keyctl.1.html.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_setperm.3.html
func KeyctlSetperm(id int, perm uint32) error ***REMOVED***
	_, err := KeyctlInt(KEYCTL_SETPERM, id, int(perm), 0, 0)
	return err
***REMOVED***

//sys	keyctlJoin(cmd int, arg2 string) (ret int, err error) = SYS_KEYCTL

// KeyctlJoinSessionKeyring implements the KEYCTL_JOIN_SESSION_KEYRING command.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_join_session_keyring.3.html
func KeyctlJoinSessionKeyring(name string) (ringid int, err error) ***REMOVED***
	return keyctlJoin(KEYCTL_JOIN_SESSION_KEYRING, name)
***REMOVED***

//sys	keyctlSearch(cmd int, arg2 int, arg3 string, arg4 string, arg5 int) (ret int, err error) = SYS_KEYCTL

// KeyctlSearch implements the KEYCTL_SEARCH command.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_search.3.html
func KeyctlSearch(ringid int, keyType, description string, destRingid int) (id int, err error) ***REMOVED***
	return keyctlSearch(KEYCTL_SEARCH, ringid, keyType, description, destRingid)
***REMOVED***

//sys	keyctlIOV(cmd int, arg2 int, payload []Iovec, arg5 int) (err error) = SYS_KEYCTL

// KeyctlInstantiateIOV implements the KEYCTL_INSTANTIATE_IOV command. This
// command is similar to KEYCTL_INSTANTIATE, except that the payload is a slice
// of Iovec (each of which represents a buffer) instead of a single buffer.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_instantiate_iov.3.html
func KeyctlInstantiateIOV(id int, payload []Iovec, ringid int) error ***REMOVED***
	return keyctlIOV(KEYCTL_INSTANTIATE_IOV, id, payload, ringid)
***REMOVED***

//sys	keyctlDH(cmd int, arg2 *KeyctlDHParams, buf []byte) (ret int, err error) = SYS_KEYCTL

// KeyctlDHCompute implements the KEYCTL_DH_COMPUTE command. This command
// computes a Diffie-Hellman shared secret based on the provide params. The
// secret is written to the provided buffer and the returned size is the number
// of bytes written (returning an error if there is insufficient space in the
// buffer). If a nil buffer is passed in, this function returns the minimum
// buffer length needed to store the appropriate data. Note that this differs
// from KEYCTL_READ's behavior which always returns the requested payload size.
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_dh_compute.3.html
func KeyctlDHCompute(params *KeyctlDHParams, buffer []byte) (size int, err error) ***REMOVED***
	return keyctlDH(KEYCTL_DH_COMPUTE, params, buffer)
***REMOVED***

func Recvmsg(fd int, p, oob []byte, flags int) (n, oobn int, recvflags int, from Sockaddr, err error) ***REMOVED***
	var msg Msghdr
	var rsa RawSockaddrAny
	msg.Name = (*byte)(unsafe.Pointer(&rsa))
	msg.Namelen = uint32(SizeofSockaddrAny)
	var iov Iovec
	if len(p) > 0 ***REMOVED***
		iov.Base = &p[0]
		iov.SetLen(len(p))
	***REMOVED***
	var dummy byte
	if len(oob) > 0 ***REMOVED***
		var sockType int
		sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
		if err != nil ***REMOVED***
			return
		***REMOVED***
		// receive at least one normal byte
		if sockType != SOCK_DGRAM && len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	***REMOVED***
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = recvmsg(fd, &msg, flags); err != nil ***REMOVED***
		return
	***REMOVED***
	oobn = int(msg.Controllen)
	recvflags = int(msg.Flags)
	// source address is only specified if the socket is unconnected
	if rsa.Addr.Family != AF_UNSPEC ***REMOVED***
		from, err = anyToSockaddr(&rsa)
	***REMOVED***
	return
***REMOVED***

func Sendmsg(fd int, p, oob []byte, to Sockaddr, flags int) (err error) ***REMOVED***
	_, err = SendmsgN(fd, p, oob, to, flags)
	return
***REMOVED***

func SendmsgN(fd int, p, oob []byte, to Sockaddr, flags int) (n int, err error) ***REMOVED***
	var ptr unsafe.Pointer
	var salen _Socklen
	if to != nil ***REMOVED***
		var err error
		ptr, salen, err = to.sockaddr()
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
	***REMOVED***
	var msg Msghdr
	msg.Name = (*byte)(ptr)
	msg.Namelen = uint32(salen)
	var iov Iovec
	if len(p) > 0 ***REMOVED***
		iov.Base = &p[0]
		iov.SetLen(len(p))
	***REMOVED***
	var dummy byte
	if len(oob) > 0 ***REMOVED***
		var sockType int
		sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		// send at least one normal byte
		if sockType != SOCK_DGRAM && len(p) == 0 ***REMOVED***
			iov.Base = &dummy
			iov.SetLen(1)
		***REMOVED***
		msg.Control = &oob[0]
		msg.SetControllen(len(oob))
	***REMOVED***
	msg.Iov = &iov
	msg.Iovlen = 1
	if n, err = sendmsg(fd, &msg, flags); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	if len(oob) > 0 && len(p) == 0 ***REMOVED***
		n = 0
	***REMOVED***
	return n, nil
***REMOVED***

// BindToDevice binds the socket associated with fd to device.
func BindToDevice(fd int, device string) (err error) ***REMOVED***
	return SetsockoptString(fd, SOL_SOCKET, SO_BINDTODEVICE, device)
***REMOVED***

//sys	ptrace(request int, pid int, addr uintptr, data uintptr) (err error)

func ptracePeek(req int, pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	// The peek requests are machine-size oriented, so we wrap it
	// to retrieve arbitrary-length data.

	// The ptrace syscall differs from glibc's ptrace.
	// Peeks returns the word in *data, not as the return value.

	var buf [sizeofPtr]byte

	// Leading edge. PEEKTEXT/PEEKDATA don't require aligned
	// access (PEEKUSER warns that it might), but if we don't
	// align our reads, we might straddle an unmapped page
	// boundary and not get the bytes leading up to the page
	// boundary.
	n := 0
	if addr%sizeofPtr != 0 ***REMOVED***
		err = ptrace(req, pid, addr-addr%sizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		n += copy(out, buf[addr%sizeofPtr:])
		out = out[n:]
	***REMOVED***

	// Remainder.
	for len(out) > 0 ***REMOVED***
		// We use an internal buffer to guarantee alignment.
		// It's not documented if this is necessary, but we're paranoid.
		err = ptrace(req, pid, addr+uintptr(n), uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		copied := copy(out, buf[0:])
		n += copied
		out = out[copied:]
	***REMOVED***

	return n, nil
***REMOVED***

func PtracePeekText(pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	return ptracePeek(PTRACE_PEEKTEXT, pid, addr, out)
***REMOVED***

func PtracePeekData(pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	return ptracePeek(PTRACE_PEEKDATA, pid, addr, out)
***REMOVED***

func PtracePeekUser(pid int, addr uintptr, out []byte) (count int, err error) ***REMOVED***
	return ptracePeek(PTRACE_PEEKUSR, pid, addr, out)
***REMOVED***

func ptracePoke(pokeReq int, peekReq int, pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	// As for ptracePeek, we need to align our accesses to deal
	// with the possibility of straddling an invalid page.

	// Leading edge.
	n := 0
	if addr%sizeofPtr != 0 ***REMOVED***
		var buf [sizeofPtr]byte
		err = ptrace(peekReq, pid, addr-addr%sizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		n += copy(buf[addr%sizeofPtr:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr-addr%sizeofPtr, word)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		data = data[n:]
	***REMOVED***

	// Interior.
	for len(data) > sizeofPtr ***REMOVED***
		word := *((*uintptr)(unsafe.Pointer(&data[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		n += sizeofPtr
		data = data[sizeofPtr:]
	***REMOVED***

	// Trailing edge.
	if len(data) > 0 ***REMOVED***
		var buf [sizeofPtr]byte
		err = ptrace(peekReq, pid, addr+uintptr(n), uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		copy(buf[0:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		n += len(data)
	***REMOVED***

	return n, nil
***REMOVED***

func PtracePokeText(pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	return ptracePoke(PTRACE_POKETEXT, PTRACE_PEEKTEXT, pid, addr, data)
***REMOVED***

func PtracePokeData(pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	return ptracePoke(PTRACE_POKEDATA, PTRACE_PEEKDATA, pid, addr, data)
***REMOVED***

func PtracePokeUser(pid int, addr uintptr, data []byte) (count int, err error) ***REMOVED***
	return ptracePoke(PTRACE_POKEUSR, PTRACE_PEEKUSR, pid, addr, data)
***REMOVED***

func PtraceGetRegs(pid int, regsout *PtraceRegs) (err error) ***REMOVED***
	return ptrace(PTRACE_GETREGS, pid, 0, uintptr(unsafe.Pointer(regsout)))
***REMOVED***

func PtraceSetRegs(pid int, regs *PtraceRegs) (err error) ***REMOVED***
	return ptrace(PTRACE_SETREGS, pid, 0, uintptr(unsafe.Pointer(regs)))
***REMOVED***

func PtraceSetOptions(pid int, options int) (err error) ***REMOVED***
	return ptrace(PTRACE_SETOPTIONS, pid, 0, uintptr(options))
***REMOVED***

func PtraceGetEventMsg(pid int) (msg uint, err error) ***REMOVED***
	var data _C_long
	err = ptrace(PTRACE_GETEVENTMSG, pid, 0, uintptr(unsafe.Pointer(&data)))
	msg = uint(data)
	return
***REMOVED***

func PtraceCont(pid int, signal int) (err error) ***REMOVED***
	return ptrace(PTRACE_CONT, pid, 0, uintptr(signal))
***REMOVED***

func PtraceSyscall(pid int, signal int) (err error) ***REMOVED***
	return ptrace(PTRACE_SYSCALL, pid, 0, uintptr(signal))
***REMOVED***

func PtraceSingleStep(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_SINGLESTEP, pid, 0, 0) ***REMOVED***

func PtraceAttach(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_ATTACH, pid, 0, 0) ***REMOVED***

func PtraceDetach(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_DETACH, pid, 0, 0) ***REMOVED***

//sys	reboot(magic1 uint, magic2 uint, cmd int, arg string) (err error)

func Reboot(cmd int) (err error) ***REMOVED***
	return reboot(LINUX_REBOOT_MAGIC1, LINUX_REBOOT_MAGIC2, cmd, "")
***REMOVED***

func ReadDirent(fd int, buf []byte) (n int, err error) ***REMOVED***
	return Getdents(fd, buf)
***REMOVED***

func direntIno(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Ino), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Ino))
***REMOVED***

func direntReclen(buf []byte) (uint64, bool) ***REMOVED***
	return readInt(buf, unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Reclen), unsafe.Sizeof(Dirent***REMOVED******REMOVED***.Reclen))
***REMOVED***

func direntNamlen(buf []byte) (uint64, bool) ***REMOVED***
	reclen, ok := direntReclen(buf)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	return reclen - uint64(unsafe.Offsetof(Dirent***REMOVED******REMOVED***.Name)), true
***REMOVED***

//sys	mount(source string, target string, fstype string, flags uintptr, data *byte) (err error)

func Mount(source string, target string, fstype string, flags uintptr, data string) (err error) ***REMOVED***
	// Certain file systems get rather angry and EINVAL if you give
	// them an empty string of data, rather than NULL.
	if data == "" ***REMOVED***
		return mount(source, target, fstype, flags, nil)
	***REMOVED***
	datap, err := BytePtrFromString(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return mount(source, target, fstype, flags, datap)
***REMOVED***

// Sendto
// Recvfrom
// Socketpair

/*
 * Direct access
 */
//sys	Acct(path string) (err error)
//sys	AddKey(keyType string, description string, payload []byte, ringid int) (id int, err error)
//sys	Adjtimex(buf *Timex) (state int, err error)
//sys	Chdir(path string) (err error)
//sys	Chroot(path string) (err error)
//sys	ClockGettime(clockid int32, time *Timespec) (err error)
//sys	Close(fd int) (err error)
//sys	CopyFileRange(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error)
//sys	Dup(oldfd int) (fd int, err error)
//sys	Dup3(oldfd int, newfd int, flags int) (err error)
//sysnb	EpollCreate(size int) (fd int, err error)
//sysnb	EpollCreate1(flag int) (fd int, err error)
//sysnb	EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error)
//sys	Eventfd(initval uint, flags int) (fd int, err error) = SYS_EVENTFD2
//sys	Exit(code int) = SYS_EXIT_GROUP
//sys	Faccessat(dirfd int, path string, mode uint32, flags int) (err error)
//sys	Fallocate(fd int, mode uint32, off int64, len int64) (err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error)
//sys	fcntl(fd int, cmd int, arg int) (val int, err error)
//sys	Fdatasync(fd int) (err error)
//sys	Flock(fd int, how int) (err error)
//sys	Fsync(fd int) (err error)
//sys	Getdents(fd int, buf []byte) (n int, err error) = SYS_GETDENTS64
//sysnb	Getpgid(pid int) (pgid int, err error)

func Getpgrp() (pid int) ***REMOVED***
	pid, _ = Getpgid(0)
	return
***REMOVED***

//sysnb	Getpid() (pid int)
//sysnb	Getppid() (ppid int)
//sys	Getpriority(which int, who int) (prio int, err error)
//sys	Getrandom(buf []byte, flags int) (n int, err error)
//sysnb	Getrusage(who int, rusage *Rusage) (err error)
//sysnb	Getsid(pid int) (sid int, err error)
//sysnb	Gettid() (tid int)
//sys	Getxattr(path string, attr string, dest []byte) (sz int, err error)
//sys	InotifyAddWatch(fd int, pathname string, mask uint32) (watchdesc int, err error)
//sysnb	InotifyInit1(flags int) (fd int, err error)
//sysnb	InotifyRmWatch(fd int, watchdesc uint32) (success int, err error)
//sysnb	Kill(pid int, sig syscall.Signal) (err error)
//sys	Klogctl(typ int, buf []byte) (n int, err error) = SYS_SYSLOG
//sys	Lgetxattr(path string, attr string, dest []byte) (sz int, err error)
//sys	Listxattr(path string, dest []byte) (sz int, err error)
//sys	Llistxattr(path string, dest []byte) (sz int, err error)
//sys	Lremovexattr(path string, attr string) (err error)
//sys	Lsetxattr(path string, attr string, data []byte, flags int) (err error)
//sys	Mkdirat(dirfd int, path string, mode uint32) (err error)
//sys	Mknodat(dirfd int, path string, mode uint32, dev int) (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys	PivotRoot(newroot string, putold string) (err error) = SYS_PIVOT_ROOT
//sysnb prlimit(pid int, resource int, newlimit *Rlimit, old *Rlimit) (err error) = SYS_PRLIMIT64
//sys   Prctl(option int, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) (err error)
//sys	Pselect(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *Sigset_t) (n int, err error) = SYS_PSELECT6
//sys	read(fd int, p []byte) (n int, err error)
//sys	Removexattr(path string, attr string) (err error)
//sys	Renameat(olddirfd int, oldpath string, newdirfd int, newpath string) (err error)
//sys	RequestKey(keyType string, description string, callback string, destRingid int) (id int, err error)
//sys	Setdomainname(p []byte) (err error)
//sys	Sethostname(p []byte) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Settimeofday(tv *Timeval) (err error)
//sys	Setns(fd int, nstype int) (err error)

// issue 1435.
// On linux Setuid and Setgid only affects the current thread, not the process.
// This does not match what most callers expect so we must return an error
// here rather than letting the caller think that the call succeeded.

func Setuid(uid int) (err error) ***REMOVED***
	return EOPNOTSUPP
***REMOVED***

func Setgid(uid int) (err error) ***REMOVED***
	return EOPNOTSUPP
***REMOVED***

//sys	Setpriority(which int, who int, prio int) (err error)
//sys	Setxattr(path string, attr string, data []byte, flags int) (err error)
//sys	Sync()
//sys	Syncfs(fd int) (err error)
//sysnb	Sysinfo(info *Sysinfo_t) (err error)
//sys	Tee(rfd int, wfd int, len int, flags int) (n int64, err error)
//sysnb	Tgkill(tgid int, tid int, sig syscall.Signal) (err error)
//sysnb	Times(tms *Tms) (ticks uintptr, err error)
//sysnb	Umask(mask int) (oldmask int)
//sysnb	Uname(buf *Utsname) (err error)
//sys	Unmount(target string, flags int) (err error) = SYS_UMOUNT2
//sys	Unshare(flags int) (err error)
//sys	Ustat(dev int, ubuf *Ustat_t) (err error)
//sys	write(fd int, p []byte) (n int, err error)
//sys	exitThread(code int) (err error) = SYS_EXIT
//sys	readlen(fd int, p *byte, np int) (n int, err error) = SYS_READ
//sys	writelen(fd int, p *byte, np int) (n int, err error) = SYS_WRITE

// mmap varies by architecture; see syscall_linux_*.go.
//sys	munmap(addr uintptr, length uintptr) (err error)

var mapper = &mmapper***REMOVED***
	active: make(map[*byte][]byte),
	mmap:   mmap,
	munmap: munmap,
***REMOVED***

func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error) ***REMOVED***
	return mapper.Mmap(fd, offset, length, prot, flags)
***REMOVED***

func Munmap(b []byte) (err error) ***REMOVED***
	return mapper.Munmap(b)
***REMOVED***

//sys	Madvise(b []byte, advice int) (err error)
//sys	Mprotect(b []byte, prot int) (err error)
//sys	Mlock(b []byte) (err error)
//sys	Mlockall(flags int) (err error)
//sys	Msync(b []byte, flags int) (err error)
//sys	Munlock(b []byte) (err error)
//sys	Munlockall() (err error)

// Vmsplice splices user pages from a slice of Iovecs into a pipe specified by fd,
// using the specified flags.
func Vmsplice(fd int, iovs []Iovec, flags int) (int, error) ***REMOVED***
	n, _, errno := Syscall6(
		SYS_VMSPLICE,
		uintptr(fd),
		uintptr(unsafe.Pointer(&iovs[0])),
		uintptr(len(iovs)),
		uintptr(flags),
		0,
		0,
	)
	if errno != 0 ***REMOVED***
		return 0, syscall.Errno(errno)
	***REMOVED***

	return int(n), nil
***REMOVED***

/*
 * Unimplemented
 */
// AfsSyscall
// Alarm
// ArchPrctl
// Brk
// Capget
// Capset
// ClockGetres
// ClockNanosleep
// ClockSettime
// Clone
// CreateModule
// DeleteModule
// EpollCtlOld
// EpollPwait
// EpollWaitOld
// Execve
// Fgetxattr
// Flistxattr
// Fork
// Fremovexattr
// Fsetxattr
// Futex
// GetKernelSyms
// GetMempolicy
// GetRobustList
// GetThreadArea
// Getitimer
// Getpmsg
// IoCancel
// IoDestroy
// IoGetevents
// IoSetup
// IoSubmit
// IoprioGet
// IoprioSet
// KexecLoad
// LookupDcookie
// Mbind
// MigratePages
// Mincore
// ModifyLdt
// Mount
// MovePages
// MqGetsetattr
// MqNotify
// MqOpen
// MqTimedreceive
// MqTimedsend
// MqUnlink
// Mremap
// Msgctl
// Msgget
// Msgrcv
// Msgsnd
// Newfstatat
// Nfsservctl
// Personality
// Pselect6
// Ptrace
// Putpmsg
// QueryModule
// Quotactl
// Readahead
// Readv
// RemapFilePages
// RestartSyscall
// RtSigaction
// RtSigpending
// RtSigprocmask
// RtSigqueueinfo
// RtSigreturn
// RtSigsuspend
// RtSigtimedwait
// SchedGetPriorityMax
// SchedGetPriorityMin
// SchedGetaffinity
// SchedGetparam
// SchedGetscheduler
// SchedRrGetInterval
// SchedSetaffinity
// SchedSetparam
// SchedYield
// Security
// Semctl
// Semget
// Semop
// Semtimedop
// SetMempolicy
// SetRobustList
// SetThreadArea
// SetTidAddress
// Shmat
// Shmctl
// Shmdt
// Shmget
// Sigaltstack
// Signalfd
// Swapoff
// Swapon
// Sysfs
// TimerCreate
// TimerDelete
// TimerGetoverrun
// TimerGettime
// TimerSettime
// Timerfd
// Tkill (obsolete)
// Tuxcall
// Umount2
// Uselib
// Utimensat
// Vfork
// Vhangup
// Vserver
// Waitid
// _Sysctl
