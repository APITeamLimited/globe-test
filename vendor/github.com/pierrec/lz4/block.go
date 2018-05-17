package lz4

import (
	"encoding/binary"
	"errors"
)

// block represents a frame data block.
// Used when compressing or decompressing frame blocks concurrently.
type block struct ***REMOVED***
	compressed bool
	zdata      []byte // compressed data
	data       []byte // decompressed data
	offset     int    // offset within the data as with block dependency the 64Kb window is prepended to it
	checksum   uint32 // compressed data checksum
	err        error  // error while [de]compressing
***REMOVED***

var (
	// ErrInvalidSource is returned by UncompressBlock when a compressed block is corrupted.
	ErrInvalidSource = errors.New("lz4: invalid source")
	// ErrShortBuffer is returned by UncompressBlock, CompressBlock or CompressBlockHC when
	// the supplied buffer for [de]compression is too small.
	ErrShortBuffer = errors.New("lz4: short buffer")
)

// CompressBlockBound returns the maximum size of a given buffer of size n, when not compressible.
func CompressBlockBound(n int) int ***REMOVED***
	return n + n/255 + 16
***REMOVED***

// UncompressBlock decompresses the source buffer into the destination one,
// starting at the di index and returning the decompressed size.
//
// The destination buffer must be sized appropriately.
//
// An error is returned if the source data is invalid or the destination buffer is too small.
func UncompressBlock(src, dst []byte, di int) (int, error) ***REMOVED***
	si, sn, di0 := 0, len(src), di
	if sn == 0 ***REMOVED***
		return 0, nil
	***REMOVED***

	for ***REMOVED***
		// literals and match lengths (token)
		lLen := int(src[si] >> 4)
		mLen := int(src[si] & 0xF)
		if si++; si == sn ***REMOVED***
			return di, ErrInvalidSource
		***REMOVED***

		// literals
		if lLen > 0 ***REMOVED***
			if lLen == 0xF ***REMOVED***
				for src[si] == 0xFF ***REMOVED***
					lLen += 0xFF
					if si++; si == sn ***REMOVED***
						return di - di0, ErrInvalidSource
					***REMOVED***
				***REMOVED***
				lLen += int(src[si])
				if si++; si == sn ***REMOVED***
					return di - di0, ErrInvalidSource
				***REMOVED***
			***REMOVED***
			if len(dst)-di < lLen || si+lLen > sn ***REMOVED***
				return di - di0, ErrShortBuffer
			***REMOVED***
			di += copy(dst[di:], src[si:si+lLen])

			if si += lLen; si >= sn ***REMOVED***
				return di - di0, nil
			***REMOVED***
		***REMOVED***

		if si += 2; si >= sn ***REMOVED***
			return di, ErrInvalidSource
		***REMOVED***
		offset := int(src[si-2]) | int(src[si-1])<<8
		if di-offset < 0 || offset == 0 ***REMOVED***
			return di - di0, ErrInvalidSource
		***REMOVED***

		// match
		if mLen == 0xF ***REMOVED***
			for src[si] == 0xFF ***REMOVED***
				mLen += 0xFF
				if si++; si == sn ***REMOVED***
					return di - di0, ErrInvalidSource
				***REMOVED***
			***REMOVED***
			mLen += int(src[si])
			if si++; si == sn ***REMOVED***
				return di - di0, ErrInvalidSource
			***REMOVED***
		***REMOVED***
		// minimum match length is 4
		mLen += 4
		if len(dst)-di <= mLen ***REMOVED***
			return di - di0, ErrShortBuffer
		***REMOVED***

		// copy the match (NB. match is at least 4 bytes long)
		if mLen >= offset ***REMOVED***
			bytesToCopy := offset * (mLen / offset)
			// Efficiently copy the match dst[di-offset:di] into the slice
			// dst[di:di+bytesToCopy]
			expanded := dst[di-offset : di+bytesToCopy]
			n := offset
			for n <= bytesToCopy+offset ***REMOVED***
				copy(expanded[n:], expanded[:n])
				n *= 2
			***REMOVED***
			di += bytesToCopy
			mLen -= bytesToCopy
		***REMOVED***

		di += copy(dst[di:], dst[di-offset:di-offset+mLen])
	***REMOVED***
***REMOVED***

