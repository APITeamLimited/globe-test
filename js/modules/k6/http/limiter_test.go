package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlotLimiter(t *testing.T) ***REMOVED***
	l := NewSlotLimiter(1)
	l.Begin()
	done := false
	go func() ***REMOVED***
		done = true
		l.End()
	***REMOVED***()
	l.Begin()
	assert.True(t, done)
	l.End()
***REMOVED***

func TestMultiSlotLimiter(t *testing.T) ***REMOVED***
	t.Run("0", func(t *testing.T) ***REMOVED***
		l := NewMultiSlotLimiter(0)
		assert.Nil(t, l.Slot("test"))
	***REMOVED***)
	t.Run("1", func(t *testing.T) ***REMOVED***
		l := NewMultiSlotLimiter(1)
		assert.Equal(t, l.Slot("test"), l.Slot("test"))
		assert.NotNil(t, l.Slot("test"))
	***REMOVED***)
***REMOVED***
