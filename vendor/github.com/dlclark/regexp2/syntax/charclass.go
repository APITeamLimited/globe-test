package syntax

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"unicode"
	"unicode/utf8"
)

// CharSet combines start-end rune ranges and unicode categories representing a set of characters
type CharSet struct ***REMOVED***
	ranges     []singleRange
	categories []category
	sub        *CharSet //optional subtractor
	negate     bool
	anything   bool
***REMOVED***

type category struct ***REMOVED***
	negate bool
	cat    string
***REMOVED***

type singleRange struct ***REMOVED***
	first rune
	last  rune
***REMOVED***

const (
	spaceCategoryText = " "
	wordCategoryText  = "W"
)

var (
	ecmaSpace = []rune***REMOVED***0x0009, 0x000e, 0x0020, 0x0021, 0x00a0, 0x00a1, 0x1680, 0x1681, 0x2000, 0x200b, 0x2028, 0x202a, 0x202f, 0x2030, 0x205f, 0x2060, 0x3000, 0x3001, 0xfeff, 0xff00***REMOVED***
	ecmaWord  = []rune***REMOVED***0x0030, 0x003a, 0x0041, 0x005b, 0x005f, 0x0060, 0x0061, 0x007b***REMOVED***
	ecmaDigit = []rune***REMOVED***0x0030, 0x003a***REMOVED***
)

var (
	AnyClass          = getCharSetFromOldString([]rune***REMOVED***0***REMOVED***, false)
	ECMAAnyClass      = getCharSetFromOldString([]rune***REMOVED***0, 0x000a, 0x000b, 0x000d, 0x000e***REMOVED***, false)
	NoneClass         = getCharSetFromOldString(nil, false)
	ECMAWordClass     = getCharSetFromOldString(ecmaWord, false)
	NotECMAWordClass  = getCharSetFromOldString(ecmaWord, true)
	ECMASpaceClass    = getCharSetFromOldString(ecmaSpace, false)
	NotECMASpaceClass = getCharSetFromOldString(ecmaSpace, true)
	ECMADigitClass    = getCharSetFromOldString(ecmaDigit, false)
	NotECMADigitClass = getCharSetFromOldString(ecmaDigit, true)

	WordClass     = getCharSetFromCategoryString(false, false, wordCategoryText)
	NotWordClass  = getCharSetFromCategoryString(true, false, wordCategoryText)
	SpaceClass    = getCharSetFromCategoryString(false, false, spaceCategoryText)
	NotSpaceClass = getCharSetFromCategoryString(true, false, spaceCategoryText)
	DigitClass    = getCharSetFromCategoryString(false, false, "Nd")
	NotDigitClass = getCharSetFromCategoryString(false, true, "Nd")
)

var unicodeCategories = func() map[string]*unicode.RangeTable ***REMOVED***
	retVal := make(map[string]*unicode.RangeTable)
	for k, v := range unicode.Scripts ***REMOVED***
		retVal[k] = v
	***REMOVED***
	for k, v := range unicode.Categories ***REMOVED***
		retVal[k] = v
	***REMOVED***
	for k, v := range unicode.Properties ***REMOVED***
		retVal[k] = v
	***REMOVED***
	return retVal
***REMOVED***()

func getCharSetFromCategoryString(negateSet bool, negateCat bool, cats ...string) func() *CharSet ***REMOVED***
	if negateCat && negateSet ***REMOVED***
		panic("BUG!  You should only negate the set OR the category in a constant setup, but not both")
	***REMOVED***

	c := CharSet***REMOVED***negate: negateSet***REMOVED***

	c.categories = make([]category, len(cats))
	for i, cat := range cats ***REMOVED***
		c.categories[i] = category***REMOVED***cat: cat, negate: negateCat***REMOVED***
	***REMOVED***
	return func() *CharSet ***REMOVED***
		//make a copy each time
		local := c
		//return that address
		return &local
	***REMOVED***
***REMOVED***

