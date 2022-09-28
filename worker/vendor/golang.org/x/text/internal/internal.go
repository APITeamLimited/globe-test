// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package internal contains non-exported functionality that are used by
// packages in the text repository.
package internal // import "golang.org/x/text/internal"

import (
	"sort"

	"golang.org/x/text/language"
)

// SortTags sorts tags in place.
func SortTags(tags []language.Tag) ***REMOVED***
	sort.Sort(sorter(tags))
***REMOVED***

type sorter []language.Tag

func (s sorter) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sorter) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func (s sorter) Less(i, j int) bool ***REMOVED***
	return s[i].String() < s[j].String()
***REMOVED***

// UniqueTags sorts and filters duplicate tags in place and returns a slice with
// only unique tags.
func UniqueTags(tags []language.Tag) []language.Tag ***REMOVED***
	if len(tags) <= 1 ***REMOVED***
		return tags
	***REMOVED***
	SortTags(tags)
	k := 0
	for i := 1; i < len(tags); i++ ***REMOVED***
		if tags[k].String() < tags[i].String() ***REMOVED***
			k++
			tags[k] = tags[i]
		***REMOVED***
	***REMOVED***
	return tags[:k+1]
***REMOVED***
