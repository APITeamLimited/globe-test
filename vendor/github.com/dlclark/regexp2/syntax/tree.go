package syntax

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
)

type RegexTree struct ***REMOVED***
	root       *regexNode
	caps       map[int]int
	capnumlist []int
	captop     int
	Capnames   map[string]int
	Caplist    []string
	options    RegexOptions
***REMOVED***

// It is built into a parsed tree for a regular expression.

// Implementation notes:
//
// Since the node tree is a temporary data structure only used
// during compilation of the regexp to integer codes, it's
// designed for clarity and convenience rather than
// space efficiency.
//
// RegexNodes are built into a tree, linked by the n.children list.
// Each node also has a n.parent and n.ichild member indicating
// its parent and which child # it is in its parent's list.
//
// RegexNodes come in as many types as there are constructs in
// a regular expression, for example, "concatenate", "alternate",
// "one", "rept", "group". There are also node types for basic
// peephole optimizations, e.g., "onerep", "notsetrep", etc.
//
// Because perl 5 allows "lookback" groups that scan backwards,
// each node also gets a "direction". Normally the value of
// boolean n.backward = false.
//
// During parsing, top-level nodes are also stacked onto a parse
// stack (a stack of trees). For this purpose we have a n.next
// pointer. [Note that to save a few bytes, we could overload the
// n.parent pointer instead.]
//
// On the parse stack, each tree has a "role" - basically, the
// nonterminal in the grammar that the parser has currently
// assigned to the tree. That code is stored in n.role.
//
// Finally, some of the different kinds of nodes have data.
// Two integers (for the looping constructs) are stored in
// n.operands, an an object (either a string or a set)
// is stored in n.data
type regexNode struct ***REMOVED***
	t        nodeType
	children []*regexNode
	str      []rune
	set      *CharSet
	ch       rune
	m        int
	n        int
	options  RegexOptions
	next     *regexNode
***REMOVED***

type nodeType int32

const (
	// The following are leaves, and correspond to primitive operations

	ntOnerep      nodeType = 0  // lef,back char,min,max    a ***REMOVED***n***REMOVED***
	ntNotonerep            = 1  // lef,back char,min,max    .***REMOVED***n***REMOVED***
	ntSetrep               = 2  // lef,back set,min,max     [\d]***REMOVED***n***REMOVED***
	ntOneloop              = 3  // lef,back char,min,max    a ***REMOVED***,n***REMOVED***
	ntNotoneloop           = 4  // lef,back char,min,max    .***REMOVED***,n***REMOVED***
	ntSetloop              = 5  // lef,back set,min,max     [\d]***REMOVED***,n***REMOVED***
	ntOnelazy              = 6  // lef,back char,min,max    a ***REMOVED***,n***REMOVED***?
	ntNotonelazy           = 7  // lef,back char,min,max    .***REMOVED***,n***REMOVED***?
	ntSetlazy              = 8  // lef,back set,min,max     [\d]***REMOVED***,n***REMOVED***?
	ntOne                  = 9  // lef      char            a
	ntNotone               = 10 // lef      char            [^a]
	ntSet                  = 11 // lef      set             [a-z\s]  \w \s \d
	ntMulti                = 12 // lef      string          abcd
	ntRef                  = 13 // lef      group           \#
	ntBol                  = 14 //                          ^
	ntEol                  = 15 //                          $
	ntBoundary             = 16 //                          \b
	ntNonboundary          = 17 //                          \B
	ntBeginning            = 18 //                          \A
	ntStart                = 19 //                          \G
	ntEndZ                 = 20 //                          \Z
	ntEnd                  = 21 //                          \Z

	// Interior nodes do not correspond to primitive operations, but
	// control structures compositing other operations

	// Concat and alternate take n children, and can run forward or backwards

	ntNothing     = 22 //          []
	ntEmpty       = 23 //          ()
	ntAlternate   = 24 //          a|b
	ntConcatenate = 25 //          ab
	ntLoop        = 26 // m,x      * + ? ***REMOVED***,***REMOVED***
	ntLazyloop    = 27 // m,x      *? +? ?? ***REMOVED***,***REMOVED***?
	ntCapture     = 28 // n        ()
	ntGroup       = 29 //          (?:)
	ntRequire     = 30 //          (?=) (?<=)
	ntPrevent     = 31 //          (?!) (?<!)
	ntGreedy      = 32 //          (?>) (?<)
	ntTestref     = 33 //          (?(n) | )
	ntTestgroup   = 34 //          (?(...) | )

	ntECMABoundary    = 41 //                          \b
	ntNonECMABoundary = 42 //                          \B
)

