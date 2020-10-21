package dynamic

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/desc"
)

// ExtensionRegistry is a registry of known extension fields. This is used to parse
// extension fields encountered when de-serializing a dynamic message.
type ExtensionRegistry struct ***REMOVED***
	includeDefault bool
	mu             sync.RWMutex
	exts           map[string]map[int32]*desc.FieldDescriptor
***REMOVED***

// NewExtensionRegistryWithDefaults is a registry that includes all "default" extensions,
// which are those that are statically linked into the current program (e.g. registered by
// protoc-generated code via proto.RegisterExtension). Extensions explicitly added to the
// registry will override any default extensions that are for the same extendee and have the
// same tag number and/or name.
func NewExtensionRegistryWithDefaults() *ExtensionRegistry ***REMOVED***
	return &ExtensionRegistry***REMOVED***includeDefault: true***REMOVED***
***REMOVED***

// AddExtensionDesc adds the given extensions to the registry.
func (r *ExtensionRegistry) AddExtensionDesc(exts ...*proto.ExtensionDesc) error ***REMOVED***
	flds := make([]*desc.FieldDescriptor, len(exts))
	for i, ext := range exts ***REMOVED***
		fd, err := desc.LoadFieldDescriptorForExtension(ext)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		flds[i] = fd
	***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.exts == nil ***REMOVED***
		r.exts = map[string]map[int32]*desc.FieldDescriptor***REMOVED******REMOVED***
	***REMOVED***
	for _, fd := range flds ***REMOVED***
		r.putExtensionLocked(fd)
	***REMOVED***
	return nil
***REMOVED***

// AddExtension adds the given extensions to the registry. The given extensions
// will overwrite any previously added extensions that are for the same extendee
// message and same extension tag number.
func (r *ExtensionRegistry) AddExtension(exts ...*desc.FieldDescriptor) error ***REMOVED***
	for _, ext := range exts ***REMOVED***
		if !ext.IsExtension() ***REMOVED***
			return fmt.Errorf("given field is not an extension: %s", ext.GetFullyQualifiedName())
		***REMOVED***
	***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.exts == nil ***REMOVED***
		r.exts = map[string]map[int32]*desc.FieldDescriptor***REMOVED******REMOVED***
	***REMOVED***
	for _, ext := range exts ***REMOVED***
		r.putExtensionLocked(ext)
	***REMOVED***
	return nil
***REMOVED***

// AddExtensionsFromFile adds to the registry all extension fields defined in the given file descriptor.
func (r *ExtensionRegistry) AddExtensionsFromFile(fd *desc.FileDescriptor) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	r.addExtensionsFromFileLocked(fd, false, nil)
***REMOVED***

// AddExtensionsFromFileRecursively adds to the registry all extension fields defined in the give file
// descriptor and also recursively adds all extensions defined in that file's dependencies. This adds
// extensions from the entire transitive closure for the given file.
func (r *ExtensionRegistry) AddExtensionsFromFileRecursively(fd *desc.FileDescriptor) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	already := map[*desc.FileDescriptor]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	r.addExtensionsFromFileLocked(fd, true, already)
***REMOVED***

func (r *ExtensionRegistry) addExtensionsFromFileLocked(fd *desc.FileDescriptor, recursive bool, alreadySeen map[*desc.FileDescriptor]struct***REMOVED******REMOVED***) ***REMOVED***
	if _, ok := alreadySeen[fd]; ok ***REMOVED***
		return
	***REMOVED***

	if r.exts == nil ***REMOVED***
		r.exts = map[string]map[int32]*desc.FieldDescriptor***REMOVED******REMOVED***
	***REMOVED***
	for _, ext := range fd.GetExtensions() ***REMOVED***
		r.putExtensionLocked(ext)
	***REMOVED***
	for _, msg := range fd.GetMessageTypes() ***REMOVED***
		r.addExtensionsFromMessageLocked(msg)
	***REMOVED***

	if recursive ***REMOVED***
		alreadySeen[fd] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, dep := range fd.GetDependencies() ***REMOVED***
			r.addExtensionsFromFileLocked(dep, recursive, alreadySeen)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *ExtensionRegistry) addExtensionsFromMessageLocked(md *desc.MessageDescriptor) ***REMOVED***
	for _, ext := range md.GetNestedExtensions() ***REMOVED***
		r.putExtensionLocked(ext)
	***REMOVED***
	for _, msg := range md.GetNestedMessageTypes() ***REMOVED***
		r.addExtensionsFromMessageLocked(msg)
	***REMOVED***
***REMOVED***

