package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"socialai/internal/ai"
	"socialai/internal/auth"
	"socialai/internal/media"
	"socialai/internal/model"
	"socialai/internal/repository"
)

type contextKey string

const usernameKey contextKey = "username"

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{2,31}$`)

type Settings struct {
	AllowedOrigins []string
	MaxUploadBytes int64
	TrustProxy     bool
}

type generatedImage struct {
	ID        string
	User      string
	Prompt    string
	Object    media.Object
	CreatedAt time.Time
}

type loginWindow struct {
	Count int
	Start time.Time
}

type Server struct {
	repository repository.Repository
	storage    media.Storage
	generator  ai.Generator
	tokens     *auth.TokenManager
	settings   Settings
	logger     *slog.Logger
	dummyHash  string
	generated  map[string]generatedImage
	generatedM sync.Mutex
	logins     map[string]loginWindow
	loginsM    sync.Mutex
}

func New(repo repository.Repository, storage media.Storage, generator ai.Generator, tokens *auth.TokenManager, settings Settings, logger *slog.Logger) *Server {
	dummy, _ := auth.HashPassword("invalid-password-value")
	if logger == nil {
		logger = slog.Default()
	}
	return &Server{repository: repo, storage: storage, generator: generator, tokens: tokens, settings: settings, logger: logger,
		dummyHash: dummy, generated: make(map[string]generatedImage), logins: make(map[string]loginWindow)}
}

func (s *Server) Handler(mediaHandler http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("POST /api/auth/signup", s.signup)
	mux.HandleFunc("POST /api/auth/signin", s.signin)
	mux.HandleFunc("POST /signup", s.signup)
	mux.HandleFunc("POST /signin", s.signin)
	mux.HandleFunc("GET /api/posts", s.authenticated(s.search))
	mux.HandleFunc("POST /api/posts", s.authenticated(s.upload))
	mux.HandleFunc("DELETE /api/posts/{id}", s.authenticated(s.deletePost))
	mux.HandleFunc("GET /search", s.authenticated(s.search))
	mux.HandleFunc("POST /upload", s.authenticated(s.upload))
	mux.HandleFunc("DELETE /post/{id}", s.authenticated(s.deletePost))
	mux.HandleFunc("POST /api/ai/images", s.authenticated(s.generateImage))
	mux.HandleFunc("DELETE /api/ai/images/{id}", s.authenticated(s.discardImage))
	mux.HandleFunc("POST /api/ai/images/{id}/publish", s.authenticated(s.publishImage))
	if mediaHandler != nil {
		mux.Handle("GET /media/", http.StripPrefix("/media/", mediaHandler))
	}
	return s.recover(s.securityHeaders(s.cors(mux)))
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))
	if !usernamePattern.MatchString(input.Username) {
		writeError(w, http.StatusBadRequest, "invalid_username", "Username must be 3-32 letters, numbers, underscores, or dashes.")
		return
	}
	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_password", err.Error())
		return
	}
	user := model.User{Username: input.Username, PasswordHash: hash, CreatedAt: time.Now().UTC()}
	if err := s.repository.CreateUser(r.Context(), user); err != nil {
		if errors.Is(err, repository.ErrConflict) {
			writeError(w, http.StatusConflict, "username_taken", "That username is unavailable.")
			return
		}
		s.logger.Error("create user", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not create the account.")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"username": user.Username, "created_at": user.CreatedAt})
}

func (s *Server) signin(w http.ResponseWriter, r *http.Request) {
	client := s.clientIP(r)
	if !s.allowLogin(client) {
		w.Header().Set("Retry-After", "60")
		writeError(w, http.StatusTooManyRequests, "rate_limited", "Too many sign-in attempts. Try again shortly.")
		return
	}
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	username := strings.ToLower(strings.TrimSpace(input.Username))
	user, err := s.repository.GetUser(r.Context(), username)
	if err != nil {
		_ = auth.CheckPassword(s.dummyHash, input.Password)
	}
	if err != nil || !auth.CheckPassword(user.PasswordHash, input.Password) {
		writeError(w, http.StatusUnauthorized, "invalid_credentials", "Username or password is incorrect.")
		return
	}
	token, err := s.tokens.Issue(user.Username)
	if err != nil {
		s.logger.Error("issue token", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not sign in.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"token": token, "username": user.Username})
}

func (s *Server) search(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	posts, err := s.repository.SearchPosts(r.Context(), model.SearchFilter{
		User: strings.TrimSpace(r.URL.Query().Get("user")), Keywords: strings.TrimSpace(r.URL.Query().Get("keywords")), Limit: limit,
	})
	if err != nil {
		s.logger.Error("search posts", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not search posts.")
		return
	}
	if posts == nil {
		posts = []model.Post{}
	}
	writeJSON(w, http.StatusOK, posts)
}

func (s *Server) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.settings.MaxUploadBytes)
	if err := r.ParseMultipartForm(8 << 20); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_upload", "Upload is malformed or exceeds the configured size limit.")
		return
	}
	message := strings.TrimSpace(r.FormValue("message"))
	if len(message) == 0 || len(message) > 500 {
		writeError(w, http.StatusBadRequest, "invalid_message", "Message must contain 1 to 500 characters.")
		return
	}
	file, _, err := r.FormFile("media_file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "media_required", "Choose an image or video to upload.")
		return
	}
	defer file.Close()
	contentType, extension, reader, err := inspectMedia(file)
	if err != nil {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media", err.Error())
		return
	}
	id, err := randomID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not create the post.")
		return
	}
	object, err := s.storage.Save(r.Context(), id+extension, contentType, reader)
	if err != nil {
		s.logger.Error("save media", "error", err)
		writeError(w, http.StatusInternalServerError, "storage_error", "Could not store the media.")
		return
	}
	post := model.Post{ID: id, User: currentUser(r), Message: message, URL: object.URL, Type: mediaType(contentType), CreatedAt: time.Now().UTC()}
	if err := s.repository.CreatePost(r.Context(), post); err != nil {
		_ = s.storage.Delete(r.Context(), object.Key)
		s.logger.Error("create post", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not create the post.")
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (s *Server) deletePost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	post, err := s.repository.GetPost(r.Context(), id)
	if err != nil || post.User != currentUser(r) {
		writeError(w, http.StatusNotFound, "not_found", "Post was not found.")
		return
	}
	if err := s.repository.DeletePost(r.Context(), id, currentUser(r)); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "Post was not found.")
			return
		}
		s.logger.Error("delete post", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not delete the post.")
		return
	}
	if parsed, err := url.Parse(post.URL); err == nil {
		if err := s.storage.Delete(r.Context(), path.Base(parsed.Path)); err != nil {
			s.logger.Warn("delete post media", "post_id", id, "error", err)
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) generateImage(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Prompt string `json:"prompt"`
	}
	if err := decodeJSON(w, r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	input.Prompt = strings.TrimSpace(input.Prompt)
	if len(input.Prompt) < 3 || len(input.Prompt) > 1000 {
		writeError(w, http.StatusBadRequest, "invalid_prompt", "Prompt must contain 3 to 1000 characters.")
		return
	}
	image, err := s.generator.Generate(r.Context(), input.Prompt)
	if err != nil {
		s.logger.Error("generate image", "error", err)
		writeError(w, http.StatusBadGateway, "generation_failed", "Image generation failed. Try again.")
		return
	}
	id, err := randomID()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not save the image.")
		return
	}
	object, err := s.storage.Save(r.Context(), id+image.Extension, image.ContentType, bytes.NewReader(image.Data))
	if err != nil {
		s.logger.Error("save generated image", "error", err)
		writeError(w, http.StatusInternalServerError, "storage_error", "Could not save the generated image.")
		return
	}
	now := time.Now().UTC()
	s.generatedM.Lock()
	expired := s.purgeGenerated(now)
	s.generated[id] = generatedImage{ID: id, User: currentUser(r), Prompt: input.Prompt, Object: object, CreatedAt: now}
	s.generatedM.Unlock()
	s.deleteGeneratedObjects(r.Context(), expired)
	writeJSON(w, http.StatusCreated, map[string]string{"id": id, "url": object.URL, "prompt": input.Prompt})
}

func (s *Server) discardImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.generatedM.Lock()
	generated, ok := s.generated[id]
	if !ok || generated.User != currentUser(r) {
		s.generatedM.Unlock()
		writeError(w, http.StatusNotFound, "not_found", "Generated image was not found.")
		return
	}
	delete(s.generated, id)
	s.generatedM.Unlock()
	if err := s.storage.Delete(r.Context(), generated.Object.Key); err != nil {
		s.generatedM.Lock()
		if _, exists := s.generated[id]; !exists {
			s.generated[id] = generated
		}
		s.generatedM.Unlock()
		s.logger.Error("discard generated image", "image_id", id, "error", err)
		writeError(w, http.StatusInternalServerError, "storage_error", "Could not discard the generated image.")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) publishImage(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var input struct {
		Message string `json:"message"`
	}
	if r.ContentLength != 0 {
		if err := decodeJSON(w, r, &input); err != nil {
			writeError(w, http.StatusBadRequest, "invalid_request", err.Error())
			return
		}
	}
	message := strings.TrimSpace(input.Message)
	if len(message) > 500 {
		writeError(w, http.StatusBadRequest, "invalid_message", "Message must contain at most 500 characters.")
		return
	}
	s.generatedM.Lock()
	generated, ok := s.generated[id]
	if !ok || generated.User != currentUser(r) {
		s.generatedM.Unlock()
		writeError(w, http.StatusNotFound, "not_found", "Generated image was not found or expired.")
		return
	}
	if time.Since(generated.CreatedAt) > time.Hour {
		delete(s.generated, id)
		s.generatedM.Unlock()
		s.deleteGeneratedObjects(r.Context(), []generatedImage{generated})
		writeError(w, http.StatusNotFound, "not_found", "Generated image was not found or expired.")
		return
	}
	delete(s.generated, id)
	s.generatedM.Unlock()
	if message == "" {
		message = generated.Prompt
	}
	post := model.Post{ID: id, User: generated.User, Message: message, URL: generated.Object.URL, Type: "image", CreatedAt: time.Now().UTC()}
	if err := s.repository.CreatePost(r.Context(), post); err != nil {
		s.generatedM.Lock()
		if _, exists := s.generated[id]; !exists {
			s.generated[id] = generated
		}
		s.generatedM.Unlock()
		s.logger.Error("publish generated image", "error", err)
		writeError(w, http.StatusInternalServerError, "internal_error", "Could not publish the image.")
		return
	}
	writeJSON(w, http.StatusCreated, post)
}

func (s *Server) authenticated(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.Fields(header)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			writeError(w, http.StatusUnauthorized, "authentication_required", "A bearer token is required.")
			return
		}
		claims, err := s.tokens.Validate(parts[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid_token", "The token is invalid or expired.")
			return
		}
		next(w, r.WithContext(context.WithValue(r.Context(), usernameKey, claims.Subject)))
	}
}

func (s *Server) cors(next http.Handler) http.Handler {
	allowed := make(map[string]bool, len(s.settings.AllowedOrigins))
	for _, origin := range s.settings.AllowedOrigins {
		allowed[origin] = true
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && (allowed[origin] || allowed["*"]) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "600")
		}
		if r.Method == http.MethodOptions {
			if origin != "" && !allowed[origin] && !allowed["*"] {
				writeError(w, http.StatusForbidden, "origin_not_allowed", "This origin is not allowed.")
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				s.logger.Error("panic recovered", "value", recovered)
				writeError(w, http.StatusInternalServerError, "internal_error", "Unexpected server error.")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func inspectMedia(file multipart.File) (string, string, io.Reader, error) {
	header := make([]byte, 512)
	n, err := io.ReadFull(file, header)
	if err != nil && err != io.ErrUnexpectedEOF {
		return "", "", nil, errors.New("could not read media")
	}
	header = header[:n]
	contentType := strings.Split(http.DetectContentType(header), ";")[0]
	extensions := map[string]string{
		"image/jpeg": ".jpg", "image/png": ".png", "image/gif": ".gif", "image/webp": ".webp",
		"video/mp4": ".mp4", "video/webm": ".webm", "video/quicktime": ".mov",
	}
	extension, ok := extensions[contentType]
	if !ok {
		return "", "", nil, errors.New("Only JPEG, PNG, GIF, WebP, MP4, WebM, and QuickTime media are supported")
	}
	return contentType, extension, io.MultiReader(bytes.NewReader(header), file), nil
}

func mediaType(contentType string) string {
	if strings.HasPrefix(contentType, "video/") {
		return "video"
	}
	return "image"
}

func currentUser(r *http.Request) string {
	username, _ := r.Context().Value(usernameKey).(string)
	return username
}

func randomID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

func decodeJSON(w http.ResponseWriter, r *http.Request, target any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return errors.New("Request body must be valid JSON.")
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return errors.New("Request body must contain one JSON object.")
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}

func (s *Server) clientIP(r *http.Request) string {
	if s.settings.TrustProxy {
		if forwarded := strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]); forwarded != "" {
			return forwarded
		}
	}
	host := r.RemoteAddr
	if parsed, _, err := net.SplitHostPort(host); err == nil {
		host = parsed
	}
	if host = strings.Trim(host, "[]"); host != "" {
		return host
	}
	return "unknown"
}

func (s *Server) allowLogin(client string) bool {
	s.loginsM.Lock()
	defer s.loginsM.Unlock()
	now := time.Now()
	window := s.logins[client]
	if window.Start.IsZero() || now.Sub(window.Start) >= time.Minute {
		window = loginWindow{Start: now}
	}
	if window.Count >= 10 {
		return false
	}
	window.Count++
	s.logins[client] = window
	return true
}

func (s *Server) purgeGenerated(now time.Time) []generatedImage {
	expired := make([]generatedImage, 0)
	for id, generated := range s.generated {
		if now.Sub(generated.CreatedAt) > time.Hour {
			delete(s.generated, id)
			expired = append(expired, generated)
		}
	}
	return expired
}

func (s *Server) deleteGeneratedObjects(ctx context.Context, generated []generatedImage) {
	for _, image := range generated {
		if err := s.storage.Delete(ctx, image.Object.Key); err != nil {
			s.logger.Warn("delete expired generated image", "image_id", image.ID, "error", err)
		}
	}
}
