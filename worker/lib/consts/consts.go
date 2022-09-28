package consts

import (
	"fmt"
	"runtime"
	"runtime/debug"
)

// Version contains the current semantic version of k6.
const Version = "0.39.0"

// VersionDetails can be set externally as part of the build process
var VersionDetails = "" //nolint:gochecknoglobals

// FullVersion returns the maximally full version and build information for
// the currently running k6 executable.
func FullVersion() string ***REMOVED***
	goVersionArch := fmt.Sprintf("%s, %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	if VersionDetails != "" ***REMOVED***
		return fmt.Sprintf("%s (%s, %s)", Version, VersionDetails, goVersionArch)
	***REMOVED***

	if buildInfo, ok := debug.ReadBuildInfo(); ok ***REMOVED***
		return fmt.Sprintf("%s (%s, %s)", Version, buildInfo.Main.Version, goVersionArch)
	***REMOVED***

	return fmt.Sprintf("%s (dev build, %s)", Version, goVersionArch)
***REMOVED***
