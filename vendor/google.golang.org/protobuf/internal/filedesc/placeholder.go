// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/pragma"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	emptyNames           = new(Names)
	emptyEnumRanges      = new(EnumRanges)
	emptyFieldRanges     = new(FieldRanges)
	emptyFieldNumbers    = new(FieldNumbers)
	emptySourceLocations = new(SourceLocations)

	emptyFiles      = new(FileImports)
	emptyMessages   = new(Messages)
	emptyFields     = new(Fields)
	emptyOneofs     = new(Oneofs)
	emptyEnums      = new(Enums)
	emptyEnumValues = new(EnumValues)
	emptyExtensions = new(Extensions)
	emptyServices   = new(Services)
)

// PlaceholderFile is a placeholder, representing only the file path.
type PlaceholderFile string

func (f PlaceholderFile) ParentFile() protoreflect.FileDescriptor       ***REMOVED*** return f ***REMOVED***
func (f PlaceholderFile) Parent() protoreflect.Descriptor               ***REMOVED*** return nil ***REMOVED***
func (f PlaceholderFile) Index() int                                    ***REMOVED*** return 0 ***REMOVED***
func (f PlaceholderFile) Syntax() protoreflect.Syntax                   ***REMOVED*** return 0 ***REMOVED***
func (f PlaceholderFile) Name() protoreflect.Name                       ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) FullName() protoreflect.FullName               ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) IsPlaceholder() bool                           ***REMOVED*** return true ***REMOVED***
func (f PlaceholderFile) Options() protoreflect.ProtoMessage            ***REMOVED*** return descopts.File ***REMOVED***
func (f PlaceholderFile) Path() string                                  ***REMOVED*** return string(f) ***REMOVED***
func (f PlaceholderFile) Package() protoreflect.FullName                ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) Imports() protoreflect.FileImports             ***REMOVED*** return emptyFiles ***REMOVED***
func (f PlaceholderFile) Messages() protoreflect.MessageDescriptors     ***REMOVED*** return emptyMessages ***REMOVED***
func (f PlaceholderFile) Enums() protoreflect.EnumDescriptors           ***REMOVED*** return emptyEnums ***REMOVED***
func (f PlaceholderFile) Extensions() protoreflect.ExtensionDescriptors ***REMOVED*** return emptyExtensions ***REMOVED***
func (f PlaceholderFile) Services() protoreflect.ServiceDescriptors     ***REMOVED*** return emptyServices ***REMOVED***
func (f PlaceholderFile) SourceLocations() protoreflect.SourceLocations ***REMOVED*** return emptySourceLocations ***REMOVED***
func (f PlaceholderFile) ProtoType(protoreflect.FileDescriptor)         ***REMOVED*** return ***REMOVED***
func (f PlaceholderFile) ProtoInternal(pragma.DoNotImplement)           ***REMOVED*** return ***REMOVED***

// PlaceholderEnum is a placeholder, representing only the full name.
type PlaceholderEnum protoreflect.FullName

func (e PlaceholderEnum) ParentFile() protoreflect.FileDescriptor   ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnum) Parent() protoreflect.Descriptor           ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnum) Index() int                                ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnum) Syntax() protoreflect.Syntax               ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnum) Name() protoreflect.Name                   ***REMOVED*** return protoreflect.FullName(e).Name() ***REMOVED***
func (e PlaceholderEnum) FullName() protoreflect.FullName           ***REMOVED*** return protoreflect.FullName(e) ***REMOVED***
func (e PlaceholderEnum) IsPlaceholder() bool                       ***REMOVED*** return true ***REMOVED***
func (e PlaceholderEnum) Options() protoreflect.ProtoMessage        ***REMOVED*** return descopts.Enum ***REMOVED***
func (e PlaceholderEnum) Values() protoreflect.EnumValueDescriptors ***REMOVED*** return emptyEnumValues ***REMOVED***
func (e PlaceholderEnum) ReservedNames() protoreflect.Names         ***REMOVED*** return emptyNames ***REMOVED***
func (e PlaceholderEnum) ReservedRanges() protoreflect.EnumRanges   ***REMOVED*** return emptyEnumRanges ***REMOVED***
func (e PlaceholderEnum) ProtoType(protoreflect.EnumDescriptor)     ***REMOVED*** return ***REMOVED***
func (e PlaceholderEnum) ProtoInternal(pragma.DoNotImplement)       ***REMOVED*** return ***REMOVED***

