// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:generate go run makexml.go -output xml.go

// Package cldr provides a parser for LDML and related XML formats.
// This package is intended to be used by the table generation tools
// for the various internationalization-related packages.
// As the XML types are generated from the CLDR DTD, and as the CLDR standard
// is periodically amended, this package may change considerably over time.
// This mostly means that data may appear and disappear between versions.
// That is, old code should keep compiling for newer versions, but data
// may have moved or changed.
// CLDR version 22 is the first version supported by this package.
// Older versions may not work.
package cldr // import "golang.org/x/text/unicode/cldr"

import (
	"fmt"
	"sort"
)

// CLDR provides access to parsed data of the Unicode Common Locale Data Repository.
type CLDR struct ***REMOVED***
	parent   map[string][]string
	locale   map[string]*LDML
	resolved map[string]*LDML
	bcp47    *LDMLBCP47
	supp     *SupplementalData
***REMOVED***

func makeCLDR() *CLDR ***REMOVED***
	return &CLDR***REMOVED***
		parent:   make(map[string][]string),
		locale:   make(map[string]*LDML),
		resolved: make(map[string]*LDML),
		bcp47:    &LDMLBCP47***REMOVED******REMOVED***,
		supp:     &SupplementalData***REMOVED******REMOVED***,
	***REMOVED***
***REMOVED***

// BCP47 returns the parsed BCP47 LDML data. If no such data was parsed, nil is returned.
func (cldr *CLDR) BCP47() *LDMLBCP47 ***REMOVED***
	return nil
***REMOVED***

// Draft indicates the draft level of an element.
type Draft int

const (
	Approved Draft = iota
	Contributed
	Provisional
	Unconfirmed
)

var drafts = []string***REMOVED***"unconfirmed", "provisional", "contributed", "approved", ""***REMOVED***

// ParseDraft returns the Draft value corresponding to the given string. The
// empty string corresponds to Approved.
func ParseDraft(level string) (Draft, error) ***REMOVED***
	if level == "" ***REMOVED***
		return Approved, nil
	***REMOVED***
	for i, s := range drafts ***REMOVED***
		if level == s ***REMOVED***
			return Unconfirmed - Draft(i), nil
		***REMOVED***
	***REMOVED***
	return Approved, fmt.Errorf("cldr: unknown draft level %q", level)
***REMOVED***

func (d Draft) String() string ***REMOVED***
	return drafts[len(drafts)-1-int(d)]
***REMOVED***

// SetDraftLevel sets which draft levels to include in the evaluated LDML.
// Any draft element for which the draft level is higher than lev will be excluded.
// If multiple draft levels are available for a single element, the one with the
// lowest draft level will be selected, unless preferDraft is true, in which case
// the highest draft will be chosen.
// It is assumed that the underlying LDML is canonicalized.
func (cldr *CLDR) SetDraftLevel(lev Draft, preferDraft bool) ***REMOVED***
	// TODO: implement
	cldr.resolved = make(map[string]*LDML)
***REMOVED***

// RawLDML returns the LDML XML for id in unresolved form.
// id must be one of the strings returned by Locales.
func (cldr *CLDR) RawLDML(loc string) *LDML ***REMOVED***
	return cldr.locale[loc]
***REMOVED***

// LDML returns the fully resolved LDML XML for loc, which must be one of
// the strings returned by Locales.
func (cldr *CLDR) LDML(loc string) (*LDML, error) ***REMOVED***
	return cldr.resolve(loc)
***REMOVED***

// Supplemental returns the parsed supplemental data. If no such data was parsed,
// nil is returned.
func (cldr *CLDR) Supplemental() *SupplementalData ***REMOVED***
	return cldr.supp
***REMOVED***

// Locales returns the locales for which there exist files.
// Valid sublocales for which there is no file are not included.
// The root locale is always sorted first.
func (cldr *CLDR) Locales() []string ***REMOVED***
	loc := []string***REMOVED***"root"***REMOVED***
	hasRoot := false
	for l, _ := range cldr.locale ***REMOVED***
		if l == "root" ***REMOVED***
			hasRoot = true
			continue
		***REMOVED***
		loc = append(loc, l)
	***REMOVED***
	sort.Strings(loc[1:])
	if !hasRoot ***REMOVED***
		return loc[1:]
	***REMOVED***
	return loc
***REMOVED***

// Get fills in the fields of x based on the XPath path.
func Get(e Elem, path string) (res Elem, err error) ***REMOVED***
	return walkXPath(e, path)
***REMOVED***
