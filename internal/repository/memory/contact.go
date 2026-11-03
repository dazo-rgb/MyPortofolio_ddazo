package memory

import (
	"context"
	"sync"

	"github.com/dazo-rgb/ddazo-portfolio/internal/domain"
)

// MessageRepository is a development-only in-memory adapter.
// Messages are lost when the server restarts.
type MessageRepository struct {
	mu       sync.RWMutex
	messages []domain.ContactMessage
}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{messages: make([]domain.ContactMessage, 0)}
}

func (r *MessageRepository) Save(_ context.Context, message domain.ContactMessage) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, message)
	return nil
}
