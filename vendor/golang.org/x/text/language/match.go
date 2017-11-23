// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package language

import "errors"

// A MatchOption configures a Matcher.
type MatchOption func(*matcher)

// PreferSameScript will, in the absence of a match, result in the first
// preferred tag with the same script as a supported tag to match this supported
// tag. The default is currently true, but this may change in the future.
func PreferSameScript(preferSame bool) MatchOption ***REMOVED***
	return func(m *matcher) ***REMOVED*** m.preferSameScript = preferSame ***REMOVED***
***REMOVED***

// TODO(v1.0.0): consider making Matcher a concrete type, instead of interface.
// There doesn't seem to be too much need for multiple types.
// Making it a concrete type allows MatchStrings to be a method, which will
// improve its discoverability.

// MatchStrings parses and matches the given strings until one of them matches
// the language in the Matcher. A string may be an Accept-Language header as
// handled by ParseAcceptLanguage. The default language is returned if no
// other language matched.
func MatchStrings(m Matcher, lang ...string) (tag Tag, index int) ***REMOVED***
	for _, accept := range lang ***REMOVED***
		desired, _, err := ParseAcceptLanguage(accept)
		if err != nil ***REMOVED***
			continue
		***REMOVED***
		if tag, index, conf := m.Match(desired...); conf != No ***REMOVED***
			return tag, index
		***REMOVED***
	***REMOVED***
	tag, index, _ = m.Match()
	return
***REMOVED***

// Matcher is the interface that wraps the Match method.
//
// Match returns the best match for any of the given tags, along with
// a unique index associated with the returned tag and a confidence
// score.
type Matcher interface ***REMOVED***
	Match(t ...Tag) (tag Tag, index int, c Confidence)
***REMOVED***

// Comprehends reports the confidence score for a speaker of a given language
// to being able to comprehend the written form of an alternative language.
func Comprehends(speaker, alternative Tag) Confidence ***REMOVED***
	_, _, c := NewMatcher([]Tag***REMOVED***alternative***REMOVED***).Match(speaker)
	return c
***REMOVED***

// NewMatcher returns a Matcher that matches an ordered list of preferred tags
// against a list of supported tags based on written intelligibility, closeness
// of dialect, equivalence of subtags and various other rules. It is initialized
// with the list of supported tags. The first element is used as the default
// value in case no match is found.
//
// Its Match method matches the first of the given Tags to reach a certain
// confidence threshold. The tags passed to Match should therefore be specified
// in order of preference. Extensions are ignored for matching.
//
// The index returned by the Match method corresponds to the index of the
// matched tag in t, but is augmented with the Unicode extension ('u')of the
// corresponding preferred tag. This allows user locale options to be passed
// transparently.
func NewMatcher(t []Tag, options ...MatchOption) Matcher ***REMOVED***
	return newMatcher(t, options)
***REMOVED***

func (m *matcher) Match(want ...Tag) (t Tag, index int, c Confidence) ***REMOVED***
	match, w, c := m.getBest(want...)
	if match != nil ***REMOVED***
		t, index = match.tag, match.index
	***REMOVED*** else ***REMOVED***
		// TODO: this should be an option
		t = m.default_.tag
		if m.preferSameScript ***REMOVED***
		outer:
			for _, w := range want ***REMOVED***
				script, _ := w.Script()
				if script.scriptID == 0 ***REMOVED***
					// Don't do anything if there is no script, such as with
					// private subtags.
					continue
				***REMOVED***
				for i, h := range m.supported ***REMOVED***
					if script.scriptID == h.maxScript ***REMOVED***
						t, index = h.tag, i
						break outer
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// TODO: select first language tag based on script.
	***REMOVED***
	if w.region != 0 && t.region != 0 && t.region.contains(w.region) ***REMOVED***
		t, _ = Raw.Compose(t, Region***REMOVED***w.region***REMOVED***)
	***REMOVED***
	// Copy options from the user-provided tag into the result tag. This is hard
	// to do after the fact, so we do it here.
	// TODO: add in alternative variants to -u-va-.
	// TODO: add preferred region to -u-rg-.
	if e := w.Extensions(); len(e) > 0 ***REMOVED***
		t, _ = Raw.Compose(t, e)
	***REMOVED***
	return t, index, c
