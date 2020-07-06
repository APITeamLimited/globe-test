package parser

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"unicode/utf16"
)

type _chr struct ***REMOVED***
	value rune
	width int
***REMOVED***

var matchIdentifier = regexp.MustCompile(`^[$_\p***REMOVED***L***REMOVED***][$_\p***REMOVED***L***REMOVED***\d***REMOVED***]*$`)

func isDecimalDigit(chr rune) bool ***REMOVED***
	return '0' <= chr && chr <= '9'
***REMOVED***

func IsIdentifier(s string) bool ***REMOVED***
	return matchIdentifier.MatchString(s)
***REMOVED***

func digitValue(chr rune) int ***REMOVED***
	switch ***REMOVED***
	case '0' <= chr && chr <= '9':
		return int(chr - '0')
	case 'a' <= chr && chr <= 'f':
		return int(chr - 'a' + 10)
	case 'A' <= chr && chr <= 'F':
		return int(chr - 'A' + 10)
	***REMOVED***
	return 16 // Larger than any legal digit value
***REMOVED***

func isDigit(chr rune, base int) bool ***REMOVED***
	return digitValue(chr) < base
***REMOVED***

func isIdentifierStart(chr rune) bool ***REMOVED***
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && unicode.IsLetter(chr)
***REMOVED***

func isIdentifierPart(chr rune) bool ***REMOVED***
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && (unicode.IsLetter(chr) || unicode.IsDigit(chr))
***REMOVED***

func (self *_parser) scanIdentifier() (string, error) ***REMOVED***
	offset := self.chrOffset
	parse := false
	for isIdentifierPart(self.chr) ***REMOVED***
		if self.chr == '\\' ***REMOVED***
			distance := self.chrOffset - offset
			self.read()
			if self.chr != 'u' ***REMOVED***
				return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", self.chr, string(self.chr))
			***REMOVED***
			parse = true
			var value rune
			for j := 0; j < 4; j++ ***REMOVED***
				self.read()
				decimal, ok := hex2decimal(byte(self.chr))
				if !ok ***REMOVED***
					return "", fmt.Errorf("Invalid identifier escape character: %c (%s)", self.chr, string(self.chr))
				***REMOVED***
				value = value<<4 | decimal
			***REMOVED***
			if value == '\\' ***REMOVED***
				return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
			***REMOVED*** else if distance == 0 ***REMOVED***
				if !isIdentifierStart(value) ***REMOVED***
					return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
				***REMOVED***
			***REMOVED*** else if distance > 0 ***REMOVED***
				if !isIdentifierPart(value) ***REMOVED***
					return "", fmt.Errorf("Invalid identifier escape value: %c (%s)", value, string(value))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		self.read()
	***REMOVED***
	literal := string(self.str[offset:self.chrOffset])
	if parse ***REMOVED***
		return parseStringLiteral(literal)
	***REMOVED***
	return literal, nil
***REMOVED***

// 7.2
func isLineWhiteSpace(chr rune) bool ***REMOVED***
	switch chr ***REMOVED***
	case '\u0009', '\u000b', '\u000c', '\u0020', '\u00a0', '\ufeff':
		return true
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return false
	case '\u0085':
		return false
	***REMOVED***
	return unicode.IsSpace(chr)
***REMOVED***

// 7.3
func isLineTerminator(chr rune) bool ***REMOVED***
	switch chr ***REMOVED***
	case '\u000a', '\u000d', '\u2028', '\u2029':
		return true
	***REMOVED***
	return false
***REMOVED***

