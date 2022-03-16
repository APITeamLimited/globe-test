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
	"encoding/binary"
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

func EpollCreate(size int) (fd int, err error) ***REMOVED***
	if size <= 0 ***REMOVED***
		return -1, EINVAL
	***REMOVED***
	return EpollCreate1(0)
***REMOVED***

//sys	FanotifyInit(flags uint, event_f_flags uint) (fd int, err error)
//sys	fanotifyMark(fd int, flags uint, mask uint64, dirFd int, pathname *byte) (err error)

func FanotifyMark(fd int, flags uint, mask uint64, dirFd int, pathname string) (err error) ***REMOVED***
	if pathname == "" ***REMOVED***
		return fanotifyMark(fd, flags, mask, dirFd, nil)
	***REMOVED***
	p, err := BytePtrFromString(pathname)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return fanotifyMark(fd, flags, mask, dirFd, p)
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

func InotifyInit() (fd int, err error) ***REMOVED***
	return InotifyInit1(0)
***REMOVED***

//sys	ioctl(fd int, req uint, arg uintptr) (err error) = SYS_IOCTL
//sys	ioctlPtr(fd int, req uint, arg unsafe.Pointer) (err error) = SYS_IOCTL

// ioctl itself should not be exposed directly, but additional get/set functions
// for specific types are permissible. These are defined in ioctl.go and
// ioctl_linux.go.
//
// The third argument to ioctl is often a pointer but sometimes an integer.
// Callers should use ioctlPtr when the third argument is a pointer and ioctl
// when the third argument is an integer.
//
// TODO: some existing code incorrectly uses ioctl when it should use ioctlPtr.

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

//sys	openat2(dirfd int, path string, open_how *OpenHow, size int) (fd int, err error)

func Openat2(dirfd int, path string, how *OpenHow) (fd int, err error) ***REMOVED***
	return openat2(dirfd, path, how, SizeofOpenHow)
***REMOVED***

func Pipe(p []int) error ***REMOVED***
	return Pipe2(p, 0)
***REMOVED***

//sysnb	pipe2(p *[2]_C_int, flags int) (err error)

func Pipe2(p []int, flags int) error ***REMOVED***
	if len(p) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	var pp [2]_C_int
	err := pipe2(&pp, flags)
	if err == nil ***REMOVED***
		p[0] = int(pp[0])
		p[1] = int(pp[1])
	***REMOVED***
	return err
***REMOVED***

//sys	ppoll(fds *PollFd, nfds int, timeout *Timespec, sigmask *Sigset_t) (n int, err error)

func Ppoll(fds []PollFd, timeout *Timespec, sigmask *Sigset_t) (n int, err error) ***REMOVED***
	if len(fds) == 0 ***REMOVED***
		return ppoll(nil, 0, timeout, sigmask)
	***REMOVED***
	return ppoll(&fds[0], len(fds), timeout, sigmask)
***REMOVED***

func Poll(fds []PollFd, timeout int) (n int, err error) ***REMOVED***
	var ts *Timespec
	if timeout >= 0 ***REMOVED***
		ts = new(Timespec)
		*ts = NsecToTimespec(int64(timeout) * 1e6)
	***REMOVED***
	return Ppoll(fds, ts, nil)
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
	return UtimesNanoAt(AT_FDCWD, path, ts, 0)
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

func Futimesat(dirfd int, path string, tv []Timeval) error ***REMOVED***
	if tv == nil ***REMOVED***
		return futimesat(dirfd, path, nil)
	***REMOVED***
	if len(tv) != 2 ***REMOVED***
		return EINVAL
	***REMOVED***
	return futimesat(dirfd, path, (*[2]Timeval)(unsafe.Pointer(&tv[0])))
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
	sa.raw.Addr = sa.Addr
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
	sa.raw.Addr = sa.Addr
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

// SockaddrLinklayer implements the Sockaddr interface for AF_PACKET type sockets.
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
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrLinklayer, nil
***REMOVED***

// SockaddrNetlink implements the Sockaddr interface for AF_NETLINK type sockets.
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

// SockaddrHCI implements the Sockaddr interface for AF_BLUETOOTH type sockets
// using the HCI protocol.
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

// SockaddrL2 implements the Sockaddr interface for AF_BLUETOOTH type sockets
// using the L2CAP protocol.
type SockaddrL2 struct ***REMOVED***
	PSM      uint16
	CID      uint16
	Addr     [6]uint8
	AddrType uint8
	raw      RawSockaddrL2
***REMOVED***

func (sa *SockaddrL2) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_BLUETOOTH
	psm := (*[2]byte)(unsafe.Pointer(&sa.raw.Psm))
	psm[0] = byte(sa.PSM)
	psm[1] = byte(sa.PSM >> 8)
	for i := 0; i < len(sa.Addr); i++ ***REMOVED***
		sa.raw.Bdaddr[i] = sa.Addr[len(sa.Addr)-1-i]
	***REMOVED***
	cid := (*[2]byte)(unsafe.Pointer(&sa.raw.Cid))
	cid[0] = byte(sa.CID)
	cid[1] = byte(sa.CID >> 8)
	sa.raw.Bdaddr_type = sa.AddrType
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2, nil
***REMOVED***

// SockaddrRFCOMM implements the Sockaddr interface for AF_BLUETOOTH type sockets
// using the RFCOMM protocol.
//
// Server example:
//
//      fd, _ := Socket(AF_BLUETOOTH, SOCK_STREAM, BTPROTO_RFCOMM)
//      _ = unix.Bind(fd, &unix.SockaddrRFCOMM***REMOVED***
//      	Channel: 1,
//      	Addr:    [6]uint8***REMOVED***0, 0, 0, 0, 0, 0***REMOVED***, // BDADDR_ANY or 00:00:00:00:00:00
//      ***REMOVED***)
//      _ = Listen(fd, 1)
//      nfd, sa, _ := Accept(fd)
//      fmt.Printf("conn addr=%v fd=%d", sa.(*unix.SockaddrRFCOMM).Addr, nfd)
//      Read(nfd, buf)
//
// Client example:
//
//      fd, _ := Socket(AF_BLUETOOTH, SOCK_STREAM, BTPROTO_RFCOMM)
//      _ = Connect(fd, &SockaddrRFCOMM***REMOVED***
//      	Channel: 1,
//      	Addr:    [6]byte***REMOVED***0x11, 0x22, 0x33, 0xaa, 0xbb, 0xcc***REMOVED***, // CC:BB:AA:33:22:11
//      ***REMOVED***)
//      Write(fd, []byte(`hello`))
type SockaddrRFCOMM struct ***REMOVED***
	// Addr represents a bluetooth address, byte ordering is little-endian.
	Addr [6]uint8

	// Channel is a designated bluetooth channel, only 1-30 are available for use.
	// Since Linux 2.6.7 and further zero value is the first available channel.
	Channel uint8

	raw RawSockaddrRFCOMM
***REMOVED***

