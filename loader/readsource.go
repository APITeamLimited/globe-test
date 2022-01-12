/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2019 Load Impact
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

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

	"go.k6.io/k6/lib/fsext"
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
