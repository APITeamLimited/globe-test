// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build freebsd || netbsd
// +build freebsd netbsd

package unix

import (
	"strings"
	"unsafe"
)

// Derive extattr namespace and attribute name

func xattrnamespace(fullattr string) (ns int, attr string, err error) ***REMOVED***
	s := strings.IndexByte(fullattr, '.')
	if s == -1 ***REMOVED***
		return -1, "", ENOATTR
	***REMOVED***

	namespace := fullattr[0:s]
	attr = fullattr[s+1:]

	switch namespace ***REMOVED***
	case "user":
		return EXTATTR_NAMESPACE_USER, attr, nil
	case "system":
		return EXTATTR_NAMESPACE_SYSTEM, attr, nil
	default:
		return -1, "", ENOATTR
	***REMOVED***
***REMOVED***

func initxattrdest(dest []byte, idx int) (d unsafe.Pointer) ***REMOVED***
	if len(dest) > idx ***REMOVED***
		return unsafe.Pointer(&dest[idx])
	***REMOVED*** else ***REMOVED***
		return unsafe.Pointer(_zero)
	***REMOVED***
***REMOVED***

// FreeBSD and NetBSD implement their own syscalls to handle extended attributes

func Getxattr(file string, attr string, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	return ExtattrGetFile(file, nsid, a, uintptr(d), destsize)
***REMOVED***

func Fgetxattr(fd int, attr string, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	return ExtattrGetFd(fd, nsid, a, uintptr(d), destsize)
***REMOVED***

func Lgetxattr(link string, attr string, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsize := len(dest)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return -1, err
	***REMOVED***

	return ExtattrGetLink(link, nsid, a, uintptr(d), destsize)
***REMOVED***

// flags are unused on FreeBSD

func Fsetxattr(fd int, attr string, data []byte, flags int) (err error) ***REMOVED***
	var d unsafe.Pointer
	if len(data) > 0 ***REMOVED***
		d = unsafe.Pointer(&data[0])
	***REMOVED***
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = ExtattrSetFd(fd, nsid, a, uintptr(d), datasiz)
	return
***REMOVED***

func Setxattr(file string, attr string, data []byte, flags int) (err error) ***REMOVED***
	var d unsafe.Pointer
	if len(data) > 0 ***REMOVED***
		d = unsafe.Pointer(&data[0])
	***REMOVED***
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = ExtattrSetFile(file, nsid, a, uintptr(d), datasiz)
	return
***REMOVED***

func Lsetxattr(link string, attr string, data []byte, flags int) (err error) ***REMOVED***
	var d unsafe.Pointer
	if len(data) > 0 ***REMOVED***
		d = unsafe.Pointer(&data[0])
	***REMOVED***
	datasiz := len(data)

	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	_, err = ExtattrSetLink(link, nsid, a, uintptr(d), datasiz)
	return
***REMOVED***

func Removexattr(file string, attr string) (err error) ***REMOVED***
	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = ExtattrDeleteFile(file, nsid, a)
	return
***REMOVED***

func Fremovexattr(fd int, attr string) (err error) ***REMOVED***
	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = ExtattrDeleteFd(fd, nsid, a)
	return
***REMOVED***

func Lremovexattr(link string, attr string) (err error) ***REMOVED***
	nsid, a, err := xattrnamespace(attr)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	err = ExtattrDeleteLink(link, nsid, a)
	return
***REMOVED***

func Listxattr(file string, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	// FreeBSD won't allow you to list xattrs from multiple namespaces
	s := 0
	for _, nsid := range [...]int***REMOVED***EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM***REMOVED*** ***REMOVED***
		stmp, e := ExtattrListFile(file, nsid, uintptr(d), destsiz)

		/* Errors accessing system attrs are ignored so that
		 * we can implement the Linux-like behavior of omitting errors that
		 * we don't have read permissions on
		 *
		 * Linux will still error if we ask for user attributes on a file that
		 * we don't have read permissions on, so don't ignore those errors
		 */
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER ***REMOVED***
			continue
		***REMOVED*** else if e != nil ***REMOVED***
			return s, e
		***REMOVED***

		s += stmp
		destsiz -= s
		if destsiz < 0 ***REMOVED***
			destsiz = 0
		***REMOVED***
		d = initxattrdest(dest, s)
	***REMOVED***

	return s, nil
***REMOVED***

func Flistxattr(fd int, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	s := 0
	for _, nsid := range [...]int***REMOVED***EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM***REMOVED*** ***REMOVED***
		stmp, e := ExtattrListFd(fd, nsid, uintptr(d), destsiz)
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER ***REMOVED***
			continue
		***REMOVED*** else if e != nil ***REMOVED***
			return s, e
		***REMOVED***

		s += stmp
		destsiz -= s
		if destsiz < 0 ***REMOVED***
			destsiz = 0
		***REMOVED***
		d = initxattrdest(dest, s)
	***REMOVED***

	return s, nil
***REMOVED***

func Llistxattr(link string, dest []byte) (sz int, err error) ***REMOVED***
	d := initxattrdest(dest, 0)
	destsiz := len(dest)

	s := 0
	for _, nsid := range [...]int***REMOVED***EXTATTR_NAMESPACE_USER, EXTATTR_NAMESPACE_SYSTEM***REMOVED*** ***REMOVED***
		stmp, e := ExtattrListLink(link, nsid, uintptr(d), destsiz)
		if e != nil && e == EPERM && nsid != EXTATTR_NAMESPACE_USER ***REMOVED***
			continue
		***REMOVED*** else if e != nil ***REMOVED***
			return s, e
		***REMOVED***

		s += stmp
		destsiz -= s
		if destsiz < 0 ***REMOVED***
			destsiz = 0
		***REMOVED***
		d = initxattrdest(dest, s)
	***REMOVED***

	return s, nil
***REMOVED***
