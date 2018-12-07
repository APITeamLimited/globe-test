// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"errors"
	"fmt"
	"io"
	"strings"

	a "golang.org/x/net/html/atom"
)

// A parser implements the HTML5 parsing algorithm:
// https://html.spec.whatwg.org/multipage/syntax.html#tree-construction
type parser struct ***REMOVED***
	// tokenizer provides the tokens for the parser.
	tokenizer *Tokenizer
	// tok is the most recently read token.
	tok Token
	// Self-closing tags like <hr/> are treated as start tags, except that
	// hasSelfClosingToken is set while they are being processed.
	hasSelfClosingToken bool
	// doc is the document root element.
	doc *Node
	// The stack of open elements (section 12.2.4.2) and active formatting
	// elements (section 12.2.4.3).
	oe, afe nodeStack
	// Element pointers (section 12.2.4.4).
	head, form *Node
	// Other parsing state flags (section 12.2.4.5).
	scripting, framesetOK bool
	// The stack of template insertion modes
	templateStack insertionModeStack
	// im is the current insertion mode.
	im insertionMode
	// originalIM is the insertion mode to go back to after completing a text
	// or inTableText insertion mode.
	originalIM insertionMode
	// fosterParenting is whether new elements should be inserted according to
	// the foster parenting rules (section 12.2.6.1).
	fosterParenting bool
	// quirks is whether the parser is operating in "quirks mode."
	quirks bool
	// fragment is whether the parser is parsing an HTML fragment.
	fragment bool
	// context is the context element when parsing an HTML fragment
	// (section 12.4).
	context *Node
***REMOVED***

func (p *parser) top() *Node ***REMOVED***
	if n := p.oe.top(); n != nil ***REMOVED***
		return n
	***REMOVED***
	return p.doc
***REMOVED***

// Stop tags for use in popUntil. These come from section 12.2.4.2.
var (
	defaultScopeStopTags = map[string][]a.Atom***REMOVED***
		"":     ***REMOVED***a.Applet, a.Caption, a.Html, a.Table, a.Td, a.Th, a.Marquee, a.Object, a.Template***REMOVED***,
		"math": ***REMOVED***a.AnnotationXml, a.Mi, a.Mn, a.Mo, a.Ms, a.Mtext***REMOVED***,
		"svg":  ***REMOVED***a.Desc, a.ForeignObject, a.Title***REMOVED***,
	***REMOVED***
)

type scope int

const (
	defaultScope scope = iota
	listItemScope
	buttonScope
	tableScope
	tableRowScope
	tableBodyScope
	selectScope
)

// popUntil pops the stack of open elements at the highest element whose tag
// is in matchTags, provided there is no higher element in the scope's stop
// tags (as defined in section 12.2.4.2). It returns whether or not there was
// such an element. If there was not, popUntil leaves the stack unchanged.
//
// For example, the set of stop tags for table scope is: "html", "table". If
// the stack was:
// ["html", "body", "font", "table", "b", "i", "u"]
// then popUntil(tableScope, "font") would return false, but
// popUntil(tableScope, "i") would return true and the stack would become:
// ["html", "body", "font", "table", "b"]
//
// If an element's tag is in both the stop tags and matchTags, then the stack
// will be popped and the function returns true (provided, of course, there was
// no higher element in the stack that was also in the stop tags). For example,
// popUntil(tableScope, "table") returns true and leaves:
// ["html", "body", "font"]
func (p *parser) popUntil(s scope, matchTags ...a.Atom) bool ***REMOVED***
	if i := p.indexOfElementInScope(s, matchTags...); i != -1 ***REMOVED***
		p.oe = p.oe[:i]
		return true
	***REMOVED***
	return false
***REMOVED***

// indexOfElementInScope returns the index in p.oe of the highest element whose
// tag is in matchTags that is in scope. If no matching element is in scope, it
// returns -1.
func (p *parser) indexOfElementInScope(s scope, matchTags ...a.Atom) int ***REMOVED***
	for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		tagAtom := p.oe[i].DataAtom
		if p.oe[i].Namespace == "" ***REMOVED***
			for _, t := range matchTags ***REMOVED***
				if t == tagAtom ***REMOVED***
					return i
				***REMOVED***
			***REMOVED***
			switch s ***REMOVED***
			case defaultScope:
				// No-op.
			case listItemScope:
				if tagAtom == a.Ol || tagAtom == a.Ul ***REMOVED***
					return -1
				***REMOVED***
			case buttonScope:
				if tagAtom == a.Button ***REMOVED***
					return -1
				***REMOVED***
			case tableScope:
				if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template ***REMOVED***
					return -1
				***REMOVED***
			case selectScope:
				if tagAtom != a.Optgroup && tagAtom != a.Option ***REMOVED***
					return -1
				***REMOVED***
			default:
				panic("unreachable")
			***REMOVED***
		***REMOVED***
		switch s ***REMOVED***
		case defaultScope, listItemScope, buttonScope:
			for _, t := range defaultScopeStopTags[p.oe[i].Namespace] ***REMOVED***
				if t == tagAtom ***REMOVED***
					return -1
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return -1
***REMOVED***

// elementInScope is like popUntil, except that it doesn't modify the stack of
// open elements.
func (p *parser) elementInScope(s scope, matchTags ...a.Atom) bool ***REMOVED***
	return p.indexOfElementInScope(s, matchTags...) != -1
***REMOVED***

// clearStackToContext pops elements off the stack of open elements until a
// scope-defined element is found.
func (p *parser) clearStackToContext(s scope) ***REMOVED***
	for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		tagAtom := p.oe[i].DataAtom
		switch s ***REMOVED***
		case tableScope:
			if tagAtom == a.Html || tagAtom == a.Table || tagAtom == a.Template ***REMOVED***
				p.oe = p.oe[:i+1]
				return
			***REMOVED***
		case tableRowScope:
			if tagAtom == a.Html || tagAtom == a.Tr || tagAtom == a.Template ***REMOVED***
				p.oe = p.oe[:i+1]
				return
			***REMOVED***
		case tableBodyScope:
			if tagAtom == a.Html || tagAtom == a.Tbody || tagAtom == a.Tfoot || tagAtom == a.Thead || tagAtom == a.Template ***REMOVED***
				p.oe = p.oe[:i+1]
				return
			***REMOVED***
		default:
			panic("unreachable")
		***REMOVED***
	***REMOVED***
***REMOVED***

// generateImpliedEndTags pops nodes off the stack of open elements as long as
// the top node has a tag name of dd, dt, li, optgroup, option, p, rb, rp, rt or rtc.
// If exceptions are specified, nodes with that name will not be popped off.
func (p *parser) generateImpliedEndTags(exceptions ...string) ***REMOVED***
	var i int
loop:
	for i = len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		n := p.oe[i]
		if n.Type == ElementNode ***REMOVED***
			switch n.DataAtom ***REMOVED***
			case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc:
				for _, except := range exceptions ***REMOVED***
					if n.Data == except ***REMOVED***
						break loop
					***REMOVED***
				***REMOVED***
				continue
			***REMOVED***
		***REMOVED***
		break
	***REMOVED***

	p.oe = p.oe[:i+1]
***REMOVED***

// addChild adds a child node n to the top element, and pushes n onto the stack
// of open elements if it is an element node.
func (p *parser) addChild(n *Node) ***REMOVED***
	if p.shouldFosterParent() ***REMOVED***
		p.fosterParent(n)
	***REMOVED*** else ***REMOVED***
		p.top().AppendChild(n)
	***REMOVED***

	if n.Type == ElementNode ***REMOVED***
		p.oe = append(p.oe, n)
	***REMOVED***
***REMOVED***

