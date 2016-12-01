/**
 * @module k6/html
 */

/**
 * Parses an HTML string into a Selection.
 *
 * @param  ***REMOVED***string***REMOVED***    src HTML source.
 * @return ***REMOVED***Selection***REMOVED***
 * @throws ***REMOVED***Error***REMOVED***         If src is not valid HTML.
 */
export function parseHTML(src) ***REMOVED***
	return new Selection(__jsapi__.HTMLParse(src));
***REMOVED***;

export class Selection ***REMOVED***
	/**
	 * Represents a set of nodes in a DOM tree.
	 *
	 * Selections have a jQuery-compatible API, but with two caveats:
	 *
	 * - CSS and screen layout are not processed, thus calls like css() and offset() are unavailable.
	 * - DOM trees are read-only, you can't set attributes or otherwise modify nodes.
	 *
	 * (Note that the read-only nature of the DOM trees is purely to avoid a maintenance burden on code
	 * with seemingly no practical use - if a compelling use case is presented, modification can
	 * easily be implemented.)
	 *
	 * @memberOf module:k6/html
	 */
	constructor(impl) ***REMOVED***
		this.impl = impl;
	***REMOVED***

	/**
	 * Extends the selection with another set of elements.
	 *
	 * @param ***REMOVED***string|Selection***REMOVED*** arg Selection or selector
	 * @return ***REMOVED***module:k6/html.Selection***REMOVED***
	 */
	add(arg) ***REMOVED***
		if (typeof arg === "string") ***REMOVED***
			return new Selection(this.impl.Add(arg));
		***REMOVED*** else if (arg instanceof Selection) ***REMOVED***
			return new Selection(__jsapi__.HTMLSelectionAddSelection(this.impl, arg.impl));
		***REMOVED***
		throw new TypeError("add() argument must be a string or Selection")
	***REMOVED***

	/**
	 * Finds children by a selector.
	 *
	 * @param  ***REMOVED***string***REMOVED***    sel CSS selector.
	 * @return ***REMOVED***module:k6/html.Selection***REMOVED***
	 */
	find(sel) ***REMOVED***
		return new Selection(this.impl.Find(sel));
	***REMOVED***

	/**
	 * Returns the combined text content of all selected nodes.
	 * @return ***REMOVED***string***REMOVED***
	 */
	text() ***REMOVED*** return this.impl.Text(); ***REMOVED***
***REMOVED***;

export default ***REMOVED***
	parseHTML: parseHTML,
	Selection: Selection,
***REMOVED***
