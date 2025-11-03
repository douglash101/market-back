package product_market

import (
	"encoding/json"
	"net/http"

	"market/pkg/httpx"
)

type Handler struct {
	usecase UseCase
}

func NewHandler(uc UseCase) *Handler {
	return &Handler{
		usecase: uc,
	}
}

func (h *Handler) CreateProductMarketHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto ProductMarketCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httpx.SendBadRequest(w, "Invalid JSON format")
		return
	}

	productMarket, err := h.usecase.CreateProductMarket(&dto)
	if err != nil {
		httpx.SendInternalServerError(w, err.Error(), err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(productMarket)
}

func (h *Handler) GetProductMarketsByProviderIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get provider ID from URL path
	providerID := r.PathValue("provider_id")
	if providerID == "" {
		httpx.SendBadRequest(w, "Provider ID is required")
		return
	}

	productMarkets, err := h.usecase.FindByProviderID(providerID)
	if err != nil {
		httpx.SendInternalServerError(w, err.Error(), err)
		return
	}

	json.NewEncoder(w).Encode(productMarkets)
}
