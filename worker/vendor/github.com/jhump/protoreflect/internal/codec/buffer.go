package codec

import (
	"fmt"
	"io"
)

// Buffer is a reader and a writer that wraps a slice of bytes and also
// provides API for decoding and encoding the protobuf binary format.
//
// Its operation is similar to that of a bytes.Buffer: writing pushes
// data to the end of the buffer while reading pops data from the head
// of the buffer. So the same buffer can be used to both read and write.
type Buffer struct ***REMOVED***
	buf   []byte
	index int

	// tmp is used when another byte slice is needed, such as when
	// serializing messages, since we need to know the length before
	// we can write the length prefix; by caching this, including
	// after it is grown by serialization operations, we reduce the
	// number of allocations needed
	tmp []byte

	deterministic bool
***REMOVED***

// NewBuffer creates a new buffer with the given slice of bytes as the
// buffer's initial contents.
func NewBuffer(buf []byte) *Buffer ***REMOVED***
	return &Buffer***REMOVED***buf: buf***REMOVED***
***REMOVED***

// SetDeterministic sets this buffer to encode messages deterministically. This
// is useful for tests. But the overhead is non-zero, so it should not likely be
// used outside of tests. When true, map fields in a message must have their
// keys sorted before serialization to ensure deterministic output. Otherwise,
// values in a map field will be serialized in map iteration order.
func (cb *Buffer) SetDeterministic(deterministic bool) ***REMOVED***
	cb.deterministic = deterministic
***REMOVED***

// IsDeterministic returns whether or not this buffer is configured to encode
// messages deterministically.
func (cb *Buffer) IsDeterministic() bool ***REMOVED***
	return cb.deterministic
***REMOVED***

// Reset resets this buffer back to empty. Any subsequent writes/encodes
// to the buffer will allocate a new backing slice of bytes.
func (cb *Buffer) Reset() ***REMOVED***
	cb.buf = []byte(nil)
	cb.index = 0
***REMOVED***

// Bytes returns the slice of bytes remaining in the buffer. Note that
// this does not perform a copy: if the contents of the returned slice
// are modified, the modifications will be visible to subsequent reads
// via the buffer.
func (cb *Buffer) Bytes() []byte ***REMOVED***
	return cb.buf[cb.index:]
***REMOVED***

// String returns the remaining bytes in the buffer as a string.
func (cb *Buffer) String() string ***REMOVED***
	return string(cb.Bytes())
***REMOVED***

// EOF returns true if there are no more bytes remaining to read.
func (cb *Buffer) EOF() bool ***REMOVED***
	return cb.index >= len(cb.buf)
***REMOVED***

// Skip attempts to skip the given number of bytes in the input. If
// the input has fewer bytes than the given count, io.ErrUnexpectedEOF
// is returned and the buffer is unchanged. Otherwise, the given number
// of bytes are skipped and nil is returned.
func (cb *Buffer) Skip(count int) error ***REMOVED***
	if count < 0 ***REMOVED***
		return fmt.Errorf("proto: bad byte length %d", count)
	***REMOVED***
	newIndex := cb.index + count
	if newIndex < cb.index || newIndex > len(cb.buf) ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	cb.index = newIndex
	return nil
***REMOVED***

// Len returns the remaining number of bytes in the buffer.
func (cb *Buffer) Len() int ***REMOVED***
	return len(cb.buf) - cb.index
***REMOVED***

// Read implements the io.Reader interface. If there are no bytes
// remaining in the buffer, it will return 0, io.EOF. Otherwise,
// it reads max(len(dest), cb.Len()) bytes from input and copies
// them into dest. It returns the number of bytes copied and a nil
// error in this case.
func (cb *Buffer) Read(dest []byte) (int, error) ***REMOVED***
	if cb.index == len(cb.buf) ***REMOVED***
		return 0, io.EOF
	***REMOVED***
	copied := copy(dest, cb.buf[cb.index:])
	cb.index += copied
	return copied, nil
***REMOVED***

var _ io.Reader = (*Buffer)(nil)

// Write implements the io.Writer interface. It always returns
// len(data), nil.
func (cb *Buffer) Write(data []byte) (int, error) ***REMOVED***
	cb.buf = append(cb.buf, data...)
	return len(data), nil
***REMOVED***

var _ io.Writer = (*Buffer)(nil)
