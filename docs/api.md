# ZemlyGo API Documentation

Production-level Uber-style ride-hailing REST & WebSocket API built with Go, Gin, PostgreSQL, Neon Auth, and WebSockets.

---

## Overview

- **Base URL:** `http://localhost:8080`
- **Protocol:** HTTP/1.1 & WebSocket
- **Database:** [Neon Serverless PostgreSQL](https://neon.tech) (with auto-migration & connection pooling)
- **Authentication:** [Neon Auth](https://neon.tech/docs/guides/auth) Bearer JWT Token in `Authorization` header (`Authorization: Bearer <token>`)
- **Maps Engine:** Road network routing with turn-by-turn polyline geometry (OSRM / OpenStreetMap)

### Standard Response Envelope

```json
{
  "success": true,
  "message": "Operation description",
  "data": {}
}
```

Error response:
```json
{
  "success": false,
  "error": "Error description message"
}
```

---

## 1. System Health & Metrics

### Ping
`GET /ping`
- **Response:** `200 OK`
```json
{"message": "PONG"}
```

### Health Check
`GET /health`
- **Response:** `200 OK`
```json
{
  "success": true,
  "message": "ZemlyGo API is online",
  "data": {
    "status": "healthy",
    "env": "development",
    "auth": "neon-auth"
  }
}
```

### Operational Metrics & Stats
`GET /metrics` or `GET /api/v1/admin/stats`
- **Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "total_rides": 1284,
    "active_rides": 42,
    "completed_rides": 1180,
    "online_drivers": 68,
    "total_revenue": 34920.50,
    "total_users": 3810
  }
}
```

---

## 2. Maps, Routing & Places

### 1. Calculate Road Route with Polyline Geometry
`POST /api/v1/maps/route`
- **Payload:**
```json
{
  "pickup_lat": 40.7128,
  "pickup_lng": -74.0060,
  "destination_lat": 40.7589,
  "destination_lng": -73.9851
}
```
- **Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "distance_km": 5.42,
    "duration_min": 14,
    "summary": "Broadway and 5th Ave",
    "source": "osrm",
    "polyline": [
      [40.7128, -74.0060],
      [40.7135, -74.0051],
      [40.7150, -74.0040],
      [40.7589, -73.9851]
    ],
    "steps": [
      {
        "instruction": "depart onto Broadway",
        "distance_km": 1.2,
        "duration_sec": 180
      },
      {
        "instruction": "turn right onto 5th Ave",
        "distance_km": 4.22,
        "duration_sec": 660
      }
    ]
  }
}
```

### 2. Search Places & Autocomplete
`GET /api/v1/maps/places?q=Airport`
- **Response:** `200 OK`
```json
{
  "success": true,
  "data": [
    {
      "place_id": "place_1",
      "name": "JFK International Airport",
      "address": "Queens, NY 11430",
      "lat": 40.6413,
      "lng": -73.7781
    }
  ]
}
```

### 3. Geocode (Address to Coordinates)
`GET /api/v1/maps/geocode?q=Times+Square`

### 4. Reverse Geocode (Coordinates to Address)
`GET /api/v1/maps/reverse?lat=40.7580&lng=-73.9855`

---

## 3. Dynamic Pricing & Surge Heatmap

### Surge Heatmap
`GET /api/v1/pricing/heatmap`
- **Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "zones": [
      {
        "zone_name": "Midtown / Times Square",
        "center_lat": 40.7580,
        "center_lng": -73.9855,
        "multiplier": 1.8,
        "demand_level": "surge"
      },
      {
        "zone_name": "Downtown / Financial District",
        "center_lat": 40.7074,
        "center_lng": -74.0113,
        "multiplier": 1.4,
        "demand_level": "surge"
      }
    ],
    "timestamp": "2026-09-11T12:00:00Z"
  }
}
```

---

## 4. In-App Trip Chat

### Send Message on Active Ride
`POST /api/v1/rides/:id/messages` *(Auth Required)*
- **Payload:**
```json
{
  "message": "I have arrived by the main entrance in a silver Camry."
}
```
- **Real-Time Delivery:** Automatically pushes `ride.chat_message` over WebSocket to the counterparty.

### Get Trip Chat History
`GET /api/v1/rides/:id/messages?limit=50&offset=0` *(Auth Required)*

---

## 5. Emergency SOS & Safety Center

### Trigger Emergency SOS Alert
`POST /api/v1/rides/:id/emergency` *(Auth Required)*
- **Payload:**
```json
{
  "lat": 40.7300,
  "lng": -73.9950,
  "reason": "Suspicious behavior or medical emergency"
}
```
- **Broadcast:** Instantly publishes critical `ride.emergency` event to admin dispatchers and counterparty.

### View Ride Emergency Alerts
`GET /api/v1/rides/:id/emergency` *(Auth Required)*

---

## 6. Driver Earnings & Analytics

### Get Driver Earnings
`GET /api/v1/drivers/earnings?period=weekly` *(Driver Role: period = today | weekly | all)*
- **Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "driver_id": 1,
    "period": "weekly",
    "total_earnings": 892.40,
    "fares_total": 782.40,
    "tips_total": 110.00,
    "completed_rides": 34,
    "trips": [
      {
        "ride_id": 42,
        "fare": 26.50,
        "tip": 5.00,
        "total": 31.50,
        "ride_type": "standard",
        "date": "2026-09-11T10:30:00Z"
      }
    ]
  }
}
```

---

## 7. Ride Tips & Settlement

### Add Tip to Driver
`POST /api/v1/rides/:id/tip` *(Auth Required)*
- **Payload:**
```json
{
  "tip": 5.00
}
```
- **Real-Time Notification:** Sends `ride.tipped` WebSocket event to driver.

---

## 8. Real-time WebSocket (`/ws`)

Connect URL: `ws://localhost:8080/ws?token=<jwt_token>`

### Real-Time Event Catalog

| Event Type | Direction | Description |
|------------|-----------|-------------|
| `ride.requested` | Server -> Drivers | Broadcast ride request to nearby candidate drivers |
| `ride.accepted` | Server -> Rider | Driver accepted ride (includes driver & vehicle info) |
| `ride.driver_arrived` | Server -> Rider | Driver has arrived at pickup address |
| `ride.in_progress` | Server -> Rider | Trip has started |
| `ride.driver_location`| Server -> Rider | Live GPS tracking stream during active ride |
| `ride.chat_message` | Server -> Counterparty | In-app trip chat message |
| `ride.emergency` | Server -> Dispatch/All | Critical SOS safety distress signal |
| `ride.completed` | Server -> Both | Ride completed with final fare |
| `ride.cancelled` | Server -> Counterparty | Trip cancellation notice |
| `payment.completed`| Server -> Both | Payment settled confirmation |
| `ride.tipped` | Server -> Driver | Tip received notification |
| `ride.rated` | Server -> Recipient | Review/score received |
| `ping` / `pong` | Both | Keep-alive heartbeat |
