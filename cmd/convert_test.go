/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package cmd

import (
	"bytes"
	"os"
	"regexp"
	"testing"

	"io/ioutil"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func TestIntegrationConvertCmd(t *testing.T) ***REMOVED***
	var tmpFile, err = ioutil.TempFile("", "")
	if err != nil ***REMOVED***
		t.Fatalf("Couldn't create temporary file: %s", err)
	***REMOVED***
	harFile := tmpFile.Name()
	defer os.Remove(harFile)
	tmpFile.Close()
	t.Run("Correlate", func(t *testing.T) ***REMOVED***
		har, err := ioutil.ReadFile("testdata/example.har")
		assert.NoError(t, err)

		expectedTestPlan, err := ioutil.ReadFile("testdata/example.js")
		assert.NoError(t, err)

		defaultFs = afero.NewMemMapFs()

		err = afero.WriteFile(defaultFs, harFile, har, 0644)
		assert.NoError(t, err)

		buf := &bytes.Buffer***REMOVED******REMOVED***
		defaultWriter = buf

		assert.NoError(t, convertCmd.Flags().Set("correlate", "true"))
		assert.NoError(t, convertCmd.Flags().Set("no-batch", "true"))
		assert.NoError(t, convertCmd.Flags().Set("enable-status-code-checks", "true"))
		assert.NoError(t, convertCmd.Flags().Set("return-on-failed-check", "true"))

		err = convertCmd.RunE(convertCmd, []string***REMOVED***harFile***REMOVED***)

		// reset the convertCmd to default flags. There must be a nicer and less error prone way to do this...
		assert.NoError(t, convertCmd.Flags().Set("correlate", "false"))
		assert.NoError(t, convertCmd.Flags().Set("no-batch", "false"))
		assert.NoError(t, convertCmd.Flags().Set("enable-status-code-checks", "false"))
		assert.NoError(t, convertCmd.Flags().Set("return-on-failed-check", "false"))

		//Sanitizing to avoid windows problems with carriage returns
		re := regexp.MustCompile(`\r`)
		expected := re.ReplaceAllString(string(expectedTestPlan), ``)
		result := re.ReplaceAllString(buf.String(), ``)

		if assert.NoError(t, err) ***REMOVED***
			// assert.Equal suppresses the diff it is too big, so we add it as the test error message manually as well.
			diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff***REMOVED***
				A:        difflib.SplitLines(expected),
				B:        difflib.SplitLines(result),
				FromFile: "Expected",
				FromDate: "",
				ToFile:   "Actual",
				ToDate:   "",
				Context:  1,
			***REMOVED***)

			assert.Equal(t, expected, result, diff)
		***REMOVED***
	***REMOVED***)
	t.Run("Stdout", func(t *testing.T) ***REMOVED***
		defaultFs = afero.NewMemMapFs()
		err := afero.WriteFile(defaultFs, harFile, []byte(testHAR), 0644)
		assert.NoError(t, err)

		buf := &bytes.Buffer***REMOVED******REMOVED***
		defaultWriter = buf

		err = convertCmd.RunE(convertCmd, []string***REMOVED***harFile***REMOVED***)
		assert.NoError(t, err)
		assert.Equal(t, testHARConvertResult, buf.String())
	***REMOVED***)
	t.Run("Output file", func(t *testing.T) ***REMOVED***
		defaultFs = afero.NewMemMapFs()
		err := afero.WriteFile(defaultFs, harFile, []byte(testHAR), 0644)
		assert.NoError(t, err)

		err = convertCmd.Flags().Set("output", "/output.js")
		assert.NoError(t, err)
		err = convertCmd.RunE(convertCmd, []string***REMOVED***harFile***REMOVED***)
		assert.NoError(t, err)

		output, err := afero.ReadFile(defaultFs, "/output.js")
		assert.NoError(t, err)
		assert.Equal(t, testHARConvertResult, string(output))
	***REMOVED***)
	// TODO: test options injection; right now that's difficult because when there are multiple
	// options, they can be emitted in different order in the JSON
***REMOVED***
