//go:build js && wasm

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

// HTTPFS is a file system that reads files from an HTTP server.
// It implements fs.FS and is used by the WASM client to load assets.
type HTTPFS struct {
	baseURL string
	client  *http.Client
}

func NewHTTPFS(baseURL string) *HTTPFS {
	return &HTTPFS{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (h *HTTPFS) Open(name string) (fs.File, error) {
	t0 := time.Now()
	url := h.baseURL + "/" + strings.TrimLeft(name, "/")
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP request failed: %s", resp.Status)
	}
	f := &httpFile{
		ReadCloser: resp.Body,
		name:       path.Base(name),
		size:       resp.ContentLength,
		modTime:    time.Now(), // Since we don't have the actual modTime
	}
	// Log the call, argument, size in mb, and loading time.
	dt := time.Since(t0)
	mb := float64(f.size) / 1e6
	log.Printf("HTTPFS.Open(%q) %d bytes (%.2f MB) in %s", name, f.size, mb, dt)
	return f, nil
}

type httpFile struct {
	io.ReadCloser
	name    string
	size    int64
	modTime time.Time
}

func (f *httpFile) Stat() (fs.FileInfo, error) {
	return &httpFileInfo{
		name:    f.name,
		size:    f.size,
		modTime: f.modTime,
	}, nil
}

type httpFileInfo struct {
	name    string
	size    int64
	modTime time.Time
}

func (fi *httpFileInfo) Name() string       { return fi.name }
func (fi *httpFileInfo) Size() int64        { return fi.size }
func (fi *httpFileInfo) Mode() fs.FileMode  { return 0444 } // Read-only
func (fi *httpFileInfo) ModTime() time.Time { return fi.modTime }
func (fi *httpFileInfo) IsDir() bool        { return false }
func (fi *httpFileInfo) Sys() interface{}   { return nil }
