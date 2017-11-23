// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

// Socket control messages

package unix

import "unsafe"

// Round the length of a raw sockaddr up to align it properly.
func cmsgAlignOf(salen int) int ***REMOVED***
	salign := sizeofPtr
	// NOTE: It seems like 64-bit Darwin, DragonFly BSD and
	// Solaris kernels still require 32-bit aligned access to
	// network subsystem.
	if darwin64Bit || dragonfly64Bit || solaris64Bit ***REMOVED***
		salign = 4
	***REMOVED***
	return (salen + salign - 1) & ^(salign - 1)
***REMOVED***

// CmsgLen returns the value to store in the Len field of the Cmsghdr
// structure, taking into account any necessary alignment.
func CmsgLen(datalen int) int ***REMOVED***
	return cmsgAlignOf(SizeofCmsghdr) + datalen
***REMOVED***

// CmsgSpace returns the number of bytes an ancillary element with
// payload of the passed data length occupies.
func CmsgSpace(datalen int) int ***REMOVED***
	return cmsgAlignOf(SizeofCmsghdr) + cmsgAlignOf(datalen)
***REMOVED***

func cmsgData(h *Cmsghdr) unsafe.Pointer ***REMOVED***
	return unsafe.Pointer(uintptr(unsafe.Pointer(h)) + uintptr(cmsgAlignOf(SizeofCmsghdr)))
***REMOVED***

// SocketControlMessage represents a socket control message.
type SocketControlMessage struct ***REMOVED***
	Header Cmsghdr
	Data   []byte
***REMOVED***

// ParseSocketControlMessage parses b as an array of socket control
// messages.
func ParseSocketControlMessage(b []byte) ([]SocketControlMessage, error) ***REMOVED***
	var msgs []SocketControlMessage
	i := 0
	for i+CmsgLen(0) <= len(b) ***REMOVED***
		h, dbuf, err := socketControlMessageHeaderAndData(b[i:])
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		m := SocketControlMessage***REMOVED***Header: *h, Data: dbuf***REMOVED***
		msgs = append(msgs, m)
		i += cmsgAlignOf(int(h.Len))
	***REMOVED***
	return msgs, nil
***REMOVED***

func socketControlMessageHeaderAndData(b []byte) (*Cmsghdr, []byte, error) ***REMOVED***
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	if h.Len < SizeofCmsghdr || uint64(h.Len) > uint64(len(b)) ***REMOVED***
		return nil, nil, EINVAL
	***REMOVED***
	return h, b[cmsgAlignOf(SizeofCmsghdr):h.Len], nil
***REMOVED***

// UnixRights encodes a set of open file descriptors into a socket
// control message for sending to another process.
func UnixRights(fds ...int) []byte ***REMOVED***
	datalen := len(fds) * 4
	b := make([]byte, CmsgSpace(datalen))
	h := (*Cmsghdr)(unsafe.Pointer(&b[0]))
	h.Level = SOL_SOCKET
	h.Type = SCM_RIGHTS
	h.SetLen(CmsgLen(datalen))
	data := cmsgData(h)
	for _, fd := range fds ***REMOVED***
		*(*int32)(data) = int32(fd)
		data = unsafe.Pointer(uintptr(data) + 4)
	***REMOVED***
	return b
***REMOVED***

// ParseUnixRights decodes a socket control message that contains an
// integer array of open file descriptors from another process.
func ParseUnixRights(m *SocketControlMessage) ([]int, error) ***REMOVED***
	if m.Header.Level != SOL_SOCKET ***REMOVED***
		return nil, EINVAL
	***REMOVED***
	if m.Header.Type != SCM_RIGHTS ***REMOVED***
		return nil, EINVAL
	***REMOVED***
	fds := make([]int, len(m.Data)>>2)
	for i, j := 0, 0; i < len(m.Data); i += 4 ***REMOVED***
		fds[j] = int(*(*int32)(unsafe.Pointer(&m.Data[i])))
		j++
	***REMOVED***
	return fds, nil
***REMOVED***
