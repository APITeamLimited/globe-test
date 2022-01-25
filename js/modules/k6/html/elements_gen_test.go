/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
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
	"testing"
)

var textTests = []struct ***REMOVED***
	id       string
	property string
	data     string
***REMOVED******REMOVED***
	***REMOVED***"a1", "download", "file:///path/name"***REMOVED***,
	***REMOVED***"a1", "referrerPolicy", "no-referrer"***REMOVED***,
	***REMOVED***"a1", "href", "http://test.url"***REMOVED***,
	***REMOVED***"a1", "target", "__blank"***REMOVED***,
	***REMOVED***"a1", "type", "text/html"***REMOVED***,
	***REMOVED***"a1", "accessKey", "w"***REMOVED***,
	***REMOVED***"a1", "hrefLang", "es"***REMOVED***,
	***REMOVED***"a1", "toString", "http://test.url"***REMOVED***,
	***REMOVED***"a2", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"a2", "accessKey", ""***REMOVED***,
	***REMOVED***"audio1", "src", "foo.wav"***REMOVED***,
	***REMOVED***"audio1", "crossOrigin", "anonymous"***REMOVED***,
	***REMOVED***"audio1", "currentSrc", "foo.wav"***REMOVED***,
	***REMOVED***"audio1", "mediaGroup", "testgroup"***REMOVED***,
	***REMOVED***"base1", "href", "foo.html"***REMOVED***,
	***REMOVED***"base1", "target", "__any"***REMOVED***,
	***REMOVED***"btn1", "accessKey", "e"***REMOVED***,
	***REMOVED***"btn1", "type", "button"***REMOVED***,
	***REMOVED***"btn2", "type", "submit"***REMOVED***,
	***REMOVED***"btn3", "type", "submit"***REMOVED***,
	***REMOVED***"data1", "value", "121"***REMOVED***,
	***REMOVED***"data2", "value", ""***REMOVED***,
	***REMOVED***"embed1", "type", "video/avi"***REMOVED***,
	***REMOVED***"embed1", "src", "movie.avi"***REMOVED***,
	***REMOVED***"embed1", "width", "640"***REMOVED***,
	***REMOVED***"embed1", "height", "480"***REMOVED***,
	***REMOVED***"fset1", "name", "fset1_name"***REMOVED***,
	***REMOVED***"form1", "target", "__self"***REMOVED***,
	***REMOVED***"form1", "action", "submit_url"***REMOVED***,
	***REMOVED***"form1", "enctype", "text/plain"***REMOVED***,
	***REMOVED***"form1", "encoding", "text/plain"***REMOVED***,
	***REMOVED***"form1", "acceptCharset", "ISO-8859-1"***REMOVED***,
	***REMOVED***"form1", "target", "__self"***REMOVED***,
	***REMOVED***"form1", "autocomplete", "off"***REMOVED***,
	***REMOVED***"form2", "enctype", "application/x-www-form-urlencoded"***REMOVED***,
	***REMOVED***"form2", "autocomplete", "on"***REMOVED***,
	***REMOVED***"iframe1", "referrerPolicy", "no-referrer"***REMOVED***,
	***REMOVED***"iframe2", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"iframe3", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"iframe1", "width", "640"***REMOVED***,
	***REMOVED***"iframe1", "height", "480"***REMOVED***,
	***REMOVED***"iframe1", "name", "frame_name"***REMOVED***,
	***REMOVED***"iframe1", "src", "testframe.html"***REMOVED***,
	***REMOVED***"img1", "src", "test.png"***REMOVED***,
	***REMOVED***"img1", "currentSrc", "test.png"***REMOVED***,
	***REMOVED***"img1", "sizes", "100vw,50vw"***REMOVED***,
	***REMOVED***"img1", "srcset", "large.jpg 1024w,medium.jpg 640w"***REMOVED***,
	***REMOVED***"img1", "alt", "alt text"***REMOVED***,
	***REMOVED***"img1", "crossOrigin", "anonymous"***REMOVED***,
	***REMOVED***"img1", "name", "img_name"***REMOVED***,
	***REMOVED***"img1", "useMap", "#map_name"***REMOVED***,
	***REMOVED***"img1", "referrerPolicy", "origin"***REMOVED***,
	***REMOVED***"img2", "crossOrigin", "use-credentials"***REMOVED***,
	***REMOVED***"img2", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"img3", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"input1", "name", "input1_name"***REMOVED***,
	***REMOVED***"input1", "type", "button"***REMOVED***,
	***REMOVED***"input1", "value", "input1-val"***REMOVED***,
	***REMOVED***"input1", "defaultValue", "input1-val"***REMOVED***,
	***REMOVED***"input2", "type", "text"***REMOVED***,
	***REMOVED***"input2", "value", ""***REMOVED***,
	***REMOVED***"input5", "alt", "input_img"***REMOVED***,
	***REMOVED***"input5", "src", "input.png"***REMOVED***,
	***REMOVED***"input5", "width", "80"***REMOVED***,
	***REMOVED***"input5", "height", "40"***REMOVED***,
	***REMOVED***"input6", "accept", ".jpg,.png"***REMOVED***,
	***REMOVED***"input7", "autocomplete", "off"***REMOVED***,
	***REMOVED***"input7", "pattern", "..."***REMOVED***,
	***REMOVED***"input7", "placeholder", "help text"***REMOVED***,
	***REMOVED***"input7", "min", "2017-01-01"***REMOVED***,
	***REMOVED***"input7", "max", "2017-12-12"***REMOVED***,
	***REMOVED***"input7", "dirName", "input7.dir"***REMOVED***,
	***REMOVED***"input7", "accessKey", "s"***REMOVED***,
	***REMOVED***"input7", "step", "0.1"***REMOVED***,
	***REMOVED***"kg1", "challenge", "cx1"***REMOVED***,
	***REMOVED***"kg1", "keytype", "DSA"***REMOVED***,
	***REMOVED***"kg1", "name", "kg1_name"***REMOVED***,
	***REMOVED***"kg2", "challenge", ""***REMOVED***,
	***REMOVED***"kg2", "keytype", "RSA"***REMOVED***,
	***REMOVED***"kg2", "type", "keygen"***REMOVED***,
	***REMOVED***"label1", "htmlFor", "input1_name"***REMOVED***,
	***REMOVED***"legend1", "accessKey", "l"***REMOVED***,
	***REMOVED***"li1", "type", "disc"***REMOVED***,
	***REMOVED***"li2", "type", ""***REMOVED***,
	***REMOVED***"link1", "crossOrigin", "use-credentials"***REMOVED***,
	***REMOVED***"link1", "referrerPolicy", "no-referrer"***REMOVED***,
	***REMOVED***"link1", "href", "test.css"***REMOVED***,
	***REMOVED***"link1", "hreflang", "pl"***REMOVED***,
	***REMOVED***"link1", "media", "print"***REMOVED***,
	***REMOVED***"link1", "rel", "alternate author"***REMOVED***,
	***REMOVED***"link1", "target", "__self"***REMOVED***,
	***REMOVED***"link1", "type", "stylesheet"***REMOVED***,
	***REMOVED***"link2", "referrerPolicy", ""***REMOVED***,
	***REMOVED***"map1", "name", "map1_name"***REMOVED***,
	***REMOVED***"meta1", "name", "author"***REMOVED***,
	***REMOVED***"meta1", "content", "author name"***REMOVED***,
	***REMOVED***"meta2", "httpEquiv", "refresh"***REMOVED***,
	***REMOVED***"meta2", "content", "1;www.test.com"***REMOVED***,
	***REMOVED***"meta2", "content", "1;www.test.com"***REMOVED***,
	***REMOVED***"ins1", "cite", "cite.html"***REMOVED***,
	***REMOVED***"ins1", "datetime", "2017-01-01"***REMOVED***,
	***REMOVED***"object1", "data", "test.png"***REMOVED***,
	***REMOVED***"object1", "type", "image/png"***REMOVED***,
	***REMOVED***"object1", "name", "obj1_name"***REMOVED***,
	***REMOVED***"object1", "width", "150"***REMOVED***,
	***REMOVED***"object1", "height", "75"***REMOVED***,
	***REMOVED***"object1", "useMap", "#map1_name"***REMOVED***,
	***REMOVED***"ol1", "type", "a"***REMOVED***,
	***REMOVED***"optgroup1", "label", "optlabel"***REMOVED***,
	***REMOVED***"out1", "htmlFor", "input1"***REMOVED***,
	***REMOVED***"out1", "name", "out1_name"***REMOVED***,
	***REMOVED***"out1", "type", "output"***REMOVED***,
	***REMOVED***"par1", "name", "param1_name"***REMOVED***,
	***REMOVED***"par1", "value", "param1_val"***REMOVED***,
	***REMOVED***"pre1", "name", "pre1_name"***REMOVED***,
	***REMOVED***"pre1", "value", "pre1_val"***REMOVED***,
	***REMOVED***"quote1", "cite", "http://cite.com/url"***REMOVED***,
	***REMOVED***"script1", "crossOrigin", "use-credentials"***REMOVED***,
	***REMOVED***"script1", "type", "text/javascript"***REMOVED***,
	***REMOVED***"script1", "src", "script.js"***REMOVED***,
	***REMOVED***"script1", "charset", "ISO-8859-1"***REMOVED***,
	***REMOVED***"select1", "name", "sel1_name"***REMOVED***,
	***REMOVED***"src1", "keySystem", "keysys"***REMOVED***,
	***REMOVED***"src1", "media", "(min-width: 600px)"***REMOVED***,
	***REMOVED***"src1", "sizes", "100vw,50vw"***REMOVED***,
	***REMOVED***"src1", "srcset", "large.jpg 1024w,medium.jpg 640w"***REMOVED***,
	***REMOVED***"src1", "src", "test.png"***REMOVED***,
	***REMOVED***"src1", "type", "image/png"***REMOVED***,
	***REMOVED***"td1", "headers", "th1"***REMOVED***,
	***REMOVED***"th1", "abbr", "hdr"***REMOVED***,
	***REMOVED***"th1", "scope", "row"***REMOVED***,
	***REMOVED***"txtarea1", "accessKey", "k"***REMOVED***,
	***REMOVED***"txtarea1", "autocomplete", "off"***REMOVED***,
	***REMOVED***"txtarea1", "autocapitalize", "words"***REMOVED***,
	***REMOVED***"txtarea1", "wrap", "hard"***REMOVED***,
	***REMOVED***"txtarea2", "autocomplete", "on"***REMOVED***,
	***REMOVED***"txtarea2", "autocapitalize", "sentences"***REMOVED***,
	***REMOVED***"txtarea2", "wrap", "soft"***REMOVED***,
	***REMOVED***"track1", "kind", "metadata"***REMOVED***,
	***REMOVED***"track1", "src", "foo.en.vtt"***REMOVED***,
	***REMOVED***"track1", "label", "English"***REMOVED***,
	***REMOVED***"track1", "srclang", "en"***REMOVED***,
	***REMOVED***"track2", "kind", "subtitle"***REMOVED***,
	***REMOVED***"track2", "src", "foo.sv.vtt"***REMOVED***,
	***REMOVED***"track2", "srclang", "sv"***REMOVED***,
	***REMOVED***"track2", "label", "Svenska"***REMOVED***,
	***REMOVED***"time1", "datetime", "2017-01-01"***REMOVED***,
	***REMOVED***"ul1", "type", "circle"***REMOVED***,
