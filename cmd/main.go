package main

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/kiloMIA/documed/internal/logger"
	"github.com/kiloMIA/documed/internal/repo"
	"github.com/kiloMIA/documed/internal/repo/postgre"
	"github.com/kiloMIA/documed/internal/service"
	"github.com/kiloMIA/documed/internal/transport/rest"
)

func main() {
	logger := logger.CreateLogger()
	defer logger.Sync()
	logger.Info("Logger Initialized")

	dbpool := postgre.CreateDB(logger)
	defer dbpool.Close()
	logger.Info("Database Initialized")

	repository := repo.NewRepository(dbpool, logger)
	logger.Info("Repository Initialized")

	sessionStore := sessions.NewCookieStore([]byte(os.Getenv("COOKIE_STORE")))
	authService := service.NewAuthService(repository.User, logger)
	logger.Info("Auth Service Initialized")

	serviceLayer := service.NewService(authService)
	logger.Info("Service Layer Initialized")

	restTransport := rest.NewTransport(serviceLayer, sessionStore, logger)
	logger.Info("Transport Layer Initialized")
}
