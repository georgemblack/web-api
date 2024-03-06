package repo

import (
	"testing"
	"time"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
)

func prep(t *testing.T) (FirestoreService, error) {
	// Configure env to look for google application credentials in correct location
	t.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/workspaces/web-api/google-application-credentials.json")

	config, err := conf.LoadConfig()
	if err != nil {
		t.Errorf("failed to load config; %s", err)
	}
	service, err := NewFirestoreService(config)
	if err != nil {
		t.Errorf("failed to create firestore service; %s", err)
	}
	return service, nil
}

func TestAddGetLike(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Add like
	expected := types.Like{
		ID:        "",
		Timestamp: time.Now(),
		Title:     "test title",
		URL:       "http://test.com",
	}
	id, err := service.AddLike(expected)
	if err != nil {
		t.Errorf("failed to post like; %s", err)
	}

	// Read like
	actual, err := service.GetLike(id)
	if err != nil {
		t.Errorf("failed to get like; %s", err)
	}

	// Delete like
	err = service.DeleteLike(id)
	if err != nil {
		t.Errorf("failed to delete like; %s", err)
	}

	// Compare
	if expected.Timestamp.Unix() != actual.Timestamp.Unix() {
		t.Errorf("expected timestamp %s, got %s", expected.Timestamp, actual.Timestamp)
	}
	if actual.Title != expected.Title {
		t.Errorf("expected title %s, got %s", expected.Title, actual.Title)
	}
	if actual.URL != expected.URL {
		t.Errorf("expected url %s, got %s", expected.URL, actual.URL)
	}
}

func TestAddGetPost(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Add post
	expected := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    true,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now(),
	}
	id, err := service.AddPost(expected)
	if err != nil {
		t.Errorf("failed to post post; %s", err)
	}

	// Read post
	actual, err := service.GetPost(id)
	if err != nil {
		t.Errorf("failed to get post; %s", err)
	}

	// Delete post
	err = service.DeletePost(id)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Compare
	if actual.Draft != expected.Draft {
		t.Errorf("expected draft %t, got %t", expected.Draft, actual.Draft)
	}
	if actual.Listed != expected.Listed {
		t.Errorf("expected listed %t, got %t", expected.Listed, actual.Listed)
	}
	if actual.Title != expected.Title {
		t.Errorf("expected title %s, got %s", expected.Title, actual.Title)
	}
	if actual.Slug != expected.Slug {
		t.Errorf("expected slug %s, got %s", expected.Slug, actual.Slug)
	}
	if actual.Content != expected.Content {
		t.Errorf("expected content %s, got %s", expected.Content, actual.Content)
	}
	for i, _ := range expected.Tags {
		if actual.Tags[i] != expected.Tags[i] {
			t.Errorf("expected tag %s, got %s", expected.Tags[i], actual.Tags[i])
		}
	}
	if actual.Published.Unix() != expected.Published.Unix() {
		t.Errorf("expected published %s, got %s", expected.Published, actual.Published)
	}
}

func TestGetPosts(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Add posts
	expected := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    true,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now(),
	}
	first, err := service.AddPost(expected)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	second, err := service.AddPost(expected)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Read posts
	posts, err := service.GetPosts(PostFilters{})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Delete posts
	err = service.DeletePost(first)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}
	err = service.DeletePost(second)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Validate posts contains at least two posts
	if len(posts) < 2 {
		t.Errorf("expected at least two posts, got %d", len(posts))
	}

	// Validate posts contains first
	found := false
	for _, post := range posts {
		if post.ID == first {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find post %s", first)
	}

	// Validate posts contains second
	found = false
	for _, post := range posts {
		if post.ID == second {
			found = true
		}
	}
	if !found {
		t.Errorf("expected to find post %s", second)
	}

	// Validate that second appears before first
	for _, post := range posts {
		if post.ID == second {
			break
		}
		if post.ID == first {
			t.Errorf("expected second post to appear before first")
		}
	}
}

func TestGetPostsWithFilters(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Add unlisted post
	unlisted := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    false,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now(),
	}
	id, err := service.AddPost(unlisted)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Read posts
	listed := true
	posts, err := service.GetPosts(PostFilters{
		Listed: &listed,
	})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Delete post
	err = service.DeletePost(id)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Validate posts does not contain unlisted post
	for _, post := range posts {
		if post.ID == id {
			t.Errorf("expected not to find post %s", id)
		}
	}
}
