/*
Package regexp2 is a regexp package that has an interface similar to Go's framework regexp engine but uses a
more feature full regex engine behind the scenes.

It doesn't have constant time guarantees, but it allows backtracking and is compatible with Perl5 and .NET.
You'll likely be better off with the RE2 engine from the regexp package and should only use this if you
need to write very complex patterns or require compatibility with .NET.
*/
package regexp2

import (
	"errors"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/dlclark/regexp2/syntax"
)

// Default timeout used when running regexp matches -- "forever"
var DefaultMatchTimeout = time.Duration(math.MaxInt64)

// Regexp is the representation of a compiled regular expression.
// A Regexp is safe for concurrent use by multiple goroutines.
type Regexp struct ***REMOVED***
	//timeout when trying to find matches
	MatchTimeout time.Duration

	// read-only after Compile
	pattern string       // as passed to Compile
	options RegexOptions // options

	caps     map[int]int    // capnum->index
	capnames map[string]int //capture group name -> index
	capslist []string       //sorted list of capture group names
	capsize  int            // size of the capture array

	code *syntax.Code // compiled program

	// cache of machines for running regexp
	muRun  sync.Mutex
	runner []*runner
***REMOVED***

// Compile parses a regular expression and returns, if successful,
// a Regexp object that can be used to match against text.
func Compile(expr string, opt RegexOptions) (*Regexp, error) ***REMOVED***
	// parse it
	tree, err := syntax.Parse(expr, syntax.RegexOptions(opt))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// translate it to code
	code, err := syntax.Write(tree)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// return it
	return &Regexp***REMOVED***
		pattern:      expr,
		options:      opt,
		caps:         code.Caps,
		capnames:     tree.Capnames,
		capslist:     tree.Caplist,
		capsize:      code.Capsize,
		code:         code,
		MatchTimeout: DefaultMatchTimeout,
	***REMOVED***, nil
***REMOVED***

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(str string, opt RegexOptions) *Regexp ***REMOVED***
	regexp, error := Compile(str, opt)
	if error != nil ***REMOVED***
		panic(`regexp2: Compile(` + quote(str) + `): ` + error.Error())
	***REMOVED***
	return regexp
***REMOVED***

// Escape adds backslashes to any special characters in the input string
func Escape(input string) string ***REMOVED***
	return syntax.Escape(input)
***REMOVED***

// Unescape removes any backslashes from previously-escaped special characters in the input string
func Unescape(input string) (string, error) ***REMOVED***
	return syntax.Unescape(input)
***REMOVED***

// String returns the source text used to compile the regular expression.
func (re *Regexp) String() string ***REMOVED***
	return re.pattern
***REMOVED***

func quote(s string) string ***REMOVED***
	if strconv.CanBackquote(s) ***REMOVED***
		return "`" + s + "`"
	***REMOVED***
	return strconv.Quote(s)
***REMOVED***

// RegexOptions impact the runtime and parsing behavior
// for each specific regex.  They are setable in code as well
// as in the regex pattern itself.
type RegexOptions int32

const (
	None                    RegexOptions = 0x0
	IgnoreCase                           = 0x0001 // "i"
	Multiline                            = 0x0002 // "m"
	ExplicitCapture                      = 0x0004 // "n"
	Compiled                             = 0x0008 // "c"
	Singleline                           = 0x0010 // "s"
	IgnorePatternWhitespace              = 0x0020 // "x"
	RightToLeft                          = 0x0040 // "r"
	Debug                                = 0x0080 // "d"
	ECMAScript                           = 0x0100 // "e"
	RE2                                  = 0x0200 // RE2 (regexp package) compatibility mode
	Unicode                              = 0x0400 // "u"
)

func (re *Regexp) RightToLeft() bool ***REMOVED***
	return re.options&RightToLeft != 0
***REMOVED***

func (re *Regexp) Debug() bool ***REMOVED***
	return re.options&Debug != 0
***REMOVED***

// Replace searches the input string and replaces each match found with the replacement text.
// Count will limit the number of matches attempted and startAt will allow
// us to skip past possible matches at the start of the input (left or right depending on RightToLeft option).
// Set startAt and count to -1 to go through the whole string
func (re *Regexp) Replace(input, replacement string, startAt, count int) (string, error) ***REMOVED***
	data, err := syntax.NewReplacerData(replacement, re.caps, re.capsize, re.capnames, syntax.RegexOptions(re.options))
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	//TODO: cache ReplacerData

	return replace(re, data, nil, input, startAt, count)
***REMOVED***

// ReplaceFunc searches the input string and replaces each match found using the string from the evaluator
// Count will limit the number of matches attempted and startAt will allow
// us to skip past possible matches at the start of the input (left or right depending on RightToLeft option).
// Set startAt and count to -1 to go through the whole string.
func (re *Regexp) ReplaceFunc(input string, evaluator MatchEvaluator, startAt, count int) (string, error) ***REMOVED***
	return replace(re, nil, evaluator, input, startAt, count)
***REMOVED***

// FindStringMatch searches the input string for a Regexp match
func (re *Regexp) FindStringMatch(s string) (*Match, error) ***REMOVED***
	// convert string to runes
	return re.run(false, -1, getRunes(s))
***REMOVED***

// FindRunesMatch searches the input rune slice for a Regexp match
func (re *Regexp) FindRunesMatch(r []rune) (*Match, error) ***REMOVED***
	return re.run(false, -1, r)
***REMOVED***

