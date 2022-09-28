// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package gridfs

import (
	"context"
	"errors"
	"io"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrWrongIndex is used when the chunk retrieved from the server does not have the expected index.
var ErrWrongIndex = errors.New("chunk index does not match expected index")

// ErrWrongSize is used when the chunk retrieved from the server does not have the expected size.
var ErrWrongSize = errors.New("chunk size does not match expected size")

var errNoMoreChunks = errors.New("no more chunks remaining")

// DownloadStream is a io.Reader that can be used to download a file from a GridFS bucket.
type DownloadStream struct ***REMOVED***
	numChunks     int32
	chunkSize     int32
	cursor        *mongo.Cursor
	done          bool
	closed        bool
	buffer        []byte // store up to 1 chunk if the user provided buffer isn't big enough
	bufferStart   int
	bufferEnd     int
	expectedChunk int32 // index of next expected chunk
	readDeadline  time.Time
	fileLen       int64

	// The pointer returned by GetFile. This should not be used in the actual DownloadStream code outside of the
	// newDownloadStream constructor because the values can be mutated by the user after calling GetFile. Instead,
	// any values needed in the code should be stored separately and copied over in the constructor.
	file *File
***REMOVED***

// File represents a file stored in GridFS. This type can be used to access file information when downloading using the
// DownloadStream.GetFile method.
type File struct ***REMOVED***
	// ID is the file's ID. This will match the file ID specified when uploading the file. If an upload helper that
	// does not require a file ID was used, this field will be a primitive.ObjectID.
	ID interface***REMOVED******REMOVED***

	// Length is the length of this file in bytes.
	Length int64

	// ChunkSize is the maximum number of bytes for each chunk in this file.
	ChunkSize int32

	// UploadDate is the time this file was added to GridFS in UTC. This field is set by the driver and is not configurable.
	// The Metadata field can be used to store a custom date.
	UploadDate time.Time

	// Name is the name of this file.
	Name string

	// Metadata is additional data that was specified when creating this file. This field can be unmarshalled into a
	// custom type using the bson.Unmarshal family of functions.
	Metadata bson.Raw
***REMOVED***

var _ bson.Unmarshaler = (*File)(nil)

// unmarshalFile is a temporary type used to unmarshal documents from the files collection and can be transformed into
// a File instance. This type exists to avoid adding BSON struct tags to the exported File type.
type unmarshalFile struct ***REMOVED***
	ID         interface***REMOVED******REMOVED*** `bson:"_id"`
	Length     int64       `bson:"length"`
	ChunkSize  int32       `bson:"chunkSize"`
	UploadDate time.Time   `bson:"uploadDate"`
	Name       string      `bson:"filename"`
	Metadata   bson.Raw    `bson:"metadata"`
***REMOVED***

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (f *File) UnmarshalBSON(data []byte) error ***REMOVED***
	var temp unmarshalFile
	if err := bson.Unmarshal(data, &temp); err != nil ***REMOVED***
		return err
	***REMOVED***

	f.ID = temp.ID
	f.Length = temp.Length
	f.ChunkSize = temp.ChunkSize
	f.UploadDate = temp.UploadDate
	f.Name = temp.Name
	f.Metadata = temp.Metadata
	return nil
***REMOVED***

func newDownloadStream(cursor *mongo.Cursor, chunkSize int32, file *File) *DownloadStream ***REMOVED***
	numChunks := int32(math.Ceil(float64(file.Length) / float64(chunkSize)))

	return &DownloadStream***REMOVED***
		numChunks: numChunks,
		chunkSize: chunkSize,
		cursor:    cursor,
		buffer:    make([]byte, chunkSize),
		done:      cursor == nil,
		fileLen:   file.Length,
		file:      file,
	***REMOVED***
***REMOVED***

// Close closes this download stream.
func (ds *DownloadStream) Close() error ***REMOVED***
	if ds.closed ***REMOVED***
		return ErrStreamClosed
	***REMOVED***

	ds.closed = true
	if ds.cursor != nil ***REMOVED***
		return ds.cursor.Close(context.Background())
	***REMOVED***
	return nil
***REMOVED***

// SetReadDeadline sets the read deadline for this download stream.
func (ds *DownloadStream) SetReadDeadline(t time.Time) error ***REMOVED***
	if ds.closed ***REMOVED***
		return ErrStreamClosed
	***REMOVED***

	ds.readDeadline = t
	return nil
***REMOVED***

// Read reads the file from the server and writes it to a destination byte slice.
func (ds *DownloadStream) Read(p []byte) (int, error) ***REMOVED***
	if ds.closed ***REMOVED***
		return 0, ErrStreamClosed
	***REMOVED***

	if ds.done ***REMOVED***
		return 0, io.EOF
	***REMOVED***

	ctx, cancel := deadlineContext(ds.readDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	bytesCopied := 0
	var err error
	for bytesCopied < len(p) ***REMOVED***
		if ds.bufferStart >= ds.bufferEnd ***REMOVED***
			// Buffer is empty and can load in data from new chunk.
			err = ds.fillBuffer(ctx)
			if err != nil ***REMOVED***
				if err == errNoMoreChunks ***REMOVED***
					if bytesCopied == 0 ***REMOVED***
						ds.done = true
						return 0, io.EOF
					***REMOVED***
					return bytesCopied, nil
				***REMOVED***
				return bytesCopied, err
			***REMOVED***
		***REMOVED***

		copied := copy(p[bytesCopied:], ds.buffer[ds.bufferStart:ds.bufferEnd])

		bytesCopied += copied
		ds.bufferStart += copied
	***REMOVED***

	return len(p), nil
***REMOVED***

// Skip skips a given number of bytes in the file.
func (ds *DownloadStream) Skip(skip int64) (int64, error) ***REMOVED***
	if ds.closed ***REMOVED***
		return 0, ErrStreamClosed
	***REMOVED***

	if ds.done ***REMOVED***
		return 0, nil
	***REMOVED***

	ctx, cancel := deadlineContext(ds.readDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	var skipped int64
	var err error

	for skipped < skip ***REMOVED***
		if ds.bufferStart >= ds.bufferEnd ***REMOVED***
			// Buffer is empty and can load in data from new chunk.
			err = ds.fillBuffer(ctx)
			if err != nil ***REMOVED***
				if err == errNoMoreChunks ***REMOVED***
					return skipped, nil
				***REMOVED***
				return skipped, err
			***REMOVED***
		***REMOVED***

		toSkip := skip - skipped
		// Cap the amount to skip to the remaining bytes in the buffer to be consumed.
		bufferRemaining := ds.bufferEnd - ds.bufferStart
		if toSkip > int64(bufferRemaining) ***REMOVED***
			toSkip = int64(bufferRemaining)
		***REMOVED***

		skipped += toSkip
		ds.bufferStart += int(toSkip)
	***REMOVED***

	return skip, nil
***REMOVED***

// GetFile returns a File object representing the file being downloaded.
func (ds *DownloadStream) GetFile() *File ***REMOVED***
	return ds.file
***REMOVED***

func (ds *DownloadStream) fillBuffer(ctx context.Context) error ***REMOVED***
	if !ds.cursor.Next(ctx) ***REMOVED***
		ds.done = true
		// Check for cursor error, otherwise there are no more chunks.
		if ds.cursor.Err() != nil ***REMOVED***
			_ = ds.cursor.Close(ctx)
			return ds.cursor.Err()
		***REMOVED***
		return errNoMoreChunks
	***REMOVED***

	chunkIndex, err := ds.cursor.Current.LookupErr("n")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	var chunkIndexInt32 int32
	if chunkIndexInt64, ok := chunkIndex.Int64OK(); ok ***REMOVED***
		chunkIndexInt32 = int32(chunkIndexInt64)
	***REMOVED*** else ***REMOVED***
		chunkIndexInt32 = chunkIndex.Int32()
	***REMOVED***

	if chunkIndexInt32 != ds.expectedChunk ***REMOVED***
		return ErrWrongIndex
	***REMOVED***

	ds.expectedChunk++
	data, err := ds.cursor.Current.LookupErr("data")
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	_, dataBytes := data.Binary()
	copied := copy(ds.buffer, dataBytes)

	bytesLen := int32(len(dataBytes))
	if ds.expectedChunk == ds.numChunks ***REMOVED***
		// final chunk can be fewer than ds.chunkSize bytes
		bytesDownloaded := int64(ds.chunkSize) * (int64(ds.expectedChunk) - int64(1))
		bytesRemaining := ds.fileLen - bytesDownloaded

		if int64(bytesLen) != bytesRemaining ***REMOVED***
			return ErrWrongSize
		***REMOVED***
	***REMOVED*** else if bytesLen != ds.chunkSize ***REMOVED***
		// all intermediate chunks must have size ds.chunkSize
		return ErrWrongSize
	***REMOVED***

	ds.bufferStart = 0
	ds.bufferEnd = copied

	return nil
***REMOVED***
