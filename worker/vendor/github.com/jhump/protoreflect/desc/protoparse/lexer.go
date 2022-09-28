package protoparse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

type runeReader struct ***REMOVED***
	rr     *bufio.Reader
	marked []rune
	unread []rune
	err    error
***REMOVED***

func (rr *runeReader) readRune() (r rune, size int, err error) ***REMOVED***
	if rr.err != nil ***REMOVED***
		return 0, 0, rr.err
	***REMOVED***
	if len(rr.unread) > 0 ***REMOVED***
		r := rr.unread[len(rr.unread)-1]
		rr.unread = rr.unread[:len(rr.unread)-1]
		if rr.marked != nil ***REMOVED***
			rr.marked = append(rr.marked, r)
		***REMOVED***
		return r, utf8.RuneLen(r), nil
	***REMOVED***
	r, sz, err := rr.rr.ReadRune()
	if err != nil ***REMOVED***
		rr.err = err
	***REMOVED*** else if rr.marked != nil ***REMOVED***
		rr.marked = append(rr.marked, r)
	***REMOVED***
	return r, sz, err
***REMOVED***

func (rr *runeReader) unreadRune(r rune) ***REMOVED***
	if rr.marked != nil ***REMOVED***
		if rr.marked[len(rr.marked)-1] != r ***REMOVED***
			panic("unread rune is not the same as last marked rune!")
		***REMOVED***
		rr.marked = rr.marked[:len(rr.marked)-1]
	***REMOVED***
	rr.unread = append(rr.unread, r)
***REMOVED***

func (rr *runeReader) startMark(initial rune) ***REMOVED***
	rr.marked = []rune***REMOVED***initial***REMOVED***
***REMOVED***

func (rr *runeReader) endMark() string ***REMOVED***
	m := string(rr.marked)
	rr.marked = rr.marked[:0]
	return m
***REMOVED***

type protoLex struct ***REMOVED***
	filename string
	input    *runeReader
	errs     *errorHandler
	res      *ast.FileNode

	lineNo int
	colNo  int
	offset int

	prevSym ast.TerminalNode
	eof     ast.TerminalNode

	prevLineNo int
	prevColNo  int
	prevOffset int
	comments   []ast.Comment
	ws         []rune
***REMOVED***

var utf8Bom = []byte***REMOVED***0xEF, 0xBB, 0xBF***REMOVED***

func newLexer(in io.Reader, filename string, errs *errorHandler) *protoLex ***REMOVED***
	br := bufio.NewReader(in)

	// if file has UTF8 byte order marker preface, consume it
	marker, err := br.Peek(3)
	if err == nil && bytes.Equal(marker, utf8Bom) ***REMOVED***
		_, _ = br.Discard(3)
	***REMOVED***

	return &protoLex***REMOVED***
		input:    &runeReader***REMOVED***rr: br***REMOVED***,
		filename: filename,
		errs:     errs,
	***REMOVED***
***REMOVED***

var keywords = map[string]int***REMOVED***
	"syntax":     _SYNTAX,
	"import":     _IMPORT,
	"weak":       _WEAK,
	"public":     _PUBLIC,
	"package":    _PACKAGE,
	"option":     _OPTION,
	"true":       _TRUE,
	"false":      _FALSE,
	"inf":        _INF,
	"nan":        _NAN,
	"repeated":   _REPEATED,
	"optional":   _OPTIONAL,
	"required":   _REQUIRED,
	"double":     _DOUBLE,
	"float":      _FLOAT,
	"int32":      _INT32,
	"int64":      _INT64,
	"uint32":     _UINT32,
	"uint64":     _UINT64,
	"sint32":     _SINT32,
	"sint64":     _SINT64,
	"fixed32":    _FIXED32,
	"fixed64":    _FIXED64,
	"sfixed32":   _SFIXED32,
	"sfixed64":   _SFIXED64,
	"bool":       _BOOL,
	"string":     _STRING,
	"bytes":      _BYTES,
	"group":      _GROUP,
	"oneof":      _ONEOF,
	"map":        _MAP,
	"extensions": _EXTENSIONS,
	"to":         _TO,
	"max":        _MAX,
	"reserved":   _RESERVED,
	"enum":       _ENUM,
	"message":    _MESSAGE,
	"extend":     _EXTEND,
	"service":    _SERVICE,
	"rpc":        _RPC,
	"stream":     _STREAM,
	"returns":    _RETURNS,
***REMOVED***

