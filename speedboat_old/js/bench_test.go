package js

import (
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"sync"
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

func BenchmarkRunScriptParallelMultipleVMs(b *testing.B) ***REMOVED***
	for n := 0; n < b.N; n++ ***REMOVED***
		b.StopTimer()

		start := sync.WaitGroup***REMOVED******REMOVED***
		start.Add(1)

		end := sync.WaitGroup***REMOVED******REMOVED***
		for i := 0; i < 100; i++ ***REMOVED***
			end.Add(1)
			go func() ***REMOVED***
				defer end.Done()

				vm := otto.New()

				var i int
				vm.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
					i += 1
					v, err := call.Otto.ToValue(i)
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					return v
				***REMOVED***)

				script, err := vm.Compile("inline", `fn();`)
				if err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
				start.Wait()
				if _, err := vm.Run(script); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***()
		***REMOVED***

		b.StartTimer()
		start.Done()
		end.Wait()
	***REMOVED***
***REMOVED***

func BenchmarkRunScriptParallelClonedVMs(b *testing.B) ***REMOVED***
	vm := otto.New()

	var i int
	vm.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
		i += 1
		v, err := call.Otto.ToValue(i)
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***
		return v
	***REMOVED***)

	script, err := vm.Compile("inline", `fn();`)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	for n := 0; n < b.N; n++ ***REMOVED***
		b.StopTimer()

		start := sync.WaitGroup***REMOVED******REMOVED***
		start.Add(1)

		end := sync.WaitGroup***REMOVED******REMOVED***
		for i := 0; i < 100; i++ ***REMOVED***
			end.Add(1)
			go func() ***REMOVED***
				defer end.Done()

				myVM := vm.Copy()

				start.Wait()
				if _, err := myVM.Run(script); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***()
		***REMOVED***

		b.StartTimer()
		start.Done()
		end.Wait()
	***REMOVED***
***REMOVED***

func BenchmarkRunScriptParallelClonedVMsNoSharedState(b *testing.B) ***REMOVED***
	vm := otto.New()

	script, err := vm.Compile("inline", `fn();`)
	if err != nil ***REMOVED***
		b.Fatal(err)
	***REMOVED***

	for n := 0; n < b.N; n++ ***REMOVED***
		b.StopTimer()

		start := sync.WaitGroup***REMOVED******REMOVED***
		start.Add(1)

		end := sync.WaitGroup***REMOVED******REMOVED***
		for i := 0; i < 100; i++ ***REMOVED***
			end.Add(1)
			go func() ***REMOVED***
				defer end.Done()

				myVM := vm.Copy()

				var i int
				myVM.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
					i += 1
					v, err := call.Otto.ToValue(i)
					if err != nil ***REMOVED***
						panic(err)
					***REMOVED***
					return v
				***REMOVED***)

				start.Wait()
				if _, err := myVM.Run(script); err != nil ***REMOVED***
					b.Fatal(err)
				***REMOVED***
			***REMOVED***()
		***REMOVED***

		b.StartTimer()
		start.Done()
		end.Wait()
	***REMOVED***
***REMOVED***
