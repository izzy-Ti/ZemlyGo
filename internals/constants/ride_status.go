package constants

const (
	RideStatusRequested  = "requested"
	RideStatusAccepted   = "accepted"
	RideStatusArrived    = "arrived"
	RideStatusInProgress = "in_progress"
	RideStatusCompleted  = "completed"
	RideStatusCancelled  = "cancelled"
)

var ActiveRideStatuses = []string{
	RideStatusRequested,
	RideStatusAccepted,
	RideStatusArrived,
	RideStatusInProgress,
}

func IsActiveStatus(status string) bool {
	for _, s := range ActiveRideStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsValidRideStatusTransition checks if state change adheres to the Uber ride state machine
func IsValidRideStatusTransition(current, next string) bool {
	switch current {
	case RideStatusRequested:
		return next == RideStatusAccepted || next == RideStatusCancelled
	case RideStatusAccepted:
		return next == RideStatusArrived || next == RideStatusCancelled
	case RideStatusArrived:
		return next == RideStatusInProgress || next == RideStatusCancelled
	case RideStatusInProgress:
		return next == RideStatusCompleted || next == RideStatusCancelled
	case RideStatusCompleted, RideStatusCancelled:
		// Terminal states cannot transition further
		return false
	default:
		return false
	}
}
