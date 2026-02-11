package route

import (
	"english-learning/configs"
	handler "english-learning/internal/modules/user/transport/http"
	"english-learning/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// Register registers all user routes on the given router.
func Register(r *gin.Engine, cfg *configs.Config, h *handler.UserHandler) {
	group := r.Group("/users")
	group.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.Get)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}
}
