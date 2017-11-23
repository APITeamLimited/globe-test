// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"log"

	"golang.org/x/text/internal/gen"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/cldr"
)

func main() ***REMOVED***
	r := gen.OpenCLDRCoreZip()
	defer r.Close()

	d := &cldr.Decoder***REMOVED******REMOVED***
	data, err := d.DecodeZip(r)
	if err != nil ***REMOVED***
		log.Fatalf("DecodeZip: %v", err)
	***REMOVED***

	w := gen.NewCodeWriter()
	defer w.WriteGoFile("tables.go", "internal")

	// Create parents table.
	parents := make([]uint16, language.NumCompactTags)
	for _, loc := range data.Locales() ***REMOVED***
		tag := language.MustParse(loc)
		index, ok := language.CompactIndex(tag)
		if !ok ***REMOVED***
			continue
		***REMOVED***
		parentIndex := 0 // und
		for p := tag.Parent(); p != language.Und; p = p.Parent() ***REMOVED***
			if x, ok := language.CompactIndex(p); ok ***REMOVED***
				parentIndex = x
				break
			***REMOVED***
		***REMOVED***
		parents[index] = uint16(parentIndex)
	***REMOVED***

	w.WriteComment(`
	Parent maps a compact index of a tag to the compact index of the parent of
	this tag.`)
	w.WriteVar("Parent", parents)
***REMOVED***
