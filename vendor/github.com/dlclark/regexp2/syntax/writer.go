package syntax

import (
	"bytes"
	"fmt"
	"math"
	"os"
)

func Write(tree *RegexTree) (*Code, error) ***REMOVED***
	w := writer***REMOVED***
		intStack:   make([]int, 0, 32),
		emitted:    make([]int, 2),
		stringhash: make(map[string]int),
		sethash:    make(map[string]int),
	***REMOVED***

	code, err := w.codeFromTree(tree)

	if tree.options&Debug > 0 && code != nil ***REMOVED***
		os.Stdout.WriteString(code.Dump())
		os.Stdout.WriteString("\n")
	***REMOVED***

	return code, err
***REMOVED***

type writer struct ***REMOVED***
	emitted []int

	intStack    []int
	curpos      int
	stringhash  map[string]int
	stringtable [][]rune
	sethash     map[string]int
	settable    []*CharSet
	counting    bool
	count       int
	trackcount  int
	caps        map[int]int
***REMOVED***

const (
	beforeChild nodeType = 64
	afterChild           = 128
	//MaxPrefixSize is the largest number of runes we'll use for a BoyerMoyer prefix
	MaxPrefixSize = 50
)

// The top level RegexCode generator. It does a depth-first walk
// through the tree and calls EmitFragment to emits code before
// and after each child of an interior node, and at each leaf.
//
// It runs two passes, first to count the size of the generated
// code, and second to generate the code.
//
// We should time it against the alternative, which is
// to just generate the code and grow the array as we go.
func (w *writer) codeFromTree(tree *RegexTree) (*Code, error) ***REMOVED***
	var (
		curNode  *regexNode
		curChild int
		capsize  int
	)
	// construct sparse capnum mapping if some numbers are unused

	if tree.capnumlist == nil || tree.captop == len(tree.capnumlist) ***REMOVED***
		capsize = tree.captop
		w.caps = nil
	***REMOVED*** else ***REMOVED***
		capsize = len(tree.capnumlist)
		w.caps = tree.caps
		for i := 0; i < len(tree.capnumlist); i++ ***REMOVED***
			w.caps[tree.capnumlist[i]] = i
		***REMOVED***
	***REMOVED***

	w.counting = true

	for ***REMOVED***
		if !w.counting ***REMOVED***
			w.emitted = make([]int, w.count)
		***REMOVED***

		curNode = tree.root
		curChild = 0

		w.emit1(Lazybranch, 0)

		for ***REMOVED***
			if len(curNode.children) == 0 ***REMOVED***
				w.emitFragment(curNode.t, curNode, 0)
			***REMOVED*** else if curChild < len(curNode.children) ***REMOVED***
				w.emitFragment(curNode.t|beforeChild, curNode, curChild)

				curNode = curNode.children[curChild]

				w.pushInt(curChild)
				curChild = 0
				continue
			***REMOVED***

			if w.emptyStack() ***REMOVED***
				break
			***REMOVED***

			curChild = w.popInt()
			curNode = curNode.next

			w.emitFragment(curNode.t|afterChild, curNode, curChild)
			curChild++
		***REMOVED***

		w.patchJump(0, w.curPos())
		w.emit(Stop)

		if !w.counting ***REMOVED***
			break
		***REMOVED***

		w.counting = false
	***REMOVED***

	fcPrefix := getFirstCharsPrefix(tree)
	prefix := getPrefix(tree)
	rtl := (tree.options & RightToLeft) != 0

	var bmPrefix *BmPrefix
	//TODO: benchmark string prefixes
	if prefix != nil && len(prefix.PrefixStr) > 0 && MaxPrefixSize > 0 ***REMOVED***
		if len(prefix.PrefixStr) > MaxPrefixSize ***REMOVED***
			// limit prefix changes to 10k
			prefix.PrefixStr = prefix.PrefixStr[:MaxPrefixSize]
		***REMOVED***
		bmPrefix = newBmPrefix(prefix.PrefixStr, prefix.CaseInsensitive, rtl)
	***REMOVED*** else ***REMOVED***
		bmPrefix = nil
	***REMOVED***

	return &Code***REMOVED***
		Codes:       w.emitted,
		Strings:     w.stringtable,
		Sets:        w.settable,
		TrackCount:  w.trackcount,
		Caps:        w.caps,
		Capsize:     capsize,
		FcPrefix:    fcPrefix,
		BmPrefix:    bmPrefix,
		Anchors:     getAnchors(tree),
		RightToLeft: rtl,
	***REMOVED***, nil
