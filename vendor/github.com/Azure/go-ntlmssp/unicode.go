package ntlmssp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"unicode/utf16"
)

// helper func's for dealing with Windows Unicode (UTF16LE)

func fromUnicode(d []byte) (string, error) ***REMOVED***
	if len(d)%2 > 0 ***REMOVED***
		return "", errors.New("Unicode (UTF 16 LE) specified, but uneven data length")
	***REMOVED***
	s := make([]uint16, len(d)/2)
	err := binary.Read(bytes.NewReader(d), binary.LittleEndian, &s)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(utf16.Decode(s)), nil
***REMOVED***

func toUnicode(s string) []byte ***REMOVED***
	uints := utf16.Encode([]rune(s))
	b := bytes.Buffer***REMOVED******REMOVED***
	binary.Write(&b, binary.LittleEndian, &uints)
	return b.Bytes()
***REMOVED***
