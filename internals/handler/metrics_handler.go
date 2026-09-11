package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
	"gorm.io/gorm"
)

type MetricsHandler struct {
	db             *gorm.DB
	pricingService *service.PricingService
}

func NewMetricsHandler(db *gorm.DB, pricingService *service.PricingService) *MetricsHandler {
	return &MetricsHandler{
		db:             db,
		pricingService: pricingService,
	}
}

func (h *MetricsHandler) GetStats(c *gin.Context) {
	var totalRides, activeRides, completedRides, onlineDrivers, totalUsers int64
	var totalRevenue float64

	h.db.Model(&domain.Ride{}).Count(&totalRides)
	h.db.Model(&domain.Ride{}).Where("status IN ?", constants.ActiveRideStatuses).Count(&activeRides)
	h.db.Model(&domain.Ride{}).Where("status = ?", constants.RideStatusCompleted).Count(&completedRides)
	h.db.Model(&domain.Drivers{}).Where("is_online = ? AND is_approved = ?", true, true).Count(&onlineDrivers)
	h.db.Model(&domain.Users{}).Count(&totalUsers)

	h.db.Model(&domain.Payment{}).Where("status = ?", "completed").Select("COALESCE(SUM(amount + tip), 0)").Scan(&totalRevenue)

	stats := dto.SystemStatsDTO{
		TotalRides:     totalRides,
		ActiveRides:    activeRides,
		CompletedRides: completedRides,
		OnlineDrivers:  onlineDrivers,
		TotalRevenue:   totalRevenue,
		TotalUsers:     totalUsers,
	}

	c.JSON(http.StatusOK, dto.SuccessResponse(stats, "System metrics retrieved"))
}

func (h *MetricsHandler) GetSurgeHeatmap(c *gin.Context) {
	heatmap := h.pricingService.GetSurgeHeatmap()
	c.JSON(http.StatusOK, dto.SuccessResponse(heatmap, "Surge pricing heatmap retrieved"))
}
