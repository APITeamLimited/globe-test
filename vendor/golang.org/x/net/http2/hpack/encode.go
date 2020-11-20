// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpack

import (
	"io"
)

const (
	uint32Max              = ^uint32(0)
	initialHeaderTableSize = 4096
)

type Encoder struct ***REMOVED***
	dynTab dynamicTable
	// minSize is the minimum table size set by
	// SetMaxDynamicTableSize after the previous Header Table Size
	// Update.
	minSize uint32
	// maxSizeLimit is the maximum table size this encoder
	// supports. This will protect the encoder from too large
	// size.
	maxSizeLimit uint32
	// tableSizeUpdate indicates whether "Header Table Size
	// Update" is required.
	tableSizeUpdate bool
	w               io.Writer
	buf             []byte
***REMOVED***

// NewEncoder returns a new Encoder which performs HPACK encoding. An
// encoded data is written to w.
func NewEncoder(w io.Writer) *Encoder ***REMOVED***
	e := &Encoder***REMOVED***
		minSize:         uint32Max,
		maxSizeLimit:    initialHeaderTableSize,
		tableSizeUpdate: false,
		w:               w,
	***REMOVED***
	e.dynTab.table.init()
	e.dynTab.setMaxSize(initialHeaderTableSize)
	return e
***REMOVED***

// WriteField encodes f into a single Write to e's underlying Writer.
// This function may also produce bytes for "Header Table Size Update"
// if necessary. If produced, it is done before encoding f.
func (e *Encoder) WriteField(f HeaderField) error ***REMOVED***
	e.buf = e.buf[:0]

	if e.tableSizeUpdate ***REMOVED***
		e.tableSizeUpdate = false
		if e.minSize < e.dynTab.maxSize ***REMOVED***
			e.buf = appendTableSize(e.buf, e.minSize)
		***REMOVED***
		e.minSize = uint32Max
		e.buf = appendTableSize(e.buf, e.dynTab.maxSize)
	***REMOVED***

	idx, nameValueMatch := e.searchTable(f)
	if nameValueMatch ***REMOVED***
		e.buf = appendIndexed(e.buf, idx)
	***REMOVED*** else ***REMOVED***
		indexing := e.shouldIndex(f)
		if indexing ***REMOVED***
			e.dynTab.add(f)
		***REMOVED***

		if idx == 0 ***REMOVED***
			e.buf = appendNewName(e.buf, f, indexing)
		***REMOVED*** else ***REMOVED***
			e.buf = appendIndexedName(e.buf, f, idx, indexing)
		***REMOVED***
	***REMOVED***
	n, err := e.w.Write(e.buf)
	if err == nil && n != len(e.buf) ***REMOVED***
		err = io.ErrShortWrite
	***REMOVED***
	return err
***REMOVED***

// searchTable searches f in both stable and dynamic header tables.
// The static header table is searched first. Only when there is no
// exact match for both name and value, the dynamic header table is
// then searched. If there is no match, i is 0. If both name and value
// match, i is the matched index and nameValueMatch becomes true. If
// only name matches, i points to that index and nameValueMatch
// becomes false.
func (e *Encoder) searchTable(f HeaderField) (i uint64, nameValueMatch bool) ***REMOVED***
	i, nameValueMatch = staticTable.search(f)
	if nameValueMatch ***REMOVED***
		return i, true
	***REMOVED***

	j, nameValueMatch := e.dynTab.table.search(f)
	if nameValueMatch || (i == 0 && j != 0) ***REMOVED***
		return j + uint64(staticTable.len()), nameValueMatch
	***REMOVED***

	return i, false
***REMOVED***

// SetMaxDynamicTableSize changes the dynamic header table size to v.
// The actual size is bounded by the value passed to
// SetMaxDynamicTableSizeLimit.
func (e *Encoder) SetMaxDynamicTableSize(v uint32) ***REMOVED***
	if v > e.maxSizeLimit ***REMOVED***
		v = e.maxSizeLimit
	***REMOVED***
	if v < e.minSize ***REMOVED***
		e.minSize = v
	***REMOVED***
	e.tableSizeUpdate = true
	e.dynTab.setMaxSize(v)
***REMOVED***

// SetMaxDynamicTableSizeLimit changes the maximum value that can be
// specified in SetMaxDynamicTableSize to v. By default, it is set to
// 4096, which is the same size of the default dynamic header table
// size described in HPACK specification. If the current maximum
// dynamic header table size is strictly greater than v, "Header Table
// Size Update" will be done in the next WriteField call and the
// maximum dynamic header table size is truncated to v.
func (e *Encoder) SetMaxDynamicTableSizeLimit(v uint32) ***REMOVED***
	e.maxSizeLimit = v
	if e.dynTab.maxSize > v ***REMOVED***
		e.tableSizeUpdate = true
		e.dynTab.setMaxSize(v)
	***REMOVED***
