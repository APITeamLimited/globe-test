// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The trie in this file is used to associate the first full character
// in a UTF-8 string to a collation element.
// All but the last byte in a UTF-8 byte sequence are
// used to look up offsets in the index table to be used for the next byte.
// The last byte is used to index into a table of collation elements.
// This file contains the code for the generation of the trie.

package build

import (
	"fmt"
	"hash/fnv"
	"io"
	"reflect"
)

const (
	blockSize   = 64
	blockOffset = 2 // Subtract 2 blocks to compensate for the 0x80 added to continuation bytes.
)

type trieHandle struct ***REMOVED***
	lookupStart uint16 // offset in table for first byte
	valueStart  uint16 // offset in table for first byte
***REMOVED***

type trie struct ***REMOVED***
	index  []uint16
	values []uint32
***REMOVED***

// trieNode is the intermediate trie structure used for generating a trie.
type trieNode struct ***REMOVED***
	index    []*trieNode
	value    []uint32
	b        byte
	refValue uint16
	refIndex uint16
***REMOVED***

func newNode() *trieNode ***REMOVED***
	return &trieNode***REMOVED***
		index: make([]*trieNode, 64),
		value: make([]uint32, 128), // root node size is 128 instead of 64
	***REMOVED***
***REMOVED***

func (n *trieNode) isInternal() bool ***REMOVED***
	return n.value != nil
***REMOVED***

func (n *trieNode) insert(r rune, value uint32) ***REMOVED***
	const maskx = 0x3F // mask out two most-significant bits
	str := string(r)
	if len(str) == 1 ***REMOVED***
		n.value[str[0]] = value
		return
	***REMOVED***
	for i := 0; i < len(str)-1; i++ ***REMOVED***
		b := str[i] & maskx
		if n.index == nil ***REMOVED***
			n.index = make([]*trieNode, blockSize)
		***REMOVED***
		nn := n.index[b]
		if nn == nil ***REMOVED***
			nn = &trieNode***REMOVED******REMOVED***
			nn.b = b
			n.index[b] = nn
		***REMOVED***
		n = nn
	***REMOVED***
	if n.value == nil ***REMOVED***
		n.value = make([]uint32, blockSize)
	***REMOVED***
	b := str[len(str)-1] & maskx
	n.value[b] = value
***REMOVED***

type trieBuilder struct ***REMOVED***
	t *trie

	roots []*trieHandle

	lookupBlocks []*trieNode
	valueBlocks  []*trieNode

	lookupBlockIdx map[uint32]*trieNode
	valueBlockIdx  map[uint32]*trieNode
***REMOVED***

func newTrieBuilder() *trieBuilder ***REMOVED***
	index := &trieBuilder***REMOVED******REMOVED***
	index.lookupBlocks = make([]*trieNode, 0)
	index.valueBlocks = make([]*trieNode, 0)
	index.lookupBlockIdx = make(map[uint32]*trieNode)
	index.valueBlockIdx = make(map[uint32]*trieNode)
	// The third nil is the default null block.  The other two blocks
	// are used to guarantee an offset of at least 3 for each block.
	index.lookupBlocks = append(index.lookupBlocks, nil, nil, nil)
	index.t = &trie***REMOVED******REMOVED***
	return index
***REMOVED***

func (b *trieBuilder) computeOffsets(n *trieNode) *trieNode ***REMOVED***
	hasher := fnv.New32()
	if n.index != nil ***REMOVED***
		for i, nn := range n.index ***REMOVED***
			var vi, vv uint16
			if nn != nil ***REMOVED***
				nn = b.computeOffsets(nn)
				n.index[i] = nn
				vi = nn.refIndex
				vv = nn.refValue
			***REMOVED***
			hasher.Write([]byte***REMOVED***byte(vi >> 8), byte(vi)***REMOVED***)
			hasher.Write([]byte***REMOVED***byte(vv >> 8), byte(vv)***REMOVED***)
		***REMOVED***
		h := hasher.Sum32()
		nn, ok := b.lookupBlockIdx[h]
		if !ok ***REMOVED***
			n.refIndex = uint16(len(b.lookupBlocks)) - blockOffset
			b.lookupBlocks = append(b.lookupBlocks, n)
			b.lookupBlockIdx[h] = n
		***REMOVED*** else ***REMOVED***
			n = nn
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, v := range n.value ***REMOVED***
			hasher.Write([]byte***REMOVED***byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)***REMOVED***)
		***REMOVED***
		h := hasher.Sum32()
		nn, ok := b.valueBlockIdx[h]
		if !ok ***REMOVED***
			n.refValue = uint16(len(b.valueBlocks)) - blockOffset
			n.refIndex = n.refValue
			b.valueBlocks = append(b.valueBlocks, n)
			b.valueBlockIdx[h] = n
		***REMOVED*** else ***REMOVED***
			n = nn
		***REMOVED***
	***REMOVED***
	return n
***REMOVED***

func (b *trieBuilder) addStartValueBlock(n *trieNode) uint16 ***REMOVED***
	hasher := fnv.New32()
	for _, v := range n.value[:2*blockSize] ***REMOVED***
		hasher.Write([]byte***REMOVED***byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)***REMOVED***)
	***REMOVED***
	h := hasher.Sum32()
	nn, ok := b.valueBlockIdx[h]
	if !ok ***REMOVED***
		n.refValue = uint16(len(b.valueBlocks))
		n.refIndex = n.refValue
		b.valueBlocks = append(b.valueBlocks, n)
		// Add a dummy block to accommodate the double block size.
		b.valueBlocks = append(b.valueBlocks, nil)
		b.valueBlockIdx[h] = n
	***REMOVED*** else ***REMOVED***
		n = nn
	***REMOVED***
	return n.refValue
