// Package difflib is a partial port of Python difflib module.
//
// It provides tools to compare sequences of strings and generate textual diffs.
//
// The following class and functions have been ported:
//
// - SequenceMatcher
//
// - unified_diff
//
// - context_diff
//
// Getting unified diffs was the main goal of the port. Keep in mind this code
// is mostly suitable to output text differences in a human friendly way, there
// are no guarantees generated diffs are consumable by patch(1).
package difflib

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

func min(a, b int) int ***REMOVED***
	if a < b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func max(a, b int) int ***REMOVED***
	if a > b ***REMOVED***
		return a
	***REMOVED***
	return b
***REMOVED***

func calculateRatio(matches, length int) float64 ***REMOVED***
	if length > 0 ***REMOVED***
		return 2.0 * float64(matches) / float64(length)
	***REMOVED***
	return 1.0
***REMOVED***

type Match struct ***REMOVED***
	A    int
	B    int
	Size int
***REMOVED***

type OpCode struct ***REMOVED***
	Tag byte
	I1  int
	I2  int
	J1  int
	J2  int
***REMOVED***

// SequenceMatcher compares sequence of strings. The basic
// algorithm predates, and is a little fancier than, an algorithm
// published in the late 1980's by Ratcliff and Obershelp under the
// hyperbolic name "gestalt pattern matching".  The basic idea is to find
// the longest contiguous matching subsequence that contains no "junk"
// elements (R-O doesn't address junk).  The same idea is then applied
// recursively to the pieces of the sequences to the left and to the right
// of the matching subsequence.  This does not yield minimal edit
// sequences, but does tend to yield matches that "look right" to people.
//
// SequenceMatcher tries to compute a "human-friendly diff" between two
// sequences.  Unlike e.g. UNIX(tm) diff, the fundamental notion is the
// longest *contiguous* & junk-free matching subsequence.  That's what
// catches peoples' eyes.  The Windows(tm) windiff has another interesting
// notion, pairing up elements that appear uniquely in each sequence.
// That, and the method here, appear to yield more intuitive difference
// reports than does diff.  This method appears to be the least vulnerable
// to synching up on blocks of "junk lines", though (like blank lines in
// ordinary text files, or maybe "<P>" lines in HTML files).  That may be
// because this is the only method of the 3 that has a *concept* of
// "junk" <wink>.
//
// Timing:  Basic R-O is cubic time worst case and quadratic time expected
// case.  SequenceMatcher is quadratic time for the worst case and has
// expected-case behavior dependent in a complicated way on how many
// elements the sequences have in common; best case time is linear.
type SequenceMatcher struct ***REMOVED***
	a              []string
	b              []string
	b2j            map[string][]int
	IsJunk         func(string) bool
	autoJunk       bool
	bJunk          map[string]struct***REMOVED******REMOVED***
	matchingBlocks []Match
	fullBCount     map[string]int
	bPopular       map[string]struct***REMOVED******REMOVED***
	opCodes        []OpCode
***REMOVED***

func NewMatcher(a, b []string) *SequenceMatcher ***REMOVED***
	m := SequenceMatcher***REMOVED***autoJunk: true***REMOVED***
	m.SetSeqs(a, b)
	return &m
***REMOVED***

func NewMatcherWithJunk(a, b []string, autoJunk bool,
	isJunk func(string) bool) *SequenceMatcher ***REMOVED***

	m := SequenceMatcher***REMOVED***IsJunk: isJunk, autoJunk: autoJunk***REMOVED***
	m.SetSeqs(a, b)
	return &m
***REMOVED***

// Set two sequences to be compared.
func (m *SequenceMatcher) SetSeqs(a, b []string) ***REMOVED***
	m.SetSeq1(a)
	m.SetSeq2(b)
***REMOVED***

// Set the first sequence to be compared. The second sequence to be compared is
// not changed.
//
// SequenceMatcher computes and caches detailed information about the second
// sequence, so if you want to compare one sequence S against many sequences,
// use .SetSeq2(s) once and call .SetSeq1(x) repeatedly for each of the other
// sequences.
//
// See also SetSeqs() and SetSeq2().
func (m *SequenceMatcher) SetSeq1(a []string) ***REMOVED***
	if &a == &m.a ***REMOVED***
		return
	***REMOVED***
	m.a = a
	m.matchingBlocks = nil
	m.opCodes = nil
