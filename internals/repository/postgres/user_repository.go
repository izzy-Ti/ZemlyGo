package postgres

import (
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) interfaces.UserRepository {
	return &UserRepository{db: db}
}
func (s *UserRepository) Create(user *domain.Users) error
func (s *UserRepository) GetByID(id uint) (*domain.Users, error) {
	var user *domain.Users
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserRepository) GetByEmail(email string) (*domain.Users, error) {
	var user *domain.Users
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserRepository) GetBySupabaseID(supabaseID string) (*domain.Users, error) {
	var user *domain.Users
	if err := s.db.Where("supabase_id = ?", supabaseID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (s *UserRepository) Update(user *domain.Users) error {
	if err := s.db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
func (s *UserRepository) Delete(id uint) error {
	if err := s.db.Where("id = ?", id).Delete(&domain.Users{}).Error; err != nil {
		return err
	}
	return nil
}
func (s *UserRepository) List(limit, offset int) ([]domain.Users, error) {
	var users []domain.Users
	if err := s.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (s *UserRepository) UpdateRole(userID uint, role string) error