***REMOVED***

// shouldIndex reports whether f should be indexed.
func (e *Encoder) shouldIndex(f HeaderField) bool ***REMOVED***
	return !f.Sensitive && f.Size() <= e.dynTab.maxSize
***REMOVED***

// appendIndexed appends index i, as encoded in "Indexed Header Field"
// representation, to dst and returns the extended buffer.
func appendIndexed(dst []byte, i uint64) []byte ***REMOVED***
	first := len(dst)
	dst = appendVarInt(dst, 7, i)
	dst[first] |= 0x80
	return dst
***REMOVED***

// appendNewName appends f, as encoded in one of "Literal Header field
// - New Name" representation variants, to dst and returns the
// extended buffer.
//
// If f.Sensitive is true, "Never Indexed" representation is used. If
// f.Sensitive is false and indexing is true, "Incremental Indexing"
// representation is used.
func appendNewName(dst []byte, f HeaderField, indexing bool) []byte ***REMOVED***
	dst = append(dst, encodeTypeByte(indexing, f.Sensitive))
	dst = appendHpackString(dst, f.Name)
	return appendHpackString(dst, f.Value)
***REMOVED***

// appendIndexedName appends f and index i referring indexed name
// entry, as encoded in one of "Literal Header field - Indexed Name"
// representation variants, to dst and returns the extended buffer.
//
// If f.Sensitive is true, "Never Indexed" representation is used. If
// f.Sensitive is false and indexing is true, "Incremental Indexing"
// representation is used.
func appendIndexedName(dst []byte, f HeaderField, i uint64, indexing bool) []byte ***REMOVED***
	first := len(dst)
	var n byte
	if indexing ***REMOVED***
		n = 6
	***REMOVED*** else ***REMOVED***
		n = 4
	***REMOVED***
	dst = appendVarInt(dst, n, i)
	dst[first] |= encodeTypeByte(indexing, f.Sensitive)
	return appendHpackString(dst, f.Value)
***REMOVED***

// appendTableSize appends v, as encoded in "Header Table Size Update"
// representation, to dst and returns the extended buffer.
func appendTableSize(dst []byte, v uint32) []byte ***REMOVED***
	first := len(dst)
	dst = appendVarInt(dst, 5, uint64(v))
	dst[first] |= 0x20
	return dst
***REMOVED***

// appendVarInt appends i, as encoded in variable integer form using n
// bit prefix, to dst and returns the extended buffer.
//
// See
// http://http2.github.io/http2-spec/compression.html#integer.representation
func appendVarInt(dst []byte, n byte, i uint64) []byte ***REMOVED***
	k := uint64((1 << n) - 1)
	if i < k ***REMOVED***
		return append(dst, byte(i))
	***REMOVED***
	dst = append(dst, byte(k))
	i -= k
	for ; i >= 128; i >>= 7 ***REMOVED***
		dst = append(dst, byte(0x80|(i&0x7f)))
	***REMOVED***
	return append(dst, byte(i))
***REMOVED***

// appendHpackString appends s, as encoded in "String Literal"
// representation, to dst and returns the extended buffer.
//
// s will be encoded in Huffman codes only when it produces strictly
// shorter byte string.
func appendHpackString(dst []byte, s string) []byte ***REMOVED***
	huffmanLength := HuffmanEncodeLength(s)
	if huffmanLength < uint64(len(s)) ***REMOVED***
		first := len(dst)
		dst = appendVarInt(dst, 7, huffmanLength)
		dst = AppendHuffmanString(dst, s)
		dst[first] |= 0x80
	***REMOVED*** else ***REMOVED***
		dst = appendVarInt(dst, 7, uint64(len(s)))
		dst = append(dst, s...)
	***REMOVED***
	return dst
***REMOVED***

// encodeTypeByte returns type byte. If sensitive is true, type byte
// for "Never Indexed" representation is returned. If sensitive is
// false and indexing is true, type byte for "Incremental Indexing"
// representation is returned. Otherwise, type byte for "Without
// Indexing" is returned.
func encodeTypeByte(indexing, sensitive bool) byte ***REMOVED***
	if sensitive ***REMOVED***
		return 0x10
	***REMOVED***
	if indexing ***REMOVED***
		return 0x40
	***REMOVED***
	return 0
***REMOVED***
