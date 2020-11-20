// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type writer interface ***REMOVED***
	io.Writer
	io.ByteWriter
	WriteString(string) (int, error)
***REMOVED***

// Render renders the parse tree n to the given writer.
//
// Rendering is done on a 'best effort' basis: calling Parse on the output of
// Render will always result in something similar to the original tree, but it
// is not necessarily an exact clone unless the original tree was 'well-formed'.
// 'Well-formed' is not easily specified; the HTML5 specification is
// complicated.
//
// Calling Parse on arbitrary input typically results in a 'well-formed' parse
// tree. However, it is possible for Parse to yield a 'badly-formed' parse tree.
// For example, in a 'well-formed' parse tree, no <a> element is a child of
// another <a> element: parsing "<a><a>" results in two sibling elements.
// Similarly, in a 'well-formed' parse tree, no <a> element is a child of a
// <table> element: parsing "<p><table><a>" results in a <p> with two sibling
// children; the <a> is reparented to the <table>'s parent. However, calling
// Parse on "<a><table><a>" does not return an error, but the result has an <a>
// element with an <a> child, and is therefore not 'well-formed'.
//
// Programmatically constructed trees are typically also 'well-formed', but it
// is possible to construct a tree that looks innocuous but, when rendered and
// re-parsed, results in a different tree. A simple example is that a solitary
// text node would become a tree containing <html>, <head> and <body> elements.
// Another example is that the programmatic equivalent of "a<head>b</head>c"
// becomes "<html><head><head/><body>abc</body></html>".
func Render(w io.Writer, n *Node) error ***REMOVED***
	if x, ok := w.(writer); ok ***REMOVED***
		return render(x, n)
	***REMOVED***
	buf := bufio.NewWriter(w)
	if err := render(buf, n); err != nil ***REMOVED***
		return err
	***REMOVED***
	return buf.Flush()
***REMOVED***

// plaintextAbort is returned from render1 when a <plaintext> element
// has been rendered. No more end tags should be rendered after that.
var plaintextAbort = errors.New("html: internal error (plaintext abort)")

func render(w writer, n *Node) error ***REMOVED***
	err := render1(w, n)
	if err == plaintextAbort ***REMOVED***
		err = nil
	***REMOVED***
	return err
***REMOVED***

func render1(w writer, n *Node) error ***REMOVED***
	// Render non-element nodes; these are the easy cases.
	switch n.Type ***REMOVED***
	case ErrorNode:
		return errors.New("html: cannot render an ErrorNode node")
	case TextNode:
		return escape(w, n.Data)
	case DocumentNode:
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			if err := render1(w, c); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		return nil
	case ElementNode:
		// No-op.
	case CommentNode:
		if _, err := w.WriteString("<!--"); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := w.WriteString(n.Data); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := w.WriteString("-->"); err != nil ***REMOVED***
			return err
		***REMOVED***
		return nil
	case DoctypeNode:
		if _, err := w.WriteString("<!DOCTYPE "); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := w.WriteString(n.Data); err != nil ***REMOVED***
			return err
		***REMOVED***
		if n.Attr != nil ***REMOVED***
			var p, s string
			for _, a := range n.Attr ***REMOVED***
				switch a.Key ***REMOVED***
				case "public":
					p = a.Val
				case "system":
					s = a.Val
				***REMOVED***
			***REMOVED***
			if p != "" ***REMOVED***
				if _, err := w.WriteString(" PUBLIC "); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := writeQuoted(w, p); err != nil ***REMOVED***
					return err
				***REMOVED***
				if s != "" ***REMOVED***
					if err := w.WriteByte(' '); err != nil ***REMOVED***
						return err
					***REMOVED***
					if err := writeQuoted(w, s); err != nil ***REMOVED***
						return err
					***REMOVED***
				***REMOVED***
			***REMOVED*** else if s != "" ***REMOVED***
				if _, err := w.WriteString(" SYSTEM "); err != nil ***REMOVED***
					return err
				***REMOVED***
				if err := writeQuoted(w, s); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		return w.WriteByte('>')
	case RawNode:
		_, err := w.WriteString(n.Data)
		return err
	default:
		return errors.New("html: unknown node type")
	***REMOVED***

	// Render the <xxx> opening tag.
	if err := w.WriteByte('<'); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.WriteString(n.Data); err != nil ***REMOVED***
		return err
	***REMOVED***
	for _, a := range n.Attr ***REMOVED***
		if err := w.WriteByte(' '); err != nil ***REMOVED***
			return err
		***REMOVED***
		if a.Namespace != "" ***REMOVED***
			if _, err := w.WriteString(a.Namespace); err != nil ***REMOVED***
				return err
			***REMOVED***
			if err := w.WriteByte(':'); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
		if _, err := w.WriteString(a.Key); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := w.WriteString(`="`); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := escape(w, a.Val); err != nil ***REMOVED***
			return err
		***REMOVED***
		if err := w.WriteByte('"'); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	if voidElements[n.Data] ***REMOVED***
		if n.FirstChild != nil ***REMOVED***
			return fmt.Errorf("html: void element <%s> has child nodes", n.Data)
		***REMOVED***
		_, err := w.WriteString("/>")
		return err
	***REMOVED***
	if err := w.WriteByte('>'); err != nil ***REMOVED***
		return err
	***REMOVED***

	// Add initial newline where there is danger of a newline beging ignored.
	if c := n.FirstChild; c != nil && c.Type == TextNode && strings.HasPrefix(c.Data, "\n") ***REMOVED***
		switch n.Data ***REMOVED***
		case "pre", "listing", "textarea":
			if err := w.WriteByte('\n'); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Render any child nodes.
	switch n.Data ***REMOVED***
	case "iframe", "noembed", "noframes", "noscript", "plaintext", "script", "style", "xmp":
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			if c.Type == TextNode ***REMOVED***
				if _, err := w.WriteString(c.Data); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				if err := render1(w, c); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if n.Data == "plaintext" ***REMOVED***
			// Don't render anything else. <plaintext> must be the
			// last element in the file, with no closing tag.
			return plaintextAbort
		***REMOVED***
	default:
		for c := n.FirstChild; c != nil; c = c.NextSibling ***REMOVED***
			if err := render1(w, c); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// Render the </xxx> closing tag.
	if _, err := w.WriteString("</"); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.WriteString(n.Data); err != nil ***REMOVED***
		return err
	***REMOVED***
	return w.WriteByte('>')
***REMOVED***

// writeQuoted writes s to w surrounded by quotes. Normally it will use double
// quotes, but if s contains a double quote, it will use single quotes.
// It is used for writing the identifiers in a doctype declaration.
// In valid HTML, they can't contain both types of quotes.
func writeQuoted(w writer, s string) error ***REMOVED***
	var q byte = '"'
	if strings.Contains(s, `"`) ***REMOVED***
		q = '\''
	***REMOVED***
	if err := w.WriteByte(q); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err := w.WriteString(s); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err := w.WriteByte(q); err != nil ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

// Section 12.1.2, "Elements", gives this list of void elements. Void elements
// are those that can't have any contents.
var voidElements = map[string]bool***REMOVED***
	"area":   true,
	"base":   true,
	"br":     true,
	"col":    true,
	"embed":  true,
	"hr":     true,
	"img":    true,
	"input":  true,
	"keygen": true,
	"link":   true,
	"meta":   true,
	"param":  true,
	"source": true,
	"track":  true,
	"wbr":    true,
***REMOVED***
