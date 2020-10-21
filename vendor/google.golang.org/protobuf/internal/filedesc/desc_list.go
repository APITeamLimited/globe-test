// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/descfmt"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type FileImports []pref.FileImport

func (p *FileImports) Len() int                            ***REMOVED*** return len(*p) ***REMOVED***
func (p *FileImports) Get(i int) pref.FileImport           ***REMOVED*** return (*p)[i] ***REMOVED***
func (p *FileImports) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *FileImports) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***

type Names struct ***REMOVED***
	List []pref.Name
	once sync.Once
	has  map[pref.Name]int // protected by once
***REMOVED***

func (p *Names) Len() int                            ***REMOVED*** return len(p.List) ***REMOVED***
func (p *Names) Get(i int) pref.Name                 ***REMOVED*** return p.List[i] ***REMOVED***
func (p *Names) Has(s pref.Name) bool                ***REMOVED*** return p.lazyInit().has[s] > 0 ***REMOVED***
func (p *Names) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *Names) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***
func (p *Names) lazyInit() *Names ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.List) > 0 ***REMOVED***
			p.has = make(map[pref.Name]int, len(p.List))
			for _, s := range p.List ***REMOVED***
				p.has[s] = p.has[s] + 1
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return p
***REMOVED***

// CheckValid reports any errors with the set of names with an error message
// that completes the sentence: "ranges is invalid because it has ..."
func (p *Names) CheckValid() error ***REMOVED***
	for s, n := range p.lazyInit().has ***REMOVED***
		switch ***REMOVED***
		case n > 1:
			return errors.New("duplicate name: %q", s)
		case false && !s.IsValid():
			// NOTE: The C++ implementation does not validate the identifier.
			// See https://github.com/protocolbuffers/protobuf/issues/6335.
			return errors.New("invalid name: %q", s)
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type EnumRanges struct ***REMOVED***
	List   [][2]pref.EnumNumber // start inclusive; end inclusive
	once   sync.Once
	sorted [][2]pref.EnumNumber // protected by once
***REMOVED***

func (p *EnumRanges) Len() int                     ***REMOVED*** return len(p.List) ***REMOVED***
func (p *EnumRanges) Get(i int) [2]pref.EnumNumber ***REMOVED*** return p.List[i] ***REMOVED***
func (p *EnumRanges) Has(n pref.EnumNumber) bool ***REMOVED***
	for ls := p.lazyInit().sorted; len(ls) > 0; ***REMOVED***
		i := len(ls) / 2
		switch r := enumRange(ls[i]); ***REMOVED***
		case n < r.Start():
			ls = ls[:i] // search lower
		case n > r.End():
			ls = ls[i+1:] // search upper
		default:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
func (p *EnumRanges) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *EnumRanges) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***
func (p *EnumRanges) lazyInit() *EnumRanges ***REMOVED***
	p.once.Do(func() ***REMOVED***
		p.sorted = append(p.sorted, p.List...)
		sort.Slice(p.sorted, func(i, j int) bool ***REMOVED***
			return p.sorted[i][0] < p.sorted[j][0]
		***REMOVED***)
	***REMOVED***)
	return p
***REMOVED***

// CheckValid reports any errors with the set of names with an error message
// that completes the sentence: "ranges is invalid because it has ..."
func (p *EnumRanges) CheckValid() error ***REMOVED***
	var rp enumRange
	for i, r := range p.lazyInit().sorted ***REMOVED***
		r := enumRange(r)
		switch ***REMOVED***
		case !(r.Start() <= r.End()):
			return errors.New("invalid range: %v", r)
		case !(rp.End() < r.Start()) && i > 0:
			return errors.New("overlapping ranges: %v with %v", rp, r)
		***REMOVED***
		rp = r
	***REMOVED***
	return nil
***REMOVED***

type enumRange [2]protoreflect.EnumNumber

