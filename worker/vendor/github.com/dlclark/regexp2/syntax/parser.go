package syntax

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"unicode"
)

type RegexOptions int32

const (
	IgnoreCase              RegexOptions = 0x0001 // "i"
	Multiline                            = 0x0002 // "m"
	ExplicitCapture                      = 0x0004 // "n"
	Compiled                             = 0x0008 // "c"
	Singleline                           = 0x0010 // "s"
	IgnorePatternWhitespace              = 0x0020 // "x"
	RightToLeft                          = 0x0040 // "r"
	Debug                                = 0x0080 // "d"
	ECMAScript                           = 0x0100 // "e"
	RE2                                  = 0x0200 // RE2 compat mode
	Unicode                              = 0x0400 // "u"
)

func optionFromCode(ch rune) RegexOptions ***REMOVED***
	// case-insensitive
	switch ch ***REMOVED***
	case 'i', 'I':
		return IgnoreCase
	case 'r', 'R':
		return RightToLeft
	case 'm', 'M':
		return Multiline
	case 'n', 'N':
		return ExplicitCapture
	case 's', 'S':
		return Singleline
	case 'x', 'X':
		return IgnorePatternWhitespace
	case 'd', 'D':
		return Debug
	case 'e', 'E':
		return ECMAScript
	case 'u', 'U':
		return Unicode
	default:
		return 0
	***REMOVED***
***REMOVED***

// An Error describes a failure to parse a regular expression
// and gives the offending expression.
type Error struct ***REMOVED***
	Code ErrorCode
	Expr string
	Args []interface***REMOVED******REMOVED***
***REMOVED***

func (e *Error) Error() string ***REMOVED***
	if len(e.Args) == 0 ***REMOVED***
		return "error parsing regexp: " + e.Code.String() + " in `" + e.Expr + "`"
	***REMOVED***
	return "error parsing regexp: " + fmt.Sprintf(e.Code.String(), e.Args...) + " in `" + e.Expr + "`"
***REMOVED***

// An ErrorCode describes a failure to parse a regular expression.
type ErrorCode string

const (
	// internal issue
	ErrInternalError ErrorCode = "regexp/syntax: internal error"
	// Parser errors
	ErrUnterminatedComment        = "unterminated comment"
	ErrInvalidCharRange           = "invalid character class range"
	ErrInvalidRepeatSize          = "invalid repeat count"
	ErrInvalidUTF8                = "invalid UTF-8"
	ErrCaptureGroupOutOfRange     = "capture group number out of range"
	ErrUnexpectedParen            = "unexpected )"
	ErrMissingParen               = "missing closing )"
	ErrMissingBrace               = "missing closing ***REMOVED***"
	ErrInvalidRepeatOp            = "invalid nested repetition operator"
	ErrMissingRepeatArgument      = "missing argument to repetition operator"
	ErrConditionalExpression      = "illegal conditional (?(...)) expression"
	ErrTooManyAlternates          = "too many | in (?()|)"
	ErrUnrecognizedGrouping       = "unrecognized grouping construct: (%v"
	ErrInvalidGroupName           = "invalid group name: group names must begin with a word character and have a matching terminator"
	ErrCapNumNotZero              = "capture number cannot be zero"
	ErrUndefinedBackRef           = "reference to undefined group number %v"
	ErrUndefinedNameRef           = "reference to undefined group name %v"
	ErrAlternationCantCapture     = "alternation conditions do not capture and cannot be named"
	ErrAlternationCantHaveComment = "alternation conditions cannot be comments"
	ErrMalformedReference         = "(?(%v) ) malformed"
	ErrUndefinedReference         = "(?(%v) ) reference to undefined group"
	ErrIllegalEndEscape           = "illegal \\ at end of pattern"
	ErrMalformedSlashP            = "malformed \\p***REMOVED***X***REMOVED*** character escape"
	ErrIncompleteSlashP           = "incomplete \\p***REMOVED***X***REMOVED*** character escape"
	ErrUnknownSlashP              = "unknown unicode category, script, or property '%v'"
	ErrUnrecognizedEscape         = "unrecognized escape sequence \\%v"
	ErrMissingControl             = "missing control character"
	ErrUnrecognizedControl        = "unrecognized control character"
	ErrTooFewHex                  = "insufficient hexadecimal digits"
	ErrInvalidHex                 = "hex values may not be larger than 0x10FFFF"
	ErrMalformedNameRef           = "malformed \\k<...> named back reference"
	ErrBadClassInCharRange        = "cannot include class \\%v in character range"
	ErrUnterminatedBracket        = "unterminated [] set"
	ErrSubtractionMustBeLast      = "a subtraction must be the last element in a character class"
	ErrReversedCharRange          = "[%c-%c] range in reverse order"
)

func (e ErrorCode) String() string ***REMOVED***
	return string(e)
***REMOVED***

type parser struct ***REMOVED***
	stack         *regexNode
	group         *regexNode
	alternation   *regexNode
	concatenation *regexNode
	unit          *regexNode

	patternRaw string
	pattern    []rune

	currentPos  int
	specialCase *unicode.SpecialCase

	autocap  int
	capcount int
	captop   int
	capsize  int

	caps     map[int]int
	capnames map[string]int

	capnumlist  []int
	capnamelist []string

	options         RegexOptions
	optionsStack    []RegexOptions
	ignoreNextParen bool
***REMOVED***

const (
	maxValueDiv10 int = math.MaxInt32 / 10
	maxValueMod10     = math.MaxInt32 % 10
)

// Parse converts a regex string into a parse tree
func Parse(re string, op RegexOptions) (*RegexTree, error) ***REMOVED***
	p := parser***REMOVED***
		options: op,
		caps:    make(map[int]int),
	***REMOVED***
	p.setPattern(re)

	if err := p.countCaptures(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p.reset(op)
	root, err := p.scanRegex()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	tree := &RegexTree***REMOVED***
		root:       root,
		caps:       p.caps,
		capnumlist: p.capnumlist,
		captop:     p.captop,
		Capnames:   p.capnames,
		Caplist:    p.capnamelist,
		options:    op,
	***REMOVED***

	if tree.options&Debug > 0 ***REMOVED***
		os.Stdout.WriteString(tree.Dump())
	***REMOVED***

	return tree, nil
***REMOVED***

func (p *parser) setPattern(pattern string) ***REMOVED***
	p.patternRaw = pattern
	p.pattern = make([]rune, 0, len(pattern))

	//populate our rune array to handle utf8 encoding
	for _, r := range pattern ***REMOVED***
		p.pattern = append(p.pattern, r)
	***REMOVED***
***REMOVED***
func (p *parser) getErr(code ErrorCode, args ...interface***REMOVED******REMOVED***) error ***REMOVED***
	return &Error***REMOVED***Code: code, Expr: p.patternRaw, Args: args***REMOVED***
***REMOVED***

func (p *parser) noteCaptureSlot(i, pos int) ***REMOVED***
	if _, ok := p.caps[i]; !ok ***REMOVED***
		// the rhs of the hashtable isn't used in the parser
		p.caps[i] = pos
		p.capcount++

		if p.captop <= i ***REMOVED***
			if i == math.MaxInt32 ***REMOVED***
				p.captop = i
			***REMOVED*** else ***REMOVED***
				p.captop = i + 1
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *parser) noteCaptureName(name string, pos int) ***REMOVED***
	if p.capnames == nil ***REMOVED***
		p.capnames = make(map[string]int)
	***REMOVED***

	if _, ok := p.capnames[name]; !ok ***REMOVED***
		p.capnames[name] = pos
		p.capnamelist = append(p.capnamelist, name)
	***REMOVED***
***REMOVED***

func (p *parser) assignNameSlots() ***REMOVED***
	if p.capnames != nil ***REMOVED***
		for _, name := range p.capnamelist ***REMOVED***
			for p.isCaptureSlot(p.autocap) ***REMOVED***
				p.autocap++
			***REMOVED***
			pos := p.capnames[name]
			p.capnames[name] = p.autocap
			p.noteCaptureSlot(p.autocap, pos)

			p.autocap++
		***REMOVED***
	***REMOVED***

	// if the caps array has at least one gap, construct the list of used slots
	if p.capcount < p.captop ***REMOVED***
		p.capnumlist = make([]int, p.capcount)
		i := 0

		for k := range p.caps ***REMOVED***
			p.capnumlist[i] = k
			i++
		***REMOVED***

		sort.Ints(p.capnumlist)
	***REMOVED***

	// merge capsnumlist into capnamelist
	if p.capnames != nil || p.capnumlist != nil ***REMOVED***
		var oldcapnamelist []string
		var next int
		var k int

		if p.capnames == nil ***REMOVED***
			oldcapnamelist = nil
			p.capnames = make(map[string]int)
			p.capnamelist = []string***REMOVED******REMOVED***
			next = -1
		***REMOVED*** else ***REMOVED***
			oldcapnamelist = p.capnamelist
			p.capnamelist = []string***REMOVED******REMOVED***
			next = p.capnames[oldcapnamelist[0]]
		***REMOVED***

		for i := 0; i < p.capcount; i++ ***REMOVED***
			j := i
			if p.capnumlist != nil ***REMOVED***
				j = p.capnumlist[i]
			***REMOVED***

			if next == j ***REMOVED***
				p.capnamelist = append(p.capnamelist, oldcapnamelist[k])
				k++

				if k == len(oldcapnamelist) ***REMOVED***
					next = -1
				***REMOVED*** else ***REMOVED***
					next = p.capnames[oldcapnamelist[k]]
				***REMOVED***

			***REMOVED*** else ***REMOVED***
				//feature: culture?
				str := strconv.Itoa(j)
				p.capnamelist = append(p.capnamelist, str)
				p.capnames[str] = j
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *parser) consumeAutocap() int ***REMOVED***
	r := p.autocap
	p.autocap++
	return r
