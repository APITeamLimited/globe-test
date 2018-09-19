//+build appengine

package gjson

func getBytes(json []byte, path string) Result ***REMOVED***
	return Get(string(json), path)
***REMOVED***
func fillIndex(json string, c *parseContext) ***REMOVED***
	// noop. Use zero for the Index value.
***REMOVED***
