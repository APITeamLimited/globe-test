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
	// 	"fmt"
	"context"
	"testing"

	"github.com/dop251/goja"
	"github.com/loadimpact/k6/js/common"
	"github.com/stretchr/testify/assert"
)

const testHTMLElem = `
<html>
<head>
	<title>This is the title</title>
</head>
<body>
	<h1 id="top">Lorem ipsum</h1>
	<empty></empty>
	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
	<p>Lorem ipsum dolor sit amet, consectetur adipiscing elit.</p>
	pretext
	<div id="div_elem" class="class1 class2" lang="en">
		innerfirst
		<h2 id="h2_elem" class="class2">Nullam id nisi eget ex pharetra imperdiet.</h2>
		<span id="span1"><b>test content</b></span>
		<span id="span2">Maecenas augue ligula, aliquet sit amet maximus ut, vestibulum et magna</span>
		innerlast
	</div>
	aftertext
	<footer>This is the footer.</footer>
</body>
`

func TestElement(t *testing.T) ***REMOVED***
	rt := goja.New()
	rt.SetFieldNameMapper(common.FieldNameMapper***REMOVED******REMOVED***)

	ctx := common.WithRuntime(context.Background(), rt)
	rt.Set("src", testHTMLElem)
	rt.Set("html", common.Bind(rt, &HTML***REMOVED******REMOVED***, &ctx))
	compileProtoElem()

	_, err := common.RunString(rt, `let doc = html.parseHTML(src)`)

	assert.NoError(t, err)
	assert.IsType(t, Selection***REMOVED******REMOVED***, rt.Get("doc").Export())

	t.Run("NodeName", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("#top").get(0).nodeName`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "h1", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("NodeType", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("#top").get(0).nodeType`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "1", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("NodeValue", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("#top").get(0).firstChild.nodeValue`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "Lorem ipsum", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("InnerHtml", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("span").get(0).innerHTML`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "<b>test content</b>", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("TextContent", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("b").get(0).textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "test content", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("OwnerDocument", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).ownerDocument.nodeName`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "#document", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("Attributes", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).attributes`)
		if assert.NoError(t, err) ***REMOVED***
			attrs := v.Export().(map[string]Attribute)
			assert.Equal(t, "div_elem", attrs["id"].Value)
		***REMOVED***
	***REMOVED***)
	t.Run("FirstChild", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).firstChild.nodeValue`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "innerfirst")
		***REMOVED***
	***REMOVED***)
	t.Run("LastChild", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).lastChild.nodeValue`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "innerlast")
		***REMOVED***
	***REMOVED***)
	t.Run("ChildElementCount", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).childElementCount`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, int64(6), v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("FirstElementChild", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).firstElementChild.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "Nullam id nisi ")
		***REMOVED***
	***REMOVED***)
	t.Run("LastElementChild", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).lastElementChild.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "Maecenas augue ligula")
		***REMOVED***
	***REMOVED***)
	t.Run("PreviousSibling", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).previousSibling.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "pretext")
		***REMOVED***
	***REMOVED***)
	t.Run("NextSibling", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).nextSibling.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "aftertext")
		***REMOVED***
	***REMOVED***)
	t.Run("PreviousElementSibling", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).previousElementSibling.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "consectetur adipiscing elit")
		***REMOVED***
	***REMOVED***)
	t.Run("NextElementSibling", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).nextElementSibling.textContent`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Contains(t, v.Export(), "This is the footer.")
		***REMOVED***
	***REMOVED***)
	t.Run("ParentElement", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).parentElement.nodeName`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "body", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("ParentNode", func(t *testing.T) ***REMOVED***
		nodeVal, err1 := common.RunString(rt, `doc.find("html").get(0).parentNode.nodeName`)
		nilVal, err2 := common.RunString(rt, `doc.find("html").get(0).parentElement`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, "#document", nodeVal.String())
			assert.Equal(t, nil, nilVal.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("ChildNodes", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).childNodes`)
		if assert.NoError(t, err) ***REMOVED***
			nodes := valToElementList(v)
			assert.Equal(t, 7, len(nodes))
			assert.Contains(t, nodes[0].TextContent(), "innerfirst")
			assert.Contains(t, nodes[6].TextContent(), "innerlast")
		***REMOVED***
	***REMOVED***)
	t.Run("Children", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).children`)
		if assert.NoError(t, err) ***REMOVED***
			nodes := valToElementList(v)
			assert.Equal(t, 3, len(nodes))
			assert.Contains(t, nodes[0].TextContent(), "Nullam id nisi eget ex")
			assert.Contains(t, nodes[2].TextContent(), "Maecenas augue ligula")
		***REMOVED***
	***REMOVED***)
	t.Run("ClassList", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).classList`)
		if assert.NoError(t, err) ***REMOVED***
			clsNames := v.Export().([]string)
			assert.Equal(t, 2, len(clsNames))
			assert.Equal(t, []string***REMOVED***"class1", "class2"***REMOVED***, clsNames)
		***REMOVED***
	***REMOVED***)
	t.Run("ClassName", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).className`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "class1 class2", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("Lang", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).lang`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "en", v.String())
		***REMOVED***
	***REMOVED***)
	t.Run("ToString", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("div").get(0).toString()`)
		v2, err2 := common.RunString(rt, `doc.find("div").get(0).previousSibling.toString()`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, "[object html.Node]", v1.String())
			assert.Equal(t, "[object #text]", v2.String())
		***REMOVED***
	***REMOVED***)
	t.Run("HasAttribute", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("div").get(0).hasAttribute("id")`)
		v2, err2 := common.RunString(rt, `doc.find("div").get(0).hasAttribute("noattr")`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, true, v1.Export())
			assert.Equal(t, false, v2.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("GetAttribute", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).getAttribute("id")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "div_elem", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("GetAttributeNode", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).getAttributeNode("id")`)
		if assert.NoError(t, err) && assert.IsType(t, Attribute***REMOVED******REMOVED***, v.Export()) ***REMOVED***
			assert.Equal(t, "div_elem", v.Export().(Attribute).Value)
		***REMOVED***
	***REMOVED***)
	t.Run("HasAttributes", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("h1").get(0).hasAttributes()`)
		v2, err2 := common.RunString(rt, `doc.find("footer").get(0).hasAttributes()`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, true, v1.Export())
			assert.Equal(t, false, v2.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("HasChildNodes", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("p").get(0).hasChildNodes()`)
		v2, err2 := common.RunString(rt, `doc.find("empty").get(0).hasChildNodes()`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, true, v1.Export())
			assert.Equal(t, false, v2.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("IsSameNode", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("p").get(0).isSameNode(doc.find("p").get(1))`)
		v2, err2 := common.RunString(rt, `doc.find("p").get(0).isSameNode(doc.find("p").get(0))`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, false, v1.Export())
			assert.Equal(t, true, v2.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("IsEqualNode", func(t *testing.T) ***REMOVED***
		v1, err1 := common.RunString(rt, `doc.find("p").get(0).isEqualNode(doc.find("p").get(1))`)
		v2, err2 := common.RunString(rt, `doc.find("p").get(0).isEqualNode(doc.find("p").get(0))`)
		if assert.NoError(t, err1) && assert.NoError(t, err2) ***REMOVED***
			assert.Equal(t, true, v1.Export())
			assert.Equal(t, true, v2.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("GetElementsByClassName", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).getElementsByClassName("class2")`)
		if assert.NoError(t, err) ***REMOVED***
			elems := valToElementList(v)
			assert.Equal(t, []string***REMOVED***"div_elem", "h2_elem"***REMOVED***, []string***REMOVED***elems[0].Id().String(), elems[1].Id().String()***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("GetElementsByTagName", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).getElementsByTagName("span")`)
		if assert.NoError(t, err) ***REMOVED***
			elems := valToElementList(v)
			assert.Equal(t, 2, len(elems))
		***REMOVED***
	***REMOVED***)
	t.Run("QuerySelector", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).querySelector("#div_elem").id`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, "div_elem", v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("QuerySelectorAll", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).querySelectorAll("span")`)
		if assert.NoError(t, err) ***REMOVED***
			elems := valToElementList(v)
			assert.Equal(t, []string***REMOVED***"span1", "span2"***REMOVED***, []string***REMOVED***elems[0].Id().String(), elems[1].Id().String()***REMOVED***)
		***REMOVED***
	***REMOVED***)
	t.Run("Contains", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("body").get(0).contains(doc.find("body").get(0))`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)
	t.Run("Matches", func(t *testing.T) ***REMOVED***
		v, err := common.RunString(rt, `doc.find("div").get(0).matches("#div_elem")`)
		if assert.NoError(t, err) ***REMOVED***
			assert.Equal(t, true, v.Export())
		***REMOVED***
	***REMOVED***)
***REMOVED***
