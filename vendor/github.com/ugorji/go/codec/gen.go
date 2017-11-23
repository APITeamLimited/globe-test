// Copyright (c) 2012-2015 Ugorji Nwoke. All rights reserved.
// Use of this source code is governed by a MIT license found in the LICENSE file.

package codec

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
	"unicode"
	"unicode/utf8"
)

// ---------------------------------------------------
// codecgen supports the full cycle of reflection-based codec:
//    - RawExt
//    - Raw
//    - Builtins
//    - Extensions
//    - (Binary|Text|JSON)(Unm|M)arshal
//    - generic by-kind
//
// This means that, for dynamic things, we MUST use reflection to at least get the reflect.Type.
// In those areas, we try to only do reflection or interface-conversion when NECESSARY:
//    - Extensions, only if Extensions are configured.
//
// However, codecgen doesn't support the following:
//   - Canonical option. (codecgen IGNORES it currently)
//     This is just because it has not been implemented.
//
// During encode/decode, Selfer takes precedence.
// A type implementing Selfer will know how to encode/decode itself statically.
//
// The following field types are supported:
//     array: [n]T
//     slice: []T
//     map: map[K]V
//     primitive: [u]int[n], float(32|64), bool, string
//     struct
//
// ---------------------------------------------------
// Note that a Selfer cannot call (e|d).(En|De)code on itself,
// as this will cause a circular reference, as (En|De)code will call Selfer methods.
// Any type that implements Selfer must implement completely and not fallback to (En|De)code.
//
// In addition, code in this file manages the generation of fast-path implementations of
// encode/decode of slices/maps of primitive keys/values.
//
// Users MUST re-generate their implementations whenever the code shape changes.
// The generated code will panic if it was generated with a version older than the supporting library.
// ---------------------------------------------------
//
// codec framework is very feature rich.
// When encoding or decoding into an interface, it depends on the runtime type of the interface.
// The type of the interface may be a named type, an extension, etc.
// Consequently, we fallback to runtime codec for encoding/decoding interfaces.
// In addition, we fallback for any value which cannot be guaranteed at runtime.
// This allows us support ANY value, including any named types, specifically those which
// do not implement our interfaces (e.g. Selfer).
//
// This explains some slowness compared to other code generation codecs (e.g. msgp).
// This reduction in speed is only seen when your refers to interfaces,
// e.g. type T struct ***REMOVED*** A interface***REMOVED******REMOVED***; B []interface***REMOVED******REMOVED***; C map[string]interface***REMOVED******REMOVED*** ***REMOVED***
//
// codecgen will panic if the file was generated with an old version of the library in use.
//
// Note:
//   It was a conscious decision to have gen.go always explicitly call EncodeNil or TryDecodeAsNil.
//   This way, there isn't a function call overhead just to see that we should not enter a block of code.

// GenVersion is the current version of codecgen.
//
// NOTE: Increment this value each time codecgen changes fundamentally.
// Fundamental changes are:
//   - helper methods change (signature change, new ones added, some removed, etc)
//   - codecgen command line changes
//
// v1: Initial Version
// v2:
// v3: Changes for Kubernetes:
//     changes in signature of some unpublished helper methods and codecgen cmdline arguments.
// v4: Removed separator support from (en|de)cDriver, and refactored codec(gen)
// v5: changes to support faster json decoding. Let encoder/decoder maintain state of collections.
const GenVersion = 5

const (
	genCodecPkg        = "codec1978"
	genTempVarPfx      = "yy"
	genTopLevelVarName = "x"

	// ignore canBeNil parameter, and always set to true.
	// This is because nil can appear anywhere, so we should always check.
	genAnythingCanBeNil = true

	// if genUseOneFunctionForDecStructMap, make a single codecDecodeSelferFromMap function;
	// else make codecDecodeSelferFromMap***REMOVED***LenPrefix,CheckBreak***REMOVED*** so that conditionals
	// are not executed a lot.
	//
	// From testing, it didn't make much difference in runtime, so keep as true (one function only)
	genUseOneFunctionForDecStructMap = true
)

type genStructMapStyle uint8

const (
	genStructMapStyleConsolidated genStructMapStyle = iota
	genStructMapStyleLenPrefix
	genStructMapStyleCheckBreak
)

var (
	genAllTypesSamePkgErr  = errors.New("All types must be in the same package")
	genExpectArrayOrMapErr = errors.New("unexpected type. Expecting array/map/slice")
	genBase64enc           = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789__")
	genQNameRegex          = regexp.MustCompile(`[A-Za-z_.]+`)
	genCheckVendor         bool
)

// genRunner holds some state used during a Gen run.
type genRunner struct ***REMOVED***
	w io.Writer      // output
	c uint64         // counter used for generating varsfx
	t []reflect.Type // list of types to run selfer on

	tc reflect.Type     // currently running selfer on this type
	te map[uintptr]bool // types for which the encoder has been created
	td map[uintptr]bool // types for which the decoder has been created
	cp string           // codec import path

	im  map[string]reflect.Type // imports to add
	imn map[string]string       // package names of imports to add
	imc uint64                  // counter for import numbers

	is map[reflect.Type]struct***REMOVED******REMOVED*** // types seen during import search
	bp string                    // base PkgPath, for which we are generating for

	cpfx   string // codec package prefix
	unsafe bool   // is unsafe to be used in generated code?

	tm map[reflect.Type]struct***REMOVED******REMOVED*** // types for which enc/dec must be generated
	ts []reflect.Type            // types for which enc/dec must be generated

	xs string // top level variable/constant suffix
	hn string // fn helper type name

	ti *TypeInfos
	// rr *rand.Rand // random generator for file-specific types
***REMOVED***

// Gen will write a complete go file containing Selfer implementations for each
// type passed. All the types must be in the same package.
//
// Library users: *DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINOUSLY WITHOUT NOTICE.*
func Gen(w io.Writer, buildTags, pkgName, uid string, useUnsafe bool, ti *TypeInfos, typ ...reflect.Type) ***REMOVED***
	// All types passed to this method do not have a codec.Selfer method implemented directly.
	// codecgen already checks the AST and skips any types that define the codec.Selfer methods.
	// Consequently, there's no need to check and trim them if they implement codec.Selfer

	if len(typ) == 0 ***REMOVED***
		return
	***REMOVED***
	x := genRunner***REMOVED***
		unsafe: useUnsafe,
		w:      w,
		t:      typ,
		te:     make(map[uintptr]bool),
		td:     make(map[uintptr]bool),
		im:     make(map[string]reflect.Type),
		imn:    make(map[string]string),
		is:     make(map[reflect.Type]struct***REMOVED******REMOVED***),
		tm:     make(map[reflect.Type]struct***REMOVED******REMOVED***),
		ts:     []reflect.Type***REMOVED******REMOVED***,
		bp:     genImportPath(typ[0]),
		xs:     uid,
		ti:     ti,
	***REMOVED***
	if x.ti == nil ***REMOVED***
		x.ti = defTypeInfos
	***REMOVED***
	if x.xs == "" ***REMOVED***
		rr := rand.New(rand.NewSource(time.Now().UnixNano()))
		x.xs = strconv.FormatInt(rr.Int63n(9999), 10)
	***REMOVED***

	// gather imports first:
	x.cp = genImportPath(reflect.TypeOf(x))
	x.imn[x.cp] = genCodecPkg
	for _, t := range typ ***REMOVED***
		// fmt.Printf("###########: PkgPath: '%v', Name: '%s'\n", genImportPath(t), t.Name())
		if genImportPath(t) != x.bp ***REMOVED***
			panic(genAllTypesSamePkgErr)
		***REMOVED***
		x.genRefPkgs(t)
	***REMOVED***
	if buildTags != "" ***REMOVED***
		x.line("// +build " + buildTags)
		x.line("")
	***REMOVED***
	x.line(`

// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED BY codecgen.
// ************************************************************

`)
	x.line("package " + pkgName)
	x.line("")
	x.line("import (")
	if x.cp != x.bp ***REMOVED***
		x.cpfx = genCodecPkg + "."
		x.linef("%s \"%s\"", genCodecPkg, x.cp)
	***REMOVED***
	// use a sorted set of im keys, so that we can get consistent output
	imKeys := make([]string, 0, len(x.im))
	for k, _ := range x.im ***REMOVED***
		imKeys = append(imKeys, k)
	***REMOVED***
	sort.Strings(imKeys)
	for _, k := range imKeys ***REMOVED*** // for k, _ := range x.im ***REMOVED***
		x.linef("%s \"%s\"", x.imn[k], k)
	***REMOVED***
	// add required packages
	for _, k := range [...]string***REMOVED***"reflect", "unsafe", "runtime", "fmt", "errors"***REMOVED*** ***REMOVED***
		if _, ok := x.im[k]; !ok ***REMOVED***
			if k == "unsafe" && !x.unsafe ***REMOVED***
				continue
			***REMOVED***
			x.line("\"" + k + "\"")
		***REMOVED***
	***REMOVED***
	x.line(")")
	x.line("")

	x.line("const (")
	x.linef("// ----- content types ----")
	x.linef("codecSelferC_UTF8%s = %v", x.xs, int64(c_UTF8))
	x.linef("codecSelferC_RAW%s = %v", x.xs, int64(c_RAW))
	x.linef("// ----- value types used ----")
	x.linef("codecSelferValueTypeArray%s = %v", x.xs, int64(valueTypeArray))
	x.linef("codecSelferValueTypeMap%s = %v", x.xs, int64(valueTypeMap))
	x.linef("// ----- containerStateValues ----")
	x.linef("codecSelfer_containerMapKey%s = %v", x.xs, int64(containerMapKey))
	x.linef("codecSelfer_containerMapValue%s = %v", x.xs, int64(containerMapValue))
	x.linef("codecSelfer_containerMapEnd%s = %v", x.xs, int64(containerMapEnd))
	x.linef("codecSelfer_containerArrayElem%s = %v", x.xs, int64(containerArrayElem))
	x.linef("codecSelfer_containerArrayEnd%s = %v", x.xs, int64(containerArrayEnd))
	x.line(")")
	x.line("var (")
	x.line("codecSelferBitsize" + x.xs + " = uint8(reflect.TypeOf(uint(0)).Bits())")
	x.line("codecSelferOnlyMapOrArrayEncodeToStructErr" + x.xs + " = errors.New(`only encoded map or array can be decoded into a struct`)")
	x.line(")")
	x.line("")

	if x.unsafe ***REMOVED***
		x.line("type codecSelferUnsafeString" + x.xs + " struct ***REMOVED*** Data uintptr; Len int***REMOVED***")
		x.line("")
	***REMOVED***
	x.hn = "codecSelfer" + x.xs
	x.line("type " + x.hn + " struct***REMOVED******REMOVED***")
	x.line("")

	x.varsfxreset()
	x.line("func init() ***REMOVED***")
	x.linef("if %sGenVersion != %v ***REMOVED***", x.cpfx, GenVersion)
	x.line("_, file, _, _ := runtime.Caller(0)")
	x.line(`err := fmt.Errorf("codecgen version mismatch: current: %v, need %v. Re-generate file: %v", `)
	x.linef(`%v, %sGenVersion, file)`, GenVersion, x.cpfx)
	x.line("panic(err)")
	x.linef("***REMOVED***")
	x.line("if false ***REMOVED*** // reference the types, but skip this branch at build/run time")
	var n int
	// for k, t := range x.im ***REMOVED***
	for _, k := range imKeys ***REMOVED***
		t := x.im[k]
		x.linef("var v%v %s.%s", n, x.imn[k], t.Name())
		n++
	***REMOVED***
	if x.unsafe ***REMOVED***
		x.linef("var v%v unsafe.Pointer", n)
		n++
	***REMOVED***
	if n > 0 ***REMOVED***
		x.out("_")
		for i := 1; i < n; i++ ***REMOVED***
			x.out(", _")
		***REMOVED***
		x.out(" = v0")
		for i := 1; i < n; i++ ***REMOVED***
			x.outf(", v%v", i)
		***REMOVED***
	***REMOVED***
	x.line("***REMOVED*** ") // close if false
	x.line("***REMOVED***")  // close init
	x.line("")

	// generate rest of type info
	for _, t := range typ ***REMOVED***
		x.tc = t
		x.selfer(true)
		x.selfer(false)
	***REMOVED***

	for _, t := range x.ts ***REMOVED***
		rtid := reflect.ValueOf(t).Pointer()
		// generate enc functions for all these slice/map types.
		x.varsfxreset()
		x.linef("func (x %s) enc%s(v %s%s, e *%sEncoder) ***REMOVED***", x.hn, x.genMethodNameT(t), x.arr2str(t, "*"), x.genTypeName(t), x.cpfx)
		x.genRequiredMethodVars(true)
		switch t.Kind() ***REMOVED***
		case reflect.Array, reflect.Slice, reflect.Chan:
			x.encListFallback("v", t)
		case reflect.Map:
			x.encMapFallback("v", t)
		default:
			panic(genExpectArrayOrMapErr)
		***REMOVED***
		x.line("***REMOVED***")
		x.line("")

		// generate dec functions for all these slice/map types.
		x.varsfxreset()
		x.linef("func (x %s) dec%s(v *%s, d *%sDecoder) ***REMOVED***", x.hn, x.genMethodNameT(t), x.genTypeName(t), x.cpfx)
		x.genRequiredMethodVars(false)
		switch t.Kind() ***REMOVED***
		case reflect.Array, reflect.Slice, reflect.Chan:
			x.decListFallback("v", rtid, t)
		case reflect.Map:
			x.decMapFallback("v", rtid, t)
		default:
			panic(genExpectArrayOrMapErr)
		***REMOVED***
		x.line("***REMOVED***")
		x.line("")
	***REMOVED***

	x.line("")
