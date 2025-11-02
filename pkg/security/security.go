package security

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	USER_KEY = "USER"
)

type UserAuth struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
}

func GetUser(ctx context.Context) (*UserAuth, error) {
	user, ok := ctx.Value(USER_KEY).(UserAuth)
	if !ok {
		return nil, fmt.Errorf("user ID not found in context")
	}

	return &user, nil
}

func CryptoPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}
