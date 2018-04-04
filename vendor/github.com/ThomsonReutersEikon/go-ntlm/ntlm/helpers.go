//Copyright 2013 Thomson Reuters Global Resources. BSD License please see License file for more information

package ntlm

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"unicode/utf16"
)

// Concatenate two byte slices into a new slice
func concat(ar ...[]byte) []byte ***REMOVED***
	return bytes.Join(ar, nil)
***REMOVED***

// Create a 0 initialized slice of bytes
func zeroBytes(length int) []byte ***REMOVED***
	return make([]byte, length, length)
***REMOVED***

func randomBytes(length int) []byte ***REMOVED***
	randombytes := make([]byte, length)
	_, err := rand.Read(randombytes)
	if err != nil ***REMOVED***
	***REMOVED*** // TODO: What to do with err here
	return randombytes
***REMOVED***

// Zero pad the input byte slice to the given size
// bytes - input byte slice
// offset - where to start taking the bytes from the input slice
// size - size of the output byte slize
func zeroPaddedBytes(bytes []byte, offset int, size int) []byte ***REMOVED***
	newSlice := zeroBytes(size)
	for i := 0; i < size && i+offset < len(bytes); i++ ***REMOVED***
		newSlice[i] = bytes[i+offset]
	***REMOVED***
	return newSlice
***REMOVED***

func MacsEqual(slice1, slice2 []byte) bool ***REMOVED***
	if len(slice1) != len(slice2) ***REMOVED***
		return false
	***REMOVED***
	for i := 0; i < len(slice1); i++ ***REMOVED***
		// bytes between 4 and 7 (inclusive) contains random
		// data that should be ignored while comparing the
		// macs
		if (i < 4 || i > 7) && slice1[i] != slice2[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

func utf16FromString(s string) []byte ***REMOVED***
	encoded := utf16.Encode([]rune(s))
	// TODO: I'm sure there is an easier way to do the conversion from utf16 to bytes
	result := zeroBytes(len(encoded) * 2)
	for i := 0; i < len(encoded); i++ ***REMOVED***
		result[i*2] = byte(encoded[i])
		result[i*2+1] = byte(encoded[i] << 8)
	***REMOVED***
	return result
***REMOVED***

// Convert a UTF16 string to UTF8 string for Go usage
func utf16ToString(bytes []byte) string ***REMOVED***
	var data []uint16

	// NOTE: This is definitely not the best way to do this, but when I tried using a buffer.Read I could not get it to work
	for offset := 0; offset < len(bytes); offset = offset + 2 ***REMOVED***
		i := binary.LittleEndian.Uint16(bytes[offset : offset+2])
		data = append(data, i)
	***REMOVED***

	return string(utf16.Decode(data))
***REMOVED***

func uint32ToBytes(v uint32) []byte ***REMOVED***
	bytes := make([]byte, 4)
	bytes[0] = byte(v & 0xff)
	bytes[1] = byte((v >> 8) & 0xff)
	bytes[2] = byte((v >> 16) & 0xff)
	bytes[3] = byte((v >> 24) & 0xff)
	return bytes
***REMOVED***
