package service

import (
	"strconv"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type RideService struct {
	rideRepo       interfaces.RideRepository
	driverRepo     interfaces.DriverRepository
	userRepo       interfaces.UserRepository
	vehicleRepo    interfaces.VehicleRepository
	paymentRepo    interfaces.PaymentRepository
	pricingService *PricingService
	matchingService *MatchingService
	hub            *realtime.Hub
}

func NewRideService(
	rideRepo interfaces.RideRepository,
	driverRepo interfaces.DriverRepository,
	userRepo interfaces.UserRepository,
	vehicleRepo interfaces.VehicleRepository,
	paymentRepo interfaces.PaymentRepository,
	pricingService *PricingService,
	matchingService *MatchingService,
	hub *realtime.Hub,
) *RideService {
	return &RideService{
		rideRepo:        rideRepo,
		driverRepo:      driverRepo,
		userRepo:        userRepo,
		vehicleRepo:     vehicleRepo,
		paymentRepo:     paymentRepo,
		pricingService:  pricingService,
		matchingService: matchingService,
		hub:             hub,
	}
}

func (s *RideService) RequestRide(riderID uint, req dto.RequestRideRequest) (*domain.Ride, error) {
	// 1. Ensure rider has no active ride
	activeRide, _ := s.rideRepo.GetActiveRideByRider(riderID)
	if activeRide != nil {
		return nil, constants.ErrActiveRideExists
	}

	rideType := req.RideType
	if !constants.IsValidRideType(rideType) {
		rideType = constants.RideTypeStandard
	}

	// 2. Calculate distance and fare
	distKm := s.pricingService.CalculateDistance(req.PickupLat, req.PickupLng, req.DestinationLat, req.DestinationLng)
	fare := s.pricingService.CalculateFare(rideType, distKm, 1.0)

	ride := &domain.Ride{
		RiderID:            riderID,
		PickupLat:          req.PickupLat,
		PickupLng:          req.PickupLng,
		PickupAddress:      req.PickupAddress,
		DestinationLat:     req.DestinationLat,
		DestinationLng:     req.DestinationLng,
		DestinationAddress: req.DestinationAddress,
		RideType:           rideType,
		Status:             constants.RideStatusRequested,
		Fare:               fare,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.rideRepo.Create(ride); err != nil {
		return nil, err
	}

	// 3. Find candidate drivers and broadcast ride request asynchronously
	go func(r *domain.Ride) {
		drivers, err := s.matchingService.FindEligibleDrivers(r.PickupLat, r.PickupLng, 10.0)
		if err == nil && len(drivers) > 0 {
			s.matchingService.BroadcastRideRequest(r, drivers)
		}
	}(ride)

	return ride, nil
}

func (s *RideService) AcceptRide(driverID uint, rideID uint) (*domain.Ride, error) {
	// 1. Verify driver profile & eligibility
	driver, err := s.driverRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return nil, constants.ErrNotFound
	}
	if !driver.IsApproved {
		return nil, constants.ErrDriverNotApproved
	}
	if !driver.IsOnline {
		return nil, constants.ErrDriverOffline
	}

	// 2. Verify driver does not have another active ride
	activeRide, _ := s.rideRepo.GetActiveRideByDriver(driverID)
	if activeRide != nil {
		return nil, constants.ErrActiveRideExists
	}

	// 3. Verify ride existence and status
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}
	if ride.Status != constants.RideStatusRequested {
		return nil, constants.ErrRideAlreadyAccepted
	}

	// 4. Assign driver and transition status
	if err := s.rideRepo.AssignDriver(rideID, driverID); err != nil {
		return nil, err
	}

	ride.DriverID = &driverID
	ride.Status = constants.RideStatusAccepted

	// 5. Notify rider via WebSocket
	riderUserIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
	vehicle, _ := s.vehicleRepo.GetActiveVehicle(driverID)
	driverUser, _ := s.userRepo.GetByID(driver.UserId)

	s.hub.SendToUser(riderUserIDStr, realtime.Event{
		Type: "ride.accepted",
		Data: map[string]interface{}{
			"ride_id":   ride.ID,
			"driver_id": driverID,
			"driver":    driverUser,
			"vehicle":   vehicle,
			"status":    ride.Status,
		},
	})

	return ride, nil
}