***REMOVED***

// The main RegexCode generator. It does a depth-first walk
// through the tree and calls EmitFragment to emits code before
// and after each child of an interior node, and at each leaf.
func (w *writer) emitFragment(nodetype nodeType, node *regexNode, curIndex int) error ***REMOVED***
	bits := InstOp(0)

	if nodetype <= ntRef ***REMOVED***
		if (node.options & RightToLeft) != 0 ***REMOVED***
			bits |= Rtl
		***REMOVED***
		if (node.options & IgnoreCase) != 0 ***REMOVED***
			bits |= Ci
		***REMOVED***
	***REMOVED***
	ntBits := nodeType(bits)

	switch nodetype ***REMOVED***
	case ntConcatenate | beforeChild, ntConcatenate | afterChild, ntEmpty:
		break

	case ntAlternate | beforeChild:
		if curIndex < len(node.children)-1 ***REMOVED***
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
		***REMOVED***

	case ntAlternate | afterChild:
		if curIndex < len(node.children)-1 ***REMOVED***
			lbPos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(lbPos, w.curPos())
		***REMOVED*** else ***REMOVED***
			for i := 0; i < curIndex; i++ ***REMOVED***
				w.patchJump(w.popInt(), w.curPos())
			***REMOVED***
		***REMOVED***
		break

	case ntTestref | beforeChild:
		if curIndex == 0 ***REMOVED***
			w.emit(Setjump)
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
			w.emit1(Testref, w.mapCapnum(node.m))
			w.emit(Forejump)
		***REMOVED***

	case ntTestref | afterChild:
		if curIndex == 0 ***REMOVED***
			branchpos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(branchpos, w.curPos())
			w.emit(Forejump)
			if len(node.children) <= 1 ***REMOVED***
				w.patchJump(w.popInt(), w.curPos())
			***REMOVED***
		***REMOVED*** else if curIndex == 1 ***REMOVED***
			w.patchJump(w.popInt(), w.curPos())
		***REMOVED***

	case ntTestgroup | beforeChild:
		if curIndex == 0 ***REMOVED***
			w.emit(Setjump)
			w.emit(Setmark)
			w.pushInt(w.curPos())
			w.emit1(Lazybranch, 0)
		***REMOVED***

	case ntTestgroup | afterChild:
		if curIndex == 0 ***REMOVED***
			w.emit(Getmark)
			w.emit(Forejump)
		***REMOVED*** else if curIndex == 1 ***REMOVED***
			Branchpos := w.popInt()
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
			w.patchJump(Branchpos, w.curPos())
			w.emit(Getmark)
			w.emit(Forejump)
			if len(node.children) <= 2 ***REMOVED***
				w.patchJump(w.popInt(), w.curPos())
			***REMOVED***
		***REMOVED*** else if curIndex == 2 ***REMOVED***
			w.patchJump(w.popInt(), w.curPos())
		***REMOVED***

	case ntLoop | beforeChild, ntLazyloop | beforeChild:

		if node.n < math.MaxInt32 || node.m > 1 ***REMOVED***
			if node.m == 0 ***REMOVED***
				w.emit1(Nullcount, 0)
			***REMOVED*** else ***REMOVED***
				w.emit1(Setcount, 1-node.m)
			***REMOVED***
		***REMOVED*** else if node.m == 0 ***REMOVED***
			w.emit(Nullmark)
		***REMOVED*** else ***REMOVED***
			w.emit(Setmark)
		***REMOVED***

		if node.m == 0 ***REMOVED***
			w.pushInt(w.curPos())
			w.emit1(Goto, 0)
		***REMOVED***
		w.pushInt(w.curPos())

	case ntLoop | afterChild, ntLazyloop | afterChild:

		startJumpPos := w.curPos()
		lazy := (nodetype - (ntLoop | afterChild))

		if node.n < math.MaxInt32 || node.m > 1 ***REMOVED***
			if node.n == math.MaxInt32 ***REMOVED***
				w.emit2(InstOp(Branchcount+lazy), w.popInt(), math.MaxInt32)
			***REMOVED*** else ***REMOVED***
				w.emit2(InstOp(Branchcount+lazy), w.popInt(), node.n-node.m)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			w.emit1(InstOp(Branchmark+lazy), w.popInt())
		***REMOVED***

		if node.m == 0 ***REMOVED***
			w.patchJump(w.popInt(), startJumpPos)
		***REMOVED***

	case ntGroup | beforeChild, ntGroup | afterChild:

	case ntCapture | beforeChild:
		w.emit(Setmark)

	case ntCapture | afterChild:
		w.emit2(Capturemark, w.mapCapnum(node.m), w.mapCapnum(node.n))

	case ntRequire | beforeChild:
		// NOTE: the following line causes lookahead/lookbehind to be
		// NON-BACKTRACKING. It can be commented out with (*)
		w.emit(Setjump)

		w.emit(Setmark)

	case ntRequire | afterChild:
		w.emit(Getmark)

		// NOTE: the following line causes lookahead/lookbehind to be
		// NON-BACKTRACKING. It can be commented out with (*)
		w.emit(Forejump)

	case ntPrevent | beforeChild:
		w.emit(Setjump)
		w.pushInt(w.curPos())
		w.emit1(Lazybranch, 0)

	case ntPrevent | afterChild:
		w.emit(Backjump)
		w.patchJump(w.popInt(), w.curPos())
		w.emit(Forejump)

	case ntGreedy | beforeChild:
		w.emit(Setjump)

	case ntGreedy | afterChild:
		w.emit(Forejump)

	case ntOne, ntNotone:
		w.emit1(InstOp(node.t|ntBits), int(node.ch))

	case ntNotoneloop, ntNotonelazy, ntOneloop, ntOnelazy:
		if node.m > 0 ***REMOVED***
			if node.t == ntOneloop || node.t == ntOnelazy ***REMOVED***
				w.emit2(Onerep|bits, int(node.ch), node.m)
			***REMOVED*** else ***REMOVED***
				w.emit2(Notonerep|bits, int(node.ch), node.m)
			***REMOVED***
		***REMOVED***
		if node.n > node.m ***REMOVED***
			if node.n == math.MaxInt32 ***REMOVED***
				w.emit2(InstOp(node.t|ntBits), int(node.ch), math.MaxInt32)
			***REMOVED*** else ***REMOVED***
				w.emit2(InstOp(node.t|ntBits), int(node.ch), node.n-node.m)
			***REMOVED***
		***REMOVED***

	case ntSetloop, ntSetlazy:
		if node.m > 0 ***REMOVED***
			w.emit2(Setrep|bits, w.setCode(node.set), node.m)
		***REMOVED***
		if node.n > node.m ***REMOVED***
			if node.n == math.MaxInt32 ***REMOVED***
				w.emit2(InstOp(node.t|ntBits), w.setCode(node.set), math.MaxInt32)
			***REMOVED*** else ***REMOVED***
				w.emit2(InstOp(node.t|ntBits), w.setCode(node.set), node.n-node.m)
			***REMOVED***
		***REMOVED***

	case ntMulti:
		w.emit1(InstOp(node.t|ntBits), w.stringCode(node.str))

	case ntSet:
		w.emit1(InstOp(node.t|ntBits), w.setCode(node.set))

	case ntRef:
		w.emit1(InstOp(node.t|ntBits), w.mapCapnum(node.m))

	case ntNothing, ntBol, ntEol, ntBoundary, ntNonboundary, ntECMABoundary, ntNonECMABoundary, ntBeginning, ntStart, ntEndZ, ntEnd:
		w.emit(InstOp(node.t))

	default:
		return fmt.Errorf("unexpected opcode in regular expression generation: %v", nodetype)
	***REMOVED***

	return nil
