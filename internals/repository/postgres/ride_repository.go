package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type RideRepository struct {
	db *gorm.DB
}

func NewRideRepository(db *gorm.DB) interfaces.RideRepository {
	return &RideRepository{db: db}
}

func (r *RideRepository) Create(ride *domain.Ride) error {
	return r.db.Create(ride).Error
}

func (r *RideRepository) GetByID(id uint) (*domain.Ride, error) {
	var ride domain.Ride
	if err := r.db.Where("id = ?", id).First(&ride).Error; err != nil {
		return nil, err
	}
	return &ride, nil
}

func (r *RideRepository) Update(ride *domain.Ride) error {
	return r.db.Save(ride).Error
}

func (r *RideRepository) UpdateStatus(rideID uint, status string) error {
	return r.db.Model(&domain.Ride{}).Where("id = ?", rideID).Update("status", status).Error
}

func (r *RideRepository) AssignDriver(rideID uint, driverID uint) error {
	return r.db.Model(&domain.Ride{}).
		Where("id = ?", rideID).
		Updates(map[string]interface{}{
			"driver_id": driverID,
			"status":    constants.RideStatusAccepted,
		}).Error
}

func (r *RideRepository) GetActiveRideByDriver(driverID uint) (*domain.Ride, error) {
	var ride domain.Ride
	if err := r.db.Where("driver_id = ? AND status IN ?", driverID, constants.ActiveRideStatuses).First(&ride).Error; err != nil {
		return nil, err
	}
	return &ride, nil
}

func (r *RideRepository) GetActiveRideByRider(riderID uint) (*domain.Ride, error) {
	var ride domain.Ride
	if err := r.db.Where("rider_id = ? AND status IN ?", riderID, constants.ActiveRideStatuses).First(&ride).Error; err != nil {
		return nil, err
	}
	return &ride, nil
}

func (r *RideRepository) GetRidesByRider(riderID uint, limit, offset int) ([]domain.Ride, error) {
	var rides []domain.Ride
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Where("rider_id = ?", riderID).Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

func (r *RideRepository) GetRidesByDriver(driverID uint, limit, offset int) ([]domain.Ride, error) {
	var rides []domain.Ride
	if err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Where("driver_id = ?", driverID).Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}

func (r *RideRepository) CancelRide(rideID uint) error {
	return r.db.Model(&domain.Ride{}).Where("id = ?", rideID).Update("status", constants.RideStatusCancelled).Error
}

func (r *RideRepository) CompleteRide(rideID uint, fare float64) error {
	return r.db.Model(&domain.Ride{}).
		Where("id = ?", rideID).
		Updates(map[string]interface{}{
			"status": constants.RideStatusCompleted,
			"fare":   fare,
		}).Error
}