func (r enumRange) Start() protoreflect.EnumNumber ***REMOVED*** return r[0] ***REMOVED*** // inclusive
func (r enumRange) End() protoreflect.EnumNumber   ***REMOVED*** return r[1] ***REMOVED*** // inclusive
func (r enumRange) String() string ***REMOVED***
	if r.Start() == r.End() ***REMOVED***
		return fmt.Sprintf("%d", r.Start())
	***REMOVED***
	return fmt.Sprintf("%d to %d", r.Start(), r.End())
***REMOVED***

type FieldRanges struct ***REMOVED***
	List   [][2]pref.FieldNumber // start inclusive; end exclusive
	once   sync.Once
	sorted [][2]pref.FieldNumber // protected by once
***REMOVED***

func (p *FieldRanges) Len() int                      ***REMOVED*** return len(p.List) ***REMOVED***
func (p *FieldRanges) Get(i int) [2]pref.FieldNumber ***REMOVED*** return p.List[i] ***REMOVED***
func (p *FieldRanges) Has(n pref.FieldNumber) bool ***REMOVED***
	for ls := p.lazyInit().sorted; len(ls) > 0; ***REMOVED***
		i := len(ls) / 2
		switch r := fieldRange(ls[i]); ***REMOVED***
		case n < r.Start():
			ls = ls[:i] // search lower
		case n > r.End():
			ls = ls[i+1:] // search upper
		default:
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
func (p *FieldRanges) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *FieldRanges) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***
func (p *FieldRanges) lazyInit() *FieldRanges ***REMOVED***
	p.once.Do(func() ***REMOVED***
		p.sorted = append(p.sorted, p.List...)
		sort.Slice(p.sorted, func(i, j int) bool ***REMOVED***
			return p.sorted[i][0] < p.sorted[j][0]
		***REMOVED***)
	***REMOVED***)
	return p
***REMOVED***

// CheckValid reports any errors with the set of ranges with an error message
// that completes the sentence: "ranges is invalid because it has ..."
func (p *FieldRanges) CheckValid(isMessageSet bool) error ***REMOVED***
	var rp fieldRange
	for i, r := range p.lazyInit().sorted ***REMOVED***
		r := fieldRange(r)
		switch ***REMOVED***
		case !isValidFieldNumber(r.Start(), isMessageSet):
			return errors.New("invalid field number: %d", r.Start())
		case !isValidFieldNumber(r.End(), isMessageSet):
			return errors.New("invalid field number: %d", r.End())
		case !(r.Start() <= r.End()):
			return errors.New("invalid range: %v", r)
		case !(rp.End() < r.Start()) && i > 0:
			return errors.New("overlapping ranges: %v with %v", rp, r)
		***REMOVED***
		rp = r
	***REMOVED***
	return nil
***REMOVED***

// isValidFieldNumber reports whether the field number is valid.
// Unlike the FieldNumber.IsValid method, it allows ranges that cover the
// reserved number range.
func isValidFieldNumber(n protoreflect.FieldNumber, isMessageSet bool) bool ***REMOVED***
	if isMessageSet ***REMOVED***
		return protowire.MinValidNumber <= n && n <= math.MaxInt32
	***REMOVED***
	return protowire.MinValidNumber <= n && n <= protowire.MaxValidNumber
***REMOVED***

// CheckOverlap reports an error if p and q overlap.
func (p *FieldRanges) CheckOverlap(q *FieldRanges) error ***REMOVED***
	rps := p.lazyInit().sorted
	rqs := q.lazyInit().sorted
	for pi, qi := 0, 0; pi < len(rps) && qi < len(rqs); ***REMOVED***
		rp := fieldRange(rps[pi])
		rq := fieldRange(rqs[qi])
		if !(rp.End() < rq.Start() || rq.End() < rp.Start()) ***REMOVED***
			return errors.New("overlapping ranges: %v with %v", rp, rq)
		***REMOVED***
		if rp.Start() < rq.Start() ***REMOVED***
			pi++
		***REMOVED*** else ***REMOVED***
			qi++
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

