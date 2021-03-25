package require

// TestingT is an interface wrapper around *testing.T
type TestingT interface ***REMOVED***
	Errorf(format string, args ...interface***REMOVED******REMOVED***)
	FailNow()
***REMOVED***

type tHelper interface ***REMOVED***
	Helper()
***REMOVED***

// ComparisonAssertionFunc is a common function prototype when comparing two values.  Can be useful
// for table driven tests.
type ComparisonAssertionFunc func(TestingT, interface***REMOVED******REMOVED***, interface***REMOVED******REMOVED***, ...interface***REMOVED******REMOVED***)

// ValueAssertionFunc is a common function prototype when validating a single value.  Can be useful
// for table driven tests.
type ValueAssertionFunc func(TestingT, interface***REMOVED******REMOVED***, ...interface***REMOVED******REMOVED***)

// BoolAssertionFunc is a common function prototype when validating a bool value.  Can be useful
// for table driven tests.
type BoolAssertionFunc func(TestingT, bool, ...interface***REMOVED******REMOVED***)

// ValuesAssertionFunc is a common function prototype when validating an error value.  Can be useful
// for table driven tests.
type ErrorAssertionFunc func(TestingT, error, ...interface***REMOVED******REMOVED***)

//go:generate go run ../_codegen/main.go -output-package=require -template=require.go.tmpl -include-format-funcs
