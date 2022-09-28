// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package driver

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt"
	"go.mongodb.org/mongo-driver/x/mongo/driver/mongocrypt/options"
)

const (
	defaultKmsPort    = 443
	defaultKmsTimeout = 10 * time.Second
)

// CollectionInfoFn is a callback used to retrieve collection information.
type CollectionInfoFn func(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error)

// KeyRetrieverFn is a callback used to retrieve keys from the key vault.
type KeyRetrieverFn func(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error)

// MarkCommandFn is a callback used to add encryption markings to a command.
type MarkCommandFn func(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)

// CryptOptions specifies options to configure a Crypt instance.
type CryptOptions struct ***REMOVED***
	MongoCrypt           *mongocrypt.MongoCrypt
	CollInfoFn           CollectionInfoFn
	KeyFn                KeyRetrieverFn
	MarkFn               MarkCommandFn
	TLSConfig            map[string]*tls.Config
	BypassAutoEncryption bool
	BypassQueryAnalysis  bool
***REMOVED***

// Crypt is an interface implemented by types that can encrypt and decrypt instances of
// bsoncore.Document.
//
// Users should rely on the driver's crypt type (used by default) for encryption and decryption
// unless they are perfectly confident in another implementation of Crypt.
type Crypt interface ***REMOVED***
	// Encrypt encrypts the given command.
	Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error)
	// Decrypt decrypts the given command response.
	Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error)
	// CreateDataKey creates a data key using the given KMS provider and options.
	CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error)
	// EncryptExplicit encrypts the given value with the given options.
	EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error)
	// DecryptExplicit decrypts the given encrypted value.
	DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error)
	// Close cleans up any resources associated with the Crypt instance.
	Close()
	// BypassAutoEncryption returns true if auto-encryption should be bypassed.
	BypassAutoEncryption() bool
	// RewrapDataKey attempts to rewrap the document data keys matching the filter, preparing the re-wrapped documents
	// to be returned as a slice of bsoncore.Document.
	RewrapDataKey(ctx context.Context, filter []byte, opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error)
***REMOVED***

// crypt consumes the libmongocrypt.MongoCrypt type to iterate the mongocrypt state machine and perform encryption
// and decryption.
type crypt struct ***REMOVED***
	mongoCrypt *mongocrypt.MongoCrypt
	collInfoFn CollectionInfoFn
	keyFn      KeyRetrieverFn
	markFn     MarkCommandFn
	tlsConfig  map[string]*tls.Config

	bypassAutoEncryption bool
***REMOVED***

// NewCrypt creates a new Crypt instance configured with the given AutoEncryptionOptions.
func NewCrypt(opts *CryptOptions) Crypt ***REMOVED***
	return &crypt***REMOVED***
		mongoCrypt:           opts.MongoCrypt,
		collInfoFn:           opts.CollInfoFn,
		keyFn:                opts.KeyFn,
		markFn:               opts.MarkFn,
		tlsConfig:            opts.TLSConfig,
		bypassAutoEncryption: opts.BypassAutoEncryption,
	***REMOVED***
***REMOVED***

// Encrypt encrypts the given command.
func (c *crypt) Encrypt(ctx context.Context, db string, cmd bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	if c.bypassAutoEncryption ***REMOVED***
		return cmd, nil
	***REMOVED***

	cryptCtx, err := c.mongoCrypt.CreateEncryptionContext(db, cmd)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, db)
***REMOVED***

