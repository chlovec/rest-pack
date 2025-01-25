package main

import (
	"log"
	"net/http"

	"github.com/chlovec/rest-pack/api"
	"github.com/chlovec/rest-pack/example_rest_api/config"
	"github.com/chlovec/rest-pack/example_rest_api/services/product"
)

func main() {
	logger := log.Default()
	initServer(logger)
}

func initServer(logger *log.Logger) {
	// Initialize server
	addr := config.Envs.ServerAddress
	apiServer := api.NewAPIServer(addr, "/api/v1", logger)

	// Register routes
	handler := product.NewHandler(logger)
	apiServer.RegisterRoute("/products", handler.ListProducts, http.MethodGet)

	// start server
	err := apiServer.Start()
	if err != nil {
		log.Fatalf("error starting server:\n%v", err)
	}
}