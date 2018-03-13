// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import (
	"encoding/xml"
	"regexp"
	"strconv"
)

// Elem is implemented by every XML element.
type Elem interface ***REMOVED***
	setEnclosing(Elem)
	setName(string)
	enclosing() Elem

	GetCommon() *Common
***REMOVED***

type hidden struct ***REMOVED***
	CharData string `xml:",chardata"`
	Alias    *struct ***REMOVED***
		Common
		Source string `xml:"source,attr"`
		Path   string `xml:"path,attr"`
	***REMOVED*** `xml:"alias"`
	Def *struct ***REMOVED***
		Common
		Choice string `xml:"choice,attr,omitempty"`
		Type   string `xml:"type,attr,omitempty"`
	***REMOVED*** `xml:"default"`
***REMOVED***

// Common holds several of the most common attributes and sub elements
// of an XML element.
type Common struct ***REMOVED***
	XMLName         xml.Name
	name            string
	enclElem        Elem
	Type            string `xml:"type,attr,omitempty"`
	Reference       string `xml:"reference,attr,omitempty"`
	Alt             string `xml:"alt,attr,omitempty"`
	ValidSubLocales string `xml:"validSubLocales,attr,omitempty"`
	Draft           string `xml:"draft,attr,omitempty"`
	hidden
***REMOVED***

// Default returns the default type to select from the enclosed list
// or "" if no default value is specified.
func (e *Common) Default() string ***REMOVED***
	if e.Def == nil ***REMOVED***
		return ""
	***REMOVED***
	if e.Def.Choice != "" ***REMOVED***
		return e.Def.Choice
	***REMOVED*** else if e.Def.Type != "" ***REMOVED***
		// Type is still used by the default element in collation.
		return e.Def.Type
	***REMOVED***
	return ""
***REMOVED***

// Element returns the XML element name.
func (e *Common) Element() string ***REMOVED***
	return e.name
***REMOVED***

// GetCommon returns e. It is provided such that Common implements Elem.
func (e *Common) GetCommon() *Common ***REMOVED***
	return e
***REMOVED***

// Data returns the character data accumulated for this element.
func (e *Common) Data() string ***REMOVED***
	e.CharData = charRe.ReplaceAllStringFunc(e.CharData, replaceUnicode)
	return e.CharData
***REMOVED***

func (e *Common) setName(s string) ***REMOVED***
	e.name = s
***REMOVED***

func (e *Common) enclosing() Elem ***REMOVED***
	return e.enclElem
***REMOVED***

func (e *Common) setEnclosing(en Elem) ***REMOVED***
	e.enclElem = en
***REMOVED***

// Escape characters that can be escaped without further escaping the string.
var charRe = regexp.MustCompile(`&#x[0-9a-fA-F]*;|\\u[0-9a-fA-F]***REMOVED***4***REMOVED***|\\U[0-9a-fA-F]***REMOVED***8***REMOVED***|\\x[0-9a-fA-F]***REMOVED***2***REMOVED***|\\[0-7]***REMOVED***3***REMOVED***|\\[abtnvfr]`)

// replaceUnicode converts hexadecimal Unicode codepoint notations to a one-rune string.
// It assumes the input string is correctly formatted.
func replaceUnicode(s string) string ***REMOVED***
	if s[1] == '#' ***REMOVED***
		r, _ := strconv.ParseInt(s[3:len(s)-1], 16, 32)
		return string(r)
	***REMOVED***
	r, _, _, _ := strconv.UnquoteChar(s, 0)
	return string(r)
***REMOVED***
