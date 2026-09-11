package constants

const (
	RideTypeStandard = "standard"
	RideTypeVan      = "van"
	RideTypeFamily   = "family"
)

func IsValidRideType(t string) bool {
	switch t {
	case RideTypeStandard, RideTypeVan, RideTypeFamily:
		return true
	default:
		return false
	}
}
