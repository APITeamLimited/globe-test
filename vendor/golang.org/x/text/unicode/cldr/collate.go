// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cldr

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// RuleProcessor can be passed to Collator's Process method, which
// parses the rules and calls the respective method for each rule found.
type RuleProcessor interface ***REMOVED***
	Reset(anchor string, before int) error
	Insert(level int, str, context, extend string) error
	Index(id string)
***REMOVED***

const (
	// cldrIndex is a Unicode-reserved sentinel value used to mark the start
	// of a grouping within an index.
	// We ignore any rule that starts with this rune.
	// See http://unicode.org/reports/tr35/#Collation_Elements for details.
	cldrIndex = "\uFDD0"

	// specialAnchor is the format in which to represent logical reset positions,
	// such as "first tertiary ignorable".
	specialAnchor = "<%s/>"
)

// Process parses the rules for the tailorings of this collation
// and calls the respective methods of p for each rule found.
func (c Collation) Process(p RuleProcessor) (err error) ***REMOVED***
	if len(c.Cr) > 0 ***REMOVED***
		if len(c.Cr) > 1 ***REMOVED***
			return fmt.Errorf("multiple cr elements, want 0 or 1")
		***REMOVED***
		return processRules(p, c.Cr[0].Data())
	***REMOVED***
	if c.Rules.Any != nil ***REMOVED***
		return c.processXML(p)
	***REMOVED***
	return errors.New("no tailoring data")
***REMOVED***

// processRules parses rules in the Collation Rule Syntax defined in
// http://www.unicode.org/reports/tr35/tr35-collation.html#Collation_Tailorings.
func processRules(p RuleProcessor, s string) (err error) ***REMOVED***
	chk := func(s string, e error) string ***REMOVED***
		if err == nil ***REMOVED***
			err = e
		***REMOVED***
		return s
	***REMOVED***
	i := 0 // Save the line number for use after the loop.
	scanner := bufio.NewScanner(strings.NewReader(s))
	for ; scanner.Scan() && err == nil; i++ ***REMOVED***
		for s := skipSpace(scanner.Text()); s != "" && s[0] != '#'; s = skipSpace(s) ***REMOVED***
			level := 5
			var ch byte
			switch ch, s = s[0], s[1:]; ch ***REMOVED***
			case '&': // followed by <anchor> or '[' <key> ']'
				if s = skipSpace(s); consume(&s, '[') ***REMOVED***
					s = chk(parseSpecialAnchor(p, s))
				***REMOVED*** else ***REMOVED***
					s = chk(parseAnchor(p, 0, s))
				***REMOVED***
			case '<': // sort relation '<'***REMOVED***1,4***REMOVED***, optionally followed by '*'.
				for level = 1; consume(&s, '<'); level++ ***REMOVED***
				***REMOVED***
				if level > 4 ***REMOVED***
					err = fmt.Errorf("level %d > 4", level)
				***REMOVED***
				fallthrough
			case '=': // identity relation, optionally followed by *.
				if consume(&s, '*') ***REMOVED***
					s = chk(parseSequence(p, level, s))
				***REMOVED*** else ***REMOVED***
					s = chk(parseOrder(p, level, s))
				***REMOVED***
			default:
				chk("", fmt.Errorf("illegal operator %q", ch))
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if chk("", scanner.Err()); err != nil ***REMOVED***
		return fmt.Errorf("%d: %v", i, err)
	***REMOVED***
	return nil
***REMOVED***

// parseSpecialAnchor parses the anchor syntax which is either of the form
//    ['before' <level>] <anchor>
// or
//    [<label>]
// The starting should already be consumed.
func parseSpecialAnchor(p RuleProcessor, s string) (tail string, err error) ***REMOVED***
	i := strings.IndexByte(s, ']')
	if i == -1 ***REMOVED***
		return "", errors.New("unmatched bracket")
	***REMOVED***
	a := strings.TrimSpace(s[:i])
	s = s[i+1:]
	if strings.HasPrefix(a, "before ") ***REMOVED***
		l, err := strconv.ParseUint(skipSpace(a[len("before "):]), 10, 3)
		if err != nil ***REMOVED***
			return s, err
		***REMOVED***
		return parseAnchor(p, int(l), s)
	***REMOVED***
	return s, p.Reset(fmt.Sprintf(specialAnchor, a), 0)
