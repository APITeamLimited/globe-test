// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package gridfs // import "go.mongodb.org/mongo-driver/mongo/gridfs"

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// TODO: add sessions options

// DefaultChunkSize is the default size of each file chunk.
const DefaultChunkSize int32 = 255 * 1024 // 255 KiB

// ErrFileNotFound occurs if a user asks to download a file with a file ID that isn't found in the files collection.
var ErrFileNotFound = errors.New("file with given parameters not found")

// ErrMissingChunkSize occurs when downloading a file if the files collection document is missing the "chunkSize" field.
var ErrMissingChunkSize = errors.New("files collection document does not contain a 'chunkSize' field")

// Bucket represents a GridFS bucket.
type Bucket struct ***REMOVED***
	db         *mongo.Database
	chunksColl *mongo.Collection // collection to store file chunks
	filesColl  *mongo.Collection // collection to store file metadata

	name      string
	chunkSize int32
	wc        *writeconcern.WriteConcern
	rc        *readconcern.ReadConcern
	rp        *readpref.ReadPref

	firstWriteDone bool
	readBuf        []byte
	writeBuf       []byte

	readDeadline  time.Time
	writeDeadline time.Time
***REMOVED***

// Upload contains options to upload a file to a bucket.
type Upload struct ***REMOVED***
	chunkSize int32
	metadata  bson.D
***REMOVED***

// NewBucket creates a GridFS bucket.
func NewBucket(db *mongo.Database, opts ...*options.BucketOptions) (*Bucket, error) ***REMOVED***
	b := &Bucket***REMOVED***
		name:      "fs",
		chunkSize: DefaultChunkSize,
		db:        db,
		wc:        db.WriteConcern(),
		rc:        db.ReadConcern(),
		rp:        db.ReadPreference(),
	***REMOVED***

	bo := options.MergeBucketOptions(opts...)
	if bo.Name != nil ***REMOVED***
		b.name = *bo.Name
	***REMOVED***
	if bo.ChunkSizeBytes != nil ***REMOVED***
		b.chunkSize = *bo.ChunkSizeBytes
	***REMOVED***
	if bo.WriteConcern != nil ***REMOVED***
		b.wc = bo.WriteConcern
	***REMOVED***
	if bo.ReadConcern != nil ***REMOVED***
		b.rc = bo.ReadConcern
	***REMOVED***
	if bo.ReadPreference != nil ***REMOVED***
		b.rp = bo.ReadPreference
	***REMOVED***

	var collOpts = options.Collection().SetWriteConcern(b.wc).SetReadConcern(b.rc).SetReadPreference(b.rp)

	b.chunksColl = db.Collection(b.name+".chunks", collOpts)
	b.filesColl = db.Collection(b.name+".files", collOpts)
	b.readBuf = make([]byte, b.chunkSize)
	b.writeBuf = make([]byte, b.chunkSize)

	return b, nil
***REMOVED***

// SetWriteDeadline sets the write deadline for this bucket.
func (b *Bucket) SetWriteDeadline(t time.Time) error ***REMOVED***
	b.writeDeadline = t
	return nil
***REMOVED***

// SetReadDeadline sets the read deadline for this bucket
func (b *Bucket) SetReadDeadline(t time.Time) error ***REMOVED***
	b.readDeadline = t
	return nil
***REMOVED***

// OpenUploadStream creates a file ID new upload stream for a file given the filename.
func (b *Bucket) OpenUploadStream(filename string, opts ...*options.UploadOptions) (*UploadStream, error) ***REMOVED***
	return b.OpenUploadStreamWithID(primitive.NewObjectID(), filename, opts...)
***REMOVED***

