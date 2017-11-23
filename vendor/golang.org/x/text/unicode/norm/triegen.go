// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// Trie table generator.
// Used by make*tables tools to generate a go file with trie data structures
// for mapping UTF-8 to a 16-bit value. All but the last byte in a UTF-8 byte
// sequence are used to lookup offsets in the index table to be used for the
// next byte. The last byte is used to index into a table with 16-bit values.

package main

import (
	"fmt"
	"io"
)

const maxSparseEntries = 16

type normCompacter struct ***REMOVED***
	sparseBlocks [][]uint64
	sparseOffset []uint16
	sparseCount  int
	name         string
***REMOVED***

func mostFrequentStride(a []uint64) int ***REMOVED***
	counts := make(map[int]int)
	var v int
	for _, x := range a ***REMOVED***
		if stride := int(x) - v; v != 0 && stride >= 0 ***REMOVED***
			counts[stride]++
		***REMOVED***
		v = int(x)
	***REMOVED***
	var maxs, maxc int
	for stride, cnt := range counts ***REMOVED***
		if cnt > maxc || (cnt == maxc && stride < maxs) ***REMOVED***
			maxs, maxc = stride, cnt
		***REMOVED***
	***REMOVED***
	return maxs
***REMOVED***

func countSparseEntries(a []uint64) int ***REMOVED***
	stride := mostFrequentStride(a)
	var v, count int
	for _, tv := range a ***REMOVED***
		if int(tv)-v != stride ***REMOVED***
			if tv != 0 ***REMOVED***
				count++
			***REMOVED***
		***REMOVED***
		v = int(tv)
	***REMOVED***
	return count
***REMOVED***

func (c *normCompacter) Size(v []uint64) (sz int, ok bool) ***REMOVED***
	if n := countSparseEntries(v); n <= maxSparseEntries ***REMOVED***
		return (n+1)*4 + 2, true
	***REMOVED***
	return 0, false
***REMOVED***

func (c *normCompacter) Store(v []uint64) uint32 ***REMOVED***
	h := uint32(len(c.sparseOffset))
	c.sparseBlocks = append(c.sparseBlocks, v)
	c.sparseOffset = append(c.sparseOffset, uint16(c.sparseCount))
	c.sparseCount += countSparseEntries(v) + 1
	return h
***REMOVED***

func (c *normCompacter) Handler() string ***REMOVED***
	return c.name + "Sparse.lookup"
***REMOVED***

func (c *normCompacter) Print(w io.Writer) (retErr error) ***REMOVED***
	p := func(f string, x ...interface***REMOVED******REMOVED***) ***REMOVED***
		if _, err := fmt.Fprintf(w, f, x...); retErr == nil && err != nil ***REMOVED***
			retErr = err
		***REMOVED***
	***REMOVED***

	ls := len(c.sparseBlocks)
	p("// %sSparseOffset: %d entries, %d bytes\n", c.name, ls, ls*2)
	p("var %sSparseOffset = %#v\n\n", c.name, c.sparseOffset)

	ns := c.sparseCount
	p("// %sSparseValues: %d entries, %d bytes\n", c.name, ns, ns*4)
	p("var %sSparseValues = [%d]valueRange ***REMOVED***", c.name, ns)
	for i, b := range c.sparseBlocks ***REMOVED***
		p("\n// Block %#x, offset %#x", i, c.sparseOffset[i])
		var v int
		stride := mostFrequentStride(b)
		n := countSparseEntries(b)
		p("\n***REMOVED***value:%#04x,lo:%#02x***REMOVED***,", stride, uint8(n))
		for i, nv := range b ***REMOVED***
			if int(nv)-v != stride ***REMOVED***
				if v != 0 ***REMOVED***
					p(",hi:%#02x***REMOVED***,", 0x80+i-1)
				***REMOVED***
				if nv != 0 ***REMOVED***
					p("\n***REMOVED***value:%#04x,lo:%#02x", nv, 0x80+i)
				***REMOVED***
			***REMOVED***
			v = int(nv)
		***REMOVED***
		if v != 0 ***REMOVED***
			p(",hi:%#02x***REMOVED***,", 0x80+len(b)-1)
		***REMOVED***
	***REMOVED***
	p("\n***REMOVED***\n\n")
	return
***REMOVED***
