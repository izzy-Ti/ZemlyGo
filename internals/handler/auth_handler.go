package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Sync(c *gin.Context) {
	var req dto.SyncUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid request: "+err.Error()))
		return
	}

	authID := req.AuthID
	if authID == "" {
		if authIDVal, exists := c.Get("auth_id"); exists {
			authID, _ = authIDVal.(string)
		}
	}
	if authID == "" {
		authID = "usr_" + req.Email
	}

	user, err := h.authService.SyncUser(authID, req.Name, req.Email, req.Phone, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	userDTO := dto.UserDTO{
		ID:             user.ID,
		AuthID:         user.AuthID,
		SupabaseUserID: user.AuthID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(userDTO, "User synchronized successfully with Neon Auth"))
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}

	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse("Local user account not found; please sync first"))
		return
	}

	profile, err := h.userService.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse("User profile not found"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(profile, "Profile retrieved"))
}

func (h *AuthHandler) DevToken(c *gin.Context) {
	var req dto.DevTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid request: "+err.Error()))
		return
	}

	res, err := h.authService.GenerateDevToken(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(res, "Dev token generated"))
}