***REMOVED***

// To avoid recursion, we use a simple integer stack.
// This is the push.
func (w *writer) pushInt(i int) ***REMOVED***
	w.intStack = append(w.intStack, i)
***REMOVED***

// Returns true if the stack is empty.
func (w *writer) emptyStack() bool ***REMOVED***
	return len(w.intStack) == 0
***REMOVED***

// This is the pop.
func (w *writer) popInt() int ***REMOVED***
	//get our item
	idx := len(w.intStack) - 1
	i := w.intStack[idx]
	//trim our slice
	w.intStack = w.intStack[:idx]
	return i
***REMOVED***

// Returns the current position in the emitted code.
func (w *writer) curPos() int ***REMOVED***
	return w.curpos
***REMOVED***

// Fixes up a jump instruction at the specified offset
// so that it jumps to the specified jumpDest.
func (w *writer) patchJump(offset, jumpDest int) ***REMOVED***
	w.emitted[offset+1] = jumpDest
***REMOVED***

// Returns an index in the set table for a charset
// uses a map to eliminate duplicates.
func (w *writer) setCode(set *CharSet) int ***REMOVED***
	if w.counting ***REMOVED***
		return 0
	***REMOVED***

	buf := &bytes.Buffer***REMOVED******REMOVED***

	set.mapHashFill(buf)
	hash := buf.String()
	i, ok := w.sethash[hash]
	if !ok ***REMOVED***
		i = len(w.sethash)
		w.sethash[hash] = i
		w.settable = append(w.settable, set)
	***REMOVED***
	return i
