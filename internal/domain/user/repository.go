package user

import (
	"database/sql"
	"market/pkg/database"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	FindByEmail(email string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
	Save(user *User) error
}

type userRepository struct {
	db             *database.PostgresDB
	log            *zap.SugaredLogger
	createStatment *sql.Stmt
}

func NewRepository(
	log *zap.SugaredLogger,
) Repository {

	dbInstance := database.GetInstance(log)

	insert := `INSERT INTO public.users
		(id, email, "password", "name", status, email_verified, last_login, created_at, updated_at)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);`

	createStatment, err := dbInstance.Prepare(insert)
	if err != nil {
		log.Errorw("error on create statment")
	}

	return &userRepository{
		db:             dbInstance,
		log:            log,
		createStatment: createStatment,
	}
}
func (u *userRepository) FindByEmail(email string) (*User, error) {
	sql := `SELECT id, email, "password", "name", status, email_verified, last_login, created_at, updated_at
	FROM users WHERE email = $1 LIMIT 1`
	row, err := u.db.Query(sql, email)

	if err != nil {
		u.log.Errorw("error on execute FindByEmail: %v", err)
		return nil, err
	}

	defer row.Close()

	var user User
	if row.Next() {
		err = row.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Status,
			&user.EmailVerified,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			u.log.Errorw("error scanning findByEmail %v", err)
			return nil, err
		}

		return &user, nil
	}
	return nil, nil
}

func (u *userRepository) FindByID(id uuid.UUID) (*User, error) {
	sql := `SELECT id, email, "password", "name", status, email_verified, last_login, created_at, updated_at
	FROM users WHERE id = $1 LIMIT 1`
	row, err := u.db.Query(sql, id)

	if err != nil {
		u.log.Errorw("error on execute FindByID: %v", err)
		return nil, err
	}

	defer row.Close()

	var user User
	if row.Next() {
		err = row.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.Status,
			&user.EmailVerified,
			&user.LastLogin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			u.log.Errorw("error scanning FindByID %v", err)
			return nil, err
		}

		return &user, nil
	}
	return nil, nil
}

func (u *userRepository) Save(user *User) error {
	_, err := u.createStatment.Exec(
		user.ID,
		user.Email,
		user.Password,
		user.Name,
		user.Status,
		user.EmailVerified,
		user.LastLogin,
	)

	if err != nil {
		u.log.Errorw("error on save new user")
		return err
	}

	return nil
}
