// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// DataKeyOptions specifies options for creating a new data key.
type DataKeyOptions struct ***REMOVED***
	KeyAltNames []string
	KeyMaterial []byte
	MasterKey   bsoncore.Document
***REMOVED***

// DataKey creates a new DataKeyOptions instance.
func DataKey() *DataKeyOptions ***REMOVED***
	return &DataKeyOptions***REMOVED******REMOVED***
***REMOVED***

// SetKeyAltNames specifies alternate key names.
func (dko *DataKeyOptions) SetKeyAltNames(names []string) *DataKeyOptions ***REMOVED***
	dko.KeyAltNames = names
	return dko
***REMOVED***

// SetMasterKey specifies the master key.
func (dko *DataKeyOptions) SetMasterKey(key bsoncore.Document) *DataKeyOptions ***REMOVED***
	dko.MasterKey = key
	return dko
***REMOVED***

// SetKeyMaterial specifies the key material.
func (dko *DataKeyOptions) SetKeyMaterial(keyMaterial []byte) *DataKeyOptions ***REMOVED***
	dko.KeyMaterial = keyMaterial
	return dko
***REMOVED***

// QueryType describes the type of query the result of Encrypt is used for.
type QueryType int

// These constants specify valid values for QueryType
const (
	QueryTypeEquality QueryType = 1
)

// ExplicitEncryptionOptions specifies options for configuring an explicit encryption context.
type ExplicitEncryptionOptions struct ***REMOVED***
	KeyID            *primitive.Binary
	KeyAltName       *string
	Algorithm        string
	QueryType        string
	ContentionFactor *int64
***REMOVED***

// ExplicitEncryption creates a new ExplicitEncryptionOptions instance.
func ExplicitEncryption() *ExplicitEncryptionOptions ***REMOVED***
	return &ExplicitEncryptionOptions***REMOVED******REMOVED***
***REMOVED***

// SetKeyID sets the key identifier.
func (eeo *ExplicitEncryptionOptions) SetKeyID(keyID primitive.Binary) *ExplicitEncryptionOptions ***REMOVED***
	eeo.KeyID = &keyID
	return eeo
***REMOVED***

// SetKeyAltName sets the key alternative name.
func (eeo *ExplicitEncryptionOptions) SetKeyAltName(keyAltName string) *ExplicitEncryptionOptions ***REMOVED***
	eeo.KeyAltName = &keyAltName
	return eeo
***REMOVED***

// SetAlgorithm specifies an encryption algorithm.
func (eeo *ExplicitEncryptionOptions) SetAlgorithm(algorithm string) *ExplicitEncryptionOptions ***REMOVED***
	eeo.Algorithm = algorithm
	return eeo
***REMOVED***

// SetQueryType specifies the query type.
func (eeo *ExplicitEncryptionOptions) SetQueryType(queryType string) *ExplicitEncryptionOptions ***REMOVED***
	eeo.QueryType = queryType
	return eeo
***REMOVED***

// SetContentionFactor specifies the contention factor.
func (eeo *ExplicitEncryptionOptions) SetContentionFactor(contentionFactor int64) *ExplicitEncryptionOptions ***REMOVED***
	eeo.ContentionFactor = &contentionFactor
	return eeo
***REMOVED***

// RewrapManyDataKeyOptions represents all possible options used to decrypt and encrypt all matching data keys with a
// possibly new masterKey.
type RewrapManyDataKeyOptions struct ***REMOVED***
	// Provider identifies the new KMS provider. If omitted, encrypting uses the current KMS provider.
	Provider *string

	// MasterKey identifies the new masterKey. If omitted, rewraps with the current masterKey.
	MasterKey bsoncore.Document
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
func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey bsoncore.Document) *RewrapManyDataKeyOptions ***REMOVED***
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
