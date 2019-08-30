package httpext

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

type reader func([]byte) (int, error)

func (r reader) Read(a []byte) (int, error) ***REMOVED***
	return ((func([]byte) (int, error))(r))(a)
***REMOVED***

const badReadMsg = "bad read error for test"
const badCloseMsg = "bad close error for test"

func badReadBody() io.Reader ***REMOVED***
	return reader(func(_ []byte) (int, error) ***REMOVED***
		return 0, errors.New(badReadMsg)
	***REMOVED***)
***REMOVED***

type closer func() error

func (c closer) Close() error ***REMOVED***
	return ((func() error)(c))()
***REMOVED***

func badCloseBody() io.ReadCloser ***REMOVED***
	return struct ***REMOVED***
		io.Reader
		io.Closer
	***REMOVED******REMOVED***
		Reader: reader(func(_ []byte) (int, error) ***REMOVED***
			return 0, io.EOF
		***REMOVED***),
		Closer: closer(func() error ***REMOVED***
			return errors.New(badCloseMsg)
		***REMOVED***),
	***REMOVED***
***REMOVED***

func TestCompressionBodyError(t *testing.T) ***REMOVED***
	var algos = []CompressionType***REMOVED***CompressionTypeGzip***REMOVED***
	t.Run("bad read body", func(t *testing.T) ***REMOVED***
		_, _, err := compressBody(algos, ioutil.NopCloser(badReadBody()))
		require.Error(t, err)
		require.Equal(t, err.Error(), badReadMsg)
	***REMOVED***)

	t.Run("bad close body", func(t *testing.T) ***REMOVED***
		_, _, err := compressBody(algos, badCloseBody())
		require.Error(t, err)
		require.Equal(t, err.Error(), badCloseMsg)
	***REMOVED***)
***REMOVED***

func TestMakeRequestError(t *testing.T) ***REMOVED***
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	t.Run("bad compression algorithm body", func(t *testing.T) ***REMOVED***
		var req, err = http.NewRequest("GET", "https://wont.be.used", nil)

		require.NoError(t, err)
		var badCompressionType = CompressionType(13)
		require.False(t, badCompressionType.IsACompressionType())
		var preq = &ParsedHTTPRequest***REMOVED***
			Req:          req,
			Body:         new(bytes.Buffer),
			Compressions: []CompressionType***REMOVED***badCompressionType***REMOVED***,
		***REMOVED***
		_, err = MakeRequest(ctx, preq)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unknown compressionType CompressionType(13)")
	***REMOVED***)
***REMOVED***

func TestURL(t *testing.T) ***REMOVED***
	t.Run("Clean", func(t *testing.T) ***REMOVED***
		testCases := []struct ***REMOVED***
			url      string
			expected string
		***REMOVED******REMOVED***
			***REMOVED***"https://example.com/", "https://example.com/"***REMOVED***,
			***REMOVED***"https://example.com/$***REMOVED******REMOVED***", "https://example.com/$***REMOVED******REMOVED***"***REMOVED***,
			***REMOVED***"https://user@example.com/", "https://****@example.com/"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/", "https://****:****@example.com/"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/path?a=1&b=2", "https://****:****@example.com/path?a=1&b=2"***REMOVED***,
			***REMOVED***"https://user:pass@example.com/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***", "https://****:****@example.com/$***REMOVED******REMOVED***/$***REMOVED******REMOVED***"***REMOVED***,
			***REMOVED***"@malformed/url", "@malformed/url"***REMOVED***,
			***REMOVED***"not a url", "not a url"***REMOVED***,
		***REMOVED***

		for _, tc := range testCases ***REMOVED***
			tc := tc
			t.Run(tc.url, func(t *testing.T) ***REMOVED***
				u, err := url.Parse(tc.url)
				require.NoError(t, err)
				ut := URL***REMOVED***u: u, URL: tc.url***REMOVED***
				require.Equal(t, tc.expected, ut.Clean())
			***REMOVED***)
		***REMOVED***
	***REMOVED***)
***REMOVED***
