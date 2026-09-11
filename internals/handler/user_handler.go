package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	profile, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse("User not found"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(profile, "Profile retrieved"))
}

func (h *UserHandler) UpdateMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid request body: "+err.Error()))
		return
	}

	updated, err := h.userService.UpdateProfile(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(updated, "Profile updated"))
}

func (h *UserHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.userService.ListUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(users, "Users list retrieved"))
}
