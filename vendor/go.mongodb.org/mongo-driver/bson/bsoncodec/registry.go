// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsoncodec

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"go.mongodb.org/mongo-driver/bson/bsontype"
)

// ErrNilType is returned when nil is passed to either LookupEncoder or LookupDecoder.
var ErrNilType = errors.New("cannot perform a decoder lookup on <nil>")

// ErrNotPointer is returned when a non-pointer type is provided to LookupDecoder.
var ErrNotPointer = errors.New("non-pointer provided to LookupDecoder")

// ErrNoEncoder is returned when there wasn't an encoder available for a type.
type ErrNoEncoder struct ***REMOVED***
	Type reflect.Type
***REMOVED***

func (ene ErrNoEncoder) Error() string ***REMOVED***
	if ene.Type == nil ***REMOVED***
		return "no encoder found for <nil>"
	***REMOVED***
	return "no encoder found for " + ene.Type.String()
***REMOVED***

// ErrNoDecoder is returned when there wasn't a decoder available for a type.
type ErrNoDecoder struct ***REMOVED***
	Type reflect.Type
***REMOVED***

func (end ErrNoDecoder) Error() string ***REMOVED***
	return "no decoder found for " + end.Type.String()
***REMOVED***

// ErrNoTypeMapEntry is returned when there wasn't a type available for the provided BSON type.
type ErrNoTypeMapEntry struct ***REMOVED***
	Type bsontype.Type
***REMOVED***

func (entme ErrNoTypeMapEntry) Error() string ***REMOVED***
	return "no type map entry found for " + entme.Type.String()
***REMOVED***

// ErrNotInterface is returned when the provided type is not an interface.
var ErrNotInterface = errors.New("The provided type is not an interface")

// A RegistryBuilder is used to build a Registry. This type is not goroutine
// safe.
type RegistryBuilder struct ***REMOVED***
	typeEncoders      map[reflect.Type]ValueEncoder
	interfaceEncoders []interfaceValueEncoder
	kindEncoders      map[reflect.Kind]ValueEncoder

	typeDecoders      map[reflect.Type]ValueDecoder
	interfaceDecoders []interfaceValueDecoder
	kindDecoders      map[reflect.Kind]ValueDecoder

	typeMap map[bsontype.Type]reflect.Type
***REMOVED***

// A Registry is used to store and retrieve codecs for types and interfaces. This type is the main
// typed passed around and Encoders and Decoders are constructed from it.
type Registry struct ***REMOVED***
	typeEncoders map[reflect.Type]ValueEncoder
	typeDecoders map[reflect.Type]ValueDecoder

	interfaceEncoders []interfaceValueEncoder
	interfaceDecoders []interfaceValueDecoder

	kindEncoders map[reflect.Kind]ValueEncoder
	kindDecoders map[reflect.Kind]ValueDecoder

	typeMap map[bsontype.Type]reflect.Type

	mu sync.RWMutex
***REMOVED***

// NewRegistryBuilder creates a new empty RegistryBuilder.
func NewRegistryBuilder() *RegistryBuilder ***REMOVED***
	return &RegistryBuilder***REMOVED***
		typeEncoders: make(map[reflect.Type]ValueEncoder),
		typeDecoders: make(map[reflect.Type]ValueDecoder),

		interfaceEncoders: make([]interfaceValueEncoder, 0),
		interfaceDecoders: make([]interfaceValueDecoder, 0),

		kindEncoders: make(map[reflect.Kind]ValueEncoder),
		kindDecoders: make(map[reflect.Kind]ValueDecoder),

		typeMap: make(map[bsontype.Type]reflect.Type),
	***REMOVED***
***REMOVED***

func buildDefaultRegistry() *Registry ***REMOVED***
	rb := NewRegistryBuilder()
	defaultValueEncoders.RegisterDefaultEncoders(rb)
	defaultValueDecoders.RegisterDefaultDecoders(rb)
	return rb.Build()
***REMOVED***

// RegisterCodec will register the provided ValueCodec for the provided type.
func (rb *RegistryBuilder) RegisterCodec(t reflect.Type, codec ValueCodec) *RegistryBuilder ***REMOVED***
	rb.RegisterTypeEncoder(t, codec)
	rb.RegisterTypeDecoder(t, codec)
	return rb
***REMOVED***

// RegisterTypeEncoder will register the provided ValueEncoder for the provided type.
//
// The type will be used directly, so an encoder can be registered for a type and a different encoder can be registered
// for a pointer to that type.
//
// If the given type is an interface, the encoder will be called when marshalling a type that is that interface. It
// will not be called when marshalling a non-interface type that implements the interface.
func (rb *RegistryBuilder) RegisterTypeEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder ***REMOVED***
	rb.typeEncoders[t] = enc
	return rb
