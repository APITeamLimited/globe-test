// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"math/bits"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	preg "google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoiface"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

var errDecode = errors.New("cannot parse invalid wire-format data")

type unmarshalOptions struct ***REMOVED***
	flags    protoiface.UnmarshalInputFlags
	resolver interface ***REMOVED***
		FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error)
		FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error)
	***REMOVED***
***REMOVED***

func (o unmarshalOptions) Options() proto.UnmarshalOptions ***REMOVED***
	return proto.UnmarshalOptions***REMOVED***
		Merge:          true,
		AllowPartial:   true,
		DiscardUnknown: o.DiscardUnknown(),
		Resolver:       o.resolver,
	***REMOVED***
***REMOVED***

func (o unmarshalOptions) DiscardUnknown() bool ***REMOVED*** return o.flags&piface.UnmarshalDiscardUnknown != 0 ***REMOVED***

func (o unmarshalOptions) IsDefault() bool ***REMOVED***
	return o.flags == 0 && o.resolver == preg.GlobalTypes
***REMOVED***

var lazyUnmarshalOptions = unmarshalOptions***REMOVED***
	resolver: preg.GlobalTypes,
***REMOVED***

type unmarshalOutput struct ***REMOVED***
	n           int // number of bytes consumed
	initialized bool
***REMOVED***

// unmarshal is protoreflect.Methods.Unmarshal.
func (mi *MessageInfo) unmarshal(in piface.UnmarshalInput) (piface.UnmarshalOutput, error) ***REMOVED***
	var p pointer
	if ms, ok := in.Message.(*messageState); ok ***REMOVED***
		p = ms.pointer()
	***REMOVED*** else ***REMOVED***
		p = in.Message.(*messageReflectWrapper).pointer()
	***REMOVED***
	out, err := mi.unmarshalPointer(in.Buf, p, 0, unmarshalOptions***REMOVED***
		flags:    in.Flags,
		resolver: in.Resolver,
	***REMOVED***)
	var flags piface.UnmarshalOutputFlags
	if out.initialized ***REMOVED***
		flags |= piface.UnmarshalInitialized
	***REMOVED***
	return piface.UnmarshalOutput***REMOVED***
		Flags: flags,
	***REMOVED***, err
***REMOVED***

// errUnknown is returned during unmarshaling to indicate a parse error that
// should result in a field being placed in the unknown fields section (for example,
// when the wire type doesn't match) as opposed to the entire unmarshal operation
// failing (for example, when a field extends past the available input).
//
// This is a sentinel error which should never be visible to the user.
var errUnknown = errors.New("unknown")