***REMOVED***

func parseAnchor(p RuleProcessor, level int, s string) (tail string, err error) ***REMOVED***
	anchor, s, err := scanString(s)
	if err != nil ***REMOVED***
		return s, err
	***REMOVED***
	return s, p.Reset(anchor, level)
***REMOVED***

func parseOrder(p RuleProcessor, level int, s string) (tail string, err error) ***REMOVED***
	var value, context, extend string
	if value, s, err = scanString(s); err != nil ***REMOVED***
		return s, err
	***REMOVED***
	if strings.HasPrefix(value, cldrIndex) ***REMOVED***
		p.Index(value[len(cldrIndex):])
		return
	***REMOVED***
	if consume(&s, '|') ***REMOVED***
		if context, s, err = scanString(s); err != nil ***REMOVED***
			return s, errors.New("missing string after context")
		***REMOVED***
	***REMOVED***
	if consume(&s, '/') ***REMOVED***
		if extend, s, err = scanString(s); err != nil ***REMOVED***
			return s, errors.New("missing string after extension")
		***REMOVED***
	***REMOVED***
	return s, p.Insert(level, value, context, extend)
***REMOVED***

// scanString scans a single input string.
func scanString(s string) (str, tail string, err error) ***REMOVED***
	if s = skipSpace(s); s == "" ***REMOVED***
		return s, s, errors.New("missing string")
	***REMOVED***
	buf := [16]byte***REMOVED******REMOVED*** // small but enough to hold most cases.
	value := buf[:0]
	for s != "" ***REMOVED***
		if consume(&s, '\'') ***REMOVED***
			i := strings.IndexByte(s, '\'')
			if i == -1 ***REMOVED***
				return "", "", errors.New(`unmatched single quote`)
			***REMOVED***
			if i == 0 ***REMOVED***
				value = append(value, '\'')
			***REMOVED*** else ***REMOVED***
				value = append(value, s[:i]...)
			***REMOVED***
			s = s[i+1:]
			continue
		***REMOVED***
		r, sz := utf8.DecodeRuneInString(s)
		if unicode.IsSpace(r) || strings.ContainsRune("&<=#", r) ***REMOVED***
			break
		***REMOVED***
		value = append(value, s[:sz]...)
		s = s[sz:]
	***REMOVED***
	return string(value), skipSpace(s), nil
***REMOVED***

func parseSequence(p RuleProcessor, level int, s string) (tail string, err error) ***REMOVED***
	if s = skipSpace(s); s == "" ***REMOVED***
		return s, errors.New("empty sequence")
	***REMOVED***
	last := rune(0)
	for s != "" ***REMOVED***
		r, sz := utf8.DecodeRuneInString(s)
		s = s[sz:]

		if r == '-' ***REMOVED***
			// We have a range. The first element was already written.
			if last == 0 ***REMOVED***
				return s, errors.New("range without starter value")
			***REMOVED***
			r, sz = utf8.DecodeRuneInString(s)
			s = s[sz:]
			if r == utf8.RuneError || r < last ***REMOVED***
				return s, fmt.Errorf("invalid range %q-%q", last, r)
			***REMOVED***
			for i := last + 1; i <= r; i++ ***REMOVED***
				if err := p.Insert(level, string(i), "", ""); err != nil ***REMOVED***
					return s, err
				***REMOVED***
			***REMOVED***
			last = 0
			continue
		***REMOVED***

		if unicode.IsSpace(r) || unicode.IsPunct(r) ***REMOVED***
			break
		***REMOVED***

		// normal case
		if err := p.Insert(level, string(r), "", ""); err != nil ***REMOVED***
			return s, err
		***REMOVED***
		last = r
	***REMOVED***
	return s, nil
***REMOVED***

func skipSpace(s string) string ***REMOVED***
	return strings.TrimLeftFunc(s, unicode.IsSpace)