***REMOVED***

// RegisterHookEncoder will register an encoder for the provided interface type t. This encoder will be called when
// marshalling a type if the type implements t or a pointer to the type implements t. If the provided type is not
// an interface (i.e. t.Kind() != reflect.Interface), this method will panic.
func (rb *RegistryBuilder) RegisterHookEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder ***REMOVED***
	if t.Kind() != reflect.Interface ***REMOVED***
		panicStr := fmt.Sprintf("RegisterHookEncoder expects a type with kind reflect.Interface, "+
			"got type %s with kind %s", t, t.Kind())
		panic(panicStr)
	***REMOVED***

	for idx, encoder := range rb.interfaceEncoders ***REMOVED***
		if encoder.i == t ***REMOVED***
			rb.interfaceEncoders[idx].ve = enc
			return rb
		***REMOVED***
	***REMOVED***

	rb.interfaceEncoders = append(rb.interfaceEncoders, interfaceValueEncoder***REMOVED***i: t, ve: enc***REMOVED***)
	return rb
***REMOVED***

// RegisterTypeDecoder will register the provided ValueDecoder for the provided type.
//
// The type will be used directly, so a decoder can be registered for a type and a different decoder can be registered
// for a pointer to that type.
//
// If the given type is an interface, the decoder will be called when unmarshalling into a type that is that interface.
// It will not be called when unmarshalling into a non-interface type that implements the interface.
func (rb *RegistryBuilder) RegisterTypeDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder ***REMOVED***
	rb.typeDecoders[t] = dec
	return rb
***REMOVED***

// RegisterHookDecoder will register an decoder for the provided interface type t. This decoder will be called when
// unmarshalling into a type if the type implements t or a pointer to the type implements t. If the provided type is not
// an interface (i.e. t.Kind() != reflect.Interface), this method will panic.
func (rb *RegistryBuilder) RegisterHookDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder ***REMOVED***
	if t.Kind() != reflect.Interface ***REMOVED***
		panicStr := fmt.Sprintf("RegisterHookDecoder expects a type with kind reflect.Interface, "+
			"got type %s with kind %s", t, t.Kind())
		panic(panicStr)
	***REMOVED***

	for idx, decoder := range rb.interfaceDecoders ***REMOVED***
		if decoder.i == t ***REMOVED***
			rb.interfaceDecoders[idx].vd = dec
			return rb
		***REMOVED***
	***REMOVED***

	rb.interfaceDecoders = append(rb.interfaceDecoders, interfaceValueDecoder***REMOVED***i: t, vd: dec***REMOVED***)
	return rb
***REMOVED***

// RegisterEncoder registers the provided type and encoder pair.
//
// Deprecated: Use RegisterTypeEncoder or RegisterHookEncoder instead.
func (rb *RegistryBuilder) RegisterEncoder(t reflect.Type, enc ValueEncoder) *RegistryBuilder ***REMOVED***
	if t == tEmpty ***REMOVED***
		rb.typeEncoders[t] = enc
		return rb
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Interface:
		for idx, ir := range rb.interfaceEncoders ***REMOVED***
			if ir.i == t ***REMOVED***
				rb.interfaceEncoders[idx].ve = enc
				return rb
			***REMOVED***
		***REMOVED***

		rb.interfaceEncoders = append(rb.interfaceEncoders, interfaceValueEncoder***REMOVED***i: t, ve: enc***REMOVED***)
	default:
		rb.typeEncoders[t] = enc
	***REMOVED***
	return rb
***REMOVED***

// RegisterDecoder registers the provided type and decoder pair.
//
// Deprecated: Use RegisterTypeDecoder or RegisterHookDecoder instead.
func (rb *RegistryBuilder) RegisterDecoder(t reflect.Type, dec ValueDecoder) *RegistryBuilder ***REMOVED***
	if t == nil ***REMOVED***
		rb.typeDecoders[nil] = dec
		return rb
	***REMOVED***
	if t == tEmpty ***REMOVED***
		rb.typeDecoders[t] = dec
		return rb
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Interface:
		for idx, ir := range rb.interfaceDecoders ***REMOVED***
			if ir.i == t ***REMOVED***
				rb.interfaceDecoders[idx].vd = dec
				return rb
			***REMOVED***
		***REMOVED***

		rb.interfaceDecoders = append(rb.interfaceDecoders, interfaceValueDecoder***REMOVED***i: t, vd: dec***REMOVED***)
	default:
		rb.typeDecoders[t] = dec
	***REMOVED***
	return rb
***REMOVED***

