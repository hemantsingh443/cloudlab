package http

import (
	"net/http"
)

func RegisterRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", HealthHandler)
	mux.HandleFunc("/bucket/", BucketHandler)
	mux.HandleFunc("/object/", ObjectHandler)

	return mux
}
