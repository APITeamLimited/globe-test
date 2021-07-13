package protoparse

import (
	"bytes"
	"fmt"
	"math"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/internal"
	"github.com/jhump/protoreflect/desc/protoparse/ast"
	"github.com/jhump/protoreflect/dynamic"
)

// NB: To process options, we need descriptors, but we may not have rich
// descriptors when trying to interpret options for unlinked parsed files.
// So we define minimal interfaces that can be backed by both rich descriptors
// as well as their poorer cousins, plain ol' descriptor protos.

type descriptorish interface ***REMOVED***
	GetFile() fileDescriptorish
	GetFullyQualifiedName() string
	AsProto() proto.Message
***REMOVED***

type fileDescriptorish interface ***REMOVED***
	descriptorish
	GetFileOptions() *dpb.FileOptions
	GetPackage() string
	FindSymbol(name string) desc.Descriptor
	GetPublicDependencies() []fileDescriptorish
	GetDependencies() []fileDescriptorish
	GetMessageTypes() []msgDescriptorish
	GetExtensions() []fldDescriptorish
	GetEnumTypes() []enumDescriptorish
	GetServices() []svcDescriptorish
***REMOVED***

type msgDescriptorish interface ***REMOVED***
	descriptorish
	GetMessageOptions() *dpb.MessageOptions
	GetFields() []fldDescriptorish
	GetOneOfs() []oneofDescriptorish
	GetExtensionRanges() []extRangeDescriptorish
	GetNestedMessageTypes() []msgDescriptorish
	GetNestedExtensions() []fldDescriptorish
	GetNestedEnumTypes() []enumDescriptorish
***REMOVED***

type fldDescriptorish interface ***REMOVED***
	descriptorish
	GetFieldOptions() *dpb.FieldOptions
	GetMessageType() *desc.MessageDescriptor
	GetEnumType() *desc.EnumDescriptor
	AsFieldDescriptorProto() *dpb.FieldDescriptorProto
***REMOVED***

type oneofDescriptorish interface ***REMOVED***
	descriptorish
	GetOneOfOptions() *dpb.OneofOptions
***REMOVED***

type enumDescriptorish interface ***REMOVED***
	descriptorish
	GetEnumOptions() *dpb.EnumOptions
	GetValues() []enumValDescriptorish
***REMOVED***

type enumValDescriptorish interface ***REMOVED***
	descriptorish
	GetEnumValueOptions() *dpb.EnumValueOptions
***REMOVED***

type svcDescriptorish interface ***REMOVED***
	descriptorish
	GetServiceOptions() *dpb.ServiceOptions
	GetMethods() []methodDescriptorish
***REMOVED***

type methodDescriptorish interface ***REMOVED***
	descriptorish
	GetMethodOptions() *dpb.MethodOptions
***REMOVED***

// The hierarchy of descriptorish implementations backed by
// rich descriptors:

type richFileDescriptorish struct ***REMOVED***
	*desc.FileDescriptor
***REMOVED***

func (d richFileDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d
***REMOVED***

func (d richFileDescriptorish) GetPublicDependencies() []fileDescriptorish ***REMOVED***
	deps := d.FileDescriptor.GetPublicDependencies()
	ret := make([]fileDescriptorish, len(deps))
	for i, d := range deps ***REMOVED***
		ret[i] = richFileDescriptorish***REMOVED***FileDescriptor: d***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richFileDescriptorish) GetDependencies() []fileDescriptorish ***REMOVED***
	deps := d.FileDescriptor.GetDependencies()
	ret := make([]fileDescriptorish, len(deps))
	for i, d := range deps ***REMOVED***
		ret[i] = richFileDescriptorish***REMOVED***FileDescriptor: d***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richFileDescriptorish) GetMessageTypes() []msgDescriptorish ***REMOVED***
	msgs := d.FileDescriptor.GetMessageTypes()
	ret := make([]msgDescriptorish, len(msgs))
	for i, m := range msgs ***REMOVED***
		ret[i] = richMsgDescriptorish***REMOVED***MessageDescriptor: m***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richFileDescriptorish) GetExtensions() []fldDescriptorish ***REMOVED***
	flds := d.FileDescriptor.GetExtensions()
	ret := make([]fldDescriptorish, len(flds))
	for i, f := range flds ***REMOVED***
		ret[i] = richFldDescriptorish***REMOVED***FieldDescriptor: f***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richFileDescriptorish) GetEnumTypes() []enumDescriptorish ***REMOVED***
	ens := d.FileDescriptor.GetEnumTypes()
	ret := make([]enumDescriptorish, len(ens))
	for i, en := range ens ***REMOVED***
		ret[i] = richEnumDescriptorish***REMOVED***EnumDescriptor: en***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richFileDescriptorish) GetServices() []svcDescriptorish ***REMOVED***
	svcs := d.FileDescriptor.GetServices()
	ret := make([]svcDescriptorish, len(svcs))
	for i, s := range svcs ***REMOVED***
		ret[i] = richSvcDescriptorish***REMOVED***ServiceDescriptor: s***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type richMsgDescriptorish struct ***REMOVED***
	*desc.MessageDescriptor
***REMOVED***

