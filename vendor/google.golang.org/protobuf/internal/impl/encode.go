// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"math"
	"sort"
	"sync/atomic"

	"google.golang.org/protobuf/internal/flags"
	proto "google.golang.org/protobuf/proto"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

type marshalOptions struct ***REMOVED***
	flags piface.MarshalInputFlags
***REMOVED***

func (o marshalOptions) Options() proto.MarshalOptions ***REMOVED***
	return proto.MarshalOptions***REMOVED***
		AllowPartial:  true,
		Deterministic: o.Deterministic(),
		UseCachedSize: o.UseCachedSize(),
	***REMOVED***
***REMOVED***

func (o marshalOptions) Deterministic() bool ***REMOVED*** return o.flags&piface.MarshalDeterministic != 0 ***REMOVED***
func (o marshalOptions) UseCachedSize() bool ***REMOVED*** return o.flags&piface.MarshalUseCachedSize != 0 ***REMOVED***

// size is protoreflect.Methods.Size.
func (mi *MessageInfo) size(in piface.SizeInput) piface.SizeOutput ***REMOVED***
	var p pointer
	if ms, ok := in.Message.(*messageState); ok ***REMOVED***
		p = ms.pointer()
	***REMOVED*** else ***REMOVED***
		p = in.Message.(*messageReflectWrapper).pointer()
	***REMOVED***
	size := mi.sizePointer(p, marshalOptions***REMOVED***
		flags: in.Flags,
	***REMOVED***)
	return piface.SizeOutput***REMOVED***Size: size***REMOVED***
***REMOVED***

func (mi *MessageInfo) sizePointer(p pointer, opts marshalOptions) (size int) ***REMOVED***
	mi.init()
	if p.IsNil() ***REMOVED***
		return 0
	***REMOVED***
	if opts.UseCachedSize() && mi.sizecacheOffset.IsValid() ***REMOVED***
		if size := atomic.LoadInt32(p.Apply(mi.sizecacheOffset).Int32()); size >= 0 ***REMOVED***
			return int(size)
		***REMOVED***
	***REMOVED***
	return mi.sizePointerSlow(p, opts)
***REMOVED***

func (mi *MessageInfo) sizePointerSlow(p pointer, opts marshalOptions) (size int) ***REMOVED***
	if flags.ProtoLegacy && mi.isMessageSet ***REMOVED***
		size = sizeMessageSet(mi, p, opts)
		if mi.sizecacheOffset.IsValid() ***REMOVED***
			atomic.StoreInt32(p.Apply(mi.sizecacheOffset).Int32(), int32(size))
		***REMOVED***
		return size
	***REMOVED***
	if mi.extensionOffset.IsValid() ***REMOVED***
		e := p.Apply(mi.extensionOffset).Extensions()
		size += mi.sizeExtensions(e, opts)
	***REMOVED***
	for _, f := range mi.orderedCoderFields ***REMOVED***
		if f.funcs.size == nil ***REMOVED***
			continue
		***REMOVED***
		fptr := p.Apply(f.offset)
		if f.isPointer && fptr.Elem().IsNil() ***REMOVED***
			continue
		***REMOVED***
		size += f.funcs.size(fptr, f, opts)
	***REMOVED***
	if mi.unknownOffset.IsValid() ***REMOVED***
		if u := mi.getUnknownBytes(p); u != nil ***REMOVED***
			size += len(*u)
		***REMOVED***
	***REMOVED***
	if mi.sizecacheOffset.IsValid() ***REMOVED***
		if size > math.MaxInt32 ***REMOVED***
			// The size is too large for the int32 sizecache field.
			// We will need to recompute the size when encoding;
			// unfortunately expensive, but better than invalid output.
			atomic.StoreInt32(p.Apply(mi.sizecacheOffset).Int32(), -1)
		***REMOVED*** else ***REMOVED***
			atomic.StoreInt32(p.Apply(mi.sizecacheOffset).Int32(), int32(size))
		***REMOVED***
	***REMOVED***
	return size
