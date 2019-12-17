package goja

import (
	"fmt"
	"github.com/go-sourcemap/sourcemap"
	"sort"
	"strings"
	"sync"
)

type Position struct ***REMOVED***
	Line, Col int
***REMOVED***

type SrcFile struct ***REMOVED***
	name string
	src  string

	lineOffsets       []int
	lineOffsetsLock   sync.Mutex
	lastScannedOffset int
	sourceMap         *sourcemap.Consumer
***REMOVED***

func NewSrcFile(name, src string, sourceMap *sourcemap.Consumer) *SrcFile ***REMOVED***
	return &SrcFile***REMOVED***
		name:      name,
		src:       src,
		sourceMap: sourceMap,
	***REMOVED***
***REMOVED***

func (f *SrcFile) Position(offset int) Position ***REMOVED***
	var line int
	var lineOffsets []int
	f.lineOffsetsLock.Lock()
	if offset > f.lastScannedOffset ***REMOVED***
		line = f.scanTo(offset)
		lineOffsets = f.lineOffsets
		f.lineOffsetsLock.Unlock()
	***REMOVED*** else ***REMOVED***
		lineOffsets = f.lineOffsets
		f.lineOffsetsLock.Unlock()
		line = sort.Search(len(lineOffsets), func(x int) bool ***REMOVED*** return lineOffsets[x] > offset ***REMOVED***) - 1
	***REMOVED***

	var lineStart int
	if line >= 0 ***REMOVED***
		lineStart = lineOffsets[line]
	***REMOVED***

	row := line + 2
	col := offset - lineStart + 1

	if f.sourceMap != nil ***REMOVED***
		if _, _, row, col, ok := f.sourceMap.Source(row, col); ok ***REMOVED***
			return Position***REMOVED***
				Line: row,
				Col:  col,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return Position***REMOVED***
		Line: row,
		Col:  col,
	***REMOVED***
***REMOVED***

func (f *SrcFile) scanTo(offset int) int ***REMOVED***
	o := f.lastScannedOffset
	for o < offset ***REMOVED***
		p := strings.Index(f.src[o:], "\n")
		if p == -1 ***REMOVED***
			f.lastScannedOffset = len(f.src)
			return len(f.lineOffsets) - 1
		***REMOVED***
		o = o + p + 1
		f.lineOffsets = append(f.lineOffsets, o)
	***REMOVED***
	f.lastScannedOffset = o

	if o == offset ***REMOVED***
		return len(f.lineOffsets) - 1
	***REMOVED***

	return len(f.lineOffsets) - 2
***REMOVED***

func (p Position) String() string ***REMOVED***
	return fmt.Sprintf("%d:%d", p.Line, p.Col)
***REMOVED***
