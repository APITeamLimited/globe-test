// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases

import "golang.org/x/text/transform"

// A context is used for iterating over source bytes, fetching case info and
// writing to a destination buffer.
//
// Casing operations may need more than one rune of context to decide how a rune
// should be cased. Casing implementations should call checkpoint on context
// whenever it is known to be safe to return the runes processed so far.
//
// It is recommended for implementations to not allow for more than 30 case
// ignorables as lookahead (analogous to the limit in norm) and to use state if
// unbounded lookahead is needed for cased runes.
type context struct ***REMOVED***
	dst, src []byte
	atEOF    bool

	pDst int // pDst points past the last written rune in dst.
	pSrc int // pSrc points to the start of the currently scanned rune.

	// checkpoints safe to return in Transform, where nDst <= pDst and nSrc <= pSrc.
	nDst, nSrc int
	err        error

	sz   int  // size of current rune
	info info // case information of currently scanned rune

	// State preserved across calls to Transform.
	isMidWord bool // false if next cased letter needs to be title-cased.
***REMOVED***

func (c *context) Reset() ***REMOVED***
	c.isMidWord = false
***REMOVED***

// ret returns the return values for the Transform method. It checks whether
// there were insufficient bytes in src to complete and introduces an error
// accordingly, if necessary.
func (c *context) ret() (nDst, nSrc int, err error) ***REMOVED***
	if c.err != nil || c.nSrc == len(c.src) ***REMOVED***
		return c.nDst, c.nSrc, c.err
	***REMOVED***
	// This point is only reached by mappers if there was no short destination
	// buffer. This means that the source buffer was exhausted and that c.sz was
	// set to 0 by next.
	if c.atEOF && c.pSrc == len(c.src) ***REMOVED***
		return c.pDst, c.pSrc, nil
	***REMOVED***
	return c.nDst, c.nSrc, transform.ErrShortSrc
***REMOVED***

// retSpan returns the return values for the Span method. It checks whether
// there were insufficient bytes in src to complete and introduces an error
// accordingly, if necessary.
func (c *context) retSpan() (n int, err error) ***REMOVED***
	_, nSrc, err := c.ret()
	return nSrc, err
***REMOVED***

// checkpoint sets the return value buffer points for Transform to the current
// positions.
func (c *context) checkpoint() ***REMOVED***
	if c.err == nil ***REMOVED***
		c.nDst, c.nSrc = c.pDst, c.pSrc+c.sz
	***REMOVED***
***REMOVED***

// unreadRune causes the last rune read by next to be reread on the next
// invocation of next. Only one unreadRune may be called after a call to next.
func (c *context) unreadRune() ***REMOVED***
	c.sz = 0
***REMOVED***

func (c *context) next() bool ***REMOVED***
	c.pSrc += c.sz
	if c.pSrc == len(c.src) || c.err != nil ***REMOVED***
		c.info, c.sz = 0, 0
		return false
	***REMOVED***
	v, sz := trie.lookup(c.src[c.pSrc:])
	c.info, c.sz = info(v), sz
	if c.sz == 0 ***REMOVED***
		if c.atEOF ***REMOVED***
			// A zero size means we have an incomplete rune. If we are atEOF,
			// this means it is an illegal rune, which we will consume one
			// byte at a time.
			c.sz = 1
		***REMOVED*** else ***REMOVED***
			c.err = transform.ErrShortSrc
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// writeBytes adds bytes to dst.
func (c *context) writeBytes(b []byte) bool ***REMOVED***
	if len(c.dst)-c.pDst < len(b) ***REMOVED***
		c.err = transform.ErrShortDst
		return false
	***REMOVED***
	// This loop is faster than using copy.
	for _, ch := range b ***REMOVED***
		c.dst[c.pDst] = ch
		c.pDst++
	***REMOVED***
	return true
***REMOVED***

// writeString writes the given string to dst.
func (c *context) writeString(s string) bool ***REMOVED***
	if len(c.dst)-c.pDst < len(s) ***REMOVED***
		c.err = transform.ErrShortDst
		return false
	***REMOVED***
	// This loop is faster than using copy.
	for i := 0; i < len(s); i++ ***REMOVED***
		c.dst[c.pDst] = s[i]
		c.pDst++
	***REMOVED***
	return true
***REMOVED***

// copy writes the current rune to dst.
func (c *context) copy() bool ***REMOVED***
	return c.writeBytes(c.src[c.pSrc : c.pSrc+c.sz])
***REMOVED***

// copyXOR copies the current rune to dst and modifies it by applying the XOR
// pattern of the case info. It is the responsibility of the caller to ensure
// that this is a rune with a XOR pattern defined.
func (c *context) copyXOR() bool ***REMOVED***
	if !c.copy() ***REMOVED***
		return false
	***REMOVED***
	if c.info&xorIndexBit == 0 ***REMOVED***
		// Fast path for 6-bit XOR pattern, which covers most cases.
		c.dst[c.pDst-1] ^= byte(c.info >> xorShift)
	***REMOVED*** else ***REMOVED***
		// Interpret XOR bits as an index.
		// TODO: test performance for unrolling this loop. Verify that we have
		// at least two bytes and at most three.
		idx := c.info >> xorShift
		for p := c.pDst - 1; ; p-- ***REMOVED***
			c.dst[p] ^= xorData[idx]
			idx--
			if xorData[idx] == 0 ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// hasPrefix returns true if src[pSrc:] starts with the given string.
func (c *context) hasPrefix(s string) bool ***REMOVED***
	b := c.src[c.pSrc:]
	if len(b) < len(s) ***REMOVED***
		return false
	***REMOVED***
	for i, c := range b[:len(s)] ***REMOVED***
		if c != s[i] ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// caseType returns an info with only the case bits, normalized to either