***REMOVED***

// CountCaptures is a prescanner for deducing the slots used for
// captures by doing a partial tokenization of the pattern.
func (p *parser) countCaptures() error ***REMOVED***
	var ch rune

	p.noteCaptureSlot(0, 0)

	p.autocap = 1

	for p.charsRight() > 0 ***REMOVED***
		pos := p.textpos()
		ch = p.moveRightGetChar()
		switch ch ***REMOVED***
		case '\\':
			if p.charsRight() > 0 ***REMOVED***
				p.scanBackslash(true)
			***REMOVED***

		case '#':
			if p.useOptionX() ***REMOVED***
				p.moveLeft()
				p.scanBlank()
			***REMOVED***

		case '[':
			p.scanCharSet(false, true)

		case ')':
			if !p.emptyOptionsStack() ***REMOVED***
				p.popOptions()
			***REMOVED***

		case '(':
			if p.charsRight() >= 2 && p.rightChar(1) == '#' && p.rightChar(0) == '?' ***REMOVED***
				p.moveLeft()
				p.scanBlank()
			***REMOVED*** else ***REMOVED***
				p.pushOptions()
				if p.charsRight() > 0 && p.rightChar(0) == '?' ***REMOVED***
					// we have (?...
					p.moveRight(1)

					if p.charsRight() > 1 && (p.rightChar(0) == '<' || p.rightChar(0) == '\'') ***REMOVED***
						// named group: (?<... or (?'...

						p.moveRight(1)
						ch = p.rightChar(0)

						if ch != '0' && IsWordChar(ch) ***REMOVED***
							if ch >= '1' && ch <= '9' ***REMOVED***
								dec, err := p.scanDecimal()
								if err != nil ***REMOVED***
									return err
								***REMOVED***
								p.noteCaptureSlot(dec, pos)
							***REMOVED*** else ***REMOVED***
								p.noteCaptureName(p.scanCapname(), pos)
							***REMOVED***
						***REMOVED***
					***REMOVED*** else if p.useRE2() && p.charsRight() > 2 && (p.rightChar(0) == 'P' && p.rightChar(1) == '<') ***REMOVED***
						// RE2-compat (?P<)
						p.moveRight(2)
						ch = p.rightChar(0)
						if IsWordChar(ch) ***REMOVED***
							p.noteCaptureName(p.scanCapname(), pos)
						***REMOVED***

					***REMOVED*** else ***REMOVED***
						// (?...

						// get the options if it's an option construct (?cimsx-cimsx...)
						p.scanOptions()

						if p.charsRight() > 0 ***REMOVED***
							if p.rightChar(0) == ')' ***REMOVED***
								// (?cimsx-cimsx)
								p.moveRight(1)
								p.popKeepOptions()
							***REMOVED*** else if p.rightChar(0) == '(' ***REMOVED***
								// alternation construct: (?(foo)yes|no)
								// ignore the next paren so we don't capture the condition
								p.ignoreNextParen = true

								// break from here so we don't reset ignoreNextParen
								continue
							***REMOVED***
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					if !p.useOptionN() && !p.ignoreNextParen ***REMOVED***
						p.noteCaptureSlot(p.consumeAutocap(), pos)
					***REMOVED***
				***REMOVED***
			***REMOVED***

			p.ignoreNextParen = false

		***REMOVED***
	***REMOVED***

	p.assignNameSlots()
	return nil
***REMOVED***

func (p *parser) reset(topopts RegexOptions) ***REMOVED***
	p.currentPos = 0
	p.autocap = 1
	p.ignoreNextParen = false

	if len(p.optionsStack) > 0 ***REMOVED***
		p.optionsStack = p.optionsStack[:0]
	***REMOVED***

	p.options = topopts
	p.stack = nil
***REMOVED***

