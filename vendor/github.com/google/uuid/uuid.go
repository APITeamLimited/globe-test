// Copyright 2018 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

// A UUID is a 128 bit (16 byte) Universal Unique IDentifier as defined in RFC
// 4122.
type UUID [16]byte

// A Version represents a UUID's version.
type Version byte

// A Variant represents a UUID's variant.
type Variant byte

// Constants returned by Variant.
const (
	Invalid   = Variant(iota) // Invalid UUID
	RFC4122                   // The variant specified in RFC4122
	Reserved                  // Reserved, NCS backward compatibility.
	Microsoft                 // Reserved, Microsoft Corporation backward compatibility.
	Future                    // Reserved for future definition.
)

const randPoolSize = 16 * 16

var (
	rander      = rand.Reader // random function
	poolEnabled = false
	poolMu      sync.Mutex
	poolPos     = randPoolSize     // protected with poolMu
	pool        [randPoolSize]byte // protected with poolMu
)

type invalidLengthError struct***REMOVED*** len int ***REMOVED***

func (err invalidLengthError) Error() string ***REMOVED***
	return fmt.Sprintf("invalid UUID length: %d", err.len)
***REMOVED***

// IsInvalidLengthError is matcher function for custom error invalidLengthError
func IsInvalidLengthError(err error) bool ***REMOVED***
	_, ok := err.(invalidLengthError)
	return ok
***REMOVED***

// Parse decodes s into a UUID or returns an error.  Both the standard UUID
// forms of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx and
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx are decoded as well as the
// Microsoft encoding ***REMOVED***xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx***REMOVED*** and the raw hex
// encoding: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.
func Parse(s string) (UUID, error) ***REMOVED***
	var uuid UUID
	switch len(s) ***REMOVED***
	// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36:

	// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9:
		if strings.ToLower(s[:9]) != "urn:uuid:" ***REMOVED***
			return uuid, fmt.Errorf("invalid urn prefix: %q", s[:9])
		***REMOVED***
		s = s[9:]

	// ***REMOVED***xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx***REMOVED***
	case 36 + 2:
		s = s[1:]

	// xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
	case 32:
		var ok bool
		for i := range uuid ***REMOVED***
			uuid[i], ok = xtob(s[i*2], s[i*2+1])
			if !ok ***REMOVED***
				return uuid, errors.New("invalid UUID format")
			***REMOVED***
		***REMOVED***
		return uuid, nil
	default:
		return uuid, invalidLengthError***REMOVED***len(s)***REMOVED***
	***REMOVED***
	// s is now at least 36 bytes long
	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if s[8] != '-' || s[13] != '-' || s[18] != '-' || s[23] != '-' ***REMOVED***
		return uuid, errors.New("invalid UUID format")
	***REMOVED***
	for i, x := range [16]int***REMOVED***
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34***REMOVED*** ***REMOVED***
		v, ok := xtob(s[x], s[x+1])
		if !ok ***REMOVED***
			return uuid, errors.New("invalid UUID format")
		***REMOVED***
		uuid[i] = v
	***REMOVED***
	return uuid, nil
***REMOVED***

// ParseBytes is like Parse, except it parses a byte slice instead of a string.
func ParseBytes(b []byte) (UUID, error) ***REMOVED***
	var uuid UUID
	switch len(b) ***REMOVED***
	case 36: // xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	case 36 + 9: // urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
		if !bytes.Equal(bytes.ToLower(b[:9]), []byte("urn:uuid:")) ***REMOVED***
			return uuid, fmt.Errorf("invalid urn prefix: %q", b[:9])
		***REMOVED***
		b = b[9:]
	case 36 + 2: // ***REMOVED***xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx***REMOVED***
		b = b[1:]
	case 32: // xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
		var ok bool
		for i := 0; i < 32; i += 2 ***REMOVED***
			uuid[i/2], ok = xtob(b[i], b[i+1])
			if !ok ***REMOVED***
				return uuid, errors.New("invalid UUID format")
			***REMOVED***
		***REMOVED***
		return uuid, nil
	default:
		return uuid, invalidLengthError***REMOVED***len(b)***REMOVED***
	***REMOVED***
	// s is now at least 36 bytes long
	// it must be of the form  xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if b[8] != '-' || b[13] != '-' || b[18] != '-' || b[23] != '-' ***REMOVED***
		return uuid, errors.New("invalid UUID format")
	***REMOVED***
	for i, x := range [16]int***REMOVED***
		0, 2, 4, 6,
		9, 11,
		14, 16,
		19, 21,
		24, 26, 28, 30, 32, 34***REMOVED*** ***REMOVED***
		v, ok := xtob(b[x], b[x+1])
		if !ok ***REMOVED***
			return uuid, errors.New("invalid UUID format")
		***REMOVED***
		uuid[i] = v
	***REMOVED***
	return uuid, nil
