package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"socialai/internal/ai"
	"socialai/internal/auth"
	"socialai/internal/media"
	"socialai/internal/model"
	"socialai/internal/repository"
)

type testApp struct {
	server *httptest.Server
	repo   *repository.Memory
	tokens *auth.TokenManager
}

func newTestApp(t *testing.T) testApp {
	t.Helper()
	repo := repository.NewMemory()
	storage, err := media.NewLocal(filepath.Join(t.TempDir(), "media"), "http://assets.test")
	if err != nil {
		t.Fatal(err)
	}
	tokens := auth.NewTokenManager([]byte("01234567890123456789012345678901"), time.Hour)
	server := New(repo, storage, ai.Local{}, tokens, Settings{AllowedOrigins: []string{"http://localhost:5173"}, MaxUploadBytes: 2 << 20}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	httpServer := httptest.NewServer(server.Handler(storage.Handler()))
	t.Cleanup(httpServer.Close)
	return testApp{server: httpServer, repo: repo, tokens: tokens}
}

func (a testApp) json(t *testing.T, method, path, token, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, a.server.URL+path, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func decode[T any](t *testing.T, res *http.Response) T {
	t.Helper()
	defer res.Body.Close()
	var value T
	if err := json.NewDecoder(res.Body).Decode(&value); err != nil {
		t.Fatal(err)
	}
	return value
}

func TestAuthPostSearchDeleteAndAIFlow(t *testing.T) {
	app := newTestApp(t)
	res := app.json(t, http.MethodPost, "/api/auth/signup", "", `{"username":"Nova_7","password":"correct horse battery staple"}`)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("signup status = %d", res.StatusCode)
	}
	res.Body.Close()
	res = app.json(t, http.MethodPost, "/api/auth/signin", "", `{"username":"nova_7","password":"correct horse battery staple"}`)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("signin status = %d", res.StatusCode)
	}
	session := decode[struct {
		Token string `json:"token"`
	}](t, res)
	if session.Token == "" {
		t.Fatal("signin returned no token")
	}

	var upload bytes.Buffer
	writer := multipart.NewWriter(&upload)
	_ = writer.WriteField("message", "Sunset over the ocean")
	part, err := writer.CreateFormFile("media_file", "sunset.png")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write(append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 600)...))
	_ = writer.Close()
	req, _ := http.NewRequest(http.MethodPost, app.server.URL+"/api/posts", &upload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+session.Token)
	res, err = http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusCreated {
		t.Fatalf("upload status=%v err=%v", status(res), err)
	}
	post := decode[model.Post](t, res)
	if post.User != "nova_7" || post.Type != "image" {
		t.Fatalf("uploaded post = %+v", post)
	}

	res = app.json(t, http.MethodGet, "/api/posts?keywords=sunset+ocean", session.Token, "")
	posts := decode[[]model.Post](t, res)
	if res.StatusCode != http.StatusOK || len(posts) != 1 || posts[0].ID != post.ID {
		t.Fatalf("search status=%d posts=%+v", res.StatusCode, posts)
	}

	otherToken, _ := app.tokens.Issue("someone_else")
	res = app.json(t, http.MethodDelete, "/api/posts/"+post.ID, otherToken, "")
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("cross-user delete status = %d", res.StatusCode)
	}
	res.Body.Close()
	res = app.json(t, http.MethodDelete, "/api/posts/"+post.ID, session.Token, "")
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("owner delete status = %d", res.StatusCode)
	}
	res.Body.Close()

	res = app.json(t, http.MethodPost, "/api/ai/images", session.Token, `{"prompt":"a lantern city in the clouds"}`)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("generate status = %d", res.StatusCode)
	}
	generated := decode[struct{ ID, URL string }](t, res)
	if generated.ID == "" || !strings.Contains(generated.URL, "/media/") {
		t.Fatalf("generated image = %+v", generated)
	}
	res = app.json(t, http.MethodPost, "/api/ai/images/"+generated.ID+"/publish", session.Token, `{"message":"Cloud lanterns"}`)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("publish status = %d", res.StatusCode)
	}
	published := decode[model.Post](t, res)
	if published.Message != "Cloud lanterns" || published.Type != "image" {
		t.Fatalf("published post = %+v", published)
	}
}

