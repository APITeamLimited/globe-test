package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) ***REMOVED***
	c := Collector***REMOVED******REMOVED***
	stat := Stat***REMOVED***Name: "test"***REMOVED***
	c.Add(Point***REMOVED***Stat: &stat, Values: Values***REMOVED***"value": 12345***REMOVED******REMOVED***)
	assert.Equal(t, 1, len(c.batch))
	assert.Equal(t, &stat, c.batch[0].Stat)
	assert.Equal(t, 12345.0, c.batch[0].Values["value"])
***REMOVED***

func TestAddNoStat(t *testing.T) ***REMOVED***
	c := Collector***REMOVED******REMOVED***
	c.Add(Point***REMOVED***Values: Values***REMOVED***"value": 12345***REMOVED******REMOVED***)
	assert.Equal(t, 0, len(c.batch))
***REMOVED***

func TestAddNoValues(t *testing.T) ***REMOVED***
	c := Collector***REMOVED******REMOVED***
	c.Add(Point***REMOVED***Stat: &Stat***REMOVED***Name: "test"***REMOVED******REMOVED***)
	assert.Equal(t, 0, len(c.batch))
***REMOVED***

func TestAddFixesTime(t *testing.T) ***REMOVED***
	c := Collector***REMOVED******REMOVED***
	c.Add(Point***REMOVED***Stat: &Stat***REMOVED***Name: "test"***REMOVED***, Values: Values***REMOVED***"value": 12345***REMOVED******REMOVED***)
	assert.False(t, c.batch[0].Time.IsZero())
***REMOVED***

func TestDrain(t *testing.T) ***REMOVED***
	c := Collector***REMOVED******REMOVED***
	c.Add(Point***REMOVED***Stat: &Stat***REMOVED***Name: "test"***REMOVED***, Values: Values***REMOVED***"value": 12345***REMOVED******REMOVED***)
	batch := c.drain()
	assert.Equal(t, 1, len(batch))
***REMOVED***
