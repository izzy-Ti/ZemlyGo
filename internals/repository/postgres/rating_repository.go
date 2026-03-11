package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type RatingRepo struct {
	db *gorm.DB
}

func NewRatingRepo(db *gorm.DB) interfaces.RatingRepository {
	return &RatingRepo{db: db}
}

func (r *RatingRepo) Create(rating *domain.Rating) error {
	if err := r.db.Save(rating).Error; err != nil {
		return err
	}
	return nil
}
func (r *RatingRepo) GetByID(id uint) (*domain.Rating, error) {
	var rate *domain.Rating
	if err := r.db.Where("id = ?", id).First(&rate).Error; err != nil {
		return nil, err
	}
	return rate, nil
}
func (r *RatingRepo) GetByRideID(rideID uint) ([]domain.Rating, error) {
	var rates []domain.Rating
	if err := r.db.Where("from_user_iD = ?", rideID).Find(rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}
func (r *RatingRepo) GetByToUserID(toUserID uint) ([]domain.Rating, error) {
	var rates []domain.Rating
	if err := r.db.Where("to_user_iD = ?", toUserID).Find(rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}
func (r *RatingRepo) Update(rating *domain.Rating) error {
	if err := r.db.Save(rating).Error; err != nil {
		return err
	}
	return nil
}
func (r *RatingRepo) Delete(id uint) error {
	if err := r.db.Where("id = ?", id).Delete(&domain.Rating{}).Error; err != nil {
		return err
	}
	return nil
}
