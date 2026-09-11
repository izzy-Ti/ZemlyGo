package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type RatingHandler struct {
	ratingService *service.RatingService
}

func NewRatingHandler(ratingService *service.RatingService) *RatingHandler {
	return &RatingHandler{ratingService: ratingService}
}

func (h *RatingHandler) Rate(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	fromUserID := userIDVal.(uint)

	idParam := c.Param("id")
	rideID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	var req dto.CreateRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid rating score (must be 1-5): "+err.Error()))
		return
	}

	rating, err := h.ratingService.SubmitRating(uint(rideID), fromUserID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(rating, "Rating submitted successfully"))
}

func (h *RatingHandler) GetDriverRatings(c *gin.Context) {
	idParam := c.Param("id")
	driverID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid driver ID parameter"))
		return
	}

	summary, err := h.ratingService.GetDriverRatingSummary(uint(driverID))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(summary, "Driver rating summary retrieved"))
}