func (self *_parser) scan() (tkn token.Token, literal string, idx file.Idx) ***REMOVED***

	self.implicitSemicolon = false

	for ***REMOVED***
		self.skipWhiteSpace()

		idx = self.idxOf(self.chrOffset)
		insertSemicolon := false

		switch chr := self.chr; ***REMOVED***
		case isIdentifierStart(chr):
			var err error
			literal, err = self.scanIdentifier()
			if err != nil ***REMOVED***
				tkn = token.ILLEGAL
				break
			***REMOVED***
			if len(literal) > 1 ***REMOVED***
				// Keywords are longer than 1 character, avoid lookup otherwise
				var strict bool
				tkn, strict = token.IsKeyword(literal)

				switch tkn ***REMOVED***

				case 0: // Not a keyword
					if literal == "true" || literal == "false" ***REMOVED***
						self.insertSemicolon = true
						tkn = token.BOOLEAN
						return
					***REMOVED*** else if literal == "null" ***REMOVED***
						self.insertSemicolon = true
						tkn = token.NULL
						return
					***REMOVED***

				case token.KEYWORD:
					tkn = token.KEYWORD
					if strict ***REMOVED***
						// TODO If strict and in strict mode, then this is not a break
						break
					***REMOVED***
					return

				case
					token.THIS,
					token.BREAK,
					token.THROW, // A newline after a throw is not allowed, but we need to detect it
					token.RETURN,
					token.CONTINUE,
					token.DEBUGGER:
					self.insertSemicolon = true
					return

				default:
					return

				***REMOVED***
			***REMOVED***
			self.insertSemicolon = true
			tkn = token.IDENTIFIER
			return
		case '0' <= chr && chr <= '9':
			self.insertSemicolon = true
			tkn, literal = self.scanNumericLiteral(false)
			return
		default:
			self.read()
			switch chr ***REMOVED***
			case -1:
				if self.insertSemicolon ***REMOVED***
					self.insertSemicolon = false
					self.implicitSemicolon = true
				***REMOVED***
				tkn = token.EOF
			case '\r', '\n', '\u2028', '\u2029':
				self.insertSemicolon = false
				self.implicitSemicolon = true
				continue
			case ':':
				tkn = token.COLON
			case '.':
				if digitValue(self.chr) < 10 ***REMOVED***
					insertSemicolon = true
					tkn, literal = self.scanNumericLiteral(true)
				***REMOVED*** else ***REMOVED***
					tkn = token.PERIOD
				***REMOVED***
			case ',':
				tkn = token.COMMA
			case ';':
				tkn = token.SEMICOLON
			case '(':
				tkn = token.LEFT_PARENTHESIS
			case ')':
				tkn = token.RIGHT_PARENTHESIS
				insertSemicolon = true
			case '[':
				tkn = token.LEFT_BRACKET
			case ']':
				tkn = token.RIGHT_BRACKET
				insertSemicolon = true
			case '***REMOVED***':
				tkn = token.LEFT_BRACE
			case '***REMOVED***':
				tkn = token.RIGHT_BRACE
				insertSemicolon = true
			case '+':
				tkn = self.switch3(token.PLUS, token.ADD_ASSIGN, '+', token.INCREMENT)
				if tkn == token.INCREMENT ***REMOVED***
					insertSemicolon = true
				***REMOVED***
			case '-':
				tkn = self.switch3(token.MINUS, token.SUBTRACT_ASSIGN, '-', token.DECREMENT)
				if tkn == token.DECREMENT ***REMOVED***
					insertSemicolon = true
				***REMOVED***
			case '*':
				tkn = self.switch2(token.MULTIPLY, token.MULTIPLY_ASSIGN)
			case '/':
				if self.chr == '/' ***REMOVED***
					self.skipSingleLineComment()
					continue
				***REMOVED*** else if self.chr == '*' ***REMOVED***
					self.skipMultiLineComment()
					continue
				***REMOVED*** else ***REMOVED***
					// Could be division, could be RegExp literal
					tkn = self.switch2(token.SLASH, token.QUOTIENT_ASSIGN)
					insertSemicolon = true
				***REMOVED***
			case '%':
				tkn = self.switch2(token.REMAINDER, token.REMAINDER_ASSIGN)
			case '^':
				tkn = self.switch2(token.EXCLUSIVE_OR, token.EXCLUSIVE_OR_ASSIGN)
			case '<':
				tkn = self.switch4(token.LESS, token.LESS_OR_EQUAL, '<', token.SHIFT_LEFT, token.SHIFT_LEFT_ASSIGN)
			case '>':
				tkn = self.switch6(token.GREATER, token.GREATER_OR_EQUAL, '>', token.SHIFT_RIGHT, token.SHIFT_RIGHT_ASSIGN, '>', token.UNSIGNED_SHIFT_RIGHT, token.UNSIGNED_SHIFT_RIGHT_ASSIGN)
			case '=':
				tkn = self.switch2(token.ASSIGN, token.EQUAL)
				if tkn == token.EQUAL && self.chr == '=' ***REMOVED***
					self.read()
					tkn = token.STRICT_EQUAL
				***REMOVED***
			case '!':
				tkn = self.switch2(token.NOT, token.NOT_EQUAL)
				if tkn == token.NOT_EQUAL && self.chr == '=' ***REMOVED***
					self.read()
					tkn = token.STRICT_NOT_EQUAL
				***REMOVED***
			case '&':
				if self.chr == '^' ***REMOVED***
					self.read()
					tkn = self.switch2(token.AND_NOT, token.AND_NOT_ASSIGN)
				***REMOVED*** else ***REMOVED***
					tkn = self.switch3(token.AND, token.AND_ASSIGN, '&', token.LOGICAL_AND)
				***REMOVED***
			case '|':
				tkn = self.switch3(token.OR, token.OR_ASSIGN, '|', token.LOGICAL_OR)
			case '~':
				tkn = token.BITWISE_NOT
			case '?':
				tkn = token.QUESTION_MARK
			case '"', '\'':
				insertSemicolon = true
				tkn = token.STRING
				var err error
				literal, err = self.scanString(self.chrOffset - 1)
				if err != nil ***REMOVED***
					tkn = token.ILLEGAL
				***REMOVED***
			default:
				self.errorUnexpected(idx, chr)
				tkn = token.ILLEGAL
			***REMOVED***
		***REMOVED***
		self.insertSemicolon = insertSemicolon
		return
	***REMOVED***
