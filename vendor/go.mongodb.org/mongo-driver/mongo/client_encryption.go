// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	mcopts "go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

// ClientEncryption is used to create data keys and explicitly encrypt and decrypt BSON values.
type ClientEncryption struct ***REMOVED***
	crypt          driver.Crypt
	keyVaultClient *Client
	keyVaultColl   *Collection
***REMOVED***

// NewClientEncryption creates a new ClientEncryption instance configured with the given options.
func NewClientEncryption(keyVaultClient *Client, opts ...*options.ClientEncryptionOptions) (*ClientEncryption, error) ***REMOVED***
	if keyVaultClient == nil ***REMOVED***
		return nil, errors.New("keyVaultClient must not be nil")
	***REMOVED***

	ce := &ClientEncryption***REMOVED***
		keyVaultClient: keyVaultClient,
	***REMOVED***
	ceo := options.MergeClientEncryptionOptions(opts...)

	// create keyVaultColl
	db, coll := splitNamespace(ceo.KeyVaultNamespace)
	ce.keyVaultColl = ce.keyVaultClient.Database(db).Collection(coll, keyVaultCollOpts)

	kmsProviders, err := transformBsoncoreDocument(bson.DefaultRegistry, ceo.KmsProviders, true, "kmsProviders")
	if err != nil ***REMOVED***
		return nil, fmt.Errorf("error creating KMS providers map: %v", err)
	***REMOVED***

	mc, err := mongocrypt.NewMongoCrypt(mcopts.MongoCrypt().
		SetKmsProviders(kmsProviders).
		// Explicitly disable loading the crypt_shared library for the Crypt used for
		// ClientEncryption because it's only needed for AutoEncryption and we don't expect users to
		// have the crypt_shared library installed if they're using ClientEncryption.
		SetCryptSharedLibDisabled(true))
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// create Crypt
	kr := keyRetriever***REMOVED***coll: ce.keyVaultColl***REMOVED***
	cir := collInfoRetriever***REMOVED***client: ce.keyVaultClient***REMOVED***
	ce.crypt = driver.NewCrypt(&driver.CryptOptions***REMOVED***
		MongoCrypt: mc,
		KeyFn:      kr.cryptKeys,
		CollInfoFn: cir.cryptCollInfo,
		TLSConfig:  ceo.TLSConfig,
	***REMOVED***)

	return ce, nil
***REMOVED***

// AddKeyAltName adds a keyAltName to the keyAltNames array of the key document in the key vault collection with the
// given UUID (BSON binary subtype 0x04). Returns the previous version of the key document.
func (ce *ClientEncryption) AddKeyAltName(ctx context.Context, id primitive.Binary, keyAltName string) *SingleResult ***REMOVED***
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	keyAltNameDoc := bsoncore.NewDocumentBuilder().AppendString("keyAltNames", keyAltName).Build()
	update := bsoncore.NewDocumentBuilder().AppendDocument("$addToSet", keyAltNameDoc).Build()
	return ce.keyVaultColl.FindOneAndUpdate(ctx, filter, update)
***REMOVED***