// shouldFosterParent returns whether the next node to be added should be
// foster parented.
func (p *parser) shouldFosterParent() bool ***REMOVED***
	if p.fosterParenting ***REMOVED***
		switch p.top().DataAtom ***REMOVED***
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// fosterParent adds a child node according to the foster parenting rules.
// Section 12.2.6.1, "foster parenting".
func (p *parser) fosterParent(n *Node) ***REMOVED***
	var table, parent, prev, template *Node
	var i int
	for i = len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		if p.oe[i].DataAtom == a.Table ***REMOVED***
			table = p.oe[i]
			break
		***REMOVED***
	***REMOVED***

	var j int
	for j = len(p.oe) - 1; j >= 0; j-- ***REMOVED***
		if p.oe[j].DataAtom == a.Template ***REMOVED***
			template = p.oe[j]
			break
		***REMOVED***
	***REMOVED***

	if template != nil && (table == nil || j > i) ***REMOVED***
		template.AppendChild(n)
		return
	***REMOVED***

	if table == nil ***REMOVED***
		// The foster parent is the html element.
		parent = p.oe[0]
	***REMOVED*** else ***REMOVED***
		parent = table.Parent
	***REMOVED***
	if parent == nil ***REMOVED***
		parent = p.oe[i-1]
	***REMOVED***

	if table != nil ***REMOVED***
		prev = table.PrevSibling
	***REMOVED*** else ***REMOVED***
		prev = parent.LastChild
	***REMOVED***
	if prev != nil && prev.Type == TextNode && n.Type == TextNode ***REMOVED***
		prev.Data += n.Data
		return
	***REMOVED***

	parent.InsertBefore(n, table)
***REMOVED***

// addText adds text to the preceding node if it is a text node, or else it
// calls addChild with a new text node.
func (p *parser) addText(text string) ***REMOVED***
	if text == "" ***REMOVED***
		return
	***REMOVED***

	if p.shouldFosterParent() ***REMOVED***
		p.fosterParent(&Node***REMOVED***
			Type: TextNode,
			Data: text,
		***REMOVED***)
		return
	***REMOVED***

	t := p.top()
	if n := t.LastChild; n != nil && n.Type == TextNode ***REMOVED***
		n.Data += text
		return
	***REMOVED***
	p.addChild(&Node***REMOVED***
		Type: TextNode,
		Data: text,
	***REMOVED***)
***REMOVED***

// addElement adds a child element based on the current token.
func (p *parser) addElement() ***REMOVED***
	p.addChild(&Node***REMOVED***
		Type:     ElementNode,
		DataAtom: p.tok.DataAtom,
		Data:     p.tok.Data,
		Attr:     p.tok.Attr,
	***REMOVED***)
***REMOVED***

// Section 12.2.4.3.
func (p *parser) addFormattingElement() ***REMOVED***
	tagAtom, attr := p.tok.DataAtom, p.tok.Attr
	p.addElement()

	// Implement the Noah's Ark clause, but with three per family instead of two.
	identicalElements := 0
findIdenticalElements:
	for i := len(p.afe) - 1; i >= 0; i-- ***REMOVED***
		n := p.afe[i]
		if n.Type == scopeMarkerNode ***REMOVED***
			break
		***REMOVED***
		if n.Type != ElementNode ***REMOVED***
			continue
		***REMOVED***
		if n.Namespace != "" ***REMOVED***
			continue
		***REMOVED***
		if n.DataAtom != tagAtom ***REMOVED***
			continue
		***REMOVED***
		if len(n.Attr) != len(attr) ***REMOVED***
			continue
		***REMOVED***
	compareAttributes:
		for _, t0 := range n.Attr ***REMOVED***
			for _, t1 := range attr ***REMOVED***
				if t0.Key == t1.Key && t0.Namespace == t1.Namespace && t0.Val == t1.Val ***REMOVED***
					// Found a match for this attribute, continue with the next attribute.
					continue compareAttributes
				***REMOVED***
			***REMOVED***
			// If we get here, there is no attribute that matches a.
			// Therefore the element is not identical to the new one.
			continue findIdenticalElements
		***REMOVED***

		identicalElements++
		if identicalElements >= 3 ***REMOVED***
			p.afe.remove(n)
		***REMOVED***
	***REMOVED***

	p.afe = append(p.afe, p.top())
***REMOVED***

// Section 12.2.4.3.
func (p *parser) clearActiveFormattingElements() ***REMOVED***
	for ***REMOVED***
		n := p.afe.pop()
		if len(p.afe) == 0 || n.Type == scopeMarkerNode ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

// Section 12.2.4.3.
func (p *parser) reconstructActiveFormattingElements() ***REMOVED***
	n := p.afe.top()
	if n == nil ***REMOVED***
		return
	***REMOVED***
	if n.Type == scopeMarkerNode || p.oe.index(n) != -1 ***REMOVED***
		return
	***REMOVED***
	i := len(p.afe) - 1
	for n.Type != scopeMarkerNode && p.oe.index(n) == -1 ***REMOVED***
		if i == 0 ***REMOVED***
			i = -1
			break
		***REMOVED***
		i--
		n = p.afe[i]
	***REMOVED***
	for ***REMOVED***
		i++
		clone := p.afe[i].clone()
		p.addChild(clone)
		p.afe[i] = clone
		if i == len(p.afe)-1 ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// Section 12.2.5.
func (p *parser) acknowledgeSelfClosingTag() ***REMOVED***
	p.hasSelfClosingToken = false
***REMOVED***

// An insertion mode (section 12.2.4.1) is the state transition function from
// a particular state in the HTML5 parser's state machine. It updates the
// parser's fields depending on parser.tok (where ErrorToken means EOF).
// It returns whether the token was consumed.
type insertionMode func(*parser) bool

// setOriginalIM sets the insertion mode to return to after completing a text or
// inTableText insertion mode.
// Section 12.2.4.1, "using the rules for".
func (p *parser) setOriginalIM() ***REMOVED***
	if p.originalIM != nil ***REMOVED***
		panic("html: bad parser state: originalIM was set twice")
	***REMOVED***
	p.originalIM = p.im
***REMOVED***

// Section 12.2.4.1, "reset the insertion mode".
func (p *parser) resetInsertionMode() ***REMOVED***
	for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		n := p.oe[i]
		last := i == 0
		if last && p.context != nil ***REMOVED***
			n = p.context
		***REMOVED***

		switch n.DataAtom ***REMOVED***
		case a.Select:
			if !last ***REMOVED***
				for ancestor, first := n, p.oe[0]; ancestor != first; ***REMOVED***
					if ancestor == first ***REMOVED***
						break
					***REMOVED***
					ancestor = p.oe[p.oe.index(ancestor)-1]
					switch ancestor.DataAtom ***REMOVED***
					case a.Template:
						p.im = inSelectIM
						return
					case a.Table:
						p.im = inSelectInTableIM
						return
					***REMOVED***
				***REMOVED***
			***REMOVED***
			p.im = inSelectIM
		case a.Td, a.Th:
			// TODO: remove this divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.im = inCellIM
		case a.Tr:
			p.im = inRowIM
		case a.Tbody, a.Thead, a.Tfoot:
			p.im = inTableBodyIM
		case a.Caption:
			p.im = inCaptionIM
		case a.Colgroup:
			p.im = inColumnGroupIM
		case a.Table:
			p.im = inTableIM
		case a.Template:
			// TODO: remove this divergence from the HTML5 spec.
			if n.Namespace != "" ***REMOVED***
				continue
			***REMOVED***
			p.im = p.templateStack.top()
		case a.Head:
			// TODO: remove this divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.im = inHeadIM
		case a.Body:
			p.im = inBodyIM
		case a.Frameset:
			p.im = inFramesetIM
		case a.Html:
			if p.head == nil ***REMOVED***
				p.im = beforeHeadIM
			***REMOVED*** else ***REMOVED***
				p.im = afterHeadIM
			***REMOVED***
		default:
			if last ***REMOVED***
				p.im = inBodyIM
				return
			***REMOVED***
			continue
		***REMOVED***
		return
	***REMOVED***
***REMOVED***

const whitespace = " \t\r\n\f"

