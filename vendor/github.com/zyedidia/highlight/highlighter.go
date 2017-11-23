package highlight

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// RunePos returns the rune index of a given byte index
// This could cause problems if the byte index is between code points
func runePos(p int, str string) int ***REMOVED***
	if p < 0 ***REMOVED***
		return 0
	***REMOVED***
	if p >= len(str) ***REMOVED***
		return utf8.RuneCountInString(str)
	***REMOVED***
	return utf8.RuneCountInString(str[:p])
***REMOVED***

func combineLineMatch(src, dst LineMatch) LineMatch ***REMOVED***
	for k, v := range src ***REMOVED***
		if g, ok := dst[k]; ok ***REMOVED***
			if g == 0 ***REMOVED***
				dst[k] = v
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			dst[k] = v
		***REMOVED***
	***REMOVED***
	return dst
***REMOVED***

// A State represents the region at the end of a line
type State *region

// LineStates is an interface for a buffer-like object which can also store the states and matches for every line
type LineStates interface ***REMOVED***
	Line(n int) string
	LinesNum() int
	State(lineN int) State
	SetState(lineN int, s State)
	SetMatch(lineN int, m LineMatch)
***REMOVED***

// A Highlighter contains the information needed to highlight a string
type Highlighter struct ***REMOVED***
	lastRegion *region
	Def        *Def
***REMOVED***

// NewHighlighter returns a new highlighter from the given syntax definition
func NewHighlighter(def *Def) *Highlighter ***REMOVED***
	h := new(Highlighter)
	h.Def = def
	return h
***REMOVED***

// LineMatch represents the syntax highlighting matches for one line. Each index where the coloring is changed is marked with that
// color's group (represented as one byte)
type LineMatch map[int]Group

func findIndex(regex *regexp.Regexp, skip *regexp.Regexp, str []rune, canMatchStart, canMatchEnd bool) []int ***REMOVED***
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") ***REMOVED***
		if !canMatchStart ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	if strings.Contains(regexStr, "$") ***REMOVED***
		if !canMatchEnd ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	var strbytes []byte
	if skip != nil ***REMOVED***
		strbytes = skip.ReplaceAllFunc([]byte(string(str)), func(match []byte) []byte ***REMOVED***
			res := make([]byte, utf8.RuneCount(match))
			return res
		***REMOVED***)
	***REMOVED*** else ***REMOVED***
		strbytes = []byte(string(str))
	***REMOVED***

	match := regex.FindIndex(strbytes)
	if match == nil ***REMOVED***
		return nil
	***REMOVED***
	// return []int***REMOVED***match.Index, match.Index + match.Length***REMOVED***
	return []int***REMOVED***runePos(match[0], string(str)), runePos(match[1], string(str))***REMOVED***
***REMOVED***

func findAllIndex(regex *regexp.Regexp, str []rune, canMatchStart, canMatchEnd bool) [][]int ***REMOVED***
	regexStr := regex.String()
	if strings.Contains(regexStr, "^") ***REMOVED***
		if !canMatchStart ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	if strings.Contains(regexStr, "$") ***REMOVED***
		if !canMatchEnd ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***
	matches := regex.FindAllIndex([]byte(string(str)), -1)
	for i, m := range matches ***REMOVED***
		matches[i][0] = runePos(m[0], string(str))
		matches[i][1] = runePos(m[1], string(str))
	***REMOVED***
	return matches
***REMOVED***

