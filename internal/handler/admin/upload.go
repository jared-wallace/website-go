package admin

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// UploadImage accepts a multipart POST with an "image" file field, validates
// the content type via magic-byte detection, generates a random filename, and
// writes the file to the configured image directory. Returns JSON with the
// public URL and a ready-to-paste markdown image tag.
func (h *AdminHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Extend the write deadline for potentially large uploads (D-14).
	if rc := http.NewResponseController(w); rc != nil {
		_ = rc.SetWriteDeadline(time.Now().Add(30 * time.Second))
	}

	// 5 MB hard limit on request body (D-04).
	const maxUpload = 5 << 20                              // 5 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload+512) // +512 for multipart overhead

	if err := r.ParseMultipartForm(maxUpload); err != nil {
		http.Error(w, "file too large or bad request", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "missing image field", http.StatusBadRequest)
		return
	}
	defer func() { _ = file.Close() }()

	// Sniff the first 512 bytes for magic-byte MIME detection (D-12).
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	mime := http.DetectContentType(buf[:n])

	ext := extensionFromMIME(mime)
	if ext == "" {
		http.Error(w, "only JPEG and PNG accepted", http.StatusUnsupportedMediaType)
		return
	}

	// Rewind past the sniffed bytes so the full file is written to disk.
	if _, seekErr := file.Seek(0, io.SeekStart); seekErr != nil {
		http.Error(w, "failed to process file", http.StatusInternalServerError)
		return
	}

	// Server-generated random hex filename -- never trust the client name.
	b := make([]byte, 16)
	if _, randErr := rand.Read(b); randErr != nil {
		slog.Error("crypto/rand failed", "error", randErr)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	filename := hex.EncodeToString(b) + ext

	dst, err := os.Create(filepath.Join(h.imageDir, filename)) //nolint:gosec // filename is server-generated random hex, not user input
	if err != nil {
		slog.Error("failed to create image file", "error", err)
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	if _, err = io.Copy(dst, file); err != nil {
		slog.Error("failed to write image file", "error", err)
		_ = dst.Close()
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	if err = dst.Close(); err != nil {
		slog.Error("failed to close image file", "error", err)
		http.Error(w, "failed to save image", http.StatusInternalServerError)
		return
	}

	imgURL := "/images/" + filename
	resp := map[string]string{
		"url":         imgURL,
		"markdownTag": "![alt](" + imgURL + ")",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode upload response", "error", err)
	}
}

// extensionFromMIME returns the file extension for allowed image MIME types.
// Returns empty string for anything we do not accept.
func extensionFromMIME(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}