// Section 12.2.6.4.1.
func initialIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 ***REMOVED***
			// It was all whitespace, so ignore it.
			return true
		***REMOVED***
	case CommentToken:
		p.doc.AppendChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		n, quirks := parseDoctype(p.tok.Data)
		p.doc.AppendChild(n)
		p.quirks = quirks
		p.im = beforeHTMLIM
		return true
	***REMOVED***
	p.quirks = true
	p.im = beforeHTMLIM
	return false
***REMOVED***

// Section 12.2.6.4.2.
func beforeHTMLIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case DoctypeToken:
		// Ignore the token.
		return true
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 ***REMOVED***
			// It was all whitespace, so ignore it.
			return true
		***REMOVED***
	case StartTagToken:
		if p.tok.DataAtom == a.Html ***REMOVED***
			p.addElement()
			p.im = beforeHeadIM
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTagToken, a.Html, a.Html.String())
			return false
		default:
			// Ignore the token.
			return true
		***REMOVED***
	case CommentToken:
		p.doc.AppendChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	***REMOVED***
	p.parseImpliedToken(StartTagToken, a.Html, a.Html.String())
	return false
***REMOVED***

// Section 12.2.6.4.3.
func beforeHeadIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		p.tok.Data = strings.TrimLeft(p.tok.Data, whitespace)
		if len(p.tok.Data) == 0 ***REMOVED***
			// It was all whitespace, so ignore it.
			return true
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Head:
			p.addElement()
			p.head = p.top()
			p.im = inHeadIM
			return true
		case a.Html:
			return inBodyIM(p)
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Head, a.Body, a.Html, a.Br:
			p.parseImpliedToken(StartTagToken, a.Head, a.Head.String())
			return false
		default:
			// Ignore the token.
			return true
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		// Ignore the token.
		return true
	***REMOVED***

	p.parseImpliedToken(StartTagToken, a.Head, a.Head.String())
	return false
***REMOVED***

// Section 12.2.6.4.4.
func inHeadIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) ***REMOVED***
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" ***REMOVED***
				return true
			***REMOVED***
			p.tok.Data = s
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Base, a.Basefont, a.Bgsound, a.Command, a.Link, a.Meta:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Script, a.Title, a.Noscript, a.Noframes, a.Style:
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
			return true
		case a.Head:
			// Ignore the token.
			return true
		case a.Template:
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
			p.im = inTemplateIM
			p.templateStack = append(p.templateStack, inTemplateIM)
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Head:
			p.oe.pop()
			p.im = afterHeadIM
			return true
		case a.Body, a.Html, a.Br:
			p.parseImpliedToken(EndTagToken, a.Head, a.Head.String())
			return false
		case a.Template:
			if !p.oe.contains(a.Template) ***REMOVED***
				return true
			***REMOVED***
			// TODO: remove this divergence from the HTML5 spec.
			//
			// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
			p.generateImpliedEndTags()
			for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
				if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template ***REMOVED***
					p.oe = p.oe[:i]
					break
				***REMOVED***
			***REMOVED***
			p.clearActiveFormattingElements()
			p.templateStack.pop()
			p.resetInsertionMode()
			return true
		default:
			// Ignore the token.
			return true
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		// Ignore the token.
		return true
	***REMOVED***

	p.parseImpliedToken(EndTagToken, a.Head, a.Head.String())
	return false
***REMOVED***

// Section 12.2.6.4.6.
func afterHeadIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) ***REMOVED***
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" ***REMOVED***
				return true
			***REMOVED***
			p.tok.Data = s
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Body:
			p.addElement()
			p.framesetOK = false
			p.im = inBodyIM
			return true
		case a.Frameset:
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			p.oe = append(p.oe, p.head)
			defer p.oe.remove(p.head)
			return inHeadIM(p)
		case a.Head:
			// Ignore the token.
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Body, a.Html, a.Br:
			// Drop down to creating an implied <body> tag.
		case a.Template:
			return inHeadIM(p)
		default:
			// Ignore the token.
			return true
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		// Ignore the token.
		return true
	***REMOVED***

	p.parseImpliedToken(StartTagToken, a.Body, a.Body.String())
	p.framesetOK = true
	return false
***REMOVED***

// copyAttributes copies attributes of src not found on dst to dst.
func copyAttributes(dst *Node, src Token) ***REMOVED***
	if len(src.Attr) == 0 ***REMOVED***
		return
	***REMOVED***
	attr := map[string]string***REMOVED******REMOVED***
	for _, t := range dst.Attr ***REMOVED***
		attr[t.Key] = t.Val
	***REMOVED***
	for _, t := range src.Attr ***REMOVED***
		if _, ok := attr[t.Key]; !ok ***REMOVED***
			dst.Attr = append(dst.Attr, t)
			attr[t.Key] = t.Val
		***REMOVED***
	***REMOVED***
***REMOVED***