func (sa *SockaddrRFCOMM) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_BLUETOOTH
	sa.raw.Channel = sa.Channel
	sa.raw.Bdaddr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrRFCOMM, nil
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

// SockaddrCANJ1939 implements the Sockaddr interface for AF_CAN using J1939
// protocol (https://en.wikipedia.org/wiki/SAE_J1939). For more information
// on the purposes of the fields, check the official linux kernel documentation
// available here: https://www.kernel.org/doc/Documentation/networking/j1939.rst
type SockaddrCANJ1939 struct ***REMOVED***
	Ifindex int
	Name    uint64
	PGN     uint32
	Addr    uint8
	raw     RawSockaddrCAN
***REMOVED***

func (sa *SockaddrCANJ1939) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Ifindex < 0 || sa.Ifindex > 0x7fffffff ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_CAN
	sa.raw.Ifindex = int32(sa.Ifindex)
	n := (*[8]byte)(unsafe.Pointer(&sa.Name))
	for i := 0; i < 8; i++ ***REMOVED***
		sa.raw.Addr[i] = n[i]
	***REMOVED***
	p := (*[4]byte)(unsafe.Pointer(&sa.PGN))
	for i := 0; i < 4; i++ ***REMOVED***
		sa.raw.Addr[i+8] = p[i]
	***REMOVED***
	sa.raw.Addr[12] = sa.Addr
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
	//  - VMADDR_CID_LOCAL: refers to local communication (loopback).
	//  - VMADDR_CID_HOST: refers to other processes on the host.
	CID   uint32
	Port  uint32
	Flags uint8
	raw   RawSockaddrVM
***REMOVED***

func (sa *SockaddrVM) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_VSOCK
	sa.raw.Port = sa.Port
	sa.raw.Cid = sa.CID
	sa.raw.Flags = sa.Flags

	return unsafe.Pointer(&sa.raw), SizeofSockaddrVM, nil
***REMOVED***

type SockaddrXDP struct ***REMOVED***
	Flags        uint16
	Ifindex      uint32
	QueueID      uint32
	SharedUmemFD uint32
	raw          RawSockaddrXDP
***REMOVED***

func (sa *SockaddrXDP) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_XDP
	sa.raw.Flags = sa.Flags
	sa.raw.Ifindex = sa.Ifindex
	sa.raw.Queue_id = sa.QueueID
	sa.raw.Shared_umem_fd = sa.SharedUmemFD

	return unsafe.Pointer(&sa.raw), SizeofSockaddrXDP, nil
***REMOVED***

// This constant mirrors the #define of PX_PROTO_OE in
// linux/if_pppox.h. We're defining this by hand here instead of
// autogenerating through mkerrors.sh because including
// linux/if_pppox.h causes some declaration conflicts with other
// includes (linux/if_pppox.h includes linux/in.h, which conflicts
// with netinet/in.h). Given that we only need a single zero constant
// out of that file, it's cleaner to just define it by hand here.
const px_proto_oe = 0

type SockaddrPPPoE struct ***REMOVED***
	SID    uint16
	Remote []byte
	Dev    string
	raw    RawSockaddrPPPoX
***REMOVED***

func (sa *SockaddrPPPoE) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if len(sa.Remote) != 6 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	if len(sa.Dev) > IFNAMSIZ-1 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***

	*(*uint16)(unsafe.Pointer(&sa.raw[0])) = AF_PPPOX
	// This next field is in host-endian byte order. We can't use the
	// same unsafe pointer cast as above, because this value is not
	// 32-bit aligned and some architectures don't allow unaligned
	// access.
	//
	// However, the value of px_proto_oe is 0, so we can use
	// encoding/binary helpers to write the bytes without worrying
	// about the ordering.
	binary.BigEndian.PutUint32(sa.raw[2:6], px_proto_oe)
	// This field is deliberately big-endian, unlike the previous
	// one. The kernel expects SID to be in network byte order.
	binary.BigEndian.PutUint16(sa.raw[6:8], sa.SID)
	copy(sa.raw[8:14], sa.Remote)
	for i := 14; i < 14+IFNAMSIZ; i++ ***REMOVED***
		sa.raw[i] = 0
	***REMOVED***
	copy(sa.raw[14:], sa.Dev)
	return unsafe.Pointer(&sa.raw), SizeofSockaddrPPPoX, nil
***REMOVED***

// SockaddrTIPC implements the Sockaddr interface for AF_TIPC type sockets.
// For more information on TIPC, see: http://tipc.sourceforge.net/.
type SockaddrTIPC struct ***REMOVED***
	// Scope is the publication scopes when binding service/service range.
	// Should be set to TIPC_CLUSTER_SCOPE or TIPC_NODE_SCOPE.
	Scope int

	// Addr is the type of address used to manipulate a socket. Addr must be
	// one of:
	//  - *TIPCSocketAddr: "id" variant in the C addr union
	//  - *TIPCServiceRange: "nameseq" variant in the C addr union
	//  - *TIPCServiceName: "name" variant in the C addr union
	//
	// If nil, EINVAL will be returned when the structure is used.
	Addr TIPCAddr

	raw RawSockaddrTIPC
***REMOVED***

// TIPCAddr is implemented by types that can be used as an address for
// SockaddrTIPC. It is only implemented by *TIPCSocketAddr, *TIPCServiceRange,
// and *TIPCServiceName.
type TIPCAddr interface ***REMOVED***
	tipcAddrtype() uint8
	tipcAddr() [12]byte
***REMOVED***

func (sa *TIPCSocketAddr) tipcAddr() [12]byte ***REMOVED***
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCSocketAddr***REMOVED******REMOVED***)]byte)(unsafe.Pointer(sa)))[:])
	return out
***REMOVED***

func (sa *TIPCSocketAddr) tipcAddrtype() uint8 ***REMOVED*** return TIPC_SOCKET_ADDR ***REMOVED***

func (sa *TIPCServiceRange) tipcAddr() [12]byte ***REMOVED***
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCServiceRange***REMOVED******REMOVED***)]byte)(unsafe.Pointer(sa)))[:])
	return out
***REMOVED***

func (sa *TIPCServiceRange) tipcAddrtype() uint8 ***REMOVED*** return TIPC_SERVICE_RANGE ***REMOVED***

func (sa *TIPCServiceName) tipcAddr() [12]byte ***REMOVED***
	var out [12]byte
	copy(out[:], (*(*[unsafe.Sizeof(TIPCServiceName***REMOVED******REMOVED***)]byte)(unsafe.Pointer(sa)))[:])
	return out
***REMOVED***

func (sa *TIPCServiceName) tipcAddrtype() uint8 ***REMOVED*** return TIPC_SERVICE_ADDR ***REMOVED***

func (sa *SockaddrTIPC) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	if sa.Addr == nil ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	sa.raw.Family = AF_TIPC
	sa.raw.Scope = int8(sa.Scope)
	sa.raw.Addrtype = sa.Addr.tipcAddrtype()
	sa.raw.Addr = sa.Addr.tipcAddr()
	return unsafe.Pointer(&sa.raw), SizeofSockaddrTIPC, nil
