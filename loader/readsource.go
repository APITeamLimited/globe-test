package loader

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"

	"github.com/APITeamLimited/k6-worker/lib/fsext"
)

// ReadSource Reads a source file from any supported destination.
func ReadSource(
	logger logrus.FieldLogger, src, pwd string, filesystems map[string]afero.Fs, stdin io.Reader,
) (*SourceData, error) ***REMOVED***
	if src == "-" ***REMOVED***
		data, err := ioutil.ReadAll(stdin)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		// TODO: don't do it in this way ...
		err = afero.WriteFile(filesystems["file"].(fsext.CacheLayerGetter).GetCachingFs(), "/-", data, 0o644)
		if err != nil ***REMOVED***
			return nil, fmt.Errorf("caching data read from -: %w", err)
		***REMOVED***
		return &SourceData***REMOVED***URL: &url.URL***REMOVED***Path: "/-", Scheme: "file"***REMOVED***, Data: data***REMOVED***, err
	***REMOVED***
	var srcLocalPath string
	if filepath.IsAbs(src) ***REMOVED***
		srcLocalPath = src
	***REMOVED*** else ***REMOVED***
		srcLocalPath = filepath.Join(pwd, src)
	***REMOVED***
	// All paths should start with a / in all fses. This is mostly for windows where it will start
	// with a volume name : C:\something.js
	srcLocalPath = filepath.Clean(afero.FilePathSeparator + srcLocalPath)
	if ok, _ := afero.Exists(filesystems["file"], srcLocalPath); ok ***REMOVED***
		// there is file on the local disk ... lets use it :)
		return Load(logger, filesystems, &url.URL***REMOVED***Scheme: "file", Path: filepath.ToSlash(srcLocalPath)***REMOVED***, src)
	***REMOVED***

	pwdURL := &url.URL***REMOVED***Scheme: "file", Path: filepath.ToSlash(filepath.Clean(pwd)) + "/"***REMOVED***
	srcURL, err := Resolve(pwdURL, filepath.ToSlash(src))
	if err != nil ***REMOVED***
		var noSchemeError noSchemeRemoteModuleResolutionError
		if errors.As(err, &noSchemeError) ***REMOVED***
			// TODO maybe try to wrap the original error here as well, without butchering the message
			return nil, fmt.Errorf(nothingWorkedLoadedMsg, noSchemeError.moduleSpecifier, noSchemeError.err)
		***REMOVED***
		return nil, err
	***REMOVED***
	result, err := Load(logger, filesystems, srcURL, src)
	var noSchemeError noSchemeRemoteModuleResolutionError
	if errors.As(err, &noSchemeError) ***REMOVED***
		// TODO maybe try to wrap the original error here as well, without butchering the message
		return nil, fmt.Errorf(nothingWorkedLoadedMsg, noSchemeError.moduleSpecifier, noSchemeError.err)
	***REMOVED***

	return result, err
***REMOVED***