***REMOVED***

// MustParse is like Parse but panics if the string cannot be parsed.
// It simplifies safe initialization of global variables holding compiled UUIDs.
func MustParse(s string) UUID ***REMOVED***
	uuid, err := Parse(s)
	if err != nil ***REMOVED***
		panic(`uuid: Parse(` + s + `): ` + err.Error())
	***REMOVED***
	return uuid
***REMOVED***

// FromBytes creates a new UUID from a byte slice. Returns an error if the slice
// does not have a length of 16. The bytes are copied from the slice.
func FromBytes(b []byte) (uuid UUID, err error) ***REMOVED***
	err = uuid.UnmarshalBinary(b)
	return uuid, err
***REMOVED***

// Must returns uuid if err is nil and panics otherwise.
func Must(uuid UUID, err error) UUID ***REMOVED***
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return uuid
***REMOVED***

// String returns the string form of uuid, xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// , or "" if uuid is invalid.
func (uuid UUID) String() string ***REMOVED***
	var buf [36]byte
	encodeHex(buf[:], uuid)
	return string(buf[:])
***REMOVED***

// URN returns the RFC 2141 URN form of uuid,
// urn:uuid:xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx,  or "" if uuid is invalid.
func (uuid UUID) URN() string ***REMOVED***
	var buf [36 + 9]byte
	copy(buf[:], "urn:uuid:")
	encodeHex(buf[9:], uuid)
	return string(buf[:])
***REMOVED***

func encodeHex(dst []byte, uuid UUID) ***REMOVED***
	hex.Encode(dst, uuid[:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], uuid[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], uuid[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], uuid[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:], uuid[10:])
***REMOVED***

// Variant returns the variant encoded in uuid.
func (uuid UUID) Variant() Variant ***REMOVED***
	switch ***REMOVED***
	case (uuid[8] & 0xc0) == 0x80:
		return RFC4122
	case (uuid[8] & 0xe0) == 0xc0:
		return Microsoft
	case (uuid[8] & 0xe0) == 0xe0:
		return Future
	default:
		return Reserved
	***REMOVED***
***REMOVED***

// Version returns the version of uuid.
func (uuid UUID) Version() Version ***REMOVED***
	return Version(uuid[6] >> 4)
***REMOVED***

func (v Version) String() string ***REMOVED***
	if v > 15 ***REMOVED***
		return fmt.Sprintf("BAD_VERSION_%d", v)
	***REMOVED***
	return fmt.Sprintf("VERSION_%d", v)
***REMOVED***

func (v Variant) String() string ***REMOVED***
	switch v ***REMOVED***
	case RFC4122:
		return "RFC4122"
	case Reserved:
		return "Reserved"
	case Microsoft:
		return "Microsoft"
	case Future:
		return "Future"
	case Invalid:
		return "Invalid"
	***REMOVED***
	return fmt.Sprintf("BadVariant%d", int(v))
***REMOVED***

// SetRand sets the random number generator to r, which implements io.Reader.
// If r.Read returns an error when the package requests random data then
// a panic will be issued.
//
// Calling SetRand with nil sets the random number generator to the default
// generator.
func SetRand(r io.Reader) ***REMOVED***
	if r == nil ***REMOVED***
		rander = rand.Reader
		return
	***REMOVED***
	rander = r
***REMOVED***

// EnableRandPool enables internal randomness pool used for Random
// (Version 4) UUID generation. The pool contains random bytes read from
// the random number generator on demand in batches. Enabling the pool
// may improve the UUID generation throughput significantly.
//
// Since the pool is stored on the Go heap, this feature may be a bad fit
// for security sensitive applications.
//
// Both EnableRandPool and DisableRandPool are not thread-safe and should
// only be called when there is no possibility that New or any other
// UUID Version 4 generation function will be called concurrently.
func EnableRandPool() ***REMOVED***
	poolEnabled = true
***REMOVED***

// DisableRandPool disables the randomness pool if it was previously
// enabled with EnableRandPool.
//
// Both EnableRandPool and DisableRandPool are not thread-safe and should
// only be called when there is no possibility that New or any other
// UUID Version 4 generation function will be called concurrently.
func DisableRandPool() ***REMOVED***
	poolEnabled = false
	defer poolMu.Unlock()
	poolMu.Lock()
	poolPos = randPoolSize
***REMOVED***
