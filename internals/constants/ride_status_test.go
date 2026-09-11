package constants

import "testing"

func TestIsValidRideStatusTransition(t *testing.T) {
	tests := []struct {
		current string
		next    string
		valid   bool
	}{
		// Happy path
		{RideStatusRequested, RideStatusAccepted, true},
		{RideStatusAccepted, RideStatusArrived, true},
		{RideStatusArrived, RideStatusInProgress, true},
		{RideStatusInProgress, RideStatusCompleted, true},

		// Cancellations
		{RideStatusRequested, RideStatusCancelled, true},
		{RideStatusAccepted, RideStatusCancelled, true},
		{RideStatusArrived, RideStatusCancelled, true},
		{RideStatusInProgress, RideStatusCancelled, true},

		// Invalid transitions
		{RideStatusRequested, RideStatusInProgress, false},
		{RideStatusRequested, RideStatusCompleted, false},
		{RideStatusCompleted, RideStatusCancelled, false},
		{RideStatusCancelled, RideStatusRequested, false},
		{RideStatusCompleted, RideStatusInProgress, false},
		{"unknown_status", RideStatusAccepted, false},
	}

	for _, tt := range tests {
		got := IsValidRideStatusTransition(tt.current, tt.next)
		if got != tt.valid {
			t.Errorf("IsValidRideStatusTransition(%s, %s) = %v; expected %v", tt.current, tt.next, got, tt.valid)
		}
	}
}

func TestIsActiveStatus(t *testing.T) {
	if !IsActiveStatus(RideStatusRequested) {
		t.Errorf("expected %s to be active", RideStatusRequested)
	}
	if !IsActiveStatus(RideStatusInProgress) {
		t.Errorf("expected %s to be active", RideStatusInProgress)
	}
	if IsActiveStatus(RideStatusCompleted) {
		t.Errorf("expected %s NOT to be active", RideStatusCompleted)
	}
	if IsActiveStatus(RideStatusCancelled) {
		t.Errorf("expected %s NOT to be active", RideStatusCancelled)
	}
}