func getCharSetFromOldString(setText []rune, negate bool) func() *CharSet ***REMOVED***
	c := CharSet***REMOVED******REMOVED***
	if len(setText) > 0 ***REMOVED***
		fillFirst := false
		l := len(setText)
		if negate ***REMOVED***
			if setText[0] == 0 ***REMOVED***
				setText = setText[1:]
			***REMOVED*** else ***REMOVED***
				l++
				fillFirst = true
			***REMOVED***
		***REMOVED***

		if l%2 == 0 ***REMOVED***
			c.ranges = make([]singleRange, l/2)
		***REMOVED*** else ***REMOVED***
			c.ranges = make([]singleRange, l/2+1)
		***REMOVED***

		first := true
		if fillFirst ***REMOVED***
			c.ranges[0] = singleRange***REMOVED***first: 0***REMOVED***
			first = false
		***REMOVED***

		i := 0
		for _, r := range setText ***REMOVED***
			if first ***REMOVED***
				// lower bound in a new range
				c.ranges[i] = singleRange***REMOVED***first: r***REMOVED***
				first = false
			***REMOVED*** else ***REMOVED***
				c.ranges[i].last = r - 1
				i++
				first = true
			***REMOVED***
		***REMOVED***
		if !first ***REMOVED***
			c.ranges[i].last = utf8.MaxRune
		***REMOVED***
	***REMOVED***

	return func() *CharSet ***REMOVED***
		local := c
		return &local
	***REMOVED***
***REMOVED***

// Copy makes a deep copy to prevent accidental mutation of a set
func (c CharSet) Copy() CharSet ***REMOVED***
	ret := CharSet***REMOVED***
		anything: c.anything,
		negate:   c.negate,
	***REMOVED***

	ret.ranges = append(ret.ranges, c.ranges...)
	ret.categories = append(ret.categories, c.categories...)

	if c.sub != nil ***REMOVED***
		sub := c.sub.Copy()
		ret.sub = &sub
	***REMOVED***

	return ret
***REMOVED***

// gets a human-readable description for a set string
func (c CharSet) String() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	buf.WriteRune('[')

	if c.IsNegated() ***REMOVED***
		buf.WriteRune('^')
	***REMOVED***

	for _, r := range c.ranges ***REMOVED***

		buf.WriteString(CharDescription(r.first))
		if r.first != r.last ***REMOVED***
			if r.last-r.first != 1 ***REMOVED***
				//groups that are 1 char apart skip the dash
				buf.WriteRune('-')
			***REMOVED***
			buf.WriteString(CharDescription(r.last))
		***REMOVED***
	***REMOVED***

	for _, c := range c.categories ***REMOVED***
		buf.WriteString(c.String())
	***REMOVED***

	if c.sub != nil ***REMOVED***
		buf.WriteRune('-')
		buf.WriteString(c.sub.String())
	***REMOVED***

	buf.WriteRune(']')

	return buf.String()
***REMOVED***

// mapHashFill converts a charset into a buffer for use in maps
func (c CharSet) mapHashFill(buf *bytes.Buffer) ***REMOVED***
	if c.negate ***REMOVED***
		buf.WriteByte(0)
	***REMOVED*** else ***REMOVED***
		buf.WriteByte(1)
	***REMOVED***

	binary.Write(buf, binary.LittleEndian, len(c.ranges))
	binary.Write(buf, binary.LittleEndian, len(c.categories))
	for _, r := range c.ranges ***REMOVED***
		buf.WriteRune(r.first)
		buf.WriteRune(r.last)
	***REMOVED***
	for _, ct := range c.categories ***REMOVED***
		buf.WriteString(ct.cat)
		if ct.negate ***REMOVED***
			buf.WriteByte(1)
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(0)
		***REMOVED***
	***REMOVED***

	if c.sub != nil ***REMOVED***
		c.sub.mapHashFill(buf)
	***REMOVED***
***REMOVED***

// CharIn returns true if the rune is in our character set (either ranges or categories).
// It handles negations and subtracted sub-charsets.
func (c CharSet) CharIn(ch rune) bool ***REMOVED***
	val := false
	// in s && !s.subtracted

	//check ranges
	for _, r := range c.ranges ***REMOVED***
		if ch < r.first ***REMOVED***
			continue
		***REMOVED***
		if ch <= r.last ***REMOVED***
			val = true
			break
		***REMOVED***
	***REMOVED***

	//check categories if we haven't already found a range
	if !val && len(c.categories) > 0 ***REMOVED***
		for _, ct := range c.categories ***REMOVED***
			// special categories...then unicode
			if ct.cat == spaceCategoryText ***REMOVED***
				if unicode.IsSpace(ch) ***REMOVED***
					// we found a space so we're done
					// negate means this is a "bad" thing
					val = !ct.negate
					break
				***REMOVED*** else if ct.negate ***REMOVED***
					val = true
					break
				***REMOVED***
			***REMOVED*** else if ct.cat == wordCategoryText ***REMOVED***
				if IsWordChar(ch) ***REMOVED***
					val = !ct.negate
					break
				***REMOVED*** else if ct.negate ***REMOVED***
					val = true
					break
				***REMOVED***
			***REMOVED*** else if unicode.Is(unicodeCategories[ct.cat], ch) ***REMOVED***
				// if we're in this unicode category then we're done
				// if negate=true on this category then we "failed" our test
				// otherwise we're good that we found it
				val = !ct.negate
				break
			***REMOVED*** else if ct.negate ***REMOVED***
				val = true
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// negate the whole char set
	if c.negate ***REMOVED***
		val = !val
	***REMOVED***

	// get subtracted recurse
	if val && c.sub != nil ***REMOVED***
		val = !c.sub.CharIn(ch)
	***REMOVED***

	//log.Printf("Char '%v' in %v == %v", string(ch), c.String(), val)
	return val