***REMOVED***

func (x *genRunner) checkForSelfer(t reflect.Type, varname string) bool ***REMOVED***
	// return varname != genTopLevelVarName && t != x.tc
	// the only time we checkForSelfer is if we are not at the TOP of the generated code.
	return varname != genTopLevelVarName
***REMOVED***

func (x *genRunner) arr2str(t reflect.Type, s string) string ***REMOVED***
	if t.Kind() == reflect.Array ***REMOVED***
		return s
	***REMOVED***
	return ""
***REMOVED***

func (x *genRunner) genRequiredMethodVars(encode bool) ***REMOVED***
	x.line("var h " + x.hn)
	if encode ***REMOVED***
		x.line("z, r := " + x.cpfx + "GenHelperEncoder(e)")
	***REMOVED*** else ***REMOVED***
		x.line("z, r := " + x.cpfx + "GenHelperDecoder(d)")
	***REMOVED***
	x.line("_, _, _ = h, z, r")
***REMOVED***

func (x *genRunner) genRefPkgs(t reflect.Type) ***REMOVED***
	if _, ok := x.is[t]; ok ***REMOVED***
		return
	***REMOVED***
	// fmt.Printf(">>>>>>: PkgPath: '%v', Name: '%s'\n", genImportPath(t), t.Name())
	x.is[t] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	tpkg, tname := genImportPath(t), t.Name()
	if tpkg != "" && tpkg != x.bp && tpkg != x.cp && tname != "" && tname[0] >= 'A' && tname[0] <= 'Z' ***REMOVED***
		if _, ok := x.im[tpkg]; !ok ***REMOVED***
			x.im[tpkg] = t
			if idx := strings.LastIndex(tpkg, "/"); idx < 0 ***REMOVED***
				x.imn[tpkg] = tpkg
			***REMOVED*** else ***REMOVED***
				x.imc++
				x.imn[tpkg] = "pkg" + strconv.FormatUint(x.imc, 10) + "_" + genGoIdentifier(tpkg[idx+1:], false)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Array, reflect.Slice, reflect.Ptr, reflect.Chan:
		x.genRefPkgs(t.Elem())
	case reflect.Map:
		x.genRefPkgs(t.Elem())
		x.genRefPkgs(t.Key())
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ ***REMOVED***
			if fname := t.Field(i).Name; fname != "" && fname[0] >= 'A' && fname[0] <= 'Z' ***REMOVED***
				x.genRefPkgs(t.Field(i).Type)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

func (x *genRunner) line(s string) ***REMOVED***
	x.out(s)
	if len(s) == 0 || s[len(s)-1] != '\n' ***REMOVED***
		x.out("\n")
	***REMOVED***
***REMOVED***

func (x *genRunner) varsfx() string ***REMOVED***
	x.c++
	return strconv.FormatUint(x.c, 10)
***REMOVED***

func (x *genRunner) varsfxreset() ***REMOVED***
	x.c = 0
***REMOVED***

