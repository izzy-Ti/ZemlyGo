package dto

import "time"

type EmergencyAlertRequest struct {
	Lat    float64 `json:"lat" binding:"required"`
	Lng    float64 `json:"lng" binding:"required"`
	Reason string  `json:"reason"`
}

type EmergencyAlertDTO struct {
	ID        uint      `json:"id"`
	RideID    uint      `json:"ride_id"`
	UserID    uint      `json:"user_id"`
	Lat       float64   `json:"lat"`
	Lng       float64   `json:"lng"`
	Reason    string    `json:"reason"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
