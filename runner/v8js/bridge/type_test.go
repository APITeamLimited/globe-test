package bridge

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBridgeType(t *testing.T) ***REMOVED***
	tp := BridgeType(reflect.TypeOf(""))
	assert.Equal(t, reflect.String, tp.Kind)
	assert.Equal(t, nil, tp.Spec)
	assert.Equal(t, "", tp.JSONKey)
***REMOVED***

func TestBridgeTypeInvalid(t *testing.T) ***REMOVED***
	assert.Panics(t, func() ***REMOVED*** BridgeType(reflect.TypeOf(func() ***REMOVED******REMOVED***)) ***REMOVED***)
***REMOVED***

func TestBridgeTypeStruct(t *testing.T) ***REMOVED***
	tp := BridgeType(reflect.TypeOf(struct ***REMOVED***
		F1 string `json:"f1"`
		F2 int    `json:"f2"`
	***REMOVED******REMOVED******REMOVED***))
	assert.Contains(t, tp.Spec, "F1")
	assert.Equal(t, reflect.String, tp.Spec["F1"].Kind)
	assert.Equal(t, "f1", tp.Spec["F1"].JSONKey)
	assert.Contains(t, tp.Spec, "F2")
	assert.Equal(t, reflect.Int, tp.Spec["F2"].Kind)
	assert.Equal(t, "f2", tp.Spec["F2"].JSONKey)
***REMOVED***

func TestBridgeTypeStructNoTagExcluded(t *testing.T) ***REMOVED***
	tp := BridgeType(reflect.TypeOf(struct ***REMOVED***
		F1 string `json:"f1"`
		F2 int
	***REMOVED******REMOVED******REMOVED***))
	assert.Equal(t, 1, len(tp.Spec))
***REMOVED***