***REMOVED***

// Returns an index in the string table for a string.
// uses a map to eliminate duplicates.
func (w *writer) stringCode(str []rune) int ***REMOVED***
	if w.counting ***REMOVED***
		return 0
	***REMOVED***

	hash := string(str)
	i, ok := w.stringhash[hash]
	if !ok ***REMOVED***
		i = len(w.stringhash)
		w.stringhash[hash] = i
		w.stringtable = append(w.stringtable, str)
	***REMOVED***

	return i
***REMOVED***

// When generating code on a regex that uses a sparse set
// of capture slots, we hash them to a dense set of indices
// for an array of capture slots. Instead of doing the hash
// at match time, it's done at compile time, here.
func (w *writer) mapCapnum(capnum int) int ***REMOVED***
	if capnum == -1 ***REMOVED***
		return -1
	***REMOVED***

	if w.caps != nil ***REMOVED***
		return w.caps[capnum]
	***REMOVED***

	return capnum
***REMOVED***

// Emits a zero-argument operation. Note that the emit
// functions all run in two modes: they can emit code, or
// they can just count the size of the code.
func (w *writer) emit(op InstOp) ***REMOVED***
	if w.counting ***REMOVED***
		w.count++
		if opcodeBacktracks(op) ***REMOVED***
			w.trackcount++
		***REMOVED***
		return
	***REMOVED***
	w.emitted[w.curpos] = int(op)
	w.curpos++
***REMOVED***

// Emits a one-argument operation.
func (w *writer) emit1(op InstOp, opd1 int) ***REMOVED***
	if w.counting ***REMOVED***
		w.count += 2
		if opcodeBacktracks(op) ***REMOVED***
			w.trackcount++
		***REMOVED***
		return
	***REMOVED***
	w.emitted[w.curpos] = int(op)
	w.curpos++
	w.emitted[w.curpos] = opd1
	w.curpos++
***REMOVED***

// Emits a two-argument operation.
func (w *writer) emit2(op InstOp, opd1, opd2 int) ***REMOVED***
	if w.counting ***REMOVED***
		w.count += 3
		if opcodeBacktracks(op) ***REMOVED***
			w.trackcount++
		***REMOVED***
		return
	***REMOVED***
	w.emitted[w.curpos] = int(op)
	w.curpos++
	w.emitted[w.curpos] = opd1
	w.curpos++
	w.emitted[w.curpos] = opd2
	w.curpos++
***REMOVED***
