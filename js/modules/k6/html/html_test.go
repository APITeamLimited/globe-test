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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/js/common"
	"go.k6.io/k6/js/modulestest"
	"go.k6.io/k6/metrics"
)

const testHTML = `
<html>
<head>
	<title>This is the title</title>
</head>
<body>
	<h1 id="top" data-test="dataval" data-num-a="123" data-num-b="1.5" data-not-num-a="1.50" data-not-num-b="1.1e02">Lorem ipsum</h1>

	<p data-test-b="true" data-opts='***REMOVED***"id":101***REMOVED***' data-test-empty="">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec ac dui erat. Pellentesque eu euismod odio, eget fringilla ante. In vitae nulla at est tincidunt gravida sit amet maximus arcu. Sed accumsan tristique massa, blandit sodales quam malesuada eu. Morbi vitae luctus augue. Nunc nec ligula quam. Cras fringilla nulla leo, at dignissim enim accumsan vitae. Sed eu cursus sapien, a rhoncus lorem. Etiam sed massa egestas, bibendum quam sit amet, eleifend ipsum. Maecenas mi ante, consectetur at tincidunt id, suscipit nec sem. Integer congue elit vel ligula commodo ultricies. Suspendisse condimentum laoreet ligula at aliquet.</p>
	<p>Nullam id nisi eget ex pharetra imperdiet. Maecenas augue ligula, aliquet sit amet maximus ut, vestibulum et magna. Nam in arcu sed tortor volutpat porttitor sed eget dolor. Duis rhoncus est id dui porttitor, id molestie ex imperdiet. Proin purus ligula, pretium eleifend felis a, tempor feugiat mi. Cras rutrum pulvinar neque, eu dictum arcu. Cras purus metus, fermentum eget malesuada sit amet, dignissim non dui.</p>

	<form id="form1">
		<input id="text_input" name="text_input" type="text" value="input-text-value"/>
		<select id="select_one" name="select_one">
			<option value="not this option">no</option>
			<option value="yes this option" selected>yes</option>
		</select>
		<select id="select_text" name="select_text">
			<option>no text</option>
			<option selected>yes text</option>
		</select>
		<select id="select_multi" name="select_multi" multiple>
			<option>option 1</option>
			<option selected>option 2</option>
			<option selected>option 3</option>
		</select>
		<textarea id="textarea" name="textarea" multiple>Lorem ipsum dolor sit amet</textarea>
	</form>

	<footer>This is the footer.</footer>
</body>
`

func getTestModuleInstance(t testing.TB) (*goja.Runtime, *ModuleInstance) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	root := New()
	mockVU := &modulestest.VU***REMOVED***
		RuntimeField: rt,
		InitEnvField: &common.InitEnvironment***REMOVED***
			Registry: metrics.NewRegistry(),
		***REMOVED***,
		CtxField:   ctx,
		StateField: nil,
	***REMOVED***
	mi, ok := root.NewModuleInstance(mockVU).(*ModuleInstance)
	require.True(t, ok)

	require.NoError(t, rt.Set("html", mi.Exports().Default))

	return rt, mi
***REMOVED***

func getTestRuntimeAndModuleInstanceWithDoc(t testing.TB, html string) (*goja.Runtime, *ModuleInstance) ***REMOVED***
	t.Helper()

	rt, mi := getTestModuleInstance(t)
	require.NoError(t, rt.Set("src", html))

	_, err := rt.RunString(`var doc = html.parseHTML(src)`)

	require.NoError(t, err)
	require.IsType(t, Selection***REMOVED******REMOVED***, rt.Get("doc").Export())

	return rt, mi
***REMOVED***

func getTestRuntimeWithDoc(t testing.TB, html string) *goja.Runtime ***REMOVED***
	t.Helper()

	rt, _ := getTestRuntimeAndModuleInstanceWithDoc(t, html)

	return rt
***REMOVED***