// OpenUploadStreamWithID creates a new upload stream for a file given the file ID and filename.
func (b *Bucket) OpenUploadStreamWithID(fileID interface***REMOVED******REMOVED***, filename string, opts ...*options.UploadOptions) (*UploadStream, error) ***REMOVED***
	ctx, cancel := deadlineContext(b.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	if err := b.checkFirstWrite(ctx); err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	upload, err := b.parseUploadOptions(opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return newUploadStream(upload, fileID, filename, b.chunksColl, b.filesColl), nil
***REMOVED***

// UploadFromStream creates a fileID and uploads a file given a source stream.
//
// If this upload requires a custom write deadline to be set on the bucket, it cannot be done concurrently with other
// write operations operations on this bucket that also require a custom deadline.
func (b *Bucket) UploadFromStream(filename string, source io.Reader, opts ...*options.UploadOptions) (primitive.ObjectID, error) ***REMOVED***
	fileID := primitive.NewObjectID()
	err := b.UploadFromStreamWithID(fileID, filename, source, opts...)
	return fileID, err
***REMOVED***

// UploadFromStreamWithID uploads a file given a source stream.
//
// If this upload requires a custom write deadline to be set on the bucket, it cannot be done concurrently with other
// write operations operations on this bucket that also require a custom deadline.
func (b *Bucket) UploadFromStreamWithID(fileID interface***REMOVED******REMOVED***, filename string, source io.Reader, opts ...*options.UploadOptions) error ***REMOVED***
	us, err := b.OpenUploadStreamWithID(fileID, filename, opts...)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	err = us.SetWriteDeadline(b.writeDeadline)
	if err != nil ***REMOVED***
		_ = us.Close()
		return err
	***REMOVED***

	for ***REMOVED***
		n, err := source.Read(b.readBuf)
		if err != nil && err != io.EOF ***REMOVED***
			_ = us.Abort() // upload considered aborted if source stream returns an error
			return err
		***REMOVED***

		if n > 0 ***REMOVED***
			_, err := us.Write(b.readBuf[:n])
			if err != nil ***REMOVED***
				return err
			***REMOVED***
		***REMOVED***

		if n == 0 || err == io.EOF ***REMOVED***
			break
		***REMOVED***
	***REMOVED***

	return us.Close()
***REMOVED***

// OpenDownloadStream creates a stream from which the contents of the file can be read.
func (b *Bucket) OpenDownloadStream(fileID interface***REMOVED******REMOVED***) (*DownloadStream, error) ***REMOVED***
	return b.openDownloadStream(bson.D***REMOVED***
		***REMOVED***"_id", fileID***REMOVED***,
	***REMOVED***)
***REMOVED***

// DownloadToStream downloads the file with the specified fileID and writes it to the provided io.Writer.
// Returns the number of bytes written to the stream and an error, or nil if there was no error.
//
// If this download requires a custom read deadline to be set on the bucket, it cannot be done concurrently with other
// read operations operations on this bucket that also require a custom deadline.
func (b *Bucket) DownloadToStream(fileID interface***REMOVED******REMOVED***, stream io.Writer) (int64, error) ***REMOVED***
	ds, err := b.OpenDownloadStream(fileID)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return b.downloadToStream(ds, stream)
***REMOVED***

// OpenDownloadStreamByName opens a download stream for the file with the given filename.
func (b *Bucket) OpenDownloadStreamByName(filename string, opts ...*options.NameOptions) (*DownloadStream, error) ***REMOVED***
	var numSkip int32 = -1
	var sortOrder int32 = 1

	nameOpts := options.MergeNameOptions(opts...)
	if nameOpts.Revision != nil ***REMOVED***
		numSkip = *nameOpts.Revision
	***REMOVED***

	if numSkip < 0 ***REMOVED***
		sortOrder = -1
		numSkip = (-1 * numSkip) - 1
	***REMOVED***

	findOpts := options.Find().SetSkip(int64(numSkip)).SetSort(bson.D***REMOVED******REMOVED***"uploadDate", sortOrder***REMOVED******REMOVED***)

	return b.openDownloadStream(bson.D***REMOVED******REMOVED***"filename", filename***REMOVED******REMOVED***, findOpts)
***REMOVED***

// DownloadToStreamByName downloads the file with the given name to the given io.Writer.
//
// If this download requires a custom read deadline to be set on the bucket, it cannot be done concurrently with other
// read operations operations on this bucket that also require a custom deadline.
func (b *Bucket) DownloadToStreamByName(filename string, stream io.Writer, opts ...*options.NameOptions) (int64, error) ***REMOVED***
	ds, err := b.OpenDownloadStreamByName(filename, opts...)
	if err != nil ***REMOVED***
		return 0, err
	***REMOVED***

	return b.downloadToStream(ds, stream)
***REMOVED***

// Delete deletes all chunks and metadata associated with the file with the given file ID.
//
// If this operation requires a custom write deadline to be set on the bucket, it cannot be done concurrently with other
// write operations operations on this bucket that also require a custom deadline.
func (b *Bucket) Delete(fileID interface***REMOVED******REMOVED***) error ***REMOVED***
	// delete document in files collection and then chunks to minimize race conditions

	ctx, cancel := deadlineContext(b.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	res, err := b.filesColl.DeleteOne(ctx, bson.D***REMOVED******REMOVED***"_id", fileID***REMOVED******REMOVED***)
	if err == nil && res.DeletedCount == 0 ***REMOVED***
		err = ErrFileNotFound
	***REMOVED***
	if err != nil ***REMOVED***
		_ = b.deleteChunks(ctx, fileID) // can attempt to delete chunks even if no docs in files collection matched
		return err
	***REMOVED***

	return b.deleteChunks(ctx, fileID)
***REMOVED***

// Find returns the files collection documents that match the given filter.
//
// If this download requires a custom read deadline to be set on the bucket, it cannot be done concurrently with other
// read operations operations on this bucket that also require a custom deadline.
func (b *Bucket) Find(filter interface***REMOVED******REMOVED***, opts ...*options.GridFSFindOptions) (*mongo.Cursor, error) ***REMOVED***
	ctx, cancel := deadlineContext(b.readDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	gfsOpts := options.MergeGridFSFindOptions(opts...)
	find := options.Find()
	if gfsOpts.AllowDiskUse != nil ***REMOVED***
		find.SetAllowDiskUse(*gfsOpts.AllowDiskUse)
	***REMOVED***
	if gfsOpts.BatchSize != nil ***REMOVED***
		find.SetBatchSize(*gfsOpts.BatchSize)
	***REMOVED***
	if gfsOpts.Limit != nil ***REMOVED***
		find.SetLimit(int64(*gfsOpts.Limit))
	***REMOVED***
	if gfsOpts.MaxTime != nil ***REMOVED***
		find.SetMaxTime(*gfsOpts.MaxTime)
	***REMOVED***
	if gfsOpts.NoCursorTimeout != nil ***REMOVED***
		find.SetNoCursorTimeout(*gfsOpts.NoCursorTimeout)
	***REMOVED***
	if gfsOpts.Skip != nil ***REMOVED***
		find.SetSkip(int64(*gfsOpts.Skip))
	***REMOVED***
	if gfsOpts.Sort != nil ***REMOVED***
		find.SetSort(gfsOpts.Sort)
	***REMOVED***

	return b.filesColl.Find(ctx, filter, find)
***REMOVED***

// Rename renames the stored file with the specified file ID.
//
// If this operation requires a custom write deadline to be set on the bucket, it cannot be done concurrently with other
// write operations operations on this bucket that also require a custom deadline
func (b *Bucket) Rename(fileID interface***REMOVED******REMOVED***, newFilename string) error ***REMOVED***
	ctx, cancel := deadlineContext(b.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	res, err := b.filesColl.UpdateOne(ctx,
		bson.D***REMOVED******REMOVED***"_id", fileID***REMOVED******REMOVED***,
		bson.D***REMOVED******REMOVED***"$set", bson.D***REMOVED******REMOVED***"filename", newFilename***REMOVED******REMOVED******REMOVED******REMOVED***,
	)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if res.MatchedCount == 0 ***REMOVED***
		return ErrFileNotFound
	***REMOVED***

	return nil
***REMOVED***

// Drop drops the files and chunks collections associated with this bucket.
//
// If this operation requires a custom write deadline to be set on the bucket, it cannot be done concurrently with other
// write operations operations on this bucket that also require a custom deadline
func (b *Bucket) Drop() error ***REMOVED***
	ctx, cancel := deadlineContext(b.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	err := b.filesColl.Drop(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return b.chunksColl.Drop(ctx)
***REMOVED***

// GetFilesCollection returns a handle to the collection that stores the file documents for this bucket.
func (b *Bucket) GetFilesCollection() *mongo.Collection ***REMOVED***
	return b.filesColl
***REMOVED***

// GetChunksCollection returns a handle to the collection that stores the file chunks for this bucket.
func (b *Bucket) GetChunksCollection() *mongo.Collection ***REMOVED***
	return b.chunksColl
***REMOVED***

func (b *Bucket) openDownloadStream(filter interface***REMOVED******REMOVED***, opts ...*options.FindOptions) (*DownloadStream, error) ***REMOVED***
	ctx, cancel := deadlineContext(b.readDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	cursor, err := b.findFile(ctx, filter, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	// Unmarshal the data into a File instance, which can be passed to newDownloadStream. The _id value has to be
	// parsed out separately because "_id" will not match the File.ID field and we want to avoid exposing BSON tags
	// in the File type. After parsing it, use RawValue.Unmarshal to ensure File.ID is set to the appropriate value.
	var foundFile File
	if err = cursor.Decode(&foundFile); err != nil ***REMOVED***
		return nil, fmt.Errorf("error decoding files collection document: %v", err)
	***REMOVED***

	if foundFile.Length == 0 ***REMOVED***
		return newDownloadStream(nil, foundFile.ChunkSize, &foundFile), nil
	***REMOVED***

	// For a file with non-zero length, chunkSize must exist so we know what size to expect when downloading chunks.
	if _, err := cursor.Current.LookupErr("chunkSize"); err != nil ***REMOVED***
		return nil, ErrMissingChunkSize
	***REMOVED***

	chunksCursor, err := b.findChunks(ctx, foundFile.ID)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	// The chunk size can be overridden for individual files, so the expected chunk size should be the "chunkSize"
	// field from the files collection document, not the bucket's chunk size.
	return newDownloadStream(chunksCursor, foundFile.ChunkSize, &foundFile), nil
***REMOVED***

func deadlineContext(deadline time.Time) (context.Context, context.CancelFunc) ***REMOVED***
	if deadline.Equal(time.Time***REMOVED******REMOVED***) ***REMOVED***
		return context.Background(), nil
	***REMOVED***

	return context.WithDeadline(context.Background(), deadline)
***REMOVED***

func (b *Bucket) downloadToStream(ds *DownloadStream, stream io.Writer) (int64, error) ***REMOVED***
	err := ds.SetReadDeadline(b.readDeadline)
	if err != nil ***REMOVED***
		_ = ds.Close()
		return 0, err
	***REMOVED***

	copied, err := io.Copy(stream, ds)
	if err != nil ***REMOVED***
		_ = ds.Close()
		return 0, err
	***REMOVED***

	return copied, ds.Close()
***REMOVED***

func (b *Bucket) deleteChunks(ctx context.Context, fileID interface***REMOVED******REMOVED***) error ***REMOVED***
	_, err := b.chunksColl.DeleteMany(ctx, bson.D***REMOVED******REMOVED***"files_id", fileID***REMOVED******REMOVED***)
	return err
***REMOVED***

func (b *Bucket) findFile(ctx context.Context, filter interface***REMOVED******REMOVED***, opts ...*options.FindOptions) (*mongo.Cursor, error) ***REMOVED***
	cursor, err := b.filesColl.Find(ctx, filter, opts...)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	if !cursor.Next(ctx) ***REMOVED***
		_ = cursor.Close(ctx)
		return nil, ErrFileNotFound
	***REMOVED***

	return cursor, nil
***REMOVED***

func (b *Bucket) findChunks(ctx context.Context, fileID interface***REMOVED******REMOVED***) (*mongo.Cursor, error) ***REMOVED***
	chunksCursor, err := b.chunksColl.Find(ctx,
		bson.D***REMOVED******REMOVED***"files_id", fileID***REMOVED******REMOVED***,
		options.Find().SetSort(bson.D***REMOVED******REMOVED***"n", 1***REMOVED******REMOVED***)) // sort by chunk index
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	return chunksCursor, nil
***REMOVED***

// returns true if the 2 index documents are equal
func numericalIndexDocsEqual(expected, actual bsoncore.Document) (bool, error) ***REMOVED***
	if bytes.Equal(expected, actual) ***REMOVED***
		return true, nil
	***REMOVED***

	actualElems, err := actual.Elements()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***
	expectedElems, err := expected.Elements()
	if err != nil ***REMOVED***
		return false, err
	***REMOVED***

	if len(actualElems) != len(expectedElems) ***REMOVED***
		return false, nil
	***REMOVED***

	for idx, expectedElem := range expectedElems ***REMOVED***
		actualElem := actualElems[idx]
		if actualElem.Key() != expectedElem.Key() ***REMOVED***
			return false, nil
		***REMOVED***

		actualVal := actualElem.Value()
		expectedVal := expectedElem.Value()
		actualInt, actualOK := actualVal.AsInt64OK()
		expectedInt, expectedOK := expectedVal.AsInt64OK()

		//GridFS indexes always have numeric values
		if !actualOK || !expectedOK ***REMOVED***
			return false, nil
		***REMOVED***

		if actualInt != expectedInt ***REMOVED***
			return false, nil
		***REMOVED***
	***REMOVED***
	return true, nil
***REMOVED***

// Create an index if it doesn't already exist
func createNumericalIndexIfNotExists(ctx context.Context, iv mongo.IndexView, model mongo.IndexModel) error ***REMOVED***
	c, err := iv.List(ctx)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer func() ***REMOVED***
		_ = c.Close(ctx)
	***REMOVED***()

	modelKeysBytes, err := bson.Marshal(model.Keys)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	modelKeysDoc := bsoncore.Document(modelKeysBytes)

	for c.Next(ctx) ***REMOVED***
		keyElem, err := c.Current.LookupErr("key")
		if err != nil ***REMOVED***
			return err
		***REMOVED***

		keyElemDoc := keyElem.Document()

		found, err := numericalIndexDocsEqual(modelKeysDoc, bsoncore.Document(keyElemDoc))
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		if found ***REMOVED***
			return nil
		***REMOVED***
	***REMOVED***

	_, err = iv.CreateOne(ctx, model)
	return err
***REMOVED***

// create indexes on the files and chunks collection if needed
func (b *Bucket) createIndexes(ctx context.Context) error ***REMOVED***
	// must use primary read pref mode to check if files coll empty
	cloned, err := b.filesColl.Clone(options.Collection().SetReadPreference(readpref.Primary()))
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	docRes := cloned.FindOne(ctx, bson.D***REMOVED******REMOVED***, options.FindOne().SetProjection(bson.D***REMOVED******REMOVED***"_id", 1***REMOVED******REMOVED***))

	_, err = docRes.DecodeBytes()
	if err != mongo.ErrNoDocuments ***REMOVED***
		// nil, or error that occurred during the FindOne operation
		return err
	***REMOVED***

	filesIv := b.filesColl.Indexes()
	chunksIv := b.chunksColl.Indexes()

	filesModel := mongo.IndexModel***REMOVED***
		Keys: bson.D***REMOVED***
			***REMOVED***"filename", int32(1)***REMOVED***,
			***REMOVED***"uploadDate", int32(1)***REMOVED***,
		***REMOVED***,
	***REMOVED***

	chunksModel := mongo.IndexModel***REMOVED***
		Keys: bson.D***REMOVED***
			***REMOVED***"files_id", int32(1)***REMOVED***,
			***REMOVED***"n", int32(1)***REMOVED***,
		***REMOVED***,
		Options: options.Index().SetUnique(true),
	***REMOVED***

	if err = createNumericalIndexIfNotExists(ctx, filesIv, filesModel); err != nil ***REMOVED***
		return err
	***REMOVED***
	if err = createNumericalIndexIfNotExists(ctx, chunksIv, chunksModel); err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***

func (b *Bucket) checkFirstWrite(ctx context.Context) error ***REMOVED***
	if !b.firstWriteDone ***REMOVED***
		// before the first write operation, must determine if files collection is empty
		// if so, create indexes if they do not already exist

		if err := b.createIndexes(ctx); err != nil ***REMOVED***
			return err
		***REMOVED***
		b.firstWriteDone = true
	***REMOVED***

	return nil
***REMOVED***

func (b *Bucket) parseUploadOptions(opts ...*options.UploadOptions) (*Upload, error) ***REMOVED***
	upload := &Upload***REMOVED***
		chunkSize: b.chunkSize, // upload chunk size defaults to bucket's value
	***REMOVED***

	uo := options.MergeUploadOptions(opts...)
	if uo.ChunkSizeBytes != nil ***REMOVED***
		upload.chunkSize = *uo.ChunkSizeBytes
	***REMOVED***
	if uo.Registry == nil ***REMOVED***
		uo.Registry = bson.DefaultRegistry
	***REMOVED***
	if uo.Metadata != nil ***REMOVED***
		raw, err := bson.MarshalWithRegistry(uo.Registry, uo.Metadata)
		if err != nil ***REMOVED***
			return nil, err
		***REMOVED***
		var doc bson.D
		unMarErr := bson.UnmarshalWithRegistry(uo.Registry, raw, &doc)
		if unMarErr != nil ***REMOVED***
			return nil, unMarErr
		***REMOVED***
		upload.metadata = doc
	***REMOVED***

	return upload, nil
***REMOVED***
