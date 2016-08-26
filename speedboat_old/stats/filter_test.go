package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterBlank(t *testing.T) ***REMOVED***
	f := MakeFilter(nil, nil)
	assert.True(t, f.Check(Sample***REMOVED***Stat: &Stat***REMOVED***Name: "test"***REMOVED******REMOVED***))
***REMOVED***

func TestFilterOnly(t *testing.T) ***REMOVED***
	f := MakeFilter(nil, []string***REMOVED***"a"***REMOVED***)
	assert.True(t, f.Check(Sample***REMOVED***Stat: &Stat***REMOVED***Name: "a"***REMOVED******REMOVED***))
	assert.False(t, f.Check(Sample***REMOVED***Stat: &Stat***REMOVED***Name: "b"***REMOVED******REMOVED***))
***REMOVED***

func TestFilterExclude(t *testing.T) ***REMOVED***
	f := MakeFilter([]string***REMOVED***"a"***REMOVED***, nil)
	assert.False(t, f.Check(Sample***REMOVED***Stat: &Stat***REMOVED***Name: "a"***REMOVED******REMOVED***))
	assert.True(t, f.Check(Sample***REMOVED***Stat: &Stat***REMOVED***Name: "b"***REMOVED******REMOVED***))
***REMOVED***
