package postgres

import (
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type EmergencyRepository struct {
	db *gorm.DB
}

func NewEmergencyRepository(db *gorm.DB) interfaces.EmergencyRepository {
	return &EmergencyRepository{db: db}
}

func (r *EmergencyRepository) Create(alert *domain.EmergencyAlert) error {
	return r.db.Create(alert).Error
}

func (r *EmergencyRepository) GetByID(id uint) (*domain.EmergencyAlert, error) {
	var alert domain.EmergencyAlert
	if err := r.db.Where("id = ?", id).First(&alert).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

func (r *EmergencyRepository) GetByRideID(rideID uint) ([]domain.EmergencyAlert, error) {
	var alerts []domain.EmergencyAlert
	err := r.db.Where("ride_id = ?", rideID).Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}

func (r *EmergencyRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&domain.EmergencyAlert{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": time.Now(),
		}).Error
}

func (r *EmergencyRepository) ListOpenAlerts() ([]domain.EmergencyAlert, error) {
	var alerts []domain.EmergencyAlert
	err := r.db.Where("status != ?", "resolved").Order("created_at DESC").Find(&alerts).Error
	return alerts, err
}
