// Copyright 2019+ Klaus Post. All rights reserved.
// License information can be found in the LICENSE file.
// Based on work by Yann Collet, released under BSD License.

package zstd

import (
	"errors"
	"runtime"
)

// DOption is an option for creating a decoder.
type DOption func(*decoderOptions) error

// options retains accumulated state of multiple options.
type decoderOptions struct ***REMOVED***
	lowMem         bool
	concurrent     int
	maxDecodedSize uint64
	maxWindowSize  uint64
	dicts          []dict
***REMOVED***

func (o *decoderOptions) setDefault() ***REMOVED***
	*o = decoderOptions***REMOVED***
		// use less ram: true for now, but may change.
		lowMem:        true,
		concurrent:    runtime.GOMAXPROCS(0),
		maxWindowSize: MaxWindowSize,
	***REMOVED***
	o.maxDecodedSize = 1 << 63
***REMOVED***

// WithDecoderLowmem will set whether to use a lower amount of memory,
// but possibly have to allocate more while running.
func WithDecoderLowmem(b bool) DOption ***REMOVED***
	return func(o *decoderOptions) error ***REMOVED*** o.lowMem = b; return nil ***REMOVED***
***REMOVED***

// WithDecoderConcurrency will set the concurrency,
// meaning the maximum number of decoders to run concurrently.
// The value supplied must be at least 1.
// By default this will be set to GOMAXPROCS.
func WithDecoderConcurrency(n int) DOption ***REMOVED***
	return func(o *decoderOptions) error ***REMOVED***
		if n <= 0 ***REMOVED***
			return errors.New("concurrency must be at least 1")
		***REMOVED***
		o.concurrent = n
		return nil
	***REMOVED***
***REMOVED***

// WithDecoderMaxMemory allows to set a maximum decoded size for in-memory
// non-streaming operations or maximum window size for streaming operations.
// This can be used to control memory usage of potentially hostile content.
// Maximum and default is 1 << 63 bytes.
func WithDecoderMaxMemory(n uint64) DOption ***REMOVED***
	return func(o *decoderOptions) error ***REMOVED***
		if n == 0 ***REMOVED***
			return errors.New("WithDecoderMaxMemory must be at least 1")
		***REMOVED***
		if n > 1<<63 ***REMOVED***
			return errors.New("WithDecoderMaxmemory must be less than 1 << 63")
		***REMOVED***
		o.maxDecodedSize = n
		return nil
	***REMOVED***
***REMOVED***

// WithDecoderDicts allows to register one or more dictionaries for the decoder.
// If several dictionaries with the same ID is provided the last one will be used.
func WithDecoderDicts(dicts ...[]byte) DOption ***REMOVED***
	return func(o *decoderOptions) error ***REMOVED***
		for _, b := range dicts ***REMOVED***
			d, err := loadDict(b)
			if err != nil ***REMOVED***
				return err
			***REMOVED***
			o.dicts = append(o.dicts, *d)
		***REMOVED***
		return nil
	***REMOVED***
***REMOVED***

// WithDecoderMaxWindow allows to set a maximum window size for decodes.
// This allows rejecting packets that will cause big memory usage.
// The Decoder will likely allocate more memory based on the WithDecoderLowmem setting.
// If WithDecoderMaxMemory is set to a lower value, that will be used.
// Default is 512MB, Maximum is ~3.75 TB as per zstandard spec.
func WithDecoderMaxWindow(size uint64) DOption ***REMOVED***
	return func(o *decoderOptions) error ***REMOVED***
		if size < MinWindowSize ***REMOVED***
			return errors.New("WithMaxWindowSize must be at least 1KB, 1024 bytes")
		***REMOVED***
		if size > (1<<41)+7*(1<<38) ***REMOVED***
			return errors.New("WithMaxWindowSize must be less than (1<<41) + 7*(1<<38) ~ 3.75TB")
		***REMOVED***
		o.maxWindowSize = size
		return nil
	***REMOVED***
***REMOVED***
