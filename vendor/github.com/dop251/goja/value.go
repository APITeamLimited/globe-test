package goja

import (
	"hash/maphash"
	"math"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/dop251/goja/ftoa"
	"github.com/dop251/goja/unistring"
)

var (
	// Not goroutine-safe, do not use for anything other than package level init
	pkgHasher maphash.Hash

	hashFalse = randomHash()
	hashTrue  = randomHash()
	hashNull  = randomHash()
	hashUndef = randomHash()
)

// Not goroutine-safe, do not use for anything other than package level init
func randomHash() uint64 ***REMOVED***
	pkgHasher.WriteByte(0)
	return pkgHasher.Sum64()
***REMOVED***

var (
	valueFalse    Value = valueBool(false)
	valueTrue     Value = valueBool(true)
	_null         Value = valueNull***REMOVED******REMOVED***
	_NaN          Value = valueFloat(math.NaN())
	_positiveInf  Value = valueFloat(math.Inf(+1))
	_negativeInf  Value = valueFloat(math.Inf(-1))
	_positiveZero Value = valueInt(0)
	negativeZero        = math.Float64frombits(0 | (1 << 63))
	_negativeZero Value = valueFloat(negativeZero)
	_epsilon            = valueFloat(2.2204460492503130808472633361816e-16)
	_undefined    Value = valueUndefined***REMOVED******REMOVED***
)

