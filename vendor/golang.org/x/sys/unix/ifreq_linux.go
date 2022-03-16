// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux
// +build linux

package unix

import (
	"bytes"
	"unsafe"
)

// Helpers for dealing with ifreq since it contains a union and thus requires a
// lot of unsafe.Pointer casts to use properly.

// An Ifreq is a type-safe wrapper around the raw ifreq struct. An Ifreq
// contains an interface name and a union of arbitrary data which can be
// accessed using the Ifreq's methods. To create an Ifreq, use the NewIfreq
// function.
//
// Use the Name method to access the stored interface name. The union data
// fields can be get and set using the following methods:
//   - Uint16/SetUint16: flags
//   - Uint32/SetUint32: ifindex, metric, mtu
type Ifreq struct***REMOVED*** raw ifreq ***REMOVED***

// NewIfreq creates an Ifreq with the input network interface name after
// validating the name does not exceed IFNAMSIZ-1 (trailing NULL required)
// bytes.
func NewIfreq(name string) (*Ifreq, error) ***REMOVED***
	// Leave room for terminating NULL byte.
	if len(name) >= IFNAMSIZ ***REMOVED***
		return nil, EINVAL
	***REMOVED***

	var ifr ifreq
	copy(ifr.Ifrn[:], name)

	return &Ifreq***REMOVED***raw: ifr***REMOVED***, nil
***REMOVED***

// TODO(mdlayher): get/set methods for hardware address sockaddr, char array, etc.

// Name returns the interface name associated with the Ifreq.
func (ifr *Ifreq) Name() string ***REMOVED***
	// BytePtrToString requires a NULL terminator or the program may crash. If
	// one is not present, just return the empty string.
	if !bytes.Contains(ifr.raw.Ifrn[:], []byte***REMOVED***0x00***REMOVED***) ***REMOVED***
		return ""
	***REMOVED***

	return BytePtrToString(&ifr.raw.Ifrn[0])
***REMOVED***

// According to netdevice(7), only AF_INET addresses are returned for numerous
// sockaddr ioctls. For convenience, we expose these as Inet4Addr since the Port
// field and other data is always empty.

// Inet4Addr returns the Ifreq union data from an embedded sockaddr as a C
// in_addr/Go []byte (4-byte IPv4 address) value. If the sockaddr family is not
// AF_INET, an error is returned.
func (ifr *Ifreq) Inet4Addr() ([]byte, error) ***REMOVED***
	raw := *(*RawSockaddrInet4)(unsafe.Pointer(&ifr.raw.Ifru[:SizeofSockaddrInet4][0]))
	if raw.Family != AF_INET ***REMOVED***
		// Cannot safely interpret raw.Addr bytes as an IPv4 address.
		return nil, EINVAL
	***REMOVED***

	return raw.Addr[:], nil
***REMOVED***

// SetInet4Addr sets a C in_addr/Go []byte (4-byte IPv4 address) value in an
// embedded sockaddr within the Ifreq's union data. v must be 4 bytes in length
// or an error will be returned.
func (ifr *Ifreq) SetInet4Addr(v []byte) error ***REMOVED***
	if len(v) != 4 ***REMOVED***
		return EINVAL
	***REMOVED***

	var addr [4]byte
	copy(addr[:], v)

	ifr.clear()
	*(*RawSockaddrInet4)(
		unsafe.Pointer(&ifr.raw.Ifru[:SizeofSockaddrInet4][0]),
	) = RawSockaddrInet4***REMOVED***
		// Always set IP family as ioctls would require it anyway.
		Family: AF_INET,
		Addr:   addr,
	***REMOVED***

	return nil
***REMOVED***

// Uint16 returns the Ifreq union data as a C short/Go uint16 value.
func (ifr *Ifreq) Uint16() uint16 ***REMOVED***
	return *(*uint16)(unsafe.Pointer(&ifr.raw.Ifru[:2][0]))
***REMOVED***

// SetUint16 sets a C short/Go uint16 value as the Ifreq's union data.
func (ifr *Ifreq) SetUint16(v uint16) ***REMOVED***
	ifr.clear()
	*(*uint16)(unsafe.Pointer(&ifr.raw.Ifru[:2][0])) = v
***REMOVED***

// Uint32 returns the Ifreq union data as a C int/Go uint32 value.
func (ifr *Ifreq) Uint32() uint32 ***REMOVED***
	return *(*uint32)(unsafe.Pointer(&ifr.raw.Ifru[:4][0]))
***REMOVED***

// SetUint32 sets a C int/Go uint32 value as the Ifreq's union data.
func (ifr *Ifreq) SetUint32(v uint32) ***REMOVED***
	ifr.clear()
	*(*uint32)(unsafe.Pointer(&ifr.raw.Ifru[:4][0])) = v
***REMOVED***

// clear zeroes the ifreq's union field to prevent trailing garbage data from
// being sent to the kernel if an ifreq is reused.
func (ifr *Ifreq) clear() ***REMOVED***
	for i := range ifr.raw.Ifru ***REMOVED***
		ifr.raw.Ifru[i] = 0
	***REMOVED***
***REMOVED***

// TODO(mdlayher): export as IfreqData? For now we can provide helpers such as
// IoctlGetEthtoolDrvinfo which use these APIs under the hood.

// An ifreqData is an Ifreq which carries pointer data. To produce an ifreqData,
// use the Ifreq.withData method.
type ifreqData struct ***REMOVED***
	name [IFNAMSIZ]byte
	// A type separate from ifreq is required in order to comply with the
	// unsafe.Pointer rules since the "pointer-ness" of data would not be
	// preserved if it were cast into the byte array of a raw ifreq.
	data unsafe.Pointer
	// Pad to the same size as ifreq.
	_ [len(ifreq***REMOVED******REMOVED***.Ifru) - SizeofPtr]byte
***REMOVED***

// withData produces an ifreqData with the pointer p set for ioctls which require
// arbitrary pointer data.
func (ifr Ifreq) withData(p unsafe.Pointer) ifreqData ***REMOVED***
	return ifreqData***REMOVED***
		name: ifr.raw.Ifrn,
		data: p,
	***REMOVED***
***REMOVED***
