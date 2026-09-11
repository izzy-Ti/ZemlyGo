package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type EmergencyRepository interface {
	Create(alert *domain.EmergencyAlert) error
	GetByID(id uint) (*domain.EmergencyAlert, error)
	GetByRideID(rideID uint) ([]domain.EmergencyAlert, error)
	UpdateStatus(id uint, status string) error
	ListOpenAlerts() ([]domain.EmergencyAlert, error)
}