***REMOVED***

// SockaddrL2TPIP implements the Sockaddr interface for IPPROTO_L2TP/AF_INET sockets.
type SockaddrL2TPIP struct ***REMOVED***
	Addr   [4]byte
	ConnId uint32
	raw    RawSockaddrL2TPIP
***REMOVED***

func (sa *SockaddrL2TPIP) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_INET
	sa.raw.Conn_id = sa.ConnId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2TPIP, nil
***REMOVED***

// SockaddrL2TPIP6 implements the Sockaddr interface for IPPROTO_L2TP/AF_INET6 sockets.
type SockaddrL2TPIP6 struct ***REMOVED***
	Addr   [16]byte
	ZoneId uint32
	ConnId uint32
	raw    RawSockaddrL2TPIP6
***REMOVED***

func (sa *SockaddrL2TPIP6) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_INET6
	sa.raw.Conn_id = sa.ConnId
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), SizeofSockaddrL2TPIP6, nil
***REMOVED***

// SockaddrIUCV implements the Sockaddr interface for AF_IUCV sockets.
type SockaddrIUCV struct ***REMOVED***
	UserID string
	Name   string
	raw    RawSockaddrIUCV
***REMOVED***

func (sa *SockaddrIUCV) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Family = AF_IUCV
	// These are EBCDIC encoded by the kernel, but we still need to pad them
	// with blanks. Initializing with blanks allows the caller to feed in either
	// a padded or an unpadded string.
	for i := 0; i < 8; i++ ***REMOVED***
		sa.raw.Nodeid[i] = ' '
		sa.raw.User_id[i] = ' '
		sa.raw.Name[i] = ' '
	***REMOVED***
	if len(sa.UserID) > 8 || len(sa.Name) > 8 ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	for i, b := range []byte(sa.UserID[:]) ***REMOVED***
		sa.raw.User_id[i] = int8(b)
	***REMOVED***
	for i, b := range []byte(sa.Name[:]) ***REMOVED***
		sa.raw.Name[i] = int8(b)
	***REMOVED***
	return unsafe.Pointer(&sa.raw), SizeofSockaddrIUCV, nil
***REMOVED***

type SockaddrNFC struct ***REMOVED***
	DeviceIdx   uint32
	TargetIdx   uint32
	NFCProtocol uint32
	raw         RawSockaddrNFC
***REMOVED***

func (sa *SockaddrNFC) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Sa_family = AF_NFC
	sa.raw.Dev_idx = sa.DeviceIdx
	sa.raw.Target_idx = sa.TargetIdx
	sa.raw.Nfc_protocol = sa.NFCProtocol
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNFC, nil
***REMOVED***

type SockaddrNFCLLCP struct ***REMOVED***
	DeviceIdx      uint32
	TargetIdx      uint32
	NFCProtocol    uint32
	DestinationSAP uint8
	SourceSAP      uint8
	ServiceName    string
	raw            RawSockaddrNFCLLCP
***REMOVED***

func (sa *SockaddrNFCLLCP) sockaddr() (unsafe.Pointer, _Socklen, error) ***REMOVED***
	sa.raw.Sa_family = AF_NFC
	sa.raw.Dev_idx = sa.DeviceIdx
	sa.raw.Target_idx = sa.TargetIdx
	sa.raw.Nfc_protocol = sa.NFCProtocol
	sa.raw.Dsap = sa.DestinationSAP
	sa.raw.Ssap = sa.SourceSAP
	if len(sa.ServiceName) > len(sa.raw.Service_name) ***REMOVED***
		return nil, 0, EINVAL
	***REMOVED***
	copy(sa.raw.Service_name[:], sa.ServiceName)
	sa.raw.SetServiceNameLen(len(sa.ServiceName))
	return unsafe.Pointer(&sa.raw), SizeofSockaddrNFCLLCP, nil
***REMOVED***

var socketProtocol = func(fd int) (int, error) ***REMOVED***
	return GetsockoptInt(fd, SOL_SOCKET, SO_PROTOCOL)
***REMOVED***

