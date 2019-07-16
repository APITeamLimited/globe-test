// Copyright 2013 The Go Authors. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

package httputil

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// StaticServer serves static files.
type StaticServer struct ***REMOVED***
	// Dir specifies the location of the directory containing the files to serve.
	Dir string

	// MaxAge specifies the maximum age for the cache control and expiration
	// headers.
	MaxAge time.Duration

	// Error specifies the function used to generate error responses. If Error
	// is nil, then http.Error is used to generate error responses.
	Error Error

	// MIMETypes is a map from file extensions to MIME types.
	MIMETypes map[string]string

	mu    sync.Mutex
	etags map[string]string
***REMOVED***

func (ss *StaticServer) resolve(fname string) string ***REMOVED***
	if path.IsAbs(fname) ***REMOVED***
		panic("Absolute path not allowed when creating a StaticServer handler")
	***REMOVED***
	dir := ss.Dir
	if dir == "" ***REMOVED***
		dir = "."
	***REMOVED***
	fname = filepath.FromSlash(fname)
	return filepath.Join(dir, fname)
***REMOVED***

func (ss *StaticServer) mimeType(fname string) string ***REMOVED***
	ext := path.Ext(fname)
	var mimeType string
	if ss.MIMETypes != nil ***REMOVED***
		mimeType = ss.MIMETypes[ext]
	***REMOVED***
	if mimeType == "" ***REMOVED***
		mimeType = mime.TypeByExtension(ext)
	***REMOVED***
	if mimeType == "" ***REMOVED***
		mimeType = "application/octet-stream"
	***REMOVED***
	return mimeType
***REMOVED***

func (ss *StaticServer) openFile(fname string) (io.ReadCloser, int64, string, error) ***REMOVED***
	f, err := os.Open(fname)
	if err != nil ***REMOVED***
		return nil, 0, "", err
	***REMOVED***
	fi, err := f.Stat()
	if err != nil ***REMOVED***
		f.Close()
		return nil, 0, "", err
	***REMOVED***
	const modeType = os.ModeDir | os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	if fi.Mode()&modeType != 0 ***REMOVED***
		f.Close()
		return nil, 0, "", errors.New("not a regular file")
	***REMOVED***
	return f, fi.Size(), ss.mimeType(fname), nil
***REMOVED***

// FileHandler returns a handler that serves a single file. The file is
// specified by a slash separated path relative to the static server's Dir
// field.
func (ss *StaticServer) FileHandler(fileName string) http.Handler ***REMOVED***
	id := fileName
	fileName = ss.resolve(fileName)
	return &staticHandler***REMOVED***
		ss:   ss,
		id:   func(_ string) string ***REMOVED*** return id ***REMOVED***,
		open: func(_ string) (io.ReadCloser, int64, string, error) ***REMOVED*** return ss.openFile(fileName) ***REMOVED***,
	***REMOVED***
***REMOVED***

// DirectoryHandler returns a handler that serves files from a directory tree.
// The directory is specified by a slash separated path relative to the static
// server's Dir field.
func (ss *StaticServer) DirectoryHandler(prefix, dirName string) http.Handler ***REMOVED***
	if !strings.HasSuffix(prefix, "/") ***REMOVED***
		prefix += "/"
	***REMOVED***
	idBase := dirName
	dirName = ss.resolve(dirName)
	return &staticHandler***REMOVED***
		ss: ss,
		id: func(p string) string ***REMOVED***
			if !strings.HasPrefix(p, prefix) ***REMOVED***
				return "."
			***REMOVED***
			return path.Join(idBase, p[len(prefix):])
		***REMOVED***,
		open: func(p string) (io.ReadCloser, int64, string, error) ***REMOVED***
			if !strings.HasPrefix(p, prefix) ***REMOVED***
				return nil, 0, "", errors.New("request url does not match directory prefix")
			***REMOVED***
			p = p[len(prefix):]
			return ss.openFile(filepath.Join(dirName, filepath.FromSlash(p)))
		***REMOVED***,
	***REMOVED***
***REMOVED***

