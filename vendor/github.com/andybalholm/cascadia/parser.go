// Package cascadia is an implementation of CSS selectors.
package cascadia

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// a parser for CSS selectors
type parser struct ***REMOVED***
	s string // the source text
	i int    // the current position
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
		for i = start; i < p.i+6 && i < len(p.s) && hexDigit(p.s[i]); i++ ***REMOVED***
			// empty
		***REMOVED***
		v, _ := strconv.ParseUint(p.s[start:i], 16, 21)
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
func (p *parser) parseTypeSelector() (result Selector, err error) ***REMOVED***
	tag, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return typeSelector(tag), nil
***REMOVED***

// parseIDSelector parses a selector that matches by id attribute.
func (p *parser) parseIDSelector() (Selector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return nil, fmt.Errorf("expected id selector (#id), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '#' ***REMOVED***
		return nil, fmt.Errorf("expected id selector (#id), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	id, err := p.parseName()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return attributeEqualsSelector("id", id), nil
***REMOVED***

// parseClassSelector parses a selector that matches by class attribute.
func (p *parser) parseClassSelector() (Selector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return nil, fmt.Errorf("expected class selector (.class), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '.' ***REMOVED***
		return nil, fmt.Errorf("expected class selector (.class), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	class, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return attributeIncludesSelector("class", class), nil
***REMOVED***

// parseAttributeSelector parses a selector that matches by attribute value.
func (p *parser) parseAttributeSelector() (Selector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return nil, fmt.Errorf("expected attribute selector ([attribute]), found EOF instead")
	***REMOVED***
	if p.s[p.i] != '[' ***REMOVED***
		return nil, fmt.Errorf("expected attribute selector ([attribute]), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	p.skipWhitespace()
	key, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return nil, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	if p.s[p.i] == ']' ***REMOVED***
		p.i++
		return attributeExistsSelector(key), nil
	***REMOVED***

	if p.i+2 >= len(p.s) ***REMOVED***
		return nil, errors.New("unexpected EOF in attribute selector")
	***REMOVED***

	op := p.s[p.i : p.i+2]
	if op[0] == '=' ***REMOVED***
		op = "="
	***REMOVED*** else if op[1] != '=' ***REMOVED***
		return nil, fmt.Errorf(`expected equality operator, found "%s" instead`, op)
	***REMOVED***
	p.i += len(op)

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return nil, errors.New("unexpected EOF in attribute selector")
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
		return nil, err
	***REMOVED***

	p.skipWhitespace()
	if p.i >= len(p.s) ***REMOVED***
		return nil, errors.New("unexpected EOF in attribute selector")
	***REMOVED***
	if p.s[p.i] != ']' ***REMOVED***
		return nil, fmt.Errorf("expected ']', found '%c' instead", p.s[p.i])
	***REMOVED***
	p.i++

	switch op ***REMOVED***
	case "=":
		return attributeEqualsSelector(key, val), nil
	case "!=":
		return attributeNotEqualSelector(key, val), nil
	case "~=":
		return attributeIncludesSelector(key, val), nil
	case "|=":
		return attributeDashmatchSelector(key, val), nil
	case "^=":
		return attributePrefixSelector(key, val), nil
	case "$=":
		return attributeSuffixSelector(key, val), nil
	case "*=":
		return attributeSubstringSelector(key, val), nil
	case "#=":
		return attributeRegexSelector(key, rx), nil
	***REMOVED***

	return nil, fmt.Errorf("attribute operator %q is not supported", op)
***REMOVED***

var errExpectedParenthesis = errors.New("expected '(' but didn't find it")
var errExpectedClosingParenthesis = errors.New("expected ')' but didn't find it")
var errUnmatchedParenthesis = errors.New("unmatched '('")

// parsePseudoclassSelector parses a pseudoclass selector like :not(p).
func (p *parser) parsePseudoclassSelector() (Selector, error) ***REMOVED***
	if p.i >= len(p.s) ***REMOVED***
		return nil, fmt.Errorf("expected pseudoclass selector (:pseudoclass), found EOF instead")
	***REMOVED***
	if p.s[p.i] != ':' ***REMOVED***
		return nil, fmt.Errorf("expected attribute selector (:pseudoclass), found '%c' instead", p.s[p.i])
	***REMOVED***

	p.i++
	name, err := p.parseIdentifier()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	name = toLowerASCII(name)

	switch name ***REMOVED***
	case "not", "has", "haschild":
		if !p.consumeParenthesis() ***REMOVED***
			return nil, errExpectedParenthesis
		***REMOVED***
		sel, parseErr := p.parseSelectorGroup()
		if parseErr != nil ***REMOVED***
			return nil, parseErr
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return nil, errExpectedClosingParenthesis
		***REMOVED***

		switch name ***REMOVED***
		case "not":
			return negatedSelector(sel), nil
		case "has":
			return hasDescendantSelector(sel), nil
		case "haschild":
			return hasChildSelector(sel), nil
		***REMOVED***

	case "contains", "containsown":
		if !p.consumeParenthesis() ***REMOVED***
			return nil, errExpectedParenthesis
		***REMOVED***
		if p.i == len(p.s) ***REMOVED***
			return nil, errUnmatchedParenthesis
		***REMOVED***
		var val string
		switch p.s[p.i] ***REMOVED***
		case '\'', '"':
			val, err = p.parseString()
		default:
			val, err = p.parseIdentifier()
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		val = strings.ToLower(val)
		p.skipWhitespace()
		if p.i >= len(p.s) ***REMOVED***
			return nil, errors.New("unexpected EOF in pseudo selector")
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return nil, errExpectedClosingParenthesis
		***REMOVED***

		switch name ***REMOVED***
		case "contains":
			return textSubstrSelector(val), nil
		case "containsown":
			return ownTextSubstrSelector(val), nil
		***REMOVED***

	case "matches", "matchesown":
		if !p.consumeParenthesis() ***REMOVED***
			return nil, errExpectedParenthesis
		***REMOVED***
		rx, err := p.parseRegex()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if p.i >= len(p.s) ***REMOVED***
			return nil, errors.New("unexpected EOF in pseudo selector")
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return nil, errExpectedClosingParenthesis
		***REMOVED***

		switch name ***REMOVED***
		case "matches":
			return textRegexSelector(rx), nil
		case "matchesown":
			return ownTextRegexSelector(rx), nil
		***REMOVED***

	case "nth-child", "nth-last-child", "nth-of-type", "nth-last-of-type":
		if !p.consumeParenthesis() ***REMOVED***
			return nil, errExpectedParenthesis
		***REMOVED***
		a, b, err := p.parseNth()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if !p.consumeClosingParenthesis() ***REMOVED***
			return nil, errExpectedClosingParenthesis
		***REMOVED***
		if a == 0 ***REMOVED***
			switch name ***REMOVED***
			case "nth-child":
				return simpleNthChildSelector(b, false), nil
			case "nth-of-type":
				return simpleNthChildSelector(b, true), nil
			case "nth-last-child":
				return simpleNthLastChildSelector(b, false), nil
			case "nth-last-of-type":
				return simpleNthLastChildSelector(b, true), nil
			***REMOVED***
		***REMOVED***
		return nthChildSelector(a, b,
				name == "nth-last-child" || name == "nth-last-of-type",
				name == "nth-of-type" || name == "nth-last-of-type"),
			nil

	case "first-child":
		return simpleNthChildSelector(1, false), nil
	case "last-child":
		return simpleNthLastChildSelector(1, false), nil
	case "first-of-type":
		return simpleNthChildSelector(1, true), nil
	case "last-of-type":
		return simpleNthLastChildSelector(1, true), nil
	case "only-child":
		return onlyChildSelector(false), nil
	case "only-of-type":
		return onlyChildSelector(true), nil
	case "input":
		return inputSelector, nil
	case "empty":
		return emptyElementSelector, nil
	case "root":
		return rootSelector, nil
	***REMOVED***

	return nil, fmt.Errorf("unknown pseudoclass :%s", name)
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
func (p *parser) parseSimpleSelectorSequence() (Selector, error) ***REMOVED***
	var result Selector

	if p.i >= len(p.s) ***REMOVED***
		return nil, errors.New("expected selector, found EOF instead")
	***REMOVED***

	switch p.s[p.i] ***REMOVED***
	case '*':
		// It's the universal selector. Just skip over it, since it doesn't affect the meaning.
		p.i++
	case '#', '.', '[', ':':
		// There's no type selector. Wait to process the other till the main loop.
	default:
		r, err := p.parseTypeSelector()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result = r
	***REMOVED***

loop:
	for p.i < len(p.s) ***REMOVED***
		var ns Selector
		var err error
		switch p.s[p.i] ***REMOVED***
		case '#':
			ns, err = p.parseIDSelector()
		case '.':
			ns, err = p.parseClassSelector()
		case '[':
			ns, err = p.parseAttributeSelector()
		case ':':
			ns, err = p.parsePseudoclassSelector()
		default:
			break loop
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if result == nil ***REMOVED***
			result = ns
		***REMOVED*** else ***REMOVED***
			result = intersectionSelector(result, ns)
		***REMOVED***
	***REMOVED***

	if result == nil ***REMOVED***
		result = func(n *html.Node) bool ***REMOVED***
			return n.Type == html.ElementNode
		***REMOVED***
	***REMOVED***

	return result, nil
***REMOVED***

// parseSelector parses a selector that may include combinators.
func (p *parser) parseSelector() (result Selector, err error) ***REMOVED***
	p.skipWhitespace()
	result, err = p.parseSimpleSelectorSequence()
	if err != nil ***REMOVED***
		return
	***REMOVED***

	for ***REMOVED***
		var combinator byte
		if p.skipWhitespace() ***REMOVED***
			combinator = ' '
		***REMOVED***
		if p.i >= len(p.s) ***REMOVED***
			return
		***REMOVED***

		switch p.s[p.i] ***REMOVED***
		case '+', '>', '~':
			combinator = p.s[p.i]
			p.i++
			p.skipWhitespace()
		case ',', ')':
			// These characters can't begin a selector, but they can legally occur after one.
			return
		***REMOVED***

		if combinator == 0 ***REMOVED***
			return
		***REMOVED***

		c, err := p.parseSimpleSelectorSequence()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***

		switch combinator ***REMOVED***
		case ' ':
			result = descendantSelector(result, c)
		case '>':
			result = childSelector(result, c)
		case '+':
			result = siblingSelector(result, c, true)
		case '~':
			result = siblingSelector(result, c, false)
		***REMOVED***
	***REMOVED***

	panic("unreachable")
***REMOVED***

// parseSelectorGroup parses a group of selectors, separated by commas.
func (p *parser) parseSelectorGroup() (result Selector, err error) ***REMOVED***
	result, err = p.parseSelector()
	if err != nil ***REMOVED***
		return
	***REMOVED***

	for p.i < len(p.s) ***REMOVED***
		if p.s[p.i] != ',' ***REMOVED***
			return result, nil
		***REMOVED***
		p.i++
		c, err := p.parseSelector()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		result = unionSelector(result, c)
	***REMOVED***

	return
***REMOVED***
