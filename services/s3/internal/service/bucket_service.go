package service

import (
	"cloudlab/s3/internal/models"
	"cloudlab/s3/internal/storage"
	"io"
	"time"
)

func CreateBucket(bucketName string) error {
	return storage.CreateBucket(bucketName)
}

func UploadObject(bucketName, objectKey string, data io.Reader, contentType string, size int64) error {

	err := storage.WriteObject(bucketName, objectKey, data)
	if err != nil {
		return err
	}

	metadata := models.ObjectMetadata{
		Bucket:      bucketName,
		Key:         objectKey,
		Size:        size,
		ContentType: contentType,
		CreatedAt:   time.Now(),
	}
  
	return storage.WriteMetadata(metadata, bucketName, objectKey)
}

func GetObject(bucketName, objectKey string) (io.ReadCloser, error) {
	return storage.ReadObject(bucketName, objectKey)
}
