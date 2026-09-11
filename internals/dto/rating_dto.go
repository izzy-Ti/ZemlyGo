package dto

import "time"

type CreateRatingRequest struct {
	Score   int    `json:"score" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type RatingDTO struct {
	ID         uint      `json:"id"`
	RideID     uint      `json:"ride_id"`
	FromUserID uint      `json:"from_user_id"`
	ToUserID   uint      `json:"to_user_id"`
	Score      int       `json:"score"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}

type DriverRatingSummary struct {
	DriverID     uint    `json:"driver_id"`
	AverageScore float64 `json:"average_score"`
	TotalRatings int     `json:"total_ratings"`
}