***REMOVED***

func genValueBlock(t *trie, n *trieNode) ***REMOVED***
	if n != nil ***REMOVED***
		for _, v := range n.value ***REMOVED***
			t.values = append(t.values, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

func genLookupBlock(t *trie, n *trieNode) ***REMOVED***
	for _, nn := range n.index ***REMOVED***
		v := uint16(0)
		if nn != nil ***REMOVED***
			if n.index != nil ***REMOVED***
				v = nn.refIndex
			***REMOVED*** else ***REMOVED***
				v = nn.refValue
			***REMOVED***
		***REMOVED***
		t.index = append(t.index, v)
	***REMOVED***
***REMOVED***

func (b *trieBuilder) addTrie(n *trieNode) *trieHandle ***REMOVED***
	h := &trieHandle***REMOVED******REMOVED***
	b.roots = append(b.roots, h)
	h.valueStart = b.addStartValueBlock(n)
	if len(b.roots) == 1 ***REMOVED***
		// We insert a null block after the first start value block.
		// This ensures that continuation bytes UTF-8 sequences of length
		// greater than 2 will automatically hit a null block if there
		// was an undefined entry.
		b.valueBlocks = append(b.valueBlocks, nil)
	***REMOVED***
	n = b.computeOffsets(n)
	// Offset by one extra block as the first byte starts at 0xC0 instead of 0x80.
	h.lookupStart = n.refIndex - 1
	return h
***REMOVED***

// generate generates and returns the trie for n.
func (b *trieBuilder) generate() (t *trie, err error) ***REMOVED***
	t = b.t
	if len(b.valueBlocks) >= 1<<16 ***REMOVED***
		return nil, fmt.Errorf("maximum number of value blocks exceeded (%d > %d)", len(b.valueBlocks), 1<<16)
	***REMOVED***
	if len(b.lookupBlocks) >= 1<<16 ***REMOVED***
		return nil, fmt.Errorf("maximum number of lookup blocks exceeded (%d > %d)", len(b.lookupBlocks), 1<<16)
	***REMOVED***
	genValueBlock(t, b.valueBlocks[0])
	genValueBlock(t, &trieNode***REMOVED***value: make([]uint32, 64)***REMOVED***)
	for i := 2; i < len(b.valueBlocks); i++ ***REMOVED***
		genValueBlock(t, b.valueBlocks[i])
	***REMOVED***
	n := &trieNode***REMOVED***index: make([]*trieNode, 64)***REMOVED***
	genLookupBlock(t, n)
	genLookupBlock(t, n)
	genLookupBlock(t, n)
	for i := 3; i < len(b.lookupBlocks); i++ ***REMOVED***
		genLookupBlock(t, b.lookupBlocks[i])
	***REMOVED***
	return b.t, nil
***REMOVED***

func (t *trie) printArrays(w io.Writer, name string) (n, size int, err error) ***REMOVED***
	p := func(f string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	nv := len(t.values)
	p("// %sValues: %d entries, %d bytes\n", name, nv, nv*4)
	p("// Block 2 is the null block.\n")
	p("var %sValues = [%d]uint32 ***REMOVED***", name, nv)
	var printnewline bool
	for i, v := range t.values ***REMOVED***
		if i%blockSize == 0 ***REMOVED***
			p("\n\t// Block %#x, offset %#x", i/blockSize, i)
		***REMOVED***
		if i%4 == 0 ***REMOVED***
			printnewline = true
		***REMOVED***
		if v != 0 ***REMOVED***
			if printnewline ***REMOVED***
				p("\n\t")
				printnewline = false
			***REMOVED***
			p("%#04x:%#08x, ", i, v)
		***REMOVED***
	***REMOVED***
	p("\n***REMOVED***\n\n")
	ni := len(t.index)
	p("// %sLookup: %d entries, %d bytes\n", name, ni, ni*2)
	p("// Block 0 is the null block.\n")
	p("var %sLookup = [%d]uint16 ***REMOVED***", name, ni)
	printnewline = false
	for i, v := range t.index ***REMOVED***
		if i%blockSize == 0 ***REMOVED***
			p("\n\t// Block %#x, offset %#x", i/blockSize, i)
		***REMOVED***
		if i%8 == 0 ***REMOVED***
			printnewline = true
		***REMOVED***
		if v != 0 ***REMOVED***
			if printnewline ***REMOVED***
				p("\n\t")
				printnewline = false
			***REMOVED***
			p("%#03x:%#02x, ", i, v)
		***REMOVED***
	***REMOVED***
	p("\n***REMOVED***\n\n")
	return n, nv*4 + ni*2, err
***REMOVED***

func (t *trie) printStruct(w io.Writer, handle *trieHandle, name string) (n, sz int, err error) ***REMOVED***
	const msg = "trie***REMOVED*** %sLookup[%d:], %sValues[%d:], %sLookup[:], %sValues[:]***REMOVED***"
	n, err = fmt.Fprintf(w, msg, name, handle.lookupStart*blockSize, name, handle.valueStart*blockSize, name, name)
	sz += int(reflect.TypeOf(trie***REMOVED******REMOVED***).Size())
	return
***REMOVED***