func (p *parser) scanRegex() (*regexNode, error) ***REMOVED***
	ch := '@' // nonspecial ch, means at beginning
	isQuant := false

	p.startGroup(newRegexNodeMN(ntCapture, p.options, 0, -1))

	for p.charsRight() > 0 ***REMOVED***
		wasPrevQuantifier := isQuant
		isQuant = false

		if err := p.scanBlank(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		startpos := p.textpos()

		// move past all of the normal characters.  We'll stop when we hit some kind of control character,
		// or if IgnorePatternWhiteSpace is on, we'll stop when we see some whitespace.
		if p.useOptionX() ***REMOVED***
			for p.charsRight() > 0 ***REMOVED***
				ch = p.rightChar(0)
				//UGLY: clean up, this is ugly
				if !(!isStopperX(ch) || (ch == '***REMOVED***' && !p.isTrueQuantifier())) ***REMOVED***
					break
				***REMOVED***
				p.moveRight(1)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for p.charsRight() > 0 ***REMOVED***
				ch = p.rightChar(0)
				if !(!isSpecial(ch) || ch == '***REMOVED***' && !p.isTrueQuantifier()) ***REMOVED***
					break
				***REMOVED***
				p.moveRight(1)
			***REMOVED***
		***REMOVED***

		endpos := p.textpos()

		p.scanBlank()

		if p.charsRight() == 0 ***REMOVED***
			ch = '!' // nonspecial, means at end
		***REMOVED*** else if ch = p.rightChar(0); isSpecial(ch) ***REMOVED***
			isQuant = isQuantifier(ch)
			p.moveRight(1)
		***REMOVED*** else ***REMOVED***
			ch = ' ' // nonspecial, means at ordinary char
		***REMOVED***

		if startpos < endpos ***REMOVED***
			cchUnquantified := endpos - startpos
			if isQuant ***REMOVED***
				cchUnquantified--
			***REMOVED***
			wasPrevQuantifier = false

			if cchUnquantified > 0 ***REMOVED***
				p.addToConcatenate(startpos, cchUnquantified, false)
			***REMOVED***

			if isQuant ***REMOVED***
				p.addUnitOne(p.charAt(endpos - 1))
			***REMOVED***
		***REMOVED***

		switch ch ***REMOVED***
		case '!':
			goto BreakOuterScan

		case ' ':
			goto ContinueOuterScan

		case '[':
			cc, err := p.scanCharSet(p.useOptionI(), false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			p.addUnitSet(cc)

		case '(':
			p.pushOptions()

			if grouper, err := p.scanGroupOpen(); err != nil ***REMOVED***
				return nil, err
			***REMOVED*** else if grouper == nil ***REMOVED***
				p.popKeepOptions()
			***REMOVED*** else ***REMOVED***
				p.pushGroup()
				p.startGroup(grouper)
			***REMOVED***

			continue

		case '|':
			p.addAlternate()
			goto ContinueOuterScan

		case ')':
			if p.emptyStack() ***REMOVED***
				return nil, p.getErr(ErrUnexpectedParen)
			***REMOVED***

			if err := p.addGroup(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if err := p.popGroup(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			p.popOptions()

			if p.unit == nil ***REMOVED***
				goto ContinueOuterScan
			***REMOVED***

		case '\\':
			n, err := p.scanBackslash(false)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			p.addUnitNode(n)

		case '^':
			if p.useOptionM() ***REMOVED***
				p.addUnitType(ntBol)
			***REMOVED*** else ***REMOVED***
				p.addUnitType(ntBeginning)
			***REMOVED***

		case '$':
			if p.useOptionM() ***REMOVED***
				p.addUnitType(ntEol)
			***REMOVED*** else ***REMOVED***
				p.addUnitType(ntEndZ)
			***REMOVED***

		case '.':
			if p.useOptionE() ***REMOVED***
				p.addUnitSet(ECMAAnyClass())
			***REMOVED*** else if p.useOptionS() ***REMOVED***
				p.addUnitSet(AnyClass())
			***REMOVED*** else ***REMOVED***
				p.addUnitNotone('\n')
			***REMOVED***

		case '***REMOVED***', '*', '+', '?':
			if p.unit == nil ***REMOVED***
				if wasPrevQuantifier ***REMOVED***
					return nil, p.getErr(ErrInvalidRepeatOp)
				***REMOVED*** else ***REMOVED***
					return nil, p.getErr(ErrMissingRepeatArgument)
				***REMOVED***
			***REMOVED***
			p.moveLeft()

		default:
			return nil, p.getErr(ErrInternalError)
		***REMOVED***

		if err := p.scanBlank(); err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if p.charsRight() > 0 ***REMOVED***
			isQuant = p.isTrueQuantifier()
		***REMOVED***
		if p.charsRight() == 0 || !isQuant ***REMOVED***
			//maintain odd C# assignment order -- not sure if required, could clean up?
			p.addConcatenate()
			goto ContinueOuterScan
		***REMOVED***

		ch = p.moveRightGetChar()

		// Handle quantifiers
		for p.unit != nil ***REMOVED***
			var min, max int
			var lazy bool

			switch ch ***REMOVED***
			case '*':
				min = 0
				max = math.MaxInt32

			case '?':
				min = 0
				max = 1

			case '+':
				min = 1
				max = math.MaxInt32

			case '***REMOVED***':
				***REMOVED***
					var err error
					startpos = p.textpos()
					if min, err = p.scanDecimal(); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					max = min
					if startpos < p.textpos() ***REMOVED***
						if p.charsRight() > 0 && p.rightChar(0) == ',' ***REMOVED***
							p.moveRight(1)
							if p.charsRight() == 0 || p.rightChar(0) == '***REMOVED***' ***REMOVED***
								max = math.MaxInt32
							***REMOVED*** else ***REMOVED***
								if max, err = p.scanDecimal(); err != nil ***REMOVED***
									return nil, err
								***REMOVED***
							***REMOVED***
						***REMOVED***
					***REMOVED***

					if startpos == p.textpos() || p.charsRight() == 0 || p.moveRightGetChar() != '***REMOVED***' ***REMOVED***
						p.addConcatenate()
						p.textto(startpos - 1)
						goto ContinueOuterScan
					***REMOVED***
				***REMOVED***

			default:
				return nil, p.getErr(ErrInternalError)
			***REMOVED***

			if err := p.scanBlank(); err != nil ***REMOVED***
				return nil, err
			***REMOVED***

			if p.charsRight() == 0 || p.rightChar(0) != '?' ***REMOVED***
				lazy = false
			***REMOVED*** else ***REMOVED***
				p.moveRight(1)
				lazy = true
			***REMOVED***

			if min > max ***REMOVED***
				return nil, p.getErr(ErrInvalidRepeatSize)
			***REMOVED***

			p.addConcatenate3(lazy, min, max)
		***REMOVED***

	ContinueOuterScan:
	***REMOVED***

BreakOuterScan:
	;

	if !p.emptyStack() ***REMOVED***
		return nil, p.getErr(ErrMissingParen)
	***REMOVED***

	if err := p.addGroup(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return p.unit, nil

***REMOVED***

/*
 * Simple parsing for replacement patterns
 */
func (p *parser) scanReplacement() (*regexNode, error) ***REMOVED***
	var c, startpos int

	p.concatenation = newRegexNode(ntConcatenate, p.options)

	for ***REMOVED***
		c = p.charsRight()
		if c == 0 ***REMOVED***
			break
		***REMOVED***

		startpos = p.textpos()

		for c > 0 && p.rightChar(0) != '$' ***REMOVED***
			p.moveRight(1)
			c--
		***REMOVED***

		p.addToConcatenate(startpos, p.textpos()-startpos, true)

		if c > 0 ***REMOVED***
			if p.moveRightGetChar() == '$' ***REMOVED***
				n, err := p.scanDollar()
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				p.addUnitNode(n)
			***REMOVED***
			p.addConcatenate()
		***REMOVED***
	***REMOVED***

	return p.concatenation, nil
***REMOVED***

/*
 * Scans $ patterns recognized within replacement patterns
 */
func (p *parser) scanDollar() (*regexNode, error) ***REMOVED***
	if p.charsRight() == 0 ***REMOVED***
		return newRegexNodeCh(ntOne, p.options, '$'), nil
	***REMOVED***

	ch := p.rightChar(0)
	angled := false
	backpos := p.textpos()
	lastEndPos := backpos

	// Note angle

	if ch == '***REMOVED***' && p.charsRight() > 1 ***REMOVED***
		angled = true
		p.moveRight(1)
		ch = p.rightChar(0)
	***REMOVED***

	// Try to parse backreference: \1 or \***REMOVED***1***REMOVED*** or \***REMOVED***cap***REMOVED***

	if ch >= '0' && ch <= '9' ***REMOVED***
		if !angled && p.useOptionE() ***REMOVED***
			capnum := -1
			newcapnum := int(ch - '0')
			p.moveRight(1)
			if p.isCaptureSlot(newcapnum) ***REMOVED***
				capnum = newcapnum
				lastEndPos = p.textpos()
			***REMOVED***

			for p.charsRight() > 0 ***REMOVED***
				ch = p.rightChar(0)
				if ch < '0' || ch > '9' ***REMOVED***
					break
				***REMOVED***
				digit := int(ch - '0')
				if newcapnum > maxValueDiv10 || (newcapnum == maxValueDiv10 && digit > maxValueMod10) ***REMOVED***
					return nil, p.getErr(ErrCaptureGroupOutOfRange)
				***REMOVED***

				newcapnum = newcapnum*10 + digit

				p.moveRight(1)
				if p.isCaptureSlot(newcapnum) ***REMOVED***
					capnum = newcapnum
					lastEndPos = p.textpos()
				***REMOVED***
			***REMOVED***
			p.textto(lastEndPos)
			if capnum >= 0 ***REMOVED***
				return newRegexNodeM(ntRef, p.options, capnum), nil
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			capnum, err := p.scanDecimal()
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			if !angled || p.charsRight() > 0 && p.moveRightGetChar() == '***REMOVED***' ***REMOVED***
				if p.isCaptureSlot(capnum) ***REMOVED***
					return newRegexNodeM(ntRef, p.options, capnum), nil
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if angled && IsWordChar(ch) ***REMOVED***
		capname := p.scanCapname()

		if p.charsRight() > 0 && p.moveRightGetChar() == '***REMOVED***' ***REMOVED***
			if p.isCaptureName(capname) ***REMOVED***
				return newRegexNodeM(ntRef, p.options, p.captureSlotFromName(capname)), nil
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if !angled ***REMOVED***
		capnum := 1

		switch ch ***REMOVED***
		case '$':
			p.moveRight(1)
			return newRegexNodeCh(ntOne, p.options, '$'), nil
		case '&':
			capnum = 0
		case '`':
			capnum = replaceLeftPortion
		case '\'':
			capnum = replaceRightPortion
		case '+':
			capnum = replaceLastGroup
		case '_':
			capnum = replaceWholeString
		***REMOVED***

		if capnum != 1 ***REMOVED***
			p.moveRight(1)
			return newRegexNodeM(ntRef, p.options, capnum), nil
		***REMOVED***
	***REMOVED***

	// unrecognized $: literalize

	p.textto(backpos)
	return newRegexNodeCh(ntOne, p.options, '$'), nil
***REMOVED***

// scanGroupOpen scans chars following a '(' (not counting the '('), and returns
// a RegexNode for the type of group scanned, or nil if the group
// simply changed options (?cimsx-cimsx) or was a comment (#...).
func (p *parser) scanGroupOpen() (*regexNode, error) ***REMOVED***
	var ch rune
	var nt nodeType
	var err error
	close := '>'
	start := p.textpos()

	// just return a RegexNode if we have:
	// 1. "(" followed by nothing
	// 2. "(x" where x != ?
	// 3. "(?)"
	if p.charsRight() == 0 || p.rightChar(0) != '?' || (p.rightChar(0) == '?' && (p.charsRight() > 1 && p.rightChar(1) == ')')) ***REMOVED***
		if p.useOptionN() || p.ignoreNextParen ***REMOVED***
			p.ignoreNextParen = false
			return newRegexNode(ntGroup, p.options), nil
		***REMOVED***
		return newRegexNodeMN(ntCapture, p.options, p.consumeAutocap(), -1), nil
	***REMOVED***

	p.moveRight(1)

	for ***REMOVED***
		if p.charsRight() == 0 ***REMOVED***
			break
		***REMOVED***

		switch ch = p.moveRightGetChar(); ch ***REMOVED***
		case ':':
			nt = ntGroup

		case '=':
			p.options &= ^RightToLeft
			nt = ntRequire

		case '!':
			p.options &= ^RightToLeft
			nt = ntPrevent

		case '>':
			nt = ntGreedy

		case '\'':
			close = '\''
			fallthrough

		case '<':
			if p.charsRight() == 0 ***REMOVED***
				goto BreakRecognize
			***REMOVED***

			switch ch = p.moveRightGetChar(); ch ***REMOVED***
			case '=':
				if close == '\'' ***REMOVED***
					goto BreakRecognize
				***REMOVED***

				p.options |= RightToLeft
				nt = ntRequire

			case '!':
				if close == '\'' ***REMOVED***
					goto BreakRecognize
				***REMOVED***

				p.options |= RightToLeft
				nt = ntPrevent

			default:
				p.moveLeft()
				capnum := -1
				uncapnum := -1
				proceed := false

				// grab part before -

				if ch >= '0' && ch <= '9' ***REMOVED***
					if capnum, err = p.scanDecimal(); err != nil ***REMOVED***
						return nil, err
					***REMOVED***

					if !p.isCaptureSlot(capnum) ***REMOVED***
						capnum = -1
					***REMOVED***

					// check if we have bogus characters after the number
					if p.charsRight() > 0 && !(p.rightChar(0) == close || p.rightChar(0) == '-') ***REMOVED***
						return nil, p.getErr(ErrInvalidGroupName)
					***REMOVED***
					if capnum == 0 ***REMOVED***
						return nil, p.getErr(ErrCapNumNotZero)
					***REMOVED***
				***REMOVED*** else if IsWordChar(ch) ***REMOVED***
					capname := p.scanCapname()

					if p.isCaptureName(capname) ***REMOVED***
						capnum = p.captureSlotFromName(capname)
					***REMOVED***

					// check if we have bogus character after the name
					if p.charsRight() > 0 && !(p.rightChar(0) == close || p.rightChar(0) == '-') ***REMOVED***
						return nil, p.getErr(ErrInvalidGroupName)
					***REMOVED***
				***REMOVED*** else if ch == '-' ***REMOVED***
					proceed = true
				***REMOVED*** else ***REMOVED***
					// bad group name - starts with something other than a word character and isn't a number
					return nil, p.getErr(ErrInvalidGroupName)
				***REMOVED***

				// grab part after - if any

				if (capnum != -1 || proceed == true) && p.charsRight() > 0 && p.rightChar(0) == '-' ***REMOVED***
					p.moveRight(1)

					//no more chars left, no closing char, etc
					if p.charsRight() == 0 ***REMOVED***
						return nil, p.getErr(ErrInvalidGroupName)
					***REMOVED***

					ch = p.rightChar(0)
					if ch >= '0' && ch <= '9' ***REMOVED***
						if uncapnum, err = p.scanDecimal(); err != nil ***REMOVED***
							return nil, err
						***REMOVED***

						if !p.isCaptureSlot(uncapnum) ***REMOVED***
							return nil, p.getErr(ErrUndefinedBackRef, uncapnum)
						***REMOVED***

						// check if we have bogus characters after the number
						if p.charsRight() > 0 && p.rightChar(0) != close ***REMOVED***
							return nil, p.getErr(ErrInvalidGroupName)
						***REMOVED***
					***REMOVED*** else if IsWordChar(ch) ***REMOVED***
						uncapname := p.scanCapname()

						if !p.isCaptureName(uncapname) ***REMOVED***
							return nil, p.getErr(ErrUndefinedNameRef, uncapname)
						***REMOVED***
						uncapnum = p.captureSlotFromName(uncapname)

						// check if we have bogus character after the name
						if p.charsRight() > 0 && p.rightChar(0) != close ***REMOVED***
							return nil, p.getErr(ErrInvalidGroupName)
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						// bad group name - starts with something other than a word character and isn't a number
						return nil, p.getErr(ErrInvalidGroupName)
					***REMOVED***
				***REMOVED***

				// actually make the node

				if (capnum != -1 || uncapnum != -1) && p.charsRight() > 0 && p.moveRightGetChar() == close ***REMOVED***
					return newRegexNodeMN(ntCapture, p.options, capnum, uncapnum), nil
				***REMOVED***
				goto BreakRecognize
			***REMOVED***

		case '(':
			// alternation construct (?(...) | )

			parenPos := p.textpos()
			if p.charsRight() > 0 ***REMOVED***
				ch = p.rightChar(0)

				// check if the alternation condition is a backref
				if ch >= '0' && ch <= '9' ***REMOVED***
					var capnum int
					if capnum, err = p.scanDecimal(); err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					if p.charsRight() > 0 && p.moveRightGetChar() == ')' ***REMOVED***
						if p.isCaptureSlot(capnum) ***REMOVED***
							return newRegexNodeM(ntTestref, p.options, capnum), nil
						***REMOVED***
						return nil, p.getErr(ErrUndefinedReference, capnum)
					***REMOVED***

					return nil, p.getErr(ErrMalformedReference, capnum)

				***REMOVED*** else if IsWordChar(ch) ***REMOVED***
					capname := p.scanCapname()

					if p.isCaptureName(capname) && p.charsRight() > 0 && p.moveRightGetChar() == ')' ***REMOVED***
						return newRegexNodeM(ntTestref, p.options, p.captureSlotFromName(capname)), nil
					***REMOVED***
				***REMOVED***
			***REMOVED***
			// not a backref
			nt = ntTestgroup
			p.textto(parenPos - 1)   // jump to the start of the parentheses
			p.ignoreNextParen = true // but make sure we don't try to capture the insides

			charsRight := p.charsRight()
			if charsRight >= 3 && p.rightChar(1) == '?' ***REMOVED***
				rightchar2 := p.rightChar(2)
				// disallow comments in the condition
				if rightchar2 == '#' ***REMOVED***
					return nil, p.getErr(ErrAlternationCantHaveComment)
				***REMOVED***

				// disallow named capture group (?<..>..) in the condition
				if rightchar2 == '\'' ***REMOVED***
					return nil, p.getErr(ErrAlternationCantCapture)
				***REMOVED***

				if charsRight >= 4 && (rightchar2 == '<' && p.rightChar(3) != '!' && p.rightChar(3) != '=') ***REMOVED***
					return nil, p.getErr(ErrAlternationCantCapture)
				***REMOVED***
			***REMOVED***

		case 'P':
			if p.useRE2() ***REMOVED***
				// support for P<name> syntax
				if p.charsRight() < 3 ***REMOVED***
					goto BreakRecognize
				***REMOVED***

				ch = p.moveRightGetChar()
				if ch != '<' ***REMOVED***
					goto BreakRecognize
				***REMOVED***

				ch = p.moveRightGetChar()
				p.moveLeft()

				if IsWordChar(ch) ***REMOVED***
					capnum := -1
					capname := p.scanCapname()

					if p.isCaptureName(capname) ***REMOVED***
						capnum = p.captureSlotFromName(capname)
					***REMOVED***

					// check if we have bogus character after the name
					if p.charsRight() > 0 && p.rightChar(0) != '>' ***REMOVED***
						return nil, p.getErr(ErrInvalidGroupName)
					***REMOVED***

					// actually make the node

					if capnum != -1 && p.charsRight() > 0 && p.moveRightGetChar() == '>' ***REMOVED***
						return newRegexNodeMN(ntCapture, p.options, capnum, -1), nil
					***REMOVED***
					goto BreakRecognize

				***REMOVED*** else ***REMOVED***
					// bad group name - starts with something other than a word character and isn't a number
					return nil, p.getErr(ErrInvalidGroupName)
				***REMOVED***
			***REMOVED***
			// if we're not using RE2 compat mode then
			// we just behave like normal
			fallthrough

		default:
			p.moveLeft()

			nt = ntGroup
			// disallow options in the children of a testgroup node
			if p.group.t != ntTestgroup ***REMOVED***
				p.scanOptions()
			***REMOVED***
			if p.charsRight() == 0 ***REMOVED***
				goto BreakRecognize
			***REMOVED***

			if ch = p.moveRightGetChar(); ch == ')' ***REMOVED***
				return nil, nil
			***REMOVED***

			if ch != ':' ***REMOVED***
				goto BreakRecognize
			***REMOVED***

		***REMOVED***

		return newRegexNode(nt, p.options), nil
	***REMOVED***

BreakRecognize:

	// break Recognize comes here

	return nil, p.getErr(ErrUnrecognizedGrouping, string(p.pattern[start:p.textpos()]))
***REMOVED***

// scans backslash specials and basics
func (p *parser) scanBackslash(scanOnly bool) (*regexNode, error) ***REMOVED***

	if p.charsRight() == 0 ***REMOVED***
		return nil, p.getErr(ErrIllegalEndEscape)
	***REMOVED***

	switch ch := p.rightChar(0); ch ***REMOVED***
	case 'b', 'B', 'A', 'G', 'Z', 'z':
		p.moveRight(1)
		return newRegexNode(p.typeFromCode(ch), p.options), nil

	case 'w':
		p.moveRight(1)
		if p.useOptionE() || p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, ECMAWordClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, WordClass()), nil

	case 'W':
		p.moveRight(1)
		if p.useOptionE() || p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, NotECMAWordClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, NotWordClass()), nil

	case 's':
		p.moveRight(1)
		if p.useOptionE() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, ECMASpaceClass()), nil
		***REMOVED*** else if p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, RE2SpaceClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, SpaceClass()), nil

	case 'S':
		p.moveRight(1)
		if p.useOptionE() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, NotECMASpaceClass()), nil
		***REMOVED*** else if p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, NotRE2SpaceClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, NotSpaceClass()), nil

	case 'd':
		p.moveRight(1)
		if p.useOptionE() || p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, ECMADigitClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, DigitClass()), nil

	case 'D':
		p.moveRight(1)
		if p.useOptionE() || p.useRE2() ***REMOVED***
			return newRegexNodeSet(ntSet, p.options, NotECMADigitClass()), nil
		***REMOVED***
		return newRegexNodeSet(ntSet, p.options, NotDigitClass()), nil

	case 'p', 'P':
		p.moveRight(1)
		prop, err := p.parseProperty()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		cc := &CharSet***REMOVED******REMOVED***
		cc.addCategory(prop, (ch != 'p'), p.useOptionI(), p.patternRaw)
		if p.useOptionI() ***REMOVED***
			cc.addLowercase()
		***REMOVED***

		return newRegexNodeSet(ntSet, p.options, cc), nil

	default:
		return p.scanBasicBackslash(scanOnly)
	***REMOVED***
