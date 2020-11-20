// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"

	"golang.org/x/net/html/atom"
)

// A TokenType is the type of a Token.
type TokenType uint32

const (
	// ErrorToken means that an error occurred during tokenization.
	ErrorToken TokenType = iota
	// TextToken means a text node.
	TextToken
	// A StartTagToken looks like <a>.
	StartTagToken
	// An EndTagToken looks like </a>.
	EndTagToken
	// A SelfClosingTagToken tag looks like <br/>.
	SelfClosingTagToken
	// A CommentToken looks like <!--x-->.
	CommentToken
	// A DoctypeToken looks like <!DOCTYPE x>
	DoctypeToken
)

// ErrBufferExceeded means that the buffering limit was exceeded.
var ErrBufferExceeded = errors.New("max buffer exceeded")

// String returns a string representation of the TokenType.
func (t TokenType) String() string ***REMOVED***
	switch t ***REMOVED***
	case ErrorToken:
		return "Error"
	case TextToken:
		return "Text"
	case StartTagToken:
		return "StartTag"
	case EndTagToken:
		return "EndTag"
	case SelfClosingTagToken:
		return "SelfClosingTag"
	case CommentToken:
		return "Comment"
	case DoctypeToken:
		return "Doctype"
	***REMOVED***
	return "Invalid(" + strconv.Itoa(int(t)) + ")"
***REMOVED***

// An Attribute is an attribute namespace-key-value triple. Namespace is
// non-empty for foreign attributes like xlink, Key is alphabetic (and hence
// does not contain escapable characters like '&', '<' or '>'), and Val is
// unescaped (it looks like "a<b" rather than "a&lt;b").
//
// Namespace is only used by the parser, not the tokenizer.
type Attribute struct ***REMOVED***
	Namespace, Key, Val string
***REMOVED***

// A Token consists of a TokenType and some Data (tag name for start and end
// tags, content for text, comments and doctypes). A tag Token may also contain
// a slice of Attributes. Data is unescaped for all Tokens (it looks like "a<b"
// rather than "a&lt;b"). For tag Tokens, DataAtom is the atom for Data, or
// zero if Data is not a known tag name.
type Token struct ***REMOVED***
	Type     TokenType
	DataAtom atom.Atom
	Data     string
	Attr     []Attribute
***REMOVED***

// tagString returns a string representation of a tag Token's Data and Attr.
func (t Token) tagString() string ***REMOVED***
	if len(t.Attr) == 0 ***REMOVED***
		return t.Data
	***REMOVED***
	buf := bytes.NewBufferString(t.Data)
	for _, a := range t.Attr ***REMOVED***
		buf.WriteByte(' ')
		buf.WriteString(a.Key)
		buf.WriteString(`="`)
		escape(buf, a.Val)
		buf.WriteByte('"')
	***REMOVED***
	return buf.String()
***REMOVED***

// String returns a string representation of the Token.
func (t Token) String() string ***REMOVED***
	switch t.Type ***REMOVED***
	case ErrorToken:
		return ""
	case TextToken:
		return EscapeString(t.Data)
	case StartTagToken:
		return "<" + t.tagString() + ">"
	case EndTagToken:
		return "</" + t.tagString() + ">"
	case SelfClosingTagToken:
		return "<" + t.tagString() + "/>"
	case CommentToken:
		return "<!--" + t.Data + "-->"
	case DoctypeToken:
		return "<!DOCTYPE " + t.Data + ">"
	***REMOVED***
	return "Invalid(" + strconv.Itoa(int(t.Type)) + ")"
***REMOVED***

// span is a range of bytes in a Tokenizer's buffer. The start is inclusive,
// the end is exclusive.
type span struct ***REMOVED***
	start, end int
***REMOVED***

