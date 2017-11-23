package regexp2

import (
	"bytes"
	"fmt"
)

// Match is a single regex result match that contains groups and repeated captures
// 	-Groups
//    -Capture
type Match struct ***REMOVED***
	Group //embeded group 0

	regex       *Regexp
	otherGroups []Group

	// input to the match
	textpos   int
	textstart int

	capcount   int
	caps       []int
	sparseCaps map[int]int

	// output from the match
	matches    [][]int
	matchcount []int

	// whether we've done any balancing with this match.  If we
	// have done balancing, we'll need to do extra work in Tidy().
	balancing bool
***REMOVED***

// Group is an explicit or implit (group 0) matched group within the pattern
type Group struct ***REMOVED***
	Capture // the last capture of this group is embeded for ease of use

	Name     string    // group name
	Captures []Capture // captures of this group
***REMOVED***

// Capture is a single capture of text within the larger original string
type Capture struct ***REMOVED***
	// the original string
	text []rune
	// the position in the original string where the first character of
	// captured substring was found.
	Index int
	// the length of the captured substring.
	Length int
***REMOVED***

// String returns the captured text as a String
func (c *Capture) String() string ***REMOVED***
	return string(c.text[c.Index : c.Index+c.Length])
***REMOVED***

// Runes returns the captured text as a rune slice
func (c *Capture) Runes() []rune ***REMOVED***
	return c.text[c.Index : c.Index+c.Length]
***REMOVED***

func newMatch(regex *Regexp, capcount int, text []rune, startpos int) *Match ***REMOVED***
	m := Match***REMOVED***
		regex:      regex,
		matchcount: make([]int, capcount),
		matches:    make([][]int, capcount),
		textstart:  startpos,
		balancing:  false,
	***REMOVED***
	m.Name = "0"
	m.text = text
	m.matches[0] = make([]int, 2)
	return &m
***REMOVED***

func newMatchSparse(regex *Regexp, caps map[int]int, capcount int, text []rune, startpos int) *Match ***REMOVED***
	m := newMatch(regex, capcount, text, startpos)
	m.sparseCaps = caps
	return m
***REMOVED***

func (m *Match) reset(text []rune, textstart int) ***REMOVED***
	m.text = text
	m.textstart = textstart
	for i := 0; i < len(m.matchcount); i++ ***REMOVED***
		m.matchcount[i] = 0
	***REMOVED***
	m.balancing = false
***REMOVED***

func (m *Match) tidy(textpos int) ***REMOVED***

	interval := m.matches[0]
	m.Index = interval[0]
	m.Length = interval[1]
	m.textpos = textpos
	m.capcount = m.matchcount[0]
	//copy our root capture to the list
	m.Group.Captures = []Capture***REMOVED***m.Group.Capture***REMOVED***

	if m.balancing ***REMOVED***
		// The idea here is that we want to compact all of our unbalanced captures.  To do that we
		// use j basically as a count of how many unbalanced captures we have at any given time
		// (really j is an index, but j/2 is the count).  First we skip past all of the real captures
		// until we find a balance captures.  Then we check each subsequent entry.  If it's a balance
		// capture (it's negative), we decrement j.  If it's a real capture, we increment j and copy
		// it down to the last free position.
		for cap := 0; cap < len(m.matchcount); cap++ ***REMOVED***
			limit := m.matchcount[cap] * 2
			matcharray := m.matches[cap]

			var i, j int

			for i = 0; i < limit; i++ ***REMOVED***
				if matcharray[i] < 0 ***REMOVED***
					break
				***REMOVED***
			***REMOVED***

			for j = i; i < limit; i++ ***REMOVED***
				if matcharray[i] < 0 ***REMOVED***
					// skip negative values
					j--
				***REMOVED*** else ***REMOVED***
					// but if we find something positive (an actual capture), copy it back to the last
					// unbalanced position.
					if i != j ***REMOVED***
						matcharray[j] = matcharray[i]
					***REMOVED***
					j++
				***REMOVED***
			***REMOVED***

			m.matchcount[cap] = j / 2
		***REMOVED***

		m.balancing = false
	***REMOVED***
