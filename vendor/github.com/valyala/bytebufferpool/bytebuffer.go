package bytebufferpool

import "io"

// ByteBuffer provides byte buffer, which can be used for minimizing
// memory allocations.
//
// ByteBuffer may be used with functions appending data to the given []byte
// slice. See example code for details.
//
// Use Get for obtaining an empty byte buffer.
type ByteBuffer struct ***REMOVED***

	// B is a byte buffer to use in append-like workloads.
	// See example code for details.
	B []byte
***REMOVED***

// Len returns the size of the byte buffer.
func (b *ByteBuffer) Len() int ***REMOVED***
	return len(b.B)
***REMOVED***

// ReadFrom implements io.ReaderFrom.
//
// The function appends all the data read from r to b.
func (b *ByteBuffer) ReadFrom(r io.Reader) (int64, error) ***REMOVED***
	p := b.B
	nStart := int64(len(p))
	nMax := int64(cap(p))
	n := nStart
	if nMax == 0 ***REMOVED***
		nMax = 64
		p = make([]byte, nMax)
	***REMOVED*** else ***REMOVED***
		p = p[:nMax]
	***REMOVED***
	for ***REMOVED***
		if n == nMax ***REMOVED***
			nMax *= 2
			bNew := make([]byte, nMax)
			copy(bNew, p)
			p = bNew
		***REMOVED***
		nn, err := r.Read(p[n:])
		n += int64(nn)
		if err != nil ***REMOVED***
			b.B = p[:n]
			n -= nStart
			if err == io.EOF ***REMOVED***
				return n, nil
			***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***
***REMOVED***

// WriteTo implements io.WriterTo.
func (b *ByteBuffer) WriteTo(w io.Writer) (int64, error) ***REMOVED***
	n, err := w.Write(b.B)
	return int64(n), err
***REMOVED***

// Bytes returns b.B, i.e. all the bytes accumulated in the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
func (b *ByteBuffer) Bytes() []byte ***REMOVED***
	return b.B
***REMOVED***

// Write implements io.Writer - it appends p to ByteBuffer.B
func (b *ByteBuffer) Write(p []byte) (int, error) ***REMOVED***
	b.B = append(b.B, p...)
	return len(p), nil
***REMOVED***

// WriteByte appends the byte c to the buffer.
//
// The purpose of this function is bytes.Buffer compatibility.
//
// The function always returns nil.
func (b *ByteBuffer) WriteByte(c byte) error ***REMOVED***
	b.B = append(b.B, c)
	return nil
***REMOVED***

// WriteString appends s to ByteBuffer.B.
func (b *ByteBuffer) WriteString(s string) (int, error) ***REMOVED***
	b.B = append(b.B, s...)
	return len(s), nil
***REMOVED***

// Set sets ByteBuffer.B to p.
func (b *ByteBuffer) Set(p []byte) ***REMOVED***
	b.B = append(b.B[:0], p...)
***REMOVED***

// SetString sets ByteBuffer.B to s.
func (b *ByteBuffer) SetString(s string) ***REMOVED***
	b.B = append(b.B[:0], s...)
***REMOVED***

// String returns string representation of ByteBuffer.B.
func (b *ByteBuffer) String() string ***REMOVED***
	return string(b.B)
***REMOVED***

// Reset makes ByteBuffer.B empty.
func (b *ByteBuffer) Reset() ***REMOVED***
	b.B = b.B[:0]
***REMOVED***
