package dto

import "time"

type TripEarningItem struct {
	RideID   uint      `json:"ride_id"`
	Fare     float64   `json:"fare"`
	Tip      float64   `json:"tip"`
	Total    float64   `json:"total"`
	RideType string    `json:"ride_type"`
	Date     time.Time `json:"date"`
}

type DriverEarningsDTO struct {
	DriverID       uint              `json:"driver_id"`
	Period         string            `json:"period"` // "today", "weekly", "all"
	TotalEarnings  float64           `json:"total_earnings"`
	FaresTotal     float64           `json:"fares_total"`
	TipsTotal      float64           `json:"tips_total"`
	CompletedRides int               `json:"completed_rides"`
	Trips          []TripEarningItem `json:"trips"`
}

type AddTipRequest struct {
	Tip float64 `json:"tip" binding:"required,gt=0"`
}

type SurgeZone struct {
	ZoneName    string  `json:"zone_name"`
	CenterLat   float64 `json:"center_lat"`
	CenterLng   float64 `json:"center_lng"`
	Multiplier  float64 `json:"multiplier"`
	DemandLevel string  `json:"demand_level"` // "normal", "high", "surge"
}

type SurgeHeatmapResponse struct {
	Zones     []SurgeZone `json:"zones"`
	Timestamp time.Time   `json:"timestamp"`
}

type SystemStatsDTO struct {
	TotalRides     int64   `json:"total_rides"`
	ActiveRides    int64   `json:"active_rides"`
	CompletedRides int64   `json:"completed_rides"`
	OnlineDrivers  int64   `json:"online_drivers"`
	TotalRevenue   float64 `json:"total_revenue"`
	TotalUsers     int64   `json:"total_users"`
}