func anyToSockaddr(fd int, rsa *RawSockaddrAny) (Sockaddr, error) ***REMOVED***
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
		sa.Addr = pp.Addr
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
		bytes := (*[len(pp.Path)]byte)(unsafe.Pointer(&pp.Path[0]))[0:n]
		sa.Name = string(bytes)
		return sa, nil

	case AF_INET:
		proto, err := socketProtocol(fd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch proto ***REMOVED***
		case IPPROTO_L2TP:
			pp := (*RawSockaddrL2TPIP)(unsafe.Pointer(rsa))
			sa := new(SockaddrL2TPIP)
			sa.ConnId = pp.Conn_id
			sa.Addr = pp.Addr
			return sa, nil
		default:
			pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
			sa := new(SockaddrInet4)
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			sa.Port = int(p[0])<<8 + int(p[1])
			sa.Addr = pp.Addr
			return sa, nil
		***REMOVED***

	case AF_INET6:
		proto, err := socketProtocol(fd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch proto ***REMOVED***
		case IPPROTO_L2TP:
			pp := (*RawSockaddrL2TPIP6)(unsafe.Pointer(rsa))
			sa := new(SockaddrL2TPIP6)
			sa.ConnId = pp.Conn_id
			sa.ZoneId = pp.Scope_id
			sa.Addr = pp.Addr
			return sa, nil
		default:
			pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
			sa := new(SockaddrInet6)
			p := (*[2]byte)(unsafe.Pointer(&pp.Port))
			sa.Port = int(p[0])<<8 + int(p[1])
			sa.ZoneId = pp.Scope_id
			sa.Addr = pp.Addr
			return sa, nil
		***REMOVED***

	case AF_VSOCK:
		pp := (*RawSockaddrVM)(unsafe.Pointer(rsa))
		sa := &SockaddrVM***REMOVED***
			CID:   pp.Cid,
			Port:  pp.Port,
			Flags: pp.Flags,
		***REMOVED***
		return sa, nil
	case AF_BLUETOOTH:
		proto, err := socketProtocol(fd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// only BTPROTO_L2CAP and BTPROTO_RFCOMM can accept connections
		switch proto ***REMOVED***
		case BTPROTO_L2CAP:
			pp := (*RawSockaddrL2)(unsafe.Pointer(rsa))
			sa := &SockaddrL2***REMOVED***
				PSM:      pp.Psm,
				CID:      pp.Cid,
				Addr:     pp.Bdaddr,
				AddrType: pp.Bdaddr_type,
			***REMOVED***
			return sa, nil
		case BTPROTO_RFCOMM:
			pp := (*RawSockaddrRFCOMM)(unsafe.Pointer(rsa))
			sa := &SockaddrRFCOMM***REMOVED***
				Channel: pp.Channel,
				Addr:    pp.Bdaddr,
			***REMOVED***
			return sa, nil
		***REMOVED***
	case AF_XDP:
		pp := (*RawSockaddrXDP)(unsafe.Pointer(rsa))
		sa := &SockaddrXDP***REMOVED***
			Flags:        pp.Flags,
			Ifindex:      pp.Ifindex,
			QueueID:      pp.Queue_id,
			SharedUmemFD: pp.Shared_umem_fd,
		***REMOVED***
		return sa, nil
	case AF_PPPOX:
		pp := (*RawSockaddrPPPoX)(unsafe.Pointer(rsa))
		if binary.BigEndian.Uint32(pp[2:6]) != px_proto_oe ***REMOVED***
			return nil, EINVAL
		***REMOVED***
		sa := &SockaddrPPPoE***REMOVED***
			SID:    binary.BigEndian.Uint16(pp[6:8]),
			Remote: pp[8:14],
		***REMOVED***
		for i := 14; i < 14+IFNAMSIZ; i++ ***REMOVED***
			if pp[i] == 0 ***REMOVED***
				sa.Dev = string(pp[14:i])
				break
			***REMOVED***
		***REMOVED***
		return sa, nil
	case AF_TIPC:
		pp := (*RawSockaddrTIPC)(unsafe.Pointer(rsa))

		sa := &SockaddrTIPC***REMOVED***
			Scope: int(pp.Scope),
		***REMOVED***

		// Determine which union variant is present in pp.Addr by checking
		// pp.Addrtype.
		switch pp.Addrtype ***REMOVED***
		case TIPC_SERVICE_RANGE:
			sa.Addr = (*TIPCServiceRange)(unsafe.Pointer(&pp.Addr))
		case TIPC_SERVICE_ADDR:
			sa.Addr = (*TIPCServiceName)(unsafe.Pointer(&pp.Addr))
		case TIPC_SOCKET_ADDR:
			sa.Addr = (*TIPCSocketAddr)(unsafe.Pointer(&pp.Addr))
		default:
			return nil, EINVAL
		***REMOVED***

		return sa, nil
	case AF_IUCV:
		pp := (*RawSockaddrIUCV)(unsafe.Pointer(rsa))

		var user [8]byte
		var name [8]byte

		for i := 0; i < 8; i++ ***REMOVED***
			user[i] = byte(pp.User_id[i])
			name[i] = byte(pp.Name[i])
		***REMOVED***

		sa := &SockaddrIUCV***REMOVED***
			UserID: string(user[:]),
			Name:   string(name[:]),
		***REMOVED***
		return sa, nil

	case AF_CAN:
		proto, err := socketProtocol(fd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		pp := (*RawSockaddrCAN)(unsafe.Pointer(rsa))

		switch proto ***REMOVED***
		case CAN_J1939:
			sa := &SockaddrCANJ1939***REMOVED***
				Ifindex: int(pp.Ifindex),
			***REMOVED***
			name := (*[8]byte)(unsafe.Pointer(&sa.Name))
			for i := 0; i < 8; i++ ***REMOVED***
				name[i] = pp.Addr[i]
			***REMOVED***
			pgn := (*[4]byte)(unsafe.Pointer(&sa.PGN))
			for i := 0; i < 4; i++ ***REMOVED***
				pgn[i] = pp.Addr[i+8]
			***REMOVED***
			addr := (*[1]byte)(unsafe.Pointer(&sa.Addr))
			addr[0] = pp.Addr[12]
			return sa, nil
		default:
			sa := &SockaddrCAN***REMOVED***
				Ifindex: int(pp.Ifindex),
			***REMOVED***
			rx := (*[4]byte)(unsafe.Pointer(&sa.RxID))
			for i := 0; i < 4; i++ ***REMOVED***
				rx[i] = pp.Addr[i]
			***REMOVED***
			tx := (*[4]byte)(unsafe.Pointer(&sa.TxID))
			for i := 0; i < 4; i++ ***REMOVED***
				tx[i] = pp.Addr[i+4]
			***REMOVED***
			return sa, nil
		***REMOVED***
	case AF_NFC:
		proto, err := socketProtocol(fd)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		switch proto ***REMOVED***
		case NFC_SOCKPROTO_RAW:
			pp := (*RawSockaddrNFC)(unsafe.Pointer(rsa))
			sa := &SockaddrNFC***REMOVED***
				DeviceIdx:   pp.Dev_idx,
				TargetIdx:   pp.Target_idx,
				NFCProtocol: pp.Nfc_protocol,
			***REMOVED***
			return sa, nil
		case NFC_SOCKPROTO_LLCP:
			pp := (*RawSockaddrNFCLLCP)(unsafe.Pointer(rsa))
			if uint64(pp.Service_name_len) > uint64(len(pp.Service_name)) ***REMOVED***
				return nil, EINVAL
			***REMOVED***
			sa := &SockaddrNFCLLCP***REMOVED***
				DeviceIdx:      pp.Dev_idx,
				TargetIdx:      pp.Target_idx,
				NFCProtocol:    pp.Nfc_protocol,
				DestinationSAP: pp.Dsap,
				SourceSAP:      pp.Ssap,
				ServiceName:    string(pp.Service_name[:pp.Service_name_len]),
			***REMOVED***
			return sa, nil
		default:
			return nil, EINVAL
		***REMOVED***
	***REMOVED***
	return nil, EAFNOSUPPORT
***REMOVED***

func Accept(fd int) (nfd int, sa Sockaddr, err error) ***REMOVED***
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, 0)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	sa, err = anyToSockaddr(fd, &rsa)
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
	sa, err = anyToSockaddr(fd, &rsa)
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
	return anyToSockaddr(fd, &rsa)
***REMOVED***

func GetsockoptIPMreqn(fd, level, opt int) (*IPMreqn, error) ***REMOVED***
	var value IPMreqn
	vallen := _Socklen(SizeofIPMreqn)
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

// GetsockoptString returns the string value of the socket option opt for the
// socket associated with fd at the given socket level.
func GetsockoptString(fd, level, opt int) (string, error) ***REMOVED***
	buf := make([]byte, 256)
	vallen := _Socklen(len(buf))
	err := getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil ***REMOVED***
		if err == ERANGE ***REMOVED***
			buf = make([]byte, vallen)
			err = getsockopt(fd, level, opt, unsafe.Pointer(&buf[0]), &vallen)
		***REMOVED***
		if err != nil ***REMOVED***
			return "", err
		***REMOVED***
	***REMOVED***
	return string(buf[:vallen-1]), nil
***REMOVED***

func GetsockoptTpacketStats(fd, level, opt int) (*TpacketStats, error) ***REMOVED***
	var value TpacketStats
	vallen := _Socklen(SizeofTpacketStats)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func GetsockoptTpacketStatsV3(fd, level, opt int) (*TpacketStatsV3, error) ***REMOVED***
	var value TpacketStatsV3
	vallen := _Socklen(SizeofTpacketStatsV3)
	err := getsockopt(fd, level, opt, unsafe.Pointer(&value), &vallen)
	return &value, err
***REMOVED***

func SetsockoptIPMreqn(fd, level, opt int, mreq *IPMreqn) (err error) ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
***REMOVED***

func SetsockoptPacketMreq(fd, level, opt int, mreq *PacketMreq) error ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(mreq), unsafe.Sizeof(*mreq))
***REMOVED***

// SetsockoptSockFprog attaches a classic BPF or an extended BPF program to a
// socket to filter incoming packets.  See 'man 7 socket' for usage information.
func SetsockoptSockFprog(fd, level, opt int, fprog *SockFprog) error ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(fprog), unsafe.Sizeof(*fprog))
***REMOVED***

func SetsockoptCanRawFilter(fd, level, opt int, filter []CanFilter) error ***REMOVED***
	var p unsafe.Pointer
	if len(filter) > 0 ***REMOVED***
		p = unsafe.Pointer(&filter[0])
	***REMOVED***
	return setsockopt(fd, level, opt, p, uintptr(len(filter)*SizeofCanFilter))
***REMOVED***

func SetsockoptTpacketReq(fd, level, opt int, tp *TpacketReq) error ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(tp), unsafe.Sizeof(*tp))
***REMOVED***

