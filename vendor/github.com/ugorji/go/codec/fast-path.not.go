// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

// +build notfastpath

package codec

import "reflect"

const fastpathEnabled = false

// The generated fast-path code is very large, and adds a few seconds to the build time.
// This causes test execution, execution of small tools which use codec, etc
// to take a long time.
//
// To mitigate, we now support the notfastpath tag.
// This tag disables fastpath during build, allowing for faster build, test execution,
// short-program runs, etc.

func fastpathDecodeTypeSwitch(iv interface***REMOVED******REMOVED***, d *Decoder) bool      ***REMOVED*** return false ***REMOVED***
func fastpathEncodeTypeSwitch(iv interface***REMOVED******REMOVED***, e *Encoder) bool      ***REMOVED*** return false ***REMOVED***
func fastpathEncodeTypeSwitchSlice(iv interface***REMOVED******REMOVED***, e *Encoder) bool ***REMOVED*** return false ***REMOVED***
func fastpathEncodeTypeSwitchMap(iv interface***REMOVED******REMOVED***, e *Encoder) bool   ***REMOVED*** return false ***REMOVED***
func fastpathDecodeSetZeroTypeSwitch(iv interface***REMOVED******REMOVED***) bool           ***REMOVED*** return false ***REMOVED***

type fastpathT struct***REMOVED******REMOVED***
type fastpathE struct ***REMOVED***
	rtid  uintptr
	rt    reflect.Type
	encfn func(*Encoder, *codecFnInfo, reflect.Value)
	decfn func(*Decoder, *codecFnInfo, reflect.Value)
***REMOVED***
type fastpathA [0]fastpathE

func (x fastpathA) index(rtid uintptr) int ***REMOVED*** return -1 ***REMOVED***

var fastpathAV fastpathA
var fastpathTV fastpathT

// ----
type TestMammoth2Wrapper struct***REMOVED******REMOVED*** // to allow testMammoth work in notfastpath mode
