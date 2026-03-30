package images

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestCacheHit(t *testing.T) {
	dir := t.TempDir()
	indexPath := filepath.Join(dir, "index.json")
	idx := map[string]Entry{
		"https://example.com/img.jpg": {FileName: "abc.jpg", PageUID: "20230130_3961"},
	}
	data, _ := json.MarshalIndent(idx, "", "  ")
	os.WriteFile(indexPath, data, 0o644)

	// Create the cached file so CopyTo won't fail.
	uidDir := filepath.Join(dir, "20230130_3961")
	os.MkdirAll(uidDir, 0o755)
	os.WriteFile(filepath.Join(uidDir, "abc.jpg"), []byte("fake image"), 0o644)

	cache := NewCache(dir)
	entry, err := cache.Get("https://example.com/img.jpg", "20230130_3961")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry.FileName != "abc.jpg" {
		t.Errorf("FileName = %q, want abc.jpg", entry.FileName)
	}
}

func TestCacheMissDownloads(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		w.Write([]byte("fake jpeg data"))
	}))
	defer ts.Close()

	dir := t.TempDir()
	cache := NewCache(dir)
	entry, err := cache.Get(ts.URL+"/photo.jpg", "20230130_3961")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry.FileName == "" {
		t.Error("expected non-empty FileName")
	}
	if entry.PageUID != "20230130_3961" {
		t.Errorf("PageUID = %q, want 20230130_3961", entry.PageUID)
	}

	// Verify file was written.
	cached := filepath.Join(dir, "20230130_3961", entry.FileName)
	if _, err := os.Stat(cached); err != nil {
		t.Errorf("cached file not found: %v", err)
	}

	// Verify index was updated.
	indexData, _ := os.ReadFile(filepath.Join(dir, "index.json"))
	var idx map[string]Entry
	json.Unmarshal(indexData, &idx)
	if _, ok := idx[ts.URL+"/photo.jpg"]; !ok {
		t.Error("expected URL in index.json")
	}
}

func TestCleanshotRedirect(t *testing.T) {
	directURL := ""
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/share/abc+":
			w.Header().Set("Location", directURL+"/direct/photo.jpg?response-content-disposition=attachment%3Bfilename%3Dscreenshot.png")
			w.WriteHeader(http.StatusFound)
		case r.URL.Path == "/direct/photo.jpg":
			w.Header().Set("Content-Type", "image/png")
			w.Write([]byte("png data"))
		}
	}))
	defer ts.Close()
	directURL = ts.URL

	dir := t.TempDir()
	cache := NewCache(dir)
	entry, err := cache.Get(ts.URL+"/share/abc", "20230130_3961")
	if err != nil {
		t.Fatalf("Get() error: %v", err)
	}
	if entry.FileName == "" {
		t.Error("expected non-empty FileName")
	}
}

func TestCopyTo(t *testing.T) {
	dir := t.TempDir()
	uidDir := filepath.Join(dir, "20230130_3961")
	os.MkdirAll(uidDir, 0o755)
	os.WriteFile(filepath.Join(uidDir, "abc.jpg"), []byte("image data"), 0o644)

	cache := NewCache(dir)
	destDir := t.TempDir()
	err := cache.CopyTo(Entry{FileName: "abc.jpg", PageUID: "20230130_3961"}, destDir)
	if err != nil {
		t.Fatalf("CopyTo() error: %v", err)
	}

	destFile := filepath.Join(destDir, "abc.jpg")
	data, err := os.ReadFile(destFile)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	if string(data) != "image data" {
		t.Errorf("data = %q, want 'image data'", data)
	}
}
