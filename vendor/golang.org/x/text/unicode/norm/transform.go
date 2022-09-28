// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import (
	"unicode/utf8"

	"golang.org/x/text/transform"
)

// Reset implements the Reset method of the transform.Transformer interface.
func (Form) Reset() ***REMOVED******REMOVED***

// Transform implements the Transform method of the transform.Transformer
// interface. It may need to write segments of up to MaxSegmentSize at once.
// Users should either catch ErrShortDst and allow dst to grow or have dst be at
// least of size MaxTransformChunkSize to be guaranteed of progress.
func (f Form) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	// Cap the maximum number of src bytes to check.
	b := src
	eof := atEOF
	if ns := len(dst); ns < len(b) ***REMOVED***
		err = transform.ErrShortDst
		eof = false
		b = b[:ns]
	***REMOVED***
	i, ok := formTable[f].quickSpan(inputBytes(b), 0, len(b), eof)
	n := copy(dst, b[:i])
	if !ok ***REMOVED***
		nDst, nSrc, err = f.transform(dst[n:], src[n:], atEOF)
		return nDst + n, nSrc + n, err
	***REMOVED***

	if err == nil && n < len(src) && !atEOF ***REMOVED***
		err = transform.ErrShortSrc
	***REMOVED***
	return n, n, err
***REMOVED***

func flushTransform(rb *reorderBuffer) bool ***REMOVED***
	// Write out (must fully fit in dst, or else it is an ErrShortDst).
	if len(rb.out) < rb.nrune*utf8.UTFMax ***REMOVED***
		return false
	***REMOVED***
	rb.out = rb.out[rb.flushCopy(rb.out):]
	return true
***REMOVED***

var errs = []error***REMOVED***nil, transform.ErrShortDst, transform.ErrShortSrc***REMOVED***

// transform implements the transform.Transformer interface. It is only called
// when quickSpan does not pass for a given string.
func (f Form) transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	// TODO: get rid of reorderBuffer. See CL 23460044.
	rb := reorderBuffer***REMOVED******REMOVED***
	rb.init(f, src)
	for ***REMOVED***
		// Load segment into reorder buffer.
		rb.setFlusher(dst[nDst:], flushTransform)
		end := decomposeSegment(&rb, nSrc, atEOF)
		if end < 0 ***REMOVED***
			return nDst, nSrc, errs[-end]
		***REMOVED***
		nDst = len(dst) - len(rb.out)
		nSrc = end

		// Next quickSpan.
		end = rb.nsrc
		eof := atEOF
		if n := nSrc + len(dst) - nDst; n < end ***REMOVED***
			err = transform.ErrShortDst
			end = n
			eof = false
		***REMOVED***
		end, ok := rb.f.quickSpan(rb.src, nSrc, end, eof)
		n := copy(dst[nDst:], rb.src.bytes[nSrc:end])
		nSrc += n
		nDst += n
		if ok ***REMOVED***
			if err == nil && n < rb.nsrc && !atEOF ***REMOVED***
				err = transform.ErrShortSrc
			***REMOVED***
			return nDst, nSrc, err
		***REMOVED***
	***REMOVED***
***REMOVED***
