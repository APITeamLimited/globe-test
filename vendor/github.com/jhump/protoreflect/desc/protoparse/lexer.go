package protoparse

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

type runeReader struct ***REMOVED***
	rr     *bufio.Reader
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
		return r, utf8.RuneLen(r), nil
	***REMOVED***
	r, sz, err := rr.rr.ReadRune()
	if err != nil ***REMOVED***
		rr.err = err
	***REMOVED***
	return r, sz, err
***REMOVED***

func (rr *runeReader) unreadRune(r rune) ***REMOVED***
	rr.unread = append(rr.unread, r)
***REMOVED***

func lexError(l protoLexer, pos *SourcePos, err string) ***REMOVED***
	pl := l.(*protoLex)
	_ = pl.errs.handleErrorWithPos(pos, err)
***REMOVED***

type protoLex struct ***REMOVED***
	filename string
	input    *runeReader
	errs     *errorHandler
	res      *fileNode

	lineNo int
	colNo  int
	offset int

	prevSym terminalNode

	prevLineNo int
	prevColNo  int
	prevOffset int
	comments   []comment
***REMOVED***

func newLexer(in io.Reader, filename string, errs *errorHandler) *protoLex ***REMOVED***
	return &protoLex***REMOVED***
		input:    &runeReader***REMOVED***rr: bufio.NewReader(in)***REMOVED***,
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
	return l.prevSym.start()
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

	for ***REMOVED***
		c, n, err := l.input.readRune()
		if err == io.EOF ***REMOVED***
			// we're not actually returning a rune, but this will associate
			// accumulated comments as a trailing comment on last symbol
			// (if appropriate)
			l.setRune(lval)
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
		if strings.ContainsRune("\n\r\t ", c) ***REMOVED***
			continue
		***REMOVED***

		if c == '.' ***REMOVED***
			// decimal literals could start with a dot
			cn, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				l.setDot(lval)
				return int(c)
			***REMOVED***
			if cn >= '0' && cn <= '9' ***REMOVED***
				l.adjustPos(cn)
				token := []rune***REMOVED***c, cn***REMOVED***
				token = l.readNumber(token, false, true)
				f, err := strconv.ParseFloat(string(token), 64)
				if err != nil ***REMOVED***
					l.setError(lval, err)
					return _ERROR
				***REMOVED***
				l.setFloat(lval, f)
				return _FLOAT_LIT
			***REMOVED***
			l.input.unreadRune(cn)
			l.setDot(lval)
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
			if c == '0' ***REMOVED***
				cn, _, err := l.input.readRune()
				if err != nil ***REMOVED***
					l.setInt(lval, 0)
					return _INT_LIT
				***REMOVED***
				if cn == 'x' || cn == 'X' ***REMOVED***
					cnn, _, err := l.input.readRune()
					if err != nil ***REMOVED***
						l.input.unreadRune(cn)
						l.setInt(lval, 0)
						return _INT_LIT
					***REMOVED***
					if (cnn >= '0' && cnn <= '9') || (cnn >= 'a' && cnn <= 'f') || (cnn >= 'A' && cnn <= 'F') ***REMOVED***
						// hexadecimal!
						l.adjustPos(cn, cnn)
						token := []rune***REMOVED***cnn***REMOVED***
						token = l.readHexNumber(token)
						ui, err := strconv.ParseUint(string(token), 16, 64)
						if err != nil ***REMOVED***
							l.setError(lval, err)
							return _ERROR
						***REMOVED***
						l.setInt(lval, ui)
						return _INT_LIT
					***REMOVED***
					l.input.unreadRune(cnn)
					l.input.unreadRune(cn)
					l.setInt(lval, 0)
					return _INT_LIT
				***REMOVED*** else ***REMOVED***
					l.input.unreadRune(cn)
				***REMOVED***
			***REMOVED***
			token := []rune***REMOVED***c***REMOVED***
			token = l.readNumber(token, true, true)
			numstr := string(token)
			if strings.Contains(numstr, ".") || strings.Contains(numstr, "e") || strings.Contains(numstr, "E") ***REMOVED***
				// floating point!
				f, err := strconv.ParseFloat(numstr, 64)
				if err != nil ***REMOVED***
					l.setError(lval, err)
					return _ERROR
				***REMOVED***
				l.setFloat(lval, f)
				return _FLOAT_LIT
			***REMOVED***
			// integer! (decimal or octal)
			ui, err := strconv.ParseUint(numstr, 0, 64)
			if err != nil ***REMOVED***
				if numErr, ok := err.(*strconv.NumError); ok && numErr.Err == strconv.ErrRange ***REMOVED***
					// if it's too big to be an int, parse it as a float
					var f float64
					f, err = strconv.ParseFloat(numstr, 64)
					if err == nil ***REMOVED***
						l.setFloat(lval, f)
						return _FLOAT_LIT
					***REMOVED***
				***REMOVED***
				l.setError(lval, err)
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
				l.setRune(lval)
				return int(c)
			***REMOVED***
			if cn == '/' ***REMOVED***
				l.adjustPos(cn)
				hitNewline, txt := l.skipToEndOfLineComment()
				commentPos := l.posRange()
				commentPos.end.Col++
				if hitNewline ***REMOVED***
					// we don't do this inside of skipToEndOfLineComment
					// because we want to know the length of previous
					// line for calculation above
					l.adjustPos('\n')
				***REMOVED***
				l.comments = append(l.comments, comment***REMOVED***posRange: commentPos, text: txt***REMOVED***)
				continue
			***REMOVED***
			if cn == '*' ***REMOVED***
				l.adjustPos(cn)
				if txt, ok := l.skipToEndOfBlockComment(); !ok ***REMOVED***
					l.setError(lval, errors.New("block comment never terminates, unexpected EOF"))
					return _ERROR
				***REMOVED*** else ***REMOVED***
					l.comments = append(l.comments, comment***REMOVED***posRange: l.posRange(), text: txt***REMOVED***)
				***REMOVED***
				continue
			***REMOVED***
			l.input.unreadRune(cn)
		***REMOVED***

		l.setRune(lval)
		return int(c)
	***REMOVED***
