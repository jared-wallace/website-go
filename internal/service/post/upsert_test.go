package post

import (
	"context"
	"errors"
	"testing"

	"github.com/jared-wallace/website-go/internal/model"
	postrepo "github.com/jared-wallace/website-go/internal/repository/post"
)

func TestUpsertBySlug_NewPost(t *testing.T) {
	repo := &mockRepo{returnErr: postrepo.ErrNotFound, returnPost: nil}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.UpsertBySlug(context.Background(), "Test Title", "test-slug", "# Hello")
	if err != nil {
		t.Fatalf("UpsertBySlug (new): unexpected error: %v", err)
	}
	if repo.createPost == nil {
		t.Fatal("UpsertBySlug (new): repo.Create was not called")
	}
	if repo.createPost.Published {
		t.Error("UpsertBySlug (new): new post should be unpublished draft")
	}
	if repo.createPost.Slug != "test-slug" {
		t.Errorf("UpsertBySlug (new): slug = %q, want %q", repo.createPost.Slug, "test-slug")
	}
	if repo.createPost.Title != "Test Title" {
		t.Errorf("UpsertBySlug (new): title = %q, want %q", repo.createPost.Title, "Test Title")
	}
}

func TestUpsertBySlug_ExistingPost(t *testing.T) {
	existing := newTestPost(7, "Old Title", "test-slug")
	existing.Tags = "existing-tags"
	repo := &mockRepo{returnPost: existing, returnErr: nil}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.UpsertBySlug(context.Background(), "New Title", "test-slug", "# Updated")
	if err != nil {
		t.Fatalf("UpsertBySlug (existing): unexpected error: %v", err)
	}
	if repo.updatePost == nil {
		t.Fatal("UpsertBySlug (existing): repo.Update was not called")
	}
	if repo.updatePost.ID != 7 {
		t.Errorf("UpsertBySlug (existing): ID = %d, want 7", repo.updatePost.ID)
	}
	if repo.updatePost.Tags != "existing-tags" {
		t.Errorf("UpsertBySlug (existing): Tags = %q, want %q", repo.updatePost.Tags, "existing-tags")
	}
	if repo.updatePost.Title != "New Title" {
		t.Errorf("UpsertBySlug (existing): Title = %q, want %q", repo.updatePost.Title, "New Title")
	}
}

func TestUpsertBySlug_FindError(t *testing.T) {
	sentinel := errors.New("db connection lost")
	repo := &mockRepo{returnErr: sentinel, returnPost: nil}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.UpsertBySlug(context.Background(), "Title", "slug", "body")
	if !errors.Is(err, sentinel) {
		t.Errorf("UpsertBySlug (find error): got %v, want %v", err, sentinel)
	}
}

// newTestPost is a helper that creates a model.Post with minimal fields set.
func newTestPost(id int64, title, slug string) *model.Post {
	return &model.Post{ID: id, Title: title, Slug: slug}
}
