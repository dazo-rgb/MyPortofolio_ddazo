package repository

import (
	"dazo-portfolio/internal/domain"
	"dazo-portfolio/internal/usecase"
	"testing"
)

func TestRepositoryFiltersFeaturedProjects(t *testing.T) {
	repo := NewInMemoryProjectRepository([]domain.Project{{Slug: "featured", Featured: true}, {Slug: "hidden"}})
	got := repo.ListFeatured()
	if len(got) != 1 || got[0].Slug != "featured" {
		t.Fatalf("ListFeatured()=%#v", got)
	}
}
func TestRepositoryReturnsProjectAndNotFound(t *testing.T) {
	repo := NewInMemoryProjectRepository([]domain.Project{{Slug: "demo"}})
	got, err := repo.GetBySlug("demo")
	if err != nil || got.Slug != "demo" {
		t.Fatalf("GetBySlug()=%#v,%v", got, err)
	}
	if _, err := repo.GetBySlug("missing"); err != usecase.ErrProjectNotFound {
		t.Fatalf("missing error=%v", err)
	}
}
func TestRepositoryDeepCopiesInputAndOutput(t *testing.T) {
	input := []domain.Project{{Slug: "demo", Stack: []string{"Go"}, Features: []string{"Search"}, TechnicalWork: []string{"API"}, EngineeringNotes: []string{"Read-only"}}}
	repo := NewInMemoryProjectRepository(input)
	input[0].Stack[0] = "changed"
	input[0].Features[0] = "changed"
	got, err := repo.GetBySlug("demo")
	if err != nil {
		t.Fatal(err)
	}
	if got.Stack[0] != "Go" || got.Features[0] != "Search" {
		t.Fatalf("input slices leaked: %#v", got)
	}
	got.Stack[0] = "changed"
	got.Features[0] = "changed"
	got.TechnicalWork[0] = "changed"
	got.EngineeringNotes[0] = "changed"
	again, err := repo.GetBySlug("demo")
	if err != nil {
		t.Fatal(err)
	}
	if again.Stack[0] != "Go" || again.Features[0] != "Search" || again.TechnicalWork[0] != "API" || again.EngineeringNotes[0] != "Read-only" {
		t.Fatalf("output slices leaked: %#v", again)
	}
}
func TestListFeaturedReturnsIndependentSlices(t *testing.T) {
	repo := NewInMemoryProjectRepository([]domain.Project{{Slug: "demo", Featured: true, Stack: []string{"Go"}}})
	first := repo.ListFeatured()
	first[0].Stack[0] = "changed"
	second := repo.ListFeatured()
	if second[0].Stack[0] != "Go" {
		t.Fatalf("ListFeatured leaked slice: %#v", second[0].Stack)
	}
}
