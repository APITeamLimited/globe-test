// +build codecgen.exec

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

// DO NOT EDIT. THIS FILE IS AUTO-GENERATED FROM gen-dec-(map|array).go.tmpl

const genDecMapTmpl = `
***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** := ****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED***
***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** := z.DecReadMapStart()
if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** == codecSelferDecContainerLenNil***REMOVED******REMOVED***xs***REMOVED******REMOVED*** ***REMOVED***
	****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED*** = nil
***REMOVED*** else ***REMOVED***
if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
	***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** := z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, z.DecBasicHandle().MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
	***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(map[***REMOVED******REMOVED*** .KTyp ***REMOVED******REMOVED***]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
	****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***
***REMOVED***
var ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .KTyp ***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mdn"***REMOVED******REMOVED*** ***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED*** bool
if z.DecBasicHandle().MapValueReset ***REMOVED***
	***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED******REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true
	***REMOVED******REMOVED***else if decElemKindIntf***REMOVED******REMOVED***if !z.DecBasicHandle().InterfaceReset ***REMOVED*** ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true ***REMOVED***
	***REMOVED******REMOVED***else if not decElemKindImmutable***REMOVED******REMOVED******REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** = true
	***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED***
if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** != 0 ***REMOVED***
	***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > 0 
	for ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** := 0; (***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** && ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***) || !(***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** || z.DecCheckBreak()); ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED***
	z.DecReadMapElemKey()
	***REMOVED******REMOVED*** $x := printf "%vmk%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVarK $x -***REMOVED******REMOVED***
	***REMOVED******REMOVED*** if eq .KTyp "interface***REMOVED******REMOVED***" ***REMOVED******REMOVED******REMOVED******REMOVED***/* // special case if a byte array. */ -***REMOVED******REMOVED***
    if ***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***.([]byte); ***REMOVED******REMOVED***var "bok"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED*** = string(***REMOVED******REMOVED***var "bv"***REMOVED******REMOVED***)
	***REMOVED***
    ***REMOVED******REMOVED*** end -***REMOVED******REMOVED***
    ***REMOVED******REMOVED***if decElemKindPtr -***REMOVED******REMOVED***
	***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = true
    ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***var "mg"***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***if decElemKindPtr -***REMOVED******REMOVED***
        ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] 
		if ***REMOVED******REMOVED***var "mok"***REMOVED******REMOVED*** ***REMOVED***
			***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** = false
		***REMOVED***
        ***REMOVED******REMOVED***else -***REMOVED******REMOVED***
        ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***]
        ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
	***REMOVED*** ***REMOVED******REMOVED***if not decElemKindImmutable***REMOVED******REMOVED***else ***REMOVED*** ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***decElemZero***REMOVED******REMOVED*** ***REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***
	z.DecReadMapElemValue()
	***REMOVED******REMOVED***var "mdn"***REMOVED******REMOVED*** = false
	***REMOVED******REMOVED*** $x := printf "%vmv%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** $y := printf "%vmdn%v" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x $y -***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***var "mdn"***REMOVED******REMOVED*** ***REMOVED***
		if z.DecBasicHandle().DeleteOnNilMapValue ***REMOVED*** delete(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***) ***REMOVED*** else ***REMOVED*** ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] = ***REMOVED******REMOVED***decElemZero***REMOVED******REMOVED*** ***REMOVED***
	***REMOVED*** else if ***REMOVED******REMOVED***if decElemKindPtr***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "ms"***REMOVED******REMOVED*** && ***REMOVED******REMOVED***end***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** != nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[***REMOVED******REMOVED***var "mk"***REMOVED******REMOVED***] = ***REMOVED******REMOVED***var "mv"***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***
***REMOVED*** // else len==0: TODO: Should we clear map entries?
z.DecReadMapEnd()
***REMOVED***
`

