// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ucd provides a parser for Unicode Character Database files, the
// format of which is defined in http://www.unicode.org/reports/tr44/. See
// http://www.unicode.org/Public/UCD/latest/ucd/ for example files.
//
// It currently does not support substitutions of missing fields.
package ucd // import "golang.org/x/text/internal/ucd"

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// UnicodeData.txt fields.
const (
	CodePoint = iota
	Name
	GeneralCategory
	CanonicalCombiningClass
	BidiClass
	DecompMapping
	DecimalValue
	DigitValue
	NumericValue
	BidiMirrored
	Unicode1Name
	ISOComment
	SimpleUppercaseMapping
	SimpleLowercaseMapping
	SimpleTitlecaseMapping
)

// Parse calls f for each entry in the given reader of a UCD file. It will close
// the reader upon return. It will call log.Fatal if any error occurred.
//
// This implements the most common usage pattern of using Parser.
func Parse(r io.ReadCloser, f func(p *Parser)) ***REMOVED***
	defer r.Close()

	p := New(r)
	for p.Next() ***REMOVED***
		f(p)
	***REMOVED***
	if err := p.Err(); err != nil ***REMOVED***
		r.Close() // os.Exit will cause defers not to be called.
		log.Fatal(err)
	***REMOVED***
***REMOVED***

// An Option is used to configure a Parser.
type Option func(p *Parser)

func keepRanges(p *Parser) ***REMOVED***
	p.keepRanges = true
***REMOVED***

var (
	// KeepRanges prevents the expansion of ranges. The raw ranges can be
	// obtained by calling Range(0) on the parser.
	KeepRanges Option = keepRanges
)

// The Part option register a handler for lines starting with a '@'. The text
// after a '@' is available as the first field. Comments are handled as usual.
func Part(f func(p *Parser)) Option ***REMOVED***
	return func(p *Parser) ***REMOVED***
		p.partHandler = f
	***REMOVED***
***REMOVED***

// The CommentHandler option passes comments that are on a line by itself to
// a given handler.
func CommentHandler(f func(s string)) Option ***REMOVED***
	return func(p *Parser) ***REMOVED***
		p.commentHandler = f
	***REMOVED***
***REMOVED***

// A Parser parses Unicode Character Database (UCD) files.
type Parser struct ***REMOVED***
	scanner *bufio.Scanner

	keepRanges bool // Don't expand rune ranges in field 0.

	err     error
	comment string
	field   []string
	// parsedRange is needed in case Range(0) is called more than once for one
	// field. In some cases this requires scanning ahead.
	line                 int
	parsedRange          bool
	rangeStart, rangeEnd rune

	partHandler    func(p *Parser)
	commentHandler func(s string)
***REMOVED***

