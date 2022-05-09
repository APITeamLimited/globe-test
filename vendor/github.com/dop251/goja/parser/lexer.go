package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"golang.org/x/text/unicode/rangetable"

	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
	"github.com/dop251/goja/unistring"
)

var (
	unicodeRangeIdNeg      = rangetable.Merge(unicode.Pattern_Syntax, unicode.Pattern_White_Space)
	unicodeRangeIdStartPos = rangetable.Merge(unicode.Letter, unicode.Nl, unicode.Other_ID_Start)
	unicodeRangeIdContPos  = rangetable.Merge(unicodeRangeIdStartPos, unicode.Mn, unicode.Mc, unicode.Nd, unicode.Pc, unicode.Other_ID_Continue)
)

func isDecimalDigit(chr rune) bool ***REMOVED***
	return '0' <= chr && chr <= '9'
***REMOVED***

func IsIdentifier(s string) bool ***REMOVED***
	if s == "" ***REMOVED***
		return false
	***REMOVED***
	r, size := utf8.DecodeRuneInString(s)
	if !isIdentifierStart(r) ***REMOVED***
		return false
	***REMOVED***
	for _, r := range s[size:] ***REMOVED***
		if !isIdentifierPart(r) ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
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

func isIdStartUnicode(r rune) bool ***REMOVED***
	return unicode.Is(unicodeRangeIdStartPos, r) && !unicode.Is(unicodeRangeIdNeg, r)
***REMOVED***

func isIdPartUnicode(r rune) bool ***REMOVED***
	return unicode.Is(unicodeRangeIdContPos, r) && !unicode.Is(unicodeRangeIdNeg, r) || r == '\u200C' || r == '\u200D'
***REMOVED***

func isIdentifierStart(chr rune) bool ***REMOVED***
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		chr >= utf8.RuneSelf && isIdStartUnicode(chr)
***REMOVED***

func isIdentifierPart(chr rune) bool ***REMOVED***
	return chr == '$' || chr == '_' || chr == '\\' ||
		'a' <= chr && chr <= 'z' || 'A' <= chr && chr <= 'Z' ||
		'0' <= chr && chr <= '9' ||
		chr >= utf8.RuneSelf && isIdPartUnicode(chr)
***REMOVED***

func (self *_parser) scanIdentifier() (string, unistring.String, bool, string) ***REMOVED***
	offset := self.chrOffset
	hasEscape := false
	isUnicode := false
	length := 0
	for isIdentifierPart(self.chr) ***REMOVED***
		r := self.chr
		length++
		if r == '\\' ***REMOVED***
			hasEscape = true
			distance := self.chrOffset - offset
			self.read()
			if self.chr != 'u' ***REMOVED***
				return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", self.chr, string(self.chr))
			***REMOVED***
			var value rune
			if self._peek() == '***REMOVED***' ***REMOVED***
				self.read()
				value = -1
				for value <= utf8.MaxRune ***REMOVED***
					self.read()
					if self.chr == '***REMOVED***' ***REMOVED***
						break
					***REMOVED***
					decimal, ok := hex2decimal(byte(self.chr))
					if !ok ***REMOVED***
						return "", "", false, "Invalid Unicode escape sequence"
					***REMOVED***
					if value == -1 ***REMOVED***
						value = decimal
					***REMOVED*** else ***REMOVED***
						value = value<<4 | decimal
					***REMOVED***
				***REMOVED***
				if value == -1 ***REMOVED***
					return "", "", false, "Invalid Unicode escape sequence"
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				for j := 0; j < 4; j++ ***REMOVED***
					self.read()
					decimal, ok := hex2decimal(byte(self.chr))
					if !ok ***REMOVED***
						return "", "", false, fmt.Sprintf("Invalid identifier escape character: %c (%s)", self.chr, string(self.chr))
					***REMOVED***
					value = value<<4 | decimal
				***REMOVED***
			***REMOVED***
			if value == '\\' ***REMOVED***
				return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
			***REMOVED*** else if distance == 0 ***REMOVED***
				if !isIdentifierStart(value) ***REMOVED***
					return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
				***REMOVED***
			***REMOVED*** else if distance > 0 ***REMOVED***
				if !isIdentifierPart(value) ***REMOVED***
					return "", "", false, fmt.Sprintf("Invalid identifier escape value: %c (%s)", value, string(value))
				***REMOVED***
			***REMOVED***
			r = value
		***REMOVED***
		if r >= utf8.RuneSelf ***REMOVED***
			isUnicode = true
			if r > 0xFFFF ***REMOVED***
				length++
			***REMOVED***
		***REMOVED***
		self.read()
	***REMOVED***

	literal := self.str[offset:self.chrOffset]
	var parsed unistring.String
	if hasEscape || isUnicode ***REMOVED***
		var err string
		// TODO strict
		parsed, err = parseStringLiteral(literal, length, isUnicode, false)
		if err != "" ***REMOVED***
			return "", "", false, err
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		parsed = unistring.String(literal)
	***REMOVED***

	return literal, parsed, hasEscape, ""
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