***REMOVED***

type scriptRegionFlags uint8

const (
	isList = 1 << iota
	scriptInFrom
	regionInFrom
)

func (t *Tag) setUndefinedLang(id langID) ***REMOVED***
	if t.lang == 0 ***REMOVED***
		t.lang = id
	***REMOVED***
***REMOVED***

func (t *Tag) setUndefinedScript(id scriptID) ***REMOVED***
	if t.script == 0 ***REMOVED***
		t.script = id
	***REMOVED***
***REMOVED***

func (t *Tag) setUndefinedRegion(id regionID) ***REMOVED***
	if t.region == 0 || t.region.contains(id) ***REMOVED***
		t.region = id
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
	id.remakeString()
	return id, nil
***REMOVED***

// specializeRegion attempts to specialize a group region.
func specializeRegion(t *Tag) bool ***REMOVED***
	if i := regionInclusion[t.region]; i < nRegionGroups ***REMOVED***
		x := likelyRegionGroup[i]
		if langID(x.lang) == t.lang && scriptID(x.script) == t.script ***REMOVED***
			t.region = regionID(x.region)
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func addTags(t Tag) (Tag, error) ***REMOVED***
	// We leave private use identifiers alone.
	if t.private() ***REMOVED***
		return t, nil
	***REMOVED***
	if t.script != 0 && t.region != 0 ***REMOVED***
		if t.lang != 0 ***REMOVED***
			// already fully specified
			specializeRegion(&t)
			return t, nil
		***REMOVED***
		// Search matches for und-script-region. Note that for these cases
		// region will never be a group so there is no need to check for this.
		list := likelyRegion[t.region : t.region+1]
		if x := list[0]; x.flags&isList != 0 ***REMOVED***
			list = likelyRegionList[x.lang : x.lang+uint16(x.script)]
		***REMOVED***
		for _, x := range list ***REMOVED***
			// Deviating from the spec. See match_test.go for details.
			if scriptID(x.script) == t.script ***REMOVED***
				t.setUndefinedLang(langID(x.lang))
				return t, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if t.lang != 0 ***REMOVED***
		// Search matches for lang-script and lang-region, where lang != und.
		if t.lang < langNoIndexOffset ***REMOVED***
			x := likelyLang[t.lang]
			if x.flags&isList != 0 ***REMOVED***
				list := likelyLangList[x.region : x.region+uint16(x.script)]
				if t.script != 0 ***REMOVED***
					for _, x := range list ***REMOVED***
						if scriptID(x.script) == t.script && x.flags&scriptInFrom != 0 ***REMOVED***
							t.setUndefinedRegion(regionID(x.region))
							return t, nil
						***REMOVED***
					***REMOVED***
				***REMOVED*** else if t.region != 0 ***REMOVED***
					count := 0
					goodScript := true
					tt := t
					for _, x := range list ***REMOVED***
						// We visit all entries for which the script was not
						// defined, including the ones where the region was not
						// defined. This allows for proper disambiguation within
						// regions.
						if x.flags&scriptInFrom == 0 && t.region.contains(regionID(x.region)) ***REMOVED***
							tt.region = regionID(x.region)
							tt.setUndefinedScript(scriptID(x.script))
							goodScript = goodScript && tt.script == scriptID(x.script)
							count++
						***REMOVED***
					***REMOVED***
					if count == 1 ***REMOVED***
						return tt, nil
					***REMOVED***
					// Even if we fail to find a unique Region, we might have
					// an unambiguous script.
					if goodScript ***REMOVED***
						t.script = tt.script
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Search matches for und-script.
		if t.script != 0 ***REMOVED***
			x := likelyScript[t.script]
			if x.region != 0 ***REMOVED***
				t.setUndefinedRegion(regionID(x.region))
				t.setUndefinedLang(langID(x.lang))
				return t, nil
			***REMOVED***
		***REMOVED***
		// Search matches for und-region. If und-script-region exists, it would
		// have been found earlier.
		if t.region != 0 ***REMOVED***
			if i := regionInclusion[t.region]; i < nRegionGroups ***REMOVED***
				x := likelyRegionGroup[i]
				if x.region != 0 ***REMOVED***
					t.setUndefinedLang(langID(x.lang))
					t.setUndefinedScript(scriptID(x.script))
					t.region = regionID(x.region)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				x := likelyRegion[t.region]
				if x.flags&isList != 0 ***REMOVED***
					x = likelyRegionList[x.lang]
				***REMOVED***
				if x.script != 0 && x.flags != scriptInFrom ***REMOVED***
					t.setUndefinedLang(langID(x.lang))
					t.setUndefinedScript(scriptID(x.script))
					return t, nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Search matches for lang.
	if t.lang < langNoIndexOffset ***REMOVED***
		x := likelyLang[t.lang]
		if x.flags&isList != 0 ***REMOVED***
			x = likelyLangList[x.region]
		***REMOVED***
		if x.region != 0 ***REMOVED***
			t.setUndefinedScript(scriptID(x.script))
			t.setUndefinedRegion(regionID(x.region))
		***REMOVED***
		specializeRegion(&t)
		if t.lang == 0 ***REMOVED***
			t.lang = _en // default language
		***REMOVED***
		return t, nil
	***REMOVED***
	return t, ErrMissingLikelyTagsData
