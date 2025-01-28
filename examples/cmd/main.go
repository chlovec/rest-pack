package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/chlovec/rest-pack/api"
	"github.com/chlovec/rest-pack/db"
	"github.com/chlovec/rest-pack/examples/config"
	"github.com/chlovec/rest-pack/examples/services/product"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// initServer(log.Default())
	logger := log.Default()
	apiServer, err := initServer(
		sql.Open, 
		config.Envs.ServerAddress, 
		config.GetDataSourceName(), 
		0,
		logger,
	)
	if err != nil {
		logger.Fatalf("error initializing db: %v", err)
	}
	startServer(apiServer, logger)
}

func initServer(
	sqlOpen func(driverName, dataSourceName string) (*sql.DB, error), serverAddress string, 
	dataSourceName string, 
	timeout time.Duration, 
	logger *log.Logger,
) (*api.APIServer, error) {
	apiServer := api.NewAPIServer(serverAddress, "/api/v1", logger)

	// Create db
	mysqlDB, err := db.InitDB(sqlOpen, "mysql", dataSourceName, timeout)
	if err != nil {
		return nil, err
	}
	logger.Println("Initialized DB!")

	// Create product store, product handler and register routes
	store := product.NewStore(mysqlDB)
	handler := product.NewHandler(logger, store)
	apiServer.RegisterRoute("/products", handler.ListProducts, http.MethodGet)
	apiServer.RegisterRoute("/products", handler.CreateProduct, http.MethodPost)

	return apiServer, nil
}

func startServer(apiServer *api.APIServer, logger *log.Logger) {
	err := apiServer.Start()
	if err != nil {
		logger.Fatalf("error starting server:\n%v", err)
	}
}