***REMOVED***

// Scans \-style backreferences and character escapes
func (p *parser) scanBasicBackslash(scanOnly bool) (*regexNode, error) ***REMOVED***
	if p.charsRight() == 0 ***REMOVED***
		return nil, p.getErr(ErrIllegalEndEscape)
	***REMOVED***
	angled := false
	k := false
	close := '\x00'

	backpos := p.textpos()
	ch := p.rightChar(0)

	// Allow \k<foo> instead of \<foo>, which is now deprecated.

	// According to ECMAScript specification, \k<name> is only parsed as a named group reference if
	// there is at least one group name in the regexp.
	// See https://www.ecma-international.org/ecma-262/#sec-isvalidregularexpressionliteral, step 7.
	// Note, during the first (scanOnly) run we may not have all group names scanned, but that's ok.
	if ch == 'k' && (!p.useOptionE() || len(p.capnames) > 0) ***REMOVED***
		if p.charsRight() >= 2 ***REMOVED***
			p.moveRight(1)
			ch = p.moveRightGetChar()

			if ch == '<' || (!p.useOptionE() && ch == '\'') ***REMOVED*** // No support for \k'name' in ECMAScript
				angled = true
				if ch == '\'' ***REMOVED***
					close = '\''
				***REMOVED*** else ***REMOVED***
					close = '>'
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if !angled || p.charsRight() <= 0 ***REMOVED***
			return nil, p.getErr(ErrMalformedNameRef)
		***REMOVED***

		ch = p.rightChar(0)
		k = true

	***REMOVED*** else if !p.useOptionE() && (ch == '<' || ch == '\'') && p.charsRight() > 1 ***REMOVED*** // Note angle without \g
		angled = true
		if ch == '\'' ***REMOVED***
			close = '\''
		***REMOVED*** else ***REMOVED***
			close = '>'
		***REMOVED***

		p.moveRight(1)
		ch = p.rightChar(0)
	***REMOVED***

	// Try to parse backreference: \<1> or \<cap>

	if angled && ch >= '0' && ch <= '9' ***REMOVED***
		capnum, err := p.scanDecimal()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if p.charsRight() > 0 && p.moveRightGetChar() == close ***REMOVED***
			if p.isCaptureSlot(capnum) ***REMOVED***
				return newRegexNodeM(ntRef, p.options, capnum), nil
			***REMOVED***
			return nil, p.getErr(ErrUndefinedBackRef, capnum)
		***REMOVED***
	***REMOVED*** else if !angled && ch >= '1' && ch <= '9' ***REMOVED*** // Try to parse backreference or octal: \1
		capnum, err := p.scanDecimal()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		if scanOnly ***REMOVED***
			return nil, nil
		***REMOVED***

		if p.isCaptureSlot(capnum) ***REMOVED***
			return newRegexNodeM(ntRef, p.options, capnum), nil
		***REMOVED***
		if capnum <= 9 && !p.useOptionE() ***REMOVED***
			return nil, p.getErr(ErrUndefinedBackRef, capnum)
		***REMOVED***

	***REMOVED*** else if angled ***REMOVED***
		capname := p.scanCapname()

		if capname != "" && p.charsRight() > 0 && p.moveRightGetChar() == close ***REMOVED***

			if scanOnly ***REMOVED***
				return nil, nil
			***REMOVED***

			if p.isCaptureName(capname) ***REMOVED***
				return newRegexNodeM(ntRef, p.options, p.captureSlotFromName(capname)), nil
			***REMOVED***
			return nil, p.getErr(ErrUndefinedNameRef, capname)
		***REMOVED*** else ***REMOVED***
			if k ***REMOVED***
				return nil, p.getErr(ErrMalformedNameRef)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Not backreference: must be char code

	p.textto(backpos)
	ch, err := p.scanCharEscape()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if scanOnly ***REMOVED***
		return nil, nil
	***REMOVED***

	if p.useOptionI() ***REMOVED***
		ch = unicode.ToLower(ch)
	***REMOVED***

	return newRegexNodeCh(ntOne, p.options, ch), nil
