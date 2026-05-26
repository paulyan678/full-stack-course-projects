package media

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStorageRoundTripAndTraversal(t *testing.T) {
	directory := t.TempDir()
	storage, err := NewLocal(directory, "http://localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	object, err := storage.Save(context.Background(), "photo.png", "image/png", bytes.NewBufferString("image"))
	if err != nil || object.URL != "http://localhost:8080/media/photo.png" {
		t.Fatalf("save object = %+v, %v", object, err)
	}
	if raw, err := os.ReadFile(filepath.Join(directory, "photo.png")); err != nil || string(raw) != "image" {
		t.Fatalf("read stored object = %q, %v", raw, err)
	}
	if _, err := storage.Save(context.Background(), "../escape", "text/plain", bytes.NewBufferString("bad")); err == nil {
		t.Fatal("path traversal key accepted")
	}
	listing := httptest.NewRecorder()
	storage.Handler().ServeHTTP(listing, httptest.NewRequest(http.MethodGet, "/", nil))
	if listing.Code != http.StatusNotFound {
		t.Fatalf("media directory listing status = %d", listing.Code)
	}
	if err := storage.Delete(context.Background(), object.Key); err != nil {
		t.Fatal(err)
	}
}