***REMOVED***

func (self *_parser) switch2(tkn0, tkn1 token.Token) token.Token ***REMOVED***
	if self.chr == '=' ***REMOVED***
		self.read()
		return tkn1
	***REMOVED***
	return tkn0
***REMOVED***

func (self *_parser) switch3(tkn0, tkn1 token.Token, chr2 rune, tkn2 token.Token) token.Token ***REMOVED***
	if self.chr == '=' ***REMOVED***
		self.read()
		return tkn1
	***REMOVED***
	if self.chr == chr2 ***REMOVED***
		self.read()
		return tkn2
	***REMOVED***
	return tkn0
***REMOVED***

func (self *_parser) switch4(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token) token.Token ***REMOVED***
	if self.chr == '=' ***REMOVED***
		self.read()
		return tkn1
	***REMOVED***
	if self.chr == chr2 ***REMOVED***
		self.read()
		if self.chr == '=' ***REMOVED***
			self.read()
			return tkn3
		***REMOVED***
		return tkn2
	***REMOVED***
	return tkn0
***REMOVED***

func (self *_parser) switch6(tkn0, tkn1 token.Token, chr2 rune, tkn2, tkn3 token.Token, chr3 rune, tkn4, tkn5 token.Token) token.Token ***REMOVED***
	if self.chr == '=' ***REMOVED***
		self.read()
		return tkn1
	***REMOVED***
	if self.chr == chr2 ***REMOVED***
		self.read()
		if self.chr == '=' ***REMOVED***
			self.read()
			return tkn3
		***REMOVED***
		if self.chr == chr3 ***REMOVED***
			self.read()
			if self.chr == '=' ***REMOVED***
				self.read()
				return tkn5
			***REMOVED***
			return tkn4
		***REMOVED***
		return tkn2
	***REMOVED***
	return tkn0
***REMOVED***

func (self *_parser) chrAt(index int) _chr ***REMOVED***
	value, width := utf8.DecodeRuneInString(self.str[index:])
	return _chr***REMOVED***
		value: value,
		width: width,
	***REMOVED***
***REMOVED***

func (self *_parser) _peek() rune ***REMOVED***
	if self.offset+1 < self.length ***REMOVED***
		return rune(self.str[self.offset+1])
	***REMOVED***
	return -1
***REMOVED***

func (self *_parser) read() ***REMOVED***
	if self.offset < self.length ***REMOVED***
		self.chrOffset = self.offset
		chr, width := rune(self.str[self.offset]), 1
		if chr >= utf8.RuneSelf ***REMOVED*** // !ASCII
			chr, width = utf8.DecodeRuneInString(self.str[self.offset:])
			if chr == utf8.RuneError && width == 1 ***REMOVED***
				self.error(self.chrOffset, "Invalid UTF-8 character")
			***REMOVED***
		***REMOVED***
		self.offset += width
		self.chr = chr
	***REMOVED*** else ***REMOVED***
		self.chrOffset = self.length
		self.chr = -1 // EOF
	***REMOVED***
***REMOVED***

