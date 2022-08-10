package cmd

import (
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHAR = `
***REMOVED***
	"log": ***REMOVED***
		"version": "1.2",
		"creator": ***REMOVED***
		"name": "WebInspector",
		"version": "537.36"
		***REMOVED***,
		"pages": [
		***REMOVED***
			"startedDateTime": "2018-01-21T19:48:40.432Z",
			"id": "page_2",
			"title": "https://golang.org/",
			"pageTimings": ***REMOVED***
			"onContentLoad": 590.3389999875799,
			"onLoad": 1593.1009999476373
			***REMOVED***
		***REMOVED***
		],
		"entries": [
		***REMOVED***
			"startedDateTime": "2018-01-21T19:48:40.587Z",
			"time": 147.5899999756366,
			"request": ***REMOVED***
				"method": "GET",
				"url": "https://golang.org/",
				"httpVersion": "http/2.0+quic/39",
				"headers": [
					***REMOVED***
					"name": "pragma",
					"value": "no-cache"
					***REMOVED***
				],
				"queryString": [],
				"cookies": [],
				"headersSize": -1,
				"bodySize": 0
			***REMOVED***,
			"cache": ***REMOVED******REMOVED***,
			"timings": ***REMOVED***
				"blocked": 0.43399997614324004,
				"dns": -1,
				"ssl": -1,
				"connect": -1,
				"send": 0.12700003571808005,
				"wait": 149.02899996377528,
				"receive": 0,
				"_blocked_queueing": -1
			***REMOVED***,
			"serverIPAddress": "172.217.22.177",
			"pageref": "page_2"
		***REMOVED***
		]
	***REMOVED***
***REMOVED***
`

const testHARConvertResult = `import ***REMOVED*** group, sleep ***REMOVED*** from 'k6';
import http from 'k6/http';

// Version: 1.2
// Creator: WebInspector

export let options = ***REMOVED***
    maxRedirects: 0,
***REMOVED***;

export default function() ***REMOVED***

	group("page_2 - https://golang.org/", function() ***REMOVED***
		let req, res;
		req = [***REMOVED***
			"method": "get",
			"url": "https://golang.org/",
			"params": ***REMOVED***
				"headers": ***REMOVED***
					"pragma": "no-cache"
				***REMOVED***
			***REMOVED***
		***REMOVED***];
		res = http.batch(req);
		// Random sleep between 20s and 40s
		sleep(Math.floor(Math.random()*20+20));
	***REMOVED***);

***REMOVED***
`

func TestConvertCmdCorrelate(t *testing.T) ***REMOVED***
	t.Parallel()
	har, err := ioutil.ReadFile("testdata/example.har")
	require.NoError(t, err)

	expectedTestPlan, err := ioutil.ReadFile("testdata/example.js")
	require.NoError(t, err)

	testState := newGlobalTestState(t)
	require.NoError(t, afero.WriteFile(testState.fs, "correlate.har", har, 0o644))
	testState.args = []string***REMOVED***
		"k6", "convert", "--output=result.js", "--correlate=true", "--no-batch=true",
		"--enable-status-code-checks=true", "--return-on-failed-check=true", "correlate.har",
	***REMOVED***

	newRootCommand(testState.globalState).execute()

	result, err := afero.ReadFile(testState.fs, "result.js")
	require.NoError(t, err)

	// Sanitizing to avoid windows problems with carriage returns
	re := regexp.MustCompile(`\r`)
	expected := re.ReplaceAllString(string(expectedTestPlan), ``)
	resultStr := re.ReplaceAllString(string(result), ``)

	if assert.NoError(t, err) ***REMOVED***
		// assert.Equal suppresses the diff it is too big, so we add it as the test error message manually as well.
		diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff***REMOVED***
			A:        difflib.SplitLines(expected),
			B:        difflib.SplitLines(resultStr),
			FromFile: "Expected",
			FromDate: "",
			ToFile:   "Actual",
			ToDate:   "",
			Context:  1,
		***REMOVED***)

		assert.Equal(t, expected, resultStr, diff)
	***REMOVED***
***REMOVED***

func TestConvertCmdStdout(t *testing.T) ***REMOVED***
	t.Parallel()
	testState := newGlobalTestState(t)
	require.NoError(t, afero.WriteFile(testState.fs, "stdout.har", []byte(testHAR), 0o644))
	testState.args = []string***REMOVED***"k6", "convert", "stdout.har"***REMOVED***

	newRootCommand(testState.globalState).execute()
	assert.Equal(t, testHARConvertResult, testState.stdOut.String())
***REMOVED***

func TestConvertCmdOutputFile(t *testing.T) ***REMOVED***
	t.Parallel()

	testState := newGlobalTestState(t)
	require.NoError(t, afero.WriteFile(testState.fs, "output.har", []byte(testHAR), 0o644))
	testState.args = []string***REMOVED***"k6", "convert", "--output", "result.js", "output.har"***REMOVED***

	newRootCommand(testState.globalState).execute()

	output, err := afero.ReadFile(testState.fs, "result.js")
	assert.NoError(t, err)
	assert.Equal(t, testHARConvertResult, string(output))
***REMOVED***

// TODO: test options injection; right now that's difficult because when there are multiple
// options, they can be emitted in different order in the JSON