func (r *ExtensionRegistry) putExtensionLocked(fd *desc.FieldDescriptor) ***REMOVED***
	msgName := fd.GetOwner().GetFullyQualifiedName()
	m := r.exts[msgName]
	if m == nil ***REMOVED***
		m = map[int32]*desc.FieldDescriptor***REMOVED******REMOVED***
		r.exts[msgName] = m
	***REMOVED***
	m[fd.GetNumber()] = fd
***REMOVED***

// FindExtension queries for the extension field with the given extendee name (must be a fully-qualified
// message name) and tag number. If no extension is known, nil is returned.
func (r *ExtensionRegistry) FindExtension(messageName string, tagNumber int32) *desc.FieldDescriptor ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	r.mu.RLock()
	defer r.mu.RUnlock()
	fd := r.exts[messageName][tagNumber]
	if fd == nil && r.includeDefault ***REMOVED***
		ext := getDefaultExtensions(messageName)[tagNumber]
		if ext != nil ***REMOVED***
			fd, _ = desc.LoadFieldDescriptorForExtension(ext)
		***REMOVED***
	***REMOVED***
	return fd
***REMOVED***

// FindExtensionByName queries for the extension field with the given extendee name (must be a fully-qualified
// message name) and field name (must also be a fully-qualified extension name). If no extension is known, nil
// is returned.
func (r *ExtensionRegistry) FindExtensionByName(messageName string, fieldName string) *desc.FieldDescriptor ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, fd := range r.exts[messageName] ***REMOVED***
		if fd.GetFullyQualifiedName() == fieldName ***REMOVED***
			return fd
		***REMOVED***
	***REMOVED***
	if r.includeDefault ***REMOVED***
		for _, ext := range getDefaultExtensions(messageName) ***REMOVED***
			fd, _ := desc.LoadFieldDescriptorForExtension(ext)
			if fd.GetFullyQualifiedName() == fieldName ***REMOVED***
				return fd
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// FindExtensionByJSONName queries for the extension field with the given extendee name (must be a fully-qualified
// message name) and JSON field name (must also be a fully-qualified name). If no extension is known, nil is returned.
// The fully-qualified JSON name is the same as the extension's normal fully-qualified name except that the last
// component uses the field's JSON name (if present).
func (r *ExtensionRegistry) FindExtensionByJSONName(messageName string, fieldName string) *desc.FieldDescriptor ***REMOVED***
	if r == nil ***REMOVED***
		return nil
	***REMOVED***
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, fd := range r.exts[messageName] ***REMOVED***
		if fd.GetFullyQualifiedJSONName() == fieldName ***REMOVED***
			return fd
		***REMOVED***
	***REMOVED***
	if r.includeDefault ***REMOVED***
		for _, ext := range getDefaultExtensions(messageName) ***REMOVED***
			fd, _ := desc.LoadFieldDescriptorForExtension(ext)
			if fd.GetFullyQualifiedJSONName() == fieldName ***REMOVED***
				return fd
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func getDefaultExtensions(messageName string) map[int32]*proto.ExtensionDesc ***REMOVED***
	t := proto.MessageType(messageName)
	if t != nil ***REMOVED***
		msg := reflect.Zero(t).Interface().(proto.Message)
		return proto.RegisteredExtensions(msg)
	***REMOVED***
	return nil
***REMOVED***

// AllExtensionsForType returns all known extension fields for the given extendee name (must be a
// fully-qualified message name).
func (r *ExtensionRegistry) AllExtensionsForType(messageName string) []*desc.FieldDescriptor ***REMOVED***
	if r == nil ***REMOVED***
		return []*desc.FieldDescriptor(nil)
	***REMOVED***
	r.mu.RLock()
	defer r.mu.RUnlock()
	flds := r.exts[messageName]
	var ret []*desc.FieldDescriptor
	if r.includeDefault ***REMOVED***
		exts := getDefaultExtensions(messageName)
		if len(exts) > 0 || len(flds) > 0 ***REMOVED***
			ret = make([]*desc.FieldDescriptor, 0, len(exts)+len(flds))
		***REMOVED***
		for tag, ext := range exts ***REMOVED***
			if _, ok := flds[tag]; ok ***REMOVED***
				// skip default extension and use the one explicitly registered instead
				continue
			***REMOVED***
			fd, _ := desc.LoadFieldDescriptorForExtension(ext)
			if fd != nil ***REMOVED***
				ret = append(ret, fd)
			***REMOVED***
		***REMOVED***
	***REMOVED*** else if len(flds) > 0 ***REMOVED***
		ret = make([]*desc.FieldDescriptor, 0, len(flds))
	***REMOVED***

	for _, ext := range flds ***REMOVED***
		ret = append(ret, ext)
	***REMOVED***
	return ret
***REMOVED***
