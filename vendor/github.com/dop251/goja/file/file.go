// Package file encapsulates the file abstractions used by the ast & parser.
//
package file

import (
	"fmt"
	"strings"
)

// Idx is a compact encoding of a source position within a file set.
// It can be converted into a Position for a more convenient, but much
// larger, representation.
type Idx int

// Position describes an arbitrary source position
// including the filename, line, and column location.
type Position struct ***REMOVED***
	Filename string // The filename where the error occurred, if any
	Offset   int    // The src offset
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
func (self *Position) String() string ***REMOVED***
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
func (self *FileSet) Position(idx Idx) *Position ***REMOVED***
	position := &Position***REMOVED******REMOVED***
	for _, file := range self.files ***REMOVED***
		if idx <= Idx(file.base+len(file.src)) ***REMOVED***
			offset := int(idx) - file.base
			src := file.src[:offset]
			position.Filename = file.name
			position.Offset = offset
			position.Line = 1 + strings.Count(src, "\n")
			if index := strings.LastIndex(src, "\n"); index >= 0 ***REMOVED***
				position.Column = offset - index
			***REMOVED*** else ***REMOVED***
				position.Column = 1 + len(src)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return position
***REMOVED***

type File struct ***REMOVED***
	name string
	src  string
	base int // This will always be 1 or greater
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
