package configdir

import "os"

var hasVendorName = true
var systemSettingFolders = []string***REMOVED***os.Getenv("PROGRAMDATA")***REMOVED***
var globalSettingFolder = os.Getenv("APPDATA")
var cacheFolder = os.Getenv("LOCALAPPDATA")
