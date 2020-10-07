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
)

type iterationKind int

const (
	iterationKindKey iterationKind = iota
	iterationKindValue
	iterationKindKeyValue
)

type global struct ***REMOVED***
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

	IteratorPrototype       *Object
	ArrayIteratorPrototype  *Object
	MapIteratorPrototype    *Object
	SetIteratorPrototype    *Object
	StringIteratorPrototype *Object

	ErrorPrototype          *Object
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

	symbolRegistry map[unistring.String]*valueSymbol

	typeInfoCache   map[reflect.Type]*reflectTypeInfo
	fieldNameMapper FieldNameMapper

	vm    *vm
	hash  *maphash.Hash
	idSeq uint64

	// Contains a list of ids of finalized weak keys so that the runtime could pick it up and remove from
	// all weak collections using the weakKeys map. The runtime picks it up either when the topmost function
	// returns (i.e. the callstack becomes empty) or every 10000 'ticks' (vm instructions).
	// It is implemented this way to avoid circular references which at the time of writing (go 1.15) causes
	// the whole structure to become not garbage-collectable.
	weakRefTracker *weakRefTracker

	// Contains a list of weak collections that contain the key with the id.
	weakKeys map[uint64]*weakCollections
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
	return f.prg.src.name
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

func (f *StackFrame) Position() Position ***REMOVED***
	if f.prg == nil || f.prg.src == nil ***REMOVED***
		return Position***REMOVED***
			0,
			0,
		***REMOVED***
	***REMOVED***
	return f.prg.src.Position(f.prg.sourceOffset(f.pc))
***REMOVED***

func (f *StackFrame) Write(b *bytes.Buffer) ***REMOVED***
	if f.prg != nil ***REMOVED***
		if n := f.prg.funcName; n != "" ***REMOVED***
			b.WriteString(n.String())
			b.WriteString(" (")
		***REMOVED***
		if n := f.prg.src.name; n != "" ***REMOVED***
			b.WriteString(n)
		***REMOVED*** else ***REMOVED***
			b.WriteString("<eval>")
		***REMOVED***
		b.WriteByte(':')
		b.WriteString(f.Position().String())
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

type InterruptedError struct ***REMOVED***
	Exception
	iface interface***REMOVED******REMOVED***
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

	o._putSym(symIterator, valueProp(r.newNativeFunc(r.returnThis, nil, "[Symbol.iterator]", nil, 0), true, false, true))
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

	r.global.FunctionPrototype = r.newNativeFunc(func(FunctionCall) Value ***REMOVED***
		return _undefined
	***REMOVED***, nil, " ", nil, 0)

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

	r.global.thrower = r.newNativeFunc(r.builtin_thrower, nil, "thrower", nil, 0)
	r.global.throwerProperty = &valueProperty***REMOVED***
		getterFunc: r.global.thrower,
		setterFunc: r.global.thrower,
		accessor:   true,
	***REMOVED***
***REMOVED***

func (r *Runtime) typeErrorResult(throw bool, args ...interface***REMOVED******REMOVED***) ***REMOVED***
	if throw ***REMOVED***
		panic(r.NewTypeError(args...))
	***REMOVED***
***REMOVED***

func (r *Runtime) newError(typ *Object, format string, args ...interface***REMOVED******REMOVED***) Value ***REMOVED***
	msg := fmt.Sprintf(format, args...)
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

func (r *Runtime) newFunc(name unistring.String, len int, strict bool) (f *funcObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f = &funcObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	v.self = f
	f.prototype = r.global.FunctionPrototype
	f.init(name, len)
	if strict ***REMOVED***
		f._put("caller", r.global.throwerProperty)
		f._put("arguments", r.global.throwerProperty)
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) newNativeFuncObj(v *Object, call func(FunctionCall) Value, construct func(args []Value, proto *Object) *Object, name unistring.String, proto *Object, length int) *nativeFuncObject ***REMOVED***
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

func (r *Runtime) newNativeConstructor(call func(ConstructorCall) *Object, name unistring.String, length int) *Object ***REMOVED***
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
		return f.defaultConstruct(call, c.Arguments)
	***REMOVED***

	f.construct = func(args []Value, proto *Object) *Object ***REMOVED***
		return f.defaultConstruct(call, args)
	***REMOVED***

	v.self = f
	f.init(name, length)

	proto := r.NewObject()
	proto.self._putProp("constructor", v, true, false, true)
	f._putProp("prototype", proto, true, false, false)

	return v