***REMOVED***

// Set the second sequence to be compared. The first sequence to be compared is
// not changed.
func (m *SequenceMatcher) SetSeq2(b []string) ***REMOVED***
	if &b == &m.b ***REMOVED***
		return
	***REMOVED***
	m.b = b
	m.matchingBlocks = nil
	m.opCodes = nil
	m.fullBCount = nil
	m.chainB()
***REMOVED***

func (m *SequenceMatcher) chainB() ***REMOVED***
	// Populate line -> index mapping
	b2j := map[string][]int***REMOVED******REMOVED***
	for i, s := range m.b ***REMOVED***
		indices := b2j[s]
		indices = append(indices, i)
		b2j[s] = indices
	***REMOVED***

	// Purge junk elements
	m.bJunk = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	if m.IsJunk != nil ***REMOVED***
		junk := m.bJunk
		for s, _ := range b2j ***REMOVED***
			if m.IsJunk(s) ***REMOVED***
				junk[s] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		for s, _ := range junk ***REMOVED***
			delete(b2j, s)
		***REMOVED***
	***REMOVED***

	// Purge remaining popular elements
	popular := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	n := len(m.b)
	if m.autoJunk && n >= 200 ***REMOVED***
		ntest := n/100 + 1
		for s, indices := range b2j ***REMOVED***
			if len(indices) > ntest ***REMOVED***
				popular[s] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
		for s, _ := range popular ***REMOVED***
			delete(b2j, s)
		***REMOVED***
	***REMOVED***
	m.bPopular = popular
	m.b2j = b2j
***REMOVED***

func (m *SequenceMatcher) isBJunk(s string) bool ***REMOVED***
	_, ok := m.bJunk[s]
	return ok
***REMOVED***

// Find longest matching block in a[alo:ahi] and b[blo:bhi].
//
// If IsJunk is not defined:
//
// Return (i,j,k) such that a[i:i+k] is equal to b[j:j+k], where
//     alo <= i <= i+k <= ahi
//     blo <= j <= j+k <= bhi
// and for all (i',j',k') meeting those conditions,
//     k >= k'
//     i <= i'
//     and if i == i', j <= j'
//
// In other words, of all maximal matching blocks, return one that
// starts earliest in a, and of all those maximal matching blocks that
// start earliest in a, return the one that starts earliest in b.
//
// If IsJunk is defined, first the longest matching block is
// determined as above, but with the additional restriction that no
// junk element appears in the block.  Then that block is extended as
// far as possible by matching (only) junk elements on both sides.  So
// the resulting block never matches on junk except as identical junk
// happens to be adjacent to an "interesting" match.
//
// If no blocks match, return (alo, blo, 0).
func (m *SequenceMatcher) findLongestMatch(alo, ahi, blo, bhi int) Match ***REMOVED***
	// CAUTION:  stripping common prefix or suffix would be incorrect.
	// E.g.,
	//    ab
	//    acab
	// Longest matching block is "ab", but if common prefix is
	// stripped, it's "a" (tied with "b").  UNIX(tm) diff does so
	// strip, so ends up claiming that ab is changed to acab by
	// inserting "ca" in the middle.  That's minimal but unintuitive:
	// "it's obvious" that someone inserted "ac" at the front.
	// Windiff ends up at the same place as diff, but by pairing up
	// the unique 'b's and then matching the first two 'a's.
	besti, bestj, bestsize := alo, blo, 0

	// find longest junk-free match
	// during an iteration of the loop, j2len[j] = length of longest
	// junk-free match ending with a[i-1] and b[j]
	j2len := map[int]int***REMOVED******REMOVED***
	for i := alo; i != ahi; i++ ***REMOVED***
		// look at all instances of a[i] in b; note that because
		// b2j has no junk keys, the loop is skipped if a[i] is junk
		newj2len := map[int]int***REMOVED******REMOVED***
		for _, j := range m.b2j[m.a[i]] ***REMOVED***
			// a[i] matches b[j]
			if j < blo ***REMOVED***
				continue
			***REMOVED***
			if j >= bhi ***REMOVED***
				break
			***REMOVED***
			k := j2len[j-1] + 1
			newj2len[j] = k
			if k > bestsize ***REMOVED***
				besti, bestj, bestsize = i-k+1, j-k+1, k
			***REMOVED***
		***REMOVED***
		j2len = newj2len
	***REMOVED***

	// Extend the best by non-junk elements on each end.  In particular,
	// "popular" non-junk elements aren't in b2j, which greatly speeds
	// the inner loop above, but also means "the best" match so far
	// doesn't contain any junk *or* popular non-junk elements.
	for besti > alo && bestj > blo && !m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] ***REMOVED***
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	***REMOVED***
	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		!m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] ***REMOVED***
		bestsize += 1
	***REMOVED***

	// Now that we have a wholly interesting match (albeit possibly
	// empty!), we may as well suck up the matching junk on each
	// side of it too.  Can't think of a good reason not to, and it
	// saves post-processing the (possibly considerable) expense of
	// figuring out what to do with it.  In the case of an empty
	// interesting match, this is clearly the right thing to do,
	// because no other kind of match is possible in the regions.
	for besti > alo && bestj > blo && m.isBJunk(m.b[bestj-1]) &&
		m.a[besti-1] == m.b[bestj-1] ***REMOVED***
		besti, bestj, bestsize = besti-1, bestj-1, bestsize+1
	***REMOVED***
	for besti+bestsize < ahi && bestj+bestsize < bhi &&
		m.isBJunk(m.b[bestj+bestsize]) &&
		m.a[besti+bestsize] == m.b[bestj+bestsize] ***REMOVED***
		bestsize += 1
	***REMOVED***

	return Match***REMOVED***A: besti, B: bestj, Size: bestsize***REMOVED***
