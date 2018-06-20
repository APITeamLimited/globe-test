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

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/loadimpact/k6/js/common"
	"github.com/loadimpact/k6/stats"
	"github.com/stretchr/testify/assert"
)

const testGetFormHTML = `
<html>
<head>
	<title>This is the title</title>
</head>
<body>
	<form method="get" id="form1">
		<input name="input_with_value" type="text" value="value"/>
		<input name="input_without_value" type="text"/>
		<select name="select_one">
			<option value="not this option">no</option>
			<option value="yes this option" selected>yes</option>
		</select>
		<select name="select_multi" multiple>
			<option>option 1</option>
			<option selected>option 2</option>
			<option selected>option 3</option>
		</select>
		<textarea name="textarea" multiple>Lorem ipsum dolor sit amet</textarea>
	</form>
</body>
`

func myFormHandler(w http.ResponseWriter, r *http.Request) ***REMOVED***
	var body []byte
	var err error
	if r.URL.RawQuery != "" ***REMOVED***
		body, err = json.Marshal(struct ***REMOVED***
			Query url.Values `json:"query"`
		***REMOVED******REMOVED***
			Query: r.URL.Query(),
		***REMOVED***)
		if err != nil ***REMOVED***
			body = []byte(`***REMOVED***"error": "failed serializing json"***REMOVED***`)
		***REMOVED***
		w.Header().Set("Content-Type", "application/json")
	***REMOVED*** else ***REMOVED***
		w.Header().Set("Content-Type", "text/html")
		body = []byte(testGetFormHTML)
	***REMOVED***
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
	w.WriteHeader(200)
	_, _ = w.Write(body)
***REMOVED***