func (mi *MessageInfo) unmarshalPointer(b []byte, p pointer, groupTag protowire.Number, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	mi.init()
	if flags.ProtoLegacy && mi.isMessageSet ***REMOVED***
		return unmarshalMessageSet(mi, b, p, opts)
	***REMOVED***
	initialized := true
	var requiredMask uint64
	var exts *map[int32]ExtensionField
	start := len(b)
	for len(b) > 0 ***REMOVED***
		// Parse the tag (field number and wire type).
		var tag uint64
		if b[0] < 0x80 ***REMOVED***
			tag = uint64(b[0])
			b = b[1:]
		***REMOVED*** else if len(b) >= 2 && b[1] < 128 ***REMOVED***
			tag = uint64(b[0]&0x7f) + uint64(b[1])<<7
			b = b[2:]
		***REMOVED*** else ***REMOVED***
			var n int
			tag, n = protowire.ConsumeVarint(b)
			if n < 0 ***REMOVED***
				return out, errDecode
			***REMOVED***
			b = b[n:]
		***REMOVED***
		var num protowire.Number
		if n := tag >> 3; n < uint64(protowire.MinValidNumber) || n > uint64(protowire.MaxValidNumber) ***REMOVED***
			return out, errDecode
		***REMOVED*** else ***REMOVED***
			num = protowire.Number(n)
		***REMOVED***
		wtyp := protowire.Type(tag & 7)

		if wtyp == protowire.EndGroupType ***REMOVED***
			if num != groupTag ***REMOVED***
				return out, errDecode
			***REMOVED***
			groupTag = 0
			break
		***REMOVED***

		var f *coderFieldInfo
		if int(num) < len(mi.denseCoderFields) ***REMOVED***
			f = mi.denseCoderFields[num]
		***REMOVED*** else ***REMOVED***
			f = mi.coderFields[num]
		***REMOVED***
		var n int
		err := errUnknown
		switch ***REMOVED***
		case f != nil:
			if f.funcs.unmarshal == nil ***REMOVED***
				break
			***REMOVED***
			var o unmarshalOutput
			o, err = f.funcs.unmarshal(b, p.Apply(f.offset), wtyp, f, opts)
			n = o.n
			if err != nil ***REMOVED***
				break
			***REMOVED***
			requiredMask |= f.validation.requiredBit
			if f.funcs.isInit != nil && !o.initialized ***REMOVED***
				initialized = false
			***REMOVED***
		default:
			// Possible extension.
			if exts == nil && mi.extensionOffset.IsValid() ***REMOVED***
				exts = p.Apply(mi.extensionOffset).Extensions()
				if *exts == nil ***REMOVED***
					*exts = make(map[int32]ExtensionField)
				***REMOVED***
			***REMOVED***
			if exts == nil ***REMOVED***
				break
			***REMOVED***
			var o unmarshalOutput
			o, err = mi.unmarshalExtension(b, num, wtyp, *exts, opts)
			if err != nil ***REMOVED***
				break
			***REMOVED***
			n = o.n
			if !o.initialized ***REMOVED***
				initialized = false
			***REMOVED***
		***REMOVED***
		if err != nil ***REMOVED***
			if err != errUnknown ***REMOVED***
				return out, err
			***REMOVED***
			n = protowire.ConsumeFieldValue(num, wtyp, b)
			if n < 0 ***REMOVED***
				return out, errDecode
			***REMOVED***
			if !opts.DiscardUnknown() && mi.unknownOffset.IsValid() ***REMOVED***
				u := mi.mutableUnknownBytes(p)
				*u = protowire.AppendTag(*u, num, wtyp)
				*u = append(*u, b[:n]...)
			***REMOVED***
		***REMOVED***
		b = b[n:]
	***REMOVED***
	if groupTag != 0 ***REMOVED***
		return out, errDecode
	***REMOVED***
	if mi.numRequiredFields > 0 && bits.OnesCount64(requiredMask) != int(mi.numRequiredFields) ***REMOVED***
		initialized = false
	***REMOVED***
	if initialized ***REMOVED***
		out.initialized = true
	***REMOVED***
	out.n = start - len(b)
	return out, nil
***REMOVED***

func (mi *MessageInfo) unmarshalExtension(b []byte, num protowire.Number, wtyp protowire.Type, exts map[int32]ExtensionField, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	x := exts[int32(num)]
	xt := x.Type()
	if xt == nil ***REMOVED***
		var err error
		xt, err = opts.resolver.FindExtensionByNumber(mi.Desc.FullName(), num)
		if err != nil ***REMOVED***
			if err == preg.NotFound ***REMOVED***
				return out, errUnknown
			***REMOVED***
			return out, errors.New("%v: unable to resolve extension %v: %v", mi.Desc.FullName(), num, err)
		***REMOVED***
	***REMOVED***
	xi := getExtensionFieldInfo(xt)
	if xi.funcs.unmarshal == nil ***REMOVED***
		return out, errUnknown
	***REMOVED***
	if flags.LazyUnmarshalExtensions ***REMOVED***
		if opts.IsDefault() && x.canLazy(xt) ***REMOVED***
			out, valid := skipExtension(b, xi, num, wtyp, opts)
			switch valid ***REMOVED***
			case ValidationValid:
				if out.initialized ***REMOVED***
					x.appendLazyBytes(xt, xi, num, wtyp, b[:out.n])
					exts[int32(num)] = x
					return out, nil
				***REMOVED***
			case ValidationInvalid:
				return out, errDecode
			case ValidationUnknown:
			***REMOVED***
		***REMOVED***
	***REMOVED***
	ival := x.Value()
	if !ival.IsValid() && xi.unmarshalNeedsValue ***REMOVED***
		// Create a new message, list, or map value to fill in.
		// For enums, create a prototype value to let the unmarshal func know the
		// concrete type.
		ival = xt.New()
	***REMOVED***
	v, out, err := xi.funcs.unmarshal(b, ival, num, wtyp, opts)
	if err != nil ***REMOVED***
		return out, err
	***REMOVED***
	if xi.funcs.isInit == nil ***REMOVED***
		out.initialized = true
	***REMOVED***
	x.Set(xt, v)
	exts[int32(num)] = x
	return out, nil
***REMOVED***

func skipExtension(b []byte, xi *extensionFieldInfo, num protowire.Number, wtyp protowire.Type, opts unmarshalOptions) (out unmarshalOutput, _ ValidationStatus) ***REMOVED***
	if xi.validation.mi == nil ***REMOVED***
		return out, ValidationUnknown
	***REMOVED***
	xi.validation.mi.init()
	switch xi.validation.typ ***REMOVED***
	case validationTypeMessage:
		if wtyp != protowire.BytesType ***REMOVED***
			return out, ValidationUnknown
		***REMOVED***
		v, n := protowire.ConsumeBytes(b)
		if n < 0 ***REMOVED***
			return out, ValidationUnknown
		***REMOVED***
		out, st := xi.validation.mi.validate(v, 0, opts)
		out.n = n
		return out, st
	case validationTypeGroup:
		if wtyp != protowire.StartGroupType ***REMOVED***
			return out, ValidationUnknown
		***REMOVED***
		out, st := xi.validation.mi.validate(b, num, opts)
		return out, st
	default:
		return out, ValidationUnknown
	***REMOVED***
***REMOVED***
