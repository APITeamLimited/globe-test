package assert

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

//go:generate go run ../_codegen/main.go -output-package=assert -template=assertion_forward.go.tmpl -include-format-funcs
