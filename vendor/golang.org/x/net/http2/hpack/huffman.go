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
	// This relies on the maximum huffman code length being 30 (See tables.go huffmanCodeLen array)
	// So if a uint64 buffer has less than 32 valid bits can always accommodate another huffmanCode.
	var (
		x uint64 // buffer
		n uint   // number valid of bits present in x
	)
	for i := 0; i < len(s); i++ ***REMOVED***
		c := s[i]
		n += uint(huffmanCodeLen[c])
		x <<= huffmanCodeLen[c] % 64
		x |= uint64(huffmanCodes[c])
		if n >= 32 ***REMOVED***
			n %= 32             // Normally would be -= 32 but %= 32 informs compiler 0 <= n <= 31 for upcoming shift
			y := uint32(x >> n) // Compiler doesn't combine memory writes if y isn't uint32
			dst = append(dst, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
		***REMOVED***
	***REMOVED***
	// Add padding bits if necessary
	if over := n % 8; over > 0 ***REMOVED***
		const (
			eosCode    = 0x3fffffff
			eosNBits   = 30
			eosPadByte = eosCode >> (eosNBits - 8)
		)
		pad := 8 - over
		x = (x << pad) | (eosPadByte >> over)
		n += pad // 8 now divides into n exactly
	***REMOVED***
	// n in (0, 8, 16, 24, 32)
	switch n / 8 ***REMOVED***
	case 0:
		return dst
	case 1:
		return append(dst, byte(x))
	case 2:
		y := uint16(x)
		return append(dst, byte(y>>8), byte(y))
	case 3:
		y := uint16(x >> 8)
		return append(dst, byte(y>>8), byte(y), byte(x))
	***REMOVED***
	//	case 4:
	y := uint32(x)
	return append(dst, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
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
