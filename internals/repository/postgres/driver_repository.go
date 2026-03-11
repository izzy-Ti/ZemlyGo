package postgres

import (
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
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
func (d *DriverRepo) GetByID(id uint) (*domain.Drivers, error) {
	var Rider domain.Drivers

	if err := d.db.Where("id = ?", id).First(&Rider).Error; err != nil {
		return nil, err
	}
	return &Rider, nil
}
func (d *DriverRepo) GetByUserID(userID uint) (*domain.Drivers, error) {
	var Rider domain.Drivers

	if err := d.db.Where("user_id = ?", userID).First(&Rider).Error; err != nil {
		return nil, err
	}
	return &Rider, nil
}
func (d *DriverRepo) Update(driver *domain.Drivers) error {
	if err := d.db.Save(&driver).Error; err != nil {
		return err
	}
	return nil
}
func (d *DriverRepo) Delete(id uint) error {
	if err := d.db.Where("id = ?", id).Delete(&domain.Drivers{}).Error; err != nil {
		return err
	}
	return nil
}
func (d *DriverRepo) SetOnline(driverID uint, online bool) error {
	var Driver domain.Drivers
	if err := d.db.Where("id = ?", driverID).Find(&Driver).Error; err != nil {
		return err
	}
	Driver.IsOnline = online

	if err := d.db.Save(Driver).Error; err != nil {
		return err
	}
	return nil

}
func (d *DriverRepo) UpdateLocation(driverID uint, lat, lng float64) error {
	var driver domain.DriverLocation
	if err := d.db.Where("driver_id = ?", driverID).First(&driver).Error; err != nil {
		return err
	}
	driver.Lat = lat
	driver.Lng = lng
	driver.UpdatedAt = time.Now()

	if err := d.db.Save(&driver).Error; err != nil {
		return err
	}
	return nil
}
func (d *DriverRepo) GetAvailableDrivers() ([]domain.Drivers, error) {
	var drivers []domain.Drivers

	err := d.db.Where("is_online = ?", true).Find(&drivers).Error
	if err != nil {
		return nil, err
	}
	return drivers, nil
}
func (d *DriverRepo) GetNearbyDrivers(lat, lng float64, radiusKm float64) ([]domain.Drivers, error) {
	drivers, err := d.GetAvailableDrivers()
	if err != nil {
		return nil, err
	}
	var DriverLoc []domain.DriverLocation

	for _, driver := range drivers {
		var location domain.DriverLocation
		if err := d.db.Where("id = ?", driver.ID).First(&location).Error; err != nil {
			return nil, err
		}
		DriverLoc = append(DriverLoc, location)
	}

	var nearBy []domain.DriverLocation
	for _, driver := range DriverLoc {
		distance := utils.Haversine(lat, lng, driver.Lat, driver.Lng)
		if distance <= radiusKm {
			nearBy = append(nearBy, driver)
		}
	}
	var nearByDriver []domain.Drivers
	for _, driver := range nearBy {
		err := d.db.Where("id = ?", driver.DriverID).Find(&nearByDriver).Error
		if err != nil {
			return nil, err
		}
	}
	return nearByDriver, nil
}
