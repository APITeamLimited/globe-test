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

package mem

import "sort"

type DirMap map[string]*FileData

func (m DirMap) Len() int           ***REMOVED*** return len(m) ***REMOVED***
func (m DirMap) Add(f *FileData)    ***REMOVED*** m[f.name] = f ***REMOVED***
func (m DirMap) Remove(f *FileData) ***REMOVED*** delete(m, f.name) ***REMOVED***
func (m DirMap) Files() (files []*FileData) ***REMOVED***
	for _, f := range m ***REMOVED***
		files = append(files, f)
	***REMOVED***
	sort.Sort(filesSorter(files))
	return files
***REMOVED***

// implement sort.Interface for []*FileData
type filesSorter []*FileData

func (s filesSorter) Len() int           ***REMOVED*** return len(s) ***REMOVED***
func (s filesSorter) Swap(i, j int)      ***REMOVED*** s[i], s[j] = s[j], s[i] ***REMOVED***
func (s filesSorter) Less(i, j int) bool ***REMOVED*** return s[i].name < s[j].name ***REMOVED***

func (m DirMap) Names() (names []string) ***REMOVED***
	for x := range m ***REMOVED***
		names = append(names, x)
	***REMOVED***
	return names
***REMOVED***
