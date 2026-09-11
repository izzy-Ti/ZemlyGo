package utils

import (
	"math"
	"testing"
)

func TestHaversine(t *testing.T) {
	// Distance between Empire State Building (40.748817, -73.985428) and Times Square (40.758896, -73.985130)
	lat1, lng1 := 40.748817, -73.985428
	lat2, lng2 := 40.758896, -73.985130

	dist := Haversine(lat1, lng1, lat2, lng2)

	// Distance is approximately 1.12 km
	if math.Abs(dist-1.12) > 0.05 {
		t.Errorf("expected approx 1.12 km, got %f", dist)
	}

	// Distance from point to itself should be 0
	zeroDist := Haversine(lat1, lng1, lat1, lng1)
	if zeroDist != 0 {
		t.Errorf("expected 0 distance, got %f", zeroDist)
	}
}
