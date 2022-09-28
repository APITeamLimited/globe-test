// Package token defines constants representing the lexical tokens of JavaScript (ECMA5).
package token

import (
	"strconv"
)

// Token is the set of lexical tokens in JavaScript (ECMA5).
type Token int

// String returns the string corresponding to the token.
// For operators, delimiters, and keywords the string is the actual
// token string (e.g., for the token PLUS, the String() is
// "+"). For all other tokens the string corresponds to the token
// name (e.g. for the token IDENTIFIER, the string is "IDENTIFIER").
//
func (tkn Token) String() string ***REMOVED***
	if tkn == 0 ***REMOVED***
		return "UNKNOWN"
	***REMOVED***
	if tkn < Token(len(token2string)) ***REMOVED***
		return token2string[tkn]
	***REMOVED***
	return "token(" + strconv.Itoa(int(tkn)) + ")"
***REMOVED***

//lint:ignore U1000 This is not used for anything
func (tkn Token) precedence(in bool) int ***REMOVED***

	switch tkn ***REMOVED***
	case LOGICAL_OR:
		return 1

	case LOGICAL_AND:
		return 2

	case OR, OR_ASSIGN:
		return 3

	case EXCLUSIVE_OR:
		return 4

	case AND, AND_ASSIGN:
		return 5

	case EQUAL,
		NOT_EQUAL,
		STRICT_EQUAL,
		STRICT_NOT_EQUAL:
		return 6

	case LESS, GREATER, LESS_OR_EQUAL, GREATER_OR_EQUAL, INSTANCEOF:
		return 7

	case IN:
		if in ***REMOVED***
			return 7
		***REMOVED***
		return 0

	case SHIFT_LEFT, SHIFT_RIGHT, UNSIGNED_SHIFT_RIGHT:
		fallthrough
	case SHIFT_LEFT_ASSIGN, SHIFT_RIGHT_ASSIGN, UNSIGNED_SHIFT_RIGHT_ASSIGN:
		return 8

	case PLUS, MINUS, ADD_ASSIGN, SUBTRACT_ASSIGN:
		return 9

	case MULTIPLY, SLASH, REMAINDER, MULTIPLY_ASSIGN, QUOTIENT_ASSIGN, REMAINDER_ASSIGN:
		return 11
	***REMOVED***
	return 0
***REMOVED***

type _keyword struct ***REMOVED***
	token         Token
	futureKeyword bool
	strict        bool
***REMOVED***

// IsKeyword returns the keyword token if literal is a keyword, a KEYWORD token
// if the literal is a future keyword (const, let, class, super, ...), or 0 if the literal is not a keyword.
//
// If the literal is a keyword, IsKeyword returns a second value indicating if the literal
// is considered a future keyword in strict-mode only.
//
// 7.6.1.2 Future Reserved Words:
//
//       const
//       class
//       enum
//       export
//       extends
//       import
//       super
//
// 7.6.1.2 Future Reserved Words (strict):
//
//       implements
//       interface
//       let
//       package
//       private
//       protected
//       public
//       static
//
func IsKeyword(literal string) (Token, bool) ***REMOVED***
	if keyword, exists := keywordTable[literal]; exists ***REMOVED***
		if keyword.futureKeyword ***REMOVED***
			return KEYWORD, keyword.strict
		***REMOVED***
		return keyword.token, false
	***REMOVED***
	return 0, false
***REMOVED***

func IsId(tkn Token) bool ***REMOVED***
	return tkn >= IDENTIFIER
***REMOVED***

func IsUnreservedWord(tkn Token) bool ***REMOVED***
	return tkn > ESCAPED_RESERVED_WORD
***REMOVED***