// cLower, cUpper, cTitle or cUncased.
func (c *context) caseType() info ***REMOVED***
	cm := c.info & 0x7
	if cm < 4 ***REMOVED***
		return cm
	***REMOVED***
	if cm >= cXORCase ***REMOVED***
		// xor the last bit of the rune with the case type bits.
		b := c.src[c.pSrc+c.sz-1]
		return info(b&1) ^ cm&0x3
	***REMOVED***
	if cm == cIgnorableCased ***REMOVED***
		return cLower
	***REMOVED***
	return cUncased
***REMOVED***

// lower writes the lowercase version of the current rune to dst.
func lower(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cLower ***REMOVED***
		return c.copy()
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		return c.copyXOR()
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	offset := 2 + e[0]&lengthMask // size of header + fold string
	if nLower := (e[1] >> lengthBits) & lengthMask; nLower != noChange ***REMOVED***
		return c.writeString(e[offset : offset+nLower])
	***REMOVED***
	return c.copy()
***REMOVED***

func isLower(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cLower ***REMOVED***
		return true
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	if nLower := (e[1] >> lengthBits) & lengthMask; nLower != noChange ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	return true
***REMOVED***

// upper writes the uppercase version of the current rune to dst.
func upper(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cUpper ***REMOVED***
		return c.copy()
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		return c.copyXOR()
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	offset := 2 + e[0]&lengthMask // size of header + fold string
	// Get length of first special case mapping.
	n := (e[1] >> lengthBits) & lengthMask
	if ct == cTitle ***REMOVED***
		// The first special case mapping is for lower. Set n to the second.
		if n == noChange ***REMOVED***
			n = 0
		***REMOVED***
		n, e = e[1]&lengthMask, e[n:]
	***REMOVED***
	if n != noChange ***REMOVED***
		return c.writeString(e[offset : offset+n])
	***REMOVED***
	return c.copy()
***REMOVED***

// isUpper writes the isUppercase version of the current rune to dst.
func isUpper(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cUpper ***REMOVED***
		return true
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	// Get length of first special case mapping.
	n := (e[1] >> lengthBits) & lengthMask
	if ct == cTitle ***REMOVED***
		n = e[1] & lengthMask
	***REMOVED***
	if n != noChange ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	return true
***REMOVED***

// title writes the title case version of the current rune to dst.
func title(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cTitle ***REMOVED***
		return c.copy()
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		if ct == cLower ***REMOVED***
			return c.copyXOR()
		***REMOVED***
		return c.copy()
	***REMOVED***
	// Get the exception data.
	e := exceptions[c.info>>exceptionShift:]
	offset := 2 + e[0]&lengthMask // size of header + fold string

	nFirst := (e[1] >> lengthBits) & lengthMask
	if nTitle := e[1] & lengthMask; nTitle != noChange ***REMOVED***
		if nFirst != noChange ***REMOVED***
			e = e[nFirst:]
		***REMOVED***
		return c.writeString(e[offset : offset+nTitle])
	***REMOVED***
	if ct == cLower && nFirst != noChange ***REMOVED***
		// Use the uppercase version instead.
		return c.writeString(e[offset : offset+nFirst])
	***REMOVED***
	// Already in correct case.
	return c.copy()
***REMOVED***

// isTitle reports whether the current rune is in title case.
func isTitle(c *context) bool ***REMOVED***
	ct := c.caseType()
	if c.info&hasMappingMask == 0 || ct == cTitle ***REMOVED***
		return true
	***REMOVED***
	if c.info&exceptionBit == 0 ***REMOVED***
		if ct == cLower ***REMOVED***
			c.err = transform.ErrEndOfSpan
			return false
		***REMOVED***
		return true
	***REMOVED***
	// Get the exception data.
	e := exceptions[c.info>>exceptionShift:]
	if nTitle := e[1] & lengthMask; nTitle != noChange ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	nFirst := (e[1] >> lengthBits) & lengthMask
	if ct == cLower && nFirst != noChange ***REMOVED***
		c.err = transform.ErrEndOfSpan
		return false
	***REMOVED***
	return true
***REMOVED***

// foldFull writes the foldFull version of the current rune to dst.
func foldFull(c *context) bool ***REMOVED***
	if c.info&hasMappingMask == 0 ***REMOVED***
		return c.copy()
	***REMOVED***
	ct := c.caseType()
	if c.info&exceptionBit == 0 ***REMOVED***
		if ct != cLower || c.info&inverseFoldBit != 0 ***REMOVED***
			return c.copyXOR()
		***REMOVED***
		return c.copy()
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	n := e[0] & lengthMask
	if n == 0 ***REMOVED***
		if ct == cLower ***REMOVED***
			return c.copy()
		***REMOVED***
		n = (e[1] >> lengthBits) & lengthMask
	***REMOVED***
	return c.writeString(e[2 : 2+n])
***REMOVED***

// isFoldFull reports whether the current run is mapped to foldFull
func isFoldFull(c *context) bool ***REMOVED***
	if c.info&hasMappingMask == 0 ***REMOVED***
		return true
	***REMOVED***
	ct := c.caseType()
	if c.info&exceptionBit == 0 ***REMOVED***
		if ct != cLower || c.info&inverseFoldBit != 0 ***REMOVED***
			c.err = transform.ErrEndOfSpan
			return false
		***REMOVED***
		return true
	***REMOVED***
	e := exceptions[c.info>>exceptionShift:]
	n := e[0] & lengthMask
	if n == 0 && ct == cLower ***REMOVED***
		return true
	***REMOVED***
	c.err = transform.ErrEndOfSpan
	return false
***REMOVED***
