package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/izzy-Ti/ZemlyGo/configs"
	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/routes"
	"github.com/izzy-Ti/ZemlyGo/internals/service"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	pricingService := service.NewPricingService()
	mapService := service.NewMapService("osrm", "")
	rideHandler := &RideHandler{
		pricingService: pricingService,
	}

	p := &routes.Providers{
		WebsocketHandler: &WebsocketHandler{},
		AuthHandler:      &AuthHandler{},
		UserHandler:      &UserHandler{},
		DriverHandler:    &DriverHandler{},
		VehicleHandler:   &VehicleHandler{},
		RideHandler:      rideHandler,
		PaymentHandler:   &PaymentHandler{},
		RatingHandler:    &RatingHandler{},
		MapHandler:       NewMapHandler(mapService),
		ChatHandler:      &ChatHandler{},
		SafetyHandler:    &SafetyHandler{},
		MetricsHandler:   &MetricsHandler{pricingService: pricingService},
		AuthMiddleware: func(c *gin.Context) {
			c.Set("user_id", uint(1))
			c.Set("user_role", "rider")
			c.Next()
		},
	}

	cfg := &configs.Config{
		Port: "8080",
		Env:  "test",
	}

	routes.SetupRoutes(r, p, cfg)
	return r
}

func TestHealthEndpoint(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var res dto.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestPingEndpoint(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestRideEstimateEndpoint(t *testing.T) {
	r := setupTestRouter()

	reqBody := dto.EstimateFareRequest{
		PickupLat:      40.7128,
		PickupLng:      -74.0060,
		DestinationLat: 40.7589,
		DestinationLng: -73.9851,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rides/estimate", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res dto.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestMapRouteEndpoint(t *testing.T) {
	r := setupTestRouter()

	reqBody := dto.RouteRequest{
		PickupLat:      40.7128,
		PickupLng:      -74.0060,
		DestinationLat: 40.7589,
		DestinationLng: -73.9851,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/maps/route", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res dto.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestMapPlacesEndpoint(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/maps/places?q=Times+Square", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res dto.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestSurgeHeatmapEndpoint(t *testing.T) {
	r := setupTestRouter()

	req, _ := http.NewRequest(http.MethodGet, "/api/v1/pricing/heatmap", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var res dto.APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}
