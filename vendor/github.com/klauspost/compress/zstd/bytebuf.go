// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"fmt"
	"io"
	"io/ioutil"
)

type byteBuffer interface ***REMOVED***
	// Read up to 8 bytes.
	// Returns io.ErrUnexpectedEOF if this cannot be satisfied.
	readSmall(n int) ([]byte, error)

	// Read >8 bytes.
	// MAY use the destination slice.
	readBig(n int, dst []byte) ([]byte, error)

	// Read a single byte.
	readByte() (byte, error)

	// Skip n bytes.
	skipN(n int) error
***REMOVED***

// in-memory buffer
type byteBuf []byte

func (b *byteBuf) readSmall(n int) ([]byte, error) ***REMOVED***
	if debugAsserts && n > 8 ***REMOVED***
		panic(fmt.Errorf("small read > 8 (%d). use readBig", n))
	***REMOVED***
	bb := *b
	if len(bb) < n ***REMOVED***
		return nil, io.ErrUnexpectedEOF
	***REMOVED***
	r := bb[:n]
	*b = bb[n:]
	return r, nil
***REMOVED***

func (b *byteBuf) readBig(n int, dst []byte) ([]byte, error) ***REMOVED***
	bb := *b
	if len(bb) < n ***REMOVED***
		return nil, io.ErrUnexpectedEOF
	***REMOVED***
	r := bb[:n]
	*b = bb[n:]
	return r, nil
***REMOVED***

func (b *byteBuf) readByte() (byte, error) ***REMOVED***
	bb := *b
	if len(bb) < 1 ***REMOVED***
		return 0, nil
	***REMOVED***
	r := bb[0]
	*b = bb[1:]
	return r, nil
***REMOVED***

func (b *byteBuf) skipN(n int) error ***REMOVED***
	bb := *b
	if len(bb) < n ***REMOVED***
		return io.ErrUnexpectedEOF
	***REMOVED***
	*b = bb[n:]
	return nil
***REMOVED***

// wrapper around a reader.
type readerWrapper struct ***REMOVED***
	r   io.Reader
	tmp [8]byte
***REMOVED***

func (r *readerWrapper) readSmall(n int) ([]byte, error) ***REMOVED***
	if debugAsserts && n > 8 ***REMOVED***
		panic(fmt.Errorf("small read > 8 (%d). use readBig", n))
	***REMOVED***
	n2, err := io.ReadFull(r.r, r.tmp[:n])
	// We only really care about the actual bytes read.
	if err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			return nil, io.ErrUnexpectedEOF
		***REMOVED***
		if debugDecoder ***REMOVED***
			println("readSmall: got", n2, "want", n, "err", err)
		***REMOVED***
		return nil, err
	***REMOVED***
	return r.tmp[:n], nil
***REMOVED***

func (r *readerWrapper) readBig(n int, dst []byte) ([]byte, error) ***REMOVED***
	if cap(dst) < n ***REMOVED***
		dst = make([]byte, n)
	***REMOVED***
	n2, err := io.ReadFull(r.r, dst[:n])
	if err == io.EOF && n > 0 ***REMOVED***
		err = io.ErrUnexpectedEOF
	***REMOVED***
	return dst[:n2], err
***REMOVED***

func (r *readerWrapper) readByte() (byte, error) ***REMOVED***
	n2, err := r.r.Read(r.tmp[:1])
	if err != nil ***REMOVED***
		if err == io.EOF ***REMOVED***
			err = io.ErrUnexpectedEOF
		***REMOVED***
		return 0, err
	***REMOVED***
	if n2 != 1 ***REMOVED***
		return 0, io.ErrUnexpectedEOF
	***REMOVED***
	return r.tmp[0], nil
***REMOVED***

func (r *readerWrapper) skipN(n int) error ***REMOVED***
	n2, err := io.CopyN(ioutil.Discard, r.r, int64(n))
	if n2 != int64(n) ***REMOVED***
		err = io.ErrUnexpectedEOF
	***REMOVED***
	return err
***REMOVED***