***REMOVED***

var intTests = []struct ***REMOVED***
	id       string
	property string
	data     int
***REMOVED******REMOVED***
	***REMOVED***"img1", "width", 100***REMOVED***,
	***REMOVED***"img1", "height", 50***REMOVED***,
	***REMOVED***"input7", "maxLength", 10***REMOVED***,
	***REMOVED***"input7", "size", 5***REMOVED***,
	***REMOVED***"li1", "value", 0***REMOVED***,
	***REMOVED***"li2", "value", 10***REMOVED***,
	***REMOVED***"meter1", "min", 90***REMOVED***,
	***REMOVED***"meter1", "max", 110***REMOVED***,
	***REMOVED***"meter1", "low", 95***REMOVED***,
	***REMOVED***"meter1", "high", 105***REMOVED***,
	***REMOVED***"meter1", "optimum", 100***REMOVED***,
	***REMOVED***"object1", "tabIndex", 6***REMOVED***,
	***REMOVED***"ol1", "start", 1***REMOVED***,
	***REMOVED***"td1", "colSpan", 2***REMOVED***,
	***REMOVED***"td1", "rowSpan", 3***REMOVED***,
	***REMOVED***"th1", "colSpan", 1***REMOVED***,
	***REMOVED***"th1", "colSpan", 1***REMOVED***,
	***REMOVED***"txtarea1", "rows", 10***REMOVED***,
	***REMOVED***"txtarea1", "cols", 12***REMOVED***,
	***REMOVED***"txtarea1", "maxLength", 128***REMOVED***,
	***REMOVED***"txtarea1", "tabIndex", 4***REMOVED***,
