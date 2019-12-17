package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	WhitespaceChars = " \f\n\r\t\v\u00a0\u1680\u2000\u2001\u2002\u2003\u2004\u2005\u2006\u2007\u2008\u2009\u200a\u2028\u2029\u202f\u205f\u3000\ufeff"
)

type _RegExp_parser struct ***REMOVED***
	str    string
	length int

	chr       rune // The current character
	chrOffset int  // The offset of current character
	offset    int  // The offset after current character (may be greater than 1)

	errors  []error
	invalid bool // The input is an invalid JavaScript RegExp

	goRegexp *bytes.Buffer
***REMOVED***

// TransformRegExp transforms a JavaScript pattern into  a Go "regexp" pattern.
//
// re2 (Go) cannot do backtracking, so the presence of a lookahead (?=) (?!) or
// backreference (\1, \2, ...) will cause an error.
//
// re2 (Go) has a different definition for \s: [\t\n\f\r ].
// The JavaScript definition, on the other hand, also includes \v, Unicode "Separator, Space", etc.
//
// If the pattern is invalid (not valid even in JavaScript), then this function
// returns the empty string and an error.
//
// If the pattern is valid, but incompatible (contains a lookahead or backreference),
// then this function returns the transformation (a non-empty string) AND an error.
func TransformRegExp(pattern string) (string, error) ***REMOVED***

	if pattern == "" ***REMOVED***
		return "", nil
	***REMOVED***

	// TODO If without \, if without (?=, (?!, then another shortcut

	parser := _RegExp_parser***REMOVED***
		str:      pattern,
		length:   len(pattern),
		goRegexp: bytes.NewBuffer(make([]byte, 0, 3*len(pattern)/2)),
	***REMOVED***
	parser.read() // Pull in the first character
	parser.scan()
	var err error
	if len(parser.errors) > 0 ***REMOVED***
		err = parser.errors[0]
	***REMOVED***
	if parser.invalid ***REMOVED***
		return "", err
	***REMOVED***

	// Might not be re2 compatible, but is still a valid JavaScript RegExp
	return parser.goRegexp.String(), err
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
			self.error(-1, "Unmatched ')'")
			self.invalid = true
			self.pass()
		case '.':
			self.goRegexp.WriteString("[^\\r\\n]")
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
			if str[1] == '=' || str[1] == '!' ***REMOVED***
				self.error(-1, "re2: Invalid (%s) <lookahead>", self.str[self.chrOffset:self.chrOffset+2])
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
			self.goRegexp.WriteString("[^\\r\\n]")
			self.read()
		default:
			self.pass()
			continue
		***REMOVED***
	***REMOVED***
	if self.chr != ')' ***REMOVED***
		self.error(-1, "Unterminated group")
		self.invalid = true
		return
	***REMOVED***
	self.pass()
***REMOVED***

