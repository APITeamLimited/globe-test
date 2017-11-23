//+build go1.8,!openbsd

package osext

import "os"

func executable() (string, error) ***REMOVED***
	return os.Executable()
***REMOVED***