func (d richMsgDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.MessageDescriptor.GetFile()***REMOVED***
***REMOVED***

func (d richMsgDescriptorish) GetFields() []fldDescriptorish ***REMOVED***
	flds := d.MessageDescriptor.GetFields()
	ret := make([]fldDescriptorish, len(flds))
	for i, f := range flds ***REMOVED***
		ret[i] = richFldDescriptorish***REMOVED***FieldDescriptor: f***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richMsgDescriptorish) GetOneOfs() []oneofDescriptorish ***REMOVED***
	oos := d.MessageDescriptor.GetOneOfs()
	ret := make([]oneofDescriptorish, len(oos))
	for i, oo := range oos ***REMOVED***
		ret[i] = richOneOfDescriptorish***REMOVED***OneOfDescriptor: oo***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richMsgDescriptorish) GetExtensionRanges() []extRangeDescriptorish ***REMOVED***
	md := d.MessageDescriptor
	mdFqn := md.GetFullyQualifiedName()
	extrs := md.AsDescriptorProto().GetExtensionRange()
	ret := make([]extRangeDescriptorish, len(extrs))
	for i, extr := range extrs ***REMOVED***
		ret[i] = extRangeDescriptorish***REMOVED***
			er:   extr,
			qual: mdFqn,
			file: richFileDescriptorish***REMOVED***FileDescriptor: md.GetFile()***REMOVED***,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richMsgDescriptorish) GetNestedMessageTypes() []msgDescriptorish ***REMOVED***
	msgs := d.MessageDescriptor.GetNestedMessageTypes()
	ret := make([]msgDescriptorish, len(msgs))
	for i, m := range msgs ***REMOVED***
		ret[i] = richMsgDescriptorish***REMOVED***MessageDescriptor: m***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richMsgDescriptorish) GetNestedExtensions() []fldDescriptorish ***REMOVED***
	flds := d.MessageDescriptor.GetNestedExtensions()
	ret := make([]fldDescriptorish, len(flds))
	for i, f := range flds ***REMOVED***
		ret[i] = richFldDescriptorish***REMOVED***FieldDescriptor: f***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d richMsgDescriptorish) GetNestedEnumTypes() []enumDescriptorish ***REMOVED***
	ens := d.MessageDescriptor.GetNestedEnumTypes()
	ret := make([]enumDescriptorish, len(ens))
	for i, en := range ens ***REMOVED***
		ret[i] = richEnumDescriptorish***REMOVED***EnumDescriptor: en***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type richFldDescriptorish struct ***REMOVED***
	*desc.FieldDescriptor
***REMOVED***

