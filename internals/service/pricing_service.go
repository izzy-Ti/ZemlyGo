package service

import (
	"math"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

type PricingConfig struct {
	BaseFare float64
	PerKm    float64
	MinFare  float64
}

var DefaultRates = map[string]PricingConfig{
	constants.RideTypeStandard: {BaseFare: 5.00, PerKm: 1.50, MinFare: 7.00},
	constants.RideTypeVan:      {BaseFare: 8.00, PerKm: 2.50, MinFare: 12.00},
	constants.RideTypeFamily:   {BaseFare: 7.00, PerKm: 2.20, MinFare: 10.00},
}

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (p *PricingService) CalculateDistance(pLat, pLng, dLat, dLng float64) float64 {
	dist := utils.Haversine(pLat, pLng, dLat, dLng)
	return math.Round(dist*100) / 100
}

func (p *PricingService) CalculateFare(rideType string, distanceKm float64, surgeMultiplier float64) float64 {
	rate, exists := DefaultRates[rideType]
	if !exists {
		rate = DefaultRates[constants.RideTypeStandard]
	}

	if surgeMultiplier <= 0 {
		surgeMultiplier = 1.0
	}

	rawFare := (rate.BaseFare + (distanceKm * rate.PerKm)) * surgeMultiplier
	if rawFare < rate.MinFare {
		rawFare = rate.MinFare
	}

	return math.Round(rawFare*100) / 100
}

func (p *PricingService) EstimateAllOptions(pLat, pLng, dLat, dLng float64) dto.EstimateFareResponse {
	distanceKm := p.CalculateDistance(pLat, pLng, dLat, dLng)
	durationMin := int(math.Ceil(distanceKm/30.0*60)) + 3 // 30km/h avg + 3 min traffic

	rideTypes := []string{
		constants.RideTypeStandard,
		constants.RideTypeVan,
		constants.RideTypeFamily,
	}

	var options []dto.FareOption
	for _, rt := range rideTypes {
		fare := p.CalculateFare(rt, distanceKm, 1.0)
		options = append(options, dto.FareOption{
			RideType:    rt,
			Fare:        fare,
			DistanceKm:  distanceKm,
			DurationMin: durationMin,
		})
	}

	return dto.EstimateFareResponse{
		DistanceKm: distanceKm,
		Options:    options,
	}
}

func (p *PricingService) GetSurgeHeatmap() *dto.SurgeHeatmapResponse {
	zones := []dto.SurgeZone{
		{ZoneName: "Downtown / Financial District", CenterLat: 40.7074, CenterLng: -74.0113, Multiplier: 1.4, DemandLevel: "surge"},
		{ZoneName: "Midtown / Times Square", CenterLat: 40.7580, CenterLng: -73.9855, Multiplier: 1.8, DemandLevel: "surge"},
		{ZoneName: "JFK International Airport", CenterLat: 40.6413, CenterLng: -73.7781, Multiplier: 1.2, DemandLevel: "high"},
		{ZoneName: "Brooklyn Heights", CenterLat: 40.6960, CenterLng: -73.9933, Multiplier: 1.0, DemandLevel: "normal"},
		{ZoneName: "Upper West Side", CenterLat: 40.7870, CenterLng: -73.9754, Multiplier: 1.1, DemandLevel: "normal"},
	}

	return &dto.SurgeHeatmapResponse{
		Zones:     zones,
		Timestamp: time.Now(),
	}
}
