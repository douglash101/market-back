package main

import (
	"market/internal/domain/attachment"
	"market/internal/domain/product"
	"market/internal/domain/user"
	"market/internal/routes"
	"market/pkg/cloud"
	"market/pkg/config"
	"market/pkg/database"
	"market/pkg/logger"
	"net/http"

	_ "market/docs"
)

func main() {

	// Load environment configuration
	config.Load()

	// Initialize logger
	log := logger.NewLogger()

	// Initialize client database connection
	db := database.GetInstance(log)
	defer db.Close()

	cloud.NewCloudInstance(cloud.AWS_PROVIDER)

	// Initialize routes with handlers
	routeInstance := routes.NewRoutes(
		user.NewHandler(user.NewService(log)),
		product.NewHandler(product.NewService(log)),
		attachment.NewHandler(attachment.NewService(log)),
	)

	log.Infof("🙏 Starting server on port %s 🙏", config.Get().SERVER_PORT)
	err := http.ListenAndServe(
		config.Get().SERVER_PORT,
		routeInstance,
	)

	if err != nil {
		log.Errorf("❌ failed to start server: %v", err)
		panic(err)
	}
}
