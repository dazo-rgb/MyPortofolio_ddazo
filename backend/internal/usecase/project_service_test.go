package usecase

import (
	"dazo-portfolio/internal/domain"
	"errors"
	"testing"
)

type stubProjectRepository struct {
	projects []domain.Project
	project  domain.Project
	err      error
	gotSlug  string
}

func (s *stubProjectRepository) ListFeatured() []domain.Project { return s.projects }
func (s *stubProjectRepository) GetBySlug(slug string) (domain.Project, error) {
	s.gotSlug = slug
	return s.project, s.err
}
func TestProjectServiceUsesRepositoryPort(t *testing.T) {
	repo := &stubProjectRepository{projects: []domain.Project{{Slug: "demo", Featured: true}}, project: domain.Project{Slug: "demo"}}
	svc := NewProjectService(repo)
	if got := svc.ListFeatured(); len(got) != 1 || got[0].Slug != "demo" {
		t.Fatalf("ListFeatured()=%#v", got)
	}
	got, err := svc.GetBySlug(" demo ")
	if err != nil || got.Slug != "demo" || repo.gotSlug != "demo" {
		t.Fatalf("GetBySlug()=%#v,%v slug=%q", got, err, repo.gotSlug)
	}
}
func TestProjectServiceRejectsInvalidSlugsBeforeRepository(t *testing.T) {
	repo := &stubProjectRepository{}
	svc := NewProjectService(repo)
	for _, slug := range []string{"", "../demo", "Demo", "two words", "demo/other", "demo--other"} {
		if _, err := svc.GetBySlug(slug); !errors.Is(err, ErrInvalidSlug) {
			t.Errorf("slug %q error=%v", slug, err)
		}
	}
	if repo.gotSlug != "" {
		t.Fatalf("repository called with invalid slug %q", repo.gotSlug)
	}
}
func TestProjectServicePropagatesRepositoryError(t *testing.T) {
	want := errors.New("storage unavailable")
	svc := NewProjectService(&stubProjectRepository{err: want})
	_, err := svc.GetBySlug("demo")
	if !errors.Is(err, want) {
		t.Fatalf("error=%v want %v", err, want)
	}
}
