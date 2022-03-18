package desc

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/internal"
)

// Descriptor is the common interface implemented by all descriptor objects.
type Descriptor interface ***REMOVED***
	// GetName returns the name of the object described by the descriptor. This will
	// be a base name that does not include enclosing message names or the package name.
	// For file descriptors, this indicates the path and name to the described file.
	GetName() string
	// GetFullyQualifiedName returns the fully-qualified name of the object described by
	// the descriptor. This will include the package name and any enclosing message names.
	// For file descriptors, this returns the path and name to the described file (same as
	// GetName).
	GetFullyQualifiedName() string
	// GetParent returns the enclosing element in a proto source file. If the described
	// object is a top-level object, this returns the file descriptor. Otherwise, it returns
	// the element in which the described object was declared. File descriptors have no
	// parent and return nil.
	GetParent() Descriptor
	// GetFile returns the file descriptor in which this element was declared. File
	// descriptors return themselves.
	GetFile() *FileDescriptor
	// GetOptions returns the options proto containing options for the described element.
	GetOptions() proto.Message
	// GetSourceInfo returns any source code information that was present in the file
	// descriptor. Source code info is optional. If no source code info is available for
	// the element (including if there is none at all in the file descriptor) then this
	// returns nil
	GetSourceInfo() *dpb.SourceCodeInfo_Location
	// AsProto returns the underlying descriptor proto for this descriptor.
	AsProto() proto.Message
***REMOVED***

type sourceInfoRecomputeFunc = internal.SourceInfoComputeFunc

// FileDescriptor describes a proto source file.
type FileDescriptor struct ***REMOVED***
	proto      *dpb.FileDescriptorProto
	symbols    map[string]Descriptor
	deps       []*FileDescriptor
	publicDeps []*FileDescriptor
	weakDeps   []*FileDescriptor
	messages   []*MessageDescriptor
	enums      []*EnumDescriptor
	extensions []*FieldDescriptor
	services   []*ServiceDescriptor
	fieldIndex map[string]map[int32]*FieldDescriptor
	isProto3   bool
	sourceInfo internal.SourceInfoMap
	sourceInfoRecomputeFunc
***REMOVED***

func (fd *FileDescriptor) recomputeSourceInfo() ***REMOVED***
	internal.PopulateSourceInfoMap(fd.proto, fd.sourceInfo)
***REMOVED***

func (fd *FileDescriptor) registerField(field *FieldDescriptor) ***REMOVED***
	fields := fd.fieldIndex[field.owner.GetFullyQualifiedName()]
	if fields == nil ***REMOVED***
		fields = map[int32]*FieldDescriptor***REMOVED******REMOVED***
		fd.fieldIndex[field.owner.GetFullyQualifiedName()] = fields
	***REMOVED***
	fields[field.GetNumber()] = field
***REMOVED***

// GetName returns the name of the file, as it was given to the protoc invocation
// to compile it, possibly including path (relative to a directory in the proto
// import path).
func (fd *FileDescriptor) GetName() string ***REMOVED***
	return fd.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the name of the file, same as GetName. It is
