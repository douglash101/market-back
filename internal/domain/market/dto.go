package market

import "github.com/google/uuid"

type MarketCreateDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type MarketFoundDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}
