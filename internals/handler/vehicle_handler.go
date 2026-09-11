package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type VehicleHandler struct {
	vehicleService *service.VehicleService
	driverService  *service.DriverService
}

func NewVehicleHandler(vehicleService *service.VehicleService, driverService *service.DriverService) *VehicleHandler {
	return &VehicleHandler{
		vehicleService: vehicleService,
		driverService:  driverService,
	}
}

func (h *VehicleHandler) Create(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	var req dto.CreateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid vehicle details: "+err.Error()))
		return
	}

	vehicle, err := h.vehicleService.CreateVehicle(driver.ID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse(vehicle, "Vehicle registered successfully"))
}

func (h *VehicleHandler) List(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	vehicles, err := h.vehicleService.GetVehiclesByDriver(driver.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(vehicles, "Vehicles retrieved"))
}

func (h *VehicleHandler) SetActive(c *gin.Context) {
	driver := h.getDriver(c)
	if driver == nil {
		return
	}

	idParam := c.Param("id")
	vehicleID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid vehicle ID parameter"))
		return
	}

	if err := h.vehicleService.SetActiveVehicle(driver.ID, uint(vehicleID)); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(nil, "Vehicle set as active"))
}

func (h *VehicleHandler) getDriver(c *gin.Context) *dto.DriverDTO {
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
