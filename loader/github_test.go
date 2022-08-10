package loader

import (
	"net/url"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.k6.io/k6/lib/testutils"
)

func TestGithub(t *testing.T) ***REMOVED***
	t.Parallel()
	logger := logrus.New()
	logger.SetOutput(testutils.NewTestOutput(t))
	path := "github.com/github/gitignore/Go.gitignore"
	expectedEndSrc := "https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore"
	name, loader, parts := pickLoader(path)
	assert.Equal(t, "github", name)
	assert.Equal(t, []string***REMOVED***"github", "gitignore", "Go.gitignore"***REMOVED***, parts)
	src, err := loader(logger, path, parts)
	assert.NoError(t, err)
	assert.Equal(t, expectedEndSrc, src)

	root := &url.URL***REMOVED***Scheme: "https", Host: "example.com", Path: "/something/"***REMOVED***
	resolvedURL, err := Resolve(root, path)
	require.NoError(t, err)
	require.Empty(t, resolvedURL.Scheme)
	require.Equal(t, path, resolvedURL.Opaque)
	t.Run("not cached", func(t *testing.T) ***REMOVED***
		t.Parallel()
		data, err := Load(logger, map[string]afero.Fs***REMOVED***"https": afero.NewMemMapFs()***REMOVED***, resolvedURL, path)
		require.NoError(t, err)
		assert.Equal(t, data.URL, resolvedURL)
		assert.Equal(t, path, data.URL.String())
		assert.NotEmpty(t, data.Data)
	***REMOVED***)

	t.Run("cached", func(t *testing.T) ***REMOVED***
		t.Parallel()
		fs := afero.NewMemMapFs()
		testData := []byte("test data")

		err := afero.WriteFile(fs, "/github.com/github/gitignore/Go.gitignore", testData, 0o644)
		require.NoError(t, err)

		data, err := Load(logger, map[string]afero.Fs***REMOVED***"https": fs***REMOVED***, resolvedURL, path)
		require.NoError(t, err)
		assert.Equal(t, path, data.URL.String())
		assert.Equal(t, data.Data, testData)
	***REMOVED***)

	t.Run("relative", func(t *testing.T) ***REMOVED***
		t.Parallel()
		tests := map[string]string***REMOVED***
			"./something.else":  "github.com/github/gitignore/something.else",
			"../something.else": "github.com/github/something.else",
			"/something.else":   "github.com/something.else",
		***REMOVED***
		for relative, expected := range tests ***REMOVED***
			relativeURL, err := Resolve(Dir(resolvedURL), relative)
			require.NoError(t, err)
			assert.Equal(t, expected, relativeURL.String())
		***REMOVED***
	***REMOVED***)

	t.Run("dir", func(t *testing.T) ***REMOVED***
		t.Parallel()
		require.Equal(t, &url.URL***REMOVED***Opaque: "github.com/github/gitignore"***REMOVED***, Dir(resolvedURL))
	***REMOVED***)
***REMOVED***
