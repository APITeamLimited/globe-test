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
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

// SourceData wraps a source file; data and filename.
type SourceData struct ***REMOVED***
	Data []byte
	URL  *url.URL
***REMOVED***

type loaderFunc func(path string, parts []string) (string, error)

//nolint: gochecknoglobals
var (
	loaders = []struct ***REMOVED***
		name string
		fn   loaderFunc
		expr *regexp.Regexp
	***REMOVED******REMOVED***
		***REMOVED***"cdnjs", cdnjs, regexp.MustCompile(`^cdnjs.com/libraries/([^/]+)(?:/([(\d\.)]+-?[^/]*))?(?:/(.*))?$`)***REMOVED***,
		***REMOVED***"github", github, regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*)$`)***REMOVED***,
	***REMOVED***
	httpsSchemeCouldntBeLoadedMsg = `The moduleSpecifier "%s" couldn't be retrieved from` +
		` the resolved url "%s". Error : "%s"`
	fileSchemeCouldntBeLoadedMsg = `The moduleSpecifier "%s" couldn't be found on ` +
		`local disk. Make sure that you've specified the right path to the file. If you're ` +
		`running k6 using the Docker image make sure you have mounted the ` +
		`local directory (-v /local/path/:/inside/docker/path) containing ` +
		`your script and modules so that they're accessible by k6 from ` +
		`inside of the container, see ` +
		`https://docs.k6.io/v1.0/docs/modules#section-using-local-modules-with-docker.`
	errNoLoaderMatched = errors.New("no loader matched")
)

// Resolve a relative path to an absolute one.
func Resolve(pwd *url.URL, moduleSpecifier string) (*url.URL, error) ***REMOVED***
	if moduleSpecifier == "" ***REMOVED***
		return nil, errors.New("local or remote path required")
	***REMOVED***

	if moduleSpecifier[0] == '.' || moduleSpecifier[0] == '/' || filepath.IsAbs(moduleSpecifier) ***REMOVED***
		if pwd.Opaque != "" ***REMOVED*** // this is a loader reference
			parts := strings.SplitN(pwd.Opaque, "/", 2)
			if moduleSpecifier[0] == '/' ***REMOVED***
				return &url.URL***REMOVED***Opaque: path.Join(parts[0], moduleSpecifier)***REMOVED***, nil
			***REMOVED***
			return &url.URL***REMOVED***Opaque: path.Join(parts[0], path.Join(path.Dir(parts[1]+"/"), moduleSpecifier))***REMOVED***, nil
		***REMOVED***

		// The file is in format like C:/something/path.js. But this will be decoded as scheme `C`
		// ... which is not what we want we want it to be decode as file:///C:/something/path.js
		if filepath.VolumeName(moduleSpecifier) != "" ***REMOVED***
			moduleSpecifier = "/" + moduleSpecifier
		***REMOVED***

		// we always want for the pwd to end in a slash, but filepath/path.Clean strips it so we read
		// it if it's missing
		var finalPwd = pwd
		if pwd.Opaque != "" ***REMOVED***
			if !strings.HasSuffix(pwd.Opaque, "/") ***REMOVED***
				finalPwd = &url.URL***REMOVED***Opaque: pwd.Opaque + "/"***REMOVED***
			***REMOVED***
		***REMOVED*** else if !strings.HasSuffix(pwd.Path, "/") ***REMOVED***
			finalPwd = &url.URL***REMOVED******REMOVED***
			*finalPwd = *pwd
			finalPwd.Path += "/"
		***REMOVED***
		return finalPwd.Parse(moduleSpecifier)
	***REMOVED***

	if strings.Contains(moduleSpecifier, "://") ***REMOVED***
		u, err := url.Parse(moduleSpecifier)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		if u.Scheme != "file" && u.Scheme != "https" ***REMOVED***
			return nil,
				errors.Errorf("only supported schemes for imports are file and https, %s has `%s`",
					moduleSpecifier, u.Scheme)
		***REMOVED***
		if u.Scheme == "file" && pwd.Scheme == "https" ***REMOVED***
			return nil, errors.Errorf("origin (%s) not allowed to load local file: %s", pwd, moduleSpecifier)
		***REMOVED***
		return u, err
	***REMOVED***
	// here we only care if a loader is pickable, if it is and later there is an error in the loading
	// from it we don't want to try another resolve
	_, loader, _ := pickLoader(moduleSpecifier)
	if loader == nil ***REMOVED***
		u, err := url.Parse("https://" + moduleSpecifier)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		u.Scheme = ""
		return u, nil
	***REMOVED***
	return &url.URL***REMOVED***Opaque: moduleSpecifier***REMOVED***, nil
***REMOVED***

// Dir returns the directory for the path.
func Dir(old *url.URL) *url.URL ***REMOVED***
	if old.Opaque != "" ***REMOVED*** // loader
		return &url.URL***REMOVED***Opaque: path.Join(old.Opaque, "../")***REMOVED***
	***REMOVED***
	return old.ResolveReference(&url.URL***REMOVED***Path: "./"***REMOVED***)
***REMOVED***

