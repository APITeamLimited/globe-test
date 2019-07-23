// Code generated by "stringer -type=blockType,literalsBlockType,seqCompMode,tableIndex"; DO NOT EDIT.

package zstd

import "strconv"

func _() ***REMOVED***
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct***REMOVED******REMOVED***
	_ = x[blockTypeRaw-0]
	_ = x[blockTypeRLE-1]
	_ = x[blockTypeCompressed-2]
	_ = x[blockTypeReserved-3]
***REMOVED***

const _blockType_name = "blockTypeRawblockTypeRLEblockTypeCompressedblockTypeReserved"

var _blockType_index = [...]uint8***REMOVED***0, 12, 24, 43, 60***REMOVED***

func (i blockType) String() string ***REMOVED***
	if i >= blockType(len(_blockType_index)-1) ***REMOVED***
		return "blockType(" + strconv.FormatInt(int64(i), 10) + ")"
	***REMOVED***
	return _blockType_name[_blockType_index[i]:_blockType_index[i+1]]
***REMOVED***
func _() ***REMOVED***
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct***REMOVED******REMOVED***
	_ = x[literalsBlockRaw-0]
	_ = x[literalsBlockRLE-1]
	_ = x[literalsBlockCompressed-2]
	_ = x[literalsBlockTreeless-3]
***REMOVED***

const _literalsBlockType_name = "literalsBlockRawliteralsBlockRLEliteralsBlockCompressedliteralsBlockTreeless"

var _literalsBlockType_index = [...]uint8***REMOVED***0, 16, 32, 55, 76***REMOVED***

func (i literalsBlockType) String() string ***REMOVED***
	if i >= literalsBlockType(len(_literalsBlockType_index)-1) ***REMOVED***
		return "literalsBlockType(" + strconv.FormatInt(int64(i), 10) + ")"
	***REMOVED***
	return _literalsBlockType_name[_literalsBlockType_index[i]:_literalsBlockType_index[i+1]]
***REMOVED***
func _() ***REMOVED***
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct***REMOVED******REMOVED***
	_ = x[compModePredefined-0]
	_ = x[compModeRLE-1]
	_ = x[compModeFSE-2]
	_ = x[compModeRepeat-3]
***REMOVED***

const _seqCompMode_name = "compModePredefinedcompModeRLEcompModeFSEcompModeRepeat"

var _seqCompMode_index = [...]uint8***REMOVED***0, 18, 29, 40, 54***REMOVED***

func (i seqCompMode) String() string ***REMOVED***
	if i >= seqCompMode(len(_seqCompMode_index)-1) ***REMOVED***
		return "seqCompMode(" + strconv.FormatInt(int64(i), 10) + ")"
	***REMOVED***
	return _seqCompMode_name[_seqCompMode_index[i]:_seqCompMode_index[i+1]]
***REMOVED***
func _() ***REMOVED***
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct***REMOVED******REMOVED***
	_ = x[tableLiteralLengths-0]
	_ = x[tableOffsets-1]
	_ = x[tableMatchLengths-2]
***REMOVED***

const _tableIndex_name = "tableLiteralLengthstableOffsetstableMatchLengths"

var _tableIndex_index = [...]uint8***REMOVED***0, 19, 31, 48***REMOVED***

func (i tableIndex) String() string ***REMOVED***
	if i >= tableIndex(len(_tableIndex_index)-1) ***REMOVED***
		return "tableIndex(" + strconv.FormatInt(int64(i), 10) + ")"
	***REMOVED***
	return _tableIndex_name[_tableIndex_index[i]:_tableIndex_index[i+1]]
***REMOVED***
