package admin_test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createMultipartRequest builds a multipart POST request with a file field named "image".
func createMultipartRequest(t *testing.T, filename string, content []byte) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("Write: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/admin/images/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

// uploadTestSetup returns a testSetup with imageDir set to a temporary directory.
func uploadTestSetup(t *testing.T) (*testSetup, string) {
	t.Helper()
	ts := newTestSetup(t)
	dir := t.TempDir()
	ts.handler.SetImageDir(dir)
	return ts, dir
}

func TestUploadImage_ValidJPEG(t *testing.T) {
	ts, dir := uploadTestSetup(t)

	// Real JPEG magic bytes followed by padding
	content := make([]byte, 1024)
	copy(content, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	req := createMultipartRequest(t, "photo.jpg", content)
	rec := ts.serve(ts.handler.UploadImage, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ValidJPEG: got status %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ValidJPEG: JSON decode: %v; body: %s", err, rec.Body.String())
	}

	url, ok := resp["url"]
	if !ok {
		t.Fatal("ValidJPEG: response missing 'url' key")
	}
	tag, ok := resp["markdownTag"]
	if !ok {
		t.Fatal("ValidJPEG: response missing 'markdownTag' key")
	}

	if !strings.HasSuffix(url, ".jpg") {
		t.Errorf("ValidJPEG: url %q does not end with .jpg", url)
	}

	// Filename should be 32 hex chars + ".jpg"
	fname := filepath.Base(url)
	if len(fname) != 36 { // 32 hex + 4 ".jpg"
		t.Errorf("ValidJPEG: filename %q length = %d, want 36", fname, len(fname))
	}

	if !strings.Contains(tag, url) {
		t.Errorf("ValidJPEG: markdownTag %q does not contain url %q", tag, url)
	}

	// File should exist on disk
	diskPath := filepath.Join(dir, fname)
	if _, err := os.Stat(diskPath); os.IsNotExist(err) {
		t.Errorf("ValidJPEG: file not written to disk at %s", diskPath)
	}
}

func TestUploadImage_ValidPNG(t *testing.T) {
	ts, _ := uploadTestSetup(t)

	content := make([]byte, 1024)
	// Full 8-byte PNG signature required for http.DetectContentType
	copy(content, []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})

	req := createMultipartRequest(t, "screenshot.png", content)
	rec := ts.serve(ts.handler.UploadImage, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ValidPNG: got status %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("ValidPNG: JSON decode: %v", err)
	}

	if !strings.HasSuffix(resp["url"], ".png") {
		t.Errorf("ValidPNG: url %q does not end with .png", resp["url"])
	}
}

func TestUploadImage_SpoofedMIME(t *testing.T) {
	ts, _ := uploadTestSetup(t)

	content := []byte("<html><body>evil</body></html>")
	req := createMultipartRequest(t, "evil.jpg", content)
	rec := ts.serve(ts.handler.UploadImage, req)

	if rec.Code != http.StatusUnsupportedMediaType {
		t.Errorf("SpoofedMIME: got status %d, want %d", rec.Code, http.StatusUnsupportedMediaType)
	}

	body := strings.TrimSpace(rec.Body.String())
	if body != "only JPEG and PNG accepted" {
		t.Errorf("SpoofedMIME: body = %q, want %q", body, "only JPEG and PNG accepted")
	}
}

func TestUploadImage_TooLarge(t *testing.T) {
	ts, _ := uploadTestSetup(t)

	// 6 MB payload
	content := make([]byte, 6*1024*1024)
	req := createMultipartRequest(t, "huge.jpg", content)
	rec := ts.serve(ts.handler.UploadImage, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("TooLarge: got status %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestUploadImage_RandomFilename(t *testing.T) {
	ts, _ := uploadTestSetup(t)

	content := make([]byte, 1024)
	copy(content, []byte{0xFF, 0xD8, 0xFF, 0xE0})

	req1 := createMultipartRequest(t, "same.jpg", content)
	rec1 := ts.serve(ts.handler.UploadImage, req1)

	req2 := createMultipartRequest(t, "same.jpg", content)
	rec2 := ts.serve(ts.handler.UploadImage, req2)

	if rec1.Code != http.StatusOK || rec2.Code != http.StatusOK {
		t.Fatalf("RandomFilename: uploads failed: %d, %d", rec1.Code, rec2.Code)
	}

	var resp1, resp2 map[string]string
	json.Unmarshal(rec1.Body.Bytes(), &resp1)
	json.Unmarshal(rec2.Body.Bytes(), &resp2)

	if resp1["url"] == resp2["url"] {
		t.Errorf("RandomFilename: two uploads produced same url: %q", resp1["url"])
	}
}

func TestUploadImage_NoFile(t *testing.T) {
	ts, _ := uploadTestSetup(t)

	// POST without any file field
	req := httptest.NewRequest(http.MethodPost, "/admin/images/upload", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := ts.serve(ts.handler.UploadImage, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("NoFile: got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
