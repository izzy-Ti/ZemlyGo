package service

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/izzy-Ti/ZemlyGo/internals/dto"
	"github.com/izzy-Ti/ZemlyGo/internals/utils"
)

type MapService struct {
	client   *http.Client
	provider string
	apiKey   string
}

func NewMapService(provider, apiKey string) *MapService {
	if provider == "" {
		provider = "osrm"
	}
	return &MapService{
		client:   &http.Client{Timeout: 4 * time.Second},
		provider: provider,
		apiKey:   apiKey,
	}
}

// OSRM Response struct
type osrmRouteResponse struct {
	Code   string `json:"code"`
	Routes []struct {
		Distance float64 `json:"distance"` // in meters
		Duration float64 `json:"duration"` // in seconds
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"` // [[lng, lat], ...]
		} `json:"geometry"`
		Legs []struct {
			Summary string `json:"summary"`
			Steps   []struct {
				Maneuver struct {
					Type     string `json:"type"`
					Modifier string `json:"modifier"`
				} `json:"maneuver"`
				Distance float64 `json:"distance"`
				Duration float64 `json:"duration"`
				Name     string  `json:"name"`
			} `json:"steps"`
		} `json:"legs"`
	} `json:"routes"`
}

func (m *MapService) GetRoute(pLat, pLng, dLat, dLng float64) (*dto.RouteResponse, error) {
	// Attempt real road routing via OSRM
	osrmURL := fmt.Sprintf(
		"https://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson&steps=true",
		pLng, pLat, dLng, dLat,
	)

	req, err := http.NewRequest(http.MethodGet, osrmURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "ZemlyGo-UberBackend/1.0")
		resp, err := m.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			var osrmResp osrmRouteResponse
			if err := json.Unmarshal(body, &osrmResp); err == nil && len(osrmResp.Routes) > 0 {
				route := osrmResp.Routes[0]
				distKm := math.Round((route.Distance/1000.0)*100) / 100
				durMin := int(math.Ceil(route.Duration / 60.0))

				// Convert GeoJSON [lng, lat] to frontend map [lat, lng]
				var polyline [][]float64
				for _, coord := range route.Geometry.Coordinates {
					if len(coord) >= 2 {
						polyline = append(polyline, []float64{coord[1], coord[0]})
					}
				}

				var steps []dto.RouteStep
				var summary string
				if len(route.Legs) > 0 {
					summary = route.Legs[0].Summary
					for _, step := range route.Legs[0].Steps {
						instr := step.Maneuver.Type
						if step.Name != "" {
							instr += " onto " + step.Name
						}
						steps = append(steps, dto.RouteStep{
							Instruction: instr,
							DistanceKm:  math.Round((step.Distance/1000.0)*100) / 100,
							DurationSec: math.Round(step.Duration),
						})
					}
				}

				return &dto.RouteResponse{
					DistanceKm:  distKm,
					DurationMin: durMin,
					Polyline:    polyline,
					Steps:       steps,
					Summary:     summary,
					Source:      "osrm",
				}, nil
			}
		}
	}

	// Resilient Fallback: High-resolution smooth geodesic road interpolation
	return m.generateInterpolatedRoute(pLat, pLng, dLat, dLng), nil
}

func (m *MapService) generateInterpolatedRoute(pLat, pLng, dLat, dLng float64) *dto.RouteResponse {
	distKm := utils.Haversine(pLat, pLng, dLat, dLng)
	distKm = math.Round(distKm*100) / 100
	// 30 km/h city average
	durationMin := int(math.Ceil(distKm/30.0*60)) + 2

	// Generate 25 waypoints simulating street turns
	numPoints := 25
	var polyline [][]float64
	for i := 0; i <= numPoints; i++ {
		t := float64(i) / float64(numPoints)
		// Base linear interpolation
		lat := pLat + (dLat-pLat)*t
		lng := pLng + (dLng-pLng)*t

		// Add subtle realistic roadway curve oscillation
		if i > 0 && i < numPoints {
			deviation := math.Sin(t*math.Pi) * 0.003
			lat += deviation
			lng -= deviation * 0.5
		}

		polyline = append(polyline, []float64{
			math.Round(lat*1000000) / 1000000,
			math.Round(lng*1000000) / 1000000,
		})
	}

	return &dto.RouteResponse{
		DistanceKm:  distKm,
		DurationMin: durationMin,
		Polyline:    polyline,
		Summary:     "Direct road corridor",
		Source:      "interpolated",
	}
}

