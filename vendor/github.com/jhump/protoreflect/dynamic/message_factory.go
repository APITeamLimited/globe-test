package dynamic

import (
	"reflect"
	"sync"

	"github.com/golang/protobuf/proto"

	"github.com/jhump/protoreflect/desc"
)

// MessageFactory can be used to create new empty message objects. A default instance
// (without extension registry or known-type registry specified) will always return
// dynamic messages (e.g. type will be *dynamic.Message) except for "well-known" types.
// The well-known types include primitive wrapper types and a handful of other special
// types defined in standard protobuf definitions, like Any, Duration, and Timestamp.
type MessageFactory struct ***REMOVED***
	er  *ExtensionRegistry
	ktr *KnownTypeRegistry
***REMOVED***

// NewMessageFactoryWithExtensionRegistry creates a new message factory where any
// dynamic messages produced will use the given extension registry to recognize and
// parse extension fields.
func NewMessageFactoryWithExtensionRegistry(er *ExtensionRegistry) *MessageFactory ***REMOVED***
	return NewMessageFactoryWithRegistries(er, nil)
***REMOVED***

// NewMessageFactoryWithKnownTypeRegistry creates a new message factory where the
// known types, per the given registry, will be returned as normal protobuf messages
// (e.g. generated structs, instead of dynamic messages).
func NewMessageFactoryWithKnownTypeRegistry(ktr *KnownTypeRegistry) *MessageFactory ***REMOVED***
	return NewMessageFactoryWithRegistries(nil, ktr)
***REMOVED***

// NewMessageFactoryWithDefaults creates a new message factory where all "default" types
// (those for which protoc-generated code is statically linked into the Go program) are
// known types. If any dynamic messages are produced, they will recognize and parse all
// "default" extension fields. This is the equivalent of:
//   NewMessageFactoryWithRegistries(
//       NewExtensionRegistryWithDefaults(),
//       NewKnownTypeRegistryWithDefaults())
func NewMessageFactoryWithDefaults() *MessageFactory ***REMOVED***
	return NewMessageFactoryWithRegistries(NewExtensionRegistryWithDefaults(), NewKnownTypeRegistryWithDefaults())
***REMOVED***

// NewMessageFactoryWithRegistries creates a new message factory with the given extension
// and known type registries.
func NewMessageFactoryWithRegistries(er *ExtensionRegistry, ktr *KnownTypeRegistry) *MessageFactory ***REMOVED***
	return &MessageFactory***REMOVED***
		er:  er,
		ktr: ktr,
	***REMOVED***
***REMOVED***

// NewMessage creates a new empty message that corresponds to the given descriptor.
// If the given descriptor describes a "known type" then that type is instantiated.
// Otherwise, an empty dynamic message is returned.
func (f *MessageFactory) NewMessage(md *desc.MessageDescriptor) proto.Message ***REMOVED***
	var ktr *KnownTypeRegistry
	if f != nil ***REMOVED***
		ktr = f.ktr
	***REMOVED***
	if m := ktr.CreateIfKnown(md.GetFullyQualifiedName()); m != nil ***REMOVED***
		return m
	***REMOVED***
	return NewMessageWithMessageFactory(md, f)
***REMOVED***

// NewDynamicMessage creates a new empty dynamic message that corresponds to the given
// descriptor. This is like f.NewMessage(md) except the known type registry is not
// consulted so the return value is always a dynamic message.
//
// This is also like dynamic.NewMessage(md) except that the returned message will use
// this factory when creating other messages, like during de-serialization of fields
// that are themselves message types.
func (f *MessageFactory) NewDynamicMessage(md *desc.MessageDescriptor) *Message ***REMOVED***
	return NewMessageWithMessageFactory(md, f)
***REMOVED***

// GetKnownTypeRegistry returns the known type registry that this factory uses to
// instantiate known (e.g. generated) message types.
func (f *MessageFactory) GetKnownTypeRegistry() *KnownTypeRegistry ***REMOVED***
	if f == nil ***REMOVED***
		return nil
	***REMOVED***
	return f.ktr
***REMOVED***

// GetExtensionRegistry returns the extension registry that this factory uses to
// create dynamic messages. The registry is used by dynamic messages to recognize
// and parse extension fields during de-serialization.
func (f *MessageFactory) GetExtensionRegistry() *ExtensionRegistry ***REMOVED***
	if f == nil ***REMOVED***
		return nil
	***REMOVED***
	return f.er
