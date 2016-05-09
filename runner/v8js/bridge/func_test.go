package bridge

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBridgeFuncEmpty(t *testing.T) ***REMOVED***
	assert.NotPanics(t, func() ***REMOVED*** BridgeFunc(func() ***REMOVED******REMOVED***) ***REMOVED***)
***REMOVED***

func TestBridgeFuncInvalid(t *testing.T) ***REMOVED***
	assert.Panics(t, func() ***REMOVED*** BridgeFunc(struct***REMOVED******REMOVED******REMOVED******REMOVED***) ***REMOVED***)
***REMOVED***

func TestBridgeFuncArgs(t *testing.T) ***REMOVED***
	fn := BridgeFunc(func(i int, a string) ***REMOVED******REMOVED***)
	assert.Equal(t, 2, len(fn.In))
***REMOVED***

func TestBridgeFuncReturns(t *testing.T) ***REMOVED***
	fn := BridgeFunc(func() (int, string) ***REMOVED*** return 0, "" ***REMOVED***)
	assert.Equal(t, 2, len(fn.Out))
***REMOVED***
