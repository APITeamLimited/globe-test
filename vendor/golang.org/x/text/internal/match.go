// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package internal

// This file contains matchers that implement CLDR inheritance.
//
//     See https://unicode.org/reports/tr35/#Locale_Inheritance.
//
// Some of the inheritance described in this document is already handled by
// the cldr package.

import (
	"golang.org/x/text/language"
)

// TODO: consider if (some of the) matching algorithm needs to be public after
// getting some feel about what is generic and what is specific.

// NewInheritanceMatcher returns a matcher that matches based on the inheritance
// chain.
//
// The matcher uses canonicalization and the parent relationship to find a
// match. The resulting match will always be either Und or a language with the
// same language and script as the requested language. It will not match
// languages for which there is understood to be mutual or one-directional
// intelligibility.
//
// A Match will indicate an Exact match if the language matches after
// canonicalization and High if the matched tag is a parent.
func NewInheritanceMatcher(t []language.Tag) *InheritanceMatcher ***REMOVED***
	tags := &InheritanceMatcher***REMOVED***make(map[language.Tag]int)***REMOVED***
	for i, tag := range t ***REMOVED***
		ct, err := language.All.Canonicalize(tag)
		if err != nil ***REMOVED***
			ct = tag
		***REMOVED***
		tags.index[ct] = i
	***REMOVED***
	return tags
***REMOVED***

type InheritanceMatcher struct ***REMOVED***
	index map[language.Tag]int
***REMOVED***

func (m InheritanceMatcher) Match(want ...language.Tag) (language.Tag, int, language.Confidence) ***REMOVED***
	for _, t := range want ***REMOVED***
		ct, err := language.All.Canonicalize(t)
		if err != nil ***REMOVED***
			ct = t
		***REMOVED***
		conf := language.Exact
		for ***REMOVED***
			if index, ok := m.index[ct]; ok ***REMOVED***
				return ct, index, conf
			***REMOVED***
			if ct == language.Und ***REMOVED***
				break
			***REMOVED***
			ct = ct.Parent()
			conf = language.High
		***REMOVED***
	***REMOVED***
	return language.Und, 0, language.No
***REMOVED***
