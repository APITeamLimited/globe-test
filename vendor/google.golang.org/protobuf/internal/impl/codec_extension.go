// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/internal/errors"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

type extensionFieldInfo struct ***REMOVED***
	wiretag             uint64
	tagsize             int
	unmarshalNeedsValue bool
	funcs               valueCoderFuncs
	validation          validationInfo
***REMOVED***

var legacyExtensionFieldInfoCache sync.Map // map[protoreflect.ExtensionType]*extensionFieldInfo

func getExtensionFieldInfo(xt pref.ExtensionType) *extensionFieldInfo ***REMOVED***
	if xi, ok := xt.(*ExtensionInfo); ok ***REMOVED***
		xi.lazyInit()
		return xi.info
	***REMOVED***
	return legacyLoadExtensionFieldInfo(xt)
***REMOVED***

// legacyLoadExtensionFieldInfo dynamically loads a *ExtensionInfo for xt.
func legacyLoadExtensionFieldInfo(xt pref.ExtensionType) *extensionFieldInfo ***REMOVED***
	if xi, ok := legacyExtensionFieldInfoCache.Load(xt); ok ***REMOVED***
		return xi.(*extensionFieldInfo)
	***REMOVED***
	e := makeExtensionFieldInfo(xt.TypeDescriptor())
	if e, ok := legacyMessageTypeCache.LoadOrStore(xt, e); ok ***REMOVED***
		return e.(*extensionFieldInfo)
	***REMOVED***
	return e
***REMOVED***

func makeExtensionFieldInfo(xd pref.ExtensionDescriptor) *extensionFieldInfo ***REMOVED***
	var wiretag uint64
	if !xd.IsPacked() ***REMOVED***
		wiretag = protowire.EncodeTag(xd.Number(), wireTypes[xd.Kind()])
	***REMOVED*** else ***REMOVED***
		wiretag = protowire.EncodeTag(xd.Number(), protowire.BytesType)
	***REMOVED***
	e := &extensionFieldInfo***REMOVED***
		wiretag: wiretag,
		tagsize: protowire.SizeVarint(wiretag),
		funcs:   encoderFuncsForValue(xd),
	***REMOVED***
	// Does the unmarshal function need a value passed to it?
	// This is true for composite types, where we pass in a message, list, or map to fill in,
	// and for enums, where we pass in a prototype value to specify the concrete enum type.
	switch xd.Kind() ***REMOVED***
	case pref.MessageKind, pref.GroupKind, pref.EnumKind:
		e.unmarshalNeedsValue = true
	default:
		if xd.Cardinality() == pref.Repeated ***REMOVED***
			e.unmarshalNeedsValue = true
		***REMOVED***
	***REMOVED***
	return e
***REMOVED***

type lazyExtensionValue struct ***REMOVED***
	atomicOnce uint32 // atomically set if value is valid
	mu         sync.Mutex
	xi         *extensionFieldInfo
	value      pref.Value
	b          []byte
	fn         func() pref.Value
***REMOVED***

type ExtensionField struct ***REMOVED***
	typ pref.ExtensionType

	// value is either the value of GetValue,
	// or a *lazyExtensionValue that then returns the value of GetValue.
	value pref.Value
	lazy  *lazyExtensionValue
***REMOVED***

func (f *ExtensionField) appendLazyBytes(xt pref.ExtensionType, xi *extensionFieldInfo, num protowire.Number, wtyp protowire.Type, b []byte) ***REMOVED***
	if f.lazy == nil ***REMOVED***
		f.lazy = &lazyExtensionValue***REMOVED***xi: xi***REMOVED***
	***REMOVED***
	f.typ = xt
	f.lazy.xi = xi
	f.lazy.b = protowire.AppendTag(f.lazy.b, num, wtyp)
	f.lazy.b = append(f.lazy.b, b...)
***REMOVED***