func SetsockoptTpacketReq3(fd, level, opt int, tp *TpacketReq3) error ***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(tp), unsafe.Sizeof(*tp))
***REMOVED***

func SetsockoptTCPRepairOpt(fd, level, opt int, o []TCPRepairOpt) (err error) ***REMOVED***
	if len(o) == 0 ***REMOVED***
		return EINVAL
	***REMOVED***
	return setsockopt(fd, level, opt, unsafe.Pointer(&o[0]), uintptr(SizeofTCPRepairOpt*len(o)))
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

// KeyctlRestrictKeyring implements the KEYCTL_RESTRICT_KEYRING command. This
// command limits the set of keys that can be linked to the keyring, regardless
// of keyring permissions. The command requires the "setattr" permission.
//
// When called with an empty keyType the command locks the keyring, preventing
// any further keys from being linked to the keyring.
//
// The "asymmetric" keyType defines restrictions requiring key payloads to be
// DER encoded X.509 certificates signed by keys in another keyring. Restrictions
// for "asymmetric" include "builtin_trusted", "builtin_and_secondary_trusted",
// "key_or_keyring:<key>", and "key_or_keyring:<key>:chain".
//
// As of Linux 4.12, only the "asymmetric" keyType defines type-specific
// restrictions.
//
// See the full documentation at:
// http://man7.org/linux/man-pages/man3/keyctl_restrict_keyring.3.html
// http://man7.org/linux/man-pages/man2/keyctl.2.html
func KeyctlRestrictKeyring(ringid int, keyType string, restriction string) error ***REMOVED***
	if keyType == "" ***REMOVED***
		return keyctlRestrictKeyring(KEYCTL_RESTRICT_KEYRING, ringid)
	***REMOVED***
	return keyctlRestrictKeyringByType(KEYCTL_RESTRICT_KEYRING, ringid, keyType, restriction)
***REMOVED***

