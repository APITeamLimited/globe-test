package js

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHTML(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED*** parseHTML(html); ***REMOVED***
	`))
***REMOVED***

func TestHTMLText(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "k6";
	import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.text() === "This is a test snippet.");
	***REMOVED***
	`))
***REMOVED***

func TestHTMLFindText(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "k6";
	import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.find('#what').text() === "test snippet");
	***REMOVED***
	`))
***REMOVED***

func TestHTMLAddSelector(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "k6";
	import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";
	let html = "<span id='sub'>This</span> is a <span id='obj'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.find('#sub').add('#obj').text() === "Thistest snippet");
	***REMOVED***
	`))
***REMOVED***

func TestHTMLAddSelection(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "k6";
	import ***REMOVED*** parseHTML ***REMOVED*** from "k6/html";
	let html = "<span id='sub'>This</span> is a <span id='obj'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.find('#sub').add(doc.find('#obj')).text() === "Thistest snippet");
	***REMOVED***
	`))
***REMOVED***
