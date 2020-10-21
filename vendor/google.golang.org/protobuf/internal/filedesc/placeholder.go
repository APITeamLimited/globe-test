// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"google.golang.org/protobuf/internal/descopts"
	"google.golang.org/protobuf/internal/pragma"
	pref "google.golang.org/protobuf/reflect/protoreflect"
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

func (f PlaceholderFile) ParentFile() pref.FileDescriptor       ***REMOVED*** return f ***REMOVED***
func (f PlaceholderFile) Parent() pref.Descriptor               ***REMOVED*** return nil ***REMOVED***
func (f PlaceholderFile) Index() int                            ***REMOVED*** return 0 ***REMOVED***
func (f PlaceholderFile) Syntax() pref.Syntax                   ***REMOVED*** return 0 ***REMOVED***
func (f PlaceholderFile) Name() pref.Name                       ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) FullName() pref.FullName               ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) IsPlaceholder() bool                   ***REMOVED*** return true ***REMOVED***
func (f PlaceholderFile) Options() pref.ProtoMessage            ***REMOVED*** return descopts.File ***REMOVED***
func (f PlaceholderFile) Path() string                          ***REMOVED*** return string(f) ***REMOVED***
func (f PlaceholderFile) Package() pref.FullName                ***REMOVED*** return "" ***REMOVED***
func (f PlaceholderFile) Imports() pref.FileImports             ***REMOVED*** return emptyFiles ***REMOVED***
func (f PlaceholderFile) Messages() pref.MessageDescriptors     ***REMOVED*** return emptyMessages ***REMOVED***
func (f PlaceholderFile) Enums() pref.EnumDescriptors           ***REMOVED*** return emptyEnums ***REMOVED***
func (f PlaceholderFile) Extensions() pref.ExtensionDescriptors ***REMOVED*** return emptyExtensions ***REMOVED***
func (f PlaceholderFile) Services() pref.ServiceDescriptors     ***REMOVED*** return emptyServices ***REMOVED***
func (f PlaceholderFile) SourceLocations() pref.SourceLocations ***REMOVED*** return emptySourceLocations ***REMOVED***
func (f PlaceholderFile) ProtoType(pref.FileDescriptor)         ***REMOVED*** return ***REMOVED***
func (f PlaceholderFile) ProtoInternal(pragma.DoNotImplement)   ***REMOVED*** return ***REMOVED***

// PlaceholderEnum is a placeholder, representing only the full name.
type PlaceholderEnum pref.FullName

func (e PlaceholderEnum) ParentFile() pref.FileDescriptor     ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnum) Parent() pref.Descriptor             ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnum) Index() int                          ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnum) Syntax() pref.Syntax                 ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnum) Name() pref.Name                     ***REMOVED*** return pref.FullName(e).Name() ***REMOVED***
func (e PlaceholderEnum) FullName() pref.FullName             ***REMOVED*** return pref.FullName(e) ***REMOVED***
func (e PlaceholderEnum) IsPlaceholder() bool                 ***REMOVED*** return true ***REMOVED***
func (e PlaceholderEnum) Options() pref.ProtoMessage          ***REMOVED*** return descopts.Enum ***REMOVED***
func (e PlaceholderEnum) Values() pref.EnumValueDescriptors   ***REMOVED*** return emptyEnumValues ***REMOVED***
func (e PlaceholderEnum) ReservedNames() pref.Names           ***REMOVED*** return emptyNames ***REMOVED***
func (e PlaceholderEnum) ReservedRanges() pref.EnumRanges     ***REMOVED*** return emptyEnumRanges ***REMOVED***
func (e PlaceholderEnum) ProtoType(pref.EnumDescriptor)       ***REMOVED*** return ***REMOVED***
func (e PlaceholderEnum) ProtoInternal(pragma.DoNotImplement) ***REMOVED*** return ***REMOVED***

// PlaceholderEnumValue is a placeholder, representing only the full name.
type PlaceholderEnumValue pref.FullName