func (h *Highlighter) highlightRegion(highlights LineMatch, start int, canMatchEnd bool, lineNum int, line []rune, curRegion *region, statesOnly bool) LineMatch ***REMOVED***
	// highlights := make(LineMatch)

	if start == 0 ***REMOVED***
		if !statesOnly ***REMOVED***
			if _, ok := highlights[0]; !ok ***REMOVED***
				highlights[0] = curRegion.group
			***REMOVED***
		***REMOVED***
	***REMOVED***

	loc := findIndex(curRegion.end, curRegion.skip, line, start == 0, canMatchEnd)
	if loc != nil ***REMOVED***
		if !statesOnly ***REMOVED***
			highlights[start+loc[0]] = curRegion.limitGroup
		***REMOVED***
		if curRegion.parent == nil ***REMOVED***
			if !statesOnly ***REMOVED***
				highlights[start+loc[1]] = 0
				h.highlightRegion(highlights, start, false, lineNum, line[:loc[0]], curRegion, statesOnly)
			***REMOVED***
			h.highlightEmptyRegion(highlights, start+loc[1], canMatchEnd, lineNum, line[loc[1]:], statesOnly)
			return highlights
		***REMOVED***
		if !statesOnly ***REMOVED***
			highlights[start+loc[1]] = curRegion.parent.group
			h.highlightRegion(highlights, start, false, lineNum, line[:loc[0]], curRegion, statesOnly)
		***REMOVED***
		h.highlightRegion(highlights, start+loc[1], canMatchEnd, lineNum, line[loc[1]:], curRegion.parent, statesOnly)
		return highlights
	***REMOVED***

	if len(line) == 0 || statesOnly ***REMOVED***
		if canMatchEnd ***REMOVED***
			h.lastRegion = curRegion
		***REMOVED***

		return highlights
	***REMOVED***

	firstLoc := []int***REMOVED***len(line), 0***REMOVED***

	var firstRegion *region
	for _, r := range curRegion.rules.regions ***REMOVED***
		loc := findIndex(r.start, nil, line, start == 0, canMatchEnd)
		if loc != nil ***REMOVED***
			if loc[0] < firstLoc[0] ***REMOVED***
				firstLoc = loc
				firstRegion = r
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if firstLoc[0] != len(line) ***REMOVED***
		highlights[start+firstLoc[0]] = firstRegion.limitGroup
		h.highlightRegion(highlights, start, false, lineNum, line[:firstLoc[0]], curRegion, statesOnly)
		h.highlightRegion(highlights, start+firstLoc[1], canMatchEnd, lineNum, line[firstLoc[1]:], firstRegion, statesOnly)
		return highlights
	***REMOVED***

	fullHighlights := make([]Group, len([]rune(string(line))))
	for i := 0; i < len(fullHighlights); i++ ***REMOVED***
		fullHighlights[i] = curRegion.group
	***REMOVED***

	for _, p := range curRegion.rules.patterns ***REMOVED***
		matches := findAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches ***REMOVED***
			for i := m[0]; i < m[1]; i++ ***REMOVED***
				fullHighlights[i] = p.group
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for i, h := range fullHighlights ***REMOVED***
		if i == 0 || h != fullHighlights[i-1] ***REMOVED***
			// if _, ok := highlights[start+i]; !ok ***REMOVED***
			highlights[start+i] = h
			// ***REMOVED***
		***REMOVED***
	***REMOVED***

	if canMatchEnd ***REMOVED***
		h.lastRegion = curRegion
	***REMOVED***

	return highlights
***REMOVED***

func (h *Highlighter) highlightEmptyRegion(highlights LineMatch, start int, canMatchEnd bool, lineNum int, line []rune, statesOnly bool) LineMatch ***REMOVED***
	if len(line) == 0 ***REMOVED***
		if canMatchEnd ***REMOVED***
			h.lastRegion = nil
		***REMOVED***
		return highlights
	***REMOVED***

	firstLoc := []int***REMOVED***len(line), 0***REMOVED***
	var firstRegion *region
	for _, r := range h.Def.rules.regions ***REMOVED***
		loc := findIndex(r.start, nil, line, start == 0, canMatchEnd)
		if loc != nil ***REMOVED***
			if loc[0] < firstLoc[0] ***REMOVED***
				firstLoc = loc
				firstRegion = r
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if firstLoc[0] != len(line) ***REMOVED***
		if !statesOnly ***REMOVED***
			highlights[start+firstLoc[0]] = firstRegion.limitGroup
		***REMOVED***
		h.highlightEmptyRegion(highlights, start, false, lineNum, line[:firstLoc[0]], statesOnly)
		h.highlightRegion(highlights, start+firstLoc[1], canMatchEnd, lineNum, line[firstLoc[1]:], firstRegion, statesOnly)
		return highlights
	***REMOVED***

	if statesOnly ***REMOVED***
		if canMatchEnd ***REMOVED***
			h.lastRegion = nil
		***REMOVED***

		return highlights
	***REMOVED***

	fullHighlights := make([]Group, len(line))
	for _, p := range h.Def.rules.patterns ***REMOVED***
		matches := findAllIndex(p.regex, line, start == 0, canMatchEnd)
		for _, m := range matches ***REMOVED***
			for i := m[0]; i < m[1]; i++ ***REMOVED***
				fullHighlights[i] = p.group
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for i, h := range fullHighlights ***REMOVED***
		if i == 0 || h != fullHighlights[i-1] ***REMOVED***
			// if _, ok := highlights[start+i]; !ok ***REMOVED***
			highlights[start+i] = h
			// ***REMOVED***
		***REMOVED***
	***REMOVED***

	if canMatchEnd ***REMOVED***
		h.lastRegion = nil
	***REMOVED***

	return highlights
