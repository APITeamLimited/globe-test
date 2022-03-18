// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.

package zstd

import (
	"errors"
	"io"
	"sync"
)

// ZipMethodWinZip is the method for Zstandard compressed data inside Zip files for WinZip.
// See https://www.winzip.com/win/en/comp_info.html
const ZipMethodWinZip = 93

// ZipMethodPKWare is the original method number used by PKWARE to indicate Zstandard compression.
// Deprecated: This has been deprecated by PKWARE, use ZipMethodWinZip instead for compression.
// See https://pkware.cachefly.net/webdocs/APPNOTE/APPNOTE-6.3.9.TXT
const ZipMethodPKWare = 20

var zipReaderPool sync.Pool

// newZipReader creates a pooled zip decompressor.
func newZipReader(r io.Reader) io.ReadCloser ***REMOVED***
	dec, ok := zipReaderPool.Get().(*Decoder)
	if ok ***REMOVED***
		dec.Reset(r)
	***REMOVED*** else ***REMOVED***
		d, err := NewReader(r, WithDecoderConcurrency(1), WithDecoderLowmem(true))
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		dec = d
	***REMOVED***
	return &pooledZipReader***REMOVED***dec: dec***REMOVED***
***REMOVED***

type pooledZipReader struct ***REMOVED***
	mu  sync.Mutex // guards Close and Read
	dec *Decoder
***REMOVED***

func (r *pooledZipReader) Read(p []byte) (n int, err error) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.dec == nil ***REMOVED***
		return 0, errors.New("read after close or EOF")
	***REMOVED***
	dec, err := r.dec.Read(p)
	if err == io.EOF ***REMOVED***
		err = r.dec.Reset(nil)
		zipReaderPool.Put(r.dec)
		r.dec = nil
	***REMOVED***
	return dec, err
***REMOVED***

func (r *pooledZipReader) Close() error ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	var err error
	if r.dec != nil ***REMOVED***
		err = r.dec.Reset(nil)
		zipReaderPool.Put(r.dec)
		r.dec = nil
	***REMOVED***
	return err
***REMOVED***

type pooledZipWriter struct ***REMOVED***
	mu   sync.Mutex // guards Close and Read
	enc  *Encoder
	pool *sync.Pool
***REMOVED***

func (w *pooledZipWriter) Write(p []byte) (n int, err error) ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.enc == nil ***REMOVED***
		return 0, errors.New("Write after Close")
	***REMOVED***
	return w.enc.Write(p)
***REMOVED***

func (w *pooledZipWriter) Close() error ***REMOVED***
	w.mu.Lock()
	defer w.mu.Unlock()
	var err error
	if w.enc != nil ***REMOVED***
		err = w.enc.Close()
		w.pool.Put(w.enc)
		w.enc = nil
	***REMOVED***
	return err
***REMOVED***

// ZipCompressor returns a compressor that can be registered with zip libraries.
// The provided encoder options will be used on all encodes.
func ZipCompressor(opts ...EOption) func(w io.Writer) (io.WriteCloser, error) ***REMOVED***
	var pool sync.Pool
	return func(w io.Writer) (io.WriteCloser, error) ***REMOVED***
		enc, ok := pool.Get().(*Encoder)
		if ok ***REMOVED***
			enc.Reset(w)
		***REMOVED*** else ***REMOVED***
			var err error
			enc, err = NewWriter(w, opts...)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		return &pooledZipWriter***REMOVED***enc: enc, pool: &pool***REMOVED***, nil
	***REMOVED***
***REMOVED***

// ZipDecompressor returns a decompressor that can be registered with zip libraries.
// See ZipCompressor for example.
func ZipDecompressor() func(r io.Reader) io.ReadCloser ***REMOVED***
	return newZipReader
***REMOVED***
