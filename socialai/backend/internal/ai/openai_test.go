package ai

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIGeneratesFromBase64(t *testing.T) {
	png := append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 600)...)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/images/generations" || r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("unexpected request path=%q auth=%q", r.URL.Path, r.Header.Get("Authorization"))
		}
		var request struct {
			Model  string `json:"model"`
			Prompt string `json:"prompt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatal(err)
		}
		if request.Model != "gpt-image-2" || request.Prompt != "a moon garden" {
			t.Errorf("unexpected request: %+v", request)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{map[string]string{"b64_json": base64.StdEncoding.EncodeToString(png)}}})
	}))
	defer server.Close()

	client, err := NewOpenAI("test-key", "gpt-image-2", server.URL)
	if err != nil {
		t.Fatal(err)
	}
	image, err := client.Generate(context.Background(), "a moon garden")
	if err != nil {
		t.Fatal(err)
	}
	if string(image.Data) != string(png) || image.ContentType != "image/png" || image.Extension != ".png" {
		t.Fatalf("image = %+v", image)
	}
}

func TestOpenAIRejectsInvalidBase64(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":[{"b64_json":"not base64!"}]}`))
	}))
	defer server.Close()
	client, _ := NewOpenAI("test-key", "gpt-image-2", server.URL)
	if _, err := client.Generate(context.Background(), "prompt"); err == nil {
		t.Fatal("invalid base64 image was accepted")
	}
}

func TestOpenAIRejectsBase64ThatIsNotARasterImage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"data": []any{map[string]string{
			"b64_json": base64.StdEncoding.EncodeToString([]byte("not actually an image")),
		}}})
	}))
	defer server.Close()
	client, _ := NewOpenAI("test-key", "gpt-image-2", server.URL)
	if _, err := client.Generate(context.Background(), "prompt"); err == nil {
		t.Fatal("non-image base64 payload was accepted")
	}
}

func TestOpenAIImageDownloadRejectsLocalAddresses(t *testing.T) {
	if _, err := safeImageDial(context.Background(), "tcp", "127.0.0.1:443"); err == nil {
		t.Fatal("loopback generated image address was accepted")
	}
	if _, err := safeImageDial(context.Background(), "tcp", "10.0.0.1:443"); err == nil {
		t.Fatal("private generated image address was accepted")
	}
}