***REMOVED***

// Return list of triples describing matching subsequences.
//
// Each triple is of the form (i, j, n), and means that
// a[i:i+n] == b[j:j+n].  The triples are monotonically increasing in
// i and in j. It's also guaranteed that if (i, j, n) and (i', j', n') are
// adjacent triples in the list, and the second is not the last triple in the
// list, then i+n != i' or j+n != j'. IOW, adjacent triples never describe
// adjacent equal blocks.
//
// The last triple is a dummy, (len(a), len(b), 0), and is the only
// triple with n==0.
func (m *SequenceMatcher) GetMatchingBlocks() []Match ***REMOVED***
	if m.matchingBlocks != nil ***REMOVED***
		return m.matchingBlocks
	***REMOVED***

	var matchBlocks func(alo, ahi, blo, bhi int, matched []Match) []Match
	matchBlocks = func(alo, ahi, blo, bhi int, matched []Match) []Match ***REMOVED***
		match := m.findLongestMatch(alo, ahi, blo, bhi)
		i, j, k := match.A, match.B, match.Size
		if match.Size > 0 ***REMOVED***
			if alo < i && blo < j ***REMOVED***
				matched = matchBlocks(alo, i, blo, j, matched)
			***REMOVED***
			matched = append(matched, match)
			if i+k < ahi && j+k < bhi ***REMOVED***
				matched = matchBlocks(i+k, ahi, j+k, bhi, matched)
			***REMOVED***
		***REMOVED***
		return matched
	***REMOVED***
	matched := matchBlocks(0, len(m.a), 0, len(m.b), nil)

	// It's possible that we have adjacent equal blocks in the
	// matching_blocks list now.
	nonAdjacent := []Match***REMOVED******REMOVED***
	i1, j1, k1 := 0, 0, 0
	for _, b := range matched ***REMOVED***
		// Is this block adjacent to i1, j1, k1?
		i2, j2, k2 := b.A, b.B, b.Size
		if i1+k1 == i2 && j1+k1 == j2 ***REMOVED***
			// Yes, so collapse them -- this just increases the length of
			// the first block by the length of the second, and the first
			// block so lengthened remains the block to compare against.
			k1 += k2
		***REMOVED*** else ***REMOVED***
			// Not adjacent.  Remember the first block (k1==0 means it's
			// the dummy we started with), and make the second block the
			// new block to compare against.
			if k1 > 0 ***REMOVED***
				nonAdjacent = append(nonAdjacent, Match***REMOVED***i1, j1, k1***REMOVED***)
			***REMOVED***
			i1, j1, k1 = i2, j2, k2
		***REMOVED***
	***REMOVED***
	if k1 > 0 ***REMOVED***
		nonAdjacent = append(nonAdjacent, Match***REMOVED***i1, j1, k1***REMOVED***)
	***REMOVED***

	nonAdjacent = append(nonAdjacent, Match***REMOVED***len(m.a), len(m.b), 0***REMOVED***)
	m.matchingBlocks = nonAdjacent
	return m.matchingBlocks