***REMOVED***

func (l *protoLex) posRange() posRange ***REMOVED***
	return posRange***REMOVED***
		start: SourcePos***REMOVED***
			Filename: l.filename,
			Offset:   l.prevOffset,
			Line:     l.prevLineNo + 1,
			Col:      l.prevColNo + 1,
		***REMOVED***,
		end: l.cur(),
	***REMOVED***
***REMOVED***

func (l *protoLex) newBasicNode() basicNode ***REMOVED***
	return basicNode***REMOVED***
		posRange: l.posRange(),
		leading:  l.comments,
	***REMOVED***
***REMOVED***

func (l *protoLex) setPrev(n terminalNode, isDot bool) ***REMOVED***
	nStart := n.start().Line
	if _, ok := n.(*basicNode); ok ***REMOVED***
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
	if l.prevSym != nil && len(n.leadingComments()) > 0 && l.prevSym.end().Line < nStart ***REMOVED***
		// we may need to re-attribute the first comment to
		// instead be previous node's trailing comment
		prevEnd := l.prevSym.end().Line
		comments := n.leadingComments()
		c := comments[0]
		commentStart := c.start.Line
		if commentStart == prevEnd ***REMOVED***
			// comment is on same line as previous symbol
			n.popLeadingComment()
			l.prevSym.pushTrailingComment(c)
		***REMOVED*** else if commentStart == prevEnd+1 ***REMOVED***
			// comment is right after previous symbol; see if it is detached
			// and if so re-attribute
			singleLineStyle := strings.HasPrefix(c.text, "//")
			line := c.end.Line
			groupEnd := -1
			for i := 1; i < len(comments); i++ ***REMOVED***
				c := comments[i]
				newGroup := false
				if !singleLineStyle || c.start.Line > line+1 ***REMOVED***
					// we've found a gap between comments, which means the
					// previous comments were detached
					newGroup = true
				***REMOVED*** else ***REMOVED***
					line = c.end.Line
					singleLineStyle = strings.HasPrefix(comments[i].text, "//")
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
				if c1.start.Line <= prevEnd+1 && c2.end.Line < nStart-1 ***REMOVED***
					groupEnd = len(comments)
				***REMOVED***
			***REMOVED***

			for i := 0; i < groupEnd; i++ ***REMOVED***
				l.prevSym.pushTrailingComment(n.popLeadingComment())
			***REMOVED***
		***REMOVED***
	***REMOVED***

	l.prevSym = n
***REMOVED***

func (l *protoLex) setString(lval *protoSymType, val string) ***REMOVED***
	lval.s = &stringLiteralNode***REMOVED***basicNode: l.newBasicNode(), val: val***REMOVED***
	l.setPrev(lval.s, false)
***REMOVED***

