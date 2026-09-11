package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) interfaces.ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) Create(msg *domain.ChatMessage) error {
	return r.db.Create(msg).Error
}

func (r *ChatRepository) GetMessagesByRideID(rideID uint, limit, offset int) ([]domain.ChatMessage, error) {
	var msgs []domain.ChatMessage
	if limit <= 0 {
		limit = 50
	}
	err := r.db.Where("ride_id = ?", rideID).
		Order("created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
