// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package compact defines a compact representation of language tags.
//
// Common language tags (at least all for which locale information is defined
// in CLDR) are assigned a unique index. Each Tag is associated with such an
// ID for selecting language-related resources (such as translations) as well
// as one for selecting regional defaults (currency, number formatting, etc.)
//
// It may want to export this functionality at some point, but at this point
// this is only available for use within x/text.
package compact // import "golang.org/x/text/internal/language/compact"

import (
	"sort"
	"strings"

	"golang.org/x/text/internal/language"
)

// ID is an integer identifying a single tag.
type ID uint16

func getCoreIndex(t language.Tag) (id ID, ok bool) ***REMOVED***
	cci, ok := language.GetCompactCore(t)
	if !ok ***REMOVED***
		return 0, false
	***REMOVED***
	i := sort.Search(len(coreTags), func(i int) bool ***REMOVED***
		return cci <= coreTags[i]
	***REMOVED***)
	if i == len(coreTags) || coreTags[i] != cci ***REMOVED***
		return 0, false
	***REMOVED***
	return ID(i), true
***REMOVED***

// Parent returns the ID of the parent or the root ID if id is already the root.
func (id ID) Parent() ID ***REMOVED***
	return parents[id]
***REMOVED***

// Tag converts id to an internal language Tag.
func (id ID) Tag() language.Tag ***REMOVED***
	if int(id) >= len(coreTags) ***REMOVED***
		return specialTags[int(id)-len(coreTags)]
	***REMOVED***
	return coreTags[id].Tag()
***REMOVED***

var specialTags []language.Tag

func init() ***REMOVED***
	tags := strings.Split(specialTagsStr, " ")
	specialTags = make([]language.Tag, len(tags))
	for i, t := range tags ***REMOVED***
		specialTags[i] = language.MustParse(t)
	***REMOVED***
***REMOVED***
