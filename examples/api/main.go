package main

import (
	"bytes"
	"log"
	"net/http"

	"github.com/chlovec/rest-pack/api"
)

func main() {
	// Set up a logger
	var logBuffer bytes.Buffer
	logger := log.New(&logBuffer, "", log.LstdFlags)

	// setup and start server
	apiServer := api.NewAPIServer(":8080", "/api/v1", logger)
	apiServer.RegisterRoute("/chats", getHandler, http.MethodGet)
	apiServer.Start()
}

func getHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger) {
	logger.Println("Handling get request")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Get route hit"))
}