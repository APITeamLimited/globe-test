package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	WhitespaceChars = " \f\n\r\t\v\u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff"
)

type regexpParseError struct ***REMOVED***
	offset int
	err    string
***REMOVED***

type RegexpErrorIncompatible struct ***REMOVED***
	regexpParseError
***REMOVED***
type RegexpSyntaxError struct ***REMOVED***
	regexpParseError
***REMOVED***

func (s regexpParseError) Error() string ***REMOVED***
	return s.err
***REMOVED***

type _RegExp_parser struct ***REMOVED***
	str    string
	length int

	chr       rune // The current character
	chrOffset int  // The offset of current character
	offset    int  // The offset after current character (may be greater than 1)

	err error

	goRegexp   strings.Builder
	passOffset int
***REMOVED***

// TransformRegExp transforms a JavaScript pattern into  a Go "regexp" pattern.
//
// re2 (Go) cannot do backtracking, so the presence of a lookahead (?=) (?!) or
// backreference (\1, \2, ...) will cause an error.
//
// re2 (Go) has a different definition for \s: [\t\n\f\r ].
// The JavaScript definition, on the other hand, also includes \v, Unicode "Separator, Space", etc.
//
// If the pattern is valid, but incompatible (contains a lookahead or backreference),
// then this function returns an empty string an error of type RegexpErrorIncompatible.
//
// If the pattern is invalid (not valid even in JavaScript), then this function
// returns an empty string and a generic error.
func TransformRegExp(pattern string) (transformed string, err error) ***REMOVED***

	if pattern == "" ***REMOVED***
		return "", nil
	***REMOVED***

	parser := _RegExp_parser***REMOVED***
		str:    pattern,
		length: len(pattern),
	***REMOVED***
	err = parser.parse()
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	return parser.ResultString(), nil
***REMOVED***

func (self *_RegExp_parser) ResultString() string ***REMOVED***
	if self.passOffset != -1 ***REMOVED***
		return self.str[:self.passOffset]
	***REMOVED***
	return self.goRegexp.String()
***REMOVED***

func (self *_RegExp_parser) parse() (err error) ***REMOVED***
	self.read() // Pull in the first character
	self.scan()
	return self.err
***REMOVED***

func (self *_RegExp_parser) read() ***REMOVED***
	if self.offset < self.length ***REMOVED***
		self.chrOffset = self.offset
		chr, width := rune(self.str[self.offset]), 1
		if chr >= utf8.RuneSelf ***REMOVED*** // !ASCII
			chr, width = utf8.DecodeRuneInString(self.str[self.offset:])
			if chr == utf8.RuneError && width == 1 ***REMOVED***
				self.error(true, "Invalid UTF-8 character")
				return
			***REMOVED***
		***REMOVED***
		self.offset += width
		self.chr = chr
	***REMOVED*** else ***REMOVED***
		self.chrOffset = self.length
		self.chr = -1 // EOF
	***REMOVED***
***REMOVED***

func (self *_RegExp_parser) stopPassing() ***REMOVED***
	self.goRegexp.Grow(3 * len(self.str) / 2)
	self.goRegexp.WriteString(self.str[:self.passOffset])
	self.passOffset = -1
***REMOVED***

func (self *_RegExp_parser) write(p []byte) ***REMOVED***
	if self.passOffset != -1 ***REMOVED***
		self.stopPassing()
	***REMOVED***
	self.goRegexp.Write(p)
***REMOVED***

func (self *_RegExp_parser) writeByte(b byte) ***REMOVED***
	if self.passOffset != -1 ***REMOVED***
		self.stopPassing()
	***REMOVED***
	self.goRegexp.WriteByte(b)
***REMOVED***

func (self *_RegExp_parser) writeString(s string) ***REMOVED***
	if self.passOffset != -1 ***REMOVED***
		self.stopPassing()
	***REMOVED***
	self.goRegexp.WriteString(s)
***REMOVED***

func (self *_RegExp_parser) scan() ***REMOVED***
	for self.chr != -1 ***REMOVED***
		switch self.chr ***REMOVED***
		case '\\':
			self.read()
			self.scanEscape(false)
		case '(':
			self.pass()
			self.scanGroup()
		case '[':
			self.scanBracket()
		case ')':
			self.error(true, "Unmatched ')'")
			return
		case '.':
			self.writeString("[^\\r\\n]")
			self.read()
		default:
			self.pass()
		***REMOVED***
	***REMOVED***
***REMOVED***