// RegisterDefaultEncoder will registr the provided ValueEncoder to the provided
// kind.
func (rb *RegistryBuilder) RegisterDefaultEncoder(kind reflect.Kind, enc ValueEncoder) *RegistryBuilder ***REMOVED***
	rb.kindEncoders[kind] = enc
	return rb
***REMOVED***

// RegisterDefaultDecoder will register the provided ValueDecoder to the
// provided kind.
func (rb *RegistryBuilder) RegisterDefaultDecoder(kind reflect.Kind, dec ValueDecoder) *RegistryBuilder ***REMOVED***
	rb.kindDecoders[kind] = dec
	return rb
***REMOVED***

// RegisterTypeMapEntry will register the provided type to the BSON type. The primary usage for this
// mapping is decoding situations where an empty interface is used and a default type needs to be
// created and decoded into.
//
// By default, BSON documents will decode into interface***REMOVED******REMOVED*** values as bson.D. To change the default type for BSON
// documents, a type map entry for bsontype.EmbeddedDocument should be registered. For example, to force BSON documents
// to decode to bson.Raw, use the following code:
//	rb.RegisterTypeMapEntry(bsontype.EmbeddedDocument, reflect.TypeOf(bson.Raw***REMOVED******REMOVED***))
func (rb *RegistryBuilder) RegisterTypeMapEntry(bt bsontype.Type, rt reflect.Type) *RegistryBuilder ***REMOVED***
	rb.typeMap[bt] = rt
	return rb
***REMOVED***

// Build creates a Registry from the current state of this RegistryBuilder.
func (rb *RegistryBuilder) Build() *Registry ***REMOVED***
	registry := new(Registry)

	registry.typeEncoders = make(map[reflect.Type]ValueEncoder)
	for t, enc := range rb.typeEncoders ***REMOVED***
		registry.typeEncoders[t] = enc
	***REMOVED***

	registry.typeDecoders = make(map[reflect.Type]ValueDecoder)
	for t, dec := range rb.typeDecoders ***REMOVED***
		registry.typeDecoders[t] = dec
	***REMOVED***

	registry.interfaceEncoders = make([]interfaceValueEncoder, len(rb.interfaceEncoders))
	copy(registry.interfaceEncoders, rb.interfaceEncoders)

	registry.interfaceDecoders = make([]interfaceValueDecoder, len(rb.interfaceDecoders))
	copy(registry.interfaceDecoders, rb.interfaceDecoders)

	registry.kindEncoders = make(map[reflect.Kind]ValueEncoder)
	for kind, enc := range rb.kindEncoders ***REMOVED***
		registry.kindEncoders[kind] = enc
	***REMOVED***

	registry.kindDecoders = make(map[reflect.Kind]ValueDecoder)
	for kind, dec := range rb.kindDecoders ***REMOVED***
		registry.kindDecoders[kind] = dec
	***REMOVED***

	registry.typeMap = make(map[bsontype.Type]reflect.Type)
	for bt, rt := range rb.typeMap ***REMOVED***
		registry.typeMap[bt] = rt
	***REMOVED***

	return registry
***REMOVED***

// LookupEncoder inspects the registry for an encoder for the given type. The lookup precedence works as follows:
//
// 1. An encoder registered for the exact type. If the given type represents an interface, an encoder registered using
// RegisterTypeEncoder for the interface will be selected.
//
// 2. An encoder registered using RegisterHookEncoder for an interface implemented by the type or by a pointer to the
// type.
//
// 3. An encoder registered for the reflect.Kind of the value.
//
// If no encoder is found, an error of type ErrNoEncoder is returned.
func (r *Registry) LookupEncoder(t reflect.Type) (ValueEncoder, error) ***REMOVED***
	encodererr := ErrNoEncoder***REMOVED***Type: t***REMOVED***
	r.mu.RLock()
	enc, found := r.lookupTypeEncoder(t)
	r.mu.RUnlock()
	if found ***REMOVED***
		if enc == nil ***REMOVED***
			return nil, ErrNoEncoder***REMOVED***Type: t***REMOVED***
		***REMOVED***
		return enc, nil
	***REMOVED***

	enc, found = r.lookupInterfaceEncoder(t, true)
	if found ***REMOVED***
		r.mu.Lock()
		r.typeEncoders[t] = enc
		r.mu.Unlock()
		return enc, nil
	***REMOVED***

	if t == nil ***REMOVED***
		r.mu.Lock()
		r.typeEncoders[t] = nil
		r.mu.Unlock()
		return nil, encodererr
	***REMOVED***

	enc, found = r.kindEncoders[t.Kind()]
	if !found ***REMOVED***
		r.mu.Lock()
		r.typeEncoders[t] = nil
		r.mu.Unlock()
		return nil, encodererr
	***REMOVED***

	r.mu.Lock()
	r.typeEncoders[t] = enc
	r.mu.Unlock()
	return enc, nil
