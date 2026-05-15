package service

import (
	"cloudlab/s3/internal/models"
	"cloudlab/s3/internal/storage"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"time"
)

func CreateBucket(bucketName string) error {
	return storage.CreateBucket(bucketName)
}

func UploadObject(bucketName, objectKey string, data io.Reader, contentType string, size int64) error {

	hasher := sha256.New()
	teeReader := io.TeeReader(data, hasher)

	err := storage.WriteObject(bucketName, objectKey, teeReader)
	if err != nil {
		return err
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))

	metadata := models.ObjectMetadata{
		Bucket:      bucketName,
		Key:         objectKey,
		Size:        size,
		ContentType: contentType,
		Checksum:    checksum,
		CreatedAt:   time.Now(),
	}

	return storage.WriteMetadata(metadata, bucketName, objectKey)
}

func GetObject(bucketName, objectKey string) (io.ReadCloser, error) {
	return storage.ReadObject(bucketName, objectKey)
}

func ListBucketObjects(bucketName string) ([]models.ObjectMetadata, error) {
	return storage.ReadbucketMetadata(bucketName)
}