// Section 12.2.6.4.7.
func inBodyIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		d := p.tok.Data
		switch n := p.oe.top(); n.DataAtom ***REMOVED***
		case a.Pre, a.Listing:
			if n.FirstChild == nil ***REMOVED***
				// Ignore a newline at the start of a <pre> block.
				if d != "" && d[0] == '\r' ***REMOVED***
					d = d[1:]
				***REMOVED***
				if d != "" && d[0] == '\n' ***REMOVED***
					d = d[1:]
				***REMOVED***
			***REMOVED***
		***REMOVED***
		d = strings.Replace(d, "\x00", "", -1)
		if d == "" ***REMOVED***
			return true
		***REMOVED***
		p.reconstructActiveFormattingElements()
		p.addText(d)
		if p.framesetOK && strings.TrimLeft(d, whitespace) != "" ***REMOVED***
			// There were non-whitespace characters inserted.
			p.framesetOK = false
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			if p.oe.contains(a.Template) ***REMOVED***
				return true
			***REMOVED***
			copyAttributes(p.oe[0], p.tok)
		case a.Base, a.Basefont, a.Bgsound, a.Command, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Body:
			if p.oe.contains(a.Template) ***REMOVED***
				return true
			***REMOVED***
			if len(p.oe) >= 2 ***REMOVED***
				body := p.oe[1]
				if body.Type == ElementNode && body.DataAtom == a.Body ***REMOVED***
					p.framesetOK = false
					copyAttributes(body, p.tok)
				***REMOVED***
			***REMOVED***
		case a.Frameset:
			if !p.framesetOK || len(p.oe) < 2 || p.oe[1].DataAtom != a.Body ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
			body := p.oe[1]
			if body.Parent != nil ***REMOVED***
				body.Parent.RemoveChild(body)
			***REMOVED***
			p.oe = p.oe[:1]
			p.addElement()
			p.im = inFramesetIM
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Center, a.Details, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Menu, a.Nav, a.Ol, a.P, a.Section, a.Summary, a.Ul:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(buttonScope, a.P)
			switch n := p.top(); n.DataAtom ***REMOVED***
			case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
				p.oe.pop()
			***REMOVED***
			p.addElement()
		case a.Pre, a.Listing:
			p.popUntil(buttonScope, a.P)
			p.addElement()
			// The newline, if any, will be dealt with by the TextToken case.
			p.framesetOK = false
		case a.Form:
			if p.form != nil && !p.oe.contains(a.Template) ***REMOVED***
				// Ignore the token
				return true
			***REMOVED***
			p.popUntil(buttonScope, a.P)
			p.addElement()
			if !p.oe.contains(a.Template) ***REMOVED***
				p.form = p.top()
			***REMOVED***
		case a.Li:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
				node := p.oe[i]
				switch node.DataAtom ***REMOVED***
				case a.Li:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Dd, a.Dt:
			p.framesetOK = false
			for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
				node := p.oe[i]
				switch node.DataAtom ***REMOVED***
				case a.Dd, a.Dt:
					p.oe = p.oe[:i]
				case a.Address, a.Div, a.P:
					continue
				default:
					if !isSpecialElement(node) ***REMOVED***
						continue
					***REMOVED***
				***REMOVED***
				break
			***REMOVED***
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Plaintext:
			p.popUntil(buttonScope, a.P)
			p.addElement()
		case a.Button:
			p.popUntil(defaultScope, a.Button)
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
		case a.A:
			for i := len(p.afe) - 1; i >= 0 && p.afe[i].Type != scopeMarkerNode; i-- ***REMOVED***
				if n := p.afe[i]; n.Type == ElementNode && n.DataAtom == a.A ***REMOVED***
					p.inBodyEndTagFormatting(a.A)
					p.oe.remove(n)
					p.afe.remove(n)
					break
				***REMOVED***
			***REMOVED***
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.reconstructActiveFormattingElements()
			p.addFormattingElement()
		case a.Nobr:
			p.reconstructActiveFormattingElements()
			if p.elementInScope(defaultScope, a.Nobr) ***REMOVED***
				p.inBodyEndTagFormatting(a.Nobr)
				p.reconstructActiveFormattingElements()
			***REMOVED***
			p.addFormattingElement()
		case a.Applet, a.Marquee, a.Object:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.framesetOK = false
		case a.Table:
			if !p.quirks ***REMOVED***
				p.popUntil(buttonScope, a.P)
			***REMOVED***
			p.addElement()
			p.framesetOK = false
			p.im = inTableIM
			return true
		case a.Area, a.Br, a.Embed, a.Img, a.Input, a.Keygen, a.Wbr:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			if p.tok.DataAtom == a.Input ***REMOVED***
				for _, t := range p.tok.Attr ***REMOVED***
					if t.Key == "type" ***REMOVED***
						if strings.ToLower(t.Val) == "hidden" ***REMOVED***
							// Skip setting framesetOK = false
							return true
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			p.framesetOK = false
		case a.Param, a.Source, a.Track:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Hr:
			p.popUntil(buttonScope, a.P)
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			p.framesetOK = false
		case a.Image:
			p.tok.DataAtom = a.Img
			p.tok.Data = a.Img.String()
			return false
		case a.Isindex:
			if p.form != nil ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
			action := ""
			prompt := "This is a searchable index. Enter search keywords: "
			attr := []Attribute***REMOVED******REMOVED***Key: "name", Val: "isindex"***REMOVED******REMOVED***
			for _, t := range p.tok.Attr ***REMOVED***
				switch t.Key ***REMOVED***
				case "action":
					action = t.Val
				case "name":
					// Ignore the attribute.
				case "prompt":
					prompt = t.Val
				default:
					attr = append(attr, t)
				***REMOVED***
			***REMOVED***
			p.acknowledgeSelfClosingTag()
			p.popUntil(buttonScope, a.P)
			p.parseImpliedToken(StartTagToken, a.Form, a.Form.String())
			if p.form == nil ***REMOVED***
				// NOTE: The 'isindex' element has been removed,
				// and the 'template' element has not been designed to be
				// collaborative with the index element.
				//
				// Ignore the token.
				return true
			***REMOVED***
			if action != "" ***REMOVED***
				p.form.Attr = []Attribute***REMOVED******REMOVED***Key: "action", Val: action***REMOVED******REMOVED***
			***REMOVED***
			p.parseImpliedToken(StartTagToken, a.Hr, a.Hr.String())
			p.parseImpliedToken(StartTagToken, a.Label, a.Label.String())
			p.addText(prompt)
			p.addChild(&Node***REMOVED***
				Type:     ElementNode,
				DataAtom: a.Input,
				Data:     a.Input.String(),
				Attr:     attr,
			***REMOVED***)
			p.oe.pop()
			p.parseImpliedToken(EndTagToken, a.Label, a.Label.String())
			p.parseImpliedToken(StartTagToken, a.Hr, a.Hr.String())
			p.parseImpliedToken(EndTagToken, a.Form, a.Form.String())
		case a.Textarea:
			p.addElement()
			p.setOriginalIM()
			p.framesetOK = false
			p.im = textIM
		case a.Xmp:
			p.popUntil(buttonScope, a.P)
			p.reconstructActiveFormattingElements()
			p.framesetOK = false
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Iframe:
			p.framesetOK = false
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Noembed, a.Noscript:
			p.addElement()
			p.setOriginalIM()
			p.im = textIM
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectIM
			return true
		case a.Optgroup, a.Option:
			if p.top().DataAtom == a.Option ***REMOVED***
				p.oe.pop()
			***REMOVED***
			p.reconstructActiveFormattingElements()
			p.addElement()
		case a.Rb, a.Rtc:
			if p.elementInScope(defaultScope, a.Ruby) ***REMOVED***
				p.generateImpliedEndTags()
			***REMOVED***
			p.addElement()
		case a.Rp, a.Rt:
			if p.elementInScope(defaultScope, a.Ruby) ***REMOVED***
				p.generateImpliedEndTags("rtc")
			***REMOVED***
			p.addElement()
		case a.Math, a.Svg:
			p.reconstructActiveFormattingElements()
			if p.tok.DataAtom == a.Math ***REMOVED***
				adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
			***REMOVED*** else ***REMOVED***
				adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
			***REMOVED***
			adjustForeignAttributes(p.tok.Attr)
			p.addElement()
			p.top().Namespace = p.tok.Data
			if p.hasSelfClosingToken ***REMOVED***
				p.oe.pop()
				p.acknowledgeSelfClosingTag()
			***REMOVED***
			return true
		case a.Caption, a.Col, a.Colgroup, a.Frame, a.Head, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
		default:
			p.reconstructActiveFormattingElements()
			p.addElement()
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Body:
			if p.elementInScope(defaultScope, a.Body) ***REMOVED***
				p.im = afterBodyIM
			***REMOVED***
		case a.Html:
			if p.elementInScope(defaultScope, a.Body) ***REMOVED***
				p.parseImpliedToken(EndTagToken, a.Body, a.Body.String())
				return false
			***REMOVED***
			return true
		case a.Address, a.Article, a.Aside, a.Blockquote, a.Button, a.Center, a.Details, a.Dir, a.Div, a.Dl, a.Fieldset, a.Figcaption, a.Figure, a.Footer, a.Header, a.Hgroup, a.Listing, a.Menu, a.Nav, a.Ol, a.Pre, a.Section, a.Summary, a.Ul:
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.Form:
			if p.oe.contains(a.Template) ***REMOVED***
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if i == -1 ***REMOVED***
					// Ignore the token.
					return true
				***REMOVED***
				p.generateImpliedEndTags()
				if p.oe[i].DataAtom != a.Form ***REMOVED***
					// Ignore the token.
					return true
				***REMOVED***
				p.popUntil(defaultScope, a.Form)
			***REMOVED*** else ***REMOVED***
				node := p.form
				p.form = nil
				i := p.indexOfElementInScope(defaultScope, a.Form)
				if node == nil || i == -1 || p.oe[i] != node ***REMOVED***
					// Ignore the token.
					return true
				***REMOVED***
				p.generateImpliedEndTags()
				p.oe.remove(node)
			***REMOVED***
		case a.P:
			if !p.elementInScope(buttonScope, a.P) ***REMOVED***
				p.parseImpliedToken(StartTagToken, a.P, a.P.String())
			***REMOVED***
			p.popUntil(buttonScope, a.P)
		case a.Li:
			p.popUntil(listItemScope, a.Li)
		case a.Dd, a.Dt:
			p.popUntil(defaultScope, p.tok.DataAtom)
		case a.H1, a.H2, a.H3, a.H4, a.H5, a.H6:
			p.popUntil(defaultScope, a.H1, a.H2, a.H3, a.H4, a.H5, a.H6)
		case a.A, a.B, a.Big, a.Code, a.Em, a.Font, a.I, a.Nobr, a.S, a.Small, a.Strike, a.Strong, a.Tt, a.U:
			p.inBodyEndTagFormatting(p.tok.DataAtom)
		case a.Applet, a.Marquee, a.Object:
			if p.popUntil(defaultScope, p.tok.DataAtom) ***REMOVED***
				p.clearActiveFormattingElements()
			***REMOVED***
		case a.Br:
			p.tok.Type = StartTagToken
			return false
		case a.Template:
			return inHeadIM(p)
		default:
			p.inBodyEndTagOther(p.tok.DataAtom)
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case ErrorToken:
		// TODO: remove this divergence from the HTML5 spec.
		if len(p.templateStack) > 0 ***REMOVED***
			p.im = inTemplateIM
			return false
		***REMOVED*** else ***REMOVED***
			for _, e := range p.oe ***REMOVED***
				switch e.DataAtom ***REMOVED***
				case a.Dd, a.Dt, a.Li, a.Optgroup, a.Option, a.P, a.Rb, a.Rp, a.Rt, a.Rtc, a.Tbody, a.Td, a.Tfoot, a.Th,
					a.Thead, a.Tr, a.Body, a.Html:
				default:
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return true
***REMOVED***