// A Tokenizer returns a stream of HTML Tokens.
type Tokenizer struct ***REMOVED***
	// r is the source of the HTML text.
	r io.Reader
	// tt is the TokenType of the current token.
	tt TokenType
	// err is the first error encountered during tokenization. It is possible
	// for tt != Error && err != nil to hold: this means that Next returned a
	// valid token but the subsequent Next call will return an error token.
	// For example, if the HTML text input was just "plain", then the first
	// Next call would set z.err to io.EOF but return a TextToken, and all
	// subsequent Next calls would return an ErrorToken.
	// err is never reset. Once it becomes non-nil, it stays non-nil.
	err error
	// readErr is the error returned by the io.Reader r. It is separate from
	// err because it is valid for an io.Reader to return (n int, err1 error)
	// such that n > 0 && err1 != nil, and callers should always process the
	// n > 0 bytes before considering the error err1.
	readErr error
	// buf[raw.start:raw.end] holds the raw bytes of the current token.
	// buf[raw.end:] is buffered input that will yield future tokens.
	raw span
	buf []byte
	// maxBuf limits the data buffered in buf. A value of 0 means unlimited.
	maxBuf int
	// buf[data.start:data.end] holds the raw bytes of the current token's data:
	// a text token's text, a tag token's tag name, etc.
	data span
	// pendingAttr is the attribute key and value currently being tokenized.
	// When complete, pendingAttr is pushed onto attr. nAttrReturned is
	// incremented on each call to TagAttr.
	pendingAttr   [2]span
	attr          [][2]span
	nAttrReturned int
	// rawTag is the "script" in "</script>" that closes the next token. If
	// non-empty, the subsequent call to Next will return a raw or RCDATA text
	// token: one that treats "<p>" as text instead of an element.
	// rawTag's contents are lower-cased.
	rawTag string
	// textIsRaw is whether the current text token's data is not escaped.
	textIsRaw bool
	// convertNUL is whether NUL bytes in the current token's data should
	// be converted into \ufffd replacement characters.
	convertNUL bool
	// allowCDATA is whether CDATA sections are allowed in the current context.
	allowCDATA bool
***REMOVED***

// AllowCDATA sets whether or not the tokenizer recognizes <![CDATA[foo]]> as
// the text "foo". The default value is false, which means to recognize it as
// a bogus comment "<!-- [CDATA[foo]] -->" instead.
//
// Strictly speaking, an HTML5 compliant tokenizer should allow CDATA if and
// only if tokenizing foreign content, such as MathML and SVG. However,
// tracking foreign-contentness is difficult to do purely in the tokenizer,
// as opposed to the parser, due to HTML integration points: an <svg> element
// can contain a <foreignObject> that is foreign-to-SVG but not foreign-to-
// HTML. For strict compliance with the HTML5 tokenization algorithm, it is the
// responsibility of the user of a tokenizer to call AllowCDATA as appropriate.
// In practice, if using the tokenizer without caring whether MathML or SVG
// CDATA is text or comments, such as tokenizing HTML to find all the anchor
// text, it is acceptable to ignore this responsibility.
func (z *Tokenizer) AllowCDATA(allowCDATA bool) ***REMOVED***
	z.allowCDATA = allowCDATA
***REMOVED***

// NextIsNotRawText instructs the tokenizer that the next token should not be
// considered as 'raw text'. Some elements, such as script and title elements,
// normally require the next token after the opening tag to be 'raw text' that
// has no child elements. For example, tokenizing "<title>a<b>c</b>d</title>"
// yields a start tag token for "<title>", a text token for "a<b>c</b>d", and
// an end tag token for "</title>". There are no distinct start tag or end tag
// tokens for the "<b>" and "</b>".
//
// This tokenizer implementation will generally look for raw text at the right
// times. Strictly speaking, an HTML5 compliant tokenizer should not look for
// raw text if in foreign content: <title> generally needs raw text, but a
// <title> inside an <svg> does not. Another example is that a <textarea>
// generally needs raw text, but a <textarea> is not allowed as an immediate
// child of a <select>; in normal parsing, a <textarea> implies </select>, but
// one cannot close the implicit element when parsing a <select>'s InnerHTML.
// Similarly to AllowCDATA, tracking the correct moment to override raw-text-
// ness is difficult to do purely in the tokenizer, as opposed to the parser.
// For strict compliance with the HTML5 tokenization algorithm, it is the
// responsibility of the user of a tokenizer to call NextIsNotRawText as
// appropriate. In practice, like AllowCDATA, it is acceptable to ignore this
// responsibility for basic usage.
//
// Note that this 'raw text' concept is different from the one offered by the
// Tokenizer.Raw method.
func (z *Tokenizer) NextIsNotRawText() ***REMOVED***
	z.rawTag = ""
***REMOVED***

// Err returns the error associated with the most recent ErrorToken token.
// This is typically io.EOF, meaning the end of tokenization.
func (z *Tokenizer) Err() error ***REMOVED***
	if z.tt != ErrorToken ***REMOVED***
		return nil
	***REMOVED***
	return z.err
***REMOVED***