***REMOVED***

var boolTests = []struct ***REMOVED***
	idTrue   string
	idFalse  string
	property string
***REMOVED******REMOVED***
	***REMOVED***"audio1", "audio2", "autoplay"***REMOVED***,
	***REMOVED***"audio1", "audio2", "controls"***REMOVED***,
	***REMOVED***"audio1", "audio2", "loop"***REMOVED***,
	***REMOVED***"audio1", "audio2", "muted"***REMOVED***,
	***REMOVED***"audio1", "audio2", "defaultMuted"***REMOVED***,
	***REMOVED***"btn1", "btn2", "autofocus"***REMOVED***,
	***REMOVED***"btn1", "btn2", "disabled"***REMOVED***,
	***REMOVED***"fset1", "fset2", "disabled"***REMOVED***,
	***REMOVED***"form1", "form2", "noValidate"***REMOVED***,
	***REMOVED***"iframe1", "iframe2", "allowfullscreen"***REMOVED***,
	***REMOVED***"img1", "img2", "isMap"***REMOVED***,
	***REMOVED***"input1", "input2", "disabled"***REMOVED***,
	***REMOVED***"input1", "input2", "autofocus"***REMOVED***,
	***REMOVED***"input1", "input2", "required"***REMOVED***,
	***REMOVED***"input3", "input4", "checked"***REMOVED***,
	***REMOVED***"input3", "input4", "defaultChecked"***REMOVED***,
	***REMOVED***"input7", "input1", "readonly"***REMOVED***,
	***REMOVED***"input3", "input4", "multiple"***REMOVED***,
	***REMOVED***"kg1", "kg2", "autofocus"***REMOVED***,
	***REMOVED***"kg1", "kg2", "disabled"***REMOVED***,
	***REMOVED***"object1", "object2", "typeMustMatch"***REMOVED***,
	***REMOVED***"ol1", "ol2", "reversed"***REMOVED***,
	***REMOVED***"optgroup1", "optgroup2", "disabled"***REMOVED***,
	***REMOVED***"opt1", "opt2", "selected"***REMOVED***,
	***REMOVED***"opt1", "opt2", "defaultSelected"***REMOVED***,
	***REMOVED***"script1", "script2", "async"***REMOVED***,
	***REMOVED***"script1", "script2", "defer"***REMOVED***,
	***REMOVED***"script1", "script2", "noModule"***REMOVED***,
	***REMOVED***"select1", "select2", "autofocus"***REMOVED***,
	***REMOVED***"select1", "select2", "disabled"***REMOVED***,
	***REMOVED***"select1", "select2", "multiple"***REMOVED***,
	***REMOVED***"select1", "select2", "required"***REMOVED***,
	***REMOVED***"table1", "table2", "sortable"***REMOVED***,

	***REMOVED***"th1", "th2", "sorted"***REMOVED***,

	***REMOVED***"txtarea1", "txtarea2", "readOnly"***REMOVED***,
	***REMOVED***"txtarea1", "txtarea2", "required"***REMOVED***,