***REMOVED***

func (r *Registry) lookupTypeEncoder(t reflect.Type) (ValueEncoder, bool) ***REMOVED***
	enc, found := r.typeEncoders[t]
	return enc, found
***REMOVED***

func (r *Registry) lookupInterfaceEncoder(t reflect.Type, allowAddr bool) (ValueEncoder, bool) ***REMOVED***
	if t == nil ***REMOVED***
		return nil, false
	***REMOVED***
	for _, ienc := range r.interfaceEncoders ***REMOVED***
		if t.Implements(ienc.i) ***REMOVED***
			return ienc.ve, true
		***REMOVED***
		if allowAddr && t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(ienc.i) ***REMOVED***
			// if *t implements an interface, this will catch if t implements an interface further ahead
			// in interfaceEncoders
			defaultEnc, found := r.lookupInterfaceEncoder(t, false)
			if !found ***REMOVED***
				defaultEnc = r.kindEncoders[t.Kind()]
			***REMOVED***
			return newCondAddrEncoder(ienc.ve, defaultEnc), true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// LookupDecoder inspects the registry for an decoder for the given type. The lookup precedence works as follows:
//
// 1. A decoder registered for the exact type. If the given type represents an interface, a decoder registered using
// RegisterTypeDecoder for the interface will be selected.
//
// 2. A decoder registered using RegisterHookDecoder for an interface implemented by the type or by a pointer to the
// type.
//
// 3. A decoder registered for the reflect.Kind of the value.
//
// If no decoder is found, an error of type ErrNoDecoder is returned.
func (r *Registry) LookupDecoder(t reflect.Type) (ValueDecoder, error) ***REMOVED***
	if t == nil ***REMOVED***
		return nil, ErrNilType
	***REMOVED***
	decodererr := ErrNoDecoder***REMOVED***Type: t***REMOVED***
	r.mu.RLock()
	dec, found := r.lookupTypeDecoder(t)
	r.mu.RUnlock()
	if found ***REMOVED***
		if dec == nil ***REMOVED***
			return nil, ErrNoDecoder***REMOVED***Type: t***REMOVED***
		***REMOVED***
		return dec, nil
	***REMOVED***

	dec, found = r.lookupInterfaceDecoder(t, true)
	if found ***REMOVED***
		r.mu.Lock()
		r.typeDecoders[t] = dec
		r.mu.Unlock()
		return dec, nil
	***REMOVED***

	dec, found = r.kindDecoders[t.Kind()]
	if !found ***REMOVED***
		r.mu.Lock()
		r.typeDecoders[t] = nil
		r.mu.Unlock()
		return nil, decodererr
	***REMOVED***

	r.mu.Lock()
	r.typeDecoders[t] = dec
	r.mu.Unlock()
	return dec, nil
***REMOVED***

func (r *Registry) lookupTypeDecoder(t reflect.Type) (ValueDecoder, bool) ***REMOVED***
	dec, found := r.typeDecoders[t]
	return dec, found
***REMOVED***

func (r *Registry) lookupInterfaceDecoder(t reflect.Type, allowAddr bool) (ValueDecoder, bool) ***REMOVED***
	for _, idec := range r.interfaceDecoders ***REMOVED***
		if t.Implements(idec.i) ***REMOVED***
			return idec.vd, true
		***REMOVED***
		if allowAddr && t.Kind() != reflect.Ptr && reflect.PtrTo(t).Implements(idec.i) ***REMOVED***
			// if *t implements an interface, this will catch if t implements an interface further ahead
			// in interfaceDecoders
			defaultDec, found := r.lookupInterfaceDecoder(t, false)
			if !found ***REMOVED***
				defaultDec = r.kindDecoders[t.Kind()]
			***REMOVED***
			return newCondAddrDecoder(idec.vd, defaultDec), true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// LookupTypeMapEntry inspects the registry's type map for a Go type for the corresponding BSON
// type. If no type is found, ErrNoTypeMapEntry is returned.
func (r *Registry) LookupTypeMapEntry(bt bsontype.Type) (reflect.Type, error) ***REMOVED***
	t, ok := r.typeMap[bt]
	if !ok || t == nil ***REMOVED***
		return nil, ErrNoTypeMapEntry***REMOVED***Type: bt***REMOVED***
	***REMOVED***
	return t, nil
***REMOVED***

type interfaceValueEncoder struct ***REMOVED***
	i  reflect.Type
	ve ValueEncoder
***REMOVED***

type interfaceValueDecoder struct ***REMOVED***
	i  reflect.Type
	vd ValueDecoder
***REMOVED***
