package rice

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

// Box abstracts a directory for resources/files.
// It can either load files from disk, or from embedded code (when `rice --embed` was ran).
type Box struct ***REMOVED***
	name         string
	absolutePath string
	embed        *embedded.EmbeddedBox
	appendd      *appendedBox
***REMOVED***

var defaultLocateOrder = []LocateMethod***REMOVED***LocateEmbedded, LocateAppended, LocateFS***REMOVED***

func findBox(name string, order []LocateMethod) (*Box, error) ***REMOVED***
	b := &Box***REMOVED***name: name***REMOVED***

	// no support for absolute paths since gopath can be different on different machines.
	// therefore, required box must be located relative to package requiring it.
	if filepath.IsAbs(name) ***REMOVED***
		return nil, errors.New("given name/path is absolute")
	***REMOVED***

	var err error
	for _, method := range order ***REMOVED***
		switch method ***REMOVED***
		case LocateEmbedded:
			if embed := embedded.EmbeddedBoxes[name]; embed != nil ***REMOVED***
				b.embed = embed
				return b, nil
			***REMOVED***

		case LocateAppended:
			appendedBoxName := strings.Replace(name, `/`, `-`, -1)
			if appendd := appendedBoxes[appendedBoxName]; appendd != nil ***REMOVED***
				b.appendd = appendd
				return b, nil
			***REMOVED***

		case LocateFS:
			// resolve absolute directory path
			err := b.resolveAbsolutePathFromCaller()
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			// check if absolutePath exists on filesystem
			info, err := os.Stat(b.absolutePath)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			// check if absolutePath is actually a directory
			if !info.IsDir() ***REMOVED***
				err = errors.New("given name/path is not a directory")
				continue
			***REMOVED***
			return b, nil
		case LocateWorkingDirectory:
			// resolve absolute directory path
			err := b.resolveAbsolutePathFromWorkingDirectory()
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			// check if absolutePath exists on filesystem
			info, err := os.Stat(b.absolutePath)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			// check if absolutePath is actually a directory
			if !info.IsDir() ***REMOVED***
				err = errors.New("given name/path is not a directory")
				continue
			***REMOVED***
			return b, nil
		***REMOVED***
	***REMOVED***

	if err == nil ***REMOVED***
		err = fmt.Errorf("could not locate box %q", name)
	***REMOVED***

	return nil, err
***REMOVED***

// FindBox returns a Box instance for given name.
// When the given name is a relative path, it's base path will be the calling pkg/cmd's source root.
// When the given name is absolute, it's absolute. derp.
// Make sure the path doesn't contain any sensitive information as it might be placed into generated go source (embedded).
func FindBox(name string) (*Box, error) ***REMOVED***
	return findBox(name, defaultLocateOrder)
***REMOVED***

// MustFindBox returns a Box instance for given name, like FindBox does.
// It does not return an error, instead it panics when an error occurs.
func MustFindBox(name string) *Box ***REMOVED***
	box, err := findBox(name, defaultLocateOrder)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return box
***REMOVED***

// This is injected as a mutable function literal so that we can mock it out in
// tests and return a fixed test file.
var resolveAbsolutePathFromCaller = func(name string, nStackFrames int) (string, error) ***REMOVED***
	_, callingGoFile, _, ok := runtime.Caller(nStackFrames)
	if !ok ***REMOVED***
		return "", errors.New("couldn't find caller on stack")
	***REMOVED***

	// resolve to proper path
	pkgDir := filepath.Dir(callingGoFile)
	// fix for go cover
	const coverPath = "_test/_obj_test"
	if !filepath.IsAbs(pkgDir) ***REMOVED***
		if i := strings.Index(pkgDir, coverPath); i >= 0 ***REMOVED***
			pkgDir = pkgDir[:i] + pkgDir[i+len(coverPath):]            // remove coverPath
			pkgDir = filepath.Join(os.Getenv("GOPATH"), "src", pkgDir) // make absolute
		***REMOVED***
	***REMOVED***
	return filepath.Join(pkgDir, name), nil
***REMOVED***

func (b *Box) resolveAbsolutePathFromCaller() error ***REMOVED***
	path, err := resolveAbsolutePathFromCaller(b.name, 4)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.absolutePath = path
	return nil

***REMOVED***

func (b *Box) resolveAbsolutePathFromWorkingDirectory() error ***REMOVED***
	path, err := os.Getwd()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	b.absolutePath = filepath.Join(path, b.name)
	return nil
***REMOVED***

// IsEmbedded indicates wether this box was embedded into the application
func (b *Box) IsEmbedded() bool ***REMOVED***
	return b.embed != nil
***REMOVED***

// IsAppended indicates wether this box was appended to the application
func (b *Box) IsAppended() bool ***REMOVED***
	return b.appendd != nil
***REMOVED***

