// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build zos && s390x
// +build zos,s390x

// Hand edited based on ztypes_linux_s390x.go
// TODO: auto-generate.

package unix

const (
	SizeofPtr      = 0x8
	SizeofShort    = 0x2
	SizeofInt      = 0x4
	SizeofLong     = 0x8
	SizeofLongLong = 0x8
	PathMax        = 0x1000
)

const (
	SizeofSockaddrAny   = 128
	SizeofCmsghdr       = 12
	SizeofIPMreq        = 8
	SizeofIPv6Mreq      = 20
	SizeofICMPv6Filter  = 32
	SizeofIPv6MTUInfo   = 32
	SizeofLinger        = 8
	SizeofSockaddrInet4 = 16
	SizeofSockaddrInet6 = 28
	SizeofTCPInfo       = 0x68
)

type (
	_C_short     int16
	_C_int       int32
	_C_long      int64
	_C_long_long int64
)

type Timespec struct ***REMOVED***
	Sec  int64
	Nsec int64
***REMOVED***

type Timeval struct ***REMOVED***
	Sec  int64
	Usec int64
***REMOVED***

type timeval_zos struct ***REMOVED*** //correct (with padding and all)
	Sec  int64
	_    [4]byte // pad
	Usec int32
***REMOVED***

type Tms struct ***REMOVED*** //clock_t is 4-byte unsigned int in zos
	Utime  uint32
	Stime  uint32
	Cutime uint32
	Cstime uint32
***REMOVED***

type Time_t int64

type Utimbuf struct ***REMOVED***
	Actime  int64
	Modtime int64
***REMOVED***

type Utsname struct ***REMOVED***
	Sysname    [65]byte
	Nodename   [65]byte
	Release    [65]byte
	Version    [65]byte
	Machine    [65]byte
	Domainname [65]byte
***REMOVED***

type RawSockaddrInet4 struct ***REMOVED***
	Len    uint8
	Family uint8
	Port   uint16
	Addr   [4]byte /* in_addr */
	Zero   [8]uint8
***REMOVED***

type RawSockaddrInet6 struct ***REMOVED***
	Len      uint8
	Family   uint8
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte /* in6_addr */
	Scope_id uint32
***REMOVED***

type RawSockaddrUnix struct ***REMOVED***
	Len    uint8
	Family uint8
	Path   [108]int8
***REMOVED***

type RawSockaddr struct ***REMOVED***
	Len    uint8
	Family uint8
	Data   [14]uint8
***REMOVED***

type RawSockaddrAny struct ***REMOVED***
	Addr RawSockaddr
	_    [112]uint8 // pad
***REMOVED***

type _Socklen uint32

type Linger struct ***REMOVED***
	Onoff  int32
	Linger int32
***REMOVED***

type Iovec struct ***REMOVED***
	Base *byte
	Len  uint64
***REMOVED***

type IPMreq struct ***REMOVED***
	Multiaddr [4]byte /* in_addr */
	Interface [4]byte /* in_addr */
***REMOVED***

type IPv6Mreq struct ***REMOVED***
	Multiaddr [16]byte /* in6_addr */
	Interface uint32
***REMOVED***

type Msghdr struct ***REMOVED***
	Name       *byte
	Iov        *Iovec
	Control    *byte
	Flags      int32
	Namelen    int32
	Iovlen     int32
	Controllen int32
***REMOVED***

type Cmsghdr struct ***REMOVED***
	Len   int32
	Level int32
	Type  int32
***REMOVED***

type Inet4Pktinfo struct ***REMOVED***
	Addr    [4]byte /* in_addr */
	Ifindex uint32
***REMOVED***

type Inet6Pktinfo struct ***REMOVED***
	Addr    [16]byte /* in6_addr */
	Ifindex uint32
***REMOVED***

type IPv6MTUInfo struct ***REMOVED***
	Addr RawSockaddrInet6
	Mtu  uint32
***REMOVED***

type ICMPv6Filter struct ***REMOVED***
	Data [8]uint32
***REMOVED***

type TCPInfo struct ***REMOVED***
	State          uint8
	Ca_state       uint8
	Retransmits    uint8
	Probes         uint8
	Backoff        uint8
	Options        uint8
	Rto            uint32
	Ato            uint32
	Snd_mss        uint32
	Rcv_mss        uint32
	Unacked        uint32
	Sacked         uint32
	Lost           uint32
	Retrans        uint32
	Fackets        uint32
	Last_data_sent uint32
	Last_ack_sent  uint32
	Last_data_recv uint32
	Last_ack_recv  uint32
	Pmtu           uint32
	Rcv_ssthresh   uint32
	Rtt            uint32
	Rttvar         uint32
	Snd_ssthresh   uint32
	Snd_cwnd       uint32
	Advmss         uint32
	Reordering     uint32
	Rcv_rtt        uint32
	Rcv_space      uint32
	Total_retrans  uint32
***REMOVED***