***REMOVED***

// Scans X for \p***REMOVED***X***REMOVED*** or \P***REMOVED***X***REMOVED***
func (p *parser) parseProperty() (string, error) ***REMOVED***
	if p.charsRight() < 3 ***REMOVED***
		return "", p.getErr(ErrIncompleteSlashP)
	***REMOVED***
	ch := p.moveRightGetChar()
	if ch != '***REMOVED***' ***REMOVED***
		return "", p.getErr(ErrMalformedSlashP)
	***REMOVED***

	startpos := p.textpos()
	for p.charsRight() > 0 ***REMOVED***
		ch = p.moveRightGetChar()
		if !(IsWordChar(ch) || ch == '-') ***REMOVED***
			p.moveLeft()
			break
		***REMOVED***
	***REMOVED***
	capname := string(p.pattern[startpos:p.textpos()])

	if p.charsRight() == 0 || p.moveRightGetChar() != '***REMOVED***' ***REMOVED***
		return "", p.getErr(ErrIncompleteSlashP)
	***REMOVED***

	if !isValidUnicodeCat(capname) ***REMOVED***
		return "", p.getErr(ErrUnknownSlashP, capname)
	***REMOVED***

	return capname, nil
***REMOVED***

// Returns ReNode type for zero-length assertions with a \ code.
func (p *parser) typeFromCode(ch rune) nodeType ***REMOVED***
	switch ch ***REMOVED***
	case 'b':
		if p.useOptionE() ***REMOVED***
			return ntECMABoundary
		***REMOVED***
		return ntBoundary
	case 'B':
		if p.useOptionE() ***REMOVED***
			return ntNonECMABoundary
		***REMOVED***
		return ntNonboundary
	case 'A':
		return ntBeginning
	case 'G':
		return ntStart
	case 'Z':
		return ntEndZ
	case 'z':
		return ntEnd
	default:
		return ntNothing
	***REMOVED***
***REMOVED***

// Scans whitespace or x-mode comments.
func (p *parser) scanBlank() error ***REMOVED***
	if p.useOptionX() ***REMOVED***
		for ***REMOVED***
			for p.charsRight() > 0 && isSpace(p.rightChar(0)) ***REMOVED***
				p.moveRight(1)
			***REMOVED***

			if p.charsRight() == 0 ***REMOVED***
				break
			***REMOVED***

			if p.rightChar(0) == '#' ***REMOVED***
				for p.charsRight() > 0 && p.rightChar(0) != '\n' ***REMOVED***
					p.moveRight(1)
				***REMOVED***
			***REMOVED*** else if p.charsRight() >= 3 && p.rightChar(2) == '#' &&
				p.rightChar(1) == '?' && p.rightChar(0) == '(' ***REMOVED***
				for p.charsRight() > 0 && p.rightChar(0) != ')' ***REMOVED***
					p.moveRight(1)
				***REMOVED***
				if p.charsRight() == 0 ***REMOVED***
					return p.getErr(ErrUnterminatedComment)
				***REMOVED***
				p.moveRight(1)
			***REMOVED*** else ***REMOVED***
				break
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for ***REMOVED***
			if p.charsRight() < 3 || p.rightChar(2) != '#' ||
				p.rightChar(1) != '?' || p.rightChar(0) != '(' ***REMOVED***
				return nil
			***REMOVED***

			for p.charsRight() > 0 && p.rightChar(0) != ')' ***REMOVED***
				p.moveRight(1)
			***REMOVED***
			if p.charsRight() == 0 ***REMOVED***
				return p.getErr(ErrUnterminatedComment)
			***REMOVED***
			p.moveRight(1)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (p *parser) scanCapname() string ***REMOVED***
	startpos := p.textpos()

	for p.charsRight() > 0 ***REMOVED***
		if !IsWordChar(p.moveRightGetChar()) ***REMOVED***
			p.moveLeft()
			break
		***REMOVED***
	***REMOVED***

	return string(p.pattern[startpos:p.textpos()])
***REMOVED***