// readByte returns the next byte from the input stream, doing a buffered read
// from z.r into z.buf if necessary. z.buf[z.raw.start:z.raw.end] remains a contiguous byte
// slice that holds all the bytes read so far for the current token.
// It sets z.err if the underlying reader returns an error.
// Pre-condition: z.err == nil.
func (z *Tokenizer) readByte() byte ***REMOVED***
	if z.raw.end >= len(z.buf) ***REMOVED***
		// Our buffer is exhausted and we have to read from z.r. Check if the
		// previous read resulted in an error.
		if z.readErr != nil ***REMOVED***
			z.err = z.readErr
			return 0
		***REMOVED***
		// We copy z.buf[z.raw.start:z.raw.end] to the beginning of z.buf. If the length
		// z.raw.end - z.raw.start is more than half the capacity of z.buf, then we
		// allocate a new buffer before the copy.
		c := cap(z.buf)
		d := z.raw.end - z.raw.start
		var buf1 []byte
		if 2*d > c ***REMOVED***
			buf1 = make([]byte, d, 2*c)
		***REMOVED*** else ***REMOVED***
			buf1 = z.buf[:d]
		***REMOVED***
		copy(buf1, z.buf[z.raw.start:z.raw.end])
		if x := z.raw.start; x != 0 ***REMOVED***
			// Adjust the data/attr spans to refer to the same contents after the copy.
			z.data.start -= x
			z.data.end -= x
			z.pendingAttr[0].start -= x
			z.pendingAttr[0].end -= x
			z.pendingAttr[1].start -= x
			z.pendingAttr[1].end -= x
			for i := range z.attr ***REMOVED***
				z.attr[i][0].start -= x
				z.attr[i][0].end -= x
				z.attr[i][1].start -= x
				z.attr[i][1].end -= x
			***REMOVED***
		***REMOVED***
		z.raw.start, z.raw.end, z.buf = 0, d, buf1[:d]
		// Now that we have copied the live bytes to the start of the buffer,
		// we read from z.r into the remainder.
		var n int
		n, z.readErr = readAtLeastOneByte(z.r, buf1[d:cap(buf1)])
		if n == 0 ***REMOVED***
			z.err = z.readErr
			return 0
		***REMOVED***
		z.buf = buf1[:d+n]
	***REMOVED***
	x := z.buf[z.raw.end]
	z.raw.end++
	if z.maxBuf > 0 && z.raw.end-z.raw.start >= z.maxBuf ***REMOVED***
		z.err = ErrBufferExceeded
		return 0
	***REMOVED***
	return x
***REMOVED***

// Buffered returns a slice containing data buffered but not yet tokenized.
func (z *Tokenizer) Buffered() []byte ***REMOVED***
	return z.buf[z.raw.end:]
***REMOVED***

// readAtLeastOneByte wraps an io.Reader so that reading cannot return (0, nil).
// It returns io.ErrNoProgress if the underlying r.Read method returns (0, nil)
// too many times in succession.
func readAtLeastOneByte(r io.Reader, b []byte) (int, error) ***REMOVED***
	for i := 0; i < 100; i++ ***REMOVED***
		if n, err := r.Read(b); n != 0 || err != nil ***REMOVED***
			return n, err
		***REMOVED***
	***REMOVED***
	return 0, io.ErrNoProgress
***REMOVED***

// skipWhiteSpace skips past any white space.
func (z *Tokenizer) skipWhiteSpace() ***REMOVED***
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			return
		***REMOVED***
		switch c ***REMOVED***
		case ' ', '\n', '\r', '\t', '\f':
			// No-op.
		default:
			z.raw.end--
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readRawOrRCDATA reads until the next "</foo>", where "foo" is z.rawTag and
// is typically something like "script" or "textarea".
func (z *Tokenizer) readRawOrRCDATA() ***REMOVED***
	if z.rawTag == "script" ***REMOVED***
		z.readScript()
		z.textIsRaw = true
		z.rawTag = ""
		return
	***REMOVED***
loop:
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c != '<' ***REMOVED***
			continue loop
		***REMOVED***
		c = z.readByte()
		if z.err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c != '/' ***REMOVED***
			z.raw.end--
			continue loop
		***REMOVED***
		if z.readRawEndTag() || z.err != nil ***REMOVED***
			break loop
		***REMOVED***
	***REMOVED***
	z.data.end = z.raw.end
	// A textarea's or title's RCDATA can contain escaped entities.
	z.textIsRaw = z.rawTag != "textarea" && z.rawTag != "title"
	z.rawTag = ""
***REMOVED***

// readRawEndTag attempts to read a tag like "</foo>", where "foo" is z.rawTag.
// If it succeeds, it backs up the input position to reconsume the tag and
// returns true. Otherwise it returns false. The opening "</" has already been
// consumed.
func (z *Tokenizer) readRawEndTag() bool ***REMOVED***
	for i := 0; i < len(z.rawTag); i++ ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			return false
		***REMOVED***
		if c != z.rawTag[i] && c != z.rawTag[i]-('a'-'A') ***REMOVED***
			z.raw.end--
			return false
		***REMOVED***
	***REMOVED***
	c := z.readByte()
	if z.err != nil ***REMOVED***
		return false
	***REMOVED***
	switch c ***REMOVED***
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		// The 3 is 2 for the leading "</" plus 1 for the trailing character c.
		z.raw.end -= 3 + len(z.rawTag)
		return true
	***REMOVED***
	z.raw.end--
	return false
