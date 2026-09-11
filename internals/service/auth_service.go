package service

import (
	"errors"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

type AuthService struct {
	userRepo  interfaces.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo interfaces.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (a *AuthService) SyncUser(authID, name, email, phone, role string) (*domain.Users, error) {
	if authID == "" {
		return nil, errors.New("auth_id is required")
	}

	if role == "" || !constants.IsValidRole(role) {
		role = constants.RoleRider
	}

	// Check if user already exists by auth_id
	user, err := a.userRepo.GetByAuthID(authID)
	if err == nil && user != nil {
		user.Name = name
		user.Phone = phone
		if user.Role == "" {
			user.Role = role
		}
		if err := a.userRepo.Update(user); err != nil {
			return nil, err
		}
		return user, nil
	}

	// Check if user exists by email
	userByEmail, err := a.userRepo.GetByEmail(email)
	if err == nil && userByEmail != nil {
		userByEmail.AuthID = authID
		userByEmail.Name = name
		userByEmail.Phone = phone
		if err := a.userRepo.Update(userByEmail); err != nil {
			return nil, err
		}
		return userByEmail, nil
	}

	// Create new user
	newUser := &domain.Users{
		AuthID:    authID,
		Name:      name,
		Email:     email,
		Phone:     phone,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (a *AuthService) GenerateDevToken(req dto.DevTokenRequest) (*dto.DevTokenResponse, error) {
	var user *domain.Users
	var err error

	if req.UserID > 0 {
		user, err = a.userRepo.GetByID(req.UserID)
	} else if req.Email != "" {
		user, err = a.userRepo.GetByEmail(req.Email)
	}

	if user == nil {
		authID := req.AuthID
		if authID == "" {
			authID = req.SupabaseUserID
		}
		if authID == "" {
			authID = "neon-user-" + strconvFormat(time.Now().UnixNano())
		}
		role := req.Role
		if role == "" {
			role = constants.RoleRider
		}
		user = &domain.Users{
			AuthID:    authID,
			Name:      "Dev User",
			Email:     req.Email,
			Phone:     "+1234567890",
			Role:      role,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := a.userRepo.Create(user); err != nil {
			return nil, err
		}
	}

	claims := utils.JWTClaims{
		UserID:         user.ID,
		SupabaseUserID: user.AuthID,
		Email:          user.Email,
		Role:           user.Role,
	}

	tokenStr, exp, err := utils.GenerateJWT(a.jwtSecret, claims, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &dto.DevTokenResponse{
		Token:     tokenStr,
		ExpiresAt: exp,
		User: dto.UserDTO{
			ID:             user.ID,
			AuthID:         user.AuthID,
			SupabaseUserID: user.AuthID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			Role:           user.Role,
			CreatedAt:      user.CreatedAt,
		},
	}, nil
}

func strconvFormat(v int64) string {
	var s string
	for v > 0 {
		s = string(rune('0'+(v%10))) + s
		v /= 10
	}
	return s
}