type parserState struct ***REMOVED***
	tok                                token.Token
	literal                            string
	parsedLiteral                      unistring.String
	implicitSemicolon, insertSemicolon bool
	chr                                rune
	chrOffset, offset                  int
	errorCount                         int
***REMOVED***

func (self *_parser) mark(state *parserState) *parserState ***REMOVED***
	if state == nil ***REMOVED***
		state = &parserState***REMOVED******REMOVED***
	***REMOVED***
	state.tok, state.literal, state.parsedLiteral, state.implicitSemicolon, state.insertSemicolon, state.chr, state.chrOffset, state.offset =
		self.token, self.literal, self.parsedLiteral, self.implicitSemicolon, self.insertSemicolon, self.chr, self.chrOffset, self.offset

	state.errorCount = len(self.errors)
	return state
***REMOVED***

func (self *_parser) restore(state *parserState) ***REMOVED***
	self.token, self.literal, self.parsedLiteral, self.implicitSemicolon, self.insertSemicolon, self.chr, self.chrOffset, self.offset =
		state.tok, state.literal, state.parsedLiteral, state.implicitSemicolon, state.insertSemicolon, state.chr, state.chrOffset, state.offset
	self.errors = self.errors[:state.errorCount]
***REMOVED***

func (self *_parser) peek() token.Token ***REMOVED***
	implicitSemicolon, insertSemicolon, chr, chrOffset, offset := self.implicitSemicolon, self.insertSemicolon, self.chr, self.chrOffset, self.offset
	tok, _, _, _ := self.scan()
	self.implicitSemicolon, self.insertSemicolon, self.chr, self.chrOffset, self.offset = implicitSemicolon, insertSemicolon, chr, chrOffset, offset
	return tok
***REMOVED***

