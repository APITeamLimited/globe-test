package rice

import (
	"os"
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

// re-type to make exported methods invisible to user (godoc)
// they're not required for the user
// embeddedDirInfo implements os.FileInfo
type embeddedDirInfo embedded.EmbeddedDir

// Name returns the base name of the directory
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) Name() string ***REMOVED***
	return ed.Filename
***REMOVED***

// Size always returns 0
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) Size() int64 ***REMOVED***
	return 0
***REMOVED***

// Mode returns the file mode bits
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) Mode() os.FileMode ***REMOVED***
	return os.FileMode(0555 | os.ModeDir) // dr-xr-xr-x
***REMOVED***

// ModTime returns the modification time
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) ModTime() time.Time ***REMOVED***
	return ed.DirModTime
***REMOVED***

// IsDir returns the abbreviation for Mode().IsDir() (always true)
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) IsDir() bool ***REMOVED***
	return true
***REMOVED***

// Sys returns the underlying data source (always nil)
// (implementing os.FileInfo)
func (ed *embeddedDirInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***

// re-type to make exported methods invisible to user (godoc)
// they're not required for the user
// embeddedFileInfo implements os.FileInfo
type embeddedFileInfo embedded.EmbeddedFile

// Name returns the base name of the file
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) Name() string ***REMOVED***
	return ef.Filename
***REMOVED***

// Size returns the length in bytes for regular files; system-dependent for others
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) Size() int64 ***REMOVED***
	return int64(len(ef.Content))
***REMOVED***

// Mode returns the file mode bits
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) Mode() os.FileMode ***REMOVED***
	return os.FileMode(0555) // r-xr-xr-x
***REMOVED***

// ModTime returns the modification time
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) ModTime() time.Time ***REMOVED***
	return ef.FileModTime
***REMOVED***

// IsDir returns the abbreviation for Mode().IsDir() (always false)
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) IsDir() bool ***REMOVED***
	return false
***REMOVED***

// Sys returns the underlying data source (always nil)
// (implementing os.FileInfo)
func (ef *embeddedFileInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***
