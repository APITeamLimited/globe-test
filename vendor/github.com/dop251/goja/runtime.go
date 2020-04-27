package goja

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/text/collate"

	js_ast "github.com/dop251/goja/ast"
	"github.com/dop251/goja/parser"
)

const (
	sqrt1_2 float64 = math.Sqrt2 / 2
)

var (
	typeCallable = reflect.TypeOf(Callable(nil))
	typeValue    = reflect.TypeOf((*Value)(nil)).Elem()
	typeTime     = reflect.TypeOf(time.Time***REMOVED******REMOVED***)
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

	ArrayBuffer *Object

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

	ArrayBufferPrototype *Object

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

	typeInfoCache   map[reflect.Type]*reflectTypeInfo
	fieldNameMapper FieldNameMapper

	vm *vm
***REMOVED***

type stackFrame struct ***REMOVED***
	prg      *Program
	funcName string
	pc       int
***REMOVED***

func (f *stackFrame) position() Position ***REMOVED***
	return f.prg.src.Position(f.prg.sourceOffset(f.pc))
***REMOVED***

func (f *stackFrame) write(b *bytes.Buffer) ***REMOVED***
	if f.prg != nil ***REMOVED***
		if n := f.prg.funcName; n != "" ***REMOVED***
			b.WriteString(n)
			b.WriteString(" (")
		***REMOVED***
		if n := f.prg.src.name; n != "" ***REMOVED***
			b.WriteString(n)
		***REMOVED*** else ***REMOVED***
			b.WriteString("<eval>")
		***REMOVED***
		b.WriteByte(':')
		b.WriteString(f.position().String())
		b.WriteByte('(')
		b.WriteString(strconv.Itoa(f.pc))
		b.WriteByte(')')
		if f.prg.funcName != "" ***REMOVED***
			b.WriteByte(')')
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		if f.funcName != "" ***REMOVED***
			b.WriteString(f.funcName)
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
	stack []stackFrame
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
		frame.write(b)
		b.WriteByte('\n')
	***REMOVED***
***REMOVED***

func (e *Exception) writeShortStack(b *bytes.Buffer) ***REMOVED***
	if len(e.stack) > 0 && (e.stack[0].prg != nil || e.stack[0].funcName != "") ***REMOVED***
		b.WriteString(" at ")
		e.stack[0].write(b)
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
	r.globalObject.self._putProp(name, value, true, false, true)
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

	r.global.FunctionPrototype = r.newNativeFunc(nil, nil, "Empty", nil, 0)
	r.initObject()
	r.initFunction()
	r.initArray()
	r.initString()
	r.initNumber()
	r.initRegExp()
	r.initDate()
	r.initBoolean()

	r.initErrors()

	r.global.Eval = r.newNativeFunc(r.builtin_eval, nil, "eval", nil, 1)
	r.addToGlobal("eval", r.global.Eval)

	r.initGlobalObject()

	r.initMath()
	r.initJSON()

	//r.initTypedArrays()

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

func (r *Runtime) throwReferenceError(name string) ***REMOVED***
	panic(r.newError(r.global.ReferenceError, "%s is not defined", name))
***REMOVED***

func (r *Runtime) newSyntaxError(msg string, offset int) Value ***REMOVED***
	return r.builtin_new((r.global.SyntaxError), []Value***REMOVED***newStringValue(msg)***REMOVED***)
***REMOVED***

func (r *Runtime) newArray(prototype *Object) (a *arrayObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	a = &arrayObject***REMOVED******REMOVED***
	a.class = classArray
	a.val = v
	a.extensible = true
	v.self = a
	a.prototype = prototype
	a.init()
	return
***REMOVED***

func (r *Runtime) newArrayObject() *arrayObject ***REMOVED***
	return r.newArray(r.global.ArrayPrototype)
***REMOVED***

func (r *Runtime) newArrayValues(values []Value) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	a := &arrayObject***REMOVED******REMOVED***
	a.class = classArray
	a.val = v
	a.extensible = true
	v.self = a
	a.prototype = r.global.ArrayPrototype
	a.init()
	a.values = values
	a.length = int64(len(values))
	a.objCount = a.length
	return v
***REMOVED***

func (r *Runtime) newArrayLength(l int64) *Object ***REMOVED***
	a := r.newArrayValues(nil)
	a.self.putStr("length", intToValue(l), true)
	return a
***REMOVED***

func (r *Runtime) newBaseObject(proto *Object, class string) (o *baseObject) ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o = &baseObject***REMOVED******REMOVED***
	o.class = class
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = proto
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

func (r *Runtime) newFunc(name string, len int, strict bool) (f *funcObject) ***REMOVED***
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

func (r *Runtime) newNativeFuncObj(v *Object, call func(FunctionCall) Value, construct func(args []Value) *Object, name string, proto *Object, length int) *nativeFuncObject ***REMOVED***
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
		construct: construct,
	***REMOVED***
	v.self = f
	f.init(name, length)
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
	***REMOVED***
	return f
***REMOVED***

func (r *Runtime) newNativeConstructor(call func(ConstructorCall) *Object, name string, length int) *Object ***REMOVED***
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

	f.construct = func(args []Value) *Object ***REMOVED***
		return f.defaultConstruct(call, args)
	***REMOVED***

	v.self = f
	f.init(name, length)

	proto := r.NewObject()
	proto.self._putProp("constructor", v, true, false, true)
	f._putProp("prototype", proto, true, false, false)

	return v
***REMOVED***

func (r *Runtime) newNativeFunc(call func(FunctionCall) Value, construct func(args []Value) *Object, name string, proto *Object, length int) *Object ***REMOVED***
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
		construct: construct,
	***REMOVED***
	v.self = f
	f.init(name, length)
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
		proto.self._putProp("constructor", v, true, false, true)
	***REMOVED***
	return v
***REMOVED***

func (r *Runtime) newNativeFuncConstructObj(v *Object, construct func(args []Value, proto *Object) *Object, name string, proto *Object, length int) *nativeFuncObject ***REMOVED***
	f := &nativeFuncObject***REMOVED***
		baseFuncObject: baseFuncObject***REMOVED***
			baseObject: baseObject***REMOVED***
				class:      classFunction,
				val:        v,
				extensible: true,
				prototype:  r.global.FunctionPrototype,
			***REMOVED***,
		***REMOVED***,
		f: r.constructWrap(construct, proto),
		construct: func(args []Value) *Object ***REMOVED***
			return construct(args, proto)
		***REMOVED***,
	***REMOVED***

	f.init(name, length)
	if proto != nil ***REMOVED***
		f._putProp("prototype", proto, false, false, false)
	***REMOVED***
	return f
***REMOVED***

func (r *Runtime) newNativeFuncConstruct(construct func(args []Value, proto *Object) *Object, name string, prototype *Object, length int) *Object ***REMOVED***
	return r.newNativeFuncConstructProto(construct, name, prototype, r.global.FunctionPrototype, length)
***REMOVED***

func (r *Runtime) newNativeFuncConstructProto(construct func(args []Value, proto *Object) *Object, name string, prototype, proto *Object, length int) *Object ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	f := &nativeFuncObject***REMOVED******REMOVED***
	f.class = classFunction
	f.val = v
	f.extensible = true
	v.self = f
	f.prototype = proto
	f.f = r.constructWrap(construct, prototype)
	f.construct = func(args []Value) *Object ***REMOVED***
		return construct(args, prototype)
	***REMOVED***
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
		return intToValue(0)
	***REMOVED***
***REMOVED***

func (r *Runtime) builtin_newNumber(args []Value) *Object ***REMOVED***
	var v Value
	if len(args) > 0 ***REMOVED***
		v = args[0].ToNumber()
	***REMOVED*** else ***REMOVED***
		v = intToValue(0)
	***REMOVED***
	return r.newPrimitiveObject(v, r.global.NumberPrototype, classNumber)
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

func (r *Runtime) builtin_newBoolean(args []Value) *Object ***REMOVED***
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
	return r.newPrimitiveObject(v, r.global.BooleanPrototype, classBoolean)
***REMOVED***

func (r *Runtime) error_toString(call FunctionCall) Value ***REMOVED***
	obj := call.This.ToObject(r).self
	msg := obj.getStr("message")
	name := obj.getStr("name")
	var nameStr, msgStr string
	if name != nil && name != _undefined ***REMOVED***
		nameStr = name.String()
	***REMOVED***
	if msg != nil && msg != _undefined ***REMOVED***
		msgStr = msg.String()
	***REMOVED***
	if nameStr != "" && msgStr != "" ***REMOVED***
		return newStringValue(fmt.Sprintf("%s: %s", name.String(), msgStr))
	***REMOVED*** else ***REMOVED***
		if nameStr != "" ***REMOVED***
			return name.ToString()
		***REMOVED*** else ***REMOVED***
			return msg.ToString()
		***REMOVED***
	***REMOVED***
***REMOVED***

func (r *Runtime) builtin_Error(args []Value, proto *Object) *Object ***REMOVED***
	obj := r.newBaseObject(proto, classError)
	if len(args) > 0 && args[0] != _undefined ***REMOVED***
		obj._putProp("message", args[0], true, false, true)
	***REMOVED***
	return obj.val
***REMOVED***

func (r *Runtime) builtin_new(construct *Object, args []Value) *Object ***REMOVED***
repeat:
	switch f := construct.self.(type) ***REMOVED***
	case *nativeFuncObject:
		if f.construct != nil ***REMOVED***
			return f.construct(args)
		***REMOVED*** else ***REMOVED***
			panic("Not a constructor")
		***REMOVED***
	case *boundFuncObject:
		if f.construct != nil ***REMOVED***
			return f.construct(args)
		***REMOVED*** else ***REMOVED***
			panic("Not a constructor")
		***REMOVED***
	case *funcObject:
		// TODO: implement
		panic("Not implemented")
	case *lazyObject:
		construct.self = f.create(construct)
		goto repeat
	default:
		panic("Not a constructor")
	***REMOVED***
***REMOVED***

func (r *Runtime) throw(e Value) ***REMOVED***
	panic(e)
***REMOVED***

func (r *Runtime) builtin_thrower(call FunctionCall) Value ***REMOVED***
	r.typeErrorResult(true, "'caller', 'callee', and 'arguments' properties may not be accessed on strict mode functions or the arguments objects for calls to them")
	return nil
***REMOVED***

func (r *Runtime) eval(src string, direct, strict bool, this Value) Value ***REMOVED***

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
	if str, ok := call.Arguments[0].assertString(); ok ***REMOVED***
		return r.eval(str.String(), false, false, r.globalObject)
	***REMOVED***
	return call.Arguments[0]
***REMOVED***

func (r *Runtime) constructWrap(construct func(args []Value, proto *Object) *Object, proto *Object) func(call FunctionCall) Value ***REMOVED***
	return func(call FunctionCall) Value ***REMOVED***
		return construct(call.Arguments, proto)
	***REMOVED***
***REMOVED***

func (r *Runtime) toCallable(v Value) func(FunctionCall) Value ***REMOVED***
	if call, ok := r.toObject(v).self.assertCallable(); ok ***REMOVED***
		return call
	***REMOVED***
	r.typeErrorResult(true, "Value is not callable: %s", v.ToString())
	return nil
***REMOVED***

func (r *Runtime) checkObjectCoercible(v Value) ***REMOVED***
	switch v.(type) ***REMOVED***
	case valueUndefined, valueNull:
		r.typeErrorResult(true, "Value is not object coercible")
	***REMOVED***
***REMOVED***

func toUInt32(v Value) uint32 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.assertInt(); ok ***REMOVED***
		return uint32(i)
	***REMOVED***

	if f, ok := v.assertFloat(); ok ***REMOVED***
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint32(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
***REMOVED***

func toUInt16(v Value) uint16 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.assertInt(); ok ***REMOVED***
		return uint16(i)
	***REMOVED***

	if f, ok := v.assertFloat(); ok ***REMOVED***
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return uint16(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
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

func toInt32(v Value) int32 ***REMOVED***
	v = v.ToNumber()
	if i, ok := v.assertInt(); ok ***REMOVED***
		return int32(i)
	***REMOVED***

	if f, ok := v.assertFloat(); ok ***REMOVED***
		if !math.IsNaN(f) && !math.IsInf(f, 0) ***REMOVED***
			return int32(int64(f))
		***REMOVED***
	***REMOVED***
	return 0
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
	***REMOVED***
	return
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
ToValue converts a Go value into JavaScript value.

Primitive types (ints and uints, floats, string, bool) are converted to the corresponding JavaScript primitives.

func(FunctionCall) Value is treated as a native JavaScript function.

map[string]interface***REMOVED******REMOVED*** is converted into a host object that largely behaves like a JavaScript Object.

[]interface***REMOVED******REMOVED*** is converted into a host object that behaves largely like a JavaScript Array, however it's not extensible
because extending it can change the pointer so it becomes detached from the original.

*[]interface***REMOVED******REMOVED*** same as above, but the array becomes extensible.

A function is wrapped within a native JavaScript function. When called the arguments are automatically converted to
the appropriate Go types. If conversion is not possible, a TypeError is thrown.

A slice type is converted into a generic reflect based host object that behaves similar to an unexpandable Array.

Any other type is converted to a generic reflect based host object. Depending on the underlying type it behaves similar
to a Number, String, Boolean or Object.

Note that the underlying type is not lost, calling Export() returns the original Go value. This applies to all
reflect based types.
*/
func (r *Runtime) ToValue(i interface***REMOVED******REMOVED***) Value ***REMOVED***
	switch i := i.(type) ***REMOVED***
	case nil:
		return _null
	case Value:
		// TODO: prevent importing Objects from a different runtime
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
		return r.newNativeFunc(i, nil, "", nil, 0)
	case func(ConstructorCall) *Object:
		return r.newNativeConstructor(i, "", 0)
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
		return r.newNativeFunc(r.wrapReflectFunc(value), nil, "", nil, value.Type().NumIn())
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
				if v, err := r.toReflectValue(a, typ.In(n)); err == nil ***REMOVED***
					in[i] = v
					callSlice = true
					break
				***REMOVED***
			***REMOVED***
			var err error
			in[i], err = r.toReflectValue(a, t)
			if err != nil ***REMOVED***
				panic(r.newError(r.global.TypeError, "Could not convert function call parameter %v to %v", a, t))
			***REMOVED***
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

func (r *Runtime) toReflectValue(v Value, typ reflect.Type) (reflect.Value, error) ***REMOVED***
	switch typ.Kind() ***REMOVED***
	case reflect.String:
		return reflect.ValueOf(v.String()).Convert(typ), nil
	case reflect.Bool:
		return reflect.ValueOf(v.ToBoolean()).Convert(typ), nil
	case reflect.Int:
		i, _ := toInt(v)
		return reflect.ValueOf(int(i)).Convert(typ), nil
	case reflect.Int64:
		i, _ := toInt(v)
		return reflect.ValueOf(i).Convert(typ), nil
	case reflect.Int32:
		i, _ := toInt(v)
		return reflect.ValueOf(int32(i)).Convert(typ), nil
	case reflect.Int16:
		i, _ := toInt(v)
		return reflect.ValueOf(int16(i)).Convert(typ), nil
	case reflect.Int8:
		i, _ := toInt(v)
		return reflect.ValueOf(int8(i)).Convert(typ), nil
	case reflect.Uint:
		i, _ := toInt(v)
		return reflect.ValueOf(uint(i)).Convert(typ), nil
	case reflect.Uint64:
		i, _ := toInt(v)
		return reflect.ValueOf(uint64(i)).Convert(typ), nil
	case reflect.Uint32:
		i, _ := toInt(v)
		return reflect.ValueOf(uint32(i)).Convert(typ), nil
	case reflect.Uint16:
		i, _ := toInt(v)
		return reflect.ValueOf(uint16(i)).Convert(typ), nil
	case reflect.Uint8:
		i, _ := toInt(v)
		return reflect.ValueOf(uint8(i)).Convert(typ), nil
	***REMOVED***

	if typ == typeCallable ***REMOVED***
		if fn, ok := AssertFunction(v); ok ***REMOVED***
			return reflect.ValueOf(fn), nil
		***REMOVED***
	***REMOVED***

	if typ.Implements(typeValue) ***REMOVED***
		return reflect.ValueOf(v), nil
	***REMOVED***

	et := v.ExportType()
	if et == nil ***REMOVED***
		return reflect.Zero(typ), nil
	***REMOVED***
	if et.AssignableTo(typ) ***REMOVED***
		return reflect.ValueOf(v.Export()), nil
	***REMOVED*** else if et.ConvertibleTo(typ) ***REMOVED***
		return reflect.ValueOf(v.Export()).Convert(typ), nil
	***REMOVED***

	if typ == typeTime && et.Kind() == reflect.String ***REMOVED***
		time, ok := dateParse(v.String())
		if !ok ***REMOVED***
			return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert string %v to %v", v, typ)
		***REMOVED***
		return reflect.ValueOf(time), nil
	***REMOVED***

	switch typ.Kind() ***REMOVED***
	case reflect.Slice:
		if o, ok := v.(*Object); ok ***REMOVED***
			if o.self.className() == classArray ***REMOVED***
				l := int(toLength(o.self.getStr("length")))
				s := reflect.MakeSlice(typ, l, l)
				elemTyp := typ.Elem()
				for i := 0; i < l; i++ ***REMOVED***
					item := o.self.get(intToValue(int64(i)))
					itemval, err := r.toReflectValue(item, elemTyp)
					if err != nil ***REMOVED***
						return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert array element %v to %v at %d: %s", v, typ, i, err)
					***REMOVED***
					s.Index(i).Set(itemval)
				***REMOVED***
				return s, nil
			***REMOVED***
		***REMOVED***
	case reflect.Map:
		if o, ok := v.(*Object); ok ***REMOVED***
			m := reflect.MakeMap(typ)
			keyTyp := typ.Key()
			elemTyp := typ.Elem()
			needConvertKeys := !reflect.ValueOf("").Type().AssignableTo(keyTyp)
			for item, f := o.self.enumerate(false, false)(); f != nil; item, f = f() ***REMOVED***
				var kv reflect.Value
				var err error
				if needConvertKeys ***REMOVED***
					kv, err = r.toReflectValue(newStringValue(item.name), keyTyp)
					if err != nil ***REMOVED***
						return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert map key %s to %v", item.name, typ)
					***REMOVED***
				***REMOVED*** else ***REMOVED***
					kv = reflect.ValueOf(item.name)
				***REMOVED***

				ival := item.value
				if ival == nil ***REMOVED***
					ival = o.self.getStr(item.name)
				***REMOVED***
				if ival != nil ***REMOVED***
					vv, err := r.toReflectValue(ival, elemTyp)
					if err != nil ***REMOVED***
						return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert map value %v to %v at key %s", ival, typ, item.name)
					***REMOVED***
					m.SetMapIndex(kv, vv)
				***REMOVED*** else ***REMOVED***
					m.SetMapIndex(kv, reflect.Zero(elemTyp))
				***REMOVED***
			***REMOVED***
			return m, nil
		***REMOVED***
	case reflect.Struct:
		if o, ok := v.(*Object); ok ***REMOVED***
			s := reflect.New(typ).Elem()
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
						v = o.self.getStr(name)
					***REMOVED***

					if v != nil ***REMOVED***
						vv, err := r.toReflectValue(v, field.Type)
						if err != nil ***REMOVED***
							return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert struct value %v to %v for field %s: %s", v, field.Type, field.Name, err)

						***REMOVED***
						s.Field(i).Set(vv)
					***REMOVED***
				***REMOVED***
			***REMOVED***
			return s, nil
		***REMOVED***
	case reflect.Func:
		if fn, ok := AssertFunction(v); ok ***REMOVED***
			return reflect.MakeFunc(typ, r.wrapJSFunc(fn, typ)), nil
		***REMOVED***
	case reflect.Ptr:
		elemTyp := typ.Elem()
		v, err := r.toReflectValue(v, elemTyp)
		if err != nil ***REMOVED***
			return reflect.Value***REMOVED******REMOVED***, err
		***REMOVED***

		ptrVal := reflect.New(v.Type())
		ptrVal.Elem().Set(v)

		return ptrVal, nil
	***REMOVED***

	return reflect.Value***REMOVED******REMOVED***, fmt.Errorf("Could not convert %v to %v", v, typ)
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
				results[0], err = r.toReflectValue(res, typ.Out(0))
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
func (r *Runtime) ExportTo(v Value, target interface***REMOVED******REMOVED***) error ***REMOVED***
	tval := reflect.ValueOf(target)
	if tval.Kind() != reflect.Ptr || tval.IsNil() ***REMOVED***
		return errors.New("target must be a non-nil pointer")
	***REMOVED***
	vv, err := r.toReflectValue(v, tval.Elem().Type())
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	tval.Elem().Set(vv)
	return nil
***REMOVED***

// GlobalObject returns the global object.
func (r *Runtime) GlobalObject() *Object ***REMOVED***
	return r.globalObject
***REMOVED***

// Set the specified value as a property of the global object.
// The value is first converted using ToValue()
func (r *Runtime) Set(name string, value interface***REMOVED******REMOVED***) ***REMOVED***
	r.globalObject.self.putStr(name, r.ToValue(value), false)
***REMOVED***

// Get the specified property of the global object.
func (r *Runtime) Get(name string) Value ***REMOVED***
	return r.globalObject.self.getStr(name)
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
				obj.runtime.vm.clearStack()
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
	f, ok := v.assertFloat()
	return ok && math.IsNaN(f)
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
