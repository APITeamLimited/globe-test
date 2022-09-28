// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cases

import "golang.org/x/text/transform"

type caseFolder struct***REMOVED*** transform.NopResetter ***REMOVED***

// caseFolder implements the Transformer interface for doing case folding.
func (t *caseFolder) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) ***REMOVED***
	c := context***REMOVED***dst: dst, src: src, atEOF: atEOF***REMOVED***
	for c.next() ***REMOVED***
		foldFull(&c)
		c.checkpoint()
	***REMOVED***
	return c.ret()
***REMOVED***

func (t *caseFolder) Span(src []byte, atEOF bool) (n int, err error) ***REMOVED***
	c := context***REMOVED***src: src, atEOF: atEOF***REMOVED***
	for c.next() && isFoldFull(&c) ***REMOVED***
		c.checkpoint()
	***REMOVED***
	return c.retSpan()
***REMOVED***

func makeFold(o options) transform.SpanningTransformer ***REMOVED***
	// TODO: Special case folding, through option Language, Special/Turkic, or
	// both.
	// TODO: Implement Compact options.
	return &caseFolder***REMOVED******REMOVED***
***REMOVED***