***REMOVED***

var nullTests = []struct ***REMOVED***
	id       string
	property string
***REMOVED******REMOVED***
	***REMOVED***"audio2", "crossOrigin"***REMOVED***,
	***REMOVED***"img3", "crossOrigin"***REMOVED***,
	***REMOVED***"link2", "crossOrigin"***REMOVED***,
***REMOVED***

var urlTests = []struct ***REMOVED***
	id       string
	property string
	baseUrl  string
	data     string
***REMOVED******REMOVED***
	***REMOVED***"a2", "href", "http://example.com/testpath", ""***REMOVED***,
	***REMOVED***"a3", "href", "http://example.com", "http://example.com/relpath"***REMOVED***,
	***REMOVED***"a3", "href", "http://example.com/somepath", "http://example.com/relpath"***REMOVED***,
	***REMOVED***"a3", "href", "http://example.com/subdir/", "http://example.com/subdir/relpath"***REMOVED***,
	***REMOVED***"a4", "href", "http://example.com/", "http://example.com/abspath"***REMOVED***,
	***REMOVED***"a4", "href", "http://example.com/subdir/", "http://example.com/abspath"***REMOVED***,
	***REMOVED***"a5", "href", "http://example.com/path?a=no-a&c=no-c", "http://example.com/path?a=yes-a&b=yes-b"***REMOVED***,
	***REMOVED***"a6", "href", "http://example.com/path#oldfrag", "http://example.com/path#testfrag"***REMOVED***,
	***REMOVED***"a7", "href", "http://example.com/prevdir/prevpath", "http://example.com/prtpath"***REMOVED***,
	***REMOVED***"a8", "href", "http://example.com/testpath", "http://example.com/testpath"***REMOVED***,
	***REMOVED***"base1", "href", "http://example.com", "http://example.com/foo.html"***REMOVED***,
	***REMOVED***"base2", "href", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"base3", "href", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"audio1", "src", "http://example.com", "http://example.com/foo.wav"***REMOVED***,
	***REMOVED***"audio2", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"audio3", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"form1", "action", "http://example.com/", "http://example.com/submit_url"***REMOVED***,
	***REMOVED***"form2", "action", "http://example.com/", ""***REMOVED***,
	***REMOVED***"form3", "action", "http://example.com/", "http://example.com/"***REMOVED***,
	***REMOVED***"iframe1", "src", "http://example.com", "http://example.com/testframe.html"***REMOVED***,
	***REMOVED***"iframe2", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"iframe3", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"img1", "src", "http://example.com", "http://example.com/test.png"***REMOVED***,
	***REMOVED***"img2", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"img3", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"input5", "src", "http://example.com", "http://example.com/input.png"***REMOVED***,
	***REMOVED***"input5b", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"input5c", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"link1", "href", "http://example.com", "http://example.com/test.css"***REMOVED***,
	***REMOVED***"link2", "href", "http://example.com", ""***REMOVED***,
	***REMOVED***"link3", "href", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"object1", "data", "http://example.com", "http://example.com/test.png"***REMOVED***,
	***REMOVED***"object2", "data", "http://example.com", ""***REMOVED***,
	***REMOVED***"object3", "data", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"script1", "src", "http://example.com", "http://example.com/script.js"***REMOVED***,
	***REMOVED***"script2", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"script3", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"src1", "src", "http://example.com", "http://example.com/test.png"***REMOVED***,
	***REMOVED***"src2", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"src3", "src", "http://example.com", "http://example.com"***REMOVED***,
	***REMOVED***"track1", "src", "http://example.com", "http://example.com/foo.en.vtt"***REMOVED***,
	***REMOVED***"track3", "src", "http://example.com", ""***REMOVED***,
	***REMOVED***"track4", "src", "http://example.com", "http://example.com"***REMOVED***,
