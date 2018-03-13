// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"fmt"
	"io"
	"reflect"

	"golang.org/x/text/internal/colltab"
)

// table is an intermediate structure that roughly resembles the table in collate.
type table struct ***REMOVED***
	colltab.Table
	trie trie
	root *trieHandle
***REMOVED***

// print writes the table as Go compilable code to w. It prefixes the
// variable names with name. It returns the number of bytes written
// and the size of the resulting table.
func (t *table) fprint(w io.Writer, name string) (n, size int, err error) ***REMOVED***
	update := func(nn, sz int, e error) ***REMOVED***
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
		size += sz
	***REMOVED***
	// Write arrays needed for the structure.
	update(printColElems(w, t.ExpandElem, name+"ExpandElem"))
	update(printColElems(w, t.ContractElem, name+"ContractElem"))
	update(t.trie.printArrays(w, name))
	update(printArray(t.ContractTries, w, name))

	nn, e := fmt.Fprintf(w, "// Total size of %sTable is %d bytes\n", name, size)
	update(nn, 0, e)
	return
***REMOVED***

func (t *table) fprintIndex(w io.Writer, h *trieHandle, id string) (n int, err error) ***REMOVED***
	p := func(f string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	p("\t***REMOVED*** // %s\n", id)
	p("\t\tlookupOffset: 0x%x,\n", h.lookupStart)
	p("\t\tvaluesOffset: 0x%x,\n", h.valueStart)
	p("\t***REMOVED***,\n")
	return
***REMOVED***

func printColElems(w io.Writer, a []uint32, name string) (n, sz int, err error) ***REMOVED***
	p := func(f string, a ...interface***REMOVED******REMOVED***) ***REMOVED***
		nn, e := fmt.Fprintf(w, f, a...)
		n += nn
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
	***REMOVED***
	sz = len(a) * int(reflect.TypeOf(uint32(0)).Size())
	p("// %s: %d entries, %d bytes\n", name, len(a), sz)
	p("var %s = [%d]uint32 ***REMOVED***", name, len(a))
	for i, c := range a ***REMOVED***
		switch ***REMOVED***
		case i%64 == 0:
			p("\n\t// Block %d, offset 0x%x\n", i/64, i)
		case (i%64)%6 == 0:
			p("\n\t")
		***REMOVED***
		p("0x%.8X, ", c)
	***REMOVED***
	p("\n***REMOVED***\n\n")
	return
***REMOVED***