// Time returns how actual the box is.
// When the box is embedded, it's value is saved in the embedding code.
// When the box is live, this methods returns time.Now()
func (b *Box) Time() time.Time ***REMOVED***
	if b.IsEmbedded() ***REMOVED***
		return b.embed.Time
	***REMOVED***

	//++ TODO: return time for appended box

	return time.Now()
***REMOVED***

// Open opens a File from the box
// If there is an error, it will be of type *os.PathError.
func (b *Box) Open(name string) (*File, error) ***REMOVED***
	if Debug ***REMOVED***
		fmt.Printf("Open(%s)\n", name)
	***REMOVED***

	if b.IsEmbedded() ***REMOVED***
		if Debug ***REMOVED***
			fmt.Println("Box is embedded")
		***REMOVED***

		// trim prefix (paths are relative to box)
		name = strings.TrimLeft(name, "/")
		if Debug ***REMOVED***
			fmt.Printf("Trying %s\n", name)
		***REMOVED***

		// search for file
		ef := b.embed.Files[name]
		if ef == nil ***REMOVED***
			if Debug ***REMOVED***
				fmt.Println("Didn't find file in embed")
			***REMOVED***
			// file not found, try dir
			ed := b.embed.Dirs[name]
			if ed == nil ***REMOVED***
				if Debug ***REMOVED***
					fmt.Println("Didn't find dir in embed")
				***REMOVED***
				// dir not found, error out
				return nil, &os.PathError***REMOVED***
					Op:   "open",
					Path: name,
					Err:  os.ErrNotExist,
				***REMOVED***
			***REMOVED***
			if Debug ***REMOVED***
				fmt.Println("Found dir. Returning virtual dir")
			***REMOVED***
			vd := newVirtualDir(ed)
			return &File***REMOVED***virtualD: vd***REMOVED***, nil
		***REMOVED***

		// box is embedded
		if Debug ***REMOVED***
			fmt.Println("Found file. Returning virtual file")
		***REMOVED***
		vf := newVirtualFile(ef)
		return &File***REMOVED***virtualF: vf***REMOVED***, nil
	***REMOVED***

	if b.IsAppended() ***REMOVED***
		// trim prefix (paths are relative to box)
		name = strings.TrimLeft(name, "/")

		// search for file
		appendedFile := b.appendd.Files[name]
		if appendedFile == nil ***REMOVED***
			return nil, &os.PathError***REMOVED***
				Op:   "open",
				Path: name,
				Err:  os.ErrNotExist,
			***REMOVED***
		***REMOVED***

		// create new file
		f := &File***REMOVED***
			appendedF: appendedFile,
		***REMOVED***

		// if this file is a directory, we want to be able to read and seek
		if !appendedFile.dir ***REMOVED***
			// looks like malformed data in zip, error now
			if appendedFile.content == nil ***REMOVED***
				return nil, &os.PathError***REMOVED***
					Op:   "open",
					Path: "name",
					Err:  errors.New("error reading data from zip file"),
				***REMOVED***
			***REMOVED***
			// create new bytes.Reader
			f.appendedFileReader = bytes.NewReader(appendedFile.content)
		***REMOVED***

		// all done
		return f, nil
	***REMOVED***

	// perform os open
	if Debug ***REMOVED***
		fmt.Printf("Using os.Open(%s)", filepath.Join(b.absolutePath, name))
	***REMOVED***
	file, err := os.Open(filepath.Join(b.absolutePath, name))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return &File***REMOVED***realF: file***REMOVED***, nil
***REMOVED***

// Bytes returns the content of the file with given name as []byte.
func (b *Box) Bytes(name string) ([]byte, error) ***REMOVED***
	file, err := b.Open(name)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return content, nil
***REMOVED***

// MustBytes returns the content of the file with given name as []byte.
// panic's on error.
func (b *Box) MustBytes(name string) []byte ***REMOVED***
	bts, err := b.Bytes(name)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return bts
***REMOVED***

// String returns the content of the file with given name as string.
func (b *Box) String(name string) (string, error) ***REMOVED***
	// check if box is embedded, optimized fast path
	if b.IsEmbedded() ***REMOVED***
		// find file in embed
		ef := b.embed.Files[name]
		if ef == nil ***REMOVED***
			return "", os.ErrNotExist
		***REMOVED***
		// return as string
		return ef.Content, nil
	***REMOVED***

	bts, err := b.Bytes(name)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	return string(bts), nil
***REMOVED***

// MustString returns the content of the file with given name as string.
// panic's on error.
func (b *Box) MustString(name string) string ***REMOVED***
	str, err := b.String(name)
	if err != nil ***REMOVED***
		panic(err)
	***REMOVED***
	return str
***REMOVED***

// Name returns the name of the box
func (b *Box) Name() string ***REMOVED***
	return b.name
***REMOVED***
