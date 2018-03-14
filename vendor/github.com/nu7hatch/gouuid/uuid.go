// This package provides immutable UUID structs and the functions
// NewV3, NewV4, NewV5 and Parse() for generating versions 3, 4
// and 5 UUIDs as specified in RFC 4122.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"regexp"
)

// The UUID reserved variants. 
const (
	ReservedNCS       byte = 0x80
	ReservedRFC4122   byte = 0x40
	ReservedMicrosoft byte = 0x20
	ReservedFuture    byte = 0x00
)

// The following standard UUIDs are for use with NewV3() or NewV5().
var (
	NamespaceDNS, _  = ParseHex("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	NamespaceURL, _  = ParseHex("6ba7b811-9dad-11d1-80b4-00c04fd430c8")
	NamespaceOID, _  = ParseHex("6ba7b812-9dad-11d1-80b4-00c04fd430c8")
	NamespaceX500, _ = ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
)

// Pattern used to parse hex string representation of the UUID.
// FIXME: do something to consider both brackets at one time,
// current one allows to parse string with only one opening
// or closing bracket.
const hexPattern = "^(urn\\:uuid\\:)?\\***REMOVED***?([a-z0-9]***REMOVED***8***REMOVED***)-([a-z0-9]***REMOVED***4***REMOVED***)-" +
	"([1-5][a-z0-9]***REMOVED***3***REMOVED***)-([a-z0-9]***REMOVED***4***REMOVED***)-([a-z0-9]***REMOVED***12***REMOVED***)\\***REMOVED***?$"

var re = regexp.MustCompile(hexPattern)

// A UUID representation compliant with specification in
// RFC 4122 document.
type UUID [16]byte

// ParseHex creates a UUID object from given hex string
// representation. Function accepts UUID string in following
// formats:
//
//     uuid.ParseHex("6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//     uuid.ParseHex("***REMOVED***6ba7b814-9dad-11d1-80b4-00c04fd430c8***REMOVED***")
//     uuid.ParseHex("urn:uuid:6ba7b814-9dad-11d1-80b4-00c04fd430c8")
//
func ParseHex(s string) (u *UUID, err error) ***REMOVED***
	md := re.FindStringSubmatch(s)
	if md == nil ***REMOVED***
		err = errors.New("Invalid UUID string")
		return
	***REMOVED***
	hash := md[2] + md[3] + md[4] + md[5] + md[6]
	b, err := hex.DecodeString(hash)
	if err != nil ***REMOVED***
		return
	***REMOVED***
	u = new(UUID)
	copy(u[:], b)
	return
***REMOVED***

// Parse creates a UUID object from given bytes slice.
func Parse(b []byte) (u *UUID, err error) ***REMOVED***
	if len(b) != 16 ***REMOVED***
		err = errors.New("Given slice is not valid UUID sequence")
		return
	***REMOVED***
	u = new(UUID)
	copy(u[:], b)
	return
***REMOVED***

// Generate a UUID based on the MD5 hash of a namespace identifier
// and a name.
func NewV3(ns *UUID, name []byte) (u *UUID, err error) ***REMOVED***
	if ns == nil ***REMOVED***
		err = errors.New("Invalid namespace UUID")
		return
	***REMOVED***
	u = new(UUID)
	// Set all bits to MD5 hash generated from namespace and name.
	u.setBytesFromHash(md5.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(3)
	return
***REMOVED***

// Generate a random UUID.
func NewV4() (u *UUID, err error) ***REMOVED***
	u = new(UUID)
	// Set all bits to randomly (or pseudo-randomly) chosen values.
	_, err = rand.Read(u[:])
	if err != nil ***REMOVED***
		return
	***REMOVED***
	u.setVariant(ReservedRFC4122)
	u.setVersion(4)
	return
***REMOVED***

// Generate a UUID based on the SHA-1 hash of a namespace identifier
// and a name.
func NewV5(ns *UUID, name []byte) (u *UUID, err error) ***REMOVED***
	u = new(UUID)
	// Set all bits to truncated SHA1 hash generated from namespace
	// and name.
	u.setBytesFromHash(sha1.New(), ns[:], name)
	u.setVariant(ReservedRFC4122)
	u.setVersion(5)
	return
***REMOVED***

// Generate a MD5 hash of a namespace and a name, and copy it to the
// UUID slice.
func (u *UUID) setBytesFromHash(hash hash.Hash, ns, name []byte) ***REMOVED***
	hash.Write(ns[:])
	hash.Write(name)
	copy(u[:], hash.Sum([]byte***REMOVED******REMOVED***)[:16])
***REMOVED***

// Set the two most significant bits (bits 6 and 7) of the
// clock_seq_hi_and_reserved to zero and one, respectively.
func (u *UUID) setVariant(v byte) ***REMOVED***
	switch v ***REMOVED***
	case ReservedNCS:
		u[8] = (u[8] | ReservedNCS) & 0xBF
	case ReservedRFC4122:
		u[8] = (u[8] | ReservedRFC4122) & 0x7F
	case ReservedMicrosoft:
		u[8] = (u[8] | ReservedMicrosoft) & 0x3F
	***REMOVED***
***REMOVED***

// Variant returns the UUID Variant, which determines the internal
// layout of the UUID. This will be one of the constants: RESERVED_NCS,
// RFC_4122, RESERVED_MICROSOFT, RESERVED_FUTURE.
func (u *UUID) Variant() byte ***REMOVED***
	if u[8]&ReservedNCS == ReservedNCS ***REMOVED***
		return ReservedNCS
	***REMOVED*** else if u[8]&ReservedRFC4122 == ReservedRFC4122 ***REMOVED***
		return ReservedRFC4122
	***REMOVED*** else if u[8]&ReservedMicrosoft == ReservedMicrosoft ***REMOVED***
		return ReservedMicrosoft
	***REMOVED***
	return ReservedFuture
***REMOVED***

// Set the four most significant bits (bits 12 through 15) of the
// time_hi_and_version field to the 4-bit version number.
func (u *UUID) setVersion(v byte) ***REMOVED***
	u[6] = (u[6] & 0xF) | (v << 4)
***REMOVED***

// Version returns a version number of the algorithm used to
// generate the UUID sequence.
func (u *UUID) Version() uint ***REMOVED***
	return uint(u[6] >> 4)
***REMOVED***

// Returns unparsed version of the generated UUID sequence.
func (u *UUID) String() string ***REMOVED***
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
***REMOVED***
