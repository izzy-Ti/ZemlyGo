package dto

import "time"

type SyncUserRequest struct {
	Name   string `json:"name" binding:"required"`
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone" binding:"required"`
	Role   string `json:"role"`    // default "rider"
	AuthID string `json:"auth_id"` // Optional Neon Auth user ID
}

type DevTokenRequest struct {
	UserID         uint   `json:"user_id"`
	AuthID         string `json:"auth_id"`
	SupabaseUserID string `json:"supabase_user_id"` // compatibility
	Email          string `json:"email" binding:"required,email"`
	Role           string `json:"role" binding:"required"`
}

type DevTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserDTO   `json:"user"`
}
