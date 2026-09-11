package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type RideHandler struct {
	rideService    *service.RideService
	pricingService *service.PricingService
	driverService  *service.DriverService
}

func NewRideHandler(
	rideService *service.RideService,
	pricingService *service.PricingService,
	driverService *service.DriverService,
) *RideHandler {
	return &RideHandler{
		rideService:    rideService,
		pricingService: pricingService,
		driverService:  driverService,
	}
}

func (h *RideHandler) Estimate(c *gin.Context) {
	var req dto.EstimateFareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid coordinates: "+err.Error()))
		return
	}

	estimates := h.pricingService.EstimateAllOptions(
		req.PickupLat,
		req.PickupLng,
		req.DestinationLat,
		req.DestinationLng,
	)

	c.JSON(http.StatusOK, dto.SuccessResponse(estimates, "Fare estimates computed"))
}

func (h *RideHandler) RequestRide(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	riderID := userIDVal.(uint)

	var req dto.RequestRideRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride request: "+err.Error()))
		return
	}

	ride, err := h.rideService.RequestRide(riderID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusCreated, dto.SuccessResponse(enriched, "Ride requested successfully; searching for drivers"))
}

func (h *RideHandler) AcceptRide(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	ride, err := h.rideService.AcceptRide(driver.ID, rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(enriched, "Ride accepted"))
}

func (h *RideHandler) Arrive(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	ride, err := h.rideService.ArriveAtPickup(driver.ID, rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(enriched, "Arrived at pickup location"))
}

func (h *RideHandler) Start(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	ride, err := h.rideService.StartTrip(driver.ID, rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(enriched, "Trip in progress"))
}

func (h *RideHandler) Complete(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	ride, payment, err := h.rideService.CompleteTrip(driver.ID, rideID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"ride":    enriched,
		"payment": payment,
	}, "Trip completed successfully"))
}

func (h *RideHandler) Cancel(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID parameter"))
		return
	}

	var req dto.CancelRideRequest
	_ = c.ShouldBindJSON(&req)

	userRole, _ := c.Get("user_role")
	isDriver := userRole == "driver"

	// If driver, retrieve driver ID
	callerID := userID
	if isDriver {
		if drv, err := h.driverService.GetDriverByUserID(userID); err == nil && drv != nil {
			callerID = drv.ID
		}
	}

	if err := h.rideService.CancelRide(callerID, isDriver, rideID, req.Reason); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Ride cancelled"))
}

func (h *RideHandler) GetActive(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	userRole, _ := c.Get("user_role")
	isDriver := userRole == "driver"

	callerID := userID
	if isDriver {
		if drv, err := h.driverService.GetDriverByUserID(userID); err == nil && drv != nil {
			callerID = drv.ID
		}
	}

	ride, err := h.rideService.GetActiveRide(callerID, isDriver)
	if err != nil || ride == nil {
		c.JSON(http.StatusOK, dto.SuccessResponse(nil, "No active ride"))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(enriched, "Active ride found"))
}

func (h *RideHandler) GetByID(c *gin.Context) {
	rideID, err := h.parseRideID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid ride ID"))
		return
	}

	ride, err := h.rideService.GetRideByID(rideID)
	if err != nil || ride == nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse("Ride not found"))
		return
	}

	enriched := h.rideService.EnrichRideDetails(ride)
	c.JSON(http.StatusOK, dto.SuccessResponse(enriched, "Ride details retrieved"))
}

func (h *RideHandler) GetHistory(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	userRole, _ := c.Get("user_role")
	isDriver := userRole == "driver"

	callerID := userID
	if isDriver {
		if drv, err := h.driverService.GetDriverByUserID(userID); err == nil && drv != nil {
			callerID = drv.ID
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	rides, err := h.rideService.GetRideHistory(callerID, isDriver, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	var enrichedList []*dto.RideDTO
	for i := range rides {
		enrichedList = append(enrichedList, h.rideService.EnrichRideDetails(&rides[i]))
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(enrichedList, "Ride history retrieved"))
}

func (h *RideHandler) parseRideID(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	return uint(id), err
}

func (h *RideHandler) getDriver(c *gin.Context) *dto.DriverDTO {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return nil
	}
	userID := userIDVal.(uint)

	driver, err := h.driverService.GetDriverByUserID(userID)
	if err != nil {
		c.JSON(http.StatusForbidden, dto.ErrorAPIResponse("Caller is not registered as a driver"))
		return nil
	}
	return driver
}