//Scans contents of [] (not including []'s), and converts to a set.
func (p *parser) scanCharSet(caseInsensitive, scanOnly bool) (*CharSet, error) ***REMOVED***
	ch := '\x00'
	chPrev := '\x00'
	inRange := false
	firstChar := true
	closed := false

	var cc *CharSet
	if !scanOnly ***REMOVED***
		cc = &CharSet***REMOVED******REMOVED***
	***REMOVED***

	if p.charsRight() > 0 && p.rightChar(0) == '^' ***REMOVED***
		p.moveRight(1)
		if !scanOnly ***REMOVED***
			cc.negate = true
		***REMOVED***
	***REMOVED***

	for ; p.charsRight() > 0; firstChar = false ***REMOVED***
		fTranslatedChar := false
		ch = p.moveRightGetChar()
		if ch == ']' ***REMOVED***
			if !firstChar ***REMOVED***
				closed = true
				break
			***REMOVED*** else if p.useOptionE() ***REMOVED***
				if !scanOnly ***REMOVED***
					cc.addRanges(NoneClass().ranges)
				***REMOVED***
				closed = true
				break
			***REMOVED***

		***REMOVED*** else if ch == '\\' && p.charsRight() > 0 ***REMOVED***
			switch ch = p.moveRightGetChar(); ch ***REMOVED***
			case 'D', 'd':
				if !scanOnly ***REMOVED***
					if inRange ***REMOVED***
						return nil, p.getErr(ErrBadClassInCharRange, ch)
					***REMOVED***
					cc.addDigit(p.useOptionE() || p.useRE2(), ch == 'D', p.patternRaw)
				***REMOVED***
				continue

			case 'S', 's':
				if !scanOnly ***REMOVED***
					if inRange ***REMOVED***
						return nil, p.getErr(ErrBadClassInCharRange, ch)
					***REMOVED***
					cc.addSpace(p.useOptionE(), p.useRE2(), ch == 'S')
				***REMOVED***
				continue

			case 'W', 'w':
				if !scanOnly ***REMOVED***
					if inRange ***REMOVED***
						return nil, p.getErr(ErrBadClassInCharRange, ch)
					***REMOVED***

					cc.addWord(p.useOptionE() || p.useRE2(), ch == 'W')
				***REMOVED***
				continue

			case 'p', 'P':
				if !scanOnly ***REMOVED***
					if inRange ***REMOVED***
						return nil, p.getErr(ErrBadClassInCharRange, ch)
					***REMOVED***
					prop, err := p.parseProperty()
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					cc.addCategory(prop, (ch != 'p'), caseInsensitive, p.patternRaw)
				***REMOVED*** else ***REMOVED***
					p.parseProperty()
				***REMOVED***

				continue

			case '-':
				if !scanOnly ***REMOVED***
					cc.addRange(ch, ch)
				***REMOVED***
				continue

			default:
				p.moveLeft()
				var err error
				ch, err = p.scanCharEscape() // non-literal character
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				fTranslatedChar = true
				break // this break will only break out of the switch
			***REMOVED***
		***REMOVED*** else if ch == '[' ***REMOVED***
			// This is code for Posix style properties - [:Ll:] or [:IsTibetan:].
			// It currently doesn't do anything other than skip the whole thing!
			if p.charsRight() > 0 && p.rightChar(0) == ':' && !inRange ***REMOVED***
				savePos := p.textpos()

				p.moveRight(1)
				negate := false
				if p.charsRight() > 1 && p.rightChar(0) == '^' ***REMOVED***
					negate = true
					p.moveRight(1)
				***REMOVED***

				nm := p.scanCapname() // snag the name
				if !scanOnly && p.useRE2() ***REMOVED***
					// look up the name since these are valid for RE2
					// add the group based on the name
					if ok := cc.addNamedASCII(nm, negate); !ok ***REMOVED***
						return nil, p.getErr(ErrInvalidCharRange)
					***REMOVED***
				***REMOVED***
				if p.charsRight() < 2 || p.moveRightGetChar() != ':' || p.moveRightGetChar() != ']' ***REMOVED***
					p.textto(savePos)
				***REMOVED*** else if p.useRE2() ***REMOVED***
					// move on
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if inRange ***REMOVED***
			inRange = false
			if !scanOnly ***REMOVED***
				if ch == '[' && !fTranslatedChar && !firstChar ***REMOVED***
					// We thought we were in a range, but we're actually starting a subtraction.
					// In that case, we'll add chPrev to our char class, skip the opening [, and
					// scan the new character class recursively.
					cc.addChar(chPrev)
					sub, err := p.scanCharSet(caseInsensitive, false)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
					cc.addSubtraction(sub)

					if p.charsRight() > 0 && p.rightChar(0) != ']' ***REMOVED***
						return nil, p.getErr(ErrSubtractionMustBeLast)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					// a regular range, like a-z
					if chPrev > ch ***REMOVED***
						return nil, p.getErr(ErrReversedCharRange, chPrev, ch)
					***REMOVED***
					cc.addRange(chPrev, ch)
				***REMOVED***
			***REMOVED***
		***REMOVED*** else if p.charsRight() >= 2 && p.rightChar(0) == '-' && p.rightChar(1) != ']' ***REMOVED***
			// this could be the start of a range
			chPrev = ch
			inRange = true
			p.moveRight(1)
		***REMOVED*** else if p.charsRight() >= 1 && ch == '-' && !fTranslatedChar && p.rightChar(0) == '[' && !firstChar ***REMOVED***
			// we aren't in a range, and now there is a subtraction.  Usually this happens
			// only when a subtraction follows a range, like [a-z-[b]]
			if !scanOnly ***REMOVED***
				p.moveRight(1)
				sub, err := p.scanCharSet(caseInsensitive, false)
				if err != nil ***REMOVED***
					return nil, err
				***REMOVED***
				cc.addSubtraction(sub)

				if p.charsRight() > 0 && p.rightChar(0) != ']' ***REMOVED***
					return nil, p.getErr(ErrSubtractionMustBeLast)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				p.moveRight(1)
				p.scanCharSet(caseInsensitive, true)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if !scanOnly ***REMOVED***
				cc.addRange(ch, ch)
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if !closed ***REMOVED***
		return nil, p.getErr(ErrUnterminatedBracket)
	***REMOVED***

	if !scanOnly && caseInsensitive ***REMOVED***
		cc.addLowercase()
	***REMOVED***

	return cc, nil
***REMOVED***

// Scans any number of decimal digits (pegs value at 2^31-1 if too large)
func (p *parser) scanDecimal() (int, error) ***REMOVED***
	i := 0
	var d int

	for p.charsRight() > 0 ***REMOVED***
		d = int(p.rightChar(0) - '0')
		if d < 0 || d > 9 ***REMOVED***
			break
		***REMOVED***
		p.moveRight(1)

		if i > maxValueDiv10 || (i == maxValueDiv10 && d > maxValueMod10) ***REMOVED***
			return 0, p.getErr(ErrCaptureGroupOutOfRange)
		***REMOVED***

		i *= 10
		i += d
	***REMOVED***

	return int(i), nil
***REMOVED***

// Returns true for options allowed only at the top level
func isOnlyTopOption(option RegexOptions) bool ***REMOVED***
	return option == RightToLeft || option == ECMAScript || option == RE2
***REMOVED***

// Scans cimsx-cimsx option string, stops at the first unrecognized char.
func (p *parser) scanOptions() ***REMOVED***

	for off := false; p.charsRight() > 0; p.moveRight(1) ***REMOVED***
		ch := p.rightChar(0)

		if ch == '-' ***REMOVED***
			off = true
		***REMOVED*** else if ch == '+' ***REMOVED***
			off = false
		***REMOVED*** else ***REMOVED***
			option := optionFromCode(ch)
			if option == 0 || isOnlyTopOption(option) ***REMOVED***
				return
			***REMOVED***

			if off ***REMOVED***
				p.options &= ^option
			***REMOVED*** else ***REMOVED***
				p.options |= option
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Scans \ code for escape codes that map to single unicode chars.
func (p *parser) scanCharEscape() (r rune, err error) ***REMOVED***

	ch := p.moveRightGetChar()

	if ch >= '0' && ch <= '7' ***REMOVED***
		p.moveLeft()
		return p.scanOctal(), nil
	***REMOVED***

	pos := p.textpos()

	switch ch ***REMOVED***
	case 'x':
		// support for \x***REMOVED***HEX***REMOVED*** syntax from Perl and PCRE
		if p.charsRight() > 0 && p.rightChar(0) == '***REMOVED***' ***REMOVED***
			if p.useOptionE() ***REMOVED***
				return ch, nil
			***REMOVED***
			p.moveRight(1)
			return p.scanHexUntilBrace()
		***REMOVED*** else ***REMOVED***
			r, err = p.scanHex(2)
		***REMOVED***
	case 'u':
		// ECMAscript suppot \u***REMOVED***HEX***REMOVED*** only if `u` is also set
		if p.useOptionE() && p.useOptionU() && p.charsRight() > 0 && p.rightChar(0) == '***REMOVED***' ***REMOVED***
			p.moveRight(1)
			return p.scanHexUntilBrace()
		***REMOVED*** else ***REMOVED***
			r, err = p.scanHex(4)
		***REMOVED***
	case 'a':
		return '\u0007', nil
	case 'b':
		return '\b', nil
	case 'e':
		return '\u001B', nil
	case 'f':
		return '\f', nil
	case 'n':
		return '\n', nil
	case 'r':
		return '\r', nil
	case 't':
		return '\t', nil
	case 'v':
		return '\u000B', nil
	case 'c':
		r, err = p.scanControl()
	default:
		if !p.useOptionE() && !p.useRE2() && IsWordChar(ch) ***REMOVED***
			return 0, p.getErr(ErrUnrecognizedEscape, string(ch))
		***REMOVED***
		return ch, nil
	***REMOVED***
	if err != nil && p.useOptionE() ***REMOVED***
		p.textto(pos)
		return ch, nil
	***REMOVED***
	return
