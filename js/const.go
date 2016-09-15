package js

import (
	"github.com/kardianos/osext"
	"path"
)

var (
	babelDir = "."
	babel    = "babel"
)

func init() ***REMOVED***
	if dir, err := osext.ExecutableFolder(); err == nil ***REMOVED***
		babelDir = path.Join(dir, "js")
		babel = path.Join(babelDir, "node_modules", ".bin", babel)
	***REMOVED***
***REMOVED***