***REMOVED***

func (t *Tag) setTagsFrom(id Tag) ***REMOVED***
	t.lang = id.lang
	t.script = id.script
	t.region = id.region
***REMOVED***

// minimize removes the region or script subtags from t such that
// t.addLikelySubtags() == t.minimize().addLikelySubtags().
func (t Tag) minimize() (Tag, error) ***REMOVED***
	t, err := minimizeTags(t)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***
	t.remakeString()
	return t, nil
***REMOVED***

// minimizeTags mimics the behavior of the ICU 51 C implementation.
func minimizeTags(t Tag) (Tag, error) ***REMOVED***
	if t.equalTags(und) ***REMOVED***
		return t, nil
	***REMOVED***
	max, err := addTags(t)
	if err != nil ***REMOVED***
		return t, err
	***REMOVED***
	for _, id := range [...]Tag***REMOVED***
		***REMOVED***lang: t.lang***REMOVED***,
		***REMOVED***lang: t.lang, region: t.region***REMOVED***,
		***REMOVED***lang: t.lang, script: t.script***REMOVED***,
	***REMOVED*** ***REMOVED***
		if x, err := addTags(id); err == nil && max.equalTags(x) ***REMOVED***
			t.setTagsFrom(id)
			break
		***REMOVED***
	***REMOVED***
	return t, nil
***REMOVED***