func (l *protoLex) cur() SourcePos ***REMOVED***
	return SourcePos***REMOVED***
		Filename: l.filename,
		Offset:   l.offset,
		Line:     l.lineNo + 1,
		Col:      l.colNo + 1,
	***REMOVED***
***REMOVED***

func (l *protoLex) adjustPos(consumedChars ...rune) ***REMOVED***
	for _, c := range consumedChars ***REMOVED***
		switch c ***REMOVED***
		case '\n':
			// new line, back to first column
			l.colNo = 0
			l.lineNo++
		case '\r':
			// no adjustment
		case '\t':
			// advance to next tab stop
			mod := l.colNo % 8
			l.colNo += 8 - mod
		default:
			l.colNo++
		***REMOVED***
	***REMOVED***
***REMOVED***

func (l *protoLex) prev() *SourcePos ***REMOVED***
	if l.prevSym == nil ***REMOVED***
		return &SourcePos***REMOVED***
			Filename: l.filename,
			Offset:   0,
			Line:     1,
			Col:      1,
		***REMOVED***
	***REMOVED***
	return l.prevSym.Start()
***REMOVED***

func (l *protoLex) Lex(lval *protoSymType) int ***REMOVED***
	if l.errs.err != nil ***REMOVED***
		// if error reporter already returned non-nil error,
		// we can skip the rest of the input
		return 0
	***REMOVED***

	l.prevLineNo = l.lineNo
	l.prevColNo = l.colNo
	l.prevOffset = l.offset
	l.comments = nil
	l.ws = nil
	l.input.endMark() // reset, just in case

	for ***REMOVED***
		c, n, err := l.input.readRune()
		if err == io.EOF ***REMOVED***
			// we're not actually returning a rune, but this will associate
			// accumulated comments as a trailing comment on last symbol
			// (if appropriate)
			l.setRune(lval, 0)
			l.eof = lval.b
			return 0
		***REMOVED*** else if err != nil ***REMOVED***
			// we don't call setError because we don't want it wrapped
			// with a source position because it's I/O, not syntax
			lval.err = err
			_ = l.errs.handleError(err)
			return _ERROR
		***REMOVED***

		l.prevLineNo = l.lineNo
		l.prevColNo = l.colNo
		l.prevOffset = l.offset

		l.offset += n
		l.adjustPos(c)
		if strings.ContainsRune("\n\r\t\f\v ", c) ***REMOVED***
			l.ws = append(l.ws, c)
			continue
		***REMOVED***

		l.input.startMark(c)
		if c == '.' ***REMOVED***
			// decimal literals could start with a dot
			cn, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				l.setRune(lval, c)
				return int(c)
			***REMOVED***
			if cn >= '0' && cn <= '9' ***REMOVED***
				l.adjustPos(cn)
				token := l.readNumber(c, cn)
				f, err := parseFloat(token)
				if err != nil ***REMOVED***
					l.setError(lval, numError(err, "float", token))
					return _ERROR
				***REMOVED***
				l.setFloat(lval, f)
				return _FLOAT_LIT
			***REMOVED***
			l.input.unreadRune(cn)
			l.setRune(lval, c)
			return int(c)
		***REMOVED***

		if c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') ***REMOVED***
			// identifier
			token := []rune***REMOVED***c***REMOVED***
			token = l.readIdentifier(token)
			str := string(token)
			if t, ok := keywords[str]; ok ***REMOVED***
				l.setIdent(lval, str)
				return t
			***REMOVED***
			l.setIdent(lval, str)
			return _NAME
		***REMOVED***

		if c >= '0' && c <= '9' ***REMOVED***
			// integer or float literal
			token := l.readNumber(c)
			if strings.HasPrefix(token, "0x") || strings.HasPrefix(token, "0X") ***REMOVED***
				// hexadecimal
				ui, err := strconv.ParseUint(token[2:], 16, 64)
				if err != nil ***REMOVED***
					l.setError(lval, numError(err, "hexadecimal integer", token[2:]))
					return _ERROR
				***REMOVED***
				l.setInt(lval, ui)
				return _INT_LIT
			***REMOVED***
			if strings.Contains(token, ".") || strings.Contains(token, "e") || strings.Contains(token, "E") ***REMOVED***
				// floating point!
				f, err := parseFloat(token)
				if err != nil ***REMOVED***
					l.setError(lval, numError(err, "float", token))
					return _ERROR
				***REMOVED***
				l.setFloat(lval, f)
				return _FLOAT_LIT
			***REMOVED***
			// integer! (decimal or octal)
			base := 10
			if token[0] == '0' ***REMOVED***
				base = 8
			***REMOVED***
			ui, err := strconv.ParseUint(token, base, 64)
			if err != nil ***REMOVED***
				kind := "integer"
				if base == 8 ***REMOVED***
					kind = "octal integer"
				***REMOVED***
				if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange ***REMOVED***
					// if it's too big to be an int, parse it as a float
					var f float64
					kind = "float"
					f, err = parseFloat(token)
					if err == nil ***REMOVED***
						l.setFloat(lval, f)
						return _FLOAT_LIT
					***REMOVED***
				***REMOVED***
				l.setError(lval, numError(err, kind, token))
				return _ERROR
			***REMOVED***
			l.setInt(lval, ui)
			return _INT_LIT
		***REMOVED***

		if c == '\'' || c == '"' ***REMOVED***
			// string literal
			str, err := l.readStringLiteral(c)
			if err != nil ***REMOVED***
				l.setError(lval, err)
				return _ERROR
			***REMOVED***
			l.setString(lval, str)
			return _STRING_LIT
		***REMOVED***

		if c == '/' ***REMOVED***
			// comment
			cn, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				l.setRune(lval, '/')
				return int(c)
			***REMOVED***
			if cn == '/' ***REMOVED***
				l.adjustPos(cn)
				hitNewline, hasErr := l.skipToEndOfLineComment(lval)
				if hasErr ***REMOVED***
					return _ERROR
				***REMOVED***
				comment := l.newComment()
				comment.PosRange.End.Col++
				if hitNewline ***REMOVED***
					// we don't do this inside of skipToEndOfLineComment
					// because we want to know the length of previous
					// line for calculation above
					l.adjustPos('\n')
				***REMOVED***
				l.comments = append(l.comments, comment)
				continue
			***REMOVED***
			if cn == '*' ***REMOVED***
				l.adjustPos(cn)
				ok, hasErr := l.skipToEndOfBlockComment(lval)
				if hasErr ***REMOVED***
					return _ERROR
				***REMOVED***
				if !ok ***REMOVED***
					l.setError(lval, errors.New("block comment never terminates, unexpected EOF"))
					return _ERROR
				***REMOVED***
				l.comments = append(l.comments, l.newComment())
				continue
			***REMOVED***
			l.input.unreadRune(cn)
		***REMOVED***

		if c < 32 || c == 127 ***REMOVED***
			l.setError(lval, errors.New("invalid control character"))
			return _ERROR
		***REMOVED***
		if !strings.ContainsRune(";,.:=-+()***REMOVED******REMOVED***[]<>/", c) ***REMOVED***
			l.setError(lval, errors.New("invalid character"))
			return _ERROR
		***REMOVED***
		l.setRune(lval, c)
		return int(c)
	***REMOVED***
