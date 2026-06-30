package blog

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestJSONRepository_LoadEmptyFile(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "posts.json")

	// Create empty file
	if err := os.WriteFile(filePath, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	repo := NewJSONRepository(filePath)

	// Test
	posts, err := repo.Load()

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected empty posts, got %d posts", len(posts))
	}
}

func TestJSONRepository_LoadNonExistentFile(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.json")

	repo := NewJSONRepository(filePath)

	// Test
	posts, err := repo.Load()

	// Assert
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("expected empty posts for nonexistent file, got %d posts", len(posts))
	}
}

func TestJSONRepository_SaveAndLoad(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "posts.json")

	repo := NewJSONRepository(filePath)

	createdAt := time.Now()
	testPosts := []Post{
		{
			ID:        1,
			Title:     "Test Post 1",
			Content:   "Content 1",
			CreatedAt: createdAt,
		},
		{
			ID:        2,
			Title:     "Test Post 2",
			Content:   "Content 2",
			CreatedAt: createdAt,
		},
	}

	// Test Save
	err := repo.Save(testPosts)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatal("expected file to be created")
	}

	// Test Load
	loadedPosts, err := repo.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	// Assert
	if len(loadedPosts) != len(testPosts) {
		t.Errorf("expected %d posts, got %d", len(testPosts), len(loadedPosts))
	}

	for i, post := range loadedPosts {
		if post.ID != testPosts[i].ID {
			t.Errorf("post %d: expected ID %d, got %d", i, testPosts[i].ID, post.ID)
		}
		if post.Title != testPosts[i].Title {
			t.Errorf("post %d: expected title %s, got %s", i, testPosts[i].Title, post.Title)
		}
		if post.Content != testPosts[i].Content {
			t.Errorf("post %d: expected content %s, got %s", i, testPosts[i].Content, post.Content)
		}
	}
}

func TestJSONRepository_LoadInvalidJSON(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.json")

	if err := os.WriteFile(filePath, []byte("invalid json {"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	repo := NewJSONRepository(filePath)

	// Test
	_, err := repo.Load()

	// Assert
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestJSONRepository_ConcurrentOperations(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "posts.json")

	repo := NewJSONRepository(filePath)

	posts := []Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", CreatedAt: time.Now()},
	}

	err := repo.Save(posts)
	if err != nil {
		t.Fatalf("initial save failed: %v", err)
	}

	// Simulate concurrent reads and writes
	done := make(chan error, 2)

	go func() {
		for i := 0; i < 5; i++ {
			_, err := repo.Load()
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()

	go func() {
		for i := 0; i < 5; i++ {
			err := repo.Save(posts)
			if err != nil {
				done <- err
				return
			}
		}
		done <- nil
	}()

	// Assert
	for i := 0; i < 2; i++ {
		if err := <-done; err != nil {
			t.Errorf("concurrent operation failed: %v", err)
		}
	}
}

func TestJSONRepository_SaveCreatesProperJSON(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "posts.json")

	repo := NewJSONRepository(filePath)

	createdAt, _ := time.Parse(time.RFC3339, "2026-07-01T10:00:00Z")
	posts := []Post{
		{ID: 1, Title: "Test", Content: "Content", CreatedAt: createdAt},
	}

	err := repo.Save(posts)
	if err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Read file content and verify JSON structure
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	var readPosts []Post
	err = json.Unmarshal(content, &readPosts)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(readPosts) != 1 {
		t.Errorf("expected 1 post in JSON, got %d", len(readPosts))
	}
}