// Tag Matching
// CLDR defines an algorithm for finding the best match between two sets of language
// tags. The basic algorithm defines how to score a possible match and then find
// the match with the best score
// (see http://www.unicode.org/reports/tr35/#LanguageMatching).
// Using scoring has several disadvantages. The scoring obfuscates the importance of
// the various factors considered, making the algorithm harder to understand. Using
// scoring also requires the full score to be computed for each pair of tags.
//
// We will use a different algorithm which aims to have the following properties:
// - clarity on the precedence of the various selection factors, and
// - improved performance by allowing early termination of a comparison.
//
// Matching algorithm (overview)
// Input:
//   - supported: a set of supported tags
//   - default:   the default tag to return in case there is no match
//   - desired:   list of desired tags, ordered by preference, starting with
//                the most-preferred.
//
// Algorithm:
//   1) Set the best match to the lowest confidence level
//   2) For each tag in "desired":
//     a) For each tag in "supported":
//        1) compute the match between the two tags.
//        2) if the match is better than the previous best match, replace it
//           with the new match. (see next section)
//     b) if the current best match is Exact and pin is true the result will be
//        frozen to the language found thusfar, although better matches may
//        still be found for the same language.
//   3) If the best match so far is below a certain threshold, return "default".
//
// Ranking:
// We use two phases to determine whether one pair of tags are a better match
// than another pair of tags. First, we determine a rough confidence level. If the
// levels are different, the one with the highest confidence wins.
// Second, if the rough confidence levels are identical, we use a set of tie-breaker
// rules.
//
// The confidence level of matching a pair of tags is determined by finding the
// lowest confidence level of any matches of the corresponding subtags (the
// result is deemed as good as its weakest link).
// We define the following levels:
//   Exact    - An exact match of a subtag, before adding likely subtags.
//   MaxExact - An exact match of a subtag, after adding likely subtags.
//              [See Note 2].
//   High     - High level of mutual intelligibility between different subtag
//              variants.
//   Low      - Low level of mutual intelligibility between different subtag
//              variants.
//   No       - No mutual intelligibility.
//
// The following levels can occur for each type of subtag:
//   Base:    Exact, MaxExact, High, Low, No
//   Script:  Exact, MaxExact [see Note 3], Low, No
//   Region:  Exact, MaxExact, High
//   Variant: Exact, High
//   Private: Exact, No
//
// Any result with a confidence level of Low or higher is deemed a possible match.
// Once a desired tag matches any of the supported tags with a level of MaxExact
// or higher, the next desired tag is not considered (see Step 2.b).
// Note that CLDR provides languageMatching data that defines close equivalence
// classes for base languages, scripts and regions.
//
// Tie-breaking
// If we get the same confidence level for two matches, we apply a sequence of
// tie-breaking rules. The first that succeeds defines the result. The rules are
// applied in the following order.
//   1) Original language was defined and was identical.
//   2) Original region was defined and was identical.
//   3) Distance between two maximized regions was the smallest.
//   4) Original script was defined and was identical.
//   5) Distance from want tag to have tag using the parent relation [see Note 5.]
// If there is still no winner after these rules are applied, the first match
// found wins.
//
// Notes:
// [2] In practice, as matching of Exact is done in a separate phase from
//     matching the other levels, we reuse the Exact level to mean MaxExact in
//     the second phase. As a consequence, we only need the levels defined by
//     the Confidence type. The MaxExact confidence level is mapped to High in
//     the public API.
// [3] We do not differentiate between maximized script values that were derived
//     from suppressScript versus most likely tag data. We determined that in
//     ranking the two, one ranks just after the other. Moreover, the two cannot
//     occur concurrently. As a consequence, they are identical for practical
//     purposes.
// [4] In case of deprecated, macro-equivalents and legacy mappings, we assign
//     the MaxExact level to allow iw vs he to still be a closer match than
//     en-AU vs en-US, for example.
// [5] In CLDR a locale inherits fields that are unspecified for this locale
//     from its parent. Therefore, if a locale is a parent of another locale,
//     it is a strong measure for closeness, especially when no other tie
//     breaker rule applies. One could also argue it is inconsistent, for
//     example, when pt-AO matches pt (which CLDR equates with pt-BR), even
//     though its parent is pt-PT according to the inheritance rules.
//
// Implementation Details:
// There are several performance considerations worth pointing out. Most notably,
// we preprocess as much as possible (within reason) at the time of creation of a
// matcher. This includes:
//   - creating a per-language map, which includes data for the raw base language
//     and its canonicalized variant (if applicable),
//   - expanding entries for the equivalence classes defined in CLDR's
//     languageMatch data.
// The per-language map ensures that typically only a very small number of tags
// need to be considered. The pre-expansion of canonicalized subtags and
// equivalence classes reduces the amount of map lookups that need to be done at
// runtime.

// matcher keeps a set of supported language tags, indexed by language.
type matcher struct ***REMOVED***
	default_         *haveTag
	supported        []*haveTag
	index            map[langID]*matchHeader
	passSettings     bool
	preferSameScript bool
***REMOVED***

// matchHeader has the lists of tags for exact matches and matches based on
// maximized and canonicalized tags for a given language.
type matchHeader struct ***REMOVED***
	haveTags []*haveTag
	original bool
***REMOVED***

// haveTag holds a supported Tag and its maximized script and region. The maximized
// or canonicalized language is not stored as it is not needed during matching.
type haveTag struct ***REMOVED***
	tag Tag

	// index of this tag in the original list of supported tags.
	index int

	// conf is the maximum confidence that can result from matching this haveTag.
	// When conf < Exact this means it was inserted after applying a CLDR equivalence rule.
	conf Confidence

	// Maximized region and script.
	maxRegion regionID
	maxScript scriptID

	// altScript may be checked as an alternative match to maxScript. If altScript
	// matches, the confidence level for this match is Low. Theoretically there
	// could be multiple alternative scripts. This does not occur in practice.
	altScript scriptID

	// nextMax is the index of the next haveTag with the same maximized tags.
	nextMax uint16
