// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package triegen

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"text/template"
)

// print writes all the data structures as well as the code necessary to use the
// trie to w.
func (b *builder) print(w io.Writer) error ***REMOVED***
	b.Stats.NValueEntries = len(b.ValueBlocks) * blockSize
	b.Stats.NValueBytes = len(b.ValueBlocks) * blockSize * b.ValueSize
	b.Stats.NIndexEntries = len(b.IndexBlocks) * blockSize
	b.Stats.NIndexBytes = len(b.IndexBlocks) * blockSize * b.IndexSize
	b.Stats.NHandleBytes = len(b.Trie) * 2 * b.IndexSize

	// If we only have one root trie, all starter blocks are at position 0 and
	// we can access the arrays directly.
	if len(b.Trie) == 1 ***REMOVED***
		// At this point we cannot refer to the generated tables directly.
		b.ASCIIBlock = b.Name + "Values"
		b.StarterBlock = b.Name + "Index"
	***REMOVED*** else ***REMOVED***
		// Otherwise we need to have explicit starter indexes in the trie
		// structure.
		b.ASCIIBlock = "t.ascii"
		b.StarterBlock = "t.utf8Start"
	***REMOVED***

	b.SourceType = "[]byte"
	if err := lookupGen.Execute(w, b); err != nil ***REMOVED***
		return err
	***REMOVED***

	b.SourceType = "string"
	if err := lookupGen.Execute(w, b); err != nil ***REMOVED***
		return err
	***REMOVED***

	if err := trieGen.Execute(w, b); err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, c := range b.Compactions ***REMOVED***
		if err := c.c.Print(w); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func printValues(n int, values []uint64) string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	boff := n * blockSize
	fmt.Fprintf(w, "\t// Block %#x, offset %#x", n, boff)
	var newline bool
	for i, v := range values ***REMOVED***
		if i%6 == 0 ***REMOVED***
			newline = true
		***REMOVED***
		if v != 0 ***REMOVED***
			if newline ***REMOVED***
				fmt.Fprintf(w, "\n")
				newline = false
			***REMOVED***
			fmt.Fprintf(w, "\t%#02x:%#04x, ", boff+i, v)
		***REMOVED***
	***REMOVED***
	return w.String()
***REMOVED***

func printIndex(b *builder, nr int, n *node) string ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	boff := nr * blockSize
	fmt.Fprintf(w, "\t// Block %#x, offset %#x", nr, boff)
	var newline bool
	for i, c := range n.children ***REMOVED***
		if i%8 == 0 ***REMOVED***
			newline = true
		***REMOVED***
		if c != nil ***REMOVED***
			v := b.Compactions[c.index.compaction].Offset + uint32(c.index.index)
			if v != 0 ***REMOVED***
				if newline ***REMOVED***
					fmt.Fprintf(w, "\n")
					newline = false
				***REMOVED***
				fmt.Fprintf(w, "\t%#02x:%#02x, ", boff+i, v)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return w.String()
***REMOVED***

var (
	trieGen = template.Must(template.New("trie").Funcs(template.FuncMap***REMOVED***
		"printValues": printValues,
		"printIndex":  printIndex,
		"title":       strings.Title,
		"dec":         func(x int) int ***REMOVED*** return x - 1 ***REMOVED***,
		"psize": func(n int) string ***REMOVED***
			return fmt.Sprintf("%d bytes (%.2f KiB)", n, float64(n)/1024)
		***REMOVED***,
	***REMOVED***).Parse(trieTemplate))
	lookupGen = template.Must(template.New("lookup").Parse(lookupTemplate))
)