type fieldRange [2]protoreflect.FieldNumber

func (r fieldRange) Start() protoreflect.FieldNumber ***REMOVED*** return r[0] ***REMOVED***     // inclusive
func (r fieldRange) End() protoreflect.FieldNumber   ***REMOVED*** return r[1] - 1 ***REMOVED*** // inclusive
func (r fieldRange) String() string ***REMOVED***
	if r.Start() == r.End() ***REMOVED***
		return fmt.Sprintf("%d", r.Start())
	***REMOVED***
	return fmt.Sprintf("%d to %d", r.Start(), r.End())
***REMOVED***

type FieldNumbers struct ***REMOVED***
	List []pref.FieldNumber
	once sync.Once
	has  map[pref.FieldNumber]struct***REMOVED******REMOVED*** // protected by once
***REMOVED***

func (p *FieldNumbers) Len() int                   ***REMOVED*** return len(p.List) ***REMOVED***
func (p *FieldNumbers) Get(i int) pref.FieldNumber ***REMOVED*** return p.List[i] ***REMOVED***
func (p *FieldNumbers) Has(n pref.FieldNumber) bool ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.List) > 0 ***REMOVED***
			p.has = make(map[pref.FieldNumber]struct***REMOVED******REMOVED***, len(p.List))
			for _, n := range p.List ***REMOVED***
				p.has[n] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	_, ok := p.has[n]
	return ok
***REMOVED***
func (p *FieldNumbers) Format(s fmt.State, r rune)          ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *FieldNumbers) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***

type OneofFields struct ***REMOVED***
	List   []pref.FieldDescriptor
	once   sync.Once
	byName map[pref.Name]pref.FieldDescriptor        // protected by once
	byJSON map[string]pref.FieldDescriptor           // protected by once
	byNum  map[pref.FieldNumber]pref.FieldDescriptor // protected by once
***REMOVED***

func (p *OneofFields) Len() int                                         ***REMOVED*** return len(p.List) ***REMOVED***
func (p *OneofFields) Get(i int) pref.FieldDescriptor                   ***REMOVED*** return p.List[i] ***REMOVED***
func (p *OneofFields) ByName(s pref.Name) pref.FieldDescriptor          ***REMOVED*** return p.lazyInit().byName[s] ***REMOVED***
func (p *OneofFields) ByJSONName(s string) pref.FieldDescriptor         ***REMOVED*** return p.lazyInit().byJSON[s] ***REMOVED***
func (p *OneofFields) ByNumber(n pref.FieldNumber) pref.FieldDescriptor ***REMOVED*** return p.lazyInit().byNum[n] ***REMOVED***
func (p *OneofFields) Format(s fmt.State, r rune)                       ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *OneofFields) ProtoInternal(pragma.DoNotImplement)              ***REMOVED******REMOVED***

func (p *OneofFields) lazyInit() *OneofFields ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.List) > 0 ***REMOVED***
			p.byName = make(map[pref.Name]pref.FieldDescriptor, len(p.List))
			p.byJSON = make(map[string]pref.FieldDescriptor, len(p.List))
			p.byNum = make(map[pref.FieldNumber]pref.FieldDescriptor, len(p.List))
			for _, f := range p.List ***REMOVED***
				// Field names and numbers are guaranteed to be unique.
				p.byName[f.Name()] = f
				p.byJSON[f.JSONName()] = f
				p.byNum[f.Number()] = f
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return p
***REMOVED***

type SourceLocations struct ***REMOVED***
	List []pref.SourceLocation
***REMOVED***

func (p *SourceLocations) Len() int                            ***REMOVED*** return len(p.List) ***REMOVED***
func (p *SourceLocations) Get(i int) pref.SourceLocation       ***REMOVED*** return p.List[i] ***REMOVED***
func (p *SourceLocations) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***
