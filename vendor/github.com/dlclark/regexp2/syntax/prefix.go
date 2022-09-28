package syntax

import (
	"bytes"
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type Prefix struct ***REMOVED***
	PrefixStr       []rune
	PrefixSet       CharSet
	CaseInsensitive bool
***REMOVED***

// It takes a RegexTree and computes the set of chars that can start it.
func getFirstCharsPrefix(tree *RegexTree) *Prefix ***REMOVED***
	s := regexFcd***REMOVED***
		fcStack:  make([]regexFc, 32),
		intStack: make([]int, 32),
	***REMOVED***
	fc := s.regexFCFromRegexTree(tree)

	if fc == nil || fc.nullable || fc.cc.IsEmpty() ***REMOVED***
		return nil
	***REMOVED***
	fcSet := fc.getFirstChars()
	return &Prefix***REMOVED***PrefixSet: fcSet, CaseInsensitive: fc.caseInsensitive***REMOVED***
***REMOVED***

type regexFcd struct ***REMOVED***
	intStack        []int
	intDepth        int
	fcStack         []regexFc
	fcDepth         int
	skipAllChildren bool // don't process any more children at the current level
	skipchild       bool // don't process the current child.
	failed          bool
***REMOVED***

/*
 * The main FC computation. It does a shortcutted depth-first walk
 * through the tree and calls CalculateFC to emits code before
 * and after each child of an interior node, and at each leaf.
 */
func (s *regexFcd) regexFCFromRegexTree(tree *RegexTree) *regexFc ***REMOVED***
	curNode := tree.root
	curChild := 0

	for ***REMOVED***
		if len(curNode.children) == 0 ***REMOVED***
			// This is a leaf node
			s.calculateFC(curNode.t, curNode, 0)
		***REMOVED*** else if curChild < len(curNode.children) && !s.skipAllChildren ***REMOVED***
			// This is an interior node, and we have more children to analyze
			s.calculateFC(curNode.t|beforeChild, curNode, curChild)

			if !s.skipchild ***REMOVED***
				curNode = curNode.children[curChild]
				// this stack is how we get a depth first walk of the tree.
				s.pushInt(curChild)
				curChild = 0
			***REMOVED*** else ***REMOVED***
				curChild++
				s.skipchild = false
			***REMOVED***
			continue
		***REMOVED***

		// This is an interior node where we've finished analyzing all the children, or
		// the end of a leaf node.
		s.skipAllChildren = false

		if s.intIsEmpty() ***REMOVED***
			break
		***REMOVED***

		curChild = s.popInt()
		curNode = curNode.next

		s.calculateFC(curNode.t|afterChild, curNode, curChild)
		if s.failed ***REMOVED***
			return nil
		***REMOVED***

		curChild++
	***REMOVED***

	if s.fcIsEmpty() ***REMOVED***
		return nil
	***REMOVED***

	return s.popFC()
***REMOVED***

// To avoid recursion, we use a simple integer stack.
// This is the push.
func (s *regexFcd) pushInt(I int) ***REMOVED***
	if s.intDepth >= len(s.intStack) ***REMOVED***
		expanded := make([]int, s.intDepth*2)
		copy(expanded, s.intStack)
		s.intStack = expanded
	***REMOVED***

	s.intStack[s.intDepth] = I
	s.intDepth++
***REMOVED***

// True if the stack is empty.
func (s *regexFcd) intIsEmpty() bool ***REMOVED***
	return s.intDepth == 0
***REMOVED***

// This is the pop.
func (s *regexFcd) popInt() int ***REMOVED***
	s.intDepth--
	return s.intStack[s.intDepth]
***REMOVED***

// We also use a stack of RegexFC objects.
// This is the push.
func (s *regexFcd) pushFC(fc regexFc) ***REMOVED***
	if s.fcDepth >= len(s.fcStack) ***REMOVED***
		expanded := make([]regexFc, s.fcDepth*2)
		copy(expanded, s.fcStack)
		s.fcStack = expanded
	***REMOVED***

	s.fcStack[s.fcDepth] = fc
	s.fcDepth++
***REMOVED***

// True if the stack is empty.
func (s *regexFcd) fcIsEmpty() bool ***REMOVED***
	return s.fcDepth == 0
***REMOVED***

// This is the pop.
func (s *regexFcd) popFC() *regexFc ***REMOVED***
	s.fcDepth--
	return &s.fcStack[s.fcDepth]
***REMOVED***

// This is the top.
func (s *regexFcd) topFC() *regexFc ***REMOVED***
	return &s.fcStack[s.fcDepth-1]
***REMOVED***

// Called in Beforechild to prevent further processing of the current child
func (s *regexFcd) skipChild() ***REMOVED***
	s.skipchild = true
***REMOVED***

// FC computation and shortcut cases for each node type
func (s *regexFcd) calculateFC(nt nodeType, node *regexNode, CurIndex int) ***REMOVED***
	//fmt.Printf("NodeType: %v, CurIndex: %v, Desc: %v\n", nt, CurIndex, node.description())
	ci := false
	rtl := false

	if nt <= ntRef ***REMOVED***
		if (node.options & IgnoreCase) != 0 ***REMOVED***
			ci = true
		***REMOVED***
		if (node.options & RightToLeft) != 0 ***REMOVED***
			rtl = true
		***REMOVED***
	***REMOVED***

	switch nt ***REMOVED***
	case ntConcatenate | beforeChild, ntAlternate | beforeChild, ntTestref | beforeChild, ntLoop | beforeChild, ntLazyloop | beforeChild:
		break

	case ntTestgroup | beforeChild:
		if CurIndex == 0 ***REMOVED***
			s.skipChild()
		***REMOVED***
		break

	case ntEmpty:
		s.pushFC(regexFc***REMOVED***nullable: true***REMOVED***)
		break

	case ntConcatenate | afterChild:
		if CurIndex != 0 ***REMOVED***
			child := s.popFC()
			cumul := s.topFC()

			s.failed = !cumul.addFC(*child, true)
		***REMOVED***

		fc := s.topFC()
		if !fc.nullable ***REMOVED***
			s.skipAllChildren = true
		***REMOVED***
		break

	case ntTestgroup | afterChild:
		if CurIndex > 1 ***REMOVED***
			child := s.popFC()
			cumul := s.topFC()

			s.failed = !cumul.addFC(*child, false)
		***REMOVED***
		break

	case ntAlternate | afterChild, ntTestref | afterChild:
		if CurIndex != 0 ***REMOVED***
			child := s.popFC()
			cumul := s.topFC()

			s.failed = !cumul.addFC(*child, false)
		***REMOVED***
		break

	case ntLoop | afterChild, ntLazyloop | afterChild:
		if node.m == 0 ***REMOVED***
			fc := s.topFC()
			fc.nullable = true
		***REMOVED***
		break

	case ntGroup | beforeChild, ntGroup | afterChild, ntCapture | beforeChild, ntCapture | afterChild, ntGreedy | beforeChild, ntGreedy | afterChild:
		break

	case ntRequire | beforeChild, ntPrevent | beforeChild:
		s.skipChild()
		s.pushFC(regexFc***REMOVED***nullable: true***REMOVED***)
		break

	case ntRequire | afterChild, ntPrevent | afterChild:
		break

	case ntOne, ntNotone:
		s.pushFC(newRegexFc(node.ch, nt == ntNotone, false, ci))
		break

	case ntOneloop, ntOnelazy:
		s.pushFC(newRegexFc(node.ch, false, node.m == 0, ci))
		break

	case ntNotoneloop, ntNotonelazy:
		s.pushFC(newRegexFc(node.ch, true, node.m == 0, ci))
		break

	case ntMulti:
		if len(node.str) == 0 ***REMOVED***
			s.pushFC(regexFc***REMOVED***nullable: true***REMOVED***)
		***REMOVED*** else if !rtl ***REMOVED***
			s.pushFC(newRegexFc(node.str[0], false, false, ci))
		***REMOVED*** else ***REMOVED***
			s.pushFC(newRegexFc(node.str[len(node.str)-1], false, false, ci))
		***REMOVED***
		break

	case ntSet:
		s.pushFC(regexFc***REMOVED***cc: node.set.Copy(), nullable: false, caseInsensitive: ci***REMOVED***)
		break

	case ntSetloop, ntSetlazy:
		s.pushFC(regexFc***REMOVED***cc: node.set.Copy(), nullable: node.m == 0, caseInsensitive: ci***REMOVED***)
		break

	case ntRef:
		s.pushFC(regexFc***REMOVED***cc: *AnyClass(), nullable: true, caseInsensitive: false***REMOVED***)
		break

	case ntNothing, ntBol, ntEol, ntBoundary, ntNonboundary, ntECMABoundary, ntNonECMABoundary, ntBeginning, ntStart, ntEndZ, ntEnd:
		s.pushFC(regexFc***REMOVED***nullable: true***REMOVED***)
		break

	default:
		panic(fmt.Sprintf("unexpected op code: %v", nt))
	***REMOVED***
***REMOVED***

type regexFc struct ***REMOVED***
	cc              CharSet
	nullable        bool
	caseInsensitive bool
***REMOVED***

func newRegexFc(ch rune, not, nullable, caseInsensitive bool) regexFc ***REMOVED***
	r := regexFc***REMOVED***
		caseInsensitive: caseInsensitive,
		nullable:        nullable,
	***REMOVED***
	if not ***REMOVED***
		if ch > 0 ***REMOVED***
			r.cc.addRange('\x00', ch-1)
		***REMOVED***
		if ch < 0xFFFF ***REMOVED***
			r.cc.addRange(ch+1, utf8.MaxRune)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.cc.addRange(ch, ch)
	***REMOVED***
	return r
***REMOVED***

func (r *regexFc) getFirstChars() CharSet ***REMOVED***
	if r.caseInsensitive ***REMOVED***
		r.cc.addLowercase()
	***REMOVED***

	return r.cc
***REMOVED***

func (r *regexFc) addFC(fc regexFc, concatenate bool) bool ***REMOVED***
	if !r.cc.IsMergeable() || !fc.cc.IsMergeable() ***REMOVED***
		return false
	***REMOVED***

	if concatenate ***REMOVED***
		if !r.nullable ***REMOVED***
			return true
		***REMOVED***

		if !fc.nullable ***REMOVED***
			r.nullable = false
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if fc.nullable ***REMOVED***
			r.nullable = true
		***REMOVED***
	***REMOVED***

	r.caseInsensitive = r.caseInsensitive || fc.caseInsensitive
	r.cc.addSet(fc.cc)

	return true
***REMOVED***

// This is a related computation: it takes a RegexTree and computes the
// leading substring if it sees one. It's quite trivial and gives up easily.
func getPrefix(tree *RegexTree) *Prefix ***REMOVED***
	var concatNode *regexNode
	nextChild := 0

	curNode := tree.root

	for ***REMOVED***
		switch curNode.t ***REMOVED***
		case ntConcatenate:
			if len(curNode.children) > 0 ***REMOVED***
				concatNode = curNode
				nextChild = 0
			***REMOVED***

		case ntGreedy, ntCapture:
			curNode = curNode.children[0]
			concatNode = nil
			continue

		case ntOneloop, ntOnelazy:
			if curNode.m > 0 ***REMOVED***
				return &Prefix***REMOVED***
					PrefixStr:       repeat(curNode.ch, curNode.m),
					CaseInsensitive: (curNode.options & IgnoreCase) != 0,
				***REMOVED***
			***REMOVED***
			return nil

		case ntOne:
			return &Prefix***REMOVED***
				PrefixStr:       []rune***REMOVED***curNode.ch***REMOVED***,
				CaseInsensitive: (curNode.options & IgnoreCase) != 0,
			***REMOVED***

		case ntMulti:
			return &Prefix***REMOVED***
				PrefixStr:       curNode.str,
				CaseInsensitive: (curNode.options & IgnoreCase) != 0,
			***REMOVED***

		case ntBol, ntEol, ntBoundary, ntECMABoundary, ntBeginning, ntStart,
			ntEndZ, ntEnd, ntEmpty, ntRequire, ntPrevent:

		default:
			return nil
		***REMOVED***

		if concatNode == nil || nextChild >= len(concatNode.children) ***REMOVED***
			return nil
		***REMOVED***

		curNode = concatNode.children[nextChild]
		nextChild++
	***REMOVED***
***REMOVED***

// repeat the rune r, c times... up to the max of MaxPrefixSize
func repeat(r rune, c int) []rune ***REMOVED***
	if c > MaxPrefixSize ***REMOVED***
		c = MaxPrefixSize
	***REMOVED***

	ret := make([]rune, c)

	// binary growth using copy for speed
	ret[0] = r
	bp := 1
	for bp < len(ret) ***REMOVED***
		copy(ret[bp:], ret[:bp])
		bp *= 2
	***REMOVED***

	return ret
***REMOVED***

// BmPrefix precomputes the Boyer-Moore
// tables for fast string scanning. These tables allow
// you to scan for the first occurrence of a string within
// a large body of text without examining every character.
// The performance of the heuristic depends on the actual
// string and the text being searched, but usually, the longer
// the string that is being searched for, the fewer characters
// need to be examined.
type BmPrefix struct ***REMOVED***
	positive        []int
	negativeASCII   []int
	negativeUnicode [][]int
	pattern         []rune
	lowASCII        rune
	highASCII       rune
	rightToLeft     bool
	caseInsensitive bool
***REMOVED***

func newBmPrefix(pattern []rune, caseInsensitive, rightToLeft bool) *BmPrefix ***REMOVED***

	b := &BmPrefix***REMOVED***
		rightToLeft:     rightToLeft,
		caseInsensitive: caseInsensitive,
		pattern:         pattern,
	***REMOVED***

	if caseInsensitive ***REMOVED***
		for i := 0; i < len(b.pattern); i++ ***REMOVED***
			// We do the ToLower character by character for consistency.  With surrogate chars, doing
			// a ToLower on the entire string could actually change the surrogate pair.  This is more correct
			// linguistically, but since Regex doesn't support surrogates, it's more important to be
			// consistent.

			b.pattern[i] = unicode.ToLower(b.pattern[i])
		***REMOVED***
	***REMOVED***

	var beforefirst, last, bump int
	var scan, match int

	if !rightToLeft ***REMOVED***
		beforefirst = -1
		last = len(b.pattern) - 1
		bump = 1
	***REMOVED*** else ***REMOVED***
		beforefirst = len(b.pattern)
		last = 0
		bump = -1
	***REMOVED***

	// PART I - the good-suffix shift table
	//
	// compute the positive requirement:
	// if char "i" is the first one from the right that doesn't match,
	// then we know the matcher can advance by _positive[i].
	//
	// This algorithm is a simplified variant of the standard
	// Boyer-Moore good suffix calculation.

	b.positive = make([]int, len(b.pattern))

	examine := last
	ch := b.pattern[examine]
	b.positive[examine] = bump
	examine -= bump

Outerloop:
	for ***REMOVED***
		// find an internal char (examine) that matches the tail

		for ***REMOVED***
			if examine == beforefirst ***REMOVED***
				break Outerloop
			***REMOVED***
			if b.pattern[examine] == ch ***REMOVED***
				break
			***REMOVED***
			examine -= bump
		***REMOVED***

		match = last
		scan = examine

		// find the length of the match
		for ***REMOVED***
			if scan == beforefirst || b.pattern[match] != b.pattern[scan] ***REMOVED***
				// at the end of the match, note the difference in _positive
				// this is not the length of the match, but the distance from the internal match
				// to the tail suffix.
				if b.positive[match] == 0 ***REMOVED***
					b.positive[match] = match - scan
				***REMOVED***

				// System.Diagnostics.Debug.WriteLine("Set positive[" + match + "] to " + (match - scan));

				break
			***REMOVED***

			scan -= bump
			match -= bump
		***REMOVED***

		examine -= bump
	***REMOVED***

	match = last - bump

	// scan for the chars for which there are no shifts that yield a different candidate

	// The inside of the if statement used to say
	// "_positive[match] = last - beforefirst;"
	// This is slightly less aggressive in how much we skip, but at worst it
	// should mean a little more work rather than skipping a potential match.
	for match != beforefirst ***REMOVED***
		if b.positive[match] == 0 ***REMOVED***
			b.positive[match] = bump
		***REMOVED***

		match -= bump
	***REMOVED***

	// PART II - the bad-character shift table
	//
	// compute the negative requirement:
	// if char "ch" is the reject character when testing position "i",
	// we can slide up by _negative[ch];
	// (_negative[ch] = str.Length - 1 - str.LastIndexOf(ch))
	//
	// the lookup table is divided into ASCII and Unicode portions;
	// only those parts of the Unicode 16-bit code set that actually
	// appear in the string are in the table. (Maximum size with
	// Unicode is 65K; ASCII only case is 512 bytes.)

	b.negativeASCII = make([]int, 128)

	for i := 0; i < len(b.negativeASCII); i++ ***REMOVED***
		b.negativeASCII[i] = last - beforefirst
	***REMOVED***

	b.lowASCII = 127
	b.highASCII = 0

	for examine = last; examine != beforefirst; examine -= bump ***REMOVED***
		ch = b.pattern[examine]

		switch ***REMOVED***
		case ch < 128:
			if b.lowASCII > ch ***REMOVED***
				b.lowASCII = ch
			***REMOVED***

			if b.highASCII < ch ***REMOVED***
				b.highASCII = ch
			***REMOVED***

			if b.negativeASCII[ch] == last-beforefirst ***REMOVED***
				b.negativeASCII[ch] = last - examine
			***REMOVED***
		case ch <= 0xffff:
			i, j := ch>>8, ch&0xFF

			if b.negativeUnicode == nil ***REMOVED***
				b.negativeUnicode = make([][]int, 256)
			***REMOVED***

			if b.negativeUnicode[i] == nil ***REMOVED***
				newarray := make([]int, 256)

				for k := 0; k < len(newarray); k++ ***REMOVED***
					newarray[k] = last - beforefirst
				***REMOVED***

				if i == 0 ***REMOVED***
					copy(newarray, b.negativeASCII)
					//TODO: this line needed?
					b.negativeASCII = newarray
				***REMOVED***

				b.negativeUnicode[i] = newarray
			***REMOVED***

			if b.negativeUnicode[i][j] == last-beforefirst ***REMOVED***
				b.negativeUnicode[i][j] = last - examine
			***REMOVED***
		default:
			// we can't do the filter because this algo doesn't support
			// unicode chars >0xffff
			return nil
		***REMOVED***
	***REMOVED***

	return b
***REMOVED***

func (b *BmPrefix) String() string ***REMOVED***
	return string(b.pattern)
***REMOVED***

// Dump returns the contents of the filter as a human readable string
func (b *BmPrefix) Dump(indent string) string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***

	fmt.Fprintf(buf, "%sBM Pattern: %s\n%sPositive: ", indent, string(b.pattern), indent)
	for i := 0; i < len(b.positive); i++ ***REMOVED***
		buf.WriteString(strconv.Itoa(b.positive[i]))
		buf.WriteRune(' ')
	***REMOVED***
	buf.WriteRune('\n')

	if b.negativeASCII != nil ***REMOVED***
		buf.WriteString(indent)
		buf.WriteString("Negative table\n")
		for i := 0; i < len(b.negativeASCII); i++ ***REMOVED***
			if b.negativeASCII[i] != len(b.pattern) ***REMOVED***
				fmt.Fprintf(buf, "%s  %s %s\n", indent, Escape(string(rune(i))), strconv.Itoa(b.negativeASCII[i]))
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return buf.String()
***REMOVED***

// Scan uses the Boyer-Moore algorithm to find the first occurrence
// of the specified string within text, beginning at index, and
// constrained within beglimit and endlimit.
//
// The direction and case-sensitivity of the match is determined
// by the arguments to the RegexBoyerMoore constructor.
func (b *BmPrefix) Scan(text []rune, index, beglimit, endlimit int) int ***REMOVED***
	var (
		defadv, test, test2         int
		match, startmatch, endmatch int
		bump, advance               int
		chTest                      rune
		unicodeLookup               []int
	)

	if !b.rightToLeft ***REMOVED***
		defadv = len(b.pattern)
		startmatch = len(b.pattern) - 1
		endmatch = 0
		test = index + defadv - 1
		bump = 1
	***REMOVED*** else ***REMOVED***
		defadv = -len(b.pattern)
		startmatch = 0
		endmatch = -defadv - 1
		test = index + defadv
		bump = -1
	***REMOVED***

	chMatch := b.pattern[startmatch]

	for ***REMOVED***
		if test >= endlimit || test < beglimit ***REMOVED***
			return -1
		***REMOVED***

		chTest = text[test]

		if b.caseInsensitive ***REMOVED***
			chTest = unicode.ToLower(chTest)
		***REMOVED***

		if chTest != chMatch ***REMOVED***
			if chTest < 128 ***REMOVED***
				advance = b.negativeASCII[chTest]
			***REMOVED*** else if chTest < 0xffff && len(b.negativeUnicode) > 0 ***REMOVED***
				unicodeLookup = b.negativeUnicode[chTest>>8]
				if len(unicodeLookup) > 0 ***REMOVED***
					advance = unicodeLookup[chTest&0xFF]
				***REMOVED*** else ***REMOVED***
					advance = defadv
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				advance = defadv
			***REMOVED***

			test += advance
		***REMOVED*** else ***REMOVED*** // if (chTest == chMatch)
			test2 = test
			match = startmatch

			for ***REMOVED***
				if match == endmatch ***REMOVED***
					if b.rightToLeft ***REMOVED***
						return test2 + 1
					***REMOVED*** else ***REMOVED***
						return test2
					***REMOVED***
				***REMOVED***

				match -= bump
				test2 -= bump

				chTest = text[test2]

				if b.caseInsensitive ***REMOVED***
					chTest = unicode.ToLower(chTest)
				***REMOVED***

				if chTest != b.pattern[match] ***REMOVED***
					advance = b.positive[match]
					if chTest < 128 ***REMOVED***
						test2 = (match - startmatch) + b.negativeASCII[chTest]
					***REMOVED*** else if chTest < 0xffff && len(b.negativeUnicode) > 0 ***REMOVED***
						unicodeLookup = b.negativeUnicode[chTest>>8]
						if len(unicodeLookup) > 0 ***REMOVED***
							test2 = (match - startmatch) + unicodeLookup[chTest&0xFF]
						***REMOVED*** else ***REMOVED***
							test += advance
							break
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						test += advance
						break
					***REMOVED***

					if b.rightToLeft ***REMOVED***
						if test2 < advance ***REMOVED***
							advance = test2
						***REMOVED***
					***REMOVED*** else if test2 > advance ***REMOVED***
						advance = test2
					***REMOVED***

					test += advance
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// When a regex is anchored, we can do a quick IsMatch test instead of a Scan
func (b *BmPrefix) IsMatch(text []rune, index, beglimit, endlimit int) bool ***REMOVED***
	if !b.rightToLeft ***REMOVED***
		if index < beglimit || endlimit-index < len(b.pattern) ***REMOVED***
			return false
		***REMOVED***

		return b.matchPattern(text, index)
	***REMOVED*** else ***REMOVED***
		if index > endlimit || index-beglimit < len(b.pattern) ***REMOVED***
			return false
		***REMOVED***

		return b.matchPattern(text, index-len(b.pattern))
	***REMOVED***
***REMOVED***

func (b *BmPrefix) matchPattern(text []rune, index int) bool ***REMOVED***
	if len(text)-index < len(b.pattern) ***REMOVED***
		return false
	***REMOVED***

	if b.caseInsensitive ***REMOVED***
		for i := 0; i < len(b.pattern); i++ ***REMOVED***
			//Debug.Assert(textinfo.ToLower(_pattern[i]) == _pattern[i], "pattern should be converted to lower case in constructor!");
			if unicode.ToLower(text[index+i]) != b.pattern[i] ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED*** else ***REMOVED***
		for i := 0; i < len(b.pattern); i++ ***REMOVED***
			if text[index+i] != b.pattern[i] ***REMOVED***
				return false
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
***REMOVED***

type AnchorLoc int16

// where the regex can be pegged
const (
	AnchorBeginning    AnchorLoc = 0x0001
	AnchorBol                    = 0x0002
	AnchorStart                  = 0x0004
	AnchorEol                    = 0x0008
	AnchorEndZ                   = 0x0010
	AnchorEnd                    = 0x0020
	AnchorBoundary               = 0x0040
	AnchorECMABoundary           = 0x0080
)

func getAnchors(tree *RegexTree) AnchorLoc ***REMOVED***

	var concatNode *regexNode
	nextChild, result := 0, AnchorLoc(0)

	curNode := tree.root

	for ***REMOVED***
		switch curNode.t ***REMOVED***
		case ntConcatenate:
			if len(curNode.children) > 0 ***REMOVED***
				concatNode = curNode
				nextChild = 0
			***REMOVED***

		case ntGreedy, ntCapture:
			curNode = curNode.children[0]
			concatNode = nil
			continue

		case ntBol, ntEol, ntBoundary, ntECMABoundary, ntBeginning,
			ntStart, ntEndZ, ntEnd:
			return result | anchorFromType(curNode.t)

		case ntEmpty, ntRequire, ntPrevent:

		default:
			return result
		***REMOVED***

		if concatNode == nil || nextChild >= len(concatNode.children) ***REMOVED***
			return result
		***REMOVED***

		curNode = concatNode.children[nextChild]
		nextChild++
	***REMOVED***
***REMOVED***

func anchorFromType(t nodeType) AnchorLoc ***REMOVED***
	switch t ***REMOVED***
	case ntBol:
		return AnchorBol
	case ntEol:
		return AnchorEol
	case ntBoundary:
		return AnchorBoundary
	case ntECMABoundary:
		return AnchorECMABoundary
	case ntBeginning:
		return AnchorBeginning
	case ntStart:
		return AnchorStart
	case ntEndZ:
		return AnchorEndZ
	case ntEnd:
		return AnchorEnd
	default:
		return 0
	***REMOVED***
***REMOVED***

// anchorDescription returns a human-readable description of the anchors
func (anchors AnchorLoc) String() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***

	if 0 != (anchors & AnchorBeginning) ***REMOVED***
		buf.WriteString(", Beginning")
	***REMOVED***
	if 0 != (anchors & AnchorStart) ***REMOVED***
		buf.WriteString(", Start")
	***REMOVED***
	if 0 != (anchors & AnchorBol) ***REMOVED***
		buf.WriteString(", Bol")
	***REMOVED***
	if 0 != (anchors & AnchorBoundary) ***REMOVED***
		buf.WriteString(", Boundary")
	***REMOVED***
	if 0 != (anchors & AnchorECMABoundary) ***REMOVED***
		buf.WriteString(", ECMABoundary")
	***REMOVED***
	if 0 != (anchors & AnchorEol) ***REMOVED***
		buf.WriteString(", Eol")
	***REMOVED***
	if 0 != (anchors & AnchorEnd) ***REMOVED***
		buf.WriteString(", End")
	***REMOVED***
	if 0 != (anchors & AnchorEndZ) ***REMOVED***
		buf.WriteString(", EndZ")
	***REMOVED***

	// trim off comma
	if buf.Len() >= 2 ***REMOVED***
		return buf.String()[2:]
	***REMOVED***
	return "None"
***REMOVED***
