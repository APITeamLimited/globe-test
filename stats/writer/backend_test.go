package writer

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFormat(t *testing.T) ***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	v := (Backend***REMOVED******REMOVED***).Format(&stats.Sample***REMOVED***
		Stat:   &stat,
		Tags:   stats.Tags***REMOVED***"a": "b"***REMOVED***,
		Values: stats.Values***REMOVED***"value": 12345.0***REMOVED***,
	***REMOVED***)

	assert.Equal(t, "test", v["stat"])
	assert.Equal(t, time.Time***REMOVED******REMOVED***, v["time"])

	assert.IsType(t, stats.Tags***REMOVED******REMOVED***, v["tags"])
	assert.Len(t, v["tags"], 1)
	assert.Equal(t, "b", v["tags"].(stats.Tags)["a"])

	assert.IsType(t, stats.Values***REMOVED******REMOVED***, v["values"])
	assert.Len(t, v["values"], 1)
	assert.Equal(t, 12345.0, v["values"].(stats.Values)["value"])
***REMOVED***

func TestFormatNilTagsBecomeEmptyMap(t *testing.T) ***REMOVED***
	stat := stats.Stat***REMOVED***Name: "test"***REMOVED***
	v := (Backend***REMOVED******REMOVED***).Format(&stats.Sample***REMOVED***
		Stat:   &stat,
		Values: stats.Values***REMOVED***"value": 12345.0***REMOVED***,
	***REMOVED***)

	assert.IsType(t, stats.Tags***REMOVED******REMOVED***, v["tags"])
	assert.Len(t, v["tags"], 0)
***REMOVED***