func (e PlaceholderEnumValue) ParentFile() pref.FileDescriptor     ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnumValue) Parent() pref.Descriptor             ***REMOVED*** return nil ***REMOVED***
func (e PlaceholderEnumValue) Index() int                          ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) Syntax() pref.Syntax                 ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) Name() pref.Name                     ***REMOVED*** return pref.FullName(e).Name() ***REMOVED***
func (e PlaceholderEnumValue) FullName() pref.FullName             ***REMOVED*** return pref.FullName(e) ***REMOVED***
func (e PlaceholderEnumValue) IsPlaceholder() bool                 ***REMOVED*** return true ***REMOVED***
func (e PlaceholderEnumValue) Options() pref.ProtoMessage          ***REMOVED*** return descopts.EnumValue ***REMOVED***
func (e PlaceholderEnumValue) Number() pref.EnumNumber             ***REMOVED*** return 0 ***REMOVED***
func (e PlaceholderEnumValue) ProtoType(pref.EnumValueDescriptor)  ***REMOVED*** return ***REMOVED***
func (e PlaceholderEnumValue) ProtoInternal(pragma.DoNotImplement) ***REMOVED*** return ***REMOVED***

// PlaceholderMessage is a placeholder, representing only the full name.
type PlaceholderMessage pref.FullName

func (m PlaceholderMessage) ParentFile() pref.FileDescriptor             ***REMOVED*** return nil ***REMOVED***
func (m PlaceholderMessage) Parent() pref.Descriptor                     ***REMOVED*** return nil ***REMOVED***
func (m PlaceholderMessage) Index() int                                  ***REMOVED*** return 0 ***REMOVED***
func (m PlaceholderMessage) Syntax() pref.Syntax                         ***REMOVED*** return 0 ***REMOVED***
func (m PlaceholderMessage) Name() pref.Name                             ***REMOVED*** return pref.FullName(m).Name() ***REMOVED***
func (m PlaceholderMessage) FullName() pref.FullName                     ***REMOVED*** return pref.FullName(m) ***REMOVED***
func (m PlaceholderMessage) IsPlaceholder() bool                         ***REMOVED*** return true ***REMOVED***
func (m PlaceholderMessage) Options() pref.ProtoMessage                  ***REMOVED*** return descopts.Message ***REMOVED***
func (m PlaceholderMessage) IsMapEntry() bool                            ***REMOVED*** return false ***REMOVED***
func (m PlaceholderMessage) Fields() pref.FieldDescriptors               ***REMOVED*** return emptyFields ***REMOVED***
func (m PlaceholderMessage) Oneofs() pref.OneofDescriptors               ***REMOVED*** return emptyOneofs ***REMOVED***
func (m PlaceholderMessage) ReservedNames() pref.Names                   ***REMOVED*** return emptyNames ***REMOVED***
func (m PlaceholderMessage) ReservedRanges() pref.FieldRanges            ***REMOVED*** return emptyFieldRanges ***REMOVED***
func (m PlaceholderMessage) RequiredNumbers() pref.FieldNumbers          ***REMOVED*** return emptyFieldNumbers ***REMOVED***
func (m PlaceholderMessage) ExtensionRanges() pref.FieldRanges           ***REMOVED*** return emptyFieldRanges ***REMOVED***
func (m PlaceholderMessage) ExtensionRangeOptions(int) pref.ProtoMessage ***REMOVED*** panic("index out of range") ***REMOVED***
func (m PlaceholderMessage) Messages() pref.MessageDescriptors           ***REMOVED*** return emptyMessages ***REMOVED***
func (m PlaceholderMessage) Enums() pref.EnumDescriptors                 ***REMOVED*** return emptyEnums ***REMOVED***
func (m PlaceholderMessage) Extensions() pref.ExtensionDescriptors       ***REMOVED*** return emptyExtensions ***REMOVED***
func (m PlaceholderMessage) ProtoType(pref.MessageDescriptor)            ***REMOVED*** return ***REMOVED***
func (m PlaceholderMessage) ProtoInternal(pragma.DoNotImplement)         ***REMOVED*** return ***REMOVED***
