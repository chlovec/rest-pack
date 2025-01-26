package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/chlovec/rest-pack/api"
	db "github.com/chlovec/rest-pack/db/mysql"
	"github.com/chlovec/rest-pack/examples/config"
	"github.com/chlovec/rest-pack/examples/services/product"
)

func main() {
	logger := log.Default()
	initServer(logger)
}

func initServer(logger *log.Logger) {
	// Initialize server
	addr := config.Envs.ServerAddress
	apiServer := api.NewAPIServer(addr, "/api/v1", logger)

	// Create db
	mysqlDB, err := db.InitDB(sql.Open, config.GetDataSourceName())
	if err != nil {
		logger.Fatalf("error initializing DB:\n%v", err)
	}

	// Create product store, product handler and register routes
	store := product.NewStore(mysqlDB)
	handler := product.NewHandler(logger, store)
	apiServer.RegisterRoute("/products", handler.ListProducts, http.MethodGet)
	apiServer.RegisterRoute("/products", handler.CreateProduct, http.MethodPost)

	// start server
	err = apiServer.Start()
	if err != nil {
		log.Fatalf("error starting server:\n%v", err)
	}
}