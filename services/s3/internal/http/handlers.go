package http

import (
	"cloudlab/s3/internal/service"
	"encoding/json"
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