// TODO: consider the return type of lookup. It could be uint64, even if the
// internal value type is smaller. We will have to verify this with the
// performance of unicode/norm, which is very sensitive to such changes.
const trieTemplate = `***REMOVED******REMOVED***$b := .***REMOVED******REMOVED******REMOVED******REMOVED***$multi := gt (len .Trie) 1***REMOVED******REMOVED***
// ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie. Total size: ***REMOVED******REMOVED***psize .Size***REMOVED******REMOVED***. Checksum: ***REMOVED******REMOVED***printf "%08x" .Checksum***REMOVED******REMOVED***.
type ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie struct ***REMOVED*** ***REMOVED******REMOVED***if $multi***REMOVED******REMOVED***
	ascii []***REMOVED******REMOVED***.ValueType***REMOVED******REMOVED*** // index for ASCII bytes
	utf8Start  []***REMOVED******REMOVED***.IndexType***REMOVED******REMOVED*** // index for UTF-8 bytes >= 0xC0
***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED***

func new***REMOVED******REMOVED***title .Name***REMOVED******REMOVED***Trie(i int) ****REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie ***REMOVED*** ***REMOVED******REMOVED***if $multi***REMOVED******REMOVED***
	h := ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***TrieHandles[i]
	return &***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie***REMOVED*** ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Values[uint32(h.ascii)<<6:], ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[uint32(h.multi)<<6:] ***REMOVED***
***REMOVED***

type ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***TrieHandle struct ***REMOVED***
	ascii, multi ***REMOVED******REMOVED***.IndexType***REMOVED******REMOVED***
***REMOVED***

// ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***TrieHandles: ***REMOVED******REMOVED***len .Trie***REMOVED******REMOVED*** handles, ***REMOVED******REMOVED***.Stats.NHandleBytes***REMOVED******REMOVED*** bytes
var ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***TrieHandles = [***REMOVED******REMOVED***len .Trie***REMOVED******REMOVED***]***REMOVED******REMOVED***.Name***REMOVED******REMOVED***TrieHandle***REMOVED***
***REMOVED******REMOVED***range .Trie***REMOVED******REMOVED***	***REMOVED*** ***REMOVED******REMOVED***.ASCIIIndex***REMOVED******REMOVED***, ***REMOVED******REMOVED***.StarterIndex***REMOVED******REMOVED*** ***REMOVED***, // ***REMOVED******REMOVED***printf "%08x" .Checksum***REMOVED******REMOVED***: ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***else***REMOVED******REMOVED***
	return &***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED***
// lookupValue determines the type of block n and looks up the value for b.
func (t ****REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie) lookupValue(n uint32, b byte) ***REMOVED******REMOVED***.ValueType***REMOVED******REMOVED******REMOVED******REMOVED***$last := dec (len .Compactions)***REMOVED******REMOVED*** ***REMOVED***
	switch ***REMOVED*** ***REMOVED******REMOVED***range $i, $c := .Compactions***REMOVED******REMOVED***
		***REMOVED******REMOVED***if eq $i $last***REMOVED******REMOVED***default***REMOVED******REMOVED***else***REMOVED******REMOVED***case n < ***REMOVED******REMOVED***$c.Cutoff***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***:***REMOVED******REMOVED***if ne $i 0***REMOVED******REMOVED***
			n -= ***REMOVED******REMOVED***$c.Offset***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
			return ***REMOVED******REMOVED***print $b.ValueType***REMOVED******REMOVED***(***REMOVED******REMOVED***$c.Handler***REMOVED******REMOVED***)***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***

// ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Values: ***REMOVED******REMOVED***len .ValueBlocks***REMOVED******REMOVED*** blocks, ***REMOVED******REMOVED***.Stats.NValueEntries***REMOVED******REMOVED*** entries, ***REMOVED******REMOVED***.Stats.NValueBytes***REMOVED******REMOVED*** bytes
// The third block is the zero block.
var ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Values = [***REMOVED******REMOVED***.Stats.NValueEntries***REMOVED******REMOVED***]***REMOVED******REMOVED***.ValueType***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***range $i, $v := .ValueBlocks***REMOVED******REMOVED******REMOVED******REMOVED***printValues $i $v***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED***

// ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index: ***REMOVED******REMOVED***len .IndexBlocks***REMOVED******REMOVED*** blocks, ***REMOVED******REMOVED***.Stats.NIndexEntries***REMOVED******REMOVED*** entries, ***REMOVED******REMOVED***.Stats.NIndexBytes***REMOVED******REMOVED*** bytes
// Block 0 is the zero block.
var ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index = [***REMOVED******REMOVED***.Stats.NIndexEntries***REMOVED******REMOVED***]***REMOVED******REMOVED***.IndexType***REMOVED******REMOVED*** ***REMOVED***
***REMOVED******REMOVED***range $i, $v := .IndexBlocks***REMOVED******REMOVED******REMOVED******REMOVED***printIndex $b $i $v***REMOVED******REMOVED***
***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED***
`

