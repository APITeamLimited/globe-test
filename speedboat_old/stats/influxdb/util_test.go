package influxdb

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseURL(t *testing.T) ***REMOVED***
	conf, db, err := parseURL("http://username:password@hostname.local:8086/db")
	assert.Nil(t, err, "couldn't parse URL")
	assert.Equal(t, "username", conf.Username, "incorrect username")
	assert.Equal(t, "password", conf.Password, "incorrect password")
	assert.Equal(t, "http://hostname.local:8086", conf.Addr, "incorrect address")
	assert.Equal(t, "db", db, "incorrect db")
***REMOVED***

func TestParseURLNoAuth(t *testing.T) ***REMOVED***
	conf, db, err := parseURL("http://hostname.local:8086/db")
	assert.Nil(t, err, "couldn't parse URL")
	assert.Equal(t, "http://hostname.local:8086", conf.Addr, "incorrect address")
	assert.Equal(t, "db", db, "incorrect db")
***REMOVED***

func TestParseURLNoDB(t *testing.T) ***REMOVED***
	_, _, err := parseURL("http://hostname.local:8086")
	assert.NotNil(t, err, "no error reported")
***REMOVED***

func TestMakeInfluxPoint(t *testing.T) ***REMOVED***
	now := time.Now()
	pt, err := makeInfluxPoint(stats.Sample***REMOVED***
		Stat:   &stats.Stat***REMOVED***Name: "test"***REMOVED***,
		Time:   now,
		Tags:   stats.Tags***REMOVED***"a": "b"***REMOVED***,
		Values: stats.Values***REMOVED***"value": 12345***REMOVED***,
	***REMOVED***)
	assert.NoError(t, err)
	assert.Equal(t, "test", pt.Name())
	assert.Equal(t, now, pt.Time())
	assert.EqualValues(t, map[string]string***REMOVED***"a": "b"***REMOVED***, pt.Tags())
	assert.EqualValues(t, map[string]interface***REMOVED******REMOVED******REMOVED***"value": float64(12345)***REMOVED***, pt.Fields())
***REMOVED***