// This is here since the functions are so similar
func (self *_RegExp_parser) read() ***REMOVED***
	if self.offset < self.length ***REMOVED***
		self.chrOffset = self.offset
		chr, width := rune(self.str[self.offset]), 1
		if chr >= utf8.RuneSelf ***REMOVED*** // !ASCII
			chr, width = utf8.DecodeRuneInString(self.str[self.offset:])
			if chr == utf8.RuneError && width == 1 ***REMOVED***
				self.error(self.chrOffset, "Invalid UTF-8 character")
			***REMOVED***
		***REMOVED***
		self.offset += width
		self.chr = chr
	***REMOVED*** else ***REMOVED***
		self.chrOffset = self.length
		self.chr = -1 // EOF
	***REMOVED***
***REMOVED***

func (self *_parser) skipSingleLineComment() ***REMOVED***
	for self.chr != -1 ***REMOVED***
		self.read()
		if isLineTerminator(self.chr) ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func (self *_parser) skipMultiLineComment() ***REMOVED***
	self.read()
	for self.chr >= 0 ***REMOVED***
		chr := self.chr
		self.read()
		if chr == '*' && self.chr == '/' ***REMOVED***
			self.read()
			return
		***REMOVED***
	***REMOVED***

	self.errorUnexpected(0, self.chr)
***REMOVED***

func (self *_parser) skipWhiteSpace() ***REMOVED***
	for ***REMOVED***
		switch self.chr ***REMOVED***
		case ' ', '\t', '\f', '\v', '\u00a0', '\ufeff':
			self.read()
			continue
		case '\r':
			if self._peek() == '\n' ***REMOVED***
				self.read()
			***REMOVED***
			fallthrough
		case '\u2028', '\u2029', '\n':
			if self.insertSemicolon ***REMOVED***
				return
			***REMOVED***
			self.read()
			continue
		***REMOVED***
		if self.chr >= utf8.RuneSelf ***REMOVED***
			if unicode.IsSpace(self.chr) ***REMOVED***
				self.read()
				continue
			***REMOVED***
		***REMOVED***
		break
	***REMOVED***
***REMOVED***

func (self *_parser) skipLineWhiteSpace() ***REMOVED***
	for isLineWhiteSpace(self.chr) ***REMOVED***
		self.read()
	***REMOVED***
***REMOVED***

func (self *_parser) scanMantissa(base int) ***REMOVED***
	for digitValue(self.chr) < base ***REMOVED***
		self.read()
	***REMOVED***
***REMOVED***

func (self *_parser) scanEscape(quote rune) ***REMOVED***

	var length, base uint32
	switch self.chr ***REMOVED***
	//case '0', '1', '2', '3', '4', '5', '6', '7':
	//    Octal:
	//    length, base, limit = 3, 8, 255
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'', '0':
		self.read()
		return
	case '\r', '\n', '\u2028', '\u2029':
		self.scanNewline()
		return
	case 'x':
		self.read()
		length, base = 2, 16
	case 'u':
		self.read()
		length, base = 4, 16
	default:
		self.read() // Always make progress
		return
	***REMOVED***

	var value uint32
	for ; length > 0 && self.chr != quote && self.chr >= 0; length-- ***REMOVED***
		digit := uint32(digitValue(self.chr))
		if digit >= base ***REMOVED***
			break
		***REMOVED***
		value = value*base + digit
		self.read()
	***REMOVED***
***REMOVED***

func (self *_parser) scanString(offset int) (string, error) ***REMOVED***
	// " ' /
	quote := rune(self.str[offset])

	for self.chr != quote ***REMOVED***
		chr := self.chr
		if chr == '\n' || chr == '\r' || chr == '\u2028' || chr == '\u2029' || chr < 0 ***REMOVED***
			goto newline
		***REMOVED***
		self.read()
		if chr == '\\' ***REMOVED***
			if quote == '/' ***REMOVED***
				if self.chr == '\n' || self.chr == '\r' || self.chr == '\u2028' || self.chr == '\u2029' || self.chr < 0 ***REMOVED***
					goto newline
				***REMOVED***
				self.read()
			***REMOVED*** else ***REMOVED***
				self.scanEscape(quote)
			***REMOVED***
		***REMOVED*** else if chr == '[' && quote == '/' ***REMOVED***
			// Allow a slash (/) in a bracket character class ([...])
			// TODO Fix this, this is hacky...
			quote = -1
		***REMOVED*** else if chr == ']' && quote == -1 ***REMOVED***
			quote = '/'
		***REMOVED***
	***REMOVED***

	// " ' /
	self.read()

	return string(self.str[offset:self.chrOffset]), nil

