package brotli

/* Copyright 2010 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Write bits into a byte array. */

type bitWriter struct ***REMOVED***
	dst []byte

	// Data waiting to be written is the low nbits of bits.
	bits  uint64
	nbits uint
***REMOVED***

func (w *bitWriter) writeBits(nb uint, b uint64) ***REMOVED***
	w.bits |= b << w.nbits
	w.nbits += nb
	if w.nbits >= 32 ***REMOVED***
		bits := w.bits
		w.bits >>= 32
		w.nbits -= 32
		w.dst = append(w.dst,
			byte(bits),
			byte(bits>>8),
			byte(bits>>16),
			byte(bits>>24),
		)
	***REMOVED***
***REMOVED***

func (w *bitWriter) writeSingleBit(bit bool) ***REMOVED***
	if bit ***REMOVED***
		w.writeBits(1, 1)
	***REMOVED*** else ***REMOVED***
		w.writeBits(1, 0)
	***REMOVED***
***REMOVED***

func (w *bitWriter) jumpToByteBoundary() ***REMOVED***
	dst := w.dst
	for w.nbits != 0 ***REMOVED***
		dst = append(dst, byte(w.bits))
		w.bits >>= 8
		if w.nbits > 8 ***REMOVED*** // Avoid underflow
			w.nbits -= 8
		***REMOVED*** else ***REMOVED***
			w.nbits = 0
		***REMOVED***
	***REMOVED***
	w.bits = 0
	w.dst = dst
***REMOVED***

func (w *bitWriter) writeBytes(b []byte) ***REMOVED***
	if w.nbits&7 != 0 ***REMOVED***
		panic("writeBytes with unfinished bits")
	***REMOVED***
	for w.nbits != 0 ***REMOVED***
		w.dst = append(w.dst, byte(w.bits))
		w.bits >>= 8
		w.nbits -= 8
	***REMOVED***
	w.dst = append(w.dst, b...)
***REMOVED***

func (w *bitWriter) getPos() uint ***REMOVED***
	return uint(len(w.dst)<<3) + w.nbits
***REMOVED***

func (w *bitWriter) rewind(p uint) ***REMOVED***
	w.bits = uint64(w.dst[p>>3] & byte((1<<(p&7))-1))
	w.nbits = p & 7
	w.dst = w.dst[:p>>3]
***REMOVED***

func (w *bitWriter) updateBits(n_bits uint, bits uint32, pos uint) ***REMOVED***
	for n_bits > 0 ***REMOVED***
		var byte_pos uint = pos >> 3
		var n_unchanged_bits uint = pos & 7
		var n_changed_bits uint = brotli_min_size_t(n_bits, 8-n_unchanged_bits)
		var total_bits uint = n_unchanged_bits + n_changed_bits
		var mask uint32 = (^((1 << total_bits) - 1)) | ((1 << n_unchanged_bits) - 1)
		var unchanged_bits uint32 = uint32(w.dst[byte_pos]) & mask
		var changed_bits uint32 = bits & ((1 << n_changed_bits) - 1)
		w.dst[byte_pos] = byte(changed_bits<<n_unchanged_bits | unchanged_bits)
		n_bits -= n_changed_bits
		bits >>= n_changed_bits
		pos += n_changed_bits
	***REMOVED***
***REMOVED***