func (self *_parser) scan() (tkn token.Token, literal string, parsedLiteral unistring.String, idx file.Idx) ***REMOVED***

	self.implicitSemicolon = false

	for ***REMOVED***
		self.skipWhiteSpace()

		idx = self.idxOf(self.chrOffset)
		insertSemicolon := false

		switch chr := self.chr; ***REMOVED***
		case isIdentifierStart(chr):
			var err string
			var hasEscape bool
			literal, parsedLiteral, hasEscape, err = self.scanIdentifier()
			if err != "" ***REMOVED***
				tkn = token.ILLEGAL
				break
			***REMOVED***
			if len(parsedLiteral) > 1 ***REMOVED***
				// Keywords are longer than 1 character, avoid lookup otherwise
				var strict bool
				tkn, strict = token.IsKeyword(string(parsedLiteral))
				if hasEscape ***REMOVED***
					self.insertSemicolon = true
					if tkn != 0 && tkn != token.LET || parsedLiteral == "true" || parsedLiteral == "false" || parsedLiteral == "null" ***REMOVED***
						tkn = token.KEYWORD
					***REMOVED*** else ***REMOVED***
						tkn = token.IDENTIFIER
					***REMOVED***
					return
				***REMOVED***
				switch tkn ***REMOVED***

				case 0: // Not a keyword
					if parsedLiteral == "true" || parsedLiteral == "false" ***REMOVED***
						self.insertSemicolon = true
						tkn = token.BOOLEAN
						return
					***REMOVED*** else if parsedLiteral == "null" ***REMOVED***
						self.insertSemicolon = true
						tkn = token.NULL
						return
					***REMOVED***

				case token.KEYWORD:
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
					if self.chr == '.' ***REMOVED***
						self.read()
						if self.chr == '.' ***REMOVED***
							self.read()
							tkn = token.ELLIPSIS
						***REMOVED*** else ***REMOVED***
							tkn = token.ILLEGAL
						***REMOVED***
					***REMOVED*** else ***REMOVED***
						tkn = token.PERIOD
					***REMOVED***
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
				if self.chr == '*' ***REMOVED***
					self.read()
					tkn = self.switch2(token.EXPONENT, token.EXPONENT_ASSIGN)
				***REMOVED*** else ***REMOVED***
					tkn = self.switch2(token.MULTIPLY, token.MULTIPLY_ASSIGN)
				***REMOVED***
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
				if self.chr == '>' ***REMOVED***
					self.read()
					if self.implicitSemicolon ***REMOVED***
						tkn = token.ILLEGAL
					***REMOVED*** else ***REMOVED***
						tkn = token.ARROW
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					tkn = self.switch2(token.ASSIGN, token.EQUAL)
					if tkn == token.EQUAL && self.chr == '=' ***REMOVED***
						self.read()
						tkn = token.STRICT_EQUAL
					***REMOVED***
				***REMOVED***
			case '!':
				tkn = self.switch2(token.NOT, token.NOT_EQUAL)
				if tkn == token.NOT_EQUAL && self.chr == '=' ***REMOVED***
					self.read()
					tkn = token.STRICT_NOT_EQUAL
				***REMOVED***
			case '&':
				tkn = self.switch3(token.AND, token.AND_ASSIGN, '&', token.LOGICAL_AND)
			case '|':
				tkn = self.switch3(token.OR, token.OR_ASSIGN, '|', token.LOGICAL_OR)
			case '~':
				tkn = token.BITWISE_NOT
			case '?':
				if self.chr == '.' && !isDecimalDigit(self._peek()) ***REMOVED***
					self.read()
					tkn = token.QUESTION_DOT
				***REMOVED*** else if self.chr == '?' ***REMOVED***
					self.read()
					tkn = token.COALESCE
				***REMOVED*** else ***REMOVED***
					tkn = token.QUESTION_MARK
				***REMOVED***
			case '"', '\'':
				insertSemicolon = true
				tkn = token.STRING
				var err string
				literal, parsedLiteral, err = self.scanString(self.chrOffset-1, true)
				if err != "" ***REMOVED***
					tkn = token.ILLEGAL
				***REMOVED***
			case '`':
				tkn = token.BACKTICK
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

func (self *_parser) _peek() rune ***REMOVED***
	if self.offset < self.length ***REMOVED***
		return rune(self.str[self.offset])
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

func (self *_parser) scanMantissa(base int) ***REMOVED***
	for digitValue(self.chr) < base ***REMOVED***
		self.read()
	***REMOVED***
***REMOVED***

