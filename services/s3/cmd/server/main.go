package main

import (
	"log"
	"net/http"

	httpHandlers "cloudlab/s3/internal/http"
)

func main() {
	router := httpHandlers.RegisterRoutes()

	log.Println("S3 service starting at :8080")

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}

}
