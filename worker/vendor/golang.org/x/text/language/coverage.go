// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import (
	"fmt"
	"sort"

	"golang.org/x/text/internal/language"
)

// The Coverage interface is used to define the level of coverage of an
// internationalization service. Note that not all types are supported by all
// services. As lists may be generated on the fly, it is recommended that users
// of a Coverage cache the results.
type Coverage interface ***REMOVED***
	// Tags returns the list of supported tags.
	Tags() []Tag

	// BaseLanguages returns the list of supported base languages.
	BaseLanguages() []Base

	// Scripts returns the list of supported scripts.
	Scripts() []Script

	// Regions returns the list of supported regions.
	Regions() []Region
***REMOVED***

var (
	// Supported defines a Coverage that lists all supported subtags. Tags
	// always returns nil.
	Supported Coverage = allSubtags***REMOVED******REMOVED***
)

// TODO:
// - Support Variants, numbering systems.
// - CLDR coverage levels.
// - Set of common tags defined in this package.

type allSubtags struct***REMOVED******REMOVED***

// Regions returns the list of supported regions. As all regions are in a
// consecutive range, it simply returns a slice of numbers in increasing order.
// The "undefined" region is not returned.
func (s allSubtags) Regions() []Region ***REMOVED***
	reg := make([]Region, language.NumRegions)
	for i := range reg ***REMOVED***
		reg[i] = Region***REMOVED***language.Region(i + 1)***REMOVED***
	***REMOVED***
	return reg
***REMOVED***

// Scripts returns the list of supported scripts. As all scripts are in a
// consecutive range, it simply returns a slice of numbers in increasing order.
// The "undefined" script is not returned.
func (s allSubtags) Scripts() []Script ***REMOVED***
	scr := make([]Script, language.NumScripts)
	for i := range scr ***REMOVED***
		scr[i] = Script***REMOVED***language.Script(i + 1)***REMOVED***
	***REMOVED***
	return scr
***REMOVED***

// BaseLanguages returns the list of all supported base languages. It generates
// the list by traversing the internal structures.
func (s allSubtags) BaseLanguages() []Base ***REMOVED***
	bs := language.BaseLanguages()
	base := make([]Base, len(bs))
	for i, b := range bs ***REMOVED***
		base[i] = Base***REMOVED***b***REMOVED***
	***REMOVED***
	return base
***REMOVED***

// Tags always returns nil.
func (s allSubtags) Tags() []Tag ***REMOVED***
	return nil
***REMOVED***

// coverage is used by NewCoverage which is used as a convenient way for
// creating Coverage implementations for partially defined data. Very often a
// package will only need to define a subset of slices. coverage provides a
// convenient way to do this. Moreover, packages using NewCoverage, instead of
// their own implementation, will not break if later new slice types are added.
type coverage struct ***REMOVED***
	tags    func() []Tag
	bases   func() []Base
	scripts func() []Script
	regions func() []Region
***REMOVED***

func (s *coverage) Tags() []Tag ***REMOVED***
	if s.tags == nil ***REMOVED***
		return nil
	***REMOVED***
	return s.tags()
***REMOVED***

// bases implements sort.Interface and is used to sort base languages.
type bases []Base

func (b bases) Len() int ***REMOVED***
	return len(b)
***REMOVED***

func (b bases) Swap(i, j int) ***REMOVED***
	b[i], b[j] = b[j], b[i]
***REMOVED***

func (b bases) Less(i, j int) bool ***REMOVED***
	return b[i].langID < b[j].langID
***REMOVED***

// BaseLanguages returns the result from calling s.bases if it is specified or
// otherwise derives the set of supported base languages from tags.
func (s *coverage) BaseLanguages() []Base ***REMOVED***
	if s.bases == nil ***REMOVED***
		tags := s.Tags()
		if len(tags) == 0 ***REMOVED***
			return nil
		***REMOVED***
		a := make([]Base, len(tags))
		for i, t := range tags ***REMOVED***
			a[i] = Base***REMOVED***language.Language(t.lang())***REMOVED***
		***REMOVED***
		sort.Sort(bases(a))
		k := 0
		for i := 1; i < len(a); i++ ***REMOVED***
			if a[k] != a[i] ***REMOVED***
				k++
				a[k] = a[i]
			***REMOVED***
		***REMOVED***
		return a[:k+1]
	***REMOVED***
	return s.bases()
***REMOVED***

func (s *coverage) Scripts() []Script ***REMOVED***
	if s.scripts == nil ***REMOVED***
		return nil
	***REMOVED***
	return s.scripts()
***REMOVED***

func (s *coverage) Regions() []Region ***REMOVED***
	if s.regions == nil ***REMOVED***
		return nil
	***REMOVED***
	return s.regions()
***REMOVED***

// NewCoverage returns a Coverage for the given lists. It is typically used by
// packages providing internationalization services to define their level of
// coverage. A list may be of type []T or func() []T, where T is either Tag,
// Base, Script or Region. The returned Coverage derives the value for Bases
// from Tags if no func or slice for []Base is specified. For other unspecified
// types the returned Coverage will return nil for the respective methods.
func NewCoverage(list ...interface***REMOVED******REMOVED***) Coverage ***REMOVED***
	s := &coverage***REMOVED******REMOVED***
	for _, x := range list ***REMOVED***
		switch v := x.(type) ***REMOVED***
		case func() []Base:
			s.bases = v
		case func() []Script:
			s.scripts = v
		case func() []Region:
			s.regions = v
		case func() []Tag:
			s.tags = v
		case []Base:
			s.bases = func() []Base ***REMOVED*** return v ***REMOVED***
		case []Script:
			s.scripts = func() []Script ***REMOVED*** return v ***REMOVED***
		case []Region:
			s.regions = func() []Region ***REMOVED*** return v ***REMOVED***
		case []Tag:
			s.tags = func() []Tag ***REMOVED*** return v ***REMOVED***
		default:
			panic(fmt.Sprintf("language: unsupported set type %T", v))
		***REMOVED***
	***REMOVED***
	return s
***REMOVED***
