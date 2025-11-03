package product_market

import (
	"fmt"
	"market/internal/domain/market"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase interface {
	CreateProductMarket(dto *ProductMarketCreateDTO) (*ProductMarketResponseDTO, error)
	FindByProviderID(providerID string) ([]*ProductMarket, error)
}

type service struct {
	log           *zap.SugaredLogger
	repository    Repository
	marketService market.UseCase
}

func NewService(
	log *zap.SugaredLogger,
) UseCase {
	return &service{
		log:           log,
		repository:    NewRepository(log),
		marketService: market.NewService(log),
	}
}

// ProductMarket methods
func (s *service) CreateProductMarket(dto *ProductMarketCreateDTO) (*ProductMarketResponseDTO, error) {
	// Basic validation
	if dto.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than 0")
	}

	if dto.PromotionalPrice != nil && *dto.PromotionalPrice <= 0 {
		return nil, fmt.Errorf("promotional price must be greater than 0")
	}

	// Create product market entity
	productMarket := &ProductMarket{
		ID:               uuid.New(),
		ProviderID:       dto.ProviderID,
		ProductID:        dto.ProductID,
		MarketID:         dto.MarketID,
		Price:            dto.Price,
		PromotionalPrice: dto.PromotionalPrice,
		Status:           ProductMarketStatusActive,
	}

	// Save to repository
	savedProductMarket, err := s.repository.Save(productMarket)
	if err != nil {
		s.log.Errorw("error saving product market", "error", err)
		return nil, fmt.Errorf("error saving product market: %w", err)
	}

	// Convert to response DTO
	responseDTO := &ProductMarketResponseDTO{
		ID:               savedProductMarket.ID,
		ProviderID:       savedProductMarket.ProviderID,
		ProductID:        savedProductMarket.ProductID,
		MarketID:         savedProductMarket.MarketID,
		Price:            savedProductMarket.Price,
		PromotionalPrice: savedProductMarket.PromotionalPrice,
		Status:           savedProductMarket.Status,
		CreatedAt:        savedProductMarket.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        savedProductMarket.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return responseDTO, nil
}

func (s *service) FindByProviderID(providerID string) ([]*ProductMarket, error) {
	if providerID == "" {
		return nil, fmt.Errorf("provider ID is required")
	}

	productMarkets, err := s.repository.FindByProviderID(providerID)
	if err != nil {
		s.log.Errorw("error finding product markets by provider ID", "error", err, "provider_id", providerID)
		return nil, fmt.Errorf("error finding product markets: %w", err)
	}

	return productMarkets, nil
}
