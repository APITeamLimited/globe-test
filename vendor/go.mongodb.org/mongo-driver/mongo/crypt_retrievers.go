// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// keyRetriever gets keys from the key vault collection.
type keyRetriever struct ***REMOVED***
	coll *Collection
***REMOVED***

func (kr *keyRetriever) cryptKeys(ctx context.Context, filter bsoncore.Document) ([]bsoncore.Document, error) ***REMOVED***
	// Remove the explicit session from the context if one is set.
	// The explicit session may be from a different client.
	ctx = NewSessionContext(ctx, nil)
	cursor, err := kr.coll.Find(ctx, filter)
	if err != nil ***REMOVED***
		return nil, EncryptionKeyVaultError***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***
	defer cursor.Close(ctx)

	var results []bsoncore.Document
	for cursor.Next(ctx) ***REMOVED***
		cur := make([]byte, len(cursor.Current))
		copy(cur, cursor.Current)
		results = append(results, cur)
	***REMOVED***
	if err = cursor.Err(); err != nil ***REMOVED***
		return nil, EncryptionKeyVaultError***REMOVED***Wrapped: err***REMOVED***
	***REMOVED***

	return results, nil
***REMOVED***

// collInfoRetriever gets info for collections from a database.
type collInfoRetriever struct ***REMOVED***
	client *Client
***REMOVED***

func (cir *collInfoRetriever) cryptCollInfo(ctx context.Context, db string, filter bsoncore.Document) (bsoncore.Document, error) ***REMOVED***
	// Remove the explicit session from the context if one is set.
	// The explicit session may be from a different client.
	ctx = NewSessionContext(ctx, nil)
	cursor, err := cir.client.Database(db).ListCollections(ctx, filter)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	defer cursor.Close(ctx)

	if !cursor.Next(ctx) ***REMOVED***
		return nil, cursor.Err()
	***REMOVED***

	res := make([]byte, len(cursor.Current))
	copy(res, cursor.Current)
	return res, nil
***REMOVED***
