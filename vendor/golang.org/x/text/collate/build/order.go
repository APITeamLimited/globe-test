// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package build

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"unicode"

	"golang.org/x/text/internal/colltab"
	"golang.org/x/text/unicode/norm"
)

type logicalAnchor int

const (
	firstAnchor logicalAnchor = -1
	noAnchor                  = 0
	lastAnchor                = 1
)

// entry is used to keep track of a single entry in the collation element table
// during building. Examples of entries can be found in the Default Unicode
// Collation Element Table.
// See http://www.unicode.org/Public/UCA/6.0.0/allkeys.txt.
type entry struct ***REMOVED***
	str    string // same as string(runes)
	runes  []rune
	elems  []rawCE // the collation elements
	extend string  // weights of extend to be appended to elems
	before bool    // weights relative to next instead of previous.
	lock   bool    // entry is used in extension and can no longer be moved.

	// prev, next, and level are used to keep track of tailorings.
	prev, next *entry
	level      colltab.Level // next differs at this level
	skipRemove bool          // do not unlink when removed

	decompose bool // can use NFKD decomposition to generate elems
	exclude   bool // do not include in table
	implicit  bool // derived, is not included in the list
	modified  bool // entry was modified in tailoring
	logical   logicalAnchor

	expansionIndex    int // used to store index into expansion table
	contractionHandle ctHandle
	contractionIndex  int // index into contraction elements
***REMOVED***

func (e *entry) String() string ***REMOVED***
	return fmt.Sprintf("%X (%q) -> %X (ch:%x; ci:%d, ei:%d)",
		e.runes, e.str, e.elems, e.contractionHandle, e.contractionIndex, e.expansionIndex)
***REMOVED***

func (e *entry) skip() bool ***REMOVED***
	return e.contraction()
***REMOVED***

func (e *entry) expansion() bool ***REMOVED***
	return !e.decompose && len(e.elems) > 1
***REMOVED***

func (e *entry) contraction() bool ***REMOVED***
	return len(e.runes) > 1
***REMOVED***

func (e *entry) contractionStarter() bool ***REMOVED***
	return e.contractionHandle.n != 0
***REMOVED***

// nextIndexed gets the next entry that needs to be stored in the table.
// It returns the entry and the collation level at which the next entry differs
// from the current entry.
// Entries that can be explicitly derived and logical reset positions are
// examples of entries that will not be indexed.
func (e *entry) nextIndexed() (*entry, colltab.Level) ***REMOVED***
	level := e.level
	for e = e.next; e != nil && (e.exclude || len(e.elems) == 0); e = e.next ***REMOVED***
		if e.level < level ***REMOVED***
			level = e.level
		***REMOVED***
	***REMOVED***
	return e, level
***REMOVED***

// remove unlinks entry e from the sorted chain and clears the collation
// elements. e may not be at the front or end of the list. This should always
// be the case, as the front and end of the list are always logical anchors,
// which may not be removed.
func (e *entry) remove() ***REMOVED***
	if e.logical != noAnchor ***REMOVED***
		log.Fatalf("may not remove anchor %q", e.str)
	***REMOVED***
	// TODO: need to set e.prev.level to e.level if e.level is smaller?
	e.elems = nil
	if !e.skipRemove ***REMOVED***
		if e.prev != nil ***REMOVED***
			e.prev.next = e.next
		***REMOVED***
		if e.next != nil ***REMOVED***
			e.next.prev = e.prev
		***REMOVED***
	***REMOVED***
	e.skipRemove = false
***REMOVED***

// insertAfter inserts n after e.
func (e *entry) insertAfter(n *entry) ***REMOVED***
	if e == n ***REMOVED***
		panic("e == anchor")
	***REMOVED***
	if e == nil ***REMOVED***
		panic("unexpected nil anchor")
	***REMOVED***
	n.remove()
	n.decompose = false // redo decomposition test

	n.next = e.next
	n.prev = e
	if e.next != nil ***REMOVED***
		e.next.prev = n
	***REMOVED***
	e.next = n
***REMOVED***

// insertBefore inserts n before e.
func (e *entry) insertBefore(n *entry) ***REMOVED***
	if e == n ***REMOVED***
		panic("e == anchor")
	***REMOVED***
	if e == nil ***REMOVED***
		panic("unexpected nil anchor")
	***REMOVED***
	n.remove()
	n.decompose = false // redo decomposition test

	n.prev = e.prev
	n.next = e
	if e.prev != nil ***REMOVED***
		e.prev.next = n
	***REMOVED***
	e.prev = n
***REMOVED***

