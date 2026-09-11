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
	return r.db.Create(rating).Error
}

func (r *RatingRepo) GetByID(id uint) (*domain.Rating, error) {
	var rate domain.Rating
	if err := r.db.Where("id = ?", id).First(&rate).Error; err != nil {
		return nil, err
	}
	return &rate, nil
}

func (r *RatingRepo) GetByRideID(rideID uint) ([]domain.Rating, error) {
	var rates []domain.Rating
	if err := r.db.Where("ride_id = ?", rideID).Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *RatingRepo) GetByToUserID(toUserID uint) ([]domain.Rating, error) {
	var rates []domain.Rating
	if err := r.db.Where("to_user_id = ?", toUserID).Find(&rates).Error; err != nil {
		return nil, err
	}
	return rates, nil
}

func (r *RatingRepo) GetAverageRating(toUserID uint) (float64, int, error) {
	type result struct {
		AvgScore float64
		Count    int
	}
	var res result
	err := r.db.Model(&domain.Rating{}).
		Select("COALESCE(AVG(score), 0) as avg_score, COUNT(id) as count").
		Where("to_user_id = ?", toUserID).
		Scan(&res).Error
	if err != nil {
		return 0, 0, err
	}
	return res.AvgScore, res.Count, nil
}

func (r *RatingRepo) Update(rating *domain.Rating) error {
	return r.db.Save(rating).Error
}

func (r *RatingRepo) Delete(id uint) error {
	return r.db.Where("id = ?", id).Delete(&domain.Rating{}).Error
}
