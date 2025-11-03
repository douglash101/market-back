package product_market

import (
	"time"

	"github.com/google/uuid"
)

type ProductMarketStatus string

const (
	ProductMarketStatusActive   ProductMarketStatus = "active"
	ProductMarketStatusInactive ProductMarketStatus = "inactive"
	ProductMarketStatusDeleted  ProductMarketStatus = "deleted"
)

// ProductMarket representa a relação entre produto e mercado com preços
type ProductMarket struct {
	ID               uuid.UUID           `json:"id"`
	ProviderID       *string             `json:"provider_id,omitempty"`
	ProductID        uuid.UUID           `json:"product_id"`
	MarketID         uuid.UUID           `json:"market_id"`
	Price            float64             `json:"price"`
	PromotionalPrice *float64            `json:"promotional_price,omitempty"`
	Status           ProductMarketStatus `json:"status"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
}
