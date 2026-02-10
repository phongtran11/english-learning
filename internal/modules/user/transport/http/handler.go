package http

import (
	"english-learning/internal/modules/user/domain"
	"english-learning/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *services.Service
}

func NewUserHandler(service *services.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req RegisterRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	// Map DTO to Domain
	domainReq := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := h.service.Create(domainReq); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	response.Created(c, nil, response.MsgUserCreated)
}

func (h *UserHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, response.MsgInvalidID)
		return
	}

	user, err := h.service.Get(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, response.CodeNotFound, response.MsgUserNotFound)
		return
	}

	response.Success(c, ToUserResponse(user), response.MsgSuccess)
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, response.MsgInvalidID)
		return
	}

	var req UpdateUserRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	user := &domain.User{
		ID:          uint(id),
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Birthdate:   req.Birthdate,
	}

	if err := h.service.Update(user); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeServerInternalError, err.Error())
		return
	}

	response.Success(c, nil, response.MsgUserUpdated)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, response.MsgInvalidID)
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeServerInternalError, err.Error())
		return
	}

	response.Success(c, nil, response.MsgUserDeleted)
}

func (h *UserHandler) List(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	users, count, err := h.service.List(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, response.CodeServerInternalError, err.Error())
		return
	}

	resp := ToUserListResponse(users, count, page, pageSize)
	response.Success(c, resp)
}
