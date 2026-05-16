package http

import (
	"cloudlab/s3/internal/service"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bucketName := strings.TrimPrefix(r.URL.Path, "/bucket/")

	if bucketName == "" {
		http.Error(w, "bucket name required", http.StatusBadRequest)
		return
	}

	err := service.CreateBucket(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := map[string]string{
		"message": "bucket created",
		"bucket":  bucketName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func UploadObjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/object/")
	parts := strings.SplitN(path, "/", 2)

	if len(parts) != 2 {
		http.Error(w, "invalid object path", http.StatusBadRequest)
		return
	}

	bucketName := parts[0]
	objectKey := parts[1]

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

func GetObjectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/object/")
	parts := strings.SplitN(path, "/", 2)

	if len(parts) != 2 {
		http.Error(w, "invalid object path", http.StatusBadRequest)
		return
	}

	bucketName := parts[0]
	objectKey := parts[1]

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

func ObjectHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		UploadObjectHandler(w, r)

	case http.MethodGet:
		GetObjectHandler(w, r) 

	case http.MethodDelete:  
		DeleteObjectHandler(w,r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func ListBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := strings.TrimPrefix(r.URL.Path, "/bucket/")

	if bucketName == "" {
		http.Error(w, "bucket name required", http.StatusBadRequest)
		return
	}

	objects, err := service.ListBucketObjects(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(objects)
}

func BucketHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		CreateBucketHandler(w, r)

	case http.MethodGet:
		ListBucketHandler(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
} 

func DeleteObjectHandler(w http.ResponseWriter, r *http.Request) { 
	path := strings.TrimPrefix(r.URL.Path, "/object/") 
	parts := strings.SplitN(path, "/", 2) 
 
	if len(parts) != 2 { 
		http.Error(w, "invalid object path", http.StatusBadRequest) 
		return
	} 

	bucketName := parts[0] 
	objectKey := parts[1] 

	err := service.DeleteObject(bucketName, objectKey) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
		return 
	} 

	w.WriteHeader(http.StatusNoContent)
}