//sys	keyctlRestrictKeyringByType(cmd int, arg2 int, keyType string, restriction string) (err error) = SYS_KEYCTL
//sys	keyctlRestrictKeyring(cmd int, arg2 int) (err error) = SYS_KEYCTL

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
		if len(p) == 0 ***REMOVED***
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil ***REMOVED***
				return
			***REMOVED***
			// receive at least one normal byte
			if sockType != SOCK_DGRAM ***REMOVED***
				iov.Base = &dummy
				iov.SetLen(1)
			***REMOVED***
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
		from, err = anyToSockaddr(fd, &rsa)
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
		if len(p) == 0 ***REMOVED***
			var sockType int
			sockType, err = GetsockoptInt(fd, SOL_SOCKET, SO_TYPE)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
			// send at least one normal byte
			if sockType != SOCK_DGRAM ***REMOVED***
				iov.Base = &dummy
				iov.SetLen(1)
			***REMOVED***
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

	var buf [SizeofPtr]byte

	// Leading edge. PEEKTEXT/PEEKDATA don't require aligned
	// access (PEEKUSER warns that it might), but if we don't
	// align our reads, we might straddle an unmapped page
	// boundary and not get the bytes leading up to the page
	// boundary.
	n := 0
	if addr%SizeofPtr != 0 ***REMOVED***
		err = ptrace(req, pid, addr-addr%SizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		n += copy(out, buf[addr%SizeofPtr:])
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
	if addr%SizeofPtr != 0 ***REMOVED***
		var buf [SizeofPtr]byte
		err = ptrace(peekReq, pid, addr-addr%SizeofPtr, uintptr(unsafe.Pointer(&buf[0])))
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		n += copy(buf[addr%SizeofPtr:], data)
		word := *((*uintptr)(unsafe.Pointer(&buf[0])))
		err = ptrace(pokeReq, pid, addr-addr%SizeofPtr, word)
		if err != nil ***REMOVED***
			return 0, err
		***REMOVED***
		data = data[n:]
	***REMOVED***

	// Interior.
	for len(data) > SizeofPtr ***REMOVED***
		word := *((*uintptr)(unsafe.Pointer(&data[0])))
		err = ptrace(pokeReq, pid, addr+uintptr(n), word)
		if err != nil ***REMOVED***
			return n, err
		***REMOVED***
		n += SizeofPtr
		data = data[SizeofPtr:]
	***REMOVED***

	// Trailing edge.
	if len(data) > 0 ***REMOVED***
		var buf [SizeofPtr]byte
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

func PtraceInterrupt(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_INTERRUPT, pid, 0, 0) ***REMOVED***

func PtraceAttach(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_ATTACH, pid, 0, 0) ***REMOVED***

func PtraceSeize(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_SEIZE, pid, 0, 0) ***REMOVED***

func PtraceDetach(pid int) (err error) ***REMOVED*** return ptrace(PTRACE_DETACH, pid, 0, 0) ***REMOVED***

//sys	reboot(magic1 uint, magic2 uint, cmd int, arg string) (err error)

func Reboot(cmd int) (err error) ***REMOVED***
	return reboot(LINUX_REBOOT_MAGIC1, LINUX_REBOOT_MAGIC2, cmd, "")
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

//sys	mountSetattr(dirfd int, pathname string, flags uint, attr *MountAttr, size uintptr) (err error) = SYS_MOUNT_SETATTR

// MountSetattr is a wrapper for mount_setattr(2).
// https://man7.org/linux/man-pages/man2/mount_setattr.2.html
//
// Requires kernel >= 5.12.
func MountSetattr(dirfd int, pathname string, flags uint, attr *MountAttr) error ***REMOVED***
	return mountSetattr(dirfd, pathname, flags, attr, unsafe.Sizeof(*attr))
***REMOVED***

func Sendfile(outfd int, infd int, offset *int64, count int) (written int, err error) ***REMOVED***
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	return sendfile(outfd, infd, offset, count)
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
//sysnb	Capget(hdr *CapUserHeader, data *CapUserData) (err error)
//sysnb	Capset(hdr *CapUserHeader, data *CapUserData) (err error)
//sys	Chdir(path string) (err error)
//sys	Chroot(path string) (err error)
//sys	ClockGetres(clockid int32, res *Timespec) (err error)
//sys	ClockGettime(clockid int32, time *Timespec) (err error)
//sys	ClockNanosleep(clockid int32, flags int, request *Timespec, remain *Timespec) (err error)
//sys	Close(fd int) (err error)
//sys	CloseRange(first uint, last uint, flags uint) (err error)
//sys	CopyFileRange(rfd int, roff *int64, wfd int, woff *int64, len int, flags int) (n int, err error)
//sys	DeleteModule(name string, flags int) (err error)
//sys	Dup(oldfd int) (fd int, err error)

func Dup2(oldfd, newfd int) error ***REMOVED***
	return Dup3(oldfd, newfd, 0)
***REMOVED***

//sys	Dup3(oldfd int, newfd int, flags int) (err error)
//sysnb	EpollCreate1(flag int) (fd int, err error)
//sysnb	EpollCtl(epfd int, op int, fd int, event *EpollEvent) (err error)
//sys	Eventfd(initval uint, flags int) (fd int, err error) = SYS_EVENTFD2
//sys	Exit(code int) = SYS_EXIT_GROUP
//sys	Fallocate(fd int, mode uint32, off int64, len int64) (err error)
//sys	Fchdir(fd int) (err error)
//sys	Fchmod(fd int, mode uint32) (err error)
//sys	Fchownat(dirfd int, path string, uid int, gid int, flags int) (err error)
//sys	Fdatasync(fd int) (err error)
//sys	Fgetxattr(fd int, attr string, dest []byte) (sz int, err error)
//sys	FinitModule(fd int, params string, flags int) (err error)
//sys	Flistxattr(fd int, dest []byte) (sz int, err error)
//sys	Flock(fd int, how int) (err error)
//sys	Fremovexattr(fd int, attr string) (err error)
//sys	Fsetxattr(fd int, attr string, dest []byte, flags int) (err error)
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
//sys	InitModule(moduleImage []byte, params string) (err error)
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
//sys	MemfdCreate(name string, flags int) (fd int, err error)
//sys	Mkdirat(dirfd int, path string, mode uint32) (err error)
//sys	Mknodat(dirfd int, path string, mode uint32, dev int) (err error)
//sys	Nanosleep(time *Timespec, leftover *Timespec) (err error)
//sys	PerfEventOpen(attr *PerfEventAttr, pid int, cpu int, groupFd int, flags int) (fd int, err error)
//sys	PivotRoot(newroot string, putold string) (err error) = SYS_PIVOT_ROOT
//sysnb	Prlimit(pid int, resource int, newlimit *Rlimit, old *Rlimit) (err error) = SYS_PRLIMIT64
//sys	Prctl(option int, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) (err error)
//sys	Pselect(nfd int, r *FdSet, w *FdSet, e *FdSet, timeout *Timespec, sigmask *Sigset_t) (n int, err error) = SYS_PSELECT6
//sys	read(fd int, p []byte) (n int, err error)
//sys	Removexattr(path string, attr string) (err error)
//sys	Renameat2(olddirfd int, oldpath string, newdirfd int, newpath string, flags uint) (err error)
//sys	RequestKey(keyType string, description string, callback string, destRingid int) (id int, err error)
//sys	Setdomainname(p []byte) (err error)
//sys	Sethostname(p []byte) (err error)
//sysnb	Setpgid(pid int, pgid int) (err error)
//sysnb	Setsid() (pid int, err error)
//sysnb	Settimeofday(tv *Timeval) (err error)
//sys	Setns(fd int, nstype int) (err error)

// PrctlRetInt performs a prctl operation specified by option and further
// optional arguments arg2 through arg5 depending on option. It returns a
// non-negative integer that is returned by the prctl syscall.
func PrctlRetInt(option int, arg2 uintptr, arg3 uintptr, arg4 uintptr, arg5 uintptr) (int, error) ***REMOVED***
	ret, _, err := Syscall6(SYS_PRCTL, uintptr(option), uintptr(arg2), uintptr(arg3), uintptr(arg4), uintptr(arg5), 0)
	if err != 0 ***REMOVED***
		return 0, err
	***REMOVED***
	return int(ret), nil
***REMOVED***

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

// SetfsgidRetGid sets fsgid for current thread and returns previous fsgid set.
// setfsgid(2) will return a non-nil error only if its caller lacks CAP_SETUID capability.
// If the call fails due to other reasons, current fsgid will be returned.
func SetfsgidRetGid(gid int) (int, error) ***REMOVED***
	return setfsgid(gid)
***REMOVED***

// SetfsuidRetUid sets fsuid for current thread and returns previous fsuid set.
// setfsgid(2) will return a non-nil error only if its caller lacks CAP_SETUID capability
// If the call fails due to other reasons, current fsuid will be returned.
func SetfsuidRetUid(uid int) (int, error) ***REMOVED***
	return setfsuid(uid)
***REMOVED***

func Setfsgid(gid int) error ***REMOVED***
	_, err := setfsgid(gid)
	return err
***REMOVED***

func Setfsuid(uid int) error ***REMOVED***
	_, err := setfsuid(uid)
	return err
***REMOVED***

func Signalfd(fd int, sigmask *Sigset_t, flags int) (newfd int, err error) ***REMOVED***
	return signalfd(fd, sigmask, _C__NSIG/8, flags)
***REMOVED***

//sys	Setpriority(which int, who int, prio int) (err error)
//sys	Setxattr(path string, attr string, data []byte, flags int) (err error)
//sys	signalfd(fd int, sigmask *Sigset_t, maskSize uintptr, flags int) (newfd int, err error) = SYS_SIGNALFD4
//sys	Statx(dirfd int, path string, flags int, mask int, stat *Statx_t) (err error)
//sys	Sync()
//sys	Syncfs(fd int) (err error)
//sysnb	Sysinfo(info *Sysinfo_t) (err error)
//sys	Tee(rfd int, wfd int, len int, flags int) (n int64, err error)
//sysnb	TimerfdCreate(clockid int, flags int) (fd int, err error)
//sysnb	TimerfdGettime(fd int, currValue *ItimerSpec) (err error)
//sysnb	TimerfdSettime(fd int, flags int, newValue *ItimerSpec, oldValue *ItimerSpec) (err error)
//sysnb	Tgkill(tgid int, tid int, sig syscall.Signal) (err error)
//sysnb	Times(tms *Tms) (ticks uintptr, err error)
//sysnb	Umask(mask int) (oldmask int)
//sysnb	Uname(buf *Utsname) (err error)
//sys	Unmount(target string, flags int) (err error) = SYS_UMOUNT2
//sys	Unshare(flags int) (err error)
//sys	write(fd int, p []byte) (n int, err error)
//sys	exitThread(code int) (err error) = SYS_EXIT
//sys	readlen(fd int, p *byte, np int) (n int, err error) = SYS_READ
//sys	writelen(fd int, p *byte, np int) (n int, err error) = SYS_WRITE
//sys	readv(fd int, iovs []Iovec) (n int, err error) = SYS_READV
//sys	writev(fd int, iovs []Iovec) (n int, err error) = SYS_WRITEV
//sys	preadv(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr) (n int, err error) = SYS_PREADV
//sys	pwritev(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr) (n int, err error) = SYS_PWRITEV
//sys	preadv2(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr, flags int) (n int, err error) = SYS_PREADV2
//sys	pwritev2(fd int, iovs []Iovec, offs_l uintptr, offs_h uintptr, flags int) (n int, err error) = SYS_PWRITEV2

func bytes2iovec(bs [][]byte) []Iovec ***REMOVED***
	iovecs := make([]Iovec, len(bs))
	for i, b := range bs ***REMOVED***
		iovecs[i].SetLen(len(b))
		if len(b) > 0 ***REMOVED***
			iovecs[i].Base = &b[0]
		***REMOVED*** else ***REMOVED***
			iovecs[i].Base = (*byte)(unsafe.Pointer(&_zero))
		***REMOVED***
	***REMOVED***
	return iovecs
***REMOVED***

// offs2lohi splits offs into its lower and upper unsigned long. On 64-bit
// systems, hi will always be 0. On 32-bit systems, offs will be split in half.
// preadv/pwritev chose this calling convention so they don't need to add a
// padding-register for alignment on ARM.
func offs2lohi(offs int64) (lo, hi uintptr) ***REMOVED***
	return uintptr(offs), uintptr(uint64(offs) >> SizeofLong)
***REMOVED***

func Readv(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	n, err = readv(fd, iovecs)
	readvRacedetect(iovecs, n, err)
	return n, err
***REMOVED***

func Preadv(fd int, iovs [][]byte, offset int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	lo, hi := offs2lohi(offset)
	n, err = preadv(fd, iovecs, lo, hi)
	readvRacedetect(iovecs, n, err)
	return n, err
***REMOVED***

func Preadv2(fd int, iovs [][]byte, offset int64, flags int) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	lo, hi := offs2lohi(offset)
	n, err = preadv2(fd, iovecs, lo, hi, flags)
	readvRacedetect(iovecs, n, err)
	return n, err
***REMOVED***

func readvRacedetect(iovecs []Iovec, n int, err error) ***REMOVED***
	if !raceenabled ***REMOVED***
		return
	***REMOVED***
	for i := 0; n > 0 && i < len(iovecs); i++ ***REMOVED***
		m := int(iovecs[i].Len)
		if m > n ***REMOVED***
			m = n
		***REMOVED***
		n -= m
		if m > 0 ***REMOVED***
			raceWriteRange(unsafe.Pointer(iovecs[i].Base), m)
		***REMOVED***
	***REMOVED***
	if err == nil ***REMOVED***
		raceAcquire(unsafe.Pointer(&ioSync))
	***REMOVED***
***REMOVED***

func Writev(fd int, iovs [][]byte) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	n, err = writev(fd, iovecs)
	writevRacedetect(iovecs, n)
	return n, err
***REMOVED***

func Pwritev(fd int, iovs [][]byte, offset int64) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	lo, hi := offs2lohi(offset)
	n, err = pwritev(fd, iovecs, lo, hi)
	writevRacedetect(iovecs, n)
	return n, err
***REMOVED***

func Pwritev2(fd int, iovs [][]byte, offset int64, flags int) (n int, err error) ***REMOVED***
	iovecs := bytes2iovec(iovs)
	if raceenabled ***REMOVED***
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	***REMOVED***
	lo, hi := offs2lohi(offset)
	n, err = pwritev2(fd, iovecs, lo, hi, flags)
	writevRacedetect(iovecs, n)
	return n, err
***REMOVED***

func writevRacedetect(iovecs []Iovec, n int) ***REMOVED***
	if !raceenabled ***REMOVED***
		return
	***REMOVED***
	for i := 0; n > 0 && i < len(iovecs); i++ ***REMOVED***
		m := int(iovecs[i].Len)
		if m > n ***REMOVED***
			m = n
		***REMOVED***
		n -= m
		if m > 0 ***REMOVED***
			raceReadRange(unsafe.Pointer(iovecs[i].Base), m)
		***REMOVED***
	***REMOVED***
***REMOVED***

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
	var p unsafe.Pointer
	if len(iovs) > 0 ***REMOVED***
		p = unsafe.Pointer(&iovs[0])
	***REMOVED***

	n, _, errno := Syscall6(SYS_VMSPLICE, uintptr(fd), uintptr(p), uintptr(len(iovs)), uintptr(flags), 0, 0)
	if errno != 0 ***REMOVED***
		return 0, syscall.Errno(errno)
	***REMOVED***

	return int(n), nil
***REMOVED***

func isGroupMember(gid int) bool ***REMOVED***
	groups, err := Getgroups()
	if err != nil ***REMOVED***
		return false
	***REMOVED***

	for _, g := range groups ***REMOVED***
		if g == gid ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

//sys	faccessat(dirfd int, path string, mode uint32) (err error)
//sys	Faccessat2(dirfd int, path string, mode uint32, flags int) (err error)

func Faccessat(dirfd int, path string, mode uint32, flags int) (err error) ***REMOVED***
	if flags == 0 ***REMOVED***
		return faccessat(dirfd, path, mode)
	***REMOVED***

	if err := Faccessat2(dirfd, path, mode, flags); err != ENOSYS && err != EPERM ***REMOVED***
		return err
	***REMOVED***

	// The Linux kernel faccessat system call does not take any flags.
	// The glibc faccessat implements the flags itself; see
	// https://sourceware.org/git/?p=glibc.git;a=blob;f=sysdeps/unix/sysv/linux/faccessat.c;hb=HEAD
	// Because people naturally expect syscall.Faccessat to act
	// like C faccessat, we do the same.

	if flags & ^(AT_SYMLINK_NOFOLLOW|AT_EACCESS) != 0 ***REMOVED***
		return EINVAL
	***REMOVED***

	var st Stat_t
	if err := Fstatat(dirfd, path, &st, flags&AT_SYMLINK_NOFOLLOW); err != nil ***REMOVED***
		return err
	***REMOVED***

	mode &= 7
	if mode == 0 ***REMOVED***
		return nil
	***REMOVED***

	var uid int
	if flags&AT_EACCESS != 0 ***REMOVED***
		uid = Geteuid()
	***REMOVED*** else ***REMOVED***
		uid = Getuid()
	***REMOVED***

	if uid == 0 ***REMOVED***
		if mode&1 == 0 ***REMOVED***
			// Root can read and write any file.
			return nil
		***REMOVED***
		if st.Mode&0111 != 0 ***REMOVED***
			// Root can execute any file that anybody can execute.
			return nil
		***REMOVED***
		return EACCES
	***REMOVED***

	var fmode uint32
	if uint32(uid) == st.Uid ***REMOVED***
		fmode = (st.Mode >> 6) & 7
	***REMOVED*** else ***REMOVED***
		var gid int
		if flags&AT_EACCESS != 0 ***REMOVED***
			gid = Getegid()
		***REMOVED*** else ***REMOVED***
			gid = Getgid()
		***REMOVED***

		if uint32(gid) == st.Gid || isGroupMember(gid) ***REMOVED***
			fmode = (st.Mode >> 3) & 7
		***REMOVED*** else ***REMOVED***
			fmode = st.Mode & 7
		***REMOVED***
	***REMOVED***

	if fmode&mode == mode ***REMOVED***
		return nil
	***REMOVED***

	return EACCES
***REMOVED***

//sys	nameToHandleAt(dirFD int, pathname string, fh *fileHandle, mountID *_C_int, flags int) (err error) = SYS_NAME_TO_HANDLE_AT
//sys	openByHandleAt(mountFD int, fh *fileHandle, flags int) (fd int, err error) = SYS_OPEN_BY_HANDLE_AT

// fileHandle is the argument to nameToHandleAt and openByHandleAt. We
// originally tried to generate it via unix/linux/types.go with "type
// fileHandle C.struct_file_handle" but that generated empty structs
// for mips64 and mips64le. Instead, hard code it for now (it's the
// same everywhere else) until the mips64 generator issue is fixed.
type fileHandle struct ***REMOVED***
	Bytes uint32
	Type  int32
***REMOVED***

// FileHandle represents the C struct file_handle used by
// name_to_handle_at (see NameToHandleAt) and open_by_handle_at (see
// OpenByHandleAt).
type FileHandle struct ***REMOVED***
	*fileHandle
***REMOVED***

// NewFileHandle constructs a FileHandle.
func NewFileHandle(handleType int32, handle []byte) FileHandle ***REMOVED***
	const hdrSize = unsafe.Sizeof(fileHandle***REMOVED******REMOVED***)
	buf := make([]byte, hdrSize+uintptr(len(handle)))
	copy(buf[hdrSize:], handle)
	fh := (*fileHandle)(unsafe.Pointer(&buf[0]))
	fh.Type = handleType
	fh.Bytes = uint32(len(handle))
	return FileHandle***REMOVED***fh***REMOVED***
***REMOVED***

func (fh *FileHandle) Size() int   ***REMOVED*** return int(fh.fileHandle.Bytes) ***REMOVED***
func (fh *FileHandle) Type() int32 ***REMOVED*** return fh.fileHandle.Type ***REMOVED***
func (fh *FileHandle) Bytes() []byte ***REMOVED***
	n := fh.Size()
	if n == 0 ***REMOVED***
		return nil
	***REMOVED***
	return (*[1 << 30]byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&fh.fileHandle.Type)) + 4))[:n:n]