func (p *parser) inBodyEndTagFormatting(tagAtom a.Atom) ***REMOVED***
	// This is the "adoption agency" algorithm, described at
	// https://html.spec.whatwg.org/multipage/syntax.html#adoptionAgency

	// TODO: this is a fairly literal line-by-line translation of that algorithm.
	// Once the code successfully parses the comprehensive test suite, we should
	// refactor this code to be more idiomatic.

	// Steps 1-4. The outer loop.
	for i := 0; i < 8; i++ ***REMOVED***
		// Step 5. Find the formatting element.
		var formattingElement *Node
		for j := len(p.afe) - 1; j >= 0; j-- ***REMOVED***
			if p.afe[j].Type == scopeMarkerNode ***REMOVED***
				break
			***REMOVED***
			if p.afe[j].DataAtom == tagAtom ***REMOVED***
				formattingElement = p.afe[j]
				break
			***REMOVED***
		***REMOVED***
		if formattingElement == nil ***REMOVED***
			p.inBodyEndTagOther(tagAtom)
			return
		***REMOVED***
		feIndex := p.oe.index(formattingElement)
		if feIndex == -1 ***REMOVED***
			p.afe.remove(formattingElement)
			return
		***REMOVED***
		if !p.elementInScope(defaultScope, tagAtom) ***REMOVED***
			// Ignore the tag.
			return
		***REMOVED***

		// Steps 9-10. Find the furthest block.
		var furthestBlock *Node
		for _, e := range p.oe[feIndex:] ***REMOVED***
			if isSpecialElement(e) ***REMOVED***
				furthestBlock = e
				break
			***REMOVED***
		***REMOVED***
		if furthestBlock == nil ***REMOVED***
			e := p.oe.pop()
			for e != formattingElement ***REMOVED***
				e = p.oe.pop()
			***REMOVED***
			p.afe.remove(e)
			return
		***REMOVED***

		// Steps 11-12. Find the common ancestor and bookmark node.
		commonAncestor := p.oe[feIndex-1]
		bookmark := p.afe.index(formattingElement)

		// Step 13. The inner loop. Find the lastNode to reparent.
		lastNode := furthestBlock
		node := furthestBlock
		x := p.oe.index(node)
		// Steps 13.1-13.2
		for j := 0; j < 3; j++ ***REMOVED***
			// Step 13.3.
			x--
			node = p.oe[x]
			// Step 13.4 - 13.5.
			if p.afe.index(node) == -1 ***REMOVED***
				p.oe.remove(node)
				continue
			***REMOVED***
			// Step 13.6.
			if node == formattingElement ***REMOVED***
				break
			***REMOVED***
			// Step 13.7.
			clone := node.clone()
			p.afe[p.afe.index(node)] = clone
			p.oe[p.oe.index(node)] = clone
			node = clone
			// Step 13.8.
			if lastNode == furthestBlock ***REMOVED***
				bookmark = p.afe.index(node) + 1
			***REMOVED***
			// Step 13.9.
			if lastNode.Parent != nil ***REMOVED***
				lastNode.Parent.RemoveChild(lastNode)
			***REMOVED***
			node.AppendChild(lastNode)
			// Step 13.10.
			lastNode = node
		***REMOVED***

		// Step 14. Reparent lastNode to the common ancestor,
		// or for misnested table nodes, to the foster parent.
		if lastNode.Parent != nil ***REMOVED***
			lastNode.Parent.RemoveChild(lastNode)
		***REMOVED***
		switch commonAncestor.DataAtom ***REMOVED***
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			p.fosterParent(lastNode)
		default:
			commonAncestor.AppendChild(lastNode)
		***REMOVED***

		// Steps 15-17. Reparent nodes from the furthest block's children
		// to a clone of the formatting element.
		clone := formattingElement.clone()
		reparentChildren(clone, furthestBlock)
		furthestBlock.AppendChild(clone)

		// Step 18. Fix up the list of active formatting elements.
		if oldLoc := p.afe.index(formattingElement); oldLoc != -1 && oldLoc < bookmark ***REMOVED***
			// Move the bookmark with the rest of the list.
			bookmark--
		***REMOVED***
		p.afe.remove(formattingElement)
		p.afe.insert(bookmark, clone)

		// Step 19. Fix up the stack of open elements.
		p.oe.remove(formattingElement)
		p.oe.insert(p.oe.index(furthestBlock)+1, clone)
	***REMOVED***
***REMOVED***

// inBodyEndTagOther performs the "any other end tag" algorithm for inBodyIM.
// "Any other end tag" handling from 12.2.6.5 The rules for parsing tokens in foreign content
// https://html.spec.whatwg.org/multipage/syntax.html#parsing-main-inforeign
func (p *parser) inBodyEndTagOther(tagAtom a.Atom) ***REMOVED***
	for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
		if p.oe[i].DataAtom == tagAtom ***REMOVED***
			p.oe = p.oe[:i]
			break
		***REMOVED***
		if isSpecialElement(p.oe[i]) ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
***REMOVED***

// Section 12.2.6.4.8.
func textIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case ErrorToken:
		p.oe.pop()
	case TextToken:
		d := p.tok.Data
		if n := p.oe.top(); n.DataAtom == a.Textarea && n.FirstChild == nil ***REMOVED***
			// Ignore a newline at the start of a <textarea> block.
			if d != "" && d[0] == '\r' ***REMOVED***
				d = d[1:]
			***REMOVED***
			if d != "" && d[0] == '\n' ***REMOVED***
				d = d[1:]
			***REMOVED***
		***REMOVED***
		if d == "" ***REMOVED***
			return true
		***REMOVED***
		p.addText(d)
		return true
	case EndTagToken:
		p.oe.pop()
	***REMOVED***
	p.im = p.originalIM
	p.originalIM = nil
	return p.tok.Type == EndTagToken
***REMOVED***

