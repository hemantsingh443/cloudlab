package service

import "cloudlab/s3/internal/storage"

func CreateBucket(bucketName string) error {
	return storage.CreateBucket(bucketName)
}
