package http

import (
	"english-learning/internal/core/domain"
	"english-learning/internal/core/services/userservice"
	"english-learning/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *userservice.Service
}

func NewUserHandler(service *userservice.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}

	if err := h.service.Create(&req); err != nil {
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

	response.Success(c, user, response.MsgSuccess)
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, response.MsgInvalidID)
		return
	}

	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		response.Error(c, http.StatusBadRequest, response.CodeBadRequest, err.Error())
		return
	}
	user.ID = uint(id)

	if err := h.service.Update(&user); err != nil {
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

	response.SuccessList(c, users, count, page, pageSize, response.MsgSuccess)
}
