package goja

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/dop251/goja/parser"
	"regexp"
)

func (r *Runtime) newRegexpObject(proto *Object) *regexpObject ***REMOVED***
	v := &Object***REMOVED***runtime: r***REMOVED***

	o := &regexpObject***REMOVED******REMOVED***
	o.class = classRegExp
	o.val = v
	o.extensible = true
	v.self = o
	o.prototype = proto
	o.init()
	return o
***REMOVED***

func (r *Runtime) newRegExpp(pattern regexpPattern, patternStr valueString, global, ignoreCase, multiline bool, proto *Object) *Object ***REMOVED***
	o := r.newRegexpObject(proto)

	o.pattern = pattern
	o.source = patternStr
	o.global = global
	o.ignoreCase = ignoreCase
	o.multiline = multiline

	return o.val
***REMOVED***

func compileRegexp(patternStr, flags string) (p regexpPattern, global, ignoreCase, multiline bool, err error) ***REMOVED***

	if flags != "" ***REMOVED***
		invalidFlags := func() ***REMOVED***
			err = fmt.Errorf("Invalid flags supplied to RegExp constructor '%s'", flags)
		***REMOVED***
		for _, chr := range flags ***REMOVED***
			switch chr ***REMOVED***
			case 'g':
				if global ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				global = true
			case 'm':
				if multiline ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				multiline = true
			case 'i':
				if ignoreCase ***REMOVED***
					invalidFlags()
					return
				***REMOVED***
				ignoreCase = true
			default:
				invalidFlags()
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***

	re2Str, err1 := parser.TransformRegExp(patternStr)
	if /*false &&*/ err1 == nil ***REMOVED***
		re2flags := ""
		if multiline ***REMOVED***
			re2flags += "m"
		***REMOVED***
		if ignoreCase ***REMOVED***
			re2flags += "i"
		***REMOVED***
		if len(re2flags) > 0 ***REMOVED***
			re2Str = fmt.Sprintf("(?%s:%s)", re2flags, re2Str)
		***REMOVED***

		pattern, err1 := regexp.Compile(re2Str)
		if err1 != nil ***REMOVED***
			err = fmt.Errorf("Invalid regular expression (re2): %s (%v)", re2Str, err1)
			return
		***REMOVED***

		p = (*regexpWrapper)(pattern)
	***REMOVED*** else ***REMOVED***
		var opts regexp2.RegexOptions = regexp2.ECMAScript
		if multiline ***REMOVED***
			opts |= regexp2.Multiline
		***REMOVED***
		if ignoreCase ***REMOVED***
			opts |= regexp2.IgnoreCase
		***REMOVED***
		regexp2Pattern, err1 := regexp2.Compile(patternStr, opts)
		if err1 != nil ***REMOVED***
			err = fmt.Errorf("Invalid regular expression (regexp2): %s (%v)", patternStr, err1)
			return
		***REMOVED***
		p = (*regexp2Wrapper)(regexp2Pattern)
	***REMOVED***
	return
***REMOVED***

func (r *Runtime) newRegExp(patternStr valueString, flags string, proto *Object) *Object ***REMOVED***
	pattern, global, ignoreCase, multiline, err := compileRegexp(patternStr.String(), flags)
	if err != nil ***REMOVED***
		panic(r.newSyntaxError(err.Error(), -1))
	***REMOVED***
	return r.newRegExpp(pattern, patternStr, global, ignoreCase, multiline, proto)
***REMOVED***

func (r *Runtime) builtin_newRegExp(args []Value) *Object ***REMOVED***
	var pattern valueString
	var flags string
	if len(args) > 0 ***REMOVED***
		if obj, ok := args[0].(*Object); ok ***REMOVED***
			if regexp, ok := obj.self.(*regexpObject); ok ***REMOVED***
				if len(args) < 2 || args[1] == _undefined ***REMOVED***
					return regexp.clone()
				***REMOVED*** else ***REMOVED***
					return r.newRegExp(regexp.source, args[1].String(), r.global.RegExpPrototype)
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if args[0] != _undefined ***REMOVED***
			pattern = args[0].ToString()
		***REMOVED***
	***REMOVED***
	if len(args) > 1 ***REMOVED***
		if a := args[1]; a != _undefined ***REMOVED***
			flags = a.String()
		***REMOVED***
	***REMOVED***
	if pattern == nil ***REMOVED***
		pattern = stringEmpty
	***REMOVED***
	return r.newRegExp(pattern, flags, r.global.RegExpPrototype)
***REMOVED***

func (r *Runtime) builtin_RegExp(call FunctionCall) Value ***REMOVED***
	flags := call.Argument(1)
	if flags == _undefined ***REMOVED***
		if obj, ok := call.Argument(0).(*Object); ok ***REMOVED***
			if _, ok := obj.self.(*regexpObject); ok ***REMOVED***
				return call.Arguments[0]
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return r.builtin_newRegExp(call.Arguments)
***REMOVED***

func (r *Runtime) regexpproto_exec(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		return this.exec(call.Argument(0).ToString())
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.exec called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_test(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.test(call.Argument(0).ToString()) ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.test called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_toString(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		var g, i, m string
		if this.global ***REMOVED***
			g = "g"
		***REMOVED***
		if this.ignoreCase ***REMOVED***
			i = "i"
		***REMOVED***
		if this.multiline ***REMOVED***
			m = "m"
		***REMOVED***
		return newStringValue(fmt.Sprintf("/%s/%s%s%s", this.source.String(), g, i, m))
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.toString called on incompatible receiver %s", call.This)
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getSource(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		return this.source
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.source getter called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getGlobal(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.global ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.global getter called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getMultiline(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.multiline ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.multiline getter called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) regexpproto_getIgnoreCase(call FunctionCall) Value ***REMOVED***
	if this, ok := r.toObject(call.This).self.(*regexpObject); ok ***REMOVED***
		if this.ignoreCase ***REMOVED***
			return valueTrue
		***REMOVED*** else ***REMOVED***
			return valueFalse
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		r.typeErrorResult(true, "Method RegExp.prototype.ignoreCase getter called on incompatible receiver %s", call.This.ToString())
		return nil
	***REMOVED***
***REMOVED***

func (r *Runtime) initRegExp() ***REMOVED***
	r.global.RegExpPrototype = r.NewObject()
	o := r.global.RegExpPrototype.self
	o._putProp("exec", r.newNativeFunc(r.regexpproto_exec, nil, "exec", nil, 1), true, false, true)
	o._putProp("test", r.newNativeFunc(r.regexpproto_test, nil, "test", nil, 1), true, false, true)
	o._putProp("toString", r.newNativeFunc(r.regexpproto_toString, nil, "toString", nil, 0), true, false, true)
	o.putStr("source", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getSource, nil, "get source", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.putStr("global", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getGlobal, nil, "get global", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.putStr("multiline", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getMultiline, nil, "get multiline", nil, 0),
		accessor:     true,
	***REMOVED***, false)
	o.putStr("ignoreCase", &valueProperty***REMOVED***
		configurable: true,
		getterFunc:   r.newNativeFunc(r.regexpproto_getIgnoreCase, nil, "get ignoreCase", nil, 0),
		accessor:     true,
	***REMOVED***, false)

	r.global.RegExp = r.newNativeFunc(r.builtin_RegExp, r.builtin_newRegExp, "RegExp", r.global.RegExpPrototype, 2)
	r.addToGlobal("RegExp", r.global.RegExp)
***REMOVED***
