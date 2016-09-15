package js

import (
	"fmt"
)

type DirectoryTraversalError struct ***REMOVED***
	Filename string
	Root     string
***REMOVED***

func (e DirectoryTraversalError) Error() string ***REMOVED***
	return fmt.Sprintf("loading files outside your working directory is prohibited, to protect against directory traversal attacks (%s is outside %s)", e.Filename, e.Root)
***REMOVED***
