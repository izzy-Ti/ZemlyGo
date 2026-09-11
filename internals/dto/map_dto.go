package dto

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type RouteRequest struct {
	PickupLat      float64  `json:"pickup_lat" binding:"required"`
	PickupLng      float64  `json:"pickup_lng" binding:"required"`
	DestinationLat float64  `json:"destination_lat" binding:"required"`
	DestinationLng float64  `json:"destination_lng" binding:"required"`
	Waypoints      []LatLng `json:"waypoints,omitempty"`
}

type RouteStep struct {
	Instruction string  `json:"instruction"`
	DistanceKm  float64 `json:"distance_km"`
	DurationSec float64 `json:"duration_sec"`
}

type RouteResponse struct {
	DistanceKm  float64       `json:"distance_km"`
	DurationMin int           `json:"duration_min"`
	Polyline    [][]float64   `json:"polyline"` // [[lat, lng], [lat, lng], ...]
	Steps       []RouteStep   `json:"steps,omitempty"`
	Summary     string        `json:"summary,omitempty"`
	Source      string        `json:"source"` // "osrm" or "interpolated"
}

type GeocodeResponse struct {
	Lat              float64 `json:"lat"`
	Lng              float64 `json:"lng"`
	FormattedAddress string  `json:"formatted_address"`
	DisplayName      string  `json:"display_name"`
	Type             string  `json:"type,omitempty"`
}

type PlaceSuggestion struct {
	PlaceID string  `json:"place_id"`
	Name    string  `json:"name"`
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}