func newRegexNode(t nodeType, opt RegexOptions) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
	***REMOVED***
***REMOVED***

func newRegexNodeCh(t nodeType, opt RegexOptions, ch rune) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
		ch:      ch,
	***REMOVED***
***REMOVED***

func newRegexNodeStr(t nodeType, opt RegexOptions, str []rune) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
		str:     str,
	***REMOVED***
***REMOVED***

func newRegexNodeSet(t nodeType, opt RegexOptions, set *CharSet) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
		set:     set,
	***REMOVED***
***REMOVED***

func newRegexNodeM(t nodeType, opt RegexOptions, m int) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
		m:       m,
	***REMOVED***
***REMOVED***
func newRegexNodeMN(t nodeType, opt RegexOptions, m, n int) *regexNode ***REMOVED***
	return &regexNode***REMOVED***
		t:       t,
		options: opt,
		m:       m,
		n:       n,
	***REMOVED***
***REMOVED***

func (n *regexNode) writeStrToBuf(buf *bytes.Buffer) ***REMOVED***
	for i := 0; i < len(n.str); i++ ***REMOVED***
		buf.WriteRune(n.str[i])
	***REMOVED***
***REMOVED***

func (n *regexNode) addChild(child *regexNode) ***REMOVED***
	reduced := child.reduce()
	n.children = append(n.children, reduced)
	reduced.next = n
***REMOVED***

func (n *regexNode) insertChildren(afterIndex int, nodes []*regexNode) ***REMOVED***
	newChildren := make([]*regexNode, 0, len(n.children)+len(nodes))
	n.children = append(append(append(newChildren, n.children[:afterIndex]...), nodes...), n.children[afterIndex:]...)
***REMOVED***

// removes children including the start but not the end index
func (n *regexNode) removeChildren(startIndex, endIndex int) ***REMOVED***
	n.children = append(n.children[:startIndex], n.children[endIndex:]...)
***REMOVED***

// Pass type as OneLazy or OneLoop
func (n *regexNode) makeRep(t nodeType, min, max int) ***REMOVED***
	n.t += (t - ntOne)
	n.m = min
	n.n = max
***REMOVED***

func (n *regexNode) reduce() *regexNode ***REMOVED***
	switch n.t ***REMOVED***
	case ntAlternate:
		return n.reduceAlternation()

	case ntConcatenate:
		return n.reduceConcatenation()

	case ntLoop, ntLazyloop:
		return n.reduceRep()

	case ntGroup:
		return n.reduceGroup()

	case ntSet, ntSetloop:
		return n.reduceSet()

	default:
		return n
	***REMOVED***
***REMOVED***

// Basic optimization. Single-letter alternations can be replaced
// by faster set specifications, and nested alternations with no
// intervening operators can be flattened:
//
// a|b|c|def|g|h -> [a-c]|def|[gh]
// apple|(?:orange|pear)|grape -> apple|orange|pear|grape
func (n *regexNode) reduceAlternation() *regexNode ***REMOVED***
	if len(n.children) == 0 ***REMOVED***
		return newRegexNode(ntNothing, n.options)
	***REMOVED***

	wasLastSet := false
	lastNodeCannotMerge := false
	var optionsLast RegexOptions
	var i, j int

	for i, j = 0, 0; i < len(n.children); i, j = i+1, j+1 ***REMOVED***
		at := n.children[i]

		if j < i ***REMOVED***
			n.children[j] = at
		***REMOVED***

		for ***REMOVED***
			if at.t == ntAlternate ***REMOVED***
				for k := 0; k < len(at.children); k++ ***REMOVED***
					at.children[k].next = n
				***REMOVED***
				n.insertChildren(i+1, at.children)

				j--
			***REMOVED*** else if at.t == ntSet || at.t == ntOne ***REMOVED***
				// Cannot merge sets if L or I options differ, or if either are negated.
				optionsAt := at.options & (RightToLeft | IgnoreCase)

				if at.t == ntSet ***REMOVED***
					if !wasLastSet || optionsLast != optionsAt || lastNodeCannotMerge || !at.set.IsMergeable() ***REMOVED***
						wasLastSet = true
						lastNodeCannotMerge = !at.set.IsMergeable()
						optionsLast = optionsAt
						break
					***REMOVED***
				***REMOVED*** else if !wasLastSet || optionsLast != optionsAt || lastNodeCannotMerge ***REMOVED***
					wasLastSet = true
					lastNodeCannotMerge = false
					optionsLast = optionsAt
					break
				***REMOVED***

				// The last node was a Set or a One, we're a Set or One and our options are the same.
				// Merge the two nodes.
				j--
				prev := n.children[j]

				var prevCharClass *CharSet
				if prev.t == ntOne ***REMOVED***
					prevCharClass = &CharSet***REMOVED******REMOVED***
					prevCharClass.addChar(prev.ch)
				***REMOVED*** else ***REMOVED***
					prevCharClass = prev.set
				***REMOVED***

				if at.t == ntOne ***REMOVED***
					prevCharClass.addChar(at.ch)
				***REMOVED*** else ***REMOVED***
					prevCharClass.addSet(*at.set)
				***REMOVED***

				prev.t = ntSet
				prev.set = prevCharClass
			***REMOVED*** else if at.t == ntNothing ***REMOVED***
				j--
			***REMOVED*** else ***REMOVED***
				wasLastSet = false
				lastNodeCannotMerge = false
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if j < i ***REMOVED***
		n.removeChildren(j, i)
	***REMOVED***

	return n.stripEnation(ntNothing)
