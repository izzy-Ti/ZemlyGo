package domain

import "time"

type Users struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	AuthID         string    `gorm:"column:auth_id;uniqueIndex;not null" json:"auth_id"`
	SupabaseUserID string    `gorm:"-" json:"supabase_user_id,omitempty"` // compatibility alias
	Name           string    `gorm:"not null" json:"name"`
	Email          string    `gorm:"uniqueIndex;not null" json:"email"`
	Phone          string    `gorm:"not null" json:"phone"`
	Role           string    `gorm:"not null;default:rider" json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (u *Users) AfterFind() error {
	if u.SupabaseUserID == "" {
		u.SupabaseUserID = u.AuthID
	}
	return nil
}