var (
	reflectTypeInt    = reflect.TypeOf(int64(0))
	reflectTypeBool   = reflect.TypeOf(false)
	reflectTypeNil    = reflect.TypeOf(nil)
	reflectTypeFloat  = reflect.TypeOf(float64(0))
	reflectTypeMap    = reflect.TypeOf(map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	reflectTypeArray  = reflect.TypeOf([]interface***REMOVED******REMOVED******REMOVED******REMOVED***)
	reflectTypeString = reflect.TypeOf("")
)

var intCache [256]Value

type Value interface ***REMOVED***
	ToInteger() int64
	toString() valueString
	string() unistring.String
	ToString() Value
	String() string
	ToFloat() float64
	ToNumber() Value
	ToBoolean() bool
	ToObject(*Runtime) *Object
	SameAs(Value) bool
	Equals(Value) bool
	StrictEquals(Value) bool
	Export() interface***REMOVED******REMOVED***
	ExportType() reflect.Type

	baseObject(r *Runtime) *Object

	hash(hasher *maphash.Hash) uint64
***REMOVED***

type valueContainer interface ***REMOVED***
	toValue(*Runtime) Value
***REMOVED***

type typeError string
type rangeError string

type valueInt int64
type valueFloat float64
type valueBool bool
type valueNull struct***REMOVED******REMOVED***
type valueUndefined struct ***REMOVED***
	valueNull
***REMOVED***
type valueSymbol struct ***REMOVED***
	h    uintptr
	desc valueString
***REMOVED***

type valueUnresolved struct ***REMOVED***
	r   *Runtime
	ref unistring.String
***REMOVED***

type memberUnresolved struct ***REMOVED***
	valueUnresolved
***REMOVED***

type valueProperty struct ***REMOVED***
	value        Value
	writable     bool
	configurable bool
	enumerable   bool
	accessor     bool
	getterFunc   *Object
	setterFunc   *Object
***REMOVED***

func propGetter(o Value, v Value, r *Runtime) *Object ***REMOVED***
	if v == _undefined ***REMOVED***
		return nil
	***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if _, ok := obj.self.assertCallable(); ok ***REMOVED***
			return obj
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Getter must be a function: %s", v.toString())
	return nil
***REMOVED***

func propSetter(o Value, v Value, r *Runtime) *Object ***REMOVED***
	if v == _undefined ***REMOVED***
		return nil
	***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if _, ok := obj.self.assertCallable(); ok ***REMOVED***
			return obj
		***REMOVED***
	***REMOVED***
	r.typeErrorResult(true, "Setter must be a function: %s", v.toString())
	return nil
***REMOVED***

func fToStr(num float64, mode ftoa.FToStrMode, prec int) string ***REMOVED***
	var buf1 [128]byte
	return string(ftoa.FToStr(num, mode, prec, buf1[:0]))
***REMOVED***

func (i valueInt) ToInteger() int64 ***REMOVED***
	return int64(i)
***REMOVED***

func (i valueInt) toString() valueString ***REMOVED***
	return asciiString(i.String())
***REMOVED***

func (i valueInt) string() unistring.String ***REMOVED***
	return unistring.String(i.String())
***REMOVED***

func (i valueInt) ToString() Value ***REMOVED***
	return i
***REMOVED***

func (i valueInt) String() string ***REMOVED***
	return strconv.FormatInt(int64(i), 10)
***REMOVED***

func (i valueInt) ToFloat() float64 ***REMOVED***
	return float64(i)
***REMOVED***

func (i valueInt) ToBoolean() bool ***REMOVED***
	return i != 0
***REMOVED***

func (i valueInt) ToObject(r *Runtime) *Object ***REMOVED***
	return r.newPrimitiveObject(i, r.global.NumberPrototype, classNumber)
***REMOVED***

func (i valueInt) ToNumber() Value ***REMOVED***
	return i
***REMOVED***

func (i valueInt) SameAs(other Value) bool ***REMOVED***
	return i == other
***REMOVED***

func (i valueInt) Equals(other Value) bool ***REMOVED***
	switch o := other.(type) ***REMOVED***
	case valueInt:
		return i == o
	case valueFloat:
		return float64(i) == float64(o)
	case valueString:
		return o.ToNumber().Equals(i)
	case valueBool:
		return int64(i) == o.ToInteger()
	case *Object:
		return i.Equals(o.toPrimitiveNumber())
	***REMOVED***

	return false
***REMOVED***

func (i valueInt) StrictEquals(other Value) bool ***REMOVED***
	switch o := other.(type) ***REMOVED***
	case valueInt:
		return i == o
	case valueFloat:
		return float64(i) == float64(o)
	***REMOVED***

	return false
***REMOVED***

func (i valueInt) baseObject(r *Runtime) *Object ***REMOVED***
	return r.global.NumberPrototype
***REMOVED***

func (i valueInt) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return int64(i)
***REMOVED***

func (i valueInt) ExportType() reflect.Type ***REMOVED***
	return reflectTypeInt
***REMOVED***

func (i valueInt) hash(*maphash.Hash) uint64 ***REMOVED***
	return uint64(i)
***REMOVED***

func (b valueBool) ToInteger() int64 ***REMOVED***
	if b ***REMOVED***
		return 1
	***REMOVED***
	return 0
***REMOVED***

func (b valueBool) toString() valueString ***REMOVED***
	if b ***REMOVED***
		return stringTrue
	***REMOVED***
	return stringFalse
***REMOVED***

func (b valueBool) ToString() Value ***REMOVED***
	return b
***REMOVED***

func (b valueBool) String() string ***REMOVED***
	if b ***REMOVED***
		return "true"
	***REMOVED***
	return "false"
***REMOVED***

func (b valueBool) string() unistring.String ***REMOVED***
	return unistring.String(b.String())
***REMOVED***

func (b valueBool) ToFloat() float64 ***REMOVED***
	if b ***REMOVED***
		return 1.0
	***REMOVED***
	return 0
***REMOVED***

func (b valueBool) ToBoolean() bool ***REMOVED***
	return bool(b)
***REMOVED***

func (b valueBool) ToObject(r *Runtime) *Object ***REMOVED***
	return r.newPrimitiveObject(b, r.global.BooleanPrototype, "Boolean")
***REMOVED***

func (b valueBool) ToNumber() Value ***REMOVED***
	if b ***REMOVED***
		return valueInt(1)
	***REMOVED***
	return valueInt(0)
***REMOVED***

func (b valueBool) SameAs(other Value) bool ***REMOVED***
	if other, ok := other.(valueBool); ok ***REMOVED***
		return b == other
	***REMOVED***
	return false
***REMOVED***

func (b valueBool) Equals(other Value) bool ***REMOVED***
	if o, ok := other.(valueBool); ok ***REMOVED***
		return b == o
	***REMOVED***

	if b ***REMOVED***
		return other.Equals(intToValue(1))
	***REMOVED*** else ***REMOVED***
		return other.Equals(intToValue(0))
	***REMOVED***

***REMOVED***

func (b valueBool) StrictEquals(other Value) bool ***REMOVED***
	if other, ok := other.(valueBool); ok ***REMOVED***
		return b == other
	***REMOVED***
	return false
***REMOVED***

func (b valueBool) baseObject(r *Runtime) *Object ***REMOVED***
	return r.global.BooleanPrototype
***REMOVED***

func (b valueBool) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return bool(b)
***REMOVED***

func (b valueBool) ExportType() reflect.Type ***REMOVED***
	return reflectTypeBool
***REMOVED***

func (b valueBool) hash(*maphash.Hash) uint64 ***REMOVED***
	if b ***REMOVED***
		return hashTrue
	***REMOVED***

	return hashFalse
***REMOVED***

func (n valueNull) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (n valueNull) toString() valueString ***REMOVED***
	return stringNull
***REMOVED***

func (n valueNull) string() unistring.String ***REMOVED***
	return stringNull.string()
***REMOVED***

func (n valueNull) ToString() Value ***REMOVED***
	return n
***REMOVED***

func (n valueNull) String() string ***REMOVED***
	return "null"
***REMOVED***

func (u valueUndefined) toString() valueString ***REMOVED***
	return stringUndefined
***REMOVED***

func (u valueUndefined) ToString() Value ***REMOVED***
	return u
***REMOVED***

func (u valueUndefined) String() string ***REMOVED***
	return "undefined"
***REMOVED***

func (u valueUndefined) string() unistring.String ***REMOVED***
	return "undefined"
***REMOVED***

func (u valueUndefined) ToNumber() Value ***REMOVED***
	return _NaN
***REMOVED***

func (u valueUndefined) SameAs(other Value) bool ***REMOVED***
	_, same := other.(valueUndefined)
	return same
***REMOVED***

func (u valueUndefined) StrictEquals(other Value) bool ***REMOVED***
	_, same := other.(valueUndefined)
	return same
***REMOVED***

func (u valueUndefined) ToFloat() float64 ***REMOVED***
	return math.NaN()
***REMOVED***

func (u valueUndefined) hash(*maphash.Hash) uint64 ***REMOVED***
	return hashUndef
***REMOVED***

func (n valueNull) ToFloat() float64 ***REMOVED***
	return 0
***REMOVED***

func (n valueNull) ToBoolean() bool ***REMOVED***
	return false
***REMOVED***

func (n valueNull) ToObject(r *Runtime) *Object ***REMOVED***
	r.typeErrorResult(true, "Cannot convert undefined or null to object")
	return nil
	//return r.newObject()
***REMOVED***

func (n valueNull) ToNumber() Value ***REMOVED***
	return intToValue(0)
***REMOVED***

func (n valueNull) SameAs(other Value) bool ***REMOVED***
	_, same := other.(valueNull)
	return same
***REMOVED***

func (n valueNull) Equals(other Value) bool ***REMOVED***
	switch other.(type) ***REMOVED***
	case valueUndefined, valueNull:
		return true
	***REMOVED***
	return false
***REMOVED***

func (n valueNull) StrictEquals(other Value) bool ***REMOVED***
	_, same := other.(valueNull)
	return same
***REMOVED***

func (n valueNull) baseObject(*Runtime) *Object ***REMOVED***
	return nil
***REMOVED***

func (n valueNull) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (n valueNull) ExportType() reflect.Type ***REMOVED***
	return reflectTypeNil
***REMOVED***

func (n valueNull) hash(*maphash.Hash) uint64 ***REMOVED***
	return hashNull
***REMOVED***

func (p *valueProperty) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (p *valueProperty) toString() valueString ***REMOVED***
	return stringEmpty
***REMOVED***

func (p *valueProperty) string() unistring.String ***REMOVED***
	return ""
***REMOVED***

func (p *valueProperty) ToString() Value ***REMOVED***
	return _undefined
***REMOVED***

func (p *valueProperty) String() string ***REMOVED***
	return ""
***REMOVED***

func (p *valueProperty) ToFloat() float64 ***REMOVED***
	return math.NaN()
***REMOVED***

func (p *valueProperty) ToBoolean() bool ***REMOVED***
	return false
***REMOVED***

func (p *valueProperty) ToObject(*Runtime) *Object ***REMOVED***
	return nil
***REMOVED***

func (p *valueProperty) ToNumber() Value ***REMOVED***
	return nil
***REMOVED***

func (p *valueProperty) isWritable() bool ***REMOVED***
	return p.writable || p.setterFunc != nil
***REMOVED***

func (p *valueProperty) get(this Value) Value ***REMOVED***
	if p.getterFunc == nil ***REMOVED***
		if p.value != nil ***REMOVED***
			return p.value
		***REMOVED***
		return _undefined
	***REMOVED***
	call, _ := p.getterFunc.self.assertCallable()
	return call(FunctionCall***REMOVED***
		This: this,
	***REMOVED***)
***REMOVED***

func (p *valueProperty) set(this, v Value) ***REMOVED***
	if p.setterFunc == nil ***REMOVED***
		p.value = v
		return
	***REMOVED***
	call, _ := p.setterFunc.self.assertCallable()
	call(FunctionCall***REMOVED***
		This:      this,
		Arguments: []Value***REMOVED***v***REMOVED***,
	***REMOVED***)
***REMOVED***

func (p *valueProperty) SameAs(other Value) bool ***REMOVED***
	if otherProp, ok := other.(*valueProperty); ok ***REMOVED***
		return p == otherProp
	***REMOVED***
	return false
***REMOVED***

func (p *valueProperty) Equals(Value) bool ***REMOVED***
	return false
***REMOVED***

func (p *valueProperty) StrictEquals(Value) bool ***REMOVED***
	return false
***REMOVED***

func (p *valueProperty) baseObject(r *Runtime) *Object ***REMOVED***
	r.typeErrorResult(true, "BUG: baseObject() is called on valueProperty") // TODO error message
	return nil
***REMOVED***

func (p *valueProperty) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	panic("Cannot export valueProperty")
***REMOVED***

func (p *valueProperty) ExportType() reflect.Type ***REMOVED***
	panic("Cannot export valueProperty")
***REMOVED***

func (p *valueProperty) hash(*maphash.Hash) uint64 ***REMOVED***
	panic("valueProperty should never be used in maps or sets")
***REMOVED***

func floatToIntClip(n float64) int64 ***REMOVED***
	switch ***REMOVED***
	case math.IsNaN(n):
		return 0
	case n >= math.MaxInt64:
		return math.MaxInt64
	case n <= math.MinInt64:
		return math.MinInt64
	***REMOVED***
	return int64(n)
***REMOVED***

func (f valueFloat) ToInteger() int64 ***REMOVED***
	return floatToIntClip(float64(f))
***REMOVED***

func (f valueFloat) toString() valueString ***REMOVED***
	return asciiString(f.String())
***REMOVED***

func (f valueFloat) string() unistring.String ***REMOVED***
	return unistring.String(f.String())
***REMOVED***

func (f valueFloat) ToString() Value ***REMOVED***
	return f
***REMOVED***

func (f valueFloat) String() string ***REMOVED***
	return fToStr(float64(f), ftoa.ModeStandard, 0)
***REMOVED***

func (f valueFloat) ToFloat() float64 ***REMOVED***
	return float64(f)
***REMOVED***

func (f valueFloat) ToBoolean() bool ***REMOVED***
	return float64(f) != 0.0 && !math.IsNaN(float64(f))
***REMOVED***

func (f valueFloat) ToObject(r *Runtime) *Object ***REMOVED***
	return r.newPrimitiveObject(f, r.global.NumberPrototype, "Number")
***REMOVED***

func (f valueFloat) ToNumber() Value ***REMOVED***
	return f
***REMOVED***

func (f valueFloat) SameAs(other Value) bool ***REMOVED***
	switch o := other.(type) ***REMOVED***
	case valueFloat:
		this := float64(f)
		o1 := float64(o)
		if math.IsNaN(this) && math.IsNaN(o1) ***REMOVED***
			return true
		***REMOVED*** else ***REMOVED***
			ret := this == o1
			if ret && this == 0 ***REMOVED***
				ret = math.Signbit(this) == math.Signbit(o1)
			***REMOVED***
			return ret
		***REMOVED***
	case valueInt:
		this := float64(f)
		ret := this == float64(o)
		if ret && this == 0 ***REMOVED***
			ret = !math.Signbit(this)
		***REMOVED***
		return ret
	***REMOVED***

	return false
***REMOVED***

func (f valueFloat) Equals(other Value) bool ***REMOVED***
	switch o := other.(type) ***REMOVED***
	case valueFloat:
		return f == o
	case valueInt:
		return float64(f) == float64(o)
	case valueString, valueBool:
		return float64(f) == o.ToFloat()
	case *Object:
		return f.Equals(o.toPrimitiveNumber())
	***REMOVED***

	return false
***REMOVED***

func (f valueFloat) StrictEquals(other Value) bool ***REMOVED***
	switch o := other.(type) ***REMOVED***
	case valueFloat:
		return f == o
	case valueInt:
		return float64(f) == float64(o)
	***REMOVED***

	return false
***REMOVED***

func (f valueFloat) baseObject(r *Runtime) *Object ***REMOVED***
	return r.global.NumberPrototype
***REMOVED***

func (f valueFloat) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return float64(f)
***REMOVED***

func (f valueFloat) ExportType() reflect.Type ***REMOVED***
	return reflectTypeFloat
***REMOVED***

func (f valueFloat) hash(*maphash.Hash) uint64 ***REMOVED***
	if f == _negativeZero ***REMOVED***
		return 0
	***REMOVED***
	return math.Float64bits(float64(f))
***REMOVED***

func (o *Object) ToInteger() int64 ***REMOVED***
	return o.toPrimitiveNumber().ToNumber().ToInteger()
***REMOVED***

func (o *Object) toString() valueString ***REMOVED***
	return o.toPrimitiveString().toString()
***REMOVED***

func (o *Object) string() unistring.String ***REMOVED***
	return o.toPrimitiveString().string()
***REMOVED***

func (o *Object) ToString() Value ***REMOVED***
	return o.toPrimitiveString().ToString()
***REMOVED***

func (o *Object) String() string ***REMOVED***
	return o.toPrimitiveString().String()
***REMOVED***

func (o *Object) ToFloat() float64 ***REMOVED***
	return o.toPrimitiveNumber().ToFloat()
***REMOVED***

func (o *Object) ToBoolean() bool ***REMOVED***
	return true
***REMOVED***

func (o *Object) ToObject(*Runtime) *Object ***REMOVED***
	return o
***REMOVED***

func (o *Object) ToNumber() Value ***REMOVED***
	return o.toPrimitiveNumber().ToNumber()
***REMOVED***

func (o *Object) SameAs(other Value) bool ***REMOVED***
	if other, ok := other.(*Object); ok ***REMOVED***
		return o == other
	***REMOVED***
	return false
***REMOVED***

func (o *Object) Equals(other Value) bool ***REMOVED***
	if other, ok := other.(*Object); ok ***REMOVED***
		return o == other || o.self.equal(other.self)
	***REMOVED***

	switch o1 := other.(type) ***REMOVED***
	case valueInt, valueFloat, valueString:
		return o.toPrimitive().Equals(other)
	case valueBool:
		return o.Equals(o1.ToNumber())
	***REMOVED***

	return false
***REMOVED***

func (o *Object) StrictEquals(other Value) bool ***REMOVED***
	if other, ok := other.(*Object); ok ***REMOVED***
		return o == other || o.self.equal(other.self)
	***REMOVED***
	return false
***REMOVED***

func (o *Object) baseObject(*Runtime) *Object ***REMOVED***
	return o
***REMOVED***

func (o *Object) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return o.self.export(&objectExportCtx***REMOVED******REMOVED***)
***REMOVED***

func (o *Object) ExportType() reflect.Type ***REMOVED***
	return o.self.exportType()
***REMOVED***

func (o *Object) hash(*maphash.Hash) uint64 ***REMOVED***
	return o.getId()
***REMOVED***

func (o *Object) Get(name string) Value ***REMOVED***
	return o.self.getStr(unistring.NewFromString(name), nil)
***REMOVED***

func (o *Object) Keys() (keys []string) ***REMOVED***
	names := o.self.ownKeys(false, nil)
	keys = make([]string, 0, len(names))
	for _, name := range names ***REMOVED***
		keys = append(keys, name.String())
	***REMOVED***

	return
***REMOVED***

// DefineDataProperty is a Go equivalent of Object.defineProperty(o, name, ***REMOVED***value: value, writable: writable,
// configurable: configurable, enumerable: enumerable***REMOVED***)
func (o *Object) DefineDataProperty(name string, value Value, writable, configurable, enumerable Flag) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.defineOwnPropertyStr(unistring.NewFromString(name), PropertyDescriptor***REMOVED***
			Value:        value,
			Writable:     writable,
			Configurable: configurable,
			Enumerable:   enumerable,
		***REMOVED***, true)
	***REMOVED***)
