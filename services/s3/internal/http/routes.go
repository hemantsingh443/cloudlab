package http

import (
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/bucket/", CreateBucketHandler)
	mux.HandleFunc("/object/", ObjectHandler)

	return mux
}
