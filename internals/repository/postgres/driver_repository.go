package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type DriverRepo struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) interfaces.DriverRepository {
	return &DriverRepo{db: db}
}

func (d *DriverRepo) Create(driver *domain.Drivers) error {
	if err := d.db.Create(driver).Error; err != nil {
		return err
	}
	return nil
}
func (d *DriverRepo) GetByID(id uint) (*domain.Drivers, error)
func (d *DriverRepo) GetByUserID(userID uint) (*domain.Drivers, error)
func (d *DriverRepo) Update(driver *domain.Drivers) error
func (d *DriverRepo) Delete(id uint) error
func (d *DriverRepo) SetOnline(driverID uint, online bool) error
func (d *DriverRepo) UpdateLocation(driverID uint, lat, lng float64) error
func (d *DriverRepo) GetNearbyDrivers(lat, lng float64, radiusKm float64) ([]domain.Drivers, error)
func (d *DriverRepo) GetAvailableDrivers() ([]domain.Drivers, error)
func (d *DriverRepo) ApproveDriver(driverID uint) error
