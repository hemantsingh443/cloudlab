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

func ReadObject(bucketName, objectKey string) (io.ReadCloser, error) {
	objectPath := filepath.Join(DataDir, bucketName, objectKey)

	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}
