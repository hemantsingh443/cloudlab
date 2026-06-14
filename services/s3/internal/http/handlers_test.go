package http

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cloudlab/s3/internal/models"
)

func TestS3APIIntegration(t *testing.T) {
	router := RegisterRoutes()

	// Unique buckets for test run
	pathBucket := "test-bucket-path-integration"
	vhBucket := "test-bucket-vh-integration"

	// Cleanup at the end
	t.Cleanup(func() {
		os.RemoveAll(filepath.Join("data", "blobs", pathBucket))
		os.RemoveAll(filepath.Join("data", "metadata", pathBucket))
		os.RemoveAll(filepath.Join("data", "blobs", vhBucket))
		os.RemoveAll(filepath.Join("data", "metadata", vhBucket))
	})

	// Test 1: Health Check
	t.Run("Health Check", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		var resp map[string]string
		if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}
		if resp["status"] != "ok" {
			t.Errorf("expected status ok, got %q", resp["status"])
		}
	})

	// Test 2: Create Bucket (Path-Style)
	t.Run("Create Bucket Path-Style", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/"+pathBucket, nil)
		req.Host = "localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 3: Upload Object (Path-Style)
	t.Run("Upload Object Path-Style", func(t *testing.T) {
		content := "hello path style"
		req := httptest.NewRequest(http.MethodPut, "/"+pathBucket+"/hello.txt", strings.NewReader(content))
		req.Host = "localhost:8080"
		req.Header.Set("Content-Type", "text/plain")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 4: Get Object (Path-Style)
	t.Run("Get Object Path-Style", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/"+pathBucket+"/hello.txt", nil)
		req.Host = "localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		bodyBytes, _ := io.ReadAll(w.Body)
		if string(bodyBytes) != "hello path style" {
			t.Errorf("expected body 'hello path style', got %q", string(bodyBytes))
		}
	})

	// Test 5: Create Bucket (Virtual-Host Style)
	t.Run("Create Bucket Virtual-Host Style", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		req.Host = vhBucket + ".localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 6: Upload Object with Nested Path (Virtual-Host Style)
	t.Run("Upload Object Nested Path Virtual-Host Style", func(t *testing.T) {
		content := "nested file contents"
		req := httptest.NewRequest(http.MethodPut, "/assets/images/logo.png", strings.NewReader(content))
		req.Host = vhBucket + ".localhost:8080"
		req.Header.Set("Content-Type", "image/png")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	// Test 7: Get Object with Nested Path (Virtual-Host Style)
	t.Run("Get Object Nested Path Virtual-Host Style", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/assets/images/logo.png", nil)
		req.Host = vhBucket + ".localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}
		bodyBytes, _ := io.ReadAll(w.Body)
		if string(bodyBytes) != "nested file contents" {
			t.Errorf("expected body 'nested file contents', got %q", string(bodyBytes))
		}
	})

	// Test 8: List Objects (Recursive Verification)
	t.Run("List Objects Recursive", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Host = vhBucket + ".localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
		}

		var objects []models.ObjectMetadata
		if err := json.NewDecoder(w.Body).Decode(&objects); err != nil {
			t.Fatal(err)
		}

		if len(objects) != 1 {
			t.Errorf("expected 1 object, got %d", len(objects))
		} else {
			if objects[0].Key != "assets/images/logo.png" {
				t.Errorf("expected key 'assets/images/logo.png', got %q", objects[0].Key)
			}
			if objects[0].ContentType != "image/png" {
				t.Errorf("expected content type 'image/png', got %q", objects[0].ContentType)
			}
		}
	})

	// Test 9: Delete Object
	t.Run("Delete Object", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/assets/images/logo.png", nil)
		req.Host = vhBucket + ".localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status 204, got %d", w.Code)
		}

		// Double check it's gone
		getReq := httptest.NewRequest(http.MethodGet, "/assets/images/logo.png", nil)
		getReq.Host = vhBucket + ".localhost:8080"
		getW := httptest.NewRecorder()
		router.ServeHTTP(getW, getReq)

		if getW.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", getW.Code)
		}
	})
}

// TestS3HandlerInvalidRequest checks edge cases
func TestS3HandlerInvalidRequest(t *testing.T) {
	router := RegisterRoutes()

	t.Run("Empty Bucket Name Path-Style", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Host = "localhost:8080"
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d. Body: %s", w.Code, w.Body.String())
		}
		if !bytes.Contains(w.Body.Bytes(), []byte("bucket name required")) {
			t.Errorf("expected 'bucket name required' error msg, got %q", w.Body.String())
		}
	})
}
