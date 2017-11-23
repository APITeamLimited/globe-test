package configdir

import "os"

var hasVendorName = true
var systemSettingFolders = []string***REMOVED***"/Library/Application Support"***REMOVED***
var globalSettingFolder = os.Getenv("HOME") + "/Library/Application Support"
var cacheFolder = os.Getenv("HOME") + "/Library/Caches"
