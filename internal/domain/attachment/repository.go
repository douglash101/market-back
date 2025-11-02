package attachment

import (
	"database/sql"
	"market/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	Create(attachment *Attachment) (*Attachment, error)
	FindByID(id uuid.UUID, companyID uuid.UUID) (*Attachment, error)
	Update(id uuid.UUID, companyID uuid.UUID, attachment *Attachment) (*Attachment, error)
	Delete(id uuid.UUID, companyID uuid.UUID) error
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

	insert := `INSERT INTO public.attachments
		(id, company_id, url, type, description, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

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

func (o *repository) Create(attachment *Attachment) (*Attachment, error) {
	_, err := o.createStatement.Exec(
		attachment.ID,
		attachment.CompanyID,
		attachment.URL,
		attachment.Type,
		attachment.Description,
	)

	if err != nil {
		o.log.Errorw("error on execute Create", "error", err)
		return nil, err
	}
	o.log.Infow("attachment created successfully", "id", attachment.ID, "url", attachment.URL)
	return attachment, nil
}

func (o *repository) FindByID(id uuid.UUID, companyID uuid.UUID) (*Attachment, error) {
	sql := `SELECT id, company_id, url, type, description, created_at, updated_at
	FROM attachments WHERE id = $1 and company_id = $2 LIMIT 1`
	row, err := o.db.Query(sql, id, companyID)

	if err != nil {
		o.log.Errorw("error on execute FindByID", "error", err)
		return nil, err
	}

	defer row.Close()

	var attachment Attachment
	if row.Next() {
		err = row.Scan(
			&attachment.ID,
			&attachment.CompanyID,
			&attachment.URL,
			&attachment.Type,
			&attachment.Description,
			&attachment.CreatedAt,
			&attachment.UpdatedAt,
		)
		if err != nil {
			o.log.Errorw("error on scan FindByID", "error", err)
			return nil, err
		}
		return &attachment, nil
	}

	return nil, nil
}

func (o *repository) Update(id uuid.UUID, companyID uuid.UUID, attachment *Attachment) (*Attachment, error) {
	sql := `UPDATE attachments SET 
		url = $2, type = $3, description = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND company_id = $5`

	_, err := o.db.Exec(sql, id, attachment.URL, attachment.Type, attachment.Description, companyID)
	if err != nil {
		o.log.Errorw("error on execute Update", "error", err)
		return nil, err
	}

	return o.FindByID(id, companyID)
}

func (o *repository) Delete(id uuid.UUID, companyID uuid.UUID) error {
	sql := `DELETE FROM attachments WHERE id = $1 AND company_id = $2`

	_, err := o.db.Exec(sql, id)
	if err != nil {
		o.log.Errorw("error on execute Delete", "error", err)
		return err
	}

	o.log.Infow("attachment deleted successfully", "id", id)
	return nil
}