// CreateDataKey creates a new key document and inserts into the key vault collection. Returns the _id of the created
// document as a UUID (BSON binary subtype 0x04).
func (ce *ClientEncryption) CreateDataKey(ctx context.Context, kmsProvider string,
	opts ...*options.DataKeyOptions) (primitive.Binary, error) ***REMOVED***

	// translate opts to mcopts.DataKeyOptions
	dko := options.MergeDataKeyOptions(opts...)
	co := mcopts.DataKey().SetKeyAltNames(dko.KeyAltNames)
	if dko.MasterKey != nil ***REMOVED***
		keyDoc, err := transformBsoncoreDocument(ce.keyVaultClient.registry, dko.MasterKey, true, "masterKey")
		if err != nil ***REMOVED***
			return primitive.Binary***REMOVED******REMOVED***, err
		***REMOVED***
		co.SetMasterKey(keyDoc)
	***REMOVED***
	if dko.KeyMaterial != nil ***REMOVED***
		co.SetKeyMaterial(dko.KeyMaterial)
	***REMOVED***

	// create data key document
	dataKeyDoc, err := ce.crypt.CreateDataKey(ctx, kmsProvider, co)
	if err != nil ***REMOVED***
		return primitive.Binary***REMOVED******REMOVED***, err
	***REMOVED***

	// insert key into key vault
	_, err = ce.keyVaultColl.InsertOne(ctx, dataKeyDoc)
	if err != nil ***REMOVED***
		return primitive.Binary***REMOVED******REMOVED***, err
	***REMOVED***

	subtype, data := bson.Raw(dataKeyDoc).Lookup("_id").Binary()
	return primitive.Binary***REMOVED***Subtype: subtype, Data: data***REMOVED***, nil
***REMOVED***

// Encrypt encrypts a BSON value with the given key and algorithm. Returns an encrypted value (BSON binary of subtype 6).
func (ce *ClientEncryption) Encrypt(ctx context.Context, val bson.RawValue,
	opts ...*options.EncryptOptions) (primitive.Binary, error) ***REMOVED***

	eo := options.MergeEncryptOptions(opts...)
	transformed := mcopts.ExplicitEncryption()
	if eo.KeyID != nil ***REMOVED***
		transformed.SetKeyID(*eo.KeyID)
	***REMOVED***
	if eo.KeyAltName != nil ***REMOVED***
		transformed.SetKeyAltName(*eo.KeyAltName)
	***REMOVED***
	transformed.SetAlgorithm(eo.Algorithm)
	transformed.SetQueryType(eo.QueryType)

	if eo.ContentionFactor != nil ***REMOVED***
		transformed.SetContentionFactor(*eo.ContentionFactor)
	***REMOVED***

	subtype, data, err := ce.crypt.EncryptExplicit(ctx, bsoncore.Value***REMOVED***Type: val.Type, Data: val.Value***REMOVED***, transformed)
	if err != nil ***REMOVED***
		return primitive.Binary***REMOVED******REMOVED***, err
	***REMOVED***
	return primitive.Binary***REMOVED***Subtype: subtype, Data: data***REMOVED***, nil
***REMOVED***

// Decrypt decrypts an encrypted value (BSON binary of subtype 6) and returns the original BSON value.
func (ce *ClientEncryption) Decrypt(ctx context.Context, val primitive.Binary) (bson.RawValue, error) ***REMOVED***
	decrypted, err := ce.crypt.DecryptExplicit(ctx, val.Subtype, val.Data)
	if err != nil ***REMOVED***
		return bson.RawValue***REMOVED******REMOVED***, err
	***REMOVED***

	return bson.RawValue***REMOVED***Type: decrypted.Type, Value: decrypted.Data***REMOVED***, nil
***REMOVED***

// Close cleans up any resources associated with the ClientEncryption instance. This includes disconnecting the
// key-vault Client instance.
func (ce *ClientEncryption) Close(ctx context.Context) error ***REMOVED***
	ce.crypt.Close()
	return ce.keyVaultClient.Disconnect(ctx)
***REMOVED***

// DeleteKey removes the key document with the given UUID (BSON binary subtype 0x04) from the key vault collection.
// Returns the result of the internal deleteOne() operation on the key vault collection.
func (ce *ClientEncryption) DeleteKey(ctx context.Context, id primitive.Binary) (*DeleteResult, error) ***REMOVED***
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	return ce.keyVaultColl.DeleteOne(ctx, filter)
***REMOVED***

// GetKeyByAltName returns a key document in the key vault collection with the given keyAltName.
func (ce *ClientEncryption) GetKeyByAltName(ctx context.Context, keyAltName string) *SingleResult ***REMOVED***
	filter := bsoncore.NewDocumentBuilder().AppendString("keyAltNames", keyAltName).Build()
	return ce.keyVaultColl.FindOne(ctx, filter)
