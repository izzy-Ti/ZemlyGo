package service

import (
	"errors"
	"math"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/domain"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/repository/interfaces"
)

type DriverService struct {
	driverRepo  interfaces.DriverRepository
	userRepo    interfaces.UserRepository
	vehicleRepo interfaces.VehicleRepository
	rideRepo    interfaces.RideRepository
	paymentRepo interfaces.PaymentRepository
}

func NewDriverService(
	driverRepo interfaces.DriverRepository,
	userRepo interfaces.UserRepository,
	vehicleRepo interfaces.VehicleRepository,
	rideRepo interfaces.RideRepository,
	paymentRepo interfaces.PaymentRepository,
) *DriverService {
	return &DriverService{
		driverRepo:  driverRepo,
		userRepo:    userRepo,
		vehicleRepo: vehicleRepo,
		rideRepo:    rideRepo,
		paymentRepo: paymentRepo,
	}
}

func (d *DriverService) RegisterDriver(userID uint, req dto.RegisterDriverRequest) (*dto.DriverDTO, error) {
	existing, _ := d.driverRepo.GetByUserID(userID)
	if existing != nil {
		return nil, errors.New("driver registration already exists for this user")
	}

	driver := &domain.Drivers{
		UserId:     userID,
		LicenseNo:  req.LicenseNo,
		IsApproved: false,
		IsOnline:   false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := d.driverRepo.Create(driver); err != nil {
		return nil, err
	}

	_ = d.userRepo.UpdateRole(userID, constants.RoleDriver)

	return &dto.DriverDTO{
		ID:         driver.ID,
		UserID:     driver.UserId,
		LicenseNo:  driver.LicenseNo,
		IsApproved: driver.IsApproved,
		IsOnline:   driver.IsOnline,
		CreatedAt:  driver.CreatedAt,
	}, nil
}

func (d *DriverService) GetDriverByUserID(userID uint) (*dto.DriverDTO, error) {
	driver, err := d.driverRepo.GetByUserID(userID)
	if err != nil || driver == nil {
		return nil, constants.ErrNotFound
	}
	return d.enrichDriver(driver)
}

func (d *DriverService) GetDriverByID(driverID uint) (*dto.DriverDTO, error) {
	driver, err := d.driverRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return nil, constants.ErrNotFound
	}
	return d.enrichDriver(driver)
}

func (d *DriverService) SetOnline(driverID uint, online bool) error {
	driver, err := d.driverRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return constants.ErrNotFound
	}

	if online {
		if !driver.IsApproved {
			return constants.ErrDriverNotApproved
		}
		vehicle, err := d.vehicleRepo.GetActiveVehicle(driverID)
		if err != nil || vehicle == nil {
			return constants.ErrNoActiveVehicle
		}
	}

	return d.driverRepo.SetOnline(driverID, online)
}

func (d *DriverService) ApproveDriver(driverID uint) error {
	driver, err := d.driverRepo.GetByID(driverID)
	if err != nil || driver == nil {
		return constants.ErrNotFound
	}
	driver.IsApproved = true
	driver.UpdatedAt = time.Now()
	return d.driverRepo.Update(driver)
}

func (d *DriverService) GetNearbyDrivers(lat, lng, radiusKm float64) ([]dto.NearbyDriverDTO, error) {
	if radiusKm <= 0 {
		radiusKm = 10.0
	}

	drivers, err := d.driverRepo.GetNearbyDrivers(lat, lng, radiusKm)
	if err != nil {
		return nil, err
	}

	var results []dto.NearbyDriverDTO
	for _, drv := range drivers {
		v, _ := d.vehicleRepo.GetActiveVehicle(drv.ID)
		var vDTO *dto.VehicleDTO
		if v != nil {
			vDTO = &dto.VehicleDTO{
				ID:          v.ID,
				DriverID:    v.DriverID,
				PlateNumber: v.PlateNumber,
				Model:       v.Model,
				Color:       v.Color,
				Status:      v.Status,
				CreatedAt:   v.CreatedAt,
			}
		}

		results = append(results, dto.NearbyDriverDTO{
			DriverID: drv.ID,
			Vehicle:  vDTO,
		})
	}

	return results, nil
}

func (d *DriverService) GetEarnings(driverID uint, period string) (*dto.DriverEarningsDTO, error) {
	rides, err := d.rideRepo.GetRidesByDriver(driverID, 100, 0)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var cutoff time.Time
	switch period {
	case "today":
		cutoff = now.Truncate(24 * time.Hour)
	case "weekly":
		cutoff = now.AddDate(0, 0, -7)
	default:
		cutoff = time.Time{} // all time
	}

	var trips []dto.TripEarningItem
	var faresTotal, tipsTotal float64

	for _, ride := range rides {
		if ride.Status != constants.RideStatusCompleted {
			continue
		}
		if !cutoff.IsZero() && ride.CreatedAt.Before(cutoff) {
			continue
		}

		fare := ride.Fare
		tip := 0.0
		if d.paymentRepo != nil {
			if pay, _ := d.paymentRepo.GetByRideID(ride.ID); pay != nil {
				tip = pay.Tip
			}
		}

		faresTotal += fare
		tipsTotal += tip
		trips = append(trips, dto.TripEarningItem{
			RideID:   ride.ID,
			Fare:     fare,
			Tip:      tip,
			Total:    math.Round((fare+tip)*100) / 100,
			RideType: ride.RideType,
			Date:     ride.CreatedAt,
		})
	}

	totalEarnings := math.Round((faresTotal+tipsTotal)*100) / 100

	return &dto.DriverEarningsDTO{
		DriverID:       driverID,
		Period:         period,
		TotalEarnings:  totalEarnings,
		FaresTotal:     math.Round(faresTotal*100) / 100,
		TipsTotal:      math.Round(tipsTotal*100) / 100,
		CompletedRides: len(trips),
		Trips:          trips,
	}, nil
}

func (d *DriverService) enrichDriver(driver *domain.Drivers) (*dto.DriverDTO, error) {
	out := &dto.DriverDTO{
		ID:         driver.ID,
		UserID:     driver.UserId,
		LicenseNo:  driver.LicenseNo,
		IsApproved: driver.IsApproved,
		IsOnline:   driver.IsOnline,
		CreatedAt:  driver.CreatedAt,
	}

	if user, err := d.userRepo.GetByID(driver.UserId); err == nil && user != nil {
		out.User = &dto.UserDTO{
			ID:             user.ID,
			AuthID:         user.AuthID,
			SupabaseUserID: user.AuthID,
			Name:           user.Name,
			Email:          user.Email,
			Phone:          user.Phone,
			Role:           user.Role,
			CreatedAt:      user.CreatedAt,
		}
	}

	if vehicle, err := d.vehicleRepo.GetActiveVehicle(driver.ID); err == nil && vehicle != nil {
		out.Vehicle = &dto.VehicleDTO{
			ID:          vehicle.ID,
			DriverID:    vehicle.DriverID,
			PlateNumber: vehicle.PlateNumber,
			Model:       vehicle.Model,
			Color:       vehicle.Color,
			Status:      vehicle.Status,
			CreatedAt:   vehicle.CreatedAt,
		}
	}

	return out, nil
}