***REMOVED***

// HighlightString syntax highlights a string
// Use this function for simple syntax highlighting and use the other functions for
// more advanced syntax highlighting. They are optimized for quick rehighlighting of the same
// text with minor changes made
func (h *Highlighter) HighlightString(input string) []LineMatch ***REMOVED***
	lines := strings.Split(input, "\n")
	var lineMatches []LineMatch

	for i := 0; i < len(lines); i++ ***REMOVED***
		line := []rune(lines[i])
		highlights := make(LineMatch)

		if i == 0 || h.lastRegion == nil ***REMOVED***
			lineMatches = append(lineMatches, h.highlightEmptyRegion(highlights, 0, true, i, line, false))
		***REMOVED*** else ***REMOVED***
			lineMatches = append(lineMatches, h.highlightRegion(highlights, 0, true, i, line, h.lastRegion, false))
		***REMOVED***
	***REMOVED***

	return lineMatches
***REMOVED***

// HighlightStates correctly sets all states for the buffer
func (h *Highlighter) HighlightStates(input LineStates) ***REMOVED***
	for i := 0; i < input.LinesNum(); i++ ***REMOVED***
		line := []rune(input.Line(i))
		// highlights := make(LineMatch)

		if i == 0 || h.lastRegion == nil ***REMOVED***
			h.highlightEmptyRegion(nil, 0, true, i, line, true)
		***REMOVED*** else ***REMOVED***
			h.highlightRegion(nil, 0, true, i, line, h.lastRegion, true)
		***REMOVED***

		curState := h.lastRegion

		input.SetState(i, curState)
	***REMOVED***
***REMOVED***

// HighlightMatches sets the matches for each line in between startline and endline
// It sets all other matches in the buffer to nil to conserve memory
// This assumes that all the states are set correctly
func (h *Highlighter) HighlightMatches(input LineStates, startline, endline int) ***REMOVED***
	for i := startline; i < endline; i++ ***REMOVED***
		if i >= input.LinesNum() ***REMOVED***
			break
		***REMOVED***

		line := []rune(input.Line(i))
		highlights := make(LineMatch)

		var match LineMatch
		if i == 0 || input.State(i-1) == nil ***REMOVED***
			match = h.highlightEmptyRegion(highlights, 0, true, i, line, false)
		***REMOVED*** else ***REMOVED***
			match = h.highlightRegion(highlights, 0, true, i, line, input.State(i-1), false)
		***REMOVED***

		input.SetMatch(i, match)
	***REMOVED***
***REMOVED***

// ReHighlightStates will scan down from `startline` and set the appropriate end of line state
// for each line until it comes across the same state in two consecutive lines
func (h *Highlighter) ReHighlightStates(input LineStates, startline int) ***REMOVED***
	// lines := input.LineData()

	h.lastRegion = nil
	if startline > 0 ***REMOVED***
		h.lastRegion = input.State(startline - 1)
	***REMOVED***
	for i := startline; i < input.LinesNum(); i++ ***REMOVED***
		line := []rune(input.Line(i))
		// highlights := make(LineMatch)

		// var match LineMatch
		if i == 0 || h.lastRegion == nil ***REMOVED***
			h.highlightEmptyRegion(nil, 0, true, i, line, true)
		***REMOVED*** else ***REMOVED***
			h.highlightRegion(nil, 0, true, i, line, h.lastRegion, true)
		***REMOVED***
		curState := h.lastRegion
		lastState := input.State(i)

		input.SetState(i, curState)

		if curState == lastState ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// ReHighlightLine will rehighlight the state and match for a single line
func (h *Highlighter) ReHighlightLine(input LineStates, lineN int) ***REMOVED***
	line := []rune(input.Line(lineN))
	highlights := make(LineMatch)

	h.lastRegion = nil
	if lineN > 0 ***REMOVED***
		h.lastRegion = input.State(lineN - 1)
	***REMOVED***

	var match LineMatch
	if lineN == 0 || h.lastRegion == nil ***REMOVED***
		match = h.highlightEmptyRegion(highlights, 0, true, lineN, line, false)
	***REMOVED*** else ***REMOVED***
		match = h.highlightRegion(highlights, 0, true, lineN, line, h.lastRegion, false)
	***REMOVED***
	curState := h.lastRegion

	input.SetMatch(lineN, match)
	input.SetState(lineN, curState)
***REMOVED***