***REMOVED***

func parseFloat(token string) (float64, error) ***REMOVED***
	// strconv.ParseFloat allows _ to separate digits, but protobuf does not
	if strings.ContainsRune(token, '_') ***REMOVED***
		return 0, &strconv.NumError***REMOVED***
			Func: "parseFloat",
			Num:  token,
			Err:  strconv.ErrSyntax,
		***REMOVED***
	***REMOVED***
	f, err := strconv.ParseFloat(token, 64)
	if err == nil ***REMOVED***
		return f, nil
	***REMOVED***
	if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange && math.IsInf(f, 1) ***REMOVED***
		// protoc doesn't complain about float overflow and instead just uses "infinity"
		// so we mirror that behavior by just returning infinity and ignoring the error
		return f, nil
	***REMOVED***
	return f, err
***REMOVED***

func (l *protoLex) posRange() ast.PosRange ***REMOVED***
	return ast.PosRange***REMOVED***
		Start: SourcePos***REMOVED***
			Filename: l.filename,
			Offset:   l.prevOffset,
			Line:     l.prevLineNo + 1,
			Col:      l.prevColNo + 1,
		***REMOVED***,
		End: l.cur(),
	***REMOVED***
***REMOVED***

func (l *protoLex) newComment() ast.Comment ***REMOVED***
	ws := string(l.ws)
	l.ws = l.ws[:0]
	return ast.Comment***REMOVED***
		PosRange:          l.posRange(),
		LeadingWhitespace: ws,
		Text:              l.input.endMark(),
	***REMOVED***
***REMOVED***

