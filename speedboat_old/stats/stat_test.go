package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestApplyIntentDefault(t *testing.T) ***REMOVED***
	v := ApplyIntent(10.0, DefaultIntent)
	assert.IsType(t, 10.0, v)
***REMOVED***

func TestApplyIntentTime(t *testing.T) ***REMOVED***
	v := ApplyIntent(10.0, TimeIntent)
	assert.IsType(t, time.Duration(10), v)
***REMOVED***
