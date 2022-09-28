// Copyright © 2014 Steve Francia <spf@spf13.com>.
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

type Dir interface ***REMOVED***
	Len() int
	Names() []string
	Files() []*FileData
	Add(*FileData)
	Remove(*FileData)
***REMOVED***

func RemoveFromMemDir(dir *FileData, f *FileData) ***REMOVED***
	dir.memDir.Remove(f)
***REMOVED***

func AddToMemDir(dir *FileData, f *FileData) ***REMOVED***
	dir.memDir.Add(f)
***REMOVED***

func InitializeDir(d *FileData) ***REMOVED***
	if d.memDir == nil ***REMOVED***
		d.dir = true
		d.memDir = &DirMap***REMOVED******REMOVED***
	***REMOVED***
***REMOVED***
