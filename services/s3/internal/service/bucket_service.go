package service

import (
	"cloudlab/s3/internal/storage"
	"io"
)

func CreateBucket(bucketName string) error {
	return storage.CreateBucket(bucketName)
}

func UploadObject(bucketName, objectKey string, data io.Reader) error {
	return storage.WriteObject(bucketName, objectKey, data)
}
