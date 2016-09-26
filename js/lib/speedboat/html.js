export class Node ***REMOVED***
	constructor(impl) ***REMOVED***
		this.impl = impl;
	***REMOVED***
***REMOVED***

export class Selection ***REMOVED***
	constructor(impl) ***REMOVED***
		this.impl = impl;
	***REMOVED***

	add(arg) ***REMOVED***
		if (typeof arg === "string") ***REMOVED***
			return new Selection(this.impl.Add(arg));
		***REMOVED*** else if (arg instanceof Selection) ***REMOVED***
			return new Selection(__jsapi__.HTMLSelectionAddSelection(this.impl, arg.impl));
		***REMOVED***
		throw new TypeError("add() argument must be a string or Selection")
	***REMOVED***

	find(sel) ***REMOVED*** return new Selection(this.impl.Find(sel)); ***REMOVED***
	text() ***REMOVED*** return this.impl.Text(); ***REMOVED***
***REMOVED***;

export function parseHTML(src) ***REMOVED***
	return new Selection(__jsapi__.HTMLParse(src));
***REMOVED***;