// Section 12.2.6.4.9.
func inTableIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		p.tok.Data = strings.Replace(p.tok.Data, "\x00", "", -1)
		switch p.oe.top().DataAtom ***REMOVED***
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if strings.Trim(p.tok.Data, whitespace) == "" ***REMOVED***
				p.addText(p.tok.Data)
				return true
			***REMOVED***
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Caption:
			p.clearStackToContext(tableScope)
			p.afe = append(p.afe, &scopeMarker)
			p.addElement()
			p.im = inCaptionIM
			return true
		case a.Colgroup:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inColumnGroupIM
			return true
		case a.Col:
			p.parseImpliedToken(StartTagToken, a.Colgroup, a.Colgroup.String())
			return false
		case a.Tbody, a.Tfoot, a.Thead:
			p.clearStackToContext(tableScope)
			p.addElement()
			p.im = inTableBodyIM
			return true
		case a.Td, a.Th, a.Tr:
			p.parseImpliedToken(StartTagToken, a.Tbody, a.Tbody.String())
			return false
		case a.Table:
			if p.popUntil(tableScope, a.Table) ***REMOVED***
				p.resetInsertionMode()
				return false
			***REMOVED***
			// Ignore the token.
			return true
		case a.Style, a.Script, a.Template:
			return inHeadIM(p)
		case a.Input:
			for _, t := range p.tok.Attr ***REMOVED***
				if t.Key == "type" && strings.ToLower(t.Val) == "hidden" ***REMOVED***
					p.addElement()
					p.oe.pop()
					return true
				***REMOVED***
			***REMOVED***
			// Otherwise drop down to the default action.
		case a.Form:
			if p.oe.contains(a.Template) || p.form != nil ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
			p.addElement()
			p.form = p.oe.pop()
		case a.Select:
			p.reconstructActiveFormattingElements()
			switch p.top().DataAtom ***REMOVED***
			case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
				p.fosterParenting = true
			***REMOVED***
			p.addElement()
			p.fosterParenting = false
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Table:
			if p.popUntil(tableScope, a.Table) ***REMOVED***
				p.resetInsertionMode()
				return true
			***REMOVED***
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
			return true
		case a.Template:
			return inHeadIM(p)
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		// Ignore the token.
		return true
	case ErrorToken:
		return inBodyIM(p)
	***REMOVED***

	p.fosterParenting = true
	defer func() ***REMOVED*** p.fosterParenting = false ***REMOVED***()

	return inBodyIM(p)
***REMOVED***

// Section 12.2.6.4.11.
func inCaptionIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Caption) ***REMOVED***
				p.clearActiveFormattingElements()
				p.im = inTableIM
				return false
			***REMOVED*** else ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Caption:
			if p.popUntil(tableScope, a.Caption) ***REMOVED***
				p.clearActiveFormattingElements()
				p.im = inTableIM
			***REMOVED***
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Caption) ***REMOVED***
				p.clearActiveFormattingElements()
				p.im = inTableIM
				return false
			***REMOVED*** else ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
		case a.Body, a.Col, a.Colgroup, a.Html, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			// Ignore the token.
			return true
		***REMOVED***
	***REMOVED***
	return inBodyIM(p)
***REMOVED***

// Section 12.2.6.4.12.
func inColumnGroupIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) < len(p.tok.Data) ***REMOVED***
			// Add the initial whitespace to the current node.
			p.addText(p.tok.Data[:len(p.tok.Data)-len(s)])
			if s == "" ***REMOVED***
				return true
			***REMOVED***
			p.tok.Data = s
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		// Ignore the token.
		return true
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Col:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
			return true
		case a.Template:
			return inHeadIM(p)
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Colgroup:
			if p.oe.top().DataAtom == a.Colgroup ***REMOVED***
				p.oe.pop()
				p.im = inTableIM
			***REMOVED***
			return true
		case a.Col:
			// Ignore the token.
			return true
		case a.Template:
			return inHeadIM(p)
		***REMOVED***
	case ErrorToken:
		return inBodyIM(p)
	***REMOVED***
	if p.oe.top().DataAtom != a.Colgroup ***REMOVED***
		return true
	***REMOVED***
	p.oe.pop()
	p.im = inTableIM
	return false
***REMOVED***

// Section 12.2.6.4.13.
func inTableBodyIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Tr:
			p.clearStackToContext(tableBodyScope)
			p.addElement()
			p.im = inRowIM
			return true
		case a.Td, a.Th:
			p.parseImpliedToken(StartTagToken, a.Tr, a.Tr.String())
			return false
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) ***REMOVED***
				p.im = inTableIM
				return false
			***REMOVED***
			// Ignore the token.
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) ***REMOVED***
				p.clearStackToContext(tableBodyScope)
				p.oe.pop()
				p.im = inTableIM
			***REMOVED***
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Tbody, a.Thead, a.Tfoot) ***REMOVED***
				p.im = inTableIM
				return false
			***REMOVED***
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th, a.Tr:
			// Ignore the token.
			return true
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	***REMOVED***

	return inTableIM(p)
***REMOVED***

// Section 12.2.6.4.14.
func inRowIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Td, a.Th:
			p.clearStackToContext(tableRowScope)
			p.addElement()
			p.afe = append(p.afe, &scopeMarker)
			p.im = inCellIM
			return true
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Tr) ***REMOVED***
				p.im = inTableBodyIM
				return false
			***REMOVED***
			// Ignore the token.
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Tr:
			if p.popUntil(tableScope, a.Tr) ***REMOVED***
				p.im = inTableBodyIM
				return true
			***REMOVED***
			// Ignore the token.
			return true
		case a.Table:
			if p.popUntil(tableScope, a.Tr) ***REMOVED***
				p.im = inTableBodyIM
				return false
			***REMOVED***
			// Ignore the token.
			return true
		case a.Tbody, a.Tfoot, a.Thead:
			if p.elementInScope(tableScope, p.tok.DataAtom) ***REMOVED***
				p.parseImpliedToken(EndTagToken, a.Tr, a.Tr.String())
				return false
			***REMOVED***
			// Ignore the token.
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html, a.Td, a.Th:
			// Ignore the token.
			return true
		***REMOVED***
	***REMOVED***

	return inTableIM(p)
***REMOVED***

// Section 12.2.6.4.15.
func inCellIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Caption, a.Col, a.Colgroup, a.Tbody, a.Td, a.Tfoot, a.Th, a.Thead, a.Tr:
			if p.popUntil(tableScope, a.Td, a.Th) ***REMOVED***
				// Close the cell and reprocess.
				p.clearActiveFormattingElements()
				p.im = inRowIM
				return false
			***REMOVED***
			// Ignore the token.
			return true
		case a.Select:
			p.reconstructActiveFormattingElements()
			p.addElement()
			p.framesetOK = false
			p.im = inSelectInTableIM
			return true
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Td, a.Th:
			if !p.popUntil(tableScope, p.tok.DataAtom) ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
			p.clearActiveFormattingElements()
			p.im = inRowIM
			return true
		case a.Body, a.Caption, a.Col, a.Colgroup, a.Html:
			// Ignore the token.
			return true
		case a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr:
			if !p.elementInScope(tableScope, p.tok.DataAtom) ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
			// Close the cell and reprocess.
			p.popUntil(tableScope, a.Td, a.Th)
			p.clearActiveFormattingElements()
			p.im = inRowIM
			return false
		***REMOVED***
	***REMOVED***
	return inBodyIM(p)
***REMOVED***

// Section 12.2.6.4.16.
func inSelectIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		p.addText(strings.Replace(p.tok.Data, "\x00", "", -1))
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Option:
			if p.top().DataAtom == a.Option ***REMOVED***
				p.oe.pop()
			***REMOVED***
			p.addElement()
		case a.Optgroup:
			if p.top().DataAtom == a.Option ***REMOVED***
				p.oe.pop()
			***REMOVED***
			if p.top().DataAtom == a.Optgroup ***REMOVED***
				p.oe.pop()
			***REMOVED***
			p.addElement()
		case a.Select:
			p.tok.Type = EndTagToken
			return false
		case a.Input, a.Keygen, a.Textarea:
			if p.elementInScope(selectScope, a.Select) ***REMOVED***
				p.parseImpliedToken(EndTagToken, a.Select, a.Select.String())
				return false
			***REMOVED***
			// In order to properly ignore <textarea>, we need to change the tokenizer mode.
			p.tokenizer.NextIsNotRawText()
			// Ignore the token.
			return true
		case a.Script, a.Template:
			return inHeadIM(p)
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Option:
			if p.top().DataAtom == a.Option ***REMOVED***
				p.oe.pop()
			***REMOVED***
		case a.Optgroup:
			i := len(p.oe) - 1
			if p.oe[i].DataAtom == a.Option ***REMOVED***
				i--
			***REMOVED***
			if p.oe[i].DataAtom == a.Optgroup ***REMOVED***
				p.oe = p.oe[:i]
			***REMOVED***
		case a.Select:
			if p.popUntil(selectScope, a.Select) ***REMOVED***
				p.resetInsertionMode()
			***REMOVED***
		case a.Template:
			return inHeadIM(p)
		***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case DoctypeToken:
		// Ignore the token.
		return true
	case ErrorToken:
		return inBodyIM(p)
	***REMOVED***

	return true
