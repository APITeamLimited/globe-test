/*
 *
 * k6 - a next-generation load testing tool
 * Copyright (C) 2016 Load Impact
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
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"

	"strings"

	"github.com/loadimpact/k6/lib"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

type loaderFunc func(path string, parts []string) (string, error)

var loaders = []struct ***REMOVED***
	name string
	fn   loaderFunc
	expr *regexp.Regexp
***REMOVED******REMOVED***
	***REMOVED***"cdnjs", cdnjs, regexp.MustCompile(`^cdnjs.com/libraries/([^/]+)(?:/([(\d\.)]+-?[^/]*))?(?:/(.*))?$`)***REMOVED***,
	***REMOVED***"github", github, regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*)$`)***REMOVED***,
***REMOVED***

func Load(fs afero.Fs, pwd, name string) (*lib.SourceData, error) ***REMOVED***
	// We just need to make sure `import ""` doesn't crash the loader.
	if name == "" ***REMOVED***
		return nil, errors.New("local or remote path required")
	***REMOVED***

	// Do not allow the protocol to be specified, it messes everything up.
	if strings.Contains(name, "://") ***REMOVED***
		return nil, errors.New("imports should not contain a protocol")
	***REMOVED***

	// Do not allow remote-loaded scripts to lift arbitrary files off the user's machine.
	if name[0] == '/' && pwd[0] != '/' ***REMOVED***
		return nil, errors.Errorf("origin (%s) not allowed to load local file: %s", pwd, name)
	***REMOVED***

	// If the file starts with ".", resolve it as a relative path.
	if name[0] == '.' ***REMOVED***
		name = filepath.Join(pwd, name)
	***REMOVED***

	// If the resolved path starts with a "/", it's a local file.
	if name[0] == '/' ***REMOVED***
		data, err := afero.ReadFile(fs, name)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return &lib.SourceData***REMOVED***Filename: name, Data: data***REMOVED***, nil
	***REMOVED***

	// If the file is from a known service, try loading from there.
	loaderName, loader, loaderArgs := pickLoader(name)
	if loader != nil ***REMOVED***
		u, err := loader(name, loaderArgs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		data, err := fetch(u)
		if err != nil ***REMOVED***
			return nil, errors.Wrap(err, loaderName)
		***REMOVED***
		return &lib.SourceData***REMOVED***Filename: name, Data: data***REMOVED***, nil
	***REMOVED***

	// If not, load it and have a look. HTTPS is enforced, because it's 2017, HTTPS is easy,
	// running arbitrary, trivially MitM'd code (even sandboxed) is very, very bad.
	url := "https://" + name
	if !strings.ContainsRune(url, '?') ***REMOVED***
		url += "?"
	***REMOVED*** else ***REMOVED***
		url += "&"
	***REMOVED***
	url += "_k6=1"
	data, err := fetch(url)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// TODO: Parse the HTML, look for meta tags!!
	// <meta name="k6-import" content="example.com/path/to/real/file.txt" />
	// <meta name="k6-import" content="github.com/myusername/repo/file.txt" />

	return &lib.SourceData***REMOVED***Filename: name, Data: data***REMOVED***, nil
***REMOVED***

func pickLoader(path string) (string, loaderFunc, []string) ***REMOVED***
	for _, loader := range loaders ***REMOVED***
		matches := loader.expr.FindAllStringSubmatch(path, -1)
		if len(matches) > 0 ***REMOVED***
			return loader.name, loader.fn, matches[0][1:]
		***REMOVED***
	***REMOVED***
	return "", nil, nil
***REMOVED***

func fetch(u string) ([]byte, error) ***REMOVED***
	res, err := http.Get(u)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer res.Body.Close()

	if res.StatusCode != 200 ***REMOVED***
		switch res.StatusCode ***REMOVED***
		case 404:
			return nil, errors.Errorf("not found: %s", u)
		default:
			return nil, errors.Errorf("wrong status code (%d) for: %s", res.StatusCode, u)
		***REMOVED***
	***REMOVED***

	return ioutil.ReadAll(res.Body)
***REMOVED***
