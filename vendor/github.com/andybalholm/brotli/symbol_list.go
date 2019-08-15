package brotli

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Utilities for building Huffman decoding tables. */

type symbolList struct ***REMOVED***
	storage []uint16
	offset  int
***REMOVED***

func symbolListGet(sl symbolList, i int) uint16 ***REMOVED***
	return sl.storage[i+sl.offset]
***REMOVED***

func symbolListPut(sl symbolList, i int, val uint16) ***REMOVED***
	sl.storage[i+sl.offset] = val
***REMOVED***