// Load loads the provided moduleSpecifier from the given filesystems which are map of afero.Fs
// for a given scheme which is they key of the map. If the scheme is https then a request will
// be made if the files is not found in the map and written to the map.
func Load(
	filesystems map[string]afero.Fs, moduleSpecifier *url.URL, originalModuleSpecifier string,
) (*SourceData, error) ***REMOVED***
	log.WithFields(
		log.Fields***REMOVED***
			"moduleSpecifier":          moduleSpecifier,
			"original moduleSpecifier": originalModuleSpecifier,
		***REMOVED***).Debug("Loading...")

	var pathOnFs string
	switch ***REMOVED***
	case moduleSpecifier.Opaque != "": // This is loader
		pathOnFs = filepath.Join(afero.FilePathSeparator, moduleSpecifier.Opaque)
	case moduleSpecifier.Scheme == "":
		pathOnFs = path.Clean(moduleSpecifier.String())
	default:
		pathOnFs = path.Clean(moduleSpecifier.String()[len(moduleSpecifier.Scheme)+len(":/"):])
	***REMOVED***
	scheme := moduleSpecifier.Scheme
	if scheme == "" ***REMOVED***
		scheme = "https"
	***REMOVED***

	pathOnFs, err := url.PathUnescape(filepath.FromSlash(pathOnFs))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	data, err := afero.ReadFile(filesystems[scheme], pathOnFs)

	if err != nil ***REMOVED***
		if os.IsNotExist(err) ***REMOVED***
			if scheme == "https" ***REMOVED***
				var finalModuleSpecifierURL = &url.URL***REMOVED******REMOVED***

				switch ***REMOVED***
				case moduleSpecifier.Opaque != "": // This is loader
					finalModuleSpecifierURL, err = resolveUsingLoaders(moduleSpecifier.Opaque)
					if err != nil ***REMOVED***
						return nil, err
					***REMOVED***
				case moduleSpecifier.Scheme == "":
					log.WithField("url", moduleSpecifier).Warning(
						"A url was resolved but it didn't have scheme. " +
							"This will be deprecated in the future and all remote modules will " +
							"need to explicitly use `https` as scheme")
					*finalModuleSpecifierURL = *moduleSpecifier
					finalModuleSpecifierURL.Scheme = scheme
				default:
					finalModuleSpecifierURL = moduleSpecifier
				***REMOVED***
				var result *SourceData
				result, err = loadRemoteURL(finalModuleSpecifierURL)
				if err != nil ***REMOVED***
					return nil, errors.Errorf(httpsSchemeCouldntBeLoadedMsg, originalModuleSpecifier, finalModuleSpecifierURL, err)
				***REMOVED***
				result.URL = moduleSpecifier
				// TODO maybe make an afero.Fs which makes request directly and than use CacheOnReadFs
				// on top of as with the `file` scheme fs
				_ = afero.WriteFile(filesystems[scheme], pathOnFs, result.Data, 0644)
				return result, nil
			***REMOVED***
			return nil, errors.Errorf(fileSchemeCouldntBeLoadedMsg, moduleSpecifier)
		***REMOVED***
		return nil, err
	***REMOVED***

	return &SourceData***REMOVED***URL: moduleSpecifier, Data: data***REMOVED***, nil
***REMOVED***

func resolveUsingLoaders(name string) (*url.URL, error) ***REMOVED***
	_, loader, loaderArgs := pickLoader(name)
	if loader != nil ***REMOVED***
		urlString, err := loader(name, loaderArgs)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return url.Parse(urlString)
	***REMOVED***

	return nil, errNoLoaderMatched
***REMOVED***

func loadRemoteURL(u *url.URL) (*SourceData, error) ***REMOVED***
	var oldQuery = u.RawQuery
	if u.RawQuery != "" ***REMOVED***
		u.RawQuery += "&"
	***REMOVED***
	u.RawQuery += "_k6=1"

	data, err := fetch(u.String())

	u.RawQuery = oldQuery
	// If this fails, try to fetch without ?_k6=1 - some sources act weird around unknown GET args.
	if err != nil ***REMOVED***
		data, err = fetch(u.String())
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***

	// TODO: Parse the HTML, look for meta tags!!
	// <meta name="k6-import" content="example.com/path/to/real/file.txt" />
	// <meta name="k6-import" content="github.com/myusername/repo/file.txt" />

	return &SourceData***REMOVED***URL: u, Data: data***REMOVED***, nil
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
	log.WithField("url", u).Debug("Fetching source...")
	startTime := time.Now()
	res, err := http.Get(u)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer func() ***REMOVED*** _ = res.Body.Close() ***REMOVED***()

	if res.StatusCode != 200 ***REMOVED***
		switch res.StatusCode ***REMOVED***
		case 404:
			return nil, errors.Errorf("not found: %s", u)
		default:
			return nil, errors.Errorf("wrong status code (%d) for: %s", res.StatusCode, u)
		***REMOVED***
	***REMOVED***

	data, err := ioutil.ReadAll(res.Body)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	log.WithFields(log.Fields***REMOVED***
		"url": u,
		"t":   time.Since(startTime),
		"len": len(data),
	***REMOVED***).Debug("Fetched!")
	return data, nil
***REMOVED***
