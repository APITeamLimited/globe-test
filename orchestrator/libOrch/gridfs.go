package libOrch

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SetInBucket(bucket *gridfs.Bucket, filename string, data []byte, contentType string, globeTestLogsReceipt primitive.ObjectID) error {
	bucketOptions := options.GridFSUpload()
	bucketOptions.Metadata = map[string]interface{}{
		"filename":    filename,
		"contentType": contentType,
	}

	uploadStream, err := bucket.OpenUploadStreamWithID(globeTestLogsReceipt, filename, bucketOptions)
	if err != nil {
		return err
	}
	defer uploadStream.Close()
	_, err = uploadStream.Write(data)
	if err != nil {
		return err
	}

	return nil
}
