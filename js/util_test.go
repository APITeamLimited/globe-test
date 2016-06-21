package js

import (
	"github.com/robertkrimen/otto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBodyFromValueUndefined(t *testing.T) ***REMOVED***
	body, err := bodyFromValue(otto.UndefinedValue())
	assert.NoError(t, err)
	assert.Equal(t, "", body)
***REMOVED***

func TestBodyFromValueNull(t *testing.T) ***REMOVED***
	body, err := bodyFromValue(otto.NullValue())
	assert.NoError(t, err)
	assert.Equal(t, "", body)
***REMOVED***

func TestBodyFromValueString(t *testing.T) ***REMOVED***
	val, err := otto.ToValue("abc123")
	assert.NoError(t, err)
	body, err := bodyFromValue(val)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", body)
***REMOVED***

func TestBodyFromValueObject(t *testing.T) ***REMOVED***
	vm := otto.New()
	val, err := vm.ToValue(map[string]string***REMOVED***"a": "b"***REMOVED***)
	assert.NoError(t, err)
	body, err := bodyFromValue(val)
	assert.NoError(t, err)
	assert.Equal(t, "a=b", body)
***REMOVED***

func TestPutBodyInURL(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/?a=b", putBodyInURL("http://example.com/", "a=b"))
***REMOVED***

func TestPutBodyInURLWithQuery(t *testing.T) ***REMOVED***
	assert.Equal(t, "http://example.com/?aa=bb&a=b", putBodyInURL("http://example.com/?aa=bb", "a=b"))
***REMOVED***
