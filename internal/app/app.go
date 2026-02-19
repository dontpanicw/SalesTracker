package app

import (
	"fmt"
	"log"

	"github.com/yourusername/analytics-service/config"
	"github.com/yourusername/analytics-service/internal/adapter/repository/postgres"
	httpServer "github.com/yourusername/analytics-service/internal/input/http"
	"github.com/yourusername/analytics-service/internal/usecases"
	"github.com/yourusername/analytics-service/pkg/migrations"
)

type App struct {
	config *config.Config
}

func New() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &App{config: cfg}, nil
}

func (a *App) Run() error {
	// Запуск миграций
	if err := migrations.Run(a.config.DatabaseDSN); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Инициализация репозитория
	repo, err := postgres.New(a.config.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	// Инициализация use cases
	uc := usecases.New(repo)

	// Запуск HTTP сервера
	server := httpServer.NewServer(uc, a.config.ServerPort)
	
	log.Printf("Starting server on port %s", a.config.ServerPort)
	return server.Start()
}
