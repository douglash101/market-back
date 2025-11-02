package product

import (
	"github.com/google/uuid"
)

// Product DTOs
type ProductCreateDTO struct {
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Name       string    `json:"name" validate:"required,min=3,max=100"`
	ImageURL   *string   `json:"image_url,omitempty" validate:"omitempty,url,max=180"`
}