***REMOVED***

// GetKey finds a single key document with the given UUID (BSON binary subtype 0x04). Returns the result of the
// internal find() operation on the key vault collection.
func (ce *ClientEncryption) GetKey(ctx context.Context, id primitive.Binary) *SingleResult ***REMOVED***
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	return ce.keyVaultColl.FindOne(ctx, filter)
***REMOVED***

// GetKeys finds all documents in the key vault collection. Returns the result of the internal find() operation on the
// key vault collection.
func (ce *ClientEncryption) GetKeys(ctx context.Context) (*Cursor, error) ***REMOVED***
	return ce.keyVaultColl.Find(ctx, bson.D***REMOVED******REMOVED***)
***REMOVED***

// RemoveKeyAltName removes a keyAltName from the keyAltNames array of the key document in the key vault collection with
// the given UUID (BSON binary subtype 0x04). Returns the previous version of the key document.
func (ce *ClientEncryption) RemoveKeyAltName(ctx context.Context, id primitive.Binary, keyAltName string) *SingleResult ***REMOVED***
	filter := bsoncore.NewDocumentBuilder().AppendBinary("_id", id.Subtype, id.Data).Build()
	update := bson.A***REMOVED***bson.D***REMOVED******REMOVED***"$set", bson.D***REMOVED******REMOVED***"keyAltNames", bson.D***REMOVED******REMOVED***"$cond", bson.A***REMOVED***bson.D***REMOVED******REMOVED***"$eq",
		bson.A***REMOVED***"$keyAltNames", bson.A***REMOVED***keyAltName***REMOVED******REMOVED******REMOVED******REMOVED***, "$$REMOVE", bson.D***REMOVED******REMOVED***"$filter",
		bson.D***REMOVED******REMOVED***"input", "$keyAltNames"***REMOVED***, ***REMOVED***"cond", bson.D***REMOVED******REMOVED***"$ne", bson.A***REMOVED***"$$this", keyAltName***REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED******REMOVED***
	return ce.keyVaultColl.FindOneAndUpdate(ctx, filter, update)
***REMOVED***

// setRewrapManyDataKeyWriteModels will prepare the WriteModel slice for a bulk updating rewrapped documents.
func setRewrapManyDataKeyWriteModels(rewrappedDocuments []bsoncore.Document, writeModels *[]WriteModel) error ***REMOVED***
	const idKey = "_id"
	const keyMaterial = "keyMaterial"
	const masterKey = "masterKey"

	if writeModels == nil ***REMOVED***
		return fmt.Errorf("writeModels pointer not set for location referenced")
	***REMOVED***

	// Append a slice of WriteModel with the update document per each rewrappedDoc _id filter.
	for _, rewrappedDocument := range rewrappedDocuments ***REMOVED***
		// Prepare the new master key for update.
		masterKeyValue, err := rewrappedDocument.LookupErr(masterKey)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		masterKeyDoc := masterKeyValue.Document()

		// Prepare the new material key for update.
		keyMaterialValue, err := rewrappedDocument.LookupErr(keyMaterial)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		keyMaterialSubtype, keyMaterialData := keyMaterialValue.Binary()
		keyMaterialBinary := primitive.Binary***REMOVED***Subtype: keyMaterialSubtype, Data: keyMaterialData***REMOVED***

		// Prepare the _id filter for documents to update.
		id, err := rewrappedDocument.LookupErr(idKey)
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		idSubtype, idData, ok := id.BinaryOK()
		if !ok ***REMOVED***
			return fmt.Errorf("expected to assert %q as binary, got type %T", idKey, id)
		***REMOVED***
		binaryID := primitive.Binary***REMOVED***Subtype: idSubtype, Data: idData***REMOVED***

		// Append the mutable document to the slice for bulk update.
		*writeModels = append(*writeModels, NewUpdateOneModel().
			SetFilter(bson.D***REMOVED******REMOVED***idKey, binaryID***REMOVED******REMOVED***).
			SetUpdate(
				bson.D***REMOVED***
					***REMOVED***"$set", bson.D***REMOVED******REMOVED***keyMaterial, keyMaterialBinary***REMOVED***, ***REMOVED***masterKey, masterKeyDoc***REMOVED******REMOVED******REMOVED***,
					***REMOVED***"$currentDate", bson.D***REMOVED******REMOVED***"updateDate", true***REMOVED******REMOVED******REMOVED***,
				***REMOVED***,
			))
	***REMOVED***
	return nil
