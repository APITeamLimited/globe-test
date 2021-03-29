package require

// Assertions provides assertion methods around the
// TestingT interface.
type Assertions struct ***REMOVED***
	t TestingT
***REMOVED***

// New makes a new Assertions object for the specified TestingT.
func New(t TestingT) *Assertions ***REMOVED***
	return &Assertions***REMOVED***
		t: t,
	***REMOVED***
***REMOVED***

//go:generate sh -c "cd ../_codegen && go build && cd - && ../_codegen/_codegen -output-package=require -template=require_forward.go.tmpl -include-format-funcs"
