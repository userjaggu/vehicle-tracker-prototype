# Vehicle Tracker Backend Prototype

## What I Built

A backend server that receives GPS location updates from transit vehicles and produces a **GTFS-Realtime Vehicle Positions feed** — the industry-standard format consumed by OneBusAway and other transit apps.

**The core loop:**
1. Vehicles (or the future Android driver app) send GPS coordinates to the server
2. The server maintains the latest position for each active vehicle
3. The server serves a **GTFS-RT protobuf feed** at `/gtfs-rt/vehicle-positions` that any GTFS-RT consumer can read
4. Stale vehicles (no update in 5 minutes) are automatically excluded from the feed

---

## Why I'm Building This

This is a prototype for the [OneBusAway Vehicle Tracker GSoC 2026 project](https://github.com/OneBusAway/vehicle-positions). Many transit agencies in developing countries have bus routes but no AVL (Automatic Vehicle Location) hardware. Drivers carry Android smartphones — this system turns those phones into GPS trackers that produce standards-compliant real-time transit data.

**This prototype demonstrates:**
- ✅ GPS data ingestion via REST API
- ✅ GTFS-RT Vehicle Positions feed generation (protobuf + JSON)
- ✅ Staleness filtering (only active vehicles appear in the feed)
- ✅ Expanded data model (trip, route, bearing, speed, accuracy)
- ✅ System status endpoint for operational monitoring
- ✅ Comprehensive test suite (feed builder + store)
- ✅ Clean code organization following Go conventions

---

## What I'm NOT Doing (Not Yet)

These are planned for the full GSoC project:

- **Authentication**: JWT for drivers, API keys for feed consumers
- **Database**: Persistent storage (SQLite/PostgreSQL)
- **Offline mode**: v2 feature — queuing data when offline
- **Admin Dashboard**: Web UI with Leaflet/OSM map
- **Android App**: Kotlin driver app with foreground service
- **Trip management**: Start/end trip endpoints

---

## How the Code is Organized

```
vehicle-tracker-prototype/
├── main.go                     # Entry point — runs: go run main.go
├── server/
│   └── server.go               # Route registration + server startup
├── handler/
│   ├── location.go             # POST /api/v1/locations  (receives GPS updates)
│   ├── vehicles.go             # GET  /vehicles          (returns all locations)
│   ├── feed.go                 # GET  /gtfs-rt/vehicle-positions (GTFS-RT feed)
│   ├── status.go               # GET  /api/v1/status     (system health)
│   └── helpers.go              # Shared JSON response utilities
├── model/
│   └── vehicle.go              # Location struct (GPS point + trip info)
├── store/
│   ├── memory.go               # Thread-safe in-memory store with staleness
│   └── memory_test.go          # Store unit tests
├── gtfsrt/
│   ├── feed.go                 # GTFS-RT FeedMessage builder
│   └── feed_test.go            # Feed builder unit tests
├── proto/
│   ├── gtfs-realtime.proto     # Official GTFS-RT proto definition
│   └── gtfsrt/
│       └── gtfs-realtime.pb.go # Generated Go protobuf code
├── go.mod                      # Go module (protobuf dependency)
├── go.sum                      # Dependency checksums
└── README.md                   # This file
```

---

## API Endpoints

| Endpoint | Method | Purpose |
|---|---|---|
| `/api/v1/locations` | POST | Submit a vehicle GPS update |
| `/gtfs-rt/vehicle-positions` | GET | GTFS-RT feed (protobuf binary) |
| `/gtfs-rt/vehicle-positions?format=json` | GET | GTFS-RT feed (JSON, for debugging) |
| `/vehicles` | GET | All stored vehicle locations (JSON) |
| `/api/v1/status` | GET | System health and active vehicle count |
| `/location` | POST | Legacy endpoint (alias for `/api/v1/locations`) |

---

## Getting Started

### What You Need

- Go 1.22 or higher ([download here](https://go.dev/dl/))

### Running on Your Machine

```bash
# Clone/download the code
cd vehicle-tracker-prototype

# Build it
go build -o vehicle-tracker

# Run it
./vehicle-tracker

# Server starts on http://localhost:8081
```

### Running Tests

```bash
go test ./... -v
```

---

## How to Use It

### 1. Send a Vehicle Location

```bash
curl -X POST http://localhost:8081/api/v1/locations \
  -H "Content-Type: application/json" \
  -d '{
    "vehicle_id": "bus-42",
    "trip_id": "route_5_0830",
    "route_id": "5",
    "latitude": -1.2921,
    "longitude": 36.8219,
    "bearing": 180.0,
    "speed": 8.5,
    "accuracy": 12.0,
    "timestamp": 1752566400
  }'
```

Response: `{"status": "ok"}`

### 2. Get the GTFS-RT Feed (JSON for debugging)

```bash
curl "http://localhost:8081/gtfs-rt/vehicle-positions?format=json"
```

Response:
```json
{
  "header": {
    "gtfsRealtimeVersion": "2.0",
    "incrementality": "FULL_DATASET",
    "timestamp": "1752566500"
  },
  "entity": [
    {
      "id": "vehicle-bus-42",
      "vehicle": {
        "trip": {
          "tripId": "route_5_0830",
          "routeId": "5"
        },
        "vehicle": {
          "id": "bus-42",
          "label": "bus-42"
        },
        "position": {
          "latitude": -1.2921,
          "longitude": 36.8219,
          "bearing": 180,
          "speed": 8.5
        },
        "timestamp": "1752566400"
      }
    }
  ]
}
```

### 3. Get the GTFS-RT Feed (Protobuf binary — for OneBusAway)

```bash
curl -o feed.pb http://localhost:8081/gtfs-rt/vehicle-positions
# Content-Type: application/x-protobuf
```

### 4. Check System Status

```bash
curl http://localhost:8081/api/v1/status
```

Response:
```json
{
  "status": "ok",
  "active_vehicles": 3,
  "total_vehicles": 5,
  "staleness_threshold": "5m0s",
  "server_time_utc": "2026-03-05T05:00:00Z",
  "feed_endpoint": "/gtfs-rt/vehicle-positions",
  "feed_endpoint_json": "/gtfs-rt/vehicle-positions?format=json"
}
```

---

## Location Report Payload

| Field | Type | Required | Description |
|---|---|---|---|
| `vehicle_id` | string | ✅ | Unique vehicle identifier |
| `latitude` | float64 | ✅ | GPS latitude (decimal degrees) |
| `longitude` | float64 | ✅ | GPS longitude (decimal degrees) |
| `timestamp` | int64 | — | Unix timestamp of GPS fix |
| `trip_id` | string | — | GTFS trip identifier |
| `route_id` | string | — | GTFS route identifier |
| `bearing` | float32 | — | Compass heading (0–360°) |
| `speed` | float32 | — | Speed in meters/second |
| `accuracy` | float32 | — | GPS accuracy in meters |

---

## GTFS-RT Feed Details

The feed follows the [GTFS-RT specification](https://gtfs.org/documentation/realtime/proto/):

- **Feed version:** 2.0
- **Incrementality:** `FULL_DATASET` (every response is the complete state)
- **Content:** `VehiclePosition` entities only
- **Staleness:** Vehicles not reporting for 5 minutes are excluded
- **Formats:** Binary protobuf (default) or JSON (`?format=json`)
- **Proto source:** Official `gtfs-realtime.proto` from [google/transit](https://github.com/google/transit)

---

## How It Works Behind the Scenes

```
                              ┌──────────────────────────┐
  POST /api/v1/locations      │                          │
  ─────────────────────────►  │   Validate + Store       │
  {vehicle_id, lat, lng, ...} │   in-memory map          │
                              │   (sync.RWMutex)         │
                              │                          │
  GET /gtfs-rt/vehicle-       │   Filter stale vehicles  │
      positions               │   Build FeedMessage      │
  ◄─────────────────────────  │   Serialize protobuf     │
  application/x-protobuf      │                          │
                              └──────────────────────────┘
```

---

## Test Results

```
=== RUN   TestBuildFeed_EmptyInput         --- PASS
=== RUN   TestBuildFeed_SingleVehicle      --- PASS
=== RUN   TestBuildFeed_NoTripInfo         --- PASS
=== RUN   TestBuildFeed_MultipleVehicles   --- PASS
=== RUN   TestMarshal_RoundTrip            --- PASS
=== RUN   TestFeedHeader_Timestamp         --- PASS
=== RUN   TestMemoryStore_UpdateAndRetrieve      --- PASS
=== RUN   TestMemoryStore_OverwriteLocation      --- PASS
=== RUN   TestMemoryStore_ActiveLocations        --- PASS
=== RUN   TestMemoryStore_ActiveVehicleCount     --- PASS
=== RUN   TestMemoryStore_ExpandedFields         --- PASS

PASS — 11/11 tests passed
```
