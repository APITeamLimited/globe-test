// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package html implements an HTML5-compliant tokenizer and parser.

Tokenization is done by creating a Tokenizer for an io.Reader r. It is the
caller's responsibility to ensure that r provides UTF-8 encoded HTML.

	z := html.NewTokenizer(r)

Given a Tokenizer z, the HTML is tokenized by repeatedly calling z.Next(),
which parses the next token and returns its type, or an error:

	for ***REMOVED***
		tt := z.Next()
		if tt == html.ErrorToken ***REMOVED***
			// ...
			return ...
		***REMOVED***
		// Process the current token.
	***REMOVED***

There are two APIs for retrieving the current token. The high-level API is to
call Token; the low-level API is to call Text or TagName / TagAttr. Both APIs
allow optionally calling Raw after Next but before Token, Text, TagName, or
TagAttr. In EBNF notation, the valid call sequence per token is:

	Next ***REMOVED***Raw***REMOVED*** [ Token | Text | TagName ***REMOVED***TagAttr***REMOVED*** ]

Token returns an independent data structure that completely describes a token.
Entities (such as "&lt;") are unescaped, tag names and attribute keys are
lower-cased, and attributes are collected into a []Attribute. For example:

	for ***REMOVED***
		if z.Next() == html.ErrorToken ***REMOVED***
			// Returning io.EOF indicates success.
			return z.Err()
		***REMOVED***
		emitToken(z.Token())
	***REMOVED***

The low-level API performs fewer allocations and copies, but the contents of
the []byte values returned by Text, TagName and TagAttr may change on the next
call to Next. For example, to extract an HTML page's anchor text:

	depth := 0
	for ***REMOVED***
		tt := z.Next()
		switch tt ***REMOVED***
		case html.ErrorToken:
			return z.Err()
		case html.TextToken:
			if depth > 0 ***REMOVED***
				// emitBytes should copy the []byte it receives,
				// if it doesn't process it immediately.
				emitBytes(z.Text())
			***REMOVED***
		case html.StartTagToken, html.EndTagToken:
			tn, _ := z.TagName()
			if len(tn) == 1 && tn[0] == 'a' ***REMOVED***
				if tt == html.StartTagToken ***REMOVED***
					depth++
				***REMOVED*** else ***REMOVED***
					depth--
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

Parsing is done by calling Parse with an io.Reader, which returns the root of
the parse tree (the document element) as a *Node. It is the caller's
responsibility to ensure that the Reader provides UTF-8 encoded HTML. For
example, to process each anchor node in depth-first order:

	doc, err := html.Parse(r)
	if err != nil ***REMOVED***
		// ...
	***REMOVED***
	var f func(*html.Node)
	f = func(n *html.Node) ***REMOVED***
		if n.Type == html.ElementNode && n.Data == "a" ***REMOVED***
			// Do something with n...
		***REMOVED***
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			f(c)
		***REMOVED***
	***REMOVED***
	f(doc)

The relevant specifications include:
https://html.spec.whatwg.org/multipage/syntax.html and
https://html.spec.whatwg.org/multipage/syntax.html#tokenization
*/
package html // import "golang.org/x/net/html"

// The tokenization algorithm implemented by this package is not a line-by-line
// transliteration of the relatively verbose state-machine in the WHATWG
// specification. A more direct approach is used instead, where the program
// counter implies the state, such as whether it is tokenizing a tag or a text
// node. Specification compliance is verified by checking expected and actual
// outputs over a test suite rather than aiming for algorithmic fidelity.

// TODO(nigeltao): Does a DOM API belong in this package or a separate one?
// TODO(nigeltao): How does parsing interact with a JavaScript engine?
