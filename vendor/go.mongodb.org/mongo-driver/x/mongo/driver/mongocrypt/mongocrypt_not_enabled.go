// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

//go:build !cse
// +build !cse

package mongocrypt

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const cseNotSupportedMsg = "client-side encryption not enabled. add the cse build tag to support"

// MongoCrypt represents a mongocrypt_t handle.
type MongoCrypt struct***REMOVED******REMOVED***

// Version returns the version string for the loaded libmongocrypt, or an empty string
// if libmongocrypt was not loaded.
func Version() string ***REMOVED***
	return ""
***REMOVED***

// NewMongoCrypt constructs a new MongoCrypt instance configured using the provided MongoCryptOptions.
func NewMongoCrypt(opts *options.MongoCryptOptions) (*MongoCrypt, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CreateEncryptionContext creates a Context to use for encryption.
func (m *MongoCrypt) CreateEncryptionContext(db string, cmd bsoncore.Document) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CreateDecryptionContext creates a Context to use for decryption.
func (m *MongoCrypt) CreateDecryptionContext(cmd bsoncore.Document) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CreateDataKeyContext creates a Context to use for creating a data key.
func (m *MongoCrypt) CreateDataKeyContext(kmsProvider string, opts *options.DataKeyOptions) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CreateExplicitEncryptionContext creates a Context to use for explicit encryption.
func (m *MongoCrypt) CreateExplicitEncryptionContext(doc bsoncore.Document, opts *options.ExplicitEncryptionOptions) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// RewrapDataKeyContext creates a Context to use for rewrapping a data key.
func (m *MongoCrypt) RewrapDataKeyContext(filter []byte, opts *options.RewrapManyDataKeyOptions) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CreateExplicitDecryptionContext creates a Context to use for explicit decryption.
func (m *MongoCrypt) CreateExplicitDecryptionContext(doc bsoncore.Document) (*Context, error) ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CryptSharedLibVersion returns the version number for the loaded crypt_shared library, or 0 if the
// crypt_shared library was not loaded.
func (m *MongoCrypt) CryptSharedLibVersion() uint64 ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// CryptSharedLibVersionString returns the version string for the loaded crypt_shared library, or an
// empty string if the crypt_shared library was not loaded.
func (m *MongoCrypt) CryptSharedLibVersionString() string ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***

// Close cleans up any resources associated with the given MongoCrypt instance.
func (m *MongoCrypt) Close() ***REMOVED***
	panic(cseNotSupportedMsg)
***REMOVED***
