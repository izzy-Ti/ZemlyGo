package postgres

import (
	"sort"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DriverRepo struct {
	db *gorm.DB
}

func NewDriverRepository(db *gorm.DB) interfaces.DriverRepository {
	return &DriverRepo{db: db}
}

func (d *DriverRepo) Create(driver *domain.Drivers) error {
	return d.db.Create(driver).Error
}

func (d *DriverRepo) GetByID(id uint) (*domain.Drivers, error) {
	var driver domain.Drivers
	if err := d.db.Where("id = ?", id).First(&driver).Error; err != nil {
		return nil, err
	}
	return &driver, nil
}

func (d *DriverRepo) GetByUserID(userID uint) (*domain.Drivers, error) {
	var driver domain.Drivers
	if err := d.db.Where("user_id = ?", userID).First(&driver).Error; err != nil {
		return nil, err
	}
	return &driver, nil
}

func (d *DriverRepo) Update(driver *domain.Drivers) error {
	return d.db.Save(driver).Error
}

func (d *DriverRepo) Delete(id uint) error {
	return d.db.Where("id = ?", id).Delete(&domain.Drivers{}).Error
}

func (d *DriverRepo) SetOnline(driverID uint, online bool) error {
	return d.db.Model(&domain.Drivers{}).Where("id = ?", driverID).Update("is_online", online).Error
}

func (d *DriverRepo) UpdateLocation(driverID uint, lat, lng float64) error {
	loc := domain.DriverLocation{
		DriverID:  driverID,
		Lat:       lat,
		Lng:       lng,
		UpdatedAt: time.Now(),
	}

	// Upsert location record for driver
	return d.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "driver_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"lat", "lng", "updated_at"}),
	}).Create(&loc).Error
}

func (d *DriverRepo) GetAvailableDrivers() ([]domain.Drivers, error) {
	var drivers []domain.Drivers
	// Available = online AND approved
	err := d.db.Where("is_online = ? AND is_approved = ?", true, true).Find(&drivers).Error
	if err != nil {
		return nil, err
	}
	return drivers, nil
}

type driverWithDistance struct {
	driver   domain.Drivers
	distance float64
}

func (d *DriverRepo) GetNearbyDrivers(lat, lng float64, radiusKm float64) ([]domain.Drivers, error) {
	availableDrivers, err := d.GetAvailableDrivers()
	if err != nil {
		return nil, err
	}

	if len(availableDrivers) == 0 {
		return []domain.Drivers{}, nil
	}

	var driverIDs []uint
	driverMap := make(map[uint]domain.Drivers)
	for _, drv := range availableDrivers {
		driverIDs = append(driverIDs, drv.ID)
		driverMap[drv.ID] = drv
	}

	var locations []domain.DriverLocation
	if err := d.db.Where("driver_id IN ?", driverIDs).Find(&locations).Error; err != nil {
		return nil, err
	}

	var list []driverWithDistance
	for _, loc := range locations {
		dist := utils.Haversine(lat, lng, loc.Lat, loc.Lng)
		if dist <= radiusKm {
			if drv, ok := driverMap[loc.DriverID]; ok {
				list = append(list, driverWithDistance{
					driver:   drv,
					distance: dist,
				})
			}
		}
	}

	// Sort nearest first
	sort.Slice(list, func(i, j int) bool {
		return list[i].distance < list[j].distance
	})

	result := make([]domain.Drivers, len(list))
	for i, item := range list {
		result[i] = item.driver
	}

	return result, nil
}