type _Gid_t uint32

type rusage_zos struct ***REMOVED***
	Utime timeval_zos
	Stime timeval_zos
***REMOVED***

type Rusage struct ***REMOVED***
	Utime    Timeval
	Stime    Timeval
	Maxrss   int64
	Ixrss    int64
	Idrss    int64
	Isrss    int64
	Minflt   int64
	Majflt   int64
	Nswap    int64
	Inblock  int64
	Oublock  int64
	Msgsnd   int64
	Msgrcv   int64
	Nsignals int64
	Nvcsw    int64
	Nivcsw   int64
***REMOVED***

type Rlimit struct ***REMOVED***
	Cur uint64
	Max uint64
***REMOVED***

// ***REMOVED*** int, short, short ***REMOVED*** in poll.h
type PollFd struct ***REMOVED***
	Fd      int32
	Events  int16
	Revents int16
***REMOVED***

type Stat_t struct ***REMOVED*** //Linux Definition
	Dev     uint64
	Ino     uint64
	Nlink   uint64
	Mode    uint32
	Uid     uint32
	Gid     uint32
	_       int32
	Rdev    uint64
	Size    int64
	Atim    Timespec
	Mtim    Timespec
	Ctim    Timespec
	Blksize int64
	Blocks  int64
	_       [3]int64
***REMOVED***

type Stat_LE_t struct ***REMOVED***
	_            [4]byte // eye catcher
	Length       uint16
	Version      uint16
	Mode         int32
	Ino          uint32
	Dev          uint32
	Nlink        int32
	Uid          int32
	Gid          int32
	Size         int64
	Atim31       [4]byte
	Mtim31       [4]byte
	Ctim31       [4]byte
	Rdev         uint32
	Auditoraudit uint32
	Useraudit    uint32
	Blksize      int32
	Creatim31    [4]byte
	AuditID      [16]byte
	_            [4]byte // rsrvd1
	File_tag     struct ***REMOVED***
		Ccsid   uint16
		Txtflag uint16 // aggregating Txflag:1 deferred:1 rsvflags:14
	***REMOVED***
	CharsetID [8]byte
	Blocks    int64
	Genvalue  uint32
	Reftim31  [4]byte
	Fid       [8]byte
	Filefmt   byte
	Fspflag2  byte
	_         [2]byte // rsrvd2
	Ctimemsec int32
	Seclabel  [8]byte
	_         [4]byte // rsrvd3
	_         [4]byte // rsrvd4
	Atim      Time_t
	Mtim      Time_t
	Ctim      Time_t
	Creatim   Time_t
	Reftim    Time_t
	_         [24]byte // rsrvd5
***REMOVED***

type Statvfs_t struct ***REMOVED***
	ID          [4]byte
	Len         int32
	Bsize       uint64
	Blocks      uint64
	Usedspace   uint64
	Bavail      uint64
	Flag        uint64
	Maxfilesize int64
	_           [16]byte
	Frsize      uint64
	Bfree       uint64
	Files       uint32
	Ffree       uint32
	Favail      uint32
	Namemax31   uint32
	Invarsec    uint32
	_           [4]byte
	Fsid        uint64
	Namemax     uint64
***REMOVED***

type Statfs_t struct ***REMOVED***
	Type    uint32
	Bsize   uint64
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint32
	Ffree   uint32
	Fsid    uint64
	Namelen uint64
	Frsize  uint64
	Flags   uint64
***REMOVED***

type Dirent struct ***REMOVED***
	Reclen uint16
	Namlen uint16
	Ino    uint32
	Extra  uintptr
	Name   [256]byte
***REMOVED***

type FdSet struct ***REMOVED***
	Bits [64]int32
***REMOVED***

// This struct is packed on z/OS so it can't be used directly.
type Flock_t struct ***REMOVED***
	Type   int16
	Whence int16
	Start  int64
	Len    int64
	Pid    int32
***REMOVED***

type Termios struct ***REMOVED***
	Cflag uint32
	Iflag uint32
	Lflag uint32
	Oflag uint32
	Cc    [11]uint8
***REMOVED***

type Winsize struct ***REMOVED***
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
***REMOVED***

type W_Mnth struct ***REMOVED***
	Hid   [4]byte
	Size  int32
	Cur1  int32 //32bit pointer
	Cur2  int32 //^
	Devno uint32
	_     [4]byte
***REMOVED***

type W_Mntent struct ***REMOVED***
	Fstype       uint32
	Mode         uint32
	Dev          uint32
	Parentdev    uint32
	Rootino      uint32
	Status       byte
	Ddname       [9]byte
	Fstname      [9]byte
	Fsname       [45]byte
	Pathlen      uint32
	Mountpoint   [1024]byte
	Jobname      [8]byte
	PID          int32
	Parmoffset   int32
	Parmlen      int16
	Owner        [8]byte
	Quiesceowner [8]byte
	_            [38]byte
***REMOVED***