***REMOVED***

func makeHaveTag(tag Tag, index int) (haveTag, langID) ***REMOVED***
	max := tag
	if tag.lang != 0 || tag.region != 0 || tag.script != 0 ***REMOVED***
		max, _ = max.canonicalize(All)
		max, _ = addTags(max)
		max.remakeString()
	***REMOVED***
	return haveTag***REMOVED***tag, index, Exact, max.region, max.script, altScript(max.lang, max.script), 0***REMOVED***, max.lang
***REMOVED***

// altScript returns an alternative script that may match the given script with
// a low confidence.  At the moment, the langMatch data allows for at most one
// script to map to another and we rely on this to keep the code simple.
func altScript(l langID, s scriptID) scriptID ***REMOVED***
	for _, alt := range matchScript ***REMOVED***
		// TODO: also match cases where language is not the same.
		if (langID(alt.wantLang) == l || langID(alt.haveLang) == l) &&
			scriptID(alt.haveScript) == s ***REMOVED***
			return scriptID(alt.wantScript)
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

// addIfNew adds a haveTag to the list of tags only if it is a unique tag.
// Tags that have the same maximized values are linked by index.
func (h *matchHeader) addIfNew(n haveTag, exact bool) ***REMOVED***
	h.original = h.original || exact
	// Don't add new exact matches.
	for _, v := range h.haveTags ***REMOVED***
		if v.tag.equalsRest(n.tag) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	// Allow duplicate maximized tags, but create a linked list to allow quickly
	// comparing the equivalents and bail out.
	for i, v := range h.haveTags ***REMOVED***
		if v.maxScript == n.maxScript &&
			v.maxRegion == n.maxRegion &&
			v.tag.variantOrPrivateTagStr() == n.tag.variantOrPrivateTagStr() ***REMOVED***
			for h.haveTags[i].nextMax != 0 ***REMOVED***
				i = int(h.haveTags[i].nextMax)
			***REMOVED***
			h.haveTags[i].nextMax = uint16(len(h.haveTags))
			break
		***REMOVED***
	***REMOVED***
	h.haveTags = append(h.haveTags, &n)
***REMOVED***

// header returns the matchHeader for the given language. It creates one if
// it doesn't already exist.
func (m *matcher) header(l langID) *matchHeader ***REMOVED***
	if h := m.index[l]; h != nil ***REMOVED***
		return h
	***REMOVED***
	h := &matchHeader***REMOVED******REMOVED***
	m.index[l] = h
	return h
***REMOVED***

func toConf(d uint8) Confidence ***REMOVED***
	if d <= 10 ***REMOVED***
		return High
	***REMOVED***
	if d < 30 ***REMOVED***
		return Low
	***REMOVED***
	return No
***REMOVED***