// CompressBlock compresses the source buffer starting at soffet into the destination one.
// This is the fast version of LZ4 compression and also the default one.
//
// The size of the compressed data is returned. If it is 0 and no error, then the data is incompressible.
//
// An error is returned if the destination buffer is too small.
func CompressBlock(src, dst []byte, soffset int) (int, error) ***REMOVED***
	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 || soffset >= sn ***REMOVED***
		return 0, nil
	***REMOVED***
	var si, di int

	// fast scan strategy:
	// we only need a hash table to store the last sequences (4 bytes)
	var hashTable [1 << hashLog]int
	var hashShift = uint((minMatch * 8) - hashLog)

	// Initialise the hash table with the first 64Kb of the input buffer
	// (used when compressing dependent blocks)
	for si < soffset ***REMOVED***
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		si++
		hashTable[h] = si
	***REMOVED***

	anchor := si
	fma := 1 << skipStrength
	for si < sn-minMatch ***REMOVED***
		// hash the next 4 bytes (sequence)...
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		// -1 to separate existing entries from new ones
		ref := hashTable[h] - 1
		// ...and store the position of the hash in the hash table (+1 to compensate the -1 upon saving)
		hashTable[h] = si + 1
		// no need to check the last 3 bytes in the first literal 4 bytes as
		// this guarantees that the next match, if any, is compressed with
		// a lower size, since to have some compression we must have:
		// ll+ml-overlap > 1 + (ll-15)/255 + (ml-4-15)/255 + 2 (uncompressed size>compressed size)
		// => ll+ml>3+2*overlap => ll+ml>= 4+2*overlap
		// and by definition we do have:
		// ll >= 1, ml >= 4
		// => ll+ml >= 5
		// => so overlap must be 0

		// the sequence is new, out of bound (64kb) or not valid: try next sequence
		if ref < 0 || fma&(1<<skipStrength-1) < 4 ||
			(si-ref)>>winSizeLog > 0 ||
			src[ref] != src[si] ||
			src[ref+1] != src[si+1] ||
			src[ref+2] != src[si+2] ||
			src[ref+3] != src[si+3] ***REMOVED***
			// variable step: improves performance on non-compressible data
			si += fma >> skipStrength
			fma++
			continue
		***REMOVED***
		// match found
		fma = 1 << skipStrength
		lLen := si - anchor
		offset := si - ref

		// encode match length part 1
		si += minMatch
		mLen := si // match length has minMatch already
		for si <= sn && src[si] == src[si-offset] ***REMOVED***
			si++
		***REMOVED***
		mLen = si - mLen
		if mLen < 0xF ***REMOVED***
			dst[di] = byte(mLen)
		***REMOVED*** else ***REMOVED***
			dst[di] = 0xF
		***REMOVED***

		// encode literals length
		if lLen < 0xF ***REMOVED***
			dst[di] |= byte(lLen << 4)
		***REMOVED*** else ***REMOVED***
			dst[di] |= 0xF0
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF ***REMOVED***
				dst[di] = 0xFF
				if di++; di == dn ***REMOVED***
					return di, ErrShortBuffer
				***REMOVED***
			***REMOVED***
			dst[di] = byte(l)
		***REMOVED***
		if di++; di == dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***

		// literals
		if di+lLen >= dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		di += copy(dst[di:], src[anchor:anchor+lLen])
		anchor = si

		// encode offset
		if di += 2; di >= dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		// encode match length part 2
		if mLen >= 0xF ***REMOVED***
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF ***REMOVED***
				dst[di] = 0xFF
				if di++; di == dn ***REMOVED***
					return di, ErrShortBuffer
				***REMOVED***
			***REMOVED***
			dst[di] = byte(mLen)
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if anchor == 0 ***REMOVED***
		// incompressible
		return 0, nil
	***REMOVED***

	// last literals
	lLen := len(src) - anchor
	if lLen < 0xF ***REMOVED***
		dst[di] = byte(lLen << 4)
	***REMOVED*** else ***REMOVED***
		dst[di] = 0xF0
		if di++; di == dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF ***REMOVED***
			dst[di] = 0xFF
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
		***REMOVED***
		dst[di] = byte(lLen)
	***REMOVED***
	if di++; di == dn ***REMOVED***
		return di, ErrShortBuffer
	***REMOVED***

	// write literals
	src = src[anchor:]
	switch n := di + len(src); ***REMOVED***
	case n > dn:
		return di, ErrShortBuffer
	case n >= sn:
		// incompressible
		return 0, nil
	***REMOVED***
	di += copy(dst[di:], src)
	return di, nil
***REMOVED***

