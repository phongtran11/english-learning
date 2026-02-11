package route

import (
	handler "english-learning/internal/modules/auth/transport/http"

	"github.com/gin-gonic/gin"
)

// Register registers all auth routes on the given router.
func Register(r *gin.Engine, h *handler.AuthHandler) {
	group := r.Group("/auth")
	{
		group.POST("/register", h.Register)
		group.POST("/login", h.Login)
		group.POST("/refresh-token", h.RefreshToken)
		group.POST("/logout", h.Logout)
	}
}
