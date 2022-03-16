package desc

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc/sourceinfo"
	"github.com/jhump/protoreflect/internal"
)

var (
	cacheMu       sync.RWMutex
	filesCache    = map[string]*FileDescriptor***REMOVED******REMOVED***
	messagesCache = map[string]*MessageDescriptor***REMOVED******REMOVED***
	enumCache     = map[reflect.Type]*EnumDescriptor***REMOVED******REMOVED***
)

// LoadFileDescriptor creates a file descriptor using the bytes returned by
// proto.FileDescriptor. Descriptors are cached so that they do not need to be
// re-processed if the same file is fetched again later.
func LoadFileDescriptor(file string) (*FileDescriptor, error) ***REMOVED***
	return loadFileDescriptor(file, nil)
***REMOVED***

func loadFileDescriptor(file string, r *ImportResolver) (*FileDescriptor, error) ***REMOVED***
	f := getFileFromCache(file)
	if f != nil ***REMOVED***
		return f, nil
	***REMOVED***
	cacheMu.Lock()
	defer cacheMu.Unlock()
	return loadFileDescriptorLocked(file, r)
***REMOVED***

func loadFileDescriptorLocked(file string, r *ImportResolver) (*FileDescriptor, error) ***REMOVED***
	f := filesCache[file]
	if f != nil ***REMOVED***
		return f, nil
	***REMOVED***
	fd, err := internal.LoadFileDescriptor(file)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f, err = toFileDescriptorLocked(fd, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	putCacheLocked(file, f)
	return f, nil
***REMOVED***

func toFileDescriptorLocked(fd *dpb.FileDescriptorProto, r *ImportResolver) (*FileDescriptor, error) ***REMOVED***
	fd.SourceCodeInfo = sourceinfo.SourceInfoForFile(fd.GetName())
	deps := make([]*FileDescriptor, len(fd.GetDependency()))
	for i, dep := range fd.GetDependency() ***REMOVED***
		resolvedDep := r.ResolveImport(fd.GetName(), dep)
		var err error
		deps[i], err = loadFileDescriptorLocked(resolvedDep, r)
		if _, ok := err.(internal.ErrNoSuchFile); ok && resolvedDep != dep ***REMOVED***
			// try original path
			deps[i], err = loadFileDescriptorLocked(dep, r)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
	return CreateFileDescriptor(fd, deps...)
***REMOVED***

func getFileFromCache(file string) *FileDescriptor ***REMOVED***
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return filesCache[file]
***REMOVED***

func putCacheLocked(filename string, fd *FileDescriptor) ***REMOVED***
	filesCache[filename] = fd
	putMessageCacheLocked(fd.messages)
***REMOVED***

func putMessageCacheLocked(mds []*MessageDescriptor) ***REMOVED***
	for _, md := range mds ***REMOVED***
		messagesCache[md.fqn] = md
		putMessageCacheLocked(md.nested)
	***REMOVED***
***REMOVED***

// interface implemented by generated messages, which all have a Descriptor() method in
// addition to the methods of proto.Message
type protoMessage interface ***REMOVED***
	proto.Message
	Descriptor() ([]byte, []int)
***REMOVED***

// LoadMessageDescriptor loads descriptor using the encoded descriptor proto returned by
// Message.Descriptor() for the given message type. If the given type is not recognized,
// then a nil descriptor is returned.
func LoadMessageDescriptor(message string) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptor(message, nil)
***REMOVED***

func loadMessageDescriptor(message string, r *ImportResolver) (*MessageDescriptor, error) ***REMOVED***
	m := getMessageFromCache(message)
	if m != nil ***REMOVED***
		return m, nil
	***REMOVED***

	pt := proto.MessageType(message)
	if pt == nil ***REMOVED***
		return nil, nil
	***REMOVED***
	msg, err := messageFromType(pt)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cacheMu.Lock()
	defer cacheMu.Unlock()
	return loadMessageDescriptorForTypeLocked(message, msg, r)
***REMOVED***

// LoadMessageDescriptorForType loads descriptor using the encoded descriptor proto returned
// by message.Descriptor() for the given message type. If the given type is not recognized,
// then a nil descriptor is returned.
func LoadMessageDescriptorForType(messageType reflect.Type) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptorForType(messageType, nil)
***REMOVED***

func loadMessageDescriptorForType(messageType reflect.Type, r *ImportResolver) (*MessageDescriptor, error) ***REMOVED***
	m, err := messageFromType(messageType)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return loadMessageDescriptorForMessage(m, r)
***REMOVED***

// LoadMessageDescriptorForMessage loads descriptor using the encoded descriptor proto
// returned by message.Descriptor(). If the given type is not recognized, then a nil
// descriptor is returned.
func LoadMessageDescriptorForMessage(message proto.Message) (*MessageDescriptor, error) ***REMOVED***
	return loadMessageDescriptorForMessage(message, nil)
***REMOVED***

func loadMessageDescriptorForMessage(message proto.Message, r *ImportResolver) (*MessageDescriptor, error) ***REMOVED***
	// efficiently handle dynamic messages
	type descriptorable interface ***REMOVED***
		GetMessageDescriptor() *MessageDescriptor
	***REMOVED***
	if d, ok := message.(descriptorable); ok ***REMOVED***
		return d.GetMessageDescriptor(), nil
	***REMOVED***

	name := proto.MessageName(message)
	if name == "" ***REMOVED***
		return nil, nil
	***REMOVED***
	m := getMessageFromCache(name)
	if m != nil ***REMOVED***
		return m, nil
	***REMOVED***

	cacheMu.Lock()
	defer cacheMu.Unlock()
	return loadMessageDescriptorForTypeLocked(name, message.(protoMessage), nil)
***REMOVED***

func messageFromType(mt reflect.Type) (protoMessage, error) ***REMOVED***
	if mt.Kind() != reflect.Ptr ***REMOVED***
		mt = reflect.PtrTo(mt)
	***REMOVED***
	m, ok := reflect.Zero(mt).Interface().(protoMessage)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to create message from type: %v", mt)
	***REMOVED***
	return m, nil
***REMOVED***

func loadMessageDescriptorForTypeLocked(name string, message protoMessage, r *ImportResolver) (*MessageDescriptor, error) ***REMOVED***
	m := messagesCache[name]
	if m != nil ***REMOVED***
		return m, nil
	***REMOVED***

	fdb, _ := message.Descriptor()
	fd, err := internal.DecodeFileDescriptor(name, fdb)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	f, err := toFileDescriptorLocked(fd, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	putCacheLocked(fd.GetName(), f)
	return f.FindSymbol(name).(*MessageDescriptor), nil
***REMOVED***

func getMessageFromCache(message string) *MessageDescriptor ***REMOVED***
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return messagesCache[message]
***REMOVED***

// interface implemented by all generated enums
type protoEnum interface ***REMOVED***
	EnumDescriptor() ([]byte, []int)
***REMOVED***

// NB: There is no LoadEnumDescriptor that takes a fully-qualified enum name because
// it is not useful since protoc-gen-go does not expose the name anywhere in generated
// code or register it in a way that is it accessible for reflection code. This also
// means we have to cache enum descriptors differently -- we can only cache them as
// they are requested, as opposed to caching all enum types whenever a file descriptor
// is cached. This is because we need to know the generated type of the enums, and we
// don't know that at the time of caching file descriptors.

// LoadEnumDescriptorForType loads descriptor using the encoded descriptor proto returned
// by enum.EnumDescriptor() for the given enum type.
func LoadEnumDescriptorForType(enumType reflect.Type) (*EnumDescriptor, error) ***REMOVED***
	return loadEnumDescriptorForType(enumType, nil)
***REMOVED***

func loadEnumDescriptorForType(enumType reflect.Type, r *ImportResolver) (*EnumDescriptor, error) ***REMOVED***
	// we cache descriptors using non-pointer type
	if enumType.Kind() == reflect.Ptr ***REMOVED***
		enumType = enumType.Elem()
	***REMOVED***
	e := getEnumFromCache(enumType)
	if e != nil ***REMOVED***
		return e, nil
	***REMOVED***
	enum, err := enumFromType(enumType)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	cacheMu.Lock()
	defer cacheMu.Unlock()
	return loadEnumDescriptorForTypeLocked(enumType, enum, r)
***REMOVED***

// LoadEnumDescriptorForEnum loads descriptor using the encoded descriptor proto
// returned by enum.EnumDescriptor().
func LoadEnumDescriptorForEnum(enum protoEnum) (*EnumDescriptor, error) ***REMOVED***
	return loadEnumDescriptorForEnum(enum, nil)
***REMOVED***

func loadEnumDescriptorForEnum(enum protoEnum, r *ImportResolver) (*EnumDescriptor, error) ***REMOVED***
	et := reflect.TypeOf(enum)
	// we cache descriptors using non-pointer type
	if et.Kind() == reflect.Ptr ***REMOVED***
		et = et.Elem()
		enum = reflect.Zero(et).Interface().(protoEnum)
	***REMOVED***
	e := getEnumFromCache(et)
	if e != nil ***REMOVED***
		return e, nil
	***REMOVED***

	cacheMu.Lock()
	defer cacheMu.Unlock()
	return loadEnumDescriptorForTypeLocked(et, enum, r)
***REMOVED***

func enumFromType(et reflect.Type) (protoEnum, error) ***REMOVED***
	if et.Kind() != reflect.Int32 ***REMOVED***
		et = reflect.PtrTo(et)
	***REMOVED***
	e, ok := reflect.Zero(et).Interface().(protoEnum)
	if !ok ***REMOVED***
		return nil, fmt.Errorf("failed to create enum from type: %v", et)
	***REMOVED***
	return e, nil
***REMOVED***

func loadEnumDescriptorForTypeLocked(et reflect.Type, enum protoEnum, r *ImportResolver) (*EnumDescriptor, error) ***REMOVED***
	e := enumCache[et]
	if e != nil ***REMOVED***
		return e, nil
	***REMOVED***

	fdb, path := enum.EnumDescriptor()
	name := fmt.Sprintf("%v", et)
	fd, err := internal.DecodeFileDescriptor(name, fdb)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// see if we already have cached "rich" descriptor
	f, ok := filesCache[fd.GetName()]
	if !ok ***REMOVED***
		f, err = toFileDescriptorLocked(fd, r)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		putCacheLocked(fd.GetName(), f)
	***REMOVED***

	ed := findEnum(f, path)
	enumCache[et] = ed
	return ed, nil
***REMOVED***

func getEnumFromCache(et reflect.Type) *EnumDescriptor ***REMOVED***
	cacheMu.RLock()
	defer cacheMu.RUnlock()
	return enumCache[et]
***REMOVED***

func findEnum(fd *FileDescriptor, path []int) *EnumDescriptor ***REMOVED***
	if len(path) == 1 ***REMOVED***
		return fd.GetEnumTypes()[path[0]]
	***REMOVED***
	md := fd.GetMessageTypes()[path[0]]
	for _, i := range path[1 : len(path)-1] ***REMOVED***
		md = md.GetNestedMessageTypes()[i]
	***REMOVED***
	return md.GetNestedEnumTypes()[path[len(path)-1]]
***REMOVED***

// LoadFieldDescriptorForExtension loads the field descriptor that corresponds to the given
// extension description.
func LoadFieldDescriptorForExtension(ext *proto.ExtensionDesc) (*FieldDescriptor, error) ***REMOVED***
	return loadFieldDescriptorForExtension(ext, nil)
***REMOVED***

func loadFieldDescriptorForExtension(ext *proto.ExtensionDesc, r *ImportResolver) (*FieldDescriptor, error) ***REMOVED***
	file, err := loadFileDescriptor(ext.Filename, r)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	field, ok := file.FindSymbol(ext.Name).(*FieldDescriptor)
	// make sure descriptor agrees with attributes of the ExtensionDesc
	if !ok || !field.IsExtension() || field.GetOwner().GetFullyQualifiedName() != proto.MessageName(ext.ExtendedType) ||
		field.GetNumber() != ext.Field ***REMOVED***
		return nil, fmt.Errorf("file descriptor contained unexpected object with name %s", ext.Name)
	***REMOVED***
	return field, nil
***REMOVED***
