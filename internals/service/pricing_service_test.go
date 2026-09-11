package service

import (
	"testing"

	"github.com/izzy-Ti/ZemlyGo/internals/constants"
)

func TestPricingService(t *testing.T) {
	pricing := NewPricingService()

	// 1. Test standard minimum fare
	// 0.5 km: base $5.00 + (0.5 * 1.50 = 0.75) = $5.75, which is less than min fare $7.00
	fareMin := pricing.CalculateFare(constants.RideTypeStandard, 0.5, 1.0)
	if fareMin != 7.00 {
		t.Errorf("expected minimum fare $7.00, got %f", fareMin)
	}

	// 2. Test standard normal fare
	// 10 km: base $5.00 + (10 * 1.50 = 15.00) = $20.00
	fare10Km := pricing.CalculateFare(constants.RideTypeStandard, 10.0, 1.0)
	if fare10Km != 20.00 {
		t.Errorf("expected fare $20.00, got %f", fare10Km)
	}

	// 3. Test van fare
	// 10 km: base $8.00 + (10 * 2.50 = 25.00) = $33.00
	fareVan := pricing.CalculateFare(constants.RideTypeVan, 10.0, 1.0)
	if fareVan != 33.00 {
		t.Errorf("expected fare $33.00, got %f", fareVan)
	}

	// 4. Test surge multiplier
	// 10 km standard with 1.5x surge: $20.00 * 1.5 = $30.00
	fareSurge := pricing.CalculateFare(constants.RideTypeStandard, 10.0, 1.5)
	if fareSurge != 30.00 {
		t.Errorf("expected fare $30.00 with surge, got %f", fareSurge)
	}

	// 5. Test EstimateAllOptions
	estimates := pricing.EstimateAllOptions(40.7128, -74.0060, 40.7589, -73.9851)
	if estimates.DistanceKm <= 0 {
		t.Errorf("expected positive distance, got %f", estimates.DistanceKm)
	}
	if len(estimates.Options) != 3 {
		t.Errorf("expected 3 fare options, got %d", len(estimates.Options))
	}
}
