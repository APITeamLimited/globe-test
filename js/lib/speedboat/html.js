export class Node ***REMOVED***
	constructor(impl) ***REMOVED***
		this.impl = impl;
	***REMOVED***
***REMOVED***

export class Selection ***REMOVED***
	constructor(impl) ***REMOVED***
		this.impl = impl;
	***REMOVED***

	find(sel) ***REMOVED*** return new Selection(this.impl.Find(sel)); ***REMOVED***
	text() ***REMOVED*** return this.impl.Text(); ***REMOVED***
***REMOVED***;

export function parseHTML(src) ***REMOVED***
	return new Selection(__jsapi__.HTMLParse(src));
***REMOVED***;
