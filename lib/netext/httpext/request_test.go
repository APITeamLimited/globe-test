package httpext

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
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