***REMOVED***

// NameToHandleAt wraps the name_to_handle_at system call; it obtains
// a handle for a path name.
func NameToHandleAt(dirfd int, path string, flags int) (handle FileHandle, mountID int, err error) ***REMOVED***
	var mid _C_int
	// Try first with a small buffer, assuming the handle will
	// only be 32 bytes.
	size := uint32(32 + unsafe.Sizeof(fileHandle***REMOVED******REMOVED***))
	didResize := false
	for ***REMOVED***
		buf := make([]byte, size)
		fh := (*fileHandle)(unsafe.Pointer(&buf[0]))
		fh.Bytes = size - uint32(unsafe.Sizeof(fileHandle***REMOVED******REMOVED***))
		err = nameToHandleAt(dirfd, path, fh, &mid, flags)
		if err == EOVERFLOW ***REMOVED***
			if didResize ***REMOVED***
				// We shouldn't need to resize more than once
				return
			***REMOVED***
			didResize = true
			size = fh.Bytes + uint32(unsafe.Sizeof(fileHandle***REMOVED******REMOVED***))
			continue
		***REMOVED***
		if err != nil ***REMOVED***
			return
		***REMOVED***
		return FileHandle***REMOVED***fh***REMOVED***, int(mid), nil
	***REMOVED***
