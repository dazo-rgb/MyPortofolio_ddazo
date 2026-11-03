package usecase

import (
	"dazo-portfolio/internal/domain"
	"errors"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrInvalidSlug     = errors.New("invalid project slug")
)

// ProjectRepository is the output port required by the project use case.
type ProjectRepository interface {
	ListFeatured() []domain.Project
	GetBySlug(slug string) (domain.Project, error)
}