***REMOVED***

// Return list of 5-tuples describing how to turn a into b.
//
// Each tuple is of the form (tag, i1, i2, j1, j2).  The first tuple
// has i1 == j1 == 0, and remaining tuples have i1 == the i2 from the
// tuple preceding it, and likewise for j1 == the previous j2.
//
// The tags are characters, with these meanings:
//
// 'r' (replace):  a[i1:i2] should be replaced by b[j1:j2]
//
// 'd' (delete):   a[i1:i2] should be deleted, j1==j2 in this case.
//
// 'i' (insert):   b[j1:j2] should be inserted at a[i1:i1], i1==i2 in this case.
//
// 'e' (equal):    a[i1:i2] == b[j1:j2]
func (m *SequenceMatcher) GetOpCodes() []OpCode ***REMOVED***
	if m.opCodes != nil ***REMOVED***
		return m.opCodes
	***REMOVED***
	i, j := 0, 0
	matching := m.GetMatchingBlocks()
	opCodes := make([]OpCode, 0, len(matching))
	for _, m := range matching ***REMOVED***
		//  invariant:  we've pumped out correct diffs to change
		//  a[:i] into b[:j], and the next matching block is
		//  a[ai:ai+size] == b[bj:bj+size]. So we need to pump
		//  out a diff to change a[i:ai] into b[j:bj], pump out
		//  the matching block, and move (i,j) beyond the match
		ai, bj, size := m.A, m.B, m.Size
		tag := byte(0)
		if i < ai && j < bj ***REMOVED***
			tag = 'r'
		***REMOVED*** else if i < ai ***REMOVED***
			tag = 'd'
		***REMOVED*** else if j < bj ***REMOVED***
			tag = 'i'
		***REMOVED***
		if tag > 0 ***REMOVED***
			opCodes = append(opCodes, OpCode***REMOVED***tag, i, ai, j, bj***REMOVED***)
		***REMOVED***
		i, j = ai+size, bj+size
		// the list of matching blocks is terminated by a
		// sentinel with size 0
		if size > 0 ***REMOVED***
			opCodes = append(opCodes, OpCode***REMOVED***'e', ai, i, bj, j***REMOVED***)
		***REMOVED***
	***REMOVED***
	m.opCodes = opCodes
	return m.opCodes
***REMOVED***

// Isolate change clusters by eliminating ranges with no changes.
//
// Return a generator of groups with up to n lines of context.
// Each group is in the same format as returned by GetOpCodes().
func (m *SequenceMatcher) GetGroupedOpCodes(n int) [][]OpCode ***REMOVED***
	if n < 0 ***REMOVED***
		n = 3
	***REMOVED***
	codes := m.GetOpCodes()
	if len(codes) == 0 ***REMOVED***
		codes = []OpCode***REMOVED***OpCode***REMOVED***'e', 0, 1, 0, 1***REMOVED******REMOVED***
	***REMOVED***
	// Fixup leading and trailing groups if they show no changes.
	if codes[0].Tag == 'e' ***REMOVED***
		c := codes[0]
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		codes[0] = OpCode***REMOVED***c.Tag, max(i1, i2-n), i2, max(j1, j2-n), j2***REMOVED***
	***REMOVED***
	if codes[len(codes)-1].Tag == 'e' ***REMOVED***
		c := codes[len(codes)-1]
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		codes[len(codes)-1] = OpCode***REMOVED***c.Tag, i1, min(i2, i1+n), j1, min(j2, j1+n)***REMOVED***
	***REMOVED***
	nn := n + n
	groups := [][]OpCode***REMOVED******REMOVED***
	group := []OpCode***REMOVED******REMOVED***
	for _, c := range codes ***REMOVED***
		i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
		// End the current group and start a new one whenever
		// there is a large range with no changes.
		if c.Tag == 'e' && i2-i1 > nn ***REMOVED***
			group = append(group, OpCode***REMOVED***c.Tag, i1, min(i2, i1+n),
				j1, min(j2, j1+n)***REMOVED***)
			groups = append(groups, group)
			group = []OpCode***REMOVED******REMOVED***
			i1, j1 = max(i1, i2-n), max(j1, j2-n)
		***REMOVED***
		group = append(group, OpCode***REMOVED***c.Tag, i1, i2, j1, j2***REMOVED***)
	***REMOVED***
	if len(group) > 0 && !(len(group) == 1 && group[0].Tag == 'e') ***REMOVED***
		groups = append(groups, group)
	***REMOVED***
	return groups
