package storage

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

const (
	DataDir     = "data"
	BlobDir     = "data/blobs"
	MetadataDir = "data/metadata"
)

func CreateBucket(bucketName string) error {
	blobpath := filepath.Join(BlobDir, bucketName)
	metadataPath := filepath.Join(MetadataDir, bucketName)

	err := os.MkdirAll(blobpath, os.ModePerm)
	if err != nil {
		return err
	}
	return os.MkdirAll(metadataPath, os.ModePerm)
}

func WriteObject(bucketName string, objectKey string, data io.Reader) error {
	objectPath := filepath.Join(BlobDir, bucketName, objectKey)

	file, err := os.Create(objectPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func ReadObject(bucketName, objectKey string) (io.ReadCloser, error) {
	objectPath := filepath.Join(BlobDir, bucketName, objectKey)

	file, err := os.Open(objectPath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func WriteMetadata(metadata interface{}, bucketName, objectKey string) error {
	path := filepath.Join(
		MetadataDir,
		bucketName,
		objectKey+".json",
	)

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(metadata)
}