// FindStringMatchStartingAt searches the input string for a Regexp match starting at the startAt index
func (re *Regexp) FindStringMatchStartingAt(s string, startAt int) (*Match, error) ***REMOVED***
	if startAt > len(s) ***REMOVED***
		return nil, errors.New("startAt must be less than the length of the input string")
	***REMOVED***
	r, startAt := re.getRunesAndStart(s, startAt)
	if startAt == -1 ***REMOVED***
		// we didn't find our start index in the string -- that's a problem
		return nil, errors.New("startAt must align to the start of a valid rune in the input string")
	***REMOVED***

	return re.run(false, startAt, r)
***REMOVED***

// FindRunesMatchStartingAt searches the input rune slice for a Regexp match starting at the startAt index
func (re *Regexp) FindRunesMatchStartingAt(r []rune, startAt int) (*Match, error) ***REMOVED***
	return re.run(false, startAt, r)
***REMOVED***

// FindNextMatch returns the next match in the same input string as the match parameter.
// Will return nil if there is no next match or if given a nil match.
func (re *Regexp) FindNextMatch(m *Match) (*Match, error) ***REMOVED***
	if m == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// If previous match was empty, advance by one before matching to prevent
	// infinite loop
	startAt := m.textpos
	if m.Length == 0 ***REMOVED***
		if m.textpos == len(m.text) ***REMOVED***
			return nil, nil
		***REMOVED***

		if re.RightToLeft() ***REMOVED***
			startAt--
		***REMOVED*** else ***REMOVED***
			startAt++
		***REMOVED***
	***REMOVED***
	return re.run(false, startAt, m.text)
***REMOVED***

// MatchString return true if the string matches the regex
// error will be set if a timeout occurs
func (re *Regexp) MatchString(s string) (bool, error) ***REMOVED***
	m, err := re.run(true, -1, getRunes(s))
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return m != nil, nil
***REMOVED***

func (re *Regexp) getRunesAndStart(s string, startAt int) ([]rune, int) ***REMOVED***
	if startAt < 0 ***REMOVED***
		if re.RightToLeft() ***REMOVED***
			r := getRunes(s)
			return r, len(r)
		***REMOVED***
		return getRunes(s), 0
	***REMOVED***
	ret := make([]rune, len(s))
	i := 0
	runeIdx := -1
	for strIdx, r := range s ***REMOVED***
		if strIdx == startAt ***REMOVED***
			runeIdx = i
		***REMOVED***
		ret[i] = r
		i++
	***REMOVED***
	if startAt == len(s) ***REMOVED***
		runeIdx = i
	***REMOVED***
	return ret[:i], runeIdx
***REMOVED***

func getRunes(s string) []rune ***REMOVED***
	return []rune(s)
***REMOVED***

// MatchRunes return true if the runes matches the regex
// error will be set if a timeout occurs
func (re *Regexp) MatchRunes(r []rune) (bool, error) ***REMOVED***
	m, err := re.run(true, -1, r)
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	return m != nil, nil
***REMOVED***

// GetGroupNames Returns the set of strings used to name capturing groups in the expression.
func (re *Regexp) GetGroupNames() []string ***REMOVED***
	var result []string

	if re.capslist == nil ***REMOVED***
		result = make([]string, re.capsize)

		for i := 0; i < len(result); i++ ***REMOVED***
			result[i] = strconv.Itoa(i)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		result = make([]string, len(re.capslist))
		copy(result, re.capslist)
	***REMOVED***

	return result
***REMOVED***

// GetGroupNumbers returns the integer group numbers corresponding to a group name.
func (re *Regexp) GetGroupNumbers() []int ***REMOVED***
	var result []int

	if re.caps == nil ***REMOVED***
		result = make([]int, re.capsize)

		for i := 0; i < len(result); i++ ***REMOVED***
			result[i] = i
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		result = make([]int, len(re.caps))

		for k, v := range re.caps ***REMOVED***
			result[v] = k
		***REMOVED***
	***REMOVED***

	return result
***REMOVED***

// GroupNameFromNumber retrieves a group name that corresponds to a group number.
// It will return "" for and unknown group number.  Unnamed groups automatically
// receive a name that is the decimal string equivalent of its number.
func (re *Regexp) GroupNameFromNumber(i int) string ***REMOVED***
	if re.capslist == nil ***REMOVED***
		if i >= 0 && i < re.capsize ***REMOVED***
			return strconv.Itoa(i)
		***REMOVED***

		return ""
	***REMOVED***

	if re.caps != nil ***REMOVED***
		var ok bool
		if i, ok = re.caps[i]; !ok ***REMOVED***
			return ""
		***REMOVED***
	***REMOVED***

	if i >= 0 && i < len(re.capslist) ***REMOVED***
		return re.capslist[i]
	***REMOVED***

	return ""
***REMOVED***

// GroupNumberFromName returns a group number that corresponds to a group name.
// Returns -1 if the name is not a recognized group name.  Numbered groups
// automatically get a group name that is the decimal string equivalent of its number.
func (re *Regexp) GroupNumberFromName(name string) int ***REMOVED***
	// look up name if we have a hashtable of names
	if re.capnames != nil ***REMOVED***
		if k, ok := re.capnames[name]; ok ***REMOVED***
			return k
		***REMOVED***

		return -1
	***REMOVED***

	// convert to an int if it looks like a number
	result := 0
	for i := 0; i < len(name); i++ ***REMOVED***
		ch := name[i]

		if ch > '9' || ch < '0' ***REMOVED***
			return -1
		***REMOVED***

		result *= 10
		result += int(ch - '0')
	***REMOVED***

	// return int if it's in range
	if result >= 0 && result < re.capsize ***REMOVED***
		return result
	***REMOVED***

	return -1
***REMOVED***
