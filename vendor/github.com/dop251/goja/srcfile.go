package goja

import (
	"fmt"
	"sort"
	"strings"
)

type Position struct ***REMOVED***
	Line, Col int
***REMOVED***

type SrcFile struct ***REMOVED***
	name string
	src  string

	lineOffsets       []int
	lastScannedOffset int
***REMOVED***

func NewSrcFile(name, src string) *SrcFile ***REMOVED***
	return &SrcFile***REMOVED***
		name: name,
		src:  src,
	***REMOVED***
***REMOVED***

func (f *SrcFile) Position(offset int) Position ***REMOVED***
	var line int
	if offset > f.lastScannedOffset ***REMOVED***
		f.scanTo(offset)
		line = len(f.lineOffsets) - 1
	***REMOVED*** else ***REMOVED***
		if len(f.lineOffsets) > 0 ***REMOVED***
			line = sort.SearchInts(f.lineOffsets, offset)
		***REMOVED*** else ***REMOVED***
			line = -1
		***REMOVED***
	***REMOVED***

	if line >= 0 ***REMOVED***
		if f.lineOffsets[line] > offset ***REMOVED***
			line--
		***REMOVED***
	***REMOVED***

	var lineStart int
	if line >= 0 ***REMOVED***
		lineStart = f.lineOffsets[line]
	***REMOVED***
	return Position***REMOVED***
		Line: line + 2,
		Col:  offset - lineStart + 1,
	***REMOVED***
***REMOVED***

func (f *SrcFile) scanTo(offset int) ***REMOVED***
	o := f.lastScannedOffset
	for o < offset ***REMOVED***
		p := strings.Index(f.src[o:], "\n")
		if p == -1 ***REMOVED***
			o = len(f.src)
			break
		***REMOVED***
		o = o + p + 1
		f.lineOffsets = append(f.lineOffsets, o)
	***REMOVED***
	f.lastScannedOffset = o
***REMOVED***

func (p Position) String() string ***REMOVED***
	return fmt.Sprintf("%d:%d", p.Line, p.Col)
***REMOVED***
