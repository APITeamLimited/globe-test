// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package colltab

// An Iter incrementally converts chunks of the input text to collation
// elements, while ensuring that the collation elements are in normalized order
// (that is, they are in the order as if the input text were normalized first).
type Iter struct ***REMOVED***
	Weighter Weighter
	Elems    []Elem
	// N is the number of elements in Elems that will not be reordered on
	// subsequent iterations, N <= len(Elems).
	N int

	bytes []byte
	str   string
	// Because the Elems buffer may contain collation elements that are needed
	// for look-ahead, we need two positions in the text (bytes or str): one for
	// the end position in the text for the current iteration and one for the
	// start of the next call to appendNext.
	pEnd  int // end position in text corresponding to N.
	pNext int // pEnd <= pNext.
***REMOVED***

// Reset sets the position in the current input text to p and discards any
// results obtained so far.
func (i *Iter) Reset(p int) ***REMOVED***
	i.Elems = i.Elems[:0]
	i.N = 0
	i.pEnd = p
	i.pNext = p
***REMOVED***

// Len returns the length of the input text.
func (i *Iter) Len() int ***REMOVED***
	if i.bytes != nil ***REMOVED***
		return len(i.bytes)
	***REMOVED***
	return len(i.str)
***REMOVED***

// Discard removes the collation elements up to N.
func (i *Iter) Discard() ***REMOVED***
	// TODO: change this such that only modifiers following starters will have
	// to be copied.
	i.Elems = i.Elems[:copy(i.Elems, i.Elems[i.N:])]
	i.N = 0
***REMOVED***

// End returns the end position of the input text for which Next has returned
// results.
func (i *Iter) End() int ***REMOVED***
	return i.pEnd
***REMOVED***

// SetInput resets i to input s.
func (i *Iter) SetInput(s []byte) ***REMOVED***
	i.bytes = s
	i.str = ""
	i.Reset(0)
***REMOVED***

// SetInputString resets i to input s.
func (i *Iter) SetInputString(s string) ***REMOVED***
	i.str = s
	i.bytes = nil
	i.Reset(0)
***REMOVED***

func (i *Iter) done() bool ***REMOVED***
	return i.pNext >= len(i.str) && i.pNext >= len(i.bytes)
***REMOVED***

func (i *Iter) appendNext() bool ***REMOVED***
	if i.done() ***REMOVED***
		return false
	***REMOVED***
	var sz int
	if i.bytes == nil ***REMOVED***
		i.Elems, sz = i.Weighter.AppendNextString(i.Elems, i.str[i.pNext:])
	***REMOVED*** else ***REMOVED***
		i.Elems, sz = i.Weighter.AppendNext(i.Elems, i.bytes[i.pNext:])
	***REMOVED***
	if sz == 0 ***REMOVED***
		sz = 1
	***REMOVED***
	i.pNext += sz
	return true
***REMOVED***

// Next appends Elems to the internal array. On each iteration, it will either
// add starters or modifiers. In the majority of cases, an Elem with a primary
// value > 0 will have a CCC of 0. The CCC values of collation elements are also
// used to detect if the input string was not normalized and to adjust the
// result accordingly.
func (i *Iter) Next() bool ***REMOVED***
	if i.N == len(i.Elems) && !i.appendNext() ***REMOVED***
		return false
	***REMOVED***

	// Check if the current segment starts with a starter.
	prevCCC := i.Elems[len(i.Elems)-1].CCC()
	if prevCCC == 0 ***REMOVED***
		i.N = len(i.Elems)
		i.pEnd = i.pNext
		return true
	***REMOVED*** else if i.Elems[i.N].CCC() == 0 ***REMOVED***
		// set i.N to only cover part of i.Elems for which prevCCC == 0 and
		// use rest for the next call to next.
		for i.N++; i.N < len(i.Elems) && i.Elems[i.N].CCC() == 0; i.N++ ***REMOVED***
		***REMOVED***
		i.pEnd = i.pNext
		return true
	***REMOVED***

	// The current (partial) segment starts with modifiers. We need to collect
	// all successive modifiers to ensure that they are normalized.
	for ***REMOVED***
		p := len(i.Elems)
		i.pEnd = i.pNext
		if !i.appendNext() ***REMOVED***
			break
		***REMOVED***

		if ccc := i.Elems[p].CCC(); ccc == 0 || len(i.Elems)-i.N > maxCombiningCharacters ***REMOVED***
			// Leave the starter for the next iteration. This ensures that we
			// do not return sequences of collation elements that cross two
			// segments.
			//
			// TODO: handle large number of combining characters by fully
			// normalizing the input segment before iteration. This ensures
			// results are consistent across the text repo.
			i.N = p
			return true
		***REMOVED*** else if ccc < prevCCC ***REMOVED***
			i.doNorm(p, ccc) // should be rare, never occurs for NFD and FCC.
		***REMOVED*** else ***REMOVED***
			prevCCC = ccc
		***REMOVED***
	***REMOVED***

	done := len(i.Elems) != i.N
	i.N = len(i.Elems)
	return done
***REMOVED***

// nextNoNorm is the same as next, but does not "normalize" the collation
// elements.
func (i *Iter) nextNoNorm() bool ***REMOVED***
	// TODO: remove this function. Using this instead of next does not seem
	// to improve performance in any significant way. We retain this until
	// later for evaluation purposes.
	if i.done() ***REMOVED***
		return false
	***REMOVED***
	i.appendNext()
	i.N = len(i.Elems)
	return true
***REMOVED***

const maxCombiningCharacters = 30

// doNorm reorders the collation elements in i.Elems.
// It assumes that blocks of collation elements added with appendNext
// either start and end with the same CCC or start with CCC == 0.
// This allows for a single insertion point for the entire block.
// The correctness of this assumption is verified in builder.go.
func (i *Iter) doNorm(p int, ccc uint8) ***REMOVED***
	n := len(i.Elems)
	k := p
	for p--; p > i.N && ccc < i.Elems[p-1].CCC(); p-- ***REMOVED***
	***REMOVED***
	i.Elems = append(i.Elems, i.Elems[p:k]...)
	copy(i.Elems[p:], i.Elems[k:])
	i.Elems = i.Elems[:n]
***REMOVED***
