package httpext

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"

	"github.com/loadimpact/k6/lib"
)

// CompressionType is used to specify what compression is to be used to compress the body of a
// request
// The conversion and validation methods are auto-generated with https://github.com/alvaroloes/enumer:
//nolint: lll
//go:generate enumer -type=CompressionType -transform=snake -trimprefix CompressionType -output compression_type_gen.go
type CompressionType uint

const (
	// CompressionTypeGzip compresses through gzip
	CompressionTypeGzip CompressionType = iota
	// CompressionTypeDeflate compresses through flate
	CompressionTypeDeflate
	// CompressionTypeZstd compresses through zstd
	CompressionTypeZstd
	// CompressionTypeBr compresses through brotli
	CompressionTypeBr
	// TODO: add compress(lzw), maybe bzip2 and others listed at
	// https://en.wikipedia.org/wiki/HTTP_compression#Content-Encoding_tokens
)

func compressBody(algos []CompressionType, body io.ReadCloser) (*bytes.Buffer, string, error) ***REMOVED***
	var contentEncoding string
	var prevBuf io.Reader = body
	var buf *bytes.Buffer
	for _, compressionType := range algos ***REMOVED***
		if buf != nil ***REMOVED***
			prevBuf = buf
		***REMOVED***
		buf = new(bytes.Buffer)

		if contentEncoding != "" ***REMOVED***
			contentEncoding += ", "
		***REMOVED***
		contentEncoding += compressionType.String()
		var w io.WriteCloser
		switch compressionType ***REMOVED***
		case CompressionTypeGzip:
			w = gzip.NewWriter(buf)
		case CompressionTypeDeflate:
			w = zlib.NewWriter(buf)
		case CompressionTypeZstd:
			w, _ = zstd.NewWriter(buf)
		case CompressionTypeBr:
			w = brotli.NewWriter(buf)
		default:
			return nil, "", fmt.Errorf("unknown compressionType %s", compressionType)
		***REMOVED***
		// we don't close in defer because zlib will write it's checksum again if it closes twice :(
		var _, err = io.Copy(w, prevBuf)
		if err != nil ***REMOVED***
			_ = w.Close()
			return nil, "", err
		***REMOVED***

		if err = w.Close(); err != nil ***REMOVED***
			return nil, "", err
		***REMOVED***
	***REMOVED***

	return buf, contentEncoding, body.Close()
***REMOVED***

//nolint:gochecknoglobals
var decompressionErrors = [...]error***REMOVED***
	zlib.ErrChecksum, zlib.ErrDictionary, zlib.ErrHeader,
	gzip.ErrChecksum, gzip.ErrHeader,
	// TODO: handle brotli errors - currently unexported
	zstd.ErrReservedBlockType, zstd.ErrCompressedSizeTooBig, zstd.ErrBlockTooSmall, zstd.ErrMagicMismatch,
	zstd.ErrWindowSizeExceeded, zstd.ErrWindowSizeTooSmall, zstd.ErrDecoderSizeExceeded, zstd.ErrUnknownDictionary,
	zstd.ErrFrameSizeExceeded, zstd.ErrCRCMismatch, zstd.ErrDecoderClosed,
***REMOVED***

func newDecompressionError(originalErr error) K6Error ***REMOVED***
	return NewK6Error(
		responseDecompressionErrorCode,
		fmt.Sprintf("error decompressing response body (%s)", originalErr.Error()),
		originalErr,
	)
***REMOVED***

func wrapDecompressionError(err error) error ***REMOVED***
	if err == nil ***REMOVED***
		return nil
	***REMOVED***

	// TODO: something more optimized? for example, we won't get zstd errors if
	// we don't use it... maybe the code that builds the decompression readers
	// could also add an appropriate error-wrapper layer?
	for _, decErr := range decompressionErrors ***REMOVED***
		if err == decErr ***REMOVED***
			return newDecompressionError(err)
		***REMOVED***
	***REMOVED***
	if strings.HasPrefix(err.Error(), "brotli: ") ***REMOVED*** // TODO: submit an upstream patch and fix...
		return newDecompressionError(err)
	***REMOVED***
	return err
***REMOVED***

func readResponseBody(
	state *lib.State,
	respType ResponseType,
	resp *http.Response,
	respErr error,
) (interface***REMOVED******REMOVED***, error) ***REMOVED***
	if resp == nil || respErr != nil ***REMOVED***
		return nil, respErr
	***REMOVED***

	if respType == ResponseTypeNone ***REMOVED***
		_, err := io.Copy(ioutil.Discard, resp.Body)
		_ = resp.Body.Close()
		if err != nil ***REMOVED***
			respErr = err
		***REMOVED***
		return nil, respErr
	***REMOVED***

	rc := &readCloser***REMOVED***resp.Body***REMOVED***
	// Ensure that the entire response body is read and closed, e.g. in case of decoding errors
	defer func(respBody io.ReadCloser) ***REMOVED***
		_, _ = io.Copy(ioutil.Discard, respBody)
		_ = respBody.Close()
	***REMOVED***(resp.Body)

	contentEncodings := strings.Split(resp.Header.Get("Content-Encoding"), ",")
	// Transparently decompress the body if it's has a content-encoding we
	// support. If not, simply return it as it is.
	for i := len(contentEncodings) - 1; i >= 0; i-- ***REMOVED***
		contentEncoding := strings.TrimSpace(contentEncodings[i])
		if compression, err := CompressionTypeString(contentEncoding); err == nil ***REMOVED***
			var decoder io.Reader
			var err error
			switch compression ***REMOVED***
			case CompressionTypeDeflate:
				decoder, err = zlib.NewReader(rc)
			case CompressionTypeGzip:
				decoder, err = gzip.NewReader(rc)
			case CompressionTypeZstd:
				decoder, err = zstd.NewReader(rc)
			case CompressionTypeBr:
				decoder = brotli.NewReader(rc)
			default:
				// We have not implemented a compression ... :(
				err = fmt.Errorf(
					"unsupported compression type %s - this is a bug in k6, please report it",
					compression,
				)
			***REMOVED***
			if err != nil ***REMOVED***
				return nil, newDecompressionError(err)
			***REMOVED***
			rc = &readCloser***REMOVED***decoder***REMOVED***
		***REMOVED***
	***REMOVED***
	buf := state.BPool.Get()
	defer state.BPool.Put(buf)
	buf.Reset()
	_, err := io.Copy(buf, rc.Reader)
	if err != nil ***REMOVED***
		respErr = wrapDecompressionError(err)
	***REMOVED***

	err = rc.Close()
	if err != nil && respErr == nil ***REMOVED*** // Don't overwrite previous errors
		respErr = wrapDecompressionError(err)
	***REMOVED***

	var result interface***REMOVED******REMOVED***
	// Binary or string
	switch respType ***REMOVED***
	case ResponseTypeText:
		result = buf.String()
	case ResponseTypeBinary:
		// Copy the data to a new slice before we return the buffer to the pool,
		// because buf.Bytes() points to the underlying buffer byte slice.
		binData := make([]byte, buf.Len())
		copy(binData, buf.Bytes())
		result = binData
	default:
		respErr = fmt.Errorf("unknown responseType %s", respType)
	***REMOVED***

	return result, respErr
***REMOVED***
