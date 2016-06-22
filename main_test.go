package main

import (
	"github.com/loadimpact/speedboat/sampler/stream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseOutputStdoutJSON(t *testing.T) ***REMOVED***
	output, err := parseOutput("-", "json")
	assert.NoError(t, err)
	assert.IsType(t, &stream.JSONOutput***REMOVED******REMOVED***, output)
***REMOVED***

func TestParseOutputStdoutCSV(t *testing.T) ***REMOVED***
	output, err := parseOutput("-", "csv")
	assert.NoError(t, err)
	assert.IsType(t, &stream.CSVOutput***REMOVED******REMOVED***, output)
***REMOVED***

func TestParseOutputStdoutUnknown(t *testing.T) ***REMOVED***
	_, err := parseOutput("-", "not a real format")
	assert.Error(t, err)
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
