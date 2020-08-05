package goja

import "github.com/dop251/goja/unistring"

var (
	symHasInstance        = newSymbol(asciiString("Symbol.hasInstance"))
	symIsConcatSpreadable = newSymbol(asciiString("Symbol.isConcatSpreadable"))
	symIterator           = newSymbol(asciiString("Symbol.iterator"))
	symMatch              = newSymbol(asciiString("Symbol.match"))
	symReplace            = newSymbol(asciiString("Symbol.replace"))
	symSearch             = newSymbol(asciiString("Symbol.search"))
	symSpecies            = newSymbol(asciiString("Symbol.species"))
	symSplit              = newSymbol(asciiString("Symbol.split"))
	symToPrimitive        = newSymbol(asciiString("Symbol.toPrimitive"))
	symToStringTag        = newSymbol(asciiString("Symbol.toStringTag"))
	symUnscopables        = newSymbol(asciiString("Symbol.unscopables"))
)

func (r *Runtime) builtin_symbol(call FunctionCall) Value ***REMOVED***
	var desc valueString
	if arg := call.Argument(0); !IsUndefined(arg) ***REMOVED***
		desc = arg.toString()
	***REMOVED*** else ***REMOVED***
		desc = stringEmpty
	***REMOVED***
	return newSymbol(desc)
***REMOVED***

func (r *Runtime) symbolproto_tostring(call FunctionCall) Value ***REMOVED***
	sym, ok := call.This.(*valueSymbol)
	if !ok ***REMOVED***
		if obj, ok := call.This.(*Object); ok ***REMOVED***
			if v, ok := obj.self.(*primitiveValueObject); ok ***REMOVED***
				if sym1, ok := v.pValue.(*valueSymbol); ok ***REMOVED***
					sym = sym1
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	if sym == nil ***REMOVED***
		panic(r.NewTypeError("Method Symbol.prototype.toString is called on incompatible receiver"))
	***REMOVED***
	return sym.desc
***REMOVED***

func (r *Runtime) symbolproto_valueOf(call FunctionCall) Value ***REMOVED***
	_, ok := call.This.(*valueSymbol)
	if ok ***REMOVED***
		return call.This
	***REMOVED***

	if obj, ok := call.This.(*Object); ok ***REMOVED***
		if v, ok := obj.self.(*primitiveValueObject); ok ***REMOVED***
			if sym, ok := v.pValue.(*valueSymbol); ok ***REMOVED***
				return sym
			***REMOVED***
		***REMOVED***
	***REMOVED***

	panic(r.NewTypeError("Symbol.prototype.valueOf requires that 'this' be a Symbol"))
***REMOVED***

func (r *Runtime) symbol_for(call FunctionCall) Value ***REMOVED***
	key := call.Argument(0).toString()
	keyStr := key.string()
	if v := r.symbolRegistry[keyStr]; v != nil ***REMOVED***
		return v
	***REMOVED***
	if r.symbolRegistry == nil ***REMOVED***
		r.symbolRegistry = make(map[unistring.String]*valueSymbol)
	***REMOVED***
	v := newSymbol(key)
	r.symbolRegistry[keyStr] = v
	return v
***REMOVED***

func (r *Runtime) symbol_keyfor(call FunctionCall) Value ***REMOVED***
	arg := call.Argument(0)
	sym, ok := arg.(*valueSymbol)
	if !ok ***REMOVED***
		panic(r.NewTypeError("%s is not a symbol", arg.String()))
	***REMOVED***
	for key, s := range r.symbolRegistry ***REMOVED***
		if s == sym ***REMOVED***
			return stringValueFromRaw(key)
		***REMOVED***
	***REMOVED***
	return _undefined
***REMOVED***

func (r *Runtime) createSymbolProto(val *Object) objectImpl ***REMOVED***
	o := &baseObject***REMOVED***
		class:      classObject,
		val:        val,
		extensible: true,
		prototype:  r.global.ObjectPrototype,
	***REMOVED***
	o.init()

	o._putProp("constructor", r.global.Symbol, true, false, true)
	o._putProp("toString", r.newNativeFunc(r.symbolproto_tostring, nil, "toString", nil, 0), true, false, true)
	o._putProp("valueOf", r.newNativeFunc(r.symbolproto_valueOf, nil, "valueOf", nil, 0), true, false, true)
	o._putSym(symToPrimitive, valueProp(r.newNativeFunc(r.symbolproto_valueOf, nil, "[Symbol.toPrimitive]", nil, 1), false, false, true))
	o._putSym(symToStringTag, valueProp(newStringValue("Symbol"), false, false, true))

	return o
***REMOVED***

func (r *Runtime) createSymbol(val *Object) objectImpl ***REMOVED***
	o := r.newNativeFuncObj(val, r.builtin_symbol, nil, "Symbol", r.global.SymbolPrototype, 0)

	o._putProp("for", r.newNativeFunc(r.symbol_for, nil, "for", nil, 1), true, false, true)
	o._putProp("keyFor", r.newNativeFunc(r.symbol_keyfor, nil, "keyFor", nil, 1), true, false, true)

	for _, s := range []*valueSymbol***REMOVED***
		symHasInstance,
		symIsConcatSpreadable,
		symIterator,
		symMatch,
		symReplace,
		symSearch,
		symSpecies,
		symSplit,
		symToPrimitive,
		symToStringTag,
		symUnscopables,
	***REMOVED*** ***REMOVED***
		n := s.desc.(asciiString)
		n = n[len("Symbol(Symbol.") : len(n)-1]
		o._putProp(unistring.String(n), s, false, false, false)
	***REMOVED***

	return o
***REMOVED***

func (r *Runtime) initSymbol() ***REMOVED***
	r.global.SymbolPrototype = r.newLazyObject(r.createSymbolProto)

	r.global.Symbol = r.newLazyObject(r.createSymbol)
	r.addToGlobal("Symbol", r.global.Symbol)

***REMOVED***
