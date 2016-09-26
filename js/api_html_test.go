package js

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHTML(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** parseHTML ***REMOVED*** from "speedboat/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED*** parseHTML(html); ***REMOVED***
	`))
***REMOVED***

func TestHTMLText(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "speedboat";
	import ***REMOVED*** parseHTML ***REMOVED*** from "speedboat/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.text() === "This is a test snippet.");
	***REMOVED***
	`))
***REMOVED***

func TestHTMLFindText(t *testing.T) ***REMOVED***
	assert.NoError(t, runSnippet(`
	import ***REMOVED*** _assert ***REMOVED*** from "speedboat";
	import ***REMOVED*** parseHTML ***REMOVED*** from "speedboat/html";
	let html = "This is a <span id='what'>test snippet</span>.";
	export default function() ***REMOVED***
		let doc = parseHTML(html);
		_assert(doc.find('#what').text() === "test snippet");
	***REMOVED***
	`))
***REMOVED***
