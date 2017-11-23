package rice

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
)

// File implements the io.Reader, io.Seeker, io.Closer and http.File interfaces
type File struct ***REMOVED***
	// File abstracts file methods so the user doesn't see the difference between rice.virtualFile, rice.virtualDir and os.File
	// TODO: maybe use internal File interface and four implementations: *os.File, appendedFile, virtualFile, virtualDir

	// real file on disk
	realF *os.File

	// when embedded (go)
	virtualF *virtualFile
	virtualD *virtualDir

	// when appended (zip)
	appendedF          *appendedFile
	appendedFileReader *bytes.Reader
	// TODO: is appendedFileReader subject of races? Might need a lock here..
***REMOVED***

// Close is like (*os.File).Close()
// Visit http://golang.org/pkg/os/#File.Close for more information
func (f *File) Close() error ***REMOVED***
	if f.appendedF != nil ***REMOVED***
		if f.appendedFileReader == nil ***REMOVED***
			return errors.New("already closed")
		***REMOVED***
		f.appendedFileReader = nil
		return nil
	***REMOVED***
	if f.virtualF != nil ***REMOVED***
		return f.virtualF.close()
	***REMOVED***
	if f.virtualD != nil ***REMOVED***
		return f.virtualD.close()
	***REMOVED***
	return f.realF.Close()
***REMOVED***

// Stat is like (*os.File).Stat()
// Visit http://golang.org/pkg/os/#File.Stat for more information
func (f *File) Stat() (os.FileInfo, error) ***REMOVED***
	if f.appendedF != nil ***REMOVED***
		if f.appendedF.dir ***REMOVED***
			return f.appendedF.dirInfo, nil
		***REMOVED***
		if f.appendedFileReader == nil ***REMOVED***
			return nil, errors.New("file is closed")
		***REMOVED***
		return f.appendedF.zipFile.FileInfo(), nil
	***REMOVED***
	if f.virtualF != nil ***REMOVED***
		return f.virtualF.stat()
	***REMOVED***
	if f.virtualD != nil ***REMOVED***
		return f.virtualD.stat()
	***REMOVED***
	return f.realF.Stat()
***REMOVED***

// Readdir is like (*os.File).Readdir()
// Visit http://golang.org/pkg/os/#File.Readdir for more information
func (f *File) Readdir(count int) ([]os.FileInfo, error) ***REMOVED***
	if f.appendedF != nil ***REMOVED***
		if f.appendedF.dir ***REMOVED***
			fi := make([]os.FileInfo, 0, len(f.appendedF.children))
			for _, childAppendedFile := range f.appendedF.children ***REMOVED***
				if childAppendedFile.dir ***REMOVED***
					fi = append(fi, childAppendedFile.dirInfo)
				***REMOVED*** else ***REMOVED***
					fi = append(fi, childAppendedFile.zipFile.FileInfo())
				***REMOVED***
			***REMOVED***
			return fi, nil
		***REMOVED***
		//++ TODO: is os.ErrInvalid the correct error for Readdir on file?
		return nil, os.ErrInvalid
	***REMOVED***
	if f.virtualF != nil ***REMOVED***
		return f.virtualF.readdir(count)
	***REMOVED***
	if f.virtualD != nil ***REMOVED***
		return f.virtualD.readdir(count)
	***REMOVED***
	return f.realF.Readdir(count)
***REMOVED***

// Read is like (*os.File).Read()
// Visit http://golang.org/pkg/os/#File.Read for more information
func (f *File) Read(bts []byte) (int, error) ***REMOVED***
	if f.appendedF != nil ***REMOVED***
		if f.appendedFileReader == nil ***REMOVED***
			return 0, &os.PathError***REMOVED***
				Op:   "read",
				Path: filepath.Base(f.appendedF.zipFile.Name),
				Err:  errors.New("file is closed"),
			***REMOVED***
		***REMOVED***
		if f.appendedF.dir ***REMOVED***
			return 0, &os.PathError***REMOVED***
				Op:   "read",
				Path: filepath.Base(f.appendedF.zipFile.Name),
				Err:  errors.New("is a directory"),
			***REMOVED***
		***REMOVED***
		return f.appendedFileReader.Read(bts)
	***REMOVED***
	if f.virtualF != nil ***REMOVED***
		return f.virtualF.read(bts)
	***REMOVED***
	if f.virtualD != nil ***REMOVED***
		return f.virtualD.read(bts)
	***REMOVED***
	return f.realF.Read(bts)
***REMOVED***

// Seek is like (*os.File).Seek()
// Visit http://golang.org/pkg/os/#File.Seek for more information
func (f *File) Seek(offset int64, whence int) (int64, error) ***REMOVED***
	if f.appendedF != nil ***REMOVED***
		if f.appendedFileReader == nil ***REMOVED***
			return 0, &os.PathError***REMOVED***
				Op:   "seek",
				Path: filepath.Base(f.appendedF.zipFile.Name),
				Err:  errors.New("file is closed"),
			***REMOVED***
		***REMOVED***
		return f.appendedFileReader.Seek(offset, whence)
	***REMOVED***
	if f.virtualF != nil ***REMOVED***
		return f.virtualF.seek(offset, whence)
	***REMOVED***
	if f.virtualD != nil ***REMOVED***
		return f.virtualD.seek(offset, whence)
	***REMOVED***
	return f.realF.Seek(offset, whence)
***REMOVED***
