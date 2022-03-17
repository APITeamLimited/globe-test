// Package cascadia is an implementation of CSS selectors.
package cascadia

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// a parser for CSS selectors
type parser struct ***REMOVED***
	s string // the source text
	i int    // the current position

	// if `false`, parsing a pseudo-element
	// returns an error.
	acceptPseudoElements bool
***REMOVED***

// parseEscape parses a backslash escape.
func (p *parser) parseEscape() (result string, err error) ***REMOVED***
	if len(p.s) < p.i+2 || p.s[p.i] != '\\' ***REMOVED***
		return "", errors.New("invalid escape sequence")
	***REMOVED***

	start := p.i + 1
	c := p.s[start]
	switch ***REMOVED***
	case c == '\r' || c == '\n' || c == '\f':
		return "", errors.New("escaped line ending outside string")
	case hexDigit(c):
		// unicode escape (hex)
		var i int
		for i = start; i < start+6 && i < len(p.s) && hexDigit(p.s[i]); i++ ***REMOVED***
			// empty
		***REMOVED***
		v, _ := strconv.ParseUint(p.s[start:i], 16, 64)
		if len(p.s) > i ***REMOVED***
			switch p.s[i] ***REMOVED***
			case '\r':
				i++
				if len(p.s) > i && p.s[i] == '\n' ***REMOVED***
					i++
				***REMOVED***
			case ' ', '\t', '\n', '\f':
				i++
			***REMOVED***
		***REMOVED***
		p.i = i
		return string(rune(v)), nil
	***REMOVED***

	// Return the literal character after the backslash.
	result = p.s[start : start+1]
	p.i += 2
	return result, nil
***REMOVED***

// toLowerASCII returns s with all ASCII capital letters lowercased.
func toLowerASCII(s string) string ***REMOVED***
	var b []byte
	for i := 0; i < len(s); i++ ***REMOVED***
		if c := s[i]; 'A' <= c && c <= 'Z' ***REMOVED***
			if b == nil ***REMOVED***
				b = make([]byte, len(s))
				copy(b, s)
			***REMOVED***
			b[i] = s[i] + ('a' - 'A')
		***REMOVED***
	***REMOVED***

	if b == nil ***REMOVED***
		return s
	***REMOVED***

	return string(b)
***REMOVED***

func hexDigit(c byte) bool ***REMOVED***
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
***REMOVED***

// nameStart returns whether c can be the first character of an identifier
// (not counting an initial hyphen, or an escape sequence).
func nameStart(c byte) bool ***REMOVED***
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c > 127
***REMOVED***

// nameChar returns whether c can be a character within an identifier
// (not counting an escape sequence).
func nameChar(c byte) bool ***REMOVED***
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c > 127 ||
		c == '-' || '0' <= c && c <= '9'
***REMOVED***

// parseIdentifier parses an identifier.
func (p *parser) parseIdentifier() (result string, err error) ***REMOVED***
	startingDash := false
	if len(p.s) > p.i && p.s[p.i] == '-' ***REMOVED***
		startingDash = true
		p.i++
	***REMOVED***

	if len(p.s) <= p.i ***REMOVED***
		return "", errors.New("expected identifier, found EOF instead")
	***REMOVED***

	if c := p.s[p.i]; !(nameStart(c) || c == '\\') ***REMOVED***
		return "", fmt.Errorf("expected identifier, found %c instead", c)
	***REMOVED***

	result, err = p.parseName()
	if startingDash && err == nil ***REMOVED***
		result = "-" + result
	***REMOVED***
	return
***REMOVED***

// parseName parses a name (which is like an identifier, but doesn't have
// extra restrictions on the first character).
func (p *parser) parseName() (result string, err error) ***REMOVED***
	i := p.i
