package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const uploadsDir = "uploads"

// saveUpload reads the "image" file field from a multipart form (if present)
// and stores it under uploads/, returning its public path. If no file was
// selected it returns fallback unchanged, so edits don't wipe an existing photo.
func saveUpload(r *http.Request, field, fallback string) string {
	file, header, err := r.FormFile(field)
	if err != nil {
		return fallback
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		ext = ".jpg"
	}
	b := make([]byte, 8)
	rand.Read(b)
	name := hex.EncodeToString(b) + ext
	dst, err := os.Create(filepath.Join(uploadsDir, name))
	if err != nil {
		return fallback
	}
	defer dst.Close()
	io.Copy(dst, file)
	return "/uploads/" + name
}