***REMOVED***

// DefineAccessorProperty is a Go equivalent of Object.defineProperty(o, name, ***REMOVED***get: getter, set: setter,
// configurable: configurable, enumerable: enumerable***REMOVED***)
func (o *Object) DefineAccessorProperty(name string, getter, setter Value, configurable, enumerable Flag) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.defineOwnPropertyStr(unistring.NewFromString(name), PropertyDescriptor***REMOVED***
			Getter:       getter,
			Setter:       setter,
			Configurable: configurable,
			Enumerable:   enumerable,
		***REMOVED***, true)
	***REMOVED***)
***REMOVED***

func (o *Object) Set(name string, value interface***REMOVED******REMOVED***) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.setOwnStr(unistring.NewFromString(name), o.runtime.ToValue(value), true)
	***REMOVED***)
***REMOVED***

func (o *Object) Delete(name string) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.deleteStr(unistring.NewFromString(name), true)
	***REMOVED***)
***REMOVED***

// MarshalJSON returns JSON representation of the Object. It is equivalent to JSON.stringify(o).
// Note, this implements json.Marshaler so that json.Marshal() can be used without the need to Export().
func (o *Object) MarshalJSON() ([]byte, error) ***REMOVED***
	ctx := _builtinJSON_stringifyContext***REMOVED***
		r: o.runtime,
	***REMOVED***
	ex := o.runtime.vm.try(func() ***REMOVED***
		if !ctx.do(o) ***REMOVED***
			ctx.buf.WriteString("null")
		***REMOVED***
	***REMOVED***)
	if ex != nil ***REMOVED***
		return nil, ex
	***REMOVED***
	return ctx.buf.Bytes(), nil