***REMOVED***

const testGenElems = `<html><body>
	<a id="a1" download="file:///path/name" referrerpolicy="no-referrer" rel="open" href="http://test.url" target="__blank" type="text/html" accesskey="w" hreflang="es"></a>
	<a id="a2"></a>
	<a id="a3" href="relpath"></a>
	<a id="a4" href="/abspath"></a>
	<a id="a5" href="?a=yes-a&b=yes-b"></a>
	<a id="a6" href="#testfrag"></a>
	<a id="a7" href="../prtpath"></a>
	<a id="a8" href=""></a>
	<audio id="audio1" autoplay controls loop muted src="foo.wav" crossorigin="anonymous" mediagroup="testgroup"></audio>
	<audio id="audio2"></audio>
	<audio id="audio3" src=""></audio>
	<base id="base1" href="foo.html" target="__any"></base>
	<base id="base2"></base>
	<base id="base3" href="" target="__any"></base>
	<button id="btn1" accesskey="e" target="__any" autofocus disabled type="button"></button>
	<button id="btn2"></button>
	<button id="btn3" type="invalid_uses_default"></button> <button id="btn3" type="invalid_uses_default"></button>
	<ul><li><data id="data1" value="121"></data></li><li><data id="data2"></data></li></ul>
	<embed id="embed1" type="video/avi" src="movie.avi" width="640" height="480">
	<fieldset id="fset1" disabled name="fset1_name"></fieldset>
	<fieldset id="fset2"></fieldset>
	<form id="form1" name="form1_name" target="__self" enctype="text/plain" action="submit_url" accept-charset="ISO-8859-1" autocomplete="off" novalidate></form>
	<form id="form2"></form>
	<form id="form3" action=""></form>
	<iframe id="iframe1" allowfullscreen referrerpolicy="no-referrer" name="frame_name" width="640" height="480" src="testframe.html"></iframe>
	<iframe id="iframe2" referrerpolicy="use-default-when-invalid"></iframe>
	<iframe id="iframe3" src=""></iframe>
	<img id="img1" src="test.png" sizes="100vw,50vw" srcset="large.jpg 1024w,medium.jpg 640w" alt="alt text" crossorigin="anonymous" height="50" width="100" ismap name="img_name" usemap="#map_name" referrerpolicy="origin"/>
	<img id="img2" crossorigin="use-credentials" referrerpolicy="use-default-when-invalid"/>
	<img id="img3" src=""/>
	<input id="input1" name="input1_name" disabled autofocus required value="input1-val" type="button"/>
	<input id="input2"/>
	<input id="input3" type="checkbox" checked multiple/>
	<input id="input4" type="checkbox"/>
	<input id="input5" type="image" alt="input_img" src="input.png" width="80" height="40"/>
	<input id="input5b" type="image" />
	<input id="input5c" type="image" src=""/>
	<input id="input6" type="file" accept=".jpg,.png"/>
	<input id="input7" type="text" autocomplete="off" maxlength="10" size="5" pattern="..." placeholder="help text" readonly min="2017-01-01" max="2017-12-12" dirname="input7.dir" accesskey="s" step="0.1"/>
	<keygen id="kg1" autofocus challenge="cx1" disabled keytype="DSA" name="kg1_name"/>
	<keygen id="kg2"/>
	<label id="label1" for="input1_name"/>
	<legend id="legend1" accesskey="l"/>
	<li id="li1" type="disc"></li> <li id="li2" value="10" type=""></li>
	<link id="link1" crossorigin="use-credentials" referrerpolicy="no-referrer" href="test.css" hreflang="pl" media="print" rel="alternate author" target="__self" type="stylesheet"/>
	<link id="link2"/>
	<link id="link3" href=""/>
	<map id="map1" name="map1_name"></map>
	<meta id="meta1" name="author" content="author name" />
	<meta id="meta2" http-equiv="refresh" content="1;www.test.com" />
	<meter id="meter1" min="90" max="110" low="95" high="105" optimum="100"/>
	<ins id="ins1" cite="cite.html" datetime="2017-01-01"/>
	<object id="object1" name="obj1_name" data="test.png" type="image/png" width="150" height="75" tabindex="6" typemustmatch usemap="#map1_name"/>
	<object id="object2"/>
	<object id="object3" data=""/>
	<ol id="ol1" reversed start="1" type="a"></ol> <ol id="ol2"></ol>
	<optgroup id="optgroup1" disabled label="optlabel"></optgroup>
	<optgroup id="optgroup2"></optgroup>
	<option id="opt1" selected/><option id="opt2" />
	<output id="out1" for="input1" name="out1_name"/>
	<param id="par1" name="param1_name" value="param1_val"/>
	<pre id="pre1" name="pre1_name" value="pre1_val"/>
	<quote id="quote1" cite="http://cite.com/url"/>
	<script id="script1" crossorigin="use-credentials" type="text/javascript" src="script.js" charset="ISO-8859-1" defer async nomodule></script>
	<script id="script2"></script>
	<script id="script3" src=""></script>
	<select id="select1" name="sel1_name" autofocus disabled multiple required></select>
	<select id="select2"></select>
	<source id="src1" keysystem="keysys" media="(min-width: 600px)" sizes="100vw,50vw" srcset="large.jpg 1024w,medium.jpg 640w" src="test.png" type="image/png"></source>
	<source id="src2"></source>
	<source id="src3" src=""></source>
	<style id="style1" media="print"></style>
	<table id="table1" sortable><tr><td id="td1" colspan="2" rowspan="3" headers="th1"></th><th id="th1" abbr="hdr" scope="row" sorted>Header</th><th id="th2"></th></tr></table>
	<table id="table2"></table>
	<textarea id="txtarea1" value="init_txt" placeholder="display_txt" rows="10" cols="12" maxlength="128" accesskey="k" tabIndex="4" readonly required autocomplete="off" autocapitalize="words" wrap="hard"></textarea>
	<textarea id="txtarea2"></textarea>
	<time id="time1" datetime="2017-01-01"/>
	<track id="track1" kind="metadata" src="foo.en.vtt" srclang="en" label="English"></track>
	<track id="track2" src="foo.sv.vtt" srclang="sv" label="Svenska"></track>
	<track id="track3"></track>
	<track id="track4" src=""></track>
	<ul id="ul1" type="circle"/>
	`

