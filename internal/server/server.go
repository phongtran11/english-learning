package server

import (
	"english-learning/configs"
	authService "english-learning/internal/modules/auth/service"
	authHandler "english-learning/internal/modules/auth/transport/http"
	authRoute "english-learning/internal/modules/auth/transport/http/route"
	sessionPostgres "english-learning/internal/modules/session/repository/postgres"
	userPostgres "english-learning/internal/modules/user/repository/postgres"
	userService "english-learning/internal/modules/user/service"
	userHandler "english-learning/internal/modules/user/transport/http"
	userRoute "english-learning/internal/modules/user/transport/http/route"
	"english-learning/pkg/middleware"
	"english-learning/pkg/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// New creates and configures the Gin router with all routes and middleware.
func New(cfg *configs.Config, db *gorm.DB) *gin.Engine {
	r := gin.New()

	// Middleware
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())

	// Register Custom Validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validation.RegisterTagName(v)
		_ = v.RegisterValidation("date_format", validation.ValidateDateFormat);
		_ = v.RegisterValidation("phone", validation.ValidatePhone)
	}

	// Init Repositories
	userRepo := userPostgres.NewUserRepository(db)
	sessionRepo := sessionPostgres.NewSessionRepository(db)

	// Init Services
	userSvc := userService.NewService(userRepo)
	authSvc := authService.NewService(userRepo, sessionRepo, cfg.JWT.Secret)

	// Init Handlers
	userH := userHandler.NewUserHandler(userSvc)
	authH := authHandler.NewAuthHandler(authSvc)

	// Register Routes
	authRoute.Register(r, authH)
	userRoute.Register(r, cfg, userH)

	return r
}
