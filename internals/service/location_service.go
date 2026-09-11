package service

import (
	"strconv"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type LocationService struct {
	driverRepo interfaces.DriverRepository
	rideRepo   interfaces.RideRepository
	hub        *realtime.Hub
}

func NewLocationService(
	driverRepo interfaces.DriverRepository,
	rideRepo interfaces.RideRepository,
	hub *realtime.Hub,
) *LocationService {
	return &LocationService{
		driverRepo: driverRepo,
		rideRepo:   rideRepo,
		hub:        hub,
	}
}

func (l *LocationService) UpdateLocation(driverID uint, lat, lng float64) error {
	if err := l.driverRepo.UpdateLocation(driverID, lat, lng); err != nil {
		return err
	}

	// Check if this driver currently has an active ride
	activeRide, err := l.rideRepo.GetActiveRideByDriver(driverID)
	if err == nil && activeRide != nil {
		// Broadcast driver coordinates in realtime to rider
		riderIDStr := strconv.FormatUint(uint64(activeRide.RiderID), 10)
		l.hub.SendToUser(riderIDStr, realtime.Event{
			Type: "ride.driver_location",
			Data: map[string]interface{}{
				"ride_id":    activeRide.ID,
				"driver_id":  driverID,
				"lat":        lat,
				"lng":        lng,
				"updated_at": time.Now().Unix(),
			},
		})
	}

	return nil
}