***REMOVED***

func (r *Runtime) newNativeConstructOnly(v *Object, ctor func(args []Value, newTarget *Object) *Object, defaultProto *Object, name unistring.String, length int) *nativeFuncObject ***REMOVED***
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
		f: func(call FunctionCall) Value ***REMOVED***
			return ctor(call.Arguments, nil)
		***REMOVED***,
		construct: func(args []Value, newTarget *Object) *Object ***REMOVED***
			if newTarget == nil ***REMOVED***
				newTarget = v
			***REMOVED***
			return ctor(args, newTarget)
		***REMOVED***,
	***REMOVED***
	v.self = f
	f.init(name, length)
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
	f.init(name, length)
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

	f.init(name, length)
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
	***REMOVED***
	return f
***REMOVED***

func (r *Runtime) newNativeFuncConstruct(construct func(args []Value, proto *Object) *Object, name unistring.String, prototype *Object, length int) *Object ***REMOVED***
	return r.newNativeFuncConstructProto(construct, name, prototype, r.global.FunctionPrototype, length)
***REMOVED***

func (r *Runtime) newNativeFuncConstructProto(construct func(args []Value, proto *Object) *Object, name unistring.String, prototype, proto *Object, length int) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f := &nativeFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	v.self = f
	f.prototype = proto
	f.f = r.constructToCall(construct, prototype)
	f.construct = r.wrapNativeConstruct(construct, prototype)
	f.init(name, length)
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
	obj := call.This.ToObject(r).self
	name := obj.getStr("name", nil)
	if name == nil || name == _undefined ***REMOVED***
		nameStr = asciiString("Error")
	***REMOVED*** else ***REMOVED***
		nameStr = name.toString()
	***REMOVED***
	msg := obj.getStr("message", nil)
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

func (r *Runtime) builtin_Error(args []Value, proto *Object) *Object ***REMOVED***
	obj := r.newBaseObject(proto, classError)
	if len(args) > 0 && args[0] != _undefined ***REMOVED***
		obj._putProp("message", args[0], true, false, true)
	***REMOVED***
	return obj.val
***REMOVED***

func (r *Runtime) builtin_new(construct *Object, args []Value) *Object ***REMOVED***
	return r.toConstructor(construct)(args, nil)
***REMOVED***

func (r *Runtime) throw(e Value) ***REMOVED***
	panic(e)
***REMOVED***

func (r *Runtime) builtin_thrower(FunctionCall) Value ***REMOVED***
	r.typeErrorResult(true, "'caller', 'callee', and 'arguments' properties may not be accessed on strict mode functions or the arguments objects for calls to them")
	return nil
***REMOVED***

func (r *Runtime) eval(srcVal valueString, direct, strict bool, this Value) Value ***REMOVED***
	src := escapeInvalidUtf16(srcVal)
	p, err := r.compile("<eval>", src, strict, true)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***

	vm := r.vm

	vm.pushCtx()
	vm.prg = p
	vm.pc = 0
	if !direct ***REMOVED***
		vm.stash = nil
	***REMOVED***
	vm.sb = vm.sp
	vm.push(this)
	if strict ***REMOVED***
		vm.push(valueTrue)
	***REMOVED*** else ***REMOVED***
		vm.push(valueFalse)
	***REMOVED***
	vm.run()
	vm.popCtx()
	vm.halt = false
	retval := vm.stack[vm.sp-1]
	vm.sp -= 2
	return retval
***REMOVED***

func (r *Runtime) builtin_eval(call FunctionCall) Value ***REMOVED***
	if len(call.Arguments) == 0 ***REMOVED***
		return _undefined
	***REMOVED***
	if str, ok := call.Arguments[0].(valueString); ok ***REMOVED***
		return r.eval(str, false, false, r.globalObject)
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

func toIntStrict(i int64) int ***REMOVED***
	if bits.UintSize == 32 ***REMOVED***
		if i > math.MaxInt32 || i < math.MinInt32 ***REMOVED***
			panic(rangeError("Integer value overflows 32-bit int"))
		***REMOVED***
	***REMOVED***
	return int(i)
***REMOVED***

func (r *Runtime) toIndex(v Value) int ***REMOVED***
	intIdx := v.ToInteger()
	if intIdx >= 0 && intIdx < maxInt ***REMOVED***
		if bits.UintSize == 32 && intIdx >= math.MaxInt32 ***REMOVED***
			panic(r.newError(r.global.RangeError, "Index %s overflows int", v.String()))
		***REMOVED***
		return int(intIdx)
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
	return compile(name, src, strict, false)