***REMOVED***

// readScript reads until the next </script> tag, following the byzantine
// rules for escaping/hiding the closing tag.
func (z *Tokenizer) readScript() ***REMOVED***
	defer func() ***REMOVED***
		z.data.end = z.raw.end
	***REMOVED***()
	var c byte

scriptData:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c == '<' ***REMOVED***
		goto scriptDataLessThanSign
	***REMOVED***
	goto scriptData

scriptDataLessThanSign:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '/':
		goto scriptDataEndTagOpen
	case '!':
		goto scriptDataEscapeStart
	***REMOVED***
	z.raw.end--
	goto scriptData

scriptDataEndTagOpen:
	if z.readRawEndTag() || z.err != nil ***REMOVED***
		return
	***REMOVED***
	goto scriptData

scriptDataEscapeStart:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c == '-' ***REMOVED***
		goto scriptDataEscapeStartDash
	***REMOVED***
	z.raw.end--
	goto scriptData

scriptDataEscapeStartDash:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c == '-' ***REMOVED***
		goto scriptDataEscapedDashDash
	***REMOVED***
	z.raw.end--
	goto scriptData

scriptDataEscaped:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataEscapedDash
	case '<':
		goto scriptDataEscapedLessThanSign
	***REMOVED***
	goto scriptDataEscaped

scriptDataEscapedDash:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	***REMOVED***
	goto scriptDataEscaped

scriptDataEscapedDashDash:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataEscapedDashDash
	case '<':
		goto scriptDataEscapedLessThanSign
	case '>':
		goto scriptData
	***REMOVED***
	goto scriptDataEscaped

scriptDataEscapedLessThanSign:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c == '/' ***REMOVED***
		goto scriptDataEscapedEndTagOpen
	***REMOVED***
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' ***REMOVED***
		goto scriptDataDoubleEscapeStart
	***REMOVED***
	z.raw.end--
	goto scriptData

scriptDataEscapedEndTagOpen:
	if z.readRawEndTag() || z.err != nil ***REMOVED***
		return
	***REMOVED***
	goto scriptDataEscaped

scriptDataDoubleEscapeStart:
	z.raw.end--
	for i := 0; i < len("script"); i++ ***REMOVED***
		c = z.readByte()
		if z.err != nil ***REMOVED***
			return
		***REMOVED***
		if c != "script"[i] && c != "SCRIPT"[i] ***REMOVED***
			z.raw.end--
			goto scriptDataEscaped
		***REMOVED***
	***REMOVED***
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case ' ', '\n', '\r', '\t', '\f', '/', '>':
		goto scriptDataDoubleEscaped
	***REMOVED***
	z.raw.end--
	goto scriptDataEscaped

scriptDataDoubleEscaped:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataDoubleEscapedDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	***REMOVED***
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDash:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	***REMOVED***
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedDashDash:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch c ***REMOVED***
	case '-':
		goto scriptDataDoubleEscapedDashDash
	case '<':
		goto scriptDataDoubleEscapedLessThanSign
	case '>':
		goto scriptData
	***REMOVED***
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapedLessThanSign:
	c = z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c == '/' ***REMOVED***
		goto scriptDataDoubleEscapeEnd
	***REMOVED***
	z.raw.end--
	goto scriptDataDoubleEscaped

scriptDataDoubleEscapeEnd:
	if z.readRawEndTag() ***REMOVED***
		z.raw.end += len("</script>")
		goto scriptDataEscaped
	***REMOVED***
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	goto scriptDataDoubleEscaped
***REMOVED***

// readComment reads the next comment token starting with "<!--". The opening
// "<!--" has already been consumed.
func (z *Tokenizer) readComment() ***REMOVED***
	z.data.start = z.raw.end
	defer func() ***REMOVED***
		if z.data.end < z.data.start ***REMOVED***
			// It's a comment with no data, like <!-->.
			z.data.end = z.data.start
		***REMOVED***
	***REMOVED***()
	for dashCount := 2; ; ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			// Ignore up to two dashes at EOF.
			if dashCount > 2 ***REMOVED***
				dashCount = 2
			***REMOVED***
			z.data.end = z.raw.end - dashCount
			return
		***REMOVED***
		switch c ***REMOVED***
		case '-':
			dashCount++
			continue
		case '>':
			if dashCount >= 2 ***REMOVED***
				z.data.end = z.raw.end - len("-->")
				return
			***REMOVED***
		case '!':
			if dashCount >= 2 ***REMOVED***
				c = z.readByte()
				if z.err != nil ***REMOVED***
					z.data.end = z.raw.end
					return
				***REMOVED***
				if c == '>' ***REMOVED***
					z.data.end = z.raw.end - len("--!>")
					return
				***REMOVED***
			***REMOVED***
		***REMOVED***
		dashCount = 0
	***REMOVED***
