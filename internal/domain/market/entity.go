package market

import (
	"time"

	"github.com/google/uuid"
)

type MarketStatus string

const (
	MarketStatusActive   MarketStatus = "active"
	MarketStatusInactive MarketStatus = "inactive"
	MarketStatusDeleted  MarketStatus = "deleted"
)

type market struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      MarketStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// NewMarket creates a new organization with default values
func NewMarket(name, description string) *market {
	return &market{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Status:      MarketStatusActive,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Activate sets the organization status to active
func (o *market) Activate() {
	o.Status = MarketStatusActive
	o.UpdatedAt = time.Now()
}

// Deactivate sets the organization status to inactive
func (o *market) Deactivate() {
	o.Status = MarketStatusInactive
	o.UpdatedAt = time.Now()
}

// Delete sets the organization status to deleted
func (o *market) Delete() {
	o.Status = MarketStatusDeleted
	o.UpdatedAt = time.Now()
}

// IsActive checks if the organization is active
func (o *market) IsActive() bool {
	return o.Status == MarketStatusActive
}
