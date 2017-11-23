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

type fastpathT struct***REMOVED******REMOVED***
type fastpathE struct ***REMOVED***
	rtid  uintptr
	rt    reflect.Type
	encfn func(*encFnInfo, reflect.Value)
	decfn func(*decFnInfo, reflect.Value)
***REMOVED***
type fastpathA [0]fastpathE

func (x fastpathA) index(rtid uintptr) int ***REMOVED*** return -1 ***REMOVED***

var fastpathAV fastpathA
var fastpathTV fastpathT
