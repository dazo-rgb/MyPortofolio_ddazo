package domain

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidMessage = errors.New("name, valid email, and message are required")
	ErrMessageTooLong = errors.New("message must be 3000 characters or fewer")
)

type ContactMessage struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

func (m ContactMessage) Validate() error {
	if strings.TrimSpace(m.Name) == "" ||
		!strings.Contains(strings.TrimSpace(m.Email), "@") ||
		strings.TrimSpace(m.Message) == "" {
		return ErrInvalidMessage
	}
	if len(m.Message) > 3000 {
		return ErrMessageTooLong
	}
	return nil
}

type ContactRepository interface {
	Save(ctx context.Context, message ContactMessage) error
}
