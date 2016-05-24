package js

import (
	"gopkg.in/olebedev/go-duktape.v2"
)

func argNumber(c *duktape.Context, index int) float64 ***REMOVED***
	if c.GetTopIndex() < index ***REMOVED***
		return 0
	***REMOVED***

	return c.ToNumber(index)
***REMOVED***

func argString(c *duktape.Context, index int) string ***REMOVED***
	if c.GetTopIndex() < index ***REMOVED***
		return ""
	***REMOVED***

	return c.ToString(index)
***REMOVED***
