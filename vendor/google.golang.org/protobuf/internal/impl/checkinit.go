// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"sync"

	"google.golang.org/protobuf/internal/errors"
	pref "google.golang.org/protobuf/reflect/protoreflect"
	piface "google.golang.org/protobuf/runtime/protoiface"
)

func (mi *MessageInfo) checkInitialized(in piface.CheckInitializedInput) (piface.CheckInitializedOutput, error) ***REMOVED***
	var p pointer
	if ms, ok := in.Message.(*messageState); ok ***REMOVED***
		p = ms.pointer()
	***REMOVED*** else ***REMOVED***
		p = in.Message.(*messageReflectWrapper).pointer()
	***REMOVED***
	return piface.CheckInitializedOutput***REMOVED******REMOVED***, mi.checkInitializedPointer(p)
***REMOVED***

func (mi *MessageInfo) checkInitializedPointer(p pointer) error ***REMOVED***
	mi.init()
	if !mi.needsInitCheck ***REMOVED***
		return nil
	***REMOVED***
	if p.IsNil() ***REMOVED***
		for _, f := range mi.orderedCoderFields ***REMOVED***
			if f.isRequired ***REMOVED***
				return errors.RequiredNotSet(string(mi.Desc.Fields().ByNumber(f.num).FullName()))
			***REMOVED***
		***REMOVED***
		return nil
	***REMOVED***
	if mi.extensionOffset.IsValid() ***REMOVED***
		e := p.Apply(mi.extensionOffset).Extensions()
		if err := mi.isInitExtensions(e); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, f := range mi.orderedCoderFields ***REMOVED***
		if !f.isRequired && f.funcs.isInit == nil ***REMOVED***
			continue
		***REMOVED***
		fptr := p.Apply(f.offset)
		if f.isPointer && fptr.Elem().IsNil() ***REMOVED***
			if f.isRequired ***REMOVED***
				return errors.RequiredNotSet(string(mi.Desc.Fields().ByNumber(f.num).FullName()))
			***REMOVED***
			continue
		***REMOVED***
		if f.funcs.isInit == nil ***REMOVED***
			continue
		***REMOVED***
		if err := f.funcs.isInit(fptr, f); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (mi *MessageInfo) isInitExtensions(ext *map[int32]ExtensionField) error ***REMOVED***
	if ext == nil ***REMOVED***
		return nil
	***REMOVED***
	for _, x := range *ext ***REMOVED***
		ei := getExtensionFieldInfo(x.Type())
		if ei.funcs.isInit == nil ***REMOVED***
			continue
		***REMOVED***
		v := x.Value()
		if !v.IsValid() ***REMOVED***
			continue
		***REMOVED***
		if err := ei.funcs.isInit(v); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

var (
	needsInitCheckMu  sync.Mutex
	needsInitCheckMap sync.Map
)

// needsInitCheck reports whether a message needs to be checked for partial initialization.
//
// It returns true if the message transitively includes any required or extension fields.
func needsInitCheck(md pref.MessageDescriptor) bool ***REMOVED***
	if v, ok := needsInitCheckMap.Load(md); ok ***REMOVED***
		if has, ok := v.(bool); ok ***REMOVED***
			return has
		***REMOVED***
	***REMOVED***
	needsInitCheckMu.Lock()
	defer needsInitCheckMu.Unlock()
	return needsInitCheckLocked(md)
***REMOVED***

func needsInitCheckLocked(md pref.MessageDescriptor) (has bool) ***REMOVED***
	if v, ok := needsInitCheckMap.Load(md); ok ***REMOVED***
		// If has is true, we've previously determined that this message
		// needs init checks.
		//
		// If has is false, we've previously determined that it can never
		// be uninitialized.
		//
		// If has is not a bool, we've just encountered a cycle in the
		// message graph. In this case, it is safe to return false: If
		// the message does have required fields, we'll detect them later
		// in the graph traversal.
		has, ok := v.(bool)
		return ok && has
	***REMOVED***
	needsInitCheckMap.Store(md, struct***REMOVED******REMOVED******REMOVED******REMOVED***) // avoid cycles while descending into this message
	defer func() ***REMOVED***
		needsInitCheckMap.Store(md, has)
	***REMOVED***()
	if md.RequiredNumbers().Len() > 0 ***REMOVED***
		return true
	***REMOVED***
	if md.ExtensionRanges().Len() > 0 ***REMOVED***
		return true
	***REMOVED***
	for i := 0; i < md.Fields().Len(); i++ ***REMOVED***
		fd := md.Fields().Get(i)
		// Map keys are never messages, so just consider the map value.
		if fd.IsMap() ***REMOVED***
			fd = fd.MapValue()
		***REMOVED***
		fmd := fd.Message()
		if fmd != nil && needsInitCheckLocked(fmd) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***
