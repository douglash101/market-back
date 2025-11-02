package attachment

import "github.com/google/uuid"

type AttachmentCreateDTO struct {
	URL         string  `json:"url" validate:"required,url"`
	Type        *string `json:"type" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
}

type AttachmentUpdateDTO struct {
	URL         *string `json:"url" validate:"omitempty,url"`
	Type        *string `json:"type" validate:"omitempty"`
	Description *string `json:"description" validate:"omitempty"`
}

type AttachmentFoundDTO struct {
	ID          uuid.UUID `json:"id"`
	URL         string    `json:"url"`
	Type        *string   `json:"type,omitempty"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

type AttachmentListDTO struct {
	Attachments []AttachmentFoundDTO `json:"attachments"`
	Total       int                  `json:"total"`
}
