package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UserCreateDTO struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	market   string `json:"market"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserFoundDTO struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate verifica se os campos obrigatórios do UserCreateDTO estão corretamente preenchidos.
func (dto *UserCreateDTO) Validate() error {
	if strings.TrimSpace(dto.Email) == "" {
		return errors.New("email não pode ser vazio")
	}
	if strings.TrimSpace(dto.Name) == "" {
		return errors.New("nome não pode ser vazio")
	}
	if strings.TrimSpace(dto.Password) == "" {
		return errors.New("senha não pode ser vazia")
	}
	return nil
}

// Validate verifica se os campos obrigatórios do UserLoginDTO estão corretamente preenchidos.
func (dto *UserLoginDTO) Validate() error {
	if strings.TrimSpace(dto.Email) == "" {
		return errors.New("email não pode ser vazio")
	}
	if strings.TrimSpace(dto.Password) == "" {
		return errors.New("senha não pode ser vazia")
	}
	return nil
}
