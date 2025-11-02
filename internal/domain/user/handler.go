package user

import (
	"market/pkg/httpx"
	"market/pkg/security"
	"encoding/json"
	"net/http"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

// CreateUserHandler godoc
// @Summary      Registrar novo usuário
// @Description  Cria uma nova conta de usuário
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body		UserCreateDTO	true	"Dados do usuário"
// @Success      201		{object}	UserFoundDTO
// @Failure      400		{object}	map[string]string
// @Router       /auth/register [post]
func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerDTO *UserCreateDTO
	err := json.NewDecoder(r.Body).Decode(&registerDTO)

	if err != nil {
		httpx.SendBadRequest(w, "Invalid request body")
		return
	}

	userAuth, err := h.usecase.Register(registerDTO)
	if err != nil {
		httpx.SendBadRequest(w, "Failed to register user")
		return
	}

	json.NewEncoder(w).Encode(userAuth)
}

// LoginHandler godoc
// @Summary      Login de usuário
// @Description  Autentica um usuário no sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request	body		UserLoginDTO	true	"Dados de login"
// @Success      200		{object}	UserFoundDTO
// @Failure      400		{object}	map[string]string
// @Router       /auth/login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginDTO *UserLoginDTO
	err := json.NewDecoder(r.Body).Decode(&loginDTO)

	if err != nil {
		httpx.SendBadRequest(w, "Invalid request body")
		return
	}

	userAuth, err := h.usecase.Login(loginDTO)
	if err != nil {
		httpx.SendBadRequest(w, "Failed to login user")
		return
	}

	json.NewEncoder(w).Encode(userAuth)
}

// MeHandler godoc
// @Summary      Obter dados do usuário autenticado
// @Description  Retorna os dados do usuário atualmente autenticado
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200		{object}	UserFoundDTO
// @Failure      401		{object}	map[string]string
// @Router       /auth/me [get]
func (h *Handler) MeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userCtx, _ := security.GetUser(r.Context())

	userFound, err := h.usecase.Me(userCtx.UserID)
	if err != nil {
		httpx.SendBadRequest(w, "Failed to retrieve user")
		return
	}

	if err := json.NewEncoder(w).Encode(userFound); err != nil {
		httpx.SendInternalServerError(w, "Failed to encode user data", err)
		return
	}
}
