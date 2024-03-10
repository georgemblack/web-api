package repo

import (
	"fmt"
	"testing"
	"time"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/google/uuid"
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

func postIn(id string, posts []types.Post) bool {
	for _, post := range posts {
		if post.ID == id {
			return true
		}
	}
	return false
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
	for i := range expected.Tags {
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

	// Add first post
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

	// Add second post, adding one hour to 'published' time
	expected.Published = expected.Published.Add(time.Hour)
	second, err := service.AddPost(expected)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Read all posts
	posts, err := service.GetPosts(PostFilters{})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Delete newly created posts
	err = service.DeletePost(first)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}
	err = service.DeletePost(second)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Validate results contain at least two posts
	if len(posts) < 2 {
		t.Errorf("expected at least two posts, got %d", len(posts))
	}

	// Validate newly created posts are in the result
	foundFirst := false
	foundSecond := false
	for _, post := range posts {
		if post.ID == first {
			foundFirst = true
		}
		if post.ID == second {
			foundSecond = true
		}
	}
	if !foundFirst {
		t.Errorf("expected to find post: %s", first)
	}
	if !foundSecond {
		t.Errorf("expected to find post: %s", second)
	}

	// Validate second post appears before first
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

	// Add 'unlisted' post
	unlisted := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    false,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now().Add(-time.Hour),
	}
	unlistedID, err := service.AddPost(unlisted)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Fetch all 'listed' posts
	listedBool := true
	posts, err := service.GetPosts(PostFilters{Listed: &listedBool})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Validate 'unlisted' post is not in the result
	if postIn(unlistedID, posts) {
		t.Errorf("unlisted post should not be in the result: %s", unlistedID)
	}

	// Delete 'unlisted' post
	err = service.DeletePost(unlistedID)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Add 'draft' post
	draft := types.Post{
		ID:        "",
		Draft:     true,
		Listed:    false,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now().Add(-time.Hour),
	}
	draftID, err := service.AddPost(draft)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Fetch all 'published' posts
	publishedBool := true
	posts, err = service.GetPosts(PostFilters{Published: &publishedBool})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Validate 'draft' post is not in the result
	if postIn(draftID, posts) {
		t.Errorf("draft post should not be in the result: %s", draftID)
	}

	// Delete 'draft' post
	err = service.DeletePost(draftID)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Add 'future' post, with publish date in the future
	future := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    true,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now().Add(time.Hour),
	}
	futureID, err := service.AddPost(future)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Fetch all 'published' posts
	publishedBool = true
	posts, err = service.GetPosts(PostFilters{Published: &publishedBool})
	if err != nil {
		t.Errorf("failed to get posts; %s", err)
	}

	// Validate 'future' post is not in the result
	if postIn(futureID, posts) {
		t.Errorf("future post should not be in the result: %s", futureID)
	}

	// Delete 'future' post
	err = service.DeletePost(futureID)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}
}

func TestUpdatePost(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Create post
	post := types.Post{
		ID:        "",
		Draft:     false,
		Listed:    true,
		Title:     "test title",
		Slug:      "test-title",
		Content:   "#test content",
		Tags:      []string{"test", "tag"},
		Published: time.Now(),
	}
	postID, err := service.AddPost(post)
	if err != nil {
		t.Errorf("failed to add post; %s", err)
	}

	// Update post
	post.ID = postID
	post.Title = "updated title"
	post.Slug = "updated-title"
	post.Content = "#updated content"
	post.Tags = []string{}

	err = service.UpdatePost(post)
	if err != nil {
		t.Errorf("failed to update post; %s", err)
	}

	// Read post
	actual, err := service.GetPost(postID)
	if err != nil {
		t.Errorf("failed to get post; %s", err)
	}

	// Delete post
	err = service.DeletePost(postID)
	if err != nil {
		t.Errorf("failed to delete post; %s", err)
	}

	// Compare
	if actual.Draft != post.Draft {
		t.Errorf("expected draft %t, got %t", post.Draft, actual.Draft)
	}
	if actual.Listed != post.Listed {
		t.Errorf("expected listed %t, got %t", post.Listed, actual.Listed)
	}
	if actual.Title != post.Title {
		t.Errorf("expected title %s, got %s", post.Title, actual.Title)
	}
	if actual.Slug != post.Slug {
		t.Errorf("expected slug %s, got %s", post.Slug, actual.Slug)
	}
	if actual.Content != post.Content {
		t.Errorf("expected content %s, got %s", post.Content, actual.Content)
	}
	if len(actual.Tags) != 0 {
		t.Errorf("expected no tags, got %v", actual.Tags)
	}
	if actual.Published.Unix() != post.Published.Unix() {
		t.Errorf("expected published %s, got %s", post.Published, actual.Published)
	}
}

func TestGetUpdateHashList(t *testing.T) {
	service, err := prep(t)
	if err != nil {
		t.Errorf("failed to prep test; %s", err)
	}

	// Create hash list
	expected := types.HashList{
		Hashes: map[string]string{
			fmt.Sprintf("test-%s", uuid.New().String()): uuid.New().String(),
			fmt.Sprintf("test-%s", uuid.New().String()): uuid.New().String(),
		},
	}
	err = service.UpdateHashList(expected)
	if err != nil {
		t.Errorf("failed to update hash list; %s", err)
	}

	// Read hash list
	actual, err := service.GetHashList()
	if err != nil {
		t.Errorf("failed to get hash list; %s", err)
	}

	// Compare
	if len(actual.Hashes) != 2 {
		t.Errorf("expected 2 hashes in result, got %d", len(actual.Hashes))
	}
	for k, v := range expected.Hashes {
		if actual.Hashes[k] != v {
			t.Errorf("expected hash '%s', got '%s'", v, actual.Hashes[k])
		}
	}
}
