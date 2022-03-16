// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hpack

import (
	"bytes"
	"errors"
	"io"
	"sync"
)

var bufPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(bytes.Buffer) ***REMOVED***,
***REMOVED***

// HuffmanDecode decodes the string in v and writes the expanded
// result to w, returning the number of bytes written to w and the
// Write call's return value. At most one Write call is made.
func HuffmanDecode(w io.Writer, v []byte) (int, error) ***REMOVED***
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := huffmanDecode(buf, 0, v); err != nil ***REMOVED***
		return 0, err
	***REMOVED***
	return w.Write(buf.Bytes())
***REMOVED***

// HuffmanDecodeToString decodes the string in v.
func HuffmanDecodeToString(v []byte) (string, error) ***REMOVED***
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	if err := huffmanDecode(buf, 0, v); err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return buf.String(), nil
***REMOVED***

// ErrInvalidHuffman is returned for errors found decoding
// Huffman-encoded strings.
var ErrInvalidHuffman = errors.New("hpack: invalid Huffman-encoded data")

// huffmanDecode decodes v to buf.
// If maxLen is greater than 0, attempts to write more to buf than
// maxLen bytes will return ErrStringLength.
func huffmanDecode(buf *bytes.Buffer, maxLen int, v []byte) error ***REMOVED***
	rootHuffmanNode := getRootHuffmanNode()
	n := rootHuffmanNode
	// cur is the bit buffer that has not been fed into n.
	// cbits is the number of low order bits in cur that are valid.
	// sbits is the number of bits of the symbol prefix being decoded.
	cur, cbits, sbits := uint(0), uint8(0), uint8(0)
	for _, b := range v ***REMOVED***
		cur = cur<<8 | uint(b)
		cbits += 8
		sbits += 8
		for cbits >= 8 ***REMOVED***
			idx := byte(cur >> (cbits - 8))
			n = n.children[idx]
			if n == nil ***REMOVED***
				return ErrInvalidHuffman
			***REMOVED***
			if n.children == nil ***REMOVED***
				if maxLen != 0 && buf.Len() == maxLen ***REMOVED***
					return ErrStringLength
				***REMOVED***
				buf.WriteByte(n.sym)
				cbits -= n.codeLen
				n = rootHuffmanNode
				sbits = cbits
			***REMOVED*** else ***REMOVED***
				cbits -= 8
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for cbits > 0 ***REMOVED***
		n = n.children[byte(cur<<(8-cbits))]
		if n == nil ***REMOVED***
			return ErrInvalidHuffman
		***REMOVED***
		if n.children != nil || n.codeLen > cbits ***REMOVED***
			break
		***REMOVED***
		if maxLen != 0 && buf.Len() == maxLen ***REMOVED***
			return ErrStringLength
		***REMOVED***
		buf.WriteByte(n.sym)
		cbits -= n.codeLen
		n = rootHuffmanNode
		sbits = cbits
	***REMOVED***
	if sbits > 7 ***REMOVED***
		// Either there was an incomplete symbol, or overlong padding.
		// Both are decoding errors per RFC 7541 section 5.2.
		return ErrInvalidHuffman
	***REMOVED***
	if mask := uint(1<<cbits - 1); cur&mask != mask ***REMOVED***
		// Trailing bits must be a prefix of EOS per RFC 7541 section 5.2.
		return ErrInvalidHuffman
	***REMOVED***

	return nil
***REMOVED***

// incomparable is a zero-width, non-comparable type. Adding it to a struct
// makes that struct also non-comparable, and generally doesn't add
// any size (as long as it's first).
type incomparable [0]func()

type node struct ***REMOVED***
	_ incomparable

	// children is non-nil for internal nodes
	children *[256]*node

	// The following are only valid if children is nil:
	codeLen uint8 // number of bits that led to the output of sym
	sym     byte  // output symbol
***REMOVED***

func newInternalNode() *node ***REMOVED***
	return &node***REMOVED***children: new([256]*node)***REMOVED***
***REMOVED***

var (
	buildRootOnce       sync.Once
	lazyRootHuffmanNode *node
)

func getRootHuffmanNode() *node ***REMOVED***
	buildRootOnce.Do(buildRootHuffmanNode)
	return lazyRootHuffmanNode
***REMOVED***

func buildRootHuffmanNode() ***REMOVED***
	if len(huffmanCodes) != 256 ***REMOVED***
		panic("unexpected size")
	***REMOVED***
	lazyRootHuffmanNode = newInternalNode()
	// allocate a leaf node for each of the 256 symbols
	leaves := new([256]node)

	for sym, code := range huffmanCodes ***REMOVED***
		codeLen := huffmanCodeLen[sym]

		cur := lazyRootHuffmanNode
		for codeLen > 8 ***REMOVED***
			codeLen -= 8
			i := uint8(code >> codeLen)
			if cur.children[i] == nil ***REMOVED***
				cur.children[i] = newInternalNode()
			***REMOVED***
			cur = cur.children[i]
		***REMOVED***
		shift := 8 - codeLen
		start, end := int(uint8(code<<shift)), int(1<<shift)

		leaves[sym].sym = byte(sym)
		leaves[sym].codeLen = codeLen
		for i := start; i < start+end; i++ ***REMOVED***
			cur.children[i] = &leaves[sym]
		***REMOVED***
	***REMOVED***
***REMOVED***

// AppendHuffmanString appends s, as encoded in Huffman codes, to dst
// and returns the extended buffer.
func AppendHuffmanString(dst []byte, s string) []byte ***REMOVED***
	rembits := uint8(8)

	for i := 0; i < len(s); i++ ***REMOVED***
		if rembits == 8 ***REMOVED***
			dst = append(dst, 0)
		***REMOVED***
		dst, rembits = appendByteToHuffmanCode(dst, rembits, s[i])
	***REMOVED***

	if rembits < 8 ***REMOVED***
		// special EOS symbol
		code := uint32(0x3fffffff)
		nbits := uint8(30)

		t := uint8(code >> (nbits - rembits))
		dst[len(dst)-1] |= t
	***REMOVED***

	return dst
***REMOVED***

// HuffmanEncodeLength returns the number of bytes required to encode
// s in Huffman codes. The result is round up to byte boundary.
func HuffmanEncodeLength(s string) uint64 ***REMOVED***
	n := uint64(0)
	for i := 0; i < len(s); i++ ***REMOVED***
		n += uint64(huffmanCodeLen[s[i]])
	***REMOVED***
	return (n + 7) / 8
***REMOVED***

// appendByteToHuffmanCode appends Huffman code for c to dst and
// returns the extended buffer and the remaining bits in the last
// element. The appending is not byte aligned and the remaining bits
// in the last element of dst is given in rembits.
func appendByteToHuffmanCode(dst []byte, rembits uint8, c byte) ([]byte, uint8) ***REMOVED***
	code := huffmanCodes[c]
	nbits := huffmanCodeLen[c]

	for ***REMOVED***
		if rembits > nbits ***REMOVED***
			t := uint8(code << (rembits - nbits))
			dst[len(dst)-1] |= t
			rembits -= nbits
			break
		***REMOVED***

		t := uint8(code >> (nbits - rembits))
		dst[len(dst)-1] |= t

		nbits -= rembits
		rembits = 8

		if nbits == 0 ***REMOVED***
			break
		***REMOVED***

		dst = append(dst, 0)
	***REMOVED***

	return dst, rembits
***REMOVED***
