package require

// TestingT is an interface wrapper around *testing.T
type TestingT interface ***REMOVED***
	Errorf(format string, args ...interface***REMOVED******REMOVED***)
	FailNow()
***REMOVED***

//go:generate go run ../_codegen/main.go -output-package=require -template=require.go.tmpl -include-format-funcs
