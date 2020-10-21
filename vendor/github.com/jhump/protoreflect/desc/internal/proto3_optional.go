package internal

import (
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/internal/codec"
	"reflect"
	"strings"

	"github.com/jhump/protoreflect/internal"
)

// NB: We use reflection or unknown fields in case we are linked against an older
// version of the proto runtime which does not know about the proto3_optional field.
// We don't require linking with newer version (which would greatly simplify this)
// because that means pulling in v1.4+ of the protobuf runtime, which has some
// compatibility issues. (We'll be nice to users and not require they upgrade to
// that latest runtime to upgrade to newer protoreflect.)

func GetProto3Optional(fd *dpb.FieldDescriptorProto) bool ***REMOVED***
	type newerFieldDesc interface ***REMOVED***
		GetProto3Optional() bool
	***REMOVED***
	var pm proto.Message = fd
	if fd, ok := pm.(newerFieldDesc); ok ***REMOVED***
		return fd.GetProto3Optional()
	***REMOVED***

	// Field does not exist, so we have to examine unknown fields
	// (we just silently bail if we have problems parsing them)
	unk := internal.GetUnrecognized(pm)
	buf := codec.NewBuffer(unk)
	for ***REMOVED***
		tag, wt, err := buf.DecodeTagAndWireType()
		if err != nil ***REMOVED***
			return false
		***REMOVED***
		if tag == Field_proto3OptionalTag && wt == proto.WireVarint ***REMOVED***
			v, _ := buf.DecodeVarint()
			return v != 0
		***REMOVED***
		if err := buf.SkipField(wt); err != nil ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
***REMOVED***

func SetProto3Optional(fd *dpb.FieldDescriptorProto) ***REMOVED***
	rv := reflect.ValueOf(fd).Elem()
	fld := rv.FieldByName("Proto3Optional")
	if fld.IsValid() ***REMOVED***
		fld.Set(reflect.ValueOf(proto.Bool(true)))
		return
	***REMOVED***

	// Field does not exist, so we have to store as unknown field.
	var buf codec.Buffer
	if err := buf.EncodeTagAndWireType(Field_proto3OptionalTag, proto.WireVarint); err != nil ***REMOVED***
		// TODO: panic? log?
		return
	***REMOVED***
	if err := buf.EncodeVarint(1); err != nil ***REMOVED***
		// TODO: panic? log?
		return
	***REMOVED***
	internal.SetUnrecognized(fd, buf.Bytes())
***REMOVED***

// ProcessProto3OptionalFields adds synthetic oneofs to the given message descriptor
// for each proto3 optional field. It also updates the fields to have the correct
// oneof index reference.
func ProcessProto3OptionalFields(msgd *dpb.DescriptorProto) ***REMOVED***
	var allNames map[string]struct***REMOVED******REMOVED***
	for _, fd := range msgd.Field ***REMOVED***
		if GetProto3Optional(fd) ***REMOVED***
			// lazy init the set of all names
			if allNames == nil ***REMOVED***
				allNames = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
				for _, fd := range msgd.Field ***REMOVED***
					allNames[fd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
				for _, fd := range msgd.Extension ***REMOVED***
					allNames[fd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
				for _, ed := range msgd.EnumType ***REMOVED***
					allNames[ed.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					for _, evd := range ed.Value ***REMOVED***
						allNames[evd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					***REMOVED***
				***REMOVED***
				for _, fd := range msgd.NestedType ***REMOVED***
					allNames[fd.GetName()] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
				for _, n := range msgd.ReservedName ***REMOVED***
					allNames[n] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***

			// Compute a name for the synthetic oneof. This uses the same
			// algorithm as used in protoc:
			//  https://github.com/protocolbuffers/protobuf/blob/74ad62759e0a9b5a21094f3fb9bb4ebfaa0d1ab8/src/google/protobuf/compiler/parser.cc#L785-L803
			ooName := fd.GetName()
			if !strings.HasPrefix(ooName, "_") ***REMOVED***
				ooName = "_" + ooName
			***REMOVED***
			for ***REMOVED***
				_, ok := allNames[ooName]
				if !ok ***REMOVED***
					// found a unique name
					allNames[ooName] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
					break
				***REMOVED***
				ooName = "X" + ooName
			***REMOVED***

			fd.OneofIndex = proto.Int32(int32(len(msgd.OneofDecl)))
			msgd.OneofDecl = append(msgd.OneofDecl, &dpb.OneofDescriptorProto***REMOVED***Name: proto.String(ooName)***REMOVED***)
		***REMOVED***
	***REMOVED***
***REMOVED***
