// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"reflect"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoiface"
)

// ExtensionInfo implements ExtensionType.
//
// This type contains a number of exported fields for legacy compatibility.
// The only non-deprecated use of this type is through the methods of the
// ExtensionType interface.
type ExtensionInfo struct ***REMOVED***
	// An ExtensionInfo may exist in several stages of initialization.
	//
	// extensionInfoUninitialized: Some or all of the legacy exported
	// fields may be set, but none of the unexported fields have been
	// initialized. This is the starting state for an ExtensionInfo
	// in legacy generated code.
	//
	// extensionInfoDescInit: The desc field is set, but other unexported fields
	// may not be initialized. Legacy exported fields may or may not be set.
	// This is the starting state for an ExtensionInfo in newly generated code.
	//
	// extensionInfoFullInit: The ExtensionInfo is fully initialized.
	// This state is only entered after lazy initialization is complete.
	init uint32
	mu   sync.Mutex

	goType reflect.Type
	desc   extensionTypeDescriptor
	conv   Converter
	info   *extensionFieldInfo // for fast-path method implementations

	// ExtendedType is a typed nil-pointer to the parent message type that
	// is being extended. It is possible for this to be unpopulated in v2
	// since the message may no longer implement the MessageV1 interface.
	//
	// Deprecated: Use the ExtendedType method instead.
	ExtendedType protoiface.MessageV1

	// ExtensionType is the zero value of the extension type.
	//
	// For historical reasons, reflect.TypeOf(ExtensionType) and the
	// type returned by InterfaceOf may not be identical.
	//
	// Deprecated: Use InterfaceOf(xt.Zero()) instead.
	ExtensionType interface***REMOVED******REMOVED***

	// Field is the field number of the extension.
	//
	// Deprecated: Use the Descriptor().Number method instead.
	Field int32

	// Name is the fully qualified name of extension.
	//
	// Deprecated: Use the Descriptor().FullName method instead.
	Name string

	// Tag is the protobuf struct tag used in the v1 API.
	//
	// Deprecated: Do not use.
	Tag string

	// Filename is the proto filename in which the extension is defined.
	//
	// Deprecated: Use Descriptor().ParentFile().Path() instead.
	Filename string
***REMOVED***

// Stages of initialization: See the ExtensionInfo.init field.
const (
	extensionInfoUninitialized = 0
	extensionInfoDescInit      = 1
	extensionInfoFullInit      = 2
)

func InitExtensionInfo(xi *ExtensionInfo, xd protoreflect.ExtensionDescriptor, goType reflect.Type) ***REMOVED***
	xi.goType = goType
	xi.desc = extensionTypeDescriptor***REMOVED***xd, xi***REMOVED***
	xi.init = extensionInfoDescInit
***REMOVED***

func (xi *ExtensionInfo) New() protoreflect.Value ***REMOVED***
	return xi.lazyInit().New()
***REMOVED***
func (xi *ExtensionInfo) Zero() protoreflect.Value ***REMOVED***
	return xi.lazyInit().Zero()
***REMOVED***
func (xi *ExtensionInfo) ValueOf(v interface***REMOVED******REMOVED***) protoreflect.Value ***REMOVED***
	return xi.lazyInit().PBValueOf(reflect.ValueOf(v))
***REMOVED***
func (xi *ExtensionInfo) InterfaceOf(v protoreflect.Value) interface***REMOVED******REMOVED*** ***REMOVED***
	return xi.lazyInit().GoValueOf(v).Interface()
***REMOVED***
func (xi *ExtensionInfo) IsValidValue(v protoreflect.Value) bool ***REMOVED***
	return xi.lazyInit().IsValidPB(v)
***REMOVED***
func (xi *ExtensionInfo) IsValidInterface(v interface***REMOVED******REMOVED***) bool ***REMOVED***
	return xi.lazyInit().IsValidGo(reflect.ValueOf(v))
***REMOVED***
func (xi *ExtensionInfo) TypeDescriptor() protoreflect.ExtensionTypeDescriptor ***REMOVED***
	if atomic.LoadUint32(&xi.init) < extensionInfoDescInit ***REMOVED***
		xi.lazyInitSlow()
	***REMOVED***
	return &xi.desc
***REMOVED***

func (xi *ExtensionInfo) lazyInit() Converter ***REMOVED***
	if atomic.LoadUint32(&xi.init) < extensionInfoFullInit ***REMOVED***
		xi.lazyInitSlow()
	***REMOVED***
	return xi.conv
***REMOVED***

func (xi *ExtensionInfo) lazyInitSlow() ***REMOVED***
	xi.mu.Lock()
	defer xi.mu.Unlock()

	if xi.init == extensionInfoFullInit ***REMOVED***
		return
	***REMOVED***
	defer atomic.StoreUint32(&xi.init, extensionInfoFullInit)

	if xi.desc.ExtensionDescriptor == nil ***REMOVED***
		xi.initFromLegacy()
	***REMOVED***
	if !xi.desc.ExtensionDescriptor.IsPlaceholder() ***REMOVED***
		if xi.ExtensionType == nil ***REMOVED***
			xi.initToLegacy()
		***REMOVED***
		xi.conv = NewConverter(xi.goType, xi.desc.ExtensionDescriptor)
		xi.info = makeExtensionFieldInfo(xi.desc.ExtensionDescriptor)
		xi.info.validation = newValidationInfo(xi.desc.ExtensionDescriptor, xi.goType)
	***REMOVED***
***REMOVED***

type extensionTypeDescriptor struct ***REMOVED***
	protoreflect.ExtensionDescriptor
	xi *ExtensionInfo
***REMOVED***

func (xtd *extensionTypeDescriptor) Type() protoreflect.ExtensionType ***REMOVED***
	return xtd.xi
***REMOVED***
func (xtd *extensionTypeDescriptor) Descriptor() protoreflect.ExtensionDescriptor ***REMOVED***
	return xtd.ExtensionDescriptor
***REMOVED***