***REMOVED***

type wkt interface ***REMOVED***
	XXX_WellKnownType() string
***REMOVED***

var typeOfWkt = reflect.TypeOf((*wkt)(nil)).Elem()

// KnownTypeRegistry is a registry of known message types, as identified by their
// fully-qualified name. A known message type is one for which a protoc-generated
// struct exists, so a dynamic message is not necessary to represent it. A
// MessageFactory uses a KnownTypeRegistry to decide whether to create a generated
// struct or a dynamic message. The zero-value registry (including the behavior of
// a nil pointer) only knows about the "well-known types" in protobuf. These
// include only the wrapper types and a handful of other special types like Any,
// Duration, and Timestamp.
type KnownTypeRegistry struct ***REMOVED***
	excludeWkt     bool
	includeDefault bool
	mu             sync.RWMutex
	types          map[string]reflect.Type
***REMOVED***

// NewKnownTypeRegistryWithDefaults creates a new registry that knows about all
// "default" types (those for which protoc-generated code is statically linked
// into the Go program).
func NewKnownTypeRegistryWithDefaults() *KnownTypeRegistry ***REMOVED***
	return &KnownTypeRegistry***REMOVED***includeDefault: true***REMOVED***
***REMOVED***

// NewKnownTypeRegistryWithoutWellKnownTypes creates a new registry that does *not*
// include the "well-known types" in protobuf. So even well-known types would be
// represented by a dynamic message.
func NewKnownTypeRegistryWithoutWellKnownTypes() *KnownTypeRegistry ***REMOVED***
	return &KnownTypeRegistry***REMOVED***excludeWkt: true***REMOVED***
***REMOVED***

// AddKnownType adds the types of the given messages as known types.
func (r *KnownTypeRegistry) AddKnownType(kts ...proto.Message) ***REMOVED***
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.types == nil ***REMOVED***
		r.types = map[string]reflect.Type***REMOVED******REMOVED***
	***REMOVED***
	for _, kt := range kts ***REMOVED***
		r.types[proto.MessageName(kt)] = reflect.TypeOf(kt)
	***REMOVED***
***REMOVED***

// CreateIfKnown will construct an instance of the given message if it is a known type.
// If the given name is unknown, nil is returned.
func (r *KnownTypeRegistry) CreateIfKnown(messageName string) proto.Message ***REMOVED***
	msgType := r.GetKnownType(messageName)
	if msgType == nil ***REMOVED***
		return nil
	***REMOVED***

	if msgType.Kind() == reflect.Ptr ***REMOVED***
		return reflect.New(msgType.Elem()).Interface().(proto.Message)
	***REMOVED*** else ***REMOVED***
		return reflect.New(msgType).Elem().Interface().(proto.Message)
	***REMOVED***
***REMOVED***

func isWellKnownType(t reflect.Type) bool ***REMOVED***
	if t.Implements(typeOfWkt) ***REMOVED***
		return true
	***REMOVED***
	if msg, ok := reflect.Zero(t).Interface().(proto.Message); ok ***REMOVED***
		name := proto.MessageName(msg)
		_, ok := wellKnownTypeNames[name]
		return ok
	***REMOVED***
	return false
***REMOVED***

// GetKnownType will return the reflect.Type for the given message name if it is
// known. If it is not known, nil is returned.
func (r *KnownTypeRegistry) GetKnownType(messageName string) reflect.Type ***REMOVED***
	var msgType reflect.Type
	if r == nil ***REMOVED***
		// a nil registry behaves the same as zero value instance: only know of well-known types
		t := proto.MessageType(messageName)
		if t != nil && isWellKnownType(t) ***REMOVED***
			msgType = t
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if r.includeDefault ***REMOVED***
			msgType = proto.MessageType(messageName)
		***REMOVED*** else if !r.excludeWkt ***REMOVED***
			t := proto.MessageType(messageName)
			if t != nil && isWellKnownType(t) ***REMOVED***
				msgType = t
			***REMOVED***
		***REMOVED***
		if msgType == nil ***REMOVED***
			r.mu.RLock()
			msgType = r.types[messageName]
			r.mu.RUnlock()
		***REMOVED***
	***REMOVED***

	return msgType
***REMOVED***