func (l *protoLex) newTokenInfo() ast.TokenInfo ***REMOVED***
	ws := string(l.ws)
	l.ws = nil
	return ast.TokenInfo***REMOVED***
		PosRange:          l.posRange(),
		LeadingComments:   l.comments,
		LeadingWhitespace: ws,
		RawText:           l.input.endMark(),
	***REMOVED***
***REMOVED***

func (l *protoLex) setPrev(n ast.TerminalNode, isDot bool) ***REMOVED***
	nStart := n.Start().Line
	if _, ok := n.(*ast.RuneNode); ok ***REMOVED***
		// This is really gross, but there are many cases where we don't want
		// to attribute comments to punctuation (like commas, equals, semicolons)
		// and would instead prefer to attribute comments to a more meaningful
		// element in the AST.
		//
		// So if it's a simple node OTHER THAN PERIOD (since that is not just
		// punctuation but typically part of a qualified identifier), don't
		// attribute comments to it. We do that with this TOTAL HACK: adjusting
		// the start line makes leading comments appear detached so logic below
		// will naturally associated trailing comment to previous symbol
		if !isDot ***REMOVED***
			nStart += 2
		***REMOVED***
	***REMOVED***
	if l.prevSym != nil && len(n.LeadingComments()) > 0 && l.prevSym.End().Line < nStart ***REMOVED***
		// we may need to re-attribute the first comment to
		// instead be previous node's trailing comment
		prevEnd := l.prevSym.End().Line
		comments := n.LeadingComments()
		c := comments[0]
		commentStart := c.Start.Line
		if commentStart == prevEnd ***REMOVED***
			// comment is on same line as previous symbol
			n.PopLeadingComment()
			l.prevSym.PushTrailingComment(c)
		***REMOVED*** else if commentStart == prevEnd+1 ***REMOVED***
			// comment is right after previous symbol; see if it is detached
			// and if so re-attribute
			singleLineStyle := strings.HasPrefix(c.Text, "//")
			line := c.End.Line
			groupEnd := -1
			for i := 1; i < len(comments); i++ ***REMOVED***
				c := comments[i]
				newGroup := false
				if !singleLineStyle || c.Start.Line > line+1 ***REMOVED***
					// we've found a gap between comments, which means the
					// previous comments were detached
					newGroup = true
				***REMOVED*** else ***REMOVED***
					line = c.End.Line
					singleLineStyle = strings.HasPrefix(comments[i].Text, "//")
					if !singleLineStyle ***REMOVED***
						// we've found a switch from // comments to /*
						// consider that a new group which means the
						// previous comments were detached
						newGroup = true
					***REMOVED***
				***REMOVED***
				if newGroup ***REMOVED***
					groupEnd = i
					break
				***REMOVED***
			***REMOVED***

			if groupEnd == -1 ***REMOVED***
				// just one group of comments; we'll mark it as a trailing
				// comment if it immediately follows previous symbol and is
				// detached from current symbol
				c1 := comments[0]
				c2 := comments[len(comments)-1]
				if c1.Start.Line <= prevEnd+1 && c2.End.Line < nStart-1 ***REMOVED***
					groupEnd = len(comments)
				***REMOVED***
			***REMOVED***

			for i := 0; i < groupEnd; i++ ***REMOVED***
				l.prevSym.PushTrailingComment(n.PopLeadingComment())
			***REMOVED***
		***REMOVED***
	***REMOVED***

	l.prevSym = n
***REMOVED***

func (l *protoLex) setString(lval *protoSymType, val string) ***REMOVED***
	lval.s = ast.NewStringLiteralNode(val, l.newTokenInfo())
	l.setPrev(lval.s, false)
***REMOVED***

func (l *protoLex) setIdent(lval *protoSymType, val string) ***REMOVED***
	lval.id = ast.NewIdentNode(val, l.newTokenInfo())
	l.setPrev(lval.id, false)
***REMOVED***

func (l *protoLex) setInt(lval *protoSymType, val uint64) ***REMOVED***
	lval.i = ast.NewUintLiteralNode(val, l.newTokenInfo())
	l.setPrev(lval.i, false)
***REMOVED***

func (l *protoLex) setFloat(lval *protoSymType, val float64) ***REMOVED***
	lval.f = ast.NewFloatLiteralNode(val, l.newTokenInfo())
	l.setPrev(lval.f, false)
***REMOVED***

func (l *protoLex) setRune(lval *protoSymType, val rune) ***REMOVED***
	lval.b = ast.NewRuneNode(val, l.newTokenInfo())
	l.setPrev(lval.b, val == '.')
***REMOVED***