***REMOVED***

// isMatched tells if a group was matched by capnum
func (m *Match) isMatched(cap int) bool ***REMOVED***
	return cap < len(m.matchcount) && m.matchcount[cap] > 0 && m.matches[cap][m.matchcount[cap]*2-1] != (-3+1)
***REMOVED***

// matchIndex returns the index of the last specified matched group by capnum
func (m *Match) matchIndex(cap int) int ***REMOVED***
	i := m.matches[cap][m.matchcount[cap]*2-2]
	if i >= 0 ***REMOVED***
		return i
	***REMOVED***

	return m.matches[cap][-3-i]
***REMOVED***

// matchLength returns the length of the last specified matched group by capnum
func (m *Match) matchLength(cap int) int ***REMOVED***
	i := m.matches[cap][m.matchcount[cap]*2-1]
	if i >= 0 ***REMOVED***
		return i
	***REMOVED***

	return m.matches[cap][-3-i]
***REMOVED***

// Nonpublic builder: add a capture to the group specified by "c"
func (m *Match) addMatch(c, start, l int) ***REMOVED***

	if m.matches[c] == nil ***REMOVED***
		m.matches[c] = make([]int, 2)
	***REMOVED***

	capcount := m.matchcount[c]

	if capcount*2+2 > len(m.matches[c]) ***REMOVED***
		oldmatches := m.matches[c]
		newmatches := make([]int, capcount*8)
		copy(newmatches, oldmatches[:capcount*2])
		m.matches[c] = newmatches
	***REMOVED***

	m.matches[c][capcount*2] = start
	m.matches[c][capcount*2+1] = l
	m.matchcount[c] = capcount + 1
	//log.Printf("addMatch: c=%v, i=%v, l=%v ... matches: %v", c, start, l, m.matches)
***REMOVED***

// Nonpublic builder: Add a capture to balance the specified group.  This is used by the
//                     balanced match construct. (?<foo-foo2>...)
//
// If there were no such thing as backtracking, this would be as simple as calling RemoveMatch(c).
// However, since we have backtracking, we need to keep track of everything.
func (m *Match) balanceMatch(c int) ***REMOVED***
	m.balancing = true

	// we'll look at the last capture first
	capcount := m.matchcount[c]
	target := capcount*2 - 2

	// first see if it is negative, and therefore is a reference to the next available
	// capture group for balancing.  If it is, we'll reset target to point to that capture.
	if m.matches[c][target] < 0 ***REMOVED***
		target = -3 - m.matches[c][target]
	***REMOVED***

	// move back to the previous capture
	target -= 2

	// if the previous capture is a reference, just copy that reference to the end.  Otherwise, point to it.
	if target >= 0 && m.matches[c][target] < 0 ***REMOVED***
		m.addMatch(c, m.matches[c][target], m.matches[c][target+1])
	***REMOVED*** else ***REMOVED***
		m.addMatch(c, -3-target, -4-target /* == -3 - (target + 1) */)
	***REMOVED***
***REMOVED***

// Nonpublic builder: removes a group match by capnum
func (m *Match) removeMatch(c int) ***REMOVED***
	m.matchcount[c]--
***REMOVED***

// GroupCount returns the number of groups this match has matched
func (m *Match) GroupCount() int ***REMOVED***
	return len(m.matchcount)
***REMOVED***

// GroupByName returns a group based on the name of the group, or nil if the group name does not exist
func (m *Match) GroupByName(name string) *Group ***REMOVED***
	num := m.regex.GroupNumberFromName(name)
	if num < 0 ***REMOVED***
		return nil
	***REMOVED***
	return m.GroupByNumber(num)
***REMOVED***

