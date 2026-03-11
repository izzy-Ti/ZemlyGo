package postgres

import (
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

func (r *RideRepository) Creat(ride *domain.Ride) error {
	if err := r.db.Save(ride).Error; err != nil {
		return err
	}
	return nil
}
func (r *RideRepository) GetById(id uint) (*domain.Ride, error) {
	var ride *domain.Ride
	if err := r.db.Where("id = ?", id).First(&ride).Error; err != nil {
		return nil, err
	}
	return ride, nil
}
func (r *RideRepository) Update(ride *domain.Ride) error {
	if err := r.db.Save(ride); err != nil {
		return nil
	}
	return nil
}
func (r *RideRepository) UpdateStatus(rideId uint, status string) error {
	var ride domain.Ride
	if err := r.db.Where("id = ?", rideId).First(ride).Error; err != nil {
		return err
	}
	ride.Status = status
	if err := r.db.Save(ride).Error; err != nil {
		return err
	}
	return nil
}
func (r *RideRepository) AssignDriver(rideId uint, DriverId uint) error {
	var ride domain.Ride
	if err := r.db.Where("id = ?", rideId).First(ride).Error; err != nil {
		return err
	}
	ride.DriverID = &DriverId
	if err := r.db.Save(ride).Error; err != nil {
		return err
	}
	return nil
}
func (r *RideRepository) GetActiveRideByDriver(driverID uint) (*domain.Ride, error) {
	var ride *domain.Ride
	if err := r.db.Where("driver_id = ? AND status = ?", driverID, "ACTIVE").First(&ride).Error; err != nil {
		return nil, err
	}
	return ride, nil
}
func (r *RideRepository) GetActiveRideByRider(riderID uint) (*domain.Ride, error) {
	var ride *domain.Ride
	if err := r.db.Where("rider_id = ? AND status = ?", riderID, "ACTIVE").First(&ride).Error; err != nil {
		return nil, err
	}
	return ride, nil
}
func (r *RideRepository) GetRidesByRider(riderID uint, limit, offset int) ([]domain.Ride, error) {
	var rides []domain.Ride
	if err := r.db.Limit(limit).Offset(offset).Where("rider_id = ?", riderID).Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}
func (r *RideRepository) GetRidesByDriver(driverID uint, limit, offset int) ([]domain.Ride, error) {
	var rides []domain.Ride
	if err := r.db.Limit(limit).Offset(offset).Where("driver_id = ?", driverID).Find(&rides).Error; err != nil {
		return nil, err
	}
	return rides, nil
}
func (r *RideRepository) CancelRide(rideID uint) error {
	var ride domain.Ride
	if err := r.db.Where("id = ?", rideID).First(ride).Error; err != nil {
		return err
	}
	ride.Status = "CANCLED"
	if err := r.db.Save(ride).Error; err != nil {
		return err
	}
	return nil
}
func (r *RideRepository) CompleteRide(rideID uint, fare float64) error {
	var ride domain.Ride
	if err := r.db.Where("id = ?", rideID).First(ride).Error; err != nil {
		return err
	}
	ride.Status = "COMPLETE"
	if err := r.db.Save(ride).Error; err != nil {
		return err
	}
	return nil
}