func (l *protoLex) setError(lval *protoSymType, err error) ***REMOVED***
	lval.err = l.addSourceError(err)
***REMOVED***

func (l *protoLex) readNumber(sofar ...rune) string ***REMOVED***
	token := sofar
	allowExpSign := false
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if (c == '-' || c == '+') && !allowExpSign ***REMOVED***
			l.input.unreadRune(c)
			break
		***REMOVED***
		allowExpSign = false
		if c != '.' && c != '_' && (c < '0' || c > '9') &&
			(c < 'a' || c > 'z') && (c < 'A' || c > 'Z') &&
			c != '-' && c != '+' ***REMOVED***
			// no more chars in the number token
			l.input.unreadRune(c)
			break
		***REMOVED***
		if c == 'e' || c == 'E' ***REMOVED***
			// scientific notation char can be followed by
			// an exponent sign
			allowExpSign = true
		***REMOVED***
		l.adjustPos(c)
		token = append(token, c)
	***REMOVED***
	return string(token)
***REMOVED***

func numError(err error, kind, s string) error ***REMOVED***
	ne, ok := err.(*strconv.NumError)
	if !ok ***REMOVED***
		return err
	***REMOVED***
	if ne.Err == strconv.ErrRange ***REMOVED***
		return fmt.Errorf("value out of range for %s: %s", kind, s)
	***REMOVED***
	// syntax error
	return fmt.Errorf("invalid syntax in %s value: %s", kind, s)
***REMOVED***

func (l *protoLex) readIdentifier(sofar []rune) []rune ***REMOVED***
	token := sofar
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if c != '_' && (c < 'a' || c > 'z') && (c < 'A' || c > 'Z') && (c < '0' || c > '9') ***REMOVED***
			l.input.unreadRune(c)
			break
		***REMOVED***
		l.adjustPos(c)
		token = append(token, c)
	***REMOVED***
	return token
***REMOVED***

