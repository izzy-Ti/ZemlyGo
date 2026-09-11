package service

import (
	"testing"
)

func TestMapService(t *testing.T) {
	mapService := NewMapService("osrm", "")

	// 1. Test GetRoute
	pLat, pLng := 40.7128, -74.0060
	dLat, dLng := 40.7589, -73.9851

	route, err := mapService.GetRoute(pLat, pLng, dLat, dLng)
	if err != nil {
		t.Fatalf("expected GetRoute to succeed, got %v", err)
	}

	if route.DistanceKm <= 0 {
		t.Errorf("expected positive distance, got %f", route.DistanceKm)
	}

	if route.DurationMin <= 0 {
		t.Errorf("expected positive duration, got %d", route.DurationMin)
	}

	if len(route.Polyline) == 0 {
		t.Errorf("expected non-empty polyline, got 0 coordinates")
	}

	// 2. Test ReverseGeocode
	rev, err := mapService.ReverseGeocode(pLat, pLng)
	if err != nil {
		t.Fatalf("expected ReverseGeocode to succeed, got %v", err)
	}
	if rev.FormattedAddress == "" {
		t.Errorf("expected formatted address, got empty")
	}

	// 3. Test SearchPlaces
	places, err := mapService.SearchPlaces("Airport")
	if err != nil {
		t.Fatalf("expected SearchPlaces to succeed, got %v", err)
	}
	if len(places) == 0 {
		t.Errorf("expected at least one place suggestion, got 0")
	}
}
