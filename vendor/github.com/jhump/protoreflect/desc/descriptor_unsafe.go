//+build !appengine,!gopherjs,!purego
// NB: other environments where unsafe is unappropriate should use "purego" build tag
// https://github.com/golang/go/issues/23172

package desc

import (
	"sync/atomic"
	"unsafe"
)

type jsonNameMap map[string]*FieldDescriptor // loaded/stored atomically via atomic+unsafe
type memoizedDefault *interface***REMOVED******REMOVED***            // loaded/stored atomically via atomic+unsafe

// FindFieldByJSONName finds the field with the given JSON field name. If no such
// field exists then nil is returned. Only regular fields are returned, not
// extensions.
func (md *MessageDescriptor) FindFieldByJSONName(jsonName string) *FieldDescriptor ***REMOVED***
	// NB: We don't want to eagerly index JSON names because many programs won't use it.
	// So we want to do it lazily, but also make sure the result is thread-safe. So we
	// atomically load/store the map as if it were a normal pointer. We don't use other
	// mechanisms -- like sync.Mutex, sync.RWMutex, sync.Once, or atomic.Value -- to
	// do this lazily because those types cannot be copied, and we'd rather not induce
	// 'go vet' errors in programs that use descriptors and try to copy them.
	// If multiple goroutines try to access the index at the same time, before it is
	// built, they will all end up computing the index redundantly. Future reads of
	// the index will use whatever was the "last one stored" by those racing goroutines.
	// Since building the index is deterministic, this is fine: all indices computed
	// will be the same.
	addrOfJsonNames := (*unsafe.Pointer)(unsafe.Pointer(&md.jsonNames))
	jsonNames := atomic.LoadPointer(addrOfJsonNames)
	var index map[string]*FieldDescriptor
	if jsonNames == nil ***REMOVED***
		// slow path: compute the index
		index = map[string]*FieldDescriptor***REMOVED******REMOVED***
		for _, f := range md.fields ***REMOVED***
			jn := f.GetJSONName()
			index[jn] = f
		***REMOVED***
		atomic.StorePointer(addrOfJsonNames, *(*unsafe.Pointer)(unsafe.Pointer(&index)))
	***REMOVED*** else ***REMOVED***
		*(*unsafe.Pointer)(unsafe.Pointer(&index)) = jsonNames
	***REMOVED***
	return index[jsonName]
***REMOVED***

func (fd *FieldDescriptor) getDefaultValue() interface***REMOVED******REMOVED*** ***REMOVED***
	addrOfDef := (*unsafe.Pointer)(unsafe.Pointer(&fd.def))
	def := atomic.LoadPointer(addrOfDef)
	if def != nil ***REMOVED***
		return *(*interface***REMOVED******REMOVED***)(def)
	***REMOVED***
	// slow path: compute the default, potentially involves decoding value
	d := fd.determineDefault()
	atomic.StorePointer(addrOfDef, (unsafe.Pointer(&d)))
	return d
***REMOVED***
