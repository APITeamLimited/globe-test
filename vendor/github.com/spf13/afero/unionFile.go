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
	Base   File
	Layer  File
	Merger DirsMerger
	off    int
	files  []os.FileInfo
***REMOVED***

func (f *UnionFile) Close() error ***REMOVED***
	// first close base, so we have a newer timestamp in the overlay. If we'd close
	// the overlay first, we'd get a cacheStale the next time we access this file
	// -> cache would be useless ;-)
	if f.Base != nil ***REMOVED***
		f.Base.Close()
	***REMOVED***
	if f.Layer != nil ***REMOVED***
		return f.Layer.Close()
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) Read(s []byte) (int, error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		n, err := f.Layer.Read(s)
		if (err == nil || err == io.EOF) && f.Base != nil ***REMOVED***
			// advance the file position also in the base file, the next
			// call may be a write at this position (or a seek with SEEK_CUR)
			if _, seekErr := f.Base.Seek(int64(n), os.SEEK_CUR); seekErr != nil ***REMOVED***
				// only overwrite err in case the seek fails: we need to
				// report an eventual io.EOF to the caller
				err = seekErr
			***REMOVED***
		***REMOVED***
		return n, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Read(s)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) ReadAt(s []byte, o int64) (int, error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		n, err := f.Layer.ReadAt(s, o)
		if (err == nil || err == io.EOF) && f.Base != nil ***REMOVED***
			_, err = f.Base.Seek(o+int64(n), os.SEEK_SET)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.ReadAt(s, o)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Seek(o int64, w int) (pos int64, err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		pos, err = f.Layer.Seek(o, w)
		if (err == nil || err == io.EOF) && f.Base != nil ***REMOVED***
			_, err = f.Base.Seek(o, w)
		***REMOVED***
		return pos, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Seek(o, w)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Write(s []byte) (n int, err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		n, err = f.Layer.Write(s)
		if err == nil && f.Base != nil ***REMOVED*** // hmm, do we have fixed size files where a write may hit the EOF mark?
			_, err = f.Base.Write(s)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Write(s)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) WriteAt(s []byte, o int64) (n int, err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		n, err = f.Layer.WriteAt(s, o)
		if err == nil && f.Base != nil ***REMOVED***
			_, err = f.Base.WriteAt(s, o)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.WriteAt(s, o)
	***REMOVED***
	return 0, BADFD
***REMOVED***

func (f *UnionFile) Name() string ***REMOVED***
	if f.Layer != nil ***REMOVED***
		return f.Layer.Name()
	***REMOVED***
	return f.Base.Name()
***REMOVED***

// DirsMerger is how UnionFile weaves two directories together.
// It takes the FileInfo slices from the layer and the base and returns a
// single view.
type DirsMerger func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error)

var defaultUnionMergeDirsFn = func(lofi, bofi []os.FileInfo) ([]os.FileInfo, error) ***REMOVED***
	var files = make(map[string]os.FileInfo)

	for _, fi := range lofi ***REMOVED***
		files[fi.Name()] = fi
	***REMOVED***

	for _, fi := range bofi ***REMOVED***
		if _, exists := files[fi.Name()]; !exists ***REMOVED***
			files[fi.Name()] = fi
		***REMOVED***
	***REMOVED***

	rfi := make([]os.FileInfo, len(files))

	i := 0
	for _, fi := range files ***REMOVED***
		rfi[i] = fi
		i++
	***REMOVED***

	return rfi, nil

***REMOVED***

// Readdir will weave the two directories together and
// return a single view of the overlayed directories
func (f *UnionFile) Readdir(c int) (ofi []os.FileInfo, err error) ***REMOVED***
	var merge DirsMerger = f.Merger
	if merge == nil ***REMOVED***
		merge = defaultUnionMergeDirsFn
	***REMOVED***

	if f.off == 0 ***REMOVED***
		var lfi []os.FileInfo
		if f.Layer != nil ***REMOVED***
			lfi, err = f.Layer.Readdir(-1)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***
		***REMOVED***

		var bfi []os.FileInfo
		if f.Base != nil ***REMOVED***
			bfi, err = f.Base.Readdir(-1)
			if err != nil ***REMOVED***
				return nil, err
			***REMOVED***

		***REMOVED***
		merged, err := merge(lfi, bfi)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		f.files = append(f.files, merged...)
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
	if f.Layer != nil ***REMOVED***
		return f.Layer.Stat()
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Stat()
	***REMOVED***
	return nil, BADFD
***REMOVED***

func (f *UnionFile) Sync() (err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		err = f.Layer.Sync()
		if err == nil && f.Base != nil ***REMOVED***
			err = f.Base.Sync()
		***REMOVED***
		return err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Sync()
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) Truncate(s int64) (err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		err = f.Layer.Truncate(s)
		if err == nil && f.Base != nil ***REMOVED***
			err = f.Base.Truncate(s)
		***REMOVED***
		return err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.Truncate(s)
	***REMOVED***
	return BADFD
***REMOVED***

func (f *UnionFile) WriteString(s string) (n int, err error) ***REMOVED***
	if f.Layer != nil ***REMOVED***
		n, err = f.Layer.WriteString(s)
		if err == nil && f.Base != nil ***REMOVED***
			_, err = f.Base.WriteString(s)
		***REMOVED***
		return n, err
	***REMOVED***
	if f.Base != nil ***REMOVED***
		return f.Base.WriteString(s)
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
