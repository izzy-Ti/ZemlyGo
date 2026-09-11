package dto

import "time"

type CreateVehicleRequest struct {
	PlateNumber string `json:"plate_number" binding:"required"`
	Model       string `json:"model" binding:"required"`
	Color       string `json:"color" binding:"required"`
}

type VehicleDTO struct {
	ID          uint      `json:"id"`
	DriverID    uint      `json:"driver_id"`
	PlateNumber string    `json:"plate_number"`
	Model       string    `json:"model"`
	Color       string    `json:"color"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