***REMOVED***

// Section 12.2.6.4.17.
func inSelectInTableIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case StartTagToken, EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Caption, a.Table, a.Tbody, a.Tfoot, a.Thead, a.Tr, a.Td, a.Th:
			if p.tok.Type == StartTagToken || p.elementInScope(tableScope, p.tok.DataAtom) ***REMOVED***
				p.parseImpliedToken(EndTagToken, a.Select, a.Select.String())
				return false
			***REMOVED*** else ***REMOVED***
				// Ignore the token.
				return true
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return inSelectIM(p)
***REMOVED***

// Section 12.2.6.4.18.
func inTemplateIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken, CommentToken, DoctypeToken:
		return inBodyIM(p)
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Base, a.Basefont, a.Bgsound, a.Link, a.Meta, a.Noframes, a.Script, a.Style, a.Template, a.Title:
			return inHeadIM(p)
		case a.Caption, a.Colgroup, a.Tbody, a.Tfoot, a.Thead:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableIM)
			p.im = inTableIM
			return false
		case a.Col:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inColumnGroupIM)
			p.im = inColumnGroupIM
			return false
		case a.Tr:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inTableBodyIM)
			p.im = inTableBodyIM
			return false
		case a.Td, a.Th:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inRowIM)
			p.im = inRowIM
			return false
		default:
			p.templateStack.pop()
			p.templateStack = append(p.templateStack, inBodyIM)
			p.im = inBodyIM
			return false
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Template:
			return inHeadIM(p)
		default:
			// Ignore the token.
			return true
		***REMOVED***
	case ErrorToken:
		if !p.oe.contains(a.Template) ***REMOVED***
			// Ignore the token.
			return true
		***REMOVED***
		// TODO: remove this divergence from the HTML5 spec.
		//
		// See https://bugs.chromium.org/p/chromium/issues/detail?id=829668
		p.generateImpliedEndTags()
		for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
			if n := p.oe[i]; n.Namespace == "" && n.DataAtom == a.Template ***REMOVED***
				p.oe = p.oe[:i]
				break
			***REMOVED***
		***REMOVED***
		p.clearActiveFormattingElements()
		p.templateStack.pop()
		p.resetInsertionMode()
		return false
	***REMOVED***
	return false
***REMOVED***

// Section 12.2.6.4.19.
func afterBodyIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case ErrorToken:
		// Stop parsing.
		return true
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) == 0 ***REMOVED***
			// It was all whitespace.
			return inBodyIM(p)
		***REMOVED***
	case StartTagToken:
		if p.tok.DataAtom == a.Html ***REMOVED***
			return inBodyIM(p)
		***REMOVED***
	case EndTagToken:
		if p.tok.DataAtom == a.Html ***REMOVED***
			if !p.fragment ***REMOVED***
				p.im = afterAfterBodyIM
			***REMOVED***
			return true
		***REMOVED***
	case CommentToken:
		// The comment is attached to the <html> element.
		if len(p.oe) < 1 || p.oe[0].DataAtom != a.Html ***REMOVED***
			panic("html: bad parser state: <html> element not found, in the after-body insertion mode")
		***REMOVED***
		p.oe[0].AppendChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	***REMOVED***
	p.im = inBodyIM
	return false
***REMOVED***

// Section 12.2.6.4.20.
func inFramesetIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case TextToken:
		// Ignore all text but whitespace.
		s := strings.Map(func(c rune) rune ***REMOVED***
			switch c ***REMOVED***
			case ' ', '\t', '\n', '\f', '\r':
				return c
			***REMOVED***
			return -1
		***REMOVED***, p.tok.Data)
		if s != "" ***REMOVED***
			p.addText(s)
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Frameset:
			p.addElement()
		case a.Frame:
			p.addElement()
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		case a.Noframes:
			return inHeadIM(p)
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Frameset:
			if p.oe.top().DataAtom != a.Html ***REMOVED***
				p.oe.pop()
				if p.oe.top().DataAtom != a.Frameset ***REMOVED***
					p.im = afterFramesetIM
					return true
				***REMOVED***
			***REMOVED***
		***REMOVED***
	default:
		// Ignore the token.
	***REMOVED***
	return true
***REMOVED***

// Section 12.2.6.4.21.
func afterFramesetIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case TextToken:
		// Ignore all text but whitespace.
		s := strings.Map(func(c rune) rune ***REMOVED***
			switch c ***REMOVED***
			case ' ', '\t', '\n', '\f', '\r':
				return c
			***REMOVED***
			return -1
		***REMOVED***, p.tok.Data)
		if s != "" ***REMOVED***
			p.addText(s)
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		***REMOVED***
	case EndTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			p.im = afterAfterFramesetIM
			return true
		***REMOVED***
	default:
		// Ignore the token.
	***REMOVED***
	return true
***REMOVED***

// Section 12.2.6.4.22.
func afterAfterBodyIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case ErrorToken:
		// Stop parsing.
		return true
	case TextToken:
		s := strings.TrimLeft(p.tok.Data, whitespace)
		if len(s) == 0 ***REMOVED***
			// It was all whitespace.
			return inBodyIM(p)
		***REMOVED***
	case StartTagToken:
		if p.tok.DataAtom == a.Html ***REMOVED***
			return inBodyIM(p)
		***REMOVED***
	case CommentToken:
		p.doc.AppendChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
		return true
	case DoctypeToken:
		return inBodyIM(p)
	***REMOVED***
	p.im = inBodyIM
	return false
***REMOVED***