func (f *ExtensionField) canLazy(xt pref.ExtensionType) bool ***REMOVED***
	if f.typ == nil ***REMOVED***
		return true
	***REMOVED***
	if f.typ == xt && f.lazy != nil && atomic.LoadUint32(&f.lazy.atomicOnce) == 0 ***REMOVED***
		return true
	***REMOVED***
	return false
***REMOVED***

func (f *ExtensionField) lazyInit() ***REMOVED***
	f.lazy.mu.Lock()
	defer f.lazy.mu.Unlock()
	if atomic.LoadUint32(&f.lazy.atomicOnce) == 1 ***REMOVED***
		return
	***REMOVED***
	if f.lazy.xi != nil ***REMOVED***
		b := f.lazy.b
		val := f.typ.New()
		for len(b) > 0 ***REMOVED***
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
					panic(errors.New("bad tag in lazy extension decoding"))
				***REMOVED***
				b = b[n:]
			***REMOVED***
			num := protowire.Number(tag >> 3)
			wtyp := protowire.Type(tag & 7)
			var out unmarshalOutput
			var err error
			val, out, err = f.lazy.xi.funcs.unmarshal(b, val, num, wtyp, lazyUnmarshalOptions)
			if err != nil ***REMOVED***
				panic(errors.New("decode failure in lazy extension decoding: %v", err))
			***REMOVED***
			b = b[out.n:]
		***REMOVED***
		f.lazy.value = val
	***REMOVED*** else ***REMOVED***
		f.lazy.value = f.lazy.fn()
	***REMOVED***
	f.lazy.xi = nil
	f.lazy.fn = nil
	f.lazy.b = nil
	atomic.StoreUint32(&f.lazy.atomicOnce, 1)
***REMOVED***

// Set sets the type and value of the extension field.
// This must not be called concurrently.
func (f *ExtensionField) Set(t pref.ExtensionType, v pref.Value) ***REMOVED***
	f.typ = t
	f.value = v
	f.lazy = nil
***REMOVED***

// SetLazy sets the type and a value that is to be lazily evaluated upon first use.
// This must not be called concurrently.
func (f *ExtensionField) SetLazy(t pref.ExtensionType, fn func() pref.Value) ***REMOVED***
	f.typ = t
	f.lazy = &lazyExtensionValue***REMOVED***fn: fn***REMOVED***
***REMOVED***

// Value returns the value of the extension field.
// This may be called concurrently.
func (f *ExtensionField) Value() pref.Value ***REMOVED***
	if f.lazy != nil ***REMOVED***
		if atomic.LoadUint32(&f.lazy.atomicOnce) == 0 ***REMOVED***
			f.lazyInit()
		***REMOVED***
		return f.lazy.value
	***REMOVED***
	return f.value
***REMOVED***

// Type returns the type of the extension field.
// This may be called concurrently.
func (f ExtensionField) Type() pref.ExtensionType ***REMOVED***
	return f.typ
***REMOVED***

// IsSet returns whether the extension field is set.
// This may be called concurrently.
func (f ExtensionField) IsSet() bool ***REMOVED***
	return f.typ != nil
***REMOVED***

// IsLazy reports whether a field is lazily encoded.
// It is exported for testing.
func IsLazy(m pref.Message, fd pref.FieldDescriptor) bool ***REMOVED***
	var mi *MessageInfo
	var p pointer
	switch m := m.(type) ***REMOVED***
	case *messageState:
		mi = m.messageInfo()
		p = m.pointer()
	case *messageReflectWrapper:
		mi = m.messageInfo()
		p = m.pointer()
	default:
		return false
	***REMOVED***
	xd, ok := fd.(pref.ExtensionTypeDescriptor)
	if !ok ***REMOVED***
		return false
	***REMOVED***
	xt := xd.Type()
	ext := mi.extensionMap(p)
	if ext == nil ***REMOVED***
		return false
	***REMOVED***
	f, ok := (*ext)[int32(fd.Number())]
	if !ok ***REMOVED***
		return false
	***REMOVED***
	return f.typ == xt && f.lazy != nil && atomic.LoadUint32(&f.lazy.atomicOnce) == 0
***REMOVED***