func (l *protoLex) setIdent(lval *protoSymType, val string) ***REMOVED***
	lval.id = &identNode***REMOVED***basicNode: l.newBasicNode(), val: val***REMOVED***
	l.setPrev(lval.id, false)
***REMOVED***

func (l *protoLex) setInt(lval *protoSymType, val uint64) ***REMOVED***
	lval.i = &intLiteralNode***REMOVED***basicNode: l.newBasicNode(), val: val***REMOVED***
	l.setPrev(lval.i, false)
***REMOVED***

func (l *protoLex) setFloat(lval *protoSymType, val float64) ***REMOVED***
	lval.f = &floatLiteralNode***REMOVED***basicNode: l.newBasicNode(), val: val***REMOVED***
	l.setPrev(lval.f, false)
***REMOVED***

func (l *protoLex) setRune(lval *protoSymType) ***REMOVED***
	b := l.newBasicNode()
	lval.b = &b
	l.setPrev(lval.b, false)
***REMOVED***

func (l *protoLex) setDot(lval *protoSymType) ***REMOVED***
	b := l.newBasicNode()
	lval.b = &b
	l.setPrev(lval.b, true)
***REMOVED***

func (l *protoLex) setError(lval *protoSymType, err error) ***REMOVED***
	lval.err = l.addSourceError(err)
***REMOVED***

func (l *protoLex) readNumber(sofar []rune, allowDot bool, allowExp bool) []rune ***REMOVED***
	token := sofar
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if c == '.' ***REMOVED***
			if !allowDot ***REMOVED***
				l.input.unreadRune(c)
				break
			***REMOVED***
			allowDot = false
		***REMOVED*** else if c == 'e' || c == 'E' ***REMOVED***
			if !allowExp ***REMOVED***
				l.input.unreadRune(c)
				break
			***REMOVED***
			allowExp = false
			cn, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				l.input.unreadRune(c)
				break
			***REMOVED***
			if cn == '-' || cn == '+' ***REMOVED***
				cnn, _, err := l.input.readRune()
				if err != nil ***REMOVED***
					l.input.unreadRune(cn)
					l.input.unreadRune(c)
					break
				***REMOVED***
				if cnn < '0' || cnn > '9' ***REMOVED***
					l.input.unreadRune(cnn)
					l.input.unreadRune(cn)
					l.input.unreadRune(c)
					break
				***REMOVED***
				l.adjustPos(c)
				token = append(token, c)
				c, cn = cn, cnn
			***REMOVED*** else if cn < '0' || cn > '9' ***REMOVED***
				l.input.unreadRune(cn)
				l.input.unreadRune(c)
				break
			***REMOVED***
			l.adjustPos(c)
			token = append(token, c)
			c = cn
		***REMOVED*** else if c < '0' || c > '9' ***REMOVED***
			l.input.unreadRune(c)
			break
		***REMOVED***
		l.adjustPos(c)
		token = append(token, c)
	***REMOVED***
	return token
***REMOVED***

func (l *protoLex) readHexNumber(sofar []rune) []rune ***REMOVED***
	token := sofar
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			break
		***REMOVED***
		if (c < 'a' || c > 'f') && (c < 'A' || c > 'F') && (c < '0' || c > '9') ***REMOVED***
			l.input.unreadRune(c)
			break
		***REMOVED***
		l.adjustPos(c)
		token = append(token, c)
	***REMOVED***
	return token
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

func (l *protoLex) skipToEndOfLineComment() (bool, string) ***REMOVED***
	txt := []rune***REMOVED***'/', '/'***REMOVED***
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			return false, string(txt)
		***REMOVED***
		if c == '\n' ***REMOVED***
			return true, string(append(txt, '\n'))
		***REMOVED***
		l.adjustPos(c)
		txt = append(txt, c)
	***REMOVED***
***REMOVED***

func (l *protoLex) skipToEndOfBlockComment() (string, bool) ***REMOVED***
	txt := []rune***REMOVED***'/', '*'***REMOVED***
	for ***REMOVED***
		c, _, err := l.input.readRune()
		if err != nil ***REMOVED***
			return "", false
		***REMOVED***
		l.adjustPos(c)
		txt = append(txt, c)
		if c == '*' ***REMOVED***
			c, _, err := l.input.readRune()
			if err != nil ***REMOVED***
				return "", false
			***REMOVED***
			if c == '/' ***REMOVED***
				l.adjustPos(c)
				txt = append(txt, c)
				return string(txt), true
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