// newMatcher builds an index for the given supported tags and returns it as
// a matcher. It also expands the index by considering various equivalence classes
// for a given tag.
func newMatcher(supported []Tag, options []MatchOption) *matcher ***REMOVED***
	m := &matcher***REMOVED***
		index:            make(map[langID]*matchHeader),
		preferSameScript: true,
	***REMOVED***
	for _, o := range options ***REMOVED***
		o(m)
	***REMOVED***
	if len(supported) == 0 ***REMOVED***
		m.default_ = &haveTag***REMOVED******REMOVED***
		return m
	***REMOVED***
	// Add supported languages to the index. Add exact matches first to give
	// them precedence.
	for i, tag := range supported ***REMOVED***
		pair, _ := makeHaveTag(tag, i)
		m.header(tag.lang).addIfNew(pair, true)
		m.supported = append(m.supported, &pair)
	***REMOVED***
	m.default_ = m.header(supported[0].lang).haveTags[0]
	// Keep these in two different loops to support the case that two equivalent
	// languages are distinguished, such as iw and he.
	for i, tag := range supported ***REMOVED***
		pair, max := makeHaveTag(tag, i)
		if max != tag.lang ***REMOVED***
			m.header(max).addIfNew(pair, true)
		***REMOVED***
	***REMOVED***

	// update is used to add indexes in the map for equivalent languages.
	// update will only add entries to original indexes, thus not computing any
	// transitive relations.
	update := func(want, have uint16, conf Confidence) ***REMOVED***
		if hh := m.index[langID(have)]; hh != nil ***REMOVED***
			if !hh.original ***REMOVED***
				return
			***REMOVED***
			hw := m.header(langID(want))
			for _, ht := range hh.haveTags ***REMOVED***
				v := *ht
				if conf < v.conf ***REMOVED***
					v.conf = conf
				***REMOVED***
				v.nextMax = 0 // this value needs to be recomputed
				if v.altScript != 0 ***REMOVED***
					v.altScript = altScript(langID(want), v.maxScript)
				***REMOVED***
				hw.addIfNew(v, conf == Exact && hh.original)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Add entries for languages with mutual intelligibility as defined by CLDR's
	// languageMatch data.
	for _, ml := range matchLang ***REMOVED***
		update(ml.want, ml.have, toConf(ml.distance))
		if !ml.oneway ***REMOVED***
			update(ml.have, ml.want, toConf(ml.distance))
		***REMOVED***
	***REMOVED***

	// Add entries for possible canonicalizations. This is an optimization to
	// ensure that only one map lookup needs to be done at runtime per desired tag.
	// First we match deprecated equivalents. If they are perfect equivalents
	// (their canonicalization simply substitutes a different language code, but
	// nothing else), the match confidence is Exact, otherwise it is High.
	for i, lm := range langAliasMap ***REMOVED***
		// If deprecated codes match and there is no fiddling with the script or
		// or region, we consider it an exact match.
		conf := Exact
		if langAliasTypes[i] != langMacro ***REMOVED***
			if !isExactEquivalent(langID(lm.from)) ***REMOVED***
				conf = High
			***REMOVED***
			update(lm.to, lm.from, conf)
		***REMOVED***
		update(lm.from, lm.to, conf)
	***REMOVED***
	return m
***REMOVED***

// getBest gets the best matching tag in m for any of the given tags, taking into
// account the order of preference of the given tags.
func (m *matcher) getBest(want ...Tag) (got *haveTag, orig Tag, c Confidence) ***REMOVED***
	best := bestMatch***REMOVED******REMOVED***
	for i, w := range want ***REMOVED***
		var max Tag
		// Check for exact match first.
		h := m.index[w.lang]
		if w.lang != 0 ***REMOVED***
			if h == nil ***REMOVED***
				continue
			***REMOVED***
			// Base language is defined.
			max, _ = w.canonicalize(Legacy | Deprecated | Macro)
			// A region that is added through canonicalization is stronger than
			// a maximized region: set it in the original (e.g. mo -> ro-MD).
			if w.region != max.region ***REMOVED***
				w.region = max.region
			***REMOVED***
			// TODO: should we do the same for scripts?
			// See test case: en, sr, nl ; sh ; sr
			max, _ = addTags(max)
		***REMOVED*** else ***REMOVED***
			// Base language is not defined.
			if h != nil ***REMOVED***
				for i := range h.haveTags ***REMOVED***
					have := h.haveTags[i]
					if have.tag.equalsRest(w) ***REMOVED***
						return have, w, Exact
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if w.script == 0 && w.region == 0 ***REMOVED***
				// We skip all tags matching und for approximate matching, including
				// private tags.
				continue
			***REMOVED***
			max, _ = addTags(w)
			if h = m.index[max.lang]; h == nil ***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		pin := true
		for _, t := range want[i+1:] ***REMOVED***
			if w.lang == t.lang ***REMOVED***
				pin = false
				break
			***REMOVED***
		***REMOVED***
		// Check for match based on maximized tag.
		for i := range h.haveTags ***REMOVED***
			have := h.haveTags[i]
			best.update(have, w, max.script, max.region, pin)
			if best.conf == Exact ***REMOVED***
				for have.nextMax != 0 ***REMOVED***
					have = h.haveTags[have.nextMax]
					best.update(have, w, max.script, max.region, pin)
				***REMOVED***
				return best.have, best.want, best.conf
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if best.conf <= No ***REMOVED***
		if len(want) != 0 ***REMOVED***
			return nil, want[0], No
		***REMOVED***
		return nil, Tag***REMOVED******REMOVED***, No
	***REMOVED***
	return best.have, best.want, best.conf
