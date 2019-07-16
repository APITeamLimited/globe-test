package brotli

import (
	"errors"
	"io"
)

type decodeError int

func (err decodeError) Error() string ***REMOVED***
	return "brotli: " + string(decoderErrorString(int(err)))
***REMOVED***

var errExcessiveInput = errors.New("brotli: excessive input")
var errInvalidState = errors.New("brotli: invalid state")

// readBufSize is a "good" buffer size that avoids excessive round-trips
// between C and Go but doesn't waste too much memory on buffering.
// It is arbitrarily chosen to be equal to the constant used in io.Copy.
const readBufSize = 32 * 1024

// NewReader creates a new Reader reading the given reader.
func NewReader(src io.Reader) *Reader ***REMOVED***
	r := new(Reader)
	r.Reset(src)
	return r
***REMOVED***

// Reset discards the Reader's state and makes it equivalent to the result of
// its original state from NewReader, but writing to src instead.
// This permits reusing a Reader rather than allocating a new one.
// Error is always nil
func (r *Reader) Reset(src io.Reader) error ***REMOVED***
	decoderStateInit(r)
	r.src = src
	r.buf = make([]byte, readBufSize)
	return nil
***REMOVED***

func (r *Reader) Read(p []byte) (n int, err error) ***REMOVED***
	if !decoderHasMoreOutput(r) && len(r.in) == 0 ***REMOVED***
		m, readErr := r.src.Read(r.buf)
		if m == 0 ***REMOVED***
			// If readErr is `nil`, we just proxy underlying stream behavior.
			return 0, readErr
		***REMOVED***
		r.in = r.buf[:m]
	***REMOVED***

	if len(p) == 0 ***REMOVED***
		return 0, nil
	***REMOVED***

	for ***REMOVED***
		var written uint
		in_len := uint(len(r.in))
		out_len := uint(len(p))
		in_remaining := in_len
		out_remaining := out_len
		result := decoderDecompressStream(r, &in_remaining, &r.in, &out_remaining, &p)
		written = out_len - out_remaining
		n = int(written)

		switch result ***REMOVED***
		case decoderResultSuccess:
			if len(r.in) > 0 ***REMOVED***
				return n, errExcessiveInput
			***REMOVED***
			return n, nil
		case decoderResultError:
			return n, decodeError(decoderGetErrorCode(r))
		case decoderResultNeedsMoreOutput:
			if n == 0 ***REMOVED***
				return 0, io.ErrShortBuffer
			***REMOVED***
			return n, nil
		case decoderNeedsMoreInput:
		***REMOVED***

		if len(r.in) != 0 ***REMOVED***
			return 0, errInvalidState
		***REMOVED***

		// Calling r.src.Read may block. Don't block if we have data to return.
		if n > 0 ***REMOVED***
			return n, nil
		***REMOVED***

		// Top off the buffer.
		encN, err := r.src.Read(r.buf)
		if encN == 0 ***REMOVED***
			// Not enough data to complete decoding.
			if err == io.EOF ***REMOVED***
				return 0, io.ErrUnexpectedEOF
			***REMOVED***
			return 0, err
		***REMOVED***
		r.in = r.buf[:encN]
	***REMOVED***
***REMOVED***
