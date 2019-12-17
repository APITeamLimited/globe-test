package goja

import (
	"math"
	"reflect"
	"regexp"
	"strconv"
)

var (
	valueFalse    Value = valueBool(false)
	valueTrue     Value = valueBool(true)
	_null         Value = valueNull***REMOVED******REMOVED***
	_NaN          Value = valueFloat(math.NaN())
	_positiveInf  Value = valueFloat(math.Inf(+1))
	_negativeInf  Value = valueFloat(math.Inf(-1))
	_positiveZero Value
	_negativeZero Value = valueFloat(math.Float64frombits(0 | (1 << 63)))
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
	ToString() valueString
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

	assertInt() (int64, bool)
	assertString() (valueString, bool)
	assertFloat() (float64, bool)

	baseObject(r *Runtime) *Object
***REMOVED***

type valueInt int64
type valueFloat float64
type valueBool bool
type valueNull struct***REMOVED******REMOVED***
type valueUndefined struct ***REMOVED***
	valueNull
***REMOVED***

type valueUnresolved struct ***REMOVED***
	r   *Runtime
	ref string
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
	r.typeErrorResult(true, "Getter must be a function: %s", v.ToString())
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
	r.typeErrorResult(true, "Setter must be a function: %s", v.ToString())
	return nil
***REMOVED***

func (i valueInt) ToInteger() int64 ***REMOVED***
	return int64(i)
***REMOVED***

func (i valueInt) ToString() valueString ***REMOVED***
	return asciiString(i.String())
***REMOVED***

func (i valueInt) String() string ***REMOVED***
	return strconv.FormatInt(int64(i), 10)
***REMOVED***

func (i valueInt) ToFloat() float64 ***REMOVED***
	return float64(int64(i))
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
	if otherInt, ok := other.assertInt(); ok ***REMOVED***
		return int64(i) == otherInt
	***REMOVED***
	return false
***REMOVED***

func (i valueInt) Equals(other Value) bool ***REMOVED***
	if o, ok := other.assertInt(); ok ***REMOVED***
		return int64(i) == o
	***REMOVED***
	if o, ok := other.assertFloat(); ok ***REMOVED***
		return float64(i) == o
	***REMOVED***
	if o, ok := other.assertString(); ok ***REMOVED***
		return o.ToNumber().Equals(i)
	***REMOVED***
	if o, ok := other.(valueBool); ok ***REMOVED***
		return int64(i) == o.ToInteger()
	***REMOVED***
	if o, ok := other.(*Object); ok ***REMOVED***
		return i.Equals(o.self.toPrimitiveNumber())
	***REMOVED***
	return false
***REMOVED***

func (i valueInt) StrictEquals(other Value) bool ***REMOVED***
	if otherInt, ok := other.assertInt(); ok ***REMOVED***
		return int64(i) == otherInt
	***REMOVED*** else if otherFloat, ok := other.assertFloat(); ok ***REMOVED***
		return float64(i) == otherFloat
	***REMOVED***
	return false
***REMOVED***

func (i valueInt) assertInt() (int64, bool) ***REMOVED***
	return int64(i), true
***REMOVED***

func (i valueInt) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (i valueInt) assertString() (valueString, bool) ***REMOVED***
	return nil, false
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

func (o valueBool) ToInteger() int64 ***REMOVED***
	if o ***REMOVED***
		return 1
	***REMOVED***
	return 0
***REMOVED***

func (o valueBool) ToString() valueString ***REMOVED***
	if o ***REMOVED***
		return stringTrue
	***REMOVED***
	return stringFalse
***REMOVED***

func (o valueBool) String() string ***REMOVED***
	if o ***REMOVED***
		return "true"
	***REMOVED***
	return "false"
***REMOVED***

func (o valueBool) ToFloat() float64 ***REMOVED***
	if o ***REMOVED***
		return 1.0
	***REMOVED***
	return 0
***REMOVED***

func (o valueBool) ToBoolean() bool ***REMOVED***
	return bool(o)
***REMOVED***

func (o valueBool) ToObject(r *Runtime) *Object ***REMOVED***
	return r.newPrimitiveObject(o, r.global.BooleanPrototype, "Boolean")
***REMOVED***

func (o valueBool) ToNumber() Value ***REMOVED***
	if o ***REMOVED***
		return valueInt(1)
	***REMOVED***
	return valueInt(0)
***REMOVED***

func (o valueBool) SameAs(other Value) bool ***REMOVED***
	if other, ok := other.(valueBool); ok ***REMOVED***
		return o == other
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

func (o valueBool) StrictEquals(other Value) bool ***REMOVED***
	if other, ok := other.(valueBool); ok ***REMOVED***
		return o == other
	***REMOVED***
	return false
***REMOVED***

func (o valueBool) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (o valueBool) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (o valueBool) assertString() (valueString, bool) ***REMOVED***
	return nil, false
***REMOVED***

func (o valueBool) baseObject(r *Runtime) *Object ***REMOVED***
	return r.global.BooleanPrototype
***REMOVED***

func (o valueBool) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return bool(o)
***REMOVED***

func (o valueBool) ExportType() reflect.Type ***REMOVED***
	return reflectTypeBool
***REMOVED***

func (n valueNull) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (n valueNull) ToString() valueString ***REMOVED***
	return stringNull
***REMOVED***

func (n valueNull) String() string ***REMOVED***
	return "null"
***REMOVED***

func (u valueUndefined) ToString() valueString ***REMOVED***
	return stringUndefined
***REMOVED***

func (u valueUndefined) String() string ***REMOVED***
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

func (n valueNull) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (n valueNull) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (n valueNull) assertString() (valueString, bool) ***REMOVED***
	return nil, false
***REMOVED***

func (n valueNull) baseObject(r *Runtime) *Object ***REMOVED***
	return nil
***REMOVED***

func (n valueNull) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

func (n valueNull) ExportType() reflect.Type ***REMOVED***
	return reflectTypeNil
***REMOVED***

func (p *valueProperty) ToInteger() int64 ***REMOVED***
	return 0
***REMOVED***

func (p *valueProperty) ToString() valueString ***REMOVED***
	return stringEmpty
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

func (p *valueProperty) ToObject(r *Runtime) *Object ***REMOVED***
	return nil
***REMOVED***

func (p *valueProperty) ToNumber() Value ***REMOVED***
	return nil
***REMOVED***

func (p *valueProperty) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (p *valueProperty) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (p *valueProperty) assertString() (valueString, bool) ***REMOVED***
	return nil, false
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

func (p *valueProperty) Equals(other Value) bool ***REMOVED***
	return false
***REMOVED***

func (p *valueProperty) StrictEquals(other Value) bool ***REMOVED***
	return false
***REMOVED***

func (n *valueProperty) baseObject(r *Runtime) *Object ***REMOVED***
	r.typeErrorResult(true, "BUG: baseObject() is called on valueProperty") // TODO error message
	return nil
***REMOVED***

func (n *valueProperty) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	panic("Cannot export valueProperty")
***REMOVED***

func (n *valueProperty) ExportType() reflect.Type ***REMOVED***
	panic("Cannot export valueProperty")
***REMOVED***

func (f valueFloat) ToInteger() int64 ***REMOVED***
	switch ***REMOVED***
	case math.IsNaN(float64(f)):
		return 0
	case math.IsInf(float64(f), 1):
		return int64(math.MaxInt64)
	case math.IsInf(float64(f), -1):
		return int64(math.MinInt64)
	***REMOVED***
	return int64(f)
***REMOVED***

func (f valueFloat) ToString() valueString ***REMOVED***
	return asciiString(f.String())
***REMOVED***

var matchLeading0Exponent = regexp.MustCompile(`([eE][\+\-])0+([1-9])`) // 1e-07 => 1e-7

func (f valueFloat) String() string ***REMOVED***
	value := float64(f)
	if math.IsNaN(value) ***REMOVED***
		return "NaN"
	***REMOVED*** else if math.IsInf(value, 0) ***REMOVED***
		if math.Signbit(value) ***REMOVED***
			return "-Infinity"
		***REMOVED***
		return "Infinity"
	***REMOVED*** else if f == _negativeZero ***REMOVED***
		return "0"
	***REMOVED***
	exponent := math.Log10(math.Abs(value))
	if exponent >= 21 || exponent < -6 ***REMOVED***
		return matchLeading0Exponent.ReplaceAllString(strconv.FormatFloat(value, 'g', -1, 64), "$1$2")
	***REMOVED***
	return strconv.FormatFloat(value, 'f', -1, 64)
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
	if o, ok := other.assertFloat(); ok ***REMOVED***
		this := float64(f)
		if math.IsNaN(this) && math.IsNaN(o) ***REMOVED***
			return true
		***REMOVED*** else ***REMOVED***
			ret := this == o
			if ret && this == 0 ***REMOVED***
				ret = math.Signbit(this) == math.Signbit(o)
			***REMOVED***
			return ret
		***REMOVED***
	***REMOVED*** else if o, ok := other.assertInt(); ok ***REMOVED***
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
	if o, ok := other.assertFloat(); ok ***REMOVED***
		return float64(f) == o
	***REMOVED***

	if o, ok := other.assertInt(); ok ***REMOVED***
		return float64(f) == float64(o)
	***REMOVED***

	if _, ok := other.assertString(); ok ***REMOVED***
		return float64(f) == other.ToFloat()
	***REMOVED***

	if o, ok := other.(valueBool); ok ***REMOVED***
		return float64(f) == o.ToFloat()
	***REMOVED***

	if o, ok := other.(*Object); ok ***REMOVED***
		return f.Equals(o.self.toPrimitiveNumber())
	***REMOVED***

	return false
***REMOVED***

func (f valueFloat) StrictEquals(other Value) bool ***REMOVED***
	if o, ok := other.assertFloat(); ok ***REMOVED***
		return float64(f) == o
	***REMOVED*** else if o, ok := other.assertInt(); ok ***REMOVED***
		return float64(f) == float64(o)
	***REMOVED***
	return false
***REMOVED***

func (f valueFloat) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (f valueFloat) assertFloat() (float64, bool) ***REMOVED***
	return float64(f), true
***REMOVED***

func (f valueFloat) assertString() (valueString, bool) ***REMOVED***
	return nil, false
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

func (o *Object) ToInteger() int64 ***REMOVED***
	return o.self.toPrimitiveNumber().ToNumber().ToInteger()
***REMOVED***

func (o *Object) ToString() valueString ***REMOVED***
	return o.self.toPrimitiveString().ToString()
***REMOVED***

func (o *Object) String() string ***REMOVED***
	return o.self.toPrimitiveString().String()
***REMOVED***

func (o *Object) ToFloat() float64 ***REMOVED***
	return o.self.toPrimitiveNumber().ToFloat()
***REMOVED***

func (o *Object) ToBoolean() bool ***REMOVED***
	return true
***REMOVED***

func (o *Object) ToObject(r *Runtime) *Object ***REMOVED***
	return o
***REMOVED***

func (o *Object) ToNumber() Value ***REMOVED***
	return o.self.toPrimitiveNumber().ToNumber()
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

	if _, ok := other.assertInt(); ok ***REMOVED***
		return o.self.toPrimitive().Equals(other)
	***REMOVED***

	if _, ok := other.assertFloat(); ok ***REMOVED***
		return o.self.toPrimitive().Equals(other)
	***REMOVED***

	if other, ok := other.(valueBool); ok ***REMOVED***
		return o.Equals(other.ToNumber())
	***REMOVED***

	if _, ok := other.assertString(); ok ***REMOVED***
		return o.self.toPrimitive().Equals(other)
	***REMOVED***
	return false
***REMOVED***

func (o *Object) StrictEquals(other Value) bool ***REMOVED***
	if other, ok := other.(*Object); ok ***REMOVED***
		return o == other || o.self.equal(other.self)
	***REMOVED***
	return false
***REMOVED***

func (o *Object) assertInt() (int64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (o *Object) assertFloat() (float64, bool) ***REMOVED***
	return 0, false
***REMOVED***

func (o *Object) assertString() (valueString, bool) ***REMOVED***
	return nil, false
***REMOVED***

func (o *Object) baseObject(r *Runtime) *Object ***REMOVED***
	return o
***REMOVED***

func (o *Object) Export() interface***REMOVED******REMOVED*** ***REMOVED***
	return o.self.export()
***REMOVED***

func (o *Object) ExportType() reflect.Type ***REMOVED***
	return o.self.exportType()
***REMOVED***

func (o *Object) Get(name string) Value ***REMOVED***
	return o.self.getStr(name)
***REMOVED***

func (o *Object) Keys() (keys []string) ***REMOVED***
	for item, f := o.self.enumerate(false, false)(); f != nil; item, f = f() ***REMOVED***
		keys = append(keys, item.name)
	***REMOVED***

	return
***REMOVED***

// DefineDataProperty is a Go equivalent of Object.defineProperty(o, name, ***REMOVED***value: value, writable: writable,
// configurable: configurable, enumerable: enumerable***REMOVED***)
func (o *Object) DefineDataProperty(name string, value Value, writable, configurable, enumerable Flag) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.defineOwnProperty(newStringValue(name), propertyDescr***REMOVED***
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
		o.self.defineOwnProperty(newStringValue(name), propertyDescr***REMOVED***
			Getter:       getter,
			Setter:       setter,
			Configurable: configurable,
			Enumerable:   enumerable,
		***REMOVED***, true)
	***REMOVED***)
***REMOVED***

func (o *Object) Set(name string, value interface***REMOVED******REMOVED***) error ***REMOVED***
	return tryFunc(func() ***REMOVED***
		o.self.putStr(name, o.runtime.ToValue(value), true)
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

func (o valueUnresolved) ToString() valueString ***REMOVED***
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

func (o valueUnresolved) ToObject(r *Runtime) *Object ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) ToNumber() Value ***REMOVED***
	o.throw()
	return nil
***REMOVED***

func (o valueUnresolved) SameAs(other Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) Equals(other Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) StrictEquals(other Value) bool ***REMOVED***
	o.throw()
	return false
***REMOVED***

func (o valueUnresolved) assertInt() (int64, bool) ***REMOVED***
	o.throw()
	return 0, false
***REMOVED***

func (o valueUnresolved) assertFloat() (float64, bool) ***REMOVED***
	o.throw()
	return 0, false
***REMOVED***

func (o valueUnresolved) assertString() (valueString, bool) ***REMOVED***
	o.throw()
	return nil, false
***REMOVED***

func (o valueUnresolved) baseObject(r *Runtime) *Object ***REMOVED***
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

func init() ***REMOVED***
	for i := 0; i < 256; i++ ***REMOVED***
		intCache[i] = valueInt(i - 128)
	***REMOVED***
	_positiveZero = intToValue(0)
***REMOVED***