func (m *MapService) ReverseGeocode(lat, lng float64) (*dto.GeocodeResponse, error) {
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f&zoom=18&addressdetails=1", lat, lng)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "ZemlyGo-UberBackend/1.0")
		resp, err := m.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var data struct {
				DisplayName string `json:"display_name"`
				Type        string `json:"type"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&data); err == nil && data.DisplayName != "" {
				return &dto.GeocodeResponse{
					Lat:              lat,
					Lng:              lng,
					FormattedAddress: data.DisplayName,
					DisplayName:      data.DisplayName,
					Type:             data.Type,
				}, nil
			}
		}
	}

	// Fallback
	return &dto.GeocodeResponse{
		Lat:              lat,
		Lng:              lng,
		FormattedAddress: fmt.Sprintf("Location at (%.4f, %.4f)", lat, lng),
		DisplayName:      fmt.Sprintf("Pin (%.4f, %.4f)", lat, lng),
	}, nil
}

func (m *MapService) Geocode(query string) ([]dto.GeocodeResponse, error) {
	if strings.TrimSpace(query) == "" {
		return []dto.GeocodeResponse{}, nil
	}

	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s&limit=5", url.QueryEscape(query))
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err == nil {
		req.Header.Set("User-Agent", "ZemlyGo-UberBackend/1.0")
		resp, err := m.client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			var items []struct {
				Lat         string `json:"lat"`
				Lon         string `json:"lon"`
				DisplayName string `json:"display_name"`
				Type        string `json:"type"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&items); err == nil && len(items) > 0 {
				var results []dto.GeocodeResponse
				for _, item := range items {
					var lat, lng float64
					fmt.Sscanf(item.Lat, "%f", &lat)
					fmt.Sscanf(item.Lon, "%f", &lng)
					results = append(results, dto.GeocodeResponse{
						Lat:              lat,
						Lng:              lng,
						FormattedAddress: item.DisplayName,
						DisplayName:      item.DisplayName,
						Type:             item.Type,
					})
				}
				return results, nil
			}
		}
	}

	// Fallback popular landmarks
	return m.searchLocalLandmarks(query), nil
}

func (m *MapService) SearchPlaces(query string) ([]dto.PlaceSuggestion, error) {
	queryLower := strings.ToLower(strings.TrimSpace(query))
	popularPlaces := []dto.PlaceSuggestion{
		{PlaceID: "place_1", Name: "JFK International Airport", Address: "Queens, NY 11430", Lat: 40.6413, Lng: -73.7781},
		{PlaceID: "place_2", Name: "Times Square", Address: "Manhattan, NY 10036", Lat: 40.7580, Lng: -73.9855},
		{PlaceID: "place_3", Name: "Empire State Building", Address: "350 5th Ave, New York, NY 10118", Lat: 40.7484, Lng: -73.9857},
		{PlaceID: "place_4", Name: "Central Park", Address: "New York, NY", Lat: 40.785091, Lng: -73.968285},
		{PlaceID: "place_5", Name: "Grand Central Terminal", Address: "89 E 42nd St, New York, NY 10017", Lat: 40.7527, Lng: -73.9772},
		{PlaceID: "place_6", Name: "Wall Street Financial District", Address: "Lower Manhattan, NY 10005", Lat: 40.7074, Lng: -74.0113},
		{PlaceID: "place_7", Name: "Brooklyn Bridge", Address: "New York, NY 10038", Lat: 40.7061, Lng: -73.9969},
		{PlaceID: "place_8", Name: "LaGuardia Airport (LGA)", Address: "Queens, NY 11371", Lat: 40.7769, Lng: -73.8740},
	}

	if queryLower == "" {
		return popularPlaces, nil
	}

	var matched []dto.PlaceSuggestion
	for _, p := range popularPlaces {
		if strings.Contains(strings.ToLower(p.Name), queryLower) || strings.Contains(strings.ToLower(p.Address), queryLower) {
			matched = append(matched, p)
		}
	}

	// If no local landmark match, fetch dynamic places
	if len(matched) == 0 {
		geoItems, _ := m.Geocode(query)
		for i, item := range geoItems {
			matched = append(matched, dto.PlaceSuggestion{
				PlaceID: fmt.Sprintf("geo_%d", i+1),
				Name:    item.DisplayName,
				Address: item.FormattedAddress,
				Lat:     item.Lat,
				Lng:     item.Lng,
			})
		}
	}

	return matched, nil
}

func (m *MapService) searchLocalLandmarks(query string) []dto.GeocodeResponse {
	places, _ := m.SearchPlaces(query)
	var res []dto.GeocodeResponse
	for _, p := range places {
		res = append(res, dto.GeocodeResponse{
			Lat:              p.Lat,
			Lng:              p.Lng,
			FormattedAddress: p.Address,
			DisplayName:      p.Name,
		})
	}
	return res
}
