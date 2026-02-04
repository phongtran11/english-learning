package main

import (
	"english-learning/configs"
	handlers "english-learning/internal/adapters/handlers/http"
	"english-learning/internal/adapters/handlers/http/middleware"
	"english-learning/internal/adapters/repositories/postgres"
	"english-learning/internal/core/domain"
	"english-learning/internal/core/services/authservice"
	"english-learning/internal/core/services/userservice"
	"english-learning/pkg/logger"
	"english-learning/pkg/validation"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	driverpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 0. Load .env
	_ = godotenv.Load()

	// 1. Load Config
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Init Logger
	logger.InitLogger(cfg.Server.Env)
	defer logger.Log.Sync()

	// Register Custom Validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("date_format", validation.ValidateDateFormat)
	}

	logger.Infof("Loading config: %v", cfg)

	// 3. Init Database
	db, err := gorm.Open(driverpostgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return
	}

	// Auto Migrate
	if err := db.AutoMigrate(&domain.User{}, &domain.Session{}); err != nil {
		logger.Errorf("Failed to migrate database: %v", err)
		return
	}

	// 4. Init Repositories
	userRepo := postgres.NewUserRepository(db)
	sessionRepo := postgres.NewSessionRepository(db)

	// 5. Init Services
	userService := userservice.NewService(userRepo)
	authService := authservice.NewAuthService(userRepo, sessionRepo, cfg.JWT)

	// 6. Init Handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// 7. Init Router
	r := gin.Default()

	// Auth Routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.POST("/refresh", authHandler.RefreshToken)
		authGroup.POST("/logout", authHandler.Logout)
	}

	// User Routes (Protected)
	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		userGroup.POST("", userHandler.Create)
		userGroup.GET("", userHandler.List)
		userGroup.GET("/:id", userHandler.Get)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
	}

	// 8. Start Server
	logger.Infof("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
