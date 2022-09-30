package http

import (
	"fmt"
	"sort"
	"testing"

	"github.com/APITeamLimited/globe-test/worker/libWorker"
	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpectedStatuses(t *testing.T) ***REMOVED***
	t.Parallel()
	rt, _, _ := getTestModuleInstance(t)

	cases := map[string]struct ***REMOVED***
		code, err string
		expected  expectedStatuses
	***REMOVED******REMOVED***
		"good example": ***REMOVED***
			expected: expectedStatuses***REMOVED***exact: []int***REMOVED***200, 300***REMOVED***, minmax: [][2]int***REMOVED******REMOVED***200, 300***REMOVED******REMOVED******REMOVED***,
			code:     `(http.expectedStatuses(200, 300, ***REMOVED***min: 200, max:300***REMOVED***))`,
		***REMOVED***,

		"strange example": ***REMOVED***
			expected: expectedStatuses***REMOVED***exact: []int***REMOVED***200, 300***REMOVED***, minmax: [][2]int***REMOVED******REMOVED***200, 300***REMOVED******REMOVED******REMOVED***,
			code:     `(http.expectedStatuses(200, 300, ***REMOVED***min: 200, max:300, other: "attribute"***REMOVED***))`,
		***REMOVED***,

		"string status code": ***REMOVED***
			code: `(http.expectedStatuses(200, "300", ***REMOVED***min: 200, max:300***REMOVED***))`,
			err:  "argument number 2 to expectedStatuses was neither an integer nor an object like ***REMOVED***min:100, max:329***REMOVED***",
		***REMOVED***,

		"string max status code": ***REMOVED***
			code: `(http.expectedStatuses(200, 300, ***REMOVED***min: 200, max:"300"***REMOVED***))`,
			err:  "both min and max need to be integers for argument number 3",
		***REMOVED***,
		"float status code": ***REMOVED***
			err:  "argument number 2 to expectedStatuses was neither an integer nor an object like ***REMOVED***min:100, max:329***REMOVED***",
			code: `(http.expectedStatuses(200, 300.5, ***REMOVED***min: 200, max:300***REMOVED***))`,
		***REMOVED***,

		"float max status code": ***REMOVED***
			err:  "both min and max need to be integers for argument number 3",
			code: `(http.expectedStatuses(200, 300, ***REMOVED***min: 200, max:300.5***REMOVED***))`,
		***REMOVED***,
		"no arguments": ***REMOVED***
			code: `(http.expectedStatuses())`,
			err:  "no arguments",
		***REMOVED***,
	***REMOVED***

	for name, testCase := range cases ***REMOVED***
		name, testCase := name, testCase
		t.Run(name, func(t *testing.T) ***REMOVED***
			val, err := rt.RunString(testCase.code)
			if testCase.err == "" ***REMOVED***
				require.NoError(t, err)
				got := new(expectedStatuses)
				err = rt.ExportTo(val, &got)
				require.NoError(t, err)
				require.Equal(t, testCase.expected, *got)
				return // the t.Run
			***REMOVED***

			require.Error(t, err)
			exc := err.(*goja.Exception)
			require.Contains(t, exc.Error(), testCase.err)
		***REMOVED***)
	***REMOVED***
***REMOVED***

type expectedSample struct ***REMOVED***
	tags    map[string]string
	metrics []string
***REMOVED***

func TestResponseCallbackInAction(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, _, samples, rt, mii := newRuntime(t)
	sr := tb.Replacer.Replace

	HTTPMetricsWithoutFailed := []string***REMOVED***
		workerMetrics.HTTPReqsName,
		workerMetrics.HTTPReqBlockedName,
		workerMetrics.HTTPReqConnectingName,
		workerMetrics.HTTPReqDurationName,
		workerMetrics.HTTPReqReceivingName,
		workerMetrics.HTTPReqWaitingName,
		workerMetrics.HTTPReqSendingName,
		workerMetrics.HTTPReqTLSHandshakingName,
	***REMOVED***

	allHTTPMetrics := append(HTTPMetricsWithoutFailed, workerMetrics.HTTPReqFailedName)

	testCases := map[string]struct ***REMOVED***
		code            string
		expectedSamples []expectedSample
	***REMOVED******REMOVED***
		"basic": ***REMOVED***
			code: `http.request("GET", "HTTPBIN_URL/redirect/1");`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/redirect/1"),
						"name":              sr("HTTPBIN_URL/redirect/1"),
						"status":            "302",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/get"),
						"name":              sr("HTTPBIN_URL/get"),
						"status":            "200",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"overwrite per request": ***REMOVED***
			code: `
			http.setResponseCallback(http.expectedStatuses(200));
			res = http.request("GET", "HTTPBIN_URL/redirect/1");
			`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/redirect/1"),
						"name":              sr("HTTPBIN_URL/redirect/1"),
						"status":            "302",
						"group":             "",
						"expected_response": "false", // this is on purpose
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/get"),
						"name":              sr("HTTPBIN_URL/get"),
						"status":            "200",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,

		"global overwrite": ***REMOVED***
			code: `http.request("GET", "HTTPBIN_URL/redirect/1", null, ***REMOVED***responseCallback: http.expectedStatuses(200)***REMOVED***);`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/redirect/1"),
						"name":              sr("HTTPBIN_URL/redirect/1"),
						"status":            "302",
						"group":             "",
						"expected_response": "false", // this is on purpose
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/get"),
						"name":              sr("HTTPBIN_URL/get"),
						"status":            "200",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"per request overwrite with null": ***REMOVED***
			code: `http.request("GET", "HTTPBIN_URL/redirect/1", null, ***REMOVED***responseCallback: null***REMOVED***);`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method": "GET",
						"url":    sr("HTTPBIN_URL/redirect/1"),
						"name":   sr("HTTPBIN_URL/redirect/1"),
						"status": "302",
						"group":  "",
						"proto":  "HTTP/1.1",
					***REMOVED***,
					metrics: HTTPMetricsWithoutFailed,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method": "GET",
						"url":    sr("HTTPBIN_URL/get"),
						"name":   sr("HTTPBIN_URL/get"),
						"status": "200",
						"group":  "",
						"proto":  "HTTP/1.1",
					***REMOVED***,
					metrics: HTTPMetricsWithoutFailed,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
		"global overwrite with null": ***REMOVED***
			code: `
			http.setResponseCallback(null);
			res = http.request("GET", "HTTPBIN_URL/redirect/1");
			`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method": "GET",
						"url":    sr("HTTPBIN_URL/redirect/1"),
						"name":   sr("HTTPBIN_URL/redirect/1"),
						"status": "302",
						"group":  "",
						"proto":  "HTTP/1.1",
					***REMOVED***,
					metrics: HTTPMetricsWithoutFailed,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method": "GET",
						"url":    sr("HTTPBIN_URL/get"),
						"name":   sr("HTTPBIN_URL/get"),
						"status": "200",
						"group":  "",
						"proto":  "HTTP/1.1",
					***REMOVED***,
					metrics: HTTPMetricsWithoutFailed,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(name, func(t *testing.T) ***REMOVED***
			mii.defaultClient.responseCallback = defaultExpectedStatuses.match

			_, err := rt.RunString(sr(testCase.code))
			assert.NoError(t, err)
			bufSamples := workerMetrics.GetBufferedSamples(samples)

			reqsCount := 0
			for _, container := range bufSamples ***REMOVED***
				for _, sample := range container.GetSamples() ***REMOVED***
					if sample.Metric.Name == "http_reqs" ***REMOVED***
						reqsCount++
					***REMOVED***
				***REMOVED***
			***REMOVED***

			require.Equal(t, len(testCase.expectedSamples), reqsCount)

			for i, expectedSample := range testCase.expectedSamples ***REMOVED***
				assertRequestMetricsEmittedSingle(t, bufSamples[i], expectedSample.tags, expectedSample.metrics, nil)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestResponseCallbackBatch(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, _, samples, rt, mii := newRuntime(t)
	sr := tb.Replacer.Replace

	HTTPMetricsWithoutFailed := []string***REMOVED***
		workerMetrics.HTTPReqsName,
		workerMetrics.HTTPReqBlockedName,
		workerMetrics.HTTPReqConnectingName,
		workerMetrics.HTTPReqDurationName,
		workerMetrics.HTTPReqReceivingName,
		workerMetrics.HTTPReqWaitingName,
		workerMetrics.HTTPReqSendingName,
		workerMetrics.HTTPReqTLSHandshakingName,
	***REMOVED***

	allHTTPMetrics := append(HTTPMetricsWithoutFailed, workerMetrics.HTTPReqFailedName)
	// IMPORTANT: the tests here depend on the fact that the url they hit can be ordered in the same
	// order as the expectedSamples even if they are made concurrently
	testCases := map[string]struct ***REMOVED***
		code            string
		expectedSamples []expectedSample
	***REMOVED******REMOVED***
		"basic": ***REMOVED***
			code: `
	http.batch([["GET", "HTTPBIN_URL/status/200", null, ***REMOVED***responseCallback: null***REMOVED***],
			["GET", "HTTPBIN_URL/status/201"],
			["GET", "HTTPBIN_URL/status/202", null, ***REMOVED***responseCallback: http.expectedStatuses(4)***REMOVED***],
			["GET", "HTTPBIN_URL/status/405", null, ***REMOVED***responseCallback: http.expectedStatuses(405)***REMOVED***],
	]);`,
			expectedSamples: []expectedSample***REMOVED***
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method": "GET",
						"url":    sr("HTTPBIN_URL/status/200"),
						"name":   sr("HTTPBIN_URL/status/200"),
						"status": "200",
						"group":  "",
						"proto":  "HTTP/1.1",
					***REMOVED***,
					metrics: HTTPMetricsWithoutFailed,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/status/201"),
						"name":              sr("HTTPBIN_URL/status/201"),
						"status":            "201",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/status/202"),
						"name":              sr("HTTPBIN_URL/status/202"),
						"status":            "202",
						"group":             "",
						"expected_response": "false",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
				***REMOVED***
					tags: map[string]string***REMOVED***
						"method":            "GET",
						"url":               sr("HTTPBIN_URL/status/405"),
						"name":              sr("HTTPBIN_URL/status/405"),
						"status":            "405",
						"error_code":        "1405",
						"group":             "",
						"expected_response": "true",
						"proto":             "HTTP/1.1",
					***REMOVED***,
					metrics: allHTTPMetrics,
				***REMOVED***,
			***REMOVED***,
		***REMOVED***,
	***REMOVED***
	for name, testCase := range testCases ***REMOVED***
		testCase := testCase
		t.Run(name, func(t *testing.T) ***REMOVED***
			mii.defaultClient.responseCallback = defaultExpectedStatuses.match

			_, err := rt.RunString(sr(testCase.code))
			assert.NoError(t, err)
			bufSamples := workerMetrics.GetBufferedSamples(samples)

			reqsCount := 0
			for _, container := range bufSamples ***REMOVED***
				for _, sample := range container.GetSamples() ***REMOVED***
					if sample.Metric.Name == "http_reqs" ***REMOVED***
						reqsCount++
					***REMOVED***
				***REMOVED***
			***REMOVED***
			sort.Slice(bufSamples, func(i, j int) bool ***REMOVED***
				iURL, _ := bufSamples[i].GetSamples()[0].Tags.Get("url")
				jURL, _ := bufSamples[j].GetSamples()[0].Tags.Get("url")
				return iURL < jURL
			***REMOVED***)

			require.Equal(t, len(testCase.expectedSamples), reqsCount)

			for i, expectedSample := range testCase.expectedSamples ***REMOVED***
				assertRequestMetricsEmittedSingle(t, bufSamples[i], expectedSample.tags, expectedSample.metrics, nil)
			***REMOVED***
		***REMOVED***)
	***REMOVED***
***REMOVED***

func TestResponseCallbackInActionWithoutPassedTag(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, state, samples, rt, _ := newRuntime(t)
	sr := tb.Replacer.Replace
	allHTTPMetrics := []string***REMOVED***
		workerMetrics.HTTPReqsName,
		workerMetrics.HTTPReqFailedName,
		workerMetrics.HTTPReqBlockedName,
		workerMetrics.HTTPReqConnectingName,
		workerMetrics.HTTPReqDurationName,
		workerMetrics.HTTPReqReceivingName,
		workerMetrics.HTTPReqSendingName,
		workerMetrics.HTTPReqWaitingName,
		workerMetrics.HTTPReqTLSHandshakingName,
	***REMOVED***
	deleteSystemTag(state, workerMetrics.TagExpectedResponse.String())

	_, err := rt.RunString(sr(`http.request("GET", "HTTPBIN_URL/redirect/1", null, ***REMOVED***responseCallback: http.expectedStatuses(200)***REMOVED***);`))
	assert.NoError(t, err)
	bufSamples := workerMetrics.GetBufferedSamples(samples)

	reqsCount := 0
	for _, container := range bufSamples ***REMOVED***
		for _, sample := range container.GetSamples() ***REMOVED***
			if sample.Metric.Name == "http_reqs" ***REMOVED***
				reqsCount++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	require.Equal(t, 2, reqsCount)

	tags := map[string]string***REMOVED***
		"method": "GET",
		"url":    sr("HTTPBIN_URL/redirect/1"),
		"name":   sr("HTTPBIN_URL/redirect/1"),
		"status": "302",
		"group":  "",
		"proto":  "HTTP/1.1",
	***REMOVED***
	assertRequestMetricsEmittedSingle(t, bufSamples[0], tags, allHTTPMetrics, func(sample workerMetrics.Sample) ***REMOVED***
		if sample.Metric.Name == workerMetrics.HTTPReqFailedName ***REMOVED***
			require.EqualValues(t, sample.Value, 1)
		***REMOVED***
	***REMOVED***)
	tags["url"] = sr("HTTPBIN_URL/get")
	tags["name"] = tags["url"]
	tags["status"] = "200"
	assertRequestMetricsEmittedSingle(t, bufSamples[1], tags, allHTTPMetrics, func(sample workerMetrics.Sample) ***REMOVED***
		if sample.Metric.Name == workerMetrics.HTTPReqFailedName ***REMOVED***
			require.EqualValues(t, sample.Value, 0)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func TestDigestWithResponseCallback(t *testing.T) ***REMOVED***
	t.Parallel()
	tb, _, samples, rt, _ := newRuntime(t)

	urlWithCreds := tb.Replacer.Replace(
		"http://testuser:testpwd@HTTPBIN_IP:HTTPBIN_PORT/digest-auth/auth/testuser/testpwd",
	)

	allHTTPMetrics := []string***REMOVED***
		workerMetrics.HTTPReqsName,
		workerMetrics.HTTPReqFailedName,
		workerMetrics.HTTPReqBlockedName,
		workerMetrics.HTTPReqConnectingName,
		workerMetrics.HTTPReqDurationName,
		workerMetrics.HTTPReqReceivingName,
		workerMetrics.HTTPReqSendingName,
		workerMetrics.HTTPReqWaitingName,
		workerMetrics.HTTPReqTLSHandshakingName,
	***REMOVED***
	_, err := rt.RunString(fmt.Sprintf(`
		var res = http.get(%q,  ***REMOVED*** auth: "digest" ***REMOVED***);
		if (res.status !== 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
		if (res.error_code !== 0) ***REMOVED*** throw new Error("wrong error code: " + res.error_code); ***REMOVED***
	`, urlWithCreds))
	require.NoError(t, err)
	bufSamples := workerMetrics.GetBufferedSamples(samples)

	reqsCount := 0
	for _, container := range bufSamples ***REMOVED***
		for _, sample := range container.GetSamples() ***REMOVED***
			if sample.Metric.Name == "http_reqs" ***REMOVED***
				reqsCount++
			***REMOVED***
		***REMOVED***
	***REMOVED***

	require.Equal(t, 2, reqsCount)

	urlRaw := tb.Replacer.Replace(
		"http://HTTPBIN_IP:HTTPBIN_PORT/digest-auth/auth/testuser/testpwd")

	tags := map[string]string***REMOVED***
		"method":            "GET",
		"url":               urlRaw,
		"name":              urlRaw,
		"status":            "401",
		"group":             "",
		"proto":             "HTTP/1.1",
		"expected_response": "true",
		"error_code":        "1401",
	***REMOVED***
	assertRequestMetricsEmittedSingle(t, bufSamples[0], tags, allHTTPMetrics, func(sample workerMetrics.Sample) ***REMOVED***
		if sample.Metric.Name == workerMetrics.HTTPReqFailedName ***REMOVED***
			require.EqualValues(t, sample.Value, 0)
		***REMOVED***
	***REMOVED***)
	tags["status"] = "200"
	delete(tags, "error_code")
	assertRequestMetricsEmittedSingle(t, bufSamples[1], tags, allHTTPMetrics, func(sample workerMetrics.Sample) ***REMOVED***
		if sample.Metric.Name == workerMetrics.HTTPReqFailedName ***REMOVED***
			require.EqualValues(t, sample.Value, 0)
		***REMOVED***
	***REMOVED***)
***REMOVED***

func deleteSystemTag(state *libWorker.State, tag string) ***REMOVED***
	enabledTags := state.Options.SystemTags.Map()
	delete(enabledTags, tag)
	tagsList := make([]string, 0, len(enabledTags))
	for k := range enabledTags ***REMOVED***
		tagsList = append(tagsList, k)
	***REMOVED***
	state.Options.SystemTags = workerMetrics.ToSystemTagSet(tagsList)
***REMOVED***
