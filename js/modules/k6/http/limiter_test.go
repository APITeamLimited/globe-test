package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimiter(t *testing.T) ***REMOVED***
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
