// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// DO NOT EDIT. THIS FILE IS AUTO-GENERATED FROM gen-dec-(map|array).go.tmpl

const genDecMapTmpl = `
***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** := ****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED***
***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** := r.ReadMapStart()
***REMOVED******REMOVED***var "bh"***REMOVED******REMOVED*** := z.DecBasicHandle()
if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
	***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***, _ := z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "bh"***REMOVED******REMOVED***.MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
	***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(map[***REMOVED******REMOVED*** .KTyp ***REMOVED******REMOVED***]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
	****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***
***REMOVED***
var ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .KTyp ***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** ***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED*** bool
if ***REMOVED******REMOVED***var "bh"***REMOVED******REMOVED***.MapValueReset ***REMOVED***
	***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED******REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true
	***REMOVED******REMOVED***else if decElemKindIntf***REMOVED******REMOVED***if !***REMOVED******REMOVED***var "bh"***REMOVED******REMOVED***.InterfaceReset ***REMOVED*** ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true ***REMOVED***
	***REMOVED******REMOVED***else if not decElemKindImmutable***REMOVED******REMOVED******REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true
	***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED***
if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > 0  ***REMOVED***
for ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** := 0; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
	z.DecSendContainerState(codecSelfer_containerMapKey***REMOVED******REMOVED*** .Sfx ***REMOVED******REMOVED***)
	***REMOVED******REMOVED*** $x := printf "%vmk%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVarK $x ***REMOVED******REMOVED***
***REMOVED******REMOVED*** if eq .KTyp "interface***REMOVED******REMOVED***" ***REMOVED******REMOVED******REMOVED******REMOVED***/* // special case if a byte array. */***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***.([]byte); ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED*** = string(***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***)
	***REMOVED******REMOVED******REMOVED*** end ***REMOVED******REMOVED******REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED***
	***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = true***REMOVED******REMOVED***end***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED******REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] 
		if ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** ***REMOVED***
			***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = false
		***REMOVED*** ***REMOVED******REMOVED***else***REMOVED******REMOVED******REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] ***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED*** ***REMOVED******REMOVED***if not decElemKindImmutable***REMOVED******REMOVED***else ***REMOVED*** ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***decElemZero***REMOVED******REMOVED*** ***REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
	z.DecSendContainerState(codecSelfer_containerMapValue***REMOVED******REMOVED*** .Sfx ***REMOVED******REMOVED***)
	***REMOVED******REMOVED*** $x := printf "%vmv%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** && ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** != nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] = ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***
***REMOVED*** else if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** < 0  ***REMOVED***
for ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** := 0; !r.CheckBreak(); ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
	z.DecSendContainerState(codecSelfer_containerMapKey***REMOVED******REMOVED*** .Sfx ***REMOVED******REMOVED***)
	***REMOVED******REMOVED*** $x := printf "%vmk%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVarK $x ***REMOVED******REMOVED***
***REMOVED******REMOVED*** if eq .KTyp "interface***REMOVED******REMOVED***" ***REMOVED******REMOVED******REMOVED******REMOVED***/* // special case if a byte array. */***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***.([]byte); ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED*** = string(***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***)
	***REMOVED******REMOVED******REMOVED*** end ***REMOVED******REMOVED******REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED***
	***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = true ***REMOVED******REMOVED*** end ***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED******REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] 
		if ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** ***REMOVED***
			***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = false
		***REMOVED*** ***REMOVED******REMOVED***else***REMOVED******REMOVED******REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] ***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED*** ***REMOVED******REMOVED***if not decElemKindImmutable***REMOVED******REMOVED***else ***REMOVED*** ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***decElemZero***REMOVED******REMOVED*** ***REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
	z.DecSendContainerState(codecSelfer_containerMapValue***REMOVED******REMOVED*** .Sfx ***REMOVED******REMOVED***)
	***REMOVED******REMOVED*** $x := printf "%vmv%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** && ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** != nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] = ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***
***REMOVED*** // else len==0: TODO: Should we clear map entries?
z.DecSendContainerState(codecSelfer_containerMapEnd***REMOVED******REMOVED*** .Sfx ***REMOVED******REMOVED***)
`