***REMOVED***

// readUntilCloseAngle reads until the next ">".
func (z *Tokenizer) readUntilCloseAngle() ***REMOVED***
	z.data.start = z.raw.end
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return
		***REMOVED***
		if c == '>' ***REMOVED***
			z.data.end = z.raw.end - len(">")
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readMarkupDeclaration reads the next token starting with "<!". It might be
// a "<!--comment-->", a "<!DOCTYPE foo>", a "<![CDATA[section]]>" or
// "<!a bogus comment". The opening "<!" has already been consumed.
func (z *Tokenizer) readMarkupDeclaration() TokenType ***REMOVED***
	z.data.start = z.raw.end
	var c [2]byte
	for i := 0; i < 2; i++ ***REMOVED***
		c[i] = z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return CommentToken
		***REMOVED***
	***REMOVED***
	if c[0] == '-' && c[1] == '-' ***REMOVED***
		z.readComment()
		return CommentToken
	***REMOVED***
	z.raw.end -= 2
	if z.readDoctype() ***REMOVED***
		return DoctypeToken
	***REMOVED***
	if z.allowCDATA && z.readCDATA() ***REMOVED***
		z.convertNUL = true
		return TextToken
	***REMOVED***
	// It's a bogus comment.
	z.readUntilCloseAngle()
	return CommentToken
***REMOVED***

// readDoctype attempts to read a doctype declaration and returns true if
// successful. The opening "<!" has already been consumed.
func (z *Tokenizer) readDoctype() bool ***REMOVED***
	const s = "DOCTYPE"
	for i := 0; i < len(s); i++ ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return false
		***REMOVED***
		if c != s[i] && c != s[i]+('a'-'A') ***REMOVED***
			// Back up to read the fragment of "DOCTYPE" again.
			z.raw.end = z.data.start
			return false
		***REMOVED***
	***REMOVED***
	if z.skipWhiteSpace(); z.err != nil ***REMOVED***
		z.data.start = z.raw.end
		z.data.end = z.raw.end
		return true
	***REMOVED***
	z.readUntilCloseAngle()
	return true
***REMOVED***

// readCDATA attempts to read a CDATA section and returns true if
// successful. The opening "<!" has already been consumed.
func (z *Tokenizer) readCDATA() bool ***REMOVED***
	const s = "[CDATA["
	for i := 0; i < len(s); i++ ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return false
		***REMOVED***
		if c != s[i] ***REMOVED***
			// Back up to read the fragment of "[CDATA[" again.
			z.raw.end = z.data.start
			return false
		***REMOVED***
	***REMOVED***
	z.data.start = z.raw.end
	brackets := 0
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return true
		***REMOVED***
		switch c ***REMOVED***
		case ']':
			brackets++
		case '>':
			if brackets >= 2 ***REMOVED***
				z.data.end = z.raw.end - len("]]>")
				return true
			***REMOVED***
			brackets = 0
		default:
			brackets = 0
		***REMOVED***
	***REMOVED***
***REMOVED***

// startTagIn returns whether the start tag in z.buf[z.data.start:z.data.end]
// case-insensitively matches any element of ss.
func (z *Tokenizer) startTagIn(ss ...string) bool ***REMOVED***
loop:
	for _, s := range ss ***REMOVED***
		if z.data.end-z.data.start != len(s) ***REMOVED***
			continue loop
		***REMOVED***
		for i := 0; i < len(s); i++ ***REMOVED***
			c := z.buf[z.data.start+i]
			if 'A' <= c && c <= 'Z' ***REMOVED***
				c += 'a' - 'A'
			***REMOVED***
			if c != s[i] ***REMOVED***
				continue loop
			***REMOVED***
		***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

// readStartTag reads the next start tag token. The opening "<a" has already
// been consumed, where 'a' means anything in [A-Za-z].
func (z *Tokenizer) readStartTag() TokenType ***REMOVED***
	z.readTag(true)
	if z.err != nil ***REMOVED***
		return ErrorToken
	***REMOVED***
	// Several tags flag the tokenizer's next token as raw.
	c, raw := z.buf[z.data.start], false
	if 'A' <= c && c <= 'Z' ***REMOVED***
		c += 'a' - 'A'
	***REMOVED***
	switch c ***REMOVED***
	case 'i':
		raw = z.startTagIn("iframe")
	case 'n':
		raw = z.startTagIn("noembed", "noframes", "noscript")
	case 'p':
		raw = z.startTagIn("plaintext")
	case 's':
		raw = z.startTagIn("script", "style")
	case 't':
		raw = z.startTagIn("textarea", "title")
	case 'x':
		raw = z.startTagIn("xmp")
	***REMOVED***
	if raw ***REMOVED***
		z.rawTag = strings.ToLower(string(z.buf[z.data.start:z.data.end]))
	***REMOVED***
	// Look for a self-closing token like "<br/>".
	if z.err == nil && z.buf[z.raw.end-2] == '/' ***REMOVED***
		return SelfClosingTagToken
	***REMOVED***
	return StartTagToken
