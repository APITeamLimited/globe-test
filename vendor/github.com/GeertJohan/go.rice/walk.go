package rice

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Walk is like filepath.Walk()
// Visit http://golang.org/pkg/path/filepath/#Walk for more information
func (b *Box) Walk(path string, walkFn filepath.WalkFunc) error ***REMOVED***

	pathFile, err := b.Open(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer pathFile.Close()

	pathInfo, err := pathFile.Stat()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if b.IsAppended() || b.IsEmbedded() ***REMOVED***
		return b.walk(path, pathInfo, walkFn)
	***REMOVED***

	// We don't have any embedded or appended box so use live filesystem mode
	return filepath.Walk(b.absolutePath+string(os.PathSeparator)+path, func(path string, info os.FileInfo, err error) error ***REMOVED***

		// Strip out the box name from the returned paths
		path = strings.TrimPrefix(path, b.absolutePath+string(os.PathSeparator))
		return walkFn(path, info, err)

	***REMOVED***)

***REMOVED***

// walk recursively descends path.
// See walk() in $GOROOT/src/pkg/path/filepath/path.go
func (b *Box) walk(path string, info os.FileInfo, walkFn filepath.WalkFunc) error ***REMOVED***

	err := walkFn(path, info, nil)
	if err != nil ***REMOVED***
		if info.IsDir() && err == filepath.SkipDir ***REMOVED***
			return nil
		***REMOVED***
		return err
	***REMOVED***

	if !info.IsDir() ***REMOVED***
		return nil
	***REMOVED***

	names, err := b.readDirNames(path)
	if err != nil ***REMOVED***
		return walkFn(path, info, err)
	***REMOVED***

	for _, name := range names ***REMOVED***

		filename := filepath.Join(path, name)
		fileObject, err := b.Open(filename)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		defer fileObject.Close()

		fileInfo, err := fileObject.Stat()
		if err != nil ***REMOVED***
			if err := walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = b.walk(filename, fileInfo, walkFn)
			if err != nil ***REMOVED***
				if !fileInfo.IsDir() || err != filepath.SkipDir ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	return nil

***REMOVED***

// readDirNames reads the directory named by path and returns a sorted list of directory entries.
// See readDirNames() in $GOROOT/pkg/path/filepath/path.go
func (b *Box) readDirNames(path string) ([]string, error) ***REMOVED***

	f, err := b.Open(path)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer f.Close()

	stat, err := f.Stat()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !stat.IsDir() ***REMOVED***
		return nil, nil
	***REMOVED***

	infos, err := f.Readdir(0)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var names []string

	for _, info := range infos ***REMOVED***
		names = append(names, info.Name())
	***REMOVED***

	sort.Strings(names)
	return names, nil

***REMOVED***