***REMOVED***

// ClassName returns the class name
func (o *Object) ClassName() string ***REMOVED***
	return o.self.className()
***REMOVED***

func (o valueUnresolved) throw() ***REMOVED***
	o.r.throwReferenceError(o.ref)
***REMOVED***

func (o valueUnresolved) ToInteger() int64 ***REMOVED***
	o.throw()
	return 0
***REMOVED***

func (o valueUnresolved) toString() valueString ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) string() unistring.String ***REMOVED***
	o.throw()
	return ""
***REMOVED***

func (o valueUnresolved) ToString() Value ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) String() string ***REMOVED***
	o.throw()
	return ""
***REMOVED***

func (o valueUnresolved) ToFloat() float64 ***REMOVED***
	o.throw()
	return 0
***REMOVED***

func (o valueUnresolved) ToBoolean() bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) ToObject(*Runtime) *Object ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) ToNumber() Value ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) SameAs(Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) Equals(Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) StrictEquals(Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) baseObject(*Runtime) *Object ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) ExportType() reflect.Type ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) hash(*maphash.Hash) uint64 ***REMOVED***
	o.throw()
	return 0
***REMOVED***

func (s *valueSymbol) ToInteger() int64 ***REMOVED***
	panic(typeError("Cannot convert a Symbol value to a number"))
