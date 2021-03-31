// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"fmt"
	"math"
	"sort"
	"sync"

	"google.golang.org/protobuf/internal/genid"

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
	return protowire.MinValidNumber <= n && (n <= protowire.MaxValidNumber || isMessageSet)
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
	byText map[string]pref.FieldDescriptor           // protected by once
	byNum  map[pref.FieldNumber]pref.FieldDescriptor // protected by once
***REMOVED***

func (p *OneofFields) Len() int                                         ***REMOVED*** return len(p.List) ***REMOVED***
func (p *OneofFields) Get(i int) pref.FieldDescriptor                   ***REMOVED*** return p.List[i] ***REMOVED***
func (p *OneofFields) ByName(s pref.Name) pref.FieldDescriptor          ***REMOVED*** return p.lazyInit().byName[s] ***REMOVED***
func (p *OneofFields) ByJSONName(s string) pref.FieldDescriptor         ***REMOVED*** return p.lazyInit().byJSON[s] ***REMOVED***
func (p *OneofFields) ByTextName(s string) pref.FieldDescriptor         ***REMOVED*** return p.lazyInit().byText[s] ***REMOVED***
func (p *OneofFields) ByNumber(n pref.FieldNumber) pref.FieldDescriptor ***REMOVED*** return p.lazyInit().byNum[n] ***REMOVED***
func (p *OneofFields) Format(s fmt.State, r rune)                       ***REMOVED*** descfmt.FormatList(s, r, p) ***REMOVED***
func (p *OneofFields) ProtoInternal(pragma.DoNotImplement)              ***REMOVED******REMOVED***

func (p *OneofFields) lazyInit() *OneofFields ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.List) > 0 ***REMOVED***
			p.byName = make(map[pref.Name]pref.FieldDescriptor, len(p.List))
			p.byJSON = make(map[string]pref.FieldDescriptor, len(p.List))
			p.byText = make(map[string]pref.FieldDescriptor, len(p.List))
			p.byNum = make(map[pref.FieldNumber]pref.FieldDescriptor, len(p.List))
			for _, f := range p.List ***REMOVED***
				// Field names and numbers are guaranteed to be unique.
				p.byName[f.Name()] = f
				p.byJSON[f.JSONName()] = f
				p.byText[f.TextName()] = f
				p.byNum[f.Number()] = f
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return p
***REMOVED***

type SourceLocations struct ***REMOVED***
	// List is a list of SourceLocations.
	// The SourceLocation.Next field does not need to be populated
	// as it will be lazily populated upon first need.
	List []pref.SourceLocation

	// File is the parent file descriptor that these locations are relative to.
	// If non-nil, ByDescriptor verifies that the provided descriptor
	// is a child of this file descriptor.
	File pref.FileDescriptor

	once   sync.Once
	byPath map[pathKey]int
***REMOVED***

func (p *SourceLocations) Len() int                      ***REMOVED*** return len(p.List) ***REMOVED***
func (p *SourceLocations) Get(i int) pref.SourceLocation ***REMOVED*** return p.lazyInit().List[i] ***REMOVED***
func (p *SourceLocations) byKey(k pathKey) pref.SourceLocation ***REMOVED***
	if i, ok := p.lazyInit().byPath[k]; ok ***REMOVED***
		return p.List[i]
	***REMOVED***
	return pref.SourceLocation***REMOVED******REMOVED***
***REMOVED***
func (p *SourceLocations) ByPath(path pref.SourcePath) pref.SourceLocation ***REMOVED***
	return p.byKey(newPathKey(path))
***REMOVED***
func (p *SourceLocations) ByDescriptor(desc pref.Descriptor) pref.SourceLocation ***REMOVED***
	if p.File != nil && desc != nil && p.File != desc.ParentFile() ***REMOVED***
		return pref.SourceLocation***REMOVED******REMOVED*** // mismatching parent files
	***REMOVED***
	var pathArr [16]int32
	path := pathArr[:0]
	for ***REMOVED***
		switch desc.(type) ***REMOVED***
		case pref.FileDescriptor:
			// Reverse the path since it was constructed in reverse.
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 ***REMOVED***
				path[i], path[j] = path[j], path[i]
			***REMOVED***
			return p.byKey(newPathKey(path))
		case pref.MessageDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_MessageType_field_number))
			case pref.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_NestedType_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case pref.FieldDescriptor:
			isExtension := desc.(pref.FieldDescriptor).IsExtension()
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			if isExtension ***REMOVED***
				switch desc.(type) ***REMOVED***
				case pref.FileDescriptor:
					path = append(path, int32(genid.FileDescriptorProto_Extension_field_number))
				case pref.MessageDescriptor:
					path = append(path, int32(genid.DescriptorProto_Extension_field_number))
				default:
					return pref.SourceLocation***REMOVED******REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch desc.(type) ***REMOVED***
				case pref.MessageDescriptor:
					path = append(path, int32(genid.DescriptorProto_Field_field_number))
				default:
					return pref.SourceLocation***REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
		case pref.OneofDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_OneofDecl_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case pref.EnumDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_EnumType_field_number))
			case pref.MessageDescriptor:
				path = append(path, int32(genid.DescriptorProto_EnumType_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case pref.EnumValueDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.EnumDescriptor:
				path = append(path, int32(genid.EnumDescriptorProto_Value_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case pref.ServiceDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.FileDescriptor:
				path = append(path, int32(genid.FileDescriptorProto_Service_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case pref.MethodDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case pref.ServiceDescriptor:
				path = append(path, int32(genid.ServiceDescriptorProto_Method_field_number))
			default:
				return pref.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		default:
			return pref.SourceLocation***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
func (p *SourceLocations) lazyInit() *SourceLocations ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.List) > 0 ***REMOVED***
			// Collect all the indexes for a given path.
			pathIdxs := make(map[pathKey][]int, len(p.List))
			for i, l := range p.List ***REMOVED***
				k := newPathKey(l.Path)
				pathIdxs[k] = append(pathIdxs[k], i)
			***REMOVED***

			// Update the next index for all locations.
			p.byPath = make(map[pathKey]int, len(p.List))
			for k, idxs := range pathIdxs ***REMOVED***
				for i := 0; i < len(idxs)-1; i++ ***REMOVED***
					p.List[idxs[i]].Next = idxs[i+1]
				***REMOVED***
				p.List[idxs[len(idxs)-1]].Next = 0
				p.byPath[k] = idxs[0] // record the first location for this path
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return p
***REMOVED***
func (p *SourceLocations) ProtoInternal(pragma.DoNotImplement) ***REMOVED******REMOVED***

// pathKey is a comparable representation of protoreflect.SourcePath.
type pathKey struct ***REMOVED***
	arr [16]uint8 // first n-1 path segments; last element is the length
	str string    // used if the path does not fit in arr
***REMOVED***

func newPathKey(p pref.SourcePath) (k pathKey) ***REMOVED***
	if len(p) < len(k.arr) ***REMOVED***
		for i, ps := range p ***REMOVED***
			if ps < 0 || math.MaxUint8 <= ps ***REMOVED***
				return pathKey***REMOVED***str: p.String()***REMOVED***
			***REMOVED***
			k.arr[i] = uint8(ps)
		***REMOVED***
		k.arr[len(k.arr)-1] = uint8(len(p))
		return k
	***REMOVED***
	return pathKey***REMOVED***str: p.String()***REMOVED***
***REMOVED***