func (self *_parser) scanEscape(quote rune) (int, bool) ***REMOVED***

	var length, base uint32
	chr := self.chr
	switch chr ***REMOVED***
	case '0', '1', '2', '3', '4', '5', '6', '7':
		//    Octal:
		length, base = 3, 8
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"', '\'':
		self.read()
		return 1, false
	case '\r':
		self.read()
		if self.chr == '\n' ***REMOVED***
			self.read()
			return 2, false
		***REMOVED***
		return 1, false
	case '\n':
		self.read()
		return 1, false
	case '\u2028', '\u2029':
		self.read()
		return 1, true
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
	default:
		self.read() // Always make progress
	***REMOVED***

	if base > 0 ***REMOVED***
		var value uint32
		if length > 0 ***REMOVED***
			for ; length > 0 && self.chr != quote && self.chr >= 0; length-- ***REMOVED***
				digit := uint32(digitValue(self.chr))
				if digit >= base ***REMOVED***
					break
				***REMOVED***
				value = value*base + digit
				self.read()
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			for self.chr != quote && self.chr >= 0 && value < utf8.MaxRune ***REMOVED***
				if self.chr == '***REMOVED***' ***REMOVED***
					self.read()
					break
				***REMOVED***
				digit := uint32(digitValue(self.chr))
				if digit >= base ***REMOVED***
					break
				***REMOVED***
				value = value*base + digit
				self.read()
			***REMOVED***
		***REMOVED***
		chr = rune(value)
	***REMOVED***
	if chr >= utf8.RuneSelf ***REMOVED***
		if chr > 0xFFFF ***REMOVED***
			return 2, true
		***REMOVED***
		return 1, true
	***REMOVED***
	return 1, false
***REMOVED***

func (self *_parser) scanString(offset int, parse bool) (literal string, parsed unistring.String, err string) ***REMOVED***
	// " ' /
	quote := rune(self.str[offset])
	length := 0
	isUnicode := false
	for self.chr != quote ***REMOVED***
		chr := self.chr
		if chr == '\n' || chr == '\r' || chr == '\u2028' || chr == '\u2029' || chr < 0 ***REMOVED***
			goto newline
		***REMOVED***
		self.read()
		if chr == '\\' ***REMOVED***
			if self.chr == '\n' || self.chr == '\r' || self.chr == '\u2028' || self.chr == '\u2029' || self.chr < 0 ***REMOVED***
				if quote == '/' ***REMOVED***
					goto newline
				***REMOVED***
				self.scanNewline()
			***REMOVED*** else ***REMOVED***
				l, u := self.scanEscape(quote)
				length += l
				if u ***REMOVED***
					isUnicode = true
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED*** else if chr == '[' && quote == '/' ***REMOVED***
			// Allow a slash (/) in a bracket character class ([...])
			// TODO Fix this, this is hacky...
			quote = -1
		***REMOVED*** else if chr == ']' && quote == -1 ***REMOVED***
			quote = '/'
		***REMOVED***
		if chr >= utf8.RuneSelf ***REMOVED***
			isUnicode = true
			if chr > 0xFFFF ***REMOVED***
				length++
			***REMOVED***
		***REMOVED***
		length++
	***REMOVED***

	// " ' /
	self.read()
	literal = self.str[offset:self.chrOffset]
	if parse ***REMOVED***
		// TODO strict
		parsed, err = parseStringLiteral(literal[1:len(literal)-1], length, isUnicode, false)
	***REMOVED***
	return

newline:
	self.scanNewline()
	errStr := "String not terminated"
	if quote == '/' ***REMOVED***
		errStr = "Invalid regular expression: missing /"
		self.error(self.idxOf(offset), errStr)
	***REMOVED***
	return "", "", errStr
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