***REMOVED***

// CompileAST creates an internal representation of the JavaScript code that can be later run using the Runtime.RunProgram()
// method. This representation is not linked to a runtime in any way and can be run in multiple runtimes (possibly
// at the same time).
func CompileAST(prg *js_ast.Program, strict bool) (*Program, error) ***REMOVED***
	return compileAST(prg, strict, false)
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

func compile(name, src string, strict, eval bool) (p *Program, err error) ***REMOVED***
	prg, err1 := parser.ParseFile(nil, name, src, 0)
	if err1 != nil ***REMOVED***
		switch err1 := err1.(type) ***REMOVED***
		case parser.ErrorList:
			if len(err1) > 0 && err1[0].Message == "Invalid left-hand side in assignment" ***REMOVED***
				err = &CompilerReferenceError***REMOVED***
					CompilerError: CompilerError***REMOVED***
						Message: err1.Error(),
					***REMOVED***,
				***REMOVED***
				return
			***REMOVED***
		***REMOVED***
		// FIXME offset
		err = &CompilerSyntaxError***REMOVED***
			CompilerError: CompilerError***REMOVED***
				Message: err1.Error(),
			***REMOVED***,
		***REMOVED***
		return
	***REMOVED***

	p, err = compileAST(prg, strict, eval)

	return
***REMOVED***

func compileAST(prg *js_ast.Program, strict, eval bool) (p *Program, err error) ***REMOVED***
	c := newCompiler()
	c.scope.strict = strict
	c.scope.eval = eval

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

	c.compile(prg)
	p = c.p
	return
***REMOVED***

func (r *Runtime) compile(name, src string, strict, eval bool) (p *Program, err error) ***REMOVED***
	p, err = compile(name, src, strict, eval)
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
	p, err := Compile(name, src, false)

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return r.RunProgram(p)
***REMOVED***

