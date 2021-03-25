// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import "reflect"

// This file exists, so that the files for specific formats do not all import reflect.
// This just helps us ensure that reflect package is isolated to a few files.

// SetInterfaceExt sets an extension
func (h *JsonHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, makeExt(ext))
***REMOVED***

// SetInterfaceExt sets an extension
func (h *CborHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, makeExt(ext))
***REMOVED***

// SetBytesExt sets an extension
func (h *MsgpackHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, makeExt(ext))
***REMOVED***

// SetBytesExt sets an extension
func (h *SimpleHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, makeExt(ext))
***REMOVED***

// SetBytesExt sets an extension
func (h *BincHandle) SetBytesExt(rt reflect.Type, tag uint64, ext BytesExt) (err error) ***REMOVED***
	return h.SetExt(rt, tag, makeExt(ext))
***REMOVED***

// func (h *XMLHandle) SetInterfaceExt(rt reflect.Type, tag uint64, ext InterfaceExt) (err error) ***REMOVED***
// 	return h.SetExt(rt, tag, &interfaceExtWrapper***REMOVED***InterfaceExt: ext***REMOVED***)
// ***REMOVED***