func TestGenElements(t *testing.T) ***REMOVED***
	t.Parallel()

	t.Run("Test text properties", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testGenElems)

		for _, test := range textTests ***REMOVED***
			v, err := rt.RunString(`doc.find("#` + test.id + `").get(0).` + test.property + `()`)
			if err != nil ***REMOVED***
				t.Errorf("Error for property name '%s' on element id '#%s':\n%+v ", test.id, test.property, err)
			***REMOVED*** else if v.Export() != test.data ***REMOVED***
				t.Errorf("Expected '%s' for property name '%s' element id '#%s'. Got '%s'", test.data, test.property, test.id, v.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("Test bool properties", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testGenElems)

		for _, test := range boolTests ***REMOVED***
			vT, errT := rt.RunString(`doc.find("#` + test.idTrue + `").get(0).` + test.property + `()`)
			if errT != nil ***REMOVED***
				t.Errorf("Error for property name '%s' on element id '#%s':\n%+v", test.property, test.idTrue, errT)
			***REMOVED*** else if vT.Export() != true ***REMOVED*** // nolint: gosimple
				t.Errorf("Expected true for property name '%s' on element id '#%s'", test.property, test.idTrue)
			***REMOVED***

			vF, errF := rt.RunString(`doc.find("#` + test.idFalse + `").get(0).` + test.property + `()`)
			if errF != nil ***REMOVED***
				t.Errorf("Error for property name '%s' on element id '#%s':\n%+v", test.property, test.idFalse, errF)
			***REMOVED*** else if vF.Export() != false ***REMOVED*** // nolint: gosimple
				t.Errorf("Expected false for property name '%s' on element id '#%s'", test.property, test.idFalse)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("Test int64 properties", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testGenElems)

		for _, test := range intTests ***REMOVED***
			v, err := rt.RunString(`doc.find("#` + test.id + `").get(0).` + test.property + `()`)
			if err != nil ***REMOVED***
				t.Errorf("Error for property name '%s' on element id '#%s':\n%+v", test.property, test.id, err)
			***REMOVED*** else if v.Export() != int64(test.data) ***REMOVED***
				t.Errorf("Expected %d for property name '%s' on element id '#%s'. Got %d", test.data, test.property, test.id, v.ToInteger())
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("Test nullable properties", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt := getTestRuntimeWithDoc(t, testGenElems)

		for _, test := range nullTests ***REMOVED***
			v, err := rt.RunString(`doc.find("#` + test.id + `").get(0).` + test.property + `()`)
			if err != nil ***REMOVED***
				t.Errorf("Error for property name '%s' on element id '#%s':\n%+v", test.property, test.id, err)
			***REMOVED*** else if v.Export() != nil ***REMOVED***
				t.Errorf("Expected null for property name '%s' on element id '#%s'", test.property, test.id)
			***REMOVED***
		***REMOVED***
	***REMOVED***)

	t.Run("Test url properties", func(t *testing.T) ***REMOVED***
		t.Parallel()
		rt, mi := getTestRuntimeAndModuleInstanceWithDoc(t, testGenElems)

		sel, parseError := mi.parseHTML(testGenElems)
		if parseError != nil ***REMOVED***
			t.Errorf("Unable to parse html")
		***REMOVED***

		for _, test := range urlTests ***REMOVED***
			sel.URL = test.baseUrl
			rt.Set("urldoc", sel)

			v, err := rt.RunString(`urldoc.find("#` + test.id + `").get(0).` + test.property + `()`)
			if err != nil ***REMOVED***
				t.Errorf("Error for url property '%s' on element id '#%s':\n%+v", test.property, test.id, err)
			***REMOVED*** else if v.Export() != test.data ***REMOVED***
				t.Errorf("Expected '%s' for property name '%s' on element id '#%s', got '%s'", test.data, test.property, test.id, v.String())
			***REMOVED***
		***REMOVED***
	***REMOVED***)
***REMOVED***
