// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package bsonoptions

// TimeCodecOptions represents all possible options for time.Time encoding and decoding.
type TimeCodecOptions struct ***REMOVED***
	UseLocalTimeZone *bool // Specifies if we should decode into the local time zone. Defaults to false.
***REMOVED***

// TimeCodec creates a new *TimeCodecOptions
func TimeCodec() *TimeCodecOptions ***REMOVED***
	return &TimeCodecOptions***REMOVED******REMOVED***
***REMOVED***

// SetUseLocalTimeZone specifies if we should decode into the local time zone. Defaults to false.
func (t *TimeCodecOptions) SetUseLocalTimeZone(b bool) *TimeCodecOptions ***REMOVED***
	t.UseLocalTimeZone = &b
	return t
***REMOVED***

// MergeTimeCodecOptions combines the given *TimeCodecOptions into a single *TimeCodecOptions in a last one wins fashion.
func MergeTimeCodecOptions(opts ...*TimeCodecOptions) *TimeCodecOptions ***REMOVED***
	t := TimeCodec()
	for _, opt := range opts ***REMOVED***
		if opt == nil ***REMOVED***
			continue
		***REMOVED***
		if opt.UseLocalTimeZone != nil ***REMOVED***
			t.UseLocalTimeZone = opt.UseLocalTimeZone
		***REMOVED***
	***REMOVED***

	return t
***REMOVED***
