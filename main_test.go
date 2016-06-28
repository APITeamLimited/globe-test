package main

import (
	"github.com/loadimpact/speedboat/stats"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestParseBackendStdout(t *testing.T) ***REMOVED***
	output, err := parseBackend("-")
	assert.NoError(t, err)
	assert.IsType(t, &stats.JSONBackend***REMOVED******REMOVED***, output)
***REMOVED***

func TestParseStagesSimple(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
***REMOVED***

func TestParseStagesSimpleTrailingDash(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10-"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
***REMOVED***

func TestParseStagesSimpleRamp(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10-15"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
***REMOVED***

func TestParseStagesSimpleRampZeroBackref(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"-15"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 0, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
***REMOVED***

func TestParseStagesSimpleMulti(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10", "15"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
***REMOVED***

func TestParseStagesSimpleMultiRamp(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10-15", "15-20"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 20, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
***REMOVED***

func TestParseStagesSimpleMultiRampBackref(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10-15", "-20"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 15, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 20, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
***REMOVED***

func TestParseStagesFixed(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10:15s"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 15*time.Second, stages[0].Duration)
***REMOVED***

func TestParseStagesFixedFluid(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10:5s", "15"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 5*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 5*time.Second, stages[1].Duration)
***REMOVED***

func TestParseStagesFixedFluidNoTimeLeft(t *testing.T) ***REMOVED***
	stages, err := parseStages([]string***REMOVED***"10:10s", "15"***REMOVED***, 10*time.Second)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(stages))
	assert.Equal(t, 10, stages[0].StartVUs)
	assert.Equal(t, 10, stages[0].EndVUs)
	assert.Equal(t, 10*time.Second, stages[0].Duration)
	assert.Equal(t, 15, stages[1].StartVUs)
	assert.Equal(t, 15, stages[1].EndVUs)
	assert.Equal(t, 0*time.Second, stages[1].Duration)
***REMOVED***

func TestParseStagesInvalid(t *testing.T) ***REMOVED***
	_, err := parseStages([]string***REMOVED***"a"***REMOVED***, 10*time.Second)
	assert.Error(t, err)
***REMOVED***

func TestParseStagesInvalidStart(t *testing.T) ***REMOVED***
	_, err := parseStages([]string***REMOVED***"a-15"***REMOVED***, 10*time.Second)
	assert.Error(t, err)
***REMOVED***

func TestParseStagesInvalidEnd(t *testing.T) ***REMOVED***
	_, err := parseStages([]string***REMOVED***"15-a"***REMOVED***, 10*time.Second)
	assert.Error(t, err)
***REMOVED***

func TestParseStagesInvalidTime(t *testing.T) ***REMOVED***
	_, err := parseStages([]string***REMOVED***"15:a"***REMOVED***, 10*time.Second)
	assert.Error(t, err)
***REMOVED***

func TestParseStagesInvalidTimeMissingUnit(t *testing.T) ***REMOVED***
	_, err := parseStages([]string***REMOVED***"15:10"***REMOVED***, 10*time.Second)
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
