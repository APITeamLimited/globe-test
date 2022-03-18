package sourceinfo

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// These are wrappers around the various interfaces in the
// google.golang.org/protobuf/reflect/protoreflect that all
// make sure to return a FileDescriptor that includes source
// code info.

type fileDescriptor struct ***REMOVED***
	protoreflect.FileDescriptor
	locs protoreflect.SourceLocations
***REMOVED***

func (f fileDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return f
***REMOVED***

func (f fileDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	return nil
***REMOVED***

func (f fileDescriptor) Imports() protoreflect.FileImports ***REMOVED***
	return imports***REMOVED***f.FileDescriptor.Imports()***REMOVED***
***REMOVED***

func (f fileDescriptor) Messages() protoreflect.MessageDescriptors ***REMOVED***
	return messages***REMOVED***f.FileDescriptor.Messages()***REMOVED***
***REMOVED***

func (f fileDescriptor) Enums() protoreflect.EnumDescriptors ***REMOVED***
	return enums***REMOVED***f.FileDescriptor.Enums()***REMOVED***
***REMOVED***

func (f fileDescriptor) Extensions() protoreflect.ExtensionDescriptors ***REMOVED***
	return extensions***REMOVED***f.FileDescriptor.Extensions()***REMOVED***
***REMOVED***

func (f fileDescriptor) Services() protoreflect.ServiceDescriptors ***REMOVED***
	return services***REMOVED***f.FileDescriptor.Services()***REMOVED***
***REMOVED***

func (f fileDescriptor) SourceLocations() protoreflect.SourceLocations ***REMOVED***
	return f.locs
***REMOVED***

type imports struct ***REMOVED***
	protoreflect.FileImports
***REMOVED***

func (im imports) Get(i int) protoreflect.FileImport ***REMOVED***
	fi := im.FileImports.Get(i)
	return protoreflect.FileImport***REMOVED***
		FileDescriptor: getFile(fi.FileDescriptor),
		IsPublic:       fi.IsPublic,
		IsWeak:         fi.IsWeak,
	***REMOVED***
***REMOVED***

type messages struct ***REMOVED***
	protoreflect.MessageDescriptors
***REMOVED***

func (m messages) Get(i int) protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***m.MessageDescriptors.Get(i)***REMOVED***
***REMOVED***

func (m messages) ByName(n protoreflect.Name) protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***m.MessageDescriptors.ByName(n)***REMOVED***
***REMOVED***

type enums struct ***REMOVED***
	protoreflect.EnumDescriptors
***REMOVED***

func (e enums) Get(i int) protoreflect.EnumDescriptor ***REMOVED***
	return enumDescriptor***REMOVED***e.EnumDescriptors.Get(i)***REMOVED***
***REMOVED***

func (e enums) ByName(n protoreflect.Name) protoreflect.EnumDescriptor ***REMOVED***
	return enumDescriptor***REMOVED***e.EnumDescriptors.ByName(n)***REMOVED***
***REMOVED***

type extensions struct ***REMOVED***
	protoreflect.ExtensionDescriptors
***REMOVED***

func (e extensions) Get(i int) protoreflect.ExtensionDescriptor ***REMOVED***
	d := e.ExtensionDescriptors.Get(i)
	if ed, ok := d.(protoreflect.ExtensionTypeDescriptor); ok ***REMOVED***
		return extensionDescriptor***REMOVED***ed***REMOVED***
	***REMOVED***
	return fieldDescriptor***REMOVED***d***REMOVED***
***REMOVED***

func (e extensions) ByName(n protoreflect.Name) protoreflect.ExtensionDescriptor ***REMOVED***
	d := e.ExtensionDescriptors.ByName(n)
	if ed, ok := d.(protoreflect.ExtensionTypeDescriptor); ok ***REMOVED***
		return extensionDescriptor***REMOVED***ed***REMOVED***
	***REMOVED***
	return fieldDescriptor***REMOVED***d***REMOVED***
***REMOVED***

type services struct ***REMOVED***
	protoreflect.ServiceDescriptors
***REMOVED***

func (s services) Get(i int) protoreflect.ServiceDescriptor ***REMOVED***
	return serviceDescriptor***REMOVED***s.ServiceDescriptors.Get(i)***REMOVED***
***REMOVED***

func (s services) ByName(n protoreflect.Name) protoreflect.ServiceDescriptor ***REMOVED***
	return serviceDescriptor***REMOVED***s.ServiceDescriptors.ByName(n)***REMOVED***
***REMOVED***

type messageDescriptor struct ***REMOVED***
	protoreflect.MessageDescriptor
***REMOVED***

func (m messageDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(m.MessageDescriptor.ParentFile())
***REMOVED***

func (m messageDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := m.MessageDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d)
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (m messageDescriptor) Fields() protoreflect.FieldDescriptors ***REMOVED***
	return fields***REMOVED***m.MessageDescriptor.Fields()***REMOVED***
***REMOVED***

func (m messageDescriptor) Oneofs() protoreflect.OneofDescriptors ***REMOVED***
	return oneOfs***REMOVED***m.MessageDescriptor.Oneofs()***REMOVED***
***REMOVED***

func (m messageDescriptor) Enums() protoreflect.EnumDescriptors ***REMOVED***
	return enums***REMOVED***m.MessageDescriptor.Enums()***REMOVED***
***REMOVED***

func (m messageDescriptor) Messages() protoreflect.MessageDescriptors ***REMOVED***
	return messages***REMOVED***m.MessageDescriptor.Messages()***REMOVED***
***REMOVED***

func (m messageDescriptor) Extensions() protoreflect.ExtensionDescriptors ***REMOVED***
	return extensions***REMOVED***m.MessageDescriptor.Extensions()***REMOVED***
***REMOVED***

type fields struct ***REMOVED***
	protoreflect.FieldDescriptors
***REMOVED***

func (f fields) Get(i int) protoreflect.FieldDescriptor ***REMOVED***
	return fieldDescriptor***REMOVED***f.FieldDescriptors.Get(i)***REMOVED***
***REMOVED***

func (f fields) ByName(n protoreflect.Name) protoreflect.FieldDescriptor ***REMOVED***
	return fieldDescriptor***REMOVED***f.FieldDescriptors.ByName(n)***REMOVED***
***REMOVED***

func (f fields) ByJSONName(n string) protoreflect.FieldDescriptor ***REMOVED***
	return fieldDescriptor***REMOVED***f.FieldDescriptors.ByJSONName(n)***REMOVED***
***REMOVED***

func (f fields) ByTextName(n string) protoreflect.FieldDescriptor ***REMOVED***
	return fieldDescriptor***REMOVED***f.FieldDescriptors.ByTextName(n)***REMOVED***
***REMOVED***

func (f fields) ByNumber(n protoreflect.FieldNumber) protoreflect.FieldDescriptor ***REMOVED***
	return fieldDescriptor***REMOVED***f.FieldDescriptors.ByNumber(n)***REMOVED***
***REMOVED***

type oneOfs struct ***REMOVED***
	protoreflect.OneofDescriptors
***REMOVED***

func (o oneOfs) Get(i int) protoreflect.OneofDescriptor ***REMOVED***
	return oneOfDescriptor***REMOVED***o.OneofDescriptors.Get(i)***REMOVED***
***REMOVED***

func (o oneOfs) ByName(n protoreflect.Name) protoreflect.OneofDescriptor ***REMOVED***
	return oneOfDescriptor***REMOVED***o.OneofDescriptors.ByName(n)***REMOVED***
***REMOVED***

type fieldDescriptor struct ***REMOVED***
	protoreflect.FieldDescriptor
***REMOVED***

func (f fieldDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(f.FieldDescriptor.ParentFile())
***REMOVED***

func (f fieldDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := f.FieldDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d)
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (f fieldDescriptor) MapKey() protoreflect.FieldDescriptor ***REMOVED***
	fd := f.FieldDescriptor.MapKey()
	if fd == nil ***REMOVED***
		return nil
	***REMOVED***
	return fieldDescriptor***REMOVED***fd***REMOVED***
***REMOVED***

func (f fieldDescriptor) MapValue() protoreflect.FieldDescriptor ***REMOVED***
	fd := f.FieldDescriptor.MapValue()
	if fd == nil ***REMOVED***
		return nil
	***REMOVED***
	return fieldDescriptor***REMOVED***fd***REMOVED***
***REMOVED***

func (f fieldDescriptor) DefaultEnumValue() protoreflect.EnumValueDescriptor ***REMOVED***
	ed := f.FieldDescriptor.DefaultEnumValue()
	if ed == nil ***REMOVED***
		return nil
	***REMOVED***
	return enumValueDescriptor***REMOVED***ed***REMOVED***
***REMOVED***

func (f fieldDescriptor) ContainingOneof() protoreflect.OneofDescriptor ***REMOVED***
	od := f.FieldDescriptor.ContainingOneof()
	if od == nil ***REMOVED***
		return nil
	***REMOVED***
	return oneOfDescriptor***REMOVED***od***REMOVED***
***REMOVED***

func (f fieldDescriptor) ContainingMessage() protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***f.FieldDescriptor.ContainingMessage()***REMOVED***
***REMOVED***

func (f fieldDescriptor) Enum() protoreflect.EnumDescriptor ***REMOVED***
	ed := f.FieldDescriptor.Enum()
	if ed == nil ***REMOVED***
		return nil
	***REMOVED***
	return enumDescriptor***REMOVED***ed***REMOVED***
***REMOVED***

func (f fieldDescriptor) Message() protoreflect.MessageDescriptor ***REMOVED***
	md := f.FieldDescriptor.Message()
	if md == nil ***REMOVED***
		return nil
	***REMOVED***
	return messageDescriptor***REMOVED***md***REMOVED***
***REMOVED***

type oneOfDescriptor struct ***REMOVED***
	protoreflect.OneofDescriptor
***REMOVED***

func (o oneOfDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(o.OneofDescriptor.ParentFile())
***REMOVED***

func (o oneOfDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := o.OneofDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (o oneOfDescriptor) Fields() protoreflect.FieldDescriptors ***REMOVED***
	return fields***REMOVED***o.OneofDescriptor.Fields()***REMOVED***
***REMOVED***

type enumDescriptor struct ***REMOVED***
	protoreflect.EnumDescriptor
***REMOVED***

func (e enumDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(e.EnumDescriptor.ParentFile())
***REMOVED***

func (e enumDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := e.EnumDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d)
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (e enumDescriptor) Values() protoreflect.EnumValueDescriptors ***REMOVED***
	return enumValues***REMOVED***e.EnumDescriptor.Values()***REMOVED***
***REMOVED***

type enumValues struct ***REMOVED***
	protoreflect.EnumValueDescriptors
***REMOVED***

func (e enumValues) Get(i int) protoreflect.EnumValueDescriptor ***REMOVED***
	return enumValueDescriptor***REMOVED***e.EnumValueDescriptors.Get(i)***REMOVED***
***REMOVED***

func (e enumValues) ByName(n protoreflect.Name) protoreflect.EnumValueDescriptor ***REMOVED***
	return enumValueDescriptor***REMOVED***e.EnumValueDescriptors.ByName(n)***REMOVED***
***REMOVED***

func (e enumValues) ByNumber(n protoreflect.EnumNumber) protoreflect.EnumValueDescriptor ***REMOVED***
	return enumValueDescriptor***REMOVED***e.EnumValueDescriptors.ByNumber(n)***REMOVED***
***REMOVED***

type enumValueDescriptor struct ***REMOVED***
	protoreflect.EnumValueDescriptor
***REMOVED***

func (e enumValueDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(e.EnumValueDescriptor.ParentFile())
***REMOVED***

func (e enumValueDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := e.EnumValueDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.EnumDescriptor:
		return enumDescriptor***REMOVED***d***REMOVED***
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

type extensionDescriptor struct ***REMOVED***
	protoreflect.ExtensionTypeDescriptor
***REMOVED***

func (e extensionDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(e.ExtensionTypeDescriptor.ParentFile())
***REMOVED***

func (e extensionDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := e.ExtensionTypeDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.MessageDescriptor:
		return messageDescriptor***REMOVED***d***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d)
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (e extensionDescriptor) MapKey() protoreflect.FieldDescriptor ***REMOVED***
	fd := e.ExtensionTypeDescriptor.MapKey()
	if fd == nil ***REMOVED***
		return nil
	***REMOVED***
	return fieldDescriptor***REMOVED***fd***REMOVED***
***REMOVED***

func (e extensionDescriptor) MapValue() protoreflect.FieldDescriptor ***REMOVED***
	fd := e.ExtensionTypeDescriptor.MapValue()
	if fd == nil ***REMOVED***
		return nil
	***REMOVED***
	return fieldDescriptor***REMOVED***fd***REMOVED***
***REMOVED***

func (e extensionDescriptor) DefaultEnumValue() protoreflect.EnumValueDescriptor ***REMOVED***
	ed := e.ExtensionTypeDescriptor.DefaultEnumValue()
	if ed == nil ***REMOVED***
		return nil
	***REMOVED***
	return enumValueDescriptor***REMOVED***ed***REMOVED***
***REMOVED***

func (e extensionDescriptor) ContainingOneof() protoreflect.OneofDescriptor ***REMOVED***
	od := e.ExtensionTypeDescriptor.ContainingOneof()
	if od == nil ***REMOVED***
		return nil
	***REMOVED***
	return oneOfDescriptor***REMOVED***od***REMOVED***
***REMOVED***

func (e extensionDescriptor) ContainingMessage() protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***e.ExtensionTypeDescriptor.ContainingMessage()***REMOVED***
***REMOVED***

func (e extensionDescriptor) Enum() protoreflect.EnumDescriptor ***REMOVED***
	ed := e.ExtensionTypeDescriptor.Enum()
	if ed == nil ***REMOVED***
		return nil
	***REMOVED***
	return enumDescriptor***REMOVED***ed***REMOVED***
***REMOVED***

func (e extensionDescriptor) Message() protoreflect.MessageDescriptor ***REMOVED***
	md := e.ExtensionTypeDescriptor.Message()
	if md == nil ***REMOVED***
		return nil
	***REMOVED***
	return messageDescriptor***REMOVED***md***REMOVED***
***REMOVED***

func (e extensionDescriptor) Descriptor() protoreflect.ExtensionDescriptor ***REMOVED***
	return e
***REMOVED***

var _ protoreflect.ExtensionTypeDescriptor = extensionDescriptor***REMOVED******REMOVED***

type serviceDescriptor struct ***REMOVED***
	protoreflect.ServiceDescriptor
***REMOVED***

func (s serviceDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(s.ServiceDescriptor.ParentFile())
***REMOVED***

func (s serviceDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := s.ServiceDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.FileDescriptor:
		return getFile(d)
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (s serviceDescriptor) Methods() protoreflect.MethodDescriptors ***REMOVED***
	return methods***REMOVED***s.ServiceDescriptor.Methods()***REMOVED***
***REMOVED***

type methods struct ***REMOVED***
	protoreflect.MethodDescriptors
***REMOVED***

func (m methods) Get(i int) protoreflect.MethodDescriptor ***REMOVED***
	return methodDescriptor***REMOVED***m.MethodDescriptors.Get(i)***REMOVED***
***REMOVED***

func (m methods) ByName(n protoreflect.Name) protoreflect.MethodDescriptor ***REMOVED***
	return methodDescriptor***REMOVED***m.MethodDescriptors.ByName(n)***REMOVED***
***REMOVED***

type methodDescriptor struct ***REMOVED***
	protoreflect.MethodDescriptor
***REMOVED***

func (m methodDescriptor) ParentFile() protoreflect.FileDescriptor ***REMOVED***
	return getFile(m.MethodDescriptor.ParentFile())
***REMOVED***

func (m methodDescriptor) Parent() protoreflect.Descriptor ***REMOVED***
	d := m.MethodDescriptor.Parent()
	switch d := d.(type) ***REMOVED***
	case protoreflect.ServiceDescriptor:
		return serviceDescriptor***REMOVED***d***REMOVED***
	case nil:
		return nil
	default:
		panic(fmt.Sprintf("unexpected descriptor type %T", d))
	***REMOVED***
***REMOVED***

func (m methodDescriptor) Input() protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***m.MethodDescriptor.Input()***REMOVED***
***REMOVED***

func (m methodDescriptor) Output() protoreflect.MessageDescriptor ***REMOVED***
	return messageDescriptor***REMOVED***m.MethodDescriptor.Output()***REMOVED***
***REMOVED***

type extensionType struct ***REMOVED***
	protoreflect.ExtensionType
***REMOVED***

func (e extensionType) TypeDescriptor() protoreflect.ExtensionTypeDescriptor ***REMOVED***
	return extensionDescriptor***REMOVED***e.ExtensionType.TypeDescriptor()***REMOVED***
***REMOVED***
