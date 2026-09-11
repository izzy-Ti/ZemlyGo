package dto

import "time"

type CreatePaymentRequest struct {
	Method string `json:"method" binding:"required"` // card, cash, wallet
}

type PaymentDTO struct {
	ID        uint      `json:"id"`
	RideID    uint      `json:"ride_id"`
	Amount    float64   `json:"amount"`
	Method    string    `json:"method"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
