package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var uploadsDir = "uploads"

// saveUpload reads the "image" file field from a multipart form (if present)
// and stores it under uploadsDir, returning its public path. If no file was
// selected it returns fallback unchanged, so edits don't wipe an existing photo.
// Any failure is logged (not silently swallowed) since a misconfigured
// UPLOADS_DIR would otherwise fail every upload with no visible sign.
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
	dstPath := filepath.Join(uploadsDir, name)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("upload failed: create %s: %v", dstPath, err)
		return fallback
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		log.Printf("upload failed: write %s: %v", dstPath, err)
		return fallback
	}
	return "/uploads/" + name
}
