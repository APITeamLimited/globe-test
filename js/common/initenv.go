package common

import (
	"net/url"
	"path/filepath"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"go.k6.io/k6/metrics"
)

// InitEnvironment contains properties that can be accessed by Go code executed
// in the k6 init context. It can be accessed by calling common.GetInitEnv().
type InitEnvironment struct ***REMOVED***
	Logger      logrus.FieldLogger
	FileSystems map[string]afero.Fs
	CWD         *url.URL
	Registry    *metrics.Registry
	Client      *redis.Client
	// TODO: add RuntimeOptions and other properties, goja sources, etc.
	// ideally, we should leave this as the only data structure necessary for
	// executing the init context for all JS modules

***REMOVED***

// GetAbsFilePath should be used to access the FileSystems, since afero has a
// bug when opening files with relative paths - it caches them from the FS root,
// not the current working directory... So, if necessary, this method will
// transform any relative paths into absolute ones, using the CWD.
//
// TODO: refactor? It was copied from
// https://github.com/k6io/k6/blob/c51095ad7304bdd1e82cdb33c91abc331533b886/js/initcontext.go#L211-L222
func (ie *InitEnvironment) GetAbsFilePath(filename string) string ***REMOVED***
	// Here IsAbs should be enough but unfortunately it doesn't handle absolute paths starting from
	// the current drive on windows like `\users\noname\...`. Also it makes it more easy to test and
	// will probably be need for archive execution under windows if always consider '/...' as an
	// absolute path.
	if filename[0] != '/' && filename[0] != '\\' && !filepath.IsAbs(filename) ***REMOVED***
		filename = filepath.Join(ie.CWD.Path, filename)
	***REMOVED***
	filename = filepath.Clean(filename)
	if filename[0:1] != afero.FilePathSeparator ***REMOVED***
		filename = afero.FilePathSeparator + filename
	***REMOVED***
	return filename
***REMOVED***
