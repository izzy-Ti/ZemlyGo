package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type SafetyHandler struct {
	safetyService *service.SafetyService
}

func NewSafetyHandler(safetyService *service.SafetyService) *SafetyHandler {
	return &SafetyHandler{safetyService: safetyService}
}

func (h *SafetyHandler) TriggerEmergency(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	rideIDParam := c.Param("id")
	rideID, err := strconv.ParseUint(rideIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	var req dto.EmergencyAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Coordinates (lat, lng) are required: "+err.Error()))
		return
	}

	alert, err := h.safetyService.TriggerEmergencyAlert(uint(rideID), userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(alert, "EMERGENCY SOS BROADCASTED TO DISPATCH AND COUNTERPARTY"))
}

func (h *SafetyHandler) GetAlerts(c *gin.Context) {
	rideIDParam := c.Param("id")
	rideID, err := strconv.ParseUint(rideIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	alerts, err := h.safetyService.GetAlertsByRide(uint(rideID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(alerts, "Emergency alerts retrieved"))
}

func (h *SafetyHandler) ResolveAlert(c *gin.Context) {
	alertIDParam := c.Param("id")
	alertID, err := strconv.ParseUint(alertIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid alert ID parameter"))
		return
	}

	if err := h.safetyService.ResolveAlert(uint(alertID)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Emergency alert resolved"))
}