// PlaceholderEnumValue is a placeholder, representing only the full name.
type PlaceholderEnumValue protoreflect.FullName

func (e PlaceholderEnumValue) ParentFile() protoreflect.FileDescriptor    ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnumValue) Parent() protoreflect.Descriptor            ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnumValue) Index() int                                 ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) Syntax() protoreflect.Syntax                ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) Name() protoreflect.Name                    ***REMOVED*** return protoreflect.FullName(e).Name() ***REMOVED***
func (e PlaceholderEnumValue) FullName() protoreflect.FullName            ***REMOVED*** return protoreflect.FullName(e) ***REMOVED***
func (e PlaceholderEnumValue) IsPlaceholder() bool                        ***REMOVED*** return true ***REMOVED***
func (e PlaceholderEnumValue) Options() protoreflect.ProtoMessage         ***REMOVED*** return descopts.EnumValue ***REMOVED***
func (e PlaceholderEnumValue) Number() protoreflect.EnumNumber            ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) ProtoType(protoreflect.EnumValueDescriptor) ***REMOVED*** return ***REMOVED***
func (e PlaceholderEnumValue) ProtoInternal(pragma.DoNotImplement)        ***REMOVED*** return ***REMOVED***

// PlaceholderMessage is a placeholder, representing only the full name.
type PlaceholderMessage protoreflect.FullName

func (m PlaceholderMessage) ParentFile() protoreflect.FileDescriptor    ***REMOVED*** return nil ***REMOVED***
func (m PlaceholderMessage) Parent() protoreflect.Descriptor            ***REMOVED*** return nil ***REMOVED***
func (m PlaceholderMessage) Index() int                                 ***REMOVED*** return 0 ***REMOVED***
func (m PlaceholderMessage) Syntax() protoreflect.Syntax                ***REMOVED*** return 0 ***REMOVED***
func (m PlaceholderMessage) Name() protoreflect.Name                    ***REMOVED*** return protoreflect.FullName(m).Name() ***REMOVED***
func (m PlaceholderMessage) FullName() protoreflect.FullName            ***REMOVED*** return protoreflect.FullName(m) ***REMOVED***
func (m PlaceholderMessage) IsPlaceholder() bool                        ***REMOVED*** return true ***REMOVED***
func (m PlaceholderMessage) Options() protoreflect.ProtoMessage         ***REMOVED*** return descopts.Message ***REMOVED***
func (m PlaceholderMessage) IsMapEntry() bool                           ***REMOVED*** return false ***REMOVED***
func (m PlaceholderMessage) Fields() protoreflect.FieldDescriptors      ***REMOVED*** return emptyFields ***REMOVED***
func (m PlaceholderMessage) Oneofs() protoreflect.OneofDescriptors      ***REMOVED*** return emptyOneofs ***REMOVED***
func (m PlaceholderMessage) ReservedNames() protoreflect.Names          ***REMOVED*** return emptyNames ***REMOVED***
func (m PlaceholderMessage) ReservedRanges() protoreflect.FieldRanges   ***REMOVED*** return emptyFieldRanges ***REMOVED***
func (m PlaceholderMessage) RequiredNumbers() protoreflect.FieldNumbers ***REMOVED*** return emptyFieldNumbers ***REMOVED***
func (m PlaceholderMessage) ExtensionRanges() protoreflect.FieldRanges  ***REMOVED*** return emptyFieldRanges ***REMOVED***
func (m PlaceholderMessage) ExtensionRangeOptions(int) protoreflect.ProtoMessage ***REMOVED***
	panic("index out of range")
***REMOVED***
func (m PlaceholderMessage) Messages() protoreflect.MessageDescriptors     ***REMOVED*** return emptyMessages ***REMOVED***
func (m PlaceholderMessage) Enums() protoreflect.EnumDescriptors           ***REMOVED*** return emptyEnums ***REMOVED***
func (m PlaceholderMessage) Extensions() protoreflect.ExtensionDescriptors ***REMOVED*** return emptyExtensions ***REMOVED***
func (m PlaceholderMessage) ProtoType(protoreflect.MessageDescriptor)      ***REMOVED*** return ***REMOVED***
func (m PlaceholderMessage) ProtoInternal(pragma.DoNotImplement)           ***REMOVED*** return ***REMOVED***
