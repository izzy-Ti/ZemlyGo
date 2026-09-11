package service

import (
	"strconv"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type SafetyService struct {
	emergencyRepo interfaces.EmergencyRepository
	rideRepo      interfaces.RideRepository
	hub           *realtime.Hub
}

func NewSafetyService(
	emergencyRepo interfaces.EmergencyRepository,
	rideRepo interfaces.RideRepository,
	hub *realtime.Hub,
) *SafetyService {
	return &SafetyService{
		emergencyRepo: emergencyRepo,
		rideRepo:      rideRepo,
		hub:           hub,
	}
}

func (s *SafetyService) TriggerEmergencyAlert(rideID uint, userID uint, req dto.EmergencyAlertRequest) (*dto.EmergencyAlertDTO, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	alert := &domain.EmergencyAlert{
		RideID:    rideID,
		UserID:    userID,
		Lat:       req.Lat,
		Lng:       req.Lng,
		Reason:    req.Reason,
		Status:    "open",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.emergencyRepo.Create(alert); err != nil {
		return nil, err
	}

	// High-priority emergency broadcast
	alertEvent := realtime.Event{
		Type: "ride.emergency",
		Data: map[string]interface{}{
			"alert_id":   alert.ID,
			"ride_id":    rideID,
			"user_id":    userID,
			"lat":        req.Lat,
			"lng":        req.Lng,
			"reason":     req.Reason,
			"status":     "open",
			"created_at": alert.CreatedAt,
		},
	}

	// 1. Alert emergency / admin dispatchers
	s.hub.SendToRole(constants.RoleAdmin, alertEvent)

	// 2. Alert counterparty on ride
	if ride.RiderID == userID && ride.DriverID != nil {
		s.hub.SendToUser(strconv.FormatUint(uint64(*ride.DriverID), 10), alertEvent)
	} else if ride.DriverID != nil && *ride.DriverID == userID {
		s.hub.SendToUser(strconv.FormatUint(uint64(ride.RiderID), 10), alertEvent)
	}

	return &dto.EmergencyAlertDTO{
		ID:        alert.ID,
		RideID:    alert.RideID,
		UserID:    alert.UserID,
		Lat:       alert.Lat,
		Lng:       alert.Lng,
		Reason:    alert.Reason,
		Status:    alert.Status,
		CreatedAt: alert.CreatedAt,
	}, nil
}

func (s *SafetyService) ResolveAlert(alertID uint) error {
	return s.emergencyRepo.UpdateStatus(alertID, "resolved")
}

func (s *SafetyService) GetAlertsByRide(rideID uint) ([]dto.EmergencyAlertDTO, error) {
	alerts, err := s.emergencyRepo.GetByRideID(rideID)
	if err != nil {
		return nil, err
	}

	var results []dto.EmergencyAlertDTO
	for _, a := range alerts {
		results = append(results, dto.EmergencyAlertDTO{
			ID:        a.ID,
			RideID:    a.RideID,
			UserID:    a.UserID,
			Lat:       a.Lat,
			Lng:       a.Lng,
			Reason:    a.Reason,
			Status:    a.Status,
			CreatedAt: a.CreatedAt,
		})
	}
	return results, nil
}
