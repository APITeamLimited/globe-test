package rice

import (
	"archive/zip"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/daaku/go.zipexe"
	"github.com/kardianos/osext"
)

// appendedBox defines an appended box
type appendedBox struct ***REMOVED***
	Name  string                   // box name
	Files map[string]*appendedFile // appended files (*zip.File) by full path
***REMOVED***

type appendedFile struct ***REMOVED***
	zipFile  *zip.File
	dir      bool
	dirInfo  *appendedDirInfo
	children []*appendedFile
	content  []byte
***REMOVED***

// appendedBoxes is a public register of appendes boxes
var appendedBoxes = make(map[string]*appendedBox)

func init() ***REMOVED***
	// find if exec is appended
	thisFile, err := osext.Executable()
	if err != nil ***REMOVED***
		return // not appended or cant find self executable
	***REMOVED***
	closer, rd, err := zipexe.OpenCloser(thisFile)
	if err != nil ***REMOVED***
		return // not appended
	***REMOVED***
	defer closer.Close()

	for _, f := range rd.File ***REMOVED***
		// get box and file name from f.Name
		fileParts := strings.SplitN(strings.TrimLeft(filepath.ToSlash(f.Name), "/"), "/", 2)
		boxName := fileParts[0]
		var fileName string
		if len(fileParts) > 1 ***REMOVED***
			fileName = fileParts[1]
		***REMOVED***

		// find box or create new one if doesn't exist
		box := appendedBoxes[boxName]
		if box == nil ***REMOVED***
			box = &appendedBox***REMOVED***
				Name:  boxName,
				Files: make(map[string]*appendedFile),
			***REMOVED***
			appendedBoxes[boxName] = box
		***REMOVED***

		// create and add file to box
		af := &appendedFile***REMOVED***
			zipFile: f,
		***REMOVED***
		if f.Comment == "dir" ***REMOVED***
			af.dir = true
			af.dirInfo = &appendedDirInfo***REMOVED***
				name: filepath.Base(af.zipFile.Name),
				//++ TODO: use zip modtime when that is set correctly: af.zipFile.ModTime()
				time: time.Now(),
			***REMOVED***
		***REMOVED*** else ***REMOVED***
			// this is a file, we need it's contents so we can create a bytes.Reader when the file is opened
			// make a new byteslice
			af.content = make([]byte, af.zipFile.FileInfo().Size())
			// ignore reading empty files from zip (empty file still is a valid file to be read though!)
			if len(af.content) > 0 ***REMOVED***
				// open io.ReadCloser
				rc, err := af.zipFile.Open()
				if err != nil ***REMOVED***
					af.content = nil // this will cause an error when the file is being opened or seeked (which is good)
					// TODO: it's quite blunt to just log this stuff. but this is in init, so rice.Debug can't be changed yet..
					log.Printf("error opening appended file %s: %v", af.zipFile.Name, err)
				***REMOVED*** else ***REMOVED***
					_, err = rc.Read(af.content)
					rc.Close()
					if err != nil ***REMOVED***
						af.content = nil // this will cause an error when the file is being opened or seeked (which is good)
						// TODO: it's quite blunt to just log this stuff. but this is in init, so rice.Debug can't be changed yet..
						log.Printf("error reading data for appended file %s: %v", af.zipFile.Name, err)
					***REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***

		// add appendedFile to box file list
		box.Files[fileName] = af

		// add to parent dir (if any)
		dirName := filepath.Dir(fileName)
		if dirName == "." ***REMOVED***
			dirName = ""
		***REMOVED***
		if fileName != "" ***REMOVED*** // don't make box root dir a child of itself
			if dir := box.Files[dirName]; dir != nil ***REMOVED***
				dir.children = append(dir.children, af)
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

// implements os.FileInfo.
// used for Readdir()
type appendedDirInfo struct ***REMOVED***
	name string
	time time.Time
***REMOVED***

func (adi *appendedDirInfo) Name() string ***REMOVED***
	return adi.name
***REMOVED***
func (adi *appendedDirInfo) Size() int64 ***REMOVED***
	return 0
***REMOVED***
func (adi *appendedDirInfo) Mode() os.FileMode ***REMOVED***
	return os.ModeDir
***REMOVED***
func (adi *appendedDirInfo) ModTime() time.Time ***REMOVED***
	return adi.time
***REMOVED***
func (adi *appendedDirInfo) IsDir() bool ***REMOVED***
	return true
***REMOVED***
func (adi *appendedDirInfo) Sys() interface***REMOVED******REMOVED*** ***REMOVED***
	return nil
***REMOVED***