func (e *entry) encodeBase() (ce uint32, err error) ***REMOVED***
	switch ***REMOVED***
	case e.expansion():
		ce, err = makeExpandIndex(e.expansionIndex)
	default:
		if e.decompose ***REMOVED***
			log.Fatal("decompose should be handled elsewhere")
		***REMOVED***
		ce, err = makeCE(e.elems[0])
	***REMOVED***
	return
***REMOVED***

func (e *entry) encode() (ce uint32, err error) ***REMOVED***
	if e.skip() ***REMOVED***
		log.Fatal("cannot build colElem for entry that should be skipped")
	***REMOVED***
	switch ***REMOVED***
	case e.decompose:
		t1 := e.elems[0].w[2]
		t2 := 0
		if len(e.elems) > 1 ***REMOVED***
			t2 = e.elems[1].w[2]
		***REMOVED***
		ce, err = makeDecompose(t1, t2)
	case e.contractionStarter():
		ce, err = makeContractIndex(e.contractionHandle, e.contractionIndex)
	default:
		if len(e.runes) > 1 ***REMOVED***
			log.Fatal("colElem: contractions are handled in contraction trie")
		***REMOVED***
		ce, err = e.encodeBase()
	***REMOVED***
	return
***REMOVED***

// entryLess returns true if a sorts before b and false otherwise.
func entryLess(a, b *entry) bool ***REMOVED***
	if res, _ := compareWeights(a.elems, b.elems); res != 0 ***REMOVED***
		return res == -1
	***REMOVED***
	if a.logical != noAnchor ***REMOVED***
		return a.logical == firstAnchor
	***REMOVED***
	if b.logical != noAnchor ***REMOVED***
		return b.logical == lastAnchor
	***REMOVED***
	return a.str < b.str
***REMOVED***

type sortedEntries []*entry

func (s sortedEntries) Len() int ***REMOVED***
	return len(s)
***REMOVED***

func (s sortedEntries) Swap(i, j int) ***REMOVED***
	s[i], s[j] = s[j], s[i]
***REMOVED***

func (s sortedEntries) Less(i, j int) bool ***REMOVED***
	return entryLess(s[i], s[j])
***REMOVED***

type ordering struct ***REMOVED***
	id       string
	entryMap map[string]*entry
	ordered  []*entry
	handle   *trieHandle
***REMOVED***

// insert inserts e into both entryMap and ordered.
// Note that insert simply appends e to ordered.  To reattain a sorted
// order, o.sort() should be called.
func (o *ordering) insert(e *entry) ***REMOVED***
	if e.logical == noAnchor ***REMOVED***
		o.entryMap[e.str] = e
	***REMOVED*** else ***REMOVED***
		// Use key format as used in UCA rules.
		o.entryMap[fmt.Sprintf("[%s]", e.str)] = e
		// Also add index entry for XML format.
		o.entryMap[fmt.Sprintf("<%s/>", strings.Replace(e.str, " ", "_", -1))] = e
	***REMOVED***
	o.ordered = append(o.ordered, e)
***REMOVED***

// newEntry creates a new entry for the given info and inserts it into
// the index.
func (o *ordering) newEntry(s string, ces []rawCE) *entry ***REMOVED***
	e := &entry***REMOVED***
		runes: []rune(s),
		elems: ces,
		str:   s,
	***REMOVED***
	o.insert(e)
	return e
***REMOVED***

// find looks up and returns the entry for the given string.
// It returns nil if str is not in the index and if an implicit value
// cannot be derived, that is, if str represents more than one rune.
func (o *ordering) find(str string) *entry ***REMOVED***
	e := o.entryMap[str]
	if e == nil ***REMOVED***
		r := []rune(str)
		if len(r) == 1 ***REMOVED***
			const (
				firstHangul = 0xAC00
				lastHangul  = 0xD7A3
			)
			if r[0] >= firstHangul && r[0] <= lastHangul ***REMOVED***
				ce := []rawCE***REMOVED******REMOVED***
				nfd := norm.NFD.String(str)
				for _, r := range nfd ***REMOVED***
					ce = append(ce, o.find(string(r)).elems...)
				***REMOVED***
				e = o.newEntry(nfd, ce)
			***REMOVED*** else ***REMOVED***
				e = o.newEntry(string(r[0]), []rawCE***REMOVED***
					***REMOVED***w: []int***REMOVED***
						implicitPrimary(r[0]),
						defaultSecondary,
						defaultTertiary,
						int(r[0]),
					***REMOVED***,
					***REMOVED***,
				***REMOVED***)
				e.modified = true
			***REMOVED***
			e.exclude = true // do not index implicits
		***REMOVED***
	***REMOVED***
	return e