// GroupByNumber returns a group based on the number of the group, or nil if the group number does not exist
func (m *Match) GroupByNumber(num int) *Group ***REMOVED***
	// check our sparse map
	if m.sparseCaps != nil ***REMOVED***
		if newNum, ok := m.sparseCaps[num]; ok ***REMOVED***
			num = newNum
		***REMOVED***
	***REMOVED***
	if num >= len(m.matchcount) || num < 0 ***REMOVED***
		return nil
	***REMOVED***

	if num == 0 ***REMOVED***
		return &m.Group
	***REMOVED***

	m.populateOtherGroups()

	return &m.otherGroups[num-1]
***REMOVED***

// Groups returns all the capture groups, starting with group 0 (the full match)
func (m *Match) Groups() []Group ***REMOVED***
	m.populateOtherGroups()
	g := make([]Group, len(m.otherGroups)+1)
	g[0] = m.Group
	copy(g[1:], m.otherGroups)
	return g
***REMOVED***

func (m *Match) populateOtherGroups() ***REMOVED***
	// Construct all the Group objects first time called
	if m.otherGroups == nil ***REMOVED***
		m.otherGroups = make([]Group, len(m.matchcount)-1)
		for i := 0; i < len(m.otherGroups); i++ ***REMOVED***
			m.otherGroups[i] = newGroup(m.regex.GroupNameFromNumber(i+1), m.text, m.matches[i+1], m.matchcount[i+1])
		***REMOVED***
	***REMOVED***
***REMOVED***

func (m *Match) groupValueAppendToBuf(groupnum int, buf *bytes.Buffer) ***REMOVED***
	c := m.matchcount[groupnum]
	if c == 0 ***REMOVED***
		return
	***REMOVED***

	matches := m.matches[groupnum]

	index := matches[(c-1)*2]
	last := index + matches[(c*2)-1]

	for ; index < last; index++ ***REMOVED***
		buf.WriteRune(m.text[index])
	***REMOVED***
***REMOVED***

func newGroup(name string, text []rune, caps []int, capcount int) Group ***REMOVED***
	g := Group***REMOVED******REMOVED***
	g.text = text
	if capcount > 0 ***REMOVED***
		g.Index = caps[(capcount-1)*2]
		g.Length = caps[(capcount*2)-1]
	***REMOVED***
	g.Name = name
	g.Captures = make([]Capture, capcount)
	for i := 0; i < capcount; i++ ***REMOVED***
		g.Captures[i] = Capture***REMOVED***
			text:   text,
			Index:  caps[i*2],
			Length: caps[i*2+1],
		***REMOVED***
	***REMOVED***
	//log.Printf("newGroup! capcount %v, %+v", capcount, g)

	return g
***REMOVED***

func (m *Match) dump() string ***REMOVED***
	buf := &bytes.Buffer***REMOVED******REMOVED***
	buf.WriteRune('\n')
	if len(m.sparseCaps) > 0 ***REMOVED***
		for k, v := range m.sparseCaps ***REMOVED***
			fmt.Fprintf(buf, "Slot %v -> %v\n", k, v)
		***REMOVED***
	***REMOVED***

	for i, g := range m.Groups() ***REMOVED***
		fmt.Fprintf(buf, "Group %v (%v), %v caps:\n", i, g.Name, len(g.Captures))

		for _, c := range g.Captures ***REMOVED***
			fmt.Fprintf(buf, "  (%v, %v) %v\n", c.Index, c.Length, c.String())
		***REMOVED***
	***REMOVED***
	/*
		for i := 0; i < len(m.matchcount); i++ ***REMOVED***
			fmt.Fprintf(buf, "\nGroup %v (%v):\n", i, m.regex.GroupNameFromNumber(i))

			for j := 0; j < m.matchcount[i]; j++ ***REMOVED***
				text := ""

				if m.matches[i][j*2] >= 0 ***REMOVED***
					start := m.matches[i][j*2]
					text = m.text[start : start+m.matches[i][j*2+1]]
				***REMOVED***

				fmt.Fprintf(buf, "  (%v, %v) %v\n", m.matches[i][j*2], m.matches[i][j*2+1], text)
			***REMOVED***
		***REMOVED***
	*/
	return buf.String()
***REMOVED***
