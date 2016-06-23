package influxdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