func (s *RideService) ArriveAtPickup(driverID uint, rideID uint) (*domain.Ride, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, constants.ErrForbidden
	}

	if !constants.IsValidRideStatusTransition(ride.Status, constants.RideStatusArrived) {
		return nil, constants.ErrInvalidRideState
	}

	if err := s.rideRepo.UpdateStatus(rideID, constants.RideStatusArrived); err != nil {
		return nil, err
	}
	ride.Status = constants.RideStatusArrived

	riderIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
	s.hub.SendToUser(riderIDStr, realtime.Event{
		Type: "ride.driver_arrived",
		Data: map[string]interface{}{
			"ride_id": ride.ID,
			"status":  ride.Status,
		},
	})

	return ride, nil
}

func (s *RideService) StartTrip(driverID uint, rideID uint) (*domain.Ride, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, constants.ErrForbidden
	}

	if !constants.IsValidRideStatusTransition(ride.Status, constants.RideStatusInProgress) {
		return nil, constants.ErrInvalidRideState
	}

	if err := s.rideRepo.UpdateStatus(rideID, constants.RideStatusInProgress); err != nil {
		return nil, err
	}
	ride.Status = constants.RideStatusInProgress

	riderIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
	s.hub.SendToUser(riderIDStr, realtime.Event{
		Type: "ride.in_progress",
		Data: map[string]interface{}{
			"ride_id": ride.ID,
			"status":  ride.Status,
		},
	})

	return ride, nil
}

func (s *RideService) CompleteTrip(driverID uint, rideID uint) (*domain.Ride, *domain.Payment, error) {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, nil, constants.ErrNotFound
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, nil, constants.ErrForbidden
	}

	if !constants.IsValidRideStatusTransition(ride.Status, constants.RideStatusCompleted) {
		return nil, nil, constants.ErrInvalidRideState
	}

	// Confirm final fare
	if err := s.rideRepo.CompleteRide(rideID, ride.Fare); err != nil {
		return nil, nil, err
	}
	ride.Status = constants.RideStatusCompleted

	// Auto-create pending payment record
	payment := &domain.Payment{
		RideID:    ride.ID,
		Amount:    ride.Fare,
		Method:    "card",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_ = s.paymentRepo.Create(payment)

	// Notify rider and driver
	riderIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
	driverIDStr := strconv.FormatUint(uint64(driverID), 10)

	completionEvent := realtime.Event{
		Type: "ride.completed",
		Data: map[string]interface{}{
			"ride_id":    ride.ID,
			"fare":       ride.Fare,
			"payment_id": payment.ID,
			"status":     ride.Status,
		},
	}
	s.hub.SendToUser(riderIDStr, completionEvent)
	s.hub.SendToUser(driverIDStr, completionEvent)

	return ride, payment, nil
}

func (s *RideService) CancelRide(userID uint, isDriver bool, rideID uint, reason string) error {
	ride, err := s.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return constants.ErrNotFound
	}

	if isDriver {
		if ride.DriverID == nil || *ride.DriverID != userID {
			return constants.ErrForbidden
		}
	} else {
		if ride.RiderID != userID {
			return constants.ErrForbidden
		}
	}

	if !constants.IsValidRideStatusTransition(ride.Status, constants.RideStatusCancelled) {
		return constants.ErrInvalidRideState
	}

	if err := s.rideRepo.CancelRide(rideID); err != nil {
		return err
	}

	cancelEvent := realtime.Event{
		Type: "ride.cancelled",
		Data: map[string]interface{}{
			"ride_id": rideID,
			"reason":  reason,
		},
	}

	// Notify counterparty
	if isDriver {
		riderIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
		s.hub.SendToUser(riderIDStr, cancelEvent)
	} else if ride.DriverID != nil {
		driverIDStr := strconv.FormatUint(uint64(*ride.DriverID), 10)
		s.hub.SendToUser(driverIDStr, cancelEvent)
	}

	return nil
}