***REMOVED***

func (s *valueSymbol) toString() valueString ***REMOVED***
	panic(typeError("Cannot convert a Symbol value to a string"))
***REMOVED***

func (s *valueSymbol) ToString() Value ***REMOVED***
	return s
***REMOVED***

func (s *valueSymbol) String() string ***REMOVED***
	return s.desc.String()
***REMOVED***

func (s *valueSymbol) string() unistring.String ***REMOVED***
	return s.desc.string()
***REMOVED***

func (s *valueSymbol) ToFloat() float64 ***REMOVED***
	panic(typeError("Cannot convert a Symbol value to a number"))
***REMOVED***

func (s *valueSymbol) ToNumber() Value ***REMOVED***
	panic(typeError("Cannot convert a Symbol value to a number"))
***REMOVED***

func (s *valueSymbol) ToBoolean() bool ***REMOVED***
	return true
***REMOVED***

func (s *valueSymbol) ToObject(r *Runtime) *Object ***REMOVED***
	return s.baseObject(r)
***REMOVED***

func (s *valueSymbol) SameAs(other Value) bool ***REMOVED***
	if s1, ok := other.(*valueSymbol); ok ***REMOVED***
		return s == s1
	***REMOVED***
	return false
***REMOVED***

func (s *valueSymbol) Equals(o Value) bool ***REMOVED***
	return s.SameAs(o)