// Section 12.2.6.4.23.
func afterAfterFramesetIM(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case CommentToken:
		p.doc.AppendChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case TextToken:
		// Ignore all text but whitespace.
		s := strings.Map(func(c rune) rune ***REMOVED***
			switch c ***REMOVED***
			case ' ', '\t', '\n', '\f', '\r':
				return c
			***REMOVED***
			return -1
		***REMOVED***, p.tok.Data)
		if s != "" ***REMOVED***
			p.tok.Data = s
			return inBodyIM(p)
		***REMOVED***
	case StartTagToken:
		switch p.tok.DataAtom ***REMOVED***
		case a.Html:
			return inBodyIM(p)
		case a.Noframes:
			return inHeadIM(p)
		***REMOVED***
	case DoctypeToken:
		return inBodyIM(p)
	default:
		// Ignore the token.
	***REMOVED***
	return true
***REMOVED***

const whitespaceOrNUL = whitespace + "\x00"

// Section 12.2.6.5
func parseForeignContent(p *parser) bool ***REMOVED***
	switch p.tok.Type ***REMOVED***
	case TextToken:
		if p.framesetOK ***REMOVED***
			p.framesetOK = strings.TrimLeft(p.tok.Data, whitespaceOrNUL) == ""
		***REMOVED***
		p.tok.Data = strings.Replace(p.tok.Data, "\x00", "\ufffd", -1)
		p.addText(p.tok.Data)
	case CommentToken:
		p.addChild(&Node***REMOVED***
			Type: CommentNode,
			Data: p.tok.Data,
		***REMOVED***)
	case StartTagToken:
		b := breakout[p.tok.Data]
		if p.tok.DataAtom == a.Font ***REMOVED***
		loop:
			for _, attr := range p.tok.Attr ***REMOVED***
				switch attr.Key ***REMOVED***
				case "color", "face", "size":
					b = true
					break loop
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if b ***REMOVED***
			for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
				n := p.oe[i]
				if n.Namespace == "" || htmlIntegrationPoint(n) || mathMLTextIntegrationPoint(n) ***REMOVED***
					p.oe = p.oe[:i+1]
					break
				***REMOVED***
			***REMOVED***
			return false
		***REMOVED***
		switch p.top().Namespace ***REMOVED***
		case "math":
			adjustAttributeNames(p.tok.Attr, mathMLAttributeAdjustments)
		case "svg":
			// Adjust SVG tag names. The tokenizer lower-cases tag names, but
			// SVG wants e.g. "foreignObject" with a capital second "O".
			if x := svgTagNameAdjustments[p.tok.Data]; x != "" ***REMOVED***
				p.tok.DataAtom = a.Lookup([]byte(x))
				p.tok.Data = x
			***REMOVED***
			adjustAttributeNames(p.tok.Attr, svgAttributeAdjustments)
		default:
			panic("html: bad parser state: unexpected namespace")
		***REMOVED***
		adjustForeignAttributes(p.tok.Attr)
		namespace := p.top().Namespace
		p.addElement()
		p.top().Namespace = namespace
		if namespace != "" ***REMOVED***
			// Don't let the tokenizer go into raw text mode in foreign content
			// (e.g. in an SVG <title> tag).
			p.tokenizer.NextIsNotRawText()
		***REMOVED***
		if p.hasSelfClosingToken ***REMOVED***
			p.oe.pop()
			p.acknowledgeSelfClosingTag()
		***REMOVED***
	case EndTagToken:
		for i := len(p.oe) - 1; i >= 0; i-- ***REMOVED***
			if p.oe[i].Namespace == "" ***REMOVED***
				return p.im(p)
			***REMOVED***
			if strings.EqualFold(p.oe[i].Data, p.tok.Data) ***REMOVED***
				p.oe = p.oe[:i]
				break
			***REMOVED***
		***REMOVED***
		return true
	default:
		// Ignore the token.
	***REMOVED***
	return true
***REMOVED***

// Section 12.2.6.
func (p *parser) inForeignContent() bool ***REMOVED***
	if len(p.oe) == 0 ***REMOVED***
		return false
	***REMOVED***
	n := p.oe[len(p.oe)-1]
	if n.Namespace == "" ***REMOVED***
		return false
	***REMOVED***
	if mathMLTextIntegrationPoint(n) ***REMOVED***
		if p.tok.Type == StartTagToken && p.tok.DataAtom != a.Mglyph && p.tok.DataAtom != a.Malignmark ***REMOVED***
			return false
		***REMOVED***
		if p.tok.Type == TextToken ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	if n.Namespace == "math" && n.DataAtom == a.AnnotationXml && p.tok.Type == StartTagToken && p.tok.DataAtom == a.Svg ***REMOVED***
		return false
	***REMOVED***
	if htmlIntegrationPoint(n) && (p.tok.Type == StartTagToken || p.tok.Type == TextToken) ***REMOVED***
		return false
	***REMOVED***
	if p.tok.Type == ErrorToken ***REMOVED***
		return false
	***REMOVED***
	return true
***REMOVED***

// parseImpliedToken parses a token as though it had appeared in the parser's
// input.
func (p *parser) parseImpliedToken(t TokenType, dataAtom a.Atom, data string) ***REMOVED***
	realToken, selfClosing := p.tok, p.hasSelfClosingToken
	p.tok = Token***REMOVED***
		Type:     t,
		DataAtom: dataAtom,
		Data:     data,
	***REMOVED***
	p.hasSelfClosingToken = false
	p.parseCurrentToken()
	p.tok, p.hasSelfClosingToken = realToken, selfClosing
***REMOVED***

// parseCurrentToken runs the current token through the parsing routines
// until it is consumed.
func (p *parser) parseCurrentToken() ***REMOVED***
	if p.tok.Type == SelfClosingTagToken ***REMOVED***
		p.hasSelfClosingToken = true
		p.tok.Type = StartTagToken
	***REMOVED***

	consumed := false
	for !consumed ***REMOVED***
		if p.inForeignContent() ***REMOVED***
			consumed = parseForeignContent(p)
		***REMOVED*** else ***REMOVED***
			consumed = p.im(p)
		***REMOVED***
	***REMOVED***

	if p.hasSelfClosingToken ***REMOVED***
		// This is a parse error, but ignore it.
		p.hasSelfClosingToken = false
	***REMOVED***
***REMOVED***

func (p *parser) parse() error ***REMOVED***
	// Iterate until EOF. Any other error will cause an early return.
	var err error
	for err != io.EOF ***REMOVED***
		// CDATA sections are allowed only in foreign content.
		n := p.oe.top()
		p.tokenizer.AllowCDATA(n != nil && n.Namespace != "")
		// Read and parse the next token.
		p.tokenizer.Next()
		p.tok = p.tokenizer.Token()
		if p.tok.Type == ErrorToken ***REMOVED***
			err = p.tokenizer.Err()
			if err != nil && err != io.EOF ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		p.parseCurrentToken()
	***REMOVED***
	return nil
***REMOVED***

// Parse returns the parse tree for the HTML from the given Reader.
//
// It implements the HTML5 parsing algorithm
// (https://html.spec.whatwg.org/multipage/syntax.html#tree-construction),
// which is very complicated. The resultant tree can contain implicitly created
// nodes that have no explicit <tag> listed in r's data, and nodes' parents can
// differ from the nesting implied by a naive processing of start and end
// <tag>s. Conversely, explicit <tag>s in r's data can be silently dropped,
// with no corresponding node in the resulting tree.
//
// The input is assumed to be UTF-8 encoded.
func Parse(r io.Reader) (*Node, error) ***REMOVED***
	p := &parser***REMOVED***
		tokenizer: NewTokenizer(r),
		doc: &Node***REMOVED***
			Type: DocumentNode,
		***REMOVED***,
		scripting:  true,
		framesetOK: true,
		im:         initialIM,
	***REMOVED***
	err := p.parse()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return p.doc, nil
***REMOVED***

// ParseFragment parses a fragment of HTML and returns the nodes that were
// found. If the fragment is the InnerHTML for an existing element, pass that
// element in context.
//
// It has the same intricacies as Parse.
func ParseFragment(r io.Reader, context *Node) ([]*Node, error) ***REMOVED***
	contextTag := ""
	if context != nil ***REMOVED***
		if context.Type != ElementNode ***REMOVED***
			return nil, errors.New("html: ParseFragment of non-element Node")
		***REMOVED***
		// The next check isn't just context.DataAtom.String() == context.Data because
		// it is valid to pass an element whose tag isn't a known atom. For example,
		// DataAtom == 0 and Data = "tagfromthefuture" is perfectly consistent.
		if context.DataAtom != a.Lookup([]byte(context.Data)) ***REMOVED***
			return nil, fmt.Errorf("html: inconsistent Node: DataAtom=%q, Data=%q", context.DataAtom, context.Data)
		***REMOVED***
		contextTag = context.DataAtom.String()
	***REMOVED***
	p := &parser***REMOVED***
		tokenizer: NewTokenizerFragment(r, contextTag),
		doc: &Node***REMOVED***
			Type: DocumentNode,
		***REMOVED***,
		scripting: true,
		fragment:  true,
		context:   context,
	***REMOVED***

	root := &Node***REMOVED***
		Type:     ElementNode,
		DataAtom: a.Html,
		Data:     a.Html.String(),
	***REMOVED***
	p.doc.AppendChild(root)
	p.oe = nodeStack***REMOVED***root***REMOVED***
	if context != nil && context.DataAtom == a.Template ***REMOVED***
		p.templateStack = append(p.templateStack, inTemplateIM)
	***REMOVED***
	p.resetInsertionMode()

	for n := context; n != nil; n = n.Parent ***REMOVED***
		if n.Type == ElementNode && n.DataAtom == a.Form ***REMOVED***
			p.form = n
			break
		***REMOVED***
	***REMOVED***

	err := p.parse()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	parent := p.doc
	if context != nil ***REMOVED***
		parent = root
	***REMOVED***

	var result []*Node
	for c := parent.FirstChild; c != nil; ***REMOVED***
		next := c.NextSibling
		parent.RemoveChild(c)
		result = append(result, c)
		c = next
	***REMOVED***
	return result, nil
***REMOVED***