// RunProgram executes a pre-compiled (see Compile()) code in the global context.
func (r *Runtime) RunProgram(p *Program) (result Value, err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			if intr, ok := x.(*InterruptedError); ok ***REMOVED***
				err = intr
			***REMOVED*** else ***REMOVED***
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	recursive := false
	if len(r.vm.callStack) > 0 ***REMOVED***
		recursive = true
		r.vm.pushCtx()
	***REMOVED***
	r.vm.prg = p
	r.vm.pc = 0
	ex := r.vm.runTry()
	if ex == nil ***REMOVED***
		result = r.vm.pop()
	***REMOVED*** else ***REMOVED***
		err = ex
	***REMOVED***
	if recursive ***REMOVED***
		r.vm.popCtx()
		r.vm.halt = false
		r.vm.clearStack()
	***REMOVED*** else ***REMOVED***
		r.vm.stack = nil
		r.leave()
	***REMOVED***
	return
***REMOVED***

// CaptureCallStack appends the current call stack frames to the stack slice (which may be nil) up to the specified depth.
// The most recent frame will be the first one.
// If depth <= 0 or more than the number of available frames, returns the entire stack.
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
automatic argument and return value type conversions (which involves reflect).

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

Field properties are writable (if the struct is addressable) and non-configurable.
Method properties are non-writable and non-configurable.

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
prototype and all the usual methods should work. There are, however, some caveats:

- If the slice is not addressable, the array cannot be extended or shrunk. Any attempt to do so (by setting an index
beyond the current length or by modifying the length) will result in a TypeError.

- Converted Arrays may not contain holes (because Go slices cannot). This means that hasOwnProperty(n) will always
return `true` if n < length. Attempt to delete an item with an index < length will fail. Nil slice elements will be
converted to `null`. Accessing an element beyond `length` will return `undefined`.

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
		if i == nil || i.runtime == nil ***REMOVED***
			return _null
		***REMOVED***
		if i.runtime != r ***REMOVED***
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
	case func(ConstructorCall) *Object:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
		return r.newNativeConstructor(i, name, 0)
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
		obj := &Object***REMOVED***runtime: r***REMOVED***
		a := &objectGoSlice***REMOVED***
			baseObject: baseObject***REMOVED***
				val: obj,
			***REMOVED***,
			data: &i,
		***REMOVED***
		obj.self = a
		a.init()
		return obj
	case *[]interface***REMOVED******REMOVED***:
		if i == nil ***REMOVED***
			return _null
		***REMOVED***
		obj := &Object***REMOVED***runtime: r***REMOVED***
		a := &objectGoSlice***REMOVED***
			baseObject: baseObject***REMOVED***
				val: obj,
			***REMOVED***,
			data:            i,
			sliceExtensible: true,
		***REMOVED***
		obj.self = a
		a.init()
		return obj
	***REMOVED***

	origValue := reflect.ValueOf(i)
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
	case reflect.Slice:
		obj := &Object***REMOVED***runtime: r***REMOVED***
		a := &objectGoSliceReflect***REMOVED***
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
	case reflect.Func:
		name := unistring.NewFromString(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name())
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

		callSlice := false
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

			// if this is a variadic Go function, and the caller has supplied
			// exactly the number of JavaScript arguments required, and this
			// is the last JavaScript argument, try treating the it as the
			// actual set of variadic Go arguments. if that succeeds, break
			// out of the loop.
			if typ.IsVariadic() && len(call.Arguments) == nargs && i == nargs-1 ***REMOVED***
				v := reflect.New(typ.In(n)).Elem()
				if err := r.toReflectValue(a, v, &objectExportCtx***REMOVED******REMOVED***); err == nil ***REMOVED***
					in[i] = v
					callSlice = true
					break
				***REMOVED***
			***REMOVED***
			v := reflect.New(t).Elem()
			err := r.toReflectValue(a, v, &objectExportCtx***REMOVED******REMOVED***)
			if err != nil ***REMOVED***
				panic(r.newError(r.global.TypeError, "could not convert function call parameter %v to %v", a, t))
			***REMOVED***
			in[i] = v
		***REMOVED***

		var out []reflect.Value
		if callSlice ***REMOVED***
			out = value.CallSlice(in)
		***REMOVED*** else ***REMOVED***
			out = value.Call(in)
		***REMOVED***

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
	switch typ.Kind() ***REMOVED***
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
	***REMOVED***

	if typ == typeCallable ***REMOVED***
		if fn, ok := AssertFunction(v); ok ***REMOVED***
			dst.Set(reflect.ValueOf(fn))
			return nil
		***REMOVED***
	***REMOVED***

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

	***REMOVED***
		et := v.ExportType()
		if et == nil || et == reflectTypeNil ***REMOVED***
			dst.Set(reflect.Zero(typ))
			return nil
		***REMOVED***
		if et.AssignableTo(typ) ***REMOVED***
			dst.Set(reflect.ValueOf(exportValue(v, ctx)))
			return nil
		***REMOVED*** else if et.ConvertibleTo(typ) ***REMOVED***
			dst.Set(reflect.ValueOf(exportValue(v, ctx)).Convert(typ))
			return nil
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
	***REMOVED***

	switch typ.Kind() ***REMOVED***
	case reflect.Slice:
		if o, ok := v.(*Object); ok ***REMOVED***
			if o.self.className() == classArray ***REMOVED***
				if v, exists := ctx.getTyped(o.self, typ); exists ***REMOVED***
					dst.Set(reflect.ValueOf(v))
					return nil
				***REMOVED***
				l := int(toLength(o.self.getStr("length", nil)))
				if dst.IsNil() || dst.Len() != l ***REMOVED***
					dst.Set(reflect.MakeSlice(typ, l, l))
				***REMOVED***
				s := dst
				ctx.putTyped(o.self, typ, s.Interface())
				for i := 0; i < l; i++ ***REMOVED***
					item := o.self.getIdx(valueInt(int64(i)), nil)
					err := r.toReflectValue(item, s.Index(i), ctx)
					if err != nil ***REMOVED***
						return fmt.Errorf("could not convert array element %v to %v at %d: %w", v, typ, i, err)
					***REMOVED***
				***REMOVED***
				return nil
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		if o, ok := v.(*Object); ok ***REMOVED***
			if v, exists := ctx.getTyped(o.self, typ); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v))
				return nil
			***REMOVED***
			if dst.IsNil() ***REMOVED***
				dst.Set(reflect.MakeMap(typ))
			***REMOVED***
			m := dst
			ctx.putTyped(o.self, typ, m.Interface())
			keyTyp := typ.Key()
			elemTyp := typ.Elem()
			needConvertKeys := !reflect.ValueOf("").Type().AssignableTo(keyTyp)
			for _, itemName := range o.self.ownKeys(false, nil) ***REMOVED***
				var kv reflect.Value
				var err error
				if needConvertKeys ***REMOVED***
					kv = reflect.New(keyTyp).Elem()
					err = r.toReflectValue(itemName, kv, ctx)
					if err != nil ***REMOVED***
						return fmt.Errorf("could not convert map key %s to %v", itemName.String(), typ)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					kv = reflect.ValueOf(itemName.String())
				***REMOVED***

				ival := o.get(itemName, nil)
				if ival != nil ***REMOVED***
					vv := reflect.New(elemTyp).Elem()
					err := r.toReflectValue(ival, vv, ctx)
					if err != nil ***REMOVED***
						return fmt.Errorf("could not convert map value %v to %v at key %s", ival, typ, itemName.String())
					***REMOVED***
					m.SetMapIndex(kv, vv)
				***REMOVED*** else ***REMOVED***
					m.SetMapIndex(kv, reflect.Zero(elemTyp))
				***REMOVED***

			***REMOVED***
			return nil
		***REMOVED***
	case reflect.Struct:
		if o, ok := v.(*Object); ok ***REMOVED***
			t := reflect.PtrTo(typ)
			if v, exists := ctx.getTyped(o.self, t); exists ***REMOVED***
				dst.Set(reflect.ValueOf(v).Elem())
				return nil
			***REMOVED***
			s := dst
			ctx.putTyped(o.self, t, s.Addr().Interface())
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
			if v, exists := ctx.getTyped(o.self, typ); exists ***REMOVED***
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
// Exporting to an interface***REMOVED******REMOVED*** results in a value of the same type as Export() would produce.
// Exporting to numeric types uses the standard ECMAScript conversion operations, same as used when assigning
// values to non-clamped typed array items, e.g.
// https://www.ecma-international.org/ecma-262/10.0/index.html#sec-toint32
// Returns error if conversion is not possible.
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

