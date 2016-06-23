package main

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseBackendStdout(t *testing.T) ***REMOVED***
	output, err := parseBackend("-")
	assert.NoError(t, err)
	assert.IsType(t, &stats.JSONBackend***REMOVED******REMOVED***, output)
***REMOVED***

func TestGuessTypeURL(t *testing.T) ***REMOVED***
	assert.Equal(t, typeURL, guessType("http://example.com/"))
***REMOVED***

func TestGuessTypeJS(t *testing.T) ***REMOVED***
	assert.Equal(t, typeJS, guessType("script.js"))
***REMOVED***

func TestGuessTypeUnknown(t *testing.T) ***REMOVED***
	assert.Equal(t, "", guessType("script.txt"))
***REMOVED***
