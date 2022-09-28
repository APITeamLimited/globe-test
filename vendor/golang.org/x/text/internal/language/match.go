// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import "errors"

type scriptRegionFlags uint8

const (
	isList = 1 << iota
	scriptInFrom
	regionInFrom
)

func (t *Tag) setUndefinedLang(id Language) ***REMOVED***
	if t.LangID == 0 ***REMOVED***
		t.LangID = id
	***REMOVED***
***REMOVED***

func (t *Tag) setUndefinedScript(id Script) ***REMOVED***
	if t.ScriptID == 0 ***REMOVED***
		t.ScriptID = id
	***REMOVED***
***REMOVED***

func (t *Tag) setUndefinedRegion(id Region) ***REMOVED***
	if t.RegionID == 0 || t.RegionID.Contains(id) ***REMOVED***
		t.RegionID = id
	***REMOVED***
***REMOVED***

// ErrMissingLikelyTagsData indicates no information was available
// to compute likely values of missing tags.
var ErrMissingLikelyTagsData = errors.New("missing likely tags data")

// addLikelySubtags sets subtags to their most likely value, given the locale.
// In most cases this means setting fields for unknown values, but in some
// cases it may alter a value.  It returns an ErrMissingLikelyTagsData error
// if the given locale cannot be expanded.
func (t Tag) addLikelySubtags() (Tag, error) ***REMOVED***
	id, err := addTags(t)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED*** else if id.equalTags(t) ***REMOVED***
		return t, nil
	***REMOVED***
	id.RemakeString()
	return id, nil
***REMOVED***

// specializeRegion attempts to specialize a group region.
func specializeRegion(t *Tag) bool ***REMOVED***
	if i := regionInclusion[t.RegionID]; i < nRegionGroups ***REMOVED***
		x := likelyRegionGroup[i]
		if Language(x.lang) == t.LangID && Script(x.script) == t.ScriptID ***REMOVED***
			t.RegionID = Region(x.region)
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// Maximize returns a new tag with missing tags filled in.
func (t Tag) Maximize() (Tag, error) ***REMOVED***
	return addTags(t)
***REMOVED***

