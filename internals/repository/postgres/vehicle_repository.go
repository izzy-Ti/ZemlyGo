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
	return v.db.Create(vehicle).Error
}

func (v *VehicleRepo) GetByID(id uint) (*domain.Vehicle, error) {
	var vec domain.Vehicle
	if err := v.db.Where("id = ?", id).First(&vec).Error; err != nil {
		return nil, err
	}
	return &vec, nil
}

func (v *VehicleRepo) GetByDriverID(driverID uint) ([]domain.Vehicle, error) {
	var list []domain.Vehicle
	if err := v.db.Where("driver_id = ?", driverID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (v *VehicleRepo) Update(vehicle *domain.Vehicle) error {
	return v.db.Save(vehicle).Error
}

func (v *VehicleRepo) Delete(id uint) error {
	return v.db.Where("id = ?", id).Delete(&domain.Vehicle{}).Error
}

func (v *VehicleRepo) SetActiveVehicle(driverID uint, vehicleID uint) error {
	return v.db.Transaction(func(tx *gorm.DB) error {
		// Deactivate all vehicles for this driver
		if err := tx.Model(&domain.Vehicle{}).Where("driver_id = ?", driverID).Update("status", false).Error; err != nil {
			return err
		}
		// Activate the selected vehicle
		if err := tx.Model(&domain.Vehicle{}).Where("id = ? AND driver_id = ?", vehicleID, driverID).Update("status", true).Error; err != nil {
			return err
		}
		return nil
	})
}

func (v *VehicleRepo) GetActiveVehicle(driverID uint) (*domain.Vehicle, error) {
	var vec domain.Vehicle
	if err := v.db.Where("driver_id = ? AND status = ?", driverID, true).First(&vec).Error; err != nil {
		return nil, err
	}
	return &vec, nil
}

func (v *VehicleRepo) SetVehicleStatus(status bool, vehicleID uint) error {
	return v.db.Model(&domain.Vehicle{}).Where("id = ?", vehicleID).Update("status", status).Error
}