// CompressBlockHC compresses the source buffer starting at soffet into the destination one.
// CompressBlockHC compression ratio is better than CompressBlock but it is also slower.
//
// The size of the compressed data is returned. If it is 0 and no error, then the data is not compressible.
//
// An error is returned if the destination buffer is too small.
func CompressBlockHC(src, dst []byte, soffset int) (int, error) ***REMOVED***
	sn, dn := len(src)-mfLimit, len(dst)
	if sn <= 0 || dn == 0 || soffset >= sn ***REMOVED***
		return 0, nil
	***REMOVED***
	var si, di int

	// Hash Chain strategy:
	// we need a hash table and a chain table
	// the chain table cannot contain more entries than the window size (64Kb entries)
	var hashTable [1 << hashLog]int
	var chainTable [winSize]int
	var hashShift = uint((minMatch * 8) - hashLog)

	// Initialise the hash table with the first 64Kb of the input buffer
	// (used when compressing dependent blocks)
	for si < soffset ***REMOVED***
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
		chainTable[si&winMask] = hashTable[h]
		si++
		hashTable[h] = si
	***REMOVED***

	anchor := si
	for si < sn-minMatch ***REMOVED***
		// hash the next 4 bytes (sequence)...
		h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift

		// follow the chain until out of window and give the longest match
		mLen := 0
		offset := 0
		for next := hashTable[h] - 1; next > 0 && next > si-winSize; next = chainTable[next&winMask] - 1 ***REMOVED***
			// the first (mLen==0) or next byte (mLen>=minMatch) at current match length must match to improve on the match length
			if src[next+mLen] == src[si+mLen] ***REMOVED***
				for ml := 0; ; ml++ ***REMOVED***
					if src[next+ml] != src[si+ml] || si+ml > sn ***REMOVED***
						// found a longer match, keep its position and length
						if mLen < ml && ml >= minMatch ***REMOVED***
							mLen = ml
							offset = si - next
						***REMOVED***
						break
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		chainTable[si&winMask] = hashTable[h]
		hashTable[h] = si + 1

		// no match found
		if mLen == 0 ***REMOVED***
			si++
			continue
		***REMOVED***

		// match found
		// update hash/chain tables with overlaping bytes:
		// si already hashed, add everything from si+1 up to the match length
		for si, ml := si+1, si+mLen; si < ml; ***REMOVED***
			h := binary.LittleEndian.Uint32(src[si:]) * hasher >> hashShift
			chainTable[si&winMask] = hashTable[h]
			si++
			hashTable[h] = si
		***REMOVED***

		lLen := si - anchor
		si += mLen
		mLen -= minMatch // match length does not include minMatch

		if mLen < 0xF ***REMOVED***
			dst[di] = byte(mLen)
		***REMOVED*** else ***REMOVED***
			dst[di] = 0xF
		***REMOVED***

		// encode literals length
		if lLen < 0xF ***REMOVED***
			dst[di] |= byte(lLen << 4)
		***REMOVED*** else ***REMOVED***
			dst[di] |= 0xF0
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
			l := lLen - 0xF
			for ; l >= 0xFF; l -= 0xFF ***REMOVED***
				dst[di] = 0xFF
				if di++; di == dn ***REMOVED***
					return di, ErrShortBuffer
				***REMOVED***
			***REMOVED***
			dst[di] = byte(l)
		***REMOVED***
		if di++; di == dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***

		// literals
		if di+lLen >= dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		di += copy(dst[di:], src[anchor:anchor+lLen])
		anchor = si

		// encode offset
		if di += 2; di >= dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		dst[di-2], dst[di-1] = byte(offset), byte(offset>>8)

		// encode match length part 2
		if mLen >= 0xF ***REMOVED***
			for mLen -= 0xF; mLen >= 0xFF; mLen -= 0xFF ***REMOVED***
				dst[di] = 0xFF
				if di++; di == dn ***REMOVED***
					return di, ErrShortBuffer
				***REMOVED***
			***REMOVED***
			dst[di] = byte(mLen)
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if anchor == 0 ***REMOVED***
		// incompressible
		return 0, nil
	***REMOVED***

	// last literals
	lLen := len(src) - anchor
	if lLen < 0xF ***REMOVED***
		dst[di] = byte(lLen << 4)
	***REMOVED*** else ***REMOVED***
		dst[di] = 0xF0
		if di++; di == dn ***REMOVED***
			return di, ErrShortBuffer
		***REMOVED***
		lLen -= 0xF
		for ; lLen >= 0xFF; lLen -= 0xFF ***REMOVED***
			dst[di] = 0xFF
			if di++; di == dn ***REMOVED***
				return di, ErrShortBuffer
			***REMOVED***
		***REMOVED***
		dst[di] = byte(lLen)
	***REMOVED***
	if di++; di == dn ***REMOVED***
		return di, ErrShortBuffer
	***REMOVED***

	// write literals
	src = src[anchor:]
	switch n := di + len(src); ***REMOVED***
	case n > dn:
		return di, ErrShortBuffer
	case n >= sn:
		// incompressible
		return 0, nil
	***REMOVED***
	di += copy(dst[di:], src)
	return di, nil
***REMOVED***
