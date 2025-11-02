package user

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
	UserStatusDeleted  UserStatus = "deleted"
)

type User struct {
	ID            uuid.UUID  `json:"id"`
	Name          string     `json:"name"`
	Email         string     `json:"email"`
	Password      string     `json:"password"`
	EmailVerified bool       `json:"email_verified"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	Status        UserStatus `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
type UserToken struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresAt   int64  `json:"expires_at"`
}

func (t *UserToken) NewUserToken(user *User, token *Token) *UserToken {
	t.Email = user.Email
	t.Name = user.Name
	t.AccessToken = token.AccessToken
	t.TokenType = token.TokenType
	t.ExpiresAt = token.ExpiresAt

	return t
}

// Estrutura para as claims do JWT
type Claims struct {
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	UserID    uuid.UUID `json:"user_id"`
	CompanyID uuid.UUID `json:"company_id"`
	jwt.StandardClaims
}
