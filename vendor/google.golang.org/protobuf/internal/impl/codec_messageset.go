// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"sort"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/encoding/messageset"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/flags"
)

func sizeMessageSet(mi *MessageInfo, p pointer, opts marshalOptions) (size int) ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		return 0
	***REMOVED***

	ext := *p.Apply(mi.extensionOffset).Extensions()
	for _, x := range ext ***REMOVED***
		xi := getExtensionFieldInfo(x.Type())
		if xi.funcs.size == nil ***REMOVED***
			continue
		***REMOVED***
		num, _ := protowire.DecodeTag(xi.wiretag)
		size += messageset.SizeField(num)
		size += xi.funcs.size(x.Value(), protowire.SizeTag(messageset.FieldMessage), opts)
	***REMOVED***

	if u := mi.getUnknownBytes(p); u != nil ***REMOVED***
		size += messageset.SizeUnknown(*u)
	***REMOVED***

	return size
***REMOVED***

func marshalMessageSet(mi *MessageInfo, b []byte, p pointer, opts marshalOptions) ([]byte, error) ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		return b, errors.New("no support for message_set_wire_format")
	***REMOVED***

	ext := *p.Apply(mi.extensionOffset).Extensions()
	switch len(ext) ***REMOVED***
	case 0:
	case 1:
		// Fast-path for one extension: Don't bother sorting the keys.
		for _, x := range ext ***REMOVED***
			var err error
			b, err = marshalMessageSetField(mi, b, x, opts)
			if err != nil ***REMOVED***
				return b, err
			***REMOVED***
		***REMOVED***
	default:
		// Sort the keys to provide a deterministic encoding.
		// Not sure this is required, but the old code does it.
		keys := make([]int, 0, len(ext))
		for k := range ext ***REMOVED***
			keys = append(keys, int(k))
		***REMOVED***
		sort.Ints(keys)
		for _, k := range keys ***REMOVED***
			var err error
			b, err = marshalMessageSetField(mi, b, ext[int32(k)], opts)
			if err != nil ***REMOVED***
				return b, err
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if u := mi.getUnknownBytes(p); u != nil ***REMOVED***
		var err error
		b, err = messageset.AppendUnknown(b, *u)
		if err != nil ***REMOVED***
			return b, err
		***REMOVED***
	***REMOVED***

	return b, nil
***REMOVED***

func marshalMessageSetField(mi *MessageInfo, b []byte, x ExtensionField, opts marshalOptions) ([]byte, error) ***REMOVED***
	xi := getExtensionFieldInfo(x.Type())
	num, _ := protowire.DecodeTag(xi.wiretag)
	b = messageset.AppendFieldStart(b, num)
	b, err := xi.funcs.marshal(b, x.Value(), protowire.EncodeTag(messageset.FieldMessage, protowire.BytesType), opts)
	if err != nil ***REMOVED***
		return b, err
	***REMOVED***
	b = messageset.AppendFieldEnd(b)
	return b, nil
***REMOVED***

func unmarshalMessageSet(mi *MessageInfo, b []byte, p pointer, opts unmarshalOptions) (out unmarshalOutput, err error) ***REMOVED***
	if !flags.ProtoLegacy ***REMOVED***
		return out, errors.New("no support for message_set_wire_format")
	***REMOVED***

	ep := p.Apply(mi.extensionOffset).Extensions()
	if *ep == nil ***REMOVED***
		*ep = make(map[int32]ExtensionField)
	***REMOVED***
	ext := *ep
	initialized := true
	err = messageset.Unmarshal(b, true, func(num protowire.Number, v []byte) error ***REMOVED***
		o, err := mi.unmarshalExtension(v, num, protowire.BytesType, ext, opts)
		if err == errUnknown ***REMOVED***
			u := mi.mutableUnknownBytes(p)
			*u = protowire.AppendTag(*u, num, protowire.BytesType)
			*u = append(*u, v...)
			return nil
		***REMOVED***
		if !o.initialized ***REMOVED***
			initialized = false
		***REMOVED***
		return err
	***REMOVED***)
	out.n = len(b)
	out.initialized = initialized
	return out, err
***REMOVED***
