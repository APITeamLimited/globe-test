// configdir provides access to configuration folder in each platforms.
//
// System wide configuration folders:
//
//   - Windows: %PROGRAMDATA% (C:\ProgramData)
//   - Linux/BSDs: $***REMOVED***XDG_CONFIG_DIRS***REMOVED*** (/etc/xdg)
//   - MacOSX: "/Library/Application Support"
//
// User wide configuration folders:
//
//   - Windows: %APPDATA% (C:\Users\<User>\AppData\Roaming)
//   - Linux/BSDs: $***REMOVED***XDG_CONFIG_HOME***REMOVED*** ($***REMOVED***HOME***REMOVED***/.config)
//   - MacOSX: "$***REMOVED***HOME***REMOVED***/Library/Application Support"
//
// User wide cache folders:
//
//   - Windows: %LOCALAPPDATA% (C:\Users\<User>\AppData\Local)
//   - Linux/BSDs: $***REMOVED***XDG_CACHE_HOME***REMOVED*** ($***REMOVED***HOME***REMOVED***/.cache)
//   - MacOSX: "$***REMOVED***HOME***REMOVED***/Library/Caches"
//
// configdir returns paths inside the above folders.

package configdir

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type ConfigType int

const (
	System ConfigType = iota
	Global
	All
	Existing
	Local
	Cache
)

// Config represents each folder
type Config struct ***REMOVED***
	Path string
	Type ConfigType
***REMOVED***

func (c Config) Open(fileName string) (*os.File, error) ***REMOVED***
	return os.Open(filepath.Join(c.Path, fileName))
***REMOVED***

func (c Config) Create(fileName string) (*os.File, error) ***REMOVED***
	err := c.CreateParentDir(fileName)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	return os.Create(filepath.Join(c.Path, fileName))
***REMOVED***

func (c Config) ReadFile(fileName string) ([]byte, error) ***REMOVED***
	return ioutil.ReadFile(filepath.Join(c.Path, fileName))
***REMOVED***

// CreateParentDir creates the parent directory of fileName inside c. fileName
// is a relative path inside c, containing zero or more path separators.
func (c Config) CreateParentDir(fileName string) error ***REMOVED***
	return os.MkdirAll(filepath.Dir(filepath.Join(c.Path, fileName)), 0755)
***REMOVED***

func (c Config) WriteFile(fileName string, data []byte) error ***REMOVED***
	err := c.CreateParentDir(fileName)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	return ioutil.WriteFile(filepath.Join(c.Path, fileName), data, 0644)
***REMOVED***

func (c Config) MkdirAll() error ***REMOVED***
	return os.MkdirAll(c.Path, 0755)
***REMOVED***

func (c Config) Exists(fileName string) bool ***REMOVED***
	_, err := os.Stat(filepath.Join(c.Path, fileName))
	return !os.IsNotExist(err)
***REMOVED***

// ConfigDir keeps setting for querying folders.
type ConfigDir struct ***REMOVED***
	VendorName      string
	ApplicationName string
	LocalPath       string
***REMOVED***

func New(vendorName, applicationName string) ConfigDir ***REMOVED***
	return ConfigDir***REMOVED***
		VendorName:      vendorName,
		ApplicationName: applicationName,
	***REMOVED***
***REMOVED***

func (c ConfigDir) joinPath(root string) string ***REMOVED***
	if c.VendorName != "" && hasVendorName ***REMOVED***
		return filepath.Join(root, c.VendorName, c.ApplicationName)
	***REMOVED***
	return filepath.Join(root, c.ApplicationName)
***REMOVED***

func (c ConfigDir) QueryFolders(configType ConfigType) []*Config ***REMOVED***
	if configType == Cache ***REMOVED***
		return []*Config***REMOVED***c.QueryCacheFolder()***REMOVED***
	***REMOVED***
	var result []*Config
	if c.LocalPath != "" && configType != System && configType != Global ***REMOVED***
		result = append(result, &Config***REMOVED***
			Path: c.LocalPath,
			Type: Local,
		***REMOVED***)
	***REMOVED***
	if configType != System && configType != Local ***REMOVED***
		result = append(result, &Config***REMOVED***
			Path: c.joinPath(globalSettingFolder),
			Type: Global,
		***REMOVED***)
	***REMOVED***
	if configType != Global && configType != Local ***REMOVED***
		for _, root := range systemSettingFolders ***REMOVED***
			result = append(result, &Config***REMOVED***
				Path: c.joinPath(root),
				Type: System,
			***REMOVED***)
		***REMOVED***
	***REMOVED***
	if configType != Existing ***REMOVED***
		return result
	***REMOVED***
	var existing []*Config
	for _, entry := range result ***REMOVED***
		if _, err := os.Stat(entry.Path); !os.IsNotExist(err) ***REMOVED***
			existing = append(existing, entry)
		***REMOVED***
	***REMOVED***
	return existing
***REMOVED***

func (c ConfigDir) QueryFolderContainsFile(fileName string) *Config ***REMOVED***
	configs := c.QueryFolders(Existing)
	for _, config := range configs ***REMOVED***
		if _, err := os.Stat(filepath.Join(config.Path, fileName)); !os.IsNotExist(err) ***REMOVED***
			return config
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func (c ConfigDir) QueryCacheFolder() *Config ***REMOVED***
	return &Config***REMOVED***
		Path: c.joinPath(cacheFolder),
		Type: Cache,
	***REMOVED***
***REMOVED***
