package user

import (
	"github.com/google/uuid"
)

type UseCase interface {
	Register(input *UserCreateDTO) (*UserToken, error)
	Login(input *UserLoginDTO) (*UserToken, error)
	Me(id uuid.UUID) (UserFoundDTO, error)
}
