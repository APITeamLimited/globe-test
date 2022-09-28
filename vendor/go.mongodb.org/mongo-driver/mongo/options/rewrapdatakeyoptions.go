// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// RewrapManyDataKeyOptions represents all possible options used to decrypt and encrypt all matching data keys with a
// possibly new masterKey.
type RewrapManyDataKeyOptions struct ***REMOVED***
	// Provider identifies the new KMS provider. If omitted, encrypting uses the current KMS provider.
	Provider *string

	// MasterKey identifies the new masterKey. If omitted, rewraps with the current masterKey.
	MasterKey interface***REMOVED******REMOVED***
***REMOVED***

// RewrapManyDataKey creates a new RewrapManyDataKeyOptions instance.
func RewrapManyDataKey() *RewrapManyDataKeyOptions ***REMOVED***
	return new(RewrapManyDataKeyOptions)
***REMOVED***

// SetProvider sets the value for the Provider field.
func (rmdko *RewrapManyDataKeyOptions) SetProvider(provider string) *RewrapManyDataKeyOptions ***REMOVED***
	rmdko.Provider = &provider
	return rmdko
***REMOVED***

// SetMasterKey sets the value for the MasterKey field.
func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey interface***REMOVED******REMOVED***) *RewrapManyDataKeyOptions ***REMOVED***
	rmdko.MasterKey = masterKey
	return rmdko
***REMOVED***

// MergeRewrapManyDataKeyOptions combines the given RewrapManyDataKeyOptions instances into a single
// RewrapManyDataKeyOptions in a last one wins fashion.
func MergeRewrapManyDataKeyOptions(opts ...*RewrapManyDataKeyOptions) *RewrapManyDataKeyOptions ***REMOVED***
	rmdkOpts := RewrapManyDataKey()
	for _, rmdko := range opts ***REMOVED***
		if rmdko == nil ***REMOVED***
			continue
		***REMOVED***
		if provider := rmdko.Provider; provider != nil ***REMOVED***
			rmdkOpts.Provider = provider
		***REMOVED***
		if masterKey := rmdko.MasterKey; masterKey != nil ***REMOVED***
			rmdkOpts.MasterKey = masterKey
		***REMOVED***
	***REMOVED***
	return rmdkOpts
***REMOVED***