***REMOVED***

func (c category) String() string ***REMOVED***
	switch c.cat ***REMOVED***
	case spaceCategoryText:
		if c.negate ***REMOVED***
			return "\\S"
		***REMOVED***
		return "\\s"
	case wordCategoryText:
		if c.negate ***REMOVED***
			return "\\W"
		***REMOVED***
		return "\\w"
	***REMOVED***
	if _, ok := unicodeCategories[c.cat]; ok ***REMOVED***

		if c.negate ***REMOVED***
			return "\\P***REMOVED***" + c.cat + "***REMOVED***"
		***REMOVED***
		return "\\p***REMOVED***" + c.cat + "***REMOVED***"
	***REMOVED***
	return "Unknown category: " + c.cat
***REMOVED***

// CharDescription Produces a human-readable description for a single character.
func CharDescription(ch rune) string ***REMOVED***
	/*if ch == '\\' ***REMOVED***
		return "\\\\"
	***REMOVED***

	if ch > ' ' && ch <= '~' ***REMOVED***
		return string(ch)
	***REMOVED*** else if ch == '\n' ***REMOVED***
		return "\\n"
	***REMOVED*** else if ch == ' ' ***REMOVED***
		return "\\ "
	***REMOVED****/

	b := &bytes.Buffer***REMOVED******REMOVED***
	escape(b, ch, false) //fmt.Sprintf("%U", ch)
	return b.String()
***REMOVED***

// According to UTS#18 Unicode Regular Expressions (http://www.unicode.org/reports/tr18/)
// RL 1.4 Simple Word Boundaries  The class of <word_character> includes all Alphabetic
// values from the Unicode character database, from UnicodeData.txt [UData], plus the U+200C
// ZERO WIDTH NON-JOINER and U+200D ZERO WIDTH JOINER.
func IsWordChar(r rune) bool ***REMOVED***
	//"L", "Mn", "Nd", "Pc"
	return unicode.In(r,
		unicode.Categories["L"], unicode.Categories["Mn"],
		unicode.Categories["Nd"], unicode.Categories["Pc"]) || r == '\u200D' || r == '\u200C'
	//return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_'
***REMOVED***

func IsECMAWordChar(r rune) bool ***REMOVED***
	return unicode.In(r,
		unicode.Categories["L"], unicode.Categories["Mn"],
		unicode.Categories["Nd"], unicode.Categories["Pc"])

	//return 'A' <= r && r <= 'Z' || 'a' <= r && r <= 'z' || '0' <= r && r <= '9' || r == '_'
***REMOVED***

// SingletonChar will return the char from the first range without validation.
// It assumes you have checked for IsSingleton or IsSingletonInverse and will panic given bad input
func (c CharSet) SingletonChar() rune ***REMOVED***
	return c.ranges[0].first
***REMOVED***

func (c CharSet) IsSingleton() bool ***REMOVED***
	return !c.negate && //negated is multiple chars
		len(c.categories) == 0 && len(c.ranges) == 1 && // multiple ranges and unicode classes represent multiple chars
		c.sub == nil && // subtraction means we've got multiple chars
		c.ranges[0].first == c.ranges[0].last // first and last equal means we're just 1 char
***REMOVED***

func (c CharSet) IsSingletonInverse() bool ***REMOVED***
	return c.negate && //same as above, but requires negated
		len(c.categories) == 0 && len(c.ranges) == 1 && // multiple ranges and unicode classes represent multiple chars
		c.sub == nil && // subtraction means we've got multiple chars
		c.ranges[0].first == c.ranges[0].last // first and last equal means we're just 1 char
***REMOVED***

func (c CharSet) IsMergeable() bool ***REMOVED***
	return !c.IsNegated() && !c.HasSubtraction()
***REMOVED***

func (c CharSet) IsNegated() bool ***REMOVED***
	return c.negate
***REMOVED***

func (c CharSet) HasSubtraction() bool ***REMOVED***
	return c.sub != nil
***REMOVED***

func (c CharSet) IsEmpty() bool ***REMOVED***
	return len(c.ranges) == 0 && len(c.categories) == 0 && c.sub == nil