***REMOVED***

// RewrapManyDataKey decrypts and encrypts all matching data keys with a possibly new masterKey value. For all
// matching documents, this method will overwrite the "masterKey", "updateDate", and "keyMaterial". On error, some
// matching data keys may have been rewrapped.
// libmongocrypt 1.5.2 is required. An error is returned if the detected version of libmongocrypt is less than 1.5.2.
func (ce *ClientEncryption) RewrapManyDataKey(ctx context.Context, filter interface***REMOVED******REMOVED***,
	opts ...*options.RewrapManyDataKeyOptions) (*RewrapManyDataKeyResult, error) ***REMOVED***

	// libmongocrypt versions 1.5.0 and 1.5.1 have a severe bug in RewrapManyDataKey.
	// Check if the version string starts with 1.5.0 or 1.5.1. This accounts for pre-release versions, like 1.5.0-rc0.
	libmongocryptVersion := mongocrypt.Version()
	if strings.HasPrefix(libmongocryptVersion, "1.5.0") || strings.HasPrefix(libmongocryptVersion, "1.5.1") ***REMOVED***
		return nil, fmt.Errorf("RewrapManyDataKey requires libmongocrypt 1.5.2 or newer. Detected version: %v", libmongocryptVersion)
	***REMOVED***

	rmdko := options.MergeRewrapManyDataKeyOptions(opts...)
	if ctx == nil ***REMOVED***
		ctx = context.Background()
	***REMOVED***

	// Transfer rmdko options to /x/ package options to publish the mongocrypt feed.
	co := mcopts.RewrapManyDataKey()
	if rmdko.MasterKey != nil ***REMOVED***
		keyDoc, err := transformBsoncoreDocument(ce.keyVaultClient.registry, rmdko.MasterKey, true, "masterKey")
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		co.SetMasterKey(keyDoc)
	***REMOVED***
	if rmdko.Provider != nil ***REMOVED***
		co.SetProvider(*rmdko.Provider)
	***REMOVED***

	// Prepare the filters and rewrap the data key using mongocrypt.
	filterdoc, err := transformBsoncoreDocument(ce.keyVaultClient.registry, filter, true, "filter")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rewrappedDocuments, err := ce.crypt.RewrapDataKey(ctx, filterdoc, co)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if len(rewrappedDocuments) == 0 ***REMOVED***
		// If there are no documents to rewrap, then do nothing.
		return new(RewrapManyDataKeyResult), nil
	***REMOVED***

	// Prepare the WriteModel slice for bulk updating the rewrapped data keys.
	models := []WriteModel***REMOVED******REMOVED***
	if err := setRewrapManyDataKeyWriteModels(rewrappedDocuments, &models); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	bulkWriteResults, err := ce.keyVaultColl.BulkWrite(ctx, models)
	return &RewrapManyDataKeyResult***REMOVED***BulkWriteResult: bulkWriteResults***REMOVED***, err
***REMOVED***

// splitNamespace takes a namespace in the form "database.collection" and returns (database name, collection name)
func splitNamespace(ns string) (string, string) ***REMOVED***
	firstDot := strings.Index(ns, ".")
	if firstDot == -1 ***REMOVED***
		return "", ns
	***REMOVED***

	return ns[:firstDot], ns[firstDot+1:]
***REMOVED***
