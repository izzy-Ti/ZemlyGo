package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type ChatRepository interface {
	Create(msg *domain.ChatMessage) error
	GetMessagesByRideID(rideID uint, limit, offset int) ([]domain.ChatMessage, error)
}