***REMOVED***

func (c *CharSet) addDigit(ecma, negate bool, pattern string) ***REMOVED***
	if ecma ***REMOVED***
		if negate ***REMOVED***
			c.addRanges(NotECMADigitClass().ranges)
		***REMOVED*** else ***REMOVED***
			c.addRanges(ECMADigitClass().ranges)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.addCategories(category***REMOVED***cat: "Nd", negate: negate***REMOVED***)
	***REMOVED***
***REMOVED***

func (c *CharSet) addChar(ch rune) ***REMOVED***
	c.addRange(ch, ch)
***REMOVED***

func (c *CharSet) addSpace(ecma, negate bool) ***REMOVED***
	if ecma ***REMOVED***
		if negate ***REMOVED***
			c.addRanges(NotECMASpaceClass().ranges)
		***REMOVED*** else ***REMOVED***
			c.addRanges(ECMASpaceClass().ranges)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.addCategories(category***REMOVED***cat: spaceCategoryText, negate: negate***REMOVED***)
	***REMOVED***
***REMOVED***

func (c *CharSet) addWord(ecma, negate bool) ***REMOVED***
	if ecma ***REMOVED***
		if negate ***REMOVED***
			c.addRanges(NotECMAWordClass().ranges)
		***REMOVED*** else ***REMOVED***
			c.addRanges(ECMAWordClass().ranges)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		c.addCategories(category***REMOVED***cat: wordCategoryText, negate: negate***REMOVED***)
	***REMOVED***
***REMOVED***

// Add set ranges and categories into ours -- no deduping or anything
func (c *CharSet) addSet(set CharSet) ***REMOVED***
	if c.anything ***REMOVED***
		return
	***REMOVED***
	if set.anything ***REMOVED***
		c.makeAnything()
		return
	***REMOVED***
	// just append here to prevent double-canon
	c.ranges = append(c.ranges, set.ranges...)
	c.addCategories(set.categories...)
	c.canonicalize()
***REMOVED***

func (c *CharSet) makeAnything() ***REMOVED***
	c.anything = true
	c.categories = []category***REMOVED******REMOVED***
	c.ranges = AnyClass().ranges
***REMOVED***

func (c *CharSet) addCategories(cats ...category) ***REMOVED***
	// don't add dupes and remove positive+negative
	if c.anything ***REMOVED***
		// if we've had a previous positive+negative group then
		// just return, we're as broad as we can get
		return
	***REMOVED***

	for _, ct := range cats ***REMOVED***
		found := false
		for _, ct2 := range c.categories ***REMOVED***
			if ct.cat == ct2.cat ***REMOVED***
				if ct.negate != ct2.negate ***REMOVED***
					// oposite negations...this mean we just
					// take us as anything and move on
					c.makeAnything()
					return
				***REMOVED***
				found = true
				break
			***REMOVED***
		***REMOVED***

		if !found ***REMOVED***
			c.categories = append(c.categories, ct)
		***REMOVED***
	***REMOVED***
***REMOVED***

// Merges new ranges to our own
func (c *CharSet) addRanges(ranges []singleRange) ***REMOVED***
	if c.anything ***REMOVED***
		return
	***REMOVED***
	c.ranges = append(c.ranges, ranges...)
	c.canonicalize()
***REMOVED***

// Merges everything but the new ranges into our own
func (c *CharSet) addNegativeRanges(ranges []singleRange) ***REMOVED***
	if c.anything ***REMOVED***
		return
	***REMOVED***

	var hi rune

	// convert incoming ranges into opposites, assume they are in order
	for _, r := range ranges ***REMOVED***
		if hi < r.first ***REMOVED***
			c.ranges = append(c.ranges, singleRange***REMOVED***hi, r.first - 1***REMOVED***)
		***REMOVED***
		hi = r.last + 1
	***REMOVED***

	if hi < utf8.MaxRune ***REMOVED***
		c.ranges = append(c.ranges, singleRange***REMOVED***hi, utf8.MaxRune***REMOVED***)
	***REMOVED***

	c.canonicalize()
***REMOVED***

func isValidUnicodeCat(catName string) bool ***REMOVED***
	_, ok := unicodeCategories[catName]
	return ok
***REMOVED***

