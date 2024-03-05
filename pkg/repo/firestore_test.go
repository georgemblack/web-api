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

	// Create like
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

	// Create post
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