***REMOVED***

func (s *valueSymbol) StrictEquals(o Value) bool ***REMOVED***
	return s.SameAs(o)
***REMOVED***

func (s *valueSymbol) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return s.String()
***REMOVED***

func (s *valueSymbol) ExportType() reflect.Type ***REMOVED***
	return reflectTypeString
***REMOVED***

func (s *valueSymbol) baseObject(r *Runtime) *Object ***REMOVED***
	return r.newPrimitiveObject(s, r.global.SymbolPrototype, "Symbol")
***REMOVED***

func (s *valueSymbol) hash(*maphash.Hash) uint64 ***REMOVED***
	return uint64(s.h)
***REMOVED***

func exportValue(v Value, ctx *objectExportCtx) interface***REMOVED******REMOVED*** ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj.self.export(ctx)
	***REMOVED***
	return v.Export()
***REMOVED***

func newSymbol(s valueString) *valueSymbol ***REMOVED***
	r := &valueSymbol***REMOVED***
		desc: asciiString("Symbol(").concat(s).concat(asciiString(")")),
	***REMOVED***
	// This may need to be reconsidered in the future.
	// Depending on changes in Go's allocation policy and/or introduction of a compacting GC
	// this may no longer provide sufficient dispersion. The alternative, however, is a globally
	// synchronised random generator/hasher/sequencer and I don't want to go down that route just yet.
	r.h = uintptr(unsafe.Pointer(r))
	return r
***REMOVED***

func init() ***REMOVED***
	for i := 0; i < 256; i++ ***REMOVED***
		intCache[i] = valueInt(i - 128)
	***REMOVED***
	_positiveZero = intToValue(0)
***REMOVED***