***REMOVED***

// readTag reads the next tag token and its attributes. If saveAttr, those
// attributes are saved in z.attr, otherwise z.attr is set to an empty slice.
// The opening "<a" or "</a" has already been consumed, where 'a' means anything
// in [A-Za-z].
func (z *Tokenizer) readTag(saveAttr bool) ***REMOVED***
	z.attr = z.attr[:0]
	z.nAttrReturned = 0
	// Read the tag name and attribute key/value pairs.
	z.readTagName()
	if z.skipWhiteSpace(); z.err != nil ***REMOVED***
		return
	***REMOVED***
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil || c == '>' ***REMOVED***
			break
		***REMOVED***
		z.raw.end--
		z.readTagAttrKey()
		z.readTagAttrVal()
		// Save pendingAttr if saveAttr and that attribute has a non-empty key.
		if saveAttr && z.pendingAttr[0].start != z.pendingAttr[0].end ***REMOVED***
			z.attr = append(z.attr, z.pendingAttr)
		***REMOVED***
		if z.skipWhiteSpace(); z.err != nil ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// readTagName sets z.data to the "div" in "<div k=v>". The reader (z.raw.end)
// is positioned such that the first byte of the tag name (the "d" in "<div")
// has already been consumed.
func (z *Tokenizer) readTagName() ***REMOVED***
	z.data.start = z.raw.end - 1
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.data.end = z.raw.end
			return
		***REMOVED***
		switch c ***REMOVED***
		case ' ', '\n', '\r', '\t', '\f':
			z.data.end = z.raw.end - 1
			return
		case '/', '>':
			z.raw.end--
			z.data.end = z.raw.end
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readTagAttrKey sets z.pendingAttr[0] to the "k" in "<div k=v>".
// Precondition: z.err == nil.
func (z *Tokenizer) readTagAttrKey() ***REMOVED***
	z.pendingAttr[0].start = z.raw.end
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			z.pendingAttr[0].end = z.raw.end
			return
		***REMOVED***
		switch c ***REMOVED***
		case ' ', '\n', '\r', '\t', '\f', '/':
			z.pendingAttr[0].end = z.raw.end - 1
			return
		case '=', '>':
			z.raw.end--
			z.pendingAttr[0].end = z.raw.end
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// readTagAttrVal sets z.pendingAttr[1] to the "v" in "<div k=v>".
func (z *Tokenizer) readTagAttrVal() ***REMOVED***
	z.pendingAttr[1].start = z.raw.end
	z.pendingAttr[1].end = z.raw.end
	if z.skipWhiteSpace(); z.err != nil ***REMOVED***
		return
	***REMOVED***
	c := z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	if c != '=' ***REMOVED***
		z.raw.end--
		return
	***REMOVED***
	if z.skipWhiteSpace(); z.err != nil ***REMOVED***
		return
	***REMOVED***
	quote := z.readByte()
	if z.err != nil ***REMOVED***
		return
	***REMOVED***
	switch quote ***REMOVED***
	case '>':
		z.raw.end--
		return

	case '\'', '"':
		z.pendingAttr[1].start = z.raw.end
		for ***REMOVED***
			c := z.readByte()
			if z.err != nil ***REMOVED***
				z.pendingAttr[1].end = z.raw.end
				return
			***REMOVED***
			if c == quote ***REMOVED***
				z.pendingAttr[1].end = z.raw.end - 1
				return
			***REMOVED***
		***REMOVED***

	default:
		z.pendingAttr[1].start = z.raw.end - 1
		for ***REMOVED***
			c := z.readByte()
			if z.err != nil ***REMOVED***
				z.pendingAttr[1].end = z.raw.end
				return
			***REMOVED***
			switch c ***REMOVED***
			case ' ', '\n', '\r', '\t', '\f':
				z.pendingAttr[1].end = z.raw.end - 1
				return
			case '>':
				z.raw.end--
				z.pendingAttr[1].end = z.raw.end
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// Next scans the next token and returns its type.
func (z *Tokenizer) Next() TokenType ***REMOVED***
	z.raw.start = z.raw.end
	z.data.start = z.raw.end
	z.data.end = z.raw.end
	if z.err != nil ***REMOVED***
		z.tt = ErrorToken
		return z.tt
	***REMOVED***
	if z.rawTag != "" ***REMOVED***
		if z.rawTag == "plaintext" ***REMOVED***
			// Read everything up to EOF.
			for z.err == nil ***REMOVED***
				z.readByte()
			***REMOVED***
			z.data.end = z.raw.end
			z.textIsRaw = true
		***REMOVED*** else ***REMOVED***
			z.readRawOrRCDATA()
		***REMOVED***
		if z.data.end > z.data.start ***REMOVED***
			z.tt = TextToken
			z.convertNUL = true
			return z.tt
		***REMOVED***
	***REMOVED***
	z.textIsRaw = false
	z.convertNUL = false