const genDecListTmpl = `
***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***if not isArray***REMOVED******REMOVED*******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED***
***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** := z.DecSliceHelperStart() ***REMOVED******REMOVED***/* // helper, containerLenS */***REMOVED******REMOVED******REMOVED******REMOVED***if not isArray***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** bool ***REMOVED******REMOVED***/* // changed */***REMOVED******REMOVED***
_ = ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** == 0 ***REMOVED***
	***REMOVED******REMOVED***if isSlice ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = []***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** else if len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) != 0 ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:0]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***if isChan ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(***REMOVED******REMOVED*** .CTyp ***REMOVED******REMOVED***, 0)
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED*** else if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > 0 ***REMOVED***
	***REMOVED******REMOVED***if isChan ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***, _ = z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, z.DecBasicHandle().MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(***REMOVED******REMOVED*** .CTyp ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED***
	for ***REMOVED******REMOVED***var "r"***REMOVED******REMOVED*** := 0; ***REMOVED******REMOVED***var "r"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***; ***REMOVED******REMOVED***var "r"***REMOVED******REMOVED***++ ***REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "r"***REMOVED******REMOVED***)
		var ***REMOVED******REMOVED***var "t"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
		***REMOVED******REMOVED*** $x := printf "%st%s" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** <- ***REMOVED******REMOVED***var "t"***REMOVED******REMOVED***
	***REMOVED***
	***REMOVED******REMOVED*** else ***REMOVED******REMOVED***	var ***REMOVED******REMOVED***var "rr"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** int ***REMOVED******REMOVED***/* // num2read, length of slice/array/chan */***REMOVED******REMOVED***
	var ***REMOVED******REMOVED***var "rt"***REMOVED******REMOVED*** bool ***REMOVED******REMOVED***/* truncated */***REMOVED******REMOVED***
	_, _ = ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rt"***REMOVED******REMOVED***
	***REMOVED******REMOVED***var "rr"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** // len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***)
	if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > cap(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***if isArray ***REMOVED******REMOVED***z.DecArrayCannotExpand(len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***), ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***)
		***REMOVED******REMOVED*** else ***REMOVED******REMOVED******REMOVED******REMOVED***if not .Immutable ***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "rg"***REMOVED******REMOVED*** := len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) > 0
		***REMOVED******REMOVED***var "v2"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rt"***REMOVED******REMOVED*** = z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, z.DecBasicHandle().MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
		if ***REMOVED******REMOVED***var "rt"***REMOVED******REMOVED*** ***REMOVED***
			if ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** <= cap(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
				***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***]
			***REMOVED*** else ***REMOVED***
				***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make([]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make([]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
		***REMOVED***
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
		***REMOVED******REMOVED***var "rr"***REMOVED******REMOVED*** = len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED******REMOVED***if not .Immutable ***REMOVED******REMOVED***
			if ***REMOVED******REMOVED***var "rg"***REMOVED******REMOVED*** ***REMOVED*** copy(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "v2"***REMOVED******REMOVED***) ***REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED***/* end not Immutable, isArray */***REMOVED******REMOVED***
	***REMOVED*** ***REMOVED******REMOVED***if isSlice ***REMOVED******REMOVED*** else if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** != len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED***	***REMOVED******REMOVED***/* end isSlice:47 */***REMOVED******REMOVED*** 
	***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** := 0
	for ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "rr"***REMOVED******REMOVED*** ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
		***REMOVED******REMOVED*** $x := printf "%[1]vv%[2]v[%[1]vj%[2]v]" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
	***REMOVED***
	***REMOVED******REMOVED***if isArray ***REMOVED******REMOVED***for ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
		z.DecSwallow()
	***REMOVED***
	***REMOVED******REMOVED*** else ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "rt"***REMOVED******REMOVED*** ***REMOVED***
		for ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** ; ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
			***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***, ***REMOVED******REMOVED*** zero***REMOVED******REMOVED***)
			***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
			***REMOVED******REMOVED*** $x := printf "%[1]vv%[2]v[%[1]vj%[2]v]" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED*** ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* end isArray:56 */***REMOVED******REMOVED***
	***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***/* end isChan:16 */***REMOVED******REMOVED***
***REMOVED*** else ***REMOVED*** ***REMOVED******REMOVED***/* len < 0 */***REMOVED******REMOVED***
	***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** := 0
	for ; !r.CheckBreak(); ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
		***REMOVED******REMOVED***if isChan ***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
		var ***REMOVED******REMOVED***var "t"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
		***REMOVED******REMOVED*** $x := printf "%st%s" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** <- ***REMOVED******REMOVED***var "t"***REMOVED******REMOVED*** 
		***REMOVED******REMOVED*** else ***REMOVED******REMOVED***
		if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** >= len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
			***REMOVED******REMOVED***if isArray ***REMOVED******REMOVED***z.DecArrayCannotExpand(len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***), ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***+1)
			***REMOVED******REMOVED*** else ***REMOVED******REMOVED******REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***, ***REMOVED******REMOVED***zero***REMOVED******REMOVED***)// var ***REMOVED******REMOVED***var "z"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
			***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true ***REMOVED******REMOVED***end***REMOVED******REMOVED***
		***REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
		if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
			***REMOVED******REMOVED*** $x := printf "%[1]vv%[2]v[%[1]vj%[2]v]" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x ***REMOVED******REMOVED***
		***REMOVED*** else ***REMOVED***
			z.DecSwallow()
		***REMOVED***
		***REMOVED******REMOVED***end***REMOVED******REMOVED***
	***REMOVED***
	***REMOVED******REMOVED***if isSlice ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** else if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** == 0 && ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = []***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.End()
***REMOVED******REMOVED***if not isArray ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** ***REMOVED*** 
	****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***
***REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
`