***REMOVED***

// Basic optimization. Adjacent strings can be concatenated.
//
// (?:abc)(?:def) -> abcdef
func (n *regexNode) reduceConcatenation() *regexNode ***REMOVED***
	// Eliminate empties and concat adjacent strings/chars

	var optionsLast RegexOptions
	var optionsAt RegexOptions
	var i, j int

	if len(n.children) == 0 ***REMOVED***
		return newRegexNode(ntEmpty, n.options)
	***REMOVED***

	wasLastString := false

	for i, j = 0, 0; i < len(n.children); i, j = i+1, j+1 ***REMOVED***
		var at, prev *regexNode

		at = n.children[i]

		if j < i ***REMOVED***
			n.children[j] = at
		***REMOVED***

		if at.t == ntConcatenate &&
			((at.options & RightToLeft) == (n.options & RightToLeft)) ***REMOVED***
			for k := 0; k < len(at.children); k++ ***REMOVED***
				at.children[k].next = n
			***REMOVED***

			//insert at.children at i+1 index in n.children
			n.insertChildren(i+1, at.children)

			j--
		***REMOVED*** else if at.t == ntMulti || at.t == ntOne ***REMOVED***
			// Cannot merge strings if L or I options differ
			optionsAt = at.options & (RightToLeft | IgnoreCase)

			if !wasLastString || optionsLast != optionsAt ***REMOVED***
				wasLastString = true
				optionsLast = optionsAt
				continue
			***REMOVED***

			j--
			prev = n.children[j]

			if prev.t == ntOne ***REMOVED***
				prev.t = ntMulti
				prev.str = []rune***REMOVED***prev.ch***REMOVED***
			***REMOVED***

			if (optionsAt & RightToLeft) == 0 ***REMOVED***
				if at.t == ntOne ***REMOVED***
					prev.str = append(prev.str, at.ch)
				***REMOVED*** else ***REMOVED***
					prev.str = append(prev.str, at.str...)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if at.t == ntOne ***REMOVED***
					// insert at the front by expanding our slice, copying the data over, and then setting the value
					prev.str = append(prev.str, 0)
					copy(prev.str[1:], prev.str)
					prev.str[0] = at.ch
				***REMOVED*** else ***REMOVED***
					//insert at the front...this one we'll make a new slice and copy both into it
					merge := make([]rune, len(prev.str)+len(at.str))
					copy(merge, at.str)
					copy(merge[len(at.str):], prev.str)
					prev.str = merge
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if at.t == ntEmpty ***REMOVED***
			j--
		***REMOVED*** else ***REMOVED***
			wasLastString = false
		***REMOVED***
	***REMOVED***

	if j < i ***REMOVED***
		// remove indices j through i from the children
		n.removeChildren(j, i)
	***REMOVED***

	return n.stripEnation(ntEmpty)
***REMOVED***

