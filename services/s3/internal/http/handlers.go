package http

import (
	"cloudlab/s3/internal/service"
	"encoding/json"
	"io"
	"net/http"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func S3Handler(w http.ResponseWriter, r *http.Request) {
	bucket, key, _ := resolveBucketAndKey(r)

	if bucket == "" {
		http.Error(w, "bucket name required", http.StatusBadRequest)
		return
	}

	if key == "" {
		// Bucket operations
		switch r.Method {
		case http.MethodPut:
			CreateBucketHandler(w, r, bucket)
		case http.MethodGet:
			ListBucketHandler(w, r, bucket)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Object operations
	switch r.Method {
	case http.MethodPut:
		UploadObjectHandler(w, r, bucket, key)
	case http.MethodGet:
		GetObjectHandler(w, r, bucket, key)
	case http.MethodDelete:
		DeleteObjectHandler(w, r, bucket, key)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func CreateBucketHandler(w http.ResponseWriter, r *http.Request, bucketName string) {
	err := service.CreateBucket(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "bucket created",
		"bucket":  bucketName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UploadObjectHandler(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	err := service.UploadObject(bucketName, objectKey, r.Body, contentType, r.ContentLength)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "object uploaded",
		"bucket":  bucketName,
		"object":  objectKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetObjectHandler(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	file, err := service.GetObject(bucketName, objectKey)
	if err != nil {
		http.Error(w, "object not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ListBucketHandler(w http.ResponseWriter, r *http.Request, bucketName string) {
	objects, err := service.ListBucketObjects(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objects)
}

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request, bucketName, objectKey string) {
	err := service.DeleteObject(bucketName, objectKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