// FilesHandler returns a handler that serves the concatentation of the
// specified files. The files are specified by slash separated paths relative
// to the static server's Dir field.
func (ss *StaticServer) FilesHandler(fileNames ...string) http.Handler ***REMOVED***

	// todo: cache concatenated files on disk and serve from there.

	mimeType := ss.mimeType(fileNames[0])
	var buf []byte
	var openErr error

	for _, fileName := range fileNames ***REMOVED***
		p, err := ioutil.ReadFile(ss.resolve(fileName))
		if err != nil ***REMOVED***
			openErr = err
			buf = nil
			break
		***REMOVED***
		buf = append(buf, p...)
	***REMOVED***

	id := strings.Join(fileNames, " ")

	return &staticHandler***REMOVED***
		ss: ss,
		id: func(_ string) string ***REMOVED*** return id ***REMOVED***,
		open: func(p string) (io.ReadCloser, int64, string, error) ***REMOVED***
			return ioutil.NopCloser(bytes.NewReader(buf)), int64(len(buf)), mimeType, openErr
		***REMOVED***,
	***REMOVED***
***REMOVED***

type staticHandler struct ***REMOVED***
	id   func(fname string) string
	open func(p string) (io.ReadCloser, int64, string, error)
	ss   *StaticServer
***REMOVED***

func (h *staticHandler) error(w http.ResponseWriter, r *http.Request, status int, err error) ***REMOVED***
	http.Error(w, http.StatusText(status), status)
***REMOVED***

func (h *staticHandler) etag(p string) (string, error) ***REMOVED***
	id := h.id(p)

	h.ss.mu.Lock()
	if h.ss.etags == nil ***REMOVED***
		h.ss.etags = make(map[string]string)
	***REMOVED***
	etag := h.ss.etags[id]
	h.ss.mu.Unlock()

	if etag != "" ***REMOVED***
		return etag, nil
	***REMOVED***

	// todo: if a concurrent goroutine is calculating the hash, then wait for
	// it instead of computing it again here.

	rc, _, _, err := h.open(p)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	defer rc.Close()

	w := sha1.New()
	_, err = io.Copy(w, rc)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	etag = fmt.Sprintf(`"%x"`, w.Sum(nil))

	h.ss.mu.Lock()
	h.ss.etags[id] = etag
	h.ss.mu.Unlock()

	return etag, nil
***REMOVED***

func (h *staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) ***REMOVED***
	p := path.Clean(r.URL.Path)
	if p != r.URL.Path ***REMOVED***
		http.Redirect(w, r, p, 301)
		return
	***REMOVED***

	etag, err := h.etag(p)
	if err != nil ***REMOVED***
		h.error(w, r, http.StatusNotFound, err)
		return
	***REMOVED***

	maxAge := h.ss.MaxAge
	if maxAge == 0 ***REMOVED***
		maxAge = 24 * time.Hour
	***REMOVED***
	if r.FormValue("v") != "" ***REMOVED***
		maxAge = 365 * 24 * time.Hour
	***REMOVED***

	cacheControl := fmt.Sprintf("public, max-age=%d", maxAge/time.Second)

	for _, e := range header.ParseList(r.Header, "If-None-Match") ***REMOVED***
		if e == etag ***REMOVED***
			w.Header().Set("Cache-Control", cacheControl)
			w.Header().Set("Etag", etag)
			w.WriteHeader(http.StatusNotModified)
			return
		***REMOVED***
	***REMOVED***

	rc, cl, ct, err := h.open(p)
	if err != nil ***REMOVED***
		h.error(w, r, http.StatusNotFound, err)
		return
	***REMOVED***
	defer rc.Close()

	w.Header().Set("Cache-Control", cacheControl)
	w.Header().Set("Etag", etag)
	if ct != "" ***REMOVED***
		w.Header().Set("Content-Type", ct)
	***REMOVED***
	if cl != 0 ***REMOVED***
		w.Header().Set("Content-Length", strconv.FormatInt(cl, 10))
	***REMOVED***
	w.WriteHeader(http.StatusOK)
	if r.Method != "HEAD" ***REMOVED***
		io.Copy(w, rc)
	***REMOVED***
***REMOVED***