// (...)
func (self *_RegExp_parser) scanGroup() ***REMOVED***
	str := self.str[self.chrOffset:]
	if len(str) > 1 ***REMOVED*** // A possibility of (?= or (?!
		if str[0] == '?' ***REMOVED***
			ch := str[1]
			switch ***REMOVED***
			case ch == '=' || ch == '!':
				self.error(false, "re2: Invalid (%s) <lookahead>", self.str[self.chrOffset:self.chrOffset+2])
				return
			case ch == '<':
				self.error(false, "re2: Invalid (%s) <lookbehind>", self.str[self.chrOffset:self.chrOffset+2])
				return
			case ch != ':':
				self.error(true, "Invalid group")
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for self.chr != -1 && self.chr != ')' ***REMOVED***
		switch self.chr ***REMOVED***
		case '\\':
			self.read()
			self.scanEscape(false)
		case '(':
			self.pass()
			self.scanGroup()
		case '[':
			self.scanBracket()
		case '.':
			self.writeString("[^\\r\\n]")
			self.read()
		default:
			self.pass()
			continue
		***REMOVED***
	***REMOVED***
	if self.chr != ')' ***REMOVED***
		self.error(true, "Unterminated group")
		return
	***REMOVED***
	self.pass()
***REMOVED***

// [...]
func (self *_RegExp_parser) scanBracket() ***REMOVED***
	str := self.str[self.chrOffset:]
	if strings.HasPrefix(str, "[]") ***REMOVED***
		// [] -- Empty character class
		self.writeString("[^\u0000-\U0001FFFF]")
		self.offset += 1
		self.read()
		return
	***REMOVED***

	if strings.HasPrefix(str, "[^]") ***REMOVED***
		self.writeString("[\u0000-\U0001FFFF]")
		self.offset += 2
		self.read()
		return
	***REMOVED***

	self.pass()
	for self.chr != -1 ***REMOVED***
		if self.chr == ']' ***REMOVED***
			break
		***REMOVED*** else if self.chr == '\\' ***REMOVED***
			self.read()
			self.scanEscape(true)
			continue
		***REMOVED***
		self.pass()
	***REMOVED***
	if self.chr != ']' ***REMOVED***
		self.error(true, "Unterminated character class")
		return
	***REMOVED***
	self.pass()
***REMOVED***

