package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type DriverRepository interface {
	Create(driver *domain.Drivers) error
	GetByID(id uint) (*domain.Drivers, error)
	GetByUserID(userID uint) (*domain.Drivers, error)
	Update(driver *domain.Drivers) error
	Delete(id uint) error
	SetOnline(driverID uint, online bool) error
	UpdateLocation(driverID uint, lat, lng float64) error
	GetNearbyDrivers(lat, lng float64, radiusKm float64) ([]domain.Drivers, error)
	GetAvailableDrivers() ([]domain.Drivers, error)
	ApproveDriver(driverID uint) error
}
