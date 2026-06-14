package http

import (
	"net"
	"net/http"
	"strings"
)

// resolveBucketAndKey extracts the S3 bucket name and object key from the incoming request.
// It supports both Virtual-Host style routing and Path-style routing.
// isVirtualHost is returned as true if a virtual-host style subdomain was matched.
func resolveBucketAndKey(r *http.Request) (bucket string, key string, isVirtualHost bool) {
	host := r.Host
	if h, _, err := net.SplitHostPort(host); err == nil {
		host = h
	}

	// If it is an IP address or exactly "localhost", it must be path-style
	if net.ParseIP(host) != nil || host == "localhost" {
		return parsePathStyle(r.URL.Path)
	}

	// Check for AWS S3 styles (e.g. s3.amazonaws.com or bucket.s3.amazonaws.com)
	if strings.HasSuffix(host, ".amazonaws.com") {
		trimmed := strings.TrimSuffix(host, ".amazonaws.com")
		parts := strings.Split(trimmed, ".")
		if len(parts) > 0 && parts[0] == "s3" {
			// Path-style: e.g. s3.amazonaws.com or s3.us-east-1.amazonaws.com
			return parsePathStyle(r.URL.Path)
		}
		// Virtual-host style: parts[0] is the bucket name
		// For bucket.s3.us-east-1.amazonaws.com, parts is [bucket, s3, us-east-1]
		if len(parts) > 0 {
			bucket = parts[0]
			key = strings.TrimPrefix(r.URL.Path, "/")
			return bucket, key, true
		}
	}

	// Check for local virtual-host subdomains (e.g., bucket.localhost)
	if strings.HasSuffix(host, ".localhost") {
		bucket = strings.TrimSuffix(host, ".localhost")
		key = strings.TrimPrefix(r.URL.Path, "/")
		return bucket, key, true
	}
	if strings.HasSuffix(host, ".127.0.0.1") {
		bucket = strings.TrimSuffix(host, ".127.0.0.1")
		key = strings.TrimPrefix(r.URL.Path, "/")
		return bucket, key, true
	}

	// General fallback for custom domains:
	// If the host has at least one dot and does not start with "s3.",
	// assume the first segment is the bucket name (virtual-host style).
	parts := strings.Split(host, ".")
	if len(parts) > 1 && parts[0] != "s3" {
		bucket = parts[0]
		key = strings.TrimPrefix(r.URL.Path, "/")
		return bucket, key, true
	}

	return parsePathStyle(r.URL.Path)
}

func parsePathStyle(urlPath string) (bucket string, key string, isVirtualHost bool) {
	path := strings.TrimPrefix(urlPath, "/")
	if path == "" {
		return "", "", false
	}
	parts := strings.SplitN(path, "/", 2)
	bucket = parts[0]
	if len(parts) > 1 {
		key = parts[1]
	}
	return bucket, key, false
}