const genDecListTmpl = `
***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***if not isArray***REMOVED******REMOVED*******REMOVED******REMOVED***end***REMOVED******REMOVED******REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED***
***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** := z.DecSliceHelperStart() ***REMOVED******REMOVED***/* // helper, containerLenS */***REMOVED******REMOVED***
***REMOVED******REMOVED***if not isArray -***REMOVED******REMOVED***
var ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** bool ***REMOVED******REMOVED***/* // changed */***REMOVED******REMOVED***
_ = ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED***
if ***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.IsNil ***REMOVED***
	if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** != nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = nil
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED***
***REMOVED*** else ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** == 0 ***REMOVED***
	***REMOVED******REMOVED***if isSlice -***REMOVED******REMOVED***
	if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = []***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED******REMOVED******REMOVED***
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** else if len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) != 0 ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:0]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** ***REMOVED******REMOVED***else if isChan ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(***REMOVED******REMOVED*** .CTyp ***REMOVED******REMOVED***, 0)
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED***
    ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
***REMOVED*** else ***REMOVED***
	***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** := ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > 0
	var ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** int
	_ =  ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***
	***REMOVED******REMOVED***if isSlice ***REMOVED******REMOVED*** if ***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** ***REMOVED***
	if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** > cap(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** = z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, z.DecBasicHandle().MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
		if ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** <= cap(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
			***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***]
		***REMOVED*** else ***REMOVED***
			***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make([]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
		***REMOVED***
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** else if ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED*** != len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED***
	***REMOVED***
    ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
	var ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** int 
    ***REMOVED******REMOVED***/* // var ***REMOVED******REMOVED***var "dn"***REMOVED******REMOVED*** bool */ -***REMOVED******REMOVED***
	for ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** = 0; (***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** && ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < ***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***) || !(***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** || z.DecCheckBreak()); ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***++ ***REMOVED*** // bounds-check-elimination
		***REMOVED******REMOVED***if not isArray***REMOVED******REMOVED*** if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** == 0 && ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
			if ***REMOVED******REMOVED***var "hl"***REMOVED******REMOVED*** ***REMOVED***
				***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** = z.DecInferLen(***REMOVED******REMOVED***var "l"***REMOVED******REMOVED***, z.DecBasicHandle().MaxInitLen, ***REMOVED******REMOVED*** .Size ***REMOVED******REMOVED***)
			***REMOVED*** else ***REMOVED***
				***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***if isSlice***REMOVED******REMOVED***8***REMOVED******REMOVED***else if isChan***REMOVED******REMOVED***64***REMOVED******REMOVED***end***REMOVED******REMOVED***
			***REMOVED***
			***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make(***REMOVED******REMOVED***if isSlice***REMOVED******REMOVED***[]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED******REMOVED******REMOVED***else if isChan***REMOVED******REMOVED******REMOVED******REMOVED***.CTyp***REMOVED******REMOVED******REMOVED******REMOVED***end***REMOVED******REMOVED***, ***REMOVED******REMOVED***var "rl"***REMOVED******REMOVED***)
			***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true 
		***REMOVED***
        ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.ElemContainerState(***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***)
        ***REMOVED******REMOVED***/* ***REMOVED******REMOVED***var "dn"***REMOVED******REMOVED*** = r.TryDecodeAsNil() */***REMOVED******REMOVED******REMOVED******REMOVED***/* commented out, as decLineVar handles this already each time */ -***REMOVED******REMOVED***
        ***REMOVED******REMOVED***if isChan***REMOVED******REMOVED******REMOVED******REMOVED*** $x := printf "%[1]vvcx%[2]v" .TempVar .Rand ***REMOVED******REMOVED***var ***REMOVED******REMOVED***$x***REMOVED******REMOVED*** ***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***
		***REMOVED******REMOVED*** decLineVar $x -***REMOVED******REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** <- ***REMOVED******REMOVED*** $x ***REMOVED******REMOVED***
        ***REMOVED******REMOVED***else***REMOVED******REMOVED******REMOVED******REMOVED***/* // if indefinite, etc, then expand the slice if necessary */ -***REMOVED******REMOVED***
		var ***REMOVED******REMOVED***var "db"***REMOVED******REMOVED*** bool
		if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** >= len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
			***REMOVED******REMOVED***if isSlice ***REMOVED******REMOVED*** ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***, ***REMOVED******REMOVED*** zero ***REMOVED******REMOVED***)
			***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
			***REMOVED******REMOVED***else***REMOVED******REMOVED*** z.DecArrayCannotExpand(len(v), ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***+1); ***REMOVED******REMOVED***var "db"***REMOVED******REMOVED*** = true
			***REMOVED******REMOVED***end -***REMOVED******REMOVED***
		***REMOVED***
		if ***REMOVED******REMOVED***var "db"***REMOVED******REMOVED*** ***REMOVED***
			z.DecSwallow()
		***REMOVED*** else ***REMOVED***
			***REMOVED******REMOVED*** $x := printf "%[1]vv%[2]v[%[1]vj%[2]v]" .TempVar .Rand ***REMOVED******REMOVED******REMOVED******REMOVED*** decLineVar $x -***REMOVED******REMOVED***
		***REMOVED***
        ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
	***REMOVED***
	***REMOVED******REMOVED***if isSlice***REMOVED******REMOVED*** if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** < len(***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***) ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***[:***REMOVED******REMOVED***var "j"***REMOVED******REMOVED***]
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED*** else if ***REMOVED******REMOVED***var "j"***REMOVED******REMOVED*** == 0 && ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** == nil ***REMOVED***
		***REMOVED******REMOVED***var "v"***REMOVED******REMOVED*** = make([]***REMOVED******REMOVED*** .Typ ***REMOVED******REMOVED***, 0)
		***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** = true
	***REMOVED***
    ***REMOVED******REMOVED***end -***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED***var "h"***REMOVED******REMOVED***.End()
***REMOVED******REMOVED***if not isArray ***REMOVED******REMOVED***if ***REMOVED******REMOVED***var "c"***REMOVED******REMOVED*** ***REMOVED*** 
	****REMOVED******REMOVED*** .Varname ***REMOVED******REMOVED*** = ***REMOVED******REMOVED***var "v"***REMOVED******REMOVED***
***REMOVED***
***REMOVED******REMOVED***end -***REMOVED******REMOVED***
`

