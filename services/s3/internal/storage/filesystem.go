package storage

import (
	"io"
	"os"
	"path/filepath"
)

const DataDir = "data"

func CreateBucket(bucketName string) error {
	path := filepath.Join(DataDir, bucketName)
	return os.MkdirAll(path, os.ModePerm)
}

func WriteObject(bucketName string, objectKey string, data io.Reader) error {
	objectPath := filepath.Join(DataDir, bucketName, objectKey)

	file, err := os.Create(objectPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}
