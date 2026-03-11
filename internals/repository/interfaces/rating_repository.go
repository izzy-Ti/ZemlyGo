package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type RatingRepository interface {
	Create(rating *domain.Rating) error
	GetByID(id uint) (*domain.Rating, error)
	GetByRideID(rideID uint) ([]domain.Rating, error)
	GetByToUserID(toUserID uint) ([]domain.Rating, error)
	Update(rating *domain.Rating) error
	Delete(id uint) error
}