***REMOVED***

// marshal is protoreflect.Methods.Marshal.
func (mi *MessageInfo) marshal(in piface.MarshalInput) (out piface.MarshalOutput, err error) ***REMOVED***
	var p pointer
	if ms, ok := in.Message.(*messageState); ok ***REMOVED***
		p = ms.pointer()
	***REMOVED*** else ***REMOVED***
		p = in.Message.(*messageReflectWrapper).pointer()
	***REMOVED***
	b, err := mi.marshalAppendPointer(in.Buf, p, marshalOptions***REMOVED***
		flags: in.Flags,
	***REMOVED***)
	return piface.MarshalOutput***REMOVED***Buf: b***REMOVED***, err
***REMOVED***

func (mi *MessageInfo) marshalAppendPointer(b []byte, p pointer, opts marshalOptions) ([]byte, error) ***REMOVED***
	mi.init()
	if p.IsNil() ***REMOVED***
		return b, nil
	***REMOVED***
	if flags.ProtoLegacy && mi.isMessageSet ***REMOVED***
		return marshalMessageSet(mi, b, p, opts)
	***REMOVED***
	var err error
	// The old marshaler encodes extensions at beginning.
	if mi.extensionOffset.IsValid() ***REMOVED***
		e := p.Apply(mi.extensionOffset).Extensions()
		// TODO: Special handling for MessageSet?
		b, err = mi.appendExtensions(b, e, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	for _, f := range mi.orderedCoderFields ***REMOVED***
		if f.funcs.marshal == nil ***REMOVED***
			continue
		***REMOVED***
		fptr := p.Apply(f.offset)
		if f.isPointer && fptr.Elem().IsNil() ***REMOVED***
			continue
		***REMOVED***
		b, err = f.funcs.marshal(b, fptr, f, opts)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***
	if mi.unknownOffset.IsValid() && !mi.isMessageSet ***REMOVED***
		if u := mi.getUnknownBytes(p); u != nil ***REMOVED***
			b = append(b, (*u)...)
		***REMOVED***
	***REMOVED***
	return b, nil
***REMOVED***

func (mi *MessageInfo) sizeExtensions(ext *map[int32]ExtensionField, opts marshalOptions) (n int) ***REMOVED***
	if ext == nil ***REMOVED***
		return 0
	***REMOVED***
	for _, x := range *ext ***REMOVED***
		xi := getExtensionFieldInfo(x.Type())
		if xi.funcs.size == nil ***REMOVED***
			continue
		***REMOVED***
		n += xi.funcs.size(x.Value(), xi.tagsize, opts)
	***REMOVED***
	return n
***REMOVED***

func (mi *MessageInfo) appendExtensions(b []byte, ext *map[int32]ExtensionField, opts marshalOptions) ([]byte, error) ***REMOVED***
	if ext == nil ***REMOVED***
		return b, nil
	***REMOVED***

	switch len(*ext) ***REMOVED***
	case 0:
		return b, nil
	case 1:
		// Fast-path for one extension: Don't bother sorting the keys.
		var err error
		for _, x := range *ext ***REMOVED***
			xi := getExtensionFieldInfo(x.Type())
			b, err = xi.funcs.marshal(b, x.Value(), xi.wiretag, opts)
		***REMOVED***
		return b, err
	default:
		// Sort the keys to provide a deterministic encoding.
		// Not sure this is required, but the old code does it.
		keys := make([]int, 0, len(*ext))
		for k := range *ext ***REMOVED***
			keys = append(keys, int(k))
		***REMOVED***
		sort.Ints(keys)
		var err error
		for _, k := range keys ***REMOVED***
			x := (*ext)[int32(k)]
			xi := getExtensionFieldInfo(x.Type())
			b, err = xi.funcs.marshal(b, x.Value(), xi.wiretag, opts)
			if err != nil ***REMOVED***
				return b, err
			***REMOVED***
		***REMOVED***
		return b, nil
	***REMOVED***
***REMOVED***
