package product

import (
	"database/sql"
	"market/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	FindByID(id uuid.UUID) (*Product, error)
	Save(product *Product) (*Product, error)
}

type productRepository struct {
	db                *database.PostgresDB
	log               *zap.SugaredLogger
	createProductStmt *sql.Stmt
}

func NewRepository(log *zap.SugaredLogger) Repository {
	dbInstance := database.GetInstance(log)

	// Product statements
	insertProduct := `INSERT INTO products 
		(category_id, image_url, name, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	// Prepare statements
	createProductStmt, err := dbInstance.Prepare(insertProduct)
	if err != nil {
		log.Errorw("error preparing create product statement", "error", err)
	}

	return &productRepository{
		db:                dbInstance,
		log:               log,
		createProductStmt: createProductStmt,
	}
}

func (p *productRepository) FindByID(id uuid.UUID) (*Product, error) {
	sql := `SELECT id, category_id, image_url, name, status, created_at, updated_at
			FROM products WHERE id = $1 AND status != 'deleted' LIMIT 1`

	rows, err := p.db.Query(sql, id)
	if err != nil {
		p.log.Errorw("error executing FindByID", "error", err, "id", id)
		return nil, err
	}
	defer rows.Close()

	var product Product
	if rows.Next() {
		err = rows.Scan(
			&product.ID,
			&product.CategoryID,
			&product.ImageURL,
			&product.Name,
			&product.Status,
			&product.CreatedAt,
			&product.UpdatedAt,
		)

		if err != nil {
			p.log.Errorw("error scanning product by ID", "error", err, "id", id)
			return nil, err
		}

		return &product, nil
	}

	return nil, nil
}

func (p *productRepository) Save(product *Product) (*Product, error) {
	err := p.createProductStmt.QueryRow(
		product.CategoryID,
		product.ImageURL,
		product.Name,
		product.Status,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		p.log.Errorw("error saving product", "error", err)
		return nil, err
	}

	return product, nil
}
