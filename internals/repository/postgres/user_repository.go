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

func (s *UserRepository) Create(user *domain.Users) error {
	if user.AuthID == "" && user.SupabaseUserID != "" {
		user.AuthID = user.SupabaseUserID
	}
	return s.db.Create(user).Error
}

func (s *UserRepository) GetByID(id uint) (*domain.Users, error) {
	var user domain.Users
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserRepository) GetByEmail(email string) (*domain.Users, error) {
	var user domain.Users
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserRepository) GetByAuthID(authID string) (*domain.Users, error) {
	var user domain.Users
	if err := s.db.Where("auth_id = ?", authID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserRepository) GetBySupabaseID(supabaseID string) (*domain.Users, error) {
	return s.GetByAuthID(supabaseID)
}

func (s *UserRepository) Update(user *domain.Users) error {
	return s.db.Save(user).Error
}

func (s *UserRepository) Delete(id uint) error {
	return s.db.Where("id = ?", id).Delete(&domain.Users{}).Error
}

func (s *UserRepository) List(limit, offset int) ([]domain.Users, error) {
	var users []domain.Users
	if err := s.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (s *UserRepository) UpdateRole(userID uint, role string) error {
	return s.db.Model(&domain.Users{}).Where("id = ?", userID).Update("role", role).Error
}
