package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

type APIServerInterface interface {
	RegisterRoute(path string, handler func(http.ResponseWriter, *http.Request), methods ...string)
	Start(timeouts ...time.Duration) error
}

type APIServer struct {
	server       *http.Server
	apiRouter    *mux.Router
	apiSubrouter *mux.Router
	logger       *log.Logger
}

func NewAPIServer(addr string, pathPrefix string, logger *log.Logger) *APIServer {
	router := mux.NewRouter()
	subrouter := router

	if pathPrefix != "" {
		subrouter = router.PathPrefix(pathPrefix).Subrouter()
	}

	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	return &APIServer{
		apiRouter:    router,
		apiSubrouter: subrouter,
		server:       server,
		logger:       logger,
	}
}

func (s *APIServer) RegisterRoute(path string, handler func(http.ResponseWriter, *http.Request), methods ...string) {
	if path == "" {
		s.logger.Println("Cannot register a route with an empty path")
		return
	}
	if handler == nil {
		s.logger.Println("Cannot register a route with a nil handler")
		return
	}

	s.apiSubrouter.HandleFunc(path, handler).Methods(methods...)
	s.logger.Printf("Route registered: %s", path)
}

func (s *APIServer) Start(timeouts ...time.Duration) error {
	// Default timeout if none is provided
	timeout := 5 * time.Second
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}

	s.logger.Printf("Starting server on %s...", s.server.Addr)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	errChan := make(chan error, 1)

	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	select {
	case <-stop:
		s.logger.Println("Shutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		if err := s.server.Shutdown(ctx); err != nil {
			return err
		}
		s.logger.Println("Server stopped gracefully.")
		return nil
	case err := <-errChan:
		return err
	}
}
