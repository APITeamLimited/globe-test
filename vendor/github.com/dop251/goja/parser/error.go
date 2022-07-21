package parser

import (
	"fmt"
	"sort"

	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
)

const (
	err_UnexpectedToken      = "Unexpected token %v"
	err_UnexpectedEndOfInput = "Unexpected end of input"
	err_UnexpectedEscape     = "Unexpected escape"
)

//    UnexpectedNumber:  'Unexpected number',
//    UnexpectedString:  'Unexpected string',
//    UnexpectedIdentifier:  'Unexpected identifier',
//    UnexpectedReserved:  'Unexpected reserved word',
//    NewlineAfterThrow:  'Illegal newline after throw',
//    InvalidRegExp: 'Invalid regular expression',
//    UnterminatedRegExp:  'Invalid regular expression: missing /',
//    InvalidLHSInAssignment:  'Invalid left-hand side in assignment',
//    InvalidLHSInForIn:  'Invalid left-hand side in for-in',
//    MultipleDefaultsInSwitch: 'More than one default clause in switch statement',
//    NoCatchOrFinally:  'Missing catch or finally after try',
//    UnknownLabel: 'Undefined label \'%0\'',
//    Redeclaration: '%0 \'%1\' has already been declared',
//    IllegalContinue: 'Illegal continue statement',
//    IllegalBreak: 'Illegal break statement',
//    IllegalReturn: 'Illegal return statement',
//    StrictModeWith:  'Strict mode code may not include a with statement',
//    StrictCatchVariable:  'Catch variable may not be eval or arguments in strict mode',
//    StrictVarName:  'Variable name may not be eval or arguments in strict mode',
//    StrictParamName:  'Parameter name eval or arguments is not allowed in strict mode',
//    StrictParamDupe: 'Strict mode function may not have duplicate parameter names',
//    StrictFunctionName:  'Function name may not be eval or arguments in strict mode',
//    StrictOctalLiteral:  'Octal literals are not allowed in strict mode.',
//    StrictDelete:  'Delete of an unqualified identifier in strict mode.',
//    StrictDuplicateProperty:  'Duplicate data property in object literal not allowed in strict mode',
//    AccessorDataProperty:  'Object literal may not have data and accessor property with the same name',
//    AccessorGetSet:  'Object literal may not have multiple get/set accessors with the same name',
//    StrictLHSAssignment:  'Assignment to eval or arguments is not allowed in strict mode',
//    StrictLHSPostfix:  'Postfix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictLHSPrefix:  'Prefix increment/decrement may not have eval or arguments operand in strict mode',
//    StrictReservedWord:  'Use of future reserved word in strict mode'

// A SyntaxError is a description of an ECMAScript syntax error.

// An Error represents a parsing error. It includes the position where the error occurred and a message/description.
type Error struct ***REMOVED***
	Position file.Position
	Message  string
***REMOVED***

// FIXME Should this be "SyntaxError"?

func (self Error) Error() string ***REMOVED***
	filename := self.Position.Filename
	if filename == "" ***REMOVED***
		filename = "(anonymous)"
	***REMOVED***
	return fmt.Sprintf("%s: Line %d:%d %s",
		filename,
		self.Position.Line,
		self.Position.Column,
		self.Message,
	)
***REMOVED***

func (self *_parser) error(place interface***REMOVED******REMOVED***, msg string, msgValues ...interface***REMOVED******REMOVED***) *Error ***REMOVED***
	idx := file.Idx(0)
	switch place := place.(type) ***REMOVED***
	case int:
		idx = self.idxOf(place)
	case file.Idx:
		if place == 0 ***REMOVED***
			idx = self.idxOf(self.chrOffset)
		***REMOVED*** else ***REMOVED***
			idx = place
		***REMOVED***
	default:
		panic(fmt.Errorf("error(%T, ...)", place))
	***REMOVED***

	position := self.position(idx)
	msg = fmt.Sprintf(msg, msgValues...)
	self.errors.Add(position, msg)
	return self.errors[len(self.errors)-1]
***REMOVED***

func (self *_parser) errorUnexpected(idx file.Idx, chr rune) error ***REMOVED***
	if chr == -1 ***REMOVED***
		return self.error(idx, err_UnexpectedEndOfInput)
	***REMOVED***
	return self.error(idx, err_UnexpectedToken, token.ILLEGAL)
***REMOVED***

func (self *_parser) errorUnexpectedToken(tkn token.Token) error ***REMOVED***
	switch tkn ***REMOVED***
	case token.EOF:
		return self.error(file.Idx(0), err_UnexpectedEndOfInput)
	***REMOVED***
	value := tkn.String()
	switch tkn ***REMOVED***
	case token.BOOLEAN, token.NULL:
		value = self.literal
	case token.IDENTIFIER:
		return self.error(self.idx, "Unexpected identifier")
	case token.KEYWORD:
		// TODO Might be a future reserved word
		return self.error(self.idx, "Unexpected reserved word")
	case token.ESCAPED_RESERVED_WORD:
		return self.error(self.idx, "Keyword must not contain escaped characters")
	case token.NUMBER:
		return self.error(self.idx, "Unexpected number")
	case token.STRING:
		return self.error(self.idx, "Unexpected string")
	***REMOVED***
	return self.error(self.idx, err_UnexpectedToken, value)
***REMOVED***

// ErrorList is a list of *Errors.
//
type ErrorList []*Error

// Add adds an Error with given position and message to an ErrorList.
func (self *ErrorList) Add(position file.Position, msg string) ***REMOVED***
	*self = append(*self, &Error***REMOVED***position, msg***REMOVED***)
***REMOVED***

// Reset resets an ErrorList to no errors.
func (self *ErrorList) Reset() ***REMOVED*** *self = (*self)[0:0] ***REMOVED***

func (self ErrorList) Len() int      ***REMOVED*** return len(self) ***REMOVED***
func (self ErrorList) Swap(i, j int) ***REMOVED*** self[i], self[j] = self[j], self[i] ***REMOVED***
func (self ErrorList) Less(i, j int) bool ***REMOVED***
	x := &self[i].Position
	y := &self[j].Position
	if x.Filename < y.Filename ***REMOVED***
		return true
	***REMOVED***
	if x.Filename == y.Filename ***REMOVED***
		if x.Line < y.Line ***REMOVED***
			return true
		***REMOVED***
		if x.Line == y.Line ***REMOVED***
			return x.Column < y.Column
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func (self ErrorList) Sort() ***REMOVED***
	sort.Sort(self)
***REMOVED***

// Error implements the Error interface.
func (self ErrorList) Error() string ***REMOVED***
	switch len(self) ***REMOVED***
	case 0:
		return "no errors"
	case 1:
		return self[0].Error()
	***REMOVED***
	return fmt.Sprintf("%s (and %d more errors)", self[0].Error(), len(self)-1)
***REMOVED***

// Err returns an error equivalent to this ErrorList.
// If the list is empty, Err returns nil.
func (self ErrorList) Err() error ***REMOVED***
	if len(self) == 0 ***REMOVED***
		return nil
	***REMOVED***
	return self
***REMOVED***
