package ntlmssp

import (
	"errors"
)

type varField struct ***REMOVED***
	Len          uint16
	MaxLen       uint16
	BufferOffset uint32
***REMOVED***

func (f varField) ReadFrom(buffer []byte) ([]byte, error) ***REMOVED***
	if len(buffer) < int(f.BufferOffset+uint32(f.Len)) ***REMOVED***
		return nil, errors.New("Error reading data, varField extends beyond buffer")
	***REMOVED***
	return buffer[f.BufferOffset : f.BufferOffset+uint32(f.Len)], nil
***REMOVED***

func (f varField) ReadStringFrom(buffer []byte, unicode bool) (string, error) ***REMOVED***
	d, err := f.ReadFrom(buffer)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	if unicode ***REMOVED*** // UTF-16LE encoding scheme
		return fromUnicode(d)
	***REMOVED***
	// OEM encoding, close enough to ASCII, since no code page is specified
	return string(d), err
***REMOVED***

func newVarField(ptr *int, fieldsize int) varField ***REMOVED***
	f := varField***REMOVED***
		Len:          uint16(fieldsize),
		MaxLen:       uint16(fieldsize),
		BufferOffset: uint32(*ptr),
	***REMOVED***
	*ptr += fieldsize
	return f
***REMOVED***
