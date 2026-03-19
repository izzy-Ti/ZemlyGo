package service

import (
	"math"

	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
)

type RideService struct {
	hub *realtime.Hub
}

func NewRideSerive(hub *realtime.Hub) *RideService {
	return &RideService{hub: hub}
}

func (s *RideService) NotifyDriverNewRide(driverID, padd, dadd, rideType string, rideID uint, Plng, plat, dlng, dlat float64) {
	const R = 6371

	dLat := (dlat - plat) * (math.Pi / 18)
	dLon := (dlng - Plng) * (math.Pi / 18)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(dlat*(math.Pi/180))*math.Cos(dlat*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	baseFare := 50.0
	perKm := 25.0
	if rideType == "standard" {
		perKm = 25.0
	}
	if rideType == "van" {
		perKm = 40.0
	}
	if rideType == "family" {
		perKm = 40.0
	}

	fare := baseFare + (c * perKm)

	s.hub.SendToUser(driverID, realtime.Event{
		Type: "ride.requested",
		Data: map[string]interface{}{
			"ride_id": rideID,
			"pickup": map[string]interface{}{
				"lat":     plat,
				"lng":     Plng,
				"address": padd,
			},
			"destination": map[string]interface{}{
				"lat":     dlat,
				"lng":     dlng,
				"address": dadd,
			},
			"fare":      fare,
			"ride_type": rideType,
		},
	})
}