// Set the specified value as a property of the global object.
// The value is first converted using ToValue()
func (r *Runtime) Set(name string, value interface***REMOVED******REMOVED***) ***REMOVED***
	r.globalObject.self.setOwnStr(unistring.NewFromString(name), r.ToValue(value), false)
***REMOVED***

// Get the specified property of the global object.
func (r *Runtime) Get(name string) Value ***REMOVED***
	return r.globalObject.self.getStr(unistring.NewFromString(name), nil)
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

// New is an equivalent of the 'new' operator allowing to call it directly from Go.
func (r *Runtime) New(construct Value, args ...Value) (o *Object, err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			switch x := x.(type) ***REMOVED***
			case *Exception:
				err = x
			case *InterruptedError:
				err = x
			default:
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()
	return r.builtin_new(r.toObject(construct), args), nil
***REMOVED***

// Callable represents a JavaScript function that can be called from Go.
type Callable func(this Value, args ...Value) (Value, error)

// AssertFunction checks if the Value is a function and returns a Callable.
func AssertFunction(v Value) (Callable, bool) ***REMOVED***
	if obj, ok := v.(*Object); ok ***REMOVED***
		if f, ok := obj.self.assertCallable(); ok ***REMOVED***
			return func(this Value, args ...Value) (ret Value, err error) ***REMOVED***
				defer func() ***REMOVED***
					if x := recover(); x != nil ***REMOVED***
						if ex, ok := x.(*InterruptedError); ok ***REMOVED***
							err = ex
						***REMOVED*** else ***REMOVED***
							panic(x)
						***REMOVED***
					***REMOVED***
				***REMOVED***()
				ex := obj.runtime.vm.try(func() ***REMOVED***
					ret = f(FunctionCall***REMOVED***
						This:      this,
						Arguments: args,
					***REMOVED***)
				***REMOVED***)
				if ex != nil ***REMOVED***
					err = ex
				***REMOVED***
				vm := obj.runtime.vm
				vm.clearStack()
				if len(vm.callStack) == 0 ***REMOVED***
					obj.runtime.leave()
				***REMOVED***
				return
			***REMOVED***, true
		***REMOVED***
	***REMOVED***
	return nil, false
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

func tryFunc(f func()) (err error) ***REMOVED***
	defer func() ***REMOVED***
		if x := recover(); x != nil ***REMOVED***
			switch x := x.(type) ***REMOVED***
			case *Exception:
				err = x
			case *InterruptedError:
				err = x
			case Value:
				err = &Exception***REMOVED***
					val: x,
				***REMOVED***
			default:
				panic(x)
			***REMOVED***
		***REMOVED***
	***REMOVED***()

	f()

	return nil
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
		c = r.toObject(c).self.getSym(symSpecies, nil)
	***REMOVED***
	if c == nil || c == _undefined || c == _null ***REMOVED***
		c = defaultConstructor
	***REMOVED***
	return r.toConstructor(c)
***REMOVED***

func (r *Runtime) speciesConstructorObj(o, defaultConstructor *Object) *Object ***REMOVED***
	c := o.self.getStr("constructor", nil)
	if c != nil && c != _undefined ***REMOVED***
		c = r.toObject(c).self.getSym(symSpecies, nil)
	***REMOVED***
	if c == nil || c == _undefined || c == _null ***REMOVED***
		return defaultConstructor
	***REMOVED***
	return r.toObject(c)
***REMOVED***

func (r *Runtime) returnThis(call FunctionCall) Value ***REMOVED***
	return call.This
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

func (r *Runtime) getIterator(obj Value, method func(FunctionCall) Value) *Object ***REMOVED***
	if method == nil ***REMOVED***
		method = toMethod(r.getV(obj, symIterator))
		if method == nil ***REMOVED***
			panic(r.NewTypeError("object is not iterable"))
		***REMOVED***
	***REMOVED***

	return r.toObject(method(FunctionCall***REMOVED***
		This: obj,
	***REMOVED***))
***REMOVED***

func returnIter(iter *Object) ***REMOVED***
	retMethod := toMethod(iter.self.getStr("return", nil))
	if retMethod != nil ***REMOVED***
		_ = tryFunc(func() ***REMOVED***
			retMethod(FunctionCall***REMOVED***This: iter***REMOVED***)
		***REMOVED***)
	***REMOVED***
***REMOVED***

func (r *Runtime) iterate(iter *Object, step func(Value)) ***REMOVED***
	for ***REMOVED***
		res := r.toObject(toMethod(iter.self.getStr("next", nil))(FunctionCall***REMOVED***This: iter***REMOVED***))
		if nilSafe(res.self.getStr("done", nil)).ToBoolean() ***REMOVED***
			break
		***REMOVED***
		err := tryFunc(func() ***REMOVED***
			step(nilSafe(res.self.getStr("value", nil)))
		***REMOVED***)
		if err != nil ***REMOVED***
			returnIter(iter)
			panic(err)
		***REMOVED***
	***REMOVED***
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

func (r *Runtime) addWeakKey(id uint64, coll weakCollection) ***REMOVED***
	keys := r.weakKeys
	if keys == nil ***REMOVED***
		keys = make(map[uint64]*weakCollections)
		r.weakKeys = keys
	***REMOVED***
	colls := keys[id]
	if colls == nil ***REMOVED***
		colls = &weakCollections***REMOVED***
			objId: id,
		***REMOVED***
		keys[id] = colls
	***REMOVED***
	colls.add(coll)
***REMOVED***

func (r *Runtime) removeWeakKey(id uint64, coll weakCollection) ***REMOVED***
	keys := r.weakKeys
	if colls := keys[id]; colls != nil ***REMOVED***
		colls.remove(coll)
		if len(colls.colls) == 0 ***REMOVED***
			delete(keys, id)
		***REMOVED***
	***REMOVED***
***REMOVED***

// this gets inlined so a CALL is avoided on a critical path
func (r *Runtime) removeDeadKeys() ***REMOVED***
	if r.weakRefTracker != nil ***REMOVED***
		r.doRemoveDeadKeys()
	***REMOVED***
***REMOVED***

func (r *Runtime) doRemoveDeadKeys() ***REMOVED***
	r.weakRefTracker.Lock()
	list := r.weakRefTracker.list
	r.weakRefTracker.list = nil
	r.weakRefTracker.Unlock()
	for _, id := range list ***REMOVED***
		if colls := r.weakKeys[id]; colls != nil ***REMOVED***
			for _, coll := range colls.colls ***REMOVED***
				coll.removeId(id)
			***REMOVED***
			delete(r.weakKeys, id)
		***REMOVED***
	***REMOVED***
***REMOVED***

// called when the top level function returns (i.e. control is passed outside the Runtime).
func (r *Runtime) leave() ***REMOVED***
	r.removeDeadKeys()
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
		matcher := o.self.getSym(symMatch, nil)
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