// \...
func (self *_RegExp_parser) scanEscape(inClass bool) ***REMOVED***
	offset := self.chrOffset

	var length, base uint32
	switch self.chr ***REMOVED***

	case '0', '1', '2', '3', '4', '5', '6', '7':
		var value int64
		size := 0
		for ***REMOVED***
			digit := int64(digitValue(self.chr))
			if digit >= 8 ***REMOVED***
				// Not a valid digit
				break
			***REMOVED***
			value = value*8 + digit
			self.read()
			size += 1
		***REMOVED***
		if size == 1 ***REMOVED*** // The number of characters read
			if value != 0 ***REMOVED***
				// An invalid backreference
				self.error(false, "re2: Invalid \\%d <backreference>", value)
				return
			***REMOVED***
			self.passString(offset-1, self.chrOffset)
			return
		***REMOVED***
		tmp := []byte***REMOVED***'\\', 'x', '0', 0***REMOVED***
		if value >= 16 ***REMOVED***
			tmp = tmp[0:2]
		***REMOVED*** else ***REMOVED***
			tmp = tmp[0:3]
		***REMOVED***
		tmp = strconv.AppendInt(tmp, value, 16)
		self.write(tmp)
		return

	case '8', '9':
		self.read()
		self.error(false, "re2: Invalid \\%s <backreference>", self.str[offset:self.chrOffset])
		return

	case 'x':
		self.read()
		length, base = 2, 16

	case 'u':
		self.read()
		if self.chr == '***REMOVED***' ***REMOVED***
			self.read()
			length, base = 0, 16
		***REMOVED*** else ***REMOVED***
			length, base = 4, 16
		***REMOVED***

	case 'b':
		if inClass ***REMOVED***
			self.write([]byte***REMOVED***'\\', 'x', '0', '8'***REMOVED***)
			self.read()
			return
		***REMOVED***
		fallthrough

	case 'B':
		fallthrough

	case 'd', 'D', 'w', 'W':
		// This is slightly broken, because ECMAScript
		// includes \v in \s, \S, while re2 does not
		fallthrough

	case '\\':
		fallthrough

	case 'f', 'n', 'r', 't', 'v':
		self.passString(offset-1, self.offset)
		self.read()
		return

	case 'c':
		self.read()
		var value int64
		if 'a' <= self.chr && self.chr <= 'z' ***REMOVED***
			value = int64(self.chr - 'a' + 1)
		***REMOVED*** else if 'A' <= self.chr && self.chr <= 'Z' ***REMOVED***
			value = int64(self.chr - 'A' + 1)
		***REMOVED*** else ***REMOVED***
			self.writeByte('c')
			return
		***REMOVED***
		tmp := []byte***REMOVED***'\\', 'x', '0', 0***REMOVED***
		if value >= 16 ***REMOVED***
			tmp = tmp[0:2]
		***REMOVED*** else ***REMOVED***
			tmp = tmp[0:3]
		***REMOVED***
		tmp = strconv.AppendInt(tmp, value, 16)
		self.write(tmp)
		self.read()
		return
	case 's':
		if inClass ***REMOVED***
			self.writeString(WhitespaceChars)
		***REMOVED*** else ***REMOVED***
			self.writeString("[" + WhitespaceChars + "]")
		***REMOVED***
		self.read()
		return
	case 'S':
		if inClass ***REMOVED***
			self.error(false, "S in class")
			return
		***REMOVED*** else ***REMOVED***
			self.writeString("[^" + WhitespaceChars + "]")
		***REMOVED***
		self.read()
		return
	default:
		// $ is an identifier character, so we have to have
		// a special case for it here
		if self.chr == '$' || self.chr < utf8.RuneSelf && !isIdentifierPart(self.chr) ***REMOVED***
			// A non-identifier character needs escaping
			self.passString(offset-1, self.offset)
			self.read()
			return
		***REMOVED***
		// Unescape the character for re2
		self.pass()
		return
	***REMOVED***

	// Otherwise, we're a \u.... or \x...
	valueOffset := self.chrOffset

	if length > 0 ***REMOVED***
		for length := length; length > 0; length-- ***REMOVED***
			digit := uint32(digitValue(self.chr))
			if digit >= base ***REMOVED***
				// Not a valid digit
				goto skip
			***REMOVED***
			self.read()
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for self.chr != '***REMOVED***' && self.chr != -1 ***REMOVED***
			digit := uint32(digitValue(self.chr))
			if digit >= base ***REMOVED***
				// Not a valid digit
				goto skip
			***REMOVED***
			self.read()
		***REMOVED***
	***REMOVED***

	if length == 4 || length == 0 ***REMOVED***
		self.write([]byte***REMOVED***
			'\\',
			'x',
			'***REMOVED***',
		***REMOVED***)
		self.passString(valueOffset, self.chrOffset)
		if length != 0 ***REMOVED***
			self.writeByte('***REMOVED***')
		***REMOVED***
	***REMOVED*** else if length == 2 ***REMOVED***
		self.passString(offset-1, valueOffset+2)
	***REMOVED*** else ***REMOVED***
		// Should never, ever get here...
		self.error(true, "re2: Illegal branch in scanEscape")
		return
	***REMOVED***

	return

skip:
	self.passString(offset, self.chrOffset)
***REMOVED***

func (self *_RegExp_parser) pass() ***REMOVED***
	if self.passOffset == self.chrOffset ***REMOVED***
		self.passOffset = self.offset
	***REMOVED*** else ***REMOVED***
		if self.passOffset != -1 ***REMOVED***
			self.stopPassing()
		***REMOVED***
		if self.chr != -1 ***REMOVED***
			self.goRegexp.WriteRune(self.chr)
		***REMOVED***
	***REMOVED***
	self.read()
***REMOVED***

func (self *_RegExp_parser) passString(start, end int) ***REMOVED***
	if self.passOffset == start ***REMOVED***
		self.passOffset = end
		return
	***REMOVED***
	if self.passOffset != -1 ***REMOVED***
		self.stopPassing()
	***REMOVED***
	self.goRegexp.WriteString(self.str[start:end])
***REMOVED***

func (self *_RegExp_parser) error(fatal bool, msg string, msgValues ...interface***REMOVED******REMOVED***) ***REMOVED***
	if self.err != nil ***REMOVED***
		return
	***REMOVED***
	e := regexpParseError***REMOVED***
		offset: self.offset,
		err:    fmt.Sprintf(msg, msgValues...),
	***REMOVED***
	if fatal ***REMOVED***
		self.err = RegexpSyntaxError***REMOVED***e***REMOVED***
	***REMOVED*** else ***REMOVED***
		self.err = RegexpErrorIncompatible***REMOVED***e***REMOVED***
	***REMOVED***
	self.offset = self.length
	self.chr = -1
***REMOVED***
