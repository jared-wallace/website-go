package post

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"testing"

	"github.com/jared-wallace/website-go/internal/model"
)

// mockRepo records calls for assertion in tests.
type mockRepo struct {
	createPost      *model.Post
	updatePost      *model.Post
	returnPost      *model.Post
	returnErr       error  // default error for most operations
	findErr         *error // overrides returnErr for FindBySlug when non-nil
	softDeleteID    int64
	restoreID       int64
	setPublishedID  int64
	setPublishedVal bool
}

func (m *mockRepo) ListPublished(_ context.Context, _, _ int) ([]model.Post, error) {
	return nil, m.returnErr
}
func (m *mockRepo) CountPublished(_ context.Context) (int, error) { return 0, m.returnErr }
func (m *mockRepo) FindBySlug(_ context.Context, _ string) (*model.Post, error) {
	if m.findErr != nil {
		return m.returnPost, *m.findErr
	}
	return m.returnPost, m.returnErr
}
func (m *mockRepo) FindByID(_ context.Context, _ int64) (*model.Post, error) {
	return m.returnPost, m.returnErr
}
func (m *mockRepo) ListAll(_ context.Context) ([]model.Post, error) {
	return nil, m.returnErr
}
func (m *mockRepo) Create(_ context.Context, p model.Post) (*model.Post, error) {
	m.createPost = &p
	if m.returnErr != nil {
		return nil, m.returnErr
	}
	p.ID = 42
	return &p, nil
}
func (m *mockRepo) Update(_ context.Context, p model.Post) error {
	m.updatePost = &p
	return m.returnErr
}
func (m *mockRepo) SoftDelete(_ context.Context, id int64) error {
	m.softDeleteID = id
	return m.returnErr
}
func (m *mockRepo) Restore(_ context.Context, id int64) error {
	m.restoreID = id
	return m.returnErr
}
func (m *mockRepo) SetPublished(_ context.Context, id int64, published bool) error {
	m.setPublishedID = id
	m.setPublishedVal = published
	return m.returnErr
}

func (m *mockRepo) AddReaction(_ context.Context, _ int64, _ string) (bool, error) {
	return false, nil
}

func (m *mockRepo) CountReactions(_ context.Context, _ int64) (int, error) {
	return 0, nil
}

// mockRenderer wraps input in <p> tags for deterministic test output.
type mockRenderer struct{}

func (r *mockRenderer) Render(src string) template.HTML {
	return template.HTML(fmt.Sprintf("<p>%s</p>", src))
}

func TestServiceSoftDelete(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.SoftDelete(context.Background(), 7)
	if err != nil {
		t.Fatalf("SoftDelete: unexpected error: %v", err)
	}
	if repo.softDeleteID != 7 {
		t.Errorf("SoftDelete: repo called with ID %d, want 7", repo.softDeleteID)
	}
}

func TestServiceSoftDeletePropagatesError(t *testing.T) {
	sentinel := errors.New("db error")
	repo := &mockRepo{returnErr: sentinel}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.SoftDelete(context.Background(), 1)
	if !errors.Is(err, sentinel) {
		t.Errorf("SoftDelete: got error %v, want %v", err, sentinel)
	}
}

func TestServiceRestore(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.Restore(context.Background(), 3)
	if err != nil {
		t.Fatalf("Restore: unexpected error: %v", err)
	}
	if repo.restoreID != 3 {
		t.Errorf("Restore: repo called with ID %d, want 3", repo.restoreID)
	}
}

func TestServicePublish(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.Publish(context.Background(), 5)
	if err != nil {
		t.Fatalf("Publish: unexpected error: %v", err)
	}
	if repo.setPublishedID != 5 {
		t.Errorf("Publish: repo called with ID %d, want 5", repo.setPublishedID)
	}
	if !repo.setPublishedVal {
		t.Errorf("Publish: SetPublished called with false, want true")
	}
}

func TestServiceUnpublish(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	err := svc.Unpublish(context.Background(), 9)
	if err != nil {
		t.Fatalf("Unpublish: unexpected error: %v", err)
	}
	if repo.setPublishedID != 9 {
		t.Errorf("Unpublish: repo called with ID %d, want 9", repo.setPublishedID)
	}
	if repo.setPublishedVal {
		t.Errorf("Unpublish: SetPublished called with true, want false")
	}
}

func TestServiceCreate(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	body := "hello world"
	_, err := svc.Create(context.Background(), "Title", "title", body, "tag1", false)
	if err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}

	if repo.createPost == nil {
		t.Fatal("Create: repo.Create was not called")
	}
	want := template.HTML("<p>hello world</p>")
	if repo.createPost.RenderedHTML != string(want) {
		t.Errorf("Create: RenderedHTML = %q, want %q", repo.createPost.RenderedHTML, want)
	}
}

func TestServiceUpdate(t *testing.T) {
	repo := &mockRepo{}
	svc := newServiceWithRenderer(repo, &mockRenderer{})

	body := "updated body"
	err := svc.Update(context.Background(), 11, "New Title", "new-title", body, "tagA")
	if err != nil {
		t.Fatalf("Update: unexpected error: %v", err)
	}

	if repo.updatePost == nil {
		t.Fatal("Update: repo.Update was not called")
	}
	want := template.HTML("<p>updated body</p>")
	if repo.updatePost.RenderedHTML != string(want) {
		t.Errorf("Update: RenderedHTML = %q, want %q", repo.updatePost.RenderedHTML, want)
	}
}
