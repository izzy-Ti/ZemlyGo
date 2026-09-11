package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type ChatService struct {
	chatRepo   interfaces.ChatRepository
	rideRepo   interfaces.RideRepository
	driverRepo interfaces.DriverRepository
	hub        *realtime.Hub
}

func NewChatService(
	chatRepo interfaces.ChatRepository,
	rideRepo interfaces.RideRepository,
	driverRepo interfaces.DriverRepository,
	hub *realtime.Hub,
) *ChatService {
	return &ChatService{
		chatRepo:   chatRepo,
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
		hub:        hub,
	}
}

func (s *ChatService) SendMessage(rideID uint, userID uint, isDriver bool, text string) (*dto.ChatMessageDTO, error) {
	if text == "" {
		return nil, errors.New("message cannot be empty")
	}

	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	senderRole := "rider"
	var recipientUserID uint

	if isDriver {
		driver, err := s.driverRepo.GetByUserID(userID)
		if err != nil || driver == nil || ride.DriverID == nil || *ride.DriverID != driver.ID {
			return nil, constants.ErrForbidden
		}
		senderRole = "driver"
		recipientUserID = ride.RiderID
	} else {
		if ride.RiderID != userID {
			return nil, constants.ErrForbidden
		}
		if ride.DriverID != nil {
			driver, _ := s.driverRepo.GetByID(*ride.DriverID)
			if driver != nil {
				recipientUserID = driver.UserId
			}
		}
	}

	msg := &domain.ChatMessage{
		RideID:     rideID,
		SenderID:   userID,
		SenderRole: senderRole,
		Message:    text,
		CreatedAt:  time.Now(),
	}

	if err := s.chatRepo.Create(msg); err != nil {
		return nil, err
	}

	// Dispatch real-time WebSocket event to recipient
	if recipientUserID > 0 {
		recipientStr := strconv.FormatUint(uint64(recipientUserID), 10)
		s.hub.SendToUser(recipientStr, realtime.Event{
			Type: "ride.chat_message",
			Data: map[string]interface{}{
				"id":          msg.ID,
				"ride_id":     rideID,
				"sender_id":   userID,
				"sender_role": senderRole,
				"message":     text,
				"created_at":  msg.CreatedAt,
			},
		})
	}

	return &dto.ChatMessageDTO{
		ID:         msg.ID,
		RideID:     msg.RideID,
		SenderID:   msg.SenderID,
		SenderRole: msg.SenderRole,
		Message:    msg.Message,
		CreatedAt:  msg.CreatedAt,
	}, nil
}

func (s *ChatService) GetMessages(rideID uint, userID uint, isDriver bool, limit, offset int) ([]dto.ChatMessageDTO, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	// Verify caller participation
	if isDriver {
		driver, err := s.driverRepo.GetByUserID(userID)
		if err != nil || driver == nil || ride.DriverID == nil || *ride.DriverID != driver.ID {
			return nil, constants.ErrForbidden
		}
	} else {
		if ride.RiderID != userID {
			return nil, constants.ErrForbidden
		}
	}

	msgs, err := s.chatRepo.GetMessagesByRideID(rideID, limit, offset)
	if err != nil {
		return nil, err
	}

	var results []dto.ChatMessageDTO
	for _, m := range msgs {
		results = append(results, dto.ChatMessageDTO{
			ID:         m.ID,
			RideID:     m.RideID,
			SenderID:   m.SenderID,
			SenderRole: m.SenderRole,
			Message:    m.Message,
			CreatedAt:  m.CreatedAt,
		})
	}
	return results, nil
}
