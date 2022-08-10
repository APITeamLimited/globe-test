package loader

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/fsext"
	"go.k6.io/k6/lib/testutils"
)

type errorReader string

func (e errorReader) Read(_ []byte) (int, error) ***REMOVED***
	return 0, errors.New((string)(e))
***REMOVED***

var _ io.Reader = errorReader("")

func TestReadSourceSTDINError(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	_, err := ReadSource(logger, "-", "", nil, errorReader("1234"))
	require.Error(t, err)
	require.Equal(t, "1234", err.Error())
***REMOVED***

func TestReadSourceSTDINCache(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	data := []byte(`test contents`)
	r := bytes.NewReader(data)
	fs := afero.NewMemMapFs()
	sourceData, err := ReadSource(logger, "-", "/path/to/pwd",
		map[string]afero.Fs***REMOVED***"file": fsext.NewCacheOnReadFs(nil, fs, 0)***REMOVED***, r)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/-"***REMOVED***,
		Data: data,
	***REMOVED***, sourceData)
	fileData, err := afero.ReadFile(fs, "/-")
	require.NoError(t, err)
	require.Equal(t, data, fileData)
***REMOVED***

func TestReadSourceRelative(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	data := []byte(`test contents`)
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/path/to/somewhere/script.js", data, 0o644))
	sourceData, err := ReadSource(logger, "../somewhere/script.js", "/path/to/pwd", map[string]afero.Fs***REMOVED***"file": fs***REMOVED***, nil)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/path/to/somewhere/script.js"***REMOVED***,
		Data: data,
	***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceAbsolute(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	data := []byte(`test contents`)
	r := bytes.NewReader(data)
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/a/b", data, 0o644))
	require.NoError(t, afero.WriteFile(fs, "/c/a/b", []byte("wrong"), 0o644))
	sourceData, err := ReadSource(logger, "/a/b", "/c", map[string]afero.Fs***REMOVED***"file": fs***REMOVED***, r)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "file", Path: "/a/b"***REMOVED***,
		Data: data,
	***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceHttps(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	data := []byte(`test contents`)
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/github.com/something", data, 0o644))
	sourceData, err := ReadSource(logger, "https://github.com/something", "/c",
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": fs***REMOVED***, nil)
	require.NoError(t, err)
	require.Equal(t, &SourceData***REMOVED***
		URL:  &url.URL***REMOVED***Scheme: "https", Host: "github.com", Path: "/something"***REMOVED***,
		Data: data,
	***REMOVED***, sourceData)
***REMOVED***

func TestReadSourceHttpError(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	data := []byte(`test contents`)
	fs := afero.NewMemMapFs()
	require.NoError(t, afero.WriteFile(fs, "/github.com/something", data, 0o644))
	_, err := ReadSource(logger, "http://github.com/something", "/c",
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": fs***REMOVED***, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), `only supported schemes for imports are file and https`)
***REMOVED***

func TestReadSourceMissingFileError(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	fs := afero.NewMemMapFs()
	_, err := ReadSource(logger, "some file with spaces.js", "/c",
		map[string]afero.Fs***REMOVED***"file": afero.NewMemMapFs(), "https": fs***REMOVED***, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), `The moduleSpecifier "some file with spaces.js" couldn't be found on local disk.`)
***REMOVED***
