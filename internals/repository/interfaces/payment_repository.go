package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type PaymentRepository interface {
	Create(payment *domain.Payment) error
	GetByID(id uint) (*domain.Payment, error)
	GetByRideID(rideID uint) (*domain.Payment, error)
	Update(payment *domain.Payment) error
	UpdateStatus(paymentID uint, status string) error
	Delete(id uint) error
}
