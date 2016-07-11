package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormat(t *testing.T) ***REMOVED***
	stat := Stat***REMOVED***Name: "test"***REMOVED***
	v := (JSONBackend***REMOVED******REMOVED***).format(&Sample***REMOVED***
		Stat:   &stat,
		Tags:   Tags***REMOVED***"a": "b"***REMOVED***,
		Values: Values***REMOVED***"value": 12345.0***REMOVED***,
	***REMOVED***)

	assert.Equal(t, "test", v["stat"])
	assert.Equal(t, time.Time***REMOVED******REMOVED***, v["time"])

	assert.IsType(t, Tags***REMOVED******REMOVED***, v["tags"])
	assert.Len(t, v["tags"], 1)
	assert.Equal(t, "b", v["tags"].(Tags)["a"])

	assert.IsType(t, Values***REMOVED******REMOVED***, v["values"])
	assert.Len(t, v["values"], 1)
	assert.Equal(t, 12345.0, v["values"].(Values)["value"])
***REMOVED***

func TestFormatNilTagsBecomeEmptyMap(t *testing.T) ***REMOVED***
	stat := Stat***REMOVED***Name: "test"***REMOVED***
	v := (JSONBackend***REMOVED******REMOVED***).format(&Sample***REMOVED***
		Stat:   &stat,
		Values: Values***REMOVED***"value": 12345.0***REMOVED***,
	***REMOVED***)

	assert.IsType(t, Tags***REMOVED******REMOVED***, v["tags"])
	assert.Len(t, v["tags"], 0)
***REMOVED***
