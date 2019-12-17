/*
Package parser implements a parser for JavaScript.

    import (
        "github.com/dop251/goja/parser"
    )

Parse and return an AST

    filename := "" // A filename is optional
    src := `
        // Sample xyzzy example
        (function()***REMOVED***
            if (3.14159 > 0) ***REMOVED***
                console.log("Hello, World.");
                return;
            ***REMOVED***

            var xyzzy = NaN;
            console.log("Nothing happens.");
            return xyzzy;
        ***REMOVED***)();
    `

    // Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
    program, err := parser.ParseFile(nil, filename, src, 0)

Warning

The parser and AST interfaces are still works-in-progress (particularly where
node types are concerned) and may change in the future.

*/
package parser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/token"
)

// A Mode value is a set of flags (or 0). They control optional parser functionality.
type Mode uint

const (
	IgnoreRegExpErrors Mode = 1 << iota // Ignore RegExp compatibility errors (allow backtracking)
)

type _parser struct ***REMOVED***
	str    string
	length int
	base   int

	chr       rune // The current character
	chrOffset int  // The offset of current character
	offset    int  // The offset after current character (may be greater than 1)

	idx     file.Idx    // The index of token
	token   token.Token // The token
	literal string      // The literal of the token, if any

	scope             *_scope
	insertSemicolon   bool // If we see a newline, then insert an implicit semicolon
	implicitSemicolon bool // An implicit semicolon exists

	errors ErrorList

	recover struct ***REMOVED***
		// Scratch when trying to seek to the next statement, etc.
		idx   file.Idx
		count int
	***REMOVED***

	mode Mode

	file *file.File
***REMOVED***

func _newParser(filename, src string, base int) *_parser ***REMOVED***
	return &_parser***REMOVED***
		chr:    ' ', // This is set so we can start scanning by skipping whitespace
		str:    src,
		length: len(src),
		base:   base,
		file:   file.NewFile(filename, src, base),
	***REMOVED***
***REMOVED***

func newParser(filename, src string) *_parser ***REMOVED***
	return _newParser(filename, src, 1)
***REMOVED***

func ReadSource(filename string, src interface***REMOVED******REMOVED***) ([]byte, error) ***REMOVED***
	if src != nil ***REMOVED***
		switch src := src.(type) ***REMOVED***
		case string:
			return []byte(src), nil
		case []byte:
			return src, nil
		case *bytes.Buffer:
			if src != nil ***REMOVED***
				return src.Bytes(), nil
			***REMOVED***
		case io.Reader:
			var bfr bytes.Buffer
			if _, err := io.Copy(&bfr, src); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			return bfr.Bytes(), nil
		***REMOVED***
		return nil, errors.New("invalid source")
	***REMOVED***
	return ioutil.ReadFile(filename)
***REMOVED***

// ParseFile parses the source code of a single JavaScript/ECMAScript source file and returns
// the corresponding ast.Program node.
//
// If fileSet == nil, ParseFile parses source without a FileSet.
// If fileSet != nil, ParseFile first adds filename and src to fileSet.
//
// The filename argument is optional and is used for labelling errors, etc.
//
// src may be a string, a byte slice, a bytes.Buffer, or an io.Reader, but it MUST always be in UTF-8.
//
//      // Parse some JavaScript, yielding a *ast.Program and/or an ErrorList
//      program, err := parser.ParseFile(nil, "", `if (abc > 1) ***REMOVED******REMOVED***`, 0)
//
func ParseFile(fileSet *file.FileSet, filename string, src interface***REMOVED******REMOVED***, mode Mode) (*ast.Program, error) ***REMOVED***
	str, err := ReadSource(filename, src)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	***REMOVED***
		str := string(str)

		base := 1
		if fileSet != nil ***REMOVED***
			base = fileSet.AddFile(filename, str)
		***REMOVED***

		parser := _newParser(filename, str, base)
		parser.mode = mode
		return parser.parse()
	***REMOVED***
***REMOVED***

// ParseFunction parses a given parameter list and body as a function and returns the
// corresponding ast.FunctionLiteral node.
//
// The parameter list, if any, should be a comma-separated list of identifiers.
//
func ParseFunction(parameterList, body string) (*ast.FunctionLiteral, error) ***REMOVED***

	src := "(function(" + parameterList + ") ***REMOVED***\n" + body + "\n***REMOVED***)"

	parser := _newParser("", src, 1)
	program, err := parser.parse()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return program.Body[0].(*ast.ExpressionStatement).Expression.(*ast.FunctionLiteral), nil
***REMOVED***

func (self *_parser) slice(idx0, idx1 file.Idx) string ***REMOVED***
	from := int(idx0) - self.base
	to := int(idx1) - self.base
	if from >= 0 && to <= len(self.str) ***REMOVED***
		return self.str[from:to]
	***REMOVED***

	return ""
***REMOVED***

func (self *_parser) parse() (*ast.Program, error) ***REMOVED***
	self.next()
	program := self.parseProgram()
	if false ***REMOVED***
		self.errors.Sort()
	***REMOVED***
	return program, self.errors.Err()
***REMOVED***

func (self *_parser) next() ***REMOVED***
	self.token, self.literal, self.idx = self.scan()
***REMOVED***

func (self *_parser) optionalSemicolon() ***REMOVED***
	if self.token == token.SEMICOLON ***REMOVED***
		self.next()
		return
	***REMOVED***

	if self.implicitSemicolon ***REMOVED***
		self.implicitSemicolon = false
		return
	***REMOVED***

	if self.token != token.EOF && self.token != token.RIGHT_BRACE ***REMOVED***
		self.expect(token.SEMICOLON)
	***REMOVED***
***REMOVED***

func (self *_parser) semicolon() ***REMOVED***
	if self.token != token.RIGHT_PARENTHESIS && self.token != token.RIGHT_BRACE ***REMOVED***
		if self.implicitSemicolon ***REMOVED***
			self.implicitSemicolon = false
			return
		***REMOVED***

		self.expect(token.SEMICOLON)
	***REMOVED***
***REMOVED***

func (self *_parser) idxOf(offset int) file.Idx ***REMOVED***
	return file.Idx(self.base + offset)
***REMOVED***

func (self *_parser) expect(value token.Token) file.Idx ***REMOVED***
	idx := self.idx
	if self.token != value ***REMOVED***
		self.errorUnexpectedToken(self.token)
	***REMOVED***
	self.next()
	return idx
***REMOVED***

func lineCount(str string) (int, int) ***REMOVED***
	line, last := 0, -1
	pair := false
	for index, chr := range str ***REMOVED***
		switch chr ***REMOVED***
		case '\r':
			line += 1
			last = index
			pair = true
			continue
		case '\n':
			if !pair ***REMOVED***
				line += 1
			***REMOVED***
			last = index
		case '\u2028', '\u2029':
			line += 1
			last = index + 2
		***REMOVED***
		pair = false
	***REMOVED***
	return line, last
***REMOVED***

func (self *_parser) position(idx file.Idx) file.Position ***REMOVED***
	position := file.Position***REMOVED******REMOVED***
	offset := int(idx) - self.base
	str := self.str[:offset]
	position.Filename = self.file.Name()
	line, last := lineCount(str)
	position.Line = 1 + line
	if last >= 0 ***REMOVED***
		position.Column = offset - last
	***REMOVED*** else ***REMOVED***
		position.Column = 1 + len(str)
	***REMOVED***

	return position
***REMOVED***
