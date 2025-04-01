package main

import (
	"log"
	"os"

	"github.com/Hexes-rgb/employee-service/internal/config"
	"github.com/Hexes-rgb/employee-service/internal/repository/postgres"
	"github.com/Hexes-rgb/employee-service/internal/server"
	"github.com/Hexes-rgb/employee-service/internal/service"
	"github.com/Hexes-rgb/employee-service/internal/transport/rest"
)

func main() {
	logger := log.New(os.Stdout, "EMPLOYEE-SERVICE: ", log.LstdFlags|log.Lshortfile)

	cfg := config.Load()

	db, err := config.InitDB(cfg.Database, logger)
	if err != nil {
		logger.Fatalf("Database initialization failed: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Printf("Error closing database connection: %v", err)
		}
	}()

	empRepo := postgres.NewEmployeeRepo(db)
	deptRepo := postgres.NewDepartmentRepo(db)

	empService := service.NewEmployeeService(empRepo, deptRepo)
	deptService := service.NewDepartmentService(deptRepo)

	router := rest.NewRouter(empService, deptService)

	srv := server.New(cfg.Server, router, logger)
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}

	srv.WaitForShutdown()
}
