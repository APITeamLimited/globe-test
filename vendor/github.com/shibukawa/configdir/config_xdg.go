// +build !windows,!darwin

package configdir

import (
	"os"
	"path/filepath"
	"strings"
)

// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html

var hasVendorName = true
var systemSettingFolders []string
var globalSettingFolder string
var cacheFolder string

func init() ***REMOVED***
	if os.Getenv("XDG_CONFIG_HOME") != "" ***REMOVED***
		globalSettingFolder = os.Getenv("XDG_CONFIG_HOME")
	***REMOVED*** else ***REMOVED***
		globalSettingFolder = filepath.Join(os.Getenv("HOME"), ".config")
	***REMOVED***
	if os.Getenv("XDG_CONFIG_DIRS") != "" ***REMOVED***
		systemSettingFolders = strings.Split(os.Getenv("XDG_CONFIG_DIRS"), ":")
	***REMOVED*** else ***REMOVED***
		systemSettingFolders = []string***REMOVED***"/etc/xdg"***REMOVED***
	***REMOVED***
	if os.Getenv("XDG_CACHE_HOME") != "" ***REMOVED***
		cacheFolder = os.Getenv("XDG_CACHE_HOME")
	***REMOVED*** else ***REMOVED***
		cacheFolder = filepath.Join(os.Getenv("HOME"), ".cache")
	***REMOVED***
***REMOVED***