func (c *CharSet) addCategory(categoryName string, negate, caseInsensitive bool, pattern string) ***REMOVED***
	if !isValidUnicodeCat(categoryName) ***REMOVED***
		// unknown unicode category, script, or property "blah"
		panic(fmt.Errorf("Unknown unicode category, script, or property '%v'", categoryName))

	***REMOVED***

	if caseInsensitive && (categoryName == "Ll" || categoryName == "Lu" || categoryName == "Lt") ***REMOVED***
		// when RegexOptions.IgnoreCase is specified then ***REMOVED***Ll***REMOVED*** ***REMOVED***Lu***REMOVED*** and ***REMOVED***Lt***REMOVED*** cases should all match
		c.addCategories(
			category***REMOVED***cat: "Ll", negate: negate***REMOVED***,
			category***REMOVED***cat: "Lu", negate: negate***REMOVED***,
			category***REMOVED***cat: "Lt", negate: negate***REMOVED***)
	***REMOVED***
	c.addCategories(category***REMOVED***cat: categoryName, negate: negate***REMOVED***)
***REMOVED***

func (c *CharSet) addSubtraction(sub *CharSet) ***REMOVED***
	c.sub = sub
***REMOVED***

func (c *CharSet) addRange(chMin, chMax rune) ***REMOVED***
	c.ranges = append(c.ranges, singleRange***REMOVED***first: chMin, last: chMax***REMOVED***)
	c.canonicalize()
***REMOVED***

func (c *CharSet) addNamedASCII(name string, negate bool) bool ***REMOVED***
	var rs []singleRange

	switch name ***REMOVED***
	case "alnum":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'0', '9'***REMOVED***, singleRange***REMOVED***'A', 'Z'***REMOVED***, singleRange***REMOVED***'a', 'z'***REMOVED******REMOVED***
	case "alpha":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'A', 'Z'***REMOVED***, singleRange***REMOVED***'a', 'z'***REMOVED******REMOVED***
	case "ascii":
		rs = []singleRange***REMOVED***singleRange***REMOVED***0, 0x7f***REMOVED******REMOVED***
	case "blank":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'\t', '\t'***REMOVED***, singleRange***REMOVED***' ', ' '***REMOVED******REMOVED***
	case "cntrl":
		rs = []singleRange***REMOVED***singleRange***REMOVED***0, 0x1f***REMOVED***, singleRange***REMOVED***0x7f, 0x7f***REMOVED******REMOVED***
	case "digit":
		c.addDigit(false, negate, "")
	case "graph":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'!', '~'***REMOVED******REMOVED***
	case "lower":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'a', 'z'***REMOVED******REMOVED***
	case "print":
		rs = []singleRange***REMOVED***singleRange***REMOVED***' ', '~'***REMOVED******REMOVED***
	case "punct": //[!-/:-@[-`***REMOVED***-~]
		rs = []singleRange***REMOVED***singleRange***REMOVED***'!', '/'***REMOVED***, singleRange***REMOVED***':', '@'***REMOVED***, singleRange***REMOVED***'[', '`'***REMOVED***, singleRange***REMOVED***'***REMOVED***', '~'***REMOVED******REMOVED***
	case "space":
		c.addSpace(true, negate)
	case "upper":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'A', 'Z'***REMOVED******REMOVED***
	case "word":
		c.addWord(true, negate)
	case "xdigit":
		rs = []singleRange***REMOVED***singleRange***REMOVED***'0', '9'***REMOVED***, singleRange***REMOVED***'A', 'F'***REMOVED***, singleRange***REMOVED***'a', 'f'***REMOVED******REMOVED***
	default:
		return false
	***REMOVED***

	if len(rs) > 0 ***REMOVED***
		if negate ***REMOVED***
			c.addNegativeRanges(rs)
		***REMOVED*** else ***REMOVED***
			c.addRanges(rs)
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

type singleRangeSorter []singleRange

func (p singleRangeSorter) Len() int           ***REMOVED*** return len(p) ***REMOVED***
func (p singleRangeSorter) Less(i, j int) bool ***REMOVED*** return p[i].first < p[j].first ***REMOVED***
func (p singleRangeSorter) Swap(i, j int)      ***REMOVED*** p[i], p[j] = p[j], p[i] ***REMOVED***

// Logic to reduce a character class to a unique, sorted form.
func (c *CharSet) canonicalize() ***REMOVED***
	var i, j int
	var last rune

	//
	// Find and eliminate overlapping or abutting ranges
	//

	if len(c.ranges) > 1 ***REMOVED***
		sort.Sort(singleRangeSorter(c.ranges))

		done := false

		for i, j = 1, 0; ; i++ ***REMOVED***
			for last = c.ranges[j].last; ; i++ ***REMOVED***
				if i == len(c.ranges) || last == utf8.MaxRune ***REMOVED***
					done = true
					break
				***REMOVED***

				CurrentRange := c.ranges[i]
				if CurrentRange.first > last+1 ***REMOVED***
					break
				***REMOVED***

				if last < CurrentRange.last ***REMOVED***
					last = CurrentRange.last
				***REMOVED***
			***REMOVED***

			c.ranges[j] = singleRange***REMOVED***first: c.ranges[j].first, last: last***REMOVED***

			j++

			if done ***REMOVED***
				break
			***REMOVED***

			if j < i ***REMOVED***
				c.ranges[j] = c.ranges[i]
			***REMOVED***
		***REMOVED***

		c.ranges = append(c.ranges[:j], c.ranges[len(c.ranges):]...)
	***REMOVED***
