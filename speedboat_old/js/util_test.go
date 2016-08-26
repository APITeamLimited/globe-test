package js

import (
	"errors"
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBodyFromValueUndefined(t *testing.T) ***REMOVED***
	body, isForm, err := bodyFromValue(otto.UndefinedValue())
	assert.NoError(t, err)
	assert.False(t, isForm)
	assert.Equal(t, "", body)
***REMOVED***

func TestBodyFromValueNull(t *testing.T) ***REMOVED***
	body, isForm, err := bodyFromValue(otto.NullValue())
	assert.NoError(t, err)
	assert.False(t, isForm)
	assert.Equal(t, "", body)
***REMOVED***

func TestBodyFromValueString(t *testing.T) ***REMOVED***
	val, err := otto.ToValue("abc123")
	assert.NoError(t, err)
	body, isForm, err := bodyFromValue(val)
	assert.NoError(t, err)
	assert.False(t, isForm)
	assert.Equal(t, "abc123", body)
***REMOVED***

func TestBodyFromValueObject(t *testing.T) ***REMOVED***
	vm := otto.New()
	val, err := vm.ToValue(map[string]string***REMOVED***"a": "b"***REMOVED***)
	assert.NoError(t, err)
	body, isForm, err := bodyFromValue(val)
	assert.NoError(t, err)
	assert.True(t, isForm)
	assert.Equal(t, "a=b", body)
***REMOVED***

func TestPutBodyInURL(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/?a=b", putBodyInURL("http://example.com/", "a=b"))
***REMOVED***

func TestPutBodyInURLWithQuery(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/?aa=bb&a=b", putBodyInURL("http://example.com/?aa=bb", "a=b"))
***REMOVED***

func TestResolveRedirectRelative(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/blah", resolveRedirect("http://example.com", "blah"))
***REMOVED***

func TestResolveRedirectRelativeParent(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/blah", resolveRedirect("http://example.com/aaa", "../blah"))
***REMOVED***

func TestResolveRedirectAbsolute(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/blah", resolveRedirect("http://example.com/aaa", "/blah"))
***REMOVED***

func TestResolveRedirectAbsoluteURL(t *testing.T) ***REMOVED***
	assert.Equal(t, "https://google.com/", resolveRedirect("http://example.com/aaa", "https://google.com/"))
***REMOVED***

func TestMake(t *testing.T) ***REMOVED***
	vm := otto.New()

	_, err := vm.Eval(`function MyType() ***REMOVED*** this.a = 'b'; ***REMOVED***;`)
	assert.NoError(t, err, "couldn't set up VM")

	obj, err := Make(vm, "MyType")
	assert.NoError(t, err)
	assert.NotNil(t, obj)
	assert.Equal(t, "Object", obj.Class())

	aVal, err := obj.Get("a")
	assert.NoError(t, err, "couldn't get 'a'")
	a, err := aVal.ToString()
	assert.NoError(t, err, "couldn't turn a into a string")
	assert.Equal(t, "b", a, "a != 'b'")
***REMOVED***

func TestJSCustomError(t *testing.T) ***REMOVED***
	vm := otto.New()
	vm.Set("fn", func(call otto.FunctionCall) otto.Value ***REMOVED***
		e := jsCustomError(vm, "CustomError", errors.New("test error"))
		str, err := e.ToString()
		assert.NoError(t, err)
		assert.Equal(t, "CustomError: test error", str)
		return otto.UndefinedValue()
	***REMOVED***)
	_, err := vm.Eval("fn()")
	assert.NoError(t, err)
***REMOVED***
