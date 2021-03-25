// +build codecgen.exec

// Copyright (c) 2012-2018 Ugorji Nwoke. All rights reserved.
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
//   - MissingFielder implementation.
//     If a type implements MissingFielder, it is completely ignored by codecgen.
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
//
// Note:
//   codecgen-generated code depends on the variables defined by fast-path.generated.go.
//   consequently, you cannot run with tags "codecgen notfastpath".

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
// v6: removed unsafe from gen, and now uses codecgen.exec tag
// v7:
// v8: current - we now maintain compatibility with old generated code.
// v9: skipped
// v10: modified encDriver and decDriver interfaces.
// v11: remove deprecated methods of encDriver and decDriver.
// v12: removed deprecated methods from genHelper and changed container tracking logic
// v13: 20190603 removed DecodeString - use DecodeStringAsBytes instead
// v14: 20190611 refactored nil handling: TryDecodeAsNil -> selective TryNil, etc
// v15: 20190626 encDriver.EncodeString handles StringToRaw flag inside handle
// v16: 20190629 refactoring for v1.1.6
const genVersion = 16

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

	// genFastpathCanonical configures whether we support Canonical in fast path.
	// The savings is not much.
	//
	// NOTE: This MUST ALWAYS BE TRUE. fast-path.go.tmp doesn't handle it being false.
	genFastpathCanonical = true // MUST be true

	// genFastpathTrimTypes configures whether we trim uncommon fastpath types.
	genFastpathTrimTypes = true

	// genDecStructArrayInlineLoopCheck configures whether we create a next function
	// for each iteration in the loop and call it, or just inline it.
	//
	// with inlining, we get better performance but about 10% larger files.
	genDecStructArrayInlineLoopCheck = true
)

type genStructMapStyle uint8

const (
	genStructMapStyleConsolidated genStructMapStyle = iota
	genStructMapStyleLenPrefix
	genStructMapStyleCheckBreak
)

var (
	errGenAllTypesSamePkg  = errors.New("All types must be in the same package")
	errGenExpectArrayOrMap = errors.New("unexpected type. Expecting array/map/slice")

	genBase64enc  = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789__")
	genQNameRegex = regexp.MustCompile(`[A-Za-z_.]+`)
)

type genBuf struct ***REMOVED***
	buf []byte
***REMOVED***

func (x *genBuf) s(s string) *genBuf              ***REMOVED*** x.buf = append(x.buf, s...); return x ***REMOVED***
func (x *genBuf) b(s []byte) *genBuf              ***REMOVED*** x.buf = append(x.buf, s...); return x ***REMOVED***
func (x *genBuf) v() string                       ***REMOVED*** return string(x.buf) ***REMOVED***
func (x *genBuf) f(s string, args ...interface***REMOVED******REMOVED***) ***REMOVED*** x.s(fmt.Sprintf(s, args...)) ***REMOVED***
func (x *genBuf) reset() ***REMOVED***
	if x.buf != nil ***REMOVED***
		x.buf = x.buf[:0]
	***REMOVED***
***REMOVED***

// genRunner holds some state used during a Gen run.
type genRunner struct ***REMOVED***
	w io.Writer // output
	c uint64    // counter used for generating varsfx
	f uint64    // counter used for saying false

	t  []reflect.Type   // list of types to run selfer on
	tc reflect.Type     // currently running selfer on this type
	te map[uintptr]bool // types for which the encoder has been created
	td map[uintptr]bool // types for which the decoder has been created
	cp string           // codec import path

	im  map[string]reflect.Type // imports to add
	imn map[string]string       // package names of imports to add
	imc uint64                  // counter for import numbers

	is map[reflect.Type]struct***REMOVED******REMOVED*** // types seen during import search
	bp string                    // base PkgPath, for which we are generating for

	cpfx string // codec package prefix

	tm map[reflect.Type]struct***REMOVED******REMOVED*** // types for which enc/dec must be generated
	ts []reflect.Type            // types for which enc/dec must be generated

	xs string // top level variable/constant suffix
	hn string // fn helper type name

	ti *TypeInfos
	// rr *rand.Rand // random generator for file-specific types

	nx bool // no extensions
***REMOVED***

type genIfClause struct ***REMOVED***
	hasIf bool
***REMOVED***

func (g *genIfClause) end(x *genRunner) ***REMOVED***
	if g.hasIf ***REMOVED***
		x.line("***REMOVED***")
	***REMOVED***
***REMOVED***

func (g *genIfClause) c(last bool) (v string) ***REMOVED***
	if last ***REMOVED***
		if g.hasIf ***REMOVED***
			v = " ***REMOVED*** else ***REMOVED*** "
		***REMOVED***
	***REMOVED*** else if g.hasIf ***REMOVED***
		v = " ***REMOVED*** else if "
	***REMOVED*** else ***REMOVED***
		v = "if "
		g.hasIf = true
	***REMOVED***
	return
***REMOVED***

