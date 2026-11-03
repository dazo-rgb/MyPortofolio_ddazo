package usecase

import (
	"dazo-portfolio/internal/domain"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// ProjectService contains application rules for reading portfolio projects.
type ProjectService struct{ repo ProjectRepository }

func NewProjectService(repo ProjectRepository) *ProjectService { return &ProjectService{repo: repo} }
func (service *ProjectService) ListFeatured() []domain.Project { return service.repo.ListFeatured() }
func (service *ProjectService) GetBySlug(slug string) (domain.Project, error) {
	slug = strings.TrimSpace(slug)
	if !slugPattern.MatchString(slug) {
		return domain.Project{}, ErrInvalidSlug
	}
	return service.repo.GetBySlug(slug)
}
