package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type PaymentHandler struct {
	paymentService *service.PaymentService
}

func NewPaymentHandler(paymentService *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

func (h *PaymentHandler) Pay(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	riderID := userIDVal.(uint)

	idParam := c.Param("id")
	rideID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	var req dto.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Payment method is required"))
		return
	}

	payment, err := h.paymentService.ProcessPayment(uint(rideID), riderID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(payment, "Payment processed successfully"))
}

func (h *PaymentHandler) AddTip(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	riderID := userIDVal.(uint)

	idParam := c.Param("id")
	rideID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	var req dto.AddTipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Valid positive tip amount is required"))
		return
	}

	payment, err := h.paymentService.AddTip(uint(rideID), riderID, req.Tip)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(payment, "Tip added to driver payment"))
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	idParam := c.Param("id")
	rideID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	payment, err := h.paymentService.GetPaymentByRideID(uint(rideID), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(payment, "Payment details retrieved"))
}
