package gin

import (
	"net/http"
	"os"
)

type (
	onlyfilesFS struct ***REMOVED***
		fs http.FileSystem
	***REMOVED***
	neuteredReaddirFile struct ***REMOVED***
		http.File
	***REMOVED***
)

// Dir returns a http.Filesystem that can be used by http.FileServer(). It is used internally
// in router.Static().
// if listDirectory == true, then it works the same as http.Dir() otherwise it returns
// a filesystem that prevents http.FileServer() to list the directory files.
func Dir(root string, listDirectory bool) http.FileSystem ***REMOVED***
	fs := http.Dir(root)
	if listDirectory ***REMOVED***
		return fs
	***REMOVED***
	return &onlyfilesFS***REMOVED***fs***REMOVED***
***REMOVED***

// Conforms to http.Filesystem
func (fs onlyfilesFS) Open(name string) (http.File, error) ***REMOVED***
	f, err := fs.fs.Open(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return neuteredReaddirFile***REMOVED***f***REMOVED***, nil
***REMOVED***

// Overrides the http.File default implementation
func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) ***REMOVED***
	// this disables directory listing
	return nil, nil
***REMOVED***
