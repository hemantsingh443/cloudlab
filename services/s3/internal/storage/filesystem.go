package storage

import (
	"os"
	"path/filepath"
)

const DataDir = "data"

func CreateBucket(bucketName string) error {
	path := filepath.Join(DataDir, bucketName)
	return os.MkdirAll(path, os.ModePerm)
}
