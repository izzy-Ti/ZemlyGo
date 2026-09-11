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

type RatingService struct {
	ratingRepo interfaces.RatingRepository
	rideRepo   interfaces.RideRepository
	driverRepo interfaces.DriverRepository
	hub        *realtime.Hub
}

func NewRatingService(
	ratingRepo interfaces.RatingRepository,
	rideRepo interfaces.RideRepository,
	driverRepo interfaces.DriverRepository,
	hub *realtime.Hub,
) *RatingService {
	return &RatingService{
		ratingRepo: ratingRepo,
		rideRepo:   rideRepo,
		driverRepo: driverRepo,
		hub:        hub,
	}
}

func (r *RatingService) SubmitRating(rideID uint, fromUserID uint, req dto.CreateRatingRequest) (*dto.RatingDTO, error) {
	if req.Score < 1 || req.Score > 5 {
		return nil, constants.ErrInvalidRatingScore
	}

	ride, err := r.rideRepo.GetByID(rideID)
	if err != nil || ride == nil {
		return nil, constants.ErrNotFound
	}

	if ride.Status != constants.RideStatusCompleted {
		return nil, constants.ErrRideNotCompleted
	}

	var toUserID uint
	if fromUserID == ride.RiderID {
		// Rider is rating the driver
		if ride.DriverID == nil {
			return nil, errors.New("ride has no assigned driver")
		}
		driver, err := r.driverRepo.GetByID(*ride.DriverID)
		if err != nil || driver == nil {
			return nil, errors.New("driver not found")
		}
		toUserID = driver.UserId
	} else {
		// Driver is rating the rider
		driver, err := r.driverRepo.GetByUserID(fromUserID)
		if err != nil || driver == nil || ride.DriverID == nil || *ride.DriverID != driver.ID {
			return nil, constants.ErrForbidden
		}
		toUserID = ride.RiderID
	}

	rating := &domain.Rating{
		RideID:     rideID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Score:      req.Score,
		Comment:    req.Comment,
		CreatedAt:  time.Now(),
	}

	if err := r.ratingRepo.Create(rating); err != nil {
		return nil, err
	}

	// Push rating notification to recipient
	toUserIDStr := strconv.FormatUint(uint64(toUserID), 10)
	r.hub.SendToUser(toUserIDStr, realtime.Event{
		Type: "ride.rated",
		Data: map[string]interface{}{
			"ride_id": rideID,
			"score":   rating.Score,
			"comment": rating.Comment,
		},
	})

	return &dto.RatingDTO{
		ID:         rating.ID,
		RideID:     rating.RideID,
		FromUserID: rating.FromUserID,
		ToUserID:   rating.ToUserID,
		Score:      rating.Score,
		Comment:    rating.Comment,
		CreatedAt:  rating.CreatedAt,
	}, nil
}

func (r *RatingService) GetDriverRatingSummary(driverID uint) (*dto.DriverRatingSummary, error) {
	driver, err := r.driverRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return nil, constants.ErrNotFound
	}

	avgScore, count, err := r.ratingRepo.GetAverageRating(driver.UserId)
	if err != nil {
		return nil, err
	}

	return &dto.DriverRatingSummary{
		DriverID:     driverID,
		AverageScore: math.Round(avgScore*10) / 10,
		TotalRatings: count,
	}, nil
}
