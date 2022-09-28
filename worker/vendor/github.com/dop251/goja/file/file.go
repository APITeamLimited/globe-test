// Package file encapsulates the file abstractions used by the ast & parser.
//
package file

import (
	"fmt"
	"net/url"
	"path"
	"sort"
	"sync"

	"github.com/go-sourcemap/sourcemap"
)

// Idx is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
type Idx int

// Position describes an arbitrary source position
// including the filename, line, and column location.
type Position struct ***REMOVED***
	Filename string // The filename where the error occurred, if any
	Line     int    // The line number, starting at 1
	Column   int    // The column number, starting at 1 (The character count)

***REMOVED***

// A Position is valid if the line number is > 0.

func (self *Position) isValid() bool ***REMOVED***
	return self.Line > 0
***REMOVED***

// String returns a string in one of several forms:
//
//	file:line:column    A valid position with filename
//	line:column         A valid position without filename
//	file                An invalid position with filename
//	-                   An invalid position without filename
//
func (self Position) String() string ***REMOVED***
	str := self.Filename
	if self.isValid() ***REMOVED***
		if str != "" ***REMOVED***
			str += ":"
		***REMOVED***
		str += fmt.Sprintf("%d:%d", self.Line, self.Column)
	***REMOVED***
	if str == "" ***REMOVED***
		str = "-"
	***REMOVED***
	return str
***REMOVED***

// FileSet

// A FileSet represents a set of source files.
type FileSet struct ***REMOVED***
	files []*File
	last  *File
***REMOVED***

// AddFile adds a new file with the given filename and src.
//
// This an internal method, but exported for cross-package use.
func (self *FileSet) AddFile(filename, src string) int ***REMOVED***
	base := self.nextBase()
	file := &File***REMOVED***
		name: filename,
		src:  src,
		base: base,
	***REMOVED***
	self.files = append(self.files, file)
	self.last = file
	return base
***REMOVED***

func (self *FileSet) nextBase() int ***REMOVED***
	if self.last == nil ***REMOVED***
		return 1
	***REMOVED***
	return self.last.base + len(self.last.src) + 1
***REMOVED***

func (self *FileSet) File(idx Idx) *File ***REMOVED***
	for _, file := range self.files ***REMOVED***
		if idx <= Idx(file.base+len(file.src)) ***REMOVED***
			return file
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Position converts an Idx in the FileSet into a Position.
func (self *FileSet) Position(idx Idx) Position ***REMOVED***
	for _, file := range self.files ***REMOVED***
		if idx <= Idx(file.base+len(file.src)) ***REMOVED***
			return file.Position(int(idx) - file.base)
		***REMOVED***
	***REMOVED***
	return Position***REMOVED******REMOVED***
***REMOVED***

type File struct ***REMOVED***
	mu                sync.Mutex
	name              string
	src               string
	base              int // This will always be 1 or greater
	sourceMap         *sourcemap.Consumer
	lineOffsets       []int
	lastScannedOffset int
***REMOVED***

func NewFile(filename, src string, base int) *File ***REMOVED***
	return &File***REMOVED***
		name: filename,
		src:  src,
		base: base,
	***REMOVED***
***REMOVED***

func (fl *File) Name() string ***REMOVED***
	return fl.name
***REMOVED***

func (fl *File) Source() string ***REMOVED***
	return fl.src
***REMOVED***

func (fl *File) Base() int ***REMOVED***
	return fl.base
***REMOVED***

func (fl *File) SetSourceMap(m *sourcemap.Consumer) ***REMOVED***
	fl.sourceMap = m
***REMOVED***

func (fl *File) Position(offset int) Position ***REMOVED***
	var line int
	var lineOffsets []int
	fl.mu.Lock()
	if offset > fl.lastScannedOffset ***REMOVED***
		line = fl.scanTo(offset)
		lineOffsets = fl.lineOffsets
		fl.mu.Unlock()
	***REMOVED*** else ***REMOVED***
		lineOffsets = fl.lineOffsets
		fl.mu.Unlock()
		line = sort.Search(len(lineOffsets), func(x int) bool ***REMOVED*** return lineOffsets[x] > offset ***REMOVED***) - 1
	***REMOVED***

	var lineStart int
	if line >= 0 ***REMOVED***
		lineStart = lineOffsets[line]
	***REMOVED***

	row := line + 2
	col := offset - lineStart + 1

	if fl.sourceMap != nil ***REMOVED***
		if source, _, row, col, ok := fl.sourceMap.Source(row, col); ok ***REMOVED***
			return Position***REMOVED***
				Filename: ResolveSourcemapURL(fl.Name(), source).String(),
				Line:     row,
				Column:   col,
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return Position***REMOVED***
		Filename: fl.name,
		Line:     row,
		Column:   col,
	***REMOVED***
***REMOVED***

func ResolveSourcemapURL(basename, source string) *url.URL ***REMOVED***
	// if the url is absolute(has scheme) there is nothing to do
	smURL, err := url.Parse(source)
	if err == nil && !smURL.IsAbs() ***REMOVED***
		baseURL, err1 := url.Parse(basename)
		if err1 == nil && path.IsAbs(baseURL.Path) ***REMOVED***
			smURL = baseURL.ResolveReference(smURL)
		***REMOVED*** else ***REMOVED***
			// pathological case where both are not absolute paths and using Resolve as above will produce an absolute
			// one
			smURL, _ = url.Parse(path.Join(path.Dir(basename), smURL.Path))
		***REMOVED***
	***REMOVED***
	return smURL
***REMOVED***

func findNextLineStart(s string) int ***REMOVED***
	for pos, ch := range s ***REMOVED***
		switch ch ***REMOVED***
		case '\r':
			if pos < len(s)-1 && s[pos+1] == '\n' ***REMOVED***
				return pos + 2
			***REMOVED***
			return pos + 1
		case '\n':
			return pos + 1
		case '\u2028', '\u2029':
			return pos + 3
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

func (fl *File) scanTo(offset int) int ***REMOVED***
	o := fl.lastScannedOffset
	for o < offset ***REMOVED***
		p := findNextLineStart(fl.src[o:])
		if p == -1 ***REMOVED***
			fl.lastScannedOffset = len(fl.src)
			return len(fl.lineOffsets) - 1
		***REMOVED***
		o = o + p
		fl.lineOffsets = append(fl.lineOffsets, o)
	***REMOVED***
	fl.lastScannedOffset = o

	if o == offset ***REMOVED***
		return len(fl.lineOffsets) - 1
	***REMOVED***

	return len(fl.lineOffsets) - 2
***REMOVED***
