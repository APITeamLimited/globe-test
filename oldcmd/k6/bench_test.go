package main

import (
	"context"
	"testing"

	"github.com/loadimpact/k6/lib"
	"github.com/spf13/afero"
)

func BenchmarkJSRunners(b *testing.B) ***REMOVED***
	scripts := map[string]string***REMOVED***
		"Empty": `export default function() ***REMOVED******REMOVED***`,
		"HTTP":  `import http from "k6/http"; export default function() ***REMOVED*** http.get("http://localhost:8080/"); ***REMOVED***`,
	***REMOVED***
	for name, script := range scripts ***REMOVED***
		b.Run(name, func(b *testing.B) ***REMOVED***
			for _, t := range []string***REMOVED***"js", "js2"***REMOVED*** ***REMOVED***
				b.Run(t, func(b *testing.B) ***REMOVED***
					r, err := makeRunner(t, &lib.SourceData***REMOVED***
						Filename: "/script.js",
						Data:     []byte(script),
					***REMOVED***, afero.NewMemMapFs())
					if err != nil ***REMOVED***
						b.Error(err)
						return
					***REMOVED***

					b.Run("Spawn", func(b *testing.B) ***REMOVED***
						for i := 0; i < b.N; i++ ***REMOVED***
							_, err := r.NewVU()
							if err != nil ***REMOVED***
								b.Error(err)
								return
							***REMOVED***
						***REMOVED***
					***REMOVED***)

					b.Run("Run", func(b *testing.B) ***REMOVED***
						vu, err := r.NewVU()
						if err != nil ***REMOVED***
							b.Error(err)
							return
						***REMOVED***
						b.ResetTimer()

						for i := 0; i < b.N; i++ ***REMOVED***
							_, err := vu.RunOnce(context.Background())
							if err != nil ***REMOVED***
								b.Error(err)
							***REMOVED***
						***REMOVED***
					***REMOVED***)
				***REMOVED***)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***
