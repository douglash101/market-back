package user

import (
	"errors"
	"fmt"
	"market/pkg/security"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordMismatch = errors.New("password does not match")
	ErrUserNotFound     = errors.New("user not found")
)

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

func (s *service) Register(input *UserCreateDTO) (*UserToken, error) {

	err := input.Validate()
	if err != nil {
		return nil, fmt.Errorf("error", err)
	}

	userFound, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		s.log.Errorw(err.Error())
		return nil, err
	}

	if userFound != nil {
		return nil, fmt.Errorf("email %s already exists", input.Email)
	}

	password, err := security.CryptoPassword(input.Password)

	if err != nil {
		return nil, fmt.Errorf("internal error")
	}

	newUser := &User{
		ID:            uuid.New(),
		Name:          input.Name,
		Email:         input.Email,
		Password:      string(password),
		EmailVerified: false,
		LastLogin:     nil,
		Status:        UserStatusActive,
	}

	err = s.repository.Save(newUser)
	if err != nil {
		return nil, fmt.Errorf("error on save new user")
	}

	token, err := GenerateUserJWT(newUser)
	if err != nil {
		return nil, err
	}

	var userToken UserToken
	return userToken.NewUserToken(newUser, token), nil
}

func (s *service) Login(input *UserLoginDTO) (*UserToken, error) {

	userFound, err := s.repository.FindByEmail(input.Email)
	if err != nil {
		s.log.Errorw(err.Error())
		return nil, err
	}

	if userFound == nil {
		return nil, fmt.Errorf("user not exists")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrPasswordMismatch
		}
		return nil, err
	}

	token, err := GenerateUserJWT(userFound)
	if err != nil {
		return nil, err
	}

	var userToken UserToken
	return userToken.NewUserToken(userFound, token), nil
}

func (s *service) Me(id uuid.UUID) (UserFoundDTO, error) {
	found, err := s.repository.FindByID(id)
	if err != nil {
		return UserFoundDTO{}, fmt.Errorf("user not exists by id")
	}

	userFoundDTO := UserFoundDTO{
		ID:        found.ID,
		Name:      found.Name,
		Email:     found.Email,
		CreatedAt: found.CreatedAt,
		UpdatedAt: found.CreatedAt,
	}

	return userFoundDTO, nil
}

func GenerateUserJWT(u *User) (*Token, error) {

	expirationTime := time.Now().Add(8 * time.Hour) // Token expira em 8 horas

	claims := &Claims{
		Email:  u.Email,
		UserID: u.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Tempo de expiração
			IssuedAt:  time.Now().Unix(),     // Tempo de emissão
		},
	}

	// Cria o token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("sua_chave_secreta_aqui"))
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	expirationTimeSeconds := int64((8 * time.Hour).Seconds())

	return &Token{
		AccessToken: tokenString,
		ExpiresAt:   expirationTimeSeconds,
		TokenType:   "Bearer",
	}, nil
}
