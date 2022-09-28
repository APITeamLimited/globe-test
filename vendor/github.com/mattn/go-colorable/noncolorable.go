package colorable

import (
	"bytes"
	"io"
)

// NonColorable holds writer but removes escape sequence.
type NonColorable struct ***REMOVED***
	out io.Writer
***REMOVED***

// NewNonColorable returns new instance of Writer which removes escape sequence from Writer.
func NewNonColorable(w io.Writer) io.Writer ***REMOVED***
	return &NonColorable***REMOVED***out: w***REMOVED***
***REMOVED***

// Write writes data on console
func (w *NonColorable) Write(data []byte) (n int, err error) ***REMOVED***
	er := bytes.NewReader(data)
	var plaintext bytes.Buffer
loop:
	for ***REMOVED***
		c1, err := er.ReadByte()
		if err != nil ***REMOVED***
			plaintext.WriteTo(w.out)
			break loop
		***REMOVED***
		if c1 != 0x1b ***REMOVED***
			plaintext.WriteByte(c1)
			continue
		***REMOVED***
		_, err = plaintext.WriteTo(w.out)
		if err != nil ***REMOVED***
			break loop
		***REMOVED***
		c2, err := er.ReadByte()
		if err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c2 != 0x5b ***REMOVED***
			continue
		***REMOVED***

		for ***REMOVED***
			c, err := er.ReadByte()
			if err != nil ***REMOVED***
				break loop
			***REMOVED***
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return len(data), nil
***REMOVED***