// TODO: consider allowing zero-length strings after evaluating performance with
// unicode/norm.
const lookupTemplate = `
// lookup***REMOVED******REMOVED***if eq .SourceType "string"***REMOVED******REMOVED***String***REMOVED******REMOVED***end***REMOVED******REMOVED*** returns the trie value for the first UTF-8 encoding in s and
// the width in bytes of this encoding. The size will be 0 if s does not
// hold enough bytes to complete the encoding. len(s) must be greater than 0.
func (t ****REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie) lookup***REMOVED******REMOVED***if eq .SourceType "string"***REMOVED******REMOVED***String***REMOVED******REMOVED***end***REMOVED******REMOVED***(s ***REMOVED******REMOVED***.SourceType***REMOVED******REMOVED***) (v ***REMOVED******REMOVED***.ValueType***REMOVED******REMOVED***, sz int) ***REMOVED***
	c0 := s[0]
	switch ***REMOVED***
	case c0 < 0x80: // is ASCII
		return ***REMOVED******REMOVED***.ASCIIBlock***REMOVED******REMOVED***[c0], 1
	case c0 < 0xC2:
		return 0, 1  // Illegal UTF-8: not a starter, not ASCII.
	case c0 < 0xE0: // 2-byte UTF-8
		if len(s) < 2 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := ***REMOVED******REMOVED***.StarterBlock***REMOVED******REMOVED***[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c1), 2
	case c0 < 0xF0: // 3-byte UTF-8
		if len(s) < 3 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := ***REMOVED******REMOVED***.StarterBlock***REMOVED******REMOVED***[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c2), 3
	case c0 < 0xF8: // 4-byte UTF-8
		if len(s) < 4 ***REMOVED***
			return 0, 0
		***REMOVED***
		i := ***REMOVED******REMOVED***.StarterBlock***REMOVED******REMOVED***[c0]
		c1 := s[1]
		if c1 < 0x80 || 0xC0 <= c1 ***REMOVED***
			return 0, 1 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o := uint32(i)<<6 + uint32(c1)
		i = ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[o]
		c2 := s[2]
		if c2 < 0x80 || 0xC0 <= c2 ***REMOVED***
			return 0, 2 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		o = uint32(i)<<6 + uint32(c2)
		i = ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[o]
		c3 := s[3]
		if c3 < 0x80 || 0xC0 <= c3 ***REMOVED***
			return 0, 3 // Illegal UTF-8: not a continuation byte.
		***REMOVED***
		return t.lookupValue(uint32(i), c3), 4
	***REMOVED***
	// Illegal rune
	return 0, 1
***REMOVED***

// lookup***REMOVED******REMOVED***if eq .SourceType "string"***REMOVED******REMOVED***String***REMOVED******REMOVED***end***REMOVED******REMOVED***Unsafe returns the trie value for the first UTF-8 encoding in s.
// s must start with a full and valid UTF-8 encoded rune.
func (t ****REMOVED******REMOVED***.Name***REMOVED******REMOVED***Trie) lookup***REMOVED******REMOVED***if eq .SourceType "string"***REMOVED******REMOVED***String***REMOVED******REMOVED***end***REMOVED******REMOVED***Unsafe(s ***REMOVED******REMOVED***.SourceType***REMOVED******REMOVED***) ***REMOVED******REMOVED***.ValueType***REMOVED******REMOVED*** ***REMOVED***
	c0 := s[0]
	if c0 < 0x80 ***REMOVED*** // is ASCII
		return ***REMOVED******REMOVED***.ASCIIBlock***REMOVED******REMOVED***[c0]
	***REMOVED***
	i := ***REMOVED******REMOVED***.StarterBlock***REMOVED******REMOVED***[c0]
	if c0 < 0xE0 ***REMOVED*** // 2-byte UTF-8
		return t.lookupValue(uint32(i), s[1])
	***REMOVED***
	i = ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[uint32(i)<<6+uint32(s[1])]
	if c0 < 0xF0 ***REMOVED*** // 3-byte UTF-8
		return t.lookupValue(uint32(i), s[2])
	***REMOVED***
	i = ***REMOVED******REMOVED***.Name***REMOVED******REMOVED***Index[uint32(i)<<6+uint32(s[2])]
	if c0 < 0xF8 ***REMOVED*** // 4-byte UTF-8
		return t.lookupValue(uint32(i), s[3])
	***REMOVED***
	return 0
***REMOVED***
`