loop:
	for i < len(p.s) ***REMOVED***
		c := p.s[i]
		switch ***REMOVED***
		case nameChar(c):
			start := i
			for i < len(p.s) && nameChar(p.s[i]) ***REMOVED***
				i++
			***REMOVED***
			result += p.s[start:i]
		case c == '\\':
			p.i = i
			val, err := p.parseEscape()
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			i = p.i
			result += val
		default:
			break loop
		***REMOVED***
	***REMOVED***

	if result == "" ***REMOVED***
		return "", errors.New("expected name, found EOF instead")
	***REMOVED***

	p.i = i
	return result, nil
***REMOVED***

// parseString parses a single- or double-quoted string.
func (p *parser) parseString() (result string, err error) ***REMOVED***
	i := p.i
	if len(p.s) < i+2 ***REMOVED***
		return "", errors.New("expected string, found EOF instead")
	***REMOVED***

	quote := p.s[i]
	i++

loop:
	for i < len(p.s) ***REMOVED***
		switch p.s[i] ***REMOVED***
		case '\\':
			if len(p.s) > i+1 ***REMOVED***
				switch c := p.s[i+1]; c ***REMOVED***
				case '\r':
					if len(p.s) > i+2 && p.s[i+2] == '\n' ***REMOVED***
						i += 3
						continue loop
					***REMOVED***
					fallthrough
				case '\n', '\f':
					i += 2
					continue loop
				***REMOVED***
			***REMOVED***
			p.i = i
			val, err := p.parseEscape()
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			i = p.i
			result += val
		case quote:
			break loop
		case '\r', '\n', '\f':
			return "", errors.New("unexpected end of line in string")
		default:
			start := i
			for i < len(p.s) ***REMOVED***
				if c := p.s[i]; c == quote || c == '\\' || c == '\r' || c == '\n' || c == '\f' ***REMOVED***
					break
				***REMOVED***
				i++
			***REMOVED***
			result += p.s[start:i]
		***REMOVED***
	***REMOVED***

	if i >= len(p.s) ***REMOVED***
		return "", errors.New("EOF in string")
	***REMOVED***

	// Consume the final quote.
	i++

	p.i = i
	return result, nil
***REMOVED***

// parseRegex parses a regular expression; the end is defined by encountering an
// unmatched closing ')' or ']' which is not consumed
func (p *parser) parseRegex() (rx *regexp.Regexp, err error) ***REMOVED***
	i := p.i
	if len(p.s) < i+2 ***REMOVED***
		return nil, errors.New("expected regular expression, found EOF instead")
	***REMOVED***

	// number of open parens or brackets;
	// when it becomes negative, finished parsing regex
	open := 0

loop:
	for i < len(p.s) ***REMOVED***
		switch p.s[i] ***REMOVED***
		case '(', '[':
			open++
		case ')', ']':
			open--
			if open < 0 ***REMOVED***
				break loop
			***REMOVED***
		***REMOVED***
		i++
	***REMOVED***

	if i >= len(p.s) ***REMOVED***
		return nil, errors.New("EOF in regular expression")
	***REMOVED***
	rx, err = regexp.Compile(p.s[p.i:i])
	p.i = i
	return rx, err
***REMOVED***

// skipWhitespace consumes whitespace characters and comments.
// It returns true if there was actually anything to skip.
func (p *parser) skipWhitespace() bool ***REMOVED***
	i := p.i
	for i < len(p.s) ***REMOVED***
		switch p.s[i] ***REMOVED***
		case ' ', '\t', '\r', '\n', '\f':
			i++
			continue
		case '/':
			if strings.HasPrefix(p.s[i:], "/*") ***REMOVED***
				end := strings.Index(p.s[i+len("/*"):], "*/")
				if end != -1 ***REMOVED***
					i += end + len("/**/")
					continue
				***REMOVED***
			***REMOVED***
		***REMOVED***
		break
	***REMOVED***

	if i > p.i ***REMOVED***
		p.i = i
		return true
	***REMOVED***

	return false
