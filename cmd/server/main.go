package main

import (
	"english-learning/configs"
	authHandler "english-learning/internal/modules/auth/transport/http"
	userHandler "english-learning/internal/modules/user/transport/http"
	"english-learning/pkg/middleware"

	authService "english-learning/internal/modules/auth/service"
	userService "english-learning/internal/modules/user/service"

	sessionPostgres "english-learning/internal/modules/session/repository/postgres"
	userPostgres "english-learning/internal/modules/user/repository/postgres"

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

	// 3. Init Database
	db, err := gorm.Open(driverpostgres.Open(cfg.Database.DSN), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return
	}

	// 4. Init Repositories
	userRepo := userPostgres.NewUserRepository(db)
	sessionRepo := sessionPostgres.NewSessionRepository(db)

	// 5. Init Services
	userServiceInstance := userService.NewService(userRepo)
	authServiceInstance := authService.NewService(userRepo, sessionRepo, cfg.JWT.Secret)

	// 6. Init Handlers
	userH := userHandler.NewUserHandler(userServiceInstance)
	authH := authHandler.NewAuthHandler(authServiceInstance)

	// 7. Init Router
	r := gin.Default()

	// Auth Routes
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authH.Register)
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.RefreshToken)
		authGroup.POST("/logout", authH.Logout)
	}

	// User Routes (Protected)
	userGroup := r.Group("/users")
	userGroup.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		userGroup.POST("", userH.Create)
		userGroup.GET("", userH.List)
		userGroup.GET("/:id", userH.Get)
		userGroup.PUT("/:id", userH.Update)
		userGroup.DELETE("/:id", userH.Delete)
	}

	// 8. Start Server
	logger.Infof("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Errorf("Failed to start server: %v", err)
	}
}
