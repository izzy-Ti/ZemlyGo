package dto

import "time"

type UserDTO struct {
	ID             uint      `json:"id"`
	AuthID         string    `json:"auth_id"`
	SupabaseUserID string    `json:"supabase_user_id,omitempty"` // backward compatibility
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}