***REMOVED***

// Adds to the class any lowercase versions of characters already
// in the class. Used for case-insensitivity.
func (c *CharSet) addLowercase() ***REMOVED***
	if c.anything ***REMOVED***
		return
	***REMOVED***
	toAdd := []singleRange***REMOVED******REMOVED***
	for i := 0; i < len(c.ranges); i++ ***REMOVED***
		r := c.ranges[i]
		if r.first == r.last ***REMOVED***
			lower := unicode.ToLower(r.first)
			c.ranges[i] = singleRange***REMOVED***first: lower, last: lower***REMOVED***
		***REMOVED*** else ***REMOVED***
			toAdd = append(toAdd, r)
		***REMOVED***
	***REMOVED***

	for _, r := range toAdd ***REMOVED***
		c.addLowercaseRange(r.first, r.last)
	***REMOVED***
	c.canonicalize()
***REMOVED***

/**************************************************************************
    Let U be the set of Unicode character values and let L be the lowercase
    function, mapping from U to U. To perform case insensitive matching of
    character sets, we need to be able to map an interval I in U, say

        I = [chMin, chMax] = ***REMOVED*** ch : chMin <= ch <= chMax ***REMOVED***

    to a set A such that A contains L(I) and A is contained in the union of
    I and L(I).

    The table below partitions U into intervals on which L is non-decreasing.
    Thus, for any interval J = [a, b] contained in one of these intervals,
    L(J) is contained in [L(a), L(b)].

    It is also true that for any such J, [L(a), L(b)] is contained in the
    union of J and L(J). This does not follow from L being non-decreasing on
    these intervals. It follows from the nature of the L on each interval.
    On each interval, L has one of the following forms:

        (1) L(ch) = constant            (LowercaseSet)
        (2) L(ch) = ch + offset         (LowercaseAdd)
        (3) L(ch) = ch | 1              (LowercaseBor)
        (4) L(ch) = ch + (ch & 1)       (LowercaseBad)

    It is easy to verify that for any of these forms [L(a), L(b)] is
    contained in the union of [a, b] and L([a, b]).
***************************************************************************/

const (
	LowercaseSet = 0 // Set to arg.
	LowercaseAdd = 1 // Add arg.
	LowercaseBor = 2 // Bitwise or with 1.
	LowercaseBad = 3 // Bitwise and with 1 and add original.
)

type lcMap struct ***REMOVED***
	chMin, chMax rune
	op, data     int32
***REMOVED***

