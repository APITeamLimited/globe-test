// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package collate

import (
	"bytes"
	"sort"
)

const (
	maxSortBuffer  = 40960
	maxSortEntries = 4096
)

type swapper interface ***REMOVED***
	Swap(i, j int)
***REMOVED***

type sorter struct ***REMOVED***
	buf  *Buffer
	keys [][]byte
	src  swapper
***REMOVED***

func (s *sorter) init(n int) ***REMOVED***
	if s.buf == nil ***REMOVED***
		s.buf = &Buffer***REMOVED******REMOVED***
		s.buf.init()
	***REMOVED***
	if cap(s.keys) < n ***REMOVED***
		s.keys = make([][]byte, n)
	***REMOVED***
	s.keys = s.keys[0:n]
***REMOVED***

func (s *sorter) sort(src swapper) ***REMOVED***
	s.src = src
	sort.Sort(s)
***REMOVED***

func (s sorter) Len() int ***REMOVED***
	return len(s.keys)
***REMOVED***

func (s sorter) Less(i, j int) bool ***REMOVED***
	return bytes.Compare(s.keys[i], s.keys[j]) == -1
***REMOVED***

func (s sorter) Swap(i, j int) ***REMOVED***
	s.keys[i], s.keys[j] = s.keys[j], s.keys[i]
	s.src.Swap(i, j)
***REMOVED***

// A Lister can be sorted by Collator's Sort method.
type Lister interface ***REMOVED***
	Len() int
	Swap(i, j int)
	// Bytes returns the bytes of the text at index i.
	Bytes(i int) []byte
***REMOVED***

// Sort uses sort.Sort to sort the strings represented by x using the rules of c.
func (c *Collator) Sort(x Lister) ***REMOVED***
	n := x.Len()
	c.sorter.init(n)
	for i := 0; i < n; i++ ***REMOVED***
		c.sorter.keys[i] = c.Key(c.sorter.buf, x.Bytes(i))
	***REMOVED***
	c.sorter.sort(x)
***REMOVED***

// SortStrings uses sort.Sort to sort the strings in x using the rules of c.
func (c *Collator) SortStrings(x []string) ***REMOVED***
	c.sorter.init(len(x))
	for i, s := range x ***REMOVED***
		c.sorter.keys[i] = c.KeyFromString(c.sorter.buf, s)
	***REMOVED***
	c.sorter.sort(sort.StringSlice(x))
***REMOVED***