func TestAuthenticationAndCORS(t *testing.T) {
	app := newTestApp(t)
	res := app.json(t, http.MethodGet, "/api/posts", "", "")
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unauthenticated search status = %d", res.StatusCode)
	}
	res.Body.Close()
	req, _ := http.NewRequest(http.MethodOptions, app.server.URL+"/api/posts", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNoContent || res.Header.Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("preflight status=%d origin=%q", res.StatusCode, res.Header.Get("Access-Control-Allow-Origin"))
	}
}

func TestRejectsUnsupportedUpload(t *testing.T) {
	app := newTestApp(t)
	token, _ := app.tokens.Issue("nova")
	var upload bytes.Buffer
	writer := multipart.NewWriter(&upload)
	_ = writer.WriteField("message", "not media")
	part, _ := writer.CreateFormFile("media_file", "notes.txt")
	_, _ = part.Write([]byte("plain text"))
	_ = writer.Close()
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, app.server.URL+"/upload", &upload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnsupportedMediaType {
		t.Fatalf("unsupported upload status = %d", res.StatusCode)
	}
}

func TestDiscardGeneratedImageChecksOwnershipAndDeletesMedia(t *testing.T) {
	app := newTestApp(t)
	ownerToken, _ := app.tokens.Issue("nova")
	otherToken, _ := app.tokens.Issue("someone_else")
	res := app.json(t, http.MethodPost, "/api/ai/images", ownerToken, `{"prompt":"a moon garden"}`)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("generate status = %d", res.StatusCode)
	}
	generated := decode[struct{ ID, URL string }](t, res)
	parsed, err := url.Parse(generated.URL)
	if err != nil {
		t.Fatal(err)
	}
	res, err = http.Get(app.server.URL + parsed.Path)
	if err != nil || res.StatusCode != http.StatusOK {
		t.Fatalf("generated media status=%v err=%v", status(res), err)
	}
	res.Body.Close()

	res = app.json(t, http.MethodDelete, "/api/ai/images/"+generated.ID, otherToken, "")
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("cross-user discard status = %d", res.StatusCode)
	}
	res.Body.Close()
	res = app.json(t, http.MethodDelete, "/api/ai/images/"+generated.ID, ownerToken, "")
	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("owner discard status = %d", res.StatusCode)
	}
	res.Body.Close()
	res, err = http.Get(app.server.URL + parsed.Path)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("discarded media status = %d", res.StatusCode)
	}
}

func TestClientIPOnlyTrustsConfiguredProxy(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://social.test/", nil)
	req.RemoteAddr = "198.51.100.20:4567"
	req.Header.Set("X-Forwarded-For", "203.0.113.9")
	untrusted := &Server{settings: Settings{TrustProxy: false}}
	if got := untrusted.clientIP(req); got != "198.51.100.20" {
		t.Fatalf("untrusted proxy client IP = %q", got)
	}
	trusted := &Server{settings: Settings{TrustProxy: true}}
	if got := trusted.clientIP(req); got != "203.0.113.9" {
		t.Fatalf("trusted proxy client IP = %q", got)
	}
}

func TestLoginRateLimitCapsAttemptsWithinAWindow(t *testing.T) {
	server := &Server{logins: make(map[string]loginWindow)}
	for attempt := 0; attempt < 10; attempt++ {
		if !server.allowLogin("198.51.100.20") {
			t.Fatalf("attempt %d was rejected too early", attempt+1)
		}
	}
	if server.allowLogin("198.51.100.20") {
		t.Fatal("eleventh attempt was allowed")
	}
}

func status(res *http.Response) int {
	if res == nil {
		return 0
	}
	return res.StatusCode
}
