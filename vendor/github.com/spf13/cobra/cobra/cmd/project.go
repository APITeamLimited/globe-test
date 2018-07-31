package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Project contains name, license and paths to projects.
type Project struct ***REMOVED***
	absPath string
	cmdPath string
	srcPath string
	license License
	name    string
***REMOVED***

// NewProject returns Project with specified project name.
func NewProject(projectName string) *Project ***REMOVED***
	if projectName == "" ***REMOVED***
		er("can't create project with blank name")
	***REMOVED***

	p := new(Project)
	p.name = projectName

	// 1. Find already created protect.
	p.absPath = findPackage(projectName)

	// 2. If there are no created project with this path, and user is in GOPATH,
	// then use GOPATH/src/projectName.
	if p.absPath == "" ***REMOVED***
		wd, err := os.Getwd()
		if err != nil ***REMOVED***
			er(err)
		***REMOVED***
		for _, srcPath := range srcPaths ***REMOVED***
			goPath := filepath.Dir(srcPath)
			if filepathHasPrefix(wd, goPath) ***REMOVED***
				p.absPath = filepath.Join(srcPath, projectName)
				break
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// 3. If user is not in GOPATH, then use (first GOPATH)/src/projectName.
	if p.absPath == "" ***REMOVED***
		p.absPath = filepath.Join(srcPaths[0], projectName)
	***REMOVED***

	return p
***REMOVED***

// findPackage returns full path to existing go package in GOPATHs.
func findPackage(packageName string) string ***REMOVED***
	if packageName == "" ***REMOVED***
		return ""
	***REMOVED***

	for _, srcPath := range srcPaths ***REMOVED***
		packagePath := filepath.Join(srcPath, packageName)
		if exists(packagePath) ***REMOVED***
			return packagePath
		***REMOVED***
	***REMOVED***

	return ""
***REMOVED***

// NewProjectFromPath returns Project with specified absolute path to
// package.
func NewProjectFromPath(absPath string) *Project ***REMOVED***
	if absPath == "" ***REMOVED***
		er("can't create project: absPath can't be blank")
	***REMOVED***
	if !filepath.IsAbs(absPath) ***REMOVED***
		er("can't create project: absPath is not absolute")
	***REMOVED***

	// If absPath is symlink, use its destination.
	fi, err := os.Lstat(absPath)
	if err != nil ***REMOVED***
		er("can't read path info: " + err.Error())
	***REMOVED***
	if fi.Mode()&os.ModeSymlink != 0 ***REMOVED***
		path, err := os.Readlink(absPath)
		if err != nil ***REMOVED***
			er("can't read the destination of symlink: " + err.Error())
		***REMOVED***
		absPath = path
	***REMOVED***

	p := new(Project)
	p.absPath = strings.TrimSuffix(absPath, findCmdDir(absPath))
	p.name = filepath.ToSlash(trimSrcPath(p.absPath, p.SrcPath()))
	return p
***REMOVED***

// trimSrcPath trims at the beginning of absPath the srcPath.
func trimSrcPath(absPath, srcPath string) string ***REMOVED***
	relPath, err := filepath.Rel(srcPath, absPath)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
	return relPath
***REMOVED***

// License returns the License object of project.
func (p *Project) License() License ***REMOVED***
	if p.license.Text == "" && p.license.Name != "None" ***REMOVED***
		p.license = getLicense()
	***REMOVED***
	return p.license
***REMOVED***

// Name returns the name of project, e.g. "github.com/spf13/cobra"
func (p Project) Name() string ***REMOVED***
	return p.name
***REMOVED***

// CmdPath returns absolute path to directory, where all commands are located.
func (p *Project) CmdPath() string ***REMOVED***
	if p.absPath == "" ***REMOVED***
		return ""
	***REMOVED***
	if p.cmdPath == "" ***REMOVED***
		p.cmdPath = filepath.Join(p.absPath, findCmdDir(p.absPath))
	***REMOVED***
	return p.cmdPath
***REMOVED***

// findCmdDir checks if base of absPath is cmd dir and returns it or
// looks for existing cmd dir in absPath.
func findCmdDir(absPath string) string ***REMOVED***
	if !exists(absPath) || isEmpty(absPath) ***REMOVED***
		return "cmd"
	***REMOVED***

	if isCmdDir(absPath) ***REMOVED***
		return filepath.Base(absPath)
	***REMOVED***

	files, _ := filepath.Glob(filepath.Join(absPath, "c*"))
	for _, file := range files ***REMOVED***
		if isCmdDir(file) ***REMOVED***
			return filepath.Base(file)
		***REMOVED***
	***REMOVED***

	return "cmd"
***REMOVED***

// isCmdDir checks if base of name is one of cmdDir.
func isCmdDir(name string) bool ***REMOVED***
	name = filepath.Base(name)
	for _, cmdDir := range []string***REMOVED***"cmd", "cmds", "command", "commands"***REMOVED*** ***REMOVED***
		if name == cmdDir ***REMOVED***
			return true
		***REMOVED***
	***REMOVED***
	return false
***REMOVED***

// AbsPath returns absolute path of project.
func (p Project) AbsPath() string ***REMOVED***
	return p.absPath
***REMOVED***

// SrcPath returns absolute path to $GOPATH/src where project is located.
func (p *Project) SrcPath() string ***REMOVED***
	if p.srcPath != "" ***REMOVED***
		return p.srcPath
	***REMOVED***
	if p.absPath == "" ***REMOVED***
		p.srcPath = srcPaths[0]
		return p.srcPath
	***REMOVED***

	for _, srcPath := range srcPaths ***REMOVED***
		if filepathHasPrefix(p.absPath, srcPath) ***REMOVED***
			p.srcPath = srcPath
			break
		***REMOVED***
	***REMOVED***

	return p.srcPath
***REMOVED***

func filepathHasPrefix(path string, prefix string) bool ***REMOVED***
	if len(path) <= len(prefix) ***REMOVED***
		return false
	***REMOVED***
	if runtime.GOOS == "windows" ***REMOVED***
		// Paths in windows are case-insensitive.
		return strings.EqualFold(path[0:len(prefix)], prefix)
	***REMOVED***
	return path[0:len(prefix)] == prefix

***REMOVED***
