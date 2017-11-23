package rice

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/GeertJohan/go.rice/embedded"
)

//++ TODO: IDEA: merge virtualFile and virtualDir, this decreases work done by rice.File

// Error indicating some function is not implemented yet (but available to satisfy an interface)
var ErrNotImplemented = errors.New("not implemented yet")

// virtualFile is a 'stateful' virtual file.
// virtualFile wraps an *EmbeddedFile for a call to Box.Open() and virtualizes 'read cursor' (offset) and 'closing'.
// virtualFile is only internally visible and should be exposed through rice.File
type virtualFile struct ***REMOVED***
	*embedded.EmbeddedFile       // the actual embedded file, embedded to obtain methods
	offset                 int64 // read position on the virtual file
	closed                 bool  // closed when true
***REMOVED***

// create a new virtualFile for given EmbeddedFile
func newVirtualFile(ef *embedded.EmbeddedFile) *virtualFile ***REMOVED***
	vf := &virtualFile***REMOVED***
		EmbeddedFile: ef,
		offset:       0,
		closed:       false,
	***REMOVED***
	return vf
***REMOVED***

//++ TODO check for nil pointers in all these methods. When so: return os.PathError with Err: os.ErrInvalid

func (vf *virtualFile) close() error ***REMOVED***
	if vf.closed ***REMOVED***
		return &os.PathError***REMOVED***
			Op:   "close",
			Path: vf.EmbeddedFile.Filename,
			Err:  errors.New("already closed"),
		***REMOVED***
	***REMOVED***
	vf.EmbeddedFile = nil
	vf.closed = true
	return nil
***REMOVED***

func (vf *virtualFile) stat() (os.FileInfo, error) ***REMOVED***
	if vf.closed ***REMOVED***
		return nil, &os.PathError***REMOVED***
			Op:   "stat",
			Path: vf.EmbeddedFile.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	return (*embeddedFileInfo)(vf.EmbeddedFile), nil
***REMOVED***

func (vf *virtualFile) readdir(count int) ([]os.FileInfo, error) ***REMOVED***
	if vf.closed ***REMOVED***
		return nil, &os.PathError***REMOVED***
			Op:   "readdir",
			Path: vf.EmbeddedFile.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	//TODO: return proper error for a readdir() call on a file
	return nil, ErrNotImplemented
***REMOVED***

func (vf *virtualFile) read(bts []byte) (int, error) ***REMOVED***
	if vf.closed ***REMOVED***
		return 0, &os.PathError***REMOVED***
			Op:   "read",
			Path: vf.EmbeddedFile.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***

	end := vf.offset + int64(len(bts))

	if end >= int64(len(vf.Content)) ***REMOVED***
		// end of file, so return what we have + EOF
		n := copy(bts, vf.Content[vf.offset:])
		vf.offset = 0
		return n, io.EOF
	***REMOVED***

	n := copy(bts, vf.Content[vf.offset:end])
	vf.offset += int64(n)
	return n, nil

***REMOVED***

func (vf *virtualFile) seek(offset int64, whence int) (int64, error) ***REMOVED***
	if vf.closed ***REMOVED***
		return 0, &os.PathError***REMOVED***
			Op:   "seek",
			Path: vf.EmbeddedFile.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	var e error

	//++ TODO: check if this is correct implementation for seek
	switch whence ***REMOVED***
	case os.SEEK_SET:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset = offset
	case os.SEEK_CUR:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset += offset
	case os.SEEK_END:
		//++ check if new offset isn't out of bounds, set e when it is, then break out of switch
		vf.offset = int64(len(vf.EmbeddedFile.Content)) - offset
	***REMOVED***

	if e != nil ***REMOVED***
		return 0, &os.PathError***REMOVED***
			Op:   "seek",
			Path: vf.Filename,
			Err:  e,
		***REMOVED***
	***REMOVED***

	return vf.offset, nil
***REMOVED***

// virtualDir is a 'stateful' virtual directory.
// virtualDir wraps an *EmbeddedDir for a call to Box.Open() and virtualizes 'closing'.
// virtualDir is only internally visible and should be exposed through rice.File
type virtualDir struct ***REMOVED***
	*embedded.EmbeddedDir
	offset int // readdir position on the directory
	closed bool
***REMOVED***

// create a new virtualDir for given EmbeddedDir
func newVirtualDir(ed *embedded.EmbeddedDir) *virtualDir ***REMOVED***
	vd := &virtualDir***REMOVED***
		EmbeddedDir: ed,
		offset:      0,
		closed:      false,
	***REMOVED***
	return vd
***REMOVED***

func (vd *virtualDir) close() error ***REMOVED***
	//++ TODO: needs sync mutex?
	if vd.closed ***REMOVED***
		return &os.PathError***REMOVED***
			Op:   "close",
			Path: vd.EmbeddedDir.Filename,
			Err:  errors.New("already closed"),
		***REMOVED***
	***REMOVED***
	vd.closed = true
	return nil
***REMOVED***

func (vd *virtualDir) stat() (os.FileInfo, error) ***REMOVED***
	if vd.closed ***REMOVED***
		return nil, &os.PathError***REMOVED***
			Op:   "stat",
			Path: vd.EmbeddedDir.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	return (*embeddedDirInfo)(vd.EmbeddedDir), nil
***REMOVED***

func (vd *virtualDir) readdir(n int) (fi []os.FileInfo, err error) ***REMOVED***

	if vd.closed ***REMOVED***
		return nil, &os.PathError***REMOVED***
			Op:   "readdir",
			Path: vd.EmbeddedDir.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***

	// Build up the array of our contents
	var files []os.FileInfo

	// Add the child directories
	for _, child := range vd.ChildDirs ***REMOVED***
		child.Filename = filepath.Base(child.Filename)
		files = append(files, (*embeddedDirInfo)(child))
	***REMOVED***

	// Add the child files
	for _, child := range vd.ChildFiles ***REMOVED***
		child.Filename = filepath.Base(child.Filename)
		files = append(files, (*embeddedFileInfo)(child))
	***REMOVED***

	// Sort it by filename (lexical order)
	sort.Sort(SortByName(files))

	// Return all contents if that's what is requested
	if n <= 0 ***REMOVED***
		vd.offset = 0
		return files, nil
	***REMOVED***

	// If user has requested past the end of our list
	// return what we can and send an EOF
	if vd.offset+n >= len(files) ***REMOVED***
		offset := vd.offset
		vd.offset = 0
		return files[offset:], io.EOF
	***REMOVED***

	offset := vd.offset
	vd.offset += n
	return files[offset : offset+n], nil

***REMOVED***

func (vd *virtualDir) read(bts []byte) (int, error) ***REMOVED***
	if vd.closed ***REMOVED***
		return 0, &os.PathError***REMOVED***
			Op:   "read",
			Path: vd.EmbeddedDir.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	return 0, &os.PathError***REMOVED***
		Op:   "read",
		Path: vd.EmbeddedDir.Filename,
		Err:  errors.New("is a directory"),
	***REMOVED***
***REMOVED***

func (vd *virtualDir) seek(offset int64, whence int) (int64, error) ***REMOVED***
	if vd.closed ***REMOVED***
		return 0, &os.PathError***REMOVED***
			Op:   "seek",
			Path: vd.EmbeddedDir.Filename,
			Err:  errors.New("bad file descriptor"),
		***REMOVED***
	***REMOVED***
	return 0, &os.PathError***REMOVED***
		Op:   "seek",
		Path: vd.Filename,
		Err:  errors.New("is a directory"),
	***REMOVED***
***REMOVED***
