package js

import (
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkCallGoFunction(b *testing.B) ***REMOVED***
	i := 0
	vm := otto.New()
	vm.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
		i += 1
		return otto.UndefinedValue()
	***REMOVED***)
	script, err := vm.Compile("script", `fn();`)
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCallGoFunctionReturn(b *testing.B) ***REMOVED***
	i := 0
	vm := otto.New()
	vm.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
		i += 1
		v, err := otto.ToValue(i)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		return v
	***REMOVED***)
	script, err := vm.Compile("script", `fn();`)
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCallJSFunction(b *testing.B) ***REMOVED***
	vm := otto.New()

	_, err := vm.Eval(`var i = 0; function fn() ***REMOVED*** i++; ***REMOVED***;`)
	assert.Nil(b, err)

	script, err := vm.Compile("script", `fn();`)
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCallJSFunctionExplicitUndefined(b *testing.B) ***REMOVED***
	vm := otto.New()

	_, err := vm.Eval(`var i = 0; function fn() ***REMOVED*** i++; return undefined; ***REMOVED***;`)
	assert.Nil(b, err)

	script, err := vm.Compile("script", `fn();`)
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***

func BenchmarkCallJSFunctionReturn(b *testing.B) ***REMOVED***
	vm := otto.New()

	_, err := vm.Eval(`var i = 0; function fn() ***REMOVED*** i++; return i; ***REMOVED***;`)
	assert.Nil(b, err)

	script, err := vm.Compile("script", `fn();`)
	assert.Nil(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ ***REMOVED***
		if _, err := vm.Run(script); err != nil ***REMOVED***
			b.Fatal(err)
		***REMOVED***
	***REMOVED***
***REMOVED***