newline:
	self.scanNewline()
	err := "String not terminated"
	if quote == '/' ***REMOVED***
		err = "Invalid regular expression: missing /"
		self.error(self.idxOf(offset), err)
	***REMOVED***
	return "", errors.New(err)
***REMOVED***

func (self *_parser) scanNewline() ***REMOVED***
	if self.chr == '\r' ***REMOVED***
		self.read()
		if self.chr != '\n' ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
	self.read()
***REMOVED***

func hex2decimal(chr byte) (value rune, ok bool) ***REMOVED***
	***REMOVED***
		chr := rune(chr)
		switch ***REMOVED***
		case '0' <= chr && chr <= '9':
			return chr - '0', true
		case 'a' <= chr && chr <= 'f':
			return chr - 'a' + 10, true
		case 'A' <= chr && chr <= 'F':
			return chr - 'A' + 10, true
		***REMOVED***
		return
	***REMOVED***
***REMOVED***

func parseNumberLiteral(literal string) (value interface***REMOVED******REMOVED***, err error) ***REMOVED***
	// TODO Is Uint okay? What about -MAX_UINT
	value, err = strconv.ParseInt(literal, 0, 64)
	if err == nil ***REMOVED***
		return
	***REMOVED***

	parseIntErr := err // Save this first error, just in case

	value, err = strconv.ParseFloat(literal, 64)
	if err == nil ***REMOVED***
		return
	***REMOVED*** else if err.(*strconv.NumError).Err == strconv.ErrRange ***REMOVED***
		// Infinity, etc.
		return value, nil
	***REMOVED***

	err = parseIntErr

	if err.(*strconv.NumError).Err == strconv.ErrRange ***REMOVED***
		if len(literal) > 2 && literal[0] == '0' && (literal[1] == 'X' || literal[1] == 'x') ***REMOVED***
			// Could just be a very large number (e.g. 0x8000000000000000)
			var value float64
			literal = literal[2:]
			for _, chr := range literal ***REMOVED***
				digit := digitValue(chr)
				if digit >= 16 ***REMOVED***
					goto error
				***REMOVED***
				value = value*16 + float64(digit)
			***REMOVED***
			return value, nil
		***REMOVED***
	***REMOVED***

error:
	return nil, errors.New("Illegal numeric literal")
***REMOVED***

func parseStringLiteral(literal string) (string, error) ***REMOVED***
	// Best case scenario...
	if literal == "" ***REMOVED***
		return "", nil
	***REMOVED***

	// Slightly less-best case scenario...
	if !strings.ContainsRune(literal, '\\') ***REMOVED***
		return literal, nil
	***REMOVED***

	str := literal
	buffer := bytes.NewBuffer(make([]byte, 0, 3*len(literal)/2))
	var surrogate rune
