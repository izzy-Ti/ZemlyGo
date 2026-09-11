package service

import (
	"errors"
	"math"
	"strconv"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/realtime"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type PaymentService struct {
	paymentRepo interfaces.PaymentRepository
	rideRepo    interfaces.RideRepository
	hub         *realtime.Hub
}

func NewPaymentService(
	paymentRepo interfaces.PaymentRepository,
	rideRepo interfaces.RideRepository,
	hub *realtime.Hub,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		rideRepo:    rideRepo,
		hub:         hub,
	}
}

func (p *PaymentService) ProcessPayment(rideID uint, payerUserID uint, req dto.CreatePaymentRequest) (*dto.PaymentDTO, error) {
	ride, err := p.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	if ride.RiderID != payerUserID {
		return nil, constants.ErrForbidden
	}

	if ride.Status != constants.RideStatusCompleted {
		return nil, constants.ErrRideNotCompleted
	}

	payment, _ := p.paymentRepo.GetByRideID(rideID)
	if payment != nil {
		if payment.Status == "completed" {
			return nil, constants.ErrPaymentAlreadyMade
		}
		payment.Method = req.Method
		payment.Status = "completed"
		payment.UpdatedAt = time.Now()
		if err := p.paymentRepo.Update(payment); err != nil {
			return nil, err
		}
	} else {
		payment = &domain.Payment{
			RideID:    ride.ID,
			Amount:    ride.Fare,
			Method:    req.Method,
			Status:    "completed",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := p.paymentRepo.Create(payment); err != nil {
			return nil, err
		}
	}

	// Notify rider and driver via WebSocket
	riderIDStr := strconv.FormatUint(uint64(ride.RiderID), 10)
	paymentEvent := realtime.Event{
		Type: "payment.completed",
		Data: map[string]interface{}{
			"payment_id": payment.ID,
			"ride_id":    ride.ID,
			"amount":     payment.Amount,
			"method":     payment.Method,
			"status":     payment.Status,
		},
	}
	p.hub.SendToUser(riderIDStr, paymentEvent)

	if ride.DriverID != nil {
		driverIDStr := strconv.FormatUint(uint64(*ride.DriverID), 10)
		p.hub.SendToUser(driverIDStr, paymentEvent)
	}

	return &dto.PaymentDTO{
		ID:        payment.ID,
		RideID:    payment.RideID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}

func (p *PaymentService) AddTip(rideID uint, riderID uint, tipAmount float64) (*dto.PaymentDTO, error) {
	if tipAmount <= 0 {
		return nil, errors.New("tip amount must be positive")
	}

	ride, err := p.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	if ride.RiderID != riderID {
		return nil, constants.ErrForbidden
	}

	payment, err := p.paymentRepo.GetByRideID(rideID)
	if err != nil || payment == nil {
		return nil, errors.New("payment record not found")
	}

	payment.Tip = math.Round(tipAmount*100) / 100
	payment.UpdatedAt = time.Now()
	if err := p.paymentRepo.Update(payment); err != nil {
		return nil, err
	}

	// Broadcast tip event to driver
	if ride.DriverID != nil {
		driverIDStr := strconv.FormatUint(uint64(*ride.DriverID), 10)
		p.hub.SendToUser(driverIDStr, realtime.Event{
			Type: "ride.tipped",
			Data: map[string]interface{}{
				"ride_id":    rideID,
				"tip_amount": payment.Tip,
				"total":      payment.Amount + payment.Tip,
			},
		})
	}

	return &dto.PaymentDTO{
		ID:        payment.ID,
		RideID:    payment.RideID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}

func (p *PaymentService) GetPaymentByRideID(rideID uint, userID uint) (*dto.PaymentDTO, error) {
	ride, err := p.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	// Caller must be either rider or driver of this ride
	if ride.RiderID != userID && (ride.DriverID == nil || *ride.DriverID != userID) {
		return nil, constants.ErrForbidden
	}

	payment, err := p.paymentRepo.GetByRideID(rideID)
	if err != nil || payment == nil {
		return nil, constants.ErrNotFound
	}

	return &dto.PaymentDTO{
		ID:        payment.ID,
		RideID:    payment.RideID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt,
		UpdatedAt: payment.UpdatedAt,
	}, nil
}