***REMOVED***

// consumes returns whether the next byte is ch. If so, it gobbles it by
// updating s.
func consume(s *string, ch byte) (ok bool) ***REMOVED***
	if *s == "" || (*s)[0] != ch ***REMOVED***
		return false
	***REMOVED***
	*s = (*s)[1:]
	return true
***REMOVED***

// The following code parses Collation rules of CLDR version 24 and before.

var lmap = map[byte]int***REMOVED***
	'p': 1,
	's': 2,
	't': 3,
	'i': 5,
***REMOVED***

type rulesElem struct ***REMOVED***
	Rules struct ***REMOVED***
		Common
		Any []*struct ***REMOVED***
			XMLName xml.Name
			rule
		***REMOVED*** `xml:",any"`
	***REMOVED*** `xml:"rules"`
***REMOVED***

type rule struct ***REMOVED***
	Value  string `xml:",chardata"`
	Before string `xml:"before,attr"`
	Any    []*struct ***REMOVED***
		XMLName xml.Name
		rule
	***REMOVED*** `xml:",any"`
***REMOVED***

var emptyValueError = errors.New("cldr: empty rule value")

func (r *rule) value() (string, error) ***REMOVED***
	// Convert hexadecimal Unicode codepoint notation to a string.
	s := charRe.ReplaceAllStringFunc(r.Value, replaceUnicode)
	r.Value = s
	if s == "" ***REMOVED***
		if len(r.Any) != 1 ***REMOVED***
			return "", emptyValueError
		***REMOVED***
		r.Value = fmt.Sprintf(specialAnchor, r.Any[0].XMLName.Local)
		r.Any = nil
	***REMOVED*** else if len(r.Any) != 0 ***REMOVED***
		return "", fmt.Errorf("cldr: XML elements found in collation rule: %v", r.Any)
	***REMOVED***
	return r.Value, nil
***REMOVED***

func (r rule) process(p RuleProcessor, name, context, extend string) error ***REMOVED***
	v, err := r.value()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	switch name ***REMOVED***
	case "p", "s", "t", "i":
		if strings.HasPrefix(v, cldrIndex) ***REMOVED***
			p.Index(v[len(cldrIndex):])
			return nil
		***REMOVED***
		if err := p.Insert(lmap[name[0]], v, context, extend); err != nil ***REMOVED***
			return err
		***REMOVED***
	case "pc", "sc", "tc", "ic":
		level := lmap[name[0]]
		for _, s := range v ***REMOVED***
			if err := p.Insert(level, string(s), context, extend); err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***
	default:
		return fmt.Errorf("cldr: unsupported tag: %q", name)
	***REMOVED***
	return nil
***REMOVED***

// processXML parses the format of CLDR versions 24 and older.
func (c Collation) processXML(p RuleProcessor) (err error) ***REMOVED***
	// Collation is generated and defined in xml.go.
	var v string
	for _, r := range c.Rules.Any ***REMOVED***
		switch r.XMLName.Local ***REMOVED***
		case "reset":
			level := 0
			switch r.Before ***REMOVED***
			case "primary", "1":
				level = 1
			case "secondary", "2":
				level = 2
			case "tertiary", "3":
				level = 3
			case "":
			default:
				return fmt.Errorf("cldr: unknown level %q", r.Before)
			***REMOVED***
			v, err = r.value()
			if err == nil ***REMOVED***
				err = p.Reset(v, level)
			***REMOVED***
		case "x":
			var context, extend string
			for _, r1 := range r.Any ***REMOVED***
				v, err = r1.value()
				switch r1.XMLName.Local ***REMOVED***
				case "context":
					context = v
				case "extend":
					extend = v
				***REMOVED***
			***REMOVED***
			for _, r1 := range r.Any ***REMOVED***
				if t := r1.XMLName.Local; t == "context" || t == "extend" ***REMOVED***
					continue
				***REMOVED***
				r1.rule.process(p, r1.XMLName.Local, context, extend)
			***REMOVED***
		default:
			err = r.rule.process(p, r.XMLName.Local, "", "")
		***REMOVED***
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