func (l *protoLex) readStringLiteral(quote rune) (string, error) ***REMOVED***
	var buf bytes.Buffer
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			if err == io.EOF ***REMOVED***
				err = io.ErrUnexpectedEOF
			***REMOVED***
			return "", err
		***REMOVED***
		if c == '\n' ***REMOVED***
			return "", errors.New("encountered end-of-line before end of string literal")
		***REMOVED***
		l.adjustPos(c)
		if c == quote ***REMOVED***
			break
		***REMOVED***
		if c == 0 ***REMOVED***
			return "", errors.New("null character ('\\0') not allowed in string literal")
		***REMOVED***
		if c == '\\' ***REMOVED***
			// escape sequence
			c, _, err = l.input.readRune()
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			l.adjustPos(c)
			if c == 'x' || c == 'X' ***REMOVED***
				// hex escape
				c, _, err := l.input.readRune()
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				l.adjustPos(c)
				c2, _, err := l.input.readRune()
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				var hex string
				if (c2 < '0' || c2 > '9') && (c2 < 'a' || c2 > 'f') && (c2 < 'A' || c2 > 'F') ***REMOVED***
					l.input.unreadRune(c2)
					hex = string(c)
				***REMOVED*** else ***REMOVED***
					l.adjustPos(c2)
					hex = string([]rune***REMOVED***c, c2***REMOVED***)
				***REMOVED***
				i, err := strconv.ParseInt(hex, 16, 32)
				if err != nil ***REMOVED***
					return "", fmt.Errorf("invalid hex escape: \\x%q", hex)
				***REMOVED***
				buf.WriteByte(byte(i))

			***REMOVED*** else if c >= '0' && c <= '7' ***REMOVED***
				// octal escape
				c2, _, err := l.input.readRune()
				if err != nil ***REMOVED***
					return "", err
				***REMOVED***
				var octal string
				if c2 < '0' || c2 > '7' ***REMOVED***
					l.input.unreadRune(c2)
					octal = string(c)
				***REMOVED*** else ***REMOVED***
					l.adjustPos(c2)
					c3, _, err := l.input.readRune()
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
					if c3 < '0' || c3 > '7' ***REMOVED***
						l.input.unreadRune(c3)
						octal = string([]rune***REMOVED***c, c2***REMOVED***)
					***REMOVED*** else ***REMOVED***
						l.adjustPos(c3)
						octal = string([]rune***REMOVED***c, c2, c3***REMOVED***)
					***REMOVED***
				***REMOVED***
				i, err := strconv.ParseInt(octal, 8, 32)
				if err != nil ***REMOVED***
					return "", fmt.Errorf("invalid octal escape: \\%q", octal)
				***REMOVED***
				if i > 0xff ***REMOVED***
					return "", fmt.Errorf("octal escape is out range, must be between 0 and 377: \\%q", octal)
				***REMOVED***
				buf.WriteByte(byte(i))

			***REMOVED*** else if c == 'u' ***REMOVED***
				// short unicode escape
				u := make([]rune, 4)
				for i := range u ***REMOVED***
					c, _, err := l.input.readRune()
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
					l.adjustPos(c)
					u[i] = c
				***REMOVED***
				i, err := strconv.ParseInt(string(u), 16, 32)
				if err != nil ***REMOVED***
					return "", fmt.Errorf("invalid unicode escape: \\u%q", string(u))
				***REMOVED***
				buf.WriteRune(rune(i))

			***REMOVED*** else if c == 'U' ***REMOVED***
				// long unicode escape
				u := make([]rune, 8)
				for i := range u ***REMOVED***
					c, _, err := l.input.readRune()
					if err != nil ***REMOVED***
						return "", err
					***REMOVED***
					l.adjustPos(c)
					u[i] = c
				***REMOVED***
				i, err := strconv.ParseInt(string(u), 16, 32)
				if err != nil ***REMOVED***
					return "", fmt.Errorf("invalid unicode escape: \\U%q", string(u))
				***REMOVED***
				if i > 0x10ffff || i < 0 ***REMOVED***
					return "", fmt.Errorf("unicode escape is out of range, must be between 0 and 0x10ffff: \\U%q", string(u))
				***REMOVED***
				buf.WriteRune(rune(i))

			***REMOVED*** else if c == 'a' ***REMOVED***
				buf.WriteByte('\a')
			***REMOVED*** else if c == 'b' ***REMOVED***
				buf.WriteByte('\b')
			***REMOVED*** else if c == 'f' ***REMOVED***
				buf.WriteByte('\f')
			***REMOVED*** else if c == 'n' ***REMOVED***
				buf.WriteByte('\n')
			***REMOVED*** else if c == 'r' ***REMOVED***
				buf.WriteByte('\r')
			***REMOVED*** else if c == 't' ***REMOVED***
				buf.WriteByte('\t')
			***REMOVED*** else if c == 'v' ***REMOVED***
				buf.WriteByte('\v')
			***REMOVED*** else if c == '\\' ***REMOVED***
				buf.WriteByte('\\')
			***REMOVED*** else if c == '\'' ***REMOVED***
				buf.WriteByte('\'')
			***REMOVED*** else if c == '"' ***REMOVED***
				buf.WriteByte('"')
			***REMOVED*** else if c == '?' ***REMOVED***
				buf.WriteByte('?')
			***REMOVED*** else ***REMOVED***
				return "", fmt.Errorf("invalid escape sequence: %q", "\\"+string(c))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			buf.WriteRune(c)
		***REMOVED***
	***REMOVED***
	return buf.String(), nil
***REMOVED***

func (l *protoLex) skipToEndOfLineComment(lval *protoSymType) (ok, hasErr bool) ***REMOVED***
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			return false, false
		***REMOVED***
		switch c ***REMOVED***
		case '\n':
			return true, false
		case 0:
			l.setError(lval, errors.New("invalid control character"))
			return false, true
		***REMOVED***
		l.adjustPos(c)
	***REMOVED***
***REMOVED***

func (l *protoLex) skipToEndOfBlockComment(lval *protoSymType) (ok, hasErr bool) ***REMOVED***
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			return false, false
		***REMOVED***
		if c == 0 ***REMOVED***
			l.setError(lval, errors.New("invalid control character"))
			return false, true
		***REMOVED***
		l.adjustPos(c)
		if c == '*' ***REMOVED***
			c, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				return false, false
			***REMOVED***
			if c == '/' ***REMOVED***
				l.adjustPos(c)
				return true, false
			***REMOVED***
			l.input.unreadRune(c)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (l *protoLex) addSourceError(err error) ErrorWithPos ***REMOVED***
	ewp, ok := err.(ErrorWithPos)
	if !ok ***REMOVED***
		ewp = ErrorWithSourcePos***REMOVED***Pos: l.prev(), Underlying: err***REMOVED***
	***REMOVED***
	_ = l.errs.handleError(ewp)
	return ewp
***REMOVED***

func (l *protoLex) Error(s string) ***REMOVED***
	_ = l.addSourceError(errors.New(s))
***REMOVED***
