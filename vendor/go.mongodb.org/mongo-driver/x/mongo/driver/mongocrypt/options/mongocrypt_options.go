// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MongoCryptOptions specifies options to configure a MongoCrypt instance.
type MongoCryptOptions struct ***REMOVED***
	KmsProviders               bsoncore.Document
	LocalSchemaMap             map[string]bsoncore.Document
	BypassQueryAnalysis        bool
	EncryptedFieldsMap         map[string]bsoncore.Document
	CryptSharedLibDisabled     bool
	CryptSharedLibOverridePath string
***REMOVED***

// MongoCrypt creates a new MongoCryptOptions instance.
func MongoCrypt() *MongoCryptOptions ***REMOVED***
	return &MongoCryptOptions***REMOVED******REMOVED***
***REMOVED***

// SetKmsProviders specifies the KMS providers map.
func (mo *MongoCryptOptions) SetKmsProviders(kmsProviders bsoncore.Document) *MongoCryptOptions ***REMOVED***
	mo.KmsProviders = kmsProviders
	return mo
***REMOVED***

// SetLocalSchemaMap specifies the local schema map.
func (mo *MongoCryptOptions) SetLocalSchemaMap(localSchemaMap map[string]bsoncore.Document) *MongoCryptOptions ***REMOVED***
	mo.LocalSchemaMap = localSchemaMap
	return mo
***REMOVED***

// SetBypassQueryAnalysis skips the NeedMongoMarkings state.
func (mo *MongoCryptOptions) SetBypassQueryAnalysis(bypassQueryAnalysis bool) *MongoCryptOptions ***REMOVED***
	mo.BypassQueryAnalysis = bypassQueryAnalysis
	return mo
***REMOVED***

// SetEncryptedFieldsMap specifies the encrypted fields map.
func (mo *MongoCryptOptions) SetEncryptedFieldsMap(efcMap map[string]bsoncore.Document) *MongoCryptOptions ***REMOVED***
	mo.EncryptedFieldsMap = efcMap
	return mo
***REMOVED***

// SetCryptSharedLibDisabled explicitly disables loading the crypt_shared library if set to true.
func (mo *MongoCryptOptions) SetCryptSharedLibDisabled(disabled bool) *MongoCryptOptions ***REMOVED***
	mo.CryptSharedLibDisabled = disabled
	return mo
***REMOVED***

// SetCryptSharedLibOverridePath sets the override path to the crypt_shared library file. Setting
// an override path disables the default operating system dynamic library search path.
func (mo *MongoCryptOptions) SetCryptSharedLibOverridePath(path string) *MongoCryptOptions ***REMOVED***
	mo.CryptSharedLibOverridePath = path
	return mo
***REMOVED***