***REMOVED***

// Grabs and converts an ascii control character
func (p *parser) scanControl() (rune, error) ***REMOVED***
	if p.charsRight() <= 0 ***REMOVED***
		return 0, p.getErr(ErrMissingControl)
	***REMOVED***

	ch := p.moveRightGetChar()

	// \ca interpreted as \cA

	if ch >= 'a' && ch <= 'z' ***REMOVED***
		ch = (ch - ('a' - 'A'))
	***REMOVED***
	ch = (ch - '@')
	if ch >= 0 && ch < ' ' ***REMOVED***
		return ch, nil
	***REMOVED***

	return 0, p.getErr(ErrUnrecognizedControl)

***REMOVED***

// Scan hex digits until we hit a closing brace.
// Non-hex digits, hex value too large for UTF-8, or running out of chars are errors
func (p *parser) scanHexUntilBrace() (rune, error) ***REMOVED***
	// PCRE spec reads like unlimited hex digits are allowed, but unicode has a limit
	// so we can enforce that
	i := 0
	hasContent := false

	for p.charsRight() > 0 ***REMOVED***
		ch := p.moveRightGetChar()
		if ch == '***REMOVED***' ***REMOVED***
			// hit our close brace, we're done here
			// prevent \x***REMOVED******REMOVED***
			if !hasContent ***REMOVED***
				return 0, p.getErr(ErrTooFewHex)
			***REMOVED***
			return rune(i), nil
		***REMOVED***
		hasContent = true
		// no brace needs to be hex digit
		d := hexDigit(ch)
		if d < 0 ***REMOVED***
			return 0, p.getErr(ErrMissingBrace)
		***REMOVED***

		i *= 0x10
		i += d

		if i > unicode.MaxRune ***REMOVED***
			return 0, p.getErr(ErrInvalidHex)
		***REMOVED***
	***REMOVED***

	// we only make it here if we run out of digits without finding the brace
	return 0, p.getErr(ErrMissingBrace)
***REMOVED***

// Scans exactly c hex digits (c=2 for \xFF, c=4 for \uFFFF)
func (p *parser) scanHex(c int) (rune, error) ***REMOVED***

	i := 0

	if p.charsRight() >= c ***REMOVED***
		for c > 0 ***REMOVED***
			d := hexDigit(p.moveRightGetChar())
			if d < 0 ***REMOVED***
				break
			***REMOVED***
			i *= 0x10
			i += d
			c--
		***REMOVED***
	***REMOVED***

	if c > 0 ***REMOVED***
		return 0, p.getErr(ErrTooFewHex)
	***REMOVED***

	return rune(i), nil
***REMOVED***

// Returns n <= 0xF for a hex digit.
func hexDigit(ch rune) int ***REMOVED***

	if d := uint(ch - '0'); d <= 9 ***REMOVED***
		return int(d)
	***REMOVED***

	if d := uint(ch - 'a'); d <= 5 ***REMOVED***
		return int(d + 0xa)
	***REMOVED***

	if d := uint(ch - 'A'); d <= 5 ***REMOVED***
		return int(d + 0xa)
	***REMOVED***

	return -1
***REMOVED***

// Scans up to three octal digits (stops before exceeding 0377).
func (p *parser) scanOctal() rune ***REMOVED***
	// Consume octal chars only up to 3 digits and value 0377

	c := 3

	if c > p.charsRight() ***REMOVED***
		c = p.charsRight()
	***REMOVED***

	//we know the first char is good because the caller had to check
	i := 0
	d := int(p.rightChar(0) - '0')
	for c > 0 && d <= 7 && d >= 0 ***REMOVED***
		if i >= 0x20 && p.useOptionE() ***REMOVED***
			break
		***REMOVED***
		i *= 8
		i += d
		c--

		p.moveRight(1)
		if !p.rightMost() ***REMOVED***
			d = int(p.rightChar(0) - '0')
		***REMOVED***
	***REMOVED***

	// Octal codes only go up to 255.  Any larger and the behavior that Perl follows
	// is simply to truncate the high bits.
	i &= 0xFF

	return rune(i)
***REMOVED***

// Returns the current parsing position.
func (p *parser) textpos() int ***REMOVED***
	return p.currentPos
***REMOVED***

// Zaps to a specific parsing position.
func (p *parser) textto(pos int) ***REMOVED***
	p.currentPos = pos
***REMOVED***

// Returns the char at the right of the current parsing position and advances to the right.
func (p *parser) moveRightGetChar() rune ***REMOVED***
	ch := p.pattern[p.currentPos]
	p.currentPos++
	return ch
***REMOVED***

// Moves the current position to the right.
func (p *parser) moveRight(i int) ***REMOVED***
	// default would be 1
	p.currentPos += i
***REMOVED***

// Moves the current parsing position one to the left.
func (p *parser) moveLeft() ***REMOVED***
	p.currentPos--
***REMOVED***

// Returns the char left of the current parsing position.
func (p *parser) charAt(i int) rune ***REMOVED***
	return p.pattern[i]
***REMOVED***

// Returns the char i chars right of the current parsing position.
func (p *parser) rightChar(i int) rune ***REMOVED***
	// default would be 0
	return p.pattern[p.currentPos+i]
***REMOVED***

// Number of characters to the right of the current parsing position.
func (p *parser) charsRight() int ***REMOVED***
	return len(p.pattern) - p.currentPos
***REMOVED***

func (p *parser) rightMost() bool ***REMOVED***
	return p.currentPos == len(p.pattern)
***REMOVED***

// Looks up the slot number for a given name
func (p *parser) captureSlotFromName(capname string) int ***REMOVED***
	return p.capnames[capname]
***REMOVED***

// True if the capture slot was noted
func (p *parser) isCaptureSlot(i int) bool ***REMOVED***
	if p.caps != nil ***REMOVED***
		_, ok := p.caps[i]
		return ok
	***REMOVED***

	return (i >= 0 && i < p.capsize)
***REMOVED***

// Looks up the slot number for a given name
func (p *parser) isCaptureName(capname string) bool ***REMOVED***
	if p.capnames == nil ***REMOVED***
		return false
	***REMOVED***

	_, ok := p.capnames[capname]
	return ok
***REMOVED***

// option shortcuts

// True if N option disabling '(' autocapture is on.
func (p *parser) useOptionN() bool ***REMOVED***
	return (p.options & ExplicitCapture) != 0
***REMOVED***

// True if I option enabling case-insensitivity is on.
func (p *parser) useOptionI() bool ***REMOVED***
	return (p.options & IgnoreCase) != 0
***REMOVED***

// True if M option altering meaning of $ and ^ is on.
func (p *parser) useOptionM() bool ***REMOVED***
	return (p.options & Multiline) != 0
***REMOVED***

// True if S option altering meaning of . is on.
func (p *parser) useOptionS() bool ***REMOVED***
	return (p.options & Singleline) != 0
***REMOVED***

// True if X option enabling whitespace/comment mode is on.
func (p *parser) useOptionX() bool ***REMOVED***
	return (p.options & IgnorePatternWhitespace) != 0
***REMOVED***

// True if E option enabling ECMAScript behavior on.
func (p *parser) useOptionE() bool ***REMOVED***
	return (p.options & ECMAScript) != 0
***REMOVED***

// true to use RE2 compatibility parsing behavior.
func (p *parser) useRE2() bool ***REMOVED***
	return (p.options & RE2) != 0
***REMOVED***

// True if U option enabling ECMAScript's Unicode behavior on.
func (p *parser) useOptionU() bool ***REMOVED***
	return (p.options & Unicode) != 0
***REMOVED***

// True if options stack is empty.
func (p *parser) emptyOptionsStack() bool ***REMOVED***
	return len(p.optionsStack) == 0
***REMOVED***

// Finish the current quantifiable (when a quantifier is not found or is not possible)
func (p *parser) addConcatenate() ***REMOVED***
	// The first (| inside a Testgroup group goes directly to the group
	p.concatenation.addChild(p.unit)
	p.unit = nil
***REMOVED***

