package routes

import (
	"market/internal/domain/attachment"
	"market/internal/domain/product"
	"market/internal/domain/user"
	"market/pkg/middleware"
	"net/http"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Auth(handler http.HandlerFunc) http.HandlerFunc {
	return middleware.AuthMiddleware(handler)
}

func NewRoutes(
	userHandler *user.Handler,
	productHandler *product.Handler,
	attachmentHandler *attachment.Handler,
) http.Handler {
	mux := http.NewServeMux()

	// swagger route
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// auth routes
	mux.HandleFunc("POST /auth/login", userHandler.LoginHandler)
	mux.HandleFunc("POST /auth/register", userHandler.CreateUserHandler)
	mux.HandleFunc("GET /auth/me", Auth(userHandler.MeHandler))

	// product routes - clean REST endpoints
	mux.HandleFunc("GET /products/{id}", Auth(productHandler.GetProductHandler))

	// attachment routes
	mux.HandleFunc("POST /attachments", Auth(attachmentHandler.UploadAttachment))
	mux.HandleFunc("GET /attachments/{id}", Auth(attachmentHandler.GetAttachmentByID))
	mux.HandleFunc("PUT /attachments/{id}", Auth(attachmentHandler.UpdateAttachment))
	mux.HandleFunc("PATCH /attachments/{id}", Auth(attachmentHandler.UpdateAttachment))
	mux.HandleFunc("DELETE /attachments/{id}", Auth(attachmentHandler.DeleteAttachment))

	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	return corsConfig.Handler(mux)
}