S:
	for len(str) > 0 ***REMOVED***
		switch chr := str[0]; ***REMOVED***
		// We do not explicitly handle the case of the quote
		// value, which can be: " ' /
		// This assumes we're already passed a partially well-formed literal
		case chr >= utf8.RuneSelf:
			chr, size := utf8.DecodeRuneInString(str)
			buffer.WriteRune(chr)
			str = str[size:]
			continue
		case chr != '\\':
			buffer.WriteByte(chr)
			str = str[1:]
			continue
		***REMOVED***

		if len(str) <= 1 ***REMOVED***
			panic("len(str) <= 1")
		***REMOVED***
		chr := str[1]
		var value rune
		if chr >= utf8.RuneSelf ***REMOVED***
			str = str[1:]
			var size int
			value, size = utf8.DecodeRuneInString(str)
			str = str[size:] // \ + <character>
		***REMOVED*** else ***REMOVED***
			str = str[2:] // \<character>
			switch chr ***REMOVED***
			case 'b':
				value = '\b'
			case 'f':
				value = '\f'
			case 'n':
				value = '\n'
			case 'r':
				value = '\r'
			case 't':
				value = '\t'
			case 'v':
				value = '\v'
			case 'x', 'u':
				size := 0
				switch chr ***REMOVED***
				case 'x':
					size = 2
				case 'u':
					size = 4
				***REMOVED***
				if len(str) < size ***REMOVED***
					return "", fmt.Errorf("invalid escape: \\%s: len(%q) != %d", string(chr), str, size)
				***REMOVED***
				for j := 0; j < size; j++ ***REMOVED***
					decimal, ok := hex2decimal(str[j])
					if !ok ***REMOVED***
						return "", fmt.Errorf("invalid escape: \\%s: %q", string(chr), str[:size])
					***REMOVED***
					value = value<<4 | decimal
				***REMOVED***
				str = str[size:]
				if chr == 'x' ***REMOVED***
					break
				***REMOVED***
				if value > utf8.MaxRune ***REMOVED***
					panic("value > utf8.MaxRune")
				***REMOVED***
			case '0':
				if len(str) == 0 || '0' > str[0] || str[0] > '7' ***REMOVED***
					value = 0
					break
				***REMOVED***
				fallthrough
			case '1', '2', '3', '4', '5', '6', '7':
				// TODO strict
				value = rune(chr) - '0'
				j := 0
				for ; j < 2; j++ ***REMOVED***
					if len(str) < j+1 ***REMOVED***
						break
					***REMOVED***
					chr := str[j]
					if '0' > chr || chr > '7' ***REMOVED***
						break
					***REMOVED***
					decimal := rune(str[j]) - '0'
					value = (value << 3) | decimal
				***REMOVED***
				str = str[j:]
			case '\\':
				value = '\\'
			case '\'', '"':
				value = rune(chr)
			case '\r':
				if len(str) > 0 ***REMOVED***
					if str[0] == '\n' ***REMOVED***
						str = str[1:]
					***REMOVED***
				***REMOVED***
				fallthrough
			case '\n':
				continue
			default:
				value = rune(chr)
			***REMOVED***
			if surrogate != 0 ***REMOVED***
				value = utf16.DecodeRune(surrogate, value)
				surrogate = 0
			***REMOVED*** else ***REMOVED***
				if utf16.IsSurrogate(value) ***REMOVED***
					surrogate = value
					continue S
				***REMOVED***
			***REMOVED***
		***REMOVED***
		buffer.WriteRune(value)
	***REMOVED***

	return buffer.String(), nil
***REMOVED***

func (self *_parser) scanNumericLiteral(decimalPoint bool) (token.Token, string) ***REMOVED***

	offset := self.chrOffset
	tkn := token.NUMBER

	if decimalPoint ***REMOVED***
		offset--
		self.scanMantissa(10)
		goto exponent
	***REMOVED***

	if self.chr == '0' ***REMOVED***
		offset := self.chrOffset
		self.read()
		if self.chr == 'x' || self.chr == 'X' ***REMOVED***
			// Hexadecimal
			self.read()
			if isDigit(self.chr, 16) ***REMOVED***
				self.read()
			***REMOVED*** else ***REMOVED***
				return token.ILLEGAL, self.str[offset:self.chrOffset]
			***REMOVED***
			self.scanMantissa(16)

			if self.chrOffset-offset <= 2 ***REMOVED***
				// Only "0x" or "0X"
				self.error(0, "Illegal hexadecimal number")
			***REMOVED***

			goto hexadecimal
		***REMOVED*** else if self.chr == '.' ***REMOVED***
			// Float
			goto float
		***REMOVED*** else ***REMOVED***
			// Octal, Float
			if self.chr == 'e' || self.chr == 'E' ***REMOVED***
				goto exponent
			***REMOVED***
			self.scanMantissa(8)
			if self.chr == '8' || self.chr == '9' ***REMOVED***
				return token.ILLEGAL, self.str[offset:self.chrOffset]
			***REMOVED***
			goto octal
		***REMOVED***
	***REMOVED***

	self.scanMantissa(10)

float:
	if self.chr == '.' ***REMOVED***
		self.read()
		self.scanMantissa(10)
	***REMOVED***

exponent:
	if self.chr == 'e' || self.chr == 'E' ***REMOVED***
		self.read()
		if self.chr == '-' || self.chr == '+' ***REMOVED***
			self.read()
		***REMOVED***
		if isDecimalDigit(self.chr) ***REMOVED***
			self.read()
			self.scanMantissa(10)
		***REMOVED*** else ***REMOVED***
			return token.ILLEGAL, self.str[offset:self.chrOffset]
		***REMOVED***
	***REMOVED***

hexadecimal:
octal:
	if isIdentifierStart(self.chr) || isDecimalDigit(self.chr) ***REMOVED***
		return token.ILLEGAL, self.str[offset:self.chrOffset]
	***REMOVED***

	return tkn, self.str[offset:self.chrOffset]
***REMOVED***
