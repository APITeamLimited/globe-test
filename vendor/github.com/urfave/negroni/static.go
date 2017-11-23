package negroni

import (
	"net/http"
	"path"
	"strings"
)

// Static is a middleware handler that serves static files in the given
// directory/filesystem. If the file does not exist on the filesystem, it
// passes along to the next middleware in the chain. If you desire "fileserver"
// type behavior where it returns a 404 for unfound files, you should consider
// using http.FileServer from the Go stdlib.
type Static struct ***REMOVED***
	// Dir is the directory to serve static files from
	Dir http.FileSystem
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
***REMOVED***

// NewStatic returns a new instance of Static
func NewStatic(directory http.FileSystem) *Static ***REMOVED***
	return &Static***REMOVED***
		Dir:       directory,
		Prefix:    "",
		IndexFile: "index.html",
	***REMOVED***
***REMOVED***

func (s *Static) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) ***REMOVED***
	if r.Method != "GET" && r.Method != "HEAD" ***REMOVED***
		next(rw, r)
		return
	***REMOVED***
	file := r.URL.Path
	// if we have a prefix, filter requests by stripping the prefix
	if s.Prefix != "" ***REMOVED***
		if !strings.HasPrefix(file, s.Prefix) ***REMOVED***
			next(rw, r)
			return
		***REMOVED***
		file = file[len(s.Prefix):]
		if file != "" && file[0] != '/' ***REMOVED***
			next(rw, r)
			return
		***REMOVED***
	***REMOVED***
	f, err := s.Dir.Open(file)
	if err != nil ***REMOVED***
		// discard the error?
		next(rw, r)
		return
	***REMOVED***
	defer f.Close()

	fi, err := f.Stat()
	if err != nil ***REMOVED***
		next(rw, r)
		return
	***REMOVED***

	// try to serve index file
	if fi.IsDir() ***REMOVED***
		// redirect if missing trailing slash
		if !strings.HasSuffix(r.URL.Path, "/") ***REMOVED***
			http.Redirect(rw, r, r.URL.Path+"/", http.StatusFound)
			return
		***REMOVED***

		file = path.Join(file, s.IndexFile)
		f, err = s.Dir.Open(file)
		if err != nil ***REMOVED***
			next(rw, r)
			return
		***REMOVED***
		defer f.Close()

		fi, err = f.Stat()
		if err != nil || fi.IsDir() ***REMOVED***
			next(rw, r)
			return
		***REMOVED***
	***REMOVED***

	http.ServeContent(rw, r, file, fi.ModTime(), f)
***REMOVED***
