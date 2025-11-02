package product

import (
	"fmt"
	"market/internal/domain/market"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase interface {
	CreateProduct(dto *ProductCreateDTO) error
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

// Product methods
func (s *service) CreateProduct(dto *ProductCreateDTO) error {
	// Basic validation
	if dto.Name == "" {
		return fmt.Errorf("product name is required")
	}

	// Create product entity
	product := &Product{
		ID:         uuid.New(),
		CategoryID: &dto.CategoryID,
		Name:       dto.Name,
		ImageURL:   dto.ImageURL,
		Status:     ProductStatusActive,
	}

	// Save to repository
	_, err := s.repository.Save(product)
	if err != nil {
		s.log.Errorw("error saving product", "error", err)
		return fmt.Errorf("error saving product: %w", err)
	}
	return nil
}
