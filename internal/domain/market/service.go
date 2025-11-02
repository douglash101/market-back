package market

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase interface {
	Create(input *MarketCreateDTO) (*MarketFoundDTO, error)
	FindByID(id uuid.UUID) (*MarketFoundDTO, error)
}

type service struct {
	log        *zap.SugaredLogger
	repository Repository
}

func NewService(
	log *zap.SugaredLogger,
) UseCase {
	return &service{
		log:        log,
		repository: NewRepository(log),
	}
}

func (s *service) Create(input *MarketCreateDTO) (*MarketFoundDTO, error) {
	market := NewMarket(
		input.Name,
		input.Description,
	)

	createdMarket, err := s.repository.Create(market)
	if err != nil {
		s.log.Errorw("error creating market", "error", err)
		return nil, err
	}

	return &MarketFoundDTO{
		ID:          createdMarket.ID,
		Name:        createdMarket.Name,
		Description: createdMarket.Description,
		Status:      string(createdMarket.Status),
	}, nil
}

func (s *service) FindByID(id uuid.UUID) (*MarketFoundDTO, error) {
	market, err := s.repository.FindByID(id)
	if err != nil {
		s.log.Errorw("error finding market by ID", "id", id, "error", err)
		return nil, err
	}
	if market == nil {
		s.log.Warnw("market not found", "id", id)
		return nil, nil
	}
	return &MarketFoundDTO{
		ID:          market.ID,
		Name:        market.Name,
		Description: market.Description,
		Status:      string(market.Status),
		CreatedAt:   market.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   market.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