func (s *RideService) GetActiveRide(userID uint, isDriver bool) (*domain.Ride, error) {
	if isDriver {
		return s.rideRepo.GetActiveRideByDriver(userID)
	}
	return s.rideRepo.GetActiveRideByRider(userID)
}

func (s *RideService) GetRideByID(rideID uint) (*domain.Ride, error) {
	return s.rideRepo.GetByID(rideID)
}

func (s *RideService) GetRideHistory(userID uint, isDriver bool, limit, offset int) ([]domain.Ride, error) {
	if limit <= 0 {
		limit = 20
	}
	if isDriver {
		return s.rideRepo.GetRidesByDriver(userID, limit, offset)
	}
	return s.rideRepo.GetRidesByRider(userID, limit, offset)
}

func (s *RideService) EnrichRideDetails(ride *domain.Ride) *dto.RideDTO {
	if ride == nil {
		return nil
	}

	out := &dto.RideDTO{
		ID:                 ride.ID,
		RiderID:            ride.RiderID,
		DriverID:           ride.DriverID,
		PickupLat:          ride.PickupLat,
		PickupLng:          ride.PickupLng,
		PickupAddress:      ride.PickupAddress,
		DestinationLat:     ride.DestinationLat,
		DestinationLng:     ride.DestinationLng,
		DestinationAddress: ride.DestinationAddress,
		RideType:           ride.RideType,
		Status:             ride.Status,
		Fare:               ride.Fare,
		CreatedAt:          ride.CreatedAt,
		UpdatedAt:          ride.UpdatedAt,
	}

	if rider, err := s.userRepo.GetByID(ride.RiderID); err == nil && rider != nil {
		out.Rider = &dto.UserDTO{
			ID:             rider.ID,
			SupabaseUserID: rider.SupabaseUserID,
			Name:           rider.Name,
			Email:          rider.Email,
			Phone:          rider.Phone,
			Role:           rider.Role,
			CreatedAt:      rider.CreatedAt,
		}
	}

	if ride.DriverID != nil {
		if driver, err := s.driverRepo.GetByID(*ride.DriverID); err == nil && driver != nil {
			driverDTO := &dto.DriverDTO{
				ID:         driver.ID,
				UserID:     driver.UserId,
				LicenseNo:  driver.LicenseNo,
				IsApproved: driver.IsApproved,
				IsOnline:   driver.IsOnline,
				CreatedAt:  driver.CreatedAt,
			}

			if vehicle, err := s.vehicleRepo.GetActiveVehicle(driver.ID); err == nil && vehicle != nil {
				driverDTO.Vehicle = &dto.VehicleDTO{
					ID:          vehicle.ID,
					DriverID:    vehicle.DriverID,
					PlateNumber: vehicle.PlateNumber,
					Model:       vehicle.Model,
					Color:       vehicle.Color,
					Status:      vehicle.Status,
					CreatedAt:   vehicle.CreatedAt,
				}
			}

			out.Driver = driverDTO
		}
	}

	return out
}

func (s *RideService) CheckRideAuthorization(ride *domain.Ride, userID uint) bool {
	if ride.RiderID == userID {
		return true
	}
	if ride.DriverID != nil && *ride.DriverID == userID {
		return true
	}
	return false
}

func (s *RideService) NotifyDriverNewRide(driverID, padd, dadd, rideType string, rideID uint, Plng, plat, dlng, dlat float64) {
	// Kept for backward compatibility
	fare := s.pricingService.CalculateFare(rideType, s.pricingService.CalculateDistance(plat, Plng, dlat, dlng), 1.0)
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