// Gen will write a complete go file containing Selfer implementations for each
// type passed. All the types must be in the same package.
//
// Library users: DO NOT USE IT DIRECTLY. IT WILL CHANGE CONTINUOUSLY WITHOUT NOTICE.
func Gen(w io.Writer, buildTags, pkgName, uid string, noExtensions bool,
	ti *TypeInfos, typ ...reflect.Type) ***REMOVED***
	// All types passed to this method do not have a codec.Selfer method implemented directly.
	// codecgen already checks the AST and skips any types that define the codec.Selfer methods.
	// Consequently, there's no need to check and trim them if they implement codec.Selfer

	if len(typ) == 0 ***REMOVED***
		return
	***REMOVED***
	x := genRunner***REMOVED***
		w:   w,
		t:   typ,
		te:  make(map[uintptr]bool),
		td:  make(map[uintptr]bool),
		im:  make(map[string]reflect.Type),
		imn: make(map[string]string),
		is:  make(map[reflect.Type]struct***REMOVED******REMOVED***),
		tm:  make(map[reflect.Type]struct***REMOVED******REMOVED***),
		ts:  []reflect.Type***REMOVED******REMOVED***,
		bp:  genImportPath(typ[0]),
		xs:  uid,
		ti:  ti,
		nx:  noExtensions,
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
			panic(errGenAllTypesSamePkg)
		***REMOVED***
		x.genRefPkgs(t)
	***REMOVED***
	x.line("// +build go1.6")
	if buildTags != "" ***REMOVED***
		x.line("// +build " + buildTags)
	***REMOVED***
	x.line(`

// Code generated by codecgen - DO NOT EDIT.

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
	for k := range x.im ***REMOVED***
		imKeys = append(imKeys, k)
	***REMOVED***
	sort.Strings(imKeys)
	for _, k := range imKeys ***REMOVED*** // for k, _ := range x.im ***REMOVED***
		if k == x.imn[k] ***REMOVED***
			x.linef("\"%s\"", k)
		***REMOVED*** else ***REMOVED***
			x.linef("%s \"%s\"", x.imn[k], k)
		***REMOVED***
	***REMOVED***
	// add required packages
	for _, k := range [...]string***REMOVED***"runtime", "errors", "strconv"***REMOVED*** ***REMOVED*** // "reflect", "fmt"
		if _, ok := x.im[k]; !ok ***REMOVED***
			x.line("\"" + k + "\"")
		***REMOVED***
	***REMOVED***
	x.line(")")
	x.line("")

	x.line("const (")
	x.linef("// ----- content types ----")
	x.linef("codecSelferCcUTF8%s = %v", x.xs, int64(cUTF8))
	x.linef("codecSelferCcRAW%s = %v", x.xs, int64(cRAW))
	x.linef("// ----- value types used ----")
	for _, vt := range [...]valueType***REMOVED***
		valueTypeArray, valueTypeMap, valueTypeString,
		valueTypeInt, valueTypeUint, valueTypeFloat,
		valueTypeNil,
	***REMOVED*** ***REMOVED***
		x.linef("codecSelferValueType%s%s = %v", vt.String(), x.xs, int64(vt))
	***REMOVED***

	x.linef("codecSelferBitsize%s = uint8(32 << (^uint(0) >> 63))", x.xs)
	x.linef("codecSelferDecContainerLenNil%s = %d", x.xs, int64(decContainerLenNil))
	x.line(")")
	x.line("var (")
	x.line("errCodecSelferOnlyMapOrArrayEncodeToStruct" + x.xs + " = " + "\nerrors.New(`only encoded map or array can be decoded into a struct`)")
	x.line(")")
	x.line("")

	x.hn = "codecSelfer" + x.xs
	x.line("type " + x.hn + " struct***REMOVED******REMOVED***")
	x.line("")
	x.linef("func %sFalse() bool ***REMOVED*** return false ***REMOVED***", x.hn)
	x.line("")
	x.varsfxreset()
	x.line("func init() ***REMOVED***")
	x.linef("if %sGenVersion != %v ***REMOVED***", x.cpfx, genVersion)
	x.line("_, file, _, _ := runtime.Caller(0)")
	x.linef("ver := strconv.FormatInt(int64(%sGenVersion), 10)", x.cpfx)
	x.outf(`panic("codecgen version mismatch: current: %v, need " + ver + ". Re-generate file: " + file)`, genVersion)
	x.linef("***REMOVED***")
	if len(imKeys) > 0 ***REMOVED***
		x.line("if false ***REMOVED*** // reference the types, but skip this branch at build/run time")
		for _, k := range imKeys ***REMOVED***
			t := x.im[k]
			x.linef("var _ %s.%s", x.imn[k], t.Name())
		***REMOVED***
		x.line("***REMOVED*** ") // close if false
	***REMOVED***
	x.line("***REMOVED***") // close init
	x.line("")

	// generate rest of type info
	for _, t := range typ ***REMOVED***
		x.tc = t
		x.selfer(true)
		x.selfer(false)
	***REMOVED***

	for _, t := range x.ts ***REMOVED***
		rtid := rt2id(t)
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
			panic(errGenExpectArrayOrMap)
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
			panic(errGenExpectArrayOrMap)
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

// sayFalse will either say "false" or use a function call that returns false.
func (x *genRunner) sayFalse() string ***REMOVED***
	x.f++
	if x.f%2 == 0 ***REMOVED***
		return x.hn + "False()"
	***REMOVED***
	return "false"
***REMOVED***

func (x *genRunner) varsfx() string ***REMOVED***
	x.c++
	return strconv.FormatUint(x.c, 10)
***REMOVED***

func (x *genRunner) varsfxreset() ***REMOVED***
	x.c = 0
***REMOVED***

func (x *genRunner) out(s string) ***REMOVED***
	_, err := io.WriteString(x.w, s)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x *genRunner) outf(s string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	_, err := fmt.Fprintf(x.w, s, params...)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
***REMOVED***

func (x *genRunner) line(s string) ***REMOVED***
	x.out(s)
	if len(s) == 0 || s[len(s)-1] != '\n' ***REMOVED***
		x.out("\n")
	***REMOVED***
***REMOVED***

func (x *genRunner) lineIf(s string) ***REMOVED***
	if s != "" ***REMOVED***
		x.line(s)
	***REMOVED***
***REMOVED***

func (x *genRunner) linef(s string, params ...interface***REMOVED******REMOVED***) ***REMOVED***
	x.outf(s, params...)
	if len(s) == 0 || s[len(s)-1] != '\n' ***REMOVED***
		x.out("\n")
	***REMOVED***
***REMOVED***

func (x *genRunner) genTypeName(t reflect.Type) (n string) ***REMOVED***
	// defer func() ***REMOVED*** xdebugf(">>>> ####: genTypeName: t: %v, name: '%s'\n", t, n) ***REMOVED***()

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
	// and structs/arrays always use a ptr receiver (encode|decode)
	isptr := !encode || t.Kind() == reflect.Array || (t.Kind() == reflect.Struct && t != timeTyp)
	x.varsfxreset()

	fnSigPfx := "func (" + genTopLevelVarName + " "
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
		x.encVar(genTopLevelVarName, t)
	***REMOVED*** else ***REMOVED***
		x.line(") CodecDecodeSelf(d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		// do not use decVar, as there is no need to check TryDecodeAsNil
		// or way to elegantly handle that, and also setting it to a
		// non-nil value doesn't affect the pointer passed.
		// x.decVar(genTopLevelVarName, t, false)
		x.dec(genTopLevelVarName, t0, true)
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
		x.decStructMap(genTopLevelVarName, "l", rt2id(t0), t0, genStructMapStyleConsolidated)
		x.line("***REMOVED***")
		x.line("")
	***REMOVED*** else ***REMOVED***
		x.out(fnSigPfx)
		x.line(") codecDecodeSelfFromMapLenPrefix(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		x.decStructMap(genTopLevelVarName, "l", rt2id(t0), t0, genStructMapStyleLenPrefix)
		x.line("***REMOVED***")
		x.line("")

		x.out(fnSigPfx)
		x.line(") codecDecodeSelfFromMapCheckBreak(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
		x.genRequiredMethodVars(false)
		x.decStructMap(genTopLevelVarName, "l", rt2id(t0), t0, genStructMapStyleCheckBreak)
		x.line("***REMOVED***")
		x.line("")
	***REMOVED***

	// write containerArray
	x.out(fnSigPfx)
	x.line(") codecDecodeSelfFromArray(l int, d *" + x.cpfx + "Decoder) ***REMOVED***")
	x.genRequiredMethodVars(false)
	x.decStructArray(genTopLevelVarName, "l", "return", rt2id(t0), t0)
	x.line("***REMOVED***")
	x.line("")

***REMOVED***

// used for chan, array, slice, map
func (x *genRunner) xtraSM(varname string, t reflect.Type, encode, isptr bool) ***REMOVED***
	var ptrPfx, addrPfx string
	if isptr ***REMOVED***
		ptrPfx = "*"
	***REMOVED*** else ***REMOVED***
		addrPfx = "&"
	***REMOVED***
	if encode ***REMOVED***
		x.linef("h.enc%s((%s%s)(%s), e)", x.genMethodNameT(t), ptrPfx, x.genTypeName(t), varname)
	***REMOVED*** else ***REMOVED***
		x.linef("h.dec%s((*%s)(%s%s), d)", x.genMethodNameT(t), x.genTypeName(t), addrPfx, varname)
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
	var checkNil bool
	// case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
	// do not include checkNil for slice and maps, as we already checkNil below it
	switch t.Kind() ***REMOVED***
	case reflect.Ptr, reflect.Interface, reflect.Chan:
		checkNil = true
	***REMOVED***
	x.encVarChkNil(varname, t, checkNil)
***REMOVED***

func (x *genRunner) encVarChkNil(varname string, t reflect.Type, checkNil bool) ***REMOVED***
	if checkNil ***REMOVED***
		x.linef("if %s == nil ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED***", varname)
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Ptr:
		telem := t.Elem()
		tek := telem.Kind()
		if tek == reflect.Array || (tek == reflect.Struct && telem != timeTyp) ***REMOVED***
			x.enc(varname, genNonPtr(t))
			break
		***REMOVED***
		i := x.varsfx()
		x.line(genTempVarPfx + i + " := *" + varname)
		x.enc(genTempVarPfx+i, genNonPtr(t))
	case reflect.Struct, reflect.Array:
		if t == timeTyp ***REMOVED***
			x.enc(varname, t)
			break
		***REMOVED***
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

// enc will encode a variable (varname) of type t, where t represents T.
// if t is !time.Time and t is of kind reflect.Struct or reflect.Array, varname is of type *T
// (to prevent copying),
// else t is of type T
func (x *genRunner) enc(varname string, t reflect.Type) ***REMOVED***
	rtid := rt2id(t)
	ti2 := x.ti.get(rtid, t)
	// We call CodecEncodeSelf if one of the following are honored:
	//   - the type already implements Selfer, call that
	//   - the type has a Selfer implementation just created, use that
	//   - the type is in the list of the ones we will generate for, but it is not currently being generated

	mi := x.varsfx()
	// tptr := reflect.PtrTo(t)
	tk := t.Kind()
	if x.checkForSelfer(t, varname) ***REMOVED***
		if tk == reflect.Array ||
			(tk == reflect.Struct && rtid != timeTypId) ***REMOVED*** // varname is of type *T
			// if tptr.Implements(selferTyp) || t.Implements(selferTyp) ***REMOVED***
			if ti2.isFlag(tiflagSelfer) || ti2.isFlag(tiflagSelferPtr) ***REMOVED***
				x.line(varname + ".CodecEncodeSelf(e)")
				return
			***REMOVED***
		***REMOVED*** else ***REMOVED*** // varname is of type T
			if ti2.isFlag(tiflagSelfer) ***REMOVED***
				x.line(varname + ".CodecEncodeSelf(e)")
				return
			***REMOVED*** else if ti2.isFlag(tiflagSelferPtr) ***REMOVED***
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
	//   - type is time.Time, RawExt, Raw
	//   - the type implements (Text|JSON|Binary)(Unm|M)arshal

	var hasIf genIfClause
	defer hasIf.end(x) // end if block (if necessary)

	if t == timeTyp ***REMOVED***
		x.linef("%s !z.EncBasicHandle().TimeNotBuiltin ***REMOVED*** r.EncodeTime(%s)", hasIf.c(false), varname)
		// return
	***REMOVED***
	if t == rawTyp ***REMOVED***
		x.linef("%s z.EncRaw(%s)", hasIf.c(true), varname)
		return
	***REMOVED***
	if t == rawExtTyp ***REMOVED***
		x.linef("%s r.EncodeRawExt(%s)", hasIf.c(true), varname)
		return
	***REMOVED***
	// only check for extensions if extensions are configured,
	// and the type is named, and has a packagePath,
	// and this is not the CodecEncodeSelf or CodecDecodeSelf method (i.e. it is not a Selfer)
	var arrayOrStruct = tk == reflect.Array || tk == reflect.Struct // meaning varname if of type *T
	if !x.nx && varname != genTopLevelVarName && genImportPath(t) != "" && t.Name() != "" ***REMOVED***
		yy := fmt.Sprintf("%sxt%s", genTempVarPfx, mi)
		x.linef("%s %s := z.Extension(z.I2Rtid(%s)); %s != nil ***REMOVED*** z.EncExtension(%s, %s) ",
			hasIf.c(false), yy, varname, yy, varname, yy)
	***REMOVED***
	if arrayOrStruct ***REMOVED*** // varname is of type *T
		if ti2.isFlag(tiflagBinaryMarshaler) || ti2.isFlag(tiflagBinaryMarshalerPtr) ***REMOVED***
			x.linef("%s z.EncBinary() ***REMOVED*** z.EncBinaryMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED***
		if ti2.isFlag(tiflagJsonMarshaler) || ti2.isFlag(tiflagJsonMarshalerPtr) ***REMOVED***
			x.linef("%s !z.EncBinary() && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED*** else if ti2.isFlag(tiflagTextUnmarshaler) || ti2.isFlag(tiflagTextUnmarshalerPtr) ***REMOVED***
			x.linef("%s !z.EncBinary() ***REMOVED*** z.EncTextMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED***
	***REMOVED*** else ***REMOVED*** // varname is of type T
		if ti2.isFlag(tiflagBinaryMarshaler) ***REMOVED***
			x.linef("%s z.EncBinary() ***REMOVED*** z.EncBinaryMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED*** else if ti2.isFlag(tiflagBinaryMarshalerPtr) ***REMOVED***
			x.linef("%s z.EncBinary() ***REMOVED*** z.EncBinaryMarshal(&%v) ", hasIf.c(false), varname)
		***REMOVED***
		if ti2.isFlag(tiflagJsonMarshaler) ***REMOVED***
			x.linef("%s !z.EncBinary() && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED*** else if ti2.isFlag(tiflagJsonMarshalerPtr) ***REMOVED***
			x.linef("%s !z.EncBinary() && z.IsJSONHandle() ***REMOVED*** z.EncJSONMarshal(&%v) ", hasIf.c(false), varname)
		***REMOVED*** else if ti2.isFlag(tiflagTextMarshaler) ***REMOVED***
			x.linef("%s !z.EncBinary() ***REMOVED*** z.EncTextMarshal(%v) ", hasIf.c(false), varname)
		***REMOVED*** else if ti2.isFlag(tiflagTextMarshalerPtr) ***REMOVED***
			x.linef("%s !z.EncBinary() ***REMOVED*** z.EncTextMarshal(&%v) ", hasIf.c(false), varname)
		***REMOVED***
	***REMOVED***
	x.lineIf(hasIf.c(true))

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
		x.linef("r.EncodeString(string(%s))", varname)
	case reflect.Chan:
		x.xtraSM(varname, t, true, false)
		// x.encListFallback(varname, rtid, t)
	case reflect.Array:
		x.xtraSM(varname, t, true, true)
	case reflect.Slice:
		// if nil, call dedicated function
		// if a []uint8, call dedicated function
		// if a known fastpath slice, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		x.linef("if %s == nil ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED***", varname)
		if rtid == uint8SliceTypId ***REMOVED***
			x.line("r.EncodeStringBytesRaw([]byte(" + varname + "))")
		***REMOVED*** else if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newFastpathGenV(t)
			x.line("z.F." + g.MethodNamePfx("Enc", false) + "V(" + varname + ", e)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, t, true, false)
		***REMOVED***
		x.linef("***REMOVED*** // end block: if %s slice == nil", varname)
	case reflect.Map:
		// if nil, call dedicated function
		// if a known fastpath map, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		x.linef("if %s == nil ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED***", varname)
		if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newFastpathGenV(t)
			x.line("z.F." + g.MethodNamePfx("Enc", false) + "V(" + varname + ", e)")
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, t, true, false)
		***REMOVED***
		x.linef("***REMOVED*** // end block: if %s map == nil", varname)
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
		x.linef(`r.EncodeString("")`)
	default:
		x.line("r.EncodeNil()")
	***REMOVED***
***REMOVED***

func (x *genRunner) doEncOmitEmptyLine(t2 reflect.StructField, varname string, buf *genBuf) ***REMOVED***
	x.f = 0
	x.encOmitEmptyLine(t2, varname, buf)
***REMOVED***

func (x *genRunner) encOmitEmptyLine(t2 reflect.StructField, varname string, buf *genBuf) ***REMOVED***
	// smartly check omitEmpty on a struct type, as it may contain uncomparable map/slice/etc.
	// also, for maps/slices/arrays, check if len ! 0 (not if == zero value)
	varname2 := varname + "." + t2.Name
	switch t2.Type.Kind() ***REMOVED***
	case reflect.Struct:
		rtid2 := rt2id(t2.Type)
		ti2 := x.ti.get(rtid2, t2.Type)
		// fmt.Printf(">>>> structfield: omitempty: type: %s, field: %s\n", t2.Type.Name(), t2.Name)
		if ti2.rtid == timeTypId ***REMOVED***
			buf.s("!(").s(varname2).s(".IsZero())")
			break
		***REMOVED***
		if ti2.isFlag(tiflagIsZeroerPtr) || ti2.isFlag(tiflagIsZeroer) ***REMOVED***
			buf.s("!(").s(varname2).s(".IsZero())")
			break
		***REMOVED***
		if ti2.isFlag(tiflagComparable) ***REMOVED***
			buf.s(varname2).s(" != ").s(x.genZeroValueR(t2.Type))
			break
		***REMOVED***
		// buf.s("(")
		buf.s(x.sayFalse()) // buf.s("false")
		for i, n := 0, t2.Type.NumField(); i < n; i++ ***REMOVED***
			f := t2.Type.Field(i)
			if f.PkgPath != "" ***REMOVED*** // unexported
				continue
			***REMOVED***
			buf.s(" || ")
			x.encOmitEmptyLine(f, varname2, buf)
		***REMOVED***
		//buf.s(")")
	case reflect.Bool:
		buf.s(varname2)
	case reflect.Map, reflect.Slice, reflect.Array, reflect.Chan:
		buf.s("len(").s(varname2).s(") != 0")
	default:
		buf.s(varname2).s(" != ").s(x.genZeroValueR(t2.Type))
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
	x.linef("_, _ = %s, %s", sepVarname, struct2arrvar)
	x.linef("const %s bool = %v // struct tag has 'toArray'", ti2arrayvar, ti.toArray)

	tisfi := ti.sfiSrc // always use sequence from file. decStruct expects same thing.

	// var nn int
	// due to omitEmpty, we need to calculate the
	// number of non-empty things we write out first.
	// This is required as we need to pre-determine the size of the container,
	// to support length-prefixing.
	if ti.anyOmitEmpty ***REMOVED***
		x.linef("var %s = [%v]bool***REMOVED*** // should field at this index be written?", numfieldsvar, len(tisfi))

		for j, si := range tisfi ***REMOVED***
			_ = j
			if !si.omitEmpty() ***REMOVED***
				x.linef("true, // %s", si.fieldName)
				continue
			***REMOVED***
			var t2 reflect.StructField
			var omitline genBuf
			***REMOVED***
				t2typ := t
				varname3 := varname
				// go through the loop, record the t2 field explicitly,
				// and gather the omit line if embedded in pointers.
				for ij, ix := range si.is ***REMOVED***
					if uint8(ij) == si.nis ***REMOVED***
						break
					***REMOVED***
					for t2typ.Kind() == reflect.Ptr ***REMOVED***
						t2typ = t2typ.Elem()
					***REMOVED***
					t2 = t2typ.Field(int(ix))
					t2typ = t2.Type
					varname3 = varname3 + "." + t2.Name
					// do not include actual field in the omit line.
					// that is done subsequently (right after - below).
					if uint8(ij+1) < si.nis && t2typ.Kind() == reflect.Ptr ***REMOVED***
						omitline.s(varname3).s(" != nil && ")
					***REMOVED***
				***REMOVED***
			***REMOVED***
			x.doEncOmitEmptyLine(t2, varname, &omitline)
			x.linef("%s, // %s", omitline.v(), si.fieldName)
		***REMOVED***
		x.line("***REMOVED***")
		x.linef("_ = %s", numfieldsvar)
	***REMOVED***

	type genFQN struct ***REMOVED***
		i       string
		fqname  string
		nilLine genBuf
		nilVar  string
		canNil  bool
		sf      reflect.StructField
	***REMOVED***

	genFQNs := make([]genFQN, len(tisfi))
	for j, si := range tisfi ***REMOVED***
		q := &genFQNs[j]
		q.i = x.varsfx()
		q.nilVar = genTempVarPfx + "n" + q.i
		q.canNil = false
		q.fqname = varname
		***REMOVED***
			t2typ := t
			for ij, ix := range si.is ***REMOVED***
				if uint8(ij) == si.nis ***REMOVED***
					break
				***REMOVED***
				for t2typ.Kind() == reflect.Ptr ***REMOVED***
					t2typ = t2typ.Elem()
				***REMOVED***
				q.sf = t2typ.Field(int(ix))
				t2typ = q.sf.Type
				q.fqname += "." + q.sf.Name
				if t2typ.Kind() == reflect.Ptr ***REMOVED***
					if !q.canNil ***REMOVED***
						q.nilLine.f("%s == nil", q.fqname)
						q.canNil = true
					***REMOVED*** else ***REMOVED***
						q.nilLine.f(" || %s == nil", q.fqname)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	for j := range genFQNs ***REMOVED***
		q := &genFQNs[j]
		if q.canNil ***REMOVED***
			x.linef("var %s bool = %s", q.nilVar, q.nilLine.v())
		***REMOVED***
	***REMOVED***

	x.linef("if %s || %s ***REMOVED***", ti2arrayvar, struct2arrvar) // if ti.toArray
	x.linef("z.EncWriteArrayStart(%d)", len(tisfi))

	for j, si := range tisfi ***REMOVED***
		q := &genFQNs[j]
		// if the type of the field is a Selfer, or one of the ones
		if q.canNil ***REMOVED***
			x.linef("if %s ***REMOVED*** z.EncWriteArrayElem(); r.EncodeNil() ***REMOVED*** else ***REMOVED*** ", q.nilVar)
		***REMOVED***
		x.linef("z.EncWriteArrayElem()")
		if si.omitEmpty() ***REMOVED***
			x.linef("if %s[%v] ***REMOVED***", numfieldsvar, j)
		***REMOVED***
		x.encVarChkNil(q.fqname, q.sf.Type, false)
		if si.omitEmpty() ***REMOVED***
			x.linef("***REMOVED*** else ***REMOVED***")
			x.encZero(q.sf.Type)
			x.linef("***REMOVED***")
		***REMOVED***
		if q.canNil ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***
	***REMOVED***

	x.line("z.EncWriteArrayEnd()")
	x.linef("***REMOVED*** else ***REMOVED***") // if not ti.toArray
	if ti.anyOmitEmpty ***REMOVED***
		x.linef("var %snn%s int", genTempVarPfx, i)
		x.linef("for _, b := range %s ***REMOVED*** if b ***REMOVED*** %snn%s++ ***REMOVED*** ***REMOVED***", numfieldsvar, genTempVarPfx, i)
		x.linef("z.EncWriteMapStart(%snn%s)", genTempVarPfx, i)
		x.linef("%snn%s = %v", genTempVarPfx, i, 0)
	***REMOVED*** else ***REMOVED***
		x.linef("z.EncWriteMapStart(%d)", len(tisfi))
	***REMOVED***

	for j, si := range tisfi ***REMOVED***
		q := &genFQNs[j]
		if si.omitEmpty() ***REMOVED***
			x.linef("if %s[%v] ***REMOVED***", numfieldsvar, j)
		***REMOVED***
		x.linef("z.EncWriteMapElemKey()")

		// emulate EncStructFieldKey
		switch ti.keyType ***REMOVED***
		case valueTypeInt:
			x.linef("r.EncodeInt(z.M.Int(strconv.ParseInt(`%s`, 10, 64)))", si.encName)
		case valueTypeUint:
			x.linef("r.EncodeUint(z.M.Uint(strconv.ParseUint(`%s`, 10, 64)))", si.encName)
		case valueTypeFloat:
			x.linef("r.EncodeFloat64(z.M.Float(strconv.ParseFloat(`%s`, 64)))", si.encName)
		default: // string
			if si.encNameAsciiAlphaNum ***REMOVED***
				x.linef(`if z.IsJSONHandle() ***REMOVED*** z.WriteStr("\"%s\"") ***REMOVED*** else ***REMOVED*** `, si.encName)
			***REMOVED***
			x.linef("r.EncodeString(`%s`)", si.encName)
			if si.encNameAsciiAlphaNum ***REMOVED***
				x.linef("***REMOVED***")
			***REMOVED***
		***REMOVED***
		x.line("z.EncWriteMapElemValue()")
		if q.canNil ***REMOVED***
			x.line("if " + q.nilVar + " ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED*** ")
			x.encVarChkNil(q.fqname, q.sf.Type, false)
			x.line("***REMOVED***")
		***REMOVED*** else ***REMOVED***
			x.encVarChkNil(q.fqname, q.sf.Type, false)
		***REMOVED***
		if si.omitEmpty() ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***
	***REMOVED***
	x.line("z.EncWriteMapEnd()")
	x.linef("***REMOVED*** ") // end if/else ti.toArray
***REMOVED***

func (x *genRunner) encListFallback(varname string, t reflect.Type) ***REMOVED***
	x.linef("if %s == nil ***REMOVED*** r.EncodeNil(); return ***REMOVED***", varname)
	elemBytes := t.Elem().Kind() == reflect.Uint8
	if t.AssignableTo(uint8SliceTyp) ***REMOVED***
		x.linef("r.EncodeStringBytesRaw([]byte(%s))", varname)
		return
	***REMOVED***
	if t.Kind() == reflect.Array && elemBytes ***REMOVED***
		x.linef("r.EncodeStringBytesRaw(((*[%d]byte)(%s))[:])", t.Len(), varname)
		return
	***REMOVED***
	i := x.varsfx()
	if t.Kind() == reflect.Chan ***REMOVED***
		type ts struct ***REMOVED***
			Label, Chan, Slice, Sfx string
		***REMOVED***
		tm, err := template.New("").Parse(genEncChanTmpl)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		x.linef("if %s == nil ***REMOVED*** r.EncodeNil() ***REMOVED*** else ***REMOVED*** ", varname)
		x.linef("var sch%s []%s", i, x.genTypeName(t.Elem()))
		err = tm.Execute(x.w, &ts***REMOVED***"Lsch" + i, varname, "sch" + i, i***REMOVED***)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		if elemBytes ***REMOVED***
			x.linef("r.EncodeStringBytesRaw([]byte(%s))", "sch"+i)
			x.line("***REMOVED***")
			return
		***REMOVED***
		varname = "sch" + i
	***REMOVED***

	x.line("z.EncWriteArrayStart(len(" + varname + "))")
	x.linef("for _, %sv%s := range %s ***REMOVED***", genTempVarPfx, i, varname)
	x.linef("z.EncWriteArrayElem()")

	x.encVar(genTempVarPfx+"v"+i, t.Elem())
	x.line("***REMOVED***")
	x.line("z.EncWriteArrayEnd()")
	if t.Kind() == reflect.Chan ***REMOVED***
		x.line("***REMOVED***")
	***REMOVED***
***REMOVED***

func (x *genRunner) encMapFallback(varname string, t reflect.Type) ***REMOVED***
	x.linef("if %s == nil ***REMOVED*** r.EncodeNil(); return ***REMOVED***", varname)
	// NOTE: Canonical Option is not honored
	i := x.varsfx()
	x.line("z.EncWriteMapStart(len(" + varname + "))")
	x.linef("for %sk%s, %sv%s := range %s ***REMOVED***", genTempVarPfx, i, genTempVarPfx, i, varname)
	x.linef("z.EncWriteMapElemKey()")
	x.encVar(genTempVarPfx+"k"+i, t.Key())
	x.line("z.EncWriteMapElemValue()")
	x.encVar(genTempVarPfx+"v"+i, t.Elem())
	x.line("***REMOVED***")
	x.line("z.EncWriteMapEnd()")
***REMOVED***

func (x *genRunner) decVarInitPtr(varname, nilvar string, t reflect.Type, si *structFieldInfo,
	newbuf, nilbuf *genBuf) (varname3 string, t2 reflect.StructField) ***REMOVED***
	//we must accommodate anonymous fields, where the embedded field is a nil pointer in the value.
	// t2 = t.FieldByIndex(si.is)
	varname3 = varname
	t2typ := t
	t2kind := t2typ.Kind()
	var nilbufed bool
	if si != nil ***REMOVED***
		for ij, ix := range si.is ***REMOVED***
			if uint8(ij) == si.nis ***REMOVED***
				break
			***REMOVED***
			// only one-level pointers can be seen in a type
			if t2typ.Kind() == reflect.Ptr ***REMOVED***
				t2typ = t2typ.Elem()
			***REMOVED***
			t2 = t2typ.Field(int(ix))
			t2typ = t2.Type
			varname3 = varname3 + "." + t2.Name
			t2kind = t2typ.Kind()
			if t2kind != reflect.Ptr ***REMOVED***
				continue
			***REMOVED***
			if newbuf != nil ***REMOVED***
				if len(newbuf.buf) > 0 ***REMOVED***
					newbuf.s("\n")
				***REMOVED***
				newbuf.f("if %s == nil ***REMOVED*** %s = new(%s) ***REMOVED***", varname3, varname3, x.genTypeName(t2typ.Elem()))
			***REMOVED***
			if nilbuf != nil ***REMOVED***
				if !nilbufed ***REMOVED***
					nilbuf.s("if ").s(varname3).s(" != nil")
					nilbufed = true
				***REMOVED*** else ***REMOVED***
					nilbuf.s(" && ").s(varname3).s(" != nil")
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if nilbuf != nil ***REMOVED***
		if nilbufed ***REMOVED***
			nilbuf.s(" ***REMOVED*** ").s("// remove the if-true\n")
		***REMOVED***
		if nilvar != "" ***REMOVED***
			nilbuf.s(nilvar).s(" = true")
		***REMOVED*** else if tk := t2typ.Kind(); tk == reflect.Ptr ***REMOVED***
			if strings.IndexByte(varname3, '.') != -1 || strings.IndexByte(varname3, '[') != -1 ***REMOVED***
				nilbuf.s(varname3).s(" = nil")
			***REMOVED*** else ***REMOVED***
				nilbuf.s("*").s(varname3).s(" = ").s(x.genZeroValueR(t2typ.Elem()))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			nilbuf.s(varname3).s(" = ").s(x.genZeroValueR(t2typ))
		***REMOVED***
		if nilbufed ***REMOVED***
			nilbuf.s("***REMOVED***")
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// decVar takes a variable called varname, of type t
func (x *genRunner) decVarMain(varname, rand string, t reflect.Type, checkNotNil bool) ***REMOVED***
	// We only encode as nil if a nillable value.
	// This removes some of the wasted checks for TryDecodeAsNil.
	// We need to think about this more, to see what happens if omitempty, etc
	// cause a nil value to be stored when something is expected.
	// This could happen when decoding from a struct encoded as an array.
	// For that, decVar should be called with canNil=true, to force true as its value.
	var varname2 string
	if t.Kind() != reflect.Ptr ***REMOVED***
		if t.PkgPath() != "" || !x.decTryAssignPrimitive(varname, t, false) ***REMOVED***
			x.dec(varname, t, false)
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if checkNotNil ***REMOVED***
			x.linef("if %s == nil ***REMOVED*** %s = new(%s) ***REMOVED***", varname, varname, x.genTypeName(t.Elem()))
		***REMOVED***
		// Ensure we set underlying ptr to a non-nil value (so we can deref to it later).
		// There's a chance of a **T in here which is nil.
		var ptrPfx string
		for t = t.Elem(); t.Kind() == reflect.Ptr; t = t.Elem() ***REMOVED***
			ptrPfx += "*"
			if checkNotNil ***REMOVED***
				x.linef("if %s%s == nil ***REMOVED*** %s%s = new(%s)***REMOVED***", ptrPfx, varname, ptrPfx, varname, x.genTypeName(t))
			***REMOVED***
		***REMOVED***
		// Should we create temp var if a slice/map indexing? No. dec(...) can now handle it.

		if ptrPfx == "" ***REMOVED***
			x.dec(varname, t, true)
		***REMOVED*** else ***REMOVED***
			varname2 = genTempVarPfx + "z" + rand
			x.line(varname2 + " := " + ptrPfx + varname)
			x.dec(varname2, t, true)
		***REMOVED***
	***REMOVED***
***REMOVED***

// decVar takes a variable called varname, of type t
func (x *genRunner) decVar(varname, nilvar string, t reflect.Type, canBeNil, checkNotNil bool) ***REMOVED***

	// We only encode as nil if a nillable value.
	// This removes some of the wasted checks for TryDecodeAsNil.
	// We need to think about this more, to see what happens if omitempty, etc
	// cause a nil value to be stored when something is expected.
	// This could happen when decoding from a struct encoded as an array.
	// For that, decVar should be called with canNil=true, to force true as its value.

	i := x.varsfx()
	if t.Kind() == reflect.Ptr ***REMOVED***
		var buf genBuf
		x.decVarInitPtr(varname, nilvar, t, nil, nil, &buf)
		x.linef("if r.TryNil() ***REMOVED*** %s ***REMOVED*** else ***REMOVED***", buf.buf)
		x.decVarMain(varname, i, t, checkNotNil)
		x.line("***REMOVED*** ")
	***REMOVED*** else ***REMOVED***
		x.decVarMain(varname, i, t, checkNotNil)
	***REMOVED***
***REMOVED***

// dec will decode a variable (varname) of type t or ptrTo(t) if isptr==true.
// t is always a basetype (i.e. not of kind reflect.Ptr).
func (x *genRunner) dec(varname string, t reflect.Type, isptr bool) ***REMOVED***
	// assumptions:
	//   - the varname is to a pointer already. No need to take address of it
	//   - t is always a baseType T (not a *T, etc).
	rtid := rt2id(t)
	ti2 := x.ti.get(rtid, t)
	if x.checkForSelfer(t, varname) ***REMOVED***
		if ti2.isFlag(tiflagSelfer) || ti2.isFlag(tiflagSelferPtr) ***REMOVED***
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
	//   - type is time.Time, Raw, RawExt
	//   - the type implements (Text|JSON|Binary)(Unm|M)arshal

	mi := x.varsfx()

	var hasIf genIfClause
	defer hasIf.end(x)

	var ptrPfx, addrPfx string
	if isptr ***REMOVED***
		ptrPfx = "*"
	***REMOVED*** else ***REMOVED***
		addrPfx = "&"
	***REMOVED***
	if t == timeTyp ***REMOVED***
		x.linef("%s !z.DecBasicHandle().TimeNotBuiltin ***REMOVED*** %s%v = r.DecodeTime()", hasIf.c(false), ptrPfx, varname)
		// return
	***REMOVED***
	if t == rawTyp ***REMOVED***
		x.linef("%s %s%v = z.DecRaw()", hasIf.c(true), ptrPfx, varname)
		return
	***REMOVED***

	if t == rawExtTyp ***REMOVED***
		x.linef("%s r.DecodeExt(%s%v, 0, nil)", hasIf.c(true), addrPfx, varname)
		return
	***REMOVED***

	// only check for extensions if extensions are configured,
	// and the type is named, and has a packagePath,
	// and this is not the CodecEncodeSelf or CodecDecodeSelf method (i.e. it is not a Selfer)
	if !x.nx && varname != genTopLevelVarName && genImportPath(t) != "" && t.Name() != "" ***REMOVED***
		// first check if extensions are configued, before doing the interface conversion
		yy := fmt.Sprintf("%sxt%s", genTempVarPfx, mi)
		x.linef("%s %s := z.Extension(z.I2Rtid(%s)); %s != nil ***REMOVED*** z.DecExtension(%s, %s) ", hasIf.c(false), yy, varname, yy, varname, yy)
	***REMOVED***

	if ti2.isFlag(tiflagBinaryUnmarshaler) || ti2.isFlag(tiflagBinaryUnmarshalerPtr) ***REMOVED***
		x.linef("%s z.DecBinary() ***REMOVED*** z.DecBinaryUnmarshal(%s%v) ", hasIf.c(false), addrPfx, varname)
	***REMOVED***
	if ti2.isFlag(tiflagJsonUnmarshaler) || ti2.isFlag(tiflagJsonUnmarshalerPtr) ***REMOVED***
		x.linef("%s !z.DecBinary() && z.IsJSONHandle() ***REMOVED*** z.DecJSONUnmarshal(%s%v)", hasIf.c(false), addrPfx, varname)
	***REMOVED*** else if ti2.isFlag(tiflagTextUnmarshaler) || ti2.isFlag(tiflagTextUnmarshalerPtr) ***REMOVED***
		x.linef("%s !z.DecBinary() ***REMOVED*** z.DecTextUnmarshal(%s%v)", hasIf.c(false), addrPfx, varname)
	***REMOVED***

	x.lineIf(hasIf.c(true))

	if x.decTryAssignPrimitive(varname, t, isptr) ***REMOVED***
		return
	***REMOVED***

	switch t.Kind() ***REMOVED***
	case reflect.Array, reflect.Chan:
		x.xtraSM(varname, t, false, isptr)
	case reflect.Slice:
		// if a []uint8, call dedicated function
		// if a known fastpath slice, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		if rtid == uint8SliceTypId ***REMOVED***
			x.linef("%s%s = r.DecodeBytes(%s(%s[]byte)(%s), false)",
				ptrPfx, varname, ptrPfx, ptrPfx, varname)
		***REMOVED*** else if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newFastpathGenV(t)
			x.linef("z.F.%sX(%s%s, d)", g.MethodNamePfx("Dec", false), addrPfx, varname)
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, t, false, isptr)
			// x.decListFallback(varname, rtid, false, t)
		***REMOVED***
	case reflect.Map:
		// if a known fastpath map, call dedicated function
		// else write encode function in-line.
		// - if elements are primitives or Selfers, call dedicated function on each member.
		// - else call Encoder.encode(XXX) on it.
		if fastpathAV.index(rtid) != -1 ***REMOVED***
			g := x.newFastpathGenV(t)
			x.linef("z.F.%sX(%s%s, d)", g.MethodNamePfx("Dec", false), addrPfx, varname)
		***REMOVED*** else ***REMOVED***
			x.xtraSM(varname, t, false, isptr)
		***REMOVED***
	case reflect.Struct:
		if inlist ***REMOVED***
			// no need to create temp variable if isptr, or x.F or x[F]
			if isptr || strings.IndexByte(varname, '.') != -1 || strings.IndexByte(varname, '[') != -1 ***REMOVED***
				x.decStruct(varname, rtid, t)
			***REMOVED*** else ***REMOVED***
				varname2 := genTempVarPfx + "j" + mi
				x.line(varname2 + " := &" + varname)
				x.decStruct(varname2, rtid, t)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// delete(x.td, rtid)
			x.line("z.DecFallback(" + addrPfx + varname + ", false)")
		***REMOVED***
	default:
		if rtidAdded ***REMOVED***
			delete(x.te, rtid)
		***REMOVED***
		x.line("z.DecFallback(" + addrPfx + varname + ", true)")
	***REMOVED***
***REMOVED***

func (x *genRunner) decTryAssignPrimitive(varname string, t reflect.Type, isptr bool) (done bool) ***REMOVED***
	// This should only be used for exact primitives (ie un-named types).
	// Named types may be implementations of Selfer, Unmarshaler, etc.
	// They should be handled by dec(...)

	var ptr string
	if isptr ***REMOVED***
		ptr = "*"
	***REMOVED***
	switch t.Kind() ***REMOVED***
	case reflect.Int:
		x.linef("%s%s = (%s)(z.C.IntV(r.DecodeInt64(), codecSelferBitsize%s))", ptr, varname, x.genTypeName(t), x.xs)
	case reflect.Int8:
		x.linef("%s%s = (%s)(z.C.IntV(r.DecodeInt64(), 8))", ptr, varname, x.genTypeName(t))
	case reflect.Int16:
		x.linef("%s%s = (%s)(z.C.IntV(r.DecodeInt64(), 16))", ptr, varname, x.genTypeName(t))
	case reflect.Int32:
		x.linef("%s%s = (%s)(z.C.IntV(r.DecodeInt64(), 32))", ptr, varname, x.genTypeName(t))
	case reflect.Int64:
		x.linef("%s%s = (%s)(r.DecodeInt64())", ptr, varname, x.genTypeName(t))

	case reflect.Uint:
		x.linef("%s%s = (%s)(z.C.UintV(r.DecodeUint64(), codecSelferBitsize%s))", ptr, varname, x.genTypeName(t), x.xs)
	case reflect.Uint8:
		x.linef("%s%s = (%s)(z.C.UintV(r.DecodeUint64(), 8))", ptr, varname, x.genTypeName(t))
	case reflect.Uint16:
		x.linef("%s%s = (%s)(z.C.UintV(r.DecodeUint64(), 16))", ptr, varname, x.genTypeName(t))
	case reflect.Uint32:
		x.linef("%s%s = (%s)(z.C.UintV(r.DecodeUint64(), 32))", ptr, varname, x.genTypeName(t))
	case reflect.Uint64:
		x.linef("%s%s = (%s)(r.DecodeUint64())", ptr, varname, x.genTypeName(t))
	case reflect.Uintptr:
		x.linef("%s%s = (%s)(z.C.UintV(r.DecodeUint64(), codecSelferBitsize%s))", ptr, varname, x.genTypeName(t), x.xs)

	case reflect.Float32:
		x.linef("%s%s = (%s)(z.DecDecodeFloat32())", ptr, varname, x.genTypeName(t))
	case reflect.Float64:
		x.linef("%s%s = (%s)(r.DecodeFloat64())", ptr, varname, x.genTypeName(t))

	case reflect.Bool:
		x.linef("%s%s = (%s)(r.DecodeBool())", ptr, varname, x.genTypeName(t))
	case reflect.String:
		x.linef("%s%s = (%s)(string(r.DecodeStringAsBytes()))", ptr, varname, x.genTypeName(t))
	default:
		return false
	***REMOVED***
	return true
***REMOVED***

func (x *genRunner) decListFallback(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	if t.AssignableTo(uint8SliceTyp) ***REMOVED***
		x.line("*" + varname + " = r.DecodeBytes(*((*[]byte)(" + varname + ")), false)")
		return
	***REMOVED***
	if t.Kind() == reflect.Array && t.Elem().Kind() == reflect.Uint8 ***REMOVED***
		x.linef("r.DecodeBytes( ((*[%d]byte)(%s))[:], true)", t.Len(), varname)
		return
	***REMOVED***
	type tstruc struct ***REMOVED***
		TempVar   string
		Sfx       string
		Rand      string
		Varname   string
		CTyp      string
		Typ       string
		Immutable bool
		Size      int
	***REMOVED***
	telem := t.Elem()
	ts := tstruc***REMOVED***genTempVarPfx, x.xs, x.varsfx(), varname, x.genTypeName(t), x.genTypeName(telem), genIsImmutable(telem), int(telem.Size())***REMOVED***

	funcs := make(template.FuncMap)

	funcs["decLineVar"] = func(varname string) string ***REMOVED***
		x.decVar(varname, "", telem, false, true)
		return ""
	***REMOVED***
	funcs["var"] = func(s string) string ***REMOVED***
		return ts.TempVar + s + ts.Rand
	***REMOVED***
	funcs["xs"] = func() string ***REMOVED***
		return ts.Sfx
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
		x.decVar(varname, "", tkey, false, true)
		return ""
	***REMOVED***
	funcs["decLineVar"] = func(varname, decodedNilVarname string) string ***REMOVED***
		x.decVar(varname, decodedNilVarname, telem, false, true)
		return ""
	***REMOVED***
	funcs["var"] = func(s string) string ***REMOVED***
		return ts.TempVar + s + ts.Rand
	***REMOVED***
	funcs["xs"] = func() string ***REMOVED***
		return ts.Sfx
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
	tisfi := ti.sfiSrc // always use sequence from file. decStruct expects same thing.
	x.line("switch (" + kName + ") ***REMOVED***")
	var newbuf, nilbuf genBuf
	for _, si := range tisfi ***REMOVED***
		x.line("case \"" + si.encName + "\":")
		newbuf.reset()
		nilbuf.reset()
		varname3, t2 := x.decVarInitPtr(varname, "", t, si, &newbuf, &nilbuf)
		if len(newbuf.buf) > 0 ***REMOVED***
			x.linef("if r.TryNil() ***REMOVED*** %s ***REMOVED*** else ***REMOVED*** %s", nilbuf.buf, newbuf.buf)
		***REMOVED***
		x.decVarMain(varname3, x.varsfx(), t2.Type, false)
		if len(newbuf.buf) > 0 ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***
	***REMOVED***
	x.line("default:")
	// pass the slice here, so that the string will not escape, and maybe save allocation
	x.line("z.DecStructFieldNotFound(-1, " + kName + ")")
	x.line("***REMOVED*** // end switch " + kName)
***REMOVED***

func (x *genRunner) decStructMap(varname, lenvarname string, rtid uintptr, t reflect.Type, style genStructMapStyle) ***REMOVED***
	tpfx := genTempVarPfx
	ti := x.ti.get(rtid, t)
	i := x.varsfx()
	kName := tpfx + "s" + i

	switch style ***REMOVED***
	case genStructMapStyleLenPrefix:
		x.linef("for %sj%s := 0; %sj%s < %s; %sj%s++ ***REMOVED***", tpfx, i, tpfx, i, lenvarname, tpfx, i)
	case genStructMapStyleCheckBreak:
		x.linef("for %sj%s := 0; !z.DecCheckBreak(); %sj%s++ ***REMOVED***", tpfx, i, tpfx, i)
	default: // 0, otherwise.
		x.linef("var %shl%s bool = %s >= 0", tpfx, i, lenvarname) // has length
		x.linef("for %sj%s := 0; ; %sj%s++ ***REMOVED***", tpfx, i, tpfx, i)
		x.linef("if %shl%s ***REMOVED*** if %sj%s >= %s ***REMOVED*** break ***REMOVED***", tpfx, i, tpfx, i, lenvarname)
		x.line("***REMOVED*** else ***REMOVED*** if z.DecCheckBreak() ***REMOVED*** break ***REMOVED***; ***REMOVED***")
	***REMOVED***
	x.line("z.DecReadMapElemKey()")

	// emulate decstructfieldkey
	switch ti.keyType ***REMOVED***
	case valueTypeInt:
		x.linef("%s := z.StringView(strconv.AppendInt(z.DecScratchArrayBuffer()[:0], r.DecodeInt64(), 10))", kName)
	case valueTypeUint:
		x.linef("%s := z.StringView(strconv.AppendUint(z.DecScratchArrayBuffer()[:0], r.DecodeUint64(), 10))", kName)
	case valueTypeFloat:
		x.linef("%s := z.StringView(strconv.AppendFloat(z.DecScratchArrayBuffer()[:0], r.DecodeFloat64(), 'f', -1, 64))", kName)
	default: // string
		x.linef("%s := z.StringView(r.DecodeStringAsBytes())", kName)
	***REMOVED***

	x.line("z.DecReadMapElemValue()")
	x.decStructMapSwitch(kName, varname, rtid, t)

	x.line("***REMOVED*** // end for " + tpfx + "j" + i)
***REMOVED***

func (x *genRunner) decStructArray(varname, lenvarname, breakString string, rtid uintptr, t reflect.Type) ***REMOVED***
	tpfx := genTempVarPfx
	i := x.varsfx()
	ti := x.ti.get(rtid, t)
	tisfi := ti.sfiSrc // always use sequence from file. decStruct expects same thing.
	x.linef("var %sj%s int", tpfx, i)
	x.linef("var %sb%s bool", tpfx, i)                        // break
	x.linef("var %shl%s bool = %s >= 0", tpfx, i, lenvarname) // has length
	if !genDecStructArrayInlineLoopCheck ***REMOVED***
		x.linef("var %sfn%s = func() bool ***REMOVED*** ", tpfx, i)
		x.linef("%sj%s++; if %shl%s ***REMOVED*** %sb%s = %sj%s > %s ***REMOVED*** else ***REMOVED*** %sb%s = z.DecCheckBreak() ***REMOVED***;",
			tpfx, i, tpfx, i, tpfx, i,
			tpfx, i, lenvarname, tpfx, i)
		x.linef("if %sb%s ***REMOVED*** z.DecReadArrayEnd(); return true ***REMOVED***; return false", tpfx, i)
		x.linef("***REMOVED*** // end func %sfn%s", tpfx, i)
	***REMOVED***
	var newbuf, nilbuf genBuf
	for _, si := range tisfi ***REMOVED***
		if genDecStructArrayInlineLoopCheck ***REMOVED***
			x.linef("%sj%s++; if %shl%s ***REMOVED*** %sb%s = %sj%s > %s ***REMOVED*** else ***REMOVED*** %sb%s = z.DecCheckBreak() ***REMOVED***",
				tpfx, i, tpfx, i, tpfx, i,
				tpfx, i, lenvarname, tpfx, i)
			x.linef("if %sb%s ***REMOVED*** z.DecReadArrayEnd(); %s ***REMOVED***", tpfx, i, breakString)
		***REMOVED*** else ***REMOVED***
			x.linef("if %sfn%s() ***REMOVED*** %s ***REMOVED***", tpfx, i, breakString)
		***REMOVED***
		x.line("z.DecReadArrayElem()")
		newbuf.reset()
		nilbuf.reset()
		varname3, t2 := x.decVarInitPtr(varname, "", t, si, &newbuf, &nilbuf)
		if len(newbuf.buf) > 0 ***REMOVED***
			x.linef("if r.TryNil() ***REMOVED*** %s ***REMOVED*** else ***REMOVED*** %s", nilbuf.buf, newbuf.buf)
		***REMOVED***
		x.decVarMain(varname3, x.varsfx(), t2.Type, false)
		if len(newbuf.buf) > 0 ***REMOVED***
			x.line("***REMOVED***")
		***REMOVED***
	***REMOVED***
	// read remaining values and throw away.
	x.line("for ***REMOVED***")
	x.linef("%sj%s++; if %shl%s ***REMOVED*** %sb%s = %sj%s > %s ***REMOVED*** else ***REMOVED*** %sb%s = z.DecCheckBreak() ***REMOVED***",
		tpfx, i, tpfx, i, tpfx, i,
		tpfx, i, lenvarname, tpfx, i)
	x.linef("if %sb%s ***REMOVED*** break ***REMOVED***", tpfx, i)
	x.line("z.DecReadArrayElem()")
	x.linef(`z.DecStructFieldNotFound(%sj%s - 1, "")`, tpfx, i)
	x.line("***REMOVED***")
***REMOVED***

func (x *genRunner) decStruct(varname string, rtid uintptr, t reflect.Type) ***REMOVED***
	// varname MUST be a ptr, or a struct field or a slice element.
	i := x.varsfx()
	x.linef("%sct%s := r.ContainerType()", genTempVarPfx, i)
	x.linef("if %sct%s == codecSelferValueTypeNil%s ***REMOVED***", genTempVarPfx, i, x.xs)
	x.linef("*(%s) = %s***REMOVED******REMOVED***", varname, x.genTypeName(t))
	x.linef("***REMOVED*** else if %sct%s == codecSelferValueTypeMap%s ***REMOVED***", genTempVarPfx, i, x.xs)
	x.line(genTempVarPfx + "l" + i + " := z.DecReadMapStart()")
	x.linef("if %sl%s == 0 ***REMOVED***", genTempVarPfx, i)
	if genUseOneFunctionForDecStructMap ***REMOVED***
		x.line("***REMOVED*** else ***REMOVED*** ")
		x.linef("%s.codecDecodeSelfFromMap(%sl%s, d)", varname, genTempVarPfx, i)
	***REMOVED*** else ***REMOVED***
		x.line("***REMOVED*** else if " + genTempVarPfx + "l" + i + " > 0 ***REMOVED*** ")
		x.line(varname + ".codecDecodeSelfFromMapLenPrefix(" + genTempVarPfx + "l" + i + ", d)")
		x.line("***REMOVED*** else ***REMOVED***")
		x.line(varname + ".codecDecodeSelfFromMapCheckBreak(" + genTempVarPfx + "l" + i + ", d)")
	***REMOVED***
	x.line("***REMOVED***")
	x.line("z.DecReadMapEnd()")

	// else if container is array
	x.linef("***REMOVED*** else if %sct%s == codecSelferValueTypeArray%s ***REMOVED***", genTempVarPfx, i, x.xs)
	x.line(genTempVarPfx + "l" + i + " := z.DecReadArrayStart()")
	x.linef("if %sl%s != 0 ***REMOVED***", genTempVarPfx, i)
	x.linef("%s.codecDecodeSelfFromArray(%sl%s, d)", varname, genTempVarPfx, i)
	x.line("***REMOVED***")
	x.line("z.DecReadArrayEnd()")
	// else panic
	x.line("***REMOVED*** else ***REMOVED*** ")
	x.line("panic(errCodecSelferOnlyMapOrArrayEncodeToStruct" + x.xs + ")")
	x.line("***REMOVED*** ")
***REMOVED***

// --------

type fastpathGenV struct ***REMOVED***
	// fastpathGenV is either a primitive (Primitive != "") or a map (MapKey != "") or a slice
	MapKey      string
	Elem        string
	Primitive   string
	Size        int
	NoCanonical bool
***REMOVED***

func (x *genRunner) newFastpathGenV(t reflect.Type) (v fastpathGenV) ***REMOVED***
	v.NoCanonical = !genFastpathCanonical
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
		panic("unexpected type for newFastpathGenV. Requires map or slice type")
	***REMOVED***
	return
***REMOVED***

func (x *fastpathGenV) MethodNamePfx(prefix string, prim bool) string ***REMOVED***
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
		// HACK: always handle vendoring. It should be typically on in go 1.6, 1.7
		s = genStripVendor(s)
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
	case "[]byte", "[]uint8", "bytes":
		return "Bytes"
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
	Version int
	Values  []fastpathGenV
***REMOVED***

func (x genInternal) FastpathLen() (l int) ***REMOVED***
	for _, v := range x.Values ***REMOVED***
		if v.Primitive == "" && !(v.MapKey == "" && v.Elem == "uint8") ***REMOVED***
			l++
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func genInternalZeroValue(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "interface***REMOVED******REMOVED***", "interface ***REMOVED******REMOVED***":
		return "nil"
	case "[]byte", "[]uint8", "bytes":
		return "nil"
	case "bool":
		return "false"
	case "string":
		return `""`
	default:
		return "0"
	***REMOVED***
***REMOVED***

var genInternalNonZeroValueIdx [6]uint64
var genInternalNonZeroValueStrs = [...][6]string***REMOVED***
	***REMOVED***`"string-is-an-interface-1"`, "true", `"some-string-1"`, `[]byte("some-string-1")`, "11.1", "111"***REMOVED***,
	***REMOVED***`"string-is-an-interface-2"`, "false", `"some-string-2"`, `[]byte("some-string-2")`, "22.2", "77"***REMOVED***,
	***REMOVED***`"string-is-an-interface-3"`, "true", `"some-string-3"`, `[]byte("some-string-3")`, "33.3e3", "127"***REMOVED***,
***REMOVED***

// Note: last numbers must be in range: 0-127 (as they may be put into a int8, uint8, etc)

func genInternalNonZeroValue(s string) string ***REMOVED***
	var i int
	switch s ***REMOVED***
	case "interface***REMOVED******REMOVED***", "interface ***REMOVED******REMOVED***":
		i = 0
	case "bool":
		i = 1
	case "string":
		i = 2
	case "bytes", "[]byte", "[]uint8":
		i = 3
	case "float32", "float64", "float", "double":
		i = 4
	default:
		i = 5
	***REMOVED***
	genInternalNonZeroValueIdx[i]++
	idx := genInternalNonZeroValueIdx[i]
	slen := uint64(len(genInternalNonZeroValueStrs))
	return genInternalNonZeroValueStrs[idx%slen][i] // return string, to remove ambiguity
***REMOVED***

func genInternalEncCommandAsString(s string, vname string) string ***REMOVED***
	switch s ***REMOVED***
	case "uint64":
		return "e.e.EncodeUint(" + vname + ")"
	case "uint", "uint8", "uint16", "uint32":
		return "e.e.EncodeUint(uint64(" + vname + "))"
	case "int64":
		return "e.e.EncodeInt(" + vname + ")"
	case "int", "int8", "int16", "int32":
		return "e.e.EncodeInt(int64(" + vname + "))"
	case "[]byte", "[]uint8", "bytes":
		return "e.e.EncodeStringBytesRaw(" + vname + ")"
	case "string":
		return "e.e.EncodeString(" + vname + ")"
	case "float32":
		return "e.e.EncodeFloat32(" + vname + ")"
	case "float64":
		return "e.e.EncodeFloat64(" + vname + ")"
	case "bool":
		return "e.e.EncodeBool(" + vname + ")"
	// case "symbol":
	// 	return "e.e.EncodeSymbol(" + vname + ")"
	default:
		return "e.encode(" + vname + ")"
	***REMOVED***
***REMOVED***

func genInternalDecCommandAsString(s string) string ***REMOVED***
	switch s ***REMOVED***
	case "uint":
		return "uint(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "uint8":
		return "uint8(chkOvf.UintV(d.d.DecodeUint64(), 8))"
	case "uint16":
		return "uint16(chkOvf.UintV(d.d.DecodeUint64(), 16))"
	case "uint32":
		return "uint32(chkOvf.UintV(d.d.DecodeUint64(), 32))"
	case "uint64":
		return "d.d.DecodeUint64()"
	case "uintptr":
		return "uintptr(chkOvf.UintV(d.d.DecodeUint64(), uintBitsize))"
	case "int":
		return "int(chkOvf.IntV(d.d.DecodeInt64(), intBitsize))"
	case "int8":
		return "int8(chkOvf.IntV(d.d.DecodeInt64(), 8))"
	case "int16":
		return "int16(chkOvf.IntV(d.d.DecodeInt64(), 16))"
	case "int32":
		return "int32(chkOvf.IntV(d.d.DecodeInt64(), 32))"
	case "int64":
		return "d.d.DecodeInt64()"

	case "string":
		return "string(d.d.DecodeStringAsBytes())"
	case "[]byte", "[]uint8", "bytes":
		return "d.d.DecodeBytes(nil, false)"
	case "float32":
		return "float32(d.decodeFloat32())"
	case "float64":
		return "d.d.DecodeFloat64()"
	case "bool":
		return "d.d.DecodeBool()"
	default:
		panic(errors.New("gen internal: unknown type for decode: " + s))
	***REMOVED***
***REMOVED***

func genInternalSortType(s string, elem bool) string ***REMOVED***
	for _, v := range [...]string***REMOVED***
		"int",
		"uint",
		"float",
		"bool",
		"string",
		"bytes", "[]uint8", "[]byte",
	***REMOVED*** ***REMOVED***
		if v == "[]byte" || v == "[]uint8" ***REMOVED***
			v = "bytes"
		***REMOVED***
		if strings.HasPrefix(s, v) ***REMOVED***
			if v == "int" || v == "uint" || v == "float" ***REMOVED***
				v += "64"
			***REMOVED***
			if elem ***REMOVED***
				return v
			***REMOVED***
			return v + "Slice"
		***REMOVED***
	***REMOVED***
	panic("sorttype: unexpected type: " + s)
***REMOVED***

func genStripVendor(s string) string ***REMOVED***
	// HACK: Misbehaviour occurs in go 1.5. May have to re-visit this later.
	// if s contains /vendor/ OR startsWith vendor/, then return everything after it.
	const vendorStart = "vendor/"
	const vendorInline = "/vendor/"
	if i := strings.LastIndex(s, vendorInline); i >= 0 ***REMOVED***
		s = s[i+len(vendorInline):]
	***REMOVED*** else if strings.HasPrefix(s, vendorStart) ***REMOVED***
		s = s[len(vendorStart):]
	***REMOVED***
	return s
***REMOVED***

// var genInternalMu sync.Mutex
var genInternalV = genInternal***REMOVED***Version: genVersion***REMOVED***
var genInternalTmplFuncs template.FuncMap
var genInternalOnce sync.Once

func genInternalInit() ***REMOVED***
	wordSizeBytes := int(intBitsize) / 8

	typesizes := map[string]int***REMOVED***
		"interface***REMOVED******REMOVED***": 2 * wordSizeBytes,
		"string":      2 * wordSizeBytes,
		"[]byte":      3 * wordSizeBytes,
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

	// keep as slice, so it is in specific iteration order.
	// Initial order was uint64, string, interface***REMOVED******REMOVED***, int, int64, ...

	var types = [...]string***REMOVED***
		"interface***REMOVED******REMOVED***",
		"string",
		"[]byte",
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

	var primitivetypes, slicetypes, mapkeytypes, mapvaltypes []string

	primitivetypes = types[:]
	slicetypes = types[:]
	mapkeytypes = types[:]
	mapvaltypes = types[:]

	if genFastpathTrimTypes ***REMOVED***
		slicetypes = []string***REMOVED***
			"interface***REMOVED******REMOVED***",
			"string",
			"[]byte",
			"float32",
			"float64",
			"uint",
			// "uint8", // no need for fastpath of []uint8, as it is handled specially
			"uint16",
			"uint32",
			"uint64",
			// "uintptr",
			"int",
			"int8",
			"int16",
			"int32",
			"int64",
			"bool",
		***REMOVED***

		mapkeytypes = []string***REMOVED***
			//"interface***REMOVED******REMOVED***",
			"string",
			//"[]byte",
			//"float32",
			//"float64",
			"uint",
			"uint8",
			//"uint16",
			//"uint32",
			"uint64",
			//"uintptr",
			"int",
			//"int8",
			//"int16",
			//"int32",
			"int64",
			// "bool",
		***REMOVED***

		mapvaltypes = []string***REMOVED***
			"interface***REMOVED******REMOVED***",
			"string",
			"[]byte",
			"uint",
			"uint8",
			//"uint16",
			//"uint32",
			"uint64",
			// "uintptr",
			"int",
			//"int8",
			//"int16",
			//"int32",
			"int64",
			"float32",
			"float64",
			"bool",
		***REMOVED***
	***REMOVED***

	// var mapkeytypes [len(&types) - 1]string // skip bool
	// copy(mapkeytypes[:], types[:])

	// var mb []byte
	// mb = append(mb, '|')
	// for _, s := range mapkeytypes ***REMOVED***
	// 	mb = append(mb, s...)
	// 	mb = append(mb, '|')
	// ***REMOVED***
	// var mapkeytypestr = string(mb)

	var gt = genInternal***REMOVED***Version: genVersion***REMOVED***

	// For each slice or map type, there must be a (symmetrical) Encode and Decode fast-path function

	for _, s := range primitivetypes ***REMOVED***
		gt.Values = append(gt.Values,
			fastpathGenV***REMOVED***Primitive: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical***REMOVED***)
	***REMOVED***
	for _, s := range slicetypes ***REMOVED***
		// if s != "uint8" ***REMOVED*** // do not generate fast path for slice of bytes. Treat specially already.
		// 	gt.Values = append(gt.Values, fastpathGenV***REMOVED***Elem: s, Size: typesizes[s]***REMOVED***)
		// ***REMOVED***
		gt.Values = append(gt.Values,
			fastpathGenV***REMOVED***Elem: s, Size: typesizes[s], NoCanonical: !genFastpathCanonical***REMOVED***)
	***REMOVED***
	for _, s := range mapkeytypes ***REMOVED***
		// if _, ok := typesizes[s]; !ok ***REMOVED***
		// if strings.Contains(mapkeytypestr, "|"+s+"|") ***REMOVED***
		// 	gt.Values = append(gt.Values, fastpathGenV***REMOVED***MapKey: s, Elem: s, Size: 2 * typesizes[s]***REMOVED***)
		// ***REMOVED***
		for _, ms := range mapvaltypes ***REMOVED***
			gt.Values = append(gt.Values,
				fastpathGenV***REMOVED***MapKey: s, Elem: ms, Size: typesizes[s] + typesizes[ms], NoCanonical: !genFastpathCanonical***REMOVED***)
		***REMOVED***
	***REMOVED***

	funcs := make(template.FuncMap)
	// funcs["haspfx"] = strings.HasPrefix
	funcs["encmd"] = genInternalEncCommandAsString
	funcs["decmd"] = genInternalDecCommandAsString
	funcs["zerocmd"] = genInternalZeroValue
	funcs["nonzerocmd"] = genInternalNonZeroValue
	funcs["hasprefix"] = strings.HasPrefix
	funcs["sorttype"] = genInternalSortType

	genInternalV = gt
	genInternalTmplFuncs = funcs
***REMOVED***

// genInternalGoFile is used to generate source files from templates.
// It is run by the program author alone.
// Unfortunately, it has to be exported so that it can be called from a command line tool.
// *** DO NOT USE ***
func genInternalGoFile(r io.Reader, w io.Writer) (err error) ***REMOVED***
	genInternalOnce.Do(genInternalInit)

	gt := genInternalV

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
