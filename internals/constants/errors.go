package constants

import "errors"

var (
	ErrNotFound            = errors.New("record not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden: insufficient permissions")
	ErrInvalidRideState    = errors.New("invalid ride state transition")
	ErrDriverNotApproved   = errors.New("driver is not approved")
	ErrDriverOffline       = errors.New("driver is offline")
	ErrNoActiveVehicle     = errors.New("driver has no active vehicle")
	ErrActiveRideExists    = errors.New("user or driver already has an active ride")
	ErrNoDriversAvailable  = errors.New("no drivers currently available nearby")
	ErrRideAlreadyAccepted = errors.New("ride has already been accepted by another driver")
	ErrInvalidRatingScore  = errors.New("rating score must be between 1 and 5")
	ErrRideNotCompleted    = errors.New("cannot pay or rate a ride that is not completed")
	ErrPaymentAlreadyMade  = errors.New("payment for this ride has already been completed")
)
