package js

import (
	"encoding/json"
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

func argJSON(c *duktape.Context, index int, out interface***REMOVED******REMOVED***) error ***REMOVED***
	if c.GetTopIndex() < index ***REMOVED***
		return nil
	***REMOVED***

	c.JsonEncode(index)
	str := c.GetString(index)
	return json.Unmarshal([]byte(str), out)
***REMOVED***
