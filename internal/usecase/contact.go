package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/dazo-rgb/ddazo-portfolio/internal/domain"
)

type ContactService struct {
	repository domain.ContactRepository
	now        func() time.Time
}

func NewContactService(repository domain.ContactRepository) *ContactService {
	return &ContactService{
		repository: repository,
		now:        time.Now,
	}
}

func (s *ContactService) Submit(ctx context.Context, message domain.ContactMessage) error {
	message.Name = strings.TrimSpace(message.Name)
	message.Email = strings.TrimSpace(message.Email)
	message.Message = strings.TrimSpace(message.Message)

	if err := message.Validate(); err != nil {
		return err
	}

	message.CreatedAt = s.now()
	return s.repository.Save(ctx, message)
}