func (self *_parser) parseTemplateCharacters() (literal string, parsed unistring.String, finished bool, parseErr, err string) ***REMOVED***
	offset := self.chrOffset
	var end int
	length := 0
	isUnicode := false
	hasCR := false
	for ***REMOVED***
		chr := self.chr
		if chr < 0 ***REMOVED***
			goto unterminated
		***REMOVED***
		self.read()
		if chr == '`' ***REMOVED***
			finished = true
			end = self.chrOffset - 1
			break
		***REMOVED***
		if chr == '\\' ***REMOVED***
			if self.chr == '\n' || self.chr == '\r' || self.chr == '\u2028' || self.chr == '\u2029' || self.chr < 0 ***REMOVED***
				if self.chr == '\r' ***REMOVED***
					hasCR = true
				***REMOVED***
				self.scanNewline()
			***REMOVED*** else ***REMOVED***
				if self.chr == '8' || self.chr == '9' ***REMOVED***
					if parseErr == "" ***REMOVED***
						parseErr = "\\8 and \\9 are not allowed in template strings."
					***REMOVED***
				***REMOVED***
				l, u := self.scanEscape('`')
				length += l
				if u ***REMOVED***
					isUnicode = true
				***REMOVED***
			***REMOVED***
			continue
		***REMOVED***
		if chr == '$' && self.chr == '***REMOVED***' ***REMOVED***
			self.read()
			end = self.chrOffset - 2
			break
		***REMOVED***
		if chr >= utf8.RuneSelf ***REMOVED***
			isUnicode = true
			if chr > 0xFFFF ***REMOVED***
				length++
			***REMOVED***
		***REMOVED*** else if chr == '\r' ***REMOVED***
			hasCR = true
			if self.chr == '\n' ***REMOVED***
				length--
			***REMOVED***
		***REMOVED***
		length++
	***REMOVED***
	literal = self.str[offset:end]
	if hasCR ***REMOVED***
		literal = normaliseCRLF(literal)
	***REMOVED***
	if parseErr == "" ***REMOVED***
		parsed, parseErr = parseStringLiteral(literal, length, isUnicode, true)
	***REMOVED***
	self.insertSemicolon = true
	return
unterminated:
	err = err_UnexpectedEndOfInput
	finished = true
	return
***REMOVED***

func normaliseCRLF(s string) string ***REMOVED***
	var buf strings.Builder
	buf.Grow(len(s))
	for i := 0; i < len(s); i++ ***REMOVED***
		if s[i] == '\r' ***REMOVED***
			buf.WriteByte('\n')
			if i < len(s)-1 && s[i+1] == '\n' ***REMOVED***
				i++
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			buf.WriteByte(s[i])
		***REMOVED***
	***REMOVED***
	return buf.String()
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

