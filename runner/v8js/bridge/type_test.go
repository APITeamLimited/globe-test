package bridge

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestBridgeType(t *testing.T) ***REMOVED***
	tp := BridgeType(reflect.TypeOf(""))
	assert.Equal(t, reflect.String, tp.Kind)
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
	assert.Contains(t, tp.Spec, "F2")
***REMOVED***

func TestBridgeTypeStructNoTagExcluded(t *testing.T) ***REMOVED***
	tp := BridgeType(reflect.TypeOf(struct ***REMOVED***
		F1 string `json:"f1"`
		F2 int
	***REMOVED******REMOVED******REMOVED***))
	assert.Equal(t, 1, len(tp.Spec))
***REMOVED***
