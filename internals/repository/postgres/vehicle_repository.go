package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type VehicleRepo struct {
	db *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) interfaces.VehicleRepository {
	return &VehicleRepo{db: db}
}

func (v *VehicleRepo) Create(vehicle *domain.Vehicle) error {
	if err := v.db.Create(vehicle).Error; err != nil {
		return err
	}
	return nil
}
func (v *VehicleRepo) GetByID(id uint) (*domain.Vehicle, error) {
	var vec *domain.Vehicle
	if err := v.db.Where("id = ?", id).First(&vec).Error; err != nil {
		return nil, err
	}
	return vec, nil
}
func (v *VehicleRepo) GetByDriverID(driverID uint) ([]domain.Vehicle, error) {
	var vec []domain.Vehicle
	if err := v.db.Where("driver_id = ?", driverID).Find(&vec).Error; err != nil {
		return nil, err
	}
	return vec, nil
}
func (v *VehicleRepo) Update(vehicle *domain.Vehicle) error {
	if err := v.db.Save(vehicle).Error; err != nil {
		return err
	}
	return nil
}
func (v *VehicleRepo) Delete(id uint) error {
	if err := v.db.Where("id = ?", id).Delete(&domain.DriverLocation{}).Error; err != nil {
		return err
	}
	return nil
}
func (v *VehicleRepo) SetActiveVehicle(driverID uint, vehicleID uint) error {
	var vec domain.Vehicle
	if err := v.db.Where("id = ?", vehicleID).First(vec).Error; err != nil {
		return err
	}
	vec.DriverID = driverID

	if err := v.db.Save(vec).Error; err != nil {
		return err
	}
	return nil
}
func (v *VehicleRepo) GetActiveVehicle(driverID uint) (*domain.Vehicle, error) {
	var vec *domain.Vehicle
	if err := v.db.Where("driver_id = ? AND status = ?", driverID, true).Find(&vec).Error; err != nil {
		return nil, err
	}
	return vec, nil
}
func (v *VehicleRepo) SetVehicleStatus(Status bool, vehicleID uint) error {
	var vec domain.Vehicle
	if err := v.db.Where("id = ?", vehicleID).First(vec).Error; err != nil {
		return err
	}
	vec.Status = Status

	if err := v.db.Save(vec).Error; err != nil {
		return err
	}
	return nil
}
