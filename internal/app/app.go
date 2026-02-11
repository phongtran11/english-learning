package app

import (
	"english-learning/configs"
	"english-learning/internal/server"
	"english-learning/pkg/logger"
	"fmt"

	driverpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// App holds the application-level dependencies and manages the lifecycle.
type App struct {
	cfg *configs.Config
	db  *gorm.DB
}

// New initializes the application: logger, database, and returns an App instance.
func New(cfg *configs.Config) (*App, error) {
	// Init Logger
	logger.InitLogger(cfg.Server.Env)

	// Init Database
	db, err := gorm.Open(driverpostgres.Open(cfg.Database.DSN), &gorm.Config{
		Logger: logger.NewGormLogger(logger.Log, 0),
	})
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return &App{
		cfg: cfg,
		db:  db,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	srv := server.New(a.cfg, a.db)

	logger.Infof("app", "Starting server on port %s", a.cfg.Server.Port)
	if err := srv.Run(":" + a.cfg.Server.Port); err != nil {
		return fmt.Errorf("starting server: %w", err)
	}

	return nil
}

// Close performs cleanup (e.g., closing DB connections).
func (a *App) Close() {
	if a.db != nil {
		sqlDB, err := a.db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	if logger.Log != nil {
		_ = logger.Log.Sync()
	}
}
