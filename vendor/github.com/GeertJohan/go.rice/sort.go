package rice

import "os"

// SortByName allows an array of os.FileInfo objects
// to be easily sorted by filename using sort.Sort(SortByName(array))
type SortByName []os.FileInfo

func (f SortByName) Len() int           ***REMOVED*** return len(f) ***REMOVED***
func (f SortByName) Less(i, j int) bool ***REMOVED*** return f[i].Name() < f[j].Name() ***REMOVED***
func (f SortByName) Swap(i, j int)      ***REMOVED*** f[i], f[j] = f[j], f[i] ***REMOVED***

// SortByModified allows an array of os.FileInfo objects
// to be easily sorted by modified date using sort.Sort(SortByModified(array))
type SortByModified []os.FileInfo

func (f SortByModified) Len() int           ***REMOVED*** return len(f) ***REMOVED***
func (f SortByModified) Less(i, j int) bool ***REMOVED*** return f[i].ModTime().Unix() > f[j].ModTime().Unix() ***REMOVED***
func (f SortByModified) Swap(i, j int)      ***REMOVED*** f[i], f[j] = f[j], f[i] ***REMOVED***
