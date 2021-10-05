// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Plan 9 directory marshalling. See intro(5).

package plan9

import "errors"

var (
	ErrShortStat = errors.New("stat buffer too short")
	ErrBadStat   = errors.New("malformed stat buffer")
	ErrBadName   = errors.New("bad character in file name")
)

// A Qid represents a 9P server's unique identification for a file.
type Qid struct ***REMOVED***
	Path uint64 // the file server's unique identification for the file
	Vers uint32 // version number for given Path
	Type uint8  // the type of the file (plan9.QTDIR for example)
***REMOVED***

// A Dir contains the metadata for a file.
type Dir struct ***REMOVED***
	// system-modified data
	Type uint16 // server type
	Dev  uint32 // server subtype

	// file data
	Qid    Qid    // unique id from server
	Mode   uint32 // permissions
	Atime  uint32 // last read time
	Mtime  uint32 // last write time
	Length int64  // file length
	Name   string // last element of path
	Uid    string // owner name
	Gid    string // group name
	Muid   string // last modifier name
***REMOVED***

var nullDir = Dir***REMOVED***
	Type: ^uint16(0),
	Dev:  ^uint32(0),
	Qid: Qid***REMOVED***
		Path: ^uint64(0),
		Vers: ^uint32(0),
		Type: ^uint8(0),
	***REMOVED***,
	Mode:   ^uint32(0),
	Atime:  ^uint32(0),
	Mtime:  ^uint32(0),
	Length: ^int64(0),
***REMOVED***

// Null assigns special "don't touch" values to members of d to
// avoid modifying them during plan9.Wstat.
func (d *Dir) Null() ***REMOVED*** *d = nullDir ***REMOVED***

// Marshal encodes a 9P stat message corresponding to d into b
//
// If there isn't enough space in b for a stat message, ErrShortStat is returned.
func (d *Dir) Marshal(b []byte) (n int, err error) ***REMOVED***
	n = STATFIXLEN + len(d.Name) + len(d.Uid) + len(d.Gid) + len(d.Muid)
	if n > len(b) ***REMOVED***
		return n, ErrShortStat
	***REMOVED***

	for _, c := range d.Name ***REMOVED***
		if c == '/' ***REMOVED***
			return n, ErrBadName
		***REMOVED***
	***REMOVED***

	b = pbit16(b, uint16(n)-2)
	b = pbit16(b, d.Type)
	b = pbit32(b, d.Dev)
	b = pbit8(b, d.Qid.Type)
	b = pbit32(b, d.Qid.Vers)
	b = pbit64(b, d.Qid.Path)
	b = pbit32(b, d.Mode)
	b = pbit32(b, d.Atime)
	b = pbit32(b, d.Mtime)
	b = pbit64(b, uint64(d.Length))
	b = pstring(b, d.Name)
	b = pstring(b, d.Uid)
	b = pstring(b, d.Gid)
	b = pstring(b, d.Muid)

	return n, nil
***REMOVED***

// UnmarshalDir decodes a single 9P stat message from b and returns the resulting Dir.
//
// If b is too small to hold a valid stat message, ErrShortStat is returned.
//
// If the stat message itself is invalid, ErrBadStat is returned.
func UnmarshalDir(b []byte) (*Dir, error) ***REMOVED***
	if len(b) < STATFIXLEN ***REMOVED***
		return nil, ErrShortStat
	***REMOVED***
	size, buf := gbit16(b)
	if len(b) != int(size)+2 ***REMOVED***
		return nil, ErrBadStat
	***REMOVED***
	b = buf

	var d Dir
	d.Type, b = gbit16(b)
	d.Dev, b = gbit32(b)
	d.Qid.Type, b = gbit8(b)
	d.Qid.Vers, b = gbit32(b)
	d.Qid.Path, b = gbit64(b)
	d.Mode, b = gbit32(b)
	d.Atime, b = gbit32(b)
	d.Mtime, b = gbit32(b)

	n, b := gbit64(b)
	d.Length = int64(n)

	var ok bool
	if d.Name, b, ok = gstring(b); !ok ***REMOVED***
		return nil, ErrBadStat
	***REMOVED***
	if d.Uid, b, ok = gstring(b); !ok ***REMOVED***
		return nil, ErrBadStat
	***REMOVED***
	if d.Gid, b, ok = gstring(b); !ok ***REMOVED***
		return nil, ErrBadStat
	***REMOVED***
	if d.Muid, b, ok = gstring(b); !ok ***REMOVED***
		return nil, ErrBadStat
	***REMOVED***

	return &d, nil
***REMOVED***

// pbit8 copies the 8-bit number v to b and returns the remaining slice of b.
func pbit8(b []byte, v uint8) []byte ***REMOVED***
	b[0] = byte(v)
	return b[1:]
***REMOVED***

// pbit16 copies the 16-bit number v to b in little-endian order and returns the remaining slice of b.
func pbit16(b []byte, v uint16) []byte ***REMOVED***
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	return b[2:]
***REMOVED***

// pbit32 copies the 32-bit number v to b in little-endian order and returns the remaining slice of b.
func pbit32(b []byte, v uint32) []byte ***REMOVED***
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	return b[4:]
***REMOVED***

// pbit64 copies the 64-bit number v to b in little-endian order and returns the remaining slice of b.
func pbit64(b []byte, v uint64) []byte ***REMOVED***
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
	return b[8:]
***REMOVED***

// pstring copies the string s to b, prepending it with a 16-bit length in little-endian order, and
// returning the remaining slice of b..
func pstring(b []byte, s string) []byte ***REMOVED***
	b = pbit16(b, uint16(len(s)))
	n := copy(b, s)
	return b[n:]
***REMOVED***

// gbit8 reads an 8-bit number from b and returns it with the remaining slice of b.
func gbit8(b []byte) (uint8, []byte) ***REMOVED***
	return uint8(b[0]), b[1:]
***REMOVED***

// gbit16 reads a 16-bit number in little-endian order from b and returns it with the remaining slice of b.
func gbit16(b []byte) (uint16, []byte) ***REMOVED***
	return uint16(b[0]) | uint16(b[1])<<8, b[2:]
***REMOVED***

// gbit32 reads a 32-bit number in little-endian order from b and returns it with the remaining slice of b.
func gbit32(b []byte) (uint32, []byte) ***REMOVED***
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, b[4:]
***REMOVED***

// gbit64 reads a 64-bit number in little-endian order from b and returns it with the remaining slice of b.
func gbit64(b []byte) (uint64, []byte) ***REMOVED***
	lo := uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24
	hi := uint32(b[4]) | uint32(b[5])<<8 | uint32(b[6])<<16 | uint32(b[7])<<24
	return uint64(lo) | uint64(hi)<<32, b[8:]
***REMOVED***

// gstring reads a string from b, prefixed with a 16-bit length in little-endian order.
// It returns the string with the remaining slice of b and a boolean. If the length is
// greater than the number of bytes in b, the boolean will be false.
func gstring(b []byte) (string, []byte, bool) ***REMOVED***
	n, b := gbit16(b)
	if int(n) > len(b) ***REMOVED***
		return "", b, false
	***REMOVED***
	return string(b[:n]), b[n:], true
***REMOVED***
