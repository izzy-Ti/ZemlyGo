package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type RideRepository interface {
	Creat(ride *domain.Ride) error
	GetById(id uint) (*domain.Ride, error)
	Update(ride *domain.Ride) error
	UpdateStatus(rideId uint, status string) error
	AssignDriver(rideId uint, DriverId uint) error
	GetActiveRideByDriver(driverID uint) (*domain.Ride, error)
	GetActiveRideByRider(riderID uint) (*domain.Ride, error)
	GetRidesByRider(riderID uint, limit, offset int) ([]domain.Ride, error)
	GetRidesByDriver(driverID uint, limit, offset int) ([]domain.Ride, error)
	CancelRide(rideID uint) error
	CompleteRide(rideID uint, fare float64) error
}
