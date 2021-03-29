// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package filedesc

import (
	"sync"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/genid"
	"google.golang.org/protobuf/internal/strs"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// fileRaw is a data struct used when initializing a file descriptor from
// a raw FileDescriptorProto.
type fileRaw struct ***REMOVED***
	builder       Builder
	allEnums      []Enum
	allMessages   []Message
	allExtensions []Extension
	allServices   []Service
***REMOVED***

func newRawFile(db Builder) *File ***REMOVED***
	fd := &File***REMOVED***fileRaw: fileRaw***REMOVED***builder: db***REMOVED******REMOVED***
	fd.initDecls(db.NumEnums, db.NumMessages, db.NumExtensions, db.NumServices)
	fd.unmarshalSeed(db.RawDescriptor)

	// Extended message targets are eagerly resolved since registration
	// needs this information at program init time.
	for i := range fd.allExtensions ***REMOVED***
		xd := &fd.allExtensions[i]
		xd.L1.Extendee = fd.resolveMessageDependency(xd.L1.Extendee, listExtTargets, int32(i))
	***REMOVED***

	fd.checkDecls()
	return fd
***REMOVED***

// initDecls pre-allocates slices for the exact number of enums, messages
// (including map entries), extensions, and services declared in the proto file.
// This is done to avoid regrowing the slice, which would change the address
// for any previously seen declaration.
//
// The alloc methods "allocates" slices by pulling from the capacity.
func (fd *File) initDecls(numEnums, numMessages, numExtensions, numServices int32) ***REMOVED***
	fd.allEnums = make([]Enum, 0, numEnums)
	fd.allMessages = make([]Message, 0, numMessages)
	fd.allExtensions = make([]Extension, 0, numExtensions)
	fd.allServices = make([]Service, 0, numServices)
***REMOVED***

func (fd *File) allocEnums(n int) []Enum ***REMOVED***
	total := len(fd.allEnums)
	es := fd.allEnums[total : total+n]
	fd.allEnums = fd.allEnums[:total+n]
	return es
***REMOVED***
func (fd *File) allocMessages(n int) []Message ***REMOVED***
	total := len(fd.allMessages)
	ms := fd.allMessages[total : total+n]
	fd.allMessages = fd.allMessages[:total+n]
	return ms
***REMOVED***
func (fd *File) allocExtensions(n int) []Extension ***REMOVED***
	total := len(fd.allExtensions)
	xs := fd.allExtensions[total : total+n]
	fd.allExtensions = fd.allExtensions[:total+n]
	return xs
***REMOVED***
func (fd *File) allocServices(n int) []Service ***REMOVED***
	total := len(fd.allServices)
	xs := fd.allServices[total : total+n]
	fd.allServices = fd.allServices[:total+n]
	return xs
***REMOVED***

// checkDecls performs a sanity check that the expected number of expected
// declarations matches the number that were found in the descriptor proto.
func (fd *File) checkDecls() ***REMOVED***
	switch ***REMOVED***
	case len(fd.allEnums) != cap(fd.allEnums):
	case len(fd.allMessages) != cap(fd.allMessages):
	case len(fd.allExtensions) != cap(fd.allExtensions):
	case len(fd.allServices) != cap(fd.allServices):
	default:
		return
	***REMOVED***
	panic("mismatching cardinality")
***REMOVED***

func (fd *File) unmarshalSeed(b []byte) ***REMOVED***
	sb := getBuilder()
	defer putBuilder(sb)

	var prevField pref.FieldNumber
	var numEnums, numMessages, numExtensions, numServices int
	var posEnums, posMessages, posExtensions, posServices int
	b0 := b
	for len(b) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.FileDescriptorProto_Syntax_field_number:
				switch string(v) ***REMOVED***
				case "proto2":
					fd.L1.Syntax = pref.Proto2
				case "proto3":
					fd.L1.Syntax = pref.Proto3
				default:
					panic("invalid syntax")
				***REMOVED***
			case genid.FileDescriptorProto_Name_field_number:
				fd.L1.Path = sb.MakeString(v)
			case genid.FileDescriptorProto_Package_field_number:
				fd.L1.Package = pref.FullName(sb.MakeString(v))
			case genid.FileDescriptorProto_EnumType_field_number:
				if prevField != genid.FileDescriptorProto_EnumType_field_number ***REMOVED***
					if numEnums > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posEnums = len(b0) - len(b) - n - m
				***REMOVED***
				numEnums++
			case genid.FileDescriptorProto_MessageType_field_number:
				if prevField != genid.FileDescriptorProto_MessageType_field_number ***REMOVED***
					if numMessages > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posMessages = len(b0) - len(b) - n - m
				***REMOVED***
				numMessages++
			case genid.FileDescriptorProto_Extension_field_number:
				if prevField != genid.FileDescriptorProto_Extension_field_number ***REMOVED***
					if numExtensions > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posExtensions = len(b0) - len(b) - n - m
				***REMOVED***
				numExtensions++
			case genid.FileDescriptorProto_Service_field_number:
				if prevField != genid.FileDescriptorProto_Service_field_number ***REMOVED***
					if numServices > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posServices = len(b0) - len(b) - n - m
				***REMOVED***
				numServices++
			***REMOVED***
			prevField = num
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
			prevField = -1 // ignore known field numbers of unknown wire type
		***REMOVED***
	***REMOVED***

	// If syntax is missing, it is assumed to be proto2.
	if fd.L1.Syntax == 0 ***REMOVED***
		fd.L1.Syntax = pref.Proto2
	***REMOVED***

	// Must allocate all declarations before parsing each descriptor type
	// to ensure we handled all descriptors in "flattened ordering".
	if numEnums > 0 ***REMOVED***
		fd.L1.Enums.List = fd.allocEnums(numEnums)
	***REMOVED***
	if numMessages > 0 ***REMOVED***
		fd.L1.Messages.List = fd.allocMessages(numMessages)
	***REMOVED***
	if numExtensions > 0 ***REMOVED***
		fd.L1.Extensions.List = fd.allocExtensions(numExtensions)
	***REMOVED***
	if numServices > 0 ***REMOVED***
		fd.L1.Services.List = fd.allocServices(numServices)
	***REMOVED***

	if numEnums > 0 ***REMOVED***
		b := b0[posEnums:]
		for i := range fd.L1.Enums.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			fd.L1.Enums.List[i].unmarshalSeed(v, sb, fd, fd, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
	if numMessages > 0 ***REMOVED***
		b := b0[posMessages:]
		for i := range fd.L1.Messages.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			fd.L1.Messages.List[i].unmarshalSeed(v, sb, fd, fd, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
	if numExtensions > 0 ***REMOVED***
		b := b0[posExtensions:]
		for i := range fd.L1.Extensions.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			fd.L1.Extensions.List[i].unmarshalSeed(v, sb, fd, fd, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
	if numServices > 0 ***REMOVED***
		b := b0[posServices:]
		for i := range fd.L1.Services.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			fd.L1.Services.List[i].unmarshalSeed(v, sb, fd, fd, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ed *Enum) unmarshalSeed(b []byte, sb *strs.Builder, pf *File, pd pref.Descriptor, i int) ***REMOVED***
	ed.L0.ParentFile = pf
	ed.L0.Parent = pd
	ed.L0.Index = i

	var numValues int
	for b := b; len(b) > 0; ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.EnumDescriptorProto_Name_field_number:
				ed.L0.FullName = appendFullName(sb, pd.FullName(), v)
			case genid.EnumDescriptorProto_Value_field_number:
				numValues++
			***REMOVED***
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		***REMOVED***
	***REMOVED***

	// Only construct enum value descriptors for top-level enums since
	// they are needed for registration.
	if pd != pf ***REMOVED***
		return
	***REMOVED***
	ed.L1.eagerValues = true
	ed.L2 = new(EnumL2)
	ed.L2.Values.List = make([]EnumValue, numValues)
	for i := 0; len(b) > 0; ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.EnumDescriptorProto_Value_field_number:
				ed.L2.Values.List[i].unmarshalFull(v, sb, pf, ed, i)
				i++
			***REMOVED***
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (md *Message) unmarshalSeed(b []byte, sb *strs.Builder, pf *File, pd pref.Descriptor, i int) ***REMOVED***
	md.L0.ParentFile = pf
	md.L0.Parent = pd
	md.L0.Index = i

	var prevField pref.FieldNumber
	var numEnums, numMessages, numExtensions int
	var posEnums, posMessages, posExtensions int
	b0 := b
	for len(b) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.DescriptorProto_Name_field_number:
				md.L0.FullName = appendFullName(sb, pd.FullName(), v)
			case genid.DescriptorProto_EnumType_field_number:
				if prevField != genid.DescriptorProto_EnumType_field_number ***REMOVED***
					if numEnums > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posEnums = len(b0) - len(b) - n - m
				***REMOVED***
				numEnums++
			case genid.DescriptorProto_NestedType_field_number:
				if prevField != genid.DescriptorProto_NestedType_field_number ***REMOVED***
					if numMessages > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posMessages = len(b0) - len(b) - n - m
				***REMOVED***
				numMessages++
			case genid.DescriptorProto_Extension_field_number:
				if prevField != genid.DescriptorProto_Extension_field_number ***REMOVED***
					if numExtensions > 0 ***REMOVED***
						panic("non-contiguous repeated field")
					***REMOVED***
					posExtensions = len(b0) - len(b) - n - m
				***REMOVED***
				numExtensions++
			case genid.DescriptorProto_Options_field_number:
				md.unmarshalSeedOptions(v)
			***REMOVED***
			prevField = num
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
			prevField = -1 // ignore known field numbers of unknown wire type
		***REMOVED***
	***REMOVED***

	// Must allocate all declarations before parsing each descriptor type
	// to ensure we handled all descriptors in "flattened ordering".
	if numEnums > 0 ***REMOVED***
		md.L1.Enums.List = pf.allocEnums(numEnums)
	***REMOVED***
	if numMessages > 0 ***REMOVED***
		md.L1.Messages.List = pf.allocMessages(numMessages)
	***REMOVED***
	if numExtensions > 0 ***REMOVED***
		md.L1.Extensions.List = pf.allocExtensions(numExtensions)
	***REMOVED***

	if numEnums > 0 ***REMOVED***
		b := b0[posEnums:]
		for i := range md.L1.Enums.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			md.L1.Enums.List[i].unmarshalSeed(v, sb, pf, md, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
	if numMessages > 0 ***REMOVED***
		b := b0[posMessages:]
		for i := range md.L1.Messages.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			md.L1.Messages.List[i].unmarshalSeed(v, sb, pf, md, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
	if numExtensions > 0 ***REMOVED***
		b := b0[posExtensions:]
		for i := range md.L1.Extensions.List ***REMOVED***
			_, n := protowire.ConsumeVarint(b)
			v, m := protowire.ConsumeBytes(b[n:])
			md.L1.Extensions.List[i].unmarshalSeed(v, sb, pf, md, i)
			b = b[n+m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (md *Message) unmarshalSeedOptions(b []byte) ***REMOVED***
	for len(b) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.VarintType:
			v, m := protowire.ConsumeVarint(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.MessageOptions_MapEntry_field_number:
				md.L1.IsMapEntry = protowire.DecodeBool(v)
			case genid.MessageOptions_MessageSetWireFormat_field_number:
				md.L1.IsMessageSet = protowire.DecodeBool(v)
			***REMOVED***
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (xd *Extension) unmarshalSeed(b []byte, sb *strs.Builder, pf *File, pd pref.Descriptor, i int) ***REMOVED***
	xd.L0.ParentFile = pf
	xd.L0.Parent = pd
	xd.L0.Index = i

	for len(b) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.VarintType:
			v, m := protowire.ConsumeVarint(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.FieldDescriptorProto_Number_field_number:
				xd.L1.Number = pref.FieldNumber(v)
			case genid.FieldDescriptorProto_Label_field_number:
				xd.L1.Cardinality = pref.Cardinality(v)
			case genid.FieldDescriptorProto_Type_field_number:
				xd.L1.Kind = pref.Kind(v)
			***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.FieldDescriptorProto_Name_field_number:
				xd.L0.FullName = appendFullName(sb, pd.FullName(), v)
			case genid.FieldDescriptorProto_Extendee_field_number:
				xd.L1.Extendee = PlaceholderMessage(makeFullName(sb, v))
			***REMOVED***
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

func (sd *Service) unmarshalSeed(b []byte, sb *strs.Builder, pf *File, pd pref.Descriptor, i int) ***REMOVED***
	sd.L0.ParentFile = pf
	sd.L0.Parent = pd
	sd.L0.Index = i

	for len(b) > 0 ***REMOVED***
		num, typ, n := protowire.ConsumeTag(b)
		b = b[n:]
		switch typ ***REMOVED***
		case protowire.BytesType:
			v, m := protowire.ConsumeBytes(b)
			b = b[m:]
			switch num ***REMOVED***
			case genid.ServiceDescriptorProto_Name_field_number:
				sd.L0.FullName = appendFullName(sb, pd.FullName(), v)
			***REMOVED***
		default:
			m := protowire.ConsumeFieldValue(num, typ, b)
			b = b[m:]
		***REMOVED***
	***REMOVED***
***REMOVED***

var nameBuilderPool = sync.Pool***REMOVED***
	New: func() interface***REMOVED******REMOVED*** ***REMOVED*** return new(strs.Builder) ***REMOVED***,
***REMOVED***

func getBuilder() *strs.Builder ***REMOVED***
	return nameBuilderPool.Get().(*strs.Builder)
***REMOVED***
func putBuilder(b *strs.Builder) ***REMOVED***
	nameBuilderPool.Put(b)
***REMOVED***

// makeFullName converts b to a protoreflect.FullName,
// where b must start with a leading dot.
func makeFullName(sb *strs.Builder, b []byte) pref.FullName ***REMOVED***
	if len(b) == 0 || b[0] != '.' ***REMOVED***
		panic("name reference must be fully qualified")
	***REMOVED***
	return pref.FullName(sb.MakeString(b[1:]))
***REMOVED***

func appendFullName(sb *strs.Builder, prefix pref.FullName, suffix []byte) pref.FullName ***REMOVED***
	return sb.AppendFullName(prefix, pref.Name(strs.UnsafeString(suffix)))
***REMOVED***
