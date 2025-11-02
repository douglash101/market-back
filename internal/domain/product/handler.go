package product

import (
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

func (h *Handler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// // Get user from context
	// userCtx, err := security.GetUser(r.Context())
	// if err != nil {
	// 	httpx.SendUnauthorized(w, "Authentication required")
	// 	return
	// }

	// Get product ID from URL path
	// idStr := r.PathValue("id")

	// productID, err := uuid.Parse(idStr)
	// if err != nil {
	// 	httpx.SendBadRequest(w, "Invalid product ID format")
	// 	return
	// }

	// product, err := h.usecase.GetProduct(productID)
	// if err != nil {
	// 	if err.Error() == "product not found" {
	// 		httpx.SendNotFound(w, "Product not found")
	// 		return
	// 	}
	// 	if err.Error() == "access denied" {
	// 		httpx.SendForbidden(w, "Access denied")
	// 		return
	// 	}
	// 	httpx.SendInternalServerError(w, err.Error(), err)
	// 	return
	// }

	// json.NewEncoder(w).Encode(product)
}
