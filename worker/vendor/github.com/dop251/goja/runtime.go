package goja

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"hash/maphash"
	"math"
	"math/bits"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/text/collate"

	js_ast "github.com/dop251/goja/ast"
	"github.com/dop251/goja/file"
	"github.com/dop251/goja/parser"
	"github.com/dop251/goja/unistring"
)

const (
	sqrt1_2 float64 = math.Sqrt2 / 2

	deoptimiseRegexp = false
)

var (
	typeCallable = reflect.TypeOf(Callable(nil))
	typeValue    = reflect.TypeOf((*Value)(nil)).Elem()
	typeObject   = reflect.TypeOf((*Object)(nil))
	typeTime     = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
	typeBytes    = reflect.TypeOf(([]byte)(nil))
)

type iterationKind int

const (
	iterationKindKey iterationKind = iota
	iterationKindValue
	iterationKindKeyValue
)

type global struct ***REMOVED***
	stash    stash
	varNames map[unistring.String]struct***REMOVED******REMOVED***

	Object   *Object
	Array    *Object
	Function *Object
	String   *Object
	Number   *Object
	Boolean  *Object
	RegExp   *Object
	Date     *Object
	Symbol   *Object
	Proxy    *Object
	Promise  *Object

	ArrayBuffer       *Object
	DataView          *Object
	TypedArray        *Object
	Uint8Array        *Object
	Uint8ClampedArray *Object
	Int8Array         *Object
	Uint16Array       *Object
	Int16Array        *Object
	Uint32Array       *Object
	Int32Array        *Object
	Float32Array      *Object
	Float64Array      *Object

	WeakSet *Object
	WeakMap *Object
	Map     *Object
	Set     *Object

	Error          *Object
	AggregateError *Object
	TypeError      *Object
	ReferenceError *Object
	SyntaxError    *Object
	RangeError     *Object
	EvalError      *Object
	URIError       *Object

	GoError *Object

	ObjectPrototype   *Object
	ArrayPrototype    *Object
	NumberPrototype   *Object
	StringPrototype   *Object
	BooleanPrototype  *Object
	FunctionPrototype *Object
	RegExpPrototype   *Object
	DatePrototype     *Object
	SymbolPrototype   *Object

	ArrayBufferPrototype *Object
	DataViewPrototype    *Object
	TypedArrayPrototype  *Object
	WeakSetPrototype     *Object
	WeakMapPrototype     *Object
	MapPrototype         *Object
	SetPrototype         *Object
	PromisePrototype     *Object

	IteratorPrototype             *Object
	ArrayIteratorPrototype        *Object
	MapIteratorPrototype          *Object
	SetIteratorPrototype          *Object
	StringIteratorPrototype       *Object
	RegExpStringIteratorPrototype *Object

	ErrorPrototype          *Object
	AggregateErrorPrototype *Object
	TypeErrorPrototype      *Object
	SyntaxErrorPrototype    *Object
	RangeErrorPrototype     *Object
	ReferenceErrorPrototype *Object
	EvalErrorPrototype      *Object
	URIErrorPrototype       *Object

	GoErrorPrototype *Object

	Eval *Object

	thrower         *Object
	throwerProperty Value

	stdRegexpProto *guardedObject

	weakSetAdder  *Object
	weakMapAdder  *Object
	mapAdder      *Object
	setAdder      *Object
	arrayValues   *Object
	arrayToString *Object
***REMOVED***

type Flag int

const (
	FLAG_NOT_SET Flag = iota
	FLAG_FALSE
	FLAG_TRUE
)

func (f Flag) Bool() bool ***REMOVED***
	return f == FLAG_TRUE
***REMOVED***

func ToFlag(b bool) Flag ***REMOVED***
	if b ***REMOVED***
		return FLAG_TRUE
	***REMOVED***
	return FLAG_FALSE
***REMOVED***

type RandSource func() float64

type Now func() time.Time

type Runtime struct ***REMOVED***
	global          global
	globalObject    *Object
	stringSingleton *stringObject
	rand            RandSource
	now             Now
	_collator       *collate.Collator
	parserOptions   []parser.Option

	symbolRegistry map[unistring.String]*Symbol

	typeInfoCache   map[reflect.Type]*reflectTypeInfo
	fieldNameMapper FieldNameMapper

	vm    *vm
	hash  *maphash.Hash
	idSeq uint64

	jobQueue []func()

	promiseRejectionTracker PromiseRejectionTracker
***REMOVED***

type StackFrame struct ***REMOVED***
	prg      *Program
	funcName unistring.String
	pc       int
***REMOVED***

func (f *StackFrame) SrcName() string ***REMOVED***
	if f.prg == nil ***REMOVED***
		return "<native>"
	***REMOVED***
	return f.prg.src.Name()
***REMOVED***

func (f *StackFrame) FuncName() string ***REMOVED***
	if f.funcName == "" && f.prg == nil ***REMOVED***
		return "<native>"
	***REMOVED***
	if f.funcName == "" ***REMOVED***
		return "<anonymous>"
	***REMOVED***
	return f.funcName.String()
***REMOVED***

func (f *StackFrame) Position() file.Position ***REMOVED***
	if f.prg == nil || f.prg.src == nil ***REMOVED***
		return file.Position***REMOVED******REMOVED***
	***REMOVED***
	return f.prg.src.Position(f.prg.sourceOffset(f.pc))
***REMOVED***

func (f *StackFrame) WriteToValueBuilder(b *valueStringBuilder) ***REMOVED***
	if f.prg != nil ***REMOVED***
		if n := f.prg.funcName; n != "" ***REMOVED***
			b.WriteString(stringValueFromRaw(n))
			b.WriteASCII(" (")
		***REMOVED***
		p := f.Position()
		if p.Filename != "" ***REMOVED***
			b.WriteASCII(p.Filename)
		***REMOVED*** else ***REMOVED***
			b.WriteASCII("<eval>")
		***REMOVED***
		b.WriteRune(':')
		b.WriteASCII(strconv.Itoa(p.Line))
		b.WriteRune(':')
		b.WriteASCII(strconv.Itoa(p.Column))
		b.WriteRune('(')
		b.WriteASCII(strconv.Itoa(f.pc))
		b.WriteRune(')')
		if f.prg.funcName != "" ***REMOVED***
			b.WriteRune(')')
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if f.funcName != "" ***REMOVED***
			b.WriteString(stringValueFromRaw(f.funcName))
			b.WriteASCII(" (")
		***REMOVED***
		b.WriteASCII("native")
		if f.funcName != "" ***REMOVED***
			b.WriteRune(')')
		***REMOVED***
	***REMOVED***
***REMOVED***

func (f *StackFrame) Write(b *bytes.Buffer) ***REMOVED***
	if f.prg != nil ***REMOVED***
		if n := f.prg.funcName; n != "" ***REMOVED***
			b.WriteString(n.String())
			b.WriteString(" (")
		***REMOVED***
		p := f.Position()
		if p.Filename != "" ***REMOVED***
			b.WriteString(p.Filename)
		***REMOVED*** else ***REMOVED***
			b.WriteString("<eval>")
		***REMOVED***
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(p.Line))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(p.Column))
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(f.pc))
		b.WriteByte(')')
		if f.prg.funcName != "" ***REMOVED***
			b.WriteByte(')')
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if f.funcName != "" ***REMOVED***
			b.WriteString(f.funcName.String())
			b.WriteString(" (")
		***REMOVED***
		b.WriteString("native")
		if f.funcName != "" ***REMOVED***
			b.WriteByte(')')
		***REMOVED***
	***REMOVED***
***REMOVED***

type Exception struct ***REMOVED***
	val   Value
	stack []StackFrame
***REMOVED***

type uncatchableException struct ***REMOVED***
	err error
***REMOVED***

func (ue *uncatchableException) Unwrap() error ***REMOVED***
	return ue.err
***REMOVED***

type InterruptedError struct ***REMOVED***
	Exception
	iface interface***REMOVED******REMOVED***
***REMOVED***

func (e *InterruptedError) Unwrap() error ***REMOVED***
	if err, ok := e.iface.(error); ok ***REMOVED***
		return err
	***REMOVED***
	return nil
***REMOVED***

type StackOverflowError struct ***REMOVED***
	Exception
***REMOVED***

func (e *InterruptedError) Value() interface***REMOVED******REMOVED*** ***REMOVED***
	return e.iface
***REMOVED***

func (e *InterruptedError) String() string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	var b bytes.Buffer
	if e.iface != nil ***REMOVED***
		b.WriteString(fmt.Sprint(e.iface))
		b.WriteByte('\n')
	***REMOVED***
	e.writeFullStack(&b)
	return b.String()
***REMOVED***

func (e *InterruptedError) Error() string ***REMOVED***
	if e == nil || e.iface == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	var b bytes.Buffer
	b.WriteString(fmt.Sprint(e.iface))
	e.writeShortStack(&b)
	return b.String()
***REMOVED***

func (e *Exception) writeFullStack(b *bytes.Buffer) ***REMOVED***
	for _, frame := range e.stack ***REMOVED***
		b.WriteString("\tat ")
		frame.Write(b)
		b.WriteByte('\n')
	***REMOVED***
***REMOVED***

func (e *Exception) writeShortStack(b *bytes.Buffer) ***REMOVED***
	if len(e.stack) > 0 && (e.stack[0].prg != nil || e.stack[0].funcName != "") ***REMOVED***
		b.WriteString(" at ")
		e.stack[0].Write(b)
	***REMOVED***
***REMOVED***

func (e *Exception) String() string ***REMOVED***
	if e == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	var b bytes.Buffer
	if e.val != nil ***REMOVED***
		b.WriteString(e.val.String())
		b.WriteByte('\n')
	***REMOVED***
	e.writeFullStack(&b)
	return b.String()
***REMOVED***

func (e *Exception) Error() string ***REMOVED***
	if e == nil || e.val == nil ***REMOVED***
		return "<nil>"
	***REMOVED***
	var b bytes.Buffer
	b.WriteString(e.val.String())
	e.writeShortStack(&b)
	return b.String()
***REMOVED***

func (e *Exception) Value() Value ***REMOVED***
	return e.val
***REMOVED***

func (r *Runtime) addToGlobal(name string, value Value) ***REMOVED***
	r.globalObject.self._putProp(unistring.String(name), value, true, false, true)
***REMOVED***

func (r *Runtime) createIterProto(val *Object) objectImpl ***REMOVED***
	o := newBaseObjectObj(val, r.global.ObjectPrototype, classObject)

	o._putSym(SymIterator, valueProp(r.newNativeFunc(r.returnThis, nil, "[Symbol.iterator]", nil, 0), true, false, true))
	return o
***REMOVED***

func (r *Runtime) init() ***REMOVED***
	r.rand = rand.Float64
	r.now = time.Now
	r.global.ObjectPrototype = r.newBaseObject(nil, classObject).val
	r.globalObject = r.NewObject()

	r.vm = &vm***REMOVED***
		r: r,
	***REMOVED***
	r.vm.init()

	funcProto := r.newNativeFunc(func(FunctionCall) Value ***REMOVED***
		return _undefined
	***REMOVED***, nil, " ", nil, 0)
	r.global.FunctionPrototype = funcProto
	funcProtoObj := funcProto.self.(*nativeFuncObject)

	r.global.IteratorPrototype = r.newLazyObject(r.createIterProto)

	r.initObject()
	r.initFunction()
	r.initArray()
	r.initString()
	r.initGlobalObject()
	r.initNumber()
	r.initRegExp()
	r.initDate()
	r.initBoolean()
	r.initProxy()
	r.initReflect()

	r.initErrors()

	r.global.Eval = r.newNativeFunc(r.builtin_eval, nil, "eval", nil, 1)
	r.addToGlobal("eval", r.global.Eval)

	r.initMath()
	r.initJSON()

	r.initTypedArrays()
	r.initSymbol()
	r.initWeakSet()
	r.initWeakMap()
	r.initMap()
	r.initSet()
	r.initPromise()

	r.global.thrower = r.newNativeFunc(r.builtin_thrower, nil, "", nil, 0)
	r.global.throwerProperty = &valueProperty***REMOVED***
		getterFunc: r.global.thrower,
		setterFunc: r.global.thrower,
		accessor:   true,
	***REMOVED***
	r.object_freeze(FunctionCall***REMOVED***Arguments: []Value***REMOVED***r.global.thrower***REMOVED******REMOVED***)

	funcProtoObj._put("caller", &valueProperty***REMOVED***
		getterFunc:   r.global.thrower,
		setterFunc:   r.global.thrower,
		accessor:     true,
		configurable: true,
	***REMOVED***)
	funcProtoObj._put("arguments", &valueProperty***REMOVED***
		getterFunc:   r.global.thrower,
		setterFunc:   r.global.thrower,
		accessor:     true,
		configurable: true,
	***REMOVED***)
