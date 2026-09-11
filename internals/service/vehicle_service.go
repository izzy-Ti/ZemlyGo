package service

import (
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type VehicleService struct {
	vehicleRepo interfaces.VehicleRepository
	driverRepo  interfaces.DriverRepository
}

func NewVehicleService(vehicleRepo interfaces.VehicleRepository, driverRepo interfaces.DriverRepository) *VehicleService {
	return &VehicleService{
		vehicleRepo: vehicleRepo,
		driverRepo:  driverRepo,
	}
}

func (v *VehicleService) CreateVehicle(driverID uint, req dto.CreateVehicleRequest) (*dto.VehicleDTO, error) {
	// If it's the first vehicle, set it active by default; otherwise false
	existing, _ := v.vehicleRepo.GetByDriverID(driverID)
	status := len(existing) == 0

	vehicle := &domain.Vehicle{
		DriverID:    driverID,
		PlateNumber: req.PlateNumber,
		Model:       req.Model,
		Color:       req.Color,
		Status:      status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := v.vehicleRepo.Create(vehicle); err != nil {
		return nil, err
	}

	return &dto.VehicleDTO{
		ID:          vehicle.ID,
		DriverID:    vehicle.DriverID,
		PlateNumber: vehicle.PlateNumber,
		Model:       vehicle.Model,
		Color:       vehicle.Color,
		Status:      vehicle.Status,
		CreatedAt:   vehicle.CreatedAt,
	}, nil
}

func (v *VehicleService) GetVehiclesByDriver(driverID uint) ([]dto.VehicleDTO, error) {
	vehicles, err := v.vehicleRepo.GetByDriverID(driverID)
	if err != nil {
		return nil, err
	}

	var results []dto.VehicleDTO
	for _, vec := range vehicles {
		results = append(results, dto.VehicleDTO{
			ID:          vec.ID,
			DriverID:    vec.DriverID,
			PlateNumber: vec.PlateNumber,
			Model:       vec.Model,
			Color:       vec.Color,
			Status:      vec.Status,
			CreatedAt:   vec.CreatedAt,
		})
	}
	return results, nil
}

func (v *VehicleService) SetActiveVehicle(driverID uint, vehicleID uint) error {
	vec, err := v.vehicleRepo.GetByID(vehicleID)
	if err != nil || vec == nil {
		return constants.ErrNotFound
	}
	if vec.DriverID != driverID {
		return constants.ErrForbidden
	}

	return v.vehicleRepo.SetActiveVehicle(driverID, vehicleID)
}

func (v *VehicleService) GetActiveVehicle(driverID uint) (*dto.VehicleDTO, error) {
	vec, err := v.vehicleRepo.GetActiveVehicle(driverID)
	if err != nil || vec == nil {
		return nil, constants.ErrNoActiveVehicle
	}

	return &dto.VehicleDTO{
		ID:          vec.ID,
		DriverID:    vec.DriverID,
		PlateNumber: vec.PlateNumber,
		Model:       vec.Model,
		Color:       vec.Color,
		Status:      vec.Status,
		CreatedAt:   vec.CreatedAt,
	}, nil
}
