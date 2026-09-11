package domain

import "time"

type Rating struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	RideID     uint      `gorm:"column:ride_id;not null;index" json:"ride_id"`
	FromUserID uint      `gorm:"column:from_user_id;not null;index" json:"from_user_id"`
	ToUserID   uint      `gorm:"column:to_user_id;not null;index" json:"to_user_id"`
	Score      int       `gorm:"column:score;not null" json:"score"` // 1 - 5
	Comment    string    `gorm:"column:comment" json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}