***REMOVED***

// makeRootOrdering returns a newly initialized ordering value and populates
// it with a set of logical reset points that can be used as anchors.
// The anchors first_tertiary_ignorable and __END__ will always sort at
// the beginning and end, respectively. This means that prev and next are non-nil
// for any indexed entry.
func makeRootOrdering() ordering ***REMOVED***
	const max = unicode.MaxRune
	o := ordering***REMOVED***
		entryMap: make(map[string]*entry),
	***REMOVED***
	insert := func(typ logicalAnchor, s string, ce []int) ***REMOVED***
		e := &entry***REMOVED***
			elems:   []rawCE***REMOVED******REMOVED***w: ce***REMOVED******REMOVED***,
			str:     s,
			exclude: true,
			logical: typ,
		***REMOVED***
		o.insert(e)
	***REMOVED***
	insert(firstAnchor, "first tertiary ignorable", []int***REMOVED***0, 0, 0, 0***REMOVED***)
	insert(lastAnchor, "last tertiary ignorable", []int***REMOVED***0, 0, 0, max***REMOVED***)
	insert(lastAnchor, "last primary ignorable", []int***REMOVED***0, defaultSecondary, defaultTertiary, max***REMOVED***)
	insert(lastAnchor, "last non ignorable", []int***REMOVED***maxPrimary, defaultSecondary, defaultTertiary, max***REMOVED***)
	insert(lastAnchor, "__END__", []int***REMOVED***1 << maxPrimaryBits, defaultSecondary, defaultTertiary, max***REMOVED***)
	return o
***REMOVED***

// patchForInsert eleminates entries from the list with more than one collation element.
// The next and prev fields of the eliminated entries still point to appropriate
// values in the newly created list.
// It requires that sort has been called.
func (o *ordering) patchForInsert() ***REMOVED***
	for i := 0; i < len(o.ordered)-1; ***REMOVED***
		e := o.ordered[i]
		lev := e.level
		n := e.next
		for ; n != nil && len(n.elems) > 1; n = n.next ***REMOVED***
			if n.level < lev ***REMOVED***
				lev = n.level
			***REMOVED***
			n.skipRemove = true
		***REMOVED***
		for ; o.ordered[i] != n; i++ ***REMOVED***
			o.ordered[i].level = lev
			o.ordered[i].next = n
			o.ordered[i+1].prev = e
		***REMOVED***
	***REMOVED***
***REMOVED***

// clone copies all ordering of es into a new ordering value.
func (o *ordering) clone() *ordering ***REMOVED***
	o.sort()
	oo := ordering***REMOVED***
		entryMap: make(map[string]*entry),
	***REMOVED***
	for _, e := range o.ordered ***REMOVED***
		ne := &entry***REMOVED***
			runes:     e.runes,
			elems:     e.elems,
			str:       e.str,
			decompose: e.decompose,
			exclude:   e.exclude,
			logical:   e.logical,
		***REMOVED***
		oo.insert(ne)
	***REMOVED***
	oo.sort() // link all ordering.
	oo.patchForInsert()
	return &oo
***REMOVED***

// front returns the first entry to be indexed.
// It assumes that sort() has been called.
func (o *ordering) front() *entry ***REMOVED***
	e := o.ordered[0]
	if e.prev != nil ***REMOVED***
		log.Panicf("unexpected first entry: %v", e)
	***REMOVED***
	// The first entry is always a logical position, which should not be indexed.
	e, _ = e.nextIndexed()
	return e
***REMOVED***

// sort sorts all ordering based on their collation elements and initializes
// the prev, next, and level fields accordingly.
func (o *ordering) sort() ***REMOVED***
	sort.Sort(sortedEntries(o.ordered))
	l := o.ordered
	for i := 1; i < len(l); i++ ***REMOVED***
		k := i - 1
		l[k].next = l[i]
		_, l[k].level = compareWeights(l[k].elems, l[i].elems)
		l[i].prev = l[k]
	***REMOVED***
***REMOVED***

// genColElems generates a collation element array from the runes in str. This
// assumes that all collation elements have already been added to the Builder.
func (o *ordering) genColElems(str string) []rawCE ***REMOVED***
	elems := []rawCE***REMOVED******REMOVED***
	for _, r := range []rune(str) ***REMOVED***
		for _, ce := range o.find(string(r)).elems ***REMOVED***
			if ce.w[0] != 0 || ce.w[1] != 0 || ce.w[2] != 0 ***REMOVED***
				elems = append(elems, ce)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return elems
***REMOVED***
