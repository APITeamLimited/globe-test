package afero

import (
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// The UnionFile implements the afero.File interface and will be returned
// when reading a directory present at least in the overlay or opening a file
// for writing.
//
// The calls to
// Readdir() and Readdirnames() merge the file os.FileInfo / names from the
// base and the overlay - for files present in both layers, only those
// from the overlay will be used.
//
// When opening files for writing (Create() / OpenFile() with the right flags)
// the operations will be done in both layers, starting with the overlay. A
// successful read in the overlay will move the cursor position in the base layer
// by the number of bytes read.
type UnionFile struct ***REMOVED***
	base  File
	layer File
	off   int
	files []os.FileInfo
***REMOVED***

func (f *UnionFile) Close() error ***REMOVED***
	// first close base, so we have a newer timestamp in the overlay. If we'd close
	// the overlay first, we'd get a cacheStale the next time we access this file
	// -> cache would be useless ;-)
	if f.base != nil ***REMOVED***
		f.base.Close()
	***REMOVED***
	if f.layer != nil ***REMOVED***
		return f.layer.Close()
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) Read(s []byte) (int, error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		n, err := f.layer.Read(s)
		if (err == nil || err == io.EOF) && f.base != nil ***REMOVED***
			// advance the file position also in the base file, the next
			// call may be a write at this position (or a seek with SEEK_CUR)
			if _, seekErr := f.base.Seek(int64(n), os.SEEK_CUR); seekErr != nil ***REMOVED***
				// only overwrite err in case the seek fails: we need to
				// report an eventual io.EOF to the caller
				err = seekErr
			***REMOVED***
		***REMOVED***
		return n, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Read(s)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) ReadAt(s []byte, o int64) (int, error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		n, err := f.layer.ReadAt(s, o)
		if (err == nil || err == io.EOF) && f.base != nil ***REMOVED***
			_, err = f.base.Seek(o+int64(n), os.SEEK_SET)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.ReadAt(s, o)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Seek(o int64, w int) (pos int64, err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		pos, err = f.layer.Seek(o, w)
		if (err == nil || err == io.EOF) && f.base != nil ***REMOVED***
			_, err = f.base.Seek(o, w)
		***REMOVED***
		return pos, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Seek(o, w)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Write(s []byte) (n int, err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		n, err = f.layer.Write(s)
		if err == nil && f.base != nil ***REMOVED*** // hmm, do we have fixed size files where a write may hit the EOF mark?
			_, err = f.base.Write(s)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Write(s)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) WriteAt(s []byte, o int64) (n int, err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		n, err = f.layer.WriteAt(s, o)
		if err == nil && f.base != nil ***REMOVED***
			_, err = f.base.WriteAt(s, o)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.WriteAt(s, o)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Name() string ***REMOVED***
	if f.layer != nil ***REMOVED***
		return f.layer.Name()
	***REMOVED***
	return f.base.Name()
***REMOVED***

// Readdir will weave the two directories together and
// return a single view of the overlayed directories
func (f *UnionFile) Readdir(c int) (ofi []os.FileInfo, err error) ***REMOVED***
	if f.off == 0 ***REMOVED***
		var files = make(map[string]os.FileInfo)
		var rfi []os.FileInfo
		if f.layer != nil ***REMOVED***
			rfi, err = f.layer.Readdir(-1)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			for _, fi := range rfi ***REMOVED***
				files[fi.Name()] = fi
			***REMOVED***
		***REMOVED***

		if f.base != nil ***REMOVED***
			rfi, err = f.base.Readdir(-1)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
			for _, fi := range rfi ***REMOVED***
				if _, exists := files[fi.Name()]; !exists ***REMOVED***
					files[fi.Name()] = fi
				***REMOVED***
			***REMOVED***
		***REMOVED***
		for _, fi := range files ***REMOVED***
			f.files = append(f.files, fi)
		***REMOVED***
	***REMOVED***
	if c == -1 ***REMOVED***
		return f.files[f.off:], nil
	***REMOVED***
	defer func() ***REMOVED*** f.off += c ***REMOVED***()
	return f.files[f.off:c], nil
***REMOVED***

func (f *UnionFile) Readdirnames(c int) ([]string, error) ***REMOVED***
	rfi, err := f.Readdir(c)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	var names []string
	for _, fi := range rfi ***REMOVED***
		names = append(names, fi.Name())
	***REMOVED***
	return names, nil
***REMOVED***

func (f *UnionFile) Stat() (os.FileInfo, error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		return f.layer.Stat()
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Stat()
	***REMOVED***
	return nil, BADFD
***REMOVED***

func (f *UnionFile) Sync() (err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		err = f.layer.Sync()
		if err == nil && f.base != nil ***REMOVED***
			err = f.base.Sync()
		***REMOVED***
		return err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Sync()
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) Truncate(s int64) (err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		err = f.layer.Truncate(s)
		if err == nil && f.base != nil ***REMOVED***
			err = f.base.Truncate(s)
		***REMOVED***
		return err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.Truncate(s)
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) WriteString(s string) (n int, err error) ***REMOVED***
	if f.layer != nil ***REMOVED***
		n, err = f.layer.WriteString(s)
		if err == nil && f.base != nil ***REMOVED***
			_, err = f.base.WriteString(s)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.base != nil ***REMOVED***
		return f.base.WriteString(s)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func copyToLayer(base Fs, layer Fs, name string) error ***REMOVED***
	bfh, err := base.Open(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer bfh.Close()

	// First make sure the directory exists
	exists, err := Exists(layer, filepath.Dir(name))
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if !exists ***REMOVED***
		err = layer.MkdirAll(filepath.Dir(name), 0777) // FIXME?
		if err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	// Create the file on the overlay
	lfh, err := layer.Create(name)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	n, err := io.Copy(lfh, bfh)
	if err != nil ***REMOVED***
		// If anything fails, clean up the file
		layer.Remove(name)
		lfh.Close()
		return err
	***REMOVED***

	bfi, err := bfh.Stat()
	if err != nil || bfi.Size() != n ***REMOVED***
		layer.Remove(name)
		lfh.Close()
		return syscall.EIO
	***REMOVED***

	err = lfh.Close()
	if err != nil ***REMOVED***
		layer.Remove(name)
		lfh.Close()
		return err
	***REMOVED***
	return layer.Chtimes(name, bfi.ModTime(), bfi.ModTime())
***REMOVED***