***REMOVED***

// Return a measure of the sequences' similarity (float in [0,1]).
//
// Where T is the total number of elements in both sequences, and
// M is the number of matches, this is 2.0*M / T.
// Note that this is 1 if the sequences are identical, and 0 if
// they have nothing in common.
//
// .Ratio() is expensive to compute if you haven't already computed
// .GetMatchingBlocks() or .GetOpCodes(), in which case you may
// want to try .QuickRatio() or .RealQuickRation() first to get an
// upper bound.
func (m *SequenceMatcher) Ratio() float64 ***REMOVED***
	matches := 0
	for _, m := range m.GetMatchingBlocks() ***REMOVED***
		matches += m.Size
	***REMOVED***
	return calculateRatio(matches, len(m.a)+len(m.b))
***REMOVED***

// Return an upper bound on ratio() relatively quickly.
//
// This isn't defined beyond that it is an upper bound on .Ratio(), and
// is faster to compute.
func (m *SequenceMatcher) QuickRatio() float64 ***REMOVED***
	// viewing a and b as multisets, set matches to the cardinality
	// of their intersection; this counts the number of matches
	// without regard to order, so is clearly an upper bound
	if m.fullBCount == nil ***REMOVED***
		m.fullBCount = map[string]int***REMOVED******REMOVED***
		for _, s := range m.b ***REMOVED***
			m.fullBCount[s] = m.fullBCount[s] + 1
		***REMOVED***
	***REMOVED***

	// avail[x] is the number of times x appears in 'b' less the
	// number of times we've seen it in 'a' so far ... kinda
	avail := map[string]int***REMOVED******REMOVED***
	matches := 0
	for _, s := range m.a ***REMOVED***
		n, ok := avail[s]
		if !ok ***REMOVED***
			n = m.fullBCount[s]
		***REMOVED***
		avail[s] = n - 1
		if n > 0 ***REMOVED***
			matches += 1
		***REMOVED***
	***REMOVED***
	return calculateRatio(matches, len(m.a)+len(m.b))
***REMOVED***

// Return an upper bound on ratio() very quickly.
//
// This isn't defined beyond that it is an upper bound on .Ratio(), and
// is faster to compute than either .Ratio() or .QuickRatio().
func (m *SequenceMatcher) RealQuickRatio() float64 ***REMOVED***
	la, lb := len(m.a), len(m.b)
	return calculateRatio(min(la, lb), la+lb)
***REMOVED***

// Convert range to the "ed" format
func formatRangeUnified(start, stop int) string ***REMOVED***
	// Per the diff spec at http://www.unix.org/single_unix_specification/
	beginning := start + 1 // lines start numbering with one
	length := stop - start
	if length == 1 ***REMOVED***
		return fmt.Sprintf("%d", beginning)
	***REMOVED***
	if length == 0 ***REMOVED***
		beginning -= 1 // empty ranges begin at line just before the range
	***REMOVED***
	return fmt.Sprintf("%d,%d", beginning, length)
***REMOVED***

// Unified diff parameters
type UnifiedDiff struct ***REMOVED***
	A        []string // First sequence lines
	FromFile string   // First file name
	FromDate string   // First file time
	B        []string // Second sequence lines
	ToFile   string   // Second file name
	ToDate   string   // Second file time
	Eol      string   // Headers end of line, defaults to LF
	Context  int      // Number of context lines