func (p *Parser) setError(err error, msg string) ***REMOVED***
	if p.err == nil && err != nil ***REMOVED***
		if msg == "" ***REMOVED***
			p.err = fmt.Errorf("ucd:line:%d: %v", p.line, err)
		***REMOVED*** else ***REMOVED***
			p.err = fmt.Errorf("ucd:line:%d:%s: %v", p.line, msg, err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (p *Parser) getField(i int) string ***REMOVED***
	if i >= len(p.field) ***REMOVED***
		return ""
	***REMOVED***
	return p.field[i]
***REMOVED***

// Err returns a non-nil error if any error occurred during parsing.
func (p *Parser) Err() error ***REMOVED***
	return p.err
***REMOVED***

// New returns a Parser for the given Reader.
func New(r io.Reader, o ...Option) *Parser ***REMOVED***
	p := &Parser***REMOVED***
		scanner: bufio.NewScanner(r),
	***REMOVED***
	for _, f := range o ***REMOVED***
		f(p)
	***REMOVED***
	return p
***REMOVED***

// Next parses the next line in the file. It returns true if a line was parsed
// and false if it reached the end of the file.
func (p *Parser) Next() bool ***REMOVED***
	if !p.keepRanges && p.rangeStart < p.rangeEnd ***REMOVED***
		p.rangeStart++
		return true
	***REMOVED***
	p.comment = ""
	p.field = p.field[:0]
	p.parsedRange = false

	for p.scanner.Scan() && p.err == nil ***REMOVED***
		p.line++
		s := p.scanner.Text()
		if s == "" ***REMOVED***
			continue
		***REMOVED***
		if s[0] == '#' ***REMOVED***
			if p.commentHandler != nil ***REMOVED***
				p.commentHandler(strings.TrimSpace(s[1:]))
			***REMOVED***
			continue
		***REMOVED***

		// Parse line
		if i := strings.IndexByte(s, '#'); i != -1 ***REMOVED***
			p.comment = strings.TrimSpace(s[i+1:])
			s = s[:i]
		***REMOVED***
		if s[0] == '@' ***REMOVED***
			if p.partHandler != nil ***REMOVED***
				p.field = append(p.field, strings.TrimSpace(s[1:]))
				p.partHandler(p)
				p.field = p.field[:0]
			***REMOVED***
			p.comment = ""
			continue
		***REMOVED***
		for ***REMOVED***
			i := strings.IndexByte(s, ';')
			if i == -1 ***REMOVED***
				p.field = append(p.field, strings.TrimSpace(s))
				break
			***REMOVED***
			p.field = append(p.field, strings.TrimSpace(s[:i]))
			s = s[i+1:]
		***REMOVED***
		if !p.keepRanges ***REMOVED***
			p.rangeStart, p.rangeEnd = p.getRange(0)
		***REMOVED***
		return true
	***REMOVED***
	p.setError(p.scanner.Err(), "scanner failed")
	return false
***REMOVED***

func parseRune(b string) (rune, error) ***REMOVED***
	if len(b) > 2 && b[0] == 'U' && b[1] == '+' ***REMOVED***
		b = b[2:]
	***REMOVED***
	x, err := strconv.ParseUint(b, 16, 32)
	return rune(x), err
***REMOVED***

func (p *Parser) parseRune(s string) rune ***REMOVED***
	x, err := parseRune(s)
	p.setError(err, "failed to parse rune")
	return x
***REMOVED***

// Rune parses and returns field i as a rune.
func (p *Parser) Rune(i int) rune ***REMOVED***
	if i > 0 || p.keepRanges ***REMOVED***
		return p.parseRune(p.getField(i))
	***REMOVED***
	return p.rangeStart
***REMOVED***

// Runes interprets and returns field i as a sequence of runes.
func (p *Parser) Runes(i int) (runes []rune) ***REMOVED***
	add := func(s string) ***REMOVED***
		if s = strings.TrimSpace(s); len(s) > 0 ***REMOVED***
			runes = append(runes, p.parseRune(s))
		***REMOVED***
	***REMOVED***
	for b := p.getField(i); ; ***REMOVED***
		i := strings.IndexByte(b, ' ')
		if i == -1 ***REMOVED***
			add(b)
			break
		***REMOVED***
		add(b[:i])
		b = b[i+1:]
	***REMOVED***
	return
***REMOVED***

var (
	errIncorrectLegacyRange = errors.New("ucd: unmatched <* First>")

	// reRange matches one line of a legacy rune range.
	reRange = regexp.MustCompile("^([0-9A-F]*);<([^,]*), ([^>]*)>(.*)$")
)

// Range parses and returns field i as a rune range. A range is inclusive at
// both ends. If the field only has one rune, first and last will be identical.
// It supports the legacy format for ranges used in UnicodeData.txt.
func (p *Parser) Range(i int) (first, last rune) ***REMOVED***
	if !p.keepRanges ***REMOVED***
		return p.rangeStart, p.rangeStart
	***REMOVED***
	return p.getRange(i)
***REMOVED***

func (p *Parser) getRange(i int) (first, last rune) ***REMOVED***
	b := p.getField(i)
	if k := strings.Index(b, ".."); k != -1 ***REMOVED***
		return p.parseRune(b[:k]), p.parseRune(b[k+2:])
	***REMOVED***
	// The first field may not be a rune, in which case we may ignore any error
	// and set the range as 0..0.
	x, err := parseRune(b)
	if err != nil ***REMOVED***
		// Disable range parsing henceforth. This ensures that an error will be
		// returned if the user subsequently will try to parse this field as
		// a Rune.
		p.keepRanges = true
	***REMOVED***
	// Special case for UnicodeData that was retained for backwards compatibility.
	if i == 0 && len(p.field) > 1 && strings.HasSuffix(p.field[1], "First>") ***REMOVED***
		if p.parsedRange ***REMOVED***
			return p.rangeStart, p.rangeEnd
		***REMOVED***
		mf := reRange.FindStringSubmatch(p.scanner.Text())
		p.line++
		if mf == nil || !p.scanner.Scan() ***REMOVED***
			p.setError(errIncorrectLegacyRange, "")
			return x, x
		***REMOVED***
		// Using Bytes would be more efficient here, but Text is a lot easier
		// and this is not a frequent case.
		ml := reRange.FindStringSubmatch(p.scanner.Text())
		if ml == nil || mf[2] != ml[2] || ml[3] != "Last" || mf[4] != ml[4] ***REMOVED***
			p.setError(errIncorrectLegacyRange, "")
			return x, x
		***REMOVED***
		p.rangeStart, p.rangeEnd = x, p.parseRune(p.scanner.Text()[:len(ml[1])])
		p.parsedRange = true
		return p.rangeStart, p.rangeEnd
	***REMOVED***
	return x, x
***REMOVED***

// bools recognizes all valid UCD boolean values.
var bools = map[string]bool***REMOVED***
	"":      false,
	"N":     false,
	"No":    false,
	"F":     false,
	"False": false,
	"Y":     true,
	"Yes":   true,
	"T":     true,
	"True":  true,
***REMOVED***

// Bool parses and returns field i as a boolean value.
func (p *Parser) Bool(i int) bool ***REMOVED***
	f := p.getField(i)
	for s, v := range bools ***REMOVED***
		if f == s ***REMOVED***
			return v
		***REMOVED***
	***REMOVED***
	p.setError(strconv.ErrSyntax, "error parsing bool")
	return false
***REMOVED***

// Int parses and returns field i as an integer value.
func (p *Parser) Int(i int) int ***REMOVED***
	x, err := strconv.ParseInt(string(p.getField(i)), 10, 64)
	p.setError(err, "error parsing int")
	return int(x)
***REMOVED***

// Uint parses and returns field i as an unsigned integer value.
func (p *Parser) Uint(i int) uint ***REMOVED***
	x, err := strconv.ParseUint(string(p.getField(i)), 10, 64)
	p.setError(err, "error parsing uint")
	return uint(x)
***REMOVED***

// Float parses and returns field i as a decimal value.
func (p *Parser) Float(i int) float64 ***REMOVED***
	x, err := strconv.ParseFloat(string(p.getField(i)), 64)
	p.setError(err, "error parsing float")
	return x
***REMOVED***

// String parses and returns field i as a string value.
func (p *Parser) String(i int) string ***REMOVED***
	return string(p.getField(i))
***REMOVED***

// Strings parses and returns field i as a space-separated list of strings.
func (p *Parser) Strings(i int) []string ***REMOVED***
	ss := strings.Split(string(p.getField(i)), " ")
	for i, s := range ss ***REMOVED***
		ss[i] = strings.TrimSpace(s)
	***REMOVED***
	return ss
***REMOVED***

// Comment returns the comments for the current line.
func (p *Parser) Comment() string ***REMOVED***
	return string(p.comment)
***REMOVED***

var errUndefinedEnum = errors.New("ucd: undefined enum value")

// Enum interprets and returns field i as a value that must be one of the values
// in enum.
func (p *Parser) Enum(i int, enum ...string) string ***REMOVED***
	f := p.getField(i)
	for _, s := range enum ***REMOVED***
		if f == s ***REMOVED***
			return s
		***REMOVED***
	***REMOVED***
	p.setError(errUndefinedEnum, "error parsing enum")
	return ""
***REMOVED***