var lcTable = []lcMap***REMOVED***
	lcMap***REMOVED***'\u0041', '\u005A', LowercaseAdd, 32***REMOVED***,
	lcMap***REMOVED***'\u00C0', '\u00DE', LowercaseAdd, 32***REMOVED***,
	lcMap***REMOVED***'\u0100', '\u012E', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0130', '\u0130', LowercaseSet, 0x0069***REMOVED***,
	lcMap***REMOVED***'\u0132', '\u0136', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0139', '\u0147', LowercaseBad, 0***REMOVED***,
	lcMap***REMOVED***'\u014A', '\u0176', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0178', '\u0178', LowercaseSet, 0x00FF***REMOVED***,
	lcMap***REMOVED***'\u0179', '\u017D', LowercaseBad, 0***REMOVED***,
	lcMap***REMOVED***'\u0181', '\u0181', LowercaseSet, 0x0253***REMOVED***,
	lcMap***REMOVED***'\u0182', '\u0184', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0186', '\u0186', LowercaseSet, 0x0254***REMOVED***,
	lcMap***REMOVED***'\u0187', '\u0187', LowercaseSet, 0x0188***REMOVED***,
	lcMap***REMOVED***'\u0189', '\u018A', LowercaseAdd, 205***REMOVED***,
	lcMap***REMOVED***'\u018B', '\u018B', LowercaseSet, 0x018C***REMOVED***,
	lcMap***REMOVED***'\u018E', '\u018E', LowercaseSet, 0x01DD***REMOVED***,
	lcMap***REMOVED***'\u018F', '\u018F', LowercaseSet, 0x0259***REMOVED***,
	lcMap***REMOVED***'\u0190', '\u0190', LowercaseSet, 0x025B***REMOVED***,
	lcMap***REMOVED***'\u0191', '\u0191', LowercaseSet, 0x0192***REMOVED***,
	lcMap***REMOVED***'\u0193', '\u0193', LowercaseSet, 0x0260***REMOVED***,
	lcMap***REMOVED***'\u0194', '\u0194', LowercaseSet, 0x0263***REMOVED***,
	lcMap***REMOVED***'\u0196', '\u0196', LowercaseSet, 0x0269***REMOVED***,
	lcMap***REMOVED***'\u0197', '\u0197', LowercaseSet, 0x0268***REMOVED***,
	lcMap***REMOVED***'\u0198', '\u0198', LowercaseSet, 0x0199***REMOVED***,
	lcMap***REMOVED***'\u019C', '\u019C', LowercaseSet, 0x026F***REMOVED***,
	lcMap***REMOVED***'\u019D', '\u019D', LowercaseSet, 0x0272***REMOVED***,
	lcMap***REMOVED***'\u019F', '\u019F', LowercaseSet, 0x0275***REMOVED***,
	lcMap***REMOVED***'\u01A0', '\u01A4', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u01A7', '\u01A7', LowercaseSet, 0x01A8***REMOVED***,
	lcMap***REMOVED***'\u01A9', '\u01A9', LowercaseSet, 0x0283***REMOVED***,
	lcMap***REMOVED***'\u01AC', '\u01AC', LowercaseSet, 0x01AD***REMOVED***,
	lcMap***REMOVED***'\u01AE', '\u01AE', LowercaseSet, 0x0288***REMOVED***,
	lcMap***REMOVED***'\u01AF', '\u01AF', LowercaseSet, 0x01B0***REMOVED***,
	lcMap***REMOVED***'\u01B1', '\u01B2', LowercaseAdd, 217***REMOVED***,
	lcMap***REMOVED***'\u01B3', '\u01B5', LowercaseBad, 0***REMOVED***,
	lcMap***REMOVED***'\u01B7', '\u01B7', LowercaseSet, 0x0292***REMOVED***,
	lcMap***REMOVED***'\u01B8', '\u01B8', LowercaseSet, 0x01B9***REMOVED***,
	lcMap***REMOVED***'\u01BC', '\u01BC', LowercaseSet, 0x01BD***REMOVED***,
	lcMap***REMOVED***'\u01C4', '\u01C5', LowercaseSet, 0x01C6***REMOVED***,
	lcMap***REMOVED***'\u01C7', '\u01C8', LowercaseSet, 0x01C9***REMOVED***,
	lcMap***REMOVED***'\u01CA', '\u01CB', LowercaseSet, 0x01CC***REMOVED***,
	lcMap***REMOVED***'\u01CD', '\u01DB', LowercaseBad, 0***REMOVED***,
	lcMap***REMOVED***'\u01DE', '\u01EE', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u01F1', '\u01F2', LowercaseSet, 0x01F3***REMOVED***,
	lcMap***REMOVED***'\u01F4', '\u01F4', LowercaseSet, 0x01F5***REMOVED***,
	lcMap***REMOVED***'\u01FA', '\u0216', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0386', '\u0386', LowercaseSet, 0x03AC***REMOVED***,
	lcMap***REMOVED***'\u0388', '\u038A', LowercaseAdd, 37***REMOVED***,
	lcMap***REMOVED***'\u038C', '\u038C', LowercaseSet, 0x03CC***REMOVED***,
	lcMap***REMOVED***'\u038E', '\u038F', LowercaseAdd, 63***REMOVED***,
	lcMap***REMOVED***'\u0391', '\u03AB', LowercaseAdd, 32***REMOVED***,
	lcMap***REMOVED***'\u03E2', '\u03EE', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0401', '\u040F', LowercaseAdd, 80***REMOVED***,
	lcMap***REMOVED***'\u0410', '\u042F', LowercaseAdd, 32***REMOVED***,
	lcMap***REMOVED***'\u0460', '\u0480', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u0490', '\u04BE', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u04C1', '\u04C3', LowercaseBad, 0***REMOVED***,
	lcMap***REMOVED***'\u04C7', '\u04C7', LowercaseSet, 0x04C8***REMOVED***,
	lcMap***REMOVED***'\u04CB', '\u04CB', LowercaseSet, 0x04CC***REMOVED***,
	lcMap***REMOVED***'\u04D0', '\u04EA', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u04EE', '\u04F4', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u04F8', '\u04F8', LowercaseSet, 0x04F9***REMOVED***,
	lcMap***REMOVED***'\u0531', '\u0556', LowercaseAdd, 48***REMOVED***,
	lcMap***REMOVED***'\u10A0', '\u10C5', LowercaseAdd, 48***REMOVED***,
	lcMap***REMOVED***'\u1E00', '\u1EF8', LowercaseBor, 0***REMOVED***,
	lcMap***REMOVED***'\u1F08', '\u1F0F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F18', '\u1F1F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F28', '\u1F2F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F38', '\u1F3F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F48', '\u1F4D', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F59', '\u1F59', LowercaseSet, 0x1F51***REMOVED***,
	lcMap***REMOVED***'\u1F5B', '\u1F5B', LowercaseSet, 0x1F53***REMOVED***,
	lcMap***REMOVED***'\u1F5D', '\u1F5D', LowercaseSet, 0x1F55***REMOVED***,
	lcMap***REMOVED***'\u1F5F', '\u1F5F', LowercaseSet, 0x1F57***REMOVED***,
	lcMap***REMOVED***'\u1F68', '\u1F6F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F88', '\u1F8F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1F98', '\u1F9F', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1FA8', '\u1FAF', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1FB8', '\u1FB9', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1FBA', '\u1FBB', LowercaseAdd, -74***REMOVED***,
	lcMap***REMOVED***'\u1FBC', '\u1FBC', LowercaseSet, 0x1FB3***REMOVED***,
	lcMap***REMOVED***'\u1FC8', '\u1FCB', LowercaseAdd, -86***REMOVED***,
	lcMap***REMOVED***'\u1FCC', '\u1FCC', LowercaseSet, 0x1FC3***REMOVED***,
	lcMap***REMOVED***'\u1FD8', '\u1FD9', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1FDA', '\u1FDB', LowercaseAdd, -100***REMOVED***,
	lcMap***REMOVED***'\u1FE8', '\u1FE9', LowercaseAdd, -8***REMOVED***,
	lcMap***REMOVED***'\u1FEA', '\u1FEB', LowercaseAdd, -112***REMOVED***,
	lcMap***REMOVED***'\u1FEC', '\u1FEC', LowercaseSet, 0x1FE5***REMOVED***,
	lcMap***REMOVED***'\u1FF8', '\u1FF9', LowercaseAdd, -128***REMOVED***,
	lcMap***REMOVED***'\u1FFA', '\u1FFB', LowercaseAdd, -126***REMOVED***,
	lcMap***REMOVED***'\u1FFC', '\u1FFC', LowercaseSet, 0x1FF3***REMOVED***,
	lcMap***REMOVED***'\u2160', '\u216F', LowercaseAdd, 16***REMOVED***,
	lcMap***REMOVED***'\u24B6', '\u24D0', LowercaseAdd, 26***REMOVED***,
	lcMap***REMOVED***'\uFF21', '\uFF3A', LowercaseAdd, 32***REMOVED***,
