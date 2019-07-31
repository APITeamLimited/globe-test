package fsext

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/afero"
)

// Walk implements afero.Walk, but in a way that it doesn't loop to infinity and doesn't have
// problems if a given path part looks like a windows volume name
func Walk(fs afero.Fs, root string, walkFn filepath.WalkFunc) error ***REMOVED***
	info, err := fs.Stat(root)
	if err != nil ***REMOVED***
		return walkFn(root, nil, err)
	***REMOVED***
	return walk(fs, root, info, walkFn)
***REMOVED***

// readDirNames reads the directory named by dirname and returns
// a sorted list of directory entries.
// adapted from https://github.com/spf13/afero/blob/master/path.go#L27
func readDirNames(fs afero.Fs, dirname string) ([]string, error) ***REMOVED***
	f, err := fs.Open(dirname)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	infos, err := f.Readdir(-1)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	err = f.Close()

	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var names = make([]string, len(infos))
	for i, info := range infos ***REMOVED***
		names[i] = info.Name()
	***REMOVED***
	sort.Strings(names)
	return names, nil
***REMOVED***

// walk recursively descends path, calling walkFn
// adapted from https://github.com/spf13/afero/blob/master/path.go#L27
func walk(fs afero.Fs, path string, info os.FileInfo, walkFn filepath.WalkFunc) error ***REMOVED***
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

	names, err := readDirNames(fs, path)
	if err != nil ***REMOVED***
		return walkFn(path, info, err)
	***REMOVED***

	for _, name := range names ***REMOVED***
		filename := filepath.Join(path, name)
		fileInfo, err := fs.Stat(filename)
		if err != nil ***REMOVED***
			if err = walkFn(filename, fileInfo, err); err != nil && err != filepath.SkipDir ***REMOVED***
				return err
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			err = walk(fs, filename, fileInfo, walkFn)
			if err != nil ***REMOVED***
				if !fileInfo.IsDir() || err != filepath.SkipDir ***REMOVED***
					return err
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***
