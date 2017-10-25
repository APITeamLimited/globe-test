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
		<input id="text_input" type="text" value="input-text-value"/>
		<select id="select_one">
			<option value="not this option">no</option>
			<option value="yes this option" selected>yes</option>
		</select>
		<select id="select_text">
			<option>no text</option>
			<option selected>yes text</option>
		</select>
		<select id="select_multi" multiple>
			<option>option 1</option>
			<option selected>option 2</option>
			<option selected>option 3</option>
		</select>
		<textarea id="textarea" multiple>Lorem ipsum dolor sit amet</textarea>
	</form>

	<footer>This is the footer.</footer>
</body>
`

func TestParseHTML(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)
	ctx := common.WithRuntime(context.Background(), rt)
	rt.Set("src", testHTML)
	rt.Set("html", common.Bind(rt, New(), &ctx))

	// TODO: I literally cannot think of a snippet that makes goquery error.
	// I'm not sure if it's even possible without like, an invalid reader or something, which would
	// be impossible to cause from the JS side.
	_, err := common.RunString(rt, `let doc = html.parseHTML(src)`)
	assert.NoError(t, err)
	assert.IsType(t, Selection***REMOVED******REMOVED***, rt.Get("doc").Export())

	t.Run("Find", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("h1")`)
		if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 1, sel.Length())
			assert.Equal(t, "Lorem ipsum", sel.Text())
		***REMOVED***
	***REMOVED***)
	t.Run("Add", func(t *testing.T) ***REMOVED***
		t.Run("Selector", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("h1").add("footer")`)
			if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Lorem ipsumThis is the footer.", sel.Text())
			***REMOVED***
		***REMOVED***)
		t.Run("Selection", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("h1").add(doc.find("footer"))`)
			if assert.NoError(t, err) && assert.IsType(t, Selection***REMOVED******REMOVED***, v.Export()) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Lorem ipsumThis is the footer.", sel.Text())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)
	t.Run("Text", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("h1").text()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "Lorem ipsum", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Attr", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("h1").attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "top", v.Export())
		***REMOVED***
		t.Run("Default", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("h1").attr("id", "default")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "top", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Unset", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("h1").attr("class")`)
			if assert.NoError(t, err) ***REMOVED***
				assert.True(t, goja.IsUndefined(v), "v is not undefined: %v", v)
			***REMOVED***

			t.Run("Default", func(t *testing.T) ***REMOVED***
				v, err := common.RunString(rt, `doc.find("h1").attr("class", "default")`)
				if assert.NoError(t, err) ***REMOVED***
					assert.Equal(t, "default", v.Export())
				***REMOVED***
			***REMOVED***)
		***REMOVED***)
	***REMOVED***)
	t.Run("Html", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("h1").html()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "Lorem ipsum", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Val", func(t *testing.T) ***REMOVED***
		t.Run("Input", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#text_input").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "input-text-value", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select option[selected]", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#select_one option[selected]").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes this option", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Attr", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#select_one").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes this option", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Text", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#select_text").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "yes text", v.Export())
			***REMOVED***
		***REMOVED***)
		t.Run("Select Option Multiple", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#select_multi").val()`)
			if assert.NoError(t, err) ***REMOVED***
				var opts []string
				rt.ExportTo(v, &opts)
				assert.Equal(t, 2, len(opts))
				assert.Equal(t, "option 2", opts[0])
				assert.Equal(t, "option 3", opts[1])
			***REMOVED***
		***REMOVED***)
		t.Run("TextArea", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("#textarea").val()`)
			if assert.NoError(t, err) ***REMOVED***
				assert.Equal(t, "Lorem ipsum dolor sit amet", v.Export())
			***REMOVED***
		***REMOVED***)
	***REMOVED***)

	t.Run("Children", func(t *testing.T) ***REMOVED***
		t.Run("All", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("head").children()`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 1, sel.Length())
				assert.Equal(t, true, sel.Is("title"))
			***REMOVED***
		***REMOVED***)
		t.Run("With selector", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("body").children("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
				assert.Equal(t, "Nullam id nisi", sel.Last().Text()[0:14])
			***REMOVED***
		***REMOVED***)
	***REMOVED***)

	t.Run("Closest", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("textarea").closest("form").attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "form1", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Contents", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("head").contents()`)
		if assert.NoError(t, err) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 3, sel.Length())
			assert.Equal(t, "\n\t", sel.First().Text())
		***REMOVED***
	***REMOVED***)

	t.Run("Each", func(t *testing.T) ***REMOVED***
		t.Run("Func arg", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `***REMOVED*** var elems = []; doc.find("#select_multi option").each(function(idx, gqval) ***REMOVED*** elems[idx] = gqval.text() ***REMOVED***); elems ***REMOVED***`)
			if assert.NoError(t, err) ***REMOVED***
				var elems[] string
				rt.ExportTo(v, &elems)
				assert.Equal(t, 3, len(elems))
				assert.Equal(t, "option 1", elems[0])
			***REMOVED***
		***REMOVED***)

		t.Run("Invalid arg", func(t *testing.T) ***REMOVED***
			_, err := common.RunString(rt, `doc.find("#select_multi option").each("");`)
			if assert.Error(t, err) ***REMOVED***
				assert.IsType(t, &goja.InterruptedError***REMOVED******REMOVED***, err)
			***REMOVED***
		***REMOVED***)
	***REMOVED***)

	t.Run("Is", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("h1").is("h1")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Filter", func(t *testing.T) ***REMOVED***
		t.Run("String", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("body").children().filter("p")`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
			***REMOVED***
		***REMOVED***)

		t.Run("Function", func(t *testing.T) ***REMOVED***
			v, err := common.RunString(rt, `doc.find("body").children().filter(function(idx, val)***REMOVED*** return val.is("p") ***REMOVED***)`)
			if assert.NoError(t, err) ***REMOVED***
				sel := v.Export().(Selection).sel
				assert.Equal(t, 2, sel.Length())
 			***REMOVED***
		***REMOVED***)
	***REMOVED***)

	t.Run("End", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").children().filter("p").end()`)
		if assert.NoError(t, err) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 5, sel.Length())
		***REMOVED***
	***REMOVED***)

	t.Run("Eq", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").children().eq(3).attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "form1", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("First", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").children().first().attr("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "top", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Last", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").children().last().text()`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "This is the footer.", v.Export())
		***REMOVED***
	***REMOVED***)

	t.Run("Has", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").children().has("input")`)
		if assert.NoError(t, err) ***REMOVED***
			sel := v.Export().(Selection).sel
			assert.Equal(t, 1, sel.Length())
		***REMOVED***
	***REMOVED***)


***REMOVED***
