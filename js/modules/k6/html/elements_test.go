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

package html

import (
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/assert"
)

const testHTMLElems = `
<html>
<head></head>
<body>
	<a href="/testhref?querytxt#hashtext">0</a>
	<a href="http://example.com:80">1</a>
	<a href="http://example.com:81/path/file">2</a>
	<a href="https://ssl.example.com:443/">3</a>
	<a href="https://ssl.example.com:444/">4</a>
	<a href="http://username:password@example.com:80">5</a>
	<a href="http://example.com" rel="prev next" target="_self" type="rare" accesskey="q" hreflang="en-US" media="print">6</a>
	<area href="web.address.com"></area>
	<base href="/rel/path" target="_self"></base>
	
	<form id="form1" action="action_url" enctype="text/plain" target="_self">
		<label for="form_btn" id="form_btn_label"></label>
		<button id="form_btn" name="form_btn" accesskey="b" autofocus disabled></button>
		<label for="form_btn_2" id="form_btn_2_label"></label>
		<label id="wrapper_label">
			<button id="form_btn_2" type="button" formaction="override_action_url" formenctype="multipart/form-data" formmethod="post" formnovalidate formtarget="_top" value="initval"></button>
		</label>
	</form>
	<form id="form2"></form>
	<button id="named_form_btn" form="form2"></button>
	<button id="no_form_btn"></button>
	<canvas width="200"></canvas>
	<datalist id="datalist1"><option id="dl_opt_1"/><option id="dl_opt_2"/></datalist>
	<form method="post" id="fieldset_form"><fieldset><input id="test_dl_input" type="text" list="datalist1"><input type="text"><select></select><button></button><textarea></textarea></fieldset></form>
</body>
`

func TestElements(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx := common.WithRuntime(context.Background(), rt)
	rt.Set("src", testHTMLElems)
	rt.Set("html", common.Bind(rt, &HTML***REMOVED******REMOVED***, &ctx))
	// compileProtoElem()

	_, err := common.RunString(rt, `let doc = html.parseHTML(src)`)

	assert.NoError(t, err)
	assert.IsType(t, Selection***REMOVED******REMOVED***, rt.Get("doc").Export())

	t.Run("AnchorElement", func(t *testing.T) ***REMOVED***
		t.Run("Hash", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(0).hash()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "#hashtext", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Host", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(1).host()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "example.com", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(2).host()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "example.com:81", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(3).host()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "ssl.example.com", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(4).host()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "ssl.example.com:444", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Hostname", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(1).hostname()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "example.com", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Port", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(5).port()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "80", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Username", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(5).username()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "username", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Password", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(5).password()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "password", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Origin", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(5).origin()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "http://example.com:80", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Pathname", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(1).pathname()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(2).pathname()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "/path/file", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Protocol", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(4).protocol()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "https", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("RelList", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(6).relList()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, []string***REMOVED***"prev", "next"***REMOVED***, v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(5).relList()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, []string***REMOVED******REMOVED***, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Search", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(0).search()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "?querytxt", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Text", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("a").get(6).text()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "6", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("AreaElement", func(t *testing.T) ***REMOVED***
		if v, err := common.RunString(rt, `doc.find("area").get(0).toString()`); assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "web.address.com", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("ButtonElement", func(t *testing.T) ***REMOVED***
		t.Run("form", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).form().id()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "form1", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#named_form_btn").get(0).form().id()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "form2", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#no_form_btn").get(0).form()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, nil, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("formaction", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).formAction()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "action_url", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).formAction()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "override_action_url", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("formenctype", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).formEnctype()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "text/plain", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).formEnctype()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "multipart/form-data", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("formmethod", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).formMethod()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "get", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).formMethod()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "post", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("formnovalidate", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).formNoValidate()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, false, v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).formNoValidate()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("formtarget", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).formTarget()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "_self", v.Export())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).formTarget()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "_top", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("labels", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).labels()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, 1, len(v.Export().([]goja.Value)))
				assert.Equal(t, "form_btn_label", v.Export().([]goja.Value)[0].Export().(Element).Id())
			***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).labels()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, 2, len(v.Export().([]goja.Value)))
				assert.Equal(t, "wrapper_label", v.Export().([]goja.Value)[0].Export().(Element).Id())
				assert.Equal(t, "form_btn_2_label", v.Export().([]goja.Value)[1].Export().(Element).Id())
			***REMOVED***
		***REMOVED***)
		t.Run("name", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn").get(0).name()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "form_btn", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("value", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#form_btn_2").get(0).value()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "initval", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("CanvasElement", func(t *testing.T) ***REMOVED***
		t.Run("width", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("canvas").get(0).width()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(200), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("height", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("canvas").get(0).height()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(150), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("DataListElement options", func(t *testing.T) ***REMOVED***
		if v, err := common.RunString(rt, `doc.find("datalist").get(0).options()`); assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, 2, len(v.Export().([]goja.Value)))
			assert.Equal(t, "dl_opt_1", v.Export().([]goja.Value)[0].Export().(Element).Id())
			assert.Equal(t, "dl_opt_2", v.Export().([]goja.Value)[1].Export().(Element).Id())
		***REMOVED***
	***REMOVED***)
	t.Run("FieldSetElement", func(t *testing.T) ***REMOVED***
		t.Run("elements", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("fieldset").get(0).elements()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, 5, len(v.Export().([]goja.Value)))
			***REMOVED***
		***REMOVED***)
		t.Run("type", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("fieldset").get(0).type()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "fieldset", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("form", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("fieldset").get(0).form().id()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "fieldset_form", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("FormElement", func(t *testing.T) ***REMOVED***
		t.Run("elements", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#fieldset_form").get(0).elements()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, 6, len(v.Export().([]goja.Value)))
			***REMOVED***
		***REMOVED***)
		t.Run("length", func(t *testing.T) ***REMOVED***
			if v, err := common.RunString(rt, `doc.find("#fieldset_form").get(0).length()`); assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(6), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("method", func(t *testing.T) ***REMOVED***
			v1, err1 := common.RunString(rt, `doc.find("#form1").get(0).method()`)
			v2, err2 := common.RunString(rt, `doc.find("#fieldset_form").get(0).method()`)
			if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
				assert.Equal(t, "get", v1.Export())
				assert.Equal(t, "post", v2.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("InputElement", func(t *testing.T) ***REMOVED***
			t.Run("form", func(t *testing.T) ***REMOVED***
				if v, err := common.RunString(rt, `doc.find("#test_dl_input").get(0).list().options()`); assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, 2, len(v.Export().([]goja.Value)))
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
***REMOVED***
