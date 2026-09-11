package service

import (
	"errors"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type UserService struct {
	userRepo interfaces.UserRepository
}

func NewUserService(userRepo interfaces.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) GetByID(id uint) (*dto.UserDTO, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil || user == nil {
		return nil, constants.ErrNotFound
	}
	return &dto.UserDTO{
		ID:             user.ID,
		SupabaseUserID: user.SupabaseUserID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (u *UserService) UpdateProfile(id uint, req dto.UpdateUserRequest) (*dto.UserDTO, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil || user == nil {
		return nil, constants.ErrNotFound
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserDTO{
		ID:             user.ID,
		SupabaseUserID: user.SupabaseUserID,
		Name:           user.Name,
		Email:          user.Email,
		Phone:          user.Phone,
		Role:           user.Role,
		CreatedAt:      user.CreatedAt,
	}, nil
}

func (u *UserService) ListUsers(limit, offset int) ([]dto.UserDTO, error) {
	if limit <= 0 {
		limit = 20
	}
	users, err := u.userRepo.List(limit, offset)
	if err != nil {
		return nil, err
	}

	var res []dto.UserDTO
	for _, user := range users {
		res = append(res, dto.UserDTO{
			ID:             user.ID,
			SupabaseUserID: user.SupabaseUserID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			Role:           user.Role,
			CreatedAt:      user.CreatedAt,
		})
	}
	return res, nil
}

func (u *UserService) UpdateRole(userID uint, role string) error {
	if !constants.IsValidRole(role) {
		return errors.New("invalid role")
	}
	return u.userRepo.UpdateRole(userID, role)
}