***REMOVED***

// bestMatch accumulates the best match so far.
type bestMatch struct ***REMOVED***
	have            *haveTag
	want            Tag
	conf            Confidence
	pinnedRegion    regionID
	pinLanguage     bool
	sameRegionGroup bool
	// Cached results from applying tie-breaking rules.
	origLang     bool
	origReg      bool
	paradigmReg  bool
	regGroupDist uint8
	origScript   bool
***REMOVED***

// update updates the existing best match if the new pair is considered to be a
// better match. To determine if the given pair is a better match, it first
// computes the rough confidence level. If this surpasses the current match, it
// will replace it and update the tie-breaker rule cache. If there is a tie, it
// proceeds with applying a series of tie-breaker rules. If there is no
// conclusive winner after applying the tie-breaker rules, it leaves the current
// match as the preferred match.
//
// If pin is true and have and tag are a strong match, it will henceforth only
// consider matches for this language. This corresponds to the nothing that most
// users have a strong preference for the first defined language. A user can
// still prefer a second language over a dialect of the preferred language by
// explicitly specifying dialects, e.g. "en, nl, en-GB". In this case pin should
// be false.
func (m *bestMatch) update(have *haveTag, tag Tag, maxScript scriptID, maxRegion regionID, pin bool) ***REMOVED***
	// Bail if the maximum attainable confidence is below that of the current best match.
	c := have.conf
	if c < m.conf ***REMOVED***
		return
	***REMOVED***
	// Don't change the language once we already have found an exact match.
	if m.pinLanguage && tag.lang != m.want.lang ***REMOVED***
		return
	***REMOVED***
	// Pin the region group if we are comparing tags for the same language.
	if tag.lang == m.want.lang && m.sameRegionGroup ***REMOVED***
		_, sameGroup := regionGroupDist(m.pinnedRegion, have.maxRegion, have.maxScript, m.want.lang)
		if !sameGroup ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	if c == Exact && have.maxScript == maxScript ***REMOVED***
		// If there is another language and then another entry of this language,
		// don't pin anything, otherwise pin the language.
		m.pinLanguage = pin
	***REMOVED***
	if have.tag.equalsRest(tag) ***REMOVED***
	***REMOVED*** else if have.maxScript != maxScript ***REMOVED***
		// There is usually very little comprehension between different scripts.
		// In a few cases there may still be Low comprehension. This possibility
		// is pre-computed and stored in have.altScript.
		if Low < m.conf || have.altScript != maxScript ***REMOVED***
			return
		***REMOVED***
		c = Low
	***REMOVED*** else if have.maxRegion != maxRegion ***REMOVED***
		if High < c ***REMOVED***
			// There is usually a small difference between languages across regions.
			c = High
		***REMOVED***
	***REMOVED***

	// We store the results of the computations of the tie-breaker rules along
	// with the best match. There is no need to do the checks once we determine
	// we have a winner, but we do still need to do the tie-breaker computations.
	// We use "beaten" to keep track if we still need to do the checks.
	beaten := false // true if the new pair defeats the current one.
	if c != m.conf ***REMOVED***
		if c < m.conf ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	// Tie-breaker rules:
	// We prefer if the pre-maximized language was specified and identical.
	origLang := have.tag.lang == tag.lang && tag.lang != 0
	if !beaten && m.origLang != origLang ***REMOVED***
		if m.origLang ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	// We prefer if the pre-maximized region was specified and identical.
	origReg := have.tag.region == tag.region && tag.region != 0
	if !beaten && m.origReg != origReg ***REMOVED***
		if m.origReg ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	regGroupDist, sameGroup := regionGroupDist(have.maxRegion, maxRegion, maxScript, tag.lang)
	if !beaten && m.regGroupDist != regGroupDist ***REMOVED***
		if regGroupDist > m.regGroupDist ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	paradigmReg := isParadigmLocale(tag.lang, have.maxRegion)
	if !beaten && m.paradigmReg != paradigmReg ***REMOVED***
		if !paradigmReg ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	// Next we prefer if the pre-maximized script was specified and identical.
	origScript := have.tag.script == tag.script && tag.script != 0
	if !beaten && m.origScript != origScript ***REMOVED***
		if m.origScript ***REMOVED***
			return
		***REMOVED***
		beaten = true
	***REMOVED***

	// Update m to the newly found best match.
	if beaten ***REMOVED***
		m.have = have
		m.want = tag
		m.conf = c
		m.pinnedRegion = maxRegion
		m.sameRegionGroup = sameGroup
		m.origLang = origLang
		m.origReg = origReg
		m.paradigmReg = paradigmReg
		m.origScript = origScript
		m.regGroupDist = regGroupDist
	***REMOVED***