// Nested repeaters just get multiplied with each other if they're not
// too lumpy
func (n *regexNode) reduceRep() *regexNode ***REMOVED***

	u := n
	t := n.t
	min := n.m
	max := n.n

	for ***REMOVED***
		if len(u.children) == 0 ***REMOVED***
			break
		***REMOVED***

		child := u.children[0]

		// multiply reps of the same type only
		if child.t != t ***REMOVED***
			childType := child.t

			if !(childType >= ntOneloop && childType <= ntSetloop && t == ntLoop ||
				childType >= ntOnelazy && childType <= ntSetlazy && t == ntLazyloop) ***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		// child can be too lumpy to blur, e.g., (a ***REMOVED***100,105***REMOVED***) ***REMOVED***3***REMOVED*** or (a ***REMOVED***2,***REMOVED***)?
		// [but things like (a ***REMOVED***2,***REMOVED***)+ are not too lumpy...]
		if u.m == 0 && child.m > 1 || child.n < child.m*2 ***REMOVED***
			break
		***REMOVED***

		u = child
		if u.m > 0 ***REMOVED***
			if (math.MaxInt32-1)/u.m < min ***REMOVED***
				u.m = math.MaxInt32
			***REMOVED*** else ***REMOVED***
				u.m = u.m * min
			***REMOVED***
		***REMOVED***
		if u.n > 0 ***REMOVED***
			if (math.MaxInt32-1)/u.n < max ***REMOVED***
				u.n = math.MaxInt32
			***REMOVED*** else ***REMOVED***
				u.n = u.n * max
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if math.MaxInt32 == min ***REMOVED***
		return newRegexNode(ntNothing, n.options)
	***REMOVED***
	return u

***REMOVED***

// Simple optimization. If a concatenation or alternation has only
// one child strip out the intermediate node. If it has zero children,
// turn it into an empty.
func (n *regexNode) stripEnation(emptyType nodeType) *regexNode ***REMOVED***
	switch len(n.children) ***REMOVED***
	case 0:
		return newRegexNode(emptyType, n.options)
	case 1:
		return n.children[0]
	default:
		return n
	***REMOVED***
***REMOVED***

func (n *regexNode) reduceGroup() *regexNode ***REMOVED***
	u := n

	for u.t == ntGroup ***REMOVED***
		u = u.children[0]
	***REMOVED***

	return u
***REMOVED***

// Simple optimization. If a set is a singleton, an inverse singleton,
// or empty, it's transformed accordingly.
func (n *regexNode) reduceSet() *regexNode ***REMOVED***
	// Extract empty-set, one and not-one case as special

	if n.set == nil ***REMOVED***
		n.t = ntNothing
	***REMOVED*** else if n.set.IsSingleton() ***REMOVED***
		n.ch = n.set.SingletonChar()
		n.set = nil
		n.t += (ntOne - ntSet)
	***REMOVED*** else if n.set.IsSingletonInverse() ***REMOVED***
		n.ch = n.set.SingletonChar()
		n.set = nil
		n.t += (ntNotone - ntSet)
	***REMOVED***

	return n
***REMOVED***

func (n *regexNode) reverseLeft() *regexNode ***REMOVED***
	if n.options&RightToLeft != 0 && n.t == ntConcatenate && len(n.children) > 0 ***REMOVED***
		//reverse children order
		for left, right := 0, len(n.children)-1; left < right; left, right = left+1, right-1 ***REMOVED***
			n.children[left], n.children[right] = n.children[right], n.children[left]
		***REMOVED***
	***REMOVED***

	return n
***REMOVED***

func (n *regexNode) makeQuantifier(lazy bool, min, max int) *regexNode ***REMOVED***
	if min == 0 && max == 0 ***REMOVED***
		return newRegexNode(ntEmpty, n.options)
	***REMOVED***

	if min == 1 && max == 1 ***REMOVED***
		return n
	***REMOVED***

	switch n.t ***REMOVED***
	case ntOne, ntNotone, ntSet:
		if lazy ***REMOVED***
			n.makeRep(Onelazy, min, max)
		***REMOVED*** else ***REMOVED***
			n.makeRep(Oneloop, min, max)
		***REMOVED***
		return n

	default:
		var t nodeType
		if lazy ***REMOVED***
			t = ntLazyloop
		***REMOVED*** else ***REMOVED***
			t = ntLoop
		***REMOVED***
		result := newRegexNodeMN(t, n.options, min, max)
		result.addChild(n)
		return result
	***REMOVED***
