package product_market

import (
	"database/sql"
	"market/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	FindByID(id uuid.UUID) (*ProductMarket, error)
	FindByProviderID(providerID string) ([]*ProductMarket, error)
	Save(productMarket *ProductMarket) (*ProductMarket, error)
}

type productMarketRepository struct {
	db                      *database.PostgresDB
	log                     *zap.SugaredLogger
	createProductMarketStmt *sql.Stmt
	findByProviderIDStmt    *sql.Stmt
}

func NewRepository(log *zap.SugaredLogger) Repository {
	dbInstance := database.GetInstance(log)

	// ProductMarket statements
	insertProductMarket := `INSERT INTO product_markets 
		(id, provider_id, product_id, market_id, price, promotional_price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING created_at, updated_at`

	findByProviderIDQuery := `SELECT id, provider_id, product_id, market_id, price, promotional_price, status, created_at, updated_at
							 FROM product_markets WHERE provider_id = $1 AND status != 'deleted'`

	// Prepare statements
	createProductMarketStmt, err := dbInstance.Prepare(insertProductMarket)
	if err != nil {
		log.Errorw("error preparing create product market statement", "error", err)
	}

	findByProviderIDStmt, err := dbInstance.Prepare(findByProviderIDQuery)
	if err != nil {
		log.Errorw("error preparing find by provider ID statement", "error", err)
	}

	return &productMarketRepository{
		db:                      dbInstance,
		log:                     log,
		createProductMarketStmt: createProductMarketStmt,
		findByProviderIDStmt:    findByProviderIDStmt,
	}
}

func (p *productMarketRepository) FindByID(id uuid.UUID) (*ProductMarket, error) {
	sql := `SELECT id, provider_id, product_id, market_id, price, promotional_price, status, created_at, updated_at
			FROM product_markets WHERE id = $1 AND status != 'deleted' LIMIT 1`

	rows, err := p.db.Query(sql, id)
	if err != nil {
		p.log.Errorw("error executing FindByID", "error", err, "id", id)
		return nil, err
	}
	defer rows.Close()

	var productMarket ProductMarket
	if rows.Next() {
		err = rows.Scan(
			&productMarket.ID,
			&productMarket.ProviderID,
			&productMarket.ProductID,
			&productMarket.MarketID,
			&productMarket.Price,
			&productMarket.PromotionalPrice,
			&productMarket.Status,
			&productMarket.CreatedAt,
			&productMarket.UpdatedAt,
		)

		if err != nil {
			p.log.Errorw("error scanning product market by ID", "error", err, "id", id)
			return nil, err
		}

		return &productMarket, nil
	}

	return nil, nil
}

func (p *productMarketRepository) FindByProviderID(providerID string) ([]*ProductMarket, error) {
	rows, err := p.findByProviderIDStmt.Query(providerID)
	if err != nil {
		p.log.Errorw("error executing FindByProviderID", "error", err, "provider_id", providerID)
		return nil, err
	}
	defer rows.Close()

	var productMarkets []*ProductMarket
	for rows.Next() {
		var productMarket ProductMarket
		err = rows.Scan(
			&productMarket.ID,
			&productMarket.ProviderID,
			&productMarket.ProductID,
			&productMarket.MarketID,
			&productMarket.Price,
			&productMarket.PromotionalPrice,
			&productMarket.Status,
			&productMarket.CreatedAt,
			&productMarket.UpdatedAt,
		)

		if err != nil {
			p.log.Errorw("error scanning product market by provider ID", "error", err, "provider_id", providerID)
			return nil, err
		}

		productMarkets = append(productMarkets, &productMarket)
	}

	if err = rows.Err(); err != nil {
		p.log.Errorw("error iterating product markets by provider ID", "error", err, "provider_id", providerID)
		return nil, err
	}

	return productMarkets, nil
}

func (p *productMarketRepository) Save(productMarket *ProductMarket) (*ProductMarket, error) {
	err := p.createProductMarketStmt.QueryRow(
		productMarket.ID,
		productMarket.ProviderID,
		productMarket.ProductID,
		productMarket.MarketID,
		productMarket.Price,
		productMarket.PromotionalPrice,
		productMarket.Status,
	).Scan(&productMarket.CreatedAt, &productMarket.UpdatedAt)

	if err != nil {
		p.log.Errorw("error saving product market", "error", err)
		return nil, err
	}

	return productMarket, nil
}
