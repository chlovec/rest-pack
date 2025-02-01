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
	// initialize config
	config.InitConfig()

	// start server
	logger := log.Default()
	apiServer := api.NewAPIServer(config.Envs.ServerAddress, config.Envs.PathPrefix, logger)
	err := run(apiServer, sql.Open, config.GetDataSourceName(), logger, 0)
	if err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}

func run(apiServer api.APIServerInterface, sqlOpen func(driverName, dataSourceName string) (*sql.DB, error), dsn string, logger *log.Logger, timeout time.Duration) error {
	// Create db
	mysqlDB, err := db.InitDB(sqlOpen, "mysql", dsn, timeout)
	if err != nil {
		return err
	}
	logger.Println("Initialized DB!")

	// Create server and register routes
	store := product.NewStore(mysqlDB)
	handler := product.NewHandler(logger, store)
	apiServer.RegisterRoute("/products", handler.ListProducts, http.MethodGet)
	apiServer.RegisterRoute("/products", handler.CreateProduct, http.MethodPost)

	// Start server
	err = apiServer.Start(timeout)
	if err != nil {
		return err
	}

	return nil
}
