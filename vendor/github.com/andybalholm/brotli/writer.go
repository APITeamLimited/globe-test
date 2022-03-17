package brotli

import (
	"errors"
	"io"
)

const (
	BestSpeed          = 0
	BestCompression    = 11
	DefaultCompression = 6
)

// WriterOptions configures Writer.
type WriterOptions struct ***REMOVED***
	// Quality controls the compression-speed vs compression-density trade-offs.
	// The higher the quality, the slower the compression. Range is 0 to 11.
	Quality int
	// LGWin is the base 2 logarithm of the sliding window size.
	// Range is 10 to 24. 0 indicates automatic configuration based on Quality.
	LGWin int
***REMOVED***

var (
	errEncode       = errors.New("brotli: encode error")
	errWriterClosed = errors.New("brotli: Writer is closed")
)

// Writes to the returned writer are compressed and written to dst.
// It is the caller's responsibility to call Close on the Writer when done.
// Writes may be buffered and not flushed until Close.
func NewWriter(dst io.Writer) *Writer ***REMOVED***
	return NewWriterLevel(dst, DefaultCompression)
***REMOVED***

// NewWriterLevel is like NewWriter but specifies the compression level instead
// of assuming DefaultCompression.
// The compression level can be DefaultCompression or any integer value between
// BestSpeed and BestCompression inclusive.
func NewWriterLevel(dst io.Writer, level int) *Writer ***REMOVED***
	return NewWriterOptions(dst, WriterOptions***REMOVED***
		Quality: level,
	***REMOVED***)
***REMOVED***

// NewWriterOptions is like NewWriter but specifies WriterOptions
func NewWriterOptions(dst io.Writer, options WriterOptions) *Writer ***REMOVED***
	w := new(Writer)
	w.options = options
	w.Reset(dst)
	return w
***REMOVED***

// Reset discards the Writer's state and makes it equivalent to the result of
// its original state from NewWriter or NewWriterLevel, but writing to dst
// instead. This permits reusing a Writer rather than allocating a new one.
func (w *Writer) Reset(dst io.Writer) ***REMOVED***
	encoderInitState(w)
	w.params.quality = w.options.Quality
	if w.options.LGWin > 0 ***REMOVED***
		w.params.lgwin = uint(w.options.LGWin)
	***REMOVED***
	w.dst = dst
	w.err = nil
***REMOVED***

func (w *Writer) writeChunk(p []byte, op int) (n int, err error) ***REMOVED***
	if w.dst == nil ***REMOVED***
		return 0, errWriterClosed
	***REMOVED***
	if w.err != nil ***REMOVED***
		return 0, w.err
	***REMOVED***

	for ***REMOVED***
		availableIn := uint(len(p))
		nextIn := p
		success := encoderCompressStream(w, op, &availableIn, &nextIn)
		bytesConsumed := len(p) - int(availableIn)
		p = p[bytesConsumed:]
		n += bytesConsumed
		if !success ***REMOVED***
			return n, errEncode
		***REMOVED***

		if len(p) == 0 || w.err != nil ***REMOVED***
			return n, w.err
		***REMOVED***
	***REMOVED***
***REMOVED***

// Flush outputs encoded data for all input provided to Write. The resulting
// output can be decoded to match all input before Flush, but the stream is
// not yet complete until after Close.
// Flush has a negative impact on compression.
func (w *Writer) Flush() error ***REMOVED***
	_, err := w.writeChunk(nil, operationFlush)
	return err
***REMOVED***

// Close flushes remaining data to the decorated writer.
func (w *Writer) Close() error ***REMOVED***
	// If stream is already closed, it is reported by `writeChunk`.
	_, err := w.writeChunk(nil, operationFinish)
	w.dst = nil
	return err
***REMOVED***

// Write implements io.Writer. Flush or Close must be called to ensure that the
// encoded bytes are actually flushed to the underlying Writer.
func (w *Writer) Write(p []byte) (n int, err error) ***REMOVED***
	return w.writeChunk(p, operationProcess)
***REMOVED***

type nopCloser struct ***REMOVED***
	io.Writer
***REMOVED***

func (nopCloser) Close() error ***REMOVED*** return nil ***REMOVED***
