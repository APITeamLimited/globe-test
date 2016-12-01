package js

import (
	"github.com/kardianos/osext"
	"os"
	"path"
)

var (
	babelDir = "."
	babel    = "babel"
)

func init() ***REMOVED***
	gopath := os.Getenv("GOPATH")
	if gopath != "" ***REMOVED***
		babelDir = path.Join(gopath, "src", "github.com", "loadimpact", "k6", "js")
	***REMOVED*** else if dir, err := osext.ExecutableFolder(); err == nil ***REMOVED***
		babelDir = path.Join(dir, "js")
	***REMOVED***
	babel = path.Join(babelDir, "node_modules", ".bin", babel)
***REMOVED***