func parseStringLiteral(literal string, length int, unicode, strict bool) (unistring.String, string) ***REMOVED***
	var sb strings.Builder
	var chars []uint16
	if unicode ***REMOVED***
		chars = make([]uint16, 1, length+1)
		chars[0] = unistring.BOM
	***REMOVED*** else ***REMOVED***
		sb.Grow(length)
	***REMOVED***
	str := literal
	for len(str) > 0 ***REMOVED***
		switch chr := str[0]; ***REMOVED***
		// We do not explicitly handle the case of the quote
		// value, which can be: " ' /
		// This assumes we're already passed a partially well-formed literal
		case chr >= utf8.RuneSelf:
			chr, size := utf8.DecodeRuneInString(str)
			if chr <= 0xFFFF ***REMOVED***
				chars = append(chars, uint16(chr))
			***REMOVED*** else ***REMOVED***
				first, second := utf16.EncodeRune(chr)
				chars = append(chars, uint16(first), uint16(second))
			***REMOVED***
			str = str[size:]
			continue
		case chr != '\\':
			if unicode ***REMOVED***
				chars = append(chars, uint16(chr))
			***REMOVED*** else ***REMOVED***
				sb.WriteByte(chr)
			***REMOVED***
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
			if value == '\u2028' || value == '\u2029' ***REMOVED***
				continue
			***REMOVED***
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
					if str == "" || str[0] != '***REMOVED***' ***REMOVED***
						size = 4
					***REMOVED***
				***REMOVED***
				if size > 0 ***REMOVED***
					if len(str) < size ***REMOVED***
						return "", fmt.Sprintf("invalid escape: \\%s: len(%q) != %d", string(chr), str, size)
					***REMOVED***
					for j := 0; j < size; j++ ***REMOVED***
						decimal, ok := hex2decimal(str[j])
						if !ok ***REMOVED***
							return "", fmt.Sprintf("invalid escape: \\%s: %q", string(chr), str[:size])
						***REMOVED***
						value = value<<4 | decimal
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					str = str[1:]
					var val rune
					value = -1
					for ; size < len(str); size++ ***REMOVED***
						if str[size] == '***REMOVED***' ***REMOVED***
							if size == 0 ***REMOVED***
								return "", fmt.Sprintf("invalid escape: \\%s", string(chr))
							***REMOVED***
							size++
							value = val
							break
						***REMOVED***
						decimal, ok := hex2decimal(str[size])
						if !ok ***REMOVED***
							return "", fmt.Sprintf("invalid escape: \\%s: %q", string(chr), str[:size+1])
						***REMOVED***
						val = val<<4 | decimal
						if val > utf8.MaxRune ***REMOVED***
							return "", fmt.Sprintf("undefined Unicode code-point: %q", str[:size+1])
						***REMOVED***
					***REMOVED***
					if value == -1 ***REMOVED***
						return "", fmt.Sprintf("unterminated \\u***REMOVED***: %q", str)
					***REMOVED***
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
				if strict ***REMOVED***
					return "", "Octal escape sequences are not allowed in this context"
				***REMOVED***
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
		***REMOVED***
		if unicode ***REMOVED***
			if value <= 0xFFFF ***REMOVED***
				chars = append(chars, uint16(value))
			***REMOVED*** else ***REMOVED***
				first, second := utf16.EncodeRune(value)
				chars = append(chars, uint16(first), uint16(second))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if value >= utf8.RuneSelf ***REMOVED***
				return "", "Unexpected unicode character"
			***REMOVED***
			sb.WriteByte(byte(value))
		***REMOVED***
	***REMOVED***

	if unicode ***REMOVED***
		if len(chars) != length+1 ***REMOVED***
			panic(fmt.Errorf("unexpected unicode length while parsing '%s'", literal))
		***REMOVED***
		return unistring.FromUtf16(chars), ""
	***REMOVED***
	if sb.Len() != length ***REMOVED***
		panic(fmt.Errorf("unexpected length while parsing '%s'", literal))
	***REMOVED***
	return unistring.String(sb.String()), ""
***REMOVED***

func (self *_parser) scanNumericLiteral(decimalPoint bool) (token.Token, string) ***REMOVED***

	offset := self.chrOffset
	tkn := token.NUMBER

	if decimalPoint ***REMOVED***
		offset--
		self.scanMantissa(10)
	***REMOVED*** else ***REMOVED***
		if self.chr == '0' ***REMOVED***
			self.read()
			base := 0
			switch self.chr ***REMOVED***
			case 'x', 'X':
				base = 16
			case 'o', 'O':
				base = 8
			case 'b', 'B':
				base = 2
			case '.', 'e', 'E':
				// no-op
			default:
				// legacy octal
				self.scanMantissa(8)
				goto end
			***REMOVED***
			if base > 0 ***REMOVED***
				self.read()
				if !isDigit(self.chr, base) ***REMOVED***
					return token.ILLEGAL, self.str[offset:self.chrOffset]
				***REMOVED***
				self.scanMantissa(base)
				goto end
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			self.scanMantissa(10)
		***REMOVED***
		if self.chr == '.' ***REMOVED***
			self.read()
			self.scanMantissa(10)
		***REMOVED***
	***REMOVED***

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
end:
	if isIdentifierStart(self.chr) || isDecimalDigit(self.chr) ***REMOVED***
		return token.ILLEGAL, self.str[offset:self.chrOffset]
	***REMOVED***

	return tkn, self.str[offset:self.chrOffset]
***REMOVED***