// Decrypt decrypts the given command response.
func (c *crypt) Decrypt(ctx context.Context, cmdResponse bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	cryptCtx, err := c.mongoCrypt.CreateDecryptionContext(cmdResponse)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
***REMOVED***

// CreateDataKey creates a data key using the given KMS provider and options.
func (c *crypt) CreateDataKey(ctx context.Context, kmsProvider string, opts *options.DataKeyOptions) (bsoncore.Document, error) ***REMOVED***
	cryptCtx, err := c.mongoCrypt.CreateDataKeyContext(kmsProvider, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cryptCtx.Close()

	return c.executeStateMachine(ctx, cryptCtx, "")
***REMOVED***

// RewrapDataKey attempts to rewrap the document data keys matching the filter, preparing the re-wrapped documents to
// be returned as a slice of bsoncore.Document.
func (c *crypt) RewrapDataKey(ctx context.Context, filter []byte,
	opts *options.RewrapManyDataKeyOptions) ([]bsoncore.Document, error) ***REMOVED***

	cryptCtx, err := c.mongoCrypt.RewrapDataKeyContext(filter, opts)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cryptCtx.Close()

	rewrappedBSON, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if rewrappedBSON == nil ***REMOVED***
		return nil, nil
	***REMOVED***

	// mongocrypt_ctx_rewrap_many_datakey_init wraps the documents in a BSON of the form ***REMOVED*** "v": [(BSON document), ...] ***REMOVED***
	// where each BSON document in the slice is a document containing a rewrapped datakey.
	rewrappedDocumentBytes, err := rewrappedBSON.LookupErr("v")
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Parse the resulting BSON as individual documents.
	rewrappedDocsArray, ok := rewrappedDocumentBytes.ArrayOK()
	if !ok ***REMOVED***
		return nil, fmt.Errorf("expected results from mongocrypt_ctx_rewrap_many_datakey_init to be an array")
	***REMOVED***

	rewrappedDocumentValues, err := rewrappedDocsArray.Values()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	rewrappedDocuments := []bsoncore.Document***REMOVED******REMOVED***
	for _, rewrappedDocumentValue := range rewrappedDocumentValues ***REMOVED***
		if rewrappedDocumentValue.Type != bsontype.EmbeddedDocument ***REMOVED***
			// If a value in the document's array returned by mongocrypt is anything other than an embedded document,
			// then something is wrong and we should terminate the routine.
			return nil, fmt.Errorf("expected value of type %q, got: %q",
				bsontype.EmbeddedDocument.String(),
				rewrappedDocumentValue.Type.String())
		***REMOVED***
		rewrappedDocuments = append(rewrappedDocuments, rewrappedDocumentValue.Document())
	***REMOVED***
	return rewrappedDocuments, nil
***REMOVED***

// EncryptExplicit encrypts the given value with the given options.
func (c *crypt) EncryptExplicit(ctx context.Context, val bsoncore.Value, opts *options.ExplicitEncryptionOptions) (byte, []byte, error) ***REMOVED***
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendValueElement(doc, "v", val)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitEncryptionContext(doc, opts)
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil ***REMOVED***
		return 0, nil, err
	***REMOVED***

	sub, data := res.Lookup("v").Binary()
	return sub, data, nil
***REMOVED***

// DecryptExplicit decrypts the given encrypted value.
func (c *crypt) DecryptExplicit(ctx context.Context, subtype byte, data []byte) (bsoncore.Value, error) ***REMOVED***
	idx, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendBinaryElement(doc, "v", subtype, data)
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)

	cryptCtx, err := c.mongoCrypt.CreateExplicitDecryptionContext(doc)
	if err != nil ***REMOVED***
		return bsoncore.Value***REMOVED******REMOVED***, err
	***REMOVED***
	defer cryptCtx.Close()

	res, err := c.executeStateMachine(ctx, cryptCtx, "")
	if err != nil ***REMOVED***
		return bsoncore.Value***REMOVED******REMOVED***, err
	***REMOVED***

	return res.Lookup("v"), nil
***REMOVED***

// Close cleans up any resources associated with the Crypt instance.
func (c *crypt) Close() ***REMOVED***
	c.mongoCrypt.Close()
***REMOVED***

func (c *crypt) BypassAutoEncryption() bool ***REMOVED***
	return c.bypassAutoEncryption
***REMOVED***

func (c *crypt) executeStateMachine(ctx context.Context, cryptCtx *mongocrypt.Context, db string) (bsoncore.Document, error) ***REMOVED***
	var err error
	for ***REMOVED***
		state := cryptCtx.State()
		switch state ***REMOVED***
		case mongocrypt.NeedMongoCollInfo:
			err = c.collectionInfo(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoMarkings:
			err = c.markCommand(ctx, cryptCtx, db)
		case mongocrypt.NeedMongoKeys:
			err = c.retrieveKeys(ctx, cryptCtx)
		case mongocrypt.NeedKms:
			err = c.decryptKeys(cryptCtx)
		case mongocrypt.Ready:
			return cryptCtx.Finish()
		case mongocrypt.Done:
			return nil, nil
		default:
			return nil, fmt.Errorf("invalid Crypt state: %v", state)
		***REMOVED***
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
	***REMOVED***
***REMOVED***

func (c *crypt) collectionInfo(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error ***REMOVED***
	op, err := cryptCtx.NextOperation()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	collInfo, err := c.collInfoFn(ctx, db, op)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if collInfo != nil ***REMOVED***
		if err = cryptCtx.AddOperationResult(collInfo); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return cryptCtx.CompleteOperation()
***REMOVED***

func (c *crypt) markCommand(ctx context.Context, cryptCtx *mongocrypt.Context, db string) error ***REMOVED***
	op, err := cryptCtx.NextOperation()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	markedCmd, err := c.markFn(ctx, db, op)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = cryptCtx.AddOperationResult(markedCmd); err != nil ***REMOVED***
		return err
	***REMOVED***

	return cryptCtx.CompleteOperation()
***REMOVED***

func (c *crypt) retrieveKeys(ctx context.Context, cryptCtx *mongocrypt.Context) error ***REMOVED***
	op, err := cryptCtx.NextOperation()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	keys, err := c.keyFn(ctx, op)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	for _, key := range keys ***REMOVED***
		if err = cryptCtx.AddOperationResult(key); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return cryptCtx.CompleteOperation()
***REMOVED***

func (c *crypt) decryptKeys(cryptCtx *mongocrypt.Context) error ***REMOVED***
	for ***REMOVED***
		kmsCtx := cryptCtx.NextKmsContext()
		if kmsCtx == nil ***REMOVED***
			break
		***REMOVED***

		if err := c.decryptKey(kmsCtx); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	return cryptCtx.FinishKmsContexts()
***REMOVED***

func (c *crypt) decryptKey(kmsCtx *mongocrypt.KmsContext) error ***REMOVED***
	host, err := kmsCtx.HostName()
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	msg, err := kmsCtx.Message()
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// add a port to the address if it's not already present
	addr := host
	if idx := strings.IndexByte(host, ':'); idx == -1 ***REMOVED***
		addr = fmt.Sprintf("%s:%d", host, defaultKmsPort)
	***REMOVED***

	kmsProvider := kmsCtx.KMSProvider()
	tlsCfg := c.tlsConfig[kmsProvider]
	if tlsCfg == nil ***REMOVED***
		tlsCfg = &tls.Config***REMOVED***MinVersion: tls.VersionTLS12***REMOVED***
	***REMOVED***
	conn, err := tls.Dial("tcp", addr, tlsCfg)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		_ = conn.Close()
	***REMOVED***()

	if err = conn.SetWriteDeadline(time.Now().Add(defaultKmsTimeout)); err != nil ***REMOVED***
		return err
	***REMOVED***
	if _, err = conn.Write(msg); err != nil ***REMOVED***
		return err
	***REMOVED***

	for ***REMOVED***
		bytesNeeded := kmsCtx.BytesNeeded()
		if bytesNeeded == 0 ***REMOVED***
			return nil
		***REMOVED***

		res := make([]byte, bytesNeeded)
		bytesRead, err := conn.Read(res)
		if err != nil && err != io.EOF ***REMOVED***
			return err
		***REMOVED***

		if err = kmsCtx.FeedResponse(res[:bytesRead]); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
***REMOVED***
