package dto

import "time"

type SendChatMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type ChatMessageDTO struct {
	ID         uint      `json:"id"`
	RideID     uint      `json:"ride_id"`
	SenderID   uint      `json:"sender_id"`
	SenderRole string    `json:"sender_role"`
	Message    string    `json:"message"`
	CreatedAt  time.Time `json:"created_at"`
}