***REMOVED***

// Compare two sequences of lines; generate the delta as a unified diff.
//
// Unified diffs are a compact way of showing line changes and a few
// lines of context.  The number of context lines is set by 'n' which
// defaults to three.
//
// By default, the diff control lines (those with ---, +++, or @@) are
// created with a trailing newline.  This is helpful so that inputs
// created from file.readlines() result in diffs that are suitable for
// file.writelines() since both the inputs and outputs have trailing
// newlines.
//
// For inputs that do not have trailing newlines, set the lineterm
// argument to "" so that the output will be uniformly newline free.
//
// The unidiff format normally has a header for filenames and modification
// times.  Any or all of these may be specified using strings for
// 'fromfile', 'tofile', 'fromfiledate', and 'tofiledate'.
// The modification times are normally expressed in the ISO 8601 format.
func WriteUnifiedDiff(writer io.Writer, diff UnifiedDiff) error ***REMOVED***
	buf := bufio.NewWriter(writer)
	defer buf.Flush()
	wf := func(format string, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
		_, err := buf.WriteString(fmt.Sprintf(format, args...))
		return err
	***REMOVED***
	ws := func(s string) error ***REMOVED***
		_, err := buf.WriteString(s)
		return err
	***REMOVED***

	if len(diff.Eol) == 0 ***REMOVED***
		diff.Eol = "\n"
	***REMOVED***

	started := false
	m := NewMatcher(diff.A, diff.B)
	for _, g := range m.GetGroupedOpCodes(diff.Context) ***REMOVED***
		if !started ***REMOVED***
			started = true
			fromDate := ""
			if len(diff.FromDate) > 0 ***REMOVED***
				fromDate = "\t" + diff.FromDate
			***REMOVED***
			toDate := ""
			if len(diff.ToDate) > 0 ***REMOVED***
				toDate = "\t" + diff.ToDate
			***REMOVED***
			if diff.FromFile != "" || diff.ToFile != "" ***REMOVED***
				err := wf("--- %s%s%s", diff.FromFile, fromDate, diff.Eol)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
				err = wf("+++ %s%s%s", diff.ToFile, toDate, diff.Eol)
				if err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		first, last := g[0], g[len(g)-1]
		range1 := formatRangeUnified(first.I1, last.I2)
		range2 := formatRangeUnified(first.J1, last.J2)
		if err := wf("@@ -%s +%s @@%s", range1, range2, diff.Eol); err != nil ***REMOVED***
			return err
		***REMOVED***
		for _, c := range g ***REMOVED***
			i1, i2, j1, j2 := c.I1, c.I2, c.J1, c.J2
			if c.Tag == 'e' ***REMOVED***
				for _, line := range diff.A[i1:i2] ***REMOVED***
					if err := ws(" " + line); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***
			if c.Tag == 'r' || c.Tag == 'd' ***REMOVED***
				for _, line := range diff.A[i1:i2] ***REMOVED***
					if err := ws("-" + line); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
			if c.Tag == 'r' || c.Tag == 'i' ***REMOVED***
				for _, line := range diff.B[j1:j2] ***REMOVED***
					if err := ws("+" + line); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// Like WriteUnifiedDiff but returns the diff a string.
func GetUnifiedDiffString(diff UnifiedDiff) (string, error) ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	err := WriteUnifiedDiff(w, diff)
	return string(w.Bytes()), err
***REMOVED***

// Convert range to the "ed" format.
func formatRangeContext(start, stop int) string ***REMOVED***
	// Per the diff spec at http://www.unix.org/single_unix_specification/
	beginning := start + 1 // lines start numbering with one
	length := stop - start
	if length == 0 ***REMOVED***
		beginning -= 1 // empty ranges begin at line just before the range
	***REMOVED***
	if length <= 1 ***REMOVED***
		return fmt.Sprintf("%d", beginning)
	***REMOVED***
	return fmt.Sprintf("%d,%d", beginning, beginning+length-1)
***REMOVED***

type ContextDiff UnifiedDiff

// Compare two sequences of lines; generate the delta as a context diff.
//
// Context diffs are a compact way of showing line changes and a few
// lines of context. The number of context lines is set by diff.Context
// which defaults to three.
//
// By default, the diff control lines (those with *** or ---) are
// created with a trailing newline.
//
// For inputs that do not have trailing newlines, set the diff.Eol
// argument to "" so that the output will be uniformly newline free.
//
// The context diff format normally has a header for filenames and
// modification times.  Any or all of these may be specified using
// strings for diff.FromFile, diff.ToFile, diff.FromDate, diff.ToDate.
// The modification times are normally expressed in the ISO 8601 format.
// If not specified, the strings default to blanks.
func WriteContextDiff(writer io.Writer, diff ContextDiff) error ***REMOVED***
	buf := bufio.NewWriter(writer)
	defer buf.Flush()
	var diffErr error
	wf := func(format string, args ...interface***REMOVED******REMOVED***) ***REMOVED***
		_, err := buf.WriteString(fmt.Sprintf(format, args...))
		if diffErr == nil && err != nil ***REMOVED***
			diffErr = err
		***REMOVED***
	***REMOVED***
	ws := func(s string) ***REMOVED***
		_, err := buf.WriteString(s)
		if diffErr == nil && err != nil ***REMOVED***
			diffErr = err
		***REMOVED***
	***REMOVED***

	if len(diff.Eol) == 0 ***REMOVED***
		diff.Eol = "\n"
	***REMOVED***

	prefix := map[byte]string***REMOVED***
		'i': "+ ",
		'd': "- ",
		'r': "! ",
		'e': "  ",
	***REMOVED***

	started := false
	m := NewMatcher(diff.A, diff.B)
	for _, g := range m.GetGroupedOpCodes(diff.Context) ***REMOVED***
		if !started ***REMOVED***
			started = true
			fromDate := ""
			if len(diff.FromDate) > 0 ***REMOVED***
				fromDate = "\t" + diff.FromDate
			***REMOVED***
			toDate := ""
			if len(diff.ToDate) > 0 ***REMOVED***
				toDate = "\t" + diff.ToDate
			***REMOVED***
			if diff.FromFile != "" || diff.ToFile != "" ***REMOVED***
				wf("*** %s%s%s", diff.FromFile, fromDate, diff.Eol)
				wf("--- %s%s%s", diff.ToFile, toDate, diff.Eol)
			***REMOVED***
		***REMOVED***

		first, last := g[0], g[len(g)-1]
		ws("***************" + diff.Eol)

		range1 := formatRangeContext(first.I1, last.I2)
		wf("*** %s ****%s", range1, diff.Eol)
		for _, c := range g ***REMOVED***
			if c.Tag == 'r' || c.Tag == 'd' ***REMOVED***
				for _, cc := range g ***REMOVED***
					if cc.Tag == 'i' ***REMOVED***
						continue
					***REMOVED***
					for _, line := range diff.A[cc.I1:cc.I2] ***REMOVED***
						ws(prefix[cc.Tag] + line)
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***

		range2 := formatRangeContext(first.J1, last.J2)
		wf("--- %s ----%s", range2, diff.Eol)
		for _, c := range g ***REMOVED***
			if c.Tag == 'r' || c.Tag == 'i' ***REMOVED***
				for _, cc := range g ***REMOVED***
					if cc.Tag == 'd' ***REMOVED***
						continue
					***REMOVED***
					for _, line := range diff.B[cc.J1:cc.J2] ***REMOVED***
						ws(prefix[cc.Tag] + line)
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return diffErr
***REMOVED***

// Like WriteContextDiff but returns the diff a string.
func GetContextDiffString(diff ContextDiff) (string, error) ***REMOVED***
	w := &bytes.Buffer***REMOVED******REMOVED***
	err := WriteContextDiff(w, diff)
	return string(w.Bytes()), err
***REMOVED***

// Split a string on "\n" while preserving them. The output can be used
// as input for UnifiedDiff and ContextDiff structures.
func SplitLines(s string) []string ***REMOVED***
	lines := strings.SplitAfter(s, "\n")
	lines[len(lines)-1] += "\n"
	return lines
***REMOVED***
