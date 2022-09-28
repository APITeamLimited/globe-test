// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package norm

import "unicode/utf8"

type input struct ***REMOVED***
	str   string
	bytes []byte
***REMOVED***

func inputBytes(str []byte) input ***REMOVED***
	return input***REMOVED***bytes: str***REMOVED***
***REMOVED***

func inputString(str string) input ***REMOVED***
	return input***REMOVED***str: str***REMOVED***
***REMOVED***

func (in *input) setBytes(str []byte) ***REMOVED***
	in.str = ""
	in.bytes = str
***REMOVED***

func (in *input) setString(str string) ***REMOVED***
	in.str = str
	in.bytes = nil
***REMOVED***

func (in *input) _byte(p int) byte ***REMOVED***
	if in.bytes == nil ***REMOVED***
		return in.str[p]
	***REMOVED***
	return in.bytes[p]
***REMOVED***

func (in *input) skipASCII(p, max int) int ***REMOVED***
	if in.bytes == nil ***REMOVED***
		for ; p < max && in.str[p] < utf8.RuneSelf; p++ ***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ; p < max && in.bytes[p] < utf8.RuneSelf; p++ ***REMOVED***
		***REMOVED***
	***REMOVED***
	return p
***REMOVED***

func (in *input) skipContinuationBytes(p int) int ***REMOVED***
	if in.bytes == nil ***REMOVED***
		for ; p < len(in.str) && !utf8.RuneStart(in.str[p]); p++ ***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ; p < len(in.bytes) && !utf8.RuneStart(in.bytes[p]); p++ ***REMOVED***
		***REMOVED***
	***REMOVED***
	return p
***REMOVED***

func (in *input) appendSlice(buf []byte, b, e int) []byte ***REMOVED***
	if in.bytes != nil ***REMOVED***
		return append(buf, in.bytes[b:e]...)
	***REMOVED***
	for i := b; i < e; i++ ***REMOVED***
		buf = append(buf, in.str[i])
	***REMOVED***
	return buf
***REMOVED***

func (in *input) copySlice(buf []byte, b, e int) int ***REMOVED***
	if in.bytes == nil ***REMOVED***
		return copy(buf, in.str[b:e])
	***REMOVED***
	return copy(buf, in.bytes[b:e])
***REMOVED***

func (in *input) charinfoNFC(p int) (uint16, int) ***REMOVED***
	if in.bytes == nil ***REMOVED***
		return nfcData.lookupString(in.str[p:])
	***REMOVED***
	return nfcData.lookup(in.bytes[p:])
***REMOVED***

func (in *input) charinfoNFKC(p int) (uint16, int) ***REMOVED***
	if in.bytes == nil ***REMOVED***
		return nfkcData.lookupString(in.str[p:])
	***REMOVED***
	return nfkcData.lookup(in.bytes[p:])
***REMOVED***

func (in *input) hangul(p int) (r rune) ***REMOVED***
	var size int
	if in.bytes == nil ***REMOVED***
		if !isHangulString(in.str[p:]) ***REMOVED***
			return 0
		***REMOVED***
		r, size = utf8.DecodeRuneInString(in.str[p:])
	***REMOVED*** else ***REMOVED***
		if !isHangul(in.bytes[p:]) ***REMOVED***
			return 0
		***REMOVED***
		r, size = utf8.DecodeRune(in.bytes[p:])
	***REMOVED***
	if size != hangulUTF8Size ***REMOVED***
		return 0
	***REMOVED***
	return r
***REMOVED***
