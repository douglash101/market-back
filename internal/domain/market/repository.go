package market

import (
	"database/sql"
	"market/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	Create(market *market) (*market, error)
	FindByID(id uuid.UUID) (*market, error)
}

type repository struct {
	db              *database.PostgresDB
	log             *zap.SugaredLogger
	createStatement *sql.Stmt
}

func NewRepository(
	log *zap.SugaredLogger,
) Repository {

	dbInstance := database.GetInstance(log)

	insert := `INSERT INTO public.markets
		(id, name, description, created_at, updated_at)
	VALUES
		($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	createStatement, err := dbInstance.Prepare(insert)
	if err != nil {
		log.Errorw("error on create statement", "error", err)
		return nil
	}

	return &repository{
		db:              dbInstance,
		log:             log,
		createStatement: createStatement,
	}
}

func (o *repository) Create(market *market) (*market, error) {
	_, err := o.createStatement.Exec(
		market.ID,
		market.Name,
		market.Description,
	)

	if err != nil {
		o.log.Errorw("error on execute Create", "error", err)
		return nil, err
	}
	o.log.Infow("market created successfully", "id", market.ID, "name", market.Name)
	return market, nil
}

func (o *repository) FindByID(id uuid.UUID) (*market, error) {
	sql := `SELECT id, name, description, created_at, updated_at
	FROM markets WHERE id = $1 LIMIT 1`
	row, err := o.db.Query(sql, id)

	if err != nil {
		o.log.Errorw("error on execute FindByID", "error", err)
		return nil, err
	}

	defer row.Close()

	var org market
	if row.Next() {
		err = row.Scan(
			&org.ID,
			&org.Name,
			&org.Description,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			o.log.Errorw("error on scan FindByID", "error", err)
			return nil, err
		}
		return &org, nil
	}

	return nil, nil
}

func (o *repository) FindByName(name string) (*market, error) {
	sql := `SELECT id, name, description, created_at, updated_at
	FROM markets WHERE name = $1 LIMIT 1`
	row, err := o.db.Query(sql, name)

	if err != nil {
		o.log.Errorw("error on execute FindByName", "error", err)
		return nil, err
	}

	defer row.Close()

	var org market
	if row.Next() {
		err = row.Scan(
			&org.ID,
			&org.Name,
			&org.Description,
			&org.CreatedAt,
			&org.UpdatedAt,
		)
		if err != nil {
			o.log.Errorw("error on scan FindByName", "error", err)
			return nil, err
		}
		return &org, nil
	}

	return nil, nil
}