***REMOVED***

// debug functions

var typeStr = []string***REMOVED***
	"Onerep", "Notonerep", "Setrep",
	"Oneloop", "Notoneloop", "Setloop",
	"Onelazy", "Notonelazy", "Setlazy",
	"One", "Notone", "Set",
	"Multi", "Ref",
	"Bol", "Eol", "Boundary", "Nonboundary",
	"Beginning", "Start", "EndZ", "End",
	"Nothing", "Empty",
	"Alternate", "Concatenate",
	"Loop", "Lazyloop",
	"Capture", "Group", "Require", "Prevent", "Greedy",
	"Testref", "Testgroup",
	"Unknown", "Unknown", "Unknown",
	"Unknown", "Unknown", "Unknown",
	"ECMABoundary", "NonECMABoundary",
***REMOVED***

func (n *regexNode) description() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***

	buf.WriteString(typeStr[n.t])

	if (n.options & ExplicitCapture) != 0 ***REMOVED***
		buf.WriteString("-C")
	***REMOVED***
	if (n.options & IgnoreCase) != 0 ***REMOVED***
		buf.WriteString("-I")
	***REMOVED***
	if (n.options & RightToLeft) != 0 ***REMOVED***
		buf.WriteString("-L")
	***REMOVED***
	if (n.options & Multiline) != 0 ***REMOVED***
		buf.WriteString("-M")
	***REMOVED***
	if (n.options & Singleline) != 0 ***REMOVED***
		buf.WriteString("-S")
	***REMOVED***
	if (n.options & IgnorePatternWhitespace) != 0 ***REMOVED***
		buf.WriteString("-X")
	***REMOVED***
	if (n.options & ECMAScript) != 0 ***REMOVED***
		buf.WriteString("-E")
	***REMOVED***

	switch n.t ***REMOVED***
	case ntOneloop, ntNotoneloop, ntOnelazy, ntNotonelazy, ntOne, ntNotone:
		buf.WriteString("(Ch = " + CharDescription(n.ch) + ")")
		break
	case ntCapture:
		buf.WriteString("(index = " + strconv.Itoa(n.m) + ", unindex = " + strconv.Itoa(n.n) + ")")
		break
	case ntRef, ntTestref:
		buf.WriteString("(index = " + strconv.Itoa(n.m) + ")")
		break
	case ntMulti:
		fmt.Fprintf(buf, "(String = %s)", string(n.str))
		break
	case ntSet, ntSetloop, ntSetlazy:
		buf.WriteString("(Set = " + n.set.String() + ")")
		break
	***REMOVED***

	switch n.t ***REMOVED***
	case ntOneloop, ntNotoneloop, ntOnelazy, ntNotonelazy, ntSetloop, ntSetlazy, ntLoop, ntLazyloop:
		buf.WriteString("(Min = ")
		buf.WriteString(strconv.Itoa(n.m))
		buf.WriteString(", Max = ")
		if n.n == math.MaxInt32 ***REMOVED***
			buf.WriteString("inf")
		***REMOVED*** else ***REMOVED***
			buf.WriteString(strconv.Itoa(n.n))
		***REMOVED***
		buf.WriteString(")")

		break
	***REMOVED***

	return buf.String()
***REMOVED***

var padSpace = []byte("                                ")

func (t *RegexTree) Dump() string ***REMOVED***
	return t.root.dump()
***REMOVED***

func (n *regexNode) dump() string ***REMOVED***
	var stack []int
	CurNode := n
	CurChild := 0

	buf := bytes.NewBufferString(CurNode.description())
	buf.WriteRune('\n')

	for ***REMOVED***
		if CurNode.children != nil && CurChild < len(CurNode.children) ***REMOVED***
			stack = append(stack, CurChild+1)
			CurNode = CurNode.children[CurChild]
			CurChild = 0

			Depth := len(stack)
			if Depth > 32 ***REMOVED***
				Depth = 32
			***REMOVED***
			buf.Write(padSpace[:Depth])
			buf.WriteString(CurNode.description())
			buf.WriteRune('\n')
		***REMOVED*** else ***REMOVED***
			if len(stack) == 0 ***REMOVED***
				break
			***REMOVED***

			CurChild = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			CurNode = CurNode.next
		***REMOVED***
	***REMOVED***
	return buf.String()
***REMOVED***
