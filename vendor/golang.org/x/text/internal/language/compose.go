// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"sort"
	"strings"
)

// A Builder allows constructing a Tag from individual components.
// Its main user is Compose in the top-level language package.
type Builder struct ***REMOVED***
	Tag Tag

	private    string // the x extension
	variants   []string
	extensions []string
***REMOVED***

// Make returns a new Tag from the current settings.
func (b *Builder) Make() Tag ***REMOVED***
	t := b.Tag

	if len(b.extensions) > 0 || len(b.variants) > 0 ***REMOVED***
		sort.Sort(sortVariants(b.variants))
		sort.Strings(b.extensions)

		if b.private != "" ***REMOVED***
			b.extensions = append(b.extensions, b.private)
		***REMOVED***
		n := maxCoreSize + tokenLen(b.variants...) + tokenLen(b.extensions...)
		buf := make([]byte, n)
		p := t.genCoreBytes(buf)
		t.pVariant = byte(p)
		p += appendTokens(buf[p:], b.variants...)
		t.pExt = uint16(p)
		p += appendTokens(buf[p:], b.extensions...)
		t.str = string(buf[:p])
		// We may not always need to remake the string, but when or when not
		// to do so is rather tricky.
		scan := makeScanner(buf[:p])
		t, _ = parse(&scan, "")
		return t

	***REMOVED*** else if b.private != "" ***REMOVED***
		t.str = b.private
		t.RemakeString()
	***REMOVED***
	return t
***REMOVED***

// SetTag copies all the settings from a given Tag. Any previously set values
// are discarded.
func (b *Builder) SetTag(t Tag) ***REMOVED***
	b.Tag.LangID = t.LangID
	b.Tag.RegionID = t.RegionID
	b.Tag.ScriptID = t.ScriptID
	// TODO: optimize
	b.variants = b.variants[:0]
	if variants := t.Variants(); variants != "" ***REMOVED***
		for _, vr := range strings.Split(variants[1:], "-") ***REMOVED***
			b.variants = append(b.variants, vr)
		***REMOVED***
	***REMOVED***
	b.extensions, b.private = b.extensions[:0], ""
	for _, e := range t.Extensions() ***REMOVED***
		b.AddExt(e)
	***REMOVED***
***REMOVED***

// AddExt adds extension e to the tag. e must be a valid extension as returned
// by Tag.Extension. If the extension already exists, it will be discarded,
// except for a -u extension, where non-existing key-type pairs will added.
func (b *Builder) AddExt(e string) ***REMOVED***
	if e[0] == 'x' ***REMOVED***
		if b.private == "" ***REMOVED***
			b.private = e
		***REMOVED***
		return
	***REMOVED***
	for i, s := range b.extensions ***REMOVED***
		if s[0] == e[0] ***REMOVED***
			if e[0] == 'u' ***REMOVED***
				b.extensions[i] += e[1:]
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	b.extensions = append(b.extensions, e)
***REMOVED***

// SetExt sets the extension e to the tag. e must be a valid extension as
// returned by Tag.Extension. If the extension already exists, it will be
// overwritten, except for a -u extension, where the individual key-type pairs
// will be set.
func (b *Builder) SetExt(e string) ***REMOVED***
	if e[0] == 'x' ***REMOVED***
		b.private = e
		return
	***REMOVED***
	for i, s := range b.extensions ***REMOVED***
		if s[0] == e[0] ***REMOVED***
			if e[0] == 'u' ***REMOVED***
				b.extensions[i] = e + s[1:]
			***REMOVED*** else ***REMOVED***
				b.extensions[i] = e
			***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	b.extensions = append(b.extensions, e)
***REMOVED***

// AddVariant adds any number of variants.
func (b *Builder) AddVariant(v ...string) ***REMOVED***
	for _, v := range v ***REMOVED***
		if v != "" ***REMOVED***
			b.variants = append(b.variants, v)
		***REMOVED***
	***REMOVED***
***REMOVED***

// ClearVariants removes any variants previously added, including those
// copied from a Tag in SetTag.
func (b *Builder) ClearVariants() ***REMOVED***
	b.variants = b.variants[:0]
***REMOVED***

// ClearExtensions removes any extensions previously added, including those
// copied from a Tag in SetTag.
func (b *Builder) ClearExtensions() ***REMOVED***
	b.private = ""
	b.extensions = b.extensions[:0]
***REMOVED***

func tokenLen(token ...string) (n int) ***REMOVED***
	for _, t := range token ***REMOVED***
		n += len(t) + 1
	***REMOVED***
	return
***REMOVED***

func appendTokens(b []byte, token ...string) int ***REMOVED***
	p := 0
	for _, t := range token ***REMOVED***
		b[p] = '-'
		copy(b[p+1:], t)
		p += 1 + len(t)
	***REMOVED***
	return p
***REMOVED***

type sortVariants []string

func (s sortVariants) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sortVariants) Swap(i, j int) ***REMOVED***
	s[j], s[i] = s[i], s[j]
***REMOVED***

func (s sortVariants) Less(i, j int) bool ***REMOVED***
	return variantIndex[s[i]] < variantIndex[s[j]]
***REMOVED***