loop:
	for ***REMOVED***
		c := z.readByte()
		if z.err != nil ***REMOVED***
			break loop
		***REMOVED***
		if c != '<' ***REMOVED***
			continue loop
		***REMOVED***

		// Check if the '<' we have just read is part of a tag, comment
		// or doctype. If not, it's part of the accumulated text token.
		c = z.readByte()
		if z.err != nil ***REMOVED***
			break loop
		***REMOVED***
		var tokenType TokenType
		switch ***REMOVED***
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			tokenType = StartTagToken
		case c == '/':
			tokenType = EndTagToken
		case c == '!' || c == '?':
			// We use CommentToken to mean any of "<!--actual comments-->",
			// "<!DOCTYPE declarations>" and "<?xml processing instructions?>".
			tokenType = CommentToken
		default:
			// Reconsume the current character.
			z.raw.end--
			continue
		***REMOVED***

		// We have a non-text token, but we might have accumulated some text
		// before that. If so, we return the text first, and return the non-
		// text token on the subsequent call to Next.
		if x := z.raw.end - len("<a"); z.raw.start < x ***REMOVED***
			z.raw.end = x
			z.data.end = x
			z.tt = TextToken
			return z.tt
		***REMOVED***
		switch tokenType ***REMOVED***
		case StartTagToken:
			z.tt = z.readStartTag()
			return z.tt
		case EndTagToken:
			c = z.readByte()
			if z.err != nil ***REMOVED***
				break loop
			***REMOVED***
			if c == '>' ***REMOVED***
				// "</>" does not generate a token at all. Generate an empty comment
				// to allow passthrough clients to pick up the data using Raw.
				// Reset the tokenizer state and start again.
				z.tt = CommentToken
				return z.tt
			***REMOVED***
			if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' ***REMOVED***
				z.readTag(false)
				if z.err != nil ***REMOVED***
					z.tt = ErrorToken
				***REMOVED*** else ***REMOVED***
					z.tt = EndTagToken
				***REMOVED***
				return z.tt
			***REMOVED***
			z.raw.end--
			z.readUntilCloseAngle()
			z.tt = CommentToken
			return z.tt
		case CommentToken:
			if c == '!' ***REMOVED***
				z.tt = z.readMarkupDeclaration()
				return z.tt
			***REMOVED***
			z.raw.end--
			z.readUntilCloseAngle()
			z.tt = CommentToken
			return z.tt
		***REMOVED***
	***REMOVED***
	if z.raw.start < z.raw.end ***REMOVED***
		z.data.end = z.raw.end
		z.tt = TextToken
		return z.tt
	***REMOVED***
	z.tt = ErrorToken
	return z.tt
***REMOVED***

// Raw returns the unmodified text of the current token. Calling Next, Token,
// Text, TagName or TagAttr may change the contents of the returned slice.
//
// The token stream's raw bytes partition the byte stream (up until an
// ErrorToken). There are no overlaps or gaps between two consecutive token's
// raw bytes. One implication is that the byte offset of the current token is
// the sum of the lengths of all previous tokens' raw bytes.
func (z *Tokenizer) Raw() []byte ***REMOVED***
	return z.buf[z.raw.start:z.raw.end]
***REMOVED***

// convertNewlines converts "\r" and "\r\n" in s to "\n".
// The conversion happens in place, but the resulting slice may be shorter.
func convertNewlines(s []byte) []byte ***REMOVED***
	for i, c := range s ***REMOVED***
		if c != '\r' ***REMOVED***
			continue
		***REMOVED***

		src := i + 1
		if src >= len(s) || s[src] != '\n' ***REMOVED***
			s[i] = '\n'
			continue
		***REMOVED***

		dst := i
		for src < len(s) ***REMOVED***
			if s[src] == '\r' ***REMOVED***
				if src+1 < len(s) && s[src+1] == '\n' ***REMOVED***
					src++
				***REMOVED***
				s[dst] = '\n'
			***REMOVED*** else ***REMOVED***
				s[dst] = s[src]
			***REMOVED***
			src++
			dst++
		***REMOVED***
		return s[:dst]
	***REMOVED***
	return s
***REMOVED***