func TestResponse(t *testing.T) ***REMOVED***
	tb, state, samples, rt, _ := newRuntime(t)
	defer tb.Cleanup()
	root := state.Group
	sr := tb.Replacer.Replace

	tb.Mux.HandleFunc("/myforms/get", myFormHandler)

	t.Run("Html", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
			let res = http.request("GET", "HTTPBIN_URL/html");
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.body.indexOf("Herman Melville - Moby-Dick") == -1) ***REMOVED*** throw new Error("wrong body: " + res.body); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/html"), "", 200, "")

		t.Run("html", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `
				if (res.html().find("h1").text() != "Herman Melville - Moby-Dick") ***REMOVED*** throw new Error("wrong title: " + res.body); ***REMOVED***
			`)
			assert.NoError(t, err)

			t.Run("shorthand", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, `
					if (res.html("h1").text() != "Herman Melville - Moby-Dick") ***REMOVED*** throw new Error("wrong title: " + res.body); ***REMOVED***
				`)
				assert.NoError(t, err)
			***REMOVED***)

			t.Run("url", func(t *testing.T) ***REMOVED***
				_, err := common.RunString(rt, sr(`
					if (res.html().url != "HTTPBIN_URL/html") ***REMOVED*** throw new Error("url incorrect: " + res.html().url); ***REMOVED***
				`))
				assert.NoError(t, err)
			***REMOVED***)
		***REMOVED***)

		t.Run("group", func(t *testing.T) ***REMOVED***
			g, err := root.Group("my group")
			if assert.NoError(t, err) ***REMOVED***
				old := state.Group
				state.Group = g
				defer func() ***REMOVED*** state.Group = old ***REMOVED***()
			***REMOVED***

			_, err = common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/html");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				if (res.body.indexOf("Herman Melville - Moby-Dick") == -1) ***REMOVED*** throw new Error("wrong body: " + res.body); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/html"), "", 200, "::my group")
		***REMOVED***)
	***REMOVED***)
	t.Run("Json", func(t *testing.T) ***REMOVED***
		_, err := common.RunString(rt, sr(`
			let res = http.request("GET", "HTTPBIN_URL/get?a=1&b=2");
			if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			if (res.json().args.a != "1") ***REMOVED*** throw new Error("wrong ?a: " + res.json().args.a); ***REMOVED***
			if (res.json().args.b != "2") ***REMOVED*** throw new Error("wrong ?b: " + res.json().args.b); ***REMOVED***
		`))
		assert.NoError(t, err)
		assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get?a=1&b=2"), "", 200, "")

		t.Run("Invalid", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`http.request("GET", "HTTPBIN_URL/html").json();`))
			assert.EqualError(t, err, "GoError: invalid character '<' looking for beginning of value")
		***REMOVED***)
	***REMOVED***)

	t.Run("SubmitForm", func(t *testing.T) ***REMOVED***
		t.Run("withoutArgs", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/forms/post");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.submitForm()
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let data = res.json().form
				if (data.custname[0] !== "" ||
					data.extradata !== undefined ||
					data.comments[0] !== "" ||
					data.custemail[0] !== "" ||
					data.custtel[0] !== "" ||
					data.delivery[0] !== ""
				) ***REMOVED*** throw new Error("incorrect body: " + JSON.stringify(data, null, 4) ); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)

		t.Run("withFields", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/forms/post");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.submitForm(***REMOVED*** fields: ***REMOVED*** custname: "test", extradata: "test2" ***REMOVED*** ***REMOVED***)
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let data = res.json().form
				if (data.custname[0] !== "test" ||
					data.extradata[0] !== "test2" ||
					data.comments[0] !== "" ||
					data.custemail[0] !== "" ||
					data.custtel[0] !== "" ||
					data.delivery[0] !== ""
				) ***REMOVED*** throw new Error("incorrect body: " + JSON.stringify(data, null, 4) ); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)

		t.Run("withRequestParams", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/forms/post");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.submitForm(***REMOVED*** params: ***REMOVED*** headers: ***REMOVED*** "My-Fancy-Header": "SomeValue" ***REMOVED*** ***REMOVED******REMOVED***)
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let headers = res.json().headers
				if (headers["My-Fancy-Header"][0] !== "SomeValue" ) ***REMOVED*** throw new Error("incorrect headers: " + JSON.stringify(headers)); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)

		t.Run("withFormSelector", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/forms/post");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.submitForm(***REMOVED*** formSelector: 'form[method="post"]' ***REMOVED***)
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let data = res.json().form
				if (data.custname[0] !== "" ||
					data.extradata !== undefined ||
					data.comments[0] !== "" ||
					data.custemail[0] !== "" ||
					data.custtel[0] !== "" ||
					data.delivery[0] !== ""
				) ***REMOVED*** throw new Error("incorrect body: " + JSON.stringify(data, null, 4) ); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "POST", sr("HTTPBIN_URL/post"), "", 200, "")
		***REMOVED***)

		t.Run("withNonExistentForm", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/forms/post");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res.submitForm(***REMOVED*** formSelector: "#doesNotExist" ***REMOVED***)
			`))
			assert.EqualError(t, err, sr("GoError: no form found for selector '#doesNotExist' in response 'HTTPBIN_URL/forms/post'"))
		***REMOVED***)

		t.Run("withGetMethod", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/myforms/get");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.submitForm()
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let data = res.json().query
				if (data.input_with_value[0] !== "value" ||
					data.input_without_value[0] !== "" ||
					data.select_one[0] !== "yes this option" ||
					data.select_multi[0] !== "option 2,option 3" ||
					data.textarea[0] !== "Lorem ipsum dolor sit amet"
				) ***REMOVED*** throw new Error("incorrect body: " + JSON.stringify(data, null, 4) ); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/myforms/get"), "", 200, "")
		***REMOVED***)
	***REMOVED***)

	t.Run("ClickLink", func(t *testing.T) ***REMOVED***
		t.Run("withoutArgs", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/links/10/0");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.clickLink()
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/links/10/1"), "", 200, "")
		***REMOVED***)

		t.Run("withSelector", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/links/10/0");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.clickLink(***REMOVED*** selector: 'a:nth-child(4)' ***REMOVED***)
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/links/10/4"), "", 200, "")
		***REMOVED***)

		t.Run("withNonExistentLink", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL/links/10/0");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.clickLink(***REMOVED*** selector: 'a#doesNotExist' ***REMOVED***)
			`))
			assert.EqualError(t, err, sr("GoError: no element found for selector 'a#doesNotExist' in response 'HTTPBIN_URL/links/10/0'"))
		***REMOVED***)

		t.Run("withRequestParams", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, sr(`
				let res = http.request("GET", "HTTPBIN_URL");
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				res = res.clickLink(***REMOVED*** selector: 'a[href="/get"]', params: ***REMOVED*** headers: ***REMOVED*** "My-Fancy-Header": "SomeValue" ***REMOVED*** ***REMOVED*** ***REMOVED***)
				if (res.status != 200) ***REMOVED*** throw new Error("wrong status: " + res.status); ***REMOVED***
				let headers = res.json().headers
				if (headers["My-Fancy-Header"][0] !== "SomeValue" ) ***REMOVED*** throw new Error("incorrect headers: " + JSON.stringify(headers)); ***REMOVED***
			`))
			assert.NoError(t, err)
			assertRequestMetricsEmitted(t, stats.GetBufferedSamples(samples), "GET", sr("HTTPBIN_URL/get"), "", 200, "")
		***REMOVED***)
	***REMOVED***)
***REMOVED***
