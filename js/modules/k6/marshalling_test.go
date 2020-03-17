/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2018 Load Impact
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

package k6_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/loadimpact/k6/js"
	"github.com/loadimpact/k6/lib"
	"github.com/loadimpact/k6/lib/testutils/httpmultibin"
	"github.com/loadimpact/k6/lib/types"
	"github.com/loadimpact/k6/loader"
	"github.com/loadimpact/k6/stats"
)

func TestSetupDataMarshalling(t *testing.T) ***REMOVED***
	tb := httpmultibin.NewHTTPMultiBin(t)
	defer tb.Cleanup()

	script := []byte(tb.Replacer.Replace(`
		import http from "k6/http";
		import html from "k6/html";
		import ws from "k6/ws";

		function make_jar() ***REMOVED***
			let jar = http.cookieJar();
			jar.set("test", "something")
			return jar;
		***REMOVED***

		export function setup() ***REMOVED***
			let res = http.get("HTTPBIN_URL/html");
			let html_selection = html.parseHTML(res.body);
			let ws_res = ws.connect("WSBIN_URL/ws-echo", function(socket)***REMOVED***
				socket.on("open", function() ***REMOVED***
					socket.send("test")
				***REMOVED***)
				socket.on("message", function (data)***REMOVED***
					if (!data=="test") ***REMOVED***
						throw new Error ("echo'd data doesn't match our message!");
					***REMOVED***
					socket.close()
				***REMOVED***);
			***REMOVED***);

			return ***REMOVED***
				http_response: res,
				html_selection: html_selection,
				html_element: html_selection.find('html'),
				html_attribute: html_selection.find('html').children('body').get(0).attributes(),
				jar: make_jar(),
				ws_res: ws_res,
			***REMOVED***;
		***REMOVED***

		function get_non_function_properties(object) ***REMOVED***
			return Object.keys(object).filter(element =>
				typeof(object[element]) !== "function");
		***REMOVED***

		function arrays_are_equal(first, second) ***REMOVED***
			if (first.length != second.length) ***REMOVED***
				return false
			***REMOVED***
			return first.every(function(element, idx) ***REMOVED***
				return element === second[idx]
			***REMOVED***);
		***REMOVED***

		function diff_object_properties(name, first, second) ***REMOVED***
			let first_properties = get_non_function_properties(first).sort();
			let second_properties = get_non_function_properties(second).sort();
			if (!(arrays_are_equal(first_properties, second_properties))) ***REMOVED***
				console.error("for " + name + ":\n" +
					"first_properties : " + JSON.stringify(first_properties) + "\n" +
					"second_properties: " + JSON.stringify(second_properties) + "\n" +
					"are not the same");
				throw new Error("not matching " + name);
			***REMOVED***
			first_properties.
				filter(element => typeof(first[element]) === "object").
					forEach(function(element) ***REMOVED***
						diff_object_properties(name+"."+element,
											   first[element],
											   second[element]);
			***REMOVED***);
		***REMOVED***

		export default function (data) ***REMOVED***
			diff_object_properties("setupdata", data, setup());
		***REMOVED***
	`))

	runner, err := js.New(
		&loader.SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/script.js"***REMOVED***, Data: script***REMOVED***,
		nil,
		lib.RuntimeOptions***REMOVED******REMOVED***,
	)

	require.NoError(t, err)

	err = runner.SetOptions(lib.Options***REMOVED***
		SetupTimeout: types.NullDurationFrom(5 * time.Second),
		Hosts:        tb.Dialer.Hosts,
	***REMOVED***)

	require.NoError(t, err)

	samples := make(chan<- stats.SampleContainer, 100)

	if !assert.NoError(t, runner.Setup(context.Background(), samples)) ***REMOVED***
		return
	***REMOVED***
	initVU, err := runner.NewVU(1, samples)
	if assert.NoError(t, err) ***REMOVED***
		vu := initVU.Activate(&lib.VUActivationParams***REMOVED***RunContext: context.Background()***REMOVED***)
		err := vu.RunOnce()
		assert.NoError(t, err)
	***REMOVED***
***REMOVED***