***REMOVED***

func (c *CharSet) addLowercaseRange(chMin, chMax rune) ***REMOVED***
	var i, iMax, iMid int
	var chMinT, chMaxT rune
	var lc lcMap

	for i, iMax = 0, len(lcTable); i < iMax; ***REMOVED***
		iMid = (i + iMax) / 2
		if lcTable[iMid].chMax < chMin ***REMOVED***
			i = iMid + 1
		***REMOVED*** else ***REMOVED***
			iMax = iMid
		***REMOVED***
	***REMOVED***

	for ; i < len(lcTable); i++ ***REMOVED***
		lc = lcTable[i]
		if lc.chMin > chMax ***REMOVED***
			return
		***REMOVED***
		chMinT = lc.chMin
		if chMinT < chMin ***REMOVED***
			chMinT = chMin
		***REMOVED***

		chMaxT = lc.chMax
		if chMaxT > chMax ***REMOVED***
			chMaxT = chMax
		***REMOVED***

		switch lc.op ***REMOVED***
		case LowercaseSet:
			chMinT = rune(lc.data)
			chMaxT = rune(lc.data)
			break
		case LowercaseAdd:
			chMinT += lc.data
			chMaxT += lc.data
			break
		case LowercaseBor:
			chMinT |= 1
			chMaxT |= 1
			break
		case LowercaseBad:
			chMinT += (chMinT & 1)
			chMaxT += (chMaxT & 1)
			break
		***REMOVED***

		if chMinT < chMin || chMaxT > chMax ***REMOVED***
			c.addRange(chMinT, chMaxT)
		***REMOVED***
	***REMOVED***
***REMOVED***
