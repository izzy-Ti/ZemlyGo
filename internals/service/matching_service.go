package service

import (
	"strconv"

	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type MatchingService struct {
	driverRepo  interfaces.DriverRepository
	rideRepo    interfaces.RideRepository
	vehicleRepo interfaces.VehicleRepository
	hub         *realtime.Hub
}

func NewMatchingService(
	driverRepo interfaces.DriverRepository,
	rideRepo interfaces.RideRepository,
	vehicleRepo interfaces.VehicleRepository,
	hub *realtime.Hub,
) *MatchingService {
	return &MatchingService{
		driverRepo:  driverRepo,
		rideRepo:    rideRepo,
		vehicleRepo: vehicleRepo,
		hub:         hub,
	}
}

func (m *MatchingService) FindEligibleDrivers(pickupLat, pickupLng, radiusKm float64) ([]domain.Drivers, error) {
	if radiusKm <= 0 {
		radiusKm = 10.0 // Default 10km radius
	}

	nearbyDrivers, err := m.driverRepo.GetNearbyDrivers(pickupLat, pickupLng, radiusKm)
	if err != nil {
		return nil, err
	}

	var eligible []domain.Drivers
	for _, driver := range nearbyDrivers {
		// Verify driver does not have an active ride
		activeRide, _ := m.rideRepo.GetActiveRideByDriver(driver.ID)
		if activeRide != nil {
			continue
		}

		// Verify driver has an active vehicle
		vehicle, _ := m.vehicleRepo.GetActiveVehicle(driver.ID)
		if vehicle == nil {
			continue
		}

		eligible = append(eligible, driver)
	}

	return eligible, nil
}

func (m *MatchingService) BroadcastRideRequest(ride *domain.Ride, eligibleDrivers []domain.Drivers) {
	eventData := map[string]interface{}{
		"ride_id":             ride.ID,
		"rider_id":            ride.RiderID,
		"pickup_lat":          ride.PickupLat,
		"pickup_lng":          ride.PickupLng,
		"pickup_address":      ride.PickupAddress,
		"destination_lat":     ride.DestinationLat,
		"destination_lng":     ride.DestinationLng,
		"destination_address": ride.DestinationAddress,
		"fare":                ride.Fare,
		"ride_type":           ride.RideType,
		"status":              ride.Status,
	}

	event := realtime.Event{
		Type: "ride.requested",
		Data: eventData,
	}

	// Notify eligible drivers by their user_id via WebSocket
	for _, drv := range eligibleDrivers {
		driverUserIDStr := strconv.FormatUint(uint64(drv.UserId), 10)
		m.hub.SendToUser(driverUserIDStr, event)
	}
}