var (
	nul         = []byte("\x00")
	replacement = []byte("\ufffd")
)

// Text returns the unescaped text of a text, comment or doctype token. The
// contents of the returned slice may change on the next call to Next.
func (z *Tokenizer) Text() []byte ***REMOVED***
	switch z.tt ***REMOVED***
	case TextToken, CommentToken, DoctypeToken:
		s := z.buf[z.data.start:z.data.end]
		z.data.start = z.raw.end
		z.data.end = z.raw.end
		s = convertNewlines(s)
		if (z.convertNUL || z.tt == CommentToken) && bytes.Contains(s, nul) ***REMOVED***
			s = bytes.Replace(s, nul, replacement, -1)
		***REMOVED***
		if !z.textIsRaw ***REMOVED***
			s = unescape(s, false)
		***REMOVED***
		return s
	***REMOVED***
	return nil
***REMOVED***

// TagName returns the lower-cased name of a tag token (the `img` out of
// `<IMG SRC="foo">`) and whether the tag has attributes.
// The contents of the returned slice may change on the next call to Next.
func (z *Tokenizer) TagName() (name []byte, hasAttr bool) ***REMOVED***
	if z.data.start < z.data.end ***REMOVED***
		switch z.tt ***REMOVED***
		case StartTagToken, EndTagToken, SelfClosingTagToken:
			s := z.buf[z.data.start:z.data.end]
			z.data.start = z.raw.end
			z.data.end = z.raw.end
			return lower(s), z.nAttrReturned < len(z.attr)
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// TagAttr returns the lower-cased key and unescaped value of the next unparsed
// attribute for the current tag token and whether there are more attributes.
// The contents of the returned slices may change on the next call to Next.
func (z *Tokenizer) TagAttr() (key, val []byte, moreAttr bool) ***REMOVED***
	if z.nAttrReturned < len(z.attr) ***REMOVED***
		switch z.tt ***REMOVED***
		case StartTagToken, SelfClosingTagToken:
			x := z.attr[z.nAttrReturned]
			z.nAttrReturned++
			key = z.buf[x[0].start:x[0].end]
			val = z.buf[x[1].start:x[1].end]
			return lower(key), unescape(convertNewlines(val), true), z.nAttrReturned < len(z.attr)
		***REMOVED***
	***REMOVED***
	return nil, nil, false
***REMOVED***

// Token returns the current Token. The result's Data and Attr values remain
// valid after subsequent Next calls.
func (z *Tokenizer) Token() Token ***REMOVED***
	t := Token***REMOVED***Type: z.tt***REMOVED***
	switch z.tt ***REMOVED***
	case TextToken, CommentToken, DoctypeToken:
		t.Data = string(z.Text())
	case StartTagToken, SelfClosingTagToken, EndTagToken:
		name, moreAttr := z.TagName()
		for moreAttr ***REMOVED***
			var key, val []byte
			key, val, moreAttr = z.TagAttr()
			t.Attr = append(t.Attr, Attribute***REMOVED***"", atom.String(key), string(val)***REMOVED***)
		***REMOVED***
		if a := atom.Lookup(name); a != 0 ***REMOVED***
			t.DataAtom, t.Data = a, a.String()
		***REMOVED*** else ***REMOVED***
			t.DataAtom, t.Data = 0, string(name)
		***REMOVED***
	***REMOVED***
	return t
***REMOVED***

// SetMaxBuf sets a limit on the amount of data buffered during tokenization.
// A value of 0 means unlimited.
func (z *Tokenizer) SetMaxBuf(n int) ***REMOVED***
	z.maxBuf = n
***REMOVED***

// NewTokenizer returns a new HTML Tokenizer for the given Reader.
// The input is assumed to be UTF-8 encoded.
func NewTokenizer(r io.Reader) *Tokenizer ***REMOVED***
	return NewTokenizerFragment(r, "")
***REMOVED***

// NewTokenizerFragment returns a new HTML Tokenizer for the given Reader, for
// tokenizing an existing element's InnerHTML fragment. contextTag is that
// element's tag, such as "div" or "iframe".
//
// For example, how the InnerHTML "a<b" is tokenized depends on whether it is
// for a <p> tag or a <script> tag.
//
// The input is assumed to be UTF-8 encoded.
func NewTokenizerFragment(r io.Reader, contextTag string) *Tokenizer ***REMOVED***
	z := &Tokenizer***REMOVED***
		r:   r,
		buf: make([]byte, 0, 4096),
	***REMOVED***
	if contextTag != "" ***REMOVED***
		switch s := strings.ToLower(contextTag); s ***REMOVED***
		case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "title", "textarea", "xmp":
			z.rawTag = s
		***REMOVED***
	***REMOVED***
	return z
***REMOVED***
