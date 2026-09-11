package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type VehicleRepository interface {
	Create(vehicle *domain.Vehicle) error
	GetByID(id uint) (*domain.Vehicle, error)
	GetByDriverID(driverID uint) ([]domain.Vehicle, error)
	Update(vehicle *domain.Vehicle) error
	Delete(id uint) error
	SetActiveVehicle(driverID uint, vehicleID uint) error
	GetActiveVehicle(driverID uint) (*domain.Vehicle, error)
	SetVehicleStatus(status bool, vehicleID uint) error
}