***REMOVED***

func (r *Runtime) typeErrorResult(throw bool, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if throw ***REMOVED***
		panic(r.NewTypeError(args...))
	***REMOVED***
***REMOVED***

func (r *Runtime) newError(typ *Object, format string, args ...interface***REMOVED******REMOVED***) Value ***REMOVED***
	var msg string
	if len(args) > 0 ***REMOVED***
		msg = fmt.Sprintf(format, args...)
	***REMOVED*** else ***REMOVED***
		msg = format
	***REMOVED***
	return r.builtin_new(typ, []Value***REMOVED***newStringValue(msg)***REMOVED***)
***REMOVED***

func (r *Runtime) throwReferenceError(name unistring.String) ***REMOVED***
	panic(r.newError(r.global.ReferenceError, "%s is not defined", name))
***REMOVED***

func (r *Runtime) newSyntaxError(msg string, offset int) Value ***REMOVED***
	return r.builtin_new(r.global.SyntaxError, []Value***REMOVED***newStringValue(msg)***REMOVED***)
***REMOVED***

func newBaseObjectObj(obj, proto *Object, class string) *baseObject ***REMOVED***
	o := &baseObject***REMOVED***
		class:      class,
		val:        obj,
		extensible: true,
		prototype:  proto,
	***REMOVED***
	obj.self = o
	o.init()
	return o
***REMOVED***

func newGuardedObj(proto *Object, class string) *guardedObject ***REMOVED***
	return &guardedObject***REMOVED***
		baseObject: baseObject***REMOVED***
			class:      class,
			extensible: true,
			prototype:  proto,
		***REMOVED***,
	***REMOVED***
***REMOVED***

func (r *Runtime) newBaseObject(proto *Object, class string) (o *baseObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	return newBaseObjectObj(v, proto, class)
***REMOVED***

func (r *Runtime) newGuardedObject(proto *Object, class string) (o *guardedObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***
	o = newGuardedObj(proto, class)
	v.self = o
	o.val = v
	o.init()
	return
***REMOVED***

func (r *Runtime) NewObject() (v *Object) ***REMOVED***
	return r.newBaseObject(r.global.ObjectPrototype, classObject).val
***REMOVED***

// CreateObject creates an object with given prototype. Equivalent of Object.create(proto).
func (r *Runtime) CreateObject(proto *Object) *Object ***REMOVED***
	return r.newBaseObject(proto, classObject).val
***REMOVED***

func (r *Runtime) NewArray(items ...interface***REMOVED******REMOVED***) *Object ***REMOVED***
	values := make([]Value, len(items))
	for i, item := range items ***REMOVED***
		values[i] = r.ToValue(item)
	***REMOVED***
	return r.newArrayValues(values)
***REMOVED***

func (r *Runtime) NewTypeError(args ...interface***REMOVED******REMOVED***) *Object ***REMOVED***
	msg := ""
	if len(args) > 0 ***REMOVED***
		f, _ := args[0].(string)
		msg = fmt.Sprintf(f, args[1:]...)
	***REMOVED***
	return r.builtin_new(r.global.TypeError, []Value***REMOVED***newStringValue(msg)***REMOVED***)
***REMOVED***

func (r *Runtime) NewGoError(err error) *Object ***REMOVED***
	e := r.newError(r.global.GoError, err.Error()).(*Object)
	e.Set("value", err)
	return e
***REMOVED***

func (r *Runtime) newFunc(name unistring.String, length int, strict bool) (f *funcObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f = &funcObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	f.strict = strict
	v.self = f
	f.prototype = r.global.FunctionPrototype
	f.init(name, intToValue(int64(length)))
	return
***REMOVED***

func (r *Runtime) newClassFunc(name unistring.String, length int, proto *Object, derived bool) (f *classFuncObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f = &classFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	f.strict = true
	f.derived = derived
	v.self = f
	f.prototype = proto
	f.init(name, intToValue(int64(length)))
	return
***REMOVED***

func (r *Runtime) newMethod(name unistring.String, length int, strict bool) (f *methodFuncObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f = &methodFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	f.strict = strict
	v.self = f
	f.prototype = r.global.FunctionPrototype
	f.init(name, intToValue(int64(length)))
	return
***REMOVED***

func (r *Runtime) newArrowFunc(name unistring.String, length int, strict bool) (f *arrowFuncObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f = &arrowFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	f.strict = strict

	vm := r.vm

	f.newTarget = vm.newTarget
	v.self = f
	f.prototype = r.global.FunctionPrototype
	f.init(name, intToValue(int64(length)))
	return
***REMOVED***

func (r *Runtime) newNativeFuncObj(v *Object, call func(FunctionCall) Value, construct func(args []Value, proto *Object) *Object, name unistring.String, proto *Object, length Value) *nativeFuncObject ***REMOVED***
	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
		f:         call,
		construct: r.wrapNativeConstruct(construct, proto),
	***REMOVED***
	v.self = f
	f.init(name, length)
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
	***REMOVED***
	return f
***REMOVED***

func (r *Runtime) newNativeConstructor(call func(ConstructorCall) *Object, name unistring.String, length int64) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***

	f.f = func(c FunctionCall) Value ***REMOVED***
		thisObj, _ := c.This.(*Object)
		if thisObj != nil ***REMOVED***
			res := call(ConstructorCall***REMOVED***
				This:      thisObj,
				Arguments: c.Arguments,
			***REMOVED***)
			if res == nil ***REMOVED***
				return _undefined
			***REMOVED***
			return res
		***REMOVED***
		return f.defaultConstruct(call, c.Arguments, nil)
	***REMOVED***

	f.construct = func(args []Value, newTarget *Object) *Object ***REMOVED***
		return f.defaultConstruct(call, args, newTarget)
	***REMOVED***

	v.self = f
	f.init(name, intToValue(length))

	proto := r.NewObject()
	proto.self._putProp("constructor", v, true, false, true)
	f._putProp("prototype", proto, true, false, false)

	return v
***REMOVED***

func (r *Runtime) newNativeConstructOnly(v *Object, ctor func(args []Value, newTarget *Object) *Object, defaultProto *Object, name unistring.String, length int64) *nativeFuncObject ***REMOVED***
	return r.newNativeFuncAndConstruct(v, func(call FunctionCall) Value ***REMOVED***
		return ctor(call.Arguments, nil)
	***REMOVED***,
		func(args []Value, newTarget *Object) *Object ***REMOVED***
			if newTarget == nil ***REMOVED***
				newTarget = v
			***REMOVED***
			return ctor(args, newTarget)
		***REMOVED***, defaultProto, name, intToValue(length))
***REMOVED***

func (r *Runtime) newNativeFuncAndConstruct(v *Object, call func(call FunctionCall) Value, ctor func(args []Value, newTarget *Object) *Object, defaultProto *Object, name unistring.String, l Value) *nativeFuncObject ***REMOVED***
	if v == nil ***REMOVED***
		v = &Object***REMOVED***runtime: r***REMOVED***
	***REMOVED***

	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
		f:         call,
		construct: ctor,
	***REMOVED***
	v.self = f
	f.init(name, l)
	if defaultProto != nil ***REMOVED***
		f._putProp("prototype", defaultProto, false, false, false)
	***REMOVED***

	return f
***REMOVED***

func (r *Runtime) newNativeFunc(call func(FunctionCall) Value, construct func(args []Value, proto *Object) *Object, name unistring.String, proto *Object, length int) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
		f:         call,
		construct: r.wrapNativeConstruct(construct, proto),
	***REMOVED***
	v.self = f
	f.init(name, intToValue(int64(length)))
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
		proto.self._putProp("constructor", v, true, false, true)
	***REMOVED***
	return v
***REMOVED***

func (r *Runtime) newNativeFuncConstructObj(v *Object, construct func(args []Value, proto *Object) *Object, name unistring.String, proto *Object, length int) *nativeFuncObject ***REMOVED***
	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
		f:         r.constructToCall(construct, proto),
		construct: r.wrapNativeConstruct(construct, proto),
	***REMOVED***

	f.init(name, intToValue(int64(length)))
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
	***REMOVED***
	return f
***REMOVED***

func (r *Runtime) newNativeFuncConstruct(construct func(args []Value, proto *Object) *Object, name unistring.String, prototype *Object, length int64) *Object ***REMOVED***
	return r.newNativeFuncConstructProto(construct, name, prototype, r.global.FunctionPrototype, length)
***REMOVED***

func (r *Runtime) newNativeFuncConstructProto(construct func(args []Value, proto *Object) *Object, name unistring.String, prototype, proto *Object, length int64) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f := &nativeFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	v.self = f
	f.prototype = proto
	f.f = r.constructToCall(construct, prototype)
	f.construct = r.wrapNativeConstruct(construct, prototype)
	f.init(name, intToValue(length))
	if prototype != nil ***REMOVED***
		f._putProp("prototype", prototype, false, false, false)
		prototype.self._putProp("constructor", v, true, false, true)
	***REMOVED***
	return v
***REMOVED***

func (r *Runtime) newPrimitiveObject(value Value, proto *Object, class string) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o := &primitiveValueObject***REMOVED******REMOVED***
	o.class = class
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = proto
	o.pValue = value
	o.init()
	return v
***REMOVED***

func (r *Runtime) builtin_Number(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) > 0 ***REMOVED***
		return call.Arguments[0].ToNumber()
	***REMOVED*** else ***REMOVED***
		return valueInt(0)
	***REMOVED***
***REMOVED***

func (r *Runtime) builtin_newNumber(args []Value, proto *Object) *Object ***REMOVED***
	var v Value
	if len(args) > 0 ***REMOVED***
		v = args[0].ToNumber()
	***REMOVED*** else ***REMOVED***
		v = intToValue(0)
	***REMOVED***
	return r.newPrimitiveObject(v, proto, classNumber)
***REMOVED***

func (r *Runtime) builtin_Boolean(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) > 0 ***REMOVED***
		if call.Arguments[0].ToBoolean() ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		return valueFalse
	***REMOVED***
***REMOVED***

func (r *Runtime) builtin_newBoolean(args []Value, proto *Object) *Object ***REMOVED***
	var v Value
	if len(args) > 0 ***REMOVED***
		if args[0].ToBoolean() ***REMOVED***
			v = valueTrue
		***REMOVED*** else ***REMOVED***
			v = valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		v = valueFalse
	***REMOVED***
	return r.newPrimitiveObject(v, proto, classBoolean)
***REMOVED***

func (r *Runtime) error_toString(call FunctionCall) Value ***REMOVED***
	var nameStr, msgStr valueString
	obj := r.toObject(call.This)
	name := obj.self.getStr("name", nil)
	if name == nil || name == _undefined ***REMOVED***
		nameStr = asciiString("Error")
	***REMOVED*** else ***REMOVED***
		nameStr = name.toString()
	***REMOVED***
	msg := obj.self.getStr("message", nil)
	if msg == nil || msg == _undefined ***REMOVED***
		msgStr = stringEmpty
	***REMOVED*** else ***REMOVED***
		msgStr = msg.toString()
	***REMOVED***
	if nameStr.length() == 0 ***REMOVED***
		return msgStr
	***REMOVED***
	if msgStr.length() == 0 ***REMOVED***
		return nameStr
	***REMOVED***
	var sb valueStringBuilder
	sb.WriteString(nameStr)
	sb.WriteString(asciiString(": "))
	sb.WriteString(msgStr)
	return sb.String()
***REMOVED***

func (r *Runtime) builtin_new(construct *Object, args []Value) *Object ***REMOVED***
	return r.toConstructor(construct)(args, nil)
***REMOVED***

func (r *Runtime) builtin_thrower(call FunctionCall) Value ***REMOVED***
	obj := r.toObject(call.This)
	strict := true
	switch fn := obj.self.(type) ***REMOVED***
	case *funcObject:
		strict = fn.strict
	***REMOVED***
	r.typeErrorResult(strict, "'caller', 'callee', and 'arguments' properties may not be accessed on strict mode functions or the arguments objects for calls to them")
	return nil
***REMOVED***

func (r *Runtime) eval(srcVal valueString, direct, strict bool) Value ***REMOVED***
	src := escapeInvalidUtf16(srcVal)
	vm := r.vm
	inGlobal := true
	if direct ***REMOVED***
		for s := vm.stash; s != nil; s = s.outer ***REMOVED***
			if s.isVariable() ***REMOVED***
				inGlobal = false
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***
	vm.pushCtx()
	funcObj := _undefined
	if !direct ***REMOVED***
		vm.stash = &r.global.stash
		vm.privEnv = nil
	***REMOVED*** else ***REMOVED***
		if sb := vm.sb; sb > 0 ***REMOVED***
			funcObj = vm.stack[sb-1]
		***REMOVED***
	***REMOVED***
	p, err := r.compile("<eval>", src, strict, inGlobal, r.vm)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	vm.prg = p
	vm.pc = 0
	vm.args = 0
	vm.result = _undefined
	vm.push(funcObj)
	vm.sb = vm.sp
	vm.push(nil) // this
	vm.run()
	retval := vm.result
	vm.popCtx()
	vm.halt = false
	vm.sp -= 2
	return retval
***REMOVED***

func (r *Runtime) builtin_eval(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) == 0 ***REMOVED***
		return _undefined
	***REMOVED***
	if str, ok := call.Arguments[0].(valueString); ok ***REMOVED***
		return r.eval(str, false, false)
	***REMOVED***
	return call.Arguments[0]
***REMOVED***

func (r *Runtime) constructToCall(construct func(args []Value, proto *Object) *Object, proto *Object) func(call FunctionCall) Value ***REMOVED***
	return func(call FunctionCall) Value ***REMOVED***
		return construct(call.Arguments, proto)
	***REMOVED***
***REMOVED***

func (r *Runtime) wrapNativeConstruct(c func(args []Value, proto *Object) *Object, proto *Object) func(args []Value, newTarget *Object) *Object ***REMOVED***
	if c == nil ***REMOVED***
		return nil
	***REMOVED***
	return func(args []Value, newTarget *Object) *Object ***REMOVED***
		var p *Object
		if newTarget != nil ***REMOVED***
			if pp, ok := newTarget.self.getStr("prototype", nil).(*Object); ok ***REMOVED***
				p = pp
			***REMOVED***
		***REMOVED***
		if p == nil ***REMOVED***
			p = proto
		***REMOVED***
		return c(args, p)
	***REMOVED***
***REMOVED***

func (r *Runtime) toCallable(v Value) func(FunctionCall) Value ***REMOVED***
	if call, ok := r.toObject(v).self.assertCallable(); ok ***REMOVED***
		return call
	***REMOVED***
	r.typeErrorResult(true, "Value is not callable: %s", v.toString())
	return nil
***REMOVED***

func (r *Runtime) checkObjectCoercible(v Value) ***REMOVED***
	switch v.(type) ***REMOVED***
	case valueUndefined, valueNull:
		r.typeErrorResult(true, "Value is not object coercible")
	***REMOVED***
***REMOVED***

func toInt8(v Value) int8 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return int8(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int8(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint8(v Value) uint8 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return uint8(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint8(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint8Clamp(v Value) uint8 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		if i < 0 ***REMOVED***
			return 0
		***REMOVED***
		if i <= 255 ***REMOVED***
			return uint8(i)
		***REMOVED***
		return 255
	***REMOVED***

	if num, ok := v.(valueFloat); ok ***REMOVED***
		num := float64(num)
		if !math.IsNaN(num) ***REMOVED***
			if num < 0 ***REMOVED***
				return 0
			***REMOVED***
			if num > 255 ***REMOVED***
				return 255
			***REMOVED***
			f := math.Floor(num)
			f1 := f + 0.5
			if f1 < num ***REMOVED***
				return uint8(f + 1)
			***REMOVED***
			if f1 > num ***REMOVED***
				return uint8(f)
			***REMOVED***
			r := uint8(f)
			if r&1 != 0 ***REMOVED***
				return r + 1
			***REMOVED***
			return r
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toInt16(v Value) int16 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return int16(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int16(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint16(v Value) uint16 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return uint16(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint16(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toInt32(v Value) int32 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return int32(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int32(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint32(v Value) uint32 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return uint32(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint32(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toInt64(v Value) int64 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return int64(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int64(f)
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint64(v Value) uint64 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return uint64(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint64(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toInt(v Value) int ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return int(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int(f)
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUint(v Value) uint ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.(valueInt); ok ***REMOVED***
		return uint(i)
	***REMOVED***

	if f, ok := v.(valueFloat); ok ***REMOVED***
		f := float64(f)
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toFloat32(v Value) float32 ***REMOVED***
	return float32(v.ToFloat())
***REMOVED***

func toLength(v Value) int64 ***REMOVED***
	if v == nil ***REMOVED***
		return 0
	***REMOVED***
	i := v.ToInteger()
	if i < 0 ***REMOVED***
		return 0
	***REMOVED***
	if i >= maxInt ***REMOVED***
		return maxInt - 1
	***REMOVED***
	return i
***REMOVED***

func (r *Runtime) toLengthUint32(v Value) uint32 ***REMOVED***
	var intVal int64
repeat:
	switch num := v.(type) ***REMOVED***
	case valueInt:
		intVal = int64(num)
	case valueFloat:
		if v != _negativeZero ***REMOVED***
			if i, ok := floatToInt(float64(num)); ok ***REMOVED***
				intVal = i
			***REMOVED*** else ***REMOVED***
				goto fail
			***REMOVED***
		***REMOVED***
	case valueString:
		v = num.ToNumber()
		goto repeat
	default:
		// Legacy behaviour as specified in https://tc39.es/ecma262/#sec-arraysetlength (see the note)
		n2 := toUint32(v)
		n1 := v.ToNumber()
		if f, ok := n1.(valueFloat); ok ***REMOVED***
			f := float64(f)
			if f != 0 || !math.Signbit(f) ***REMOVED***
				goto fail
			***REMOVED***
		***REMOVED***
		if n1.ToInteger() != int64(n2) ***REMOVED***
			goto fail
		***REMOVED***
		return n2
	***REMOVED***
	if intVal >= 0 && intVal <= math.MaxUint32 ***REMOVED***
		return uint32(intVal)
	***REMOVED***
fail:
	panic(r.newError(r.global.RangeError, "Invalid array length"))
***REMOVED***

func toIntStrict(i int64) int ***REMOVED***
	if bits.UintSize == 32 ***REMOVED***
		if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
			panic(rangeError("Integer value overflows 32-bit int"))
		***REMOVED***
	***REMOVED***
	return int(i)
***REMOVED***

func toIntClamp(i int64) int ***REMOVED***
	if bits.UintSize == 32 ***REMOVED***
		if i > math.MaxInt32 ***REMOVED***
			return math.MaxInt32
		***REMOVED***
		if i < math.MinInt32 ***REMOVED***
			return math.MinInt32
		***REMOVED***
	***REMOVED***
	return int(i)
***REMOVED***

func (r *Runtime) toIndex(v Value) int ***REMOVED***
	num := v.ToInteger()
	if num >= 0 && num < maxInt ***REMOVED***
		if bits.UintSize == 32 && num >= math.MaxInt32 ***REMOVED***
			panic(r.newError(r.global.RangeError, "Index %s overflows int", v.String()))
		***REMOVED***
		return int(num)
	***REMOVED***
	panic(r.newError(r.global.RangeError, "Invalid index %s", v.String()))
***REMOVED***

func (r *Runtime) toBoolean(b bool) Value ***REMOVED***
	if b ***REMOVED***
		return valueTrue
	***REMOVED*** else ***REMOVED***
		return valueFalse
	***REMOVED***
***REMOVED***

// New creates an instance of a Javascript runtime that can be used to run code. Multiple instances may be created and
// used simultaneously, however it is not possible to pass JS values across runtimes.
func New() *Runtime ***REMOVED***
	r := &Runtime***REMOVED******REMOVED***
	r.init()
	return r
***REMOVED***

// Compile creates an internal representation of the JavaScript code that can be later run using the Runtime.RunProgram()
// method. This representation is not linked to a runtime in any way and can be run in multiple runtimes (possibly
// at the same time).
func Compile(name, src string, strict bool) (*Program, error) ***REMOVED***
	return compile(name, src, strict, true, nil)
***REMOVED***

// CompileAST creates an internal representation of the JavaScript code that can be later run using the Runtime.RunProgram()
// method. This representation is not linked to a runtime in any way and can be run in multiple runtimes (possibly
// at the same time).
func CompileAST(prg *js_ast.Program, strict bool) (*Program, error) ***REMOVED***
	return compileAST(prg, strict, true, nil)
***REMOVED***

// MustCompile is like Compile but panics if the code cannot be compiled.
// It simplifies safe initialization of global variables holding compiled JavaScript code.
func MustCompile(name, src string, strict bool) *Program ***REMOVED***
	prg, err := Compile(name, src, strict)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	return prg
***REMOVED***

// Parse takes a source string and produces a parsed AST. Use this function if you want to pass options
// to the parser, e.g.:
//
//  p, err := Parse("test.js", "var a = true", parser.WithDisableSourceMaps)
//  if err != nil ***REMOVED*** /* ... */ ***REMOVED***
//  prg, err := CompileAST(p, true)
//  // ...
//
// Otherwise use Compile which combines both steps.
func Parse(name, src string, options ...parser.Option) (prg *js_ast.Program, err error) ***REMOVED***
	prg, err1 := parser.ParseFile(nil, name, src, 0, options...)
	if err1 != nil ***REMOVED***
		// FIXME offset
		err = &CompilerSyntaxError***REMOVED***
			CompilerError: CompilerError***REMOVED***
				Message: err1.Error(),
			***REMOVED***,
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

func compile(name, src string, strict, inGlobal bool, evalVm *vm, parserOptions ...parser.Option) (p *Program, err error) ***REMOVED***
	prg, err := Parse(name, src, parserOptions...)
	if err != nil ***REMOVED***
		return
	***REMOVED***

	return compileAST(prg, strict, inGlobal, evalVm)
***REMOVED***

func compileAST(prg *js_ast.Program, strict, inGlobal bool, evalVm *vm) (p *Program, err error) ***REMOVED***
	c := newCompiler()

	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			p = nil
			switch x1 := x.(type) ***REMOVED***
			case *CompilerSyntaxError:
				err = x1
			default:
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	c.compile(prg, strict, inGlobal, evalVm)
	p = c.p
	return
***REMOVED***

func (r *Runtime) compile(name, src string, strict, inGlobal bool, evalVm *vm) (p *Program, err error) ***REMOVED***
	p, err = compile(name, src, strict, inGlobal, evalVm, r.parserOptions...)
	if err != nil ***REMOVED***
		switch x1 := err.(type) ***REMOVED***
		case *CompilerSyntaxError:
			err = &Exception***REMOVED***
				val: r.builtin_new(r.global.SyntaxError, []Value***REMOVED***newStringValue(x1.Error())***REMOVED***),
			***REMOVED***
		case *CompilerReferenceError:
			err = &Exception***REMOVED***
				val: r.newError(r.global.ReferenceError, x1.Message),
			***REMOVED*** // TODO proper message
		***REMOVED***
	***REMOVED***
	return
***REMOVED***

// RunString executes the given string in the global context.
func (r *Runtime) RunString(str string) (Value, error) ***REMOVED***
	return r.RunScript("", str)
***REMOVED***

// RunScript executes the given string in the global context.
func (r *Runtime) RunScript(name, src string) (Value, error) ***REMOVED***
	p, err := r.compile(name, src, false, true, nil)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return r.RunProgram(p)
***REMOVED***

// RunProgram executes a pre-compiled (see Compile()) code in the global context.
func (r *Runtime) RunProgram(p *Program) (result Value, err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			if ex, ok := x.(*uncatchableException); ok ***REMOVED***
				err = ex.err
				if len(r.vm.callStack) == 0 ***REMOVED***
					r.leaveAbrupt()
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	vm := r.vm
	recursive := false
	if len(vm.callStack) > 0 ***REMOVED***
		recursive = true
		vm.pushCtx()
		vm.stash = &r.global.stash
		vm.sb = vm.sp - 1
	***REMOVED***
	vm.prg = p
	vm.pc = 0
	vm.result = _undefined
	ex := vm.runTry()
	if ex == nil ***REMOVED***
		result = r.vm.result
	***REMOVED*** else ***REMOVED***
		err = ex
	***REMOVED***
	if recursive ***REMOVED***
		vm.popCtx()
		vm.halt = false
		vm.clearStack()
	***REMOVED*** else ***REMOVED***
		vm.stack = nil
		vm.prg = nil
		vm.funcName = ""
		r.leave()
	***REMOVED***
	return
***REMOVED***

// CaptureCallStack appends the current call stack frames to the stack slice (which may be nil) up to the specified depth.
// The most recent frame will be the first one.
// If depth <= 0 or more than the number of available frames, returns the entire stack.
// This method is not safe for concurrent use and should only be called by a Go function that is
// called from a running script.
func (r *Runtime) CaptureCallStack(depth int, stack []StackFrame) []StackFrame ***REMOVED***
	l := len(r.vm.callStack)
	var offset int
	if depth > 0 ***REMOVED***
		offset = l - depth + 1
		if offset < 0 ***REMOVED***
			offset = 0
		***REMOVED***
	***REMOVED***
	if stack == nil ***REMOVED***
		stack = make([]StackFrame, 0, l-offset+1)
	***REMOVED***
	return r.vm.captureStack(stack, offset)
***REMOVED***

// Interrupt a running JavaScript. The corresponding Go call will return an *InterruptedError containing v.
// If the interrupt propagates until the stack is empty the currently queued promise resolve/reject jobs will be cleared
// without being executed. This is the same time they would be executed otherwise.
// Note, it only works while in JavaScript code, it does not interrupt native Go functions (which includes all built-ins).
// If the runtime is currently not running, it will be immediately interrupted on the next Run*() call.
// To avoid that use ClearInterrupt()
func (r *Runtime) Interrupt(v interface***REMOVED******REMOVED***) ***REMOVED***
	r.vm.Interrupt(v)
***REMOVED***

// ClearInterrupt resets the interrupt flag. Typically this needs to be called before the runtime
// is made available for re-use if there is a chance it could have been interrupted with Interrupt().
// Otherwise if Interrupt() was called when runtime was not running (e.g. if it had already finished)
// so that Interrupt() didn't actually trigger, an attempt to use the runtime will immediately cause
// an interruption. It is up to the user to ensure proper synchronisation so that ClearInterrupt() is
// only called when the runtime has finished and there is no chance of a concurrent Interrupt() call.
func (r *Runtime) ClearInterrupt() ***REMOVED***
	r.vm.ClearInterrupt()
***REMOVED***

/*
ToValue converts a Go value into a JavaScript value of a most appropriate type. Structural types (such as structs, maps
and slices) are wrapped so that changes are reflected on the original value which can be retrieved using Value.Export().

WARNING! These wrapped Go values do not behave in the same way as native ECMAScript values. If you plan to modify
them in ECMAScript, bear in mind the following caveats:

1. If a regular JavaScript Object is assigned as an element of a wrapped Go struct, map or array, it is
Export()'ed and therefore copied. This may result in an unexpected behaviour in JavaScript:

 m := map[string]interface***REMOVED******REMOVED******REMOVED******REMOVED***
 vm.Set("m", m)
 vm.RunString(`
 var obj = ***REMOVED***test: false***REMOVED***;
 m.obj = obj; // obj gets Export()'ed, i.e. copied to a new map[string]interface***REMOVED******REMOVED*** and then this map is set as m["obj"]
 obj.test = true; // note, m.obj.test is still false
 `)
 fmt.Println(m["obj"].(map[string]interface***REMOVED******REMOVED***)["test"]) // prints "false"

2. Be careful with nested non-pointer compound types (structs, slices and arrays) if you modify them in
ECMAScript. Better avoid it at all if possible. One of the fundamental differences between ECMAScript and Go is in
the former all Objects are references whereas in Go you can have a literal struct or array. Consider the following
example:

 type S struct ***REMOVED***
     Field int
 ***REMOVED***

 a := []S***REMOVED******REMOVED***1***REMOVED***, ***REMOVED***2***REMOVED******REMOVED*** // slice of literal structs
 vm.Set("a", &a)
 vm.RunString(`
     let tmp = ***REMOVED***Field: 1***REMOVED***;
     a[0] = tmp;
     a[1] = tmp;
     tmp.Field = 2;
 `)

In ECMAScript one would expect a[0].Field and a[1].Field to be equal to 2, but this is really not possible
(or at least non-trivial without some complex reference tracking).

To cover the most common use cases and to avoid excessive memory allocation, the following 'copy-on-change' mechanism
is implemented (for both arrays and structs):

* When a nested compound value is accessed, the returned ES value becomes a reference to the literal value.
This ensures that things like 'a[0].Field = 1' work as expected and simple access to 'a[0].Field' does not result
in copying of a[0].

* The original container ('a' in our case) keeps track of the returned reference value and if a[0] is reassigned
(e.g. by direct assignment, deletion or shrinking the array) the old a[0] is copied and the earlier returned value
becomes a reference to the copy:

 let tmp = a[0];                      // no copy, tmp is a reference to a[0]
 tmp.Field = 1;                       // a[0].Field === 1 after this
 a[0] = ***REMOVED***Field: 2***REMOVED***;                   // tmp is now a reference to a copy of the old value (with Field === 1)
 a[0].Field === 2 && tmp.Field === 1; // true

* Array value swaps caused by in-place sort (using Array.prototype.sort()) do not count as re-assignments, instead
the references are adjusted to point to the new indices.

* Assignment to an inner compound value always does a copy (and sometimes type conversion):

 a[1] = tmp;    // a[1] is now a copy of tmp
 tmp.Field = 3; // does not affect a[1].Field

3. Non-addressable structs, slices and arrays get copied. This sometimes may lead to a confusion as assigning to
inner fields does not appear to work:

 a1 := []interface***REMOVED******REMOVED******REMOVED***S***REMOVED***1***REMOVED***, S***REMOVED***2***REMOVED******REMOVED***
 vm.Set("a1", &a1)
 vm.RunString(`
    a1[0].Field === 1; // true
    a1[0].Field = 2;
    a1[0].Field === 2; // FALSE, because what it really did was copy a1[0] set its Field to 2 and immediately drop it
 `)

An alternative would be making a1[0].Field a non-writable property which would probably be more in line with
ECMAScript, however it would require to manually copy the value if it does need to be modified which may be
impractical.

Note, the same applies to slices. If a slice is passed by value (not as a pointer), resizing the slice does not reflect on the original
value. Moreover, extending the slice may result in the underlying array being re-allocated and copied.
For example:

 a := []interface***REMOVED******REMOVED******REMOVED***1***REMOVED***
 vm.Set("a", a)
 vm.RunString(`a.push(2); a[0] = 0;`)
 fmt.Println(a[0]) // prints "1"

Notes on individual types:

Primitive types

Primitive types (numbers, string, bool) are converted to the corresponding JavaScript primitives.

Strings

Because of the difference in internal string representation between ECMAScript (which uses UTF-16) and Go (which uses
UTF-8) conversion from JS to Go may be lossy. In particular, code points that can be part of UTF-16 surrogate pairs
(0xD800-0xDFFF) cannot be represented in UTF-8 unless they form a valid surrogate pair and are replaced with
utf8.RuneError.

Nil

Nil is converted to null.

Functions

func(FunctionCall) Value is treated as a native JavaScript function. This increases performance because there are no
automatic argument and return value type conversions (which involves reflect). Attempting to use
the function as a constructor will result in a TypeError.

func(FunctionCall, *Runtime) Value is treated as above, except the *Runtime is also passed as a parameter.

func(ConstructorCall) *Object is treated as a native constructor, allowing to use it with the new
operator:

 func MyObject(call goja.ConstructorCall) *goja.Object ***REMOVED***
    // call.This contains the newly created object as per http://www.ecma-international.org/ecma-262/5.1/index.html#sec-13.2.2
    // call.Arguments contain arguments passed to the function

    call.This.Set("method", method)

    //...

    // If return value is a non-nil *Object, it will be used instead of call.This
    // This way it is possible to return a Go struct or a map converted
    // into goja.Value using ToValue(), however in this case
    // instanceof will not work as expected, unless you set the prototype:
    //
    // instance := &myCustomStruct***REMOVED******REMOVED***
    // instanceValue := vm.ToValue(instance).(*Object)
    // instanceValue.SetPrototype(call.This.Prototype())
    // return instanceValue
    return nil
 ***REMOVED***

 runtime.Set("MyObject", MyObject)

Then it can be used in JS as follows:

 var o = new MyObject(arg);
 var o1 = MyObject(arg); // same thing
 o instanceof MyObject && o1 instanceof MyObject; // true

When a native constructor is called directly (without the new operator) its behavior depends on
this value: if it's an Object, it is passed through, otherwise a new one is created exactly as
if it was called with the new operator. In either case call.NewTarget will be nil.

func(ConstructorCall, *Runtime) *Object is treated as above, except the *Runtime is also passed as a parameter.

Any other Go function is wrapped so that the arguments are automatically converted into the required Go types and the
return value is converted to a JavaScript value (using this method).  If conversion is not possible, a TypeError is
thrown.

Functions with multiple return values return an Array. If the last return value is an `error` it is not returned but
converted into a JS exception. If the error is *Exception, it is thrown as is, otherwise it's wrapped in a GoEerror.
Note that if there are exactly two return values and the last is an `error`, the function returns the first value as is,
not an Array.

Structs

Structs are converted to Object-like values. Fields and methods are available as properties, their values are
results of this method (ToValue()) applied to the corresponding Go value.

Field properties are writable and non-configurable. Method properties are non-writable and non-configurable.

Attempt to define a new property or delete an existing property will fail (throw in strict mode) unless it's a Symbol
property. Symbol properties only exist in the wrapper and do not affect the underlying Go value.
Note that because a wrapper is created every time a property is accessed it may lead to unexpected results such as this:

 type Field struct***REMOVED***
 ***REMOVED***
 type S struct ***REMOVED***
	Field *Field
 ***REMOVED***
 var s = S***REMOVED***
	Field: &Field***REMOVED******REMOVED***,
 ***REMOVED***
 vm := New()
 vm.Set("s", &s)
 res, err := vm.RunString(`
 var sym = Symbol(66);
 var field1 = s.Field;
 field1[sym] = true;
 var field2 = s.Field;
 field1 === field2; // true, because the equality operation compares the wrapped values, not the wrappers
 field1[sym] === true; // true
 field2[sym] === undefined; // also true
 `)

The same applies to values from maps and slices as well.

Handling of time.Time

time.Time does not get special treatment and therefore is converted just like any other `struct` providing access to
all its methods. This is done deliberately instead of converting it to a `Date` because these two types are not fully
compatible: `time.Time` includes zone, whereas JS `Date` doesn't. Doing the conversion implicitly therefore would
result in a loss of information.

If you need to convert it to a `Date`, it can be done either in JS:

 var d = new Date(goval.UnixNano()/1e6);

... or in Go:

 now := time.Now()
 vm := New()
 val, err := vm.New(vm.Get("Date").ToObject(vm), vm.ToValue(now.UnixNano()/1e6))
 if err != nil ***REMOVED***
	...
 ***REMOVED***
 vm.Set("d", val)

Note that Value.Export() for a `Date` value returns time.Time in local timezone.

Maps

Maps with string or integer key type are converted into host objects that largely behave like a JavaScript Object.

Maps with methods

If a map type has at least one method defined, the properties of the resulting Object represent methods, not map keys.
This is because in JavaScript there is no distinction between 'object.key` and `object[key]`, unlike Go.
If access to the map values is required, it can be achieved by defining another method or, if it's not possible, by
defining an external getter function.

Slices

Slices are converted into host objects that behave largely like JavaScript Array. It has the appropriate
prototype and all the usual methods should work. There is, however, a caveat: converted Arrays may not contain holes
(because Go slices cannot). This means that hasOwnProperty(n) always returns `true` if n < length. Deleting an item with
an index < length will set it to a zero value (but the property will remain). Nil slice elements are be converted to
`null`. Accessing an element beyond `length` returns `undefined`. Also see the warning above about passing slices as
values (as opposed to pointers).

Arrays

Arrays are converted similarly to slices, except the resulting Arrays are not resizable (and therefore the 'length'
property is non-writable).

Any other type is converted to a generic reflect based host object. Depending on the underlying type it behaves similar
to a Number, String, Boolean or Object.

Note that the underlying type is not lost, calling Export() returns the original Go value. This applies to all
reflect based types.
*/
func (r *Runtime) ToValue(i interface***REMOVED******REMOVED***) Value ***REMOVED***
	switch i := i.(type) ***REMOVED***
	case nil:
		return _null
	case *Object:
		if i == nil || i.self == nil ***REMOVED***
			return _null
		***REMOVED***
		if i.runtime != nil && i.runtime != r ***REMOVED***
			panic(r.NewTypeError("Illegal runtime transition of an Object"))
		***REMOVED***
		return i
	case valueContainer:
		return i.toValue(r)
	case Value:
		return i
	case string:
		return newStringValue(i)
	case bool:
		if i ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	case func(FunctionCall) Value:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
		return r.newNativeFunc(i, nil, name, nil, 0)
	case func(FunctionCall, *Runtime) Value:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
		return r.newNativeFunc(func(call FunctionCall) Value ***REMOVED***
			return i(call, r)
		***REMOVED***, nil, name, nil, 0)
	case func(ConstructorCall) *Object:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
		return r.newNativeConstructor(i, name, 0)
	case func(ConstructorCall, *Runtime) *Object:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
		return r.newNativeConstructor(func(call ConstructorCall) *Object ***REMOVED***
			return i(call, r)
		***REMOVED***, name, 0)
	case int:
		return intToValue(int64(i))
	case int8:
		return intToValue(int64(i))
	case int16:
		return intToValue(int64(i))
	case int32:
		return intToValue(int64(i))
	case int64:
		return intToValue(i)
	case uint:
		if uint64(i) <= math.MaxInt64 ***REMOVED***
			return intToValue(int64(i))
		***REMOVED*** else ***REMOVED***
			return floatToValue(float64(i))
		***REMOVED***
	case uint8:
		return intToValue(int64(i))
	case uint16:
		return intToValue(int64(i))
	case uint32:
		return intToValue(int64(i))
	case uint64:
		if i <= math.MaxInt64 ***REMOVED***
			return intToValue(int64(i))
		***REMOVED***
		return floatToValue(float64(i))
	case float32:
		return floatToValue(float64(i))
	case float64:
		return floatToValue(i)
	case map[string]interface***REMOVED******REMOVED***:
		if i == nil ***REMOVED***
			return _null
		***REMOVED***
		obj := &Object***REMOVED***runtime: r***REMOVED***
		m := &objectGoMapSimple***REMOVED***
			baseObject: baseObject***REMOVED***
				val:        obj,
				extensible: true,
			***REMOVED***,
			data: i,
		***REMOVED***
		obj.self = m
		m.init()
		return obj
	case []interface***REMOVED******REMOVED***:
		if i == nil ***REMOVED***
			return _null
		***REMOVED***
		return r.newObjectGoSlice(&i).val
	case *[]interface***REMOVED******REMOVED***:
		if i == nil ***REMOVED***
			return _null
		***REMOVED***
		return r.newObjectGoSlice(i).val
	***REMOVED***

	return r.reflectValueToValue(reflect.ValueOf(i))
***REMOVED***

func (r *Runtime) reflectValueToValue(origValue reflect.Value) Value ***REMOVED***
	value := origValue
	for value.Kind() == reflect.Ptr ***REMOVED***
		value = reflect.Indirect(value)
	***REMOVED***

	if !value.IsValid() ***REMOVED***
		return _null
	***REMOVED***

	switch value.Kind() ***REMOVED***
	case reflect.Map:
		if value.Type().NumMethod() == 0 ***REMOVED***
			switch value.Type().Key().Kind() ***REMOVED***
			case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float64, reflect.Float32:

				obj := &Object***REMOVED***runtime: r***REMOVED***
				m := &objectGoMapReflect***REMOVED***
					objectGoReflect: objectGoReflect***REMOVED***
						baseObject: baseObject***REMOVED***
							val:        obj,
							extensible: true,
						***REMOVED***,
						origValue: origValue,
						value:     value,
					***REMOVED***,
				***REMOVED***
				m.init()
				obj.self = m
				return obj
			***REMOVED***
		***REMOVED***
	case reflect.Array:
		obj := &Object***REMOVED***runtime: r***REMOVED***
		a := &objectGoArrayReflect***REMOVED***
			objectGoReflect: objectGoReflect***REMOVED***
				baseObject: baseObject***REMOVED***
					val: obj,
				***REMOVED***,
				origValue: origValue,
				value:     value,
			***REMOVED***,
		***REMOVED***
		a.init()
		obj.self = a
		return obj
	case reflect.Slice:
		obj := &Object***REMOVED***runtime: r***REMOVED***
		a := &objectGoSliceReflect***REMOVED***
			objectGoArrayReflect: objectGoArrayReflect***REMOVED***
				objectGoReflect: objectGoReflect***REMOVED***
					baseObject: baseObject***REMOVED***
						val: obj,
					***REMOVED***,
					origValue: origValue,
					value:     value,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***
		a.init()
		obj.self = a
		return obj
	case reflect.Func:
		name := unistring.NewFromString(runtime.FuncForPC(value.Pointer()).Name())
		return r.newNativeFunc(r.wrapReflectFunc(value), nil, name, nil, value.Type().NumIn())
	***REMOVED***

	obj := &Object***REMOVED***runtime: r***REMOVED***
	o := &objectGoReflect***REMOVED***
		baseObject: baseObject***REMOVED***
			val: obj,
		***REMOVED***,
		origValue: origValue,
		value:     value,
	***REMOVED***
	obj.self = o
	o.init()
	return obj
***REMOVED***

func (r *Runtime) wrapReflectFunc(value reflect.Value) func(FunctionCall) Value ***REMOVED***
	return func(call FunctionCall) Value ***REMOVED***
		typ := value.Type()
		nargs := typ.NumIn()
		var in []reflect.Value

		if l := len(call.Arguments); l < nargs ***REMOVED***
			// fill missing arguments with zero values
			n := nargs
			if typ.IsVariadic() ***REMOVED***
				n--
			***REMOVED***
			in = make([]reflect.Value, n)
			for i := l; i < n; i++ ***REMOVED***
				in[i] = reflect.Zero(typ.In(i))
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			if l > nargs && !typ.IsVariadic() ***REMOVED***
				l = nargs
			***REMOVED***
			in = make([]reflect.Value, l)
		***REMOVED***

		for i, a := range call.Arguments ***REMOVED***
			var t reflect.Type

			n := i
			if n >= nargs-1 && typ.IsVariadic() ***REMOVED***
				if n > nargs-1 ***REMOVED***
					n = nargs - 1
				***REMOVED***

				t = typ.In(n).Elem()
			***REMOVED*** else if n > nargs-1 ***REMOVED*** // ignore extra arguments
				break
			***REMOVED*** else ***REMOVED***
				t = typ.In(n)
			***REMOVED***

			v := reflect.New(t).Elem()
			err := r.toReflectValue(a, v, &objectExportCtx***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				panic(r.NewTypeError("could not convert function call parameter %d: %v", i, err))
			***REMOVED***
			in[i] = v
		***REMOVED***

		out := value.Call(in)
		if len(out) == 0 ***REMOVED***
			return _undefined
		***REMOVED***

		if last := out[len(out)-1]; last.Type().Name() == "error" ***REMOVED***
			if !last.IsNil() ***REMOVED***
				err := last.Interface()
				if _, ok := err.(*Exception); ok ***REMOVED***
					panic(err)
				***REMOVED***
				panic(r.NewGoError(last.Interface().(error)))
			***REMOVED***
			out = out[:len(out)-1]
		***REMOVED***

		switch len(out) ***REMOVED***
		case 0:
			return _undefined
		case 1:
			return r.ToValue(out[0].Interface())
		default:
			s := make([]interface***REMOVED******REMOVED***, len(out))
			for i, v := range out ***REMOVED***
				s[i] = v.Interface()
			***REMOVED***

			return r.ToValue(s)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) toReflectValue(v Value, dst reflect.Value, ctx *objectExportCtx) error ***REMOVED***
	typ := dst.Type()

	if typ == typeValue ***REMOVED***
		dst.Set(reflect.ValueOf(v))
		return nil
	***REMOVED***

	if typ == typeObject ***REMOVED***
		if obj, ok := v.(*Object); ok ***REMOVED***
			dst.Set(reflect.ValueOf(obj))
			return nil
		***REMOVED***
	***REMOVED***

	if typ == typeCallable ***REMOVED***
		if fn, ok := AssertFunction(v); ok ***REMOVED***
			dst.Set(reflect.ValueOf(fn))
			return nil
		***REMOVED***
	***REMOVED***

	et := v.ExportType()
	if et == nil || et == reflectTypeNil ***REMOVED***
		dst.Set(reflect.Zero(typ))
		return nil
	***REMOVED***

	kind := typ.Kind()
	for i := 0; ; i++ ***REMOVED***
		if et.AssignableTo(typ) ***REMOVED***
			ev := reflect.ValueOf(exportValue(v, ctx))
			for ; i > 0; i-- ***REMOVED***
				ev = ev.Elem()
			***REMOVED***
			dst.Set(ev)
			return nil
		***REMOVED***
		expKind := et.Kind()
		if expKind == kind && et.ConvertibleTo(typ) || expKind == reflect.String && typ == typeBytes ***REMOVED***
			ev := reflect.ValueOf(exportValue(v, ctx))
			for ; i > 0; i-- ***REMOVED***
				ev = ev.Elem()
			***REMOVED***
			dst.Set(ev.Convert(typ))
			return nil
		***REMOVED***
		if expKind == reflect.Ptr ***REMOVED***
			et = et.Elem()
		***REMOVED*** else ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	if typ == typeTime ***REMOVED***
		if obj, ok := v.(*Object); ok ***REMOVED***
			if d, ok := obj.self.(*dateObject); ok ***REMOVED***
				dst.Set(reflect.ValueOf(d.time()))
				return nil
			***REMOVED***
		***REMOVED***
		if et.Kind() == reflect.String ***REMOVED***
			tme, ok := dateParse(v.String())
			if !ok ***REMOVED***
				return fmt.Errorf("could not convert string %v to %v", v, typ)
			***REMOVED***
			dst.Set(reflect.ValueOf(tme))
			return nil
		***REMOVED***
	***REMOVED***

	switch kind ***REMOVED***
	case reflect.String:
		dst.Set(reflect.ValueOf(v.String()).Convert(typ))
		return nil
	case reflect.Bool:
		dst.Set(reflect.ValueOf(v.ToBoolean()).Convert(typ))
		return nil
	case reflect.Int:
		dst.Set(reflect.ValueOf(toInt(v)).Convert(typ))
		return nil
	case reflect.Int64:
		dst.Set(reflect.ValueOf(toInt64(v)).Convert(typ))
		return nil
	case reflect.Int32:
		dst.Set(reflect.ValueOf(toInt32(v)).Convert(typ))
		return nil
	case reflect.Int16:
		dst.Set(reflect.ValueOf(toInt16(v)).Convert(typ))
		return nil
	case reflect.Int8:
		dst.Set(reflect.ValueOf(toInt8(v)).Convert(typ))
		return nil
	case reflect.Uint:
		dst.Set(reflect.ValueOf(toUint(v)).Convert(typ))
		return nil
	case reflect.Uint64:
		dst.Set(reflect.ValueOf(toUint64(v)).Convert(typ))
		return nil
	case reflect.Uint32:
		dst.Set(reflect.ValueOf(toUint32(v)).Convert(typ))
		return nil
	case reflect.Uint16:
		dst.Set(reflect.ValueOf(toUint16(v)).Convert(typ))
		return nil
	case reflect.Uint8:
		dst.Set(reflect.ValueOf(toUint8(v)).Convert(typ))
		return nil
	case reflect.Float64:
		dst.Set(reflect.ValueOf(v.ToFloat()).Convert(typ))
		return nil
	case reflect.Float32:
		dst.Set(reflect.ValueOf(toFloat32(v)).Convert(typ))
		return nil
	case reflect.Slice, reflect.Array:
		if o, ok := v.(*Object); ok ***REMOVED***
			if v, exists := ctx.getTyped(o, typ); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v))
				return nil
			***REMOVED***
			return o.self.exportToArrayOrSlice(dst, typ, ctx)
		***REMOVED***
	case reflect.Map:
		if o, ok := v.(*Object); ok ***REMOVED***
			if v, exists := ctx.getTyped(o, typ); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v))
				return nil
			***REMOVED***
			return o.self.exportToMap(dst, typ, ctx)
		***REMOVED***
	case reflect.Struct:
		if o, ok := v.(*Object); ok ***REMOVED***
			t := reflect.PtrTo(typ)
			if v, exists := ctx.getTyped(o, t); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v).Elem())
				return nil
			***REMOVED***
			s := dst
			ctx.putTyped(o, t, s.Addr().Interface())
			for i := 0; i < typ.NumField(); i++ ***REMOVED***
				field := typ.Field(i)
				if ast.IsExported(field.Name) ***REMOVED***
					name := field.Name
					if r.fieldNameMapper != nil ***REMOVED***
						name = r.fieldNameMapper.FieldName(typ, field)
					***REMOVED***
					var v Value
					if field.Anonymous ***REMOVED***
						v = o
					***REMOVED*** else ***REMOVED***
						v = o.self.getStr(unistring.NewFromString(name), nil)
					***REMOVED***

					if v != nil ***REMOVED***
						err := r.toReflectValue(v, s.Field(i), ctx)
						if err != nil ***REMOVED***
							return fmt.Errorf("could not convert struct value %v to %v for field %s: %w", v, field.Type, field.Name, err)
						***REMOVED***
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return nil
		***REMOVED***
	case reflect.Func:
		if fn, ok := AssertFunction(v); ok ***REMOVED***
			dst.Set(reflect.MakeFunc(typ, r.wrapJSFunc(fn, typ)))
			return nil
		***REMOVED***
	case reflect.Ptr:
		if o, ok := v.(*Object); ok ***REMOVED***
			if v, exists := ctx.getTyped(o, typ); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v))
				return nil
			***REMOVED***
		***REMOVED***
		if dst.IsNil() ***REMOVED***
			dst.Set(reflect.New(typ.Elem()))
		***REMOVED***
		return r.toReflectValue(v, dst.Elem(), ctx)
	***REMOVED***

	return fmt.Errorf("could not convert %v to %v", v, typ)
***REMOVED***

func (r *Runtime) wrapJSFunc(fn Callable, typ reflect.Type) func(args []reflect.Value) (results []reflect.Value) ***REMOVED***
	return func(args []reflect.Value) (results []reflect.Value) ***REMOVED***
		jsArgs := make([]Value, len(args))
		for i, arg := range args ***REMOVED***
			jsArgs[i] = r.ToValue(arg.Interface())
		***REMOVED***

		results = make([]reflect.Value, typ.NumOut())
		res, err := fn(_undefined, jsArgs...)
		if err == nil ***REMOVED***
			if typ.NumOut() > 0 ***REMOVED***
				v := reflect.New(typ.Out(0)).Elem()
				err = r.toReflectValue(res, v, &objectExportCtx***REMOVED******REMOVED***)
				if err == nil ***REMOVED***
					results[0] = v
				***REMOVED***
			***REMOVED***
		***REMOVED***

		if err != nil ***REMOVED***
			if typ.NumOut() == 2 && typ.Out(1).Name() == "error" ***REMOVED***
				results[1] = reflect.ValueOf(err).Convert(typ.Out(1))
			***REMOVED*** else ***REMOVED***
				panic(err)
			***REMOVED***
		***REMOVED***

		for i, v := range results ***REMOVED***
			if !v.IsValid() ***REMOVED***
				results[i] = reflect.Zero(typ.Out(i))
			***REMOVED***
		***REMOVED***

		return
	***REMOVED***
***REMOVED***

// ExportTo converts a JavaScript value into the specified Go value. The second parameter must be a non-nil pointer.
// Returns error if conversion is not possible.
//
// Notes on specific cases:
//
// Empty interface
//
// Exporting to an interface***REMOVED******REMOVED*** results in a value of the same type as Value.Export() would produce.
//
// Numeric types
//
// Exporting to numeric types uses the standard ECMAScript conversion operations, same as used when assigning
// values to non-clamped typed array items, e.g. https://262.ecma-international.org/#sec-toint32.
//
// Functions
//
// Exporting to a 'func' creates a strictly typed 'gateway' into an ES function which can be called from Go.
// The arguments are converted into ES values using Runtime.ToValue(). If the func has no return values,
// the return value is ignored. If the func has exactly one return value, it is converted to the appropriate
// type using ExportTo(). If the func has exactly 2 return values and the second value is 'error', exceptions
// are caught and returned as *Exception. In all other cases exceptions result in a panic. Any extra return values
// are zeroed.
//
// Note, if you want to catch and return exceptions as an `error` and you don't need the return value,
// 'func(...) error' will not work as expected. The 'error' in this case is mapped to the function return value, not
// the exception which will still result in a panic. Use 'func(...) (Value, error)' instead, and ignore the Value.
//
// 'this' value will always be set to 'undefined'.
//
// For a more low-level mechanism see AssertFunction().
//
// Map types
//
// An ES Map can be exported into a Go map type. If any exported key value is non-hashable, the operation panics
// (as reflect.Value.SetMapIndex() would). Symbol.iterator is ignored.
//
// Exporting an ES Set into a map type results in the map being populated with (element) -> (zero value) key/value
// pairs. If any value is non-hashable, the operation panics (as reflect.Value.SetMapIndex() would).
// Symbol.iterator is ignored.
//
// Any other Object populates the map with own enumerable non-symbol properties.
//
// Slice types
//
// Exporting an ES Set into a slice type results in its elements being exported.
//
// Exporting any Object that implements the iterable protocol (https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Iteration_protocols#the_iterable_protocol)
// into a slice type results in the slice being populated with the results of the iteration.
//
// Array is treated as iterable (i.e. overwriting Symbol.iterator affects the result).
//
// If an object has a 'length' property and is not a function it is treated as array-like. The resulting slice
// will contain obj[0], ... obj[length-1].
//
// For any other Object an error is returned.
//
// Array types
//
// Anything that can be exported to a slice type can also be exported to an array type, as long as the lengths
// match. If they do not, an error is returned.
//
// Proxy
//
// Proxy objects are treated the same way as if they were accessed from ES code in regard to their properties
// (such as 'length' or [Symbol.iterator]). This means exporting them to slice types works, however
// exporting a proxied Map into a map type does not produce its contents, because the Proxy is not recognised
// as a Map. Same applies to a proxied Set.
//
func (r *Runtime) ExportTo(v Value, target interface***REMOVED******REMOVED***) error ***REMOVED***
	tval := reflect.ValueOf(target)
	if tval.Kind() != reflect.Ptr || tval.IsNil() ***REMOVED***
		return errors.New("target must be a non-nil pointer")
	***REMOVED***
	return r.toReflectValue(v, tval.Elem(), &objectExportCtx***REMOVED******REMOVED***)
***REMOVED***

// GlobalObject returns the global object.
func (r *Runtime) GlobalObject() *Object ***REMOVED***
	return r.globalObject
***REMOVED***

// Set the specified variable in the global context.
// Equivalent to running "name = value" in non-strict mode.
// The value is first converted using ToValue().
// Note, this is not the same as GlobalObject().Set(name, value),
// because if a global lexical binding (let or const) exists, it is set instead.
func (r *Runtime) Set(name string, value interface***REMOVED******REMOVED***) error ***REMOVED***
	return r.try(func() ***REMOVED***
		name := unistring.NewFromString(name)
		v := r.ToValue(value)
		if ref := r.global.stash.getRefByName(name, false); ref != nil ***REMOVED***
			ref.set(v)
		***REMOVED*** else ***REMOVED***
			r.globalObject.self.setOwnStr(name, v, true)
		***REMOVED***
	***REMOVED***)
***REMOVED***

// Get the specified variable in the global context.
// Equivalent to dereferencing a variable by name in non-strict mode. If variable is not defined returns nil.
// Note, this is not the same as GlobalObject().Get(name),
// because if a global lexical binding (let or const) exists, it is used instead.
// This method will panic with an *Exception if a JavaScript exception is thrown in the process.
func (r *Runtime) Get(name string) (ret Value) ***REMOVED***
	r.tryPanic(func() ***REMOVED***
		n := unistring.NewFromString(name)
		if v, exists := r.global.stash.getByName(n); exists ***REMOVED***
			ret = v
		***REMOVED*** else ***REMOVED***
			ret = r.globalObject.self.getStr(n, nil)
		***REMOVED***
	***REMOVED***)
	return
***REMOVED***

// SetRandSource sets random source for this Runtime. If not called, the default math/rand is used.
func (r *Runtime) SetRandSource(source RandSource) ***REMOVED***
	r.rand = source
***REMOVED***

// SetTimeSource sets the current time source for this Runtime.
// If not called, the default time.Now() is used.
func (r *Runtime) SetTimeSource(now Now) ***REMOVED***
	r.now = now
***REMOVED***

// SetParserOptions sets parser options to be used by RunString, RunScript and eval() within the code.
func (r *Runtime) SetParserOptions(opts ...parser.Option) ***REMOVED***
	r.parserOptions = opts
***REMOVED***

// SetMaxCallStackSize sets the maximum function call depth. When exceeded, a *StackOverflowError is thrown and
// returned by RunProgram or by a Callable call. This is useful to prevent memory exhaustion caused by an
// infinite recursion. The default value is math.MaxInt32.
// This method (as the rest of the Set* methods) is not safe for concurrent use and may only be called
// from the vm goroutine or when the vm is not running.
func (r *Runtime) SetMaxCallStackSize(size int) ***REMOVED***
	r.vm.maxCallStackSize = size
***REMOVED***

// New is an equivalent of the 'new' operator allowing to call it directly from Go.
func (r *Runtime) New(construct Value, args ...Value) (o *Object, err error) ***REMOVED***
	err = r.try(func() ***REMOVED***
		o = r.builtin_new(r.toObject(construct), args)
	***REMOVED***)
	return
***REMOVED***

// Callable represents a JavaScript function that can be called from Go.
type Callable func(this Value, args ...Value) (Value, error)

// AssertFunction checks if the Value is a function and returns a Callable.
// Note, for classes this returns a callable and a 'true', however calling it will always result in a TypeError.
// For classes use AssertConstructor().
func AssertFunction(v Value) (Callable, bool) ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if f, ok := obj.self.assertCallable(); ok ***REMOVED***
			return func(this Value, args ...Value) (ret Value, err error) ***REMOVED***
				err = obj.runtime.runWrapped(func() ***REMOVED***
					ret = f(FunctionCall***REMOVED***
						This:      this,
						Arguments: args,
					***REMOVED***)
				***REMOVED***)
				return
			***REMOVED***, true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

// Constructor is a type that can be used to call constructors. The first argument (newTarget) can be nil
// which sets it to the constructor function itself.
type Constructor func(newTarget *Object, args ...Value) (*Object, error)

// AssertConstructor checks if the Value is a constructor and returns a Constructor.
func AssertConstructor(v Value) (Constructor, bool) ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if ctor := obj.self.assertConstructor(); ctor != nil ***REMOVED***
			return func(newTarget *Object, args ...Value) (ret *Object, err error) ***REMOVED***
				err = obj.runtime.runWrapped(func() ***REMOVED***
					ret = ctor(args, newTarget)
				***REMOVED***)
				return
			***REMOVED***, true
		***REMOVED***
	***REMOVED***
	return nil, false
***REMOVED***

func (r *Runtime) runWrapped(f func()) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			if ex, ok := x.(*uncatchableException); ok ***REMOVED***
				err = ex.err
				if len(r.vm.callStack) == 0 ***REMOVED***
					r.leaveAbrupt()
				***REMOVED***
			***REMOVED*** else ***REMOVED***
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	ex := r.vm.try(f)
	if ex != nil ***REMOVED***
		err = ex
	***REMOVED***
	r.vm.clearStack()
	if len(r.vm.callStack) == 0 ***REMOVED***
		r.leave()
	***REMOVED***
	return
***REMOVED***

// IsUndefined returns true if the supplied Value is undefined. Note, it checks against the real undefined, not
// against the global object's 'undefined' property.
func IsUndefined(v Value) bool ***REMOVED***
	return v == _undefined
***REMOVED***

// IsNull returns true if the supplied Value is null.
func IsNull(v Value) bool ***REMOVED***
	return v == _null
***REMOVED***

// IsNaN returns true if the supplied value is NaN.
func IsNaN(v Value) bool ***REMOVED***
	f, ok := v.(valueFloat)
	return ok && math.IsNaN(float64(f))
***REMOVED***

// IsInfinity returns true if the supplied is (+/-)Infinity
func IsInfinity(v Value) bool ***REMOVED***
	return v == _positiveInf || v == _negativeInf
***REMOVED***

// Undefined returns JS undefined value. Note if global 'undefined' property is changed this still returns the original value.
func Undefined() Value ***REMOVED***
	return _undefined
***REMOVED***

// Null returns JS null value.
func Null() Value ***REMOVED***
	return _null
***REMOVED***

// NaN returns a JS NaN value.
func NaN() Value ***REMOVED***
	return _NaN
***REMOVED***

// PositiveInf returns a JS +Inf value.
func PositiveInf() Value ***REMOVED***
	return _positiveInf
***REMOVED***

// NegativeInf returns a JS -Inf value.
func NegativeInf() Value ***REMOVED***
	return _negativeInf
***REMOVED***

func tryFunc(f func()) (ret interface***REMOVED******REMOVED***) ***REMOVED***
	defer func() ***REMOVED***
		ret = recover()
	***REMOVED***()

	f()
	return
***REMOVED***

func (r *Runtime) try(f func()) error ***REMOVED***
	if ex := r.vm.try(f); ex != nil ***REMOVED***
		return ex
	***REMOVED***
	return nil
***REMOVED***

func (r *Runtime) tryPanic(f func()) ***REMOVED***
	if ex := r.vm.try(f); ex != nil ***REMOVED***
		panic(ex)
	***REMOVED***
***REMOVED***

func (r *Runtime) toObject(v Value, args ...interface***REMOVED******REMOVED***) *Object ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj
	***REMOVED***
	if len(args) > 0 ***REMOVED***
		panic(r.NewTypeError(args...))
	***REMOVED*** else ***REMOVED***
		var s string
		if v == nil ***REMOVED***
			s = "undefined"
		***REMOVED*** else ***REMOVED***
			s = v.String()
		***REMOVED***
		panic(r.NewTypeError("Value is not an object: %s", s))
	***REMOVED***
***REMOVED***

func (r *Runtime) toNumber(v Value) Value ***REMOVED***
	switch o := v.(type) ***REMOVED***
	case valueInt, valueFloat:
		return v
	case *Object:
		if pvo, ok := o.self.(*primitiveValueObject); ok ***REMOVED***
			return r.toNumber(pvo.pValue)
		***REMOVED***
	***REMOVED***
	panic(r.NewTypeError("Value is not a number: %s", v))
***REMOVED***

func (r *Runtime) speciesConstructor(o, defaultConstructor *Object) func(args []Value, newTarget *Object) *Object ***REMOVED***
	c := o.self.getStr("constructor", nil)
	if c != nil && c != _undefined ***REMOVED***
		c = r.toObject(c).self.getSym(SymSpecies, nil)
	***REMOVED***
	if c == nil || c == _undefined || c == _null ***REMOVED***
		c = defaultConstructor
	***REMOVED***
	return r.toConstructor(c)
***REMOVED***

func (r *Runtime) speciesConstructorObj(o, defaultConstructor *Object) *Object ***REMOVED***
	c := o.self.getStr("constructor", nil)
	if c != nil && c != _undefined ***REMOVED***
		c = r.toObject(c).self.getSym(SymSpecies, nil)
	***REMOVED***
	if c == nil || c == _undefined || c == _null ***REMOVED***
		return defaultConstructor
	***REMOVED***
	obj := r.toObject(c)
	if obj.self.assertConstructor() == nil ***REMOVED***
		panic(r.NewTypeError("Value is not a constructor"))
	***REMOVED***
	return obj
***REMOVED***

func (r *Runtime) returnThis(call FunctionCall) Value ***REMOVED***
	return call.This
***REMOVED***

func createDataProperty(o *Object, p Value, v Value) ***REMOVED***
	o.defineOwnProperty(p, PropertyDescriptor***REMOVED***
		Writable:     FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
		Configurable: FLAG_TRUE,
		Value:        v,
	***REMOVED***, false)
***REMOVED***

func createDataPropertyOrThrow(o *Object, p Value, v Value) ***REMOVED***
	o.defineOwnProperty(p, PropertyDescriptor***REMOVED***
		Writable:     FLAG_TRUE,
		Enumerable:   FLAG_TRUE,
		Configurable: FLAG_TRUE,
		Value:        v,
	***REMOVED***, true)
***REMOVED***

func toPropertyKey(key Value) Value ***REMOVED***
	return key.ToString()
***REMOVED***

func (r *Runtime) getVStr(v Value, p unistring.String) Value ***REMOVED***
	o := v.ToObject(r)
	return o.self.getStr(p, v)
***REMOVED***

func (r *Runtime) getV(v Value, p Value) Value ***REMOVED***
	o := v.ToObject(r)
	return o.get(p, v)
***REMOVED***

type iteratorRecord struct ***REMOVED***
	iterator *Object
	next     func(FunctionCall) Value
***REMOVED***

func (r *Runtime) getIterator(obj Value, method func(FunctionCall) Value) *iteratorRecord ***REMOVED***
	if method == nil ***REMOVED***
		method = toMethod(r.getV(obj, SymIterator))
		if method == nil ***REMOVED***
			panic(r.NewTypeError("object is not iterable"))
		***REMOVED***
	***REMOVED***

	iter := r.toObject(method(FunctionCall***REMOVED***
		This: obj,
	***REMOVED***))
	next := toMethod(iter.self.getStr("next", nil))
	return &iteratorRecord***REMOVED***
		iterator: iter,
		next:     next,
	***REMOVED***
***REMOVED***

func (ir *iteratorRecord) iterate(step func(Value)) ***REMOVED***
	r := ir.iterator.runtime
	for ***REMOVED***
		res := r.toObject(ir.next(FunctionCall***REMOVED***This: ir.iterator***REMOVED***))
		if nilSafe(res.self.getStr("done", nil)).ToBoolean() ***REMOVED***
			break
		***REMOVED***
		value := nilSafe(res.self.getStr("value", nil))
		ret := tryFunc(func() ***REMOVED***
			step(value)
		***REMOVED***)
		if ret != nil ***REMOVED***
			_ = tryFunc(func() ***REMOVED***
				ir.returnIter()
			***REMOVED***)
			panic(ret)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (ir *iteratorRecord) step() (value Value, ex *Exception) ***REMOVED***
	r := ir.iterator.runtime
	ex = r.vm.try(func() ***REMOVED***
		res := r.toObject(ir.next(FunctionCall***REMOVED***This: ir.iterator***REMOVED***))
		done := nilSafe(res.self.getStr("done", nil)).ToBoolean()
		if !done ***REMOVED***
			value = nilSafe(res.self.getStr("value", nil))
		***REMOVED*** else ***REMOVED***
			ir.close()
		***REMOVED***
	***REMOVED***)
	return
***REMOVED***

func (ir *iteratorRecord) returnIter() ***REMOVED***
	if ir.iterator == nil ***REMOVED***
		return
	***REMOVED***
	retMethod := toMethod(ir.iterator.self.getStr("return", nil))
	if retMethod != nil ***REMOVED***
		ir.iterator.runtime.toObject(retMethod(FunctionCall***REMOVED***This: ir.iterator***REMOVED***))
	***REMOVED***
	ir.iterator = nil
	ir.next = nil
***REMOVED***

func (ir *iteratorRecord) close() ***REMOVED***
	ir.iterator = nil
	ir.next = nil
***REMOVED***

func (r *Runtime) createIterResultObject(value Value, done bool) Value ***REMOVED***
	o := r.NewObject()
	o.self.setOwnStr("value", value, false)
	o.self.setOwnStr("done", r.toBoolean(done), false)
	return o
***REMOVED***

func (r *Runtime) newLazyObject(create func(*Object) objectImpl) *Object ***REMOVED***
	val := &Object***REMOVED***runtime: r***REMOVED***
	o := &lazyObject***REMOVED***
		val:    val,
		create: create,
	***REMOVED***
	val.self = o
	return val
***REMOVED***

func (r *Runtime) getHash() *maphash.Hash ***REMOVED***
	if r.hash == nil ***REMOVED***
		r.hash = &maphash.Hash***REMOVED******REMOVED***
	***REMOVED***
	return r.hash
***REMOVED***

// called when the top level function returns normally (i.e. control is passed outside the Runtime).
func (r *Runtime) leave() ***REMOVED***
	for ***REMOVED***
		jobs := r.jobQueue
		r.jobQueue = nil
		if len(jobs) == 0 ***REMOVED***
			break
		***REMOVED***
		for _, job := range jobs ***REMOVED***
			job()
		***REMOVED***
	***REMOVED***
***REMOVED***

// called when the top level function returns (i.e. control is passed outside the Runtime) but it was due to an interrupt
func (r *Runtime) leaveAbrupt() ***REMOVED***
	r.jobQueue = nil
	r.ClearInterrupt()
***REMOVED***

func nilSafe(v Value) Value ***REMOVED***
	if v != nil ***REMOVED***
		return v
	***REMOVED***
	return _undefined
***REMOVED***

func isArray(object *Object) bool ***REMOVED***
	self := object.self
	if proxy, ok := self.(*proxyObject); ok ***REMOVED***
		if proxy.target == nil ***REMOVED***
			panic(typeError("Cannot perform 'IsArray' on a proxy that has been revoked"))
		***REMOVED***
		return isArray(proxy.target)
	***REMOVED***
	switch self.className() ***REMOVED***
	case classArray:
		return true
	default:
		return false
	***REMOVED***
***REMOVED***

func isRegexp(v Value) bool ***REMOVED***
	if o, ok := v.(*Object); ok ***REMOVED***
		matcher := o.self.getSym(SymMatch, nil)
		if matcher != nil && matcher != _undefined ***REMOVED***
			return matcher.ToBoolean()
		***REMOVED***
		_, reg := o.self.(*regexpObject)
		return reg
	***REMOVED***
	return false
***REMOVED***

func limitCallArgs(call FunctionCall, n int) FunctionCall ***REMOVED***
	if len(call.Arguments) > n ***REMOVED***
		return FunctionCall***REMOVED***This: call.This, Arguments: call.Arguments[:n]***REMOVED***
	***REMOVED*** else ***REMOVED***
		return call
	***REMOVED***
***REMOVED***

func shrinkCap(newSize, oldCap int) int ***REMOVED***
	if oldCap > 8 ***REMOVED***
		if cap := oldCap / 2; cap >= newSize ***REMOVED***
			return cap
		***REMOVED***
	***REMOVED***
	return oldCap
***REMOVED***

func growCap(newSize, oldSize, oldCap int) int ***REMOVED***
	// Use the same algorithm as in runtime.growSlice
	doublecap := oldCap + oldCap
	if newSize > doublecap ***REMOVED***
		return newSize
	***REMOVED*** else ***REMOVED***
		if oldSize < 1024 ***REMOVED***
			return doublecap
		***REMOVED*** else ***REMOVED***
			cap := oldCap
			// Check 0 < cap to detect overflow
			// and prevent an infinite loop.
			for 0 < cap && cap < newSize ***REMOVED***
				cap += cap / 4
			***REMOVED***
			// Return the requested cap when
			// the calculation overflowed.
			if cap <= 0 ***REMOVED***
				return newSize
			***REMOVED***
			return cap
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) genId() (ret uint64) ***REMOVED***
	if r.hash == nil ***REMOVED***
		h := r.getHash()
		r.idSeq = h.Sum64()
	***REMOVED***
	if r.idSeq == 0 ***REMOVED***
		r.idSeq = 1
	***REMOVED***
	ret = r.idSeq
	r.idSeq++
	return
***REMOVED***

func (r *Runtime) setGlobal(name unistring.String, v Value, strict bool) ***REMOVED***
	if ref := r.global.stash.getRefByName(name, strict); ref != nil ***REMOVED***
		ref.set(v)
	***REMOVED*** else ***REMOVED***
		o := r.globalObject.self
		if strict ***REMOVED***
			if o.hasOwnPropertyStr(name) ***REMOVED***
				o.setOwnStr(name, v, true)
			***REMOVED*** else ***REMOVED***
				r.throwReferenceError(name)
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			o.setOwnStr(name, v, false)
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) trackPromiseRejection(p *Promise, operation PromiseRejectionOperation) ***REMOVED***
	if r.promiseRejectionTracker != nil ***REMOVED***
		r.promiseRejectionTracker(p, operation)
	***REMOVED***
***REMOVED***

func (r *Runtime) callJobCallback(job *jobCallback, this Value, args ...Value) Value ***REMOVED***
	return job.callback(FunctionCall***REMOVED***This: this, Arguments: args***REMOVED***)
***REMOVED***

func (r *Runtime) invoke(v Value, p unistring.String, args ...Value) Value ***REMOVED***
	o := v.ToObject(r)
	return r.toCallable(o.self.getStr(p, nil))(FunctionCall***REMOVED***This: v, Arguments: args***REMOVED***)
***REMOVED***

func (r *Runtime) iterableToList(iterable Value, method func(FunctionCall) Value) []Value ***REMOVED***
	iter := r.getIterator(iterable, method)
	var values []Value
	iter.iterate(func(item Value) ***REMOVED***
		values = append(values, item)
	***REMOVED***)
	return values
***REMOVED***

func (r *Runtime) putSpeciesReturnThis(o objectImpl) ***REMOVED***
	o._putSym(SymSpecies, &valueProperty***REMOVED***
		getterFunc:   r.newNativeFunc(r.returnThis, nil, "get [Symbol.species]", nil, 0),
		accessor:     true,
		configurable: true,
	***REMOVED***)
***REMOVED***

func strToArrayIdx(s unistring.String) uint32 ***REMOVED***
	if s == "" ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	l := len(s)
	if s[0] == '0' ***REMOVED***
		if l == 1 ***REMOVED***
			return 0
		***REMOVED***
		return math.MaxUint32
	***REMOVED***
	var n uint32
	if l < 10 ***REMOVED***
		// guaranteed not to overflow
		for i := 0; i < len(s); i++ ***REMOVED***
			c := s[i]
			if c < '0' || c > '9' ***REMOVED***
				return math.MaxUint32
			***REMOVED***
			n = n*10 + uint32(c-'0')
		***REMOVED***
		return n
	***REMOVED***
	if l > 10 ***REMOVED***
		// guaranteed to overflow
		return math.MaxUint32
	***REMOVED***
	c9 := s[9]
	if c9 < '0' || c9 > '9' ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	for i := 0; i < 9; i++ ***REMOVED***
		c := s[i]
		if c < '0' || c > '9' ***REMOVED***
			return math.MaxUint32
		***REMOVED***
		n = n*10 + uint32(c-'0')
	***REMOVED***
	if n >= math.MaxUint32/10+1 ***REMOVED***
		return math.MaxUint32
	***REMOVED***
	n *= 10
	n1 := n + uint32(c9-'0')
	if n1 < n ***REMOVED***
		return math.MaxUint32
	***REMOVED***

	return n1
***REMOVED***

func strToInt32(s unistring.String) (int32, bool) ***REMOVED***
	if s == "" ***REMOVED***
		return -1, false
	***REMOVED***
	neg := s[0] == '-'
	if neg ***REMOVED***
		s = s[1:]
	***REMOVED***
	l := len(s)
	if s[0] == '0' ***REMOVED***
		if l == 1 ***REMOVED***
			return 0, !neg
		***REMOVED***
		return -1, false
	***REMOVED***
	var n uint32
	if l < 10 ***REMOVED***
		// guaranteed not to overflow
		for i := 0; i < len(s); i++ ***REMOVED***
			c := s[i]
			if c < '0' || c > '9' ***REMOVED***
				return -1, false
			***REMOVED***
			n = n*10 + uint32(c-'0')
		***REMOVED***
	***REMOVED*** else if l > 10 ***REMOVED***
		// guaranteed to overflow
		return -1, false
	***REMOVED*** else ***REMOVED***
		c9 := s[9]
		if c9 >= '0' ***REMOVED***
			if !neg && c9 > '7' || c9 > '8' ***REMOVED***
				// guaranteed to overflow
				return -1, false
			***REMOVED***
			for i := 0; i < 9; i++ ***REMOVED***
				c := s[i]
				if c < '0' || c > '9' ***REMOVED***
					return -1, false
				***REMOVED***
				n = n*10 + uint32(c-'0')
			***REMOVED***
			if n >= math.MaxInt32/10+1 ***REMOVED***
				// valid number, but it overflows integer
				return 0, false
			***REMOVED***
			n = n*10 + uint32(c9-'0')
		***REMOVED*** else ***REMOVED***
			return -1, false
		***REMOVED***
	***REMOVED***
	if neg ***REMOVED***
		return int32(-n), true
	***REMOVED***
	return int32(n), true
***REMOVED***

func strToInt64(s unistring.String) (int64, bool) ***REMOVED***
	if s == "" ***REMOVED***
		return -1, false
	***REMOVED***
	neg := s[0] == '-'
	if neg ***REMOVED***
		s = s[1:]
	***REMOVED***
	l := len(s)
	if s[0] == '0' ***REMOVED***
		if l == 1 ***REMOVED***
			return 0, !neg
		***REMOVED***
		return -1, false
	***REMOVED***
	var n uint64
	if l < 19 ***REMOVED***
		// guaranteed not to overflow
		for i := 0; i < len(s); i++ ***REMOVED***
			c := s[i]
			if c < '0' || c > '9' ***REMOVED***
				return -1, false
			***REMOVED***
			n = n*10 + uint64(c-'0')
		***REMOVED***
	***REMOVED*** else if l > 19 ***REMOVED***
		// guaranteed to overflow
		return -1, false
	***REMOVED*** else ***REMOVED***
		c18 := s[18]
		if c18 >= '0' ***REMOVED***
			if !neg && c18 > '7' || c18 > '8' ***REMOVED***
				// guaranteed to overflow
				return -1, false
			***REMOVED***
			for i := 0; i < 18; i++ ***REMOVED***
				c := s[i]
				if c < '0' || c > '9' ***REMOVED***
					return -1, false
				***REMOVED***
				n = n*10 + uint64(c-'0')
			***REMOVED***
			if n >= math.MaxInt64/10+1 ***REMOVED***
				// valid number, but it overflows integer
				return 0, false
			***REMOVED***
			n = n*10 + uint64(c18-'0')
		***REMOVED*** else ***REMOVED***
			return -1, false
		***REMOVED***
	***REMOVED***
	if neg ***REMOVED***
		return int64(-n), true
	***REMOVED***
	return int64(n), true
***REMOVED***

func strToInt(s unistring.String) (int, bool) ***REMOVED***
	if bits.UintSize == 32 ***REMOVED***
		n, ok := strToInt32(s)
		return int(n), ok
	***REMOVED***
	n, ok := strToInt64(s)
	return int(n), ok
***REMOVED***

// Attempts to convert a string into a canonical integer.
// On success returns (number, true).
// If it was a canonical number, but not an integer returns (0, false). This includes -0 and overflows.
// In all other cases returns (-1, false).
// See https://262.ecma-international.org/#sec-canonicalnumericindexstring
func strToIntNum(s unistring.String) (int, bool) ***REMOVED***
	n, ok := strToInt64(s)
	if n == 0 ***REMOVED***
		return 0, ok
	***REMOVED***
	if ok && n >= -maxInt && n <= maxInt ***REMOVED***
		if bits.UintSize == 32 ***REMOVED***
			if n > math.MaxInt32 || n < math.MinInt32 ***REMOVED***
				return 0, false
			***REMOVED***
		***REMOVED***
		return int(n), true
	***REMOVED***
	str := stringValueFromRaw(s)
	if str.ToNumber().toString().SameAs(str) ***REMOVED***
		return 0, false
	***REMOVED***
	return -1, false
***REMOVED***

func strToGoIdx(s unistring.String) int ***REMOVED***
	if n, ok := strToInt(s); ok ***REMOVED***
		return n
	***REMOVED***
	return -1
***REMOVED***

func strToIdx64(s unistring.String) int64 ***REMOVED***
	if n, ok := strToInt64(s); ok ***REMOVED***
		return n
	***REMOVED***
	return -1
***REMOVED***

func assertCallable(v Value) (func(FunctionCall) Value, bool) ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		return obj.self.assertCallable()
	***REMOVED***
	return nil, false
***REMOVED***
