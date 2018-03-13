// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"reflect"
	"strings"
	"unicode"

	"golang.org/x/text/collate"
	"golang.org/x/text/internal/gen"
	"golang.org/x/text/internal/ucd"
	"golang.org/x/text/language"
	"golang.org/x/text/unicode/rangetable"
)

var versionList = flag.String("versions", "",
	"list of versions for which to generate RangeTables")

const bootstrapMessage = `No versions specified.
To bootstrap the code generation, run:
	go run gen.go --versions=4.1.0,5.0.0,6.0.0,6.1.0,6.2.0,6.3.0,7.0.0

and ensure that the latest versions are included by checking:
	http://www.unicode.org/Public/`

func getVersions() []string ***REMOVED***
	if *versionList == "" ***REMOVED***
		log.Fatal(bootstrapMessage)
	***REMOVED***

	c := collate.New(language.Und, collate.Numeric)
	versions := strings.Split(*versionList, ",")
	c.SortStrings(versions)

	// Ensure that at least the current version is included.
	for _, v := range versions ***REMOVED***
		if v == gen.UnicodeVersion() ***REMOVED***
			return versions
		***REMOVED***
	***REMOVED***

	versions = append(versions, gen.UnicodeVersion())
	c.SortStrings(versions)
	return versions
***REMOVED***

func main() ***REMOVED***
	gen.Init()

	versions := getVersions()

	w := &bytes.Buffer***REMOVED******REMOVED***

	fmt.Fprintf(w, "//go:generate go run gen.go --versions=%s\n\n", strings.Join(versions, ","))
	fmt.Fprintf(w, "import \"unicode\"\n\n")

	vstr := func(s string) string ***REMOVED*** return strings.Replace(s, ".", "_", -1) ***REMOVED***

	fmt.Fprintf(w, "var assigned = map[string]*unicode.RangeTable***REMOVED***\n")
	for _, v := range versions ***REMOVED***
		fmt.Fprintf(w, "\t%q: assigned%s,\n", v, vstr(v))
	***REMOVED***
	fmt.Fprintf(w, "***REMOVED***\n\n")

	var size int
	for _, v := range versions ***REMOVED***
		assigned := []rune***REMOVED******REMOVED***

		r := gen.Open("http://www.unicode.org/Public/", "", v+"/ucd/UnicodeData.txt")
		ucd.Parse(r, func(p *ucd.Parser) ***REMOVED***
			assigned = append(assigned, p.Rune(0))
		***REMOVED***)

		rt := rangetable.New(assigned...)
		sz := int(reflect.TypeOf(unicode.RangeTable***REMOVED******REMOVED***).Size())
		sz += int(reflect.TypeOf(unicode.Range16***REMOVED******REMOVED***).Size()) * len(rt.R16)
		sz += int(reflect.TypeOf(unicode.Range32***REMOVED******REMOVED***).Size()) * len(rt.R32)

		fmt.Fprintf(w, "// size %d bytes (%d KiB)\n", sz, sz/1024)
		fmt.Fprintf(w, "var assigned%s = ", vstr(v))
		print(w, rt)

		size += sz
	***REMOVED***

	fmt.Fprintf(w, "// Total size %d bytes (%d KiB)\n", size, size/1024)

	gen.WriteVersionedGoFile("tables.go", "rangetable", w.Bytes())
***REMOVED***

func print(w io.Writer, rt *unicode.RangeTable) ***REMOVED***
	fmt.Fprintln(w, "&unicode.RangeTable***REMOVED***")
	fmt.Fprintln(w, "\tR16: []unicode.Range16***REMOVED***")
	for _, r := range rt.R16 ***REMOVED***
		fmt.Fprintf(w, "\t\t***REMOVED***%#04x, %#04x, %d***REMOVED***,\n", r.Lo, r.Hi, r.Stride)
	***REMOVED***
	fmt.Fprintln(w, "\t***REMOVED***,")
	fmt.Fprintln(w, "\tR32: []unicode.Range32***REMOVED***")
	for _, r := range rt.R32 ***REMOVED***
		fmt.Fprintf(w, "\t\t***REMOVED***%#08x, %#08x, %d***REMOVED***,\n", r.Lo, r.Hi, r.Stride)
	***REMOVED***
	fmt.Fprintln(w, "\t***REMOVED***,")
	fmt.Fprintf(w, "\tLatinOffset: %d,\n", rt.LatinOffset)
	fmt.Fprintf(w, "***REMOVED***\n\n")
***REMOVED***
