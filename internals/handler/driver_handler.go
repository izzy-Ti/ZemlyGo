package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type DriverHandler struct {
	driverService   *service.DriverService
	locationService *service.LocationService
}

func NewDriverHandler(driverService *service.DriverService, locationService *service.LocationService) *DriverHandler {
	return &DriverHandler{
		driverService:   driverService,
		locationService: locationService,
	}
}

func (h *DriverHandler) Register(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	var req dto.RegisterDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid request: "+err.Error()))
		return
	}

	driver, err := h.driverService.RegisterDriver(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(driver, "Driver registration submitted, pending approval"))
}

func (h *DriverHandler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorAPIResponse("User not authenticated"))
		return
	}
	userID := userIDVal.(uint)

	driver, err := h.driverService.GetDriverByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorAPIResponse("Driver profile not found"))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(driver, "Driver profile retrieved"))
}

func (h *DriverHandler) SetOnline(c *gin.Context) {
	driver := h.getDriverFromContext(c)
	if driver == nil {
		return
	}

	if err := h.driverService.SetOnline(driver.ID, true); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{"is_online": true}, "Driver is now online"))
}

func (h *DriverHandler) SetOffline(c *gin.Context) {
	driver := h.getDriverFromContext(c)
	if driver == nil {
		return
	}

	if err := h.driverService.SetOnline(driver.ID, false); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{"is_online": false}, "Driver is now offline"))
}

func (h *DriverHandler) UpdateLocation(c *gin.Context) {
	driver := h.getDriverFromContext(c)
	if driver == nil {
		return
	}

	var req dto.UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid coordinates: "+err.Error()))
		return
	}

	if err := h.locationService.UpdateLocation(driver.ID, req.Lat, req.Lng); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(gin.H{
		"lat": req.Lat,
		"lng": req.Lng,
	}, "Driver location updated"))
}

func (h *DriverHandler) GetEarnings(c *gin.Context) {
	driver := h.getDriverFromContext(c)
	if driver == nil {
		return
	}

	period := c.DefaultQuery("period", "all") // "today", "weekly", "all"

	earnings, err := h.driverService.GetEarnings(driver.ID, period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Failed to calculate earnings: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(earnings, "Driver earnings retrieved"))
}

func (h *DriverHandler) Approve(c *gin.Context) {
	idParam := c.Param("id")
	driverID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid driver id parameter"))
		return
	}

	if err := h.driverService.ApproveDriver(uint(driverID)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Driver approved successfully"))
}

func (h *DriverHandler) GetNearby(c *gin.Context) {
	lat, err1 := strconv.ParseFloat(c.Query("lat"), 64)
	lng, err2 := strconv.ParseFloat(c.Query("lng"), 64)
	radius, _ := strconv.ParseFloat(c.DefaultQuery("radius", "10"), 64)

	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("lat and lng query parameters are required"))
		return
	}

	drivers, err := h.driverService.GetNearbyDrivers(lat, lng, radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(drivers, "Nearby drivers retrieved"))
}

func (h *DriverHandler) getDriverFromContext(c *gin.Context) *dto.DriverDTO {
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
