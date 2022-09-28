// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"sync"

	"github.com/golang/snappy"
	"github.com/klauspost/compress/zstd"
	"go.mongodb.org/mongo-driver/x/mongo/driver/wiremessage"
)

// CompressionOpts holds settings for how to compress a payload
type CompressionOpts struct ***REMOVED***
	Compressor       wiremessage.CompressorID
	ZlibLevel        int
	ZstdLevel        int
	UncompressedSize int32
***REMOVED***

var zstdEncoders = &sync.Map***REMOVED******REMOVED***

func getZstdEncoder(l zstd.EncoderLevel) (*zstd.Encoder, error) ***REMOVED***
	v, ok := zstdEncoders.Load(l)
	if ok ***REMOVED***
		return v.(*zstd.Encoder), nil
	***REMOVED***
	encoder, err := zstd.NewWriter(nil, zstd.WithEncoderLevel(l))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	zstdEncoders.Store(l, encoder)
	return encoder, nil
***REMOVED***

// CompressPayload takes a byte slice and compresses it according to the options passed
func CompressPayload(in []byte, opts CompressionOpts) ([]byte, error) ***REMOVED***
	switch opts.Compressor ***REMOVED***
	case wiremessage.CompressorNoOp:
		return in, nil
	case wiremessage.CompressorSnappy:
		return snappy.Encode(nil, in), nil
	case wiremessage.CompressorZLib:
		var b bytes.Buffer
		w, err := zlib.NewWriterLevel(&b, opts.ZlibLevel)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		_, err = w.Write(in)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		err = w.Close()
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return b.Bytes(), nil
	case wiremessage.CompressorZstd:
		encoder, err := getZstdEncoder(zstd.EncoderLevelFromZstd(opts.ZstdLevel))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return encoder.EncodeAll(in, nil), nil
	default:
		return nil, fmt.Errorf("unknown compressor ID %v", opts.Compressor)
	***REMOVED***
***REMOVED***

// DecompressPayload takes a byte slice that has been compressed and undoes it according to the options passed
func DecompressPayload(in []byte, opts CompressionOpts) ([]byte, error) ***REMOVED***
	switch opts.Compressor ***REMOVED***
	case wiremessage.CompressorNoOp:
		return in, nil
	case wiremessage.CompressorSnappy:
		uncompressed := make([]byte, opts.UncompressedSize)
		return snappy.Decode(uncompressed, in)
	case wiremessage.CompressorZLib:
		decompressor, err := zlib.NewReader(bytes.NewReader(in))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		uncompressed := make([]byte, opts.UncompressedSize)
		_, err = io.ReadFull(decompressor, uncompressed)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return uncompressed, nil
	case wiremessage.CompressorZstd:
		w, err := zstd.NewReader(bytes.NewBuffer(in))
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		defer w.Close()
		var b bytes.Buffer
		_, err = io.Copy(&b, w)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		return b.Bytes(), nil
	default:
		return nil, fmt.Errorf("unknown compressor ID %v", opts.Compressor)
	***REMOVED***
***REMOVED***
