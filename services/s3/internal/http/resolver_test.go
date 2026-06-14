package http

import (
	"net/http"
	"net/url"
	"testing"
)

func TestResolveBucketAndKey(t *testing.T) {
	tests := []struct {
		name                 string
		host                 string
		path                 string
		expectedBucket       string
		expectedKey          string
		expectedVirtualHost  bool
	}{
		{
			name:                "Path-style localhost with port",
			host:                "localhost:8080",
			path:                "/mybucket/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: false,
		},
		{
			name:                "Path-style IP address with port",
			host:                "127.0.0.1:8080",
			path:                "/mybucket/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: false,
		},
		{
			name:                "Virtual-host style localhost",
			host:                "mybucket.localhost:8080",
			path:                "/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: true,
		},
		{
			name:                "Virtual-host style localhost with nested key",
			host:                "mybucket.localhost:8080",
			path:                "/images/summer/pic.jpg",
			expectedBucket:      "mybucket",
			expectedKey:         "images/summer/pic.jpg",
			expectedVirtualHost: true,
		},
		{
			name:                "Virtual-host style AWS S3 standard",
			host:                "mybucket.s3.amazonaws.com",
			path:                "/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: true,
		},
		{
			name:                "Virtual-host style AWS S3 regional",
			host:                "mybucket.s3.us-west-2.amazonaws.com",
			path:                "/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: true,
		},
		{
			name:                "Path-style AWS S3 standard",
			host:                "s3.amazonaws.com",
			path:                "/mybucket/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: false,
		},
		{
			name:                "Path-style AWS S3 regional",
			host:                "s3.us-west-2.amazonaws.com",
			path:                "/mybucket/myfile.txt",
			expectedBucket:      "mybucket",
			expectedKey:         "myfile.txt",
			expectedVirtualHost: false,
		},
		{
			name:                "Path-style nested key",
			host:                "localhost:8080",
			path:                "/mybucket/images/summer/pic.jpg",
			expectedBucket:      "mybucket",
			expectedKey:         "images/summer/pic.jpg",
			expectedVirtualHost: false,
		},
		{
			name:                "Empty path/bucket operation localhost",
			host:                "localhost:8080",
			path:                "/mybucket",
			expectedBucket:      "mybucket",
			expectedKey:         "",
			expectedVirtualHost: false,
		},
		{
			name:                "Empty path/bucket operation virtual-host",
			host:                "mybucket.localhost:8080",
			path:                "/",
			expectedBucket:      "mybucket",
			expectedKey:         "",
			expectedVirtualHost: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				Host: tt.host,
				URL: &url.URL{
					Path: tt.path,
				},
			}

			bucket, key, isVH := resolveBucketAndKey(req)
			if bucket != tt.expectedBucket {
				t.Errorf("expected bucket %q, got %q", tt.expectedBucket, bucket)
			}
			if key != tt.expectedKey {
				t.Errorf("expected key %q, got %q", tt.expectedKey, key)
			}
			if isVH != tt.expectedVirtualHost {
				t.Errorf("expected isVirtualHost %v, got %v", tt.expectedVirtualHost, isVH)
			}
		})
	}
}
