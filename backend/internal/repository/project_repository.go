package repository

import (
	"dazo-portfolio/internal/domain"
	"dazo-portfolio/internal/usecase"
)

var _ usecase.ProjectRepository = (*InMemoryProjectRepository)(nil)

// InMemoryProjectRepository is a read-only adapter for static case studies.
type InMemoryProjectRepository struct{ projects []domain.Project }

func NewInMemoryProjectRepository(projects []domain.Project) *InMemoryProjectRepository {
	copied := make([]domain.Project, len(projects))
	for i, project := range projects {
		copied[i] = cloneProject(project)
	}
	return &InMemoryProjectRepository{projects: copied}
}
func (repo *InMemoryProjectRepository) ListFeatured() []domain.Project {
	featured := make([]domain.Project, 0, len(repo.projects))
	for _, project := range repo.projects {
		if project.Featured {
			featured = append(featured, cloneProject(project))
		}
	}
	return featured
}
func (repo *InMemoryProjectRepository) GetBySlug(slug string) (domain.Project, error) {
	for _, project := range repo.projects {
		if project.Slug == slug {
			return cloneProject(project), nil
		}
	}
	return domain.Project{}, usecase.ErrProjectNotFound
}

func cloneProject(project domain.Project) domain.Project {
	project.Stack = cloneStrings(project.Stack)
	project.Features = cloneStrings(project.Features)
	project.TechnicalWork = cloneStrings(project.TechnicalWork)
	project.EngineeringNotes = cloneStrings(project.EngineeringNotes)
	return project
}

func cloneStrings(values []string) []string {
	if values == nil {
		return nil
	}
	return append([]string(nil), values...)
}
