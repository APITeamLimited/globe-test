// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package gridfs

import (
	"errors"

	"context"
	"time"

	"math"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UploadBufferSize is the size in bytes of one stream batch. Chunks will be written to the db after the sum of chunk
// lengths is equal to the batch size.
const UploadBufferSize = 16 * 1024 * 1024 // 16 MiB

// ErrStreamClosed is an error returned if an operation is attempted on a closed/aborted stream.
var ErrStreamClosed = errors.New("stream is closed or aborted")

// UploadStream is used to upload a file in chunks. This type implements the io.Writer interface and a file can be
// uploaded using the Write method. After an upload is complete, the Close method must be called to write file
// metadata.
type UploadStream struct ***REMOVED***
	*Upload // chunk size and metadata
	FileID  interface***REMOVED******REMOVED***

	chunkIndex    int
	chunksColl    *mongo.Collection // collection to store file chunks
	filename      string
	filesColl     *mongo.Collection // collection to store file metadata
	closed        bool
	buffer        []byte
	bufferIndex   int
	fileLen       int64
	writeDeadline time.Time
***REMOVED***

// NewUploadStream creates a new upload stream.
func newUploadStream(upload *Upload, fileID interface***REMOVED******REMOVED***, filename string, chunks, files *mongo.Collection) *UploadStream ***REMOVED***
	return &UploadStream***REMOVED***
		Upload: upload,
		FileID: fileID,

		chunksColl: chunks,
		filename:   filename,
		filesColl:  files,
		buffer:     make([]byte, UploadBufferSize),
	***REMOVED***
***REMOVED***

// Close writes file metadata to the files collection and cleans up any resources associated with the UploadStream.
func (us *UploadStream) Close() error ***REMOVED***
	if us.closed ***REMOVED***
		return ErrStreamClosed
	***REMOVED***

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	if us.bufferIndex != 0 ***REMOVED***
		if err := us.uploadChunks(ctx, true); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***

	if err := us.createFilesCollDoc(ctx); err != nil ***REMOVED***
		return err
	***REMOVED***

	us.closed = true
	return nil
***REMOVED***

// SetWriteDeadline sets the write deadline for this stream.
func (us *UploadStream) SetWriteDeadline(t time.Time) error ***REMOVED***
	if us.closed ***REMOVED***
		return ErrStreamClosed
	***REMOVED***

	us.writeDeadline = t
	return nil
***REMOVED***

// Write transfers the contents of a byte slice into this upload stream. If the stream's underlying buffer fills up,
// the buffer will be uploaded as chunks to the server. Implements the io.Writer interface.
func (us *UploadStream) Write(p []byte) (int, error) ***REMOVED***
	if us.closed ***REMOVED***
		return 0, ErrStreamClosed
	***REMOVED***

	var ctx context.Context

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	origLen := len(p)
	for ***REMOVED***
		if len(p) == 0 ***REMOVED***
			break
		***REMOVED***

		n := copy(us.buffer[us.bufferIndex:], p) // copy as much as possible
		p = p[n:]
		us.bufferIndex += n

		if us.bufferIndex == UploadBufferSize ***REMOVED***
			err := us.uploadChunks(ctx, false)
			if err != nil ***REMOVED***
				return 0, err
			***REMOVED***
		***REMOVED***
	***REMOVED***
	return origLen, nil
***REMOVED***

// Abort closes the stream and deletes all file chunks that have already been written.
func (us *UploadStream) Abort() error ***REMOVED***
	if us.closed ***REMOVED***
		return ErrStreamClosed
	***REMOVED***

	ctx, cancel := deadlineContext(us.writeDeadline)
	if cancel != nil ***REMOVED***
		defer cancel()
	***REMOVED***

	_, err := us.chunksColl.DeleteMany(ctx, bson.D***REMOVED******REMOVED***"files_id", us.FileID***REMOVED******REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	us.closed = true
	return nil
***REMOVED***

// uploadChunks uploads the current buffer as a series of chunks to the bucket
// if uploadPartial is true, any data at the end of the buffer that is smaller than a chunk will be uploaded as a partial
// chunk. if it is false, the data will be moved to the front of the buffer.
// uploadChunks sets us.bufferIndex to the next available index in the buffer after uploading
func (us *UploadStream) uploadChunks(ctx context.Context, uploadPartial bool) error ***REMOVED***
	chunks := float64(us.bufferIndex) / float64(us.chunkSize)
	numChunks := int(math.Ceil(chunks))
	if !uploadPartial ***REMOVED***
		numChunks = int(math.Floor(chunks))
	***REMOVED***

	docs := make([]interface***REMOVED******REMOVED***, numChunks)

	begChunkIndex := us.chunkIndex
	for i := 0; i < us.bufferIndex; i += int(us.chunkSize) ***REMOVED***
		endIndex := i + int(us.chunkSize)
		if us.bufferIndex-i < int(us.chunkSize) ***REMOVED***
			// partial chunk
			if !uploadPartial ***REMOVED***
				break
			***REMOVED***
			endIndex = us.bufferIndex
		***REMOVED***
		chunkData := us.buffer[i:endIndex]
		docs[us.chunkIndex-begChunkIndex] = bson.D***REMOVED***
			***REMOVED***"_id", primitive.NewObjectID()***REMOVED***,
			***REMOVED***"files_id", us.FileID***REMOVED***,
			***REMOVED***"n", int32(us.chunkIndex)***REMOVED***,
			***REMOVED***"data", primitive.Binary***REMOVED***Subtype: 0x00, Data: chunkData***REMOVED******REMOVED***,
		***REMOVED***
		us.chunkIndex++
		us.fileLen += int64(len(chunkData))
	***REMOVED***

	_, err := us.chunksColl.InsertMany(ctx, docs)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	// copy any remaining bytes to beginning of buffer and set buffer index
	bytesUploaded := numChunks * int(us.chunkSize)
	if bytesUploaded != UploadBufferSize && !uploadPartial ***REMOVED***
		copy(us.buffer[0:], us.buffer[bytesUploaded:us.bufferIndex])
	***REMOVED***
	us.bufferIndex = UploadBufferSize - bytesUploaded
	return nil
***REMOVED***

func (us *UploadStream) createFilesCollDoc(ctx context.Context) error ***REMOVED***
	doc := bson.D***REMOVED***
		***REMOVED***"_id", us.FileID***REMOVED***,
		***REMOVED***"length", us.fileLen***REMOVED***,
		***REMOVED***"chunkSize", us.chunkSize***REMOVED***,
		***REMOVED***"uploadDate", primitive.DateTime(time.Now().UnixNano() / int64(time.Millisecond))***REMOVED***,
		***REMOVED***"filename", us.filename***REMOVED***,
	***REMOVED***

	if us.metadata != nil ***REMOVED***
		doc = append(doc, bson.E***REMOVED***"metadata", us.metadata***REMOVED***)
	***REMOVED***

	_, err := us.filesColl.InsertOne(ctx, doc)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