***REMOVED***

func isParadigmLocale(lang langID, r regionID) bool ***REMOVED***
	for _, e := range paradigmLocales ***REMOVED***
		if langID(e[0]) == lang && (r == regionID(e[1]) || r == regionID(e[2])) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// regionGroupDist computes the distance between two regions based on their
// CLDR grouping.
func regionGroupDist(a, b regionID, script scriptID, lang langID) (dist uint8, same bool) ***REMOVED***
	const defaultDistance = 4

	aGroup := uint(regionToGroups[a]) << 1
	bGroup := uint(regionToGroups[b]) << 1
	for _, ri := range matchRegion ***REMOVED***
		if langID(ri.lang) == lang && (ri.script == 0 || scriptID(ri.script) == script) ***REMOVED***
			group := uint(1 << (ri.group &^ 0x80))
			if 0x80&ri.group == 0 ***REMOVED***
				if aGroup&bGroup&group != 0 ***REMOVED*** // Both regions are in the group.
					return ri.distance, ri.distance == defaultDistance
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if (aGroup|bGroup)&group == 0 ***REMOVED*** // Both regions are not in the group.
					return ri.distance, ri.distance == defaultDistance
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return defaultDistance, true
***REMOVED***

func (t Tag) variants() string ***REMOVED***
	if t.pVariant == 0 ***REMOVED***
		return ""
	***REMOVED***
	return t.str[t.pVariant:t.pExt]
***REMOVED***

// variantOrPrivateTagStr returns variants or private use tags.
func (t Tag) variantOrPrivateTagStr() string ***REMOVED***
	if t.pExt > 0 ***REMOVED***
		return t.str[t.pVariant:t.pExt]
	***REMOVED***
	return t.str[t.pVariant:]
***REMOVED***

// equalsRest compares everything except the language.
func (a Tag) equalsRest(b Tag) bool ***REMOVED***
	// TODO: don't include extensions in this comparison. To do this efficiently,
	// though, we should handle private tags separately.
	return a.script == b.script && a.region == b.region && a.variantOrPrivateTagStr() == b.variantOrPrivateTagStr()
***REMOVED***

// isExactEquivalent returns true if canonicalizing the language will not alter
// the script or region of a tag.
func isExactEquivalent(l langID) bool ***REMOVED***
	for _, o := range notEquivalent ***REMOVED***
		if o == l ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

var notEquivalent []langID

func init() ***REMOVED***
	// Create a list of all languages for which canonicalization may alter the
	// script or region.
	for _, lm := range langAliasMap ***REMOVED***
		tag := Tag***REMOVED***lang: langID(lm.from)***REMOVED***
		if tag, _ = tag.canonicalize(All); tag.script != 0 || tag.region != 0 ***REMOVED***
			notEquivalent = append(notEquivalent, langID(lm.from))
		***REMOVED***
	***REMOVED***
	// Maximize undefined regions of paradigm locales.
	for i, v := range paradigmLocales ***REMOVED***
		max, _ := addTags(Tag***REMOVED***lang: langID(v[0])***REMOVED***)
		if v[1] == 0 ***REMOVED***
			paradigmLocales[i][1] = uint16(max.region)
		***REMOVED***
		if v[2] == 0 ***REMOVED***
			paradigmLocales[i][2] = uint16(max.region)
		***REMOVED***
	***REMOVED***
***REMOVED***
