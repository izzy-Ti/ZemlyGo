package interfaces

import "github.com/izzy-Ti/ZemlyGo/internals/domain"

type UserRepository interface {
	Create(user *domain.Users) error
	GetByID(id uint) (*domain.Users, error)
	GetByEmail(email string) (*domain.Users, error)
	GetBySupabaseID(supabaseID string) (*domain.Users, error)

	Update(user *domain.Users) error
	Delete(id uint) error

	List(limit, offset int) ([]domain.Users, error)

	UpdateRole(userID uint, role string) error
}
