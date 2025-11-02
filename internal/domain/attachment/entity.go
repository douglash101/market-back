package attachment

import (
	"time"

	"github.com/google/uuid"
)

type Attachment struct {
	ID          uuid.UUID `json:"id" db:"id"`
	CompanyID   uuid.UUID `json:"company_id" db:"company_id"`
	URL         string    `json:"url" db:"url"`
	Type        *string   `json:"type" db:"type"`
	Description *string   `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

func NewAttachment(url string, attachmentType, description *string, companyID uuid.UUID) *Attachment {
	return &Attachment{
		ID:          uuid.New(),
		CompanyID:   companyID,
		URL:         url,
		Type:        attachmentType,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
