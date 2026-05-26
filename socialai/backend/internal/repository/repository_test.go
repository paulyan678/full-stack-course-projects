package repository

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"socialai/internal/model"
)

func exerciseRepository(t *testing.T, repo Repository) {
	t.Helper()
	ctx := context.Background()
	user := model.User{Username: "ivy", PasswordHash: "hash", CreatedAt: time.Now()}
	if err := repo.CreateUser(ctx, user); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreateUser(ctx, user); !errors.Is(err, ErrConflict) {
		t.Fatalf("duplicate user error = %v", err)
	}
	old := model.Post{ID: "one", User: "ivy", Message: "a quiet ocean", Type: "image", CreatedAt: time.Now().Add(-time.Hour)}
	newer := model.Post{ID: "two", User: "max", Message: "a loud ocean sunset", Type: "video", CreatedAt: time.Now()}
	if err := repo.CreatePost(ctx, old); err != nil {
		t.Fatal(err)
	}
	if err := repo.CreatePost(ctx, newer); err != nil {
		t.Fatal(err)
	}
	posts, err := repo.SearchPosts(ctx, model.SearchFilter{Keywords: "ocean sunset"})
	if err != nil || len(posts) != 1 || posts[0].ID != "two" {
		t.Fatalf("keyword search = %+v, %v", posts, err)
	}
	if err := repo.DeletePost(ctx, "one", "max"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-user delete error = %v", err)
	}
	if err := repo.DeletePost(ctx, "one", "ivy"); err != nil {
		t.Fatal(err)
	}
}

func TestMemoryRepository(t *testing.T) { exerciseRepository(t, NewMemory()) }

func TestFileRepositoryPersists(t *testing.T) {
	file := filepath.Join(t.TempDir(), "data", "socialai.json")
	repo, err := NewFile(file)
	if err != nil {
		t.Fatal(err)
	}
	exerciseRepository(t, repo)
	reloaded, err := NewFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := reloaded.GetUser(context.Background(), "IVY"); err != nil {
		t.Fatalf("persisted user missing: %v", err)
	}
	posts, err := reloaded.SearchPosts(context.Background(), model.SearchFilter{})
	if err != nil || len(posts) != 1 || posts[0].ID != "two" {
		t.Fatalf("persisted posts = %+v, %v", posts, err)
	}
}

func TestFileRepositoryRollsBackFailedPersistence(t *testing.T) {
	ctx := context.Background()
	directory := t.TempDir()
	validPath := filepath.Join(directory, "socialai.json")
	repo, err := NewFile(validPath)
	if err != nil {
		t.Fatal(err)
	}
	baseUser := model.User{Username: "ivy", PasswordHash: "hash", CreatedAt: time.Now()}
	if err := repo.CreateUser(ctx, baseUser); err != nil {
		t.Fatal(err)
	}
	basePost := model.Post{ID: "existing", User: "ivy", Message: "kept", CreatedAt: time.Now()}
	if err := repo.CreatePost(ctx, basePost); err != nil {
		t.Fatal(err)
	}

	blocker := filepath.Join(directory, "not-a-directory")
	if err := os.WriteFile(blocker, []byte("block"), 0o600); err != nil {
		t.Fatal(err)
	}
	repo.path = filepath.Join(blocker, "socialai.json")
	if err := repo.CreateUser(ctx, model.User{Username: "ghost", PasswordHash: "hash"}); err == nil {
		t.Fatal("create user unexpectedly succeeded when persistence failed")
	}
	if _, err := repo.GetUser(ctx, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("failed user remained in memory: %v", err)
	}
	if err := repo.CreatePost(ctx, model.Post{ID: "ghost", User: "ivy", Message: "ghost"}); err == nil {
		t.Fatal("create post unexpectedly succeeded when persistence failed")
	}
	if _, err := repo.GetPost(ctx, "ghost"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("failed post remained in memory: %v", err)
	}
	if err := repo.DeletePost(ctx, "existing", "ivy"); err == nil {
		t.Fatal("delete unexpectedly succeeded when persistence failed")
	}
	if post, err := repo.GetPost(ctx, "existing"); err != nil || post.Message != "kept" {
		t.Fatalf("failed delete was not rolled back: post=%+v err=%v", post, err)
	}
}
