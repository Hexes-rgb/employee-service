package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Hexes-rgb/employee-service/internal/config"
)

type Server struct {
	httpServer *http.Server
	logger     *log.Logger
}

func New(cfg config.ServerConfig, handler http.Handler, logger *log.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         ":" + cfg.Port,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			ErrorLog:     logger,
		},
		logger: logger,
	}
}

func (s *Server) Run() error {
	go func() {
		s.logger.Printf("Server starting on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Server error: %v", err)
		}
	}()

	return nil
}

func (s *Server) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	s.logger.Println("Shutdown signal received")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.logger.Printf("Server shutdown error: %v", err)
	}

	s.logger.Println("Server stopped gracefully")
}
