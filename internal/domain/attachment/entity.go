package attachment

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	URL         string    `json:"url" db:"url"`
	Type        *string   `json:"type" db:"type"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewAttachment(url string, attachmentType, description *string) *Attachment {
	return &Attachment{
		ID:          uuid.New(),
		URL:         url,
		Type:        attachmentType,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