func (x *genRunner) out(s string) ***REMOVED***
	if _, err := io.WriteString(x.w, s); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x *genRunner) linef(s string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	x.line(fmt.Sprintf(s, params...))
***REMOVED***

func (x *genRunner) outf(s string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	x.out(fmt.Sprintf(s, params...))
***REMOVED***

func (x *genRunner) genTypeName(t reflect.Type) (n string) ***REMOVED***
	// defer func() ***REMOVED*** fmt.Printf(">>>> ####: genTypeName: t: %v, name: '%s'\n", t, n) ***REMOVED***()

	// if the type has a PkgPath, which doesn't match the current package,
	// then include it.
	// We cannot depend on t.String() because it includes current package,
	// or t.PkgPath because it includes full import path,
	//
	var ptrPfx string
	for t.Kind() == reflect.Ptr ***REMOVED***
		ptrPfx += "*"
		t = t.Elem()
	***REMOVED***
	if tn := t.Name(); tn != "" ***REMOVED***
		return ptrPfx + x.genTypeNamePrim(t)
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Map:
		return ptrPfx + "map[" + x.genTypeName(t.Key()) + "]" + x.genTypeName(t.Elem())
	case reflect.Slice:
		return ptrPfx + "[]" + x.genTypeName(t.Elem())
	case reflect.Array:
		return ptrPfx + "[" + strconv.FormatInt(int64(t.Len()), 10) + "]" + x.genTypeName(t.Elem())
	case reflect.Chan:
		return ptrPfx + t.ChanDir().String() + " " + x.genTypeName(t.Elem())
	default:
		if t == intfTyp ***REMOVED***
			return ptrPfx + "interface***REMOVED******REMOVED***"
		***REMOVED*** else ***REMOVED***
			return ptrPfx + x.genTypeNamePrim(t)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (x *genRunner) genTypeNamePrim(t reflect.Type) (n string) ***REMOVED***
	if t.Name() == "" ***REMOVED***
		return t.String()
	***REMOVED*** else if genImportPath(t) == "" || genImportPath(t) == genImportPath(x.tc) ***REMOVED***
		return t.Name()
	***REMOVED*** else ***REMOVED***
		return x.imn[genImportPath(t)] + "." + t.Name()
		// return t.String() // best way to get the package name inclusive
	***REMOVED***
***REMOVED***

func (x *genRunner) genZeroValueR(t reflect.Type) string ***REMOVED***
	// if t is a named type, w
	switch t.Kind() ***REMOVED***
	case reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func,
		reflect.Slice, reflect.Map, reflect.Invalid:
		return "nil"
	case reflect.Bool:
		return "false"
	case reflect.String:
		return `""`
	case reflect.Struct, reflect.Array:
		return x.genTypeName(t) + "***REMOVED******REMOVED***"
	default: // all numbers
		return "0"
	***REMOVED***
***REMOVED***

func (x *genRunner) genMethodNameT(t reflect.Type) (s string) ***REMOVED***
	return genMethodNameT(t, x.tc)
***REMOVED***

func (x *genRunner) selfer(encode bool) ***REMOVED***
	t := x.tc
	t0 := t
	// always make decode use a pointer receiver,
	// and structs always use a ptr receiver (encode|decode)
	isptr := !encode || t.Kind() == reflect.Struct
	x.varsfxreset()
	fnSigPfx := "func (x "
	if isptr ***REMOVED***
		fnSigPfx += "*"
	***REMOVED***
	fnSigPfx += x.genTypeName(t)

	x.out(fnSigPfx)
	if isptr ***REMOVED***
		t = reflect.PtrTo(t)
	***REMOVED***
	if encode ***REMOVED***
		x.line(") CodecEncodeSelf(e *" + x.cpfx + "Encoder) ***REMOVED***")
		x.genRequiredMethodVars(true)
		// x.enc(genTopLevelVarName, t)
		x.encVar(genTopLevelVarName, t)
	***REMOVED*** else ***REMOVED***
		x.line(") CodecDecodeSelf(d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		// do not use decVar, as there is no need to check TryDecodeAsNil
		// or way to elegantly handle that, and also setting it to a
		// non-nil value doesn't affect the pointer passed.
		// x.decVar(genTopLevelVarName, t, false)
		x.dec(genTopLevelVarName, t0)
	***REMOVED***
	x.line("***REMOVED***")
	x.line("")

	if encode || t0.Kind() != reflect.Struct ***REMOVED***
		return
	***REMOVED***

	// write is containerMap
	if genUseOneFunctionForDecStructMap ***REMOVED***
		x.out(fnSigPfx)
		x.line(") codecDecodeSelfFromMap(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		x.decStructMap(genTopLevelVarName, "l", reflect.ValueOf(t0).Pointer(), t0, genStructMapStyleConsolidated)
		x.line("***REMOVED***")
		x.line("")
	***REMOVED*** else ***REMOVED***
		x.out(fnSigPfx)
		x.line(") codecDecodeSelfFromMapLenPrefix(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		x.decStructMap(genTopLevelVarName, "l", reflect.ValueOf(t0).Pointer(), t0, genStructMapStyleLenPrefix)
		x.line("***REMOVED***")
		x.line("")

		x.out(fnSigPfx)
		x.line(") codecDecodeSelfFromMapCheckBreak(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		x.decStructMap(genTopLevelVarName, "l", reflect.ValueOf(t0).Pointer(), t0, genStructMapStyleCheckBreak)
		x.line("***REMOVED***")
		x.line("")
	***REMOVED***

	// write containerArray
	x.out(fnSigPfx)
	x.line(") codecDecodeSelfFromArray(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
	x.genRequiredMethodVars(false)
	x.decStructArray(genTopLevelVarName, "l", "return", reflect.ValueOf(t0).Pointer(), t0)
	x.line("***REMOVED***")
	x.line("")

***REMOVED***

// used for chan, array, slice, map
func (x *genRunner) xtraSM(varname string, encode bool, t reflect.Type) ***REMOVED***
	if encode ***REMOVED***
		x.linef("h.enc%s((%s%s)(%s), e)", x.genMethodNameT(t), x.arr2str(t, "*"), x.genTypeName(t), varname)
	***REMOVED*** else ***REMOVED***
		x.linef("h.dec%s((*%s)(%s), d)", x.genMethodNameT(t), x.genTypeName(t), varname)
	***REMOVED***
	x.registerXtraT(t)
***REMOVED***

func (x *genRunner) registerXtraT(t reflect.Type) ***REMOVED***
	// recursively register the types
	if _, ok := x.tm[t]; ok ***REMOVED***
		return
	***REMOVED***
	var tkey reflect.Type
	switch t.Kind() ***REMOVED***
	case reflect.Chan, reflect.Slice, reflect.Array:
	case reflect.Map:
		tkey = t.Key()
	default:
		return
	***REMOVED***
	x.tm[t] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
	x.ts = append(x.ts, t)
	// check if this refers to any xtra types eg. a slice of array: add the array
	x.registerXtraT(t.Elem())
	if tkey != nil ***REMOVED***
		x.registerXtraT(tkey)
	***REMOVED***
***REMOVED***

// encVar will encode a variable.
// The parameter, t, is the reflect.Type of the variable itself
func (x *genRunner) encVar(varname string, t reflect.Type) ***REMOVED***
	// fmt.Printf(">>>>>> varname: %s, t: %v\n", varname, t)
	var checkNil bool
	switch t.Kind() ***REMOVED***
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
		checkNil = true
	***REMOVED***
	if checkNil ***REMOVED***
		x.linef("if %s == nil ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED*** ", varname)
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Ptr:
		switch t.Elem().Kind() ***REMOVED***
		case reflect.Struct, reflect.Array:
			x.enc(varname, genNonPtr(t))
		default:
			i := x.varsfx()
			x.line(genTempVarPfx + i + " := *" + varname)
			x.enc(genTempVarPfx+i, genNonPtr(t))
		***REMOVED***
	case reflect.Struct, reflect.Array:
		i := x.varsfx()
		x.line(genTempVarPfx + i + " := &" + varname)
		x.enc(genTempVarPfx+i, t)
	default:
		x.enc(varname, t)
	***REMOVED***

	if checkNil ***REMOVED***
		x.line("***REMOVED***")
	***REMOVED***

***REMOVED***

// enc will encode a variable (varname) of type t,
// except t is of kind reflect.Struct or reflect.Array, wherein varname is of type ptrTo(T) (to prevent copying)
func (x *genRunner) enc(varname string, t reflect.Type) ***REMOVED***
	rtid := reflect.ValueOf(t).Pointer()
	// We call CodecEncodeSelf if one of the following are honored:
	//   - the type already implements Selfer, call that
	//   - the type has a Selfer implementation just created, use that
	//   - the type is in the list of the ones we will generate for, but it is not currently being generated

	mi := x.varsfx()
	tptr := reflect.PtrTo(t)
	tk := t.Kind()
	if x.checkForSelfer(t, varname) ***REMOVED***
		if tk == reflect.Array || tk == reflect.Struct ***REMOVED*** // varname is of type *T
			if tptr.Implements(selferTyp) || t.Implements(selferTyp) ***REMOVED***
				x.line(varname + ".CodecEncodeSelf(e)")
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED*** // varname is of type T
			if t.Implements(selferTyp) ***REMOVED***
				x.line(varname + ".CodecEncodeSelf(e)")
				return
			***REMOVED*** else if tptr.Implements(selferTyp) ***REMOVED***
				x.linef("%ssf%s := &%s", genTempVarPfx, mi, varname)
				x.linef("%ssf%s.CodecEncodeSelf(e)", genTempVarPfx, mi)
				return
			***REMOVED***
		***REMOVED***

		if _, ok := x.te[rtid]; ok ***REMOVED***
			x.line(varname + ".CodecEncodeSelf(e)")
			return
		***REMOVED***
	***REMOVED***

	inlist := false
	for _, t0 := range x.t ***REMOVED***
		if t == t0 ***REMOVED***
			inlist = true
			if x.checkForSelfer(t, varname) ***REMOVED***
				x.line(varname + ".CodecEncodeSelf(e)")
				return
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	var rtidAdded bool
	if t == x.tc ***REMOVED***
		x.te[rtid] = true
		rtidAdded = true
	***REMOVED***

	// check if
	//   - type is RawExt, Raw
	//   - the type implements (Text|JSON|Binary)(Unm|M)arshal
	x.linef("%sm%s := z.EncBinary()", genTempVarPfx, mi)
	x.linef("_ = %sm%s", genTempVarPfx, mi)
	x.line("if false ***REMOVED***")           //start if block
	defer func() ***REMOVED*** x.line("***REMOVED***") ***REMOVED***() //end if block

	if t == rawTyp ***REMOVED***
		x.linef("***REMOVED*** else ***REMOVED*** z.EncRaw(%v)", varname)
		return
	***REMOVED***
	if t == rawExtTyp ***REMOVED***
		x.linef("***REMOVED*** else ***REMOVED*** r.EncodeRawExt(%v, e)", varname)
		return
	***REMOVED***
	// HACK: Support for Builtins.
	//       Currently, only Binc supports builtins, and the only builtin type is time.Time.
	//       Have a method that returns the rtid for time.Time if Handle is Binc.
	if t == timeTyp ***REMOVED***
		vrtid := genTempVarPfx + "m" + x.varsfx()
		x.linef("***REMOVED*** else if %s := z.TimeRtidIfBinc(); %s != 0 ***REMOVED*** ", vrtid, vrtid)
		x.linef("r.EncodeBuiltin(%s, %s)", vrtid, varname)
	***REMOVED***
	// only check for extensions if the type is named, and has a packagePath.
	if genImportPath(t) != "" && t.Name() != "" ***REMOVED***
		// first check if extensions are configued, before doing the interface conversion
		x.linef("***REMOVED*** else if z.HasExtensions() && z.EncExt(%s) ***REMOVED***", varname)
	***REMOVED***
	if tk == reflect.Array || tk == reflect.Struct ***REMOVED*** // varname is of type *T
		if t.Implements(binaryMarshalerTyp) || tptr.Implements(binaryMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if %sm%s ***REMOVED*** z.EncBinaryMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED***
		if t.Implements(jsonMarshalerTyp) || tptr.Implements(jsonMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED*** else if t.Implements(textMarshalerTyp) || tptr.Implements(textMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s ***REMOVED*** z.EncTextMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // varname is of type T
		if t.Implements(binaryMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if %sm%s ***REMOVED*** z.EncBinaryMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED*** else if tptr.Implements(binaryMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if %sm%s ***REMOVED*** z.EncBinaryMarshal(&%v) ", genTempVarPfx, mi, varname)
		***REMOVED***
		if t.Implements(jsonMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED*** else if tptr.Implements(jsonMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(&%v) ", genTempVarPfx, mi, varname)
		***REMOVED*** else if t.Implements(textMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s ***REMOVED*** z.EncTextMarshal(%v) ", genTempVarPfx, mi, varname)
		***REMOVED*** else if tptr.Implements(textMarshalerTyp) ***REMOVED***
			x.linef("***REMOVED*** else if !%sm%s ***REMOVED*** z.EncTextMarshal(&%v) ", genTempVarPfx, mi, varname)
		***REMOVED***
	***REMOVED***
	x.line("***REMOVED*** else ***REMOVED***")

	switch t.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.line("r.EncodeInt(int64(" + varname + "))")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.line("r.EncodeUint(uint64(" + varname + "))")
	case reflect.Float32:
		x.line("r.EncodeFloat32(float32(" + varname + "))")
	case reflect.Float64:
		x.line("r.EncodeFloat64(float64(" + varname + "))")
	case reflect.Bool:
		x.line("r.EncodeBool(bool(" + varname + "))")
	case reflect.String:
		x.line("r.EncodeString(codecSelferC_UTF8" + x.xs + ", string(" + varname + "))")
	case reflect.Chan:
		x.xtraSM(varname, true, t)
		// x.encListFallback(varname, rtid, t)
	case reflect.Array:
		x.xtraSM(varname, true, t)
	case reflect.Slice:
		// if nil, call dedicated function
		// if a []uint8, call dedicated function
		// if a known fastpath slice, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		if rtid == uint8SliceTypId ***REMOVED***
			x.line("r.EncodeStringBytes(codecSelferC_RAW" + x.xs + ", []byte(" + varname + "))")
		***REMOVED*** else if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newGenV(t)
			x.line("z.F." + g.MethodNamePfx("Enc", false) + "V(" + varname + ", false, e)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, true, t)
			// x.encListFallback(varname, rtid, t)
		***REMOVED***
	case reflect.Map:
		// if nil, call dedicated function
		// if a known fastpath map, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		// x.line("if " + varname + " == nil ***REMOVED*** \nr.EncodeNil()\n ***REMOVED*** else ***REMOVED*** ")
		if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newGenV(t)
			x.line("z.F." + g.MethodNamePfx("Enc", false) + "V(" + varname + ", false, e)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, true, t)
			// x.encMapFallback(varname, rtid, t)
		***REMOVED***
	case reflect.Struct:
		if !inlist ***REMOVED***
			delete(x.te, rtid)
			x.line("z.EncFallback(" + varname + ")")
			break
		***REMOVED***
		x.encStruct(varname, rtid, t)
	default:
		if rtidAdded ***REMOVED***
			delete(x.te, rtid)
		***REMOVED***
		x.line("z.EncFallback(" + varname + ")")
	***REMOVED***
***REMOVED***

func (x *genRunner) encZero(t reflect.Type) ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		x.line("r.EncodeInt(0)")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		x.line("r.EncodeUint(0)")
	case reflect.Float32:
		x.line("r.EncodeFloat32(0)")
	case reflect.Float64:
		x.line("r.EncodeFloat64(0)")
	case reflect.Bool:
		x.line("r.EncodeBool(false)")
	case reflect.String:
		x.line("r.EncodeString(codecSelferC_UTF8" + x.xs + `, "")`)
	default:
		x.line("r.EncodeNil()")
	***REMOVED***
***REMOVED***

func (x *genRunner) encStruct(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	// Use knowledge from structfieldinfo (mbs, encodable fields. Ignore omitempty. )
	// replicate code in kStruct i.e. for each field, deref type to non-pointer, and call x.enc on it

	// if t === type currently running selfer on, do for all
	ti := x.ti.get(rtid, t)
	i := x.varsfx()
	sepVarname := genTempVarPfx + "sep" + i
	numfieldsvar := genTempVarPfx + "q" + i
	ti2arrayvar := genTempVarPfx + "r" + i
	struct2arrvar := genTempVarPfx + "2arr" + i

	x.line(sepVarname + " := !z.EncBinary()")
	x.linef("%s := z.EncBasicHandle().StructToArray", struct2arrvar)
	tisfi := ti.sfip // always use sequence from file. decStruct expects same thing.
	// due to omitEmpty, we need to calculate the
	// number of non-empty things we write out first.
	// This is required as we need to pre-determine the size of the container,
	// to support length-prefixing.
	x.linef("var %s [%v]bool", numfieldsvar, len(tisfi))
	x.linef("_, _, _ = %s, %s, %s", sepVarname, numfieldsvar, struct2arrvar)
	x.linef("const %s bool = %v", ti2arrayvar, ti.toArray)
	nn := 0
	for j, si := range tisfi ***REMOVED***
		if !si.omitEmpty ***REMOVED***
			nn++
			continue
		***REMOVED***
		var t2 reflect.StructField
		var omitline string
		if si.i != -1 ***REMOVED***
			t2 = t.Field(int(si.i))
		***REMOVED*** else ***REMOVED***
			t2typ := t
			varname3 := varname
			for _, ix := range si.is ***REMOVED***
				for t2typ.Kind() == reflect.Ptr ***REMOVED***
					t2typ = t2typ.Elem()
				***REMOVED***
				t2 = t2typ.Field(ix)
				t2typ = t2.Type
				varname3 = varname3 + "." + t2.Name
				if t2typ.Kind() == reflect.Ptr ***REMOVED***
					omitline += varname3 + " != nil && "
				***REMOVED***
			***REMOVED***
		***REMOVED***
		// never check omitEmpty on a struct type, as it may contain uncomparable map/slice/etc.
		// also, for maps/slices/arrays, check if len ! 0 (not if == zero value)
		switch t2.Type.Kind() ***REMOVED***
		case reflect.Struct:
			omitline += " true"
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Chan:
			omitline += "len(" + varname + "." + t2.Name + ") != 0"
		default:
			omitline += varname + "." + t2.Name + " != " + x.genZeroValueR(t2.Type)
		***REMOVED***
		x.linef("%s[%v] = %s", numfieldsvar, j, omitline)
	***REMOVED***
	x.linef("var %snn%s int", genTempVarPfx, i)
	x.linef("if %s || %s ***REMOVED***", ti2arrayvar, struct2arrvar) // if ti.toArray ***REMOVED***
	x.line("r.EncodeArrayStart(" + strconv.FormatInt(int64(len(tisfi)), 10) + ")")
	x.linef("***REMOVED*** else ***REMOVED***") // if not ti.toArray
	x.linef("%snn%s = %v", genTempVarPfx, i, nn)
	x.linef("for _, b := range %s ***REMOVED*** if b ***REMOVED*** %snn%s++ ***REMOVED*** ***REMOVED***", numfieldsvar, genTempVarPfx, i)
	x.linef("r.EncodeMapStart(%snn%s)", genTempVarPfx, i)
	x.linef("%snn%s = %v", genTempVarPfx, i, 0)
	// x.line("r.EncodeMapStart(" + strconv.FormatInt(int64(len(tisfi)), 10) + ")")
	x.line("***REMOVED***") // close if not StructToArray

	for j, si := range tisfi ***REMOVED***
		i := x.varsfx()
		isNilVarName := genTempVarPfx + "n" + i
		var labelUsed bool
		var t2 reflect.StructField
		if si.i != -1 ***REMOVED***
			t2 = t.Field(int(si.i))
		***REMOVED*** else ***REMOVED***
			t2typ := t
			varname3 := varname
			for _, ix := range si.is ***REMOVED***
				// fmt.Printf("%%%% %v, ix: %v\n", t2typ, ix)
				for t2typ.Kind() == reflect.Ptr ***REMOVED***
					t2typ = t2typ.Elem()
				***REMOVED***
				t2 = t2typ.Field(ix)
				t2typ = t2.Type
				varname3 = varname3 + "." + t2.Name
				if t2typ.Kind() == reflect.Ptr ***REMOVED***
					if !labelUsed ***REMOVED***
						x.line("var " + isNilVarName + " bool")
					***REMOVED***
					x.line("if " + varname3 + " == nil ***REMOVED*** " + isNilVarName + " = true ")
					x.line("goto LABEL" + i)
					x.line("***REMOVED***")
					labelUsed = true
					// "varname3 = new(" + x.genTypeName(t3.Elem()) + ") ***REMOVED***")
				***REMOVED***
			***REMOVED***
			// t2 = t.FieldByIndex(si.is)
		***REMOVED***
		if labelUsed ***REMOVED***
			x.line("LABEL" + i + ":")
		***REMOVED***
		// if the type of the field is a Selfer, or one of the ones

		x.linef("if %s || %s ***REMOVED***", ti2arrayvar, struct2arrvar) // if ti.toArray
		if labelUsed ***REMOVED***
			x.line("if " + isNilVarName + " ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED*** ")
		***REMOVED***
		x.linef("z.EncSendContainerState(codecSelfer_containerArrayElem%s)", x.xs)
		if si.omitEmpty ***REMOVED***
			x.linef("if %s[%v] ***REMOVED***", numfieldsvar, j)
		***REMOVED***
		x.encVar(varname+"."+t2.Name, t2.Type)
		if si.omitEmpty ***REMOVED***
			x.linef("***REMOVED*** else ***REMOVED***")
			x.encZero(t2.Type)
			x.linef("***REMOVED***")
		***REMOVED***
		if labelUsed ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***

		x.linef("***REMOVED*** else ***REMOVED***") // if not ti.toArray

		if si.omitEmpty ***REMOVED***
			x.linef("if %s[%v] ***REMOVED***", numfieldsvar, j)
		***REMOVED***
		x.linef("z.EncSendContainerState(codecSelfer_containerMapKey%s)", x.xs)
		x.line("r.EncodeString(codecSelferC_UTF8" + x.xs + ", string(\"" + si.encName + "\"))")
		x.linef("z.EncSendContainerState(codecSelfer_containerMapValue%s)", x.xs)
		if labelUsed ***REMOVED***
			x.line("if " + isNilVarName + " ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED*** ")
			x.encVar(varname+"."+t2.Name, t2.Type)
			x.line("***REMOVED***")
		***REMOVED*** else ***REMOVED***
			x.encVar(varname+"."+t2.Name, t2.Type)
		***REMOVED***
		if si.omitEmpty ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***
		x.linef("***REMOVED*** ") // end if/else ti.toArray
	***REMOVED***
	x.linef("if %s || %s ***REMOVED***", ti2arrayvar, struct2arrvar) // if ti.toArray ***REMOVED***
	x.linef("z.EncSendContainerState(codecSelfer_containerArrayEnd%s)", x.xs)
	x.line("***REMOVED*** else ***REMOVED***")
	x.linef("z.EncSendContainerState(codecSelfer_containerMapEnd%s)", x.xs)
	x.line("***REMOVED***")

***REMOVED***

func (x *genRunner) encListFallback(varname string, t reflect.Type) ***REMOVED***
	if t.AssignableTo(uint8SliceTyp) ***REMOVED***
		x.linef("r.EncodeStringBytes(codecSelferC_RAW%s, []byte(%s))", x.xs, varname)
		return
	***REMOVED***
	if t.Kind() == reflect.Array && t.Elem().Kind() == reflect.Uint8 ***REMOVED***
		x.linef("r.EncodeStringBytes(codecSelferC_RAW%s, ([%v]byte(%s))[:])", x.xs, t.Len(), varname)
		return
	***REMOVED***
	i := x.varsfx()
	g := genTempVarPfx
	x.line("r.EncodeArrayStart(len(" + varname + "))")
	if t.Kind() == reflect.Chan ***REMOVED***
		x.linef("for %si%s, %si2%s := 0, len(%s); %si%s < %si2%s; %si%s++ ***REMOVED***", g, i, g, i, varname, g, i, g, i, g, i)
		x.linef("z.EncSendContainerState(codecSelfer_containerArrayElem%s)", x.xs)
		x.linef("%sv%s := <-%s", g, i, varname)
	***REMOVED*** else ***REMOVED***
		// x.linef("for %si%s, %sv%s := range %s ***REMOVED***", genTempVarPfx, i, genTempVarPfx, i, varname)
		x.linef("for _, %sv%s := range %s ***REMOVED***", genTempVarPfx, i, varname)
		x.linef("z.EncSendContainerState(codecSelfer_containerArrayElem%s)", x.xs)
	***REMOVED***
	x.encVar(genTempVarPfx+"v"+i, t.Elem())
	x.line("***REMOVED***")
	x.linef("z.EncSendContainerState(codecSelfer_containerArrayEnd%s)", x.xs)
***REMOVED***

func (x *genRunner) encMapFallback(varname string, t reflect.Type) ***REMOVED***
	// TODO: expand this to handle canonical.
	i := x.varsfx()
	x.line("r.EncodeMapStart(len(" + varname + "))")
	x.linef("for %sk%s, %sv%s := range %s ***REMOVED***", genTempVarPfx, i, genTempVarPfx, i, varname)
	// x.line("for " + genTempVarPfx + "k" + i + ", " + genTempVarPfx + "v" + i + " := range " + varname + " ***REMOVED***")
	x.linef("z.EncSendContainerState(codecSelfer_containerMapKey%s)", x.xs)
	x.encVar(genTempVarPfx+"k"+i, t.Key())
	x.linef("z.EncSendContainerState(codecSelfer_containerMapValue%s)", x.xs)
	x.encVar(genTempVarPfx+"v"+i, t.Elem())
	x.line("***REMOVED***")
	x.linef("z.EncSendContainerState(codecSelfer_containerMapEnd%s)", x.xs)
***REMOVED***

func (x *genRunner) decVar(varname string, t reflect.Type, canBeNil bool) ***REMOVED***
	// We only encode as nil if a nillable value.
	// This removes some of the wasted checks for TryDecodeAsNil.
	// We need to think about this more, to see what happens if omitempty, etc
	// cause a nil value to be stored when something is expected.
	// This could happen when decoding from a struct encoded as an array.
	// For that, decVar should be called with canNil=true, to force true as its value.
	i := x.varsfx()
	if !canBeNil ***REMOVED***
		canBeNil = genAnythingCanBeNil || !genIsImmutable(t)
	***REMOVED***
	if canBeNil ***REMOVED***
		x.line("if r.TryDecodeAsNil() ***REMOVED***")
		if t.Kind() == reflect.Ptr ***REMOVED***
			x.line("if " + varname + " != nil ***REMOVED*** ")

			// if varname is a field of a struct (has a dot in it),
			// then just set it to nil
			if strings.IndexByte(varname, '.') != -1 ***REMOVED***
				x.line(varname + " = nil")
			***REMOVED*** else ***REMOVED***
				x.line("*" + varname + " = " + x.genZeroValueR(t.Elem()))
			***REMOVED***
			x.line("***REMOVED***")
		***REMOVED*** else ***REMOVED***
			x.line(varname + " = " + x.genZeroValueR(t))
		***REMOVED***
		x.line("***REMOVED*** else ***REMOVED***")
	***REMOVED*** else ***REMOVED***
		x.line("// cannot be nil")
	***REMOVED***
	if t.Kind() != reflect.Ptr ***REMOVED***
		if x.decTryAssignPrimitive(varname, t) ***REMOVED***
			x.line(genTempVarPfx + "v" + i + " := &" + varname)
			x.dec(genTempVarPfx+"v"+i, t)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		x.linef("if %s == nil ***REMOVED*** %s = new(%s) ***REMOVED***", varname, varname, x.genTypeName(t.Elem()))
		// Ensure we set underlying ptr to a non-nil value (so we can deref to it later).
		// There's a chance of a **T in here which is nil.
		var ptrPfx string
		for t = t.Elem(); t.Kind() == reflect.Ptr; t = t.Elem() ***REMOVED***
			ptrPfx += "*"
			x.linef("if %s%s == nil ***REMOVED*** %s%s = new(%s)***REMOVED***",
				ptrPfx, varname, ptrPfx, varname, x.genTypeName(t))
		***REMOVED***
		// if varname has [ in it, then create temp variable for this ptr thingie
		if strings.Index(varname, "[") >= 0 ***REMOVED***
			varname2 := genTempVarPfx + "w" + i
			x.line(varname2 + " := " + varname)
			varname = varname2
		***REMOVED***

		if ptrPfx == "" ***REMOVED***
			x.dec(varname, t)
		***REMOVED*** else ***REMOVED***
			x.line(genTempVarPfx + "z" + i + " := " + ptrPfx + varname)
			x.dec(genTempVarPfx+"z"+i, t)
		***REMOVED***

	***REMOVED***

	if canBeNil ***REMOVED***
		x.line("***REMOVED*** ")
	***REMOVED***
***REMOVED***

// dec will decode a variable (varname) of type ptrTo(t).
// t is always a basetype (i.e. not of kind reflect.Ptr).
func (x *genRunner) dec(varname string, t reflect.Type) ***REMOVED***
	// assumptions:
	//   - the varname is to a pointer already. No need to take address of it
	//   - t is always a baseType T (not a *T, etc).
	rtid := reflect.ValueOf(t).Pointer()
	tptr := reflect.PtrTo(t)
	if x.checkForSelfer(t, varname) ***REMOVED***
		if t.Implements(selferTyp) || tptr.Implements(selferTyp) ***REMOVED***
			x.line(varname + ".CodecDecodeSelf(d)")
			return
		***REMOVED***
		if _, ok := x.td[rtid]; ok ***REMOVED***
			x.line(varname + ".CodecDecodeSelf(d)")
			return
		***REMOVED***
	***REMOVED***

	inlist := false
	for _, t0 := range x.t ***REMOVED***
		if t == t0 ***REMOVED***
			inlist = true
			if x.checkForSelfer(t, varname) ***REMOVED***
				x.line(varname + ".CodecDecodeSelf(d)")
				return
			***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	var rtidAdded bool
	if t == x.tc ***REMOVED***
		x.td[rtid] = true
		rtidAdded = true
	***REMOVED***

	// check if
	//   - type is Raw, RawExt
	//   - the type implements (Text|JSON|Binary)(Unm|M)arshal
	mi := x.varsfx()
	x.linef("%sm%s := z.DecBinary()", genTempVarPfx, mi)
	x.linef("_ = %sm%s", genTempVarPfx, mi)
	x.line("if false ***REMOVED***")           //start if block
	defer func() ***REMOVED*** x.line("***REMOVED***") ***REMOVED***() //end if block

	if t == rawTyp ***REMOVED***
		x.linef("***REMOVED*** else ***REMOVED*** *%v = z.DecRaw()", varname)
		return
	***REMOVED***
	if t == rawExtTyp ***REMOVED***
		x.linef("***REMOVED*** else ***REMOVED*** r.DecodeExt(%v, 0, nil)", varname)
		return
	***REMOVED***

	// HACK: Support for Builtins.
	//       Currently, only Binc supports builtins, and the only builtin type is time.Time.
	//       Have a method that returns the rtid for time.Time if Handle is Binc.
	if t == timeTyp ***REMOVED***
		vrtid := genTempVarPfx + "m" + x.varsfx()
		x.linef("***REMOVED*** else if %s := z.TimeRtidIfBinc(); %s != 0 ***REMOVED*** ", vrtid, vrtid)
		x.linef("r.DecodeBuiltin(%s, %s)", vrtid, varname)
	***REMOVED***
	// only check for extensions if the type is named, and has a packagePath.
	if genImportPath(t) != "" && t.Name() != "" ***REMOVED***
		// first check if extensions are configued, before doing the interface conversion
		x.linef("***REMOVED*** else if z.HasExtensions() && z.DecExt(%s) ***REMOVED***", varname)
	***REMOVED***

	if t.Implements(binaryUnmarshalerTyp) || tptr.Implements(binaryUnmarshalerTyp) ***REMOVED***
		x.linef("***REMOVED*** else if %sm%s ***REMOVED*** z.DecBinaryUnmarshal(%v) ", genTempVarPfx, mi, varname)
	***REMOVED***
	if t.Implements(jsonUnmarshalerTyp) || tptr.Implements(jsonUnmarshalerTyp) ***REMOVED***
		x.linef("***REMOVED*** else if !%sm%s && z.IsJSONHandle() ***REMOVED*** z.DecJSONUnmarshal(%v)", genTempVarPfx, mi, varname)
	***REMOVED*** else if t.Implements(textUnmarshalerTyp) || tptr.Implements(textUnmarshalerTyp) ***REMOVED***
		x.linef("***REMOVED*** else if !%sm%s ***REMOVED*** z.DecTextUnmarshal(%v)", genTempVarPfx, mi, varname)
	***REMOVED***

	x.line("***REMOVED*** else ***REMOVED***")

	// Since these are pointers, we cannot share, and have to use them one by one
	switch t.Kind() ***REMOVED***
	case reflect.Int:
		x.line("*((*int)(" + varname + ")) = int(r.DecodeInt(codecSelferBitsize" + x.xs + "))")
		// x.line("z.DecInt((*int)(" + varname + "))")
	case reflect.Int8:
		x.line("*((*int8)(" + varname + ")) = int8(r.DecodeInt(8))")
		// x.line("z.DecInt8((*int8)(" + varname + "))")
	case reflect.Int16:
		x.line("*((*int16)(" + varname + ")) = int16(r.DecodeInt(16))")
		// x.line("z.DecInt16((*int16)(" + varname + "))")
	case reflect.Int32:
		x.line("*((*int32)(" + varname + ")) = int32(r.DecodeInt(32))")
		// x.line("z.DecInt32((*int32)(" + varname + "))")
	case reflect.Int64:
		x.line("*((*int64)(" + varname + ")) = int64(r.DecodeInt(64))")
		// x.line("z.DecInt64((*int64)(" + varname + "))")

	case reflect.Uint:
		x.line("*((*uint)(" + varname + ")) = uint(r.DecodeUint(codecSelferBitsize" + x.xs + "))")
		// x.line("z.DecUint((*uint)(" + varname + "))")
	case reflect.Uint8:
		x.line("*((*uint8)(" + varname + ")) = uint8(r.DecodeUint(8))")
		// x.line("z.DecUint8((*uint8)(" + varname + "))")
	case reflect.Uint16:
		x.line("*((*uint16)(" + varname + ")) = uint16(r.DecodeUint(16))")
		//x.line("z.DecUint16((*uint16)(" + varname + "))")
	case reflect.Uint32:
		x.line("*((*uint32)(" + varname + ")) = uint32(r.DecodeUint(32))")
		//x.line("z.DecUint32((*uint32)(" + varname + "))")
	case reflect.Uint64:
		x.line("*((*uint64)(" + varname + ")) = uint64(r.DecodeUint(64))")
		//x.line("z.DecUint64((*uint64)(" + varname + "))")
	case reflect.Uintptr:
		x.line("*((*uintptr)(" + varname + ")) = uintptr(r.DecodeUint(codecSelferBitsize" + x.xs + "))")

	case reflect.Float32:
		x.line("*((*float32)(" + varname + ")) = float32(r.DecodeFloat(true))")
		//x.line("z.DecFloat32((*float32)(" + varname + "))")
	case reflect.Float64:
		x.line("*((*float64)(" + varname + ")) = float64(r.DecodeFloat(false))")
		// x.line("z.DecFloat64((*float64)(" + varname + "))")

	case reflect.Bool:
		x.line("*((*bool)(" + varname + ")) = r.DecodeBool()")
		// x.line("z.DecBool((*bool)(" + varname + "))")
	case reflect.String:
		x.line("*((*string)(" + varname + ")) = r.DecodeString()")
		// x.line("z.DecString((*string)(" + varname + "))")
	case reflect.Array, reflect.Chan:
		x.xtraSM(varname, false, t)
		// x.decListFallback(varname, rtid, true, t)
	case reflect.Slice:
		// if a []uint8, call dedicated function
		// if a known fastpath slice, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		if rtid == uint8SliceTypId ***REMOVED***
			x.line("*" + varname + " = r.DecodeBytes(*(*[]byte)(" + varname + "), false, false)")
		***REMOVED*** else if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newGenV(t)
			x.line("z.F." + g.MethodNamePfx("Dec", false) + "X(" + varname + ", false, d)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, false, t)
			// x.decListFallback(varname, rtid, false, t)
		***REMOVED***
	case reflect.Map:
		// if a known fastpath map, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newGenV(t)
			x.line("z.F." + g.MethodNamePfx("Dec", false) + "X(" + varname + ", false, d)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, false, t)
			// x.decMapFallback(varname, rtid, t)
		***REMOVED***
	case reflect.Struct:
		if inlist ***REMOVED***
			x.decStruct(varname, rtid, t)
		***REMOVED*** else ***REMOVED***
			// delete(x.td, rtid)
			x.line("z.DecFallback(" + varname + ", false)")
		***REMOVED***
	default:
		if rtidAdded ***REMOVED***
			delete(x.te, rtid)
		***REMOVED***
		x.line("z.DecFallback(" + varname + ", true)")
	***REMOVED***
***REMOVED***

func (x *genRunner) decTryAssignPrimitive(varname string, t reflect.Type) (tryAsPtr bool) ***REMOVED***
	// This should only be used for exact primitives (ie un-named types).
	// Named types may be implementations of Selfer, Unmarshaler, etc.
	// They should be handled by dec(...)

	if t.Name() != "" ***REMOVED***
		tryAsPtr = true
		return
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Int:
		x.linef("%s = r.DecodeInt(codecSelferBitsize%s)", varname, x.xs)
	case reflect.Int8:
		x.linef("%s = r.DecodeInt(8)", varname)
	case reflect.Int16:
		x.linef("%s = r.DecodeInt(16)", varname)
	case reflect.Int32:
		x.linef("%s = r.DecodeInt(32)", varname)
	case reflect.Int64:
		x.linef("%s = r.DecodeInt(64)", varname)

	case reflect.Uint:
		x.linef("%s = r.DecodeUint(codecSelferBitsize%s)", varname, x.xs)
	case reflect.Uint8:
		x.linef("%s = r.DecodeUint(8)", varname)
	case reflect.Uint16:
		x.linef("%s = r.DecodeUint(16)", varname)
	case reflect.Uint32:
		x.linef("%s = r.DecodeUint(32)", varname)
	case reflect.Uint64:
		x.linef("%s = r.DecodeUint(64)", varname)
	case reflect.Uintptr:
		x.linef("%s = r.DecodeUint(codecSelferBitsize%s)", varname, x.xs)

	case reflect.Float32:
		x.linef("%s = r.DecodeFloat(true)", varname)
	case reflect.Float64:
		x.linef("%s = r.DecodeFloat(false)", varname)

	case reflect.Bool:
		x.linef("%s = r.DecodeBool()", varname)
	case reflect.String:
		x.linef("%s = r.DecodeString()", varname)
	default:
		tryAsPtr = true
	***REMOVED***
	return
***REMOVED***

func (x *genRunner) decListFallback(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	if t.AssignableTo(uint8SliceTyp) ***REMOVED***
		x.line("*" + varname + " = r.DecodeBytes(*((*[]byte)(" + varname + ")), false, false)")
		return
	***REMOVED***
	if t.Kind() == reflect.Array && t.Elem().Kind() == reflect.Uint8 ***REMOVED***
		x.linef("r.DecodeBytes( ((*[%s]byte)(%s))[:], false, true)", t.Len(), varname)
		return
	***REMOVED***
	type tstruc struct ***REMOVED***
		TempVar   string
		Rand      string
		Varname   string
		CTyp      string
		Typ       string
		Immutable bool
		Size      int
	***REMOVED***
	telem := t.Elem()
	ts := tstruc***REMOVED***genTempVarPfx, x.varsfx(), varname, x.genTypeName(t), x.genTypeName(telem), genIsImmutable(telem), int(telem.Size())***REMOVED***

	funcs := make(template.FuncMap)

	funcs["decLineVar"] = func(varname string) string ***REMOVED***
		x.decVar(varname, telem, false)
		return ""
	***REMOVED***
	funcs["decLine"] = func(pfx string) string ***REMOVED***
		x.decVar(ts.TempVar+pfx+ts.Rand, reflect.PtrTo(telem), false)
		return ""
	***REMOVED***
	funcs["var"] = func(s string) string ***REMOVED***
		return ts.TempVar + s + ts.Rand
	***REMOVED***
	funcs["zero"] = func() string ***REMOVED***
		return x.genZeroValueR(telem)
	***REMOVED***
	funcs["isArray"] = func() bool ***REMOVED***
		return t.Kind() == reflect.Array
	***REMOVED***
	funcs["isSlice"] = func() bool ***REMOVED***
		return t.Kind() == reflect.Slice
	***REMOVED***
	funcs["isChan"] = func() bool ***REMOVED***
		return t.Kind() == reflect.Chan
	***REMOVED***
	tm, err := template.New("").Funcs(funcs).Parse(genDecListTmpl)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if err = tm.Execute(x.w, &ts); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x *genRunner) decMapFallback(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	type tstruc struct ***REMOVED***
		TempVar string
		Sfx     string
		Rand    string
		Varname string
		KTyp    string
		Typ     string
		Size    int
	***REMOVED***
	telem := t.Elem()
	tkey := t.Key()
	ts := tstruc***REMOVED***
		genTempVarPfx, x.xs, x.varsfx(), varname, x.genTypeName(tkey),
		x.genTypeName(telem), int(telem.Size() + tkey.Size()),
	***REMOVED***

	funcs := make(template.FuncMap)
	funcs["decElemZero"] = func() string ***REMOVED***
		return x.genZeroValueR(telem)
	***REMOVED***
	funcs["decElemKindImmutable"] = func() bool ***REMOVED***
		return genIsImmutable(telem)
	***REMOVED***
	funcs["decElemKindPtr"] = func() bool ***REMOVED***
		return telem.Kind() == reflect.Ptr
	***REMOVED***
	funcs["decElemKindIntf"] = func() bool ***REMOVED***
		return telem.Kind() == reflect.Interface
	***REMOVED***
	funcs["decLineVarK"] = func(varname string) string ***REMOVED***
		x.decVar(varname, tkey, false)
		return ""
	***REMOVED***
	funcs["decLineVar"] = func(varname string) string ***REMOVED***
		x.decVar(varname, telem, false)
		return ""
	***REMOVED***
	funcs["decLineK"] = func(pfx string) string ***REMOVED***
		x.decVar(ts.TempVar+pfx+ts.Rand, reflect.PtrTo(tkey), false)
		return ""
	***REMOVED***
	funcs["decLine"] = func(pfx string) string ***REMOVED***
		x.decVar(ts.TempVar+pfx+ts.Rand, reflect.PtrTo(telem), false)
		return ""
	***REMOVED***
	funcs["var"] = func(s string) string ***REMOVED***
		return ts.TempVar + s + ts.Rand
	***REMOVED***

	tm, err := template.New("").Funcs(funcs).Parse(genDecMapTmpl)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	if err = tm.Execute(x.w, &ts); err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x *genRunner) decStructMapSwitch(kName string, varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	ti := x.ti.get(rtid, t)
	tisfi := ti.sfip // always use sequence from file. decStruct expects same thing.
	x.line("switch (" + kName + ") ***REMOVED***")
	for _, si := range tisfi ***REMOVED***
		x.line("case \"" + si.encName + "\":")
		var t2 reflect.StructField
		if si.i != -1 ***REMOVED***
			t2 = t.Field(int(si.i))
		***REMOVED*** else ***REMOVED***
			//we must accommodate anonymous fields, where the embedded field is a nil pointer in the value.
			// t2 = t.FieldByIndex(si.is)
			t2typ := t
			varname3 := varname
			for _, ix := range si.is ***REMOVED***
				for t2typ.Kind() == reflect.Ptr ***REMOVED***
					t2typ = t2typ.Elem()
				***REMOVED***
				t2 = t2typ.Field(ix)
				t2typ = t2.Type
				varname3 = varname3 + "." + t2.Name
				if t2typ.Kind() == reflect.Ptr ***REMOVED***
					x.linef("if %s == nil ***REMOVED*** %s = new(%s) ***REMOVED***", varname3, varname3, x.genTypeName(t2typ.Elem()))
				***REMOVED***
			***REMOVED***
		***REMOVED***
		x.decVar(varname+"."+t2.Name, t2.Type, false)
	***REMOVED***
	x.line("default:")
	// pass the slice here, so that the string will not escape, and maybe save allocation
	x.line("z.DecStructFieldNotFound(-1, " + kName + ")")
	x.line("***REMOVED*** // end switch " + kName)
***REMOVED***

func (x *genRunner) decStructMap(varname, lenvarname string, rtid uintptr, t reflect.Type, style genStructMapStyle) ***REMOVED***
	tpfx := genTempVarPfx
	i := x.varsfx()
	kName := tpfx + "s" + i

	// We thought to use ReadStringAsBytes, as go compiler might optimize the copy out.
	// However, using that was more expensive, as it seems that the switch expression
	// is evaluated each time.
	//
	// We could depend on decodeString using a temporary/shared buffer internally.
	// However, this model of creating a byte array, and using explicitly is faster,
	// and allows optional use of unsafe []byte->string conversion without alloc.

	// Also, ensure that the slice array doesn't escape.
	// That will help escape analysis prevent allocation when it gets better.

	// x.line("var " + kName + "Arr = [32]byte***REMOVED******REMOVED*** // default string to decode into")
	// x.line("var " + kName + "Slc = " + kName + "Arr[:] // default slice to decode into")
	// use the scratch buffer to avoid allocation (most field names are < 32).

	x.line("var " + kName + "Slc = z.DecScratchBuffer() // default slice to decode into")

	x.line("_ = " + kName + "Slc")
	switch style ***REMOVED***
	case genStructMapStyleLenPrefix:
		x.linef("for %sj%s := 0; %sj%s < %s; %sj%s++ ***REMOVED***", tpfx, i, tpfx, i, lenvarname, tpfx, i)
	case genStructMapStyleCheckBreak:
		x.linef("for %sj%s := 0; !r.CheckBreak(); %sj%s++ ***REMOVED***", tpfx, i, tpfx, i)
	default: // 0, otherwise.
		x.linef("var %shl%s bool = %s >= 0", tpfx, i, lenvarname) // has length
		x.linef("for %sj%s := 0; ; %sj%s++ ***REMOVED***", tpfx, i, tpfx, i)
		x.linef("if %shl%s ***REMOVED*** if %sj%s >= %s ***REMOVED*** break ***REMOVED***", tpfx, i, tpfx, i, lenvarname)
		x.line("***REMOVED*** else ***REMOVED*** if r.CheckBreak() ***REMOVED*** break ***REMOVED***; ***REMOVED***")
	***REMOVED***
	x.linef("z.DecSendContainerState(codecSelfer_containerMapKey%s)", x.xs)
	x.line(kName + "Slc = r.DecodeBytes(" + kName + "Slc, true, true)")
	// let string be scoped to this loop alone, so it doesn't escape.
	if x.unsafe ***REMOVED***
		x.line(kName + "SlcHdr := codecSelferUnsafeString" + x.xs + "***REMOVED***uintptr(unsafe.Pointer(&" +
			kName + "Slc[0])), len(" + kName + "Slc)***REMOVED***")
		x.line(kName + " := *(*string)(unsafe.Pointer(&" + kName + "SlcHdr))")
	***REMOVED*** else ***REMOVED***
		x.line(kName + " := string(" + kName + "Slc)")
	***REMOVED***
	x.linef("z.DecSendContainerState(codecSelfer_containerMapValue%s)", x.xs)
	x.decStructMapSwitch(kName, varname, rtid, t)

	x.line("***REMOVED*** // end for " + tpfx + "j" + i)
	x.linef("z.DecSendContainerState(codecSelfer_containerMapEnd%s)", x.xs)
***REMOVED***

func (x *genRunner) decStructArray(varname, lenvarname, breakString string, rtid uintptr, t reflect.Type) ***REMOVED***
	tpfx := genTempVarPfx
	i := x.varsfx()
	ti := x.ti.get(rtid, t)
	tisfi := ti.sfip // always use sequence from file. decStruct expects same thing.
	x.linef("var %sj%s int", tpfx, i)
	x.linef("var %sb%s bool", tpfx, i)                        // break
	x.linef("var %shl%s bool = %s >= 0", tpfx, i, lenvarname) // has length
	for _, si := range tisfi ***REMOVED***
		var t2 reflect.StructField
		if si.i != -1 ***REMOVED***
			t2 = t.Field(int(si.i))
		***REMOVED*** else ***REMOVED***
			//we must accommodate anonymous fields, where the embedded field is a nil pointer in the value.
			// t2 = t.FieldByIndex(si.is)
			t2typ := t
			varname3 := varname
			for _, ix := range si.is ***REMOVED***
				for t2typ.Kind() == reflect.Ptr ***REMOVED***
					t2typ = t2typ.Elem()
				***REMOVED***
				t2 = t2typ.Field(ix)
				t2typ = t2.Type
				varname3 = varname3 + "." + t2.Name
				if t2typ.Kind() == reflect.Ptr ***REMOVED***
					x.linef("if %s == nil ***REMOVED*** %s = new(%s) ***REMOVED***", varname3, varname3, x.genTypeName(t2typ.Elem()))
				***REMOVED***
			***REMOVED***
		***REMOVED***

		x.linef("%sj%s++; if %shl%s ***REMOVED*** %sb%s = %sj%s > %s ***REMOVED*** else ***REMOVED*** %sb%s = r.CheckBreak() ***REMOVED***",
			tpfx, i, tpfx, i, tpfx, i,
			tpfx, i, lenvarname, tpfx, i)
		x.linef("if %sb%s ***REMOVED*** z.DecSendContainerState(codecSelfer_containerArrayEnd%s); %s ***REMOVED***",
			tpfx, i, x.xs, breakString)
		x.linef("z.DecSendContainerState(codecSelfer_containerArrayElem%s)", x.xs)
		x.decVar(varname+"."+t2.Name, t2.Type, true)
	***REMOVED***
	// read remaining values and throw away.
	x.line("for ***REMOVED***")
	x.linef("%sj%s++; if %shl%s ***REMOVED*** %sb%s = %sj%s > %s ***REMOVED*** else ***REMOVED*** %sb%s = r.CheckBreak() ***REMOVED***",
		tpfx, i, tpfx, i, tpfx, i,
		tpfx, i, lenvarname, tpfx, i)
	x.linef("if %sb%s ***REMOVED*** break ***REMOVED***", tpfx, i)
	x.linef("z.DecSendContainerState(codecSelfer_containerArrayElem%s)", x.xs)
	x.linef(`z.DecStructFieldNotFound(%sj%s - 1, "")`, tpfx, i)
	x.line("***REMOVED***")
	x.linef("z.DecSendContainerState(codecSelfer_containerArrayEnd%s)", x.xs)
***REMOVED***

func (x *genRunner) decStruct(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	// if container is map
	i := x.varsfx()
	x.linef("%sct%s := r.ContainerType()", genTempVarPfx, i)
	x.linef("if %sct%s == codecSelferValueTypeMap%s ***REMOVED***", genTempVarPfx, i, x.xs)
	x.line(genTempVarPfx + "l" + i + " := r.ReadMapStart()")
	x.linef("if %sl%s == 0 ***REMOVED***", genTempVarPfx, i)
	x.linef("z.DecSendContainerState(codecSelfer_containerMapEnd%s)", x.xs)
	if genUseOneFunctionForDecStructMap ***REMOVED***
		x.line("***REMOVED*** else ***REMOVED*** ")
		x.linef("x.codecDecodeSelfFromMap(%sl%s, d)", genTempVarPfx, i)
	***REMOVED*** else ***REMOVED***
		x.line("***REMOVED*** else if " + genTempVarPfx + "l" + i + " > 0 ***REMOVED*** ")
		x.line("x.codecDecodeSelfFromMapLenPrefix(" + genTempVarPfx + "l" + i + ", d)")
		x.line("***REMOVED*** else ***REMOVED***")
		x.line("x.codecDecodeSelfFromMapCheckBreak(" + genTempVarPfx + "l" + i + ", d)")
	***REMOVED***
	x.line("***REMOVED***")

	// else if container is array
	x.linef("***REMOVED*** else if %sct%s == codecSelferValueTypeArray%s ***REMOVED***", genTempVarPfx, i, x.xs)
	x.line(genTempVarPfx + "l" + i + " := r.ReadArrayStart()")
	x.linef("if %sl%s == 0 ***REMOVED***", genTempVarPfx, i)
	x.linef("z.DecSendContainerState(codecSelfer_containerArrayEnd%s)", x.xs)
	x.line("***REMOVED*** else ***REMOVED*** ")
	x.linef("x.codecDecodeSelfFromArray(%sl%s, d)", genTempVarPfx, i)
	x.line("***REMOVED***")
	// else panic
	x.line("***REMOVED*** else ***REMOVED*** ")
	x.line("panic(codecSelferOnlyMapOrArrayEncodeToStructErr" + x.xs + ")")
	x.line("***REMOVED*** ")
***REMOVED***

// --------

type genV struct ***REMOVED***
	// genV is either a primitive (Primitive != "") or a map (MapKey != "") or a slice
	MapKey    string
	Elem      string
	Primitive string
	Size      int
***REMOVED***

func (x *genRunner) newGenV(t reflect.Type) (v genV) ***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Slice, reflect.Array:
		te := t.Elem()
		v.Elem = x.genTypeName(te)
		v.Size = int(te.Size())
	case reflect.Map:
		te, tk := t.Elem(), t.Key()
		v.Elem = x.genTypeName(te)
		v.MapKey = x.genTypeName(tk)
		v.Size = int(te.Size() + tk.Size())
	default:
		panic("unexpected type for newGenV. Requires map or slice type")
	***REMOVED***
	return
***REMOVED***

func (x *genV) MethodNamePfx(prefix string, prim bool) string ***REMOVED***
	var name []byte
	if prefix != "" ***REMOVED***
		name = append(name, prefix...)
	***REMOVED***
	if prim ***REMOVED***
		name = append(name, genTitleCaseName(x.Primitive)...)
	***REMOVED*** else ***REMOVED***
		if x.MapKey == "" ***REMOVED***
			name = append(name, "Slice"...)
		***REMOVED*** else ***REMOVED***
			name = append(name, "Map"...)
			name = append(name, genTitleCaseName(x.MapKey)...)
		***REMOVED***
		name = append(name, genTitleCaseName(x.Elem)...)
	***REMOVED***
	return string(name)

***REMOVED***

// genImportPath returns import path of a non-predeclared named typed, or an empty string otherwise.
//
// This handles the misbehaviour that occurs when 1.5-style vendoring is enabled,
// where PkgPath returns the full path, including the vendoring pre-fix that should have been stripped.
// We strip it here.
func genImportPath(t reflect.Type) (s string) ***REMOVED***
	s = t.PkgPath()
	if genCheckVendor ***REMOVED***
		// HACK: Misbehaviour occurs in go 1.5. May have to re-visit this later.
		// if s contains /vendor/ OR startsWith vendor/, then return everything after it.
		const vendorStart = "vendor/"
		const vendorInline = "/vendor/"
		if i := strings.LastIndex(s, vendorInline); i >= 0 ***REMOVED***
			s = s[i+len(vendorInline):]
		***REMOVED*** else if strings.HasPrefix(s, vendorStart) ***REMOVED***
			s = s[len(vendorStart):]
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// A go identifier is (letter|_)[letter|number|_]*
func genGoIdentifier(s string, checkFirstChar bool) string ***REMOVED***
	b := make([]byte, 0, len(s))
	t := make([]byte, 4)
	var n int
	for i, r := range s ***REMOVED***
		if checkFirstChar && i == 0 && !unicode.IsLetter(r) ***REMOVED***
			b = append(b, '_')
		***REMOVED***
		// r must be unicode_letter, unicode_digit or _
		if unicode.IsLetter(r) || unicode.IsDigit(r) ***REMOVED***
			n = utf8.EncodeRune(t, r)
			b = append(b, t[:n]...)
		***REMOVED*** else ***REMOVED***
			b = append(b, '_')
		***REMOVED***
	***REMOVED***
	return string(b)
***REMOVED***

func genNonPtr(t reflect.Type) reflect.Type ***REMOVED***
	for t.Kind() == reflect.Ptr ***REMOVED***
		t = t.Elem()
	***REMOVED***
	return t
***REMOVED***

func genTitleCaseName(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "interface***REMOVED******REMOVED***", "interface ***REMOVED******REMOVED***":
		return "Intf"
	default:
		return strings.ToUpper(s[0:1]) + s[1:]
	***REMOVED***
***REMOVED***

func genMethodNameT(t reflect.Type, tRef reflect.Type) (n string) ***REMOVED***
	var ptrPfx string
	for t.Kind() == reflect.Ptr ***REMOVED***
		ptrPfx += "Ptrto"
		t = t.Elem()
	***REMOVED***
	tstr := t.String()
	if tn := t.Name(); tn != "" ***REMOVED***
		if tRef != nil && genImportPath(t) == genImportPath(tRef) ***REMOVED***
			return ptrPfx + tn
		***REMOVED*** else ***REMOVED***
			if genQNameRegex.MatchString(tstr) ***REMOVED***
				return ptrPfx + strings.Replace(tstr, ".", "_", 1000)
			***REMOVED*** else ***REMOVED***
				return ptrPfx + genCustomTypeName(tstr)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Map:
		return ptrPfx + "Map" + genMethodNameT(t.Key(), tRef) + genMethodNameT(t.Elem(), tRef)
	case reflect.Slice:
		return ptrPfx + "Slice" + genMethodNameT(t.Elem(), tRef)
	case reflect.Array:
		return ptrPfx + "Array" + strconv.FormatInt(int64(t.Len()), 10) + genMethodNameT(t.Elem(), tRef)
	case reflect.Chan:
		var cx string
		switch t.ChanDir() ***REMOVED***
		case reflect.SendDir:
			cx = "ChanSend"
		case reflect.RecvDir:
			cx = "ChanRecv"
		default:
			cx = "Chan"
		***REMOVED***
		return ptrPfx + cx + genMethodNameT(t.Elem(), tRef)
	default:
		if t == intfTyp ***REMOVED***
			return ptrPfx + "Interface"
		***REMOVED*** else ***REMOVED***
			if tRef != nil && genImportPath(t) == genImportPath(tRef) ***REMOVED***
				if t.Name() != "" ***REMOVED***
					return ptrPfx + t.Name()
				***REMOVED*** else ***REMOVED***
					return ptrPfx + genCustomTypeName(tstr)
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				// best way to get the package name inclusive
				// return ptrPfx + strings.Replace(tstr, ".", "_", 1000)
				// return ptrPfx + genBase64enc.EncodeToString([]byte(tstr))
				if t.Name() != "" && genQNameRegex.MatchString(tstr) ***REMOVED***
					return ptrPfx + strings.Replace(tstr, ".", "_", 1000)
				***REMOVED*** else ***REMOVED***
					return ptrPfx + genCustomTypeName(tstr)
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// genCustomNameForType base64encodes the t.String() value in such a way
// that it can be used within a function name.
func genCustomTypeName(tstr string) string ***REMOVED***
	len2 := genBase64enc.EncodedLen(len(tstr))
	bufx := make([]byte, len2)
	genBase64enc.Encode(bufx, []byte(tstr))
	for i := len2 - 1; i >= 0; i-- ***REMOVED***
		if bufx[i] == '=' ***REMOVED***
			len2--
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***
	return string(bufx[:len2])
***REMOVED***

func genIsImmutable(t reflect.Type) (v bool) ***REMOVED***
	return isImmutableKind(t.Kind())
***REMOVED***

type genInternal struct ***REMOVED***
	Values []genV
	Unsafe bool
***REMOVED***

func (x genInternal) FastpathLen() (l int) ***REMOVED***
	for _, v := range x.Values ***REMOVED***
		if v.Primitive == "" ***REMOVED***
			l++
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func genInternalZeroValue(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "interface***REMOVED******REMOVED***", "interface ***REMOVED******REMOVED***":
		return "nil"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return "0"
	***REMOVED***
***REMOVED***

func genInternalEncCommandAsString(s string, vname string) string ***REMOVED***
	switch s ***REMOVED***
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return "ee.EncodeUint(uint64(" + vname + "))"
	case "int", "int8", "int16", "int32", "int64":
		return "ee.EncodeInt(int64(" + vname + "))"
	case "string":
		return "ee.EncodeString(c_UTF8, " + vname + ")"
	case "float32":
		return "ee.EncodeFloat32(" + vname + ")"
	case "float64":
		return "ee.EncodeFloat64(" + vname + ")"
	case "bool":
		return "ee.EncodeBool(" + vname + ")"
	case "symbol":
		return "ee.EncodeSymbol(" + vname + ")"
	default:
		return "e.encode(" + vname + ")"
	***REMOVED***
***REMOVED***

func genInternalDecCommandAsString(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "uint":
		return "uint(dd.DecodeUint(uintBitsize))"
	case "uint8":
		return "uint8(dd.DecodeUint(8))"
	case "uint16":
		return "uint16(dd.DecodeUint(16))"
	case "uint32":
		return "uint32(dd.DecodeUint(32))"
	case "uint64":
		return "dd.DecodeUint(64)"
	case "uintptr":
		return "uintptr(dd.DecodeUint(uintBitsize))"
	case "int":
		return "int(dd.DecodeInt(intBitsize))"
	case "int8":
		return "int8(dd.DecodeInt(8))"
	case "int16":
		return "int16(dd.DecodeInt(16))"
	case "int32":
		return "int32(dd.DecodeInt(32))"
	case "int64":
		return "dd.DecodeInt(64)"

	case "string":
		return "dd.DecodeString()"
	case "float32":
		return "float32(dd.DecodeFloat(true))"
	case "float64":
		return "dd.DecodeFloat(false)"
	case "bool":
		return "dd.DecodeBool()"
	default:
		panic(errors.New("gen internal: unknown type for decode: " + s))
	***REMOVED***
***REMOVED***

func genInternalSortType(s string, elem bool) string ***REMOVED***
	for _, v := range [...]string***REMOVED***"int", "uint", "float", "bool", "string"***REMOVED*** ***REMOVED***
		if strings.HasPrefix(s, v) ***REMOVED***
			if elem ***REMOVED***
				if v == "int" || v == "uint" || v == "float" ***REMOVED***
					return v + "64"
				***REMOVED*** else ***REMOVED***
					return v
				***REMOVED***
			***REMOVED***
			return v + "Slice"
		***REMOVED***
	***REMOVED***
	panic("sorttype: unexpected type: " + s)
***REMOVED***

// var genInternalMu sync.Mutex
var genInternalV genInternal
var genInternalTmplFuncs template.FuncMap
var genInternalOnce sync.Once

func genInternalInit() ***REMOVED***
	types := [...]string***REMOVED***
		"interface***REMOVED******REMOVED***",
		"string",
		"float32",
		"float64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"bool",
	***REMOVED***
	// keep as slice, so it is in specific iteration order.
	// Initial order was uint64, string, interface***REMOVED******REMOVED***, int, int64
	mapvaltypes := [...]string***REMOVED***
		"interface***REMOVED******REMOVED***",
		"string",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"float32",
		"float64",
		"bool",
	***REMOVED***
	wordSizeBytes := int(intBitsize) / 8

	mapvaltypes2 := map[string]int***REMOVED***
		"interface***REMOVED******REMOVED***": 2 * wordSizeBytes,
		"string":      2 * wordSizeBytes,
		"uint":        1 * wordSizeBytes,
		"uint8":       1,
		"uint16":      2,
		"uint32":      4,
		"uint64":      8,
		"uintptr":     1 * wordSizeBytes,
		"int":         1 * wordSizeBytes,
		"int8":        1,
		"int16":       2,
		"int32":       4,
		"int64":       8,
		"float32":     4,
		"float64":     8,
		"bool":        1,
	***REMOVED***
	var gt genInternal

	// For each slice or map type, there must be a (symmetrical) Encode and Decode fast-path function
	for _, s := range types ***REMOVED***
		gt.Values = append(gt.Values, genV***REMOVED***Primitive: s, Size: mapvaltypes2[s]***REMOVED***)
		if s != "uint8" ***REMOVED*** // do not generate fast path for slice of bytes. Treat specially already.
			gt.Values = append(gt.Values, genV***REMOVED***Elem: s, Size: mapvaltypes2[s]***REMOVED***)
		***REMOVED***
		if _, ok := mapvaltypes2[s]; !ok ***REMOVED***
			gt.Values = append(gt.Values, genV***REMOVED***MapKey: s, Elem: s, Size: 2 * mapvaltypes2[s]***REMOVED***)
		***REMOVED***
		for _, ms := range mapvaltypes ***REMOVED***
			gt.Values = append(gt.Values, genV***REMOVED***MapKey: s, Elem: ms, Size: mapvaltypes2[s] + mapvaltypes2[ms]***REMOVED***)
		***REMOVED***
	***REMOVED***

	funcs := make(template.FuncMap)
	// funcs["haspfx"] = strings.HasPrefix
	funcs["encmd"] = genInternalEncCommandAsString
	funcs["decmd"] = genInternalDecCommandAsString
	funcs["zerocmd"] = genInternalZeroValue
	funcs["hasprefix"] = strings.HasPrefix
	funcs["sorttype"] = genInternalSortType

	genInternalV = gt
	genInternalTmplFuncs = funcs
***REMOVED***

// genInternalGoFile is used to generate source files from templates.
// It is run by the program author alone.
// Unfortunately, it has to be exported so that it can be called from a command line tool.
// *** DO NOT USE ***
func genInternalGoFile(r io.Reader, w io.Writer, safe bool) (err error) ***REMOVED***
	genInternalOnce.Do(genInternalInit)

	gt := genInternalV
	gt.Unsafe = !safe

	t := template.New("").Funcs(genInternalTmplFuncs)

	tmplstr, err := ioutil.ReadAll(r)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	if t, err = t.Parse(string(tmplstr)); err != nil ***REMOVED***
		return
	***REMOVED***

	var out bytes.Buffer
	err = t.Execute(&out, gt)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	bout, err := format.Source(out.Bytes())
	if err != nil ***REMOVED***
		w.Write(out.Bytes()) // write out if error, so we can still see.
		// w.Write(bout) // write out if error, as much as possible, so we can still see.
		return
	***REMOVED***
	w.Write(bout)
	return
***REMOVED***
