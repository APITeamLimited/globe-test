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
	var bw [1]byte
loop:
	for ***REMOVED***
		c1, err := er.ReadByte()
		if err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c1 != 0x1b ***REMOVED***
			bw[0] = c1
			w.out.Write(bw[:])
			continue
		***REMOVED***
		c2, err := er.ReadByte()
		if err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c2 != 0x5b ***REMOVED***
			continue
		***REMOVED***

		var buf bytes.Buffer
		for ***REMOVED***
			c, err := er.ReadByte()
			if err != nil ***REMOVED***
				break loop
			***REMOVED***
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' ***REMOVED***
				break
			***REMOVED***
			buf.Write([]byte(string(c)))
		***REMOVED***
	***REMOVED***

	return len(data), nil
***REMOVED***
