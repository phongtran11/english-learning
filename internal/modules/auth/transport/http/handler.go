package http

import (
	authDomain "english-learning/internal/modules/auth/domain"
	"english-learning/pkg/response"
	"english-learning/pkg/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles HTTP requests for authentication operations.
type AuthHandler struct {
	service authDomain.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given service interface.
func NewAuthHandler(service authDomain.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, validation.FormatError(err))
		return
	}

	domainReq := &authDomain.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Register(domainReq); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, nil, response.MsgUserRegistered)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, validation.FormatError(err))
		return
	}

	domainReq := &authDomain.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	userAgent := c.Request.UserAgent()
	clientIP := c.ClientIP()

	tokenPair, err := h.service.Login(domainReq, clientIP, userAgent)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, "Invalid credentials")
		return
	}

	resp := TokenPairResponseDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	response.Success(c, resp, response.MsgLoginSuccess)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, validation.FormatError(err))
		return
	}

	tokenPair, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, response.CodeUnauthorized, err.Error())
		return
	}

	resp := TokenPairResponseDTO{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}

	response.Success(c, resp, response.MsgRefreshTokenSuccess)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req RefreshTokenRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, validation.FormatError(err))
		return
	}

	if err := h.service.Logout(req.RefreshToken); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeServerInternalError, err.Error())
		return
	}

	response.Success(c, nil, response.MsgSuccess)
}