func addTags(t Tag) (Tag, error) ***REMOVED***
	// We leave private use identifiers alone.
	if t.IsPrivateUse() ***REMOVED***
		return t, nil
	***REMOVED***
	if t.ScriptID != 0 && t.RegionID != 0 ***REMOVED***
		if t.LangID != 0 ***REMOVED***
			// already fully specified
			specializeRegion(&t)
			return t, nil
		***REMOVED***
		// Search matches for und-script-region. Note that for these cases
		// region will never be a group so there is no need to check for this.
		list := likelyRegion[t.RegionID : t.RegionID+1]
		if x := list[0]; x.flags&isList != 0 ***REMOVED***
			list = likelyRegionList[x.lang : x.lang+uint16(x.script)]
		***REMOVED***
		for _, x := range list ***REMOVED***
			// Deviating from the spec. See match_test.go for details.
			if Script(x.script) == t.ScriptID ***REMOVED***
				t.setUndefinedLang(Language(x.lang))
				return t, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if t.LangID != 0 ***REMOVED***
		// Search matches for lang-script and lang-region, where lang != und.
		if t.LangID < langNoIndexOffset ***REMOVED***
			x := likelyLang[t.LangID]
			if x.flags&isList != 0 ***REMOVED***
				list := likelyLangList[x.region : x.region+uint16(x.script)]
				if t.ScriptID != 0 ***REMOVED***
					for _, x := range list ***REMOVED***
						if Script(x.script) == t.ScriptID && x.flags&scriptInFrom != 0 ***REMOVED***
							t.setUndefinedRegion(Region(x.region))
							return t, nil
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if t.RegionID != 0 ***REMOVED***
					count := 0
					goodScript := true
					tt := t
					for _, x := range list ***REMOVED***
						// We visit all entries for which the script was not
						// defined, including the ones where the region was not
						// defined. This allows for proper disambiguation within
						// regions.
						if x.flags&scriptInFrom == 0 && t.RegionID.Contains(Region(x.region)) ***REMOVED***
							tt.RegionID = Region(x.region)
							tt.setUndefinedScript(Script(x.script))
							goodScript = goodScript && tt.ScriptID == Script(x.script)
							count++
						***REMOVED***
					***REMOVED***
					if count == 1 ***REMOVED***
						return tt, nil
					***REMOVED***
					// Even if we fail to find a unique Region, we might have
					// an unambiguous script.
					if goodScript ***REMOVED***
						t.ScriptID = tt.ScriptID
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Search matches for und-script.
		if t.ScriptID != 0 ***REMOVED***
			x := likelyScript[t.ScriptID]
			if x.region != 0 ***REMOVED***
				t.setUndefinedRegion(Region(x.region))
				t.setUndefinedLang(Language(x.lang))
				return t, nil
			***REMOVED***
		***REMOVED***
		// Search matches for und-region. If und-script-region exists, it would
		// have been found earlier.
		if t.RegionID != 0 ***REMOVED***
			if i := regionInclusion[t.RegionID]; i < nRegionGroups ***REMOVED***
				x := likelyRegionGroup[i]
				if x.region != 0 ***REMOVED***
					t.setUndefinedLang(Language(x.lang))
					t.setUndefinedScript(Script(x.script))
					t.RegionID = Region(x.region)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				x := likelyRegion[t.RegionID]
				if x.flags&isList != 0 ***REMOVED***
					x = likelyRegionList[x.lang]
				***REMOVED***
				if x.script != 0 && x.flags != scriptInFrom ***REMOVED***
					t.setUndefinedLang(Language(x.lang))
					t.setUndefinedScript(Script(x.script))
					return t, nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Search matches for lang.
	if t.LangID < langNoIndexOffset ***REMOVED***
		x := likelyLang[t.LangID]
		if x.flags&isList != 0 ***REMOVED***
			x = likelyLangList[x.region]
		***REMOVED***
		if x.region != 0 ***REMOVED***
			t.setUndefinedScript(Script(x.script))
			t.setUndefinedRegion(Region(x.region))
		***REMOVED***
		specializeRegion(&t)
		if t.LangID == 0 ***REMOVED***
			t.LangID = _en // default language
		***REMOVED***
		return t, nil
	***REMOVED***
	return t, ErrMissingLikelyTagsData
***REMOVED***

func (t *Tag) setTagsFrom(id Tag) ***REMOVED***
	t.LangID = id.LangID
	t.ScriptID = id.ScriptID
	t.RegionID = id.RegionID
***REMOVED***

// minimize removes the region or script subtags from t such that
// t.addLikelySubtags() == t.minimize().addLikelySubtags().
func (t Tag) minimize() (Tag, error) ***REMOVED***
	t, err := minimizeTags(t)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***
	t.RemakeString()
	return t, nil
***REMOVED***

// minimizeTags mimics the behavior of the ICU 51 C implementation.
func minimizeTags(t Tag) (Tag, error) ***REMOVED***
	if t.equalTags(Und) ***REMOVED***
		return t, nil
	***REMOVED***
	max, err := addTags(t)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***
	for _, id := range [...]Tag***REMOVED***
		***REMOVED***LangID: t.LangID***REMOVED***,
		***REMOVED***LangID: t.LangID, RegionID: t.RegionID***REMOVED***,
		***REMOVED***LangID: t.LangID, ScriptID: t.ScriptID***REMOVED***,
	***REMOVED*** ***REMOVED***
		if x, err := addTags(id); err == nil && max.equalTags(x) ***REMOVED***
			t.setTagsFrom(id)
			break
		***REMOVED***
	***REMOVED***
	return t, nil
***REMOVED***