func TestParseHTML(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Find", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("h1")`)
		if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 1, sel.Length())
			assert.Equal(t, "Lorem ipsum", sel.Text())
		***REMOVED***
	***REMOVED***)
	t.Run("Add", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("Selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").add("footer")`)
			if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Lorem ipsumThis is the footer.", sel.Text())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").add(doc.find("footer"))`)
			if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Lorem ipsumThis is the footer.", sel.Text())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Text", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("h1").text()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "Lorem ipsum", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Attr", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("h1").attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "top", v.Export())
		***REMOVED***
		t.Run("Default", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").attr("id", "default")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "top", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Unset", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").attr("class")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.True(t, goja.IsUndefined(v), "v is not undefined: %v", v)
			***REMOVED***

			t.Run("Default", func(t *testing.T) ***REMOVED***
				v, err := rt.RunString(`doc.find("h1").attr("class", "default")`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "default", v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
	t.Run("Html", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("h1").html()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "Lorem ipsum", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Val", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("Input", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#text_input").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "input-text-value", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select option[selected]", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_one option[selected]").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes this option", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Attr", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_one").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes this option", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Text", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_text").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes text", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Multiple", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_multi").val()`)
			var opts []string
			if assert.NoError(t, err) && rt.ExportTo(v, &opts) == nil ***REMOVED***
				assert.Equal(t, 2, len(opts))
				assert.Equal(t, "option 2", opts[0])
				assert.Equal(t, "option 3", opts[1])
			***REMOVED***
		***REMOVED***)
		t.Run("TextArea", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#textarea").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "Lorem ipsum dolor sit amet", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Children", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("All", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("head").children()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
				assert.True(t, sel.Is("title"))
			***REMOVED***
		***REMOVED***)
		t.Run("With selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Nullam id nisi", sel.Last().Text()[0:14])
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Closest", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("textarea").closest("form").attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "form1", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Contents", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("head").contents()`)
		if assert.NoError(t, err) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 3, sel.Length())
			assert.Equal(t, "\n\t", sel.First().Text())
		***REMOVED***
	***REMOVED***)
	t.Run("Each", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("Func arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`***REMOVED*** var elems = []; doc.find("#select_multi option").each(function(idx, elem) ***REMOVED*** elems[idx] = elem.innerHTML(); ***REMOVED***); elems ***REMOVED***`)
			var elems []string
			if assert.NoError(t, err) && rt.ExportTo(v, &elems) == nil ***REMOVED***
				assert.Equal(t, 3, len(elems))
				assert.Equal(t, "option 1", elems[0])
				assert.Equal(t, "option 2", elems[1])
			***REMOVED***
		***REMOVED***)
		t.Run("Invalid arg", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(`doc.find("#select_multi option").each("");`)
			if assert.Error(t, err) ***REMOVED***
				assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
				assert.Contains(t, err.Error(), "must be a function")
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Is", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").is("h1")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Function selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").is(function(idx, val)***REMOVED*** return val.text() == "Lorem ipsum" ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().first().is(doc.find("h1"))`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Filter", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().filter("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)
		t.Run("Function", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().filter(function(idx, val)***REMOVED*** return val.is("p") ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().filter(doc.find("p"))`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("End", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("body").children().filter("p").end()`)
		if assert.NoError(t, err) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 5, sel.Length())
		***REMOVED***
	***REMOVED***)
	t.Run("Eq", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("body").children().eq(3).attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "form1", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("First", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("body").children().first().attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "top", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Last", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("body").children().last().text()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "This is the footer.", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Has", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().has("input").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().has(doc.find("input")).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Map", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("Valid", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_multi option").map(function(idx, val) ***REMOVED*** return val.text() ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				mapped, ok := v.Export().([]string)
				assert.True(t, ok)
				assert.Equal(t, 3, len(mapped))
				assert.Equal(t, []string***REMOVED***"option 1", "option 2", "option 3"***REMOVED***, mapped)
			***REMOVED***
		***REMOVED***)
		t.Run("Invalid arg", func(t *testing.T) ***REMOVED***
			_, err := rt.RunString(`doc.find("#select_multi option").map("");`)
			if assert.Error(t, err) ***REMOVED***
				assert.IsType(t, &goja.Exception***REMOVED******REMOVED***, err)
				assert.Contains(t, err.Error(), "must be a function")
			***REMOVED***
		***REMOVED***)
		t.Run("Map with attr must return string", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("#select_multi").map(function(idx, val) ***REMOVED*** return val.attr("name") ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				mapped, ok := v.Export().([]string)
				assert.True(t, ok)
				assert.Equal(t, 1, len(mapped))
				assert.Equal(t, []string***REMOVED***"select_multi"***REMOVED***, mapped)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Next", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").next()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
				assert.True(t, sel.Is("p"))
			***REMOVED***
		***REMOVED***)
		t.Run("Filter arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").next("form")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("NextAll", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextAll()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 4, sel.Length())
			***REMOVED***
		***REMOVED***)
		t.Run("Filter arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextAll("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Prev", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prev()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.True(t, sel.Is("form"))
			***REMOVED***
		***REMOVED***)
		t.Run("Filter arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prev("form")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("PrevAll", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").prevAll()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 3, sel.Length())
			***REMOVED***
		***REMOVED***)
		t.Run("Filter arg", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").prevAll("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("PrevUntil", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prevUntil("h1").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prevUntil(doc.find("h1")).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("String filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").prevUntil("h1", "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").prevUntil(doc.find("h1"), "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prevUntil().size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(4), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("footer").prevUntil(null, "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("NextUntil", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil("footer").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil(doc.find("footer")).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("String filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil("footer", "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil(doc.find("footer"), "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil().size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(4), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").nextUntil(null, "p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Parent", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No filter", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parent().attr("id")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "form1", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parent("form").attr("id")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "form1", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Parents", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No filter", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parents().size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parents("body").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("ParentsUntil", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil("html").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil(doc.find("html")).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("String filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil("html", "body").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Query filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil(doc.find("html"), "body").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil().size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("All filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("textarea").parentsUntil(null, "body").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Not", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("String selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().not("p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().not(doc.find("p")).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Function selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().not(function(idx, val)***REMOVED*** return val.is("p") ***REMOVED***).size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Siblings", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No filter", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").siblings().size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(4), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").siblings("p").size()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(2), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Slice", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No filter", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().slice(1, 2)`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
				assert.True(t, sel.Is("p"))
				assert.Contains(t, sel.Text(), "Lorem ipsum dolor")
			***REMOVED***
		***REMOVED***)
		t.Run("Filtered", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().slice(3)`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, true, sel.Eq(0).Is("form"))
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Get", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No args", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().get()`)
			if assert.NoError(t, err) ***REMOVED***
				elems, ok := v.Export().([]goja.Value)

				assert.True(t, ok)
				assert.Equal(t, "h1", elems[0].Export().(Element).NodeName())
				assert.Equal(t, "p", elems[1].Export().(Element).NodeName())
				assert.Equal(t, "footer", elems[4].Export().(Element).NodeName())
			***REMOVED***
		***REMOVED***)
		t.Run("+ve index", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().get(1)`)
			if assert.NoError(t, err) ***REMOVED***
				elem, _ := v.Export().(Element)
				assert.Contains(t, elem.InnerHTML(), "Lorem ipsum dolor sit amet")
			***REMOVED***
		***REMOVED***)
		t.Run("-ve index", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().get(-1)`)
			if assert.NoError(t, err) ***REMOVED***
				elem, _ := v.Export().(Element)
				assert.Equal(t, "This is the footer.", elem.InnerHTML().String())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("ToArray", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		v, err := rt.RunString(`doc.find("p").toArray()`)
		if assert.NoError(t, err) ***REMOVED***
			arr, ok := v.Export().([]Selection)

			assert.True(t, ok)
			assert.Equal(t, 2, len(arr))
			assert.Equal(t, 1, arr[0].sel.Length())
			assert.Equal(t, 1, arr[1].sel.Length())
			assert.Contains(t, arr[0].sel.Text(), "Lorem ipsum dolor sit amet")
			assert.Contains(t, arr[1].sel.Text(), "Nullam id nisi eget")
		***REMOVED***
	***REMOVED***)
	t.Run("Index", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("No args", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").index()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(1), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("String selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("form").index("body > *")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(3), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection selector", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("body").children().index(doc.find("footer"))`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(4), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Data <h1>", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("string attr", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data("test")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "dataval", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("numeric attr 1", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data("num-a")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(123), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("numeric attr 2", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data("num-b")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, float64(1.5), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("not numeric attr 1", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data("not-num-a")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "1.50", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("not numeric attr 2", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data("not-num-b")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "1.1e02", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("dataset", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("h1").data()`)
			if assert.NoError(t, err) ***REMOVED***
				data, ok := v.Export().(map[string]interface***REMOVED******REMOVED***)

				assert.True(t, ok)
				assert.Equal(t, "dataval", data["test"])
				assert.Equal(t, float64(123), data["numA"])
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Data <p>", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testHTML)

		t.Run("boolean attr", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data("test-b")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("snakeCase attr name", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data("testB")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("empty string", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data("test-empty")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, nil, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("json attr", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data("opts").id`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(101), v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("dataset property", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data().testB`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, true, v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("dataset object", func(t *testing.T) ***REMOVED***
			v, err := rt.RunString(`doc.find("p").data().opts.id`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, int64(101), v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
***REMOVED***