***REMOVED***

// consumeParenthesis consumes an opening parenthesis and any following
// whitespace. It returns true if there was actually a parenthesis to skip.
func (p *parser) consumeParenthesis() bool ***REMOVED***
	if p.i < len(p.s) && p.s[p.i] == '(' ***REMOVED***
		p.i++
		p.skipWhitespace()
		return true
	***REMOVED***
	return false
***REMOVED***

// consumeClosingParenthesis consumes a closing parenthesis and any preceding
// whitespace. It returns true if there was actually a parenthesis to skip.
func (p *parser) consumeClosingParenthesis() bool ***REMOVED***
	i := p.i
	p.skipWhitespace()
	if p.i < len(p.s) && p.s[p.i] == ')' ***REMOVED***
		p.i++
		return true
	***REMOVED***
	p.i = i
	return false
***REMOVED***

// parseTypeSelector parses a type selector (one that matches by tag name).
func (p *parser) parseTypeSelector() (result tagSelector, err error) ***REMOVED***
	tag, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	return tagSelector***REMOVED***tag: toLowerASCII(tag)***REMOVED***, nil
***REMOVED***

// parseIDSelector parses a selector that matches by id attribute.
func (p *parser) parseIDSelector() (idSelector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return idSelector***REMOVED******REMOVED***, fmt.Errorf("expected id selector (#id), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '#' ***REMOVED***
		return idSelector***REMOVED******REMOVED***, fmt.Errorf("expected id selector (#id), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	id, err := p.parseName()
	if err != nil ***REMOVED***
		return idSelector***REMOVED******REMOVED***, err
	***REMOVED***

	return idSelector***REMOVED***id: id***REMOVED***, nil
***REMOVED***

// parseClassSelector parses a selector that matches by class attribute.
func (p *parser) parseClassSelector() (classSelector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return classSelector***REMOVED******REMOVED***, fmt.Errorf("expected class selector (.class), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '.' ***REMOVED***
		return classSelector***REMOVED******REMOVED***, fmt.Errorf("expected class selector (.class), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	class, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return classSelector***REMOVED******REMOVED***, err
	***REMOVED***

	return classSelector***REMOVED***class: class***REMOVED***, nil
***REMOVED***

// parseAttributeSelector parses a selector that matches by attribute value.
func (p *parser) parseAttributeSelector() (attrSelector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, fmt.Errorf("expected attribute selector ([attribute]), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '[' ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, fmt.Errorf("expected attribute selector ([attribute]), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	p.skipWhitespace()
	key, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, err
	***REMOVED***
	key = toLowerASCII(key)

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	if p.s[p.i] == ']' ***REMOVED***
		p.i++
		return attrSelector***REMOVED***key: key, operation: ""***REMOVED***, nil
	***REMOVED***

	if p.i+2 >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	op := p.s[p.i : p.i+2]
	if op[0] == '=' ***REMOVED***
		op = "="
	***REMOVED*** else if op[1] != '=' ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, fmt.Errorf(`expected equality operator, found "%s" instead`, op)
	***REMOVED***
	p.i += len(op)

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, errors.New("unexpected EOF in attribute selector")
	***REMOVED***
	var val string
	var rx *regexp.Regexp
	if op == "#=" ***REMOVED***
		rx, err = p.parseRegex()
	***REMOVED*** else ***REMOVED***
		switch p.s[p.i] ***REMOVED***
		case '\'', '"':
			val, err = p.parseString()
		default:
			val, err = p.parseIdentifier()
		***REMOVED***
	***REMOVED***
	if err != nil ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, err
	***REMOVED***

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	// check if the attribute contains an ignore case flag
	ignoreCase := false
	if p.s[p.i] == 'i' || p.s[p.i] == 'I' ***REMOVED***
		ignoreCase = true
		p.i++
	***REMOVED***

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	if p.s[p.i] != ']' ***REMOVED***
		return attrSelector***REMOVED******REMOVED***, fmt.Errorf("expected ']', found '%c' instead", p.s[p.i])
	***REMOVED***
	p.i++

	switch op ***REMOVED***
	case "=", "!=", "~=", "|=", "^=", "$=", "*=", "#=":
		return attrSelector***REMOVED***key: key, val: val, operation: op, regexp: rx, insensitive: ignoreCase***REMOVED***, nil
	default:
		return attrSelector***REMOVED******REMOVED***, fmt.Errorf("attribute operator %q is not supported", op)
	***REMOVED***
***REMOVED***

var (
	errExpectedParenthesis        = errors.New("expected '(' but didn't find it")
	errExpectedClosingParenthesis = errors.New("expected ')' but didn't find it")
	errUnmatchedParenthesis       = errors.New("unmatched '('")
)

// parsePseudoclassSelector parses a pseudoclass selector like :not(p) or a pseudo-element
// For backwards compatibility, both ':' and '::' prefix are allowed for pseudo-elements.
// https://drafts.csswg.org/selectors-3/#pseudo-elements
// Returning a nil `Sel` (and a nil `error`) means we found a pseudo-element.
func (p *parser) parsePseudoclassSelector() (out Sel, pseudoElement string, err error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return nil, "", fmt.Errorf("expected pseudoclass selector (:pseudoclass), found EOF instead")
	***REMOVED***
	if p.s[p.i] != ':' ***REMOVED***
		return nil, "", fmt.Errorf("expected attribute selector (:pseudoclass), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	var mustBePseudoElement bool
	if p.i >= len(p.s) ***REMOVED***
		return nil, "", fmt.Errorf("got empty pseudoclass (or pseudoelement)")
	***REMOVED***
	if p.s[p.i] == ':' ***REMOVED*** // we found a pseudo-element
		mustBePseudoElement = true
		p.i++
	***REMOVED***

	name, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return
	***REMOVED***
	name = toLowerASCII(name)
	if mustBePseudoElement && (name != "after" && name != "backdrop" && name != "before" &&
		name != "cue" && name != "first-letter" && name != "first-line" && name != "grammar-error" &&
		name != "marker" && name != "placeholder" && name != "selection" && name != "spelling-error") ***REMOVED***
		return out, "", fmt.Errorf("unknown pseudoelement :%s", name)
	***REMOVED***

	switch name ***REMOVED***
	case "not", "has", "haschild":
		if !p.consumeParenthesis() ***REMOVED***
			return out, "", errExpectedParenthesis
		***REMOVED***
		sel, parseErr := p.parseSelectorGroup()
		if parseErr != nil ***REMOVED***
			return out, "", parseErr
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return out, "", errExpectedClosingParenthesis
		***REMOVED***

		out = relativePseudoClassSelector***REMOVED***name: name, match: sel***REMOVED***

	case "contains", "containsown":
		if !p.consumeParenthesis() ***REMOVED***
			return out, "", errExpectedParenthesis
		***REMOVED***
		if p.i == len(p.s) ***REMOVED***
			return out, "", errUnmatchedParenthesis
		***REMOVED***
		var val string
		switch p.s[p.i] ***REMOVED***
		case '\'', '"':
			val, err = p.parseString()
		default:
			val, err = p.parseIdentifier()
		***REMOVED***
		if err != nil ***REMOVED***
			return out, "", err
		***REMOVED***
		val = strings.ToLower(val)
		p.skipWhitespace()
		if p.i >= len(p.s) ***REMOVED***
			return out, "", errors.New("unexpected EOF in pseudo selector")
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return out, "", errExpectedClosingParenthesis
		***REMOVED***

		out = containsPseudoClassSelector***REMOVED***own: name == "containsown", value: val***REMOVED***

	case "matches", "matchesown":
		if !p.consumeParenthesis() ***REMOVED***
			return out, "", errExpectedParenthesis
		***REMOVED***
		rx, err := p.parseRegex()
		if err != nil ***REMOVED***
			return out, "", err
		***REMOVED***
		if p.i >= len(p.s) ***REMOVED***
			return out, "", errors.New("unexpected EOF in pseudo selector")
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return out, "", errExpectedClosingParenthesis
		***REMOVED***

		out = regexpPseudoClassSelector***REMOVED***own: name == "matchesown", regexp: rx***REMOVED***

	case "nth-child", "nth-last-child", "nth-of-type", "nth-last-of-type":
		if !p.consumeParenthesis() ***REMOVED***
			return out, "", errExpectedParenthesis
		***REMOVED***
		a, b, err := p.parseNth()
		if err != nil ***REMOVED***
			return out, "", err
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return out, "", errExpectedClosingParenthesis
		***REMOVED***
		last := name == "nth-last-child" || name == "nth-last-of-type"
		ofType := name == "nth-of-type" || name == "nth-last-of-type"
		out = nthPseudoClassSelector***REMOVED***a: a, b: b, last: last, ofType: ofType***REMOVED***

	case "first-child":
		out = nthPseudoClassSelector***REMOVED***a: 0, b: 1, ofType: false, last: false***REMOVED***
	case "last-child":
		out = nthPseudoClassSelector***REMOVED***a: 0, b: 1, ofType: false, last: true***REMOVED***
	case "first-of-type":
		out = nthPseudoClassSelector***REMOVED***a: 0, b: 1, ofType: true, last: false***REMOVED***
	case "last-of-type":
		out = nthPseudoClassSelector***REMOVED***a: 0, b: 1, ofType: true, last: true***REMOVED***
	case "only-child":
		out = onlyChildPseudoClassSelector***REMOVED***ofType: false***REMOVED***
	case "only-of-type":
		out = onlyChildPseudoClassSelector***REMOVED***ofType: true***REMOVED***
	case "input":
		out = inputPseudoClassSelector***REMOVED******REMOVED***
	case "empty":
		out = emptyElementPseudoClassSelector***REMOVED******REMOVED***
	case "root":
		out = rootPseudoClassSelector***REMOVED******REMOVED***
	case "link":
		out = linkPseudoClassSelector***REMOVED******REMOVED***
	case "lang":
		if !p.consumeParenthesis() ***REMOVED***
			return out, "", errExpectedParenthesis
		***REMOVED***
		if p.i == len(p.s) ***REMOVED***
			return out, "", errUnmatchedParenthesis
		***REMOVED***
		val, err := p.parseIdentifier()
		if err != nil ***REMOVED***
			return out, "", err
		***REMOVED***
		val = strings.ToLower(val)
		p.skipWhitespace()
		if p.i >= len(p.s) ***REMOVED***
			return out, "", errors.New("unexpected EOF in pseudo selector")
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return out, "", errExpectedClosingParenthesis
		***REMOVED***
		out = langPseudoClassSelector***REMOVED***lang: val***REMOVED***
	case "enabled":
		out = enabledPseudoClassSelector***REMOVED******REMOVED***
	case "disabled":
		out = disabledPseudoClassSelector***REMOVED******REMOVED***
	case "checked":
		out = checkedPseudoClassSelector***REMOVED******REMOVED***
	case "visited", "hover", "active", "focus", "target":
		// Not applicable in a static context: never match.
		out = neverMatchSelector***REMOVED***value: ":" + name***REMOVED***
	case "after", "backdrop", "before", "cue", "first-letter", "first-line", "grammar-error", "marker", "placeholder", "selection", "spelling-error":
		return nil, name, nil
	default:
		return out, "", fmt.Errorf("unknown pseudoclass or pseudoelement :%s", name)
	***REMOVED***
	return
***REMOVED***

// parseInteger parses a  decimal integer.
func (p *parser) parseInteger() (int, error) ***REMOVED***
	i := p.i
	start := i
	for i < len(p.s) && '0' <= p.s[i] && p.s[i] <= '9' ***REMOVED***
		i++
	***REMOVED***
	if i == start ***REMOVED***
		return 0, errors.New("expected integer, but didn't find it")
	***REMOVED***
	p.i = i

	val, err := strconv.Atoi(p.s[start:i])
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return val, nil
***REMOVED***

// parseNth parses the argument for :nth-child (normally of the form an+b).
func (p *parser) parseNth() (a, b int, err error) ***REMOVED***
	// initial state
	if p.i >= len(p.s) ***REMOVED***
		goto eof
	***REMOVED***
	switch p.s[p.i] ***REMOVED***
	case '-':
		p.i++
		goto negativeA
	case '+':
		p.i++
		goto positiveA
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto positiveA
	case 'n', 'N':
		a = 1
		p.i++
		goto readN
	case 'o', 'O', 'e', 'E':
		id, nameErr := p.parseName()
		if nameErr != nil ***REMOVED***
			return 0, 0, nameErr
		***REMOVED***
		id = toLowerASCII(id)
		if id == "odd" ***REMOVED***
			return 2, 1, nil
		***REMOVED***
		if id == "even" ***REMOVED***
			return 2, 0, nil
		***REMOVED***
		return 0, 0, fmt.Errorf("expected 'odd' or 'even', but found '%s' instead", id)
	default:
		goto invalid
	***REMOVED***

positiveA:
	if p.i >= len(p.s) ***REMOVED***
		goto eof
	***REMOVED***
	switch p.s[p.i] ***REMOVED***
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		a, err = p.parseInteger()
		if err != nil ***REMOVED***
			return 0, 0, err
		***REMOVED***
		goto readA
	case 'n', 'N':
		a = 1
		p.i++
		goto readN
	default:
		goto invalid
	***REMOVED***

negativeA:
	if p.i >= len(p.s) ***REMOVED***
		goto eof
	***REMOVED***
	switch p.s[p.i] ***REMOVED***
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		a, err = p.parseInteger()
		if err != nil ***REMOVED***
			return 0, 0, err
		***REMOVED***
		a = -a
		goto readA
	case 'n', 'N':
		a = -1
		p.i++
		goto readN
	default:
		goto invalid
	***REMOVED***

readA:
	if p.i >= len(p.s) ***REMOVED***
		goto eof
	***REMOVED***
	switch p.s[p.i] ***REMOVED***
	case 'n', 'N':
		p.i++
		goto readN
	default:
		// The number we read as a is actually b.
		return 0, a, nil
	***REMOVED***

readN:
	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		goto eof
	***REMOVED***
	switch p.s[p.i] ***REMOVED***
	case '+':
		p.i++
		p.skipWhitespace()
		b, err = p.parseInteger()
		if err != nil ***REMOVED***
			return 0, 0, err
		***REMOVED***
		return a, b, nil
	case '-':
		p.i++
		p.skipWhitespace()
		b, err = p.parseInteger()
		if err != nil ***REMOVED***
			return 0, 0, err
		***REMOVED***
		return a, -b, nil
	default:
		return a, 0, nil
	***REMOVED***

eof:
	return 0, 0, errors.New("unexpected EOF while attempting to parse expression of form an+b")

invalid:
	return 0, 0, errors.New("unexpected character while attempting to parse expression of form an+b")
***REMOVED***

// parseSimpleSelectorSequence parses a selector sequence that applies to
// a single element.
func (p *parser) parseSimpleSelectorSequence() (Sel, error) ***REMOVED***
	var selectors []Sel

	if p.i >= len(p.s) ***REMOVED***
		return nil, errors.New("expected selector, found EOF instead")
	***REMOVED***

	switch p.s[p.i] ***REMOVED***
	case '*':
		// It's the universal selector. Just skip over it, since it doesn't affect the meaning.
		p.i++
		if p.i+2 < len(p.s) && p.s[p.i:p.i+2] == "|*" ***REMOVED*** // other version of universal selector
			p.i += 2
		***REMOVED***
	case '#', '.', '[', ':':
		// There's no type selector. Wait to process the other till the main loop.
	default:
		r, err := p.parseTypeSelector()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		selectors = append(selectors, r)
	***REMOVED***

	var pseudoElement string
loop:
	for p.i < len(p.s) ***REMOVED***
		var (
			ns               Sel
			newPseudoElement string
			err              error
		)
		switch p.s[p.i] ***REMOVED***
		case '#':
			ns, err = p.parseIDSelector()
		case '.':
			ns, err = p.parseClassSelector()
		case '[':
			ns, err = p.parseAttributeSelector()
		case ':':
			ns, newPseudoElement, err = p.parsePseudoclassSelector()
		default:
			break loop
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// From https://drafts.csswg.org/selectors-3/#pseudo-elements :
		// "Only one pseudo-element may appear per selector, and if present
		// it must appear after the sequence of simple selectors that
		// represents the subjects of the selector.""
		if ns == nil ***REMOVED*** // we found a pseudo-element
			if pseudoElement != "" ***REMOVED***
				return nil, fmt.Errorf("only one pseudo-element is accepted per selector, got %s and %s", pseudoElement, newPseudoElement)
			***REMOVED***
			if !p.acceptPseudoElements ***REMOVED***
				return nil, fmt.Errorf("pseudo-element %s found, but pseudo-elements support is disabled", newPseudoElement)
			***REMOVED***
			pseudoElement = newPseudoElement
		***REMOVED*** else ***REMOVED***
			if pseudoElement != "" ***REMOVED***
				return nil, fmt.Errorf("pseudo-element %s must be at the end of selector", pseudoElement)
			***REMOVED***
			selectors = append(selectors, ns)
		***REMOVED***

	***REMOVED***
	if len(selectors) == 1 && pseudoElement == "" ***REMOVED*** // no need wrap the selectors in compoundSelector
		return selectors[0], nil
	***REMOVED***
	return compoundSelector***REMOVED***selectors: selectors, pseudoElement: pseudoElement***REMOVED***, nil
***REMOVED***

// parseSelector parses a selector that may include combinators.
func (p *parser) parseSelector() (Sel, error) ***REMOVED***
	p.skipWhitespace()
	result, err := p.parseSimpleSelectorSequence()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	for ***REMOVED***
		var (
			combinator byte
			c          Sel
		)
		if p.skipWhitespace() ***REMOVED***
			combinator = ' '
		***REMOVED***
		if p.i >= len(p.s) ***REMOVED***
			return result, nil
		***REMOVED***

		switch p.s[p.i] ***REMOVED***
		case '+', '>', '~':
			combinator = p.s[p.i]
			p.i++
			p.skipWhitespace()
		case ',', ')':
			// These characters can't begin a selector, but they can legally occur after one.
			return result, nil
		***REMOVED***

		if combinator == 0 ***REMOVED***
			return result, nil
		***REMOVED***

		c, err = p.parseSimpleSelectorSequence()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result = combinedSelector***REMOVED***first: result, combinator: combinator, second: c***REMOVED***
	***REMOVED***
***REMOVED***

// parseSelectorGroup parses a group of selectors, separated by commas.
func (p *parser) parseSelectorGroup() (SelectorGroup, error) ***REMOVED***
	current, err := p.parseSelector()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	result := SelectorGroup***REMOVED***current***REMOVED***

	for p.i < len(p.s) ***REMOVED***
		if p.s[p.i] != ',' ***REMOVED***
			break
		***REMOVED***
		p.i++
		c, err := p.parseSelector()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result = append(result, c)
	***REMOVED***
	return result, nil
***REMOVED***