// Finish the current quantifiable (when a quantifier is found)
func (p *parser) addConcatenate3(lazy bool, min, max int) ***REMOVED***
	p.concatenation.addChild(p.unit.makeQuantifier(lazy, min, max))
	p.unit = nil
***REMOVED***

// Sets the current unit to a single char node
func (p *parser) addUnitOne(ch rune) ***REMOVED***
	if p.useOptionI() ***REMOVED***
		ch = unicode.ToLower(ch)
	***REMOVED***

	p.unit = newRegexNodeCh(ntOne, p.options, ch)
***REMOVED***

// Sets the current unit to a single inverse-char node
func (p *parser) addUnitNotone(ch rune) ***REMOVED***
	if p.useOptionI() ***REMOVED***
		ch = unicode.ToLower(ch)
	***REMOVED***

	p.unit = newRegexNodeCh(ntNotone, p.options, ch)
***REMOVED***

// Sets the current unit to a single set node
func (p *parser) addUnitSet(set *CharSet) ***REMOVED***
	p.unit = newRegexNodeSet(ntSet, p.options, set)
***REMOVED***

// Sets the current unit to a subtree
func (p *parser) addUnitNode(node *regexNode) ***REMOVED***
	p.unit = node
***REMOVED***

// Sets the current unit to an assertion of the specified type
func (p *parser) addUnitType(t nodeType) ***REMOVED***
	p.unit = newRegexNode(t, p.options)
***REMOVED***

// Finish the current group (in response to a ')' or end)
func (p *parser) addGroup() error ***REMOVED***
	if p.group.t == ntTestgroup || p.group.t == ntTestref ***REMOVED***
		p.group.addChild(p.concatenation.reverseLeft())
		if (p.group.t == ntTestref && len(p.group.children) > 2) || len(p.group.children) > 3 ***REMOVED***
			return p.getErr(ErrTooManyAlternates)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		p.alternation.addChild(p.concatenation.reverseLeft())
		p.group.addChild(p.alternation)
	***REMOVED***

	p.unit = p.group
	return nil
***REMOVED***

// Pops the option stack, but keeps the current options unchanged.
func (p *parser) popKeepOptions() ***REMOVED***
	lastIdx := len(p.optionsStack) - 1
	p.optionsStack = p.optionsStack[:lastIdx]
***REMOVED***

// Recalls options from the stack.
func (p *parser) popOptions() ***REMOVED***
	lastIdx := len(p.optionsStack) - 1
	// get the last item on the stack and then remove it by reslicing
	p.options = p.optionsStack[lastIdx]
	p.optionsStack = p.optionsStack[:lastIdx]
***REMOVED***

// Saves options on a stack.
func (p *parser) pushOptions() ***REMOVED***
	p.optionsStack = append(p.optionsStack, p.options)
***REMOVED***

// Add a string to the last concatenate.
func (p *parser) addToConcatenate(pos, cch int, isReplacement bool) ***REMOVED***
	var node *regexNode

	if cch == 0 ***REMOVED***
		return
	***REMOVED***

	if cch > 1 ***REMOVED***
		str := make([]rune, cch)
		copy(str, p.pattern[pos:pos+cch])

		if p.useOptionI() && !isReplacement ***REMOVED***
			// We do the ToLower character by character for consistency.  With surrogate chars, doing
			// a ToLower on the entire string could actually change the surrogate pair.  This is more correct
			// linguistically, but since Regex doesn't support surrogates, it's more important to be
			// consistent.
			for i := 0; i < len(str); i++ ***REMOVED***
				str[i] = unicode.ToLower(str[i])
			***REMOVED***
		***REMOVED***

		node = newRegexNodeStr(ntMulti, p.options, str)
	***REMOVED*** else ***REMOVED***
		ch := p.charAt(pos)

		if p.useOptionI() && !isReplacement ***REMOVED***
			ch = unicode.ToLower(ch)
		***REMOVED***

		node = newRegexNodeCh(ntOne, p.options, ch)
	***REMOVED***

	p.concatenation.addChild(node)
***REMOVED***

// Push the parser state (in response to an open paren)
func (p *parser) pushGroup() ***REMOVED***
	p.group.next = p.stack
	p.alternation.next = p.group
	p.concatenation.next = p.alternation
	p.stack = p.concatenation
***REMOVED***

// Remember the pushed state (in response to a ')')
func (p *parser) popGroup() error ***REMOVED***
	p.concatenation = p.stack
	p.alternation = p.concatenation.next
	p.group = p.alternation.next
	p.stack = p.group.next

	// The first () inside a Testgroup group goes directly to the group
	if p.group.t == ntTestgroup && len(p.group.children) == 0 ***REMOVED***
		if p.unit == nil ***REMOVED***
			return p.getErr(ErrConditionalExpression)
		***REMOVED***

		p.group.addChild(p.unit)
		p.unit = nil
	***REMOVED***
	return nil
***REMOVED***

// True if the group stack is empty.
func (p *parser) emptyStack() bool ***REMOVED***
	return p.stack == nil
***REMOVED***

// Start a new round for the parser state (in response to an open paren or string start)
func (p *parser) startGroup(openGroup *regexNode) ***REMOVED***
	p.group = openGroup
	p.alternation = newRegexNode(ntAlternate, p.options)
	p.concatenation = newRegexNode(ntConcatenate, p.options)
***REMOVED***

// Finish the current concatenation (in response to a |)
func (p *parser) addAlternate() ***REMOVED***
	// The | parts inside a Testgroup group go directly to the group

	if p.group.t == ntTestgroup || p.group.t == ntTestref ***REMOVED***
		p.group.addChild(p.concatenation.reverseLeft())
	***REMOVED*** else ***REMOVED***
		p.alternation.addChild(p.concatenation.reverseLeft())
	***REMOVED***

	p.concatenation = newRegexNode(ntConcatenate, p.options)
***REMOVED***

// For categorizing ascii characters.

const (
	Q byte = 5 // quantifier
	S      = 4 // ordinary stopper
	Z      = 3 // ScanBlank stopper
	X      = 2 // whitespace
	E      = 1 // should be escaped
)

var _category = []byte***REMOVED***
	//01  2  3  4  5  6  7  8  9  A  B  C  D  E  F  0  1  2  3  4  5  6  7  8  9  A  B  C  D  E  F
	0, 0, 0, 0, 0, 0, 0, 0, 0, X, X, X, X, X, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	// !  "  #  $  %  &  '  (  )  *  +  ,  -  .  /  0  1  2  3  4  5  6  7  8  9  :  ;  <  =  >  ?
	X, 0, 0, Z, S, 0, 0, 0, S, S, Q, Q, 0, 0, S, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, Q,
	//@A  B  C  D  E  F  G  H  I  J  K  L  M  N  O  P  Q  R  S  T  U  V  W  X  Y  Z  [  \  ]  ^  _
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, S, S, 0, S, 0,
	//'a  b  c  d  e  f  g  h  i  j  k  l  m  n  o  p  q  r  s  t  u  v  w  x  y  z  ***REMOVED***  |  ***REMOVED***  ~
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, Q, S, 0, 0, 0,
***REMOVED***

func isSpace(ch rune) bool ***REMOVED***
	return (ch <= ' ' && _category[ch] == X)
***REMOVED***

// Returns true for those characters that terminate a string of ordinary chars.
func isSpecial(ch rune) bool ***REMOVED***
	return (ch <= '|' && _category[ch] >= S)
***REMOVED***

// Returns true for those characters that terminate a string of ordinary chars.
func isStopperX(ch rune) bool ***REMOVED***
	return (ch <= '|' && _category[ch] >= X)
***REMOVED***

// Returns true for those characters that begin a quantifier.
func isQuantifier(ch rune) bool ***REMOVED***
	return (ch <= '***REMOVED***' && _category[ch] >= Q)
***REMOVED***

func (p *parser) isTrueQuantifier() bool ***REMOVED***
	nChars := p.charsRight()
	if nChars == 0 ***REMOVED***
		return false
	***REMOVED***

	startpos := p.textpos()
	ch := p.charAt(startpos)
	if ch != '***REMOVED***' ***REMOVED***
		return ch <= '***REMOVED***' && _category[ch] >= Q
	***REMOVED***

	//UGLY: this is ugly -- the original code was ugly too
	pos := startpos
	for ***REMOVED***
		nChars--
		if nChars <= 0 ***REMOVED***
			break
		***REMOVED***
		pos++
		ch = p.charAt(pos)
		if ch < '0' || ch > '9' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if nChars == 0 || pos-startpos == 1 ***REMOVED***
		return false
	***REMOVED***
	if ch == '***REMOVED***' ***REMOVED***
		return true
	***REMOVED***
	if ch != ',' ***REMOVED***
		return false
	***REMOVED***
	for ***REMOVED***
		nChars--
		if nChars <= 0 ***REMOVED***
			break
		***REMOVED***
		pos++
		ch = p.charAt(pos)
		if ch < '0' || ch > '9' ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return nChars > 0 && ch == '***REMOVED***'
***REMOVED***