***REMOVED***

// OpenByHandleAt wraps the open_by_handle_at system call; it opens a
// file via a handle as previously returned by NameToHandleAt.
func OpenByHandleAt(mountFD int, handle FileHandle, flags int) (fd int, err error) ***REMOVED***
	return openByHandleAt(mountFD, handle.fileHandle, flags)
***REMOVED***

// Klogset wraps the sys_syslog system call; it sets console_loglevel to
// the value specified by arg and passes a dummy pointer to bufp.
func Klogset(typ int, arg int) (err error) ***REMOVED***
	var p unsafe.Pointer
	_, _, errno := Syscall(SYS_SYSLOG, uintptr(typ), uintptr(p), uintptr(arg))
	if errno != 0 ***REMOVED***
		return errnoErr(errno)
	***REMOVED***
	return nil
***REMOVED***

// RemoteIovec is Iovec with the pointer replaced with an integer.
// It is used for ProcessVMReadv and ProcessVMWritev, where the pointer
// refers to a location in a different process' address space, which
// would confuse the Go garbage collector.
type RemoteIovec struct ***REMOVED***
	Base uintptr
	Len  int
***REMOVED***

//sys	ProcessVMReadv(pid int, localIov []Iovec, remoteIov []RemoteIovec, flags uint) (n int, err error) = SYS_PROCESS_VM_READV
//sys	ProcessVMWritev(pid int, localIov []Iovec, remoteIov []RemoteIovec, flags uint) (n int, err error) = SYS_PROCESS_VM_WRITEV

//sys	PidfdOpen(pid int, flags int) (fd int, err error) = SYS_PIDFD_OPEN
//sys	PidfdGetfd(pidfd int, targetfd int, flags int) (fd int, err error) = SYS_PIDFD_GETFD

//sys	shmat(id int, addr uintptr, flag int) (ret uintptr, err error)
//sys	shmctl(id int, cmd int, buf *SysvShmDesc) (result int, err error)
//sys	shmdt(addr uintptr) (err error)
//sys	shmget(key int, size int, flag int) (id int, err error)

/*
 * Unimplemented
 */
// AfsSyscall
// Alarm
// ArchPrctl
// Brk
// ClockNanosleep
// ClockSettime
// Clone
// EpollCtlOld
// EpollPwait
// EpollWaitOld
// Execve
// Fork
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
// Nfsservctl
// Personality
// Pselect6
// Ptrace
// Putpmsg
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
// SchedGetparam
// SchedGetscheduler
// SchedRrGetInterval
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
// Sigaltstack
// Swapoff
// Swapon
// Sysfs
// TimerCreate
// TimerDelete
// TimerGetoverrun
// TimerGettime
// TimerSettime
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
