package protoparse

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var errNoImportPathsForAbsoluteFilePath = errors.New("must specify at least one import path if any absolute file paths are given")

// ResolveFilenames tries to resolve fileNames into paths that are relative to
// directories in the given importPaths. The returned slice has the results in
// the same order as they are supplied in fileNames.
//
// The resulting names should be suitable for passing to Parser.ParseFiles.
//
// If no import paths are given and any file name is absolute, this returns an
// error.  If no import paths are given and all file names are relative, this
// returns the original file names. If a file name is already relative to one
// of the given import paths, it will be unchanged in the returned slice. If a
// file name given is relative to the current working directory, it will be made
// relative to one of the given import paths; but if it cannot be made relative
// (due to no matching import path), an error will be returned.
func ResolveFilenames(importPaths []string, fileNames ...string) ([]string, error) ***REMOVED***
	if len(importPaths) == 0 ***REMOVED***
		if containsAbsFilePath(fileNames) ***REMOVED***
			// We have to do this as otherwise parseProtoFiles can result in duplicate symbols.
			// For example, assume we import "foo/bar/bar.proto" in a file "/home/alice/dev/foo/bar/baz.proto"
			// as we call ParseFiles("/home/alice/dev/foo/bar/bar.proto","/home/alice/dev/foo/bar/baz.proto")
			// with "/home/alice/dev" as our current directory. Due to the recursive nature of parseProtoFiles,
			// it will discover the import "foo/bar/bar.proto" in the input file, and call parse on this,
			// adding "foo/bar/bar.proto" to the parsed results, as well as "/home/alice/dev/foo/bar/bar.proto"
			// from the input file list. This will result in a
			// 'duplicate symbol SYMBOL: already defined as field in "/home/alice/dev/foo/bar/bar.proto'
			// error being returned from ParseFiles.
			return nil, errNoImportPathsForAbsoluteFilePath
		***REMOVED***
		return fileNames, nil
	***REMOVED***
	absImportPaths, err := absoluteFilePaths(importPaths)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	resolvedFileNames := make([]string, 0, len(fileNames))
	for _, fileName := range fileNames ***REMOVED***
		resolvedFileName, err := resolveFilename(absImportPaths, fileName)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// On Windows, the resolved paths will use "\", but proto imports
		// require the use of "/". So fix up here.
		if filepath.Separator != '/' ***REMOVED***
			resolvedFileName = strings.Replace(resolvedFileName, string(filepath.Separator), "/", -1)
		***REMOVED***
		resolvedFileNames = append(resolvedFileNames, resolvedFileName)
	***REMOVED***
	return resolvedFileNames, nil
***REMOVED***

func containsAbsFilePath(filePaths []string) bool ***REMOVED***
	for _, filePath := range filePaths ***REMOVED***
		if filepath.IsAbs(filePath) ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

func absoluteFilePaths(filePaths []string) ([]string, error) ***REMOVED***
	absFilePaths := make([]string, 0, len(filePaths))
	for _, filePath := range filePaths ***REMOVED***
		absFilePath, err := canonicalize(filePath)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		absFilePaths = append(absFilePaths, absFilePath)
	***REMOVED***
	return absFilePaths, nil
***REMOVED***

func canonicalize(filePath string) (string, error) ***REMOVED***
	absPath, err := filepath.Abs(filePath)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	// this is kind of gross, but it lets us construct a resolved path even if some
	// path elements do not exist (a single call to filepath.EvalSymlinks would just
	// return an error, ENOENT, in that case).
	head := absPath
	tail := ""
	for ***REMOVED***
		noLinks, err := filepath.EvalSymlinks(head)
		if err == nil ***REMOVED***
			if tail != "" ***REMOVED***
				return filepath.Join(noLinks, tail), nil
			***REMOVED***
			return noLinks, nil
		***REMOVED***

		if tail == "" ***REMOVED***
			tail = filepath.Base(head)
		***REMOVED*** else ***REMOVED***
			tail = filepath.Join(filepath.Base(head), tail)
		***REMOVED***
		head = filepath.Dir(head)
		if head == "." ***REMOVED***
			// ran out of path elements to try to resolve
			return absPath, nil
		***REMOVED***
	***REMOVED***
***REMOVED***

const dotPrefix = "." + string(filepath.Separator)
const dotDotPrefix = ".." + string(filepath.Separator)

func resolveFilename(absImportPaths []string, fileName string) (string, error) ***REMOVED***
	if filepath.IsAbs(fileName) ***REMOVED***
		return resolveAbsFilename(absImportPaths, fileName)
	***REMOVED***

	if !strings.HasPrefix(fileName, dotPrefix) && !strings.HasPrefix(fileName, dotDotPrefix) ***REMOVED***
		// Use of . and .. are assumed to be relative to current working
		// directory. So if those aren't present, check to see if the file is
		// relative to an import path.
		for _, absImportPath := range absImportPaths ***REMOVED***
			absFileName := filepath.Join(absImportPath, fileName)
			_, err := os.Stat(absFileName)
			if err != nil ***REMOVED***
				continue
			***REMOVED***
			// found it! it was relative to this import path
			return fileName, nil
		***REMOVED***
	***REMOVED***

	// must be relative to current working dir
	return resolveAbsFilename(absImportPaths, fileName)
***REMOVED***

func resolveAbsFilename(absImportPaths []string, fileName string) (string, error) ***REMOVED***
	absFileName, err := canonicalize(fileName)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***
	for _, absImportPath := range absImportPaths ***REMOVED***
		if isDescendant(absImportPath, absFileName) ***REMOVED***
			resolvedPath, err := filepath.Rel(absImportPath, absFileName)
			if err != nil ***REMOVED***
				return "", err
			***REMOVED***
			return resolvedPath, nil
		***REMOVED***
	***REMOVED***
	return "", fmt.Errorf("%s does not reside in any import path", fileName)
***REMOVED***

// isDescendant returns true if file is a descendant of dir. Both dir and file must
// be cleaned, absolute paths.
func isDescendant(dir, file string) bool ***REMOVED***
	dir = filepath.Clean(dir)
	cur := file
	for ***REMOVED***
		d := filepath.Dir(cur)
		if d == dir ***REMOVED***
			return true
		***REMOVED***
		if d == "." || d == cur ***REMOVED***
			// we've run out of path elements
			return false
		***REMOVED***
		cur = d
	***REMOVED***
***REMOVED***