func (d richFldDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.FieldDescriptor.GetFile()***REMOVED***
***REMOVED***

func (d richFldDescriptorish) AsFieldDescriptorProto() *dpb.FieldDescriptorProto ***REMOVED***
	return d.FieldDescriptor.AsFieldDescriptorProto()
***REMOVED***

type richOneOfDescriptorish struct ***REMOVED***
	*desc.OneOfDescriptor
***REMOVED***

func (d richOneOfDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.OneOfDescriptor.GetFile()***REMOVED***
***REMOVED***

type richEnumDescriptorish struct ***REMOVED***
	*desc.EnumDescriptor
***REMOVED***

func (d richEnumDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.EnumDescriptor.GetFile()***REMOVED***
***REMOVED***

func (d richEnumDescriptorish) GetValues() []enumValDescriptorish ***REMOVED***
	vals := d.EnumDescriptor.GetValues()
	ret := make([]enumValDescriptorish, len(vals))
	for i, val := range vals ***REMOVED***
		ret[i] = richEnumValDescriptorish***REMOVED***EnumValueDescriptor: val***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type richEnumValDescriptorish struct ***REMOVED***
	*desc.EnumValueDescriptor
***REMOVED***

func (d richEnumValDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.EnumValueDescriptor.GetFile()***REMOVED***
***REMOVED***

type richSvcDescriptorish struct ***REMOVED***
	*desc.ServiceDescriptor
***REMOVED***

func (d richSvcDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.ServiceDescriptor.GetFile()***REMOVED***
***REMOVED***

func (d richSvcDescriptorish) GetMethods() []methodDescriptorish ***REMOVED***
	mtds := d.ServiceDescriptor.GetMethods()
	ret := make([]methodDescriptorish, len(mtds))
	for i, mtd := range mtds ***REMOVED***
		ret[i] = richMethodDescriptorish***REMOVED***MethodDescriptor: mtd***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type richMethodDescriptorish struct ***REMOVED***
	*desc.MethodDescriptor
***REMOVED***

func (d richMethodDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return richFileDescriptorish***REMOVED***FileDescriptor: d.MethodDescriptor.GetFile()***REMOVED***
***REMOVED***

// The hierarchy of descriptorish implementations backed by
// plain descriptor protos:

type poorFileDescriptorish struct ***REMOVED***
	*dpb.FileDescriptorProto
***REMOVED***

func (d poorFileDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d
***REMOVED***

func (d poorFileDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return d.FileDescriptorProto.GetName()
***REMOVED***

func (d poorFileDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.FileDescriptorProto
***REMOVED***

func (d poorFileDescriptorish) GetFileOptions() *dpb.FileOptions ***REMOVED***
	return d.FileDescriptorProto.GetOptions()
***REMOVED***

func (d poorFileDescriptorish) FindSymbol(name string) desc.Descriptor ***REMOVED***
	return nil
***REMOVED***

func (d poorFileDescriptorish) GetPublicDependencies() []fileDescriptorish ***REMOVED***
	return nil
***REMOVED***

func (d poorFileDescriptorish) GetDependencies() []fileDescriptorish ***REMOVED***
	return nil
***REMOVED***

func (d poorFileDescriptorish) GetMessageTypes() []msgDescriptorish ***REMOVED***
	msgs := d.FileDescriptorProto.GetMessageType()
	pkg := d.FileDescriptorProto.GetPackage()
	ret := make([]msgDescriptorish, len(msgs))
	for i, m := range msgs ***REMOVED***
		ret[i] = poorMsgDescriptorish***REMOVED***
			DescriptorProto: m,
			qual:            pkg,
			file:            d,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorFileDescriptorish) GetExtensions() []fldDescriptorish ***REMOVED***
	exts := d.FileDescriptorProto.GetExtension()
	pkg := d.FileDescriptorProto.GetPackage()
	ret := make([]fldDescriptorish, len(exts))
	for i, e := range exts ***REMOVED***
		ret[i] = poorFldDescriptorish***REMOVED***
			FieldDescriptorProto: e,
			qual:                 pkg,
			file:                 d,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorFileDescriptorish) GetEnumTypes() []enumDescriptorish ***REMOVED***
	ens := d.FileDescriptorProto.GetEnumType()
	pkg := d.FileDescriptorProto.GetPackage()
	ret := make([]enumDescriptorish, len(ens))
	for i, e := range ens ***REMOVED***
		ret[i] = poorEnumDescriptorish***REMOVED***
			EnumDescriptorProto: e,
			qual:                pkg,
			file:                d,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorFileDescriptorish) GetServices() []svcDescriptorish ***REMOVED***
	svcs := d.FileDescriptorProto.GetService()
	pkg := d.FileDescriptorProto.GetPackage()
	ret := make([]svcDescriptorish, len(svcs))
	for i, s := range svcs ***REMOVED***
		ret[i] = poorSvcDescriptorish***REMOVED***
			ServiceDescriptorProto: s,
			qual:                   pkg,
			file:                   d,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type poorMsgDescriptorish struct ***REMOVED***
	*dpb.DescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorMsgDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorMsgDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.DescriptorProto.GetName())
***REMOVED***

func qualify(qual, name string) string ***REMOVED***
	if qual == "" ***REMOVED***
		return name
	***REMOVED*** else ***REMOVED***
		return fmt.Sprintf("%s.%s", qual, name)
	***REMOVED***
***REMOVED***

func (d poorMsgDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.DescriptorProto
***REMOVED***

func (d poorMsgDescriptorish) GetMessageOptions() *dpb.MessageOptions ***REMOVED***
	return d.DescriptorProto.GetOptions()
***REMOVED***

func (d poorMsgDescriptorish) GetFields() []fldDescriptorish ***REMOVED***
	flds := d.DescriptorProto.GetField()
	ret := make([]fldDescriptorish, len(flds))
	for i, f := range flds ***REMOVED***
		ret[i] = poorFldDescriptorish***REMOVED***
			FieldDescriptorProto: f,
			qual:                 d.GetFullyQualifiedName(),
			file:                 d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorMsgDescriptorish) GetOneOfs() []oneofDescriptorish ***REMOVED***
	oos := d.DescriptorProto.GetOneofDecl()
	ret := make([]oneofDescriptorish, len(oos))
	for i, oo := range oos ***REMOVED***
		ret[i] = poorOneOfDescriptorish***REMOVED***
			OneofDescriptorProto: oo,
			qual:                 d.GetFullyQualifiedName(),
			file:                 d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorMsgDescriptorish) GetExtensionRanges() []extRangeDescriptorish ***REMOVED***
	mdFqn := d.GetFullyQualifiedName()
	extrs := d.DescriptorProto.GetExtensionRange()
	ret := make([]extRangeDescriptorish, len(extrs))
	for i, extr := range extrs ***REMOVED***
		ret[i] = extRangeDescriptorish***REMOVED***
			er:   extr,
			qual: mdFqn,
			file: d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorMsgDescriptorish) GetNestedMessageTypes() []msgDescriptorish ***REMOVED***
	msgs := d.DescriptorProto.GetNestedType()
	ret := make([]msgDescriptorish, len(msgs))
	for i, m := range msgs ***REMOVED***
		ret[i] = poorMsgDescriptorish***REMOVED***
			DescriptorProto: m,
			qual:            d.GetFullyQualifiedName(),
			file:            d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorMsgDescriptorish) GetNestedExtensions() []fldDescriptorish ***REMOVED***
	flds := d.DescriptorProto.GetExtension()
	ret := make([]fldDescriptorish, len(flds))
	for i, f := range flds ***REMOVED***
		ret[i] = poorFldDescriptorish***REMOVED***
			FieldDescriptorProto: f,
			qual:                 d.GetFullyQualifiedName(),
			file:                 d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

func (d poorMsgDescriptorish) GetNestedEnumTypes() []enumDescriptorish ***REMOVED***
	ens := d.DescriptorProto.GetEnumType()
	ret := make([]enumDescriptorish, len(ens))
	for i, en := range ens ***REMOVED***
		ret[i] = poorEnumDescriptorish***REMOVED***
			EnumDescriptorProto: en,
			qual:                d.GetFullyQualifiedName(),
			file:                d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type poorFldDescriptorish struct ***REMOVED***
	*dpb.FieldDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorFldDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorFldDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.FieldDescriptorProto.GetName())
***REMOVED***

func (d poorFldDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.FieldDescriptorProto
***REMOVED***

func (d poorFldDescriptorish) GetFieldOptions() *dpb.FieldOptions ***REMOVED***
	return d.FieldDescriptorProto.GetOptions()
***REMOVED***

func (d poorFldDescriptorish) GetMessageType() *desc.MessageDescriptor ***REMOVED***
	return nil
***REMOVED***

func (d poorFldDescriptorish) GetEnumType() *desc.EnumDescriptor ***REMOVED***
	return nil
***REMOVED***

type poorOneOfDescriptorish struct ***REMOVED***
	*dpb.OneofDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorOneOfDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorOneOfDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.OneofDescriptorProto.GetName())
***REMOVED***

func (d poorOneOfDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.OneofDescriptorProto
***REMOVED***

func (d poorOneOfDescriptorish) GetOneOfOptions() *dpb.OneofOptions ***REMOVED***
	return d.OneofDescriptorProto.GetOptions()
***REMOVED***

func (d poorFldDescriptorish) AsFieldDescriptorProto() *dpb.FieldDescriptorProto ***REMOVED***
	return d.FieldDescriptorProto
***REMOVED***

type poorEnumDescriptorish struct ***REMOVED***
	*dpb.EnumDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorEnumDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorEnumDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.EnumDescriptorProto.GetName())
***REMOVED***

func (d poorEnumDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.EnumDescriptorProto
***REMOVED***

func (d poorEnumDescriptorish) GetEnumOptions() *dpb.EnumOptions ***REMOVED***
	return d.EnumDescriptorProto.GetOptions()
***REMOVED***

func (d poorEnumDescriptorish) GetValues() []enumValDescriptorish ***REMOVED***
	vals := d.EnumDescriptorProto.GetValue()
	ret := make([]enumValDescriptorish, len(vals))
	for i, v := range vals ***REMOVED***
		ret[i] = poorEnumValDescriptorish***REMOVED***
			EnumValueDescriptorProto: v,
			qual:                     d.GetFullyQualifiedName(),
			file:                     d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type poorEnumValDescriptorish struct ***REMOVED***
	*dpb.EnumValueDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorEnumValDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorEnumValDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.EnumValueDescriptorProto.GetName())
***REMOVED***

func (d poorEnumValDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.EnumValueDescriptorProto
***REMOVED***

func (d poorEnumValDescriptorish) GetEnumValueOptions() *dpb.EnumValueOptions ***REMOVED***
	return d.EnumValueDescriptorProto.GetOptions()
***REMOVED***

type poorSvcDescriptorish struct ***REMOVED***
	*dpb.ServiceDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorSvcDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorSvcDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.ServiceDescriptorProto.GetName())
***REMOVED***

func (d poorSvcDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.ServiceDescriptorProto
***REMOVED***

func (d poorSvcDescriptorish) GetServiceOptions() *dpb.ServiceOptions ***REMOVED***
	return d.ServiceDescriptorProto.GetOptions()
***REMOVED***

func (d poorSvcDescriptorish) GetMethods() []methodDescriptorish ***REMOVED***
	mtds := d.ServiceDescriptorProto.GetMethod()
	ret := make([]methodDescriptorish, len(mtds))
	for i, m := range mtds ***REMOVED***
		ret[i] = poorMethodDescriptorish***REMOVED***
			MethodDescriptorProto: m,
			qual:                  d.GetFullyQualifiedName(),
			file:                  d.file,
		***REMOVED***
	***REMOVED***
	return ret
***REMOVED***

type poorMethodDescriptorish struct ***REMOVED***
	*dpb.MethodDescriptorProto
	qual string
	file fileDescriptorish
***REMOVED***

func (d poorMethodDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return d.file
***REMOVED***

func (d poorMethodDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(d.qual, d.MethodDescriptorProto.GetName())
***REMOVED***

func (d poorMethodDescriptorish) AsProto() proto.Message ***REMOVED***
	return d.MethodDescriptorProto
***REMOVED***

func (d poorMethodDescriptorish) GetMethodOptions() *dpb.MethodOptions ***REMOVED***
	return d.MethodDescriptorProto.GetOptions()
***REMOVED***

type extRangeDescriptorish struct ***REMOVED***
	er   *dpb.DescriptorProto_ExtensionRange
	qual string
	file fileDescriptorish
***REMOVED***

func (er extRangeDescriptorish) GetFile() fileDescriptorish ***REMOVED***
	return er.file
***REMOVED***

func (er extRangeDescriptorish) GetFullyQualifiedName() string ***REMOVED***
	return qualify(er.qual, fmt.Sprintf("%d-%d", er.er.GetStart(), er.er.GetEnd()-1))
***REMOVED***

func (er extRangeDescriptorish) AsProto() proto.Message ***REMOVED***
	return er.er
***REMOVED***

func (er extRangeDescriptorish) GetExtensionRangeOptions() *dpb.ExtensionRangeOptions ***REMOVED***
	return er.er.GetOptions()
***REMOVED***

func interpretFileOptions(l *linker, r *parseResult, fd fileDescriptorish) error ***REMOVED***
	opts := fd.GetFileOptions()
	if opts != nil ***REMOVED***
		if len(opts.UninterpretedOption) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, fd, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, md := range fd.GetMessageTypes() ***REMOVED***
		if err := interpretMessageOptions(l, r, md); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, fld := range fd.GetExtensions() ***REMOVED***
		if err := interpretFieldOptions(l, r, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ed := range fd.GetEnumTypes() ***REMOVED***
		if err := interpretEnumOptions(l, r, ed); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, sd := range fd.GetServices() ***REMOVED***
		opts := sd.GetServiceOptions()
		if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, sd, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
		for _, mtd := range sd.GetMethods() ***REMOVED***
			opts := mtd.GetMethodOptions()
			if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
				if remain, err := interpretOptions(l, r, mtd, opts, opts.UninterpretedOption); err != nil ***REMOVED***
					return err
				***REMOVED*** else ***REMOVED***
					opts.UninterpretedOption = remain
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func interpretMessageOptions(l *linker, r *parseResult, md msgDescriptorish) error ***REMOVED***
	opts := md.GetMessageOptions()
	if opts != nil ***REMOVED***
		if len(opts.UninterpretedOption) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, md, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, fld := range md.GetFields() ***REMOVED***
		if err := interpretFieldOptions(l, r, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ood := range md.GetOneOfs() ***REMOVED***
		opts := ood.GetOneOfOptions()
		if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, ood, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, fld := range md.GetNestedExtensions() ***REMOVED***
		if err := interpretFieldOptions(l, r, fld); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, er := range md.GetExtensionRanges() ***REMOVED***
		opts := er.GetExtensionRangeOptions()
		if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, er, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, nmd := range md.GetNestedMessageTypes() ***REMOVED***
		if err := interpretMessageOptions(l, r, nmd); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, ed := range md.GetNestedEnumTypes() ***REMOVED***
		if err := interpretEnumOptions(l, r, ed); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func interpretFieldOptions(l *linker, r *parseResult, fld fldDescriptorish) error ***REMOVED***
	opts := fld.GetFieldOptions()
	if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
		uo := opts.UninterpretedOption
		scope := fmt.Sprintf("field %s", fld.GetFullyQualifiedName())

		// process json_name pseudo-option
		if index, err := findOption(r, scope, uo, "json_name"); err != nil && !r.lenient ***REMOVED***
			return err
		***REMOVED*** else if index >= 0 ***REMOVED***
			opt := uo[index]
			optNode := r.getOptionNode(opt)

			// attribute source code info
			if on, ok := optNode.(*ast.OptionNode); ok ***REMOVED***
				r.interpretedOptions[on] = []int32***REMOVED***-1, internal.Field_jsonNameTag***REMOVED***
			***REMOVED***
			uo = removeOption(uo, index)
			if opt.StringValue == nil ***REMOVED***
				if err := r.errs.handleErrorWithPos(optNode.GetValue().Start(), "%s: expecting string value for json_name option", scope); err != nil ***REMOVED***
					return err
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				fld.AsFieldDescriptorProto().JsonName = proto.String(string(opt.StringValue))
			***REMOVED***
		***REMOVED***

		// and process default pseudo-option
		if index, err := processDefaultOption(r, scope, fld, uo); err != nil && !r.lenient ***REMOVED***
			return err
		***REMOVED*** else if index >= 0 ***REMOVED***
			// attribute source code info
			optNode := r.getOptionNode(uo[index])
			if on, ok := optNode.(*ast.OptionNode); ok ***REMOVED***
				r.interpretedOptions[on] = []int32***REMOVED***-1, internal.Field_defaultTag***REMOVED***
			***REMOVED***
			uo = removeOption(uo, index)
		***REMOVED***

		if len(uo) == 0 ***REMOVED***
			// no real options, only pseudo-options above? clear out options
			fld.AsFieldDescriptorProto().Options = nil
		***REMOVED*** else if remain, err := interpretOptions(l, r, fld, opts, uo); err != nil ***REMOVED***
			return err
		***REMOVED*** else ***REMOVED***
			opts.UninterpretedOption = remain
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func processDefaultOption(res *parseResult, scope string, fld fldDescriptorish, uos []*dpb.UninterpretedOption) (defaultIndex int, err error) ***REMOVED***
	found, err := findOption(res, scope, uos, "default")
	if err != nil || found == -1 ***REMOVED***
		return -1, err
	***REMOVED***
	opt := uos[found]
	optNode := res.getOptionNode(opt)
	fdp := fld.AsFieldDescriptorProto()
	if fdp.GetLabel() == dpb.FieldDescriptorProto_LABEL_REPEATED ***REMOVED***
		return -1, res.errs.handleErrorWithPos(optNode.GetName().Start(), "%s: default value cannot be set because field is repeated", scope)
	***REMOVED***
	if fdp.GetType() == dpb.FieldDescriptorProto_TYPE_GROUP || fdp.GetType() == dpb.FieldDescriptorProto_TYPE_MESSAGE ***REMOVED***
		return -1, res.errs.handleErrorWithPos(optNode.GetName().Start(), "%s: default value cannot be set because field is a message", scope)
	***REMOVED***
	val := optNode.GetValue()
	if _, ok := val.(*ast.MessageLiteralNode); ok ***REMOVED***
		return -1, res.errs.handleErrorWithPos(val.Start(), "%s: default value cannot be a message", scope)
	***REMOVED***
	mc := &messageContext***REMOVED***
		res:         res,
		file:        fld.GetFile(),
		elementName: fld.GetFullyQualifiedName(),
		elementType: descriptorType(fld.AsProto()),
		option:      opt,
	***REMOVED***
	v, err := fieldValue(res, mc, fld, val, true)
	if err != nil ***REMOVED***
		return -1, res.errs.handleError(err)
	***REMOVED***
	if str, ok := v.(string); ok ***REMOVED***
		fld.AsFieldDescriptorProto().DefaultValue = proto.String(str)
	***REMOVED*** else if b, ok := v.([]byte); ok ***REMOVED***
		fld.AsFieldDescriptorProto().DefaultValue = proto.String(encodeDefaultBytes(b))
	***REMOVED*** else ***REMOVED***
		var flt float64
		var ok bool
		if flt, ok = v.(float64); !ok ***REMOVED***
			var flt32 float32
			if flt32, ok = v.(float32); ok ***REMOVED***
				flt = float64(flt32)
			***REMOVED***
		***REMOVED***
		if ok ***REMOVED***
			if math.IsInf(flt, 1) ***REMOVED***
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("inf")
			***REMOVED*** else if ok && math.IsInf(flt, -1) ***REMOVED***
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("-inf")
			***REMOVED*** else if ok && math.IsNaN(flt) ***REMOVED***
				fld.AsFieldDescriptorProto().DefaultValue = proto.String("nan")
			***REMOVED*** else ***REMOVED***
				fld.AsFieldDescriptorProto().DefaultValue = proto.String(fmt.Sprintf("%v", v))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			fld.AsFieldDescriptorProto().DefaultValue = proto.String(fmt.Sprintf("%v", v))
		***REMOVED***
	***REMOVED***
	return found, nil
***REMOVED***

func encodeDefaultBytes(b []byte) string ***REMOVED***
	var buf bytes.Buffer
	writeEscapedBytes(&buf, b)
	return buf.String()
***REMOVED***

func interpretEnumOptions(l *linker, r *parseResult, ed enumDescriptorish) error ***REMOVED***
	opts := ed.GetEnumOptions()
	if opts != nil ***REMOVED***
		if len(opts.UninterpretedOption) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, ed, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for _, evd := range ed.GetValues() ***REMOVED***
		opts := evd.GetEnumValueOptions()
		if len(opts.GetUninterpretedOption()) > 0 ***REMOVED***
			if remain, err := interpretOptions(l, r, evd, opts, opts.UninterpretedOption); err != nil ***REMOVED***
				return err
			***REMOVED*** else ***REMOVED***
				opts.UninterpretedOption = remain
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func interpretOptions(l *linker, res *parseResult, element descriptorish, opts proto.Message, uninterpreted []*dpb.UninterpretedOption) ([]*dpb.UninterpretedOption, error) ***REMOVED***
	optsd, err := loadMessageDescriptorForOptions(l, element.GetFile(), opts)
	if err != nil ***REMOVED***
		if res.lenient ***REMOVED***
			return uninterpreted, nil
		***REMOVED***
		return nil, res.errs.handleError(err)
	***REMOVED***
	dm := dynamic.NewMessage(optsd)
	err = dm.ConvertFrom(opts)
	if err != nil ***REMOVED***
		if res.lenient ***REMOVED***
			return uninterpreted, nil
		***REMOVED***
		node := res.nodes[element.AsProto()]
		return nil, res.errs.handleError(ErrorWithSourcePos***REMOVED***Pos: node.Start(), Underlying: err***REMOVED***)
	***REMOVED***

	mc := &messageContext***REMOVED***res: res, file: element.GetFile(), elementName: element.GetFullyQualifiedName(), elementType: descriptorType(element.AsProto())***REMOVED***
	var remain []*dpb.UninterpretedOption
	for _, uo := range uninterpreted ***REMOVED***
		node := res.getOptionNode(uo)
		if !uo.Name[0].GetIsExtension() && uo.Name[0].GetNamePart() == "uninterpreted_option" ***REMOVED***
			if res.lenient ***REMOVED***
				remain = append(remain, uo)
				continue
			***REMOVED***
			// uninterpreted_option might be found reflectively, but is not actually valid for use
			if err := res.errs.handleErrorWithPos(node.GetName().Start(), "%vinvalid option 'uninterpreted_option'", mc); err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***
		mc.option = uo
		path, err := interpretField(res, mc, element, dm, uo, 0, nil)
		if err != nil ***REMOVED***
			if res.lenient ***REMOVED***
				remain = append(remain, uo)
				continue
			***REMOVED***
			return nil, err
		***REMOVED***
		if optn, ok := node.(*ast.OptionNode); ok ***REMOVED***
			res.interpretedOptions[optn] = path
		***REMOVED***
	***REMOVED***

	if res.lenient ***REMOVED***
		// If we're lenient, then we don't want to clobber the passed in message
		// and leave it partially populated. So we convert into a copy first
		optsClone := proto.Clone(opts)
		if err := dm.ConvertToDeterministic(optsClone); err != nil ***REMOVED***
			// TODO: do this in a more granular way, so we can convert individual
			// fields and leave bad ones uninterpreted instead of skipping all of
			// the work we've done so far.
			return uninterpreted, nil
		***REMOVED***
		// conversion from dynamic message above worked, so now
		// it is safe to overwrite the passed in message
		opts.Reset()
		proto.Merge(opts, optsClone)

		return remain, nil
	***REMOVED***

	if err := dm.ValidateRecursive(); err != nil ***REMOVED***
		node := res.nodes[element.AsProto()]
		if err := res.errs.handleErrorWithPos(node.Start(), "error in %s options: %v", descriptorType(element.AsProto()), err); err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// now try to convert into the passed in message and fail if not successful
	if err := dm.ConvertToDeterministic(opts); err != nil ***REMOVED***
		node := res.nodes[element.AsProto()]
		return nil, res.errs.handleError(ErrorWithSourcePos***REMOVED***Pos: node.Start(), Underlying: err***REMOVED***)
	***REMOVED***

	return nil, nil
***REMOVED***

func loadMessageDescriptorForOptions(l *linker, fd fileDescriptorish, opts proto.Message) (*desc.MessageDescriptor, error) ***REMOVED***
	// see if the file imports a custom version of descriptor.proto
	fqn := proto.MessageName(opts)
	d := findMessageDescriptorForOptions(l, fd, fqn)
	if d != nil ***REMOVED***
		return d, nil
	***REMOVED***
	// fall back to built-in options descriptors
	return desc.LoadMessageDescriptorForMessage(opts)
***REMOVED***

func findMessageDescriptorForOptions(l *linker, fd fileDescriptorish, messageName string) *desc.MessageDescriptor ***REMOVED***
	d := fd.FindSymbol(messageName)
	if d != nil ***REMOVED***
		md, _ := d.(*desc.MessageDescriptor)
		return md
	***REMOVED***

	// TODO: should this support public imports and be recursive?
	for _, dep := range fd.GetDependencies() ***REMOVED***
		d := dep.FindSymbol(messageName)
		if d != nil ***REMOVED***
			if l != nil ***REMOVED***
				l.markUsed(fd.AsProto().(*dpb.FileDescriptorProto), d.GetFile().AsFileDescriptorProto())
			***REMOVED***
			md, _ := d.(*desc.MessageDescriptor)
			return md
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func interpretField(res *parseResult, mc *messageContext, element descriptorish, dm *dynamic.Message, opt *dpb.UninterpretedOption, nameIndex int, pathPrefix []int32) (path []int32, err error) ***REMOVED***
	var fld *desc.FieldDescriptor
	nm := opt.GetName()[nameIndex]
	node := res.getOptionNamePartNode(nm)
	if nm.GetIsExtension() ***REMOVED***
		extName := nm.GetNamePart()
		if extName[0] == '.' ***REMOVED***
			extName = extName[1:] /* skip leading dot */
		***REMOVED***
		fld = findExtension(element.GetFile(), extName, false, map[fileDescriptorish]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
		if fld == nil ***REMOVED***
			return nil, res.errs.handleErrorWithPos(node.Start(),
				"%vunrecognized extension %s of %s",
				mc, extName, dm.GetMessageDescriptor().GetFullyQualifiedName())
		***REMOVED***
		if fld.GetOwner().GetFullyQualifiedName() != dm.GetMessageDescriptor().GetFullyQualifiedName() ***REMOVED***
			return nil, res.errs.handleErrorWithPos(node.Start(),
				"%vextension %s should extend %s but instead extends %s",
				mc, extName, dm.GetMessageDescriptor().GetFullyQualifiedName(), fld.GetOwner().GetFullyQualifiedName())
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		fld = dm.GetMessageDescriptor().FindFieldByName(nm.GetNamePart())
		if fld == nil ***REMOVED***
			return nil, res.errs.handleErrorWithPos(node.Start(),
				"%vfield %s of %s does not exist",
				mc, nm.GetNamePart(), dm.GetMessageDescriptor().GetFullyQualifiedName())
		***REMOVED***
	***REMOVED***

	path = append(pathPrefix, fld.GetNumber())

	if len(opt.GetName()) > nameIndex+1 ***REMOVED***
		nextnm := opt.GetName()[nameIndex+1]
		nextnode := res.getOptionNamePartNode(nextnm)
		if fld.GetType() != dpb.FieldDescriptorProto_TYPE_MESSAGE ***REMOVED***
			return nil, res.errs.handleErrorWithPos(nextnode.Start(),
				"%vcannot set field %s because %s is not a message",
				mc, nextnm.GetNamePart(), nm.GetNamePart())
		***REMOVED***
		if fld.IsRepeated() ***REMOVED***
			return nil, res.errs.handleErrorWithPos(nextnode.Start(),
				"%vcannot set field %s because %s is repeated (must use an aggregate)",
				mc, nextnm.GetNamePart(), nm.GetNamePart())
		***REMOVED***
		var fdm *dynamic.Message
		var err error
		if dm.HasField(fld) ***REMOVED***
			var v interface***REMOVED******REMOVED***
			v, err = dm.TryGetField(fld)
			fdm, _ = v.(*dynamic.Message)
		***REMOVED*** else ***REMOVED***
			fdm = dynamic.NewMessage(fld.GetMessageType())
			err = dm.TrySetField(fld, fdm)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, res.errs.handleError(ErrorWithSourcePos***REMOVED***Pos: node.Start(), Underlying: err***REMOVED***)
		***REMOVED***
		// recurse to set next part of name
		return interpretField(res, mc, element, fdm, opt, nameIndex+1, path)
	***REMOVED***

	optNode := res.getOptionNode(opt)
	if err := setOptionField(res, mc, dm, fld, node, optNode.GetValue()); err != nil ***REMOVED***
		return nil, res.errs.handleError(err)
	***REMOVED***
	if fld.IsRepeated() ***REMOVED***
		path = append(path, int32(dm.FieldLength(fld))-1)
	***REMOVED***
	return path, nil
***REMOVED***

func findExtension(fd fileDescriptorish, name string, public bool, checked map[fileDescriptorish]struct***REMOVED******REMOVED***) *desc.FieldDescriptor ***REMOVED***
	if _, ok := checked[fd]; ok ***REMOVED***
		return nil
	***REMOVED***
	checked[fd] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	d := fd.FindSymbol(name)
	if d != nil ***REMOVED***
		if fld, ok := d.(*desc.FieldDescriptor); ok ***REMOVED***
			return fld
		***REMOVED***
		return nil
	***REMOVED***

	// When public = false, we are searching only directly imported symbols. But we
	// also need to search transitive public imports due to semantics of public imports.
	if public ***REMOVED***
		for _, dep := range fd.GetPublicDependencies() ***REMOVED***
			d := findExtension(dep, name, true, checked)
			if d != nil ***REMOVED***
				return d
			***REMOVED***
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		for _, dep := range fd.GetDependencies() ***REMOVED***
			d := findExtension(dep, name, true, checked)
			if d != nil ***REMOVED***
				return d
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func setOptionField(res *parseResult, mc *messageContext, dm *dynamic.Message, fld *desc.FieldDescriptor, name ast.Node, val ast.ValueNode) error ***REMOVED***
	v := val.Value()
	if sl, ok := v.([]ast.ValueNode); ok ***REMOVED***
		// handle slices a little differently than the others
		if !fld.IsRepeated() ***REMOVED***
			return errorWithPos(val.Start(), "%vvalue is an array but field is not repeated", mc)
		***REMOVED***
		origPath := mc.optAggPath
		defer func() ***REMOVED***
			mc.optAggPath = origPath
		***REMOVED***()
		for index, item := range sl ***REMOVED***
			mc.optAggPath = fmt.Sprintf("%s[%d]", origPath, index)
			if v, err := fieldValue(res, mc, richFldDescriptorish***REMOVED***FieldDescriptor: fld***REMOVED***, item, false); err != nil ***REMOVED***
				return err
			***REMOVED*** else if err = dm.TryAddRepeatedField(fld, v); err != nil ***REMOVED***
				return errorWithPos(val.Start(), "%verror setting value: %s", mc, err)
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***

	v, err := fieldValue(res, mc, richFldDescriptorish***REMOVED***FieldDescriptor: fld***REMOVED***, val, false)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if fld.IsRepeated() ***REMOVED***
		err = dm.TryAddRepeatedField(fld, v)
	***REMOVED*** else ***REMOVED***
		if dm.HasField(fld) ***REMOVED***
			return errorWithPos(name.Start(), "%vnon-repeated option field %s already set", mc, fieldName(fld))
		***REMOVED***
		err = dm.TrySetField(fld, v)
	***REMOVED***
	if err != nil ***REMOVED***
		return errorWithPos(val.Start(), "%verror setting value: %s", mc, err)
	***REMOVED***

	return nil
***REMOVED***

func findOption(res *parseResult, scope string, opts []*dpb.UninterpretedOption, name string) (int, error) ***REMOVED***
	found := -1
	for i, opt := range opts ***REMOVED***
		if len(opt.Name) != 1 ***REMOVED***
			continue
		***REMOVED***
		if opt.Name[0].GetIsExtension() || opt.Name[0].GetNamePart() != name ***REMOVED***
			continue
		***REMOVED***
		if found >= 0 ***REMOVED***
			optNode := res.getOptionNode(opt)
			return -1, res.errs.handleErrorWithPos(optNode.GetName().Start(), "%s: option %s cannot be defined more than once", scope, name)
		***REMOVED***
		found = i
	***REMOVED***
	return found, nil
***REMOVED***

func removeOption(uo []*dpb.UninterpretedOption, indexToRemove int) []*dpb.UninterpretedOption ***REMOVED***
	if indexToRemove == 0 ***REMOVED***
		return uo[1:]
	***REMOVED*** else if int(indexToRemove) == len(uo)-1 ***REMOVED***
		return uo[:len(uo)-1]
	***REMOVED*** else ***REMOVED***
		return append(uo[:indexToRemove], uo[indexToRemove+1:]...)
	***REMOVED***
***REMOVED***

type messageContext struct ***REMOVED***
	res         *parseResult
	file        fileDescriptorish
	elementType string
	elementName string
	option      *dpb.UninterpretedOption
	optAggPath  string
***REMOVED***

func (c *messageContext) String() string ***REMOVED***
	var ctx bytes.Buffer
	if c.elementType != "file" ***REMOVED***
		_, _ = fmt.Fprintf(&ctx, "%s %s: ", c.elementType, c.elementName)
	***REMOVED***
	if c.option != nil && c.option.Name != nil ***REMOVED***
		ctx.WriteString("option ")
		writeOptionName(&ctx, c.option.Name)
		if c.res.nodes == nil ***REMOVED***
			// if we have no source position info, try to provide as much context
			// as possible (if nodes != nil, we don't need this because any errors
			// will actually have file and line numbers)
			if c.optAggPath != "" ***REMOVED***
				_, _ = fmt.Fprintf(&ctx, " at %s", c.optAggPath)
			***REMOVED***
		***REMOVED***
		ctx.WriteString(": ")
	***REMOVED***
	return ctx.String()
***REMOVED***

func writeOptionName(buf *bytes.Buffer, parts []*dpb.UninterpretedOption_NamePart) ***REMOVED***
	first := true
	for _, p := range parts ***REMOVED***
		if first ***REMOVED***
			first = false
		***REMOVED*** else ***REMOVED***
			buf.WriteByte('.')
		***REMOVED***
		nm := p.GetNamePart()
		if nm[0] == '.' ***REMOVED***
			// skip leading dot
			nm = nm[1:]
		***REMOVED***
		if p.GetIsExtension() ***REMOVED***
			buf.WriteByte('(')
			buf.WriteString(nm)
			buf.WriteByte(')')
		***REMOVED*** else ***REMOVED***
			buf.WriteString(nm)
		***REMOVED***
	***REMOVED***
***REMOVED***

func fieldName(fld *desc.FieldDescriptor) string ***REMOVED***
	if fld.IsExtension() ***REMOVED***
		return fld.GetFullyQualifiedName()
	***REMOVED*** else ***REMOVED***
		return fld.GetName()
	***REMOVED***
***REMOVED***

func valueKind(val interface***REMOVED******REMOVED***) string ***REMOVED***
	switch val := val.(type) ***REMOVED***
	case ast.Identifier:
		return "identifier"
	case bool:
		return "bool"
	case int64:
		if val < 0 ***REMOVED***
			return "negative integer"
		***REMOVED***
		return "integer"
	case uint64:
		return "integer"
	case float64:
		return "double"
	case string, []byte:
		return "string"
	case []*ast.MessageFieldNode:
		return "message"
	case []ast.ValueNode:
		return "array"
	default:
		return fmt.Sprintf("%T", val)
	***REMOVED***
***REMOVED***

func fieldValue(res *parseResult, mc *messageContext, fld fldDescriptorish, val ast.ValueNode, enumAsString bool) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	v := val.Value()
	t := fld.AsFieldDescriptorProto().GetType()
	switch t ***REMOVED***
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		if id, ok := v.(ast.Identifier); ok ***REMOVED***
			ev := fld.GetEnumType().FindValueByName(string(id))
			if ev == nil ***REMOVED***
				return nil, errorWithPos(val.Start(), "%venum %s has no value named %s", mc, fld.GetEnumType().GetFullyQualifiedName(), id)
			***REMOVED***
			if enumAsString ***REMOVED***
				return ev.GetName(), nil
			***REMOVED*** else ***REMOVED***
				return ev.GetNumber(), nil
			***REMOVED***
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting enum, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_MESSAGE, dpb.FieldDescriptorProto_TYPE_GROUP:
		if aggs, ok := v.([]*ast.MessageFieldNode); ok ***REMOVED***
			fmd := fld.GetMessageType()
			fdm := dynamic.NewMessage(fmd)
			origPath := mc.optAggPath
			defer func() ***REMOVED***
				mc.optAggPath = origPath
			***REMOVED***()
			for _, a := range aggs ***REMOVED***
				if origPath == "" ***REMOVED***
					mc.optAggPath = a.Name.Value()
				***REMOVED*** else ***REMOVED***
					mc.optAggPath = origPath + "." + a.Name.Value()
				***REMOVED***
				var ffld *desc.FieldDescriptor
				if a.Name.IsExtension() ***REMOVED***
					n := string(a.Name.Name.AsIdentifier())
					ffld = findExtension(mc.file, n, false, map[fileDescriptorish]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
					if ffld == nil ***REMOVED***
						// may need to qualify with package name
						pkg := mc.file.GetPackage()
						if pkg != "" ***REMOVED***
							ffld = findExtension(mc.file, pkg+"."+n, false, map[fileDescriptorish]struct***REMOVED******REMOVED******REMOVED******REMOVED***)
						***REMOVED***
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					ffld = fmd.FindFieldByName(a.Name.Value())
				***REMOVED***
				if ffld == nil ***REMOVED***
					return nil, errorWithPos(val.Start(), "%vfield %s not found", mc, string(a.Name.Name.AsIdentifier()))
				***REMOVED***
				if err := setOptionField(res, mc, fdm, ffld, a.Name, a.Val); err != nil ***REMOVED***
					return nil, err
				***REMOVED***
			***REMOVED***
			return fdm, nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting message, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		if b, ok := v.(bool); ok ***REMOVED***
			return b, nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting bool, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		if str, ok := v.(string); ok ***REMOVED***
			return []byte(str), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting bytes, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_STRING:
		if str, ok := v.(string); ok ***REMOVED***
			return str, nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting string, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_INT32, dpb.FieldDescriptorProto_TYPE_SINT32, dpb.FieldDescriptorProto_TYPE_SFIXED32:
		if i, ok := v.(int64); ok ***REMOVED***
			if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for int32", mc, i)
			***REMOVED***
			return int32(i), nil
		***REMOVED***
		if ui, ok := v.(uint64); ok ***REMOVED***
			if ui > math.MaxInt32 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for int32", mc, ui)
			***REMOVED***
			return int32(ui), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting int32, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_UINT32, dpb.FieldDescriptorProto_TYPE_FIXED32:
		if i, ok := v.(int64); ok ***REMOVED***
			if i > math.MaxUint32 || i < 0 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for uint32", mc, i)
			***REMOVED***
			return uint32(i), nil
		***REMOVED***
		if ui, ok := v.(uint64); ok ***REMOVED***
			if ui > math.MaxUint32 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for uint32", mc, ui)
			***REMOVED***
			return uint32(ui), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting uint32, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_INT64, dpb.FieldDescriptorProto_TYPE_SINT64, dpb.FieldDescriptorProto_TYPE_SFIXED64:
		if i, ok := v.(int64); ok ***REMOVED***
			return i, nil
		***REMOVED***
		if ui, ok := v.(uint64); ok ***REMOVED***
			if ui > math.MaxInt64 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for int64", mc, ui)
			***REMOVED***
			return int64(ui), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting int64, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_UINT64, dpb.FieldDescriptorProto_TYPE_FIXED64:
		if i, ok := v.(int64); ok ***REMOVED***
			if i < 0 ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %d is out of range for uint64", mc, i)
			***REMOVED***
			return uint64(i), nil
		***REMOVED***
		if ui, ok := v.(uint64); ok ***REMOVED***
			return ui, nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting uint64, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_DOUBLE:
		if d, ok := v.(float64); ok ***REMOVED***
			return d, nil
		***REMOVED***
		if i, ok := v.(int64); ok ***REMOVED***
			return float64(i), nil
		***REMOVED***
		if u, ok := v.(uint64); ok ***REMOVED***
			return float64(u), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting double, got %s", mc, valueKind(v))
	case dpb.FieldDescriptorProto_TYPE_FLOAT:
		if d, ok := v.(float64); ok ***REMOVED***
			if (d > math.MaxFloat32 || d < -math.MaxFloat32) && !math.IsInf(d, 1) && !math.IsInf(d, -1) && !math.IsNaN(d) ***REMOVED***
				return nil, errorWithPos(val.Start(), "%vvalue %f is out of range for float", mc, d)
			***REMOVED***
			return float32(d), nil
		***REMOVED***
		if i, ok := v.(int64); ok ***REMOVED***
			return float32(i), nil
		***REMOVED***
		if u, ok := v.(uint64); ok ***REMOVED***
			return float32(u), nil
		***REMOVED***
		return nil, errorWithPos(val.Start(), "%vexpecting float, got %s", mc, valueKind(v))
	default:
		return nil, errorWithPos(val.Start(), "%vunrecognized field type: %s", mc, t)
	***REMOVED***
***REMOVED***
