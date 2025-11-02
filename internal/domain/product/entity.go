package product

import (
	"time"

	"github.com/google/uuid"
)

type ProductStatus string

const (
	ProductStatusActive   ProductStatus = "active"
	ProductStatusInactive ProductStatus = "inactive"
	ProductStatusDeleted  ProductStatus = "deleted"
)

type ProductCategory struct {
	ID          uuid.UUID `json:"id"`
	CompanyID   uuid.UUID `json:"company_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Product representa um produto no sistema
type Product struct {
	ID         uuid.UUID     `json:"id"`
	CategoryID *uuid.UUID    `json:"category_id,omitempty"`
	ImageURL   *string       `json:"image_url,omitempty"`
	Name       string        `json:"name"`
	Unit       *string       `json:"unit,omitempty"`
	Status     ProductStatus `json:"status"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}

// ProductWithCategory representa um produto com informações da categoria
type ProductWithCategory struct {
	Product
	Category *ProductCategory `json:"category,omitempty"`
}
