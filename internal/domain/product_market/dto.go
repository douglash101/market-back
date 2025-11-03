package product_market

import (
	"github.com/google/uuid"
)

// ProductMarket DTOs
type ProductMarketCreateDTO struct {
	ProviderID       *string   `json:"provider_id,omitempty"`
	ProductID        uuid.UUID `json:"product_id" validate:"required"`
	MarketID         uuid.UUID `json:"market_id" validate:"required"`
	Price            float64   `json:"price" validate:"required,gt=0"`
	PromotionalPrice *float64  `json:"promotional_price,omitempty" validate:"omitempty,gt=0"`
}

type ProductMarketResponseDTO struct {
	ID               uuid.UUID           `json:"id"`
	ProviderID       *string             `json:"provider_id,omitempty"`
	ProductID        uuid.UUID           `json:"product_id"`
	MarketID         uuid.UUID           `json:"market_id"`
	Price            float64             `json:"price"`
	PromotionalPrice *float64            `json:"promotional_price,omitempty"`
	Status           ProductMarketStatus `json:"status"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
}
