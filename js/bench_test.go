package js

import (
	"github.com/robertkrimen/otto"
	"testing"
)

func BenchmarkOttoRun(b *testing.B) ***REMOVED***
	vm := otto.New()
	src := `1 + 1`

	b.Run("string", func(b *testing.B) ***REMOVED***
		b.ResetTimer()

		for i := 0; i < b.N; i++ ***REMOVED***
			_, err := vm.Run(src)
			if err != nil ***REMOVED***
				b.Error(err)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	b.Run("*Script", func(b *testing.B) ***REMOVED***
		script, err := vm.Compile("__snippet__", src)
		if err != nil ***REMOVED***
			b.Error(err)
			return
		***REMOVED***
		b.ResetTimer()

		for i := 0; i < b.N; i++ ***REMOVED***
			_, err := vm.Run(script)
			if err != nil ***REMOVED***
				b.Error(err)
				return
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***