// [...]
func (self *_RegExp_parser) scanBracket() ***REMOVED***
	str := self.str[self.chrOffset:]
	if strings.HasPrefix(str, "[]") ***REMOVED***
		// [] -- Empty character class
		self.goRegexp.WriteString("[^\u0000-uffff]")
		self.offset += 1
		self.read()
		return
	***REMOVED***

	if strings.HasPrefix(str, "[^]") ***REMOVED***
		self.goRegexp.WriteString("[\u0000-\uffff]")
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
		self.error(-1, "Unterminated character class")
		self.invalid = true
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
			_, err := self.goRegexp.Write([]byte***REMOVED***'\\', byte(value) + '0'***REMOVED***)
			if err != nil ***REMOVED***
				self.errors = append(self.errors, err)
			***REMOVED***
			if value != 0 ***REMOVED***
				// An invalid backreference
				self.error(-1, "re2: Invalid \\%d <backreference>", value)
			***REMOVED***
			return
		***REMOVED***
		tmp := []byte***REMOVED***'\\', 'x', '0', 0***REMOVED***
		if value >= 16 ***REMOVED***
			tmp = tmp[0:2]
		***REMOVED*** else ***REMOVED***
			tmp = tmp[0:3]
		***REMOVED***
		tmp = strconv.AppendInt(tmp, value, 16)
		_, err := self.goRegexp.Write(tmp)
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
		return

	case '8', '9':
		size := 0
		for ***REMOVED***
			digit := digitValue(self.chr)
			if digit >= 10 ***REMOVED***
				// Not a valid digit
				break
			***REMOVED***
			self.read()
			size += 1
		***REMOVED***
		err := self.goRegexp.WriteByte('\\')
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
		_, err = self.goRegexp.WriteString(self.str[offset:self.chrOffset])
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
		self.error(-1, "re2: Invalid \\%s <backreference>", self.str[offset:self.chrOffset])
		return

	case 'x':
		self.read()
		length, base = 2, 16

	case 'u':
		self.read()
		length, base = 4, 16

	case 'b':
		if inClass ***REMOVED***
			_, err := self.goRegexp.Write([]byte***REMOVED***'\\', 'x', '0', '8'***REMOVED***)
			if err != nil ***REMOVED***
				self.errors = append(self.errors, err)
			***REMOVED***
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
		err := self.goRegexp.WriteByte('\\')
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
		self.pass()
		return

	case 'c':
		self.read()
		var value int64
		if 'a' <= self.chr && self.chr <= 'z' ***REMOVED***
			value = int64(self.chr) - 'a' + 1
		***REMOVED*** else if 'A' <= self.chr && self.chr <= 'Z' ***REMOVED***
			value = int64(self.chr) - 'A' + 1
		***REMOVED*** else ***REMOVED***
			err := self.goRegexp.WriteByte('c')
			if err != nil ***REMOVED***
				self.errors = append(self.errors, err)
			***REMOVED***
			return
		***REMOVED***
		tmp := []byte***REMOVED***'\\', 'x', '0', 0***REMOVED***
		if value >= 16 ***REMOVED***
			tmp = tmp[0:2]
		***REMOVED*** else ***REMOVED***
			tmp = tmp[0:3]
		***REMOVED***
		tmp = strconv.AppendInt(tmp, value, 16)
		_, err := self.goRegexp.Write(tmp)
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
		self.read()
		return
	case 's':
		if inClass ***REMOVED***
			self.goRegexp.WriteString(WhitespaceChars)
		***REMOVED*** else ***REMOVED***
			self.goRegexp.WriteString("[" + WhitespaceChars + "]")
		***REMOVED***
		self.read()
		return
	case 'S':
		if inClass ***REMOVED***
			self.error(self.chrOffset, "S in class")
			self.invalid = true
			return
		***REMOVED*** else ***REMOVED***
			self.goRegexp.WriteString("[^" + WhitespaceChars + "]")
		***REMOVED***
		self.read()
		return
	default:
		// $ is an identifier character, so we have to have
		// a special case for it here
		if self.chr == '$' || self.chr < utf8.RuneSelf && !isIdentifierPart(self.chr) ***REMOVED***
			// A non-identifier character needs escaping
			err := self.goRegexp.WriteByte('\\')
			if err != nil ***REMOVED***
				self.errors = append(self.errors, err)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// Unescape the character for re2
		***REMOVED***
		self.pass()
		return
	***REMOVED***

	// Otherwise, we're a \u.... or \x...
	valueOffset := self.chrOffset

	var value uint32
	***REMOVED***
		length := length
		for ; length > 0; length-- ***REMOVED***
			digit := uint32(digitValue(self.chr))
			if digit >= base ***REMOVED***
				// Not a valid digit
				goto skip
			***REMOVED***
			value = value*base + digit
			self.read()
		***REMOVED***
	***REMOVED***

	if length == 4 ***REMOVED***
		_, err := self.goRegexp.Write([]byte***REMOVED***
			'\\',
			'x',
			'***REMOVED***',
			self.str[valueOffset+0],
			self.str[valueOffset+1],
			self.str[valueOffset+2],
			self.str[valueOffset+3],
			'***REMOVED***',
		***REMOVED***)
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
	***REMOVED*** else if length == 2 ***REMOVED***
		_, err := self.goRegexp.Write([]byte***REMOVED***
			'\\',
			'x',
			self.str[valueOffset+0],
			self.str[valueOffset+1],
		***REMOVED***)
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// Should never, ever get here...
		self.error(-1, "re2: Illegal branch in scanEscape")
		goto skip
	***REMOVED***

	return

skip:
	_, err := self.goRegexp.WriteString(self.str[offset:self.chrOffset])
	if err != nil ***REMOVED***
		self.errors = append(self.errors, err)
	***REMOVED***
***REMOVED***

func (self *_RegExp_parser) pass() ***REMOVED***
	if self.chr != -1 ***REMOVED***
		_, err := self.goRegexp.WriteRune(self.chr)
		if err != nil ***REMOVED***
			self.errors = append(self.errors, err)
		***REMOVED***
	***REMOVED***
	self.read()
***REMOVED***

// TODO Better error reporting, use the offset, etc.
func (self *_RegExp_parser) error(offset int, msg string, msgValues ...interface***REMOVED******REMOVED***) error ***REMOVED***
	err := fmt.Errorf(msg, msgValues...)
	self.errors = append(self.errors, err)
	return err
***REMOVED***
