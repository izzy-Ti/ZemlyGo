package constants

const (
	RoleRider  = "rider"
	RoleDriver = "driver"
	RoleAdmin  = "admin"
)

func IsValidRole(role string) bool {
	switch role {
	case RoleRider, RoleDriver, RoleAdmin:
		return true
	default:
		return false
	}
}