const genEncChanTmpl = `
***REMOVED******REMOVED***.Label***REMOVED******REMOVED***:
switch timeout***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** :=  z.EncBasicHandle().ChanRecvTimeout; ***REMOVED***
case timeout***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** == 0: // only consume available
	for ***REMOVED***
		select ***REMOVED***
		case b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** := <-***REMOVED******REMOVED***.Chan***REMOVED******REMOVED***:
			***REMOVED******REMOVED*** .Slice ***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***.Slice***REMOVED******REMOVED***, b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED***)
		default:
			break ***REMOVED******REMOVED***.Label***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
case timeout***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** > 0: // consume until timeout
	tt***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** := time.NewTimer(timeout***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED***)
	for ***REMOVED***
		select ***REMOVED***
		case b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** := <-***REMOVED******REMOVED***.Chan***REMOVED******REMOVED***:
			***REMOVED******REMOVED***.Slice***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***.Slice***REMOVED******REMOVED***, b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED***)
		case <-tt***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED***.C:
			// close(tt.C)
			break ***REMOVED******REMOVED***.Label***REMOVED******REMOVED***
		***REMOVED***
	***REMOVED***
default: // consume until close
	for b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED*** := range ***REMOVED******REMOVED***.Chan***REMOVED******REMOVED*** ***REMOVED***
		***REMOVED******REMOVED***.Slice***REMOVED******REMOVED*** = append(***REMOVED******REMOVED***.Slice***REMOVED******REMOVED***, b***REMOVED******REMOVED***.Sfx***REMOVED******REMOVED***)
	***REMOVED***
***REMOVED***
`
