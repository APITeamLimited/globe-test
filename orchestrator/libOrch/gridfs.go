package libOrch

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetInBucket(bucket *gridfs.Bucket, filename string, data []byte, contentType string, globeTestLogsId primitive.ObjectID) error ***REMOVED***
	bucketOptions := options.GridFSUpload()
	bucketOptions.Metadata = map[string]interface***REMOVED******REMOVED******REMOVED***
		"filename":    filename,
		"contentType": contentType,
	***REMOVED***

	uploadStream, err := bucket.OpenUploadStreamWithID(globeTestLogsId, filename, bucketOptions)
	if err != nil ***REMOVED***
		return err
	***REMOVED***
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	return nil
***REMOVED***
