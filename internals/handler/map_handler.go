package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

type MapHandler struct {
	mapService *service.MapService
}

func NewMapHandler(mapService *service.MapService) *MapHandler {
	return &MapHandler{mapService: mapService}
}

func (h *MapHandler) GetRoute(c *gin.Context) {
	var req dto.RouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Invalid coordinates: "+err.Error()))
		return
	}

	route, err := h.mapService.GetRoute(req.PickupLat, req.PickupLng, req.DestinationLat, req.DestinationLng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Failed to calculate route: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(route, "Route calculated successfully"))
}

func (h *MapHandler) ReverseGeocode(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	lat, err1 := strconv.ParseFloat(latStr, 64)
	lng, err2 := strconv.ParseFloat(lngStr, 64)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Valid 'lat' and 'lng' parameters are required"))
		return
	}

	result, err := h.mapService.ReverseGeocode(lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Reverse geocoding failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(result, "Address resolved"))
}

func (h *MapHandler) Geocode(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorAPIResponse("Query parameter 'q' is required"))
		return
	}

	results, err := h.mapService.Geocode(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Geocoding failed: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(results, "Locations found"))
}

func (h *MapHandler) SearchPlaces(c *gin.Context) {
	query := c.Query("q")
	suggestions, err := h.mapService.SearchPlaces(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorAPIResponse("Failed to search places: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(suggestions, "Places found"))
}
