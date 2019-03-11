package lib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextState(t *testing.T) ***REMOVED***
	st := &State***REMOVED******REMOVED***
	assert.Equal(t, st, GetState(WithState(context.Background(), st)))
***REMOVED***

func TestContextStateNil(t *testing.T) ***REMOVED***
	assert.Nil(t, GetState(context.Background()))
***REMOVED***
