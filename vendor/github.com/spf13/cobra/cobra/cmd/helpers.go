// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var srcPaths []string

func init() ***REMOVED***
	// Initialize srcPaths.
	envGoPath := os.Getenv("GOPATH")
	goPaths := filepath.SplitList(envGoPath)
	if len(goPaths) == 0 ***REMOVED***
		// Adapted from https://github.com/Masterminds/glide/pull/798/files.
		// As of Go 1.8 the GOPATH is no longer required to be set. Instead there
		// is a default value. If there is no GOPATH check for the default value.
		// Note, checking the GOPATH first to avoid invoking the go toolchain if
		// possible.

		goExecutable := os.Getenv("COBRA_GO_EXECUTABLE")
		if len(goExecutable) <= 0 ***REMOVED***
			goExecutable = "go"
		***REMOVED***

		out, err := exec.Command(goExecutable, "env", "GOPATH").Output()
		if err != nil ***REMOVED***
			er(err)
		***REMOVED***

		toolchainGoPath := strings.TrimSpace(string(out))
		goPaths = filepath.SplitList(toolchainGoPath)
		if len(goPaths) == 0 ***REMOVED***
			er("$GOPATH is not set")
		***REMOVED***
	***REMOVED***
	srcPaths = make([]string, 0, len(goPaths))
	for _, goPath := range goPaths ***REMOVED***
		srcPaths = append(srcPaths, filepath.Join(goPath, "src"))
	***REMOVED***
***REMOVED***

func er(msg interface***REMOVED******REMOVED***) ***REMOVED***
	fmt.Println("Error:", msg)
	os.Exit(1)
***REMOVED***

// isEmpty checks if a given path is empty.
// Hidden files in path are ignored.
func isEmpty(path string) bool ***REMOVED***
	fi, err := os.Stat(path)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***

	if !fi.IsDir() ***REMOVED***
		return fi.Size() == 0
	***REMOVED***

	f, err := os.Open(path)
	if err != nil ***REMOVED***
		er(err)
	***REMOVED***
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil && err != io.EOF ***REMOVED***
		er(err)
	***REMOVED***

	for _, name := range names ***REMOVED***
		if len(name) > 0 && name[0] != '.' ***REMOVED***
			return false
		***REMOVED***
	***REMOVED***
	return true
***REMOVED***

// exists checks if a file or directory exists.
func exists(path string) bool ***REMOVED***
	if path == "" ***REMOVED***
		return false
	***REMOVED***
	_, err := os.Stat(path)
	if err == nil ***REMOVED***
		return true
	***REMOVED***
	if !os.IsNotExist(err) ***REMOVED***
		er(err)
	***REMOVED***
	return false
***REMOVED***

func executeTemplate(tmplStr string, data interface***REMOVED******REMOVED***) (string, error) ***REMOVED***
	tmpl, err := template.New("").Funcs(template.FuncMap***REMOVED***"comment": commentifyString***REMOVED***).Parse(tmplStr)
	if err != nil ***REMOVED***
		return "", err
	***REMOVED***

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	return buf.String(), err
***REMOVED***

func writeStringToFile(path string, s string) error ***REMOVED***
	return writeToFile(path, strings.NewReader(s))
***REMOVED***

// writeToFile writes r to file with path only
// if file/directory on given path doesn't exist.
func writeToFile(path string, r io.Reader) error ***REMOVED***
	if exists(path) ***REMOVED***
		return fmt.Errorf("%v already exists", path)
	***REMOVED***

	dir := filepath.Dir(path)
	if dir != "" ***REMOVED***
		if err := os.MkdirAll(dir, 0777); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	file, err := os.Create(path)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer file.Close()

	_, err = io.Copy(file, r)
	return err
***REMOVED***

// commentfyString comments every line of in.
func commentifyString(in string) string ***REMOVED***
	var newlines []string
	lines := strings.Split(in, "\n")
	for _, line := range lines ***REMOVED***
		if strings.HasPrefix(line, "//") ***REMOVED***
			newlines = append(newlines, line)
		***REMOVED*** else ***REMOVED***
			if line == "" ***REMOVED***
				newlines = append(newlines, "//")
			***REMOVED*** else ***REMOVED***
				newlines = append(newlines, "// "+line)
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return strings.Join(newlines, "\n")
***REMOVED***
