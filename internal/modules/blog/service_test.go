package blog

import (
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupTestRepository(t *testing.T) (*PostService, string) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "posts.json")
	repo := NewJSONRepository(filePath)
	service := NewPostService(repo)
	return service, filePath
}

func TestPostService_CreatePost(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "First Post",
		Content:   "This is my first post",
		CreatedAt: createdAt,
	}

	// Test Create
	err := service.Create(post)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Verify
	posts, _ := service.Index()
	if len(posts) != 1 {
		t.Errorf("expected 1 post after create, got %d", len(posts))
	}

	if posts[0].Title != post.Title {
		t.Errorf("expected title %s, got %s", post.Title, posts[0].Title)
	}

	if posts[0].ID != 1 {
		t.Errorf("expected ID 1 for first post, got %d", posts[0].ID)
	}
}

func TestPostService_CreateMultiplePosts(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")

	// Create multiple posts
	for i := 1; i <= 3; i++ {
		post := Post{
			Title:     "Post " + string(rune(i)),
			Content:   "Content " + string(rune(i)),
			CreatedAt: createdAt,
		}
		if err := service.Create(post); err != nil {
			t.Fatalf("create post %d failed: %v", i, err)
		}
	}

	// Verify
	posts, _ := service.Index()
	if len(posts) != 3 {
		t.Errorf("expected 3 posts, got %d", len(posts))
	}

	// Verify IDs are auto-incremented
	for i, post := range posts {
		if post.ID != i+1 {
			t.Errorf("post %d: expected ID %d, got %d", i, i+1, post.ID)
		}
	}
}

func TestPostService_ShowPost(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "Show Test Post",
		Content:   "Content to show",
		CreatedAt: createdAt,
	}
	service.Create(post)

	// Test Show
	retrieved, err := service.Show(1)
	if err != nil {
		t.Fatalf("show failed: %v", err)
	}

	// Verify
	if retrieved == nil {
		t.Fatal("expected post to be returned, got nil")
	}

	if retrieved.Title != post.Title {
		t.Errorf("expected title %s, got %s", post.Title, retrieved.Title)
	}

	if retrieved.Content != post.Content {
		t.Errorf("expected content %s, got %s", post.Content, retrieved.Content)
	}
}

func TestPostService_ShowPostNotFound(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	// Test Show with non-existent ID
	_, err := service.Show(999)

	// Verify
	if err == nil {
		t.Error("expected error for non-existent post, got nil")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' error, got: %v", err)
	}
}

func TestPostService_EditPost(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "Original Title",
		Content:   "Original Content",
		CreatedAt: createdAt,
	}
	service.Create(post)

	// Test Edit
	updatedAt, _ := time.Parse("2006-01-02", "2026-07-02")
	updatedPost := Post{
		Title:     "Updated Title",
		Content:   "Updated Content",
		CreatedAt: updatedAt,
	}

	err := service.Edit(1, updatedPost)
	if err != nil {
		t.Fatalf("edit failed: %v", err)
	}

	// Verify
	retrieved, _ := service.Show(1)
	if retrieved.Title != updatedPost.Title {
		t.Errorf("expected title %s, got %s", updatedPost.Title, retrieved.Title)
	}

	if retrieved.Content != updatedPost.Content {
		t.Errorf("expected content %s, got %s", updatedPost.Content, retrieved.Content)
	}

	if retrieved.CreatedAt != updatedAt {
		t.Errorf("expected createdAt %v, got %v", updatedAt, retrieved.CreatedAt)
	}
}

func TestPostService_DeletePost(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "To Delete",
		Content:   "This will be deleted",
		CreatedAt: createdAt,
	}
	service.Create(post)

	// Verify post was created
	posts, _ := service.Index()
	if len(posts) != 1 {
		t.Fatal("post was not created")
	}

	// Test Delete
	err := service.Delete(1)
	if err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	// Verify
	posts, _ = service.Index()
	if len(posts) != 0 {
		t.Errorf("expected 0 posts after delete, got %d", len(posts))
	}

	_, err = service.Show(1)
	if err == nil {
		t.Error("expected error when retrieving deleted post")
	}
}

func TestPostService_DeleteNonExistent(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	// Test Delete non-existent post (should not error, just filter it out)
	err := service.Delete(999)

	// Verify - delete should not error even if post doesn't exist
	if err != nil {
		t.Errorf("delete non-existent post should not error, got: %v", err)
	}
}

func TestPostService_Index(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")

	// Create posts
	for i := 1; i <= 5; i++ {
		post := Post{
			Title:     "Post " + string(rune(i)),
			Content:   "Content " + string(rune(i)),
			CreatedAt: createdAt,
		}
		service.Create(post)
	}

	// Test Index
	posts, err := service.Index()
	if err != nil {
		t.Fatalf("index failed: %v", err)
	}

	// Verify
	if len(posts) != 5 {
		t.Errorf("expected 5 posts in index, got %d", len(posts))
	}
}

func TestPostService_IndexEmpty(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	// Test Index on empty repository
	posts, err := service.Index()
	if err != nil {
		t.Fatalf("index on empty repo failed: %v", err)
	}

	// Verify
	if len(posts) != 0 {
		t.Errorf("expected 0 posts in empty index, got %d", len(posts))
	}
}

func TestPostService_TrimSpaceOnCreate(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "Test",
		Content:   "   Content with spaces   \n\n",
		CreatedAt: createdAt,
	}

	// Test Create
	service.Create(post)

	// Verify content is trimmed
	retrieved, _ := service.Show(1)
	if retrieved.Content != "Content with spaces" {
		t.Errorf("expected trimmed content, got: %q", retrieved.Content)
	}
}

func TestPostService_EditNonExistentPost(t *testing.T) {
	// Setup
	service, _ := setupTestRepository(t)

	createdAt, _ := time.Parse("2006-01-02", "2026-07-01")
	post := Post{
		Title:     "Test",
		Content:   "Content",
		CreatedAt: createdAt,
	}

	// Test Edit non-existent post - should not error, just not update anything
	err := service.Edit(999, post)
	if err != nil {
		t.Fatalf("edit non-existent post should not error, got: %v", err)
	}

	// Verify nothing was added
	posts, _ := service.Index()
	if len(posts) != 0 {
		t.Errorf("expected 0 posts after editing non-existent, got %d", len(posts))
	}
}

func TestPostService_RepositoryError(t *testing.T) {
	// Setup with an invalid/unreachable file path
	// Use an invalid path that repository will fail to access
	invalidPath := "/invalid/nonexistent/path/that/will/fail/posts.json"
	
	repo := NewJSONRepository(invalidPath)
	service := NewPostService(repo)

	// Test Index should fail when repository cannot read
	// This will fail on save, so first try to create which should error
	post := Post{
		ID:        1,
		Title:     "Test",
		Content:   "Test",
		CreatedAt: time.Now(),
	}

	err := service.Create(post)
	if err == nil {
		t.Error("expected error when repository fails, got nil")
	}
}

