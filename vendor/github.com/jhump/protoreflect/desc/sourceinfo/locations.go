package sourceinfo

import (
	"math"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/jhump/protoreflect/desc/internal"
)

// NB: forked from google.golang.org/protobuf/internal/filedesc
type sourceLocations struct ***REMOVED***
	protoreflect.SourceLocations

	orig []*descriptorpb.SourceCodeInfo_Location
	// locs is a list of sourceLocations.
	// The SourceLocation.Next field does not need to be populated
	// as it will be lazily populated upon first need.
	locs []protoreflect.SourceLocation

	// fd is the parent file descriptor that these locations are relative to.
	// If non-nil, ByDescriptor verifies that the provided descriptor
	// is a child of this file descriptor.
	fd protoreflect.FileDescriptor

	once   sync.Once
	byPath map[pathKey]int
***REMOVED***

func (p *sourceLocations) Len() int ***REMOVED*** return len(p.orig) ***REMOVED***
func (p *sourceLocations) Get(i int) protoreflect.SourceLocation ***REMOVED***
	return p.lazyInit().locs[i]
***REMOVED***
func (p *sourceLocations) byKey(k pathKey) protoreflect.SourceLocation ***REMOVED***
	if i, ok := p.lazyInit().byPath[k]; ok ***REMOVED***
		return p.locs[i]
	***REMOVED***
	return protoreflect.SourceLocation***REMOVED******REMOVED***
***REMOVED***
func (p *sourceLocations) ByPath(path protoreflect.SourcePath) protoreflect.SourceLocation ***REMOVED***
	return p.byKey(newPathKey(path))
***REMOVED***
func (p *sourceLocations) ByDescriptor(desc protoreflect.Descriptor) protoreflect.SourceLocation ***REMOVED***
	if p.fd != nil && desc != nil && p.fd != desc.ParentFile() ***REMOVED***
		return protoreflect.SourceLocation***REMOVED******REMOVED*** // mismatching parent imports
	***REMOVED***
	var pathArr [16]int32
	path := pathArr[:0]
	for ***REMOVED***
		switch desc.(type) ***REMOVED***
		case protoreflect.FileDescriptor:
			// Reverse the path since it was constructed in reverse.
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 ***REMOVED***
				path[i], path[j] = path[j], path[i]
			***REMOVED***
			return p.byKey(newPathKey(path))
		case protoreflect.MessageDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.FileDescriptor:
				path = append(path, int32(internal.File_messagesTag))
			case protoreflect.MessageDescriptor:
				path = append(path, int32(internal.Message_nestedMessagesTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case protoreflect.FieldDescriptor:
			isExtension := desc.(protoreflect.FieldDescriptor).IsExtension()
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			if isExtension ***REMOVED***
				switch desc.(type) ***REMOVED***
				case protoreflect.FileDescriptor:
					path = append(path, int32(internal.File_extensionsTag))
				case protoreflect.MessageDescriptor:
					path = append(path, int32(internal.Message_extensionsTag))
				default:
					return protoreflect.SourceLocation***REMOVED******REMOVED***
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				switch desc.(type) ***REMOVED***
				case protoreflect.MessageDescriptor:
					path = append(path, int32(internal.Message_fieldsTag))
				default:
					return protoreflect.SourceLocation***REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
		case protoreflect.OneofDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.MessageDescriptor:
				path = append(path, int32(internal.Message_oneOfsTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case protoreflect.EnumDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.FileDescriptor:
				path = append(path, int32(internal.File_enumsTag))
			case protoreflect.MessageDescriptor:
				path = append(path, int32(internal.Message_enumsTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case protoreflect.EnumValueDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.EnumDescriptor:
				path = append(path, int32(internal.Enum_valuesTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case protoreflect.ServiceDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.FileDescriptor:
				path = append(path, int32(internal.File_servicesTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		case protoreflect.MethodDescriptor:
			path = append(path, int32(desc.Index()))
			desc = desc.Parent()
			switch desc.(type) ***REMOVED***
			case protoreflect.ServiceDescriptor:
				path = append(path, int32(internal.Service_methodsTag))
			default:
				return protoreflect.SourceLocation***REMOVED******REMOVED***
			***REMOVED***
		default:
			return protoreflect.SourceLocation***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
func (p *sourceLocations) lazyInit() *sourceLocations ***REMOVED***
	p.once.Do(func() ***REMOVED***
		if len(p.orig) > 0 ***REMOVED***
			p.locs = make([]protoreflect.SourceLocation, len(p.orig))
			// Collect all the indexes for a given path.
			pathIdxs := make(map[pathKey][]int, len(p.locs))
			for i := range p.orig ***REMOVED***
				l := asSourceLocation(p.orig[i])
				p.locs[i] = l
				k := newPathKey(l.Path)
				pathIdxs[k] = append(pathIdxs[k], i)
			***REMOVED***

			// Update the next index for all locations.
			p.byPath = make(map[pathKey]int, len(p.locs))
			for k, idxs := range pathIdxs ***REMOVED***
				for i := 0; i < len(idxs)-1; i++ ***REMOVED***
					p.locs[idxs[i]].Next = idxs[i+1]
				***REMOVED***
				p.locs[idxs[len(idxs)-1]].Next = 0
				p.byPath[k] = idxs[0] // record the first location for this path
			***REMOVED***
		***REMOVED***
	***REMOVED***)
	return p
***REMOVED***

func asSourceLocation(l *descriptorpb.SourceCodeInfo_Location) protoreflect.SourceLocation ***REMOVED***
	endLine := l.Span[0]
	endCol := l.Span[2]
	if len(l.Span) > 3 ***REMOVED***
		endLine = l.Span[2]
		endCol = l.Span[3]
	***REMOVED***
	return protoreflect.SourceLocation***REMOVED***
		Path:                    l.Path,
		StartLine:               int(l.Span[0]),
		StartColumn:             int(l.Span[1]),
		EndLine:                 int(endLine),
		EndColumn:               int(endCol),
		LeadingDetachedComments: l.LeadingDetachedComments,
		LeadingComments:         l.GetLeadingComments(),
		TrailingComments:        l.GetTrailingComments(),
	***REMOVED***
***REMOVED***

// pathKey is a comparable representation of protoreflect.SourcePath.
type pathKey struct ***REMOVED***
	arr [16]uint8 // first n-1 path segments; last element is the length
	str string    // used if the path does not fit in arr
***REMOVED***

func newPathKey(p protoreflect.SourcePath) (k pathKey) ***REMOVED***
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