// present to satisfy the Descriptor interface.
func (fd *FileDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return fd.proto.GetName()
***REMOVED***

// GetPackage returns the name of the package declared in the file.
func (fd *FileDescriptor) GetPackage() string ***REMOVED***
	return fd.proto.GetPackage()
***REMOVED***

// GetParent always returns nil: files are the root of descriptor hierarchies.
// Is it present to satisfy the Descriptor interface.
func (fd *FileDescriptor) GetParent() Descriptor ***REMOVED***
	return nil
***REMOVED***

// GetFile returns the receiver, which is a file descriptor. This is present
// to satisfy the Descriptor interface.
func (fd *FileDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return fd
***REMOVED***

// GetOptions returns the file's options. Most usages will be more interested
// in GetFileOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (fd *FileDescriptor) GetOptions() proto.Message ***REMOVED***
	return fd.proto.GetOptions()
***REMOVED***

// GetFileOptions returns the file's options.
func (fd *FileDescriptor) GetFileOptions() *dpb.FileOptions ***REMOVED***
	return fd.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns nil for files. It is present to satisfy the Descriptor
// interface.
func (fd *FileDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return nil
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsFileDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (fd *FileDescriptor) AsProto() proto.Message ***REMOVED***
	return fd.proto
***REMOVED***

// AsFileDescriptorProto returns the underlying descriptor proto.
func (fd *FileDescriptor) AsFileDescriptorProto() *dpb.FileDescriptorProto ***REMOVED***
	return fd.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (fd *FileDescriptor) String() string ***REMOVED***
	return fd.proto.String()
***REMOVED***

// IsProto3 returns true if the file declares a syntax of "proto3".
func (fd *FileDescriptor) IsProto3() bool ***REMOVED***
	return fd.isProto3
***REMOVED***

// GetDependencies returns all of this file's dependencies. These correspond to
// import statements in the file.
func (fd *FileDescriptor) GetDependencies() []*FileDescriptor ***REMOVED***
	return fd.deps
***REMOVED***

// GetPublicDependencies returns all of this file's public dependencies. These
// correspond to public import statements in the file.
func (fd *FileDescriptor) GetPublicDependencies() []*FileDescriptor ***REMOVED***
	return fd.publicDeps
***REMOVED***

// GetWeakDependencies returns all of this file's weak dependencies. These
// correspond to weak import statements in the file.
func (fd *FileDescriptor) GetWeakDependencies() []*FileDescriptor ***REMOVED***
	return fd.weakDeps
***REMOVED***

// GetMessageTypes returns all top-level messages declared in this file.
func (fd *FileDescriptor) GetMessageTypes() []*MessageDescriptor ***REMOVED***
	return fd.messages
***REMOVED***

// GetEnumTypes returns all top-level enums declared in this file.
func (fd *FileDescriptor) GetEnumTypes() []*EnumDescriptor ***REMOVED***
	return fd.enums
***REMOVED***

// GetExtensions returns all top-level extensions declared in this file.
func (fd *FileDescriptor) GetExtensions() []*FieldDescriptor ***REMOVED***
	return fd.extensions
***REMOVED***

// GetServices returns all services declared in this file.
func (fd *FileDescriptor) GetServices() []*ServiceDescriptor ***REMOVED***
	return fd.services
***REMOVED***

// FindSymbol returns the descriptor contained within this file for the
// element with the given fully-qualified symbol name. If no such element
// exists then this method returns nil.
func (fd *FileDescriptor) FindSymbol(symbol string) Descriptor ***REMOVED***
	if len(symbol) == 0 ***REMOVED***
		return nil
	***REMOVED***
	if symbol[0] == '.' ***REMOVED***
		symbol = symbol[1:]
	***REMOVED***
	if ret := fd.symbols[symbol]; ret != nil ***REMOVED***
		return ret
	***REMOVED***

	// allow accessing symbols through public imports, too
	for _, dep := range fd.GetPublicDependencies() ***REMOVED***
		if ret := dep.FindSymbol(symbol); ret != nil ***REMOVED***
			return ret
		***REMOVED***
	***REMOVED***

	// not found
	return nil
***REMOVED***

// FindMessage finds the message with the given fully-qualified name. If no
// such element exists in this file then nil is returned.
func (fd *FileDescriptor) FindMessage(msgName string) *MessageDescriptor ***REMOVED***
	if md, ok := fd.symbols[msgName].(*MessageDescriptor); ok ***REMOVED***
		return md
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindEnum finds the enum with the given fully-qualified name. If no such
// element exists in this file then nil is returned.
func (fd *FileDescriptor) FindEnum(enumName string) *EnumDescriptor ***REMOVED***
	if ed, ok := fd.symbols[enumName].(*EnumDescriptor); ok ***REMOVED***
		return ed
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindService finds the service with the given fully-qualified name. If no
// such element exists in this file then nil is returned.
func (fd *FileDescriptor) FindService(serviceName string) *ServiceDescriptor ***REMOVED***
	if sd, ok := fd.symbols[serviceName].(*ServiceDescriptor); ok ***REMOVED***
		return sd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindExtension finds the extension field for the given extended type name and
// tag number. If no such element exists in this file then nil is returned.
func (fd *FileDescriptor) FindExtension(extendeeName string, tagNumber int32) *FieldDescriptor ***REMOVED***
	if exd, ok := fd.fieldIndex[extendeeName][tagNumber]; ok && exd.IsExtension() ***REMOVED***
		return exd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindExtensionByName finds the extension field with the given fully-qualified
// name. If no such element exists in this file then nil is returned.
func (fd *FileDescriptor) FindExtensionByName(extName string) *FieldDescriptor ***REMOVED***
	if exd, ok := fd.symbols[extName].(*FieldDescriptor); ok && exd.IsExtension() ***REMOVED***
		return exd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// MessageDescriptor describes a protocol buffer message.
type MessageDescriptor struct ***REMOVED***
	proto          *dpb.DescriptorProto
	parent         Descriptor
	file           *FileDescriptor
	fields         []*FieldDescriptor
	nested         []*MessageDescriptor
	enums          []*EnumDescriptor
	extensions     []*FieldDescriptor
	oneOfs         []*OneOfDescriptor
	extRanges      extRanges
	fqn            string
	sourceInfoPath []int32
	jsonNames      jsonNameMap
	isProto3       bool
	isMapEntry     bool
***REMOVED***

func createMessageDescriptor(fd *FileDescriptor, parent Descriptor, enclosing string, md *dpb.DescriptorProto, symbols map[string]Descriptor) (*MessageDescriptor, string) ***REMOVED***
	msgName := merge(enclosing, md.GetName())
	ret := &MessageDescriptor***REMOVED***proto: md, parent: parent, file: fd, fqn: msgName***REMOVED***
	for _, f := range md.GetField() ***REMOVED***
		fld, n := createFieldDescriptor(fd, ret, msgName, f)
		symbols[n] = fld
		ret.fields = append(ret.fields, fld)
	***REMOVED***
	for _, nm := range md.NestedType ***REMOVED***
		nmd, n := createMessageDescriptor(fd, ret, msgName, nm, symbols)
		symbols[n] = nmd
		ret.nested = append(ret.nested, nmd)
	***REMOVED***
	for _, e := range md.EnumType ***REMOVED***
		ed, n := createEnumDescriptor(fd, ret, msgName, e, symbols)
		symbols[n] = ed
		ret.enums = append(ret.enums, ed)
	***REMOVED***
	for _, ex := range md.GetExtension() ***REMOVED***
		exd, n := createFieldDescriptor(fd, ret, msgName, ex)
		symbols[n] = exd
		ret.extensions = append(ret.extensions, exd)
	***REMOVED***
	for i, o := range md.GetOneofDecl() ***REMOVED***
		od, n := createOneOfDescriptor(fd, ret, i, msgName, o)
		symbols[n] = od
		ret.oneOfs = append(ret.oneOfs, od)
	***REMOVED***
	for _, r := range md.GetExtensionRange() ***REMOVED***
		// proto.ExtensionRange is inclusive (and that's how extension ranges are defined in code).
		// but protoc converts range to exclusive end in descriptor, so we must convert back
		end := r.GetEnd() - 1
		ret.extRanges = append(ret.extRanges, proto.ExtensionRange***REMOVED***
			Start: r.GetStart(),
			End:   end***REMOVED***)
	***REMOVED***
	sort.Sort(ret.extRanges)
	ret.isProto3 = fd.isProto3
	ret.isMapEntry = md.GetOptions().GetMapEntry() &&
		len(ret.fields) == 2 &&
		ret.fields[0].GetNumber() == 1 &&
		ret.fields[1].GetNumber() == 2

	return ret, msgName
***REMOVED***

func (md *MessageDescriptor) resolve(path []int32, scopes []scope) error ***REMOVED***
	md.sourceInfoPath = append([]int32(nil), path...) // defensive copy
	path = append(path, internal.Message_nestedMessagesTag)
	scopes = append(scopes, messageScope(md))
	for i, nmd := range md.nested ***REMOVED***
		if err := nmd.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	path[len(path)-1] = internal.Message_enumsTag
	for i, ed := range md.enums ***REMOVED***
		ed.resolve(append(path, int32(i)))
	***REMOVED***
	path[len(path)-1] = internal.Message_fieldsTag
	for i, fld := range md.fields ***REMOVED***
		if err := fld.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	path[len(path)-1] = internal.Message_extensionsTag
	for i, exd := range md.extensions ***REMOVED***
		if err := exd.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	path[len(path)-1] = internal.Message_oneOfsTag
	for i, od := range md.oneOfs ***REMOVED***
		od.resolve(append(path, int32(i)))
	***REMOVED***
	return nil
***REMOVED***

// GetName returns the simple (unqualified) name of the message.
func (md *MessageDescriptor) GetName() string ***REMOVED***
	return md.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the message. This
// includes the package name (if there is one) as well as the names of any
// enclosing messages.
func (md *MessageDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return md.fqn
***REMOVED***

// GetParent returns the message's enclosing descriptor. For top-level messages,
// this will be a file descriptor. Otherwise it will be the descriptor for the
// enclosing message.
func (md *MessageDescriptor) GetParent() Descriptor ***REMOVED***
	return md.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this message is defined.
func (md *MessageDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return md.file
***REMOVED***

// GetOptions returns the message's options. Most usages will be more interested
// in GetMessageOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (md *MessageDescriptor) GetOptions() proto.Message ***REMOVED***
	return md.proto.GetOptions()
***REMOVED***

// GetMessageOptions returns the message's options.
func (md *MessageDescriptor) GetMessageOptions() *dpb.MessageOptions ***REMOVED***
	return md.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the message, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// message was defined and also contains comments associated with the message
// definition.
func (md *MessageDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return md.file.sourceInfo.Get(md.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (md *MessageDescriptor) AsProto() proto.Message ***REMOVED***
	return md.proto
***REMOVED***

// AsDescriptorProto returns the underlying descriptor proto.
func (md *MessageDescriptor) AsDescriptorProto() *dpb.DescriptorProto ***REMOVED***
	return md.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (md *MessageDescriptor) String() string ***REMOVED***
	return md.proto.String()
***REMOVED***

// IsMapEntry returns true if this is a synthetic message type that represents an entry
// in a map field.
func (md *MessageDescriptor) IsMapEntry() bool ***REMOVED***
	return md.isMapEntry
***REMOVED***

// GetFields returns all of the fields for this message.
func (md *MessageDescriptor) GetFields() []*FieldDescriptor ***REMOVED***
	return md.fields
***REMOVED***

// GetNestedMessageTypes returns all of the message types declared inside this message.
func (md *MessageDescriptor) GetNestedMessageTypes() []*MessageDescriptor ***REMOVED***
	return md.nested
***REMOVED***

// GetNestedEnumTypes returns all of the enums declared inside this message.
func (md *MessageDescriptor) GetNestedEnumTypes() []*EnumDescriptor ***REMOVED***
	return md.enums
***REMOVED***

// GetNestedExtensions returns all of the extensions declared inside this message.
func (md *MessageDescriptor) GetNestedExtensions() []*FieldDescriptor ***REMOVED***
	return md.extensions
***REMOVED***

// GetOneOfs returns all of the one-of field sets declared inside this message.
func (md *MessageDescriptor) GetOneOfs() []*OneOfDescriptor ***REMOVED***
	return md.oneOfs
***REMOVED***

// IsProto3 returns true if the file in which this message is defined declares a syntax of "proto3".
func (md *MessageDescriptor) IsProto3() bool ***REMOVED***
	return md.isProto3
***REMOVED***

// GetExtensionRanges returns the ranges of extension field numbers for this message.
func (md *MessageDescriptor) GetExtensionRanges() []proto.ExtensionRange ***REMOVED***
	return md.extRanges
***REMOVED***

// IsExtendable returns true if this message has any extension ranges.
func (md *MessageDescriptor) IsExtendable() bool ***REMOVED***
	return len(md.extRanges) > 0
***REMOVED***

// IsExtension returns true if the given tag number is within any of this message's
// extension ranges.
func (md *MessageDescriptor) IsExtension(tagNumber int32) bool ***REMOVED***
	return md.extRanges.IsExtension(tagNumber)
***REMOVED***

type extRanges []proto.ExtensionRange

func (er extRanges) String() string ***REMOVED***
	var buf bytes.Buffer
	first := true
	for _, r := range er ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			buf.WriteString(",")
		***REMOVED***
		fmt.Fprintf(&buf, "%d..%d", r.Start, r.End)
	***REMOVED***
	return buf.String()
***REMOVED***

func (er extRanges) IsExtension(tagNumber int32) bool ***REMOVED***
	i := sort.Search(len(er), func(i int) bool ***REMOVED*** return er[i].End >= tagNumber ***REMOVED***)
	return i < len(er) && tagNumber >= er[i].Start
***REMOVED***

func (er extRanges) Len() int ***REMOVED***
	return len(er)
***REMOVED***

func (er extRanges) Less(i, j int) bool ***REMOVED***
	return er[i].Start < er[j].Start
***REMOVED***

func (er extRanges) Swap(i, j int) ***REMOVED***
	er[i], er[j] = er[j], er[i]
***REMOVED***

// FindFieldByName finds the field with the given name. If no such field exists
// then nil is returned. Only regular fields are returned, not extensions.
func (md *MessageDescriptor) FindFieldByName(fieldName string) *FieldDescriptor ***REMOVED***
	fqn := fmt.Sprintf("%s.%s", md.fqn, fieldName)
	if fd, ok := md.file.symbols[fqn].(*FieldDescriptor); ok && !fd.IsExtension() ***REMOVED***
		return fd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindFieldByNumber finds the field with the given tag number. If no such field
// exists then nil is returned. Only regular fields are returned, not extensions.
func (md *MessageDescriptor) FindFieldByNumber(tagNumber int32) *FieldDescriptor ***REMOVED***
	if fd, ok := md.file.fieldIndex[md.fqn][tagNumber]; ok && !fd.IsExtension() ***REMOVED***
		return fd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FieldDescriptor describes a field of a protocol buffer message.
type FieldDescriptor struct ***REMOVED***
	proto          *dpb.FieldDescriptorProto
	parent         Descriptor
	owner          *MessageDescriptor
	file           *FileDescriptor
	oneOf          *OneOfDescriptor
	msgType        *MessageDescriptor
	enumType       *EnumDescriptor
	fqn            string
	sourceInfoPath []int32
	def            memoizedDefault
	isMap          bool
***REMOVED***

func createFieldDescriptor(fd *FileDescriptor, parent Descriptor, enclosing string, fld *dpb.FieldDescriptorProto) (*FieldDescriptor, string) ***REMOVED***
	fldName := merge(enclosing, fld.GetName())
	ret := &FieldDescriptor***REMOVED***proto: fld, parent: parent, file: fd, fqn: fldName***REMOVED***
	if fld.GetExtendee() == "" ***REMOVED***
		ret.owner = parent.(*MessageDescriptor)
	***REMOVED***
	// owner for extensions, field type (be it message or enum), and one-ofs get resolved later
	return ret, fldName
***REMOVED***

func (fd *FieldDescriptor) resolve(path []int32, scopes []scope) error ***REMOVED***
	if fd.proto.OneofIndex != nil && fd.oneOf == nil ***REMOVED***
		return fmt.Errorf("could not link field %s to one-of index %d", fd.fqn, *fd.proto.OneofIndex)
	***REMOVED***
	fd.sourceInfoPath = append([]int32(nil), path...) // defensive copy
	if fd.proto.GetType() == dpb.FieldDescriptorProto_TYPE_ENUM ***REMOVED***
		if desc, err := resolve(fd.file, fd.proto.GetTypeName(), scopes); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			fd.enumType = desc.(*EnumDescriptor)
		***REMOVED***
	***REMOVED***
	if fd.proto.GetType() == dpb.FieldDescriptorProto_TYPE_MESSAGE || fd.proto.GetType() == dpb.FieldDescriptorProto_TYPE_GROUP ***REMOVED***
		if desc, err := resolve(fd.file, fd.proto.GetTypeName(), scopes); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			fd.msgType = desc.(*MessageDescriptor)
		***REMOVED***
	***REMOVED***
	if fd.proto.GetExtendee() != "" ***REMOVED***
		if desc, err := resolve(fd.file, fd.proto.GetExtendee(), scopes); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			fd.owner = desc.(*MessageDescriptor)
		***REMOVED***
	***REMOVED***
	fd.file.registerField(fd)
	fd.isMap = fd.proto.GetLabel() == dpb.FieldDescriptorProto_LABEL_REPEATED &&
		fd.proto.GetType() == dpb.FieldDescriptorProto_TYPE_MESSAGE &&
		fd.GetMessageType().IsMapEntry()
	return nil
***REMOVED***

func (fd *FieldDescriptor) determineDefault() interface***REMOVED******REMOVED*** ***REMOVED***
	if fd.IsMap() ***REMOVED***
		return map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***(nil)
	***REMOVED*** else if fd.IsRepeated() ***REMOVED***
		return []interface***REMOVED******REMOVED***(nil)
	***REMOVED*** else if fd.msgType != nil ***REMOVED***
		return nil
	***REMOVED***

	proto3 := fd.file.isProto3
	if !proto3 ***REMOVED***
		def := fd.AsFieldDescriptorProto().GetDefaultValue()
		if def != "" ***REMOVED***
			ret := parseDefaultValue(fd, def)
			if ret != nil ***REMOVED***
				return ret
			***REMOVED***
			// if we can't parse default value, fall-through to return normal default...
		***REMOVED***
	***REMOVED***

	switch fd.GetType() ***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_FIXED32,
		dpb.FieldDescriptorProto_TYPE_UINT32:
		return uint32(0)
	case dpb.FieldDescriptorProto_TYPE_SFIXED32,
		dpb.FieldDescriptorProto_TYPE_INT32,
		dpb.FieldDescriptorProto_TYPE_SINT32:
		return int32(0)
	case dpb.FieldDescriptorProto_TYPE_FIXED64,
		dpb.FieldDescriptorProto_TYPE_UINT64:
		return uint64(0)
	case dpb.FieldDescriptorProto_TYPE_SFIXED64,
		dpb.FieldDescriptorProto_TYPE_INT64,
		dpb.FieldDescriptorProto_TYPE_SINT64:
		return int64(0)
	case dpb.FieldDescriptorProto_TYPE_FLOAT:
		return float32(0.0)
	case dpb.FieldDescriptorProto_TYPE_DOUBLE:
		return float64(0.0)
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		return false
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		return []byte(nil)
	case dpb.FieldDescriptorProto_TYPE_STRING:
		return ""
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		if proto3 ***REMOVED***
			return int32(0)
		***REMOVED***
		enumVals := fd.GetEnumType().GetValues()
		if len(enumVals) > 0 ***REMOVED***
			return enumVals[0].GetNumber()
		***REMOVED*** else ***REMOVED***
			return int32(0) // WTF?
		***REMOVED***
	default:
		panic(fmt.Sprintf("Unknown field type: %v", fd.GetType()))
	***REMOVED***
***REMOVED***

func parseDefaultValue(fd *FieldDescriptor, val string) interface***REMOVED******REMOVED*** ***REMOVED***
	switch fd.GetType() ***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		vd := fd.GetEnumType().FindValueByName(val)
		if vd != nil ***REMOVED***
			return vd.GetNumber()
		***REMOVED***
		return nil
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		if val == "true" ***REMOVED***
			return true
		***REMOVED*** else if val == "false" ***REMOVED***
			return false
		***REMOVED***
		return nil
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		return []byte(unescape(val))
	case dpb.FieldDescriptorProto_TYPE_STRING:
		return val
	case dpb.FieldDescriptorProto_TYPE_FLOAT:
		if f, err := strconv.ParseFloat(val, 32); err == nil ***REMOVED***
			return float32(f)
		***REMOVED*** else ***REMOVED***
			return float32(0)
		***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_DOUBLE:
		if f, err := strconv.ParseFloat(val, 64); err == nil ***REMOVED***
			return f
		***REMOVED*** else ***REMOVED***
			return float64(0)
		***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_INT32,
		dpb.FieldDescriptorProto_TYPE_SINT32,
		dpb.FieldDescriptorProto_TYPE_SFIXED32:
		if i, err := strconv.ParseInt(val, 10, 32); err == nil ***REMOVED***
			return int32(i)
		***REMOVED*** else ***REMOVED***
			return int32(0)
		***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_UINT32,
		dpb.FieldDescriptorProto_TYPE_FIXED32:
		if i, err := strconv.ParseUint(val, 10, 32); err == nil ***REMOVED***
			return uint32(i)
		***REMOVED*** else ***REMOVED***
			return uint32(0)
		***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_INT64,
		dpb.FieldDescriptorProto_TYPE_SINT64,
		dpb.FieldDescriptorProto_TYPE_SFIXED64:
		if i, err := strconv.ParseInt(val, 10, 64); err == nil ***REMOVED***
			return i
		***REMOVED*** else ***REMOVED***
			return int64(0)
		***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_UINT64,
		dpb.FieldDescriptorProto_TYPE_FIXED64:
		if i, err := strconv.ParseUint(val, 10, 64); err == nil ***REMOVED***
			return i
		***REMOVED*** else ***REMOVED***
			return uint64(0)
		***REMOVED***
	default:
		return nil
	***REMOVED***
***REMOVED***

func unescape(s string) string ***REMOVED***
	// protoc encodes default values for 'bytes' fields using C escaping,
	// so this function reverses that escaping
	out := make([]byte, 0, len(s))
	var buf [4]byte
	for len(s) > 0 ***REMOVED***
		if s[0] != '\\' || len(s) < 2 ***REMOVED***
			// not escape sequence, or too short to be well-formed escape
			out = append(out, s[0])
			s = s[1:]
		***REMOVED*** else if s[1] == 'x' || s[1] == 'X' ***REMOVED***
			n := matchPrefix(s[2:], 2, isHex)
			if n == 0 ***REMOVED***
				// bad escape
				out = append(out, s[:2]...)
				s = s[2:]
			***REMOVED*** else ***REMOVED***
				c, err := strconv.ParseUint(s[2:2+n], 16, 8)
				if err != nil ***REMOVED***
					// shouldn't really happen...
					out = append(out, s[:2+n]...)
				***REMOVED*** else ***REMOVED***
					out = append(out, byte(c))
				***REMOVED***
				s = s[2+n:]
			***REMOVED***
		***REMOVED*** else if s[1] >= '0' && s[1] <= '7' ***REMOVED***
			n := 1 + matchPrefix(s[2:], 2, isOctal)
			c, err := strconv.ParseUint(s[1:1+n], 8, 8)
			if err != nil || c > 0xff ***REMOVED***
				out = append(out, s[:1+n]...)
			***REMOVED*** else ***REMOVED***
				out = append(out, byte(c))
			***REMOVED***
			s = s[1+n:]
		***REMOVED*** else if s[1] == 'u' ***REMOVED***
			if len(s) < 6 ***REMOVED***
				// bad escape
				out = append(out, s...)
				s = s[len(s):]
			***REMOVED*** else ***REMOVED***
				c, err := strconv.ParseUint(s[2:6], 16, 16)
				if err != nil ***REMOVED***
					// bad escape
					out = append(out, s[:6]...)
				***REMOVED*** else ***REMOVED***
					w := utf8.EncodeRune(buf[:], rune(c))
					out = append(out, buf[:w]...)
				***REMOVED***
				s = s[6:]
			***REMOVED***
		***REMOVED*** else if s[1] == 'U' ***REMOVED***
			if len(s) < 10 ***REMOVED***
				// bad escape
				out = append(out, s...)
				s = s[len(s):]
			***REMOVED*** else ***REMOVED***
				c, err := strconv.ParseUint(s[2:10], 16, 32)
				if err != nil || c > 0x10ffff ***REMOVED***
					// bad escape
					out = append(out, s[:10]...)
				***REMOVED*** else ***REMOVED***
					w := utf8.EncodeRune(buf[:], rune(c))
					out = append(out, buf[:w]...)
				***REMOVED***
				s = s[10:]
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			switch s[1] ***REMOVED***
			case 'a':
				out = append(out, '\a')
			case 'b':
				out = append(out, '\b')
			case 'f':
				out = append(out, '\f')
			case 'n':
				out = append(out, '\n')
			case 'r':
				out = append(out, '\r')
			case 't':
				out = append(out, '\t')
			case 'v':
				out = append(out, '\v')
			case '\\':
				out = append(out, '\\')
			case '\'':
				out = append(out, '\'')
			case '"':
				out = append(out, '"')
			case '?':
				out = append(out, '?')
			default:
				// invalid escape, just copy it as-is
				out = append(out, s[:2]...)
			***REMOVED***
			s = s[2:]
		***REMOVED***
	***REMOVED***
	return string(out)
***REMOVED***

func isOctal(b byte) bool ***REMOVED*** return b >= '0' && b <= '7' ***REMOVED***
func isHex(b byte) bool ***REMOVED***
	return (b >= '0' && b <= '9') || (b >= 'a' && b <= 'f') || (b >= 'A' && b <= 'F')
***REMOVED***
func matchPrefix(s string, limit int, fn func(byte) bool) int ***REMOVED***
	l := len(s)
	if l > limit ***REMOVED***
		l = limit
	***REMOVED***
	i := 0
	for ; i < l; i++ ***REMOVED***
		if !fn(s[i]) ***REMOVED***
			return i
		***REMOVED***
	***REMOVED***
	return i
***REMOVED***

// GetName returns the name of the field.
func (fd *FieldDescriptor) GetName() string ***REMOVED***
	return fd.proto.GetName()
***REMOVED***

// GetNumber returns the tag number of this field.
func (fd *FieldDescriptor) GetNumber() int32 ***REMOVED***
	return fd.proto.GetNumber()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the field. Unlike
// GetName, this includes fully qualified name of the enclosing message for
// regular fields.
//
// For extension fields, this includes the package (if there is one) as well as
// any enclosing messages. The package and/or enclosing messages are for where
// the extension is defined, not the message it extends.
//
// If this field is part of a one-of, the fully qualified name does *not*
// include the name of the one-of, only of the enclosing message.
func (fd *FieldDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return fd.fqn
***REMOVED***

// GetParent returns the fields's enclosing descriptor. For normal
// (non-extension) fields, this is the enclosing message. For extensions, this
// is the descriptor in which the extension is defined, not the message that is
// extended. The parent for an extension may be a file descriptor or a message,
// depending on where the extension is defined.
func (fd *FieldDescriptor) GetParent() Descriptor ***REMOVED***
	return fd.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this field is defined.
func (fd *FieldDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return fd.file
***REMOVED***

// GetOptions returns the field's options. Most usages will be more interested
// in GetFieldOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (fd *FieldDescriptor) GetOptions() proto.Message ***REMOVED***
	return fd.proto.GetOptions()
***REMOVED***

// GetFieldOptions returns the field's options.
func (fd *FieldDescriptor) GetFieldOptions() *dpb.FieldOptions ***REMOVED***
	return fd.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the field, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// field was defined and also contains comments associated with the field
// definition.
func (fd *FieldDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return fd.file.sourceInfo.Get(fd.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsFieldDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (fd *FieldDescriptor) AsProto() proto.Message ***REMOVED***
	return fd.proto
***REMOVED***

// AsFieldDescriptorProto returns the underlying descriptor proto.
func (fd *FieldDescriptor) AsFieldDescriptorProto() *dpb.FieldDescriptorProto ***REMOVED***
	return fd.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (fd *FieldDescriptor) String() string ***REMOVED***
	return fd.proto.String()
***REMOVED***

// GetJSONName returns the name of the field as referenced in the message's JSON
// format.
func (fd *FieldDescriptor) GetJSONName() string ***REMOVED***
	if jsonName := fd.proto.JsonName; jsonName != nil ***REMOVED***
		// if json name is present, use its value
		return *jsonName
	***REMOVED***
	// otherwise, compute the proper JSON name from the field name
	return jsonCamelCase(fd.proto.GetName())
***REMOVED***

func jsonCamelCase(s string) string ***REMOVED***
	// This mirrors the implementation in protoc/C++ runtime and in the Java runtime:
	//   https://github.com/protocolbuffers/protobuf/blob/a104dffcb6b1958a424f5fa6f9e6bdc0ab9b6f9e/src/google/protobuf/descriptor.cc#L276
	//   https://github.com/protocolbuffers/protobuf/blob/a1c886834425abb64a966231dd2c9dd84fb289b3/java/core/src/main/java/com/google/protobuf/Descriptors.java#L1286
	var buf bytes.Buffer
	prevWasUnderscore := false
	for _, r := range s ***REMOVED***
		if r == '_' ***REMOVED***
			prevWasUnderscore = true
			continue
		***REMOVED***
		if prevWasUnderscore ***REMOVED***
			r = unicode.ToUpper(r)
			prevWasUnderscore = false
		***REMOVED***
		buf.WriteRune(r)
	***REMOVED***
	return buf.String()
***REMOVED***

// GetFullyQualifiedJSONName returns the JSON format name (same as GetJSONName),
// but includes the fully qualified name of the enclosing message.
//
// If the field is an extension, it will return the package name (if there is
// one) as well as the names of any enclosing messages. The package and/or
// enclosing messages are for where the extension is defined, not the message it
// extends.
func (fd *FieldDescriptor) GetFullyQualifiedJSONName() string ***REMOVED***
	parent := fd.GetParent()
	switch parent := parent.(type) ***REMOVED***
	case *FileDescriptor:
		pkg := parent.GetPackage()
		if pkg == "" ***REMOVED***
			return fd.GetJSONName()
		***REMOVED***
		return fmt.Sprintf("%s.%s", pkg, fd.GetJSONName())
	default:
		return fmt.Sprintf("%s.%s", parent.GetFullyQualifiedName(), fd.GetJSONName())
	***REMOVED***
***REMOVED***

// GetOwner returns the message type that this field belongs to. If this is a normal
// field then this is the same as GetParent. But for extensions, this will be the
// extendee message whereas GetParent refers to where the extension was declared.
func (fd *FieldDescriptor) GetOwner() *MessageDescriptor ***REMOVED***
	return fd.owner
***REMOVED***

// IsExtension returns true if this is an extension field.
func (fd *FieldDescriptor) IsExtension() bool ***REMOVED***
	return fd.proto.GetExtendee() != ""
***REMOVED***

// GetOneOf returns the one-of field set to which this field belongs. If this field
// is not part of a one-of then this method returns nil.
func (fd *FieldDescriptor) GetOneOf() *OneOfDescriptor ***REMOVED***
	return fd.oneOf
***REMOVED***

// GetType returns the type of this field. If the type indicates an enum, the
// enum type can be queried via GetEnumType. If the type indicates a message, the
// message type can be queried via GetMessageType.
func (fd *FieldDescriptor) GetType() dpb.FieldDescriptorProto_Type ***REMOVED***
	return fd.proto.GetType()
***REMOVED***

// GetLabel returns the label for this field. The label can be required (proto2-only),
// optional (default for proto3), or required.
func (fd *FieldDescriptor) GetLabel() dpb.FieldDescriptorProto_Label ***REMOVED***
	return fd.proto.GetLabel()
***REMOVED***

// IsRequired returns true if this field has the "required" label.
func (fd *FieldDescriptor) IsRequired() bool ***REMOVED***
	return fd.proto.GetLabel() == dpb.FieldDescriptorProto_LABEL_REQUIRED
***REMOVED***

// IsRepeated returns true if this field has the "repeated" label.
func (fd *FieldDescriptor) IsRepeated() bool ***REMOVED***
	return fd.proto.GetLabel() == dpb.FieldDescriptorProto_LABEL_REPEATED
***REMOVED***

// IsProto3Optional returns true if this field has an explicit "optional" label
// and is in a "proto3" syntax file. Such fields, if they are normal fields (not
// extensions), will be nested in synthetic oneofs that contain only the single
// field.
func (fd *FieldDescriptor) IsProto3Optional() bool ***REMOVED***
	return internal.GetProto3Optional(fd.proto)
***REMOVED***

// HasPresence returns true if this field can distinguish when a value is
// present or not. Scalar fields in "proto3" syntax files, for example, return
// false since absent values are indistinguishable from zero values.
func (fd *FieldDescriptor) HasPresence() bool ***REMOVED***
	if !fd.file.isProto3 ***REMOVED***
		return true
	***REMOVED***
	return fd.msgType != nil || fd.oneOf != nil
***REMOVED***

// IsMap returns true if this is a map field. If so, it will have the "repeated"
// label its type will be a message that represents a map entry. The map entry
// message will have exactly two fields: tag #1 is the key and tag #2 is the value.
func (fd *FieldDescriptor) IsMap() bool ***REMOVED***
	return fd.isMap
***REMOVED***

// GetMapKeyType returns the type of the key field if this is a map field. If it is
// not a map field, nil is returned.
func (fd *FieldDescriptor) GetMapKeyType() *FieldDescriptor ***REMOVED***
	if fd.isMap ***REMOVED***
		return fd.msgType.FindFieldByNumber(int32(1))
	***REMOVED***
	return nil
***REMOVED***

// GetMapValueType returns the type of the value field if this is a map field. If it
// is not a map field, nil is returned.
func (fd *FieldDescriptor) GetMapValueType() *FieldDescriptor ***REMOVED***
	if fd.isMap ***REMOVED***
		return fd.msgType.FindFieldByNumber(int32(2))
	***REMOVED***
	return nil
***REMOVED***

// GetMessageType returns the type of this field if it is a message type. If
// this field is not a message type, it returns nil.
func (fd *FieldDescriptor) GetMessageType() *MessageDescriptor ***REMOVED***
	return fd.msgType
***REMOVED***

// GetEnumType returns the type of this field if it is an enum type. If this
// field is not an enum type, it returns nil.
func (fd *FieldDescriptor) GetEnumType() *EnumDescriptor ***REMOVED***
	return fd.enumType
***REMOVED***

// GetDefaultValue returns the default value for this field.
//
// If this field represents a message type, this method always returns nil (even though
// for proto2 files, the default value should be a default instance of the message type).
// If the field represents an enum type, this method returns an int32 corresponding to the
// enum value. If this field is a map, it returns a nil map[interface***REMOVED******REMOVED***]interface***REMOVED******REMOVED***. If
// this field is repeated (and not a map), it returns a nil []interface***REMOVED******REMOVED***.
//
// Otherwise, it returns the declared default value for the field or a zero value, if no
// default is declared or if the file is proto3. The type of said return value corresponds
// to the type of the field:
//  +-------------------------+-----------+
//  |       Declared Type     |  Go Type  |
//  +-------------------------+-----------+
//  | int32, sint32, sfixed32 | int32     |
//  | int64, sint64, sfixed64 | int64     |
//  | uint32, fixed32         | uint32    |
//  | uint64, fixed64         | uint64    |
//  | float                   | float32   |
//  | double                  | double32  |
//  | bool                    | bool      |
//  | string                  | string    |
//  | bytes                   | []byte    |
//  +-------------------------+-----------+
func (fd *FieldDescriptor) GetDefaultValue() interface***REMOVED******REMOVED*** ***REMOVED***
	return fd.getDefaultValue()
***REMOVED***

// EnumDescriptor describes an enum declared in a proto file.
type EnumDescriptor struct ***REMOVED***
	proto          *dpb.EnumDescriptorProto
	parent         Descriptor
	file           *FileDescriptor
	values         []*EnumValueDescriptor
	valuesByNum    sortedValues
	fqn            string
	sourceInfoPath []int32
***REMOVED***

func createEnumDescriptor(fd *FileDescriptor, parent Descriptor, enclosing string, ed *dpb.EnumDescriptorProto, symbols map[string]Descriptor) (*EnumDescriptor, string) ***REMOVED***
	enumName := merge(enclosing, ed.GetName())
	ret := &EnumDescriptor***REMOVED***proto: ed, parent: parent, file: fd, fqn: enumName***REMOVED***
	for _, ev := range ed.GetValue() ***REMOVED***
		evd, n := createEnumValueDescriptor(fd, ret, enumName, ev)
		symbols[n] = evd
		ret.values = append(ret.values, evd)
	***REMOVED***
	if len(ret.values) > 0 ***REMOVED***
		ret.valuesByNum = make(sortedValues, len(ret.values))
		copy(ret.valuesByNum, ret.values)
		sort.Stable(ret.valuesByNum)
	***REMOVED***
	return ret, enumName
***REMOVED***

type sortedValues []*EnumValueDescriptor

func (sv sortedValues) Len() int ***REMOVED***
	return len(sv)
***REMOVED***

func (sv sortedValues) Less(i, j int) bool ***REMOVED***
	return sv[i].GetNumber() < sv[j].GetNumber()
***REMOVED***

func (sv sortedValues) Swap(i, j int) ***REMOVED***
	sv[i], sv[j] = sv[j], sv[i]
***REMOVED***

func (ed *EnumDescriptor) resolve(path []int32) ***REMOVED***
	ed.sourceInfoPath = append([]int32(nil), path...) // defensive copy
	path = append(path, internal.Enum_valuesTag)
	for i, evd := range ed.values ***REMOVED***
		evd.resolve(append(path, int32(i)))
	***REMOVED***
***REMOVED***

// GetName returns the simple (unqualified) name of the enum type.
func (ed *EnumDescriptor) GetName() string ***REMOVED***
	return ed.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the enum type.
// This includes the package name (if there is one) as well as the names of any
// enclosing messages.
func (ed *EnumDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return ed.fqn
***REMOVED***

// GetParent returns the enum type's enclosing descriptor. For top-level enums,
// this will be a file descriptor. Otherwise it will be the descriptor for the
// enclosing message.
func (ed *EnumDescriptor) GetParent() Descriptor ***REMOVED***
	return ed.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this enum is defined.
func (ed *EnumDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return ed.file
***REMOVED***

// GetOptions returns the enum type's options. Most usages will be more
// interested in GetEnumOptions, which has a concrete return type. This generic
// version is present to satisfy the Descriptor interface.
func (ed *EnumDescriptor) GetOptions() proto.Message ***REMOVED***
	return ed.proto.GetOptions()
***REMOVED***

// GetEnumOptions returns the enum type's options.
func (ed *EnumDescriptor) GetEnumOptions() *dpb.EnumOptions ***REMOVED***
	return ed.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the enum type, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// enum type was defined and also contains comments associated with the enum
// definition.
func (ed *EnumDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return ed.file.sourceInfo.Get(ed.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsEnumDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (ed *EnumDescriptor) AsProto() proto.Message ***REMOVED***
	return ed.proto
***REMOVED***

// AsEnumDescriptorProto returns the underlying descriptor proto.
func (ed *EnumDescriptor) AsEnumDescriptorProto() *dpb.EnumDescriptorProto ***REMOVED***
	return ed.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (ed *EnumDescriptor) String() string ***REMOVED***
	return ed.proto.String()
***REMOVED***

// GetValues returns all of the allowed values defined for this enum.
func (ed *EnumDescriptor) GetValues() []*EnumValueDescriptor ***REMOVED***
	return ed.values
***REMOVED***

// FindValueByName finds the enum value with the given name. If no such value exists
// then nil is returned.
func (ed *EnumDescriptor) FindValueByName(name string) *EnumValueDescriptor ***REMOVED***
	fqn := fmt.Sprintf("%s.%s", ed.fqn, name)
	if vd, ok := ed.file.symbols[fqn].(*EnumValueDescriptor); ok ***REMOVED***
		return vd
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// FindValueByNumber finds the value with the given numeric value. If no such value
// exists then nil is returned. If aliases are allowed and multiple values have the
// given number, the first declared value is returned.
func (ed *EnumDescriptor) FindValueByNumber(num int32) *EnumValueDescriptor ***REMOVED***
	index := sort.Search(len(ed.valuesByNum), func(i int) bool ***REMOVED*** return ed.valuesByNum[i].GetNumber() >= num ***REMOVED***)
	if index < len(ed.valuesByNum) ***REMOVED***
		vd := ed.valuesByNum[index]
		if vd.GetNumber() == num ***REMOVED***
			return vd
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// EnumValueDescriptor describes an allowed value of an enum declared in a proto file.
type EnumValueDescriptor struct ***REMOVED***
	proto          *dpb.EnumValueDescriptorProto
	parent         *EnumDescriptor
	file           *FileDescriptor
	fqn            string
	sourceInfoPath []int32
***REMOVED***

func createEnumValueDescriptor(fd *FileDescriptor, parent *EnumDescriptor, enclosing string, evd *dpb.EnumValueDescriptorProto) (*EnumValueDescriptor, string) ***REMOVED***
	valName := merge(enclosing, evd.GetName())
	return &EnumValueDescriptor***REMOVED***proto: evd, parent: parent, file: fd, fqn: valName***REMOVED***, valName
***REMOVED***

func (vd *EnumValueDescriptor) resolve(path []int32) ***REMOVED***
	vd.sourceInfoPath = append([]int32(nil), path...) // defensive copy
***REMOVED***

// GetName returns the name of the enum value.
func (vd *EnumValueDescriptor) GetName() string ***REMOVED***
	return vd.proto.GetName()
***REMOVED***

// GetNumber returns the numeric value associated with this enum value.
func (vd *EnumValueDescriptor) GetNumber() int32 ***REMOVED***
	return vd.proto.GetNumber()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the enum value.
// Unlike GetName, this includes fully qualified name of the enclosing enum.
func (vd *EnumValueDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return vd.fqn
***REMOVED***

// GetParent returns the descriptor for the enum in which this enum value is
// defined. Most usages will prefer to use GetEnum, which has a concrete return
// type. This more generic method is present to satisfy the Descriptor interface.
func (vd *EnumValueDescriptor) GetParent() Descriptor ***REMOVED***
	return vd.parent
***REMOVED***

// GetEnum returns the enum in which this enum value is defined.
func (vd *EnumValueDescriptor) GetEnum() *EnumDescriptor ***REMOVED***
	return vd.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this enum value is
// defined.
func (vd *EnumValueDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return vd.file
***REMOVED***

// GetOptions returns the enum value's options. Most usages will be more
// interested in GetEnumValueOptions, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (vd *EnumValueDescriptor) GetOptions() proto.Message ***REMOVED***
	return vd.proto.GetOptions()
***REMOVED***

// GetEnumValueOptions returns the enum value's options.
func (vd *EnumValueDescriptor) GetEnumValueOptions() *dpb.EnumValueOptions ***REMOVED***
	return vd.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the enum value, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// enum value was defined and also contains comments associated with the enum
// value definition.
func (vd *EnumValueDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return vd.file.sourceInfo.Get(vd.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsEnumValueDescriptorProto, which has a concrete return type.
// This generic version is present to satisfy the Descriptor interface.
func (vd *EnumValueDescriptor) AsProto() proto.Message ***REMOVED***
	return vd.proto
***REMOVED***

// AsEnumValueDescriptorProto returns the underlying descriptor proto.
func (vd *EnumValueDescriptor) AsEnumValueDescriptorProto() *dpb.EnumValueDescriptorProto ***REMOVED***
	return vd.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (vd *EnumValueDescriptor) String() string ***REMOVED***
	return vd.proto.String()
***REMOVED***

// ServiceDescriptor describes an RPC service declared in a proto file.
type ServiceDescriptor struct ***REMOVED***
	proto          *dpb.ServiceDescriptorProto
	file           *FileDescriptor
	methods        []*MethodDescriptor
	fqn            string
	sourceInfoPath []int32
***REMOVED***

func createServiceDescriptor(fd *FileDescriptor, enclosing string, sd *dpb.ServiceDescriptorProto, symbols map[string]Descriptor) (*ServiceDescriptor, string) ***REMOVED***
	serviceName := merge(enclosing, sd.GetName())
	ret := &ServiceDescriptor***REMOVED***proto: sd, file: fd, fqn: serviceName***REMOVED***
	for _, m := range sd.GetMethod() ***REMOVED***
		md, n := createMethodDescriptor(fd, ret, serviceName, m)
		symbols[n] = md
		ret.methods = append(ret.methods, md)
	***REMOVED***
	return ret, serviceName
***REMOVED***

func (sd *ServiceDescriptor) resolve(path []int32, scopes []scope) error ***REMOVED***
	sd.sourceInfoPath = append([]int32(nil), path...) // defensive copy
	path = append(path, internal.Service_methodsTag)
	for i, md := range sd.methods ***REMOVED***
		if err := md.resolve(append(path, int32(i)), scopes); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// GetName returns the simple (unqualified) name of the service.
func (sd *ServiceDescriptor) GetName() string ***REMOVED***
	return sd.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the service. This
// includes the package name (if there is one).
func (sd *ServiceDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return sd.fqn
***REMOVED***

// GetParent returns the descriptor for the file in which this service is
// defined. Most usages will prefer to use GetFile, which has a concrete return
// type. This more generic method is present to satisfy the Descriptor interface.
func (sd *ServiceDescriptor) GetParent() Descriptor ***REMOVED***
	return sd.file
***REMOVED***

// GetFile returns the descriptor for the file in which this service is defined.
func (sd *ServiceDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return sd.file
***REMOVED***

// GetOptions returns the service's options. Most usages will be more interested
// in GetServiceOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (sd *ServiceDescriptor) GetOptions() proto.Message ***REMOVED***
	return sd.proto.GetOptions()
***REMOVED***

// GetServiceOptions returns the service's options.
func (sd *ServiceDescriptor) GetServiceOptions() *dpb.ServiceOptions ***REMOVED***
	return sd.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the service, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// service was defined and also contains comments associated with the service
// definition.
func (sd *ServiceDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return sd.file.sourceInfo.Get(sd.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsServiceDescriptorProto, which has a concrete return type.
// This generic version is present to satisfy the Descriptor interface.
func (sd *ServiceDescriptor) AsProto() proto.Message ***REMOVED***
	return sd.proto
***REMOVED***

// AsServiceDescriptorProto returns the underlying descriptor proto.
func (sd *ServiceDescriptor) AsServiceDescriptorProto() *dpb.ServiceDescriptorProto ***REMOVED***
	return sd.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (sd *ServiceDescriptor) String() string ***REMOVED***
	return sd.proto.String()
***REMOVED***

// GetMethods returns all of the RPC methods for this service.
func (sd *ServiceDescriptor) GetMethods() []*MethodDescriptor ***REMOVED***
	return sd.methods
***REMOVED***

// FindMethodByName finds the method with the given name. If no such method exists
// then nil is returned.
func (sd *ServiceDescriptor) FindMethodByName(name string) *MethodDescriptor ***REMOVED***
	fqn := fmt.Sprintf("%s.%s", sd.fqn, name)
	if md, ok := sd.file.symbols[fqn].(*MethodDescriptor); ok ***REMOVED***
		return md
	***REMOVED*** else ***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// MethodDescriptor describes an RPC method declared in a proto file.
type MethodDescriptor struct ***REMOVED***
	proto          *dpb.MethodDescriptorProto
	parent         *ServiceDescriptor
	file           *FileDescriptor
	inType         *MessageDescriptor
	outType        *MessageDescriptor
	fqn            string
	sourceInfoPath []int32
***REMOVED***

func createMethodDescriptor(fd *FileDescriptor, parent *ServiceDescriptor, enclosing string, md *dpb.MethodDescriptorProto) (*MethodDescriptor, string) ***REMOVED***
	// request and response types get resolved later
	methodName := merge(enclosing, md.GetName())
	return &MethodDescriptor***REMOVED***proto: md, parent: parent, file: fd, fqn: methodName***REMOVED***, methodName
***REMOVED***

func (md *MethodDescriptor) resolve(path []int32, scopes []scope) error ***REMOVED***
	md.sourceInfoPath = append([]int32(nil), path...) // defensive copy
	if desc, err := resolve(md.file, md.proto.GetInputType(), scopes); err != nil ***REMOVED***
		return err
	***REMOVED*** else ***REMOVED***
		md.inType = desc.(*MessageDescriptor)
	***REMOVED***
	if desc, err := resolve(md.file, md.proto.GetOutputType(), scopes); err != nil ***REMOVED***
		return err
	***REMOVED*** else ***REMOVED***
		md.outType = desc.(*MessageDescriptor)
	***REMOVED***
	return nil
***REMOVED***

// GetName returns the name of the method.
func (md *MethodDescriptor) GetName() string ***REMOVED***
	return md.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the method. Unlike
// GetName, this includes fully qualified name of the enclosing service.
func (md *MethodDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return md.fqn
***REMOVED***

// GetParent returns the descriptor for the service in which this method is
// defined. Most usages will prefer to use GetService, which has a concrete
// return type. This more generic method is present to satisfy the Descriptor
// interface.
func (md *MethodDescriptor) GetParent() Descriptor ***REMOVED***
	return md.parent
***REMOVED***

// GetService returns the RPC service in which this method is declared.
func (md *MethodDescriptor) GetService() *ServiceDescriptor ***REMOVED***
	return md.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this method is defined.
func (md *MethodDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return md.file
***REMOVED***

// GetOptions returns the method's options. Most usages will be more interested
// in GetMethodOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (md *MethodDescriptor) GetOptions() proto.Message ***REMOVED***
	return md.proto.GetOptions()
***REMOVED***

// GetMethodOptions returns the method's options.
func (md *MethodDescriptor) GetMethodOptions() *dpb.MethodOptions ***REMOVED***
	return md.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the method, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// method was defined and also contains comments associated with the method
// definition.
func (md *MethodDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return md.file.sourceInfo.Get(md.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsMethodDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (md *MethodDescriptor) AsProto() proto.Message ***REMOVED***
	return md.proto
***REMOVED***

// AsMethodDescriptorProto returns the underlying descriptor proto.
func (md *MethodDescriptor) AsMethodDescriptorProto() *dpb.MethodDescriptorProto ***REMOVED***
	return md.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (md *MethodDescriptor) String() string ***REMOVED***
	return md.proto.String()
***REMOVED***

// IsServerStreaming returns true if this is a server-streaming method.
func (md *MethodDescriptor) IsServerStreaming() bool ***REMOVED***
	return md.proto.GetServerStreaming()
***REMOVED***

// IsClientStreaming returns true if this is a client-streaming method.
func (md *MethodDescriptor) IsClientStreaming() bool ***REMOVED***
	return md.proto.GetClientStreaming()
***REMOVED***

// GetInputType returns the input type, or request type, of the RPC method.
func (md *MethodDescriptor) GetInputType() *MessageDescriptor ***REMOVED***
	return md.inType
***REMOVED***

// GetOutputType returns the output type, or response type, of the RPC method.
func (md *MethodDescriptor) GetOutputType() *MessageDescriptor ***REMOVED***
	return md.outType
***REMOVED***

// OneOfDescriptor describes a one-of field set declared in a protocol buffer message.
type OneOfDescriptor struct ***REMOVED***
	proto          *dpb.OneofDescriptorProto
	parent         *MessageDescriptor
	file           *FileDescriptor
	choices        []*FieldDescriptor
	fqn            string
	sourceInfoPath []int32
***REMOVED***

func createOneOfDescriptor(fd *FileDescriptor, parent *MessageDescriptor, index int, enclosing string, od *dpb.OneofDescriptorProto) (*OneOfDescriptor, string) ***REMOVED***
	oneOfName := merge(enclosing, od.GetName())
	ret := &OneOfDescriptor***REMOVED***proto: od, parent: parent, file: fd, fqn: oneOfName***REMOVED***
	for _, f := range parent.fields ***REMOVED***
		oi := f.proto.OneofIndex
		if oi != nil && *oi == int32(index) ***REMOVED***
			f.oneOf = ret
			ret.choices = append(ret.choices, f)
		***REMOVED***
	***REMOVED***
	return ret, oneOfName
***REMOVED***

func (od *OneOfDescriptor) resolve(path []int32) ***REMOVED***
	od.sourceInfoPath = append([]int32(nil), path...) // defensive copy
***REMOVED***

// GetName returns the name of the one-of.
func (od *OneOfDescriptor) GetName() string ***REMOVED***
	return od.proto.GetName()
***REMOVED***

// GetFullyQualifiedName returns the fully qualified name of the one-of. Unlike
// GetName, this includes fully qualified name of the enclosing message.
func (od *OneOfDescriptor) GetFullyQualifiedName() string ***REMOVED***
	return od.fqn
***REMOVED***

// GetParent returns the descriptor for the message in which this one-of is
// defined. Most usages will prefer to use GetOwner, which has a concrete
// return type. This more generic method is present to satisfy the Descriptor
// interface.
func (od *OneOfDescriptor) GetParent() Descriptor ***REMOVED***
	return od.parent
***REMOVED***

// GetOwner returns the message to which this one-of field set belongs.
func (od *OneOfDescriptor) GetOwner() *MessageDescriptor ***REMOVED***
	return od.parent
***REMOVED***

// GetFile returns the descriptor for the file in which this one-fof is defined.
func (od *OneOfDescriptor) GetFile() *FileDescriptor ***REMOVED***
	return od.file
***REMOVED***

// GetOptions returns the one-of's options. Most usages will be more interested
// in GetOneOfOptions, which has a concrete return type. This generic version
// is present to satisfy the Descriptor interface.
func (od *OneOfDescriptor) GetOptions() proto.Message ***REMOVED***
	return od.proto.GetOptions()
***REMOVED***

// GetOneOfOptions returns the one-of's options.
func (od *OneOfDescriptor) GetOneOfOptions() *dpb.OneofOptions ***REMOVED***
	return od.proto.GetOptions()
***REMOVED***

// GetSourceInfo returns source info for the one-of, if present in the
// descriptor. Not all descriptors will contain source info. If non-nil, the
// returned info contains information about the location in the file where the
// one-of was defined and also contains comments associated with the one-of
// definition.
func (od *OneOfDescriptor) GetSourceInfo() *dpb.SourceCodeInfo_Location ***REMOVED***
	return od.file.sourceInfo.Get(od.sourceInfoPath)
***REMOVED***

// AsProto returns the underlying descriptor proto. Most usages will be more
// interested in AsOneofDescriptorProto, which has a concrete return type. This
// generic version is present to satisfy the Descriptor interface.
func (od *OneOfDescriptor) AsProto() proto.Message ***REMOVED***
	return od.proto
***REMOVED***

// AsOneofDescriptorProto returns the underlying descriptor proto.
func (od *OneOfDescriptor) AsOneofDescriptorProto() *dpb.OneofDescriptorProto ***REMOVED***
	return od.proto
***REMOVED***

// String returns the underlying descriptor proto, in compact text format.
func (od *OneOfDescriptor) String() string ***REMOVED***
	return od.proto.String()
***REMOVED***

// GetChoices returns the fields that are part of the one-of field set. At most one of
// these fields may be set for a given message.
func (od *OneOfDescriptor) GetChoices() []*FieldDescriptor ***REMOVED***
	return od.choices
***REMOVED***

func (od *OneOfDescriptor) IsSynthetic() bool ***REMOVED***
	return len(od.choices) == 1 && od.choices[0].IsProto3Optional()
***REMOVED***

// scope represents a lexical scope in a proto file in which messages and enums
// can be declared.
type scope func(string) Descriptor

func fileScope(fd *FileDescriptor) scope ***REMOVED***
	// we search symbols in this file, but also symbols in other files that have
	// the same package as this file or a "parent" package (in protobuf,
	// packages are a hierarchy like C++ namespaces)
	prefixes := internal.CreatePrefixList(fd.proto.GetPackage())
	return func(name string) Descriptor ***REMOVED***
		for _, prefix := range prefixes ***REMOVED***
			n := merge(prefix, name)
			d := findSymbol(fd, n, false)
			if d != nil ***REMOVED***
				return d
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func messageScope(md *MessageDescriptor) scope ***REMOVED***
	return func(name string) Descriptor ***REMOVED***
		n := merge(md.fqn, name)
		if d, ok := md.file.symbols[n]; ok ***REMOVED***
			return d
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

func resolve(fd *FileDescriptor, name string, scopes []scope) (Descriptor, error) ***REMOVED***
	if strings.HasPrefix(name, ".") ***REMOVED***
		// already fully-qualified
		d := findSymbol(fd, name[1:], false)
		if d != nil ***REMOVED***
			return d, nil
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// unqualified, so we look in the enclosing (last) scope first and move
		// towards outermost (first) scope, trying to resolve the symbol
		for i := len(scopes) - 1; i >= 0; i-- ***REMOVED***
			d := scopes[i](name)
			if d != nil ***REMOVED***
				return d, nil
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil, fmt.Errorf("file %q included an unresolvable reference to %q", fd.proto.GetName(), name)
***REMOVED***

func findSymbol(fd *FileDescriptor, name string, public bool) Descriptor ***REMOVED***
	d := fd.symbols[name]
	if d != nil ***REMOVED***
		return d
	***REMOVED***

	// When public = false, we are searching only directly imported symbols. But we
	// also need to search transitive public imports due to semantics of public imports.
	var deps []*FileDescriptor
	if public ***REMOVED***
		deps = fd.publicDeps
	***REMOVED*** else ***REMOVED***
		deps = fd.deps
	***REMOVED***
	for _, dep := range deps ***REMOVED***
		d = findSymbol(dep, name, true)
		if d != nil ***REMOVED***
			return d
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func merge(a, b string) string ***REMOVED***
	if a == "" ***REMOVED***
		return b
	***REMOVED*** else ***REMOVED***
		return a + "." + b
	***REMOVED***
***REMOVED***
