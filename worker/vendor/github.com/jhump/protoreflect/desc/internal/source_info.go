package internal

import (
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

// SourceInfoMap is a map of paths in a descriptor to the corresponding source
// code info.
type SourceInfoMap map[string][]*dpb.SourceCodeInfo_Location

// Get returns the source code info for the given path. If there are
// multiple locations for the same path, the first one is returned.
func (m SourceInfoMap) Get(path []int32) *dpb.SourceCodeInfo_Location ***REMOVED***
	v := m[asMapKey(path)]
	if len(v) > 0 ***REMOVED***
		return v[0]
	***REMOVED***
	return nil
***REMOVED***

// GetAll returns all source code info for the given path.
func (m SourceInfoMap) GetAll(path []int32) []*dpb.SourceCodeInfo_Location ***REMOVED***
	return m[asMapKey(path)]
***REMOVED***

// Add stores the given source code info for the given path.
func (m SourceInfoMap) Add(path []int32, loc *dpb.SourceCodeInfo_Location) ***REMOVED***
	m[asMapKey(path)] = append(m[asMapKey(path)], loc)
***REMOVED***

// PutIfAbsent stores the given source code info for the given path only if the
// given path does not exist in the map. This method returns true when the value
// is stored, false if the path already exists.
func (m SourceInfoMap) PutIfAbsent(path []int32, loc *dpb.SourceCodeInfo_Location) bool ***REMOVED***
	k := asMapKey(path)
	if _, ok := m[k]; ok ***REMOVED***
		return false
	***REMOVED***
	m[k] = []*dpb.SourceCodeInfo_Location***REMOVED***loc***REMOVED***
	return true
***REMOVED***

func asMapKey(slice []int32) string ***REMOVED***
	// NB: arrays should be usable as map keys, but this does not
	// work due to a bug: https://github.com/golang/go/issues/22605
	//rv := reflect.ValueOf(slice)
	//arrayType := reflect.ArrayOf(rv.Len(), rv.Type().Elem())
	//array := reflect.New(arrayType).Elem()
	//reflect.Copy(array, rv)
	//return array.Interface()

	b := make([]byte, len(slice)*4)
	j := 0
	for _, s := range slice ***REMOVED***
		b[j] = byte(s)
		b[j+1] = byte(s >> 8)
		b[j+2] = byte(s >> 16)
		b[j+3] = byte(s >> 24)
		j += 4
	***REMOVED***
	return string(b)
***REMOVED***

// CreateSourceInfoMap constructs a new SourceInfoMap and populates it with the
// source code info in the given file descriptor proto.
func CreateSourceInfoMap(fd *dpb.FileDescriptorProto) SourceInfoMap ***REMOVED***
	res := SourceInfoMap***REMOVED******REMOVED***
	PopulateSourceInfoMap(fd, res)
	return res
***REMOVED***

// PopulateSourceInfoMap populates the given SourceInfoMap with information from
// the given file descriptor.
func PopulateSourceInfoMap(fd *dpb.FileDescriptorProto, m SourceInfoMap) ***REMOVED***
	for _, l := range fd.GetSourceCodeInfo().GetLocation() ***REMOVED***
		m.Add(l.Path, l)
	***REMOVED***
***REMOVED***

// NB: This wonkiness allows desc.Descriptor impl to implement an interface that
// is only usable from this package, by embedding a SourceInfoComputeFunc that
// implements the actual logic (which must live in desc package to avoid a
// dependency cycle).

// SourceInfoComputer is a single method which will be invoked to recompute
// source info. This is needed for the protoparse package, which needs to link
// descriptors without source info in order to interpret options, but then needs
// to re-compute source info after that interpretation so that final linked
// descriptors expose the right info.
type SourceInfoComputer interface ***REMOVED***
	recomputeSourceInfo()
***REMOVED***

// SourceInfoComputeFunc is the type that a desc.Descriptor will embed. It will
// be aliased in the desc package to an unexported name so it is not marked as
// an exported field in reflection and not present in Go docs.
type SourceInfoComputeFunc func()

func (f SourceInfoComputeFunc) recomputeSourceInfo() ***REMOVED***
	f()
***REMOVED***

// RecomputeSourceInfo is used to initiate recomputation of source info. This is
// is used by the protoparse package, after it interprets options.
func RecomputeSourceInfo(c SourceInfoComputer) ***REMOVED***
	c.recomputeSourceInfo()
***REMOVED